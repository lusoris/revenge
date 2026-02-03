// Package testutil provides testing utilities for database integration tests.
package testutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/validate"
)

// TestDB provides a fast, parallel-safe database for testing using the template database pattern.
// Instead of starting a new PostgreSQL instance for each test, it:
// 1. Starts one shared embedded PostgreSQL instance
// 2. Creates a template database with all migrations applied
// 3. For each test, instantly clones the template database (CREATE DATABASE ... TEMPLATE)
// 4. Drops the test database after the test completes
//
// This approach reduces test database setup from ~3-5s to ~10ms per test.
type TestDB struct {
	pool     *pgxpool.Pool
	dbName   string
	adminURL string
	t        testing.TB
}

var (
	// sharedPostgres is the single embedded PostgreSQL instance shared across all tests
	sharedPostgres     *embeddedpostgres.EmbeddedPostgres
	sharedPostgresMu   sync.Mutex
	sharedPostgresOnce sync.Once
	sharedPostgresErr  error

	// templateReady tracks if the template database has been created
	templateReady     atomic.Bool
	templateReadyOnce sync.Once
	templateReadyErr  error

	// dbCounter generates unique database names
	dbCounter atomic.Int64

	// Configuration
	sharedPort     = 15555
	adminUser      = "postgres"
	adminPassword  = "postgres"
	templateDBName = "revenge_template"
)

// startSharedPostgres starts the shared PostgreSQL instance (once)
func startSharedPostgres() error {
	sharedPostgresOnce.Do(func() {
		// Kill any orphaned postgres processes on the test port before starting
		cleanupOrphanedPostgres()

		// Use a unique RuntimePath per process to avoid conflicts
		runtimePath := fmt.Sprintf("/tmp/embedded-postgres-%d", os.Getpid())

		// Safe conversion of test port to uint32
		testPort, err := validate.SafeUint32(sharedPort)
		if err != nil {
			sharedPostgresErr = fmt.Errorf("invalid test port %d: %w", sharedPort, err)
			return
		}

		sharedPostgres = embeddedpostgres.NewDatabase(
			embeddedpostgres.DefaultConfig().
				Port(testPort).
				Username(adminUser).
				Password(adminPassword).
				Database("postgres").
				RuntimePath(runtimePath).
				StartTimeout(60*time.Second),
		)

		sharedPostgresErr = sharedPostgres.Start()
		if sharedPostgresErr == nil {
			// Register signal handlers to ensure cleanup on interrupt
			registerCleanupHandlers()
		}
	})
	return sharedPostgresErr
}

// createTemplateDB creates the template database with all migrations applied (once)
func createTemplateDB() error {
	templateReadyOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		adminURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/postgres?sslmode=disable",
			adminUser, adminPassword, sharedPort)

		// Connect as admin
		conn, err := pgx.Connect(ctx, adminURL)
		if err != nil {
			templateReadyErr = fmt.Errorf("failed to connect to postgres: %w", err)
			return
		}
		defer func() { _ = conn.Close(ctx) }()

		// Drop template if exists (in case of previous failed run)
		_, _ = conn.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", templateDBName))

		// Create template database
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", templateDBName))
		if err != nil {
			templateReadyErr = fmt.Errorf("failed to create template database: %w", err)
			return
		}

		// Connect to template database and run migrations
		templateURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
			adminUser, adminPassword, sharedPort, templateDBName)

		// Run migrations using the embedded migrations
		templateReadyErr = runMigrationsOnURL(templateURL)
		if templateReadyErr != nil {
			templateReadyErr = fmt.Errorf("failed to run migrations on template: %w", templateReadyErr)
			return
		}

		// Verify migrations were applied
		templateConn, err := pgx.Connect(ctx, templateURL)
		if err != nil {
			templateReadyErr = fmt.Errorf("failed to connect to verify template: %w", err)
			return
		}
		defer func() { _ = templateConn.Close(ctx) }()

		var version int
		err = templateConn.QueryRow(ctx, "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
		if err != nil {
			templateReadyErr = fmt.Errorf("failed to verify migration version: %w", err)
			return
		}

		templateReady.Store(true)
	})
	return templateReadyErr
}

// runMigrationsOnURL runs migrations on the given database URL
func runMigrationsOnURL(databaseURL string) error {
	// Import and use the database.MigrateUp function
	// We need to avoid circular imports, so we'll use the migrate package directly
	return runMigrationsWithMigrate(databaseURL)
}

// NewTestDB creates a new test database by cloning the template database.
// This is extremely fast (~10ms) compared to running migrations (~3-5s).
//
// Usage:
//
//	func TestSomething(t *testing.T) {
//	    t.Parallel() // Safe to run in parallel!
//	    db := testutil.NewTestDB(t)
//	    // Use db.Pool() for database operations
//	    // Database is automatically cleaned up when test ends
//	}
func NewTestDB(t testing.TB) *TestDB {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	// Start shared postgres if not already running
	if err := startSharedPostgres(); err != nil {
		t.Fatalf("failed to start shared postgres: %v", err)
	}

	// Create template database if not already created
	if err := createTemplateDB(); err != nil {
		t.Fatalf("failed to create template database: %v", err)
	}

	// Generate unique database name
	dbNum := dbCounter.Add(1)
	dbName := fmt.Sprintf("test_%d_%d", os.Getpid(), dbNum)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	adminURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/postgres?sslmode=disable",
		adminUser, adminPassword, sharedPort)

	// Connect as admin to create the test database
	adminConn, err := pgx.Connect(ctx, adminURL)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	// Clone from template (this is instant!)
	_, err = adminConn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", dbName, templateDBName))
	_ = adminConn.Close(ctx)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Connect to the new test database
	testURL := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		adminUser, adminPassword, sharedPort, dbName)

	pool, err := pgxpool.New(ctx, testURL)
	if err != nil {
		t.Fatalf("failed to create pool for test database: %v", err)
	}

	testDB := &TestDB{
		pool:     pool,
		dbName:   dbName,
		adminURL: adminURL,
		t:        t,
	}

	// Register cleanup
	t.Cleanup(func() {
		testDB.cleanup()
	})

	return testDB
}

// Pool returns the database connection pool for this test database.
func (db *TestDB) Pool() *pgxpool.Pool {
	return db.pool
}

// URL returns the connection URL for this test database.
func (db *TestDB) URL() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		adminUser, adminPassword, sharedPort, db.dbName)
}

// cleanup drops the test database
func (db *TestDB) cleanup() {
	if db.pool != nil {
		db.pool.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect as admin to drop the test database
	adminConn, err := pgx.Connect(ctx, db.adminURL)
	if err != nil {
		db.t.Logf("warning: failed to connect for cleanup: %v", err)
		return
	}
	defer func() { _ = adminConn.Close(ctx) }()

	// Force disconnect any remaining connections
	_, _ = adminConn.Exec(ctx, fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, db.dbName))

	// Drop the database
	_, err = adminConn.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", db.dbName))
	if err != nil {
		db.t.Logf("warning: failed to drop test database %s: %v", db.dbName, err)
	}
}

// StopSharedPostgres stops the shared PostgreSQL instance.
// This should be called in TestMain if you want explicit cleanup.
// If not called, the instance will be cleaned up when the process exits.
func StopSharedPostgres() {
	sharedPostgresMu.Lock()
	defer sharedPostgresMu.Unlock()

	if sharedPostgres != nil {
		_ = sharedPostgres.Stop()
		sharedPostgres = nil
	}
}

// cleanupOrphanedPostgres kills any postgres processes listening on the test port
func cleanupOrphanedPostgres() {
	// Validate port is in safe range (to prevent command injection even though it's a constant)
	if sharedPort < 1024 || sharedPort > 65535 {
		return // Invalid port, skip cleanup
	}

	// Kill any process listening on our port
	// Safe: sharedPort is validated and only numeric
	// #nosec G204 -- port is validated constant, not user input
	cmd := exec.Command("sh", "-c", fmt.Sprintf("lsof -ti:%d 2>/dev/null | xargs -r kill -9 2>/dev/null || true", sharedPort))
	_ = cmd.Run()

	// Kill any postgres processes from embedded-postgres
	cmd = exec.Command("sh", "-c", "pkill -9 -f 'embedded-postgres.*postgres' 2>/dev/null || true")
	_ = cmd.Run()

	// Kill any postgres processes in /tmp/embedded-postgres-*
	cmd = exec.Command("sh", "-c", "ps aux | grep '[/]tmp/embedded-postgres' | awk '{print $2}' | xargs -r kill -9 2>/dev/null || true")
	_ = cmd.Run()

	// Give the OS time to release the port
	time.Sleep(200 * time.Millisecond)
}

var cleanupRegistered atomic.Bool

// registerCleanupHandlers registers signal handlers to cleanup postgres on interrupt
func registerCleanupHandlers() {
	if cleanupRegistered.Swap(true) {
		return // Already registered
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		StopSharedPostgres()
		os.Exit(1)
	}()
}
