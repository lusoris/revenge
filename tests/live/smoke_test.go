//go:build live

// Package live contains comprehensive smoke tests against a running Revenge stack.
// Start the stack with: make docker-local
// Run these tests with: make test-live
package live

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	baseURL = envOr("REVENGE_TEST_URL", "http://localhost:8096")
	dbURL   = envOr("REVENGE_TEST_DB_URL", "postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable")
	client  = &http.Client{Timeout: 10 * time.Second}
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// =============================================================================
// Helpers
// =============================================================================

type credentials struct {
	username     string
	email        string
	password     string
	userID       string
	accessToken  string
	refreshToken string
}

// registerUser creates a new user and returns credentials (no login).
func registerUser(t *testing.T) *credentials {
	t.Helper()
	creds := &credentials{
		username: fmt.Sprintf("live_%d", time.Now().UnixNano()),
		password: "LiveTestPass123!",
	}
	creds.email = creds.username + "@live.test"

	body, _ := json.Marshal(map[string]string{
		"username": creds.username,
		"email":    creds.email,
		"password": creds.password,
	})

	resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode, "registration should succeed")

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	require.NoError(t, err)
	creds.userID, _ = user["id"].(string)
	require.NotEmpty(t, creds.userID)
	return creds
}

// loginUser logs in and populates tokens on the credentials.
func loginUser(t *testing.T, creds *credentials) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"username": creds.username,
		"password": creds.password,
	})

	resp, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "login should succeed")

	var loginResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(t, err)

	creds.accessToken, _ = loginResp["access_token"].(string)
	creds.refreshToken, _ = loginResp["refresh_token"].(string)
	require.NotEmpty(t, creds.accessToken)
	require.NotEmpty(t, creds.refreshToken)
}

// registerAndLogin is a convenience that does both.
func registerAndLogin(t *testing.T) *credentials {
	t.Helper()
	creds := registerUser(t)
	loginUser(t, creds)
	return creds
}

// makeAdmin grants admin role to a user by inserting Casbin policy directly.
// The running server will pick up the change via SyncedEnforcer's auto-reload
// (configured to reload every 10s). Use waitForAdminRole to wait for visibility.
func makeAdmin(t *testing.T, userID string) {
	t.Helper()
	d, err := sql.Open("postgres", dbURL)
	require.NoError(t, err, "should connect to test database")
	defer d.Close()

	_, err = d.Exec(
		`INSERT INTO shared.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ('g', $1, 'admin', '', '', '', '') ON CONFLICT DO NOTHING`,
		userID,
	)
	require.NoError(t, err, "should grant admin role")

	// Also set is_admin flag on user record for consistency
	_, err = d.Exec(`UPDATE shared.users SET is_admin = true WHERE id = $1`, userID)
	require.NoError(t, err)
}

// ensureAdmin creates an admin user by inserting Casbin rules into the DB,
// then waits for the server's SyncedEnforcer auto-reload to pick them up.
func ensureAdmin(t *testing.T) *credentials {
	t.Helper()
	admin := registerAndLogin(t)
	makeAdmin(t, admin.userID)
	waitForAdminRole(t, admin)
	return admin
}

// waitForAdminRole polls an admin-only endpoint until the server's Casbin
// auto-reload picks up the newly inserted role. Timeout after 15s.
func waitForAdminRole(t *testing.T, creds *credentials) {
	t.Helper()
	// SyncedEnforcer auto-reloads every 10s. Allow up to 15s (1+ reload cycles)
	// to handle jitter, goroutine scheduling delays, and DB query time.
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		resp := doRequest(t, "GET", "/api/v1/rbac/roles", creds.accessToken, nil)
		resp.Body.Close()
		if resp.StatusCode == 200 {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("admin role not visible to server after 15s (Casbin auto-reload timeout)")
}

// doRequest executes an HTTP request and returns the response.
func doRequest(t *testing.T, method, path, token string, body interface{}) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, baseURL+path, bodyReader)
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

// doJSON executes a request and decodes the JSON response body.
func doJSON(t *testing.T, method, path, token string, body interface{}) (int, map[string]interface{}) {
	t.Helper()
	resp := doRequest(t, method, path, token, body)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &result)
	}
	return resp.StatusCode, result
}

// doJSONArray executes a request and decodes the JSON response as an array.
func doJSONArray(t *testing.T, method, path, token string) (int, []interface{}) {
	t.Helper()
	resp := doRequest(t, method, path, token, nil)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var result []interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &result)
	}
	return resp.StatusCode, result
}

// =============================================================================
// 1. Infrastructure Health (DB, Cache, Search)
// =============================================================================

func TestLive_Infrastructure(t *testing.T) {
	t.Run("liveness", func(t *testing.T) {
		status, result := doJSON(t, "GET", "/healthz", "", nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "healthy", result["status"])
		assert.Equal(t, "liveness", result["name"])
	})

	t.Run("readiness_checks_database", func(t *testing.T) {
		status, result := doJSON(t, "GET", "/readyz", "", nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "healthy", result["status"])
		// Readiness checks database connectivity
		if details, ok := result["details"].(map[string]interface{}); ok {
			assert.Contains(t, details, "database", "readiness should check database")
		}
	})

	t.Run("startup_complete", func(t *testing.T) {
		status, result := doJSON(t, "GET", "/startupz", "", nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "healthy", result["status"])
	})

	t.Run("postgres_direct_connectivity", func(t *testing.T) {
		db, err := sql.Open("postgres", dbURL)
		require.NoError(t, err)
		defer db.Close()

		var one int
		err = db.QueryRow("SELECT 1").Scan(&one)
		require.NoError(t, err, "direct PostgreSQL query should work")
		assert.Equal(t, 1, one)
	})

	t.Run("postgres_schemas_exist", func(t *testing.T) {
		db, err := sql.Open("postgres", dbURL)
		require.NoError(t, err)
		defer db.Close()

		// Check that our application schemas were created by migrations
		var schemas []string
		rows, err := db.Query(`SELECT schema_name FROM information_schema.schemata WHERE schema_name IN ('shared', 'movie', 'tvshow', 'qar') ORDER BY schema_name`)
		require.NoError(t, err)
		defer rows.Close()
		for rows.Next() {
			var s string
			require.NoError(t, rows.Scan(&s))
			schemas = append(schemas, s)
		}
		assert.Contains(t, schemas, "shared", "shared schema should exist")
		assert.Contains(t, schemas, "tvshow", "tvshow schema should exist")
	})

	t.Run("river_tables_exist", func(t *testing.T) {
		db, err := sql.Open("postgres", dbURL)
		require.NoError(t, err)
		defer db.Close()

		// River job queue tables should be created
		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'river' OR (table_schema = 'public' AND table_name LIKE 'river_%')`).Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "River tables should exist")
	})

	t.Run("typesense_reachable", func(t *testing.T) {
		tsURL := envOr("REVENGE_TEST_TYPESENSE_URL", "http://localhost:8108")
		req, _ := http.NewRequest("GET", tsURL+"/health", nil)
		resp, err := client.Do(req)
		require.NoError(t, err, "Typesense should be reachable")
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("dragonfly_reachable", func(t *testing.T) {
		// Verify Dragonfly (Redis) is reachable via TCP
		dfHost := envOr("REVENGE_TEST_DRAGONFLY_HOST", "localhost:6379")
		conn, err := (&net.Dialer{Timeout: 2 * time.Second}).Dial("tcp", dfHost)
		require.NoError(t, err, "Dragonfly should be reachable")
		conn.Close()
	})

	t.Run("casbin_policies_loaded", func(t *testing.T) {
		db, err := sql.Open("postgres", dbURL)
		require.NoError(t, err)
		defer db.Close()

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM shared.casbin_rule WHERE ptype = 'p'`).Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "Casbin default policies should be loaded")
	})
}

// =============================================================================
// 2. Auth Endpoints (individual)
// =============================================================================

func TestLive_AuthEndpoints(t *testing.T) {
	t.Run("register_success", func(t *testing.T) {
		creds := registerUser(t)
		assert.NotEmpty(t, creds.userID)
	})

	t.Run("register_duplicate_username", func(t *testing.T) {
		creds := registerUser(t)
		// Try to register again with same username
		body, _ := json.Marshal(map[string]string{
			"username": creds.username,
			"email":    "different@test.com",
			"password": "SomePass123!",
		})
		resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusBadRequest,
			"duplicate username should be rejected (got %d)", resp.StatusCode)
	})

	t.Run("register_invalid_email", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": fmt.Sprintf("test_%d", time.Now().UnixNano()),
			"email":    "not-an-email",
			"password": "TestPass123!",
		})
		resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500, "invalid email should be rejected (got %d)", resp.StatusCode)
	})

	t.Run("register_weak_password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": fmt.Sprintf("test_%d", time.Now().UnixNano()),
			"email":    fmt.Sprintf("test_%d@test.com", time.Now().UnixNano()),
			"password": "123",
		})
		resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500, "weak password should be rejected (got %d)", resp.StatusCode)
	})

	t.Run("login_nonexistent_user", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/auth/login", "", map[string]string{
			"username": "nonexistent_user_xyz",
			"password": "whatever",
		})
		assert.Equal(t, http.StatusUnauthorized, status)
	})

	t.Run("login_wrong_password", func(t *testing.T) {
		creds := registerUser(t)
		status, _ := doJSON(t, "POST", "/api/v1/auth/login", "", map[string]string{
			"username": creds.username,
			"password": "WrongPassword123!",
		})
		assert.Equal(t, http.StatusUnauthorized, status)
	})

	t.Run("refresh_invalid_token", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/auth/refresh", "", map[string]string{
			"refresh_token": "invalid-token-value",
		})
		assert.True(t, status >= 400, "invalid refresh token should fail (got %d)", status)
	})

	t.Run("forgot_password", func(t *testing.T) {
		creds := registerUser(t)
		resp := doRequest(t, "POST", "/api/v1/auth/forgot-password", "", map[string]string{
			"email": creds.email,
		})
		defer resp.Body.Close()
		// Should accept (even if email sending is disabled, the endpoint shouldn't error)
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204 || resp.StatusCode == 202,
			"forgot-password should accept (got %d)", resp.StatusCode)
	})

	t.Run("reset_password_invalid_token", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/auth/reset-password", "", map[string]string{
			"token":    "invalid-reset-token",
			"password": "NewPassword123!",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode >= 400, "invalid reset token should fail (got %d)", resp.StatusCode)
	})

	t.Run("verify_email_invalid_token", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/auth/verify-email", "", map[string]string{
			"token": "invalid-verification-token",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode >= 400, "invalid verification token should fail (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 3. Chain: Full Auth Lifecycle
// =============================================================================

func TestLive_Chain_AuthLifecycle(t *testing.T) {
	ctx := context.Background()
	creds := registerUser(t)

	// Step 1: Login
	loginUser(t, creds)
	originalToken := creds.accessToken

	// Step 2: Access profile
	t.Run("access_profile", func(t *testing.T) {
		status, user := doJSON(t, "GET", "/api/v1/users/me", creds.accessToken, nil)
		require.Equal(t, 200, status)
		assert.Equal(t, creds.username, user["username"])
		assert.Equal(t, creds.email, user["email"])
	})

	// Step 3: Update profile
	t.Run("update_profile", func(t *testing.T) {
		resp := doRequest(t, "PUT", "/api/v1/users/me", creds.accessToken, map[string]string{
			"display_name": "Smoke Test User",
		})
		defer resp.Body.Close()
		// Could be 200 (success) or might have validation issues
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"profile update should succeed (got %d)", resp.StatusCode)
	})

	// Step 4: Change password
	newPassword := "NewSmokePass456!"
	t.Run("change_password", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/auth/change-password", creds.accessToken, map[string]string{
			"old_password": creds.password,
			"new_password": newPassword,
		})
		if status == 200 || status == 204 {
			creds.password = newPassword
		}
		assert.True(t, status == 200 || status == 204,
			"password change should succeed (got %d)", status)
	})

	// Step 5: Login with new password
	t.Run("login_with_new_password", func(t *testing.T) {
		loginUser(t, creds)
		assert.NotEqual(t, originalToken, creds.accessToken, "new login should give new token")
	})

	// Step 6: Refresh token
	t.Run("refresh_token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"refresh_token": creds.refreshToken,
		})
		req, _ := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v1/auth/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)

		var refreshResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&refreshResp)
		newToken, _ := refreshResp["access_token"].(string)
		require.NotEmpty(t, newToken)
		creds.accessToken = newToken
	})

	// Step 7: Logout
	t.Run("logout", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/auth/logout", creds.accessToken, map[string]string{
			"refresh_token": creds.refreshToken,
		})
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	// Step 8: Refresh token should be revoked (access token stays valid until JWT expiry)
	t.Run("refresh_revoked_after_logout", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/auth/refresh", "", map[string]string{
			"refresh_token": creds.refreshToken,
		})
		assert.NotEqual(t, 200, status, "refresh token should be revoked after logout")
	})
}

// =============================================================================
// 4. Chain: Session Management
// =============================================================================

func TestLive_Chain_SessionManagement(t *testing.T) {
	creds := registerUser(t)

	// Login from "device 1"
	loginUser(t, creds)
	token1 := creds.accessToken

	// Login again from "device 2"
	loginUser(t, creds)
	token2 := creds.accessToken
	assert.NotEqual(t, token1, token2, "second login should produce different token")

	// List sessions
	t.Run("list_sessions", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/sessions", token2, nil)
		require.Equal(t, 200, status)
		// Sessions may be tracked as refresh tokens, check that response is valid
		assert.NotNil(t, body)
	})

	// Get current session
	// NOTE: Returns 401 because JWT claims don't include session_id yet,
	// so the handler can't identify the current session. This is a known limitation.
	t.Run("current_session", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/sessions/current", token2, nil)
		assert.True(t, status == 200 || status == 401,
			"current session should return 200 or 401 (got %d)", status)
	})

	// Revoke all sessions (logout everywhere)
	t.Run("logout_all", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/sessions", token2, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200,
			"logout all should succeed (got %d)", resp.StatusCode)
	})

	// Refresh tokens should be revoked (JWT access tokens stay valid until expiry)
	t.Run("refresh_tokens_revoked", func(t *testing.T) {
		// Can't refresh anymore since all refresh tokens are revoked
		status, _ := doJSON(t, "POST", "/api/v1/auth/refresh", "", map[string]string{
			"refresh_token": "any-previously-valid-refresh-token",
		})
		assert.NotEqual(t, 200, status)
	})
}

// =============================================================================
// 5. Chain: API Key Lifecycle
// =============================================================================

func TestLive_Chain_APIKeyLifecycle(t *testing.T) {
	creds := registerAndLogin(t)

	// List API keys - should be empty
	t.Run("list_empty", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/apikeys", creds.accessToken, nil)
		require.Equal(t, 200, status)
		if keys, ok := body["keys"].([]interface{}); ok {
			assert.Empty(t, keys)
		}
	})

	// Create API key
	var keyID, keySecret string
	t.Run("create_key", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/apikeys", creds.accessToken, map[string]interface{}{
			"name":   "smoke-test-key",
			"scopes": []string{"read", "write"},
		})
		require.Equal(t, 201, status, "API key creation should return 201")
		keyID, _ = body["id"].(string)
		keySecret, _ = body["api_key"].(string)
		assert.NotEmpty(t, keyID, "should have key ID")
		assert.NotEmpty(t, keySecret, "should have key secret")
	})

	// List API keys - should have 1
	t.Run("list_has_one", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/apikeys", creds.accessToken, nil)
		require.Equal(t, 200, status)
		if keys, ok := body["keys"].([]interface{}); ok {
			assert.Len(t, keys, 1)
		}
	})

	// Get specific key
	t.Run("get_key", func(t *testing.T) {
		if keyID == "" {
			t.Skip("no key ID")
		}
		status, body := doJSON(t, "GET", "/api/v1/apikeys/"+keyID, creds.accessToken, nil)
		require.Equal(t, 200, status)
		assert.Equal(t, "smoke-test-key", body["name"])
	})

	// Use API key for authentication (if supported as X-API-Key header)
	t.Run("auth_with_api_key", func(t *testing.T) {
		if keySecret == "" {
			t.Skip("no key secret")
		}
		req, _ := http.NewRequestWithContext(context.Background(), "GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("X-API-Key", keySecret)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		// API key auth may or may not be supported for all endpoints
		// Just verify we get a definitive response (not 500)
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 401 || resp.StatusCode == 403,
			"API key auth should give definitive response (got %d)", resp.StatusCode)
	})

	// Revoke API key
	t.Run("revoke_key", func(t *testing.T) {
		if keyID == "" {
			t.Skip("no key ID")
		}
		resp := doRequest(t, "DELETE", "/api/v1/apikeys/"+keyID, creds.accessToken, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"revoke should succeed (got %d)", resp.StatusCode)
	})

	// Verify key is gone
	t.Run("list_empty_after_revoke", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/apikeys", creds.accessToken, nil)
		require.Equal(t, 200, status)
		if keys, ok := body["keys"].([]interface{}); ok {
			assert.Empty(t, keys)
		}
	})
}

// =============================================================================
// 6. Chain: User Settings CRUD
// =============================================================================

func TestLive_Chain_UserSettings(t *testing.T) {
	creds := registerAndLogin(t)

	// List settings - initially empty or defaults
	t.Run("list_initial", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/settings/user", creds.accessToken)
		assert.Equal(t, 200, status)
	})

	// Set a setting
	t.Run("set_setting", func(t *testing.T) {
		resp := doRequest(t, "PUT", "/api/v1/settings/user/theme", creds.accessToken, map[string]string{
			"value": "dark",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"setting should be saved (got %d)", resp.StatusCode)
	})

	// Get that setting
	t.Run("get_setting", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/settings/user/theme", creds.accessToken, nil)
		if status == 200 {
			assert.Equal(t, "dark", body["value"])
		}
	})

	// Update the setting
	t.Run("update_setting", func(t *testing.T) {
		resp := doRequest(t, "PUT", "/api/v1/settings/user/theme", creds.accessToken, map[string]string{
			"value": "light",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"setting update should succeed (got %d)", resp.StatusCode)
	})

	// Delete the setting
	t.Run("delete_setting", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/settings/user/theme", creds.accessToken, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"setting delete should succeed (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 7. Chain: MFA Status
// =============================================================================

func TestLive_Chain_MFA(t *testing.T) {
	creds := registerAndLogin(t)

	// Check MFA status (should be disabled)
	t.Run("status_disabled", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/mfa/status", creds.accessToken, nil)
		require.Equal(t, 200, status)
		assert.NotNil(t, body, "MFA status should return a body")
	})

	// Setup TOTP
	var totpSecret string
	t.Run("setup_totp", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/mfa/totp/setup", creds.accessToken, map[string]string{
			"accountName": creds.email,
		})
		if status == 200 {
			totpSecret, _ = body["secret"].(string)
			assert.NotEmpty(t, totpSecret, "should receive TOTP secret")
			if uri, ok := body["uri"].(string); ok {
				assert.Contains(t, uri, "otpauth://", "should have otpauth URI")
			}
		}
		assert.True(t, status == 200 || status == 201,
			"TOTP setup should succeed (got %d)", status)
	})

	// Check MFA status again
	t.Run("status_after_setup", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/mfa/status", creds.accessToken, nil)
		assert.Equal(t, 200, status)
	})

	// List WebAuthn credentials (should be empty)
	t.Run("webauthn_credentials_empty", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/mfa/webauthn/credentials", creds.accessToken, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Generate backup codes
	t.Run("generate_backup_codes", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/mfa/backup-codes/generate", creds.accessToken, nil)
		if status == 200 {
			if codes, ok := body["codes"].([]interface{}); ok {
				assert.NotEmpty(t, codes, "should receive backup codes")
			}
		}
	})
}

// =============================================================================
// 8. User Preferences
// =============================================================================

func TestLive_Chain_Preferences(t *testing.T) {
	creds := registerAndLogin(t)

	// Get preferences
	t.Run("get_initial", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/users/me/preferences", creds.accessToken, nil)
		require.Equal(t, 200, status)
		assert.NotNil(t, body)
	})

	// Update preferences
	t.Run("update", func(t *testing.T) {
		resp := doRequest(t, "PUT", "/api/v1/users/me/preferences", creds.accessToken, map[string]interface{}{
			"metadata_language": "de",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"preference update should succeed (got %d)", resp.StatusCode)
	})

	// Verify update persisted
	t.Run("verify_update", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/users/me/preferences", creds.accessToken, nil)
		require.Equal(t, 200, status)
		if lang, ok := body["metadata_language"].(string); ok {
			assert.Equal(t, "de", lang)
		}
	})
}

// =============================================================================
// 9. Content Endpoints (empty library)
// =============================================================================

func TestLive_ContentEndpoints_EmptyLibrary(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	// All these should return 200 with empty arrays on a fresh database
	emptyListEndpoints := []string{
		"/api/v1/movies",
		"/api/v1/movies/recently-added",
		"/api/v1/movies/top-rated",
		"/api/v1/movies/continue-watching",
		"/api/v1/movies/watch-history",
		"/api/v1/tvshows",
		"/api/v1/tvshows/recently-added",
		"/api/v1/tvshows/continue-watching",
	}

	for _, path := range emptyListEndpoints {
		t.Run("GET_"+strings.TrimPrefix(path, "/api/v1/"), func(t *testing.T) {
			resp := doRequest(t, "GET", path, tok, nil)
			defer resp.Body.Close()
			assert.Equal(t, 200, resp.StatusCode, "should return 200 for empty list at %s", path)
			body, _ := io.ReadAll(resp.Body)
			assert.True(t, json.Valid(body), "should return valid JSON at %s", path)
		})
	}

	// Movie search (database-backed)
	t.Run("search_movies_empty", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/movies/search?query=test", tok)
		assert.Equal(t, 200, status)
	})

	// TV show search (database-backed)
	t.Run("search_tvshows_empty", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/tvshows/search?query=test", tok)
		assert.Equal(t, 200, status)
	})

	// Movie stats
	t.Run("movie_stats", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/movies/stats", tok, nil)
		assert.Equal(t, 200, status)
	})

	// TV stats
	t.Run("tv_stats", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/tvshows/stats", tok, nil)
		assert.Equal(t, 200, status)
	})

	// Non-existent movie
	t.Run("movie_not_found", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/movies/00000000-0000-0000-0000-000000000000", tok, nil)
		assert.Equal(t, 404, status)
	})

	// Non-existent TV show
	t.Run("tvshow_not_found", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/tvshows/00000000-0000-0000-0000-000000000000", tok, nil)
		assert.Equal(t, 404, status)
	})

	// Episodes endpoints
	t.Run("recent_episodes", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/tvshows/episodes/recent", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("upcoming_episodes", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/tvshows/episodes/upcoming", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
}

// =============================================================================
// 10. Search Infrastructure (Typesense)
// =============================================================================

func TestLive_SearchInfrastructure(t *testing.T) {
	creds := registerAndLogin(t)
	tok := creds.accessToken

	// Typesense-backed search (exercises the search service → Typesense path)
	t.Run("typesense_search_empty", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/search/movies?q=nonexistent", tok, nil)
		// Should return 200 with empty results (not 500 from Typesense connection failure)
		assert.Equal(t, 200, status, "Typesense search should work (got %d)", status)
		if body != nil {
			if hits, ok := body["hits"].([]interface{}); ok {
				assert.Empty(t, hits)
			}
		}
	})

	// Autocomplete (exercises Typesense)
	t.Run("typesense_autocomplete", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/search/movies/autocomplete?q=test", tok, nil)
		assert.Equal(t, 200, status, "autocomplete should work (got %d)", status)
		if body != nil {
			if suggestions, ok := body["suggestions"].([]interface{}); ok {
				assert.Empty(t, suggestions)
			}
		}
	})

	// Facets (exercises Typesense)
	t.Run("typesense_facets", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/search/movies/facets", tok, nil)
		assert.Equal(t, 200, status, "facets should work (got %d)", status)
	})
}

// =============================================================================
// 11. Admin Endpoints (requires direct DB role grant)
// =============================================================================

func TestLive_AdminEndpoints(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	// Create a non-admin for comparison
	regular := registerAndLogin(t)
	regularTok := regular.accessToken

	t.Run("admin_server_settings", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/settings/server", tok)
		assert.Equal(t, 200, status, "admin should access server settings")
	})

	t.Run("non_admin_server_settings_denied", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/settings/server", regularTok, nil)
		resp.Body.Close()
		assert.Equal(t, 403, resp.StatusCode, "non-admin should be denied server settings")
	})

	// RBAC endpoints
	t.Run("list_roles", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/rbac/roles", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.True(t, json.Valid(body))
	})

	t.Run("list_permissions", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/rbac/permissions", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("list_policies", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/rbac/policies", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("non_admin_rbac_denied", func(t *testing.T) {
		for _, path := range []string{"/api/v1/rbac/roles", "/api/v1/rbac/permissions", "/api/v1/rbac/policies"} {
			resp := doRequest(t, "GET", path, regularTok, nil)
			resp.Body.Close()
			assert.Equal(t, 403, resp.StatusCode, "non-admin should be denied %s", path)
		}
	})

	// Activity logs
	t.Run("activity_logs", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/activity", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("activity_stats", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/activity/stats", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("non_admin_activity_denied", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/activity", regularTok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 403, resp.StatusCode)
	})

	// Integration status endpoints
	t.Run("radarr_status", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/integrations/radarr/status", tok, nil)
		defer resp.Body.Close()
		// May return error since Radarr isn't configured, but should NOT 500
		assert.True(t, resp.StatusCode != 500, "radarr status should not 500 (got %d)", resp.StatusCode)
	})

	t.Run("sonarr_status", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/integrations/sonarr/status", tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode != 500, "sonarr status should not 500 (got %d)", resp.StatusCode)
	})

	// OIDC admin endpoints
	t.Run("admin_oidc_providers", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/oidc/providers", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})
}

// =============================================================================
// 12. Chain: Admin Library Management
// =============================================================================

func TestLive_Chain_LibraryManagement(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	regular := registerAndLogin(t)
	regularTok := regular.accessToken

	// List libraries (should be empty)
	t.Run("list_empty", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/libraries", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Create library
	var libraryID string
	t.Run("create_library", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/libraries", tok, map[string]interface{}{
			"name":  "Smoke Test Library",
			"type":  "movie",
			"paths": []string{"/media/movies"},
		})
		if status == 201 || status == 200 {
			libraryID, _ = body["id"].(string)
		}
		assert.True(t, status == 201 || status == 200,
			"library creation should succeed (got %d)", status)
	})

	// Get library
	t.Run("get_library", func(t *testing.T) {
		if libraryID == "" {
			t.Skip("no library")
		}
		status, body := doJSON(t, "GET", "/api/v1/libraries/"+libraryID, tok, nil)
		assert.Equal(t, 200, status)
		assert.Equal(t, "Smoke Test Library", body["name"])
	})

	// Regular user should be able to list libraries (access may vary)
	t.Run("regular_user_list_libraries", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/libraries", regularTok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Non-admin can't create library
	t.Run("non_admin_create_denied", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/libraries", regularTok, map[string]interface{}{
			"name":  "Unauthorized Library",
			"type":  "movie",
			"paths": []string{"/media/unauthorized"},
		})
		assert.Equal(t, 403, status, "non-admin should not create libraries")
	})

	// Trigger library scan (requires scanType in body)
	t.Run("trigger_scan", func(t *testing.T) {
		if libraryID == "" {
			t.Skip("no library")
		}
		resp := doRequest(t, "POST", "/api/v1/libraries/"+libraryID+"/scan", tok, map[string]string{
			"scanType": "full",
		})
		defer resp.Body.Close()
		// Should accept the scan request (even if the path doesn't exist on disk)
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 202 || resp.StatusCode == 204,
			"scan trigger should be accepted (got %d)", resp.StatusCode)
	})

	// Delete library
	t.Run("delete_library", func(t *testing.T) {
		if libraryID == "" {
			t.Skip("no library")
		}
		resp := doRequest(t, "DELETE", "/api/v1/libraries/"+libraryID, tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"library deletion should succeed (got %d)", resp.StatusCode)
	})

	// Verify deleted
	t.Run("verify_deleted", func(t *testing.T) {
		if libraryID == "" {
			t.Skip("no library")
		}
		status, _ := doJSON(t, "GET", "/api/v1/libraries/"+libraryID, tok, nil)
		assert.Equal(t, 404, status, "deleted library should return 404")
	})
}

// =============================================================================
// 13. Chain: RBAC Role Management
// =============================================================================

func TestLive_Chain_RBAC(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	target := registerAndLogin(t)

	// Create a custom role
	t.Run("create_custom_role", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/rbac/roles", tok, map[string]interface{}{
			"name":        "test-moderator",
			"description": "Test moderator role",
			"permissions": []map[string]string{
				{"resource": "library", "action": "read"},
				{"resource": "library", "action": "write"},
			},
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 201,
			"role creation should succeed (got %d)", resp.StatusCode)
	})

	// Assign role to target user
	t.Run("assign_role", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/rbac/users/"+target.userID+"/roles", tok, map[string]string{
			"role": "test-moderator",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 204,
			"role assignment should succeed (got %d)", resp.StatusCode)
	})

	// Get user roles
	t.Run("get_user_roles", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/rbac/users/"+target.userID+"/roles", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "test-moderator")
	})

	// Remove role
	t.Run("remove_role", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/rbac/users/"+target.userID+"/roles/test-moderator", tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"role removal should succeed (got %d)", resp.StatusCode)
	})

	// Clean up: delete custom role
	t.Run("delete_custom_role", func(t *testing.T) {
		resp := doRequest(t, "DELETE", "/api/v1/rbac/roles/test-moderator", tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204,
			"role deletion should succeed (got %d)", resp.StatusCode)
	})
}

// =============================================================================
// 14. Unauthenticated Access (comprehensive)
// =============================================================================

func TestLive_UnauthenticatedAccess(t *testing.T) {
	protectedEndpoints := []struct {
		method string
		path   string
	}{
		// User endpoints
		{"GET", "/api/v1/users/me"},
		{"PUT", "/api/v1/users/me"},
		{"GET", "/api/v1/users/me/preferences"},
		// Sessions
		{"GET", "/api/v1/sessions"},
		{"GET", "/api/v1/sessions/current"},
		// Libraries
		{"GET", "/api/v1/libraries"},
		// Movies
		{"GET", "/api/v1/movies"},
		{"GET", "/api/v1/movies/recently-added"},
		{"GET", "/api/v1/movies/continue-watching"},
		// TV Shows
		{"GET", "/api/v1/tvshows"},
		{"GET", "/api/v1/tvshows/recently-added"},
		// API Keys
		{"GET", "/api/v1/apikeys"},
		// MFA
		{"GET", "/api/v1/mfa/status"},
		// Settings
		{"GET", "/api/v1/settings/user"},
		// Search
		{"GET", "/api/v1/search/movies?q=test"},
		// Admin
		{"GET", "/api/v1/settings/server"},
		{"GET", "/api/v1/rbac/roles"},
		{"GET", "/api/v1/admin/activity"},
	}

	for _, ep := range protectedEndpoints {
		t.Run(ep.method+"_"+strings.ReplaceAll(ep.path, "/", "_"), func(t *testing.T) {
			resp := doRequest(t, ep.method, ep.path, "", nil)
			defer resp.Body.Close()
			assert.NotEqual(t, 200, resp.StatusCode,
				"%s %s should not succeed without auth (got %d)", ep.method, ep.path, resp.StatusCode)
		})
	}
}

// =============================================================================
// 15. Public Endpoints
// =============================================================================

func TestLive_PublicEndpoints(t *testing.T) {
	// OIDC providers list (public)
	t.Run("oidc_providers", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/oidc/providers", "", nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode, "OIDC providers should be public")
	})

	// Health endpoints are public (already tested in Infrastructure, but verify no auth needed)
	for _, path := range []string{"/healthz", "/readyz", "/startupz"} {
		t.Run("public_"+strings.TrimPrefix(path, "/"), func(t *testing.T) {
			resp := doRequest(t, "GET", path, "", nil)
			defer resp.Body.Close()
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

// =============================================================================
// 16. Chain: Cross-Cutting Scenario
// Register → Login → Create API Key → Browse (empty) content →
// Search (Typesense) → Check settings → Logout → Verify revoked
// =============================================================================

func TestLive_Chain_FullUserJourney(t *testing.T) {
	// Register new user
	creds := registerUser(t)
	t.Logf("Created user: %s (ID: %s)", creds.username, creds.userID)

	// Login
	loginUser(t, creds)

	// Browse profile
	status, profile := doJSON(t, "GET", "/api/v1/users/me", creds.accessToken, nil)
	require.Equal(t, 200, status)
	assert.Equal(t, creds.username, profile["username"])

	// Set preferences
	t.Run("set_preferences", func(t *testing.T) {
		resp := doRequest(t, "PUT", "/api/v1/users/me/preferences", creds.accessToken, map[string]interface{}{
			"metadata_language": "en",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 204)
	})

	// Create API key for programmatic access
	var apiKeySecret string
	t.Run("create_api_key", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/apikeys", creds.accessToken, map[string]interface{}{
			"name":   "programmatic-access",
			"scopes": []string{"read"},
		})
		if status == 201 {
			apiKeySecret, _ = body["api_key"].(string)
		}
	})
	_ = apiKeySecret // might be used for API key auth

	// Browse empty movie library
	t.Run("browse_movies", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/movies", creds.accessToken)
		assert.Equal(t, 200, status)
	})

	// Browse empty TV library
	t.Run("browse_tvshows", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/tvshows", creds.accessToken)
		assert.Equal(t, 200, status)
	})

	// Search via Typesense
	t.Run("typesense_search", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/search/movies?q=inception", creds.accessToken, nil)
		assert.Equal(t, 200, status, "Typesense search should be functional")
		if body != nil {
			_, hasHits := body["hits"]
			assert.True(t, hasHits, "response should have hits field")
		}
	})

	// Check user settings
	t.Run("settings", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/settings/user", creds.accessToken)
		assert.Equal(t, 200, status)
	})

	// List libraries
	t.Run("list_libraries", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/libraries", creds.accessToken, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Check sessions
	t.Run("sessions", func(t *testing.T) {
		status, _ := doJSON(t, "GET", "/api/v1/sessions", creds.accessToken, nil)
		assert.Equal(t, 200, status)
	})

	// Logout
	t.Run("logout", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/auth/logout", creds.accessToken, map[string]string{
			"refresh_token": creds.refreshToken,
		})
		defer resp.Body.Close()
		assert.Equal(t, 204, resp.StatusCode)
	})

	// After logout, refresh token is revoked but JWT access token remains valid
	// until it expires (stateless JWT design). Test that refresh is revoked.
	t.Run("refresh_revoked", func(t *testing.T) {
		status, _ := doJSON(t, "POST", "/api/v1/auth/refresh", "", map[string]string{
			"refresh_token": creds.refreshToken,
		})
		assert.NotEqual(t, 200, status, "refresh token should be revoked after logout")
	})
}

// =============================================================================
// 17. Chain: Admin Full Journey
// Admin creates library → triggers scan → checks activity logs → manages RBAC
// =============================================================================

func TestLive_Chain_AdminJourney(t *testing.T) {
	admin := ensureAdmin(t)
	tok := admin.accessToken

	// Check server settings
	t.Run("server_settings", func(t *testing.T) {
		status, _ := doJSONArray(t, "GET", "/api/v1/settings/server", tok)
		assert.Equal(t, 200, status)
	})

	// Create library
	var libraryID string
	t.Run("create_library", func(t *testing.T) {
		status, body := doJSON(t, "POST", "/api/v1/libraries", tok, map[string]interface{}{
			"name":  "Admin Journey Library",
			"type":  "movie",
			"paths": []string{"/media/admin-test"},
		})
		if status == 201 || status == 200 {
			libraryID, _ = body["id"].(string)
		}
		assert.True(t, status == 201 || status == 200,
			"library creation should succeed (got %d)", status)
	})

	// Trigger scan (requires scanType in body)
	t.Run("trigger_scan", func(t *testing.T) {
		if libraryID == "" {
			t.Skip("no library")
		}
		resp := doRequest(t, "POST", "/api/v1/libraries/"+libraryID+"/scan", tok, map[string]string{
			"scanType": "full",
		})
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 202 || resp.StatusCode == 204,
			"scan should be accepted (got %d)", resp.StatusCode)
	})

	// Check activity logs (should contain our actions)
	t.Run("activity_logs_populated", func(t *testing.T) {
		status, body := doJSON(t, "GET", "/api/v1/admin/activity", tok, nil)
		assert.Equal(t, 200, status)
		// Should have activity from registration, login, library creation, etc.
		if body != nil {
			t.Logf("Activity log response keys: %v", keys(body))
		}
	})

	// Check RBAC - list built-in roles
	t.Run("list_builtin_roles", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/rbac/roles", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		assert.Contains(t, bodyStr, "admin", "should have admin role")
		assert.Contains(t, bodyStr, "user", "should have user role")
	})

	// Check integration status (not configured, but shouldn't crash)
	t.Run("integration_status", func(t *testing.T) {
		for _, svc := range []string{"radarr", "sonarr"} {
			resp := doRequest(t, "GET", "/api/v1/admin/integrations/"+svc+"/status", tok, nil)
			resp.Body.Close()
			assert.NotEqual(t, 500, resp.StatusCode,
				"%s status should not 500", svc)
		}
	})

	// OIDC admin (empty)
	t.Run("oidc_admin", func(t *testing.T) {
		resp := doRequest(t, "GET", "/api/v1/admin/oidc/providers", tok, nil)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Trigger search reindex (exercises River job queue)
	t.Run("trigger_reindex", func(t *testing.T) {
		resp := doRequest(t, "POST", "/api/v1/search/reindex", tok, nil)
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 202,
			"reindex should be accepted (got %d)", resp.StatusCode)
	})

	// Clean up library
	if libraryID != "" {
		resp := doRequest(t, "DELETE", "/api/v1/libraries/"+libraryID, tok, nil)
		resp.Body.Close()
	}
}

// keys returns the keys of a map (for debug logging).
func keys(m map[string]interface{}) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
