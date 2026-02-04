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

// helper to create a test cache with L1 only
func newTestCache(t *testing.T) *Cache {
	t.Helper()
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
	return cache
}

// TestInvalidateSession tests session cache invalidation
func TestInvalidateSession(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	tokenHash := "abc123hash"

	// Set session in cache
	sessionKey := SessionKey(tokenHash)
	err := cache.Set(ctx, sessionKey, []byte(`{"user_id":"user1"}`), 1*time.Minute)
	require.NoError(t, err)

	// Verify it exists
	exists, err := cache.Exists(ctx, sessionKey)
	require.NoError(t, err)
	assert.True(t, exists)

	// Invalidate session
	err = cache.InvalidateSession(ctx, tokenHash)
	require.NoError(t, err)

	// Should be gone from L1 (Delete clears L1)
	_, ok := cache.l1.Get(sessionKey)
	assert.False(t, ok)
}

// TestInvalidateUserSessions tests invalidation of all sessions for a user
func TestInvalidateUserSessions(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	userID := "user-uuid-123"

	// Set multiple session keys for the user
	err := cache.Set(ctx, SessionByUserKey(userID)+":session1", []byte("data1"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, SessionByUserKey(userID)+":session2", []byte("data2"), 1*time.Minute)
	require.NoError(t, err)
	// Also set a session for a different user
	err = cache.Set(ctx, SessionByUserKey("other-user")+":session1", []byte("other"), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate all sessions for the user
	err = cache.InvalidateUserSessions(ctx, userID)
	require.NoError(t, err)

	// L1 should be cleared (pattern invalidation clears entire L1)
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateRBACForUser tests RBAC cache invalidation for a user
func TestInvalidateRBACForUser(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	userID := "user-uuid-456"

	// Set RBAC cache entries for the user
	err := cache.Set(ctx, RBACEnforceKey(userID, "movies", "read"), []byte("allowed"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, RBACUserRolesKey(userID), []byte(`["admin","user"]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, RBACUserPermsKey(userID), []byte(`["read","write"]`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate RBAC for the user
	err = cache.InvalidateRBACForUser(ctx, userID)
	require.NoError(t, err)

	// L1 should be cleared (pattern invalidation clears entire L1)
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateAllRBAC tests invalidation of all RBAC cache entries
func TestInvalidateAllRBAC(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	// Set multiple RBAC entries for different users
	err := cache.Set(ctx, RBACUserRolesKey("user1"), []byte(`["admin"]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, RBACUserRolesKey("user2"), []byte(`["viewer"]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, RBACEnforceKey("user1", "movies", "read"), []byte("true"), 1*time.Minute)
	require.NoError(t, err)

	// Also set non-RBAC entry
	err = cache.Set(ctx, "other:key", []byte("value"), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate all RBAC
	err = cache.InvalidateAllRBAC(ctx)
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateServerSettings tests server settings cache invalidation
func TestInvalidateServerSettings(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	// Set server settings
	err := cache.Set(ctx, ServerSettingKey("theme.mode"), []byte("dark"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, ServerSettingKey("locale"), []byte("en-US"), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate all server settings
	err = cache.InvalidateServerSettings(ctx)
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateServerSetting tests single server setting invalidation
func TestInvalidateServerSetting(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	// Set a server setting
	key := "theme.mode"
	settingKey := ServerSettingKey(key)
	err := cache.Set(ctx, settingKey, []byte("dark"), 1*time.Minute)
	require.NoError(t, err)

	// Verify it exists
	exists, err := cache.Exists(ctx, settingKey)
	require.NoError(t, err)
	assert.True(t, exists)

	// Invalidate the specific setting
	err = cache.InvalidateServerSetting(ctx, key)
	require.NoError(t, err)

	// Should be gone
	_, ok := cache.l1.Get(settingKey)
	assert.False(t, ok)
}

// TestInvalidateUserSettings tests user settings cache invalidation
func TestInvalidateUserSettings(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	userID := "user-uuid-789"

	// Set user settings
	err := cache.Set(ctx, UserSettingKey(userID, "notifications.enabled"), []byte("true"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, UserSettingKey(userID, "theme.mode"), []byte("dark"), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate user settings
	err = cache.InvalidateUserSettings(ctx, userID)
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateMovie tests movie cache invalidation
func TestInvalidateMovie(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	movieID := "movie-uuid-123"

	// Set movie-related cache entries
	err := cache.Set(ctx, MovieKey(movieID), []byte(`{"title":"Test Movie"}`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, MovieCastKey(movieID), []byte(`[{"name":"Actor"}]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, MovieCrewKey(movieID), []byte(`[{"name":"Director"}]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, MovieGenresKey(movieID), []byte(`["Action","Drama"]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, MovieFilesKey(movieID), []byte(`[{"path":"/movie.mkv"}]`), 1*time.Minute)
	require.NoError(t, err)

	// Also set movie lists
	err = cache.Set(ctx, MovieRecentKey(10, 0), []byte(`[{"id":"movie-1"}]`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate movie
	err = cache.InvalidateMovie(ctx, movieID)
	require.NoError(t, err)

	// L1 should be cleared (movie invalidation also clears lists)
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateMovieLists tests movie lists cache invalidation
func TestInvalidateMovieLists(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	// Set movie list entries
	err := cache.Set(ctx, MovieListKey("filter-hash-1"), []byte(`[{"id":"movie-1"}]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, MovieRecentKey(10, 0), []byte(`[{"id":"movie-2"}]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, MovieTopRatedKey(100, 10, 0), []byte(`[{"id":"movie-3"}]`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate movie lists
	err = cache.InvalidateMovieLists(ctx)
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateSearch tests search cache invalidation
func TestInvalidateSearch(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	// Set search cache entries
	err := cache.Set(ctx, SearchMoviesKey("query-hash-1"), []byte(`[{"id":"movie-1"}]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, SearchAutocompleteKey("the"), []byte(`["The Matrix","The Godfather"]`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate search
	err = cache.InvalidateSearch(ctx)
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateLibrary tests library cache invalidation
func TestInvalidateLibrary(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	libraryID := "library-uuid-123"

	// Set library cache entries
	err := cache.Set(ctx, LibraryKey(libraryID), []byte(`{"name":"Movies"}`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, LibraryStatsKey(libraryID), []byte(`{"count":100}`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate library
	err = cache.InvalidateLibrary(ctx, libraryID)
	require.NoError(t, err)

	// Library invalidation uses Delete, not Invalidate pattern
	// Check specific keys are gone
	_, ok := cache.l1.Get(LibraryKey(libraryID))
	assert.False(t, ok)
	_, ok = cache.l1.Get(LibraryStatsKey(libraryID))
	assert.False(t, ok)
}

// TestInvalidateUser tests user cache invalidation
func TestInvalidateUser(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	userID := "user-uuid-999"

	// Set user cache entries
	err := cache.Set(ctx, UserKey(userID), []byte(`{"username":"testuser"}`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, ContinueWatchingKey(userID, 10), []byte(`[{"id":"movie-1"}]`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate user
	err = cache.InvalidateUser(ctx, userID)
	require.NoError(t, err)

	// L1 should be cleared (pattern invalidation clears entire L1)
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidateContinueWatching tests continue watching cache invalidation
func TestInvalidateContinueWatching(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()
	userID := "user-uuid-888"

	// Set continue watching entries
	err := cache.Set(ctx, ContinueWatchingKey(userID, 5), []byte(`[{"id":"movie-1"}]`), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, ContinueWatchingKey(userID, 10), []byte(`[{"id":"movie-1"},{"id":"movie-2"}]`), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate continue watching
	err = cache.InvalidateContinueWatching(ctx, userID)
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestInvalidatePattern tests InvalidatePattern alias
func TestInvalidatePattern(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	// Set some entries
	err := cache.Set(ctx, "prefix:1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, "prefix:2", []byte("value2"), 1*time.Minute)
	require.NoError(t, err)

	// Verify entries exist
	assert.Greater(t, cache.l1.Size(), 0)

	// Invalidate pattern
	err = cache.InvalidatePattern(ctx, "prefix:*")
	require.NoError(t, err)

	// L1 should be cleared
	assert.Equal(t, 0, cache.l1.Size())
}

// TestCacheAside tests the cache-aside pattern helper
func TestCacheAside(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	type Movie struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	loaderCalled := 0
	loader := func() (interface{}, error) {
		loaderCalled++
		return &Movie{ID: "movie-1", Title: "Test Movie"}, nil
	}

	// First call - cache miss, should call loader
	var movie1 Movie
	err := cache.CacheAside(ctx, "movie:1", 1*time.Minute, loader, &movie1)
	require.NoError(t, err)
	assert.Equal(t, 1, loaderCalled)
	assert.Equal(t, "movie-1", movie1.ID)
	assert.Equal(t, "Test Movie", movie1.Title)

	// Wait for async cache set
	time.Sleep(50 * time.Millisecond)

	// Second call - cache hit, should NOT call loader
	var movie2 Movie
	err = cache.CacheAside(ctx, "movie:1", 1*time.Minute, loader, &movie2)
	require.NoError(t, err)
	assert.Equal(t, 1, loaderCalled, "Loader should not be called on cache hit")
	assert.Equal(t, "movie-1", movie2.ID)
}

// TestCacheAside_LoaderError tests CacheAside when loader returns error
func TestCacheAside_LoaderError(t *testing.T) {
	cache := newTestCache(t)
	defer cache.Close()

	ctx := context.Background()

	loader := func() (interface{}, error) {
		return nil, assert.AnError
	}

	var result map[string]string
	err := cache.CacheAside(ctx, "error:key", 1*time.Minute, loader, &result)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

// TestNewNamedCache tests creating a named cache
func TestNewNamedCache(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	cache, err := NewNamedCache(client, 100, 1*time.Minute, "movies")
	require.NoError(t, err)
	require.NotNil(t, cache)
	defer cache.Close()

	assert.Equal(t, "movies", cache.name)
}

// TestCache_SetShortTTL tests that short TTLs skip L1
func TestCache_SetShortTTL(t *testing.T) {
	client := &Client{
		config: &config.Config{
			Cache: config.CacheConfig{
				Enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	// Create cache with 1 minute L1 TTL
	cache, err := NewCache(client, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set with TTL shorter than L1 TTL - should delete from L1
	err = cache.Set(ctx, "short_ttl", []byte("value"), 10*time.Second)
	require.NoError(t, err)

	// Should NOT be in L1 (TTL < L1 TTL causes delete)
	_, ok := cache.l1.Get("short_ttl")
	assert.False(t, ok, "Short TTL should skip L1 to prevent stale reads")

	// Set with TTL equal to L1 TTL - should be in L1
	err = cache.Set(ctx, "equal_ttl", []byte("value"), 1*time.Minute)
	require.NoError(t, err)

	_, ok = cache.l1.Get("equal_ttl")
	assert.True(t, ok, "Equal TTL should be stored in L1")

	// Set with TTL longer than L1 TTL - should be in L1
	err = cache.Set(ctx, "long_ttl", []byte("value"), 5*time.Minute)
	require.NoError(t, err)

	_, ok = cache.l1.Get("long_ttl")
	assert.True(t, ok, "Long TTL should be stored in L1")
}

// TestCache_SetSubSecondTTL tests sub-second TTL handling
func TestCache_SetSubSecondTTL(t *testing.T) {
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

	// Set with sub-second TTL (milliseconds)
	err = cache.Set(ctx, "subsecond", []byte("value"), 500*time.Millisecond)
	require.NoError(t, err)

	// Should NOT be in L1 (short TTL)
	_, ok := cache.l1.Get("subsecond")
	assert.False(t, ok, "Sub-second TTL should skip L1")
}
