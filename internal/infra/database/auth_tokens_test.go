package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthTokensTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 15437)
	defer cleanup()

	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'auth_tokens')").Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists)

	rows, err := db.Query(`SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema = 'shared' AND table_name = 'auth_tokens' ORDER BY ordinal_position`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	expected := map[string]struct {
		dataType   string
		isNullable string
	}{
		"id":                 {"uuid", "NO"},
		"user_id":            {"uuid", "NO"},
		"token_hash":         {"text", "NO"},
		"token_type":         {"character varying", "NO"},
		"device_name":        {"character varying", "YES"},
		"device_fingerprint": {"text", "YES"},
		"ip_address":         {"inet", "YES"},
		"user_agent":         {"text", "YES"},
		"expires_at":         {"timestamp with time zone", "NO"},
		"revoked_at":         {"timestamp with time zone", "YES"},
		"last_used_at":       {"timestamp with time zone", "YES"},
		"created_at":         {"timestamp with time zone", "NO"},
		"updated_at":         {"timestamp with time zone", "NO"},
	}

	found := make(map[string]bool)
	for rows.Next() {
		var colName, dataType, nullable string
		require.NoError(t, rows.Scan(&colName, &dataType, &nullable))
		exp, exists := expected[colName]
		assert.True(t, exists)
		if exists {
			assert.Equal(t, exp.dataType, dataType)
			assert.Equal(t, exp.isNullable, nullable)
		}
		found[colName] = true
	}
	for colName := range expected {
		assert.True(t, found[colName])
	}
}

func TestPasswordResetTokensTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 15438)
	defer cleanup()

	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'password_reset_tokens')").Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists)

	rows, err := db.Query(`SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema = 'shared' AND table_name = 'password_reset_tokens' ORDER BY ordinal_position`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	expected := map[string]struct {
		dataType   string
		isNullable string
	}{
		"id":         {"uuid", "NO"},
		"user_id":    {"uuid", "NO"},
		"token_hash": {"text", "NO"},
		"ip_address": {"inet", "YES"},
		"user_agent": {"text", "YES"},
		"expires_at": {"timestamp with time zone", "NO"},
		"used_at":    {"timestamp with time zone", "YES"},
		"created_at": {"timestamp with time zone", "NO"},
	}

	found := make(map[string]bool)
	for rows.Next() {
		var colName, dataType, nullable string
		require.NoError(t, rows.Scan(&colName, &dataType, &nullable))
		exp, exists := expected[colName]
		assert.True(t, exists)
		if exists {
			assert.Equal(t, exp.dataType, dataType)
			assert.Equal(t, exp.isNullable, nullable)
		}
		found[colName] = true
	}
	for colName := range expected {
		assert.True(t, found[colName])
	}
}

func TestEmailVerificationTokensTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 15439)
	defer cleanup()

	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'shared' AND table_name = 'email_verification_tokens')").Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists)

	rows, err := db.Query(`SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema = 'shared' AND table_name = 'email_verification_tokens' ORDER BY ordinal_position`)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	expected := map[string]struct {
		dataType   string
		isNullable string
	}{
		"id":          {"uuid", "NO"},
		"user_id":     {"uuid", "NO"},
		"token_hash":  {"text", "NO"},
		"email":       {"character varying", "NO"},
		"ip_address":  {"inet", "YES"},
		"user_agent":  {"text", "YES"},
		"expires_at":  {"timestamp with time zone", "NO"},
		"verified_at": {"timestamp with time zone", "YES"},
		"created_at":  {"timestamp with time zone", "NO"},
	}

	found := make(map[string]bool)
	for rows.Next() {
		var colName, dataType, nullable string
		require.NoError(t, rows.Scan(&colName, &dataType, &nullable))
		exp, exists := expected[colName]
		assert.True(t, exists)
		if exists {
			assert.Equal(t, exp.dataType, dataType)
			assert.Equal(t, exp.isNullable, nullable)
		}
		found[colName] = true
	}
	for colName := range expected {
		assert.True(t, found[colName])
	}
}
