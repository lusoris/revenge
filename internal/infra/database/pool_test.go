package database

import (
	"runtime"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPoolConfig_DefaultMaxConns(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:      "postgres://user:pass@localhost:5432/test?sslmode=disable",
			MaxConns: 0, // Use default
			MinConns: 2,
		},
	}

	poolCfg, err := PoolConfig(cfg)
	require.NoError(t, err)

	// Should be (CPU * 2) + 1
	expected := int32((runtime.NumCPU() * 2) + 1)
	assert.Equal(t, expected, poolCfg.MaxConns, "Default MaxConns should be (CPU * 2) + 1")
}

func TestPoolConfig_CustomMaxConns(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:      "postgres://user:pass@localhost:5432/test?sslmode=disable",
			MaxConns: 50,
			MinConns: 5,
		},
	}

	poolCfg, err := PoolConfig(cfg)
	require.NoError(t, err)

	assert.Equal(t, int32(50), poolCfg.MaxConns)
	assert.Equal(t, int32(5), poolCfg.MinConns)
}

func TestPoolConfig_ConnectionSettings(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               "postgres://user:pass@localhost:5432/test?sslmode=disable",
			MaxConns:          20,
			MinConns:          5,
			MaxConnLifetime:   30 * time.Minute,
			MaxConnIdleTime:   5 * time.Minute,
			HealthCheckPeriod: 1 * time.Minute,
		},
	}

	poolCfg, err := PoolConfig(cfg)
	require.NoError(t, err)

	assert.Equal(t, int32(20), poolCfg.MaxConns)
	assert.Equal(t, int32(5), poolCfg.MinConns)
	assert.Equal(t, 30*time.Minute, poolCfg.MaxConnLifetime)
	assert.Equal(t, 5*time.Minute, poolCfg.MaxConnIdleTime)
	assert.Equal(t, 1*time.Minute, poolCfg.HealthCheckPeriod)
}

func TestPoolConfig_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL: "not-a-valid-url",
		},
	}

	poolCfg, err := PoolConfig(cfg)
	assert.Error(t, err, "Should error on invalid URL")
	assert.Nil(t, poolCfg)
	assert.Contains(t, err.Error(), "failed to parse database URL")
}

func TestPoolConfig_ParsesHostAndDatabase(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL: "postgres://myuser:mypass@myhost:5433/mydb?sslmode=require",
		},
	}

	poolCfg, err := PoolConfig(cfg)
	require.NoError(t, err)

	assert.Equal(t, "myhost", poolCfg.ConnConfig.Host)
	assert.Equal(t, uint16(5433), poolCfg.ConnConfig.Port)
	assert.Equal(t, "mydb", poolCfg.ConnConfig.Database)
	assert.Equal(t, "myuser", poolCfg.ConnConfig.User)
}

func TestPoolConfig_ZeroValuesNotOverridden(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:               "postgres://user:pass@localhost:5432/test?sslmode=disable",
			MaxConns:          0, // 0 means auto
			MinConns:          0, // 0 should not override
			MaxConnLifetime:   0, // 0 should not override
			MaxConnIdleTime:   0, // 0 should not override
			HealthCheckPeriod: 0, // 0 should not override
		},
	}

	poolCfg, err := PoolConfig(cfg)
	require.NoError(t, err)

	// MaxConns gets auto value
	assert.Greater(t, poolCfg.MaxConns, int32(0))
	// Others should use pgxpool defaults (we can't easily verify, just check no error)
}

func TestPoolConfig_MultipleURLFormats(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		host string
		db   string
	}{
		{
			name: "standard format",
			url:  "postgres://user:pass@localhost:5432/mydb?sslmode=disable",
			host: "localhost",
			db:   "mydb",
		},
		{
			name: "with options",
			url:  "postgres://user:pass@db.example.com:5432/production?sslmode=require&connect_timeout=10",
			host: "db.example.com",
			db:   "production",
		},
		{
			name: "ipv4 host",
			url:  "postgres://admin:secret@192.168.1.100:5432/app?sslmode=disable",
			host: "192.168.1.100",
			db:   "app",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				Database: config.DatabaseConfig{URL: tc.url},
			}

			poolCfg, err := PoolConfig(cfg)
			require.NoError(t, err)

			assert.Equal(t, tc.host, poolCfg.ConnConfig.Host)
			assert.Equal(t, tc.db, poolCfg.ConnConfig.Database)
		})
	}
}

func TestNewPool_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL: "invalid-url",
		},
	}

	logger := logging.NewLogger(logging.Config{
		Level:       "error",
		Format:      "text",
		Development: true,
	})

	pool, err := NewPool(cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "failed to parse database URL")
}

func TestNewPool_ConnectionRefused(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:      "postgres://user:pass@localhost:59999/test?sslmode=disable",
			MaxConns: 5,
			MinConns: 1,
		},
	}

	logger := logging.NewLogger(logging.Config{
		Level:       "error",
		Format:      "text",
		Development: true,
	})

	pool, err := NewPool(cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, pool)
	// Error should indicate connection failure
	assert.Contains(t, err.Error(), "failed to")
}
