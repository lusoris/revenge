package database

import (
	"testing"
)

func TestAuthTokensTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "auth_tokens", map[string]columnSpec{
		"id":         {"uuid", "NO"},
		"user_id":    {"uuid", "NO"},
		"token_hash": {"text", "NO"},
		"token_type": {"character varying", "NO"},
		"expires_at": {"timestamp with time zone", "NO"},
		"created_at": {"timestamp with time zone", "NO"},
		"updated_at": {"timestamp with time zone", "NO"},
	})
}

func TestPasswordResetTokensTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "password_reset_tokens", map[string]columnSpec{
		"id":         {"uuid", "NO"},
		"user_id":    {"uuid", "NO"},
		"token_hash": {"text", "NO"},
		"expires_at": {"timestamp with time zone", "NO"},
		"created_at": {"timestamp with time zone", "NO"},
	})
}

func TestEmailVerificationTokensTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping table structure test in short mode")
	}

	db, cleanup := setupTestDB(t, 0)
	defer cleanup()

	assertColumnsExist(t, db, "shared", "email_verification_tokens", map[string]columnSpec{
		"id":         {"uuid", "NO"},
		"user_id":    {"uuid", "NO"},
		"token_hash": {"text", "NO"},
		"email":      {"character varying", "NO"},
		"expires_at": {"timestamp with time zone", "NO"},
		"created_at": {"timestamp with time zone", "NO"},
	})
}
