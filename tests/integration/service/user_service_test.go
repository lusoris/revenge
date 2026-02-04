//go:build integration
// +build integration

package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDatabaseURL = "postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable"

func setupUserService(t *testing.T) (*user.Service, *pgxpool.Pool, func()) {
	ctx := context.Background()

	// Create database pool
	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)

	// Create queries and repository
	queries := db.New(pool)
	repo := user.NewPostgresRepository(queries)
	mockStorage := storage.NewMockStorage()
	avatarCfg := config.AvatarConfig{
		StoragePath:  "/tmp/test-avatars",
		MaxSizeBytes: 5 * 1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
	}
	svc := user.NewService(repo, activity.NewNoopLogger(), mockStorage, avatarCfg)

	cleanup := func() {
		pool.Close()
	}

	return svc, pool, cleanup
}

func TestUserService_CreateUser(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Create user
	params := user.CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", timestamp),
		Email:        fmt.Sprintf("test_%d@example.com", timestamp),
		PasswordHash: "password123",
		DisplayName:  stringPtr("Test User"),
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, params.Username, created.Username)
	assert.Equal(t, params.Email, created.Email)
	assert.NotEmpty(t, created.PasswordHash)
	assert.NotEqual(t, "password123", created.PasswordHash, "password should be hashed")
}

func TestUserService_CreateUserDuplicateUsername(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Create first user
	params1 := user.CreateUserParams{
		Username:     fmt.Sprintf("duplicateuser_%d", timestamp),
		Email:        fmt.Sprintf("user1_%d@example.com", timestamp),
		PasswordHash: "password123",
	}

	_, err := svc.CreateUser(ctx, params1)
	require.NoError(t, err)

	// Try to create user with same username
	params2 := user.CreateUserParams{
		Username:     fmt.Sprintf("duplicateuser_%d", timestamp),
		Email:        fmt.Sprintf("user2_%d@example.com", timestamp),
		PasswordHash: "password456",
	}

	_, err = svc.CreateUser(ctx, params2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
}

func TestUserService_CreateUserDuplicateEmail(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Create first user
	params1 := user.CreateUserParams{
		Username:     fmt.Sprintf("user1_%d", timestamp),
		Email:        fmt.Sprintf("duplicate_%d@example.com", timestamp),
		PasswordHash: "password123",
	}

	_, err := svc.CreateUser(ctx, params1)
	require.NoError(t, err)

	// Try to create user with same email
	params2 := user.CreateUserParams{
		Username:     fmt.Sprintf("user2_%d", timestamp),
		Email:        fmt.Sprintf("duplicate_%d@example.com", timestamp),
		PasswordHash: "password456",
	}

	_, err = svc.CreateUser(ctx, params2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestUserService_GetUser(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Create user
	params := user.CreateUserParams{
		Username:     fmt.Sprintf("getuser_%d", timestamp),
		Email:        fmt.Sprintf("get_%d@example.com", timestamp),
		PasswordHash: "password123",
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)

	// Get user by ID
	retrieved, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Username, retrieved.Username)
	assert.Equal(t, created.Email, retrieved.Email)
}

func TestUserService_GetUserByUsername(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()
	username := fmt.Sprintf("usernametest_%d", timestamp)

	// Create user
	params := user.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("username_%d@example.com", timestamp),
		PasswordHash: "password123",
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)

	// Get user by username
	retrieved, err := svc.GetUserByUsername(ctx, username)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("emailtest_%d@example.com", timestamp)

	// Create user
	params := user.CreateUserParams{
		Username:     fmt.Sprintf("emailtest_%d", timestamp),
		Email:        email,
		PasswordHash: "password123",
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)

	// Get user by email
	retrieved, err := svc.GetUserByEmail(ctx, email)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
}

func TestUserService_UpdateUser(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Create user
	params := user.CreateUserParams{
		Username:     fmt.Sprintf("updateuser_%d", timestamp),
		Email:        fmt.Sprintf("update_%d@example.com", timestamp),
		PasswordHash: "password123",
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)

	// Update user
	newEmail := fmt.Sprintf("updated_%d@example.com", timestamp)
	newDisplayName := "Updated Name"

	updateParams := user.UpdateUserParams{
		Email:       &newEmail,
		DisplayName: &newDisplayName,
	}

	updated, err := svc.UpdateUser(ctx, created.ID, updateParams)
	require.NoError(t, err)
	assert.Equal(t, newEmail, updated.Email)
	assert.Equal(t, newDisplayName, *updated.DisplayName)
}

func TestUserService_DeleteUser(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Create user
	params := user.CreateUserParams{
		Username:     fmt.Sprintf("deleteuser_%d", timestamp),
		Email:        fmt.Sprintf("delete_%d@example.com", timestamp),
		PasswordHash: "password123",
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)

	// Delete user
	err = svc.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = svc.GetUser(ctx, created.ID)
	assert.Error(t, err)
}

func TestUserService_ListUsers(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple users
	for i := 0; i < 5; i++ {
		params := user.CreateUserParams{
			Username:     fmt.Sprintf("listuser%d_%d", i, time.Now().UnixNano()),
			Email:        fmt.Sprintf("list%d_%d@example.com", i, time.Now().UnixNano()),
			PasswordHash: "password123",
		}

		_, err := svc.CreateUser(ctx, params)
		require.NoError(t, err)
	}

	// List users
	filters := user.UserFilters{
		Limit:  10,
		Offset: 0,
	}

	users, total, err := svc.ListUsers(ctx, filters)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 5)
	assert.GreaterOrEqual(t, total, int64(5))
}

func TestUserService_PasswordHashing(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()
	password := "MySecurePassword123!"

	// Create user with plain password
	params := user.CreateUserParams{
		Username:     fmt.Sprintf("hashtest_%d", timestamp),
		Email:        fmt.Sprintf("hash_%d@example.com", timestamp),
		PasswordHash: password,
	}

	created, err := svc.CreateUser(ctx, params)
	require.NoError(t, err)

	// Verify password is hashed
	assert.NotEqual(t, password, created.PasswordHash)
	assert.Greater(t, len(created.PasswordHash), 50, "argon2id hash should be long")

	// Verify password verification works
	err = svc.VerifyPassword(created.PasswordHash, password)
	assert.NoError(t, err)

	// Verify wrong password fails
	err = svc.VerifyPassword(created.PasswordHash, "WrongPassword")
	assert.Error(t, err)
}

func TestUserService_ConcurrentCreation(t *testing.T) {
	svc, _, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()

	timestamp := time.Now().UnixNano()

	// Try to create same user concurrently
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			params := user.CreateUserParams{
				Username:     fmt.Sprintf("concurrent_%d", timestamp),
				Email:        fmt.Sprintf("concurrent%d_%d@example.com", index, timestamp),
				PasswordHash: "password123",
			}

			_, err := svc.CreateUser(ctx, params)
			errors <- err
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < 10; i++ {
		err := <-errors
		if err == nil {
			successCount++
		}
	}

	// Only one should succeed (same username)
	assert.Equal(t, 1, successCount, "only one concurrent creation with same username should succeed")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

