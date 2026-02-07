package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testLogger creates a test logger.
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestStatsReturnsAllFields tests that Stats returns all expected fields.
func TestStatsReturnsAllFields(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	pool, cleanup := setupTestPool(t)
	defer cleanup()

	stats := Stats(pool)

	expectedFields := []string{
		"acquire_count",
		"acquire_duration_ms",
		"acquired_conns",
		"canceled_acquire_count",
		"constructing_conns",
		"empty_acquire_count",
		"idle_conns",
		"max_conns",
		"total_conns",
		"new_conns_count",
		"max_lifetime_destroy_count",
		"max_idle_destroy_count",
	}

	for _, field := range expectedFields {
		_, exists := stats[field]
		assert.True(t, exists, "Stats should contain field: %s", field)
	}

	maxConns, ok := stats["max_conns"].(int32)
	assert.True(t, ok, "max_conns should be int32")
	assert.Greater(t, maxConns, int32(0), "max_conns should be > 0")
}

// TestStatsAfterQueries tests stats change after executing queries.
func TestStatsAfterQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	pool, cleanup := setupTestPool(t)
	defer cleanup()

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		var result int
		err := pool.QueryRow(ctx, "SELECT 1").Scan(&result)
		require.NoError(t, err)
	}

	stats := Stats(pool)

	acquireCount, ok := stats["acquire_count"].(int64)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, acquireCount, int64(5))
}

// TestHealthSuccess tests Health returns nil on healthy database.
func TestHealthSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	pool, cleanup := setupTestPool(t)
	defer cleanup()

	err := Health(context.Background(), pool)
	assert.NoError(t, err, "Health check should pass on healthy database")
}

// TestHealthWithCanceledContext tests Health with canceled context.
func TestHealthWithCanceledContext(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	pool, cleanup := setupTestPool(t)
	defer cleanup()

	canceledCtx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	err := Health(canceledCtx, pool)
	assert.Error(t, err, "Health should fail with canceled context")
}

// TestMigrateUpAndVersion tests MigrateUp and MigrateVersion.
func TestMigrateUpAndVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	logger := testLogger()

	err := MigrateUp(testURL, logger)
	require.NoError(t, err, "MigrateUp should succeed")

	version, dirty, err := MigrateVersion(testURL)
	require.NoError(t, err, "MigrateVersion should succeed")
	assert.False(t, dirty, "Migration should not be dirty")
	assert.Greater(t, version, uint(0), "Version should be > 0 after migrations")
}

// TestMigrateUpIdempotent tests that MigrateUp is idempotent.
func TestMigrateUpIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	logger := testLogger()

	err := MigrateUp(testURL, logger)
	require.NoError(t, err)

	err = MigrateUp(testURL, logger)
	require.NoError(t, err, "Running MigrateUp twice should succeed")
}

// TestMigrateDown tests MigrateDown rolls back one migration.
func TestMigrateDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	logger := testLogger()

	err := MigrateUp(testURL, logger)
	require.NoError(t, err)

	versionBefore, _, err := MigrateVersion(testURL)
	require.NoError(t, err)

	err = MigrateDown(testURL, logger)
	require.NoError(t, err, "MigrateDown should succeed")

	versionAfter, _, err := MigrateVersion(testURL)
	require.NoError(t, err)

	assert.Less(t, versionAfter, versionBefore, "Version should decrease after MigrateDown")
}

// TestMigrateTo tests MigrateTo migrates to specific version.
func TestMigrateTo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	logger := testLogger()

	err := MigrateTo(testURL, 1, logger)
	require.NoError(t, err, "MigrateTo should succeed")

	version, dirty, err := MigrateVersion(testURL)
	require.NoError(t, err)
	assert.Equal(t, uint(1), version, "Version should be 1")
	assert.False(t, dirty)
}

// TestMigrateVersionOnFreshDatabase tests MigrateVersion on fresh database.
func TestMigrateVersionOnFreshDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	version, dirty, err := MigrateVersion(testURL)
	require.NoError(t, err)
	assert.Equal(t, uint(0), version, "Fresh database should have version 0")
	assert.False(t, dirty)
}

// TestMigrateInvalidURL tests migration functions with invalid URL.
func TestMigrateInvalidURL(t *testing.T) {
	logger := testLogger()

	t.Run("MigrateUp with invalid URL", func(t *testing.T) {
		err := MigrateUp("invalid://url", logger)
		assert.Error(t, err, "MigrateUp should error with invalid URL")
	})

	t.Run("MigrateDown with invalid URL", func(t *testing.T) {
		err := MigrateDown("invalid://url", logger)
		assert.Error(t, err, "MigrateDown should error with invalid URL")
	})

	t.Run("MigrateVersion with invalid URL", func(t *testing.T) {
		_, _, err := MigrateVersion("invalid://url")
		assert.Error(t, err, "MigrateVersion should error with invalid URL")
	})

	t.Run("MigrateTo with invalid URL", func(t *testing.T) {
		err := MigrateTo("invalid://url", 1, logger)
		assert.Error(t, err, "MigrateTo should error with invalid URL")
	})
}

// TestMigrateConnectionRefused tests migration with unreachable database.
func TestMigrateConnectionRefused(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	logger := testLogger()
	unreachableURL := "postgres://user:pass@localhost:59998/test?sslmode=disable"

	err := MigrateUp(unreachableURL, logger)
	assert.Error(t, err, "Should error when database is unreachable")
}

// TestNewPoolIntegration tests NewPool with real database.
func TestNewPoolIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               testURL,
			MaxConns:          10,
			MinConns:          2,
			MaxConnLifetime:   time.Hour,
			MaxConnIdleTime:   30 * time.Minute,
			HealthCheckPeriod: time.Minute,
		},
	}
	logger := testLogger()

	pool, err := NewPool(cfg, logger)
	require.NoError(t, err, "NewPool should succeed")
	require.NotNil(t, pool)
	defer pool.Close()

	ctx := context.Background()
	var result int
	err = pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 1, result)
}

// TestNewPoolAndHealthCheck tests NewPool and Health together.
func TestNewPoolAndHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:      testURL,
			MaxConns: 5,
			MinConns: 1,
		},
	}
	logger := testLogger()

	pool, err := NewPool(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, pool)
	defer pool.Close()

	err = Health(context.Background(), pool)
	assert.NoError(t, err, "Health check should pass on newly created pool")

	stats := Stats(pool)
	assert.NotEmpty(t, stats)
}

// TestHealthQueryResultVerification tests that Health properly verifies query result.
func TestHealthQueryResultVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	pool, cleanup := setupTestPool(t)
	defer cleanup()

	ctx := context.Background()

	_, err := pool.Exec(ctx, `CREATE OR REPLACE FUNCTION evil_select_one() RETURNS int AS $$ BEGIN RETURN 2; END; $$ LANGUAGE plpgsql;`)
	require.NoError(t, err)

	err = Health(ctx, pool)
	assert.NoError(t, err)
}

// TestConcurrentHealthChecks tests Health is safe for concurrent use.
func TestConcurrentHealthChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	pool, cleanup := setupTestPool(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const numChecks = 20
	errChan := make(chan error, numChecks)

	for i := 0; i < numChecks; i++ {
		go func() {
			errChan <- Health(ctx, pool)
		}()
	}

	for i := 0; i < numChecks; i++ {
		err := <-errChan
		assert.NoError(t, err, "Concurrent health check should succeed")
	}
}

// TestMigrateUpVerifiesTables tests that MigrateUp actually creates tables.
func TestMigrateUpVerifiesTables(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	testURL, cleanup := freshTestDBURL(t)
	defer cleanup()

	logger := testLogger()

	err := MigrateUp(testURL, logger)
	require.NoError(t, err)

	db, err := sql.Open("pgx", testURL)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	var schemaExists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'shared')`).Scan(&schemaExists)
	require.NoError(t, err)
	assert.True(t, schemaExists, "shared schema should exist after migrations")

	var usersTableExists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'users')`).Scan(&usersTableExists)
	require.NoError(t, err)
	assert.True(t, usersTableExists, "shared.users table should exist after migrations")
}
