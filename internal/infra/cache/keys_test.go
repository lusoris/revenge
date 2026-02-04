package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionKey(t *testing.T) {
	key := SessionKey("abc123")
	assert.Equal(t, "session:abc123", key)
}

func TestSessionByUserKey(t *testing.T) {
	key := SessionByUserKey("user-uuid")
	assert.Equal(t, "session:user:user-uuid", key)
}

func TestRBACEnforceKey(t *testing.T) {
	key := RBACEnforceKey("user1", "resource", "read")
	assert.Equal(t, "rbac:enforce:user1:resource:read", key)
}

func TestRBACUserRolesKey(t *testing.T) {
	key := RBACUserRolesKey("user-uuid")
	assert.Equal(t, "rbac:roles:user-uuid", key)
}

func TestRBACUserPermsKey(t *testing.T) {
	key := RBACUserPermsKey("user-uuid")
	assert.Equal(t, "rbac:perms:user-uuid", key)
}

func TestServerSettingKey(t *testing.T) {
	key := ServerSettingKey("theme.mode")
	assert.Equal(t, "settings:server:theme.mode", key)
}

func TestUserSettingKey(t *testing.T) {
	key := UserSettingKey("user-uuid", "notifications.enabled")
	assert.Equal(t, "settings:user:user-uuid:notifications.enabled", key)
}

func TestUserKey(t *testing.T) {
	key := UserKey("user-uuid")
	assert.Equal(t, "user:user-uuid", key)
}

func TestUserByNameKey(t *testing.T) {
	key := UserByNameKey("johndoe")
	assert.Equal(t, "user:name:johndoe", key)
}

func TestMovieKey(t *testing.T) {
	key := MovieKey("movie-uuid")
	assert.Equal(t, "movie:movie-uuid", key)
}

func TestMovieMetaKey(t *testing.T) {
	key := MovieMetaKey("tmdb", "12345")
	assert.Equal(t, "movie:meta:tmdb:12345", key)
}

func TestDefaultTTLs(t *testing.T) {
	// Verify TTLs are reasonable values
	assert.Equal(t, 30*time.Second, SessionTTL)
	assert.Equal(t, 5*time.Minute, RBACPolicyTTL)
	assert.Equal(t, 30*time.Second, RBACEnforceTTL)
	assert.Equal(t, 5*time.Minute, ServerSettingsTTL)
	assert.Equal(t, 2*time.Minute, UserSettingsTTL)
	assert.Equal(t, 1*time.Minute, UserTTL)
	assert.Equal(t, 10*time.Minute, MovieMetaTTL)
}

func TestKeyPrefixes(t *testing.T) {
	// Verify key prefixes are correct
	assert.Equal(t, "session:", KeyPrefixSession)
	assert.Equal(t, "session:user:", KeyPrefixSessionByUser)
	assert.Equal(t, "rbac:policy:", KeyPrefixRBACPolicy)
	assert.Equal(t, "rbac:enforce:", KeyPrefixRBACEnforce)
	assert.Equal(t, "rbac:roles:", KeyPrefixRBACUserRoles)
	assert.Equal(t, "rbac:perms:", KeyPrefixRBACUserPerms)
	assert.Equal(t, "settings:server:", KeyPrefixServerSetting)
	assert.Equal(t, "settings:user:", KeyPrefixUserSetting)
	assert.Equal(t, "user:", KeyPrefixUser)
	assert.Equal(t, "user:name:", KeyPrefixUserByName)
	assert.Equal(t, "user:email:", KeyPrefixUserEmail)
	assert.Equal(t, "movie:", KeyPrefixMovie)
	assert.Equal(t, "movie:meta:", KeyPrefixMovieMeta)
}

func TestMovieCastKey(t *testing.T) {
	key := MovieCastKey("movie-uuid")
	assert.Equal(t, "movie:cast:movie-uuid", key)
}

func TestMovieCrewKey(t *testing.T) {
	key := MovieCrewKey("movie-uuid")
	assert.Equal(t, "movie:crew:movie-uuid", key)
}

func TestMovieGenresKey(t *testing.T) {
	key := MovieGenresKey("movie-uuid")
	assert.Equal(t, "movie:genres:movie-uuid", key)
}

func TestMovieFilesKey(t *testing.T) {
	key := MovieFilesKey("movie-uuid")
	assert.Equal(t, "movie:files:movie-uuid", key)
}

func TestMovieListKey(t *testing.T) {
	key := MovieListKey("filter-hash-123")
	assert.Equal(t, "movie:list:filter-hash-123", key)
}

func TestMovieRecentKey(t *testing.T) {
	key := MovieRecentKey(10, 0)
	assert.Equal(t, "movie:recent:10:0", key)

	key = MovieRecentKey(25, 50)
	assert.Equal(t, "movie:recent:25:50", key)
}

func TestMovieTopRatedKey(t *testing.T) {
	key := MovieTopRatedKey(100, 10, 0)
	assert.Equal(t, "movie:toprated:100:10:0", key)

	key = MovieTopRatedKey(50, 20, 40)
	assert.Equal(t, "movie:toprated:50:20:40", key)
}

func TestLibraryKey(t *testing.T) {
	key := LibraryKey("library-uuid")
	assert.Equal(t, "library:library-uuid", key)
}

func TestLibraryStatsKey(t *testing.T) {
	key := LibraryStatsKey("library-uuid")
	assert.Equal(t, "library:stats:library-uuid", key)
}

func TestSearchMoviesKey(t *testing.T) {
	key := SearchMoviesKey("query-hash")
	assert.Equal(t, "search:movies:query-hash", key)
}

func TestSearchAutocompleteKey(t *testing.T) {
	key := SearchAutocompleteKey("the matrix")
	assert.Equal(t, "search:autocomplete:the matrix", key)
}

func TestImageKey(t *testing.T) {
	key := ImageKey("poster", "w500", "/path/to/image.jpg")
	assert.Equal(t, "image:poster:w500:/path/to/image.jpg", key)
}

func TestContinueWatchingKey(t *testing.T) {
	key := ContinueWatchingKey("user-uuid", 10)
	assert.Equal(t, "user:continue:user-uuid:10", key)
}

func TestAdditionalKeyPrefixes(t *testing.T) {
	// Verify additional key prefixes
	assert.Equal(t, "movie:cast:", KeyPrefixMovieCast)
	assert.Equal(t, "movie:crew:", KeyPrefixMovieCrew)
	assert.Equal(t, "movie:genres:", KeyPrefixMovieGenres)
	assert.Equal(t, "movie:files:", KeyPrefixMovieFiles)
	assert.Equal(t, "movie:list:", KeyPrefixMovieList)
	assert.Equal(t, "movie:recent", KeyPrefixMovieRecent)
	assert.Equal(t, "movie:toprated", KeyPrefixMovieTopRated)
	assert.Equal(t, "library:", KeyPrefixLibrary)
	assert.Equal(t, "library:stats:", KeyPrefixLibraryStats)
	assert.Equal(t, "search:", KeyPrefixSearch)
	assert.Equal(t, "search:movies:", KeyPrefixSearchMovies)
	assert.Equal(t, "search:autocomplete:", KeyPrefixSearchAutocomplete)
	assert.Equal(t, "image:", KeyPrefixImage)
	assert.Equal(t, "user:continue:", KeyPrefixContinueWatching)
}

func TestAdditionalTTLs(t *testing.T) {
	// Verify additional TTL values
	assert.Equal(t, 5*time.Minute, MovieTTL)
	assert.Equal(t, 10*time.Minute, LibraryStatsTTL)
	assert.Equal(t, 30*time.Second, SearchResultsTTL)
	assert.Equal(t, 24*time.Hour, ImageMetaTTL)
	assert.Equal(t, 1*time.Minute, ContinueWatchingTTL)
	assert.Equal(t, 2*time.Minute, RecentlyAddedTTL)
	assert.Equal(t, 5*time.Minute, TopRatedTTL)
}
