package metadata

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	assert.Equal(t, rate.Limit(4.0), config.RateLimit)
	assert.Equal(t, 10, config.RateBurst)
	assert.Equal(t, 24*time.Hour, config.CacheTTL)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 3, config.RetryCount)
}

func TestNewBaseClient(t *testing.T) {
	t.Run("with defaults", func(t *testing.T) {
		config := ClientConfig{
			BaseURL: "https://api.example.com",
			APIKey:  "test-key",
		}

		client := NewBaseClient(config)

		assert.NotNil(t, client)
		assert.Equal(t, "test-key", client.GetAPIKey())
		assert.Equal(t, "https://api.example.com", client.GetBaseURL())
	})

	t.Run("with custom values", func(t *testing.T) {
		config := ClientConfig{
			BaseURL:    "https://api.custom.com",
			APIKey:     "custom-key",
			RateLimit:  rate.Limit(10.0),
			RateBurst:  20,
			CacheTTL:   1 * time.Hour,
			Timeout:    15 * time.Second,
			RetryCount: 5,
		}

		client := NewBaseClient(config)

		assert.NotNil(t, client)
		assert.Equal(t, "custom-key", client.GetAPIKey())
		assert.Equal(t, "https://api.custom.com", client.GetBaseURL())
	})
}

func TestBaseClientCache(t *testing.T) {
	config := ClientConfig{
		BaseURL:  "https://api.example.com",
		APIKey:   "test-key",
		CacheTTL: 100 * time.Millisecond,
	}
	client := NewBaseClient(config)

	t.Run("set and get", func(t *testing.T) {
		client.SetCache("test-key", "test-value")

		result := client.GetFromCache("test-key")
		assert.Equal(t, "test-value", result)
	})

	t.Run("get non-existent", func(t *testing.T) {
		result := client.GetFromCache("non-existent")
		assert.Nil(t, result)
	})

	t.Run("expiration", func(t *testing.T) {
		client.SetCache("expires", "soon")

		// Should exist immediately
		result := client.GetFromCache("expires")
		assert.Equal(t, "soon", result)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Should be gone
		result = client.GetFromCache("expires")
		assert.Nil(t, result)
	})

	t.Run("clear cache", func(t *testing.T) {
		client.SetCache("key1", "value1")
		client.SetCache("key2", "value2")

		client.ClearCache()

		assert.Nil(t, client.GetFromCache("key1"))
		assert.Nil(t, client.GetFromCache("key2"))
	})

	t.Run("set with custom TTL", func(t *testing.T) {
		client.SetCacheWithTTL("custom-ttl", "value", 50*time.Millisecond)

		// Should exist
		result := client.GetFromCache("custom-ttl")
		assert.Equal(t, "value", result)

		// Wait for expiration (custom TTL is shorter)
		time.Sleep(100 * time.Millisecond)

		// Should be gone
		result = client.GetFromCache("custom-ttl")
		assert.Nil(t, result)
	})
}

func TestBaseClientRateLimit(t *testing.T) {
	config := ClientConfig{
		BaseURL:   "https://api.example.com",
		APIKey:    "test-key",
		RateLimit: rate.Limit(100.0), // High limit for fast tests
		RateBurst: 10,
	}
	client := NewBaseClient(config)

	ctx := context.Background()

	// Should not block with high rate limit
	for i := 0; i < 5; i++ {
		err := client.WaitForRateLimit(ctx)
		require.NoError(t, err)
	}
}

func TestBaseClientRateLimitContextCancel(t *testing.T) {
	config := ClientConfig{
		BaseURL:   "https://api.example.com",
		APIKey:    "test-key",
		RateLimit: rate.Limit(0.1), // Very low rate limit
		RateBurst: 1,
	}
	client := NewBaseClient(config)

	// First request should succeed
	ctx := context.Background()
	err := client.WaitForRateLimit(ctx)
	require.NoError(t, err)

	// Second request with cancelled context should fail
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = client.WaitForRateLimit(ctx)
	assert.Error(t, err)
}

func TestCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		parts    []any
		expected string
	}{
		{
			name:     "single string",
			parts:    []any{"movie"},
			expected: "movie",
		},
		{
			name:     "string and int",
			parts:    []any{"movie", 123},
			expected: "movie:123",
		},
		{
			name:     "multiple parts",
			parts:    []any{"search", "query", "en-US", 1},
			expected: "search:query:en-US:1",
		},
		{
			name:     "with nil",
			parts:    []any{"movie", nil},
			expected: "movie:<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CacheKey(tt.parts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheEntryIsExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		entry := &CacheEntry{
			Data:      "test",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}
		assert.False(t, entry.IsExpired())
	})

	t.Run("expired", func(t *testing.T) {
		entry := &CacheEntry{
			Data:      "test",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		assert.True(t, entry.IsExpired())
	})
}
