// Package database provides integration tests for database operations
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

const testDatabaseURL = "postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable"

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err, "should connect to database")
	defer pool.Close()

	// Test ping
	err = pool.Ping(ctx)
	assert.NoError(t, err, "should ping successfully")
}

func TestDatabasePoolStats(t *testing.T) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(testDatabaseURL)
	require.NoError(t, err)

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	defer pool.Close()

	// pgxpool creates connections lazily, not eagerly
	// Initial pool may have 0 connections
	stats := pool.Stat()
	assert.LessOrEqual(t, stats.TotalConns(), int32(10), "should not exceed max connections")

	// After acquiring connections, MinConns should be respected
	conn1, err := pool.Acquire(ctx)
	require.NoError(t, err)
	conn1.Release()

	conn2, err := pool.Acquire(ctx)
	require.NoError(t, err)
	conn2.Release()

	// Now we should have at least MinConns connections
	time.Sleep(100 * time.Millisecond) // Give pool time to stabilize
	stats = pool.Stat()
	assert.GreaterOrEqual(t, stats.TotalConns(), int32(2), "should maintain min connections after usage")
}

func TestDatabaseTransactions(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	// Start transaction
	tx, err := pool.Begin(ctx)
	require.NoError(t, err)

	txQueries := queries.WithTx(tx)

	// Create a test user
	user, err := txQueries.CreateUser(ctx, db.CreateUserParams{
		Username:     "tx_test_user",
		Email:        "tx_test@example.com",
		PasswordHash: "hashed_password",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Rollback - user should not exist
	err = tx.Rollback(ctx)
	require.NoError(t, err)

	// Verify user was rolled back
	_, err = queries.GetUserByUsername(ctx, "tx_test_user")
	assert.Error(t, err, "user should not exist after rollback")
}

func TestDatabaseConcurrentOperations(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	// Create base user for concurrent operations with unique username
	username := fmt.Sprintf("concurrent_test_%d", time.Now().UnixNano())
	baseUser, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "hashed_password",
	})
	require.NoError(t, err)
	defer func() {
		// Cleanup
		_ = queries.DeleteUser(ctx, baseUser.ID)
	}()

	// Run 50 concurrent reads
	results := make(chan error, 50)
	for i := 0; i < 50; i++ {
		go func() {
			_, err := queries.GetUserByID(ctx, baseUser.ID)
			results <- err
		}()
	}

	// Verify all reads succeeded
	for i := 0; i < 50; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent read should succeed")
	}
}

func TestDatabaseConnectionPoolExhaustion(t *testing.T) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(testDatabaseURL)
	require.NoError(t, err)

	// Set very small pool for testing
	config.MaxConns = 2
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	defer pool.Close()

	// Acquire all connections
	conn1, err := pool.Acquire(ctx)
	require.NoError(t, err)
	defer conn1.Release()

	conn2, err := pool.Acquire(ctx)
	require.NoError(t, err)
	defer conn2.Release()

	// Next acquire should timeout or wait
	ctx2, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	_, err = pool.Acquire(ctx2)
	assert.Error(t, err, "should fail when pool is exhausted")
}

func TestDatabaseSchemaExists(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	// Check shared schema exists
	var exists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_namespace WHERE nspname = 'shared')").Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "shared schema should exist")

	// Check key tables exist
	tables := []string{"users", "sessions", "api_keys", "casbin_rule"}
	for _, table := range tables {
		var tableExists bool
		err = pool.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = $1)",
			table,
		).Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "table %s should exist in shared schema", table)
	}
}

func TestDatabaseNullableTypes(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)
	defer pool.Close()

	queries := db.New(pool)

	// Create user with nullable fields using unique username
	username := fmt.Sprintf("nullable_test_%d", time.Now().UnixNano())
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        fmt.Sprintf("%s@example.com", username),
		PasswordHash: "hashed_password",
	})
	require.NoError(t, err)
	defer func() {
		_ = queries.DeleteUser(ctx, user.ID)
	}()

	// Avatar URL should be null
	assert.Nil(t, user.AvatarUrl, "avatar_url should be null for new user")

	// Email verified should be null or false
	if user.EmailVerified != nil {
		assert.False(t, *user.EmailVerified, "email_verified should be false for new user")
	}
}
