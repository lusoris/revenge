package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// columnSpec defines expected column type and nullability.
type columnSpec struct {
	dataType   string
	isNullable string
}

// assertColumnsExist verifies that all expected columns exist in the given table with correct types.
// Extra columns added by newer migrations are allowed and do not cause failures.
func assertColumnsExist(t *testing.T, db *sql.DB, schema, table string, expected map[string]columnSpec) {
	t.Helper()

	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position
	`, schema, table)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	actual := make(map[string]columnSpec)
	for rows.Next() {
		var colName, dataType, nullable string
		require.NoError(t, rows.Scan(&colName, &dataType, &nullable))
		actual[colName] = columnSpec{dataType, nullable}
	}

	for colName, exp := range expected {
		got, exists := actual[colName]
		if assert.True(t, exists, "missing column: %s.%s.%s", schema, table, colName) {
			assert.Equal(t, exp.dataType, got.dataType, "column %s has wrong type", colName)
			assert.Equal(t, exp.isNullable, got.isNullable, "column %s has wrong nullable", colName)
		}
	}
}

// newMigrateInstance creates a migrate instance using the embedded migration FS.
func newMigrateInstance(t *testing.T, sqlDB *sql.DB) *migrate.Migrate {
	t.Helper()
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	require.NoError(t, err, "failed to create migration driver")
	sourceDriver, err := iofs.New(migrationsFS, "migrations/shared")
	require.NoError(t, err, "failed to create iofs source")
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	require.NoError(t, err, "failed to create migrate instance")
	return m
}

func TestMigrationsUpDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping migration test in short mode")
	}

	db, _, cleanup := setupFreshTestDB(t)
	defer cleanup()

	m := newMigrateInstance(t, db)

	t.Run("MigrateUp", func(t *testing.T) {
		err := m.Up()
		require.NoError(t, err, "failed to migrate up")

		// Core schemas and tables must exist after full migration
		coreTables := []struct{ schema, table string }{
			{"shared", "users"},
			{"shared", "sessions"},
			{"shared", "server_settings"},
			{"shared", "user_settings"},
			{"shared", "user_preferences"},
			{"shared", "user_avatars"},
			{"shared", "auth_tokens"},
			{"shared", "password_reset_tokens"},
			{"shared", "email_verification_tokens"},
		}

		for _, ct := range coreTables {
			var exists bool
			err = db.QueryRow(
				"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2)",
				ct.schema, ct.table,
			).Scan(&exists)
			require.NoError(t, err)
			assert.True(t, exists, "%s.%s should exist", ct.schema, ct.table)
		}

		// Verify default settings were inserted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM shared.server_settings").Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "default settings should be inserted")

		var serverName string
		err = db.QueryRow("SELECT value::text FROM shared.server_settings WHERE key = 'server.name'").Scan(&serverName)
		require.NoError(t, err)
		assert.Equal(t, "\"Revenge Media Server\"", serverName)
	})

	t.Run("MigrateDown", func(t *testing.T) {
		versionBefore, _, err := m.Version()
		require.NoError(t, err)

		err = m.Steps(-1)
		require.NoError(t, err, "failed to migrate down one step")

		versionAfter, _, err := m.Version()
		require.NoError(t, err)
		assert.Less(t, versionAfter, versionBefore, "version should decrease after step down")

		// Core tables from early migrations should still exist
		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'users')").Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "users table should survive stepping down one migration")

		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'sessions')").Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "sessions table should survive stepping down one migration")
	})

	t.Run("MigrateUpAgain", func(t *testing.T) {
		err := m.Up()
		require.NoError(t, err, "failed to migrate up again")

		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'server_settings')").Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "server_settings should exist after re-migration")
	})
}

func TestServerSettingsTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "server_settings", map[string]columnSpec{
		"key":        {"character varying", "NO"},
		"value":      {"jsonb", "NO"},
		"data_type":  {"character varying", "NO"},
		"category":   {"character varying", "YES"},
		"is_secret":  {"boolean", "YES"},
		"is_public":  {"boolean", "YES"},
		"created_at": {"timestamp with time zone", "NO"},
		"updated_at": {"timestamp with time zone", "NO"},
	})
}

func TestUserSettingsTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "user_settings", map[string]columnSpec{
		"user_id":    {"uuid", "NO"},
		"key":        {"character varying", "NO"},
		"value":      {"jsonb", "NO"},
		"data_type":  {"character varying", "NO"},
		"created_at": {"timestamp with time zone", "NO"},
		"updated_at": {"timestamp with time zone", "NO"},
	})
}

func TestServerSettingsDefaultValues(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping default values test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	type setting struct {
		key      string
		category string
		dataType string
		isPublic bool
	}

	expectedSettings := []setting{
		{"server.name", "general", "string", true},
		{"auth.jwt.access_token_expiry", "auth", "number", false},
		{"auth.jwt.refresh_token_expiry", "auth", "number", false},
		{"auth.password.min_length", "auth", "number", false},
		{"auth.session.max_per_user", "auth", "number", false},
		{"features.registration_enabled", "features", "boolean", true},
		{"features.oidc_enabled", "features", "boolean", true},
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
			assert.True(t, value.Valid, "value should not be null")
		})
	}
}

func TestUserPreferencesTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping user_preferences test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "user_preferences", map[string]columnSpec{
		"user_id":    {"uuid", "NO"},
		"theme":      {"character varying", "YES"},
		"created_at": {"timestamp with time zone", "NO"},
		"updated_at": {"timestamp with time zone", "NO"},
	})
}

func TestUserAvatarsTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping user_avatars test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "user_avatars", map[string]columnSpec{
		"id":              {"uuid", "NO"},
		"user_id":         {"uuid", "NO"},
		"file_path":       {"text", "NO"},
		"file_size_bytes": {"bigint", "NO"},
		"mime_type":       {"character varying", "NO"},
		"width":           {"integer", "NO"},
		"height":          {"integer", "NO"},
		"is_current":      {"boolean", "YES"},
		"created_at":      {"timestamp with time zone", "NO"},
		"updated_at":      {"timestamp with time zone", "NO"},
	})
}
