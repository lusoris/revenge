package testutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/logging"
)

const (
	// TestDBPort is the port for the embedded PostgreSQL instance
	TestDBPort uint32 = 15432
	// TestDBUser is the default test database user
	TestDBUser = "revenge_test"
	// TestDBPassword is the default test database password
	TestDBPassword = "test_pass"
	// TestDBName is the default test database name
	TestDBName = "revenge_test"
)

// TestDatabase represents an embedded PostgreSQL instance for testing.
type TestDatabase struct {
	postgres *embeddedpostgres.EmbeddedPostgres
	Pool     *pgxpool.Pool
	URL      string
	Config   *config.Config
	Logger   *slog.Logger
	port     uint32
}

// NewTestDatabase creates and starts an embedded PostgreSQL instance.
// It automatically runs migrations and returns a configured database pool.
//
// Example:
//
//	func TestMyFunction(t *testing.T) {
//	    db := testutil.NewTestDatabase(t)
//	    defer db.Close()
//
//	    // Use db.Pool for testing
//	}
func NewTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()

	// Create logger for test output
	logger := logging.NewLogger(logging.Config{
		Level:       "debug",
		Format:      "text",
		Development: true,
	})

	// Start embedded PostgreSQL
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(TestDBPort).
		Username(TestDBUser).
		Password(TestDBPassword).
		Database(TestDBName).
		Logger(nil)) // Suppress embedded-postgres logs in tests

	if err := postgres.Start(); err != nil {
		t.Fatalf("failed to start embedded postgres: %v", err)
	}

	// Build connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		TestDBUser, TestDBPassword, TestDBPort, TestDBName)

	// Run migrations
	if err := runMigrations(dbURL); err != nil {
		_ = postgres.Stop() // Best-effort cleanup
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               dbURL,
			MaxConns:          5,
			MinConns:          1,
			MaxConnLifetime:   300000000000, // 5m
			MaxConnIdleTime:   60000000000,  // 1m
			HealthCheckPeriod: 30000000000,  // 30s
		},
		Logging: config.LoggingConfig{
			Level:       "debug",
			Format:      "text",
			Development: true,
		},
	}

	// Create database pool (NewPool creates its own context internally)
	pool, err := database.NewPool(cfg, logger)
	if err != nil {
		_ = postgres.Stop() // Best-effort cleanup
		t.Fatalf("failed to create database pool: %v", err)
	}

	return &TestDatabase{
		postgres: postgres,
		Pool:     pool,
		URL:      dbURL,
		Config:   cfg,
		Logger:   logger,
		port:     TestDBPort,
	}
}

// Close stops the embedded PostgreSQL instance and closes the pool.
func (db *TestDatabase) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.postgres != nil {
		_ = db.postgres.Stop() // Best-effort cleanup
	}
}

// Reset truncates all tables in the test database, preserving schema.
// Useful for cleaning up between test cases.
func (db *TestDatabase) Reset(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Truncate all tables in shared schema
	_, err := db.Pool.Exec(ctx, `
		TRUNCATE TABLE shared.sessions CASCADE;
		TRUNCATE TABLE shared.users CASCADE;
	`)
	if err != nil {
		t.Fatalf("failed to reset database: %v", err)
	}
}

// runMigrations runs all migrations from the migrations directory.
func runMigrations(dbURL string) error {
	// Find migrations directory
	migrationsPath := findMigrationsPath()

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close() //nolint:errcheck // Deferred cleanup, error not actionable

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// findMigrationsPath finds the migrations directory relative to the test file.
func findMigrationsPath() string {
	// Try common relative paths to internal/infra/database/migrations/shared/
	paths := []string{
		"internal/infra/database/migrations/shared",
		"../internal/infra/database/migrations/shared",
		"../../internal/infra/database/migrations/shared",
		"../../../internal/infra/database/migrations/shared",
		"../../../../internal/infra/database/migrations/shared",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "internal/infra/database/migrations/shared" // Default fallback
}
