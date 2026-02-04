// Package cache provides unified L1 (otter) + L2 (rueidis) caching operations.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// KeyPrefix defines standard cache key prefixes for different data types.
const (
	// Session cache keys
	KeyPrefixSession       = "session:"
	KeyPrefixSessionByUser = "session:user:"

	// RBAC cache keys
	KeyPrefixRBACPolicy     = "rbac:policy:"
	KeyPrefixRBACEnforce    = "rbac:enforce:"
	KeyPrefixRBACUserRoles  = "rbac:roles:"
	KeyPrefixRBACUserPerms  = "rbac:perms:"

	// Settings cache keys
	KeyPrefixServerSetting = "settings:server:"
	KeyPrefixUserSetting   = "settings:user:"

	// User cache keys
	KeyPrefixUser       = "user:"
	KeyPrefixUserByName = "user:name:"
	KeyPrefixUserEmail  = "user:email:"

	// Content cache keys
	KeyPrefixMovie            = "movie:"
	KeyPrefixMovieCast        = "movie:cast:"
	KeyPrefixMovieCrew        = "movie:crew:"
	KeyPrefixMovieGenres      = "movie:genres:"
	KeyPrefixMovieFiles       = "movie:files:"
	KeyPrefixMovieList        = "movie:list:"
	KeyPrefixMovieRecent      = "movie:recent"
	KeyPrefixMovieTopRated    = "movie:toprated"
	KeyPrefixMovieMeta        = "movie:meta:"
	KeyPrefixLibrary          = "library:"
	KeyPrefixLibraryStats     = "library:stats:"
	KeyPrefixSearch           = "search:"
	KeyPrefixSearchMovies     = "search:movies:"
	KeyPrefixSearchAutocomplete = "search:autocomplete:"
	KeyPrefixImage            = "image:"
	KeyPrefixContinueWatching = "user:continue:"
)

// DefaultTTLs for different cache types.
const (
	// SessionTTL is the TTL for session cache entries.
	// Sessions are validated frequently, so use a short TTL.
	SessionTTL = 30 * time.Second

	// RBACPolicyTTL is the TTL for RBAC policy cache entries.
	// Policies change infrequently, use longer TTL.
	RBACPolicyTTL = 5 * time.Minute

	// RBACEnforceTTL is the TTL for RBAC enforcement results.
	// Short TTL to balance performance vs policy freshness.
	RBACEnforceTTL = 30 * time.Second

	// ServerSettingsTTL is the TTL for server settings.
	// Settings rarely change, use longer TTL.
	ServerSettingsTTL = 5 * time.Minute

	// UserSettingsTTL is the TTL for user settings.
	UserSettingsTTL = 2 * time.Minute

	// UserTTL is the TTL for user data cache.
	UserTTL = 1 * time.Minute

	// MovieTTL is the TTL for movie data cache.
	// Movies are read frequently, 5 min TTL.
	MovieTTL = 5 * time.Minute

	// MovieMetaTTL is the TTL for movie metadata.
	MovieMetaTTL = 10 * time.Minute

	// LibraryStatsTTL is the TTL for library statistics.
	// Stats are expensive to compute, use longer TTL.
	LibraryStatsTTL = 10 * time.Minute

	// SearchResultsTTL is the TTL for search results.
	// Search results change with index updates, use short TTL.
	SearchResultsTTL = 30 * time.Second

	// ImageMetaTTL is the TTL for image metadata (not the image bytes).
	ImageMetaTTL = 24 * time.Hour

	// ContinueWatchingTTL is the TTL for user's continue watching list.
	// Per-user, changes with watch progress updates.
	ContinueWatchingTTL = 1 * time.Minute

	// RecentlyAddedTTL is the TTL for recently added movies.
	// Homepage hot path, short TTL for freshness.
	RecentlyAddedTTL = 2 * time.Minute

	// TopRatedTTL is the TTL for top rated movies list.
	TopRatedTTL = 5 * time.Minute
)

// SessionKey returns the cache key for a session by token hash.
func SessionKey(tokenHash string) string {
	return KeyPrefixSession + tokenHash
}

// SessionByUserKey returns the cache key for sessions by user ID.
func SessionByUserKey(userID string) string {
	return KeyPrefixSessionByUser + userID
}

// RBACEnforceKey returns the cache key for an RBAC enforcement result.
func RBACEnforceKey(subject, object, action string) string {
	return fmt.Sprintf("%s%s:%s:%s", KeyPrefixRBACEnforce, subject, object, action)
}

// RBACUserRolesKey returns the cache key for a user's roles.
func RBACUserRolesKey(userID string) string {
	return KeyPrefixRBACUserRoles + userID
}

// RBACUserPermsKey returns the cache key for a user's permissions.
func RBACUserPermsKey(userID string) string {
	return KeyPrefixRBACUserPerms + userID
}

// ServerSettingKey returns the cache key for a server setting.
func ServerSettingKey(key string) string {
	return KeyPrefixServerSetting + key
}

// UserSettingKey returns the cache key for a user setting.
func UserSettingKey(userID, key string) string {
	return fmt.Sprintf("%s%s:%s", KeyPrefixUserSetting, userID, key)
}

// UserKey returns the cache key for a user by ID.
func UserKey(userID string) string {
	return KeyPrefixUser + userID
}

// UserByNameKey returns the cache key for a user by username.
func UserByNameKey(username string) string {
	return KeyPrefixUserByName + username
}

// MovieKey returns the cache key for a movie by ID.
func MovieKey(movieID string) string {
	return KeyPrefixMovie + movieID
}

// MovieCastKey returns the cache key for a movie's cast.
func MovieCastKey(movieID string) string {
	return KeyPrefixMovieCast + movieID
}

// MovieCrewKey returns the cache key for a movie's crew.
func MovieCrewKey(movieID string) string {
	return KeyPrefixMovieCrew + movieID
}

// MovieGenresKey returns the cache key for a movie's genres.
func MovieGenresKey(movieID string) string {
	return KeyPrefixMovieGenres + movieID
}

// MovieFilesKey returns the cache key for a movie's files.
func MovieFilesKey(movieID string) string {
	return KeyPrefixMovieFiles + movieID
}

// MovieListKey returns the cache key for a movie list with filters hash.
func MovieListKey(hash string) string {
	return KeyPrefixMovieList + hash
}

// MovieRecentKey returns the cache key for recently added movies.
func MovieRecentKey(limit, offset int32) string {
	return fmt.Sprintf("%s:%d:%d", KeyPrefixMovieRecent, limit, offset)
}

// MovieTopRatedKey returns the cache key for top rated movies.
func MovieTopRatedKey(minVotes, limit, offset int32) string {
	return fmt.Sprintf("%s:%d:%d:%d", KeyPrefixMovieTopRated, minVotes, limit, offset)
}

// MovieMetaKey returns the cache key for movie metadata by external ID.
func MovieMetaKey(provider, externalID string) string {
	return fmt.Sprintf("%s%s:%s", KeyPrefixMovieMeta, provider, externalID)
}

// LibraryKey returns the cache key for a library by ID.
func LibraryKey(libraryID string) string {
	return KeyPrefixLibrary + libraryID
}

// LibraryStatsKey returns the cache key for library stats by ID.
func LibraryStatsKey(libraryID string) string {
	return KeyPrefixLibraryStats + libraryID
}

// SearchMoviesKey returns the cache key for movie search results.
func SearchMoviesKey(hash string) string {
	return KeyPrefixSearchMovies + hash
}

// SearchAutocompleteKey returns the cache key for autocomplete results.
func SearchAutocompleteKey(query string) string {
	return KeyPrefixSearchAutocomplete + query
}

// ImageKey returns the cache key for an image (metadata only, not bytes).
func ImageKey(imageType, size, path string) string {
	return fmt.Sprintf("%s%s:%s:%s", KeyPrefixImage, imageType, size, path)
}

// ContinueWatchingKey returns the cache key for a user's continue watching list.
func ContinueWatchingKey(userID string, limit int32) string {
	return fmt.Sprintf("%s%s:%d", KeyPrefixContinueWatching, userID, limit)
}

// CacheAside is a helper that implements the cache-aside pattern.
// It first checks the cache, and on miss, calls the loader function.
func (c *Cache) CacheAside(ctx context.Context, key string, ttl time.Duration, loader func() (interface{}, error), dest interface{}) error {
	// Try cache first
	err := c.GetJSON(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// Cache miss - load from source
	value, err := loader()
	if err != nil {
		return err
	}

	// Store in cache (async to not block response)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := c.SetJSON(cacheCtx, key, value, ttl); setErr != nil {
			// Log error but don't fail the request
			_ = setErr
		}
	}()

	// Copy value to dest through JSON marshaling
	// This is a simple way to copy the value without reflection
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal loaded value: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal to destination: %w", err)
	}

	return nil
}

// InvalidatePattern invalidates all cache keys matching a pattern.
// Pattern supports Redis glob-style patterns: * matches any sequence of characters.
func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
	return c.Invalidate(ctx, pattern)
}

// InvalidateSession invalidates a specific session from cache.
func (c *Cache) InvalidateSession(ctx context.Context, tokenHash string) error {
	return c.Delete(ctx, SessionKey(tokenHash))
}

// InvalidateUserSessions invalidates all sessions for a user.
func (c *Cache) InvalidateUserSessions(ctx context.Context, userID string) error {
	return c.Invalidate(ctx, KeyPrefixSessionByUser+userID+"*")
}

// InvalidateRBACForUser invalidates all RBAC cache entries for a user.
func (c *Cache) InvalidateRBACForUser(ctx context.Context, userID string) error {
	// Invalidate roles, permissions, and enforcement results for this user
	patterns := []string{
		KeyPrefixRBACEnforce + userID + ":*",
		KeyPrefixRBACUserRoles + userID,
		KeyPrefixRBACUserPerms + userID,
	}

	for _, pattern := range patterns {
		if err := c.Invalidate(ctx, pattern); err != nil {
			return fmt.Errorf("failed to invalidate RBAC cache pattern %s: %w", pattern, err)
		}
	}

	return nil
}

// InvalidateAllRBAC invalidates all RBAC cache entries (e.g., after policy reload).
func (c *Cache) InvalidateAllRBAC(ctx context.Context) error {
	return c.Invalidate(ctx, "rbac:*")
}

// InvalidateServerSettings invalidates all server settings cache.
func (c *Cache) InvalidateServerSettings(ctx context.Context) error {
	return c.Invalidate(ctx, KeyPrefixServerSetting+"*")
}

// InvalidateServerSetting invalidates a specific server setting.
func (c *Cache) InvalidateServerSetting(ctx context.Context, key string) error {
	return c.Delete(ctx, ServerSettingKey(key))
}

// InvalidateUserSettings invalidates all settings for a user.
func (c *Cache) InvalidateUserSettings(ctx context.Context, userID string) error {
	return c.Invalidate(ctx, KeyPrefixUserSetting+userID+":*")
}

// InvalidateMovie invalidates all cache entries for a movie.
func (c *Cache) InvalidateMovie(ctx context.Context, movieID string) error {
	patterns := []string{
		KeyPrefixMovie + movieID,
		KeyPrefixMovieCast + movieID,
		KeyPrefixMovieCrew + movieID,
		KeyPrefixMovieGenres + movieID,
		KeyPrefixMovieFiles + movieID,
	}

	for _, pattern := range patterns {
		if err := c.Delete(ctx, pattern); err != nil {
			return fmt.Errorf("failed to invalidate movie cache %s: %w", pattern, err)
		}
	}

	// Also invalidate movie lists since they may contain this movie
	return c.InvalidateMovieLists(ctx)
}

// InvalidateMovieLists invalidates all movie list caches (recently added, top rated, etc).
func (c *Cache) InvalidateMovieLists(ctx context.Context) error {
	patterns := []string{
		KeyPrefixMovieList + "*",
		KeyPrefixMovieRecent + "*",
		KeyPrefixMovieTopRated + "*",
	}

	for _, pattern := range patterns {
		if err := c.Invalidate(ctx, pattern); err != nil {
			return fmt.Errorf("failed to invalidate movie list cache %s: %w", pattern, err)
		}
	}

	return nil
}

// InvalidateSearch invalidates all search caches.
func (c *Cache) InvalidateSearch(ctx context.Context) error {
	return c.Invalidate(ctx, KeyPrefixSearch+"*")
}

// InvalidateLibrary invalidates all cache entries for a library.
func (c *Cache) InvalidateLibrary(ctx context.Context, libraryID string) error {
	patterns := []string{
		KeyPrefixLibrary + libraryID,
		KeyPrefixLibraryStats + libraryID,
	}

	for _, pattern := range patterns {
		if err := c.Delete(ctx, pattern); err != nil {
			return fmt.Errorf("failed to invalidate library cache %s: %w", pattern, err)
		}
	}

	return nil
}

// InvalidateUser invalidates all cache entries for a user.
func (c *Cache) InvalidateUser(ctx context.Context, userID string) error {
	patterns := []string{
		KeyPrefixUser + userID,
		KeyPrefixContinueWatching + userID + ":*",
	}

	for _, pattern := range patterns {
		if err := c.Invalidate(ctx, pattern); err != nil {
			return fmt.Errorf("failed to invalidate user cache %s: %w", pattern, err)
		}
	}

	return nil
}

// InvalidateContinueWatching invalidates a user's continue watching cache.
func (c *Cache) InvalidateContinueWatching(ctx context.Context, userID string) error {
	return c.Invalidate(ctx, KeyPrefixContinueWatching+userID+":*")
}
