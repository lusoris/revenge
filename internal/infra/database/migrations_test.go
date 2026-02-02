package database

import (
	"context"
	"database/sql"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationsUpDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping migration test in short mode")
	}

	// Start embedded PostgreSQL
	embeddedPG := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15432).
		Username("postgres").
		Password("postgres").
		Database("test_migrations"))

	err := embeddedPG.Start()
	require.NoError(t, err, "failed to start embedded postgres")
	defer func() {
		_ = embeddedPG.Stop()
	}()

	// Connect to database
	db, err := sql.Open("postgres", "host=localhost port=15432 user=postgres password=postgres dbname=test_migrations sslmode=disable")
	require.NoError(t, err, "failed to connect to database")
	defer func() { _ = db.Close() }()

	// Wait for connection
	err = db.Ping()
	require.NoError(t, err, "failed to ping database")

	// Create migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err, "failed to create migration driver")

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../migrations",
		"postgres",
		driver,
	)
	require.NoError(t, err, "failed to create migrate instance")

	// Test: Migrate UP
	t.Run("MigrateUp", func(t *testing.T) {
		err := m.Up()
		require.NoError(t, err, "failed to migrate up")

		// Verify all tables exist
		var tableExists bool

		// Check shared schema
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.schemata WHERE schema_name = 'shared')").Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "shared schema should exist")

		// Check users table
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'users')").Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "users table should exist")

		// Check sessions table
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'sessions')").Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "sessions table should exist")

		// Check server_settings table (our new migration)
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'server_settings')").Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "server_settings table should exist")

		// Verify default settings were inserted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM shared.server_settings").Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "default settings should be inserted")

		// Verify specific default setting
		var serverName string
		err = db.QueryRow("SELECT value::text FROM shared.server_settings WHERE key = 'server.name'").Scan(&serverName)
		require.NoError(t, err)
		assert.Equal(t, "\"Revenge Media Server\"", serverName, "default server name should be set")
	})

	// Test: Migrate DOWN one step
	t.Run("MigrateDown", func(t *testing.T) {
		err := m.Steps(-1)
		require.NoError(t, err, "failed to migrate down one step")

		// Verify server_settings table is gone
		var tableExists bool
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'server_settings')").Scan(&tableExists)
		require.NoError(t, err)
		assert.False(t, tableExists, "server_settings table should not exist after down migration")

		// Other tables should still exist
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'sessions')").Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "sessions table should still exist")
	})

	// Test: Migrate UP again
	t.Run("MigrateUpAgain", func(t *testing.T) {
		err := m.Up()
		require.NoError(t, err, "failed to migrate up again")

		// Verify table exists again
		var tableExists bool
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'server_settings')").Scan(&tableExists)
		require.NoError(t, err)
		assert.True(t, tableExists, "server_settings table should exist again")
	})
}

func TestServerSettingsTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	// Start embedded PostgreSQL
	embeddedPG := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15433).
		Username("postgres").
		Password("postgres").
		Database("test_structure"))

	err := embeddedPG.Start()
	require.NoError(t, err)
	defer func() {
		_ = embeddedPG.Stop()
	}()

	// Connect
	db, err := sql.Open("postgres", "host=localhost port=15433 user=postgres password=postgres dbname=test_structure sslmode=disable")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	require.NoError(t, db.Ping())

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance("file://../../../migrations", "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	// Test column existence and types
	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'shared' AND table_name = 'server_settings'
		ORDER BY ordinal_position
	`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	expectedColumns := map[string]struct {
		dataType   string
		isNullable string
	}{
		"key":            {"character varying", "NO"},
		"value":          {"jsonb", "NO"},
		"description":    {"text", "YES"},
		"category":       {"character varying", "YES"},
		"data_type":      {"character varying", "NO"},
		"is_secret":      {"boolean", "YES"},
		"is_public":      {"boolean", "YES"},
		"allowed_values": {"jsonb", "YES"},
		"min_value":      {"numeric", "YES"},
		"max_value":      {"numeric", "YES"},
		"pattern":        {"text", "YES"},
		"created_at":     {"timestamp with time zone", "NO"},
		"updated_at":     {"timestamp with time zone", "NO"},
		"updated_by":     {"uuid", "YES"},
	}

	foundColumns := make(map[string]bool)
	for rows.Next() {
		var colName, dataType, nullable string
		err := rows.Scan(&colName, &dataType, &nullable)
		require.NoError(t, err)

		expected, exists := expectedColumns[colName]
		assert.True(t, exists, "unexpected column: %s", colName)
		if exists {
			assert.Equal(t, expected.dataType, dataType, "column %s has wrong type", colName)
			assert.Equal(t, expected.isNullable, nullable, "column %s has wrong nullable", colName)
		}
		foundColumns[colName] = true
	}

	// Verify all expected columns were found
	for colName := range expectedColumns {
		assert.True(t, foundColumns[colName], "missing column: %s", colName)
	}
}

func TestServerSettingsDefaultValues(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping default values test in short mode")
	}

	embeddedPG := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15434).
		Username("postgres").
		Password("postgres").
		Database("test_defaults"))

	err := embeddedPG.Start()
	require.NoError(t, err)
	defer func() {
		_ = embeddedPG.Stop()
	}()

	db, err := sql.Open("postgres", "host=localhost port=15434 user=postgres password=postgres dbname=test_defaults sslmode=disable")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	require.NoError(t, db.Ping())

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance("file://../../../migrations", "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	// Test default values
	type setting struct {
		key      string
		category string
		dataType string
		isPublic bool
		hasValue bool
	}

	expectedSettings := []setting{
		{"server.name", "general", "string", true, true},
		{"auth.jwt.access_token_expiry", "auth", "number", false, true},
		{"auth.jwt.refresh_token_expiry", "auth", "number", false, true},
		{"auth.password.min_length", "auth", "number", false, true},
		{"auth.session.max_per_user", "auth", "number", false, true},
		{"features.registration_enabled", "features", "boolean", true, true},
		{"features.oidc_enabled", "features", "boolean", true, true},
	}

	for _, expected := range expectedSettings {
		t.Run(expected.key, func(t *testing.T) {
			var category, dataType string
			var isPublic bool
			var value sql.NullString

			err := db.QueryRowContext(context.Background(),
				"SELECT category, data_type, is_public, value FROM shared.server_settings WHERE key = $1",
				expected.key,
			).Scan(&category, &dataType, &isPublic, &value)

			require.NoError(t, err, "setting %s should exist", expected.key)
			assert.Equal(t, expected.category, category)
			assert.Equal(t, expected.dataType, dataType)
			assert.Equal(t, expected.isPublic, isPublic)
			if expected.hasValue {
				assert.True(t, value.Valid, "value should not be null")
			}
		})
	}
}
