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

func TestNewCache(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, cache)
	defer cache.Close()

	assert.NotNil(t, cache.l1)
	assert.NotNil(t, cache.client)
}

func TestCache_L1Only_Get(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Cache miss
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")

	// Set in L1 directly
	cache.l1.Set("key1", []byte("value1"))

	// Get should hit L1
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)
}

func TestCache_Set(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set value
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	// Should be in L1
	val, ok := cache.l1.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, []byte("value1"), val)

	// Get should work
	val, err = cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)
}

func TestCache_Delete(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set value
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	// Verify it exists
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Delete
	err = cache.Delete(ctx, "key1")
	require.NoError(t, err)

	// Should be gone
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
}

func TestCache_Exists(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Initially should not exist
	exists, err := cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)

	// Set value
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	// Now should exist
	exists, err = cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.True(t, exists)

	// Delete
	err = cache.Delete(ctx, "key1")
	require.NoError(t, err)

	// Should not exist anymore
	exists, err = cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestCache_Invalidate(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set multiple values
	err = cache.Set(ctx, "user:1", []byte("alice"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, "user:2", []byte("bob"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, "session:abc", []byte("data"), 1*time.Minute)
	require.NoError(t, err)

	// Verify L1 has entries
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate pattern (clears entire L1)
	err = cache.Invalidate(ctx, "user:*")
	require.NoError(t, err)

	// L1 should be empty
	assert.Equal(t, 0, cache.l1.Size())
}

func TestCache_GetJSON(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	type User struct {
		ID   int
		Name string
	}

	user := User{ID: 1, Name: "Alice"}

	// Set JSON
	err = cache.SetJSON(ctx, "user:1", user, 1*time.Minute)
	require.NoError(t, err)

	// Get JSON
	var retrieved User
	err = cache.GetJSON(ctx, "user:1", &retrieved)
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}

func TestCache_SetJSON(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	type Session struct {
		Token     string
		UserID    int
		ExpiresAt time.Time
	}

	session := Session{
		Token:     "abc123",
		UserID:    42,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Set JSON
	err = cache.SetJSON(ctx, "session:abc123", session, 1*time.Minute)
	require.NoError(t, err)

	// Verify it's in L1
	data, ok := cache.l1.Get("session:abc123")
	assert.True(t, ok)
	assert.NotEmpty(t, data)

	// Get and verify
	var retrieved Session
	err = cache.GetJSON(ctx, "session:abc123", &retrieved)
	require.NoError(t, err)
	assert.Equal(t, session.Token, retrieved.Token)
	assert.Equal(t, session.UserID, retrieved.UserID)
}

func TestCache_Close(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)

	// Should not panic
	assert.NotPanics(t, func() {
		cache.Close()
	})
}
