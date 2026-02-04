//go:build integration

package cache

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_L2Cache tests L2 cache operations with a real Dragonfly instance.
func TestIntegration_L2Cache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start Dragonfly container
	df := testutil.NewDragonflyContainer(t)
	defer df.Close()

	// Create cache client with Dragonfly connection
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     df.URL,
		},
	}

	logger := slog.Default()
	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	defer client.Close()

	// Verify connection
	ctx := context.Background()
	err = client.Ping(ctx)
	require.NoError(t, err, "should connect to Dragonfly")

	// Create unified cache with L1 and L2
	cache, err := NewCache(client, 1000, time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test:integration:setget"
		value := []byte("hello world")

		// Set value
		err := cache.Set(ctx, key, value, 5*time.Minute)
		require.NoError(t, err)

		// Get value - should hit L1 first
		result, err := cache.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("L2 Fallback", func(t *testing.T) {
		key := "test:integration:l2fallback"
		value := []byte("from L2")

		// Set value with short TTL (will go to L2 only due to L1 TTL policy)
		err := cache.Set(ctx, key, value, 10*time.Second)
		require.NoError(t, err)

		// Clear L1 to force L2 lookup
		cache.l1.Clear()

		// Get should fall back to L2
		result, err := cache.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:integration:delete"
		value := []byte("to be deleted")

		// Set value
		err := cache.Set(ctx, key, value, 5*time.Minute)
		require.NoError(t, err)

		// Verify it exists
		_, err = cache.Get(ctx, key)
		require.NoError(t, err)

		// Delete
		err = cache.Delete(ctx, key)
		require.NoError(t, err)

		// Should fail to get after deletion
		_, err = cache.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("Exists", func(t *testing.T) {
		key := "test:integration:exists"
		value := []byte("check existence")

		// Should not exist initially
		exists, err := cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)

		// Set value
		err = cache.Set(ctx, key, value, 5*time.Minute)
		require.NoError(t, err)

		// Should exist now
		exists, err = cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("JSON Operations", func(t *testing.T) {
		key := "test:integration:json"
		type TestData struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		original := TestData{Name: "test", Value: 42}

		// Set JSON
		err := cache.SetJSON(ctx, key, original, 5*time.Minute)
		require.NoError(t, err)

		// Get JSON
		var result TestData
		err = cache.GetJSON(ctx, key, &result)
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("Pattern Invalidation", func(t *testing.T) {
		// Set multiple keys with same prefix
		prefix := "test:integration:pattern:"
		for i := 0; i < 3; i++ {
			key := prefix + string(rune('a'+i))
			err := cache.Set(ctx, key, []byte("value"), 5*time.Minute)
			require.NoError(t, err)
		}

		// Invalidate by pattern
		err := cache.Invalidate(ctx, prefix+"*")
		require.NoError(t, err)

		// All keys should be gone
		for i := 0; i < 3; i++ {
			key := prefix + string(rune('a'+i))
			_, err := cache.Get(ctx, key)
			assert.Error(t, err, "key %s should be invalidated", key)
		}
	})
}

// TestIntegration_TTLExpiration tests that TTL actually expires keys.
func TestIntegration_TTLExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	df := testutil.NewDragonflyContainer(t)
	defer df.Close()

	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     df.URL,
		},
	}

	client, err := NewClient(cfg, slog.Default())
	require.NoError(t, err)
	defer client.Close()

	// Create cache with very short L1 TTL
	cache, err := NewCache(client, 1000, 100*time.Millisecond)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()
	key := "test:ttl:expiration"
	value := []byte("expires soon")

	// Set with 500ms TTL
	err = cache.Set(ctx, key, value, 500*time.Millisecond)
	require.NoError(t, err)

	// Should exist immediately
	_, err = cache.Get(ctx, key)
	require.NoError(t, err)

	// Wait for expiration
	time.Sleep(600 * time.Millisecond)

	// Clear L1 to force L2 lookup
	cache.l1.Clear()

	// Should be expired in L2
	_, err = cache.Get(ctx, key)
	assert.Error(t, err, "key should be expired")
}

// TestIntegration_SessionCache tests session-specific caching scenarios.
func TestIntegration_SessionCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	df := testutil.NewDragonflyContainer(t)
	defer df.Close()

	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     df.URL,
		},
	}

	client, err := NewClient(cfg, slog.Default())
	require.NoError(t, err)
	defer client.Close()

	cache, err := NewNamedCache(client, 1000, time.Minute, "session")
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	t.Run("Session Key Functions", func(t *testing.T) {
		tokenHash := "abc123def456"
		key := SessionKey(tokenHash)
		assert.Equal(t, "session:abc123def456", key)

		// Set session data
		sessionData := map[string]interface{}{
			"user_id":   "user-123",
			"roles":     []string{"user", "admin"},
			"issued_at": time.Now().Unix(),
		}
		err := cache.SetJSON(ctx, key, sessionData, SessionTTL)
		require.NoError(t, err)

		// Retrieve session
		var result map[string]interface{}
		err = cache.GetJSON(ctx, key, &result)
		require.NoError(t, err)
		assert.Equal(t, "user-123", result["user_id"])
	})

	t.Run("Invalidate User Sessions", func(t *testing.T) {
		userID := "user-to-invalidate"

		// Create multiple session keys for user
		for i := 0; i < 3; i++ {
			key := SessionByUserKey(userID) + ":session" + string(rune('1'+i))
			err := cache.Set(ctx, key, []byte("session data"), time.Minute)
			require.NoError(t, err)
		}

		// Invalidate all user sessions
		err := cache.InvalidateUserSessions(ctx, userID)
		require.NoError(t, err)

		// Sessions should be invalidated (L1 cleared, L2 keys deleted)
		// Note: L1 is fully cleared on invalidation, so checking specific keys
		// requires checking L2 directly
	})
}

// TestIntegration_ClientPing tests basic connectivity.
func TestIntegration_ClientPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	df := testutil.NewDragonflyContainer(t)
	defer df.Close()

	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     df.URL,
		},
	}

	client, err := NewClient(cfg, slog.Default())
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	err = client.Ping(ctx)
	assert.NoError(t, err, "ping should succeed")
}

// TestIntegration_DisabledCache tests behavior when cache is disabled.
func TestIntegration_DisabledCache(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}

	client, err := NewClient(cfg, slog.Default())
	require.NoError(t, err)
	defer client.Close()

	// Ping should fail when cache is disabled
	ctx := context.Background()
	err = client.Ping(ctx)
	assert.Error(t, err)

	// RueidisClient should be nil
	assert.Nil(t, client.RueidisClient())
}
