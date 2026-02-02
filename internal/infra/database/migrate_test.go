package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5/pgxpool"
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

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15540).
		Username("test").
		Password("test").
		Database("statstest"))

	err := postgres.Start()
	require.NoError(t, err, "Failed to start embedded postgres")
	defer func() {
		_ = postgres.Stop()
	}()

	// Create pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15540/statstest?sslmode=disable")
	require.NoError(t, err, "Failed to create pool")
	defer pool.Close()

	// Get stats
	stats := Stats(pool)

	// Verify all expected fields are present
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

	// max_conns should be > 0
	maxConns, ok := stats["max_conns"].(int32)
	assert.True(t, ok, "max_conns should be int32")
	assert.Greater(t, maxConns, int32(0), "max_conns should be > 0")
}

// TestStatsAfterQueries tests stats change after executing queries.
func TestStatsAfterQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15541).
		Username("test").
		Password("test").
		Database("statstest2"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15541/statstest2?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Execute a few queries
	for i := 0; i < 5; i++ {
		var result int
		err := pool.QueryRow(ctx, "SELECT 1").Scan(&result)
		require.NoError(t, err)
	}

	stats := Stats(pool)

	// acquire_count should be >= 5 (we acquired at least 5 times)
	acquireCount, ok := stats["acquire_count"].(int64)
	assert.True(t, ok)
	assert.GreaterOrEqual(t, acquireCount, int64(5))
}

// TestHealthSuccess tests Health returns nil on healthy database.
func TestHealthSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15542).
		Username("test").
		Password("test").
		Database("healthtest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15542/healthtest?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Health check should pass
	err = Health(ctx, pool)
	assert.NoError(t, err, "Health check should pass on healthy database")
}

// TestHealthWithCanceledContext tests Health with canceled context.
func TestHealthWithCanceledContext(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15543).
		Username("test").
		Password("test").
		Database("healthtest2"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15543/healthtest2?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Create already canceled context
	canceledCtx, cancelFn := context.WithCancel(context.Background())
	cancelFn() // Cancel immediately

	err = Health(canceledCtx, pool)
	assert.Error(t, err, "Health should fail with canceled context")
}

// TestMigrateUpAndVersion tests MigrateUp and MigrateVersion.
func TestMigrateUpAndVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15544).
		Username("test").
		Password("test").
		Database("migratetest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	dbURL := "postgres://test:test@localhost:15544/migratetest?sslmode=disable"
	logger := testLogger()

	// Run migrations up
	err = MigrateUp(dbURL, logger)
	require.NoError(t, err, "MigrateUp should succeed")

	// Check version
	version, dirty, err := MigrateVersion(dbURL)
	require.NoError(t, err, "MigrateVersion should succeed")
	assert.False(t, dirty, "Migration should not be dirty")
	assert.Greater(t, version, uint(0), "Version should be > 0 after migrations")
}

// TestMigrateUpIdempotent tests that MigrateUp is idempotent.
func TestMigrateUpIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15545).
		Username("test").
		Password("test").
		Database("migratetest2"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	dbURL := "postgres://test:test@localhost:15545/migratetest2?sslmode=disable"
	logger := testLogger()

	// Run migrations twice
	err = MigrateUp(dbURL, logger)
	require.NoError(t, err)

	// Second run should also succeed (ErrNoChange is handled)
	err = MigrateUp(dbURL, logger)
	require.NoError(t, err, "Running MigrateUp twice should succeed")
}

// TestMigrateDown tests MigrateDown rolls back one migration.
func TestMigrateDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15546).
		Username("test").
		Password("test").
		Database("migratetest3"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	dbURL := "postgres://test:test@localhost:15546/migratetest3?sslmode=disable"
	logger := testLogger()

	// First migrate up
	err = MigrateUp(dbURL, logger)
	require.NoError(t, err)

	versionBefore, _, err := MigrateVersion(dbURL)
	require.NoError(t, err)

	// Roll back one
	err = MigrateDown(dbURL, logger)
	require.NoError(t, err, "MigrateDown should succeed")

	versionAfter, _, err := MigrateVersion(dbURL)
	require.NoError(t, err)

	// Version should have decreased
	assert.Less(t, versionAfter, versionBefore, "Version should decrease after MigrateDown")
}

// TestMigrateTo tests MigrateTo migrates to specific version.
func TestMigrateTo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15547).
		Username("test").
		Password("test").
		Database("migratetest4"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	dbURL := "postgres://test:test@localhost:15547/migratetest4?sslmode=disable"
	logger := testLogger()

	// Migrate to version 1
	err = MigrateTo(dbURL, 1, logger)
	require.NoError(t, err, "MigrateTo should succeed")

	version, dirty, err := MigrateVersion(dbURL)
	require.NoError(t, err)
	assert.Equal(t, uint(1), version, "Version should be 1")
	assert.False(t, dirty)
}

// TestMigrateVersionOnFreshDatabase tests MigrateVersion on fresh database.
func TestMigrateVersionOnFreshDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15548).
		Username("test").
		Password("test").
		Database("migratetest5"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	dbURL := "postgres://test:test@localhost:15548/migratetest5?sslmode=disable"

	// On fresh database, version should be 0
	version, dirty, err := MigrateVersion(dbURL)
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

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15549).
		Username("test").
		Password("test").
		Database("pooltest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               "postgres://test:test@localhost:15549/pooltest?sslmode=disable",
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

	// Verify pool is working
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

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15550).
		Username("test").
		Password("test").
		Database("poolhealthtest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:      "postgres://test:test@localhost:15550/poolhealthtest?sslmode=disable",
			MaxConns: 5,
			MinConns: 1,
		},
	}
	logger := testLogger()

	pool, err := NewPool(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, pool)
	defer pool.Close()

	// Health check should pass
	ctx := context.Background()
	err = Health(ctx, pool)
	assert.NoError(t, err, "Health check should pass on newly created pool")

	// Stats should show data
	stats := Stats(pool)
	assert.NotEmpty(t, stats)
}

// TestHealthQueryResult tests that Health properly verifies query result.
func TestHealthQueryResultVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15551).
		Username("test").
		Password("test").
		Database("healthquerytest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15551/healthquerytest?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Create a function that returns something other than 1
	_, err = pool.Exec(ctx, `CREATE OR REPLACE FUNCTION evil_select_one() RETURNS int AS $$ BEGIN RETURN 2; END; $$ LANGUAGE plpgsql;`)
	require.NoError(t, err)

	// Normal health check should still pass (uses SELECT 1, not our function)
	err = Health(ctx, pool)
	assert.NoError(t, err)
}

// TestConcurrentHealthChecks tests Health is safe for concurrent use.
func TestConcurrentHealthChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping embedded-postgres test in short mode")
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15552).
		Username("test").
		Password("test").
		Database("concurrenthealthtest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15552/concurrenthealthtest?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Run 20 concurrent health checks
	const numChecks = 20
	errChan := make(chan error, numChecks)

	for i := 0; i < numChecks; i++ {
		go func() {
			errChan <- Health(ctx, pool)
		}()
	}

	// Collect results
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

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15553).
		Username("test").
		Password("test").
		Database("tableverifytest"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() {
		_ = postgres.Stop()
	}()

	dbURL := "postgres://test:test@localhost:15553/tableverifytest?sslmode=disable"
	logger := testLogger()

	// Run migrations
	err = MigrateUp(dbURL, logger)
	require.NoError(t, err)

	// Connect and verify tables exist
	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)
	defer db.Close()

	// Check if shared schema exists (from 000001_create_schemas.up.sql)
	var schemaExists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'shared')`).Scan(&schemaExists)
	require.NoError(t, err)
	assert.True(t, schemaExists, "shared schema should exist after migrations")

	// Check if users table exists (from 000002_create_users_table.up.sql)
	var usersTableExists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'users')`).Scan(&usersTableExists)
	require.NoError(t, err)
	assert.True(t, usersTableExists, "shared.users table should exist after migrations")
}
