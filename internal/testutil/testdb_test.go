package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	// Cleanup shared postgres after all tests
	StopSharedPostgres()
	os.Exit(code)
}

func TestNewTestDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	db := NewTestDB(t)
	require.NotNil(t, db)
	require.NotNil(t, db.Pool())

	// Verify we can query the database
	ctx := context.Background()
	var result int
	err := db.Pool().QueryRow(ctx, "SELECT 1").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	// Verify migrations were applied (shared schema should exist)
	var schemaExists bool
	err = db.Pool().QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM information_schema.schemata 
			WHERE schema_name = 'shared'
		)
	`).Scan(&schemaExists)
	require.NoError(t, err)
	assert.True(t, schemaExists, "shared schema should exist after migrations")

	// Verify users table exists in shared schema
	var tableExists bool
	err = db.Pool().QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = 'shared' AND table_name = 'users'
		)
	`).Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists, "shared.users table should exist after migrations")
}

func TestNewTestDB_Parallel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	// Run multiple tests in parallel to verify isolation
	for i := 0; i < 5; i++ {
		t.Run("parallel", func(t *testing.T) {
			t.Parallel()
			db := NewTestDB(t)
			require.NotNil(t, db)

			ctx := context.Background()

			// Insert a user unique to this test (in shared.users)
			_, err := db.Pool().Exec(ctx, `
				INSERT INTO shared.users (id, username, email, password_hash, created_at, updated_at)
				VALUES (gen_random_uuid(), $1, $2, 'hash', NOW(), NOW())
			`, t.Name(), t.Name()+"@test.com")
			require.NoError(t, err)

			// Verify only this test's user exists
			var count int
			err = db.Pool().QueryRow(ctx, `
				SELECT COUNT(*) FROM shared.users WHERE username = $1
			`, t.Name()).Scan(&count)
			require.NoError(t, err)
			assert.Equal(t, 1, count, "should have exactly one user from this test")
		})
	}
}

func TestNewTestDB_URL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	db := NewTestDB(t)
	url := db.URL()
	require.NotEmpty(t, url)
	assert.Contains(t, url, "postgres://")
	assert.Contains(t, url, db.dbName)
}
