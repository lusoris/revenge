package api

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupSessionTestHandler(t *testing.T) (*Handler, testutil.DB, uuid.UUID, uuid.UUID, string) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())

	// Set up session service
	repo := session.NewRepositoryPG(queries)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTExpiry:     15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}
	sessionService := session.NewService(repo, logging.NewTestLogger(), cfg)

	// Create test user
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "testuser",
		Email:    "testuser@example.com",
	})

	// Create a session for the user
	ipAddr := netip.MustParseAddr("127.0.0.1")
	deviceInfo := session.DeviceInfo{
		DeviceName: stringPtr("Test Device"),
		IPAddress:  &ipAddr,
		UserAgent:  stringPtr("Test Agent"),
	}
	_, accessToken, refreshToken, err := sessionService.CreateSession(context.Background(), user.ID, deviceInfo, []string{"read", "write"})
	require.NoError(t, err)

	// Validate the session to get session ID
	sess, err := sessionService.ValidateSession(context.Background(), accessToken)
	require.NoError(t, err)

	handler := &Handler{
		logger:         logging.NewTestLogger(),
		sessionService: sessionService,
		cfg:            cfg,
	}

	return handler, testDB, user.ID, sess.ID, refreshToken
}

func contextWithSessionID(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

// ListSessions tests

func TestHandler_ListSessions_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _, _ := setupSessionTestHandler(t)

	ctx := context.Background()

	result, err := handler.ListSessions(ctx)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_ListSessions_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.ListSessions(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.SessionListResponse)
	require.True(t, ok)
	assert.Len(t, response.Sessions, 1)
	assert.True(t, response.Sessions[0].IsActive)
}

func TestHandler_ListSessions_MultipleSessions(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	// Create another session
	ipAddr := netip.MustParseAddr("127.0.0.2")
	deviceInfo := session.DeviceInfo{
		DeviceName: stringPtr("Device 2"),
		IPAddress:  &ipAddr,
		UserAgent:  stringPtr("Agent 2"),
	}
	_, _, _, err := handler.sessionService.CreateSession(context.Background(), userID, deviceInfo, []string{"read"})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.ListSessions(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.SessionListResponse)
	require.True(t, ok)
	assert.Len(t, response.Sessions, 2)
}

// GetCurrentSession tests

func TestHandler_GetCurrentSession_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _, _ := setupSessionTestHandler(t)

	ctx := context.Background()

	result, err := handler.GetCurrentSession(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetCurrentSessionUnauthorized)
	require.True(t, ok)
}

func TestHandler_GetCurrentSession_NoSessionID(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	// User ID but no session ID in context
	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.GetCurrentSession(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetCurrentSessionUnauthorized)
	require.True(t, ok)
}

func TestHandler_GetCurrentSession_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID, sessionID, _ := setupSessionTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	ctx = contextWithSessionID(ctx, sessionID)

	result, err := handler.GetCurrentSession(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.SessionInfo)
	require.True(t, ok)
	assert.Equal(t, sessionID, response.ID)
	assert.True(t, response.IsCurrent)
	assert.True(t, response.IsActive)
}

// LogoutCurrent tests

func TestHandler_LogoutCurrent_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _, _ := setupSessionTestHandler(t)

	ctx := context.Background()

	result, err := handler.LogoutCurrent(ctx)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_LogoutCurrent_NoSessionID(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.LogoutCurrent(ctx)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_LogoutCurrent_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID, sessionID, _ := setupSessionTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	ctx = contextWithSessionID(ctx, sessionID)

	result, err := handler.LogoutCurrent(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.LogoutCurrentNoContent)
	require.True(t, ok)

	// Verify session is revoked
	sessions, err := handler.sessionService.ListUserSessions(context.Background(), userID)
	require.NoError(t, err)
	for _, s := range sessions {
		if s.ID == sessionID {
			assert.False(t, s.IsActive, "Session should be inactive after logout")
		}
	}
}

// LogoutAll tests

func TestHandler_LogoutAll_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, _, _ := setupSessionTestHandler(t)

	ctx := context.Background()

	result, err := handler.LogoutAll(ctx)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_LogoutAll_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	// Create multiple sessions
	ipAddr := netip.MustParseAddr("127.0.0.3")
	deviceInfo := session.DeviceInfo{
		DeviceName: stringPtr("Device 2"),
		IPAddress:  &ipAddr,
		UserAgent:  stringPtr("Agent 2"),
	}
	_, _, _, err := handler.sessionService.CreateSession(context.Background(), userID, deviceInfo, []string{"read"})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.LogoutAll(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.LogoutAllNoContent)
	require.True(t, ok)

	// Verify all sessions are revoked
	sessions, err := handler.sessionService.ListUserSessions(context.Background(), userID)
	require.NoError(t, err)
	for _, s := range sessions {
		assert.False(t, s.IsActive, "All sessions should be inactive after logout all")
	}
}

// RefreshSession tests

func TestHandler_RefreshSession_InvalidToken(t *testing.T) {
	t.Parallel()
	handler, _, _, _, _ := setupSessionTestHandler(t)

	ctx := context.Background()
	req := &ogen.RefreshSessionRequest{
		RefreshToken: "invalid-token",
	}

	result, err := handler.RefreshSession(ctx, req)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_RefreshSession_Success(t *testing.T) {
	t.Parallel()
	handler, _, _, _, refreshToken := setupSessionTestHandler(t)

	ctx := context.Background()
	req := &ogen.RefreshSessionRequest{
		RefreshToken: refreshToken,
	}

	result, err := handler.RefreshSession(ctx, req)
	require.NoError(t, err)

	response, ok := result.(*ogen.RefreshSessionResponse)
	require.True(t, ok)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Greater(t, response.ExpiresIn, 0)
}

// RevokeSession tests

func TestHandler_RevokeSession_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _, sessionID, _ := setupSessionTestHandler(t)

	ctx := context.Background()
	params := ogen.RevokeSessionParams{
		SessionId: sessionID,
	}

	result, err := handler.RevokeSession(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeSessionUnauthorized)
	require.True(t, ok)
}

func TestHandler_RevokeSession_NotOwner(t *testing.T) {
	t.Parallel()
	handler, testDB, _, sessionID, _ := setupSessionTestHandler(t)

	// Create another user
	otherUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "otheruser",
		Email:    "other@example.com",
	})

	// Try to revoke first user's session as other user
	ctx := contextWithUserID(context.Background(), otherUser.ID)
	params := ogen.RevokeSessionParams{
		SessionId: sessionID,
	}

	result, err := handler.RevokeSession(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeSessionNotFound)
	require.True(t, ok, "Should not be able to revoke another user's session")
}

func TestHandler_RevokeSession_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	// Create a new session to revoke
	ipAddr := netip.MustParseAddr("127.0.0.4")
	deviceInfo := session.DeviceInfo{
		DeviceName: stringPtr("To Be Revoked"),
		IPAddress:  &ipAddr,
		UserAgent:  stringPtr("Revoke Test"),
	}
	_, token, _, err := handler.sessionService.CreateSession(context.Background(), userID, deviceInfo, []string{"read"})
	require.NoError(t, err)

	// Validate to get session ID
	newSess, err := handler.sessionService.ValidateSession(context.Background(), token)
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.RevokeSessionParams{
		SessionId: newSess.ID,
	}

	result, err := handler.RevokeSession(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeSessionNoContent)
	require.True(t, ok)

	// Verify session is revoked
	sessions, err := handler.sessionService.ListUserSessions(context.Background(), userID)
	require.NoError(t, err)
	for _, s := range sessions {
		if s.ID == newSess.ID {
			assert.False(t, s.IsActive, "Session should be inactive after revocation")
		}
	}
}

func TestHandler_RevokeSession_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, userID, _, _ := setupSessionTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.RevokeSessionParams{
		SessionId: uuid.Must(uuid.NewV7()),
	}

	result, err := handler.RevokeSession(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeSessionNotFound)
	require.True(t, ok)
}
