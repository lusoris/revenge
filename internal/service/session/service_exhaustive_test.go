package session_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
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
	userID := uuid.New()

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
	userID := uuid.New()

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
	userID := uuid.New()

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
		Return(db.SharedSession{ID: uuid.New()}, nil).
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
	userID := uuid.New()

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return len(params.Scopes) == 0
		})).
		Return(db.SharedSession{ID: uuid.New()}, nil).
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
	userID := uuid.New()

	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(0), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.MatchedBy(func(params session.CreateSessionParams) bool {
			return params.Scopes == nil
		})).
		Return(db.SharedSession{ID: uuid.New()}, nil).
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
	userID := uuid.New()

	// Return count equal to max (should still allow creation but log warning)
	mockRepo.EXPECT().
		CountActiveUserSessions(ctx, userID).
		Return(int64(10), nil).
		Once()

	mockRepo.EXPECT().
		CreateSession(ctx, mock.AnythingOfType("session.CreateSessionParams")).
		Return(db.SharedSession{ID: uuid.New()}, nil).
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

	sessionID := uuid.New()
	userID := uuid.New()
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

	sessionID := uuid.New()
	userID := uuid.New()
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

	sessionID := uuid.New()
	userID := uuid.New()
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
	userID := uuid.New()

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
	userID := uuid.New()

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
	sessionID := uuid.New()

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
	sessionID := uuid.New()

	// Repo returns error for non-existent session
	notFoundErr := sql.ErrNoRows
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
	userID := uuid.New()

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
	userID := uuid.New()

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
	userID := uuid.New()
	currentSessionID := uuid.New()

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
