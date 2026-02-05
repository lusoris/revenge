// Package testutil provides testing utilities for database integration tests.
package testutil

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	// sharedContainer is the single PostgreSQL container shared across all tests
	sharedContainer     *postgres.PostgresContainer
	sharedContainerMu   sync.Mutex
	sharedContainerOnce sync.Once
	sharedContainerErr  error
	sharedContainerURL  string

	// containerTemplateReady tracks if the template database has been created
	containerTemplateReady     bool
	containerTemplateReadyOnce sync.Once
	containerTemplateReadyErr  error
)

const (
	containerUser       = "test"
	containerPassword   = "test"
	containerDB         = "testdb"
	containerTemplateDB = "revenge_template"
)

// startSharedContainer starts the shared PostgreSQL container (once)
func startSharedContainer(ctx context.Context) error {
	sharedContainerOnce.Do(func() {
		container, err := postgres.Run(ctx,
			"postgres:17-alpine",
			postgres.WithDatabase(containerDB),
			postgres.WithUsername(containerUser),
			postgres.WithPassword(containerPassword),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(60*time.Second),
			),
		)
		if err != nil {
			sharedContainerErr = fmt.Errorf("failed to start postgres container: %w", err)
			return
		}

		sharedContainer = container

		// Get connection URL
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			sharedContainerErr = fmt.Errorf("failed to get connection string: %w", err)
			return
		}
		sharedContainerURL = connStr
	})
	return sharedContainerErr
}

// createContainerTemplateDB creates the template database with all migrations applied (once)
func createContainerTemplateDB(ctx context.Context) error {
	containerTemplateReadyOnce.Do(func() {
		// Connect to admin database
		pool, err := pgxpool.New(ctx, sharedContainerURL)
		if err != nil {
			containerTemplateReadyErr = fmt.Errorf("failed to connect to postgres: %w", err)
			return
		}
		defer pool.Close()

		// Drop template if exists (in case of previous failed run)
		_, _ = pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", containerTemplateDB))

		// Create template database
		_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", containerTemplateDB))
		if err != nil {
			containerTemplateReadyErr = fmt.Errorf("failed to create template database: %w", err)
			return
		}

		// Get template database URL
		templateURL := replaceDBName(sharedContainerURL, containerTemplateDB)

		// Run migrations on template database
		containerTemplateReadyErr = runMigrationsOnURL(templateURL)
		if containerTemplateReadyErr != nil {
			containerTemplateReadyErr = fmt.Errorf("failed to run migrations on template: %w", containerTemplateReadyErr)
			return
		}

		containerTemplateReady = true
	})
	return containerTemplateReadyErr
}

// replaceDBName replaces the database name in a connection URL
func replaceDBName(connURL, newDB string) string {
	// Parse the URL to get the old database name
	cfg, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return connURL
	}
	oldDB := cfg.ConnConfig.Database

	// Replace the database name in the URL string
	// URL format: postgres://user:password@host:port/database?options
	return strings.Replace(connURL, "/"+oldDB, "/"+newDB, 1)
}

// FastTestDB provides a fast, parallel-safe database for testing using testcontainers.
// It creates databases from templates, making test setup ~10ms instead of seconds.
type FastTestDB struct {
	pool   *pgxpool.Pool
	dbName string
	t      testing.TB
}

// containerDBCounter generates unique database names
var containerDBCounter int64
var containerDBCounterMu sync.Mutex

func nextContainerDBName() string {
	containerDBCounterMu.Lock()
	defer containerDBCounterMu.Unlock()
	containerDBCounter++
	return fmt.Sprintf("test_%d", containerDBCounter)
}

// NewFastTestDB creates a new test database using testcontainers and template databases.
// This starts a PostgreSQL container (shared across tests) and creates isolated databases.
//
// Usage:
//
//	func TestSomething(t *testing.T) {
//	    t.Parallel() // Safe to run in parallel!
//	    db := testutil.NewFastTestDB(t)
//	    // Use db.Pool() for database operations
//	    // Database is automatically cleaned up when test ends
//	}
func NewFastTestDB(t testing.TB) *FastTestDB {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	ctx := context.Background()

	// Start shared container if not already running
	if err := startSharedContainer(ctx); err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Create template database if not already created
	if err := createContainerTemplateDB(ctx); err != nil {
		t.Fatalf("failed to create template database: %v", err)
	}

	// Generate unique database name
	dbName := nextContainerDBName()

	// Connect to admin database to create test database
	adminPool, err := pgxpool.New(ctx, sharedContainerURL)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	// Clone from template (this is instant!)
	_, err = adminPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", dbName, containerTemplateDB))
	adminPool.Close()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Connect to the new test database
	testURL := replaceDBName(sharedContainerURL, dbName)
	pool, err := pgxpool.New(ctx, testURL)
	if err != nil {
		t.Fatalf("failed to create pool for test database: %v", err)
	}

	testDB := &FastTestDB{
		pool:   pool,
		dbName: dbName,
		t:      t,
	}

	// Register cleanup
	t.Cleanup(func() {
		testDB.cleanup()
	})

	return testDB
}

// Pool returns the database connection pool for this test database.
func (db *FastTestDB) Pool() *pgxpool.Pool {
	return db.pool
}

// cleanup drops the test database
func (db *FastTestDB) cleanup() {
	if db.pool != nil {
		db.pool.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect as admin to drop the test database
	adminPool, err := pgxpool.New(ctx, sharedContainerURL)
	if err != nil {
		db.t.Logf("warning: failed to connect for cleanup: %v", err)
		return
	}
	defer adminPool.Close()

	// Force disconnect any remaining connections
	_, _ = adminPool.Exec(ctx, fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, db.dbName))

	// Drop the database
	_, err = adminPool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", db.dbName))
	if err != nil {
		db.t.Logf("warning: failed to drop test database %s: %v", db.dbName, err)
	}
}

// StopSharedContainer stops the shared PostgreSQL container.
// This should be called in TestMain if you want explicit cleanup.
func StopSharedContainer() {
	sharedContainerMu.Lock()
	defer sharedContainerMu.Unlock()

	if sharedContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = sharedContainer.Terminate(ctx)
		sharedContainer = nil
	}
}
