package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// User represents a test user fixture.
type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	DisplayName  string
	QAREnabled   bool
	IsActive     bool
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Session represents a test session fixture.
type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	TokenHash        string
	RefreshTokenHash string
	IPAddress        string
	UserAgent        string
	Scopes           []string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

// CreateUser creates a test user in the database.
//
// Example:
//
//	user := testutil.CreateUser(t, db.Pool, testutil.User{
//	    Username: "testuser",
//	    Email: "test@example.com",
//	    IsAdmin: true,
//	})
func CreateUser(t *testing.T, pool *pgxpool.Pool, user User) *User {
	t.Helper()

	// Set defaults
	if user.ID == uuid.Nil {
		user.ID = uuid.Must(uuid.NewV7())
	}
	if user.PasswordHash == "" {
		user.PasswordHash = "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$hash" // Dummy hash
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = user.CreatedAt
	}

	ctx := context.Background()
	err := pool.QueryRow(ctx, `
		INSERT INTO shared.users (
			id, username, email, password_hash, display_name,
			qar_enabled, is_active, is_admin,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, username, email, password_hash, display_name,
		          qar_enabled, is_active, is_admin,
		          created_at, updated_at
	`,
		user.ID, user.Username, user.Email, user.PasswordHash, user.DisplayName,
		user.QAREnabled, user.IsActive, user.IsAdmin,
		user.CreatedAt, user.UpdatedAt,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.QAREnabled, &user.IsActive, &user.IsAdmin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	require.NoError(t, err, "failed to create test user")

	return &user
}

// CreateSession creates a test session in the database.
//
// Example:
//
//	session := testutil.CreateSession(t, db.Pool, testutil.Session{
//	    UserID: user.ID,
//	    TokenHash: "test_token_hash",
//	    Scopes: []string{"legacy:read"},
//	})
func CreateSession(t *testing.T, pool *pgxpool.Pool, session Session) *Session {
	t.Helper()

	// Set defaults
	if session.ID == uuid.Nil {
		session.ID = uuid.Must(uuid.NewV7())
	}
	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = time.Now().Add(24 * time.Hour)
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	if session.Scopes == nil {
		session.Scopes = []string{}
	}

	ctx := context.Background()
	err := pool.QueryRow(ctx, `
		INSERT INTO shared.sessions (
			id, user_id, token_hash, refresh_token_hash,
			ip_address, user_agent, scopes,
			expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, token_hash, refresh_token_hash,
		          ip_address, user_agent, scopes,
		          expires_at, created_at
	`,
		session.ID, session.UserID, session.TokenHash, session.RefreshTokenHash,
		session.IPAddress, session.UserAgent, session.Scopes,
		session.ExpiresAt, session.CreatedAt,
	).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.RefreshTokenHash,
		&session.IPAddress, &session.UserAgent, &session.Scopes,
		&session.ExpiresAt, &session.CreatedAt,
	)

	require.NoError(t, err, "failed to create test session")

	return &session
}

// DefaultUser returns a default test user for quick testing.
func DefaultUser() User {
	return User{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		IsActive:    true,
		IsAdmin:     false,
		QAREnabled:  false,
	}
}

// AdminUser returns a test admin user.
func AdminUser() User {
	return User{
		Username:    "admin",
		Email:       "admin@example.com",
		DisplayName: "Admin User",
		IsActive:    true,
		IsAdmin:     true,
		QAREnabled:  false,
	}
}

// QARUser returns a test user with QAR access.
func QARUser() User {
	return User{
		Username:    "qaruser",
		Email:       "qar@example.com",
		DisplayName: "QAR User",
		IsActive:    true,
		IsAdmin:     false,
		QAREnabled:  true,
	}
}
