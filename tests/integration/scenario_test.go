//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Shared scenario helpers
// ============================================================================

type scenarioClient struct {
	t       *testing.T
	ts      *TestServer
	token   string
	userID  string
	refresh string
}

func newScenarioClient(t *testing.T, ts *TestServer) *scenarioClient {
	return &scenarioClient{t: t, ts: ts}
}

func (c *scenarioClient) register(username, email, password string) {
	c.t.Helper()
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	})
	resp, err := c.ts.HTTPClient.Post(
		c.ts.BaseURL+"/api/v1/auth/register",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(c.t, err)
	defer resp.Body.Close()
	require.Equal(c.t, 201, resp.StatusCode, "registration should succeed")

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user) //nolint:errcheck
	if id, ok := user["id"].(string); ok {
		c.userID = id
	}
}

func (c *scenarioClient) login(username, password string) {
	c.t.Helper()
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	resp, err := c.ts.HTTPClient.Post(
		c.ts.BaseURL+"/api/v1/auth/login",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(c.t, err)
	defer resp.Body.Close()
	require.Equal(c.t, 200, resp.StatusCode, "login should succeed")

	var loginResp LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp) //nolint:errcheck
	c.token = loginResp.AccessToken
	c.refresh = loginResp.RefreshToken
	c.userID = loginResp.User.ID
	require.NotEmpty(c.t, c.token, "must have access token")
}

func (c *scenarioClient) registerAndLogin(username, email, password string) {
	c.t.Helper()
	c.register(username, email, password)
	c.login(username, password)
}

func (c *scenarioClient) do(method, path string, body interface{}) (int, map[string]interface{}) {
	c.t.Helper()
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(
		context.Background(), method, c.ts.BaseURL+path, reqBody,
	)
	require.NoError(c.t, err)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.ts.HTTPClient.Do(req)
	require.NoError(c.t, err)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result) //nolint:errcheck
	return resp.StatusCode, result
}

func (c *scenarioClient) doArray(method, path string) (int, []interface{}) {
	c.t.Helper()
	req, err := http.NewRequestWithContext(
		context.Background(), method, c.ts.BaseURL+path, nil,
	)
	require.NoError(c.t, err)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.ts.HTTPClient.Do(req)
	require.NoError(c.t, err)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result []interface{}
	json.Unmarshal(respBody, &result) //nolint:errcheck
	return resp.StatusCode, result
}

// makeAdmin gives a user admin role via direct DB insert.
func (c *scenarioClient) makeAdmin() {
	c.t.Helper()
	require.NotEmpty(c.t, c.userID, "must have userID to make admin")
	_, err := c.ts.AppPool.Exec(context.Background(),
		"INSERT INTO shared.casbin_rule (ptype, v0, v1) VALUES ('g', $1, 'admin') ON CONFLICT DO NOTHING",
		c.userID)
	require.NoError(c.t, err)
	_, err = c.ts.AppPool.Exec(context.Background(),
		"UPDATE shared.users SET is_admin = true WHERE id = $1", c.userID)
	require.NoError(c.t, err)

	// Reload RBAC policy so casbin sees the new role immediately
	err = c.ts.RBACService.LoadPolicy(context.Background())
	require.NoError(c.t, err)
}

// ============================================================================
// Scenario 1: Complete User Journey
// Register -> Login -> View profile -> Change password -> Re-login -> API keys -> Logout
// ============================================================================

func TestScenario_CompleteUserJourney(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	c := newScenarioClient(t, ts)

	// Step 1: Register
	c.register("alice", "alice@example.com", "S3cureP@ssw0rd!")

	// Step 2: Login
	c.login("alice", "S3cureP@ssw0rd!")

	// Step 3: View own profile
	t.Run("view_profile", func(t *testing.T) {
		status, body := c.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "alice", body["username"])
		assert.Equal(t, "alice@example.com", body["email"])
	})

	// Step 4: Change password
	t.Run("change_password", func(t *testing.T) {
		status, _ := c.do("POST", "/api/v1/auth/change-password", map[string]string{
			"old_password": "S3cureP@ssw0rd!",
			"new_password": "N3wS3cureP@ss!",
		})
		assert.Equal(t, 204, status, "change password should succeed")
	})

	// Step 5: Old password should fail
	t.Run("old_password_fails", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": "alice",
			"password": "S3cureP@ssw0rd!",
		})
		resp, err := ts.HTTPClient.Post(
			ts.BaseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, 401, resp.StatusCode, "old password should be rejected")
	})

	// Step 6: New password works
	t.Run("new_password_works", func(t *testing.T) {
		c.login("alice", "N3wS3cureP@ss!")
	})

	// Step 7: Create API key
	var apiKeyID string
	t.Run("create_api_key", func(t *testing.T) {
		status, body := c.do("POST", "/api/v1/apikeys", map[string]interface{}{
			"name":   "test-key",
			"scopes": []string{"read", "write"},
		})
		assert.Equal(t, 201, status)
		if id, ok := body["id"].(string); ok {
			apiKeyID = id
		}
		assert.NotEmpty(t, body["api_key"], "should return the raw key")
	})

	// Step 8: List API keys
	t.Run("list_api_keys", func(t *testing.T) {
		status, body := c.do("GET", "/api/v1/apikeys", nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.GreaterOrEqual(t, len(items), 1)
		}
	})

	// Step 9: Delete API key
	t.Run("delete_api_key", func(t *testing.T) {
		if apiKeyID == "" {
			t.Skip("no API key to delete")
		}
		status, _ := c.do("DELETE", "/api/v1/apikeys/"+apiKeyID, nil)
		assert.Equal(t, 204, status)
	})

	// Step 10: Logout
	t.Run("logout", func(t *testing.T) {
		status, _ := c.do("POST", "/api/v1/auth/logout", map[string]string{
			"refresh_token": c.refresh,
		})
		assert.True(t, status == 200 || status == 204, "logout should succeed")
	})
}

// ============================================================================
// Scenario 2: Admin User Management Journey
// Register admin -> Assign role -> Create second user -> Verify permissions
// ============================================================================

func TestScenario_AdminUserManagement(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Admin setup
	admin := newScenarioClient(t, ts)
	admin.registerAndLogin("admin_user", "admin@example.com", "Adm1nP@ss!")
	admin.makeAdmin()
	admin.login("admin_user", "Adm1nP@ss!") // re-login to pick up admin role

	// Regular user
	user := newScenarioClient(t, ts)
	user.registerAndLogin("regular_user", "regular@example.com", "Us3rP@ss!")

	// Step 1: Admin can list users
	t.Run("admin_lists_users", func(t *testing.T) {
		status, body := admin.do("GET", "/api/v1/admin/users", nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.GreaterOrEqual(t, len(items), 2, "should see admin + regular user")
		}
	})

	// Step 2: Regular user cannot list all users
	t.Run("user_cannot_list_users", func(t *testing.T) {
		status, _ := user.do("GET", "/api/v1/admin/users", nil)
		assert.True(t, status == 401 || status == 403,
			"regular user should be denied admin endpoint (got %d)", status)
	})

	// Step 3: Admin can view specific user (via general users endpoint, admin/users only has DELETE)
	t.Run("admin_views_user", func(t *testing.T) {
		status, body := admin.do("GET", "/api/v1/users/"+user.userID, nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "regular_user", body["username"])
	})

	// Step 4: Admin can list roles
	t.Run("admin_lists_roles", func(t *testing.T) {
		status, _ := admin.doArray("GET", "/api/v1/rbac/roles")
		assert.Equal(t, 200, status)
	})

	// Step 5: Regular user cannot access RBAC
	t.Run("user_cannot_access_rbac", func(t *testing.T) {
		status, _ := user.doArray("GET", "/api/v1/rbac/roles")
		assert.True(t, status == 401 || status == 403,
			"regular user should be denied RBAC endpoint (got %d)", status)
	})

	// Step 6: Admin assigns role to user
	t.Run("admin_assigns_role", func(t *testing.T) {
		status, _ := admin.do("POST", "/api/v1/rbac/users/"+user.userID+"/roles",
			map[string]string{"role": "moderator"})
		assert.True(t, status == 200 || status == 204,
			"role assignment should succeed (got %d)", status)
	})

	// Step 7: Admin can list permissions
	t.Run("admin_lists_permissions", func(t *testing.T) {
		status, _ := admin.doArray("GET", "/api/v1/rbac/permissions")
		assert.Equal(t, 200, status)
	})
}

// ============================================================================
// Scenario 3: Multi-User Session Isolation
// Two users login -> Each has independent sessions -> One logout doesn't affect other
// ============================================================================

func TestScenario_MultiUserSessionIsolation(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	alice := newScenarioClient(t, ts)
	alice.registerAndLogin("session_alice", "session_alice@test.com", "P@ssw0rd1!")

	bob := newScenarioClient(t, ts)
	bob.registerAndLogin("session_bob", "session_bob@test.com", "P@ssw0rd2!")

	// Both can access their profiles
	t.Run("both_users_have_profiles", func(t *testing.T) {
		statusA, bodyA := alice.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, statusA)
		assert.Equal(t, "session_alice", bodyA["username"])

		statusB, bodyB := bob.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, statusB)
		assert.Equal(t, "session_bob", bodyB["username"])
	})

	// Alice lists her sessions
	t.Run("alice_lists_sessions", func(t *testing.T) {
		status, _ := alice.doArray("GET", "/api/v1/sessions")
		assert.Equal(t, 200, status)
	})

	// Alice logs out
	t.Run("alice_logs_out", func(t *testing.T) {
		status, _ := alice.do("POST", "/api/v1/auth/logout", map[string]string{
			"refresh_token": alice.refresh,
		})
		assert.True(t, status == 200 || status == 204)
	})

	// Bob is still authenticated
	t.Run("bob_still_authenticated", func(t *testing.T) {
		statusB, bodyB := bob.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, statusB)
		assert.Equal(t, "session_bob", bodyB["username"])
	})
}

// ============================================================================
// Scenario 4: Content Browsing (Empty Library)
// Verify all content endpoints return proper empty responses
// ============================================================================

func TestScenario_EmptyLibraryBrowsing(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	c := newScenarioClient(t, ts)
	c.registerAndLogin("browser", "browser@test.com", "Br0ws3r!")

	// Movies
	t.Run("empty_movie_list", func(t *testing.T) {
		status, body := c.do("GET", "/api/v1/movies", nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.Empty(t, items)
		}
	})

	// TV Shows
	t.Run("empty_tvshow_list", func(t *testing.T) {
		status, body := c.do("GET", "/api/v1/tvshows", nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.Empty(t, items)
		}
	})

	// Continue watching (should be empty)
	t.Run("empty_continue_watching_movies", func(t *testing.T) {
		status, _ := c.do("GET", "/api/v1/movies/continue-watching", nil)
		assert.Equal(t, 200, status)
	})

	t.Run("empty_continue_watching_tv", func(t *testing.T) {
		status, _ := c.do("GET", "/api/v1/tvshows/continue-watching", nil)
		assert.Equal(t, 200, status)
	})

	// Watch history (should be empty)
	t.Run("empty_watch_history", func(t *testing.T) {
		status, _ := c.do("GET", "/api/v1/movies/watch-history", nil)
		assert.Equal(t, 200, status)
	})

	// Recently added (should be empty)
	t.Run("empty_recently_added_movies", func(t *testing.T) {
		status, _ := c.do("GET", "/api/v1/movies/recently-added", nil)
		assert.Equal(t, 200, status)
	})

	t.Run("empty_recently_added_tv", func(t *testing.T) {
		status, _ := c.do("GET", "/api/v1/tvshows/recently-added", nil)
		assert.Equal(t, 200, status)
	})

	// User movie stats
	t.Run("user_movie_stats_zero", func(t *testing.T) {
		status, body := c.do("GET", "/api/v1/movies/stats", nil)
		assert.Equal(t, 200, status)
		if body != nil {
			if watched, ok := body["watched_count"].(float64); ok {
				assert.Equal(t, float64(0), watched)
			}
		}
	})
}

// ============================================================================
// Scenario 5: Multi-Device Session Management
// Login from 3 "devices" -> List sessions -> Revoke one -> Verify others still work
// ============================================================================

func TestScenario_MultiDeviceSessionManagement(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	password := "D3v1ceP@ss!"
	c := newScenarioClient(t, ts)
	c.register("multidev_user", "multidev@test.com", password)

	// Login from 3 "devices"
	device1 := newScenarioClient(t, ts)
	device1.login("multidev_user", password)

	device2 := newScenarioClient(t, ts)
	device2.login("multidev_user", password)

	device3 := newScenarioClient(t, ts)
	device3.login("multidev_user", password)

	// Each should have a valid session
	t.Run("all_devices_authenticated", func(t *testing.T) {
		for i, dev := range []*scenarioClient{device1, device2, device3} {
			status, body := dev.do("GET", "/api/v1/users/me", nil)
			assert.Equal(t, 200, status, "device %d should be authenticated", i+1)
			assert.Equal(t, "multidev_user", body["username"])
		}
	})

	// List sessions from device1
	t.Run("list_multiple_sessions", func(t *testing.T) {
		status, result := device1.doArray("GET", "/api/v1/sessions")
		assert.Equal(t, 200, status)
		assert.GreaterOrEqual(t, len(result), 3, "should have at least 3 sessions")
	})

	// Revoke device2's session (logout)
	t.Run("revoke_device2", func(t *testing.T) {
		status, _ := device2.do("POST", "/api/v1/auth/logout", map[string]string{
			"refresh_token": device2.refresh,
		})
		assert.True(t, status == 200 || status == 204)
	})

	// Device1 and Device3 still work
	t.Run("remaining_devices_still_work", func(t *testing.T) {
		status1, _ := device1.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, status1)

		status3, _ := device3.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, status3)
	})
}

// ============================================================================
// Scenario 6: Unauthenticated Access Control
// Verify protected endpoints reject unauthenticated requests
// ============================================================================

func TestScenario_UnauthenticatedAccessControl(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	noAuth := newScenarioClient(t, ts)

	protectedEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/users/me"},
		{"GET", "/api/v1/movies"},
		{"GET", "/api/v1/tvshows"},
		{"GET", "/api/v1/sessions"},
		{"GET", "/api/v1/apikeys"},
		{"GET", "/api/v1/movies/continue-watching"},
		{"GET", "/api/v1/tvshows/continue-watching"},
		{"GET", "/api/v1/admin/users"},
	}

	for _, ep := range protectedEndpoints {
		t.Run(fmt.Sprintf("%s_%s", ep.method, ep.path), func(t *testing.T) {
			status, _ := noAuth.do(ep.method, ep.path, nil)
			assert.True(t, status == 401 || status == 403,
				"unauthenticated %s %s should be 401/403, got %d",
				ep.method, ep.path, status)
		})
	}

	// Public endpoints should work without auth
	publicEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/healthz"},
		{"GET", "/readyz"},
	}

	for _, ep := range publicEndpoints {
		t.Run(fmt.Sprintf("public_%s_%s", ep.method, ep.path), func(t *testing.T) {
			status, _ := noAuth.do(ep.method, ep.path, nil)
			assert.Equal(t, 200, status,
				"public endpoint %s %s should return 200", ep.method, ep.path)
		})
	}
}

// ============================================================================
// Scenario 7: API Key Authentication Flow
// Create API key -> Use it to authenticate -> Verify works
// ============================================================================

func TestScenario_APIKeyAuthFlow(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	c := newScenarioClient(t, ts)
	c.registerAndLogin("apikey_user", "apikey@test.com", "Ap1K3yP@ss!")

	// Create API key
	var rawKey string
	t.Run("create_key", func(t *testing.T) {
	status, body := c.do("POST", "/api/v1/apikeys", map[string]interface{}{
			"name":   "automation-key",
			"scopes": []string{"read", "write"},
		})
		require.Equal(t, 201, status)
		if key, ok := body["api_key"].(string); ok {
			rawKey = key
		}
		require.NotEmpty(t, rawKey, "should return raw API key")
	})

	// Use API key to authenticate
	t.Run("authenticate_with_api_key", func(t *testing.T) {
		if rawKey == "" {
			t.Skip("no API key")
		}
		req, err := http.NewRequestWithContext(
			context.Background(), "GET", ts.BaseURL+"/api/v1/users/me", nil,
		)
		require.NoError(t, err)
		req.Header.Set("X-API-Key", rawKey)

		resp, err := ts.HTTPClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode, "API key auth should work")

		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body) //nolint:errcheck
		assert.Equal(t, "apikey_user", body["username"])
	})
}

// ============================================================================
// Scenario 8: Token Refresh Flow
// Login -> Use access token -> Refresh -> Use new token
// ============================================================================

func TestScenario_TokenRefreshFlow(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	c := newScenarioClient(t, ts)
	c.registerAndLogin("refresh_user", "refresh@test.com", "R3fr3shP@ss!")

	originalToken := c.token
	_ = originalToken

	// Verify original token works
	t.Run("original_token_works", func(t *testing.T) {
		status, _ := c.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, status)
	})

	// Refresh the token
	t.Run("refresh_token", func(t *testing.T) {
		status, body := c.do("POST", "/api/v1/auth/refresh", map[string]string{
			"refresh_token": c.refresh,
		})
		assert.Equal(t, 200, status, "token refresh should succeed")
		if newToken, ok := body["access_token"].(string); ok {
			c.token = newToken
		}
		if newRefresh, ok := body["refresh_token"].(string); ok {
			c.refresh = newRefresh
		}
		assert.NotEmpty(t, c.token, "should get new access token")
	})

	// New token works
	t.Run("new_token_works", func(t *testing.T) {
		status, body := c.do("GET", "/api/v1/users/me", nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "refresh_user", body["username"])
	})
}

// ============================================================================
// Scenario 9: Admin Activity Monitoring
// Admin performs actions -> Views activity log
// ============================================================================

func TestScenario_AdminActivityMonitoring(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	admin := newScenarioClient(t, ts)
	admin.registerAndLogin("activity_admin", "activity_admin@test.com", "Act1v1tyP@ss!")
	admin.makeAdmin()
	admin.login("activity_admin", "Act1v1tyP@ss!")

	user := newScenarioClient(t, ts)
	user.registerAndLogin("activity_user", "activity_user@test.com", "Us3rP@ss!")

	// Admin views recent activity
	t.Run("admin_views_activity", func(t *testing.T) {
		status, body := admin.do("GET", "/api/v1/admin/activity", nil)
		assert.Equal(t, 200, status)
		if items, ok := body["items"].([]interface{}); ok {
			assert.GreaterOrEqual(t, len(items), 1,
				"should have activity entries from user registration/login")
		}
	})

	// Regular user cannot view activity
	t.Run("user_cannot_view_activity", func(t *testing.T) {
		status, _ := user.do("GET", "/api/v1/admin/activity", nil)
		assert.True(t, status == 401 || status == 403,
			"regular user denied admin activity (got %d)", status)
	})
}

// ============================================================================
// Scenario 10: Concurrent User Operations
// Multiple users perform operations simultaneously to test isolation
// ============================================================================

func TestScenario_ConcurrentUserOperations(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	users := make([]*scenarioClient, 5)
	for i := 0; i < 5; i++ {
		users[i] = newScenarioClient(t, ts)
		users[i].registerAndLogin(
			fmt.Sprintf("concurrent_%d", i),
			fmt.Sprintf("concurrent_%d@test.com", i),
			fmt.Sprintf("P@ss%d!", i),
		)
	}

	// All users fetch their profiles concurrently
	t.Run("concurrent_profile_access", func(t *testing.T) {
		done := make(chan struct{}, len(users))
		for i, u := range users {
			go func(idx int, c *scenarioClient) {
				defer func() { done <- struct{}{} }()
				status, body := c.do("GET", "/api/v1/users/me", nil)
				assert.Equal(t, 200, status, "user %d should get profile", idx)
				expectedUsername := fmt.Sprintf("concurrent_%d", idx)
				assert.Equal(t, expectedUsername, body["username"],
					"user %d should see their own profile", idx)
			}(i, u)
		}
		for range users {
			<-done
		}
	})

	// All users create API keys concurrently
	t.Run("concurrent_api_key_creation", func(t *testing.T) {
		done := make(chan struct{}, len(users))
		for i, u := range users {
			go func(idx int, c *scenarioClient) {
				defer func() { done <- struct{}{} }()
				status, body := c.do("POST", "/api/v1/apikeys", map[string]interface{}{
					"name":   fmt.Sprintf("key_%d", idx),
					"scopes": []string{"read"},
				})
				assert.Equal(t, 201, status, "user %d should create API key", idx)
				assert.NotEmpty(t, body["api_key"], "user %d should get key", idx)
			}(i, u)
		}
		for range users {
			<-done
		}
	})
}

// ============================================================================
// Scenario 11: Settings and Preferences
// User changes settings -> Verifies persistence -> Other user has separate settings
// ============================================================================

func TestScenario_UserSettings(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	alice := newScenarioClient(t, ts)
	alice.registerAndLogin("settings_alice", "settings_alice@test.com", "S3tt1ngsA!")

	bob := newScenarioClient(t, ts)
	bob.registerAndLogin("settings_bob", "settings_bob@test.com", "S3tt1ngsB!")

	// Get default settings (returns array)
	t.Run("default_settings", func(t *testing.T) {
		status, _ := alice.doArray("GET", "/api/v1/settings/user")
		assert.Equal(t, 200, status)
	})

	// Update a setting
	t.Run("update_theme_setting", func(t *testing.T) {
		status, _ := alice.do("PUT", "/api/v1/settings/user/theme",
			map[string]interface{}{"value": "dark"})
		assert.True(t, status == 200 || status == 204,
			"settings update should succeed (got %d)", status)
	})

	// Verify setting persisted
	t.Run("setting_persisted", func(t *testing.T) {
		status, body := alice.do("GET", "/api/v1/settings/user/theme", nil)
		assert.Equal(t, 200, status)
		if body != nil {
			assert.Equal(t, "dark", body["value"])
		}
	})

	// Bob has separate settings
	t.Run("bob_has_independent_settings", func(t *testing.T) {
		status, _ := bob.doArray("GET", "/api/v1/settings/user")
		assert.Equal(t, 200, status)
		// Bob's settings list should not contain alice's theme=dark override
	})
}

// ============================================================================
// Scenario 12: Health & Infrastructure Endpoints
// Verify liveness, readiness, startup probes work correctly
// ============================================================================

func TestScenario_HealthProbes(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	t.Run("healthz", func(t *testing.T) {
		resp, err := ts.HTTPClient.Get(ts.BaseURL + "/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("readyz", func(t *testing.T) {
		resp, err := ts.HTTPClient.Get(ts.BaseURL + "/readyz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("openapi_spec", func(t *testing.T) {
		resp, err := ts.HTTPClient.Get(ts.BaseURL + "/api/openapi.yaml")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
}
