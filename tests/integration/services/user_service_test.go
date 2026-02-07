//go:build integration
// +build integration

package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDatabaseURL = "postgres://postgres:postgres@localhost:5432/revenge_test?sslmode=disable"

func setupUserService(t *testing.T) (*user.Service, *pgxpool.Pool, func()) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)

	queries := db.New(pool)
	repo := user.NewPostgresRepository(queries)
	mockStorage := storage.NewMockStorage()
	avatarCfg := config.AvatarConfig{
		StoragePath:  "/tmp/test-avatars",
		MaxSizeBytes: 5 * 1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
	}
	svc := user.NewService(pool, repo, activity.NewNoopLogger(), mockStorage, avatarCfg)

	cleanup := func() {
		pool.Close()
	}

	return svc, pool, cleanup
}

func TestUserService_CreateAndGetUser(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)

	// Create user
	createParams := user.CreateUserParams{
		Username: username,
		Email:    email,
		PasswordHash:"TestPassword123!",
	}

	createdUser, err := svc.CreateUser(ctx, createParams)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	defer func() {
		_ = queries.DeleteUser(ctx, createdUser.ID)
	}()

	// Verify user fields
	assert.Equal(t, username, createdUser.Username)
	assert.Equal(t, email, createdUser.Email)
	assert.NotEmpty(t, createdUser.PasswordHash)
	assert.NotEqual(t, "TestPassword123!", createdUser.PasswordHash) // Should be hashed
	assert.NotEmpty(t, createdUser.ID)

	// Get user by ID
	retrievedUser, err := svc.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, retrievedUser.ID)
	assert.Equal(t, createdUser.Username, retrievedUser.Username)

	// Get user by username
	userByUsername, err := svc.GetUserByUsername(ctx, username)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, userByUsername.ID)

	// Get user by email
	userByEmail, err := svc.GetUserByEmail(ctx, email)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, userByEmail.ID)
}

func TestUserService_UpdateUser(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)

	// Create user
	createdUser, err := svc.CreateUser(ctx, user.CreateUserParams{
		Username: username,
		Email:    email,
		PasswordHash:"TestPassword123!",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, createdUser.ID)
	}()

	// Update user
	newDisplayName := "Updated Display Name"
	newEmail := fmt.Sprintf("updated_%s", email)

	updatedUser, err := svc.UpdateUser(ctx, createdUser.ID, user.UpdateUserParams{
		DisplayName: &newDisplayName,
		Email:       &newEmail,
	})
	require.NoError(t, err)

	assert.Equal(t, newDisplayName, *updatedUser.DisplayName)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, createdUser.Username, updatedUser.Username) // Username shouldn't change
}

func TestUserService_DeleteUser(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)

	// Create user
	createdUser, err := svc.CreateUser(ctx, user.CreateUserParams{
		Username: username,
		Email:    email,
		PasswordHash:"TestPassword123!",
	})
	require.NoError(t, err)

	// Delete user
	err = svc.DeleteUser(ctx, createdUser.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = queries.GetUserByID(ctx, createdUser.ID)
	assert.Error(t, err) // Should not be found
}

func TestUserService_ListUsers(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	timestamp := time.Now().UnixNano()

	// Create multiple users
	userIDs := []uuid.UUID{}
	for i := 0; i < 5; i++ {
		username := fmt.Sprintf("testuser_%d_%d", timestamp, i)
		email := fmt.Sprintf("%s@example.com", username)

		createdUser, err := svc.CreateUser(ctx, user.CreateUserParams{
			Username: username,
			Email:    email,
			PasswordHash:"TestPassword123!",
		})
		require.NoError(t, err)
		userIDs = append(userIDs, createdUser.ID)
	}

	defer func() {
		for _, id := range userIDs {
			_ = queries.DeleteUser(ctx, id)
		}
	}()

	// List users with pagination
	users, total, err := svc.ListUsers(ctx, user.UserFilters{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(users), 5) // At least our 5 users
	assert.GreaterOrEqual(t, total, int64(5))
}

func TestUserService_PasswordValidation(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)

	// Create user
	createdUser, err := svc.CreateUser(ctx, user.CreateUserParams{
		Username: username,
		Email:    email,
		PasswordHash:"CorrectPassword123!",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, createdUser.ID)
	}()

	// Validate correct password
	err = svc.VerifyPassword(createdUser.PasswordHash, "CorrectPassword123!")
	require.NoError(t, err)

	// Validate incorrect password
	err = svc.VerifyPassword(createdUser.PasswordHash, "WrongPassword")
	assert.Error(t, err)
}

func TestUserService_UpdatePassword(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.com", username)

	// Create user
	createdUser, err := svc.CreateUser(ctx, user.CreateUserParams{
		Username: username,
		Email:    email,
		PasswordHash:"OldPassword123!",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, createdUser.ID)
	}()

	// Update password
	err = svc.UpdatePassword(ctx, createdUser.ID, "OldPassword123!", "NewPassword456!")
	require.NoError(t, err)

	// Validate old password no longer works
	updatedUser, err := svc.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	err = svc.VerifyPassword(updatedUser.PasswordHash, "OldPassword123!")
	assert.Error(t, err)

	// Validate new password works
	err = svc.VerifyPassword(updatedUser.PasswordHash, "NewPassword456!")
	require.NoError(t, err)
}

func TestUserService_DuplicateUsername(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	email1 := fmt.Sprintf("%s_1@example.com", username)
	email2 := fmt.Sprintf("%s_2@example.com", username)

	// Create first user
	user1, err := svc.CreateUser(ctx, user.CreateUserParams{
		Username: username,
		Email:    email1,
		PasswordHash:"TestPassword123!",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, user1.ID)
	}()

	// Try to create second user with same username
	_, err = svc.CreateUser(ctx, user.CreateUserParams{
		Username: username, // Duplicate username
		Email:    email2,
		PasswordHash:"TestPassword123!",
	})
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestUserService_DuplicateEmail(t *testing.T) {
	svc, pool, cleanup := setupUserService(t)
	defer cleanup()

	ctx := context.Background()
	queries := db.New(pool)

	timestamp := time.Now().UnixNano()
	username1 := fmt.Sprintf("testuser_%d_1", timestamp)
	username2 := fmt.Sprintf("testuser_%d_2", timestamp)
	email := fmt.Sprintf("shared_%d@example.com", timestamp)

	// Create first user
	user1, err := svc.CreateUser(ctx, user.CreateUserParams{
		Username: username1,
		Email:    email,
		PasswordHash:"TestPassword123!",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, user1.ID)
	}()

	// Try to create second user with same email
	_, err = svc.CreateUser(ctx, user.CreateUserParams{
		Username: username2,
		Email:    email, // Duplicate email
		PasswordHash:"TestPassword123!",
	})
	assert.Error(t, err) // Should fail due to unique constraint
}
