// Package database provides advanced integration tests for constraint handling
//go:build integration
// +build integration

package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseUniqueConstraints(t *testing.T) {
	requirePostgres(t)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	// Create first user
	username := fmt.Sprintf("unique_test_%d", time.Now().UnixNano())
	user1, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "hashed_password",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, user1.ID)
	}()

	// Try to create duplicate username - should fail
	_, err = queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username, // Same username
		Email:        "different@example.com",
		PasswordHash: "hashed_password",
	})
	assert.Error(t, err, "should fail on duplicate username")
	assert.Contains(t, err.Error(), "users_username_key", "should be unique constraint error")

	// Try to create duplicate email - should fail
	_, err = queries.CreateUser(ctx, db.CreateUserParams{
		Username:     fmt.Sprintf("%s_different", username),
		Email:        fmt.Sprintf("%s@example.com", username), // Same email
		PasswordHash: "hashed_password",
	})
	assert.Error(t, err, "should fail on duplicate email")
	assert.Contains(t, err.Error(), "users_email_key", "should be unique constraint error")
}

func TestDatabaseTransactionCommit(t *testing.T) {
	requirePostgres(t)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	username := fmt.Sprintf("tx_commit_test_%d", time.Now().UnixNano())

	// Start transaction and commit
	tx, err := pool.Begin(ctx)
	require.NoError(t, err)

	txQueries := queries.WithTx(tx)
	user, err := txQueries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "hashed_password",
	})
	require.NoError(t, err)

	// Commit transaction
	err = tx.Commit(ctx)
	require.NoError(t, err)

	// Verify user exists after commit
	fetchedUser, err := queries.GetUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, username, fetchedUser.Username)

	// Cleanup
	_ = queries.DeleteUser(ctx, user.ID)
}

func TestDatabaseNullHandling(t *testing.T) {
	requirePostgres(t)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	username := fmt.Sprintf("null_test_%d", time.Now().UnixNano())
	displayName := "Test Display Name"
	timezone := "America/New_York"

	// Create user with optional fields
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "hashed_password",
		DisplayName:  &displayName,
		Timezone:     &timezone,
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, user.ID)
	}()

	// Verify optional fields were set
	assert.NotNil(t, user.DisplayName)
	assert.Equal(t, displayName, *user.DisplayName)
	assert.NotNil(t, user.Timezone)
	assert.Equal(t, timezone, *user.Timezone)
}

func TestDatabaseConcurrentUpdates(t *testing.T) {
	requirePostgres(t)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	// Create test user
	username := fmt.Sprintf("concurrent_update_test_%d", time.Now().UnixNano())
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "initial_password",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, user.ID)
	}()

	// Run 10 concurrent updates (updating email which is allowed)
	results := make(chan error, 10)
	timestamp := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		index := i
		go func() {
			newEmail := fmt.Sprintf("newemail_%d_%d@example.com", timestamp, index)
			_, err := queries.UpdateUser(ctx, db.UpdateUserParams{
				UserID: user.ID,
				Email:  &newEmail,
			})
			results <- err
		}()
	}

	// Verify all updates succeeded (last write wins)
	for i := 0; i < 10; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent update should succeed")
	}

	// Verify final state
	updatedUser, err := queries.GetUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, fmt.Sprintf("%s@example.com", username), updatedUser.Email, "email should be updated")
}

func TestDatabaseTransactionIsolation(t *testing.T) {
	requirePostgres(t)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	username := fmt.Sprintf("isolation_test_%d", time.Now().UnixNano())

	// Transaction 1: Create user but don't commit yet
	tx1, err := pool.Begin(ctx)
	require.NoError(t, err)

	tx1Queries := queries.WithTx(tx1)
	user, err := tx1Queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "hashed_password",
	})
	require.NoError(t, err)

	// Transaction 2: Try to read the uncommitted user (should not see it)
	_, err = queries.GetUserByID(ctx, user.ID)
	assert.Error(t, err, "should not see uncommitted data from other transaction")

	// Commit transaction 1
	err = tx1.Commit(ctx)
	require.NoError(t, err)

	// Now transaction 2 should see the user
	fetchedUser, err := queries.GetUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, username, fetchedUser.Username)

	// Cleanup
	_ = queries.DeleteUser(ctx, user.ID)
}

func TestDatabaseQueryTimeout(t *testing.T) {
	requirePostgres(t)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	// Create a context with very short timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	// This query should timeout
	time.Sleep(2 * time.Millisecond) // Ensure we're past the deadline
	_, err = pool.Exec(ctxTimeout, "SELECT pg_sleep(1)")
	assert.Error(t, err, "query should timeout")
	assert.Contains(t, err.Error(), "context", "should be context error")
}
