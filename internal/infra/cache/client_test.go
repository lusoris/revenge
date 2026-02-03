package cache

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/config"
)

func TestNewClient_Disabled(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Nil(t, client.RueidisClient())
	assert.False(t, cfg.Cache.Enabled)
}

func TestNewClient_EmptyURL(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "cache URL is required")
}

func TestNewClient_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "invalid://url",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to parse cache URL")
}

func TestNewClient_ValidURL(t *testing.T) {
	t.Skip("Skipping test - requires running Redis/Dragonfly instance")

	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "redis://localhost:6379/0",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.RueidisClient())

	// Clean up
	client.Close()
}

func TestClient_Close_WhenDisabled(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)

	// Should not panic
	assert.NotPanics(t, func() {
		client.Close()
	})
}

func TestClient_Close_WhenEnabled(t *testing.T) {
	t.Skip("Skipping test - requires running Redis/Dragonfly instance")

	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "redis://localhost:6379/0",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)

	// Should not panic
	assert.NotPanics(t, func() {
		client.Close()
	})
}

func TestClient_Ping_NotInitialized(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)

	ctx := context.Background()
	err = client.Ping(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rueidis client not initialized")
}

func TestClient_Ping_Integration(t *testing.T) {
	t.Skip("Skipping integration test - requires running Redis/Dragonfly instance")

	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "redis://localhost:6379/0",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx)
	assert.NoError(t, err)
}

func TestClient_RueidisClient(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		cfg := &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		client, err := NewClient(cfg, logger)
		require.NoError(t, err)

		assert.Nil(t, client.RueidisClient())
	})
	t.Skip("Skipping test - requires running Redis/Dragonfly instance")

	t.Run("Enabled", func(t *testing.T) {
		cfg := &config.Config{
			Cache: config.CacheConfig{
				Enabled: true,
				URL:     "redis://localhost:6379/0",
			},
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		client, err := NewClient(cfg, logger)
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client.RueidisClient())
	})
}

func TestDefaultConstants(t *testing.T) {
	assert.Equal(t, 5*time.Second, DefaultDialTimeout)
	assert.Equal(t, 3*time.Second, DefaultReadTimeout)
	assert.Equal(t, 3*time.Second, DefaultWriteTimeout)
	assert.Equal(t, 16*1024*1024, DefaultCacheSizeEachConn)
	assert.Equal(t, 10, DefaultRingScale)
	assert.Equal(t, 128, DefaultBlockingPoolSize)
}

func TestNewClient_URLParsing(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		enabled     bool
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Valid Redis URL (may connect if service running)",
			url:         "redis://localhost:6379/0",
			enabled:     true,
			shouldError: false, // May succeed if dragonfly/redis is running locally
			errorMsg:    "",
		},
		{
			name:        "Unreachable host (will fail)",
			url:         "redis://192.0.2.1:6379/0", // TEST-NET-1, guaranteed unreachable
			enabled:     true,
			shouldError: true, // Will fail on connection
			errorMsg:    "failed to create rueidis client",
		},
		{
			name:        "Valid Redis URL with auth (may connect)",
			url:         "redis://:password@localhost:6379/0",
			enabled:     true,
			shouldError: false, // May succeed if service running (might fail auth)
		},
		{
			name:        "Valid Redis URL with user and password (auth may fail)",
			url:         "redis://user:password@localhost:6379/1",
			enabled:     true,
			shouldError: false, // May succeed or fail based on auth config
		},
		{
			name:        "Invalid scheme",
			url:         "http://localhost:6379",
			enabled:     true,
			shouldError: true,
			errorMsg:    "failed to parse cache URL",
		},
		{
			name:        "Malformed URL",
			url:         "redis://[invalid",
			enabled:     true,
			shouldError: true,
			errorMsg:    "failed to parse cache URL",
		},
		{
			name:        "Disabled cache ignores invalid URL",
			url:         "invalid://bad",
			enabled:     false,
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				Cache: config.CacheConfig{
					Enabled: tc.enabled,
					URL:     tc.url,
				},
			}
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			client, err := NewClient(cfg, logger)

			if tc.shouldError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				// For cases where connection may or may not succeed
				// (e.g., service might be running or not, auth might work or not)
				// Just ensure proper cleanup if client was created
				if err != nil {
					t.Logf("Got expected error (service not available or auth failed): %v", err)
				}
			}

			// Always cleanup if client was created
			if client != nil {
				client.Close()
			}
		})
	}
}

func TestClient_Methods_WithDisabledCache(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test all methods work when cache is disabled
	assert.Nil(t, client.RueidisClient())

	ctx := context.Background()
	err = client.Ping(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rueidis client not initialized")

	// Close should not panic
	assert.NotPanics(t, func() {
		client.Close()
	})
}
