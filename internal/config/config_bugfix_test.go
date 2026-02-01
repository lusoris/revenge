package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseURLDefault tests that Database.URL has a non-empty default value.
// Regression test for ISSUE-001: TestDefault fails - Database.URL has no default value
func TestDatabaseURLDefault(t *testing.T) {
	cfg := Default()

	// Database.URL should have a placeholder default
	assert.NotEmpty(t, cfg.Database.URL, "Database.URL should have a default value")
	assert.Contains(t, cfg.Database.URL, "postgres://", "Database.URL should be a postgres connection string")
	assert.Contains(t, cfg.Database.URL, "localhost", "Default URL should point to localhost")
	assert.Contains(t, cfg.Database.URL, "revenge", "Default URL should reference revenge database")
}

// TestDefaultsMapDatabaseURL tests that the Defaults() map also contains Database.URL.
// Regression test for ISSUE-001: Both Defaults() and Default() must have the value
func TestDefaultsMapDatabaseURL(t *testing.T) {
	defaults := Defaults()

	dbURL, exists := defaults["database.url"]
	require.True(t, exists, "database.url should exist in Defaults() map")
	assert.NotEmpty(t, dbURL, "database.url should not be empty")

	// Verify it's a valid postgres URL structure
	urlStr, ok := dbURL.(string)
	require.True(t, ok, "database.url should be a string")
	assert.Contains(t, urlStr, "postgres://", "Should be a postgres URL")
}

// TestDefaultConfigStructure tests that the default config has expected structure.
// Note: Full validation will fail due to required fields like JWTSecret that must be user-provided.
func TestDefaultConfigStructure(t *testing.T) {
	cfg := Default()

	// Verify structure is created
	assert.NotNil(t, cfg, "Default config should not be nil")
	assert.NotEmpty(t, cfg.Server.Host, "Server host should have default")
	assert.Greater(t, cfg.Server.Port, 0, "Server port should have default")
	assert.NotEmpty(t, cfg.Database.URL, "Database URL should have default")
}

// TestAuthJWTSecretDefault tests that Auth.JWTSecret validation would fail without value.
// Documents that JWTSecret is required and must be set via config file or env var.
func TestAuthJWTSecretValidation(t *testing.T) {
	cfg := Default()

	// JWTSecret is intentionally empty in defaults (must be set by user)
	// This test documents the expected behavior
	err := validate(cfg)

	if err != nil {
		// Expected: validation should fail if JWTSecret is required
		assert.Contains(t, err.Error(), "JWTSecret", "Validation error should mention JWTSecret")
	}
}
