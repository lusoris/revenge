//go:build integration

// Package integration provides integration tests for the database layer.
// These tests require Docker to be running and are skipped by default.
// Run with: go test -tags=integration ./tests/integration/...
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/infra/database/repository"
)

// testDB holds the test database connection and cleanup function.
type testDB struct {
	pool      *pgxpool.Pool
	container testcontainers.Container
}

// setupTestDB creates a PostgreSQL container and returns a connection pool.
func setupTestDB(t *testing.T) *testDB {
	t.Helper()
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("jellyfin_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	// Run migrations
	if err := runMigrations(ctx, pool); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return &testDB{
		pool:      pool,
		container: pgContainer,
	}
}

// cleanup closes the connection pool and terminates the container.
func (db *testDB) cleanup(t *testing.T) {
	t.Helper()
	db.pool.Close()
	if err := db.container.Terminate(context.Background()); err != nil {
		t.Errorf("failed to terminate container: %v", err)
	}
}

// runMigrations applies the database schema.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Simplified migration - just create the tables we need for tests
	schema := `
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) UNIQUE,
			password_hash VARCHAR(255),
			display_name VARCHAR(255),
			is_admin BOOLEAN NOT NULL DEFAULT false,
			is_disabled BOOLEAN NOT NULL DEFAULT false,
			last_login_at TIMESTAMPTZ,
			last_activity_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(64) NOT NULL UNIQUE,
			refresh_token_hash VARCHAR(64) UNIQUE,
			device_id VARCHAR(255),
			device_name VARCHAR(255),
			client_name VARCHAR(255),
			client_version VARCHAR(50),
			ip_address INET,
			expires_at TIMESTAMPTZ NOT NULL,
			refresh_expires_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE oidc_providers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL UNIQUE,
			display_name VARCHAR(255) NOT NULL,
			issuer_url VARCHAR(512) NOT NULL,
			client_id VARCHAR(255) NOT NULL,
			client_secret_encrypted BYTEA NOT NULL,
			scopes TEXT[] NOT NULL DEFAULT ARRAY['openid', 'profile', 'email'],
			enabled BOOLEAN NOT NULL DEFAULT true,
			auto_create_users BOOLEAN NOT NULL DEFAULT true,
			default_admin BOOLEAN NOT NULL DEFAULT false,
			claim_mappings JSONB NOT NULL DEFAULT '{}',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE oidc_user_links (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			provider_id UUID NOT NULL REFERENCES oidc_providers(id) ON DELETE CASCADE,
			subject VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			last_login_at TIMESTAMPTZ,
			UNIQUE(provider_id, subject)
		);
	`

	_, err := pool.Exec(ctx, schema)
	return err
}

// =============================================================================
// USER REPOSITORY TESTS
// =============================================================================

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	t.Run("creates user successfully", func(t *testing.T) {
		email := "test@example.com"
		displayName := "Test User"
		passwordHash := "hashed_password"

		user, err := repo.Create(ctx, domain.CreateUserParams{
			Username:     "testuser",
			Email:        &email,
			DisplayName:  &displayName,
			PasswordHash: &passwordHash,
			IsAdmin:      false,
		})

		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		if user.ID == uuid.Nil {
			t.Error("expected non-nil user ID")
		}
		if user.Username != "testuser" {
			t.Errorf("expected username 'testuser', got %q", user.Username)
		}
		if user.Email == nil || *user.Email != email {
			t.Errorf("expected email %q, got %v", email, user.Email)
		}
		if user.IsAdmin {
			t.Error("expected non-admin user")
		}
	})

	t.Run("creates admin user", func(t *testing.T) {
		user, err := repo.Create(ctx, domain.CreateUserParams{
			Username: "adminuser",
			IsAdmin:  true,
		})

		if err != nil {
			t.Fatalf("failed to create admin user: %v", err)
		}

		if !user.IsAdmin {
			t.Error("expected admin user")
		}
	})

	t.Run("rejects duplicate username", func(t *testing.T) {
		_, err := repo.Create(ctx, domain.CreateUserParams{
			Username: "testuser", // Already exists
		})

		if err != domain.ErrDuplicateUsername {
			t.Errorf("expected ErrDuplicateUsername, got %v", err)
		}
	})

	t.Run("rejects duplicate email", func(t *testing.T) {
		email := "test@example.com" // Already exists
		_, err := repo.Create(ctx, domain.CreateUserParams{
			Username: "anotheruser",
			Email:    &email,
		})

		if err != domain.ErrDuplicateEmail {
			t.Errorf("expected ErrDuplicateEmail, got %v", err)
		}
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	// Create a user first
	created, err := repo.Create(ctx, domain.CreateUserParams{
		Username: "getbyid_user",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("finds existing user", func(t *testing.T) {
		user, err := repo.GetByID(ctx, created.ID)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		if user.ID != created.ID {
			t.Errorf("expected ID %v, got %v", created.ID, user.ID)
		}
		if user.Username != "getbyid_user" {
			t.Errorf("expected username 'getbyid_user', got %q", user.Username)
		}
	})

	t.Run("returns not found for non-existent user", func(t *testing.T) {
		_, err := repo.GetByID(ctx, uuid.New())
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	// Create a user
	_, err := repo.Create(ctx, domain.CreateUserParams{
		Username: "byusername_user",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("finds user by username", func(t *testing.T) {
		user, err := repo.GetByUsername(ctx, "byusername_user")
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		if user.Username != "byusername_user" {
			t.Errorf("expected username 'byusername_user', got %q", user.Username)
		}
	})

	t.Run("returns not found for non-existent username", func(t *testing.T) {
		_, err := repo.GetByUsername(ctx, "nonexistent")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	// Create a user
	created, err := repo.Create(ctx, domain.CreateUserParams{
		Username: "update_user",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("updates email and display name", func(t *testing.T) {
		newEmail := "updated@example.com"
		newDisplay := "Updated Name"

		err := repo.Update(ctx, domain.UpdateUserParams{
			ID:          created.ID,
			Email:       &newEmail,
			DisplayName: &newDisplay,
		})
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		// Verify update
		user, _ := repo.GetByID(ctx, created.ID)
		if user.Email == nil || *user.Email != newEmail {
			t.Errorf("expected email %q, got %v", newEmail, user.Email)
		}
		if user.DisplayName == nil || *user.DisplayName != newDisplay {
			t.Errorf("expected display name %q, got %v", newDisplay, user.DisplayName)
		}
	})

	t.Run("updates admin status", func(t *testing.T) {
		isAdmin := true
		err := repo.Update(ctx, domain.UpdateUserParams{
			ID:      created.ID,
			IsAdmin: &isAdmin,
		})
		if err != nil {
			t.Fatalf("failed to update admin status: %v", err)
		}

		user, _ := repo.GetByID(ctx, created.ID)
		if !user.IsAdmin {
			t.Error("expected user to be admin")
		}
	})

	t.Run("returns not found for non-existent user", func(t *testing.T) {
		newEmail := "test@test.com"
		err := repo.Update(ctx, domain.UpdateUserParams{
			ID:    uuid.New(),
			Email: &newEmail,
		})
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	// Create a user
	created, err := repo.Create(ctx, domain.CreateUserParams{
		Username: "delete_user",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("deletes user successfully", func(t *testing.T) {
		err := repo.Delete(ctx, created.ID)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		// Verify deletion
		_, err = repo.GetByID(ctx, created.ID)
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound after deletion, got %v", err)
		}
	})
}

func TestUserRepository_SetPassword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	// Create a user without password
	created, err := repo.Create(ctx, domain.CreateUserParams{
		Username: "password_user",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("sets password", func(t *testing.T) {
		err := repo.SetPassword(ctx, created.ID, "new_password_hash")
		if err != nil {
			t.Fatalf("failed to set password: %v", err)
		}

		// Verify password was set
		user, _ := repo.GetByID(ctx, created.ID)
		if user.PasswordHash == nil || *user.PasswordHash != "new_password_hash" {
			t.Error("expected password hash to be set")
		}
	})
}

func TestUserRepository_Count(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	repo := repository.NewUserRepository(db.pool)
	ctx := context.Background()

	// Create some users
	for i := 0; i < 3; i++ {
		_, err := repo.Create(ctx, domain.CreateUserParams{
			Username: fmt.Sprintf("count_user_%d", i),
			IsAdmin:  i == 0, // First user is admin
		})
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}
	}

	t.Run("counts all users", func(t *testing.T) {
		count, err := repo.Count(ctx)
		if err != nil {
			t.Fatalf("failed to count users: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 users, got %d", count)
		}
	})

	t.Run("counts admin users", func(t *testing.T) {
		count, err := repo.CountAdmins(ctx)
		if err != nil {
			t.Fatalf("failed to count admin users: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 admin user, got %d", count)
		}
	})
}

// =============================================================================
// SESSION REPOSITORY TESTS
// =============================================================================

func TestSessionRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	userRepo := repository.NewUserRepository(db.pool)
	sessionRepo := repository.NewSessionRepository(db.pool)
	ctx := context.Background()

	// Create a user first
	user, err := userRepo.Create(ctx, domain.CreateUserParams{
		Username: "session_user",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	t.Run("creates session successfully", func(t *testing.T) {
		deviceID := "device-123"
		deviceName := "Test Device"

		session, err := sessionRepo.Create(ctx, domain.CreateSessionParams{
			UserID:     user.ID,
			TokenHash:  "token_hash_123",
			DeviceID:   &deviceID,
			DeviceName: &deviceName,
			ExpiresAt:  time.Now().Add(24 * time.Hour),
		})

		if err != nil {
			t.Fatalf("failed to create session: %v", err)
		}

		if session.ID == uuid.Nil {
			t.Error("expected non-nil session ID")
		}
		if session.UserID != user.ID {
			t.Errorf("expected user ID %v, got %v", user.ID, session.UserID)
		}
		if session.TokenHash != "token_hash_123" {
			t.Errorf("expected token hash 'token_hash_123', got %q", session.TokenHash)
		}
	})
}

func TestSessionRepository_GetByTokenHash(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	userRepo := repository.NewUserRepository(db.pool)
	sessionRepo := repository.NewSessionRepository(db.pool)
	ctx := context.Background()

	// Create a user
	user, _ := userRepo.Create(ctx, domain.CreateUserParams{
		Username: "token_user",
	})

	// Create a session
	_, err := sessionRepo.Create(ctx, domain.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: "find_by_token_hash",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	t.Run("finds session by token hash", func(t *testing.T) {
		session, err := sessionRepo.GetByTokenHash(ctx, "find_by_token_hash")
		if err != nil {
			t.Fatalf("failed to get session: %v", err)
		}

		if session.TokenHash != "find_by_token_hash" {
			t.Errorf("expected token hash 'find_by_token_hash', got %q", session.TokenHash)
		}
	})

	t.Run("returns not found for non-existent token", func(t *testing.T) {
		_, err := sessionRepo.GetByTokenHash(ctx, "nonexistent")
		if err != domain.ErrSessionNotFound {
			t.Errorf("expected ErrSessionNotFound, got %v", err)
		}
	})

	t.Run("returns expired for expired session", func(t *testing.T) {
		// Create expired session
		_, err := sessionRepo.Create(ctx, domain.CreateSessionParams{
			UserID:    user.ID,
			TokenHash: "expired_token_hash",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
		})
		if err != nil {
			t.Fatalf("failed to create expired session: %v", err)
		}

		_, err = sessionRepo.GetByTokenHash(ctx, "expired_token_hash")
		if err != domain.ErrSessionExpired {
			t.Errorf("expected ErrSessionExpired, got %v", err)
		}
	})
}

func TestSessionRepository_DeleteByUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	userRepo := repository.NewUserRepository(db.pool)
	sessionRepo := repository.NewSessionRepository(db.pool)
	ctx := context.Background()

	// Create a user
	user, _ := userRepo.Create(ctx, domain.CreateUserParams{
		Username: "delete_sessions_user",
	})

	// Create multiple sessions
	for i := 0; i < 3; i++ {
		_, err := sessionRepo.Create(ctx, domain.CreateSessionParams{
			UserID:    user.ID,
			TokenHash: fmt.Sprintf("delete_user_token_%d", i),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		if err != nil {
			t.Fatalf("failed to create session: %v", err)
		}
	}

	// Verify sessions exist
	count, _ := sessionRepo.CountByUser(ctx, user.ID)
	if count != 3 {
		t.Fatalf("expected 3 sessions, got %d", count)
	}

	t.Run("deletes all user sessions", func(t *testing.T) {
		err := sessionRepo.DeleteByUser(ctx, user.ID)
		if err != nil {
			t.Fatalf("failed to delete sessions: %v", err)
		}

		count, _ := sessionRepo.CountByUser(ctx, user.ID)
		if count != 0 {
			t.Errorf("expected 0 sessions after deletion, got %d", count)
		}
	})
}

func TestSessionRepository_GetWithUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := setupTestDB(t)
	defer db.cleanup(t)

	userRepo := repository.NewUserRepository(db.pool)
	sessionRepo := repository.NewSessionRepository(db.pool)
	ctx := context.Background()

	// Create a user
	email := "withuser@example.com"
	user, _ := userRepo.Create(ctx, domain.CreateUserParams{
		Username: "with_user_test",
		Email:    &email,
		IsAdmin:  true,
	})

	// Create a session
	_, err := sessionRepo.Create(ctx, domain.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: "with_user_token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	t.Run("returns session with user info", func(t *testing.T) {
		result, err := sessionRepo.GetWithUser(ctx, "with_user_token")
		if err != nil {
			t.Fatalf("failed to get session with user: %v", err)
		}

		if result.Username != "with_user_test" {
			t.Errorf("expected username 'with_user_test', got %q", result.Username)
		}
		if !result.IsAdmin {
			t.Error("expected user to be admin")
		}
		if result.Email == nil || *result.Email != email {
			t.Errorf("expected email %q, got %v", email, result.Email)
		}
	})
}
