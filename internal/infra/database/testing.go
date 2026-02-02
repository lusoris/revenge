package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary PostgreSQL database for testing
func setupTestDB(t *testing.T, port uint32) (*sql.DB, func()) {
	t.Helper()

	// Use embeddedpostgres for tests
	pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(port).
		Database("revenge_test").
		Username("revenge_test").
		Password("revenge_test"))

	err := pg.Start()
	require.NoError(t, err, "failed to start embeddedpostgres")

	// Connect to database
	connStr := fmt.Sprintf("host=localhost port=%d user=revenge_test password=revenge_test dbname=revenge_test sslmode=disable", port)
	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err, "failed to connect to test database")

	// Wait for database to be ready
	err = db.PingContext(context.Background())
	require.NoError(t, err, "database not ready")

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err, "failed to create postgres driver")

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres",
		driver,
	)
	require.NoError(t, err, "failed to create migrate instance")

	err = m.Up()
	require.NoError(t, err, "failed to run migrations")

	cleanup := func() {
		_ = db.Close()
		_ = pg.Stop()
	}

	return db, cleanup
}
