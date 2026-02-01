package database

import (
	"testing"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
)

// TestNewPoolNoExternalContext tests that NewPool creates pool without requiring external context.
// Regression test for ISSUE-002: fx dependency injection - NewPool requires context.Context
//
// This test verifies that NewPool can be called with just cfg and logger parameters,
// making it compatible with fx dependency injection which cannot provide context.Context automatically.
func TestNewPoolNoExternalContext(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database pool test in short mode")
	}

	// Create test config with invalid URL (pool creation will fail, but that's expected)
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               "postgres://invalid:5432/test?sslmode=disable",
			MaxConns:          5,
			MinConns:          1,
			MaxConnLifetime:   300000000000,
			MaxConnIdleTime:   60000000000,
			HealthCheckPeriod: 30000000000,
		},
	}

	logger := logging.NewLogger(logging.Config{
		Level:       "error",
		Format:      "text",
		Development: true,
	})

	// The key test: NewPool should accept only cfg and logger (no context parameter)
	// This will fail to connect (invalid URL), but the signature is what we're testing
	pool, err := NewPool(cfg, logger)

	// We expect an error due to invalid URL, but not a compile error or signature mismatch
	assert.Error(t, err, "Should error with invalid database URL")
	assert.Nil(t, pool, "Pool should be nil when creation fails")

	// Verify error message indicates connection failure, not parameter issues
	assert.Contains(t, err.Error(), "failed to", "Error should indicate failure reason")
}

// TestPoolConfigParsing tests that PoolConfig correctly converts config to pgxpool config.
// Ensures the config parsing doesn't depend on external context.
func TestPoolConfigParsing(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			MaxConns:          10,
			MinConns:          2,
			MaxConnLifetime:   1800000000000, // 30m
			MaxConnIdleTime:   300000000000,  // 5m
			HealthCheckPeriod: 30000000000,   // 30s
		},
	}

	poolCfg, err := PoolConfig(cfg)
	assert.NoError(t, err, "PoolConfig should succeed with valid URL")
	assert.NotNil(t, poolCfg, "Should return pool config")

	if poolCfg != nil {
		assert.Equal(t, int32(10), poolCfg.MaxConns, "MaxConns should be set from config")
		assert.Equal(t, int32(2), poolCfg.MinConns, "MinConns should be set from config")
		assert.Equal(t, "dbname", poolCfg.ConnConfig.Database, "Database name should be parsed")
	}
}
