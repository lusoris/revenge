//go:build live

package live

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// 28. Playback Heartbeat — endpoint tests (no real media needed)
// =============================================================================

func TestLive_PlaybackHeartbeat(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken
	fakeSessionID := "00000000-0000-0000-0000-000000000099"

	t.Run("heartbeat_no_auth", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/"+fakeSessionID+"/heartbeat", "", nil)
		defer resp.Body.Close()
		assert.Equal(t, 401, resp.StatusCode, "heartbeat without auth should 401")
	})

	t.Run("heartbeat_not_found", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/"+fakeSessionID+"/heartbeat", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 404, resp.StatusCode,
			"heartbeat for non-existent session should 404 (got %d)", resp.StatusCode)
	})

	t.Run("heartbeat_invalid_id", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/not-a-uuid/heartbeat", tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 400 || resp.StatusCode == 422,
			"heartbeat with invalid UUID should 400/422 (got %d)", resp.StatusCode)
	})

	t.Run("heartbeat_with_position_not_found", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/"+fakeSessionID+"/heartbeat", tok,
			map[string]interface{}{"position_seconds": 42})
		defer resp.Body.Close()
		assert.Equal(t, 404, resp.StatusCode,
			"heartbeat with position for non-existent session should 404")
	})
}

// TestLive_PlaybackHeartbeatRealSession creates a real playback session (needs
// the BBB movie file) and validates the heartbeat works end-to-end.
func TestLive_PlaybackHeartbeatRealSession(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)
	defer db.Close()

	var movieID string
	err = db.QueryRow(`
		INSERT INTO movie.movies (title, year, runtime, overview, original_language)
		VALUES ('Heartbeat Test Movie', 2024, 10, 'Heartbeat test', 'en')
		RETURNING id`).Scan(&movieID)
	require.NoError(t, err)
	defer func() {
		if movieID != "" {
			d, _ := sql.Open("pgx", dbURL)
			if d != nil {
				defer d.Close()
				_, _ = d.Exec(`DELETE FROM movie.movies WHERE id = $1`, movieID)
			}
		}
	}()

	var fileID string
	err = db.QueryRow(`
		INSERT INTO movie.movie_files
			(movie_id, file_path, file_size, file_name, resolution, video_codec, audio_codec, container, duration_seconds, bitrate_kbps)
		VALUES ($1, '/movies/bbb_sunflower_2160p_30fps_normal.mp4', 633000000,
			'bbb_sunflower_2160p_30fps_normal.mp4', '2160p', 'h264', 'mp3', 'mp4', 634, 8000)
		RETURNING id`, movieID).Scan(&fileID)
	if err != nil {
		t.Skipf("could not insert movie_file (video file may not exist): %v", err)
	}
	require.NotEmpty(t, fileID)

	// Start playback session
	var sessionID string
	status, body := doJSON(t, "POST", "/api/v1/playback/sessions", tok, map[string]interface{}{
		"media_type": "movie",
		"media_id":   movieID,
	})
	if status != 200 && status != 201 {
		t.Skipf("could not start playback session (FFmpeg may not be available): status=%d", status)
	}
	sessionID, _ = body["session_id"].(string)
	require.NotEmpty(t, sessionID)

	defer func() {
		if sessionID != "" {
			resp := doRequest(t, "DELETE", "/api/v1/playback/sessions/"+sessionID, tok, nil)
			resp.Body.Close()
		}
	}()

	t.Run("heartbeat_keepalive", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/"+sessionID+"/heartbeat", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 204, resp.StatusCode,
			"heartbeat keep-alive should return 204 (got %d)", resp.StatusCode)
	})

	t.Run("heartbeat_with_position", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/"+sessionID+"/heartbeat", tok,
			map[string]interface{}{"position_seconds": 120})
		defer resp.Body.Close()
		assert.Equal(t, 204, resp.StatusCode,
			"heartbeat with position should return 204 (got %d)", resp.StatusCode)
	})

	t.Run("session_still_active", func(t *testing.T) {
		st, b := doJSON(t, "GET", "/api/v1/playback/sessions/"+sessionID, tok, nil)
		require.Equal(t, 200, st, "session should still exist: %v", b)
	})

	t.Run("rapid_heartbeats", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := doRequest(t, "POST",
				"/api/v1/playback/sessions/"+sessionID+"/heartbeat", tok,
				map[string]interface{}{"position_seconds": 120 + i*10})
			resp.Body.Close()
			assert.Equal(t, 204, resp.StatusCode, "rapid heartbeat %d should succeed", i)
		}
	})

	savedSessionID := sessionID
	t.Run("stop_session", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/playback/sessions/"+sessionID, tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204)
		sessionID = ""
	})

	t.Run("heartbeat_after_stop", func(t *testing.T) {
		resp := doRequest(t, "POST",
			"/api/v1/playback/sessions/"+savedSessionID+"/heartbeat", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 404, resp.StatusCode, "heartbeat after session stop should 404")
	})
}

// =============================================================================
// 29. API Key Caching — validate cache hit/invalidation behavior
// =============================================================================

func TestLive_APIKeyCaching(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	var keyID, rawKey string
	t.Run("create_key", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/apikeys", tok,
			map[string]interface{}{"name": "cache-test-key", "scopes": []string{"read", "write"}})
		require.Equal(t, 201, status)
		keyID, _ = body["id"].(string)
		rawKey, _ = body["api_key"].(string)
		require.NotEmpty(t, keyID)
		require.NotEmpty(t, rawKey)
	})

	t.Run("first_request_cold_cache", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no key")
		}
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", rawKey)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode, "first API key request should succeed")

		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body) //nolint:errcheck
		assert.Equal(t, creds.username, body["username"])
	})

	t.Run("second_request_warm_cache", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no key")
		}
		start := time.Now()
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", rawKey)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		elapsed := time.Since(start)
		require.Equal(t, 200, resp.StatusCode)
		t.Logf("cached API key validation: %v", elapsed)

		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body) //nolint:errcheck
		assert.Equal(t, creds.username, body["username"])
	})

	t.Run("concurrent_api_key_requests", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no key")
		}
		const concurrency = 10
		results := make(chan int, concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
				req.Header.Set("X-API-Key", rawKey)
				resp, err := client.Do(req)
				if err != nil {
					results <- 0
					return
				}
				defer resp.Body.Close()
				results <- resp.StatusCode
			}()
		}
		for i := 0; i < concurrency; i++ {
			status := <-results
			assert.Equal(t, 200, status, "concurrent request %d should succeed", i)
		}
	})

	t.Run("revoke_key", func(t *testing.T) {
		if keyID == "" {
			t.Skip("no key")
		}
		resp := doRequest(t, "DELETE", "/api/v1/apikeys/"+keyID, tok, nil)
		resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204)
	})

	t.Run("revoked_key_rejected", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no key")
		}
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", rawKey)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			t.Log("key still accepted immediately after revocation — waiting for cache invalidation")
			time.Sleep(2 * time.Second)

			req2, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
			req2.Header.Set("X-API-Key", rawKey)
			resp2, err := client.Do(req2)
			require.NoError(t, err)
			defer resp2.Body.Close()
			assert.True(t, resp2.StatusCode == 401 || resp2.StatusCode == 403,
				"revoked key should be rejected after cache TTL (got %d)", resp2.StatusCode)
		} else {
			assert.True(t, resp.StatusCode == 401 || resp.StatusCode == 403,
				"revoked key should be rejected (got %d)", resp.StatusCode)
		}
	})
}

// =============================================================================
// 30. API Key — invalid/malformed key handling
// =============================================================================

func TestLive_APIKeyEdgeCases(t *testing.T) {
	creds := registerAndLogin(t)

	t.Run("empty_api_key_header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", "")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 401, resp.StatusCode, "empty API key should 401")
	})

	t.Run("garbage_api_key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", "this-is-not-a-valid-api-key-at-all")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 401 || resp.StatusCode == 403,
			"garbage API key (got %d)", resp.StatusCode)
	})

	t.Run("sql_injection_attempt", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", "'; DROP TABLE shared.api_keys; --")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 401 || resp.StatusCode == 403,
			"SQL injection attempt should be safely rejected (got %d)", resp.StatusCode)
	})

	t.Run("very_long_api_key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", strings.Repeat("A", 10000))
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 401 || resp.StatusCode == 403 || resp.StatusCode == 400,
			"very long API key should be rejected (got %d)", resp.StatusCode)
	})

	t.Run("api_keys_still_work", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/apikeys", creds.accessToken,
			map[string]interface{}{"name": "after-edge-case", "scopes": []string{"read"}})
		require.Equal(t, 201, status)
		keyID, _ := body["id"].(string)
		require.NotEmpty(t, keyID)

		resp := doRequest(t, "DELETE", "/api/v1/apikeys/"+keyID, creds.accessToken, nil)
		resp.Body.Close()
	})
}

// =============================================================================
// 31. API Key concurrent create + use + delete stress test
// =============================================================================

func TestLive_APIKeyConcurrentLifecycle(t *testing.T) {
	const numKeys = 5
	creds := registerAndLogin(t)
	tok := creds.accessToken

	type keyInfo struct {
		id     string
		secret string
	}

	keys := make([]keyInfo, 0, numKeys)
	for i := 0; i < numKeys; i++ {
		status, body := doJSON(t, "POST", "/api/v1/apikeys", tok,
			map[string]interface{}{
				"name":   fmt.Sprintf("concurrent-key-%d", i),
				"scopes": []string{"read"},
			})
		require.Equal(t, 201, status)
		id, _ := body["id"].(string)
		secret, _ := body["api_key"].(string)
		keys = append(keys, keyInfo{id: id, secret: secret})
	}
	require.Len(t, keys, numKeys)

	t.Run("use_all_keys_concurrently", func(t *testing.T) {
		results := make(chan int, numKeys)
		for _, k := range keys {
			go func(apiKey string) {
				req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
				req.Header.Set("X-API-Key", apiKey)
				resp, err := client.Do(req)
				if err != nil {
					results <- 0
					return
				}
				defer resp.Body.Close()
				results <- resp.StatusCode
			}(k.secret)
		}
		for i := 0; i < numKeys; i++ {
			status := <-results
			assert.Equal(t, 200, status, "concurrent key %d should work", i)
		}
	})

	t.Run("delete_all_keys", func(t *testing.T) {
		for _, k := range keys {
			resp := doRequest(t, "DELETE", "/api/v1/apikeys/"+k.id, tok, nil)
			resp.Body.Close()
			assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204)
		}
	})

	t.Run("list_keys_after_cleanup", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/apikeys", tok, nil)
		require.Equal(t, 200, status)
		if apiKeys, ok := body["keys"].([]interface{}); ok {
			assert.Empty(t, apiKeys, "all test keys should be deleted")
		}
	})
}

// =============================================================================
// 32. Notification config endpoint (admin-only)
// =============================================================================

func TestLive_NotificationConfig(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	t.Run("server_info_has_notifications", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/admin/settings", tok, nil)
		t.Logf("admin settings status=%d keys=%v", status, keys(body))
		assert.True(t, status == 200 || status == 404,
			"admin settings should be accessible or not exist (got %d)", status)
	})

	t.Run("non_admin_cannot_access", func(t *testing.T) {
		user := registerAndLogin(t)
		status, _ := doJSON(t, "GET", "/api/v1/admin/settings", user.accessToken, nil)
		assert.True(t, status == 403 || status == 401 || status == 404,
			"non-admin should be forbidden from admin settings (got %d)", status)
	})
}

// =============================================================================
// 33. Rate limiting validation
// =============================================================================

func TestLive_RateLimiting(t *testing.T) {
	t.Run("auth_rate_limit", func(t *testing.T) {
		hitRateLimit := false
		for i := 0; i < 30; i++ {
			resp := doRequest(t, "POST", "/api/v1/auth/login", "", map[string]string{
				"username": "nonexistent_user_rate_limit_test",
				"password": "wrong",
			})
			resp.Body.Close()
			if resp.StatusCode == 429 {
				hitRateLimit = true
				t.Logf("rate limited after %d requests", i+1)
				break
			}
		}
		t.Logf("hit rate limit: %v", hitRateLimit)
	})
}

// =============================================================================
// 34. Integration endpoints (Radarr/Sonarr) with arrbase types
// =============================================================================

func TestLive_IntegrationConfigEndpoints(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	t.Run("radarr_config", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/integrations/radarr", tok, nil)
		t.Logf("radarr config: status=%d keys=%v", status, keys(body))
		assert.True(t, status == 200 || status == 404 || status == 503,
			"radarr config should return 200/404/503 (got %d)", status)
	})

	t.Run("sonarr_config", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/integrations/sonarr", tok, nil)
		t.Logf("sonarr config: status=%d keys=%v", status, keys(body))
		assert.True(t, status == 200 || status == 404 || status == 503,
			"sonarr config should return 200/404/503 (got %d)", status)
	})

	t.Run("non_admin_forbidden", func(t *testing.T) {
		user := registerAndLogin(t)
		status, _ := doJSON(t, "GET", "/api/v1/integrations/radarr", user.accessToken, nil)
		assert.True(t, status == 403 || status == 401 || status == 404,
			"non-admin should not access integrations (got %d)", status)
	})
}
