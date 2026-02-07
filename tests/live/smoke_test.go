//go:build live

// Package live contains smoke tests that run against an already-running Revenge stack.
// Start the stack with: make docker-local
// Run these tests with: make test-live
package live

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	baseURL = envOr("REVENGE_TEST_URL", "http://localhost:8096")
	client  = &http.Client{Timeout: 10 * time.Second}
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// --- Health & Infrastructure ---

func TestLive_HealthLive(t *testing.T) {
	resp, err := client.Get(baseURL + "/healthz")
	require.NoError(t, err, "liveness endpoint should be reachable")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "healthy", result["status"])
}

func TestLive_HealthReady(t *testing.T) {
	resp, err := client.Get(baseURL + "/readyz")
	require.NoError(t, err, "readiness endpoint should be reachable")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "healthy", result["status"])
}

func TestLive_HealthStartup(t *testing.T) {
	resp, err := client.Get(baseURL + "/startupz")
	require.NoError(t, err, "startup endpoint should be reachable")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// --- Auth Flow ---

func TestLive_AuthFlow(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("smoketest_%d", time.Now().UnixNano())
	email := username + "@smoke.test"
	password := "SmokeTestPass123!"

	var accessToken, refreshToken string

	t.Run("register", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		})

		resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode, "registration should succeed")

		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)
		assert.Equal(t, username, user["username"])
		assert.NotEmpty(t, user["id"])
	})

	t.Run("login", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})

		resp, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode, "login should succeed")

		var loginResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)

		accessToken, _ = loginResp["access_token"].(string)
		refreshToken, _ = loginResp["refresh_token"].(string)
		assert.NotEmpty(t, accessToken, "should have access token")
		assert.NotEmpty(t, refreshToken, "should have refresh token")
	})

	t.Run("login_invalid_password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": username,
			"password": "WrongPassword!",
		})

		resp, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get_current_user", func(t *testing.T) {
		if accessToken == "" {
			t.Skip("no access token")
		}

		req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)
		assert.Equal(t, username, user["username"])
	})

	t.Run("refresh_token", func(t *testing.T) {
		if refreshToken == "" {
			t.Skip("no refresh token")
		}

		body, _ := json.Marshal(map[string]string{
			"refresh_token": refreshToken,
		})

		resp, err := client.Post(baseURL+"/api/v1/auth/refresh", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var refreshResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&refreshResp)
		require.NoError(t, err)
		newToken, _ := refreshResp["access_token"].(string)
		assert.NotEmpty(t, newToken, "should get new access token")

		// Update token for subsequent tests
		if newToken != "" {
			accessToken = newToken
		}
	})

	t.Run("logout", func(t *testing.T) {
		if accessToken == "" || refreshToken == "" {
			t.Skip("no tokens")
		}

		body, _ := json.Marshal(map[string]string{
			"refresh_token": refreshToken,
		})

		req, _ := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v1/auth/logout", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode,
			"logout should return 204")
	})
}

// --- Unauthenticated Access ---

func TestLive_UnauthenticatedAccess(t *testing.T) {
	protectedEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/users/me"},
		{"GET", "/api/v1/libraries"},
	}

	for _, ep := range protectedEndpoints {
		t.Run(ep.method+"_"+ep.path, func(t *testing.T) {
			req, _ := http.NewRequestWithContext(context.Background(), ep.method, baseURL+ep.path, nil)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// ogen-generated security handlers return 401 or 500 when auth header is missing
			assert.NotEqual(t, http.StatusOK, resp.StatusCode,
				"%s %s should not succeed without auth", ep.method, ep.path)
		})
	}
}

// --- Server Settings ---

func TestLive_ServerSettings(t *testing.T) {
	ctx := context.Background()
	accessToken := registerAndLogin(t)

	t.Run("get_settings", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/settings/server", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// May be 200 (if admin) or 403 (if not admin)
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden,
			"should get settings or be denied (got %d)", resp.StatusCode)
	})
}

// --- Library CRUD ---

func TestLive_Libraries(t *testing.T) {
	ctx := context.Background()
	accessToken := registerAndLogin(t)

	t.Run("list_libraries", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/libraries", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, _ := io.ReadAll(resp.Body)
		// Response should be valid JSON (array or object with items)
		assert.True(t, json.Valid(respBody), "should return valid JSON")
	})
}

// --- User Preferences ---

func TestLive_UserPreferences(t *testing.T) {
	ctx := context.Background()
	accessToken := registerAndLogin(t)

	t.Run("get_preferences", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/users/me/preferences", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// --- Helpers ---

// registerAndLogin creates a new user and returns their access token.
func registerAndLogin(t *testing.T) string {
	t.Helper()

	username := fmt.Sprintf("live_%d", time.Now().UnixNano())
	password := "LiveTestPass123!"

	// Register
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"email":    username + "@live.test",
		"password": password,
	})

	resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode, "registration should succeed")

	// Login
	body, _ = json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	resp2, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp2.Body.Close()
	require.Equal(t, http.StatusOK, resp2.StatusCode, "login should succeed")

	var loginResp map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&loginResp)
	require.NoError(t, err)

	token, _ := loginResp["access_token"].(string)
	require.NotEmpty(t, token, "should have access token")
	return token
}
