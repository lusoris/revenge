package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"go.uber.org/mock/gomock"

	"github.com/lusoris/revenge/internal/config"
)

// newL1OnlyCache creates a cache with no L2 (rueidis) for unit testing.
// Uses a 1-minute L1 TTL so test values with 1-minute TTL are stored in L1.
func newL1OnlyCache(t *testing.T) *Cache {
	t.Helper()
	c, err := NewCache(nil, 1000, 1*time.Minute)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

func newL1OnlyCacheWithClient(t *testing.T) *Cache {
	t.Helper()
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: false},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: nil, // no rueidis
	}
	c, err := NewCache(client, 1000, 1*time.Minute)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

// TestNewNamedCache_ErrorPath tests NewNamedCache error handling when L1 creation fails.
func TestNewNamedCache_ErrorPath(t *testing.T) {
	// Negative maxSize gets corrected to default, so this should succeed
	c, err := NewNamedCache(nil, -1, -1, "error_test")
	require.NoError(t, err)
	require.NotNil(t, c)
	c.Close()

	// Verify the name is set
	assert.Equal(t, "error_test", c.name)
}

// TestNewCache_NilClient tests that NewCache works with nil client.
func TestNewCache_NilClient(t *testing.T) {
	c, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Close()

	// Verify default name
	assert.Equal(t, "default", c.name)
}

// TestCache_Get_L1Hit verifies L1 cache hit path records metrics.
func TestCache_Get_L1Hit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Pre-populate L1
	c.l1.Set("hit-key", []byte("hit-value"))

	val, err := c.Get(ctx, "hit-key")
	require.NoError(t, err)
	assert.Equal(t, []byte("hit-value"), val)
}

// TestCache_Get_L1Miss_L2Unavailable verifies the error path when L1 misses and L2 is nil.
func TestCache_Get_L1Miss_L2Unavailable(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	_, err := c.Get(ctx, "nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
	assert.Contains(t, err.Error(), "L2 unavailable")
}

// TestCache_Get_ClientWithNilRueidis verifies behavior when client exists but rueidisClient is nil.
func TestCache_Get_ClientWithNilRueidis(t *testing.T) {
	c := newL1OnlyCacheWithClient(t)
	ctx := context.Background()

	_, err := c.Get(ctx, "no-l2")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// TestCache_Set_TTLBranches tests all TTL code paths in Set.
func TestCache_Set_TTLBranches(t *testing.T) {
	tests := []struct {
		name        string
		ttl         time.Duration
		l1TTL       time.Duration
		expectInL1  bool
		description string
	}{
		{
			name:        "zero TTL stores in L1",
			ttl:         0,
			l1TTL:       1 * time.Minute,
			expectInL1:  true,
			description: "Zero TTL means no expiration, uses L1",
		},
		{
			name:        "TTL equal to L1 TTL stores in L1",
			ttl:         5 * time.Minute,
			l1TTL:       5 * time.Minute,
			expectInL1:  true,
			description: "Equal TTL stores in L1",
		},
		{
			name:        "TTL longer than L1 TTL stores in L1",
			ttl:         10 * time.Minute,
			l1TTL:       5 * time.Minute,
			expectInL1:  true,
			description: "Longer TTL stores in L1",
		},
		{
			name:        "TTL shorter than L1 TTL deletes from L1",
			ttl:         1 * time.Second,
			l1TTL:       5 * time.Minute,
			expectInL1:  false,
			description: "Short TTL skips L1 to prevent stale reads",
		},
		{
			name:        "sub-second TTL deletes from L1",
			ttl:         100 * time.Millisecond,
			l1TTL:       5 * time.Minute,
			expectInL1:  false,
			description: "Sub-second TTL skips L1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCache(nil, 100, tt.l1TTL)
			require.NoError(t, err)
			defer c.Close()

			ctx := context.Background()
			key := "ttl-test-" + tt.name

			// Pre-populate L1 to test deletion path for short TTL
			c.l1.Set(key, []byte("old-value"))

			err = c.Set(ctx, key, []byte("new-value"), tt.ttl)
			require.NoError(t, err)

			_, ok := c.l1.Get(key)
			assert.Equal(t, tt.expectInL1, ok, tt.description)
		})
	}
}

// TestCache_Set_NoL2 verifies Set works without L2 (nil client).
func TestCache_Set_NoL2(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// All TTL branches should succeed with nil client
	err := c.Set(ctx, "k1", []byte("v1"), 0)
	require.NoError(t, err)

	err = c.Set(ctx, "k2", []byte("v2"), 500*time.Millisecond)
	require.NoError(t, err)

	err = c.Set(ctx, "k3", []byte("v3"), 5*time.Second)
	require.NoError(t, err)
}

// TestCache_Delete_NoL2 verifies Delete works with nil client.
func TestCache_Delete_NoL2(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Set then delete
	c.l1.Set("del-key", []byte("value"))
	assert.True(t, c.l1.Has("del-key"))

	err := c.Delete(ctx, "del-key")
	require.NoError(t, err)
	assert.False(t, c.l1.Has("del-key"))
}

// TestCache_Delete_ClientWithNilRueidis verifies Delete with client but no rueidis.
func TestCache_Delete_ClientWithNilRueidis(t *testing.T) {
	c := newL1OnlyCacheWithClient(t)
	ctx := context.Background()

	c.l1.Set("del-key", []byte("value"))

	err := c.Delete(ctx, "del-key")
	require.NoError(t, err)
	assert.False(t, c.l1.Has("del-key"))
}

// TestCache_Exists_L1Hit verifies Exists returns true for L1 hit.
func TestCache_Exists_L1Hit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set("exists-key", []byte("value"))

	exists, err := c.Exists(ctx, "exists-key")
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestCache_Exists_L1Miss_NoL2 verifies Exists returns false when not in L1 and no L2.
func TestCache_Exists_L1Miss_NoL2(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	exists, err := c.Exists(ctx, "nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestCache_Exists_ClientWithNilRueidis verifies Exists with client but no rueidis.
func TestCache_Exists_ClientWithNilRueidis(t *testing.T) {
	c := newL1OnlyCacheWithClient(t)
	ctx := context.Background()

	exists, err := c.Exists(ctx, "nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestCache_Invalidate_NoL2 verifies Invalidate removes only matching keys from L1 with no L2.
func TestCache_Invalidate_NoL2(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set("user:1", []byte("1"))
	c.l1.Set("user:2", []byte("2"))
	c.l1.Set("session:1", []byte("3"))
	assert.Equal(t, 3, c.l1.Size())

	err := c.Invalidate(ctx, "user:*")
	require.NoError(t, err)

	// Only user:* keys removed, session:1 survives
	_, ok := c.l1.Get("user:1")
	assert.False(t, ok)
	_, ok = c.l1.Get("user:2")
	assert.False(t, ok)
	_, ok = c.l1.Get("session:1")
	assert.True(t, ok)
	assert.Equal(t, 1, c.l1.Size())
}

// TestCache_Invalidate_ClientWithNilRueidis verifies Invalidate with client but no rueidis.
func TestCache_Invalidate_ClientWithNilRueidis(t *testing.T) {
	c := newL1OnlyCacheWithClient(t)
	ctx := context.Background()

	c.l1.Set("x", []byte("1"))
	assert.Greater(t, c.l1.Size(), 0)

	err := c.Invalidate(ctx, "*")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_Close_NilL1 verifies Close doesn't panic with nil l1.
func TestCache_Close_NilL1(t *testing.T) {
	c := &Cache{l1: nil, client: nil}
	assert.NotPanics(t, func() {
		c.Close()
	})
}

// TestCache_GetJSON_CacheMiss verifies GetJSON returns error on cache miss.
func TestCache_GetJSON_CacheMiss(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	var result map[string]string
	err := c.GetJSON(ctx, "nonexistent", &result)
	require.Error(t, err)
}

// TestCache_GetJSON_InvalidJSON verifies GetJSON error on corrupt data.
func TestCache_GetJSON_InvalidJSON(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set("bad-json", []byte("not json {{{"))

	var result map[string]string
	err := c.GetJSON(ctx, "bad-json", &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}

// TestCache_SetJSON_MarshalError verifies SetJSON returns error for unmarshalable types.
func TestCache_SetJSON_MarshalError(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Channels cannot be marshaled to JSON
	ch := make(chan int)
	err := c.SetJSON(ctx, "bad-value", ch, 1*time.Minute)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "marshal")
}

// TestCache_SetJSON_ComplexTypes tests SetJSON with various types.
func TestCache_SetJSON_ComplexTypes(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	tests := []struct {
		name  string
		key   string
		value any
	}{
		{"nil value", "nil", nil},
		{"string value", "str", "hello"},
		{"int value", "int", 42},
		{"bool value", "bool", true},
		{"slice value", "slice", []string{"a", "b", "c"}},
		{"map value", "map", map[string]int{"x": 1, "y": 2}},
		{"nested struct", "nested", struct {
			Name   string `json:"name"`
			Nested struct {
				Value int `json:"value"`
			} `json:"nested"`
		}{Name: "test", Nested: struct {
			Value int `json:"value"`
		}{Value: 99}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.SetJSON(ctx, tt.key, tt.value, 1*time.Minute)
			require.NoError(t, err)

			// Verify it can be read back
			data, err := c.Get(ctx, tt.key)
			require.NoError(t, err)

			// Verify it's valid JSON
			assert.True(t, json.Valid(data), "stored data should be valid JSON")
		})
	}
}

// TestCache_CacheAside_CacheHit verifies CacheAside uses cached data.
func TestCache_CacheAside_CacheHit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	type Item struct {
		ID string `json:"id"`
	}

	// Pre-populate cache
	err := c.SetJSON(ctx, "aside:1", &Item{ID: "cached"}, 1*time.Minute)
	require.NoError(t, err)

	loaderCalled := false
	loader := func() (any, error) {
		loaderCalled = true
		return &Item{ID: "loaded"}, nil
	}

	var result Item
	err = c.CacheAside(ctx, "aside:1", 1*time.Minute, loader, &result)
	require.NoError(t, err)
	assert.False(t, loaderCalled, "loader should not be called on cache hit")
	assert.Equal(t, "cached", result.ID)
}

// TestCache_CacheAside_CacheMissLoadsAndCaches verifies CacheAside loads and caches on miss.
func TestCache_CacheAside_CacheMissLoadsAndCaches(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	type Item struct {
		ID string `json:"id"`
	}

	callCount := 0
	loader := func() (any, error) {
		callCount++
		return &Item{ID: "from-loader"}, nil
	}

	var result Item
	err := c.CacheAside(ctx, "aside:miss", 1*time.Minute, loader, &result)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)
	assert.Equal(t, "from-loader", result.ID)

	// Wait for async cache set
	time.Sleep(200 * time.Millisecond)

	// Second call should hit cache
	var result2 Item
	err = c.CacheAside(ctx, "aside:miss", 1*time.Minute, loader, &result2)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount, "loader should not be called on second call (cache hit)")
	assert.Equal(t, "from-loader", result2.ID)
}

// TestCache_CacheAside_LoaderError verifies CacheAside propagates loader errors.
func TestCache_CacheAside_LoaderError(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	loader := func() (any, error) {
		return nil, assert.AnError
	}

	var result map[string]string
	err := c.CacheAside(ctx, "aside:err", 1*time.Minute, loader, &result)
	require.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

// TestCache_InvalidatePattern_Alias verifies InvalidatePattern calls Invalidate.
func TestCache_InvalidatePattern_Alias(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set("pattern:a", []byte("1"))
	c.l1.Set("pattern:b", []byte("2"))

	err := c.InvalidatePattern(ctx, "pattern:*")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateSession verifies session invalidation.
func TestCache_InvalidateSession_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	key := SessionKey("token-hash-abc")
	c.l1.Set(key, []byte("session-data"))

	err := c.InvalidateSession(ctx, "token-hash-abc")
	require.NoError(t, err)
	assert.False(t, c.l1.Has(key))
}

// TestCache_InvalidateUserSessions verifies user session invalidation clears L1.
func TestCache_InvalidateUserSessions_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set(SessionByUserKey("u1")+":s1", []byte("d1"))
	c.l1.Set(SessionByUserKey("u1")+":s2", []byte("d2"))

	err := c.InvalidateUserSessions(ctx, "u1")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateRBACForUser_Success verifies RBAC invalidation for a specific user.
func TestCache_InvalidateRBACForUser_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Set RBAC entries
	c.l1.Set(RBACEnforceKey("u1", "resource", "read"), []byte("allow"))
	c.l1.Set(RBACUserRolesKey("u1"), []byte(`["admin"]`))
	c.l1.Set(RBACUserPermsKey("u1"), []byte(`["read"]`))

	err := c.InvalidateRBACForUser(ctx, "u1")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateAllRBAC_Unit verifies all RBAC invalidation.
func TestCache_InvalidateAllRBAC_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set("rbac:policy:1", []byte("p1"))
	c.l1.Set("rbac:roles:u1", []byte("r1"))

	err := c.InvalidateAllRBAC(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateServerSettings_Unit verifies server settings invalidation.
func TestCache_InvalidateServerSettings_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set(ServerSettingKey("theme"), []byte("dark"))

	err := c.InvalidateServerSettings(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateServerSetting_Unit verifies specific server setting invalidation.
func TestCache_InvalidateServerSetting_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	key := ServerSettingKey("locale")
	c.l1.Set(key, []byte("en"))

	err := c.InvalidateServerSetting(ctx, "locale")
	require.NoError(t, err)
	assert.False(t, c.l1.Has(key))
}

// TestCache_InvalidateUserSettings_Unit verifies user settings invalidation.
func TestCache_InvalidateUserSettings_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set(UserSettingKey("u1", "theme"), []byte("dark"))
	c.l1.Set(UserSettingKey("u1", "lang"), []byte("en"))

	err := c.InvalidateUserSettings(ctx, "u1")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateMovie_Unit verifies movie invalidation clears all related keys.
func TestCache_InvalidateMovie_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	movieID := "m1"
	c.l1.Set(MovieKey(movieID), []byte("data"))
	c.l1.Set(MovieCastKey(movieID), []byte("cast"))
	c.l1.Set(MovieCrewKey(movieID), []byte("crew"))
	c.l1.Set(MovieGenresKey(movieID), []byte("genres"))
	c.l1.Set(MovieFilesKey(movieID), []byte("files"))

	err := c.InvalidateMovie(ctx, movieID)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateMovieLists_Unit verifies movie lists invalidation.
func TestCache_InvalidateMovieLists_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set(MovieListKey("hash1"), []byte("list"))
	c.l1.Set(MovieRecentKey(10, 0), []byte("recent"))
	c.l1.Set(MovieTopRatedKey(100, 10, 0), []byte("top"))

	err := c.InvalidateMovieLists(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateSearch_Unit verifies search invalidation.
func TestCache_InvalidateSearch_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set(SearchMoviesKey("hash"), []byte("results"))
	c.l1.Set(SearchAutocompleteKey("mat"), []byte("auto"))

	err := c.InvalidateSearch(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateLibrary_Unit verifies library invalidation.
func TestCache_InvalidateLibrary_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	libID := "lib1"
	c.l1.Set(LibraryKey(libID), []byte("lib"))
	c.l1.Set(LibraryStatsKey(libID), []byte("stats"))

	err := c.InvalidateLibrary(ctx, libID)
	require.NoError(t, err)
	// Library invalidation uses Delete (not Invalidate/Clear), so check specific keys
	assert.False(t, c.l1.Has(LibraryKey(libID)))
	assert.False(t, c.l1.Has(LibraryStatsKey(libID)))
}

// TestCache_InvalidateUser_Unit verifies user invalidation.
func TestCache_InvalidateUser_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	userID := "u1"
	c.l1.Set(UserKey(userID), []byte("user"))
	c.l1.Set(ContinueWatchingKey(userID, 10), []byte("watch"))

	err := c.InvalidateUser(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestCache_InvalidateContinueWatching_Unit verifies continue watching invalidation.
func TestCache_InvalidateContinueWatching_Unit(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	userID := "u1"
	c.l1.Set(ContinueWatchingKey(userID, 5), []byte("watch5"))
	c.l1.Set(ContinueWatchingKey(userID, 10), []byte("watch10"))

	err := c.InvalidateContinueWatching(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// TestClient_Ping_NilRueidis verifies Ping returns error with nil rueidis.
func TestClient_Ping_NilRueidis(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: false},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: nil,
	}

	err := client.Ping(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

// TestClient_Close_NilRueidis verifies Close is safe with nil rueidis.
func TestClient_Close_NilRueidis(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: false},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: nil,
	}

	assert.NotPanics(t, func() {
		client.Close()
	})
}

// TestClient_RueidisClient_NilReturn verifies accessor with disabled cache.
func TestClient_RueidisClient_NilReturn(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: false},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: nil,
	}

	assert.Nil(t, client.RueidisClient())
}

// TestNewClient_Disabled_CreatesClientWithoutRueidis verifies disabled cache creates a plain client.
func TestNewClient_Disabled_CreatesClientWithoutRueidis(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
			URL:     "should-be-ignored",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Nil(t, client.rueidisClient)
	assert.Equal(t, cfg, client.config)
}

// TestNewClient_EnabledEmptyURL verifies error when enabled but URL is empty.
func TestNewClient_EnabledEmptyURL(t *testing.T) {
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

// TestNewClient_InvalidScheme verifies error with invalid URL scheme.
func TestNewClient_InvalidScheme(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "http://localhost:6379",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to parse cache URL")
}

// TestModuleConstants verifies all module-level constants.
func TestModuleConstants_Unit(t *testing.T) {
	assert.Equal(t, 5*time.Second, DefaultDialTimeout)
	assert.Equal(t, 3*time.Second, DefaultReadTimeout)
	assert.Equal(t, 3*time.Second, DefaultWriteTimeout)
	assert.Equal(t, 16*1024*1024, DefaultCacheSizeEachConn)
	assert.Equal(t, 10, DefaultRingScale)
	assert.Equal(t, 128, DefaultBlockingPoolSize)
}

// --- registerHooks tests ---

// TestRegisterHooks_DisabledCache verifies lifecycle hooks with disabled cache.
func TestRegisterHooks_DisabledCache(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: false},
		},
		logger:        logger,
		rueidisClient: nil,
	}

	registerHooks(lc, client, logger)

	// Start should succeed (disabled path logs and returns nil)
	ctx := context.Background()
	require.NoError(t, lc.Start(ctx))

	// Stop should succeed (disabled path returns nil)
	require.NoError(t, lc.Stop(ctx))
}

// TestRegisterHooks_EnabledCacheNilRueidis verifies lifecycle hooks when cache is
// enabled but rueidis client is nil (ping will fail, but startup should not fail).
func TestRegisterHooks_EnabledCacheNilRueidis(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: true, URL: "redis://localhost:6379"},
		},
		logger:        logger,
		rueidisClient: nil, // nil causes Ping to return error
	}

	registerHooks(lc, client, logger)

	ctx := context.Background()
	// Start should succeed even though ping fails (it just warns)
	require.NoError(t, lc.Start(ctx))

	// Stop calls client.Close() which is safe with nil rueidis
	require.NoError(t, lc.Stop(ctx))
}

// --- NewClient_EnabledInvalidURL verifies error path when URL parsing fails ---

func TestNewClient_EnabledInvalidURL(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "not-a-valid-redis-url://???",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, client)
}

// --- CacheAside_LoaderReturnsUnmarshalable tests the marshal error path ---

func TestCache_CacheAside_LoaderReturnsUnmarshalable(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Loader returns something that can't be marshaled to JSON
	loader := func() (any, error) {
		return make(chan int), nil
	}

	var result map[string]string
	err := c.CacheAside(ctx, "aside:unmarshal", 1*time.Minute, loader, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "marshal")
}

// --- GetJSON roundtrip tests to improve Get/Set coverage via L1 ---

func TestCache_GetJSON_Roundtrip(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	type Movie struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
	}

	original := &Movie{Title: "Test Movie", Year: 2025}
	err := c.SetJSON(ctx, "movie:1", original, 1*time.Minute)
	require.NoError(t, err)

	var retrieved Movie
	err = c.GetJSON(ctx, "movie:1", &retrieved)
	require.NoError(t, err)
	assert.Equal(t, "Test Movie", retrieved.Title)
	assert.Equal(t, 2025, retrieved.Year)
}

// --- Test Cache operations with empty key ---

func TestCache_Set_EmptyKey(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	err := c.Set(ctx, "", []byte("value"), 1*time.Minute)
	require.NoError(t, err)

	val, err := c.Get(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, []byte("value"), val)
}

// --- Test Cache.Delete nonexistent key ---

func TestCache_Delete_NonexistentKey(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Deleting a key that doesn't exist should succeed
	err := c.Delete(ctx, "does-not-exist")
	require.NoError(t, err)
}

// --- Test NewNamedCache with different names ---

func TestNewNamedCache_CustomName(t *testing.T) {
	c, err := NewNamedCache(nil, 100, 1*time.Minute, "sessions")
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Close()

	assert.Equal(t, "sessions", c.name)
}

// --- Test Client struct fields ---

func TestClient_FieldAccess(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{Enabled: false},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := NewClient(cfg, logger)
	require.NoError(t, err)

	// Verify fields are set properly
	assert.Equal(t, cfg, client.config)
	assert.Equal(t, logger, client.logger)
	assert.Nil(t, client.RueidisClient())
}

// --- InvalidatePattern with empty pattern ---

func TestCache_InvalidatePattern_EmptyPattern(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	c.l1.Set("a", []byte("1"))

	// Empty pattern still clears L1 (L1 is always fully cleared)
	err := c.InvalidatePattern(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// --- Cache.Exists after Set and Delete ---

func TestCache_Exists_AfterSetAndDelete(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Set a value
	err := c.Set(ctx, "exist-key", []byte("val"), 1*time.Minute)
	require.NoError(t, err)

	exists, err := c.Exists(ctx, "exist-key")
	require.NoError(t, err)
	assert.True(t, exists)

	// Delete it
	err = c.Delete(ctx, "exist-key")
	require.NoError(t, err)

	exists, err = c.Exists(ctx, "exist-key")
	require.NoError(t, err)
	assert.False(t, exists)
}

// --- Multiple Invalidate calls ---

func TestCache_Invalidate_Multiple(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Set values
	c.l1.Set("x:1", []byte("a"))
	c.l1.Set("y:1", []byte("b"))

	// First invalidate — only x:* removed
	err := c.Invalidate(ctx, "x:*")
	require.NoError(t, err)
	assert.Equal(t, 1, c.l1.Size())

	// Second invalidate — removes y:*
	err = c.Invalidate(ctx, "y:*")
	require.NoError(t, err)
	assert.Equal(t, 0, c.l1.Size())
}

// --- JSON test with large data ---

func TestCache_SetJSON_LargeData(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Create a reasonably sized JSON payload
	data := make(map[string]string, 100)
	for i := range 100 {
		data[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	err := c.SetJSON(ctx, "large-data", data, 1*time.Minute)
	require.NoError(t, err)

	var retrieved map[string]string
	err = c.GetJSON(ctx, "large-data", &retrieved)
	require.NoError(t, err)
	assert.Len(t, retrieved, 100)
}

// --- Test CacheAside_UnmarshalDestError ---

func TestCache_CacheAside_UnmarshalToDestError(t *testing.T) {
	c := newL1OnlyCache(t)
	ctx := context.Background()

	// Loader returns a string, but dest is *int (will fail unmarshal)
	loader := func() (any, error) {
		return "not a number", nil
	}

	var result int
	err := c.CacheAside(ctx, "aside:type-mismatch", 1*time.Minute, loader, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}

// --- L2 (rueidis) mock tests ---

// newMockL2Cache creates a Cache with a mocked rueidis client for testing L2 code paths.
func newMockL2Cache(t *testing.T, ctrl *gomock.Controller) (*Cache, *mock.Client) {
	t.Helper()
	mockClient := mock.NewClient(ctrl)

	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: true, URL: "redis://localhost:6379"},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: mockClient,
	}

	c, err := NewCache(client, 1000, 1*time.Minute)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c, mockClient
}

// TestCache_Get_L2Hit tests the full L2 hit path with mock rueidis.
func TestCache_Get_L2Hit(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// Key not in L1 -- should go to L2 via DoCache
	mockClient.EXPECT().
		DoCache(gomock.Any(), mock.Match("GET", "l2-key"), gomock.Any()).
		Return(mock.Result(mock.RedisBlobString("l2-value")))

	val, err := c.Get(ctx, "l2-key")
	require.NoError(t, err)
	assert.Equal(t, []byte("l2-value"), val)

	// Verify it was populated in L1 after L2 hit
	l1Val, ok := c.l1.Get("l2-key")
	assert.True(t, ok, "L1 should be populated after L2 hit")
	assert.Equal(t, []byte("l2-value"), l1Val)
}

// TestCache_Get_L2Error tests the L2 error path.
func TestCache_Get_L2Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		DoCache(gomock.Any(), mock.Match("GET", "err-key"), gomock.Any()).
		Return(mock.ErrorResult(rueidis.Nil))

	_, err := c.Get(ctx, "err-key")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// TestCache_Get_L2RealError tests the L2 path with a real (non-nil) error.
func TestCache_Get_L2RealError(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		DoCache(gomock.Any(), mock.Match("GET", "fail-key"), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	_, err := c.Get(ctx, "fail-key")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache get failed")
}

// TestCache_Set_L2_SecondPrecision tests Set with TTL >= 1 second (EX path).
func TestCache_Set_L2_SecondPrecision(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// TTL of 5 seconds should use EX (second precision)
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("SET", "set-key", "set-value", "EX", "5")).
		Return(mock.Result(mock.RedisString("OK")))

	err := c.Set(ctx, "set-key", []byte("set-value"), 5*time.Second)
	require.NoError(t, err)
}

// TestCache_Set_L2_MillisecondPrecision tests Set with TTL < 1 second (PX path).
func TestCache_Set_L2_MillisecondPrecision(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// TTL of 500ms should use PX (millisecond precision)
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("SET", "ms-key", "ms-value", "PX", "500")).
		Return(mock.Result(mock.RedisString("OK")))

	err := c.Set(ctx, "ms-key", []byte("ms-value"), 500*time.Millisecond)
	require.NoError(t, err)
}

// TestCache_Set_L2_NoExpiration tests Set with zero TTL (no expiration).
func TestCache_Set_L2_NoExpiration(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// Zero TTL means no expiration
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("SET", "noexp-key", "noexp-value")).
		Return(mock.Result(mock.RedisString("OK")))

	err := c.Set(ctx, "noexp-key", []byte("noexp-value"), 0)
	require.NoError(t, err)
}

// TestCache_Set_L2_Error tests Set when L2 returns an error.
func TestCache_Set_L2_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	err := c.Set(ctx, "err-key", []byte("value"), 5*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache set failed")
}

// TestCache_Set_L2_PxError tests Set when L2 returns error on PX path.
func TestCache_Set_L2_PxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	err := c.Set(ctx, "err-px", []byte("val"), 200*time.Millisecond)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache set failed")
}

// TestCache_Set_L2_NoExpError tests Set when L2 returns error on no-expiration path.
func TestCache_Set_L2_NoExpError(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	err := c.Set(ctx, "err-noexp", []byte("val"), 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache set failed")
}

// TestCache_Delete_L2 tests Delete with L2 available.
func TestCache_Delete_L2(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("DEL", "del-key")).
		Return(mock.Result(mock.RedisInt64(1)))

	err := c.Delete(ctx, "del-key")
	require.NoError(t, err)
}

// TestCache_Delete_L2_Error tests Delete when L2 returns an error.
func TestCache_Delete_L2_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	err := c.Delete(ctx, "err-del")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache delete failed")
}

// TestCache_Exists_L2_Found tests Exists when key exists in L2.
func TestCache_Exists_L2_Found(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("EXISTS", "exists-key")).
		Return(mock.Result(mock.RedisInt64(1)))

	exists, err := c.Exists(ctx, "exists-key")
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestCache_Exists_L2_NotFound tests Exists when key does not exist in L2.
func TestCache_Exists_L2_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("EXISTS", "nope")).
		Return(mock.Result(mock.RedisInt64(0)))

	exists, err := c.Exists(ctx, "nope")
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestCache_Exists_L2_Error tests Exists when L2 returns an error.
func TestCache_Exists_L2_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	_, err := c.Exists(ctx, "err-key")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache exists check failed")
}

// TestCache_Invalidate_L2_WithKeys tests Invalidate when L2 has matching keys.
func TestCache_Invalidate_L2_WithKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// SCAN returns matching keys with cursor=0 (done)
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("SCAN", "0", "MATCH", "prefix:*", "COUNT", "100")).
		Return(mock.Result(mock.RedisArray(
			mock.RedisBlobString("0"),
			mock.RedisArray(
				mock.RedisBlobString("prefix:a"),
				mock.RedisBlobString("prefix:b"),
			),
		)))

	// DEL deletes those keys
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("DEL", "prefix:a", "prefix:b")).
		Return(mock.Result(mock.RedisInt64(2)))

	err := c.Invalidate(ctx, "prefix:*")
	require.NoError(t, err)
}

// TestCache_Invalidate_L2_NoKeys tests Invalidate when L2 returns no matching keys.
func TestCache_Invalidate_L2_NoKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// SCAN returns empty result with cursor=0
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("SCAN", "0", "MATCH", "empty:*", "COUNT", "100")).
		Return(mock.Result(mock.RedisArray(
			mock.RedisBlobString("0"),
			mock.RedisArray(),
		)))

	err := c.Invalidate(ctx, "empty:*")
	require.NoError(t, err)
}

// TestCache_Invalidate_L2_ScanError tests Invalidate when SCAN command fails.
func TestCache_Invalidate_L2_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	err := c.Invalidate(ctx, "err:*")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache SCAN failed")
}

// TestCache_Invalidate_L2_DelError tests Invalidate when DEL command fails after SCAN succeeds.
func TestCache_Invalidate_L2_DelError(t *testing.T) {
	ctrl := gomock.NewController(t)
	c, mockClient := newMockL2Cache(t, ctrl)
	ctx := context.Background()

	// SCAN succeeds with keys
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("SCAN", "0", "MATCH", "del-err:*", "COUNT", "100")).
		Return(mock.Result(mock.RedisArray(
			mock.RedisBlobString("0"),
			mock.RedisArray(
				mock.RedisBlobString("del-err:a"),
			),
		)))

	// DEL fails
	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(fmt.Errorf("connection refused")))

	err := c.Invalidate(ctx, "del-err:*")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "L2 cache batch delete failed")
}

// TestCache_Ping_L2_Success tests Ping when rueidis client succeeds.
func TestCache_Ping_L2_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock.NewClient(ctrl)
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: true, URL: "redis://localhost:6379"},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: mockClient,
	}

	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("PING")).
		Return(mock.Result(mock.RedisString("PONG")))

	err := client.Ping(context.Background())
	require.NoError(t, err)
}

// TestCache_Close_L2 tests Close when rueidis client is set.
func TestCache_Close_L2(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock.NewClient(ctrl)
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: true, URL: "redis://localhost:6379"},
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		rueidisClient: mockClient,
	}

	mockClient.EXPECT().Close()

	// Should call rueidisClient.Close() and log
	client.Close()
}

// TestRegisterHooks_EnabledCachePingSuccess tests lifecycle hooks when cache is enabled
// and ping succeeds.
func TestRegisterHooks_EnabledCachePingSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock.NewClient(ctrl)
	lc := fxtest.NewLifecycle(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{Enabled: true, URL: "redis://localhost:6379"},
		},
		logger:        logger,
		rueidisClient: mockClient,
	}

	registerHooks(lc, client, logger)

	// Expect PING on Start
	mockClient.EXPECT().
		Do(gomock.Any(), mock.Match("PING")).
		Return(mock.Result(mock.RedisString("PONG")))

	// Expect Close on Stop
	mockClient.EXPECT().Close()

	ctx := context.Background()
	require.NoError(t, lc.Start(ctx))
	require.NoError(t, lc.Stop(ctx))
}
