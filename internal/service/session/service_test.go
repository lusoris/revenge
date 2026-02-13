package session

import (
	"context"
	"fmt"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/errors"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupTestService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}

	logger := logging.NewTestLogger()

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
		DeviceName: new("Test Device"),
		UserAgent:  new("Mozilla/5.0"),
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
	for range 3 {
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
	for range 3 {
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
	for range 3 {
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
	assert.GreaterOrEqual(t, count, 0)

	// Session should be deleted
	var dbCount int
	err = testDB.Pool().QueryRow(ctx, "SELECT COUNT(*) FROM shared.sessions WHERE token_hash = 'cleanup_test'").Scan(&dbCount)
	require.NoError(t, err)
	assert.Equal(t, 0, dbCount)
}

func TestService_SessionToInfo(t *testing.T) {
	t.Parallel()
	service, _ := setupTestService(t)

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	ipAddr := netip.MustParseAddr("192.168.1.1")

	session := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "test_hash",
		IpAddress:      ipAddr,
		DeviceName:     new("Test Device"),
		UserAgent:      new("Mozilla/5.0"),
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

// =====================================================
// Integration edge case tests for coverage improvement
// =====================================================

func TestService_RevokeAllUserSessionsExcept_Integration_CurrentSessionSurvives(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create 5 sessions
	tokens := make([]string, 5)
	for i := range 5 {
		token, _, err := service.CreateSession(ctx, userID, DeviceInfo{
			DeviceName: new(fmt.Sprintf("Device %d", i)),
		}, []string{"read"})
		require.NoError(t, err)
		tokens[i] = token
	}

	// Pick the third one as "current"
	currentSession, err := service.ValidateSession(ctx, tokens[2])
	require.NoError(t, err)

	// Revoke all except current
	err = service.RevokeAllUserSessionsExcept(ctx, userID, currentSession.ID)
	require.NoError(t, err)

	// Current session must still be valid
	validSession, err := service.ValidateSession(ctx, tokens[2])
	require.NoError(t, err)
	assert.Equal(t, currentSession.ID, validSession.ID)

	// All other sessions must be invalid
	for i, token := range tokens {
		if i == 2 {
			continue
		}
		_, err := service.ValidateSession(ctx, token)
		require.Error(t, err, "session %d should have been revoked", i)
		assert.ErrorIs(t, err, errors.ErrUnauthorized)
	}

	// CountActive should return 1
	count, err := service.repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestService_CountActiveUserSessions_AfterCreateRevokeCleanupCycle(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Initially 0
	count, err := service.repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create 3 sessions
	var tokens []string
	for range 3 {
		token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
		require.NoError(t, err)
		tokens = append(tokens, token)
	}

	// Count should be 3
	count, err = service.repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Revoke one session
	session, err := service.ValidateSession(ctx, tokens[0])
	require.NoError(t, err)
	err = service.RevokeSession(ctx, session.ID)
	require.NoError(t, err)

	// Count should be 2
	count, err = service.repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Revoke all remaining
	err = service.RevokeAllUserSessions(ctx, userID)
	require.NoError(t, err)

	// Count should be 0
	count, err = service.repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestService_DeleteExpiredAndRevokedSessions_CleanupCycle(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create an old expired session (91+ days ago)
	queries := db.New(testDB.Pool())
	_, err := queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:    userID,
		TokenHash: "expired_cleanup_" + uuid.Must(uuid.NewV7()).String()[:8],
		ExpiresAt: time.Now().Add(-91 * 24 * time.Hour),
	})
	require.NoError(t, err)

	// Create a revoked session and backdate the revocation to 31+ days ago
	session, err := queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:    userID,
		TokenHash: "revoked_cleanup_" + uuid.Must(uuid.NewV7()).String()[:8],
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	reason := "test cleanup"
	err = service.repo.RevokeSession(ctx, session.ID, &reason)
	require.NoError(t, err)

	_, err = testDB.Pool().Exec(ctx, "UPDATE shared.sessions SET revoked_at = $1 WHERE id = $2",
		time.Now().Add(-31*24*time.Hour), session.ID)
	require.NoError(t, err)

	// Run cleanup
	count, err := service.CleanupExpiredSessions(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 2, "should have cleaned up at least 2 sessions (1 expired + 1 old revoked)")
}

func TestService_GetInactiveSessions_AndRevokeInactive(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create a session
	token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	// Validate to set the last_activity_at
	_, err = service.ValidateSession(ctx, token)
	require.NoError(t, err)

	// Get inactive sessions since 1 hour from now (all sessions created now are "inactive" relative to that future threshold)
	inactiveSince := time.Now().Add(1 * time.Hour)
	inactiveSessions, err := service.repo.GetInactiveSessions(ctx, inactiveSince)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(inactiveSessions), 1, "should find at least our session as inactive")

	// Revoke inactive sessions
	err = service.repo.RevokeInactiveSessions(ctx, inactiveSince)
	require.NoError(t, err)

	// The session should now be revoked
	_, err = service.ValidateSession(ctx, token)
	require.Error(t, err, "inactive session should be revoked")
}

func TestService_UpdateSessionActivity_Integration(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create session
	token, _, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	// First validation sets initial activity
	session1, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	initialActivity := session1.LastActivityAt

	// Small delay to ensure timestamp difference
	time.Sleep(15 * time.Millisecond)

	// Second validation should update activity
	session2, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	assert.True(t, session2.LastActivityAt.After(initialActivity),
		"activity timestamp should advance after second validation")

	// Also test UpdateSessionActivityByTokenHash directly
	tokenHash := service.hashToken(token)
	time.Sleep(15 * time.Millisecond)
	err = service.repo.UpdateSessionActivityByTokenHash(ctx, tokenHash)
	require.NoError(t, err)

	session3, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	assert.True(t, session3.LastActivityAt.After(session2.LastActivityAt),
		"activity timestamp should advance after direct UpdateSessionActivityByTokenHash")
}

func TestService_CreateSession_MaxPerUserBoundary(t *testing.T) {
	t.Parallel()

	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}

	logger := logging.NewTestLogger()

	// Create service with very low max of 3 sessions per user
	service := &Service{
		repo:          repo,
		logger:        logger,
		tokenLength:   32,
		expiry:        24 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
		maxPerUser:    3,
	}

	ctx := context.Background()
	userID := createTestUser(t, testDB)
	deviceInfo := DeviceInfo{}

	// Create exactly maxPerUser sessions
	for range 3 {
		_, _, err := service.CreateSession(ctx, userID, deviceInfo, []string{"read"})
		require.NoError(t, err)
	}

	// Count should be exactly 3
	count, err := repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Creating one more should still work (warns but does not reject)
	_, _, err = service.CreateSession(ctx, userID, deviceInfo, []string{"read"})
	require.NoError(t, err)

	// Count should now be 4
	count, err = repo.CountActiveUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(4), count)
}

func TestService_RefreshSession_Integration_PreservesDeviceInfo(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	ipAddr := netip.MustParseAddr("10.0.0.42")
	deviceInfo := DeviceInfo{
		DeviceName: new("Firefox on Linux"),
		UserAgent:  new("Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0"),
		IPAddress:  &ipAddr,
	}

	// Create session with device info
	_, refreshToken, err := service.CreateSession(ctx, userID, deviceInfo, []string{"read", "write"})
	require.NoError(t, err)

	// Refresh the session
	newToken, _, err := service.RefreshSession(ctx, refreshToken)
	require.NoError(t, err)

	// Validate the new token and verify device info was preserved
	newSession, err := service.ValidateSession(ctx, newToken)
	require.NoError(t, err)
	assert.Equal(t, userID, newSession.UserID)
	require.NotNil(t, newSession.DeviceName)
	assert.Equal(t, "Firefox on Linux", *newSession.DeviceName)
	require.NotNil(t, newSession.UserAgent)
	assert.Contains(t, *newSession.UserAgent, "Firefox")
	assert.Equal(t, netip.MustParseAddr("10.0.0.42"), newSession.IpAddress)
	assert.ElementsMatch(t, []string{"read", "write"}, newSession.Scopes)
}

func TestService_RefreshSession_Integration_OldRefreshTokenInvalid(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	_, refreshToken, err := service.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	// First refresh succeeds
	_, newRefreshToken, err := service.RefreshSession(ctx, refreshToken)
	require.NoError(t, err)
	assert.NotEqual(t, refreshToken, newRefreshToken)

	// Old refresh token should be invalid
	_, _, err = service.RefreshSession(ctx, refreshToken)
	require.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)

	// New refresh token should still work
	_, _, err = service.RefreshSession(ctx, newRefreshToken)
	require.NoError(t, err)
}

func TestService_CachedService_Integration_RevokeInvalidatesCache(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	// Create L1 cache for testing
	testCache, err := cache.NewCache(nil, 1000, 15*time.Second)
	require.NoError(t, err)
	defer testCache.Close()

	cachedSvc := NewCachedService(service, testCache, logging.NewTestLogger(), 5*time.Minute)

	userID := createTestUser(t, testDB)

	// Create session via cached service
	token, _, err := cachedSvc.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
	require.NoError(t, err)

	// Validate to populate cache
	session, err := cachedSvc.ValidateSession(ctx, token)
	require.NoError(t, err)
	require.NotNil(t, session)

	// Wait for async cache write
	time.Sleep(200 * time.Millisecond)

	// Revoke the session (should invalidate cache)
	err = cachedSvc.RevokeSession(ctx, session.ID)
	require.NoError(t, err)

	// Validate should now fail (cache should have been invalidated)
	_, err = cachedSvc.ValidateSession(ctx, token)
	require.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
}

func TestService_CachedService_Integration_RevokeAllInvalidatesCache(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	testCache, err := cache.NewCache(nil, 1000, 15*time.Second)
	require.NoError(t, err)
	defer testCache.Close()

	cachedSvc := NewCachedService(service, testCache, logging.NewTestLogger(), 5*time.Minute)

	userID := createTestUser(t, testDB)

	// Create 3 sessions
	tokens := make([]string, 3)
	for i := range 3 {
		token, _, err := cachedSvc.CreateSession(ctx, userID, DeviceInfo{}, []string{"read"})
		require.NoError(t, err)
		tokens[i] = token
	}

	// Validate all to populate cache
	for _, token := range tokens {
		_, err := cachedSvc.ValidateSession(ctx, token)
		require.NoError(t, err)
	}

	// Wait for async cache writes
	time.Sleep(200 * time.Millisecond)

	// Revoke all user sessions
	err = cachedSvc.RevokeAllUserSessions(ctx, userID)
	require.NoError(t, err)

	// All tokens should be invalid now
	for _, token := range tokens {
		_, err := cachedSvc.ValidateSession(ctx, token)
		require.Error(t, err)
	}
}

func TestService_ListUserSessions_Integration_SessionInfoFields(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	ipAddr := netip.MustParseAddr("172.16.0.1")
	_, _, err := service.CreateSession(ctx, userID, DeviceInfo{
		DeviceName: new("Chrome on MacOS"),
		UserAgent:  new("Chrome/120"),
		IPAddress:  &ipAddr,
	}, []string{"read", "write"})
	require.NoError(t, err)

	sessions, err := service.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	require.Len(t, sessions, 1)

	info := sessions[0]
	assert.True(t, info.IsActive)
	assert.False(t, info.IsCurrent) // ListUserSessions passes isCurrent=false
	require.NotNil(t, info.DeviceName)
	assert.Equal(t, "Chrome on MacOS", *info.DeviceName)
	require.NotNil(t, info.UserAgent)
	assert.Equal(t, "Chrome/120", *info.UserAgent)
	require.NotNil(t, info.IPAddress)
	assert.Equal(t, "172.16.0.1", *info.IPAddress)
	assert.False(t, info.CreatedAt.IsZero())
	assert.False(t, info.LastActivityAt.IsZero())
	assert.False(t, info.ExpiresAt.IsZero())
	assert.True(t, info.ExpiresAt.After(time.Now()))
}

func TestService_GetSessionByID_Integration_NotFound(t *testing.T) {
	t.Parallel()
	_, testDB := setupTestService(t)
	ctx := context.Background()

	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}

	// Querying a random UUID should return nil, nil (not found)
	session, err := repo.GetSessionByID(ctx, uuid.Must(uuid.NewV7()))
	require.NoError(t, err)
	assert.Nil(t, session)
}

func TestService_GetSessionByTokenHash_Integration_NotFound(t *testing.T) {
	t.Parallel()
	_, testDB := setupTestService(t)
	ctx := context.Background()

	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}

	// Querying a non-existent token hash should return nil, nil
	session, err := repo.GetSessionByTokenHash(ctx, "nonexistent_hash_value")
	require.NoError(t, err)
	assert.Nil(t, session)
}

func TestService_GetSessionByRefreshTokenHash_Integration_NotFound(t *testing.T) {
	t.Parallel()
	_, testDB := setupTestService(t)
	ctx := context.Background()

	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}

	session, err := repo.GetSessionByRefreshTokenHash(ctx, "nonexistent_refresh_hash")
	require.NoError(t, err)
	assert.Nil(t, session)
}

func TestService_CreateSession_Integration_MinimalDeviceInfo(t *testing.T) {
	t.Parallel()
	service, testDB := setupTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Create session with nil device info fields
	token, refreshToken, err := service.CreateSession(ctx, userID, DeviceInfo{
		DeviceName: nil,
		UserAgent:  nil,
		IPAddress:  nil,
	}, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)

	// Validate and check fields
	session, err := service.ValidateSession(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, userID, session.UserID)
	assert.Nil(t, session.DeviceName)
	assert.Nil(t, session.UserAgent)
}
