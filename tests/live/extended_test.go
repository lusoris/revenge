//go:build live

// Package live — extended live tests covering endpoints not in smoke_test.go.
// Requires a running stack: make docker-local
// Run with: go test -tags=live -v -count=1 ./tests/live/...
//
// Metadata endpoints talk to TMDb — tests are lenient:
//   - 200 = TMDb configured and reachable (ideal)
//   - 503/500/502 = TMDb not configured or rate-limited (acceptable)
//   - 401/403 = auth bug (fail)
//   - 429 = rate-limited (skip remaining metadata tests)
package live

import (
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Rate-limit protection for external APIs (TMDb)
// =============================================================================

var metadataRateLimited atomic.Bool

// metadataGet does a GET to a metadata endpoint. If the server returns 429
// it sets a global flag to skip further metadata calls in this test run.
func metadataGet(t *testing.T, path, token string) (int, map[string]interface{}) {
	t.Helper()
	if metadataRateLimited.Load() {
		t.Skip("skipping: metadata API is rate-limited (429 received earlier)")
	}
	// Small delay between metadata calls to avoid TMDb rate limits
	time.Sleep(250 * time.Millisecond)

	status, body := doJSON(t, "GET", path, token, nil)
	if status == 429 {
		metadataRateLimited.Store(true)
		t.Skip("metadata API rate-limited (429)")
	}
	return status, body
}

// assertMetadataOK asserts the response is not an auth error.
// Metadata endpoints may return 200 (configured), 503/500 (not configured), or 404.
func assertMetadataOK(t *testing.T, status int, path string) {
	t.Helper()
	assert.True(t, status != 401 && status != 403,
		"metadata %s should not return auth error (got %d)", path, status)
}

// =============================================================================
// 1. Metadata Sub-resources — Movie
// TMDb Movie 550 = Fight Club, Collection 10 = Star Wars
// =============================================================================

func TestLive_MetadataMovieSubresources(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	subresources := []string{
		"/api/v1/metadata/movie/550/credits",
		"/api/v1/metadata/movie/550/images",
		"/api/v1/metadata/movie/550/external-ids",
		"/api/v1/metadata/movie/550/recommendations",
		"/api/v1/metadata/movie/550/similar",
	}

	for _, path := range subresources {
		t.Run(path, func(t *testing.T) {
			status, _ := metadataGet(t, path, tok)
			assertMetadataOK(t, status, path)
		})
	}
}

// =============================================================================
// 2. Metadata Sub-resources — TV Show
// TMDb TV 1396 = Breaking Bad
// =============================================================================

func TestLive_MetadataTVSubresources(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	subresources := []string{
		"/api/v1/metadata/tv/1396/credits",
		"/api/v1/metadata/tv/1396/images",
		"/api/v1/metadata/tv/1396/external-ids",
		"/api/v1/metadata/tv/1396/content-ratings",
		// Season-level
		"/api/v1/metadata/tv/1396/season/1/images",
		// Episode-level
		"/api/v1/metadata/tv/1396/season/1/episode/1/images",
	}

	for _, path := range subresources {
		t.Run(path, func(t *testing.T) {
			status, _ := metadataGet(t, path, tok)
			assertMetadataOK(t, status, path)
		})
	}
}

// =============================================================================
// 3. Metadata — Person endpoints
// TMDb Person 287 = Brad Pitt
// =============================================================================

func TestLive_MetadataPersonEndpoints(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("person_details", func(t *testing.T) {
		status, _ := metadataGet(t, "/api/v1/metadata/person/287", tok)
		assertMetadataOK(t, status, "person/287")
	})

	t.Run("person_credits", func(t *testing.T) {
		status, _ := metadataGet(t, "/api/v1/metadata/person/287/credits", tok)
		assertMetadataOK(t, status, "person/287/credits")
	})

	t.Run("person_images", func(t *testing.T) {
		status, _ := metadataGet(t, "/api/v1/metadata/person/287/images", tok)
		assertMetadataOK(t, status, "person/287/images")
	})
}

// =============================================================================
// 4. Metadata — Search person, providers
// =============================================================================

func TestLive_MetadataSearchAndProviders(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("search_person", func(t *testing.T) {
		status, _ := metadataGet(t, "/api/v1/metadata/search/person?q=brad+pitt", tok)
		assertMetadataOK(t, status, "search/person")
	})

	t.Run("metadata_providers", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/metadata/providers", tok, nil)
		assert.Equal(t, 200, status, "metadata providers endpoint should return 200")
	})
}

// =============================================================================
// 5. Admin Users (list + delete paths)
// =============================================================================

func TestLive_AdminUsersEndpoints(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	t.Run("list_users", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/admin/users", tok, nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.GreaterOrEqual(t, len(items), 1, "should have at least the admin user")
		}
	})

	// Create a throwaway user and try to soft-delete them
	t.Run("delete_user", func(t *testing.T) {
		victim := registerAndLogin(t)
		status, _ := doJSON(t, "DELETE", "/api/v1/admin/users/"+victim.userID, tok, nil)
		assert.True(t, status == 204 || status == 200,
			"admin delete user should succeed (got %d)", status)
	})

	// Non-admin cannot list users
	t.Run("non_admin_denied", func(t *testing.T) {
		user := registerAndLogin(t)
		status, _ := doJSON(t, "GET", "/api/v1/admin/users", user.accessToken, nil)
		assert.True(t, status == 401 || status == 403,
			"non-admin should be denied (got %d)", status)
	})
}

// =============================================================================
// 6. OIDC Provider Admin — enable/disable/default
// =============================================================================

func TestLive_OIDCProviderManagement(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	// Create a test provider first
	var providerID string
	t.Run("create_provider", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/admin/oidc/providers", tok, map[string]interface{}{
			"name":          "test-ext-provider",
			"display_name":  "Test Ext Provider",
			"client_id":     "test-client",
			"client_secret": "test-secret",
			"issuer_url":    "https://accounts.google.com",
		})
		if status == 409 {
			t.Skip("provider already exists (previous test run)")
		}
		require.Equal(t, 201, status, "create OIDC provider should succeed")
		if id, ok := body["id"].(string); ok {
			providerID = id
		}
		require.NotEmpty(t, providerID)
	})

	if providerID == "" {
		t.Skip("no provider created, cannot test management endpoints")
	}

	t.Run("disable_provider", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/admin/oidc/providers/"+providerID+"/disable", tok, nil)
		assert.True(t, status == 200 || status == 204,
			"disable should succeed (got %d)", status)
	})

	t.Run("enable_provider", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/admin/oidc/providers/"+providerID+"/enable", tok, nil)
		assert.True(t, status == 200 || status == 204,
			"enable should succeed (got %d)", status)
	})

	t.Run("set_default_provider", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/admin/oidc/providers/"+providerID+"/default", tok, nil)
		assert.True(t, status == 200 || status == 204,
			"set default should succeed (got %d)", status)
	})

	// Cleanup: delete the test provider
	t.Run("delete_provider", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/admin/oidc/providers/"+providerID, tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200,
			"delete provider should succeed (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 7. Admin Activity — user-scoped and resource-scoped
// =============================================================================

func TestLive_AdminActivitySubresources(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	t.Run("activity_by_user", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/admin/activity/users/"+admin.userID, tok, nil)
		assert.Equal(t, 200, status)
	})

	t.Run("activity_by_resource", func(t *testing.T) {
		fakeResource := "00000000-0000-0000-0000-000000000001"
		status, _ := doJSON(t, "GET", "/api/v1/admin/activity/resources/user/"+fakeResource, tok, nil)
		assert.Equal(t, 200, status, "resource activity should return 200 (even if empty)")
	})
}

// =============================================================================
// 8. Library Permissions
// =============================================================================

func TestLive_LibraryPermissions(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	// First, find or create a library
	var libraryID string
	t.Run("find_library", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/libraries", tok, nil)
		assert.Equal(t, 200, status)
		if libs, ok := body["libraries"].([]interface{}); ok && len(libs) > 0 {
			if item, ok := libs[0].(map[string]interface{}); ok {
				libraryID, _ = item["id"].(string)
			}
		}
		if libraryID == "" {
			// Create a library (fields match CreateLibraryRequest schema)
			createStatus, createBody := doJSON(t, "POST", "/api/v1/libraries", tok, map[string]interface{}{
				"name":              "Test Library",
				"type":              "movie",
				"paths":             []string{"/media/movies"},
				"metadata_provider": "tmdb",
			})
			require.True(t, createStatus == 201 || createStatus == 200,
				"create library should succeed (got %d): %v", createStatus, createBody)
			libraryID, _ = createBody["id"].(string)
			require.NotEmpty(t, libraryID, "library ID should not be empty")
		}
	})

	if libraryID == "" {
		t.Fatal("no library available to test permissions — find_library should have created one")
	}

	t.Run("get_library_details", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/libraries/"+libraryID, tok, nil)
		assert.Equal(t, 200, status)
		assert.NotEmpty(t, body["name"])
	})

	t.Run("list_library_permissions", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/libraries/"+libraryID+"/permissions", tok, nil)
		assert.True(t, status == 200 || status == 404,
			"library permissions should return 200 or 404 (got %d)", status)
	})

	// Grant permission to a user
	t.Run("grant_user_permission", func(t *testing.T) {
		user := registerAndLogin(t)
		status, _ := doJSON(t, "POST",
			"/api/v1/libraries/"+libraryID+"/permissions",
			tok, map[string]interface{}{
				"user_id":    user.userID,
				"permission": "view",
			})
		assert.True(t, status == 200 || status == 201,
			"grant permission should succeed (got %d)", status)
	})

	// Library scans
	t.Run("list_library_scans", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/libraries/"+libraryID+"/scans", tok, nil)
		assert.True(t, status == 200 || status == 404,
			"scans list should return 200 or 404 (got %d)", status)
	})

	// Trigger scan (won't actually find files but should accept the request)
	t.Run("trigger_scan", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/libraries/"+libraryID+"/scan", tok, map[string]interface{}{
			"scan_type": "full",
		})
		assert.True(t, status == 200 || status == 202 || status == 204 || status == 404,
			"trigger scan should be accepted (got %d)", status)
	})
}

// =============================================================================
// 9. Genres endpoint
// =============================================================================

func TestLive_GenresEndpoint(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("list_genres", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/genres", tok)
		assert.Equal(t, 200, status)
	})
}

// =============================================================================
// 10. Collections CRUD (with fake IDs)
// =============================================================================

func TestLive_CollectionsCRUD(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken
	fakeID := "00000000-0000-0000-0000-000000000001"

	t.Run("get_collection_not_found", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/collections/"+fakeID, tok, nil)
		assert.Equal(t, 404, status)
	})

	t.Run("collection_movies_not_found", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/collections/"+fakeID+"/movies", tok, nil)
		assert.True(t, status == 404 || status == 200,
			"collection movies should be 404 or 200 empty (got %d)", status)
	})
}

// =============================================================================
// 11. Search — TV shows, multi-search, reindex
// =============================================================================

func TestLive_SearchTVShows(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("search_tvshows", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/search/tvshows?q=test", tok, nil)
		assert.True(t, status == 200 || status == 503,
			"tvshow search should return 200 or 503 if search disabled (got %d)", status)
	})

	t.Run("search_tvshows_autocomplete", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/search/tvshows/autocomplete?q=test", tok, nil)
		assert.True(t, status == 200 || status == 503,
			"tvshow autocomplete should return 200 or 503 (got %d)", status)
	})

	t.Run("multi_search", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/search/multi?q=test", tok, nil)
		assert.True(t, status == 200 || status == 503,
			"multi search should return 200 or 503 (got %d)", status)
	})

	// Reindex — check it doesn't error (endpoint may or may not require admin)
	t.Run("reindex_tvshows_requires_admin", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/search/tvshows/reindex", tok, nil)
		assert.True(t, status == 200 || status == 202 || status == 401 || status == 403 || status == 503,
			"reindex should succeed or be denied (got %d)", status)
	})

	t.Run("reindex_tvshows_admin", func(t *testing.T) {
		admin := ensureAdmin(t)
		status, _ := doJSON(t, "POST", "/api/v1/search/tvshows/reindex", admin.accessToken, nil)
		assert.True(t, status == 200 || status == 202 || status == 204 || status == 503,
			"admin reindex should succeed or search disabled (got %d)", status)
	})
}

// =============================================================================
// 12. Settings — per-key user settings
// =============================================================================

func TestLive_UserSettingsPerKey(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	// Set a setting
	t.Run("set_theme", func(t *testing.T) {
		status, _ := doJSON(t, "PUT", "/api/v1/settings/user/theme", tok,
			map[string]interface{}{"value": "dark"})
		assert.True(t, status == 200 || status == 201,
			"set theme should succeed (got %d)", status)
	})

	// Get the setting back
	t.Run("get_theme", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/settings/user/theme", tok, nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "dark", body["value"])
	})

	// Delete the setting
	t.Run("delete_theme", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/settings/user/theme", tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200,
			"delete setting should succeed (got %d)", resp.StatusCode)
	})

	// Setting should be gone (404)
	t.Run("theme_deleted", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/settings/user/theme", tok, nil)
		assert.True(t, status == 404 || status == 200,
			"deleted setting should return 404 or default (got %d)", status)
	})

	// Server settings (admin only)
	t.Run("server_settings_admin", func(t *testing.T) {
		admin := ensureAdmin(t)
		status, _ := doJSON(t, "GET", "/api/v1/settings/server/app_name", admin.accessToken, nil)
		assert.True(t, status == 200 || status == 404,
			"server setting should return 200 or 404 (got %d)", status)
	})

	t.Run("server_settings_non_admin", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/settings/server/app_name", tok, nil)
		assert.True(t, status == 401 || status == 403 || status == 200,
			"non-admin server setting access (got %d)", status)
	})
}

// =============================================================================
// 13. Sessions per-ID
// =============================================================================

func TestLive_SessionPerID(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	// List sessions to get the current session ID
	var sessionID string
	t.Run("list_sessions", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/sessions", tok, nil)
		require.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok && len(items) > 0 {
			if item, ok := items[0].(map[string]interface{}); ok {
				sessionID, _ = item["id"].(string)
			}
		}
	})

	if sessionID == "" {
		t.Skip("no session found")
	}

	t.Run("get_session_by_id", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/sessions/"+sessionID, tok, nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, sessionID, body["id"])
	})

	// Revoke a different session (create one first)
	t.Run("revoke_session", func(t *testing.T) {
		// Login again to create a second session
		loginUser(t, creds)
		newTok := creds.accessToken

		// List sessions from new token
		status, body := doJSON(t, "GET", "/api/v1/sessions", newTok, nil)
		require.Equal(t, 200, status)
		items, ok := body["items"].([]interface{})
		if !ok || len(items) < 2 {
			t.Skip("need at least 2 sessions to test revocation")
		}

		// Find a session that is NOT the current one and revoke it
		for _, item := range items {
			s, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			sid, _ := s["id"].(string)
			if sid != "" && sid != sessionID {
				resp := doRequest(t, "DELETE", "/api/v1/sessions/"+sid, newTok, nil)
				resp.Body.Close()
				assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200,
					"revoke session should succeed (got %d)", resp.StatusCode)
				return
			}
		}
		t.Skip("could not find a secondary session to revoke")
	})
}

// =============================================================================
// 14. RBAC Granular — role CRUD, user role assignment
// =============================================================================

func TestLive_RBACUserRoles(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	user := registerAndLogin(t)

	// Assign role to user
	t.Run("assign_role", func(t *testing.T) {
		status, _ := doJSON(t, "POST",
			"/api/v1/rbac/users/"+user.userID+"/roles",
			tok, map[string]interface{}{"role": "moderator"})
		assert.True(t, status == 200 || status == 204 || status == 201,
			"assign role should succeed (got %d)", status)
	})

	// List user's roles
	t.Run("list_user_roles", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/rbac/users/"+user.userID+"/roles", tok)
		assert.Equal(t, 200, status)
	})

	// Remove role from user
	t.Run("remove_role", func(t *testing.T) {
		resp := doRequest(t, "DELETE",
			"/api/v1/rbac/users/"+user.userID+"/roles/moderator",
			tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200,
			"remove role should succeed (got %d)", resp.StatusCode)
	})

	// RBAC role by name
	t.Run("get_role_by_name", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/rbac/roles/admin", tok, nil)
		assert.True(t, status == 200 || status == 404 || status == 403,
			"get role should return 200, 404, or 403 (got %d)", status)
	})
}

// =============================================================================
// 15. TV Show Episode — bulk watched
// =============================================================================

func TestLive_TVShowBulkWatched(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	fakeEp1 := "00000000-0000-0000-0000-000000000001"
	fakeEp2 := "00000000-0000-0000-0000-000000000002"

	t.Run("bulk_mark_watched_not_found", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/tvshows/episodes/bulk-watched", tok,
			map[string]interface{}{
				"episode_ids": []string{fakeEp1, fakeEp2},
			})
		// Episodes don't exist so this should 404 or maybe 400
		assert.True(t, status == 404 || status == 400 || status == 200 || status == 204,
			"bulk watched with fake IDs (got %d)", status)
	})
}

// =============================================================================
// 16. OIDC User Endpoints (link/unlink)
// =============================================================================

func TestLive_OIDCUserLink(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("unlink_nonexistent_provider", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/users/me/oidc/nonexistent", tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 400,
			"unlink non-existent OIDC should 404 or 400 (got %d)", resp.StatusCode)
	})

	// The /link endpoint initiates an OAuth flow — returns auth URL or error
	t.Run("link_provider_init", func(t *testing.T) {
		// POST to initiate link, expect 200 with auth URL or 404 if provider doesn't exist
		status, _ := doJSON(t, "POST", "/api/v1/users/me/oidc/google/link", tok, nil)
		assert.True(t, status == 200 || status == 404 || status == 400 || status == 503,
			"OIDC link init should return auth URL or error (got %d)", status)
	})
}

// =============================================================================
// 17. User by ID (public user endpoint)
// =============================================================================

func TestLive_UserByID(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("get_own_user", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/users/"+creds.userID, tok, nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, creds.username, body["username"])
	})

	t.Run("get_nonexistent_user", func(t *testing.T) {
		fakeID := "00000000-0000-0000-0000-000000000099"
		status, _ := doJSON(t, "GET", "/api/v1/users/"+fakeID, tok, nil)
		assert.Equal(t, 404, status)
	})
}

// =============================================================================
// 18. OIDC Auth Flow (redirect-based)
// =============================================================================

func TestLive_OIDCAuthFlow(t *testing.T) {
	noRedirectClient := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	t.Run("auth_nonexistent_provider", func(t *testing.T) {
		resp, err := noRedirectClient.Get(baseURL + "/api/v1/oidc/auth/nonexistent")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 400,
			"OIDC auth for non-existent provider (got %d)", resp.StatusCode)
	})

	t.Run("callback_no_code", func(t *testing.T) {
		resp, err := noRedirectClient.Get(baseURL + "/api/v1/oidc/callback/nonexistent")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 400 || resp.StatusCode == 404,
			"OIDC callback without code (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 19. Playback Session by ID
// =============================================================================

func TestLive_PlaybackSessionByID(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	fakeSession := "00000000-0000-0000-0000-000000000001"

	t.Run("get_nonexistent_session", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/playback/sessions/"+fakeSession, tok, nil)
		assert.True(t, status == 404 || status == 200,
			"non-existent playback session (got %d)", status)
	})

	t.Run("delete_nonexistent_session", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/playback/sessions/"+fakeSession, tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 404,
			"delete non-existent playback session (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 20. Images endpoint (non-existent image)
// =============================================================================

func TestLive_ImagesEndpoint(t *testing.T) {
	t.Run("nonexistent_image", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/api/v1/images/poster/w500/nonexistent.jpg")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 400,
			"non-existent image should 404 (got %d)", resp.StatusCode)
	})

	// Valid image type/size combinations
	t.Run("valid_path_no_image", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/api/v1/images/backdrop/w1280/doesnotexist.jpg")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 400,
			"valid path but no image (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 21. Movies search endpoint (query param variant)
// =============================================================================

func TestLive_MovieSearchEndpoint(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("search_movies_empty", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/movies/search?query=nonexistent_xyzzy", tok, nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.Empty(t, items, "search for nonsense should return empty")
		}
	})
}

// =============================================================================
// 22. TV Show search endpoint
// =============================================================================

func TestLive_TVShowSearchEndpoint(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("search_tvshows_empty", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/tvshows/search?query=nonexistent_xyzzy", tok, nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.Empty(t, items)
		}
	})
}

// =============================================================================
// 23. API Key by ID
// =============================================================================

func TestLive_APIKeyByID(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	// Create a key
	var keyID string
	var rawKey string
	t.Run("create_key", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/apikeys", tok,
			map[string]interface{}{"name": "live-test-key", "scopes": []string{"read", "write"}})
		require.Equal(t, 201, status)
		keyID, _ = body["id"].(string)
		rawKey, _ = body["api_key"].(string)
		require.NotEmpty(t, keyID)
		require.NotEmpty(t, rawKey)
	})

	// Authenticate with the API key
	t.Run("auth_with_api_key", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no key created")
		}
		req, err := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		require.NoError(t, err)
		req.Header.Set("X-API-Key", rawKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body) //nolint:errcheck
		assert.Equal(t, creds.username, body["username"])
	})

	// Get key details (without revealing the raw key)
	t.Run("get_key_by_id", func(t *testing.T) {
		if keyID == "" {
			t.Skip("no key")
		}
		status, body := doJSON(t, "GET", "/api/v1/apikeys/"+keyID, tok, nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "live-test-key", body["name"])
	})

	// Delete key
	t.Run("delete_key", func(t *testing.T) {
		if keyID == "" {
			t.Skip("no key to delete")
		}
		resp := doRequest(t, "DELETE", "/api/v1/apikeys/"+keyID, tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200)
	})

	// Deleted key should no longer authenticate
	t.Run("deleted_key_rejected", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no key")
		}
		req, err := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		require.NoError(t, err)
		req.Header.Set("X-API-Key", rawKey)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.True(t, resp.StatusCode == 401 || resp.StatusCode == 403,
			"deleted API key should be rejected (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 24. WebAuthn credential by ID (non-existent)
// =============================================================================

func TestLive_WebAuthnCredentialByID(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("delete_nonexistent_credential", func(t *testing.T) {
		resp := doRequest(t, "DELETE",
			"/api/v1/mfa/webauthn/credentials/nonexistent-cred-id", tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 400 || resp.StatusCode == 204,
			"delete non-existent credential (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 25. Content Detail Endpoint Validation (response structure)
// Verify that detail endpoints return proper JSON structure when accessed
// with a non-existent UUID — ensures handlers are wired correctly.
// =============================================================================

func TestLive_ContentDetailStructure(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken
	fakeID := "00000000-0000-0000-0000-000000000099"

	// Movie detail
	t.Run("movie_detail_404", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/movies/"+fakeID, tok, nil)
		assert.Equal(t, 404, status, "non-existent movie should return 404")
		assert.NotEmpty(t, body["message"], "404 should have error message")
	})

	// TV Show detail
	t.Run("tvshow_detail_404", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/tvshows/"+fakeID, tok, nil)
		assert.Equal(t, 404, status, "non-existent tvshow should return 404")
		assert.NotEmpty(t, body["message"], "404 should have error message")
	})
}

// =============================================================================
// 26. WebAuthn Login Begin (no credentials registered)
// =============================================================================

func TestLive_WebAuthnLoginBegin(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	t.Run("login_begin_no_credentials", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/mfa/webauthn/login/begin", tok, nil)
		// Should fail — user has no WebAuthn credentials registered
		assert.True(t, status == 400 || status == 404 || status == 500,
			"login begin without credentials (got %d)", status)
	})
}

// =============================================================================
// 27. Response structure validation — list endpoints return total + items
// =============================================================================

func TestLive_ListResponseStructure(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	// Each endpoint has its own array field name but all have 'total'
	listEndpoints := map[string]string{
		"/api/v1/movies":         "items",
		"/api/v1/tvshows":        "items",
		"/api/v1/admin/users":    "users",
		"/api/v1/admin/activity": "entries",
	}

	for path, arrayField := range listEndpoints {
		t.Run(path, func(t *testing.T) {
			resp := doRequest(t, "GET", path, tok, nil)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var body map[string]interface{}
			err = json.Unmarshal(raw, &body)
			require.NoError(t, err, "should be valid JSON: %s", string(raw[:min(len(raw), 200)]))

			_, hasArrayField := body[arrayField]
			_, hasTotal := body["total"]
			assert.True(t, hasArrayField, "%s should have '%s' field", path, arrayField)
			assert.True(t, hasTotal, "%s should have 'total' field", path)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
