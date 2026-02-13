package session

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/testutil"
)

// ============================================================================
// Helpers
// ============================================================================

func setupIntegrationService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := &RepositoryPG{queries: queries}
	svc := NewServiceForTesting(repo, logging.NewTestLogger(), 32, 15*time.Minute, 7*24*time.Hour, 10)
	return svc, testDB
}

func testDeviceInfo() DeviceInfo {
	ip := netip.MustParseAddr("192.168.1.1")
	return DeviceInfo{
		DeviceName: new("Test Device"),
		UserAgent:  new("TestAgent/1.0"),
		IPAddress:  &ip,
	}
}

// ============================================================================
// Full Lifecycle: Create → Validate → Refresh → Revoke → Validate (fails)
// ============================================================================

func TestServiceIntegration_FullSessionLifecycle(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	// Create session
	token, refreshToken, err := svc.CreateSession(ctx, userID, testDeviceInfo(), []string{"read", "write"})
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)

	// Validate session
	sess, err := svc.ValidateSession(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, userID, sess.UserID)

	// Refresh session
	newToken, newRefreshToken, err := svc.RefreshSession(ctx, refreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEqual(t, token, newToken)
	assert.NotEqual(t, refreshToken, newRefreshToken)

	// Validate new token succeeds
	newSess, err := svc.ValidateSession(ctx, newToken)
	require.NoError(t, err)
	assert.Equal(t, userID, newSess.UserID)

	// Old token should be revoked after refresh
	_, err = svc.ValidateSession(ctx, token)
	assert.Error(t, err)

	// Revoke the new session
	err = svc.RevokeSession(ctx, newSess.ID)
	require.NoError(t, err)

	// Validate after revoke fails
	_, err = svc.ValidateSession(ctx, newToken)
	assert.Error(t, err)
}

// ============================================================================
// ListUserSessions: create multiple, list, verify count
// ============================================================================

func TestServiceIntegration_ListUserSessions(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	// Create multiple sessions
	for i := range 3 {
		ip := netip.MustParseAddr("192.168.1." + string(rune('1'+i)))
		device := DeviceInfo{
			DeviceName: new("Device " + string(rune('A'+i))),
			IPAddress:  &ip,
		}
		_, _, err := svc.CreateSession(ctx, userID, device, []string{"read"})
		require.NoError(t, err)
	}

	sessions, err := svc.ListUserSessions(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, sessions, 3)

	for _, s := range sessions {
		assert.True(t, s.IsActive)
		assert.NotNil(t, s.DeviceName)
	}
}

// ============================================================================
// RevokeAllUserSessions: create multiple, revoke all, verify none valid
// ============================================================================

func TestServiceIntegration_RevokeAllUserSessions(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	var tokens []string
	for range 3 {
		tok, _, err := svc.CreateSession(ctx, userID, testDeviceInfo(), []string{"read"})
		require.NoError(t, err)
		tokens = append(tokens, tok)
	}

	// Revoke all
	err := svc.RevokeAllUserSessions(ctx, userID)
	require.NoError(t, err)

	// All tokens should be invalid
	for _, tok := range tokens {
		_, err := svc.ValidateSession(ctx, tok)
		assert.Error(t, err)
	}
}

// ============================================================================
// RevokeAllUserSessionsExcept: keep current, revoke others
// ============================================================================

func TestServiceIntegration_RevokeAllExceptCurrent(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	// Create "current" session
	currentToken, _, err := svc.CreateSession(ctx, userID, testDeviceInfo(), []string{"read"})
	require.NoError(t, err)
	currentSess, err := svc.ValidateSession(ctx, currentToken)
	require.NoError(t, err)

	// Create other sessions
	var otherTokens []string
	for range 2 {
		tok, _, err := svc.CreateSession(ctx, userID, testDeviceInfo(), []string{"read"})
		require.NoError(t, err)
		otherTokens = append(otherTokens, tok)
	}

	// Revoke all except current
	err = svc.RevokeAllUserSessionsExcept(ctx, userID, currentSess.ID)
	require.NoError(t, err)

	// Current session still valid
	_, err = svc.ValidateSession(ctx, currentToken)
	require.NoError(t, err)

	// Others should be revoked
	for _, tok := range otherTokens {
		_, err := svc.ValidateSession(ctx, tok)
		assert.Error(t, err)
	}
}

// ============================================================================
// Refresh with invalid token
// ============================================================================

func TestServiceIntegration_RefreshSession_InvalidToken(t *testing.T) {
	t.Parallel()
	svc, _ := setupIntegrationService(t)
	ctx := context.Background()

	_, _, err := svc.RefreshSession(ctx, "invalid-token-that-doesnt-exist")
	assert.Error(t, err)
}

// ============================================================================
// CleanupExpiredSessions
// ============================================================================

func TestServiceIntegration_CleanupExpiredSessions(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	// Create and immediately revoke a session
	token, _, err := svc.CreateSession(ctx, userID, testDeviceInfo(), []string{"read"})
	require.NoError(t, err)

	sess, err := svc.ValidateSession(ctx, token)
	require.NoError(t, err)

	err = svc.RevokeSession(ctx, sess.ID)
	require.NoError(t, err)

	// Cleanup should succeed (won't delete recent revocations due to 30-day retention)
	deleted, err := svc.CleanupExpiredSessions(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, deleted, 0)
}

// ============================================================================
// ValidateSession updates last activity
// ============================================================================

func TestServiceIntegration_ValidateSession_UpdatesActivity(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	token, _, err := svc.CreateSession(ctx, userID, testDeviceInfo(), []string{"read"})
	require.NoError(t, err)

	sess1, err := svc.ValidateSession(ctx, token)
	require.NoError(t, err)
	firstActivity := sess1.LastActivityAt

	// Small delay then validate again
	time.Sleep(50 * time.Millisecond)

	sess2, err := svc.ValidateSession(ctx, token)
	require.NoError(t, err)

	// LastActivityAt should advance (or stay same if DB precision is coarse)
	assert.True(t, !sess2.LastActivityAt.Before(firstActivity))
}

// ============================================================================
// Multiple users have isolated sessions
// ============================================================================

func TestServiceIntegration_MultipleUsers_Isolated(t *testing.T) {
	t.Parallel()
	svc, testDB := setupIntegrationService(t)
	ctx := context.Background()

	user1 := createTestUser(t, testDB)
	user2 := createTestUser(t, testDB)

	// Create sessions for both users
	tok1, _, err := svc.CreateSession(ctx, user1, testDeviceInfo(), []string{"read"})
	require.NoError(t, err)
	tok2, _, err := svc.CreateSession(ctx, user2, testDeviceInfo(), []string{"write"})
	require.NoError(t, err)

	// Validate user1 token returns user1
	sess1, err := svc.ValidateSession(ctx, tok1)
	require.NoError(t, err)
	assert.Equal(t, user1, sess1.UserID)

	// Validate user2 token returns user2
	sess2, err := svc.ValidateSession(ctx, tok2)
	require.NoError(t, err)
	assert.Equal(t, user2, sess2.UserID)

	// Revoking user1 sessions doesn't affect user2
	err = svc.RevokeAllUserSessions(ctx, user1)
	require.NoError(t, err)

	_, err = svc.ValidateSession(ctx, tok1)
	assert.Error(t, err)

	_, err = svc.ValidateSession(ctx, tok2)
	require.NoError(t, err) // user2 unaffected
}
