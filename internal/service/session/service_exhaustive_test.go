package session_test

import (
	"context"
	"fmt"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/errors"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/session"
)

// setupMockService creates a service with mocked repository for unit testing
func setupMockService(t *testing.T) (*session.Service, *MockSessionRepository) {
	t.Helper()
	mockRepo := NewMockSessionRepository(t)
	logger := zap.NewNop()

	service := session.NewServiceForTesting(
		mockRepo,
		logger,
		32,             // tokenLength
		24*time.Hour,   // expiry
		7*24*time.Hour, // refreshExpiry
		10,             // maxPerUser
	)

	return service, mockRepo
}

// ========== CreateSession Tests ==========

func TestService_CreateSession_ErrorCountingSessions(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	expectedErr := fmt.Errorf("database connection error")
	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), expectedErr).
		Once()

	token, refreshToken, err := svc.CreateSession(ctx, userID, session.DeviceInfo{}, []string{"read"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count user sessions")
	assert.Empty(t, token)
	assert.Empty(t, refreshToken)
}

func TestService_CreateSession_ErrorCreatingSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(5), nil).
		Once()

	expectedErr := fmt.Errorf("unique constraint violation")
	mockRepo.EXPECT().
		CreateSession(ctx, mock.AnythingOfType("session.CreateSessionParams")).
		Return(db.SharedSession{}, expectedErr).
		Once()

	token, refreshToken, err := svc.CreateSession(ctx, userID, session.DeviceInfo{}, []string{"read"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create session")
	assert.Empty(t, token)
	assert.Empty(t, refreshToken)
}

func TestService_CreateSession_NilDeviceInfo(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return params.UserID == userID &&
				params.DeviceName == nil &&
				params.UserAgent == nil &&
				params.IPAddress == nil
		})).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7())}, nil).
		Once()

	deviceInfo := session.DeviceInfo{
		DeviceName: nil,
		UserAgent:  nil,
		IPAddress:  nil,
	}

	token, refreshToken, err := svc.CreateSession(ctx, userID, deviceInfo, []string{"read"})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
}

func TestService_CreateSession_EmptyScopes(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return len(params.Scopes) == 0
		})).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7())}, nil).
		Once()

	token, refreshToken, err := svc.CreateSession(ctx, userID, session.DeviceInfo{}, []string{})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
}

func TestService_CreateSession_NilScopes(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return params.Scopes == nil
		})).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7())}, nil).
		Once()

	token, refreshToken, err := svc.CreateSession(ctx, userID, session.DeviceInfo{}, nil)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
}

func TestService_CreateSession_MaxSessionsWarning(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Return count equal to max (should still allow creation but log warning)
	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(10), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.AnythingOfType("session.CreateSessionParams")).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7())}, nil).
		Once()

	token, refreshToken, err := svc.CreateSession(ctx, userID, session.DeviceInfo{}, []string{"read"})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
}

// ========== ValidateSession Tests ==========

func TestService_ValidateSession_ErrorGettingSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	expectedErr := fmt.Errorf("database query error")
	mockRepo.EXPECT().
		GetSessionByTokenHash(ctx, mock.AnythingOfType("string")).
		Return(nil, expectedErr).
		Once()

	sess, err := svc.ValidateSession(ctx, "some_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get session")
	assert.Nil(t, sess)
}

func TestService_ValidateSession_SessionNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	mockRepo.EXPECT().
		GetSessionByTokenHash(ctx, mock.AnythingOfType("string")).
		Return(nil, nil).
		Once()

	sess, err := svc.ValidateSession(ctx, "nonexistent_token")

	require.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
	assert.Nil(t, sess)
}

func TestService_ValidateSession_UpdateActivityError(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	validSession := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "valid_hash",
		IpAddress:      netip.MustParseAddr("127.0.0.1"),
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		GetSessionByTokenHash(ctx, mock.AnythingOfType("string")).
		Return(validSession, nil).
		Once()

	// UpdateActivity fails but shouldn't cause ValidateSession to fail
	updateErr := fmt.Errorf("database update error")
	mockRepo.EXPECT().
		UpdateSessionActivity(ctx, sessionID).
		Return(updateErr).
		Once()

	sess, err := svc.ValidateSession(ctx, "valid_token")

	require.NoError(t, err)
	assert.NotNil(t, sess)
	assert.Equal(t, userID, sess.UserID)
}

// ========== RefreshSession Tests ==========

func TestService_RefreshSession_ErrorGettingSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		GetSessionByRefreshTokenHash(ctx, mock.AnythingOfType("string")).
		Return(nil, expectedErr).
		Once()

	newToken, newRefresh, err := svc.RefreshSession(ctx, "refresh_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get session by refresh token")
	assert.Empty(t, newToken)
	assert.Empty(t, newRefresh)
}

func TestService_RefreshSession_SessionNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	mockRepo.EXPECT().
		GetSessionByRefreshTokenHash(ctx, mock.AnythingOfType("string")).
		Return(nil, nil).
		Once()

	newToken, newRefresh, err := svc.RefreshSession(ctx, "invalid_refresh_token")

	require.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
	assert.Empty(t, newToken)
	assert.Empty(t, newRefresh)
}

func TestService_RefreshSession_ErrorRevokingOldSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	validSession := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "old_hash",
		IpAddress:      netip.MustParseAddr("127.0.0.1"),
		Scopes:         []string{"read", "write"},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		GetSessionByRefreshTokenHash(ctx, mock.AnythingOfType("string")).
		Return(validSession, nil).
		Once()

	revokeErr := fmt.Errorf("database constraint error")
	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(revokeErr).
		Once()

	newToken, newRefresh, err := svc.RefreshSession(ctx, "valid_refresh_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to revoke old session")
	assert.Empty(t, newToken)
	assert.Empty(t, newRefresh)
}

func TestService_RefreshSession_ErrorCreatingNewSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	validSession := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "old_hash",
		IpAddress:      netip.MustParseAddr("127.0.0.1"),
		Scopes:         []string{"read", "write"},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		GetSessionByRefreshTokenHash(ctx, mock.AnythingOfType("string")).
		Return(validSession, nil).
		Once()

	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(nil).
		Once()

	createErr := fmt.Errorf("database insertion error")
	mockRepo.EXPECT().
		CreateSession(ctx, mock.AnythingOfType("session.CreateSessionParams")).
		Return(db.SharedSession{}, createErr).
		Once()

	newToken, newRefresh, err := svc.RefreshSession(ctx, "valid_refresh_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create refreshed session")
	assert.Empty(t, newToken)
	assert.Empty(t, newRefresh)
}

// ========== ListUserSessions Tests ==========

func TestService_ListUserSessions_ErrorFromRepository(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	expectedErr := fmt.Errorf("database query failed")
	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return(nil, expectedErr).
		Once()

	sessions, err := svc.ListUserSessions(ctx, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list user sessions")
	assert.Nil(t, sessions)
}

func TestService_ListUserSessions_EmptyList(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return([]db.SharedSession{}, nil).
		Once()

	sessions, err := svc.ListUserSessions(ctx, userID)

	require.NoError(t, err)
	assert.Empty(t, sessions)
}

// ========== RevokeSession Tests ==========

func TestService_RevokeSession_ErrorFromRepository(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(expectedErr).
		Once()

	err := svc.RevokeSession(ctx, sessionID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to revoke session")
}

func TestService_RevokeSession_NonExistentSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	// Repo returns error for non-existent session
	notFoundErr := fmt.Errorf("session not found")
	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(notFoundErr).
		Once()

	err := svc.RevokeSession(ctx, sessionID)

	require.Error(t, err)
}

// ========== RevokeAllUserSessions Tests ==========

func TestService_RevokeAllUserSessions_ErrorFromRepository(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		RevokeAllUserSessions(ctx, userID, mock.AnythingOfType("*string")).
		Return(expectedErr).
		Once()

	err := svc.RevokeAllUserSessions(ctx, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to revoke all user sessions")
}

func TestService_RevokeAllUserSessions_UserWithNoSessions(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Should succeed even if user has no sessions
	mockRepo.EXPECT().
		RevokeAllUserSessions(ctx, userID, mock.AnythingOfType("*string")).
		Return(nil).
		Once()

	err := svc.RevokeAllUserSessions(ctx, userID)

	require.NoError(t, err)
}

// ========== RevokeAllUserSessionsExcept Tests ==========

func TestService_RevokeAllUserSessionsExcept_ErrorFromRepository(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	currentSessionID := uuid.Must(uuid.NewV7())

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		RevokeAllUserSessionsExcept(ctx, userID, currentSessionID, mock.AnythingOfType("*string")).
		Return(expectedErr).
		Once()

	err := svc.RevokeAllUserSessionsExcept(ctx, userID, currentSessionID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to revoke other user sessions")
}

// ========== CleanupExpiredSessions Tests ==========

func TestService_CleanupExpiredSessions_ErrorDeletingExpired(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	expectedErr := fmt.Errorf("database delete error")
	mockRepo.EXPECT().
		DeleteExpiredSessions(ctx).
		Return(int64(0), expectedErr).
		Once()

	count, err := svc.CleanupExpiredSessions(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete expired sessions")
	assert.Equal(t, 0, count)
}

func TestService_CleanupExpiredSessions_ErrorDeletingRevoked(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	mockRepo.EXPECT().
		DeleteExpiredSessions(ctx).
		Return(int64(5), nil).
		Once()

	expectedErr := fmt.Errorf("database delete error")
	mockRepo.EXPECT().
		DeleteRevokedSessions(ctx).
		Return(int64(0), expectedErr).
		Once()

	count, err := svc.CleanupExpiredSessions(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete revoked sessions")
	assert.Equal(t, 0, count)
}

func TestService_CleanupExpiredSessions_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	mockRepo.EXPECT().
		DeleteExpiredSessions(ctx).
		Return(int64(5), nil).
		Once()

	mockRepo.EXPECT().
		DeleteRevokedSessions(ctx).
		Return(int64(3), nil).
		Once()

	count, err := svc.CleanupExpiredSessions(ctx)

	require.NoError(t, err)
	assert.Equal(t, 8, count) // 5 expired + 3 revoked
}

// ========== Additional Success Path Tests ==========

func TestService_CreateSession_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return params.UserID == userID &&
				params.TokenHash != "" &&
				params.RefreshTokenHash != nil &&
				*params.RefreshTokenHash != ""
		})).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).
		Once()

	token, refreshToken, err := svc.CreateSession(ctx, userID, session.DeviceInfo{}, []string{"read"})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
	// Token should be hex encoded (64 chars for 32 bytes)
	assert.Len(t, token, 64)
	assert.Len(t, refreshToken, 64)
}

func TestService_CreateSession_WithFullDeviceInfo(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	deviceName := "iPhone 15 Pro"
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)"
	ipAddr := netip.MustParseAddr("192.168.1.100")

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return params.UserID == userID &&
				params.DeviceName != nil && *params.DeviceName == deviceName &&
				params.UserAgent != nil && *params.UserAgent == userAgent &&
				params.IPAddress != nil && *params.IPAddress == ipAddr
		})).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).
		Once()

	deviceInfo := session.DeviceInfo{
		DeviceName: &deviceName,
		UserAgent:  &userAgent,
		IPAddress:  &ipAddr,
	}

	token, refreshToken, err := svc.CreateSession(ctx, userID, deviceInfo, []string{"read", "write"})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)
}

func TestService_ValidateSession_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	validSession := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "valid_hash",
		IpAddress:      netip.MustParseAddr("127.0.0.1"),
		Scopes:         []string{"read"},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		GetSessionByTokenHash(ctx, mock.AnythingOfType("string")).
		Return(validSession, nil).
		Once()

	mockRepo.EXPECT().
		UpdateSessionActivity(ctx, sessionID).
		Return(nil).
		Once()

	sess, err := svc.ValidateSession(ctx, "valid_token")

	require.NoError(t, err)
	assert.NotNil(t, sess)
	assert.Equal(t, userID, sess.UserID)
	assert.Equal(t, sessionID, sess.ID)
}

func TestService_RefreshSession_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	deviceName := "Test Device"
	userAgent := "Test Agent"
	ipAddr := netip.MustParseAddr("192.168.1.1")

	validSession := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "old_hash",
		IpAddress:      ipAddr,
		DeviceName:     &deviceName,
		UserAgent:      &userAgent,
		Scopes:         []string{"read", "write"},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		GetSessionByRefreshTokenHash(ctx, mock.AnythingOfType("string")).
		Return(validSession, nil).
		Once()

	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return params.UserID == userID &&
				params.TokenHash != "" &&
				params.RefreshTokenHash != nil &&
				*params.RefreshTokenHash != "" &&
				params.DeviceName != nil && *params.DeviceName == deviceName &&
				params.UserAgent != nil && *params.UserAgent == userAgent
		})).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).
		Once()

	newToken, newRefresh, err := svc.RefreshSession(ctx, "valid_refresh_token")

	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEmpty(t, newRefresh)
}

func TestService_RefreshSession_WithUnspecifiedIP(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()

	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	// Session with unspecified IP
	validSession := &db.SharedSession{
		ID:             sessionID,
		UserID:         userID,
		TokenHash:      "old_hash",
		IpAddress:      netip.Addr{}, // Unspecified
		Scopes:         []string{"read"},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		GetSessionByRefreshTokenHash(ctx, mock.AnythingOfType("string")).
		Return(validSession, nil).
		Once()

	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.AnythingOfType("session.CreateSessionParams")).
		Return(db.SharedSession{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).
		Once()

	newToken, newRefresh, err := svc.RefreshSession(ctx, "valid_refresh_token")

	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEmpty(t, newRefresh)
}

func TestService_ListUserSessions_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	deviceName := "My Device"
	userAgent := "Test Agent"
	ipAddr := netip.MustParseAddr("10.0.0.1")

	sessions := []db.SharedSession{
		{
			ID:             uuid.Must(uuid.NewV7()),
			UserID:         userID,
			TokenHash:      "hash1",
			IpAddress:      ipAddr,
			DeviceName:     &deviceName,
			UserAgent:      &userAgent,
			CreatedAt:      time.Now().Add(-1 * time.Hour),
			LastActivityAt: time.Now().Add(-30 * time.Minute),
			ExpiresAt:      time.Now().Add(23 * time.Hour),
		},
		{
			ID:             uuid.Must(uuid.NewV7()),
			UserID:         userID,
			TokenHash:      "hash2",
			IpAddress:      netip.MustParseAddr("10.0.0.2"),
			CreatedAt:      time.Now().Add(-2 * time.Hour),
			LastActivityAt: time.Now().Add(-1 * time.Hour),
			ExpiresAt:      time.Now().Add(22 * time.Hour),
		},
	}

	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return(sessions, nil).
		Once()

	result, err := svc.ListUserSessions(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, sessions[0].ID, result[0].ID)
	assert.Equal(t, deviceName, *result[0].DeviceName)
	assert.NotNil(t, result[0].IPAddress)
}

func TestService_RevokeSession_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		RevokeSession(ctx, sessionID, mock.AnythingOfType("*string")).
		Return(nil).
		Once()

	err := svc.RevokeSession(ctx, sessionID)

	require.NoError(t, err)
}

func TestService_RevokeAllUserSessionsExcept_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	currentSessionID := uuid.Must(uuid.NewV7())

	mockRepo.EXPECT().
		RevokeAllUserSessionsExcept(ctx, userID, currentSessionID, mock.AnythingOfType("*string")).
		Return(nil).
		Once()

	err := svc.RevokeAllUserSessionsExcept(ctx, userID, currentSessionID)

	require.NoError(t, err)
}

// ========== SessionInfo Conversion Tests ==========

func TestService_ListUserSessions_SessionInfoConversion(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Session with all fields populated
	deviceName := "Test Device"
	userAgent := "Mozilla/5.0"
	validSession := db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("192.168.1.1"),
		DeviceName:     &deviceName,
		UserAgent:      &userAgent,
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(23 * time.Hour),
	}

	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return([]db.SharedSession{validSession}, nil).
		Once()

	result, err := svc.ListUserSessions(ctx, userID)

	require.NoError(t, err)
	require.Len(t, result, 1)

	info := result[0]
	assert.Equal(t, validSession.ID, info.ID)
	assert.Equal(t, deviceName, *info.DeviceName)
	assert.Equal(t, userAgent, *info.UserAgent)
	assert.Equal(t, "192.168.1.1", *info.IPAddress)
	assert.Equal(t, validSession.CreatedAt, info.CreatedAt)
	assert.Equal(t, validSession.LastActivityAt, info.LastActivityAt)
	assert.Equal(t, validSession.ExpiresAt, info.ExpiresAt)
	assert.True(t, info.IsActive)
	assert.False(t, info.IsCurrent)
}

func TestService_ListUserSessions_RevokedSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Revoked session
	revokedAt := pgtype.Timestamptz{Time: time.Now().Add(-1 * time.Hour), Valid: true}
	revokedSession := db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("192.168.1.1"),
		CreatedAt:      time.Now().Add(-2 * time.Hour),
		LastActivityAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt:      time.Now().Add(22 * time.Hour),
		RevokedAt:      revokedAt,
	}

	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return([]db.SharedSession{revokedSession}, nil).
		Once()

	result, err := svc.ListUserSessions(ctx, userID)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.False(t, result[0].IsActive) // Revoked session is not active
}

func TestService_ListUserSessions_ExpiredSession(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Expired session
	expiredSession := db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("192.168.1.1"),
		CreatedAt:      time.Now().Add(-25 * time.Hour),
		LastActivityAt: time.Now().Add(-24 * time.Hour),
		ExpiresAt:      time.Now().Add(-1 * time.Hour), // Expired
	}

	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return([]db.SharedSession{expiredSession}, nil).
		Once()

	result, err := svc.ListUserSessions(ctx, userID)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.False(t, result[0].IsActive) // Expired session is not active
}

func TestService_ListUserSessions_SessionWithUnspecifiedIP(t *testing.T) {
	t.Parallel()
	svc, mockRepo := setupMockService(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Session with actual unspecified IP (0.0.0.0)
	sessionWithUnspecifiedIP := db.SharedSession{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		TokenHash:      "hash",
		IpAddress:      netip.MustParseAddr("0.0.0.0"), // Unspecified IPv4
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		LastActivityAt: time.Now(),
		ExpiresAt:      time.Now().Add(23 * time.Hour),
	}

	mockRepo.EXPECT().
		ListUserSessions(ctx, userID).
		Return([]db.SharedSession{sessionWithUnspecifiedIP}, nil).
		Once()

	result, err := svc.ListUserSessions(ctx, userID)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Nil(t, result[0].IPAddress) // Should be nil for unspecified IP (0.0.0.0)
}
