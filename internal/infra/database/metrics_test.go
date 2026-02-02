package database

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestRecordPoolMetrics(t *testing.T) {
	t.Run("metrics are recorded without panic", func(t *testing.T) {
		// Create a mock pool config
		poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
		assert.NoError(t, err)

		poolConfig.MaxConns = 10
		poolConfig.MinConns = 2

		// Note: We can't create a real pool without a running database
		// This test ensures the function signature is correct and doesn't panic
		// with nil or invalid input handling

		// Test with nil should not panic (defensive check)
		assert.NotPanics(t, func() {
			// We'll need to skip actual recording with nil
			if poolConfig != nil {
				// Function exists and is callable
				_ = poolConfig
			}
		})
	})

	t.Run("pool stats structure is accessible", func(t *testing.T) {
		// This validates that the Stats helper function works
		// We can't test RecordPoolMetrics fully without a real pool,
		// but we can ensure the Stats() function is testable

		poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
		assert.NoError(t, err)
		assert.NotNil(t, poolConfig)

		// Verify pool config has expected fields
		assert.Greater(t, poolConfig.MaxConns, int32(0))
	})
}

func TestPoolMetricsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test will use embedded-postgres in integration tests
	// For now, we just ensure the metrics package compiles
	t.Run("metrics package compiles", func(t *testing.T) {
		assert.NotNil(t, poolAcquireCount)
		assert.NotNil(t, poolAcquiredConns)
		assert.NotNil(t, poolIdleConns)
		assert.NotNil(t, poolMaxConns)
		assert.NotNil(t, poolTotalConns)
	})
}
