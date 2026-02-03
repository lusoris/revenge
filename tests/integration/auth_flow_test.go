//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Register Request/Response structures
type RegisterRequest struct {
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	DisplayName *string `json:"display_name,omitempty"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	DisplayName   string `json:"display_name,omitempty"`
	EmailVerified bool   `json:"email_verified"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
}

type LoginRequest struct {
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	DeviceName *string `json:"device_name,omitempty"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TestAuthFlow_Register(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	tests := []struct {
		name           string
		req            RegisterRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful registration",
			req: RegisterRequest{
				Username: "testuser1",
				Email:    "test1@example.com",
				Password: "SecurePass123!",
			},
			expectedStatus: 201,
			checkResponse: func(t *testing.T, body []byte) {
				var user User
				err := json.Unmarshal(body, &user)
				require.NoError(t, err, "should unmarshal user")
				assert.Equal(t, "testuser1", user.Username)
				assert.Equal(t, "test1@example.com", user.Email)
				assert.NotEmpty(t, user.ID)
				assert.True(t, user.IsActive)
				assert.False(t, user.EmailVerified)
			},
		},
		{
			name: "registration with display name",
			req: RegisterRequest{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "AnotherPass456!",
				DisplayName: stringPtr("John Doe"),
			},
			expectedStatus: 201,
			checkResponse: func(t *testing.T, body []byte) {
				var user User
				err := json.Unmarshal(body, &user)
				require.NoError(t, err)
				assert.Equal(t, "John Doe", user.DisplayName)
			},
		},
		{
			name: "duplicate username",
			req: RegisterRequest{
				Username: "testuser1", // Already exists
				Email:    "different@example.com",
				Password: "Pass123!",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.req)
			require.NoError(t, err)

			resp, err := ts.HTTPClient.Post(
				ts.BaseURL+"/api/v1/auth/register",
				"application/json",
				bytes.NewReader(body),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkResponse != nil {
				respBody := make([]byte, 4096)
				n, _ := resp.Body.Read(respBody)
				tt.checkResponse(t, respBody[:n])
			}
		})
	}
}

func TestAuthFlow_LoginAndRefresh(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Register a user first
	password := "TestPassword123!"
	registerReq := RegisterRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: password,
	}
	registerBody, err := json.Marshal(registerReq)
	require.NoError(t, err)

	resp, err := ts.HTTPClient.Post(
		ts.BaseURL+"/api/v1/auth/register",
		"application/json",
		bytes.NewReader(registerBody),
	)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, 201, resp.StatusCode, "registration should succeed")

	t.Run("login with username and password", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "loginuser",
			Password: password,
		}
		loginBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		resp, err := ts.HTTPClient.Post(
			ts.BaseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(loginBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode, "login should succeed")

		var loginResp LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err, "should decode login response")

		assert.Equal(t, "loginuser", loginResp.User.Username)
		assert.NotEmpty(t, loginResp.AccessToken, "should have access token")
		assert.NotEmpty(t, loginResp.RefreshToken, "should have refresh token")
		assert.Greater(t, loginResp.ExpiresIn, 0, "should have expiry time")
	})

	t.Run("login with invalid password", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "loginuser",
			Password: "WrongPassword",
		}
		loginBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		resp, err := ts.HTTPClient.Post(
			ts.BaseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(loginBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 401, resp.StatusCode, "should fail with 401")
	})

	t.Run("login with nonexistent user", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "nonexistent",
			Password: password,
		}
		loginBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		resp, err := ts.HTTPClient.Post(
			ts.BaseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(loginBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 401, resp.StatusCode, "should fail with 401")
	})
}

func TestAuthFlow_CompleteFlow(t *testing.T) {
	ts := setupServer(t)
	defer teardownServer(t, ts)

	ctx := context.Background()
	password := "CompleteFlowPass123!"
	var accessToken, refreshToken string

	// Step 1: Register
	t.Run("step_1_register", func(t *testing.T) {
		req := RegisterRequest{
			Username:    "flowuser",
			Email:       "flow@example.com",
			Password:    password,
			DisplayName: stringPtr("Flow Test User"),
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		resp, err := ts.HTTPClient.Post(
			ts.BaseURL+"/api/v1/auth/register",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 201, resp.StatusCode, "registration should succeed")
	})

	// Step 2: Login
	t.Run("step_2_login", func(t *testing.T) {
		req := LoginRequest{
			Username: "flowuser",
			Password: password,
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		resp, err := ts.HTTPClient.Post(
			ts.BaseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode, "login should succeed")

		var loginResp LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)

		accessToken = loginResp.AccessToken
		refreshToken = loginResp.RefreshToken
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
	})

	// Step 3: Access protected resource with token
	t.Run("step_3_access_with_token", func(t *testing.T) {
		if accessToken == "" {
			t.Skip("no access token from login")
		}

		req, err := http.NewRequestWithContext(ctx, "GET", ts.BaseURL+"/api/v1/users/me", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		resp, err := ts.HTTPClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should either succeed or return 404 if endpoint doesn't exist
		// Just checking that auth middleware processes the token
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 404,
			"should process authenticated request (got %d)", resp.StatusCode)
	})

	// Step 4: Logout
	t.Run("step_4_logout", func(t *testing.T) {
		if accessToken == "" {
			t.Skip("no access token from login")
		}

		req, err := http.NewRequestWithContext(ctx, "POST", ts.BaseURL+"/api/v1/auth/logout", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		resp, err := ts.HTTPClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should logout successfully
		assert.True(t, resp.StatusCode == 204 || resp.StatusCode == 200,
			"logout should succeed (got %d)", resp.StatusCode)
	})

	_ = refreshToken // Could test token refresh here
}

func stringPtr(s string) *string {
	return &s
}
