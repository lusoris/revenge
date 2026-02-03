package database

import (
	"testing"

	"github.com/lusoris/revenge/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestPoolConfig_IntegerOverflowProtection verifies our security fixes
func TestPoolConfig_IntegerOverflowProtection(t *testing.T) {
	tests := []struct {
		name        string
		maxConns    int
		minConns    int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "normal values",
			maxConns:    25,
			minConns:    5,
			expectError: false,
		},
		{
			name:        "max int32 value",
			maxConns:    2147483647, // math.MaxInt32
			minConns:    5,
			expectError: false,
		},
		{
			name:        "zero max (auto mode)",
			maxConns:    0,
			minConns:    2,
			expectError: false,
		},
		{
			name:        "minimum value",
			maxConns:    1,
			minConns:    1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Database: config.DatabaseConfig{
					URL:      "postgres://test:test@localhost:5432/test",
					MaxConns: tt.maxConns,
					MinConns: tt.minConns,
				},
			}

			poolCfg, err := PoolConfig(cfg)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, poolCfg)

				// Verify conversions worked correctly
				if tt.maxConns > 0 {
					assert.Equal(t, int32(tt.maxConns), poolCfg.MaxConns)
				} else {
					// Auto mode - should be (CPU * 2) + 1
					assert.Greater(t, poolCfg.MaxConns, int32(0))
				}

				if tt.minConns > 0 {
					assert.Equal(t, int32(tt.minConns), poolCfg.MinConns)
				}
			}
		})
	}
}

// TestPoolConfig_SafeConversionWorking ensures SafeInt32 catches overflows
func TestPoolConfig_SafeConversionWorking(t *testing.T) {
	// This test verifies that our SafeInt32 function is actually being used
	// On 64-bit systems, we can test with values > int32 max

	// Note: On 32-bit systems, int max = int32 max, so this test is skipped
	const ptrSize = 32 << (^uint(0) >> 63) // 32 or 64
	if ptrSize == 32 {
		t.Skip("Skipping on 32-bit systems where int == int32")
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:      "postgres://test:test@localhost:5432/test",
			MaxConns: 3000000000, // > int32 max (2147483647)
			MinConns: 2,
		},
	}

	poolCfg, err := PoolConfig(cfg)

	// Should get an error because 3000000000 > MaxInt32
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "overflows int32")
	assert.Nil(t, poolCfg)
}
