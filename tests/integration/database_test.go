//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseMigrations(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	ctx := context.Background()

	// Verify shared schema exists
	var schemaExists bool
	err := ts.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.schemata
			WHERE schema_name = 'shared'
		)
	`).Scan(&schemaExists)
	require.NoError(t, err, "should be able to query schema existence")
	assert.True(t, schemaExists, "shared schema should exist")

	// Verify users table exists with correct structure
	var tableExists bool
	err = ts.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'shared' AND table_name = 'users'
		)
	`).Scan(&tableExists)
	require.NoError(t, err, "should be able to query table existence")
	assert.True(t, tableExists, "users table should exist")

	// Verify users table columns
	type columnInfo struct {
		Name       string
		DataType   string
		IsNullable string
	}

	rows, err := ts.DB.Pool.Query(ctx, `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'shared' AND table_name = 'users'
		ORDER BY ordinal_position
	`)
	require.NoError(t, err, "should be able to query table columns")
	defer rows.Close()

	var columns []columnInfo
	for rows.Next() {
		var col columnInfo
		err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable)
		require.NoError(t, err, "should be able to scan column info")
		columns = append(columns, col)
	}

	// Verify expected columns exist
	expectedColumns := map[string]bool{
		"id":            true,
		"username":      true,
		"email":         true,
		"password_hash": true,
		"created_at":    true,
		"updated_at":    true,
	}

	for _, col := range columns {
		if expectedColumns[col.Name] {
			t.Logf("Found expected column: %s (%s)", col.Name, col.DataType)
			delete(expectedColumns, col.Name)
		}
	}

	assert.Empty(t, expectedColumns, "all expected columns should exist")

	// Verify sessions table exists
	err = ts.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'shared' AND table_name = 'sessions'
		)
	`).Scan(&tableExists)
	require.NoError(t, err, "should be able to query sessions table existence")
	assert.True(t, tableExists, "sessions table should exist")
}

func TestDatabaseHealthCheck(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	ctx := context.Background()

	// Verify database is healthy
	err := ts.DB.Pool.Ping(ctx)
	assert.NoError(t, err, "database should be healthy")

	// Close pool to simulate failure
	ts.DB.Pool.Close()

	// Wait a moment
	time.Sleep(50 * time.Millisecond)

	// Verify ping now fails
	err = ts.DB.Pool.Ping(ctx)
	assert.Error(t, err, "ping should fail after pool close")
}

func TestDatabaseGracefulShutdown(t *testing.T) {
	// Setup test server
	ts := setupServer(t)

	ctx := context.Background()

	// Verify database is accessible
	err := ts.DB.Pool.Ping(ctx)
	require.NoError(t, err, "database should be accessible initially")

	// Close pool gracefully
	ts.DB.Pool.Close()

	// Verify pool is closed
	err = ts.DB.Pool.Ping(ctx)
	assert.Error(t, err, "pool should be closed")

	// Cleanup container
	ts.DB.Close()
	ts.DB = nil

	// Stop app
	ts.App.RequireStop()
	ts.App = nil
}

func TestDatabaseConnectionPooling(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	ctx := context.Background()

	// Acquire multiple connections
	const numConns = 5
	conns := make([]*pgxpool.Conn, numConns)

	for i := 0; i < numConns; i++ {
		conn, err := ts.DB.Pool.Acquire(ctx)
		require.NoError(t, err, "should be able to acquire connection %d", i)
		conns[i] = conn
	}

	// Verify all connections work
	for i, conn := range conns {
		err := conn.Ping(ctx)
		assert.NoError(t, err, "connection %d should work", i)
	}

	// Release all connections
	for _, conn := range conns {
		conn.Release()
	}

	// Verify pool is still healthy
	err := ts.DB.Pool.Ping(ctx)
	assert.NoError(t, err, "pool should be healthy after releasing connections")
}

func TestDatabaseTransactions(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)
	resetDatabase(t, ts)

	ctx := context.Background()

	// Start transaction
	tx, err := ts.DB.Pool.Begin(ctx)
	require.NoError(t, err, "should be able to start transaction")

	// Insert test data within transaction
	_, err = tx.Exec(ctx, `
		INSERT INTO shared.users (username, email, password_hash)
		VALUES ('txuser', 'tx@example.com', 'hash123')
	`)
	require.NoError(t, err, "should be able to insert within transaction")

	// Rollback transaction
	err = tx.Rollback(ctx)
	require.NoError(t, err, "should be able to rollback transaction")

	// Verify data was not persisted
	var count int
	err = ts.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM shared.users WHERE username = 'txuser'").Scan(&count)
	require.NoError(t, err, "should be able to query")
	assert.Equal(t, 0, count, "data should not be persisted after rollback")

	// Start new transaction
	tx, err = ts.DB.Pool.Begin(ctx)
	require.NoError(t, err, "should be able to start new transaction")

	// Insert test data
	_, err = tx.Exec(ctx, `
		INSERT INTO shared.users (username, email, password_hash)
		VALUES ('txuser2', 'tx2@example.com', 'hash456')
	`)
	require.NoError(t, err, "should be able to insert within transaction")

	// Commit transaction
	err = tx.Commit(ctx)
	require.NoError(t, err, "should be able to commit transaction")

	// Verify data was persisted
	err = ts.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM shared.users WHERE username = 'txuser2'").Scan(&count)
	require.NoError(t, err, "should be able to query")
	assert.Equal(t, 1, count, "data should be persisted after commit")
}

func TestDatabaseConcurrentQueries(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)
	resetDatabase(t, ts)

	ctx := context.Background()

	// Insert test data
	for i := 0; i < 10; i++ {
		_, err := ts.DB.Pool.Exec(ctx, `
			INSERT INTO shared.users (username, email, password_hash)
			VALUES ($1, $2, $3)
		`, "user"+string(rune(i+'0')), "user"+string(rune(i+'0'))+"@example.com", "hash"+string(rune(i+'0')))
		require.NoError(t, err, "should be able to insert test data")
	}

	// Run concurrent queries
	const numQueries = 20
	results := make(chan error, numQueries)

	for i := 0; i < numQueries; i++ {
		go func() {
			var count int
			err := ts.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM shared.users").Scan(&count)
			if err != nil {
				results <- err
				return
			}
			if count != 10 {
				results <- assert.AnError
				return
			}
			results <- nil
		}()
	}

	// Collect results
	for i := 0; i < numQueries; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent query %d should succeed", i)
	}
}

func TestDatabaseConstraints(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)
	resetDatabase(t, ts)

	ctx := context.Background()

	// Insert test user
	_, err := ts.DB.Pool.Exec(ctx, `
		INSERT INTO shared.users (username, email, password_hash)
		VALUES ('testuser', 'test@example.com', 'hash123')
	`)
	require.NoError(t, err, "should be able to insert user")

	// Try to insert duplicate username (should fail due to unique constraint)
	_, err = ts.DB.Pool.Exec(ctx, `
		INSERT INTO shared.users (username, email, password_hash)
		VALUES ('testuser', 'other@example.com', 'hash456')
	`)
	assert.Error(t, err, "duplicate username should fail")
	assert.Contains(t, err.Error(), "unique", "error should mention unique constraint")

	// Try to insert duplicate email (should fail due to unique constraint)
	_, err = ts.DB.Pool.Exec(ctx, `
		INSERT INTO shared.users (username, email, password_hash)
		VALUES ('otheruser', 'test@example.com', 'hash789')
	`)
	assert.Error(t, err, "duplicate email should fail")
	assert.Contains(t, err.Error(), "unique", "error should mention unique constraint")
}

func TestDatabaseIndexes(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	ctx := context.Background()

	// Query for indexes on users table
	rows, err := ts.DB.Pool.Query(ctx, `
		SELECT indexname, indexdef
		FROM pg_indexes
		WHERE schemaname = 'shared' AND tablename = 'users'
	`)
	require.NoError(t, err, "should be able to query indexes")
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var indexName, indexDef string
		err := rows.Scan(&indexName, &indexDef)
		require.NoError(t, err, "should be able to scan index info")
		indexes = append(indexes, indexName)
		t.Logf("Found index: %s -> %s", indexName, indexDef)
	}

	// Verify primary key index exists
	assert.NotEmpty(t, indexes, "should have at least one index (primary key)")
}
