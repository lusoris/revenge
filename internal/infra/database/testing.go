package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

var (
	// sharedPG is the single embedded postgres shared across all database tests.
	sharedPG     *embeddedpostgres.EmbeddedPostgres
	sharedPGOnce sync.Once
	sharedPGErr  error
	sharedPGPort uint32 = 15600

	// testDBCounter generates unique database names.
	testDBCounter atomic.Int64
)

// startSharedPG starts the shared embedded postgres for the database package tests.
func startSharedPG() error {
	sharedPGOnce.Do(func() {
		runtimePath := fmt.Sprintf("/tmp/embedded-postgres-db-%d", os.Getpid())

		sharedPG = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
			Port(sharedPGPort).
			Username("test").
			Password("test").
			Database("postgres").
			RuntimePath(runtimePath).
			StartTimeout(60 * time.Second))

		sharedPGErr = sharedPG.Start()
	})
	return sharedPGErr
}

// StopSharedPG stops the shared embedded postgres. Call from TestMain.
func StopSharedPG() {
	if sharedPG != nil {
		_ = sharedPG.Stop()
	}
	runtimePath := fmt.Sprintf("/tmp/embedded-postgres-db-%d", os.Getpid())
	_ = os.RemoveAll(runtimePath)
}

// adminConnStr returns the connection string for the admin database.
func adminConnStr() string {
	return fmt.Sprintf("host=localhost port=%d user=test password=test dbname=postgres sslmode=disable", sharedPGPort)
}

// createTestDatabase creates a new database on the shared postgres and returns its name and cleanup function.
func createTestDatabase(t *testing.T, prefix string) (dbName string, cleanup func()) {
	t.Helper()

	err := startSharedPG()
	require.NoError(t, err, "failed to start shared embedded postgres")

	dbNum := testDBCounter.Add(1)
	dbName = fmt.Sprintf("%s_%d_%d", prefix, os.Getpid(), dbNum)

	adminConn := adminConnStr()
	adminDB, err := sql.Open("pgx", adminConn)
	require.NoError(t, err, "failed to connect to admin database")

	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err, "failed to create test database")
	_ = adminDB.Close()

	cleanup = func() {
		adminDB2, err := sql.Open("pgx", adminConn)
		if err == nil {
			_, _ = adminDB2.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", dbName))
			_ = adminDB2.Close()
		}
	}

	return dbName, cleanup
}

// dbURL returns a postgres URL for the given database name on the shared postgres.
func dbURL(dbName string) string {
	return fmt.Sprintf("postgres://test:test@localhost:%d/%s?sslmode=disable", sharedPGPort, dbName)
}

// dbDSN returns a DSN connection string for the given database name on the shared postgres.
func dbDSN(dbName string) string {
	return fmt.Sprintf("host=localhost port=%d user=test password=test dbname=%s sslmode=disable", sharedPGPort, dbName)
}

// setupTestDB creates a fresh database on the shared embedded postgres with all migrations applied.
func setupTestDB(t *testing.T, _ uint32) (*sql.DB, func()) {
	t.Helper()

	dbName, dropDB := createTestDatabase(t, "test")

	db, err := sql.Open("pgx", dbDSN(dbName))
	require.NoError(t, err, "failed to connect to test database")

	err = db.PingContext(context.Background())
	require.NoError(t, err, "database not ready")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err, "failed to create postgres driver")

	sourceDriver, err := iofs.New(migrationsFS, "migrations/shared")
	require.NoError(t, err, "failed to create source driver from embedded migrations")

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	require.NoError(t, err, "failed to create migrate instance")

	err = m.Up()
	require.NoError(t, err, "failed to run migrations")

	cleanup := func() {
		_ = db.Close()
		dropDB()
	}

	return db, cleanup
}

// setupFreshTestDB creates a fresh database WITHOUT migrations (for testing migration operations).
// Returns the sql.DB, a connection URL, and a cleanup function.
func setupFreshTestDB(t *testing.T) (*sql.DB, string, func()) {
	t.Helper()

	dbName, dropDB := createTestDatabase(t, "test_fresh")

	connURL := dbURL(dbName)
	db, err := sql.Open("pgx", dbDSN(dbName))
	require.NoError(t, err, "failed to connect to test database")

	err = db.PingContext(context.Background())
	require.NoError(t, err, "database not ready")

	cleanup := func() {
		_ = db.Close()
		dropDB()
	}

	return db, connURL, cleanup
}

// setupTestPool creates a fresh database on the shared embedded postgres and returns a pgxpool.Pool.
func setupTestPool(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	dbName, dropDB := createTestDatabase(t, "test_pool")

	pool, err := pgxpool.New(context.Background(), dbURL(dbName))
	require.NoError(t, err, "failed to create pool")

	cleanup := func() {
		pool.Close()
		dropDB()
	}

	return pool, cleanup
}

// freshTestDBURL creates a fresh database on the shared postgres and returns only the URL.
// Useful for testing functions that open their own connections (e.g., NewPool, MigrateUp).
func freshTestDBURL(t *testing.T) (string, func()) {
	t.Helper()

	dbName, dropDB := createTestDatabase(t, "test_url")
	return dbURL(dbName), dropDB
}
