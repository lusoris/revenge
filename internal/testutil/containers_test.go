package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPostgreSQLContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Start container
	pg := NewPostgreSQLContainer(t)
	defer pg.Close()

	// Verify container is running and pool is connected
	require.NotNil(t, pg.container, "container should not be nil")
	require.NotNil(t, pg.Pool, "pool should not be nil")
	require.NotEmpty(t, pg.URL, "URL should not be empty")
	require.NotNil(t, pg.Config, "config should not be nil")

	// Test database connectivity
	ctx := context.Background()
	err := pg.Pool.Ping(ctx)
	require.NoError(t, err, "should be able to ping database")

	// Verify migrations ran successfully
	var schemaExists bool
	err = pg.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.schemata
			WHERE schema_name = 'shared'
		)
	`).Scan(&schemaExists)
	require.NoError(t, err, "should be able to query schema existence")
	assert.True(t, schemaExists, "shared schema should exist after migrations")

	// Verify users table exists
	var tableExists bool
	err = pg.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'shared' AND table_name = 'users'
		)
	`).Scan(&tableExists)
	require.NoError(t, err, "should be able to query table existence")
	assert.True(t, tableExists, "users table should exist after migrations")
}

func TestPostgreSQLContainer_Reset(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Start container
	pg := NewPostgreSQLContainer(t)
	defer pg.Close()

	ctx := context.Background()

	// Insert test data
	_, err := pg.Pool.Exec(ctx, `
		INSERT INTO shared.users (username, email, password_hash)
		VALUES ('testuser', 'test@example.com', 'hash123')
	`)
	require.NoError(t, err, "should be able to insert test data")

	// Verify data exists
	var count int
	err = pg.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM shared.users").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "should have 1 user before reset")

	// Reset database
	pg.Reset(t)

	// Verify data is gone
	err = pg.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM shared.users").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "should have 0 users after reset")

	// Verify schema still exists
	var schemaExists bool
	err = pg.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.schemata
			WHERE schema_name = 'shared'
		)
	`).Scan(&schemaExists)
	require.NoError(t, err)
	assert.True(t, schemaExists, "shared schema should still exist after reset")
}

func TestPostgreSQLContainer_MultipleConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Start container
	pg := NewPostgreSQLContainer(t)
	defer pg.Close()

	ctx := context.Background()

	// Acquire multiple connections
	conn1, err := pg.Pool.Acquire(ctx)
	require.NoError(t, err, "should be able to acquire first connection")
	defer conn1.Release()

	conn2, err := pg.Pool.Acquire(ctx)
	require.NoError(t, err, "should be able to acquire second connection")
	defer conn2.Release()

	// Verify both connections work
	err = conn1.Ping(ctx)
	assert.NoError(t, err, "first connection should work")

	err = conn2.Ping(ctx)
	assert.NoError(t, err, "second connection should work")
}

func TestPostgreSQLContainer_Configuration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Start container
	pg := NewPostgreSQLContainer(t)
	defer pg.Close()

	// Verify configuration
	assert.Equal(t, "debug", pg.Config.Logging.Level, "logging level should be debug")
	assert.Equal(t, "text", pg.Config.Logging.Format, "logging format should be text")
	assert.True(t, pg.Config.Logging.Development, "development mode should be enabled")

	// Verify pool statistics (actual values from the pool, not input config)
	stats := pg.Pool.Stat()
	assert.Equal(t, int32(10), stats.MaxConns(), "max connections should be 10")
	// MinConns is not directly exposed in pgxpool.Stat, so verify via config
	assert.Equal(t, pg.URL, pg.Config.Database.URL, "database URL should match")
}

func TestPostgreSQLContainer_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Start container
	pg := NewPostgreSQLContainer(t)

	// Verify pool is open
	ctx := context.Background()
	err := pg.Pool.Ping(ctx)
	require.NoError(t, err, "pool should be open before close")

	// Close container
	pg.Close()

	// Verify pool is closed (ping should fail)
	err = pg.Pool.Ping(ctx)
	assert.Error(t, err, "pool should be closed after close")
}
