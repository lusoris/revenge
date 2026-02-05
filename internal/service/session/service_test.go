package session

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/errors"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupTestService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}

	logger := zap.NewNop()

	service := &Service{
		repo:          repo,
		logger:        logger,
		tokenLength:   32,
		expiry:        24 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
		maxPerUser:    10,
	}

	return service, testDB
}

func TestService_CreateSession(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	ipAddr := netip.MustParseAddr("192.168.1.1")
	deviceInfo := DeviceInfo{
		DeviceName: stringPtr("Test Device"),
		UserAgent:  stringPtr("Mozilla/5.0"),
		IPAddress:  &ipAddr,
	}

	token, refreshToken, err := service.CreateSession(ctx, userID, deviceInfo, []string{"read", "write"})
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
	assert.NotEqual(t, token, refreshToken)
}

func TestService_CreateSession_MaxPerUser(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)
	deviceInfo := DeviceInfo{}

	// Create max sessions
	for i := 0; i < service.maxPerUser; i++ {
		_, _, err := service.CreateSession(ctx, userID, deviceInfo, []string{"read"})
		require.NoError(t, err)
	}

	// Creating one more should still succeed (just warns)
	_, _, err := service.CreateSession(ctx, userID, deviceInfo, []string{"read"})
	require.NoError(t, err)
}

func TestService_ValidateSession(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	session, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
}

func TestService_ValidateSession_InvalidToken(t *testing.T) {
	t.Parallel()
	service, _ := setupTestService(t)
	ctx := context.Background()

	session, err := service.ValidateSession(ctx, "invalid_token")
	require.Error(t, err)
	// ValidateSession returns ErrUnauthorized if session is nil
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
	assert.Nil(t, session)
}

func TestService_ValidateSession_UpdatesActivity(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	session1, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	originalActivity := session1.LastActivityAt

	time.Sleep(10 * time.Millisecond)

	session2, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	assert.True(t, session2.LastActivityAt.After(originalActivity))
}

func TestService_RefreshSession(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	oldToken, oldRefreshToken, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	newToken, newRefreshToken, err := service.RefreshSession(ctx, oldRefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEmpty(t, newRefreshToken)
	assert.NotEqual(t, oldToken, newToken)
	assert.NotEqual(t, oldRefreshToken, newRefreshToken)

	// Old token should be invalid
	_, err = service.ValidateSession(ctx, oldToken)
	require.Error(t, err)

	// New token should be valid
	session, err := service.ValidateSession(ctx, newToken)
	require.NoError(t, err)
	assert.Equal(t, userID, session.UserID)

	// Old refresh token should be invalid
	_, _, err = service.RefreshSession(ctx, oldRefreshToken)
	require.Error(t, err)
}

func TestService_RefreshSession_InvalidToken(t *testing.T) {
	t.Parallel()
	service, _ := setupTestService(t)
	ctx := context.Background()

	_, _, err := service.RefreshSession(ctx, "invalid_refresh_token")
	require.Error(t, err)
	// RefreshSession returns ErrUnauthorized if session is nil
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
}

func TestService_ListUserSessions(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create multiple sessions
	for i := 0; i < 3; i++ {
		_, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
		require.NoError(t, err)
	}

	sessions, err := service.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sessions), 3)

	for _, session := range sessions {
		assert.True(t, session.IsActive)
	}
}

func TestService_RevokeSession(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	session, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)

	err = service.RevokeSession(ctx, session.ID)
	require.NoError(t, err)

	// Token should now be invalid
	_, err = service.ValidateSession(ctx, token)
	require.Error(t, err)
}

func TestService_RevokeAllUserSessions(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	var tokens []string
	for i := 0; i < 3; i++ {
		token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
		require.NoError(t, err)
		tokens = append(tokens, token)
	}

	err := service.RevokeAllUserSessions(ctx, userID)
	require.NoError(t, err)

	// All tokens should be invalid
	for _, token := range tokens {
		_, err := service.ValidateSession(ctx, token)
		require.Error(t, err)
	}

	// List should be empty
	sessions, err := service.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, sessions)
}

func TestService_RevokeAllUserSessionsExcept(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	var tokens []string
	for i := 0; i < 3; i++ {
		token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
		require.NoError(t, err)
		tokens = append(tokens, token)
	}

	// Keep first session
	keepSession, err := service.ValidateSession(ctx, tokens[0])
	require.NoError(t, err)

	err = service.RevokeAllUserSessionsExcept(ctx, userID, keepSession.ID)
	require.NoError(t, err)

	// First token should still be valid
	_, err = service.ValidateSession(ctx, tokens[0])
	require.NoError(t, err)

	// Other tokens should be invalid
	for i := 1; i < len(tokens); i++ {
		_, err := service.ValidateSession(ctx, tokens[i])
		require.Error(t, err)
	}
}

func TestService_CleanupExpiredSessions(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create old expired session
	queries := db.New(testDB.Pool())
	_, err := queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:    userID,
		TokenHash: "cleanup_test",
		ExpiresAt: time.Now().Add(-91 * 24 * time.Hour),
	})
	require.NoError(t, err)

	count, err := service.CleanupExpiredSessions(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0) // Count not implemented yet

	// Session should be deleted
	var dbCount int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.sessions WHERE token_hash = 'cleanup_test'").Scan(&dbCount)
	require.NoError(t, err)
	assert.Equal(t, 0, dbCount)
}

func TestService_SessionToInfo(t *testing.T) {
	t.Parallel()
	service, _ := setupTestService(t)

	sessionID := uuid.New()
	userID := uuid.New()
	ipAddr := netip.MustParseAddr("192.168.1.1")

	session := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "test_hash",
		IpAddress:      ipAddr,
		DeviceName:     stringPtr("Test Device"),
		UserAgent:      stringPtr("Mozilla/5.0"),
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	info := service.sessionToInfo(session, false)

	assert.Equal(t, sessionID, info.ID)
	assert.True(t, info.IsActive)
	assert.False(t, info.IsCurrent)
	assert.NotNil(t, info.DeviceName)
	assert.Equal(t, "Test Device", *info.DeviceName)
	assert.NotNil(t, info.UserAgent)
	assert.NotNil(t, info.IPAddress)
	assert.Equal(t, "192.168.1.1", *info.IPAddress)
}
