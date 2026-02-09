package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/mfa"
)

// ============================================================================
// Internal mock implementations for same-package testing
// ============================================================================

// mockRepo is a mock implementation of the Repository interface for unit tests.
type mockRepo struct {
	mock.Mock
}

func newMockRepo(t *testing.T) *mockRepo {
	m := &mockRepo{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *mockRepo) CreateUser(ctx context.Context, params db.CreateUserParams) (db.SharedUser, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(db.SharedUser), args.Error(1)
}

func (m *mockRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*db.SharedUser, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.SharedUser), args.Error(1)
}

func (m *mockRepo) GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.SharedUser), args.Error(1)
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*db.SharedUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.SharedUser), args.Error(1)
}

func (m *mockRepo) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

func (m *mockRepo) UpdateUserEmailVerified(ctx context.Context, userID uuid.UUID, verified bool) error {
	args := m.Called(ctx, userID, verified)
	return args.Error(0)
}

func (m *mockRepo) UpdateUserLastLogin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockRepo) CreateAuthToken(ctx context.Context, params CreateAuthTokenParams) (AuthToken, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(AuthToken), args.Error(1)
}

func (m *mockRepo) GetAuthTokenByHash(ctx context.Context, tokenHash string) (AuthToken, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).(AuthToken), args.Error(1)
}

func (m *mockRepo) GetAuthTokensByUserID(ctx context.Context, userID uuid.UUID) ([]AuthToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]AuthToken), args.Error(1)
}

func (m *mockRepo) GetAuthTokensByDeviceFingerprint(ctx context.Context, userID uuid.UUID, deviceFingerprint string) ([]AuthToken, error) {
	args := m.Called(ctx, userID, deviceFingerprint)
	return args.Get(0).([]AuthToken), args.Error(1)
}

func (m *mockRepo) UpdateAuthTokenLastUsed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) RevokeAuthToken(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) RevokeAuthTokenByHash(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *mockRepo) RevokeAllUserAuthTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockRepo) RevokeAllUserAuthTokensExcept(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID) error {
	args := m.Called(ctx, userID, exceptID)
	return args.Error(0)
}

func (m *mockRepo) DeleteExpiredAuthTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockRepo) DeleteRevokedAuthTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockRepo) CountActiveAuthTokensByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) CreatePasswordResetToken(ctx context.Context, params CreatePasswordResetTokenParams) (PasswordResetToken, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(PasswordResetToken), args.Error(1)
}

func (m *mockRepo) GetPasswordResetToken(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).(PasswordResetToken), args.Error(1)
}

func (m *mockRepo) MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) InvalidateUserPasswordResetTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockRepo) DeleteExpiredPasswordResetTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockRepo) DeleteUsedPasswordResetTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockRepo) CreateEmailVerificationToken(ctx context.Context, params CreateEmailVerificationTokenParams) (EmailVerificationToken, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(EmailVerificationToken), args.Error(1)
}

func (m *mockRepo) GetEmailVerificationToken(ctx context.Context, tokenHash string) (EmailVerificationToken, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).(EmailVerificationToken), args.Error(1)
}

func (m *mockRepo) MarkEmailVerificationTokenUsed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) InvalidateUserEmailVerificationTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockRepo) InvalidateEmailVerificationTokensByEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *mockRepo) DeleteExpiredEmailVerificationTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockRepo) DeleteVerifiedEmailTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockRepo) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*db.SharedSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.SharedSession), args.Error(1)
}

func (m *mockRepo) MarkSessionMFAVerified(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *mockRepo) RecordFailedLoginAttempt(ctx context.Context, username, ipAddress string) error {
	args := m.Called(ctx, username, ipAddress)
	return args.Error(0)
}

func (m *mockRepo) CountFailedLoginAttemptsByUsername(ctx context.Context, username string, since time.Time) (int64, error) {
	args := m.Called(ctx, username, since)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) CountFailedLoginAttemptsByIP(ctx context.Context, ipAddress string, since time.Time) (int64, error) {
	args := m.Called(ctx, ipAddress, since)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) ClearFailedLoginAttemptsByUsername(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func (m *mockRepo) DeleteOldFailedLoginAttempts(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// mockTokenMgr is a mock implementation of the TokenManager interface.
type mockTokenMgr struct {
	mock.Mock
}

func newMockTokenMgr(t *testing.T) *mockTokenMgr {
	m := &mockTokenMgr{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *mockTokenMgr) GenerateAccessToken(userID uuid.UUID, username string) (string, error) {
	args := m.Called(userID, username)
	return args.String(0), args.Error(1)
}

func (m *mockTokenMgr) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockTokenMgr) ValidateAccessToken(token string) (*Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Claims), args.Error(1)
}

func (m *mockTokenMgr) HashRefreshToken(token string) string {
	args := m.Called(token)
	return args.String(0)
}

func (m *mockTokenMgr) ExtractClaims(token string) (*Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Claims), args.Error(1)
}

// ============================================================================
// Test helper: creates a Service with mocks, directly setting internal fields
// ============================================================================

type testHarness struct {
	svc      *Service
	repo     *mockRepo
	tokenMgr *mockTokenMgr
}

func newTestHarness(t *testing.T) *testHarness {
	t.Helper()

	repo := newMockRepo(t)
	tm := newMockTokenMgr(t)

	svc := NewServiceForTesting(
		nil, // pool - nil for unit tests; transaction-based methods cannot be tested
		repo,
		tm,
		activity.NewNoopLogger(),
		15*time.Minute,
		7*24*time.Hour,
	)

	return &testHarness{
		svc:      svc,
		repo:     repo,
		tokenMgr: tm,
	}
}

// newTestHarnessWithLockout creates a harness with lockout enabled.
func newTestHarnessWithLockout(t *testing.T, threshold int, window time.Duration) *testHarness {
	t.Helper()

	repo := newMockRepo(t)
	tm := newMockTokenMgr(t)

	svc := &Service{
		pool:             nil,
		repo:             repo,
		tokenManager:     tm,
		hasher:           nil, // set below
		activityLogger:   activity.NewNoopLogger(),
		emailService:     nil,
		logger:           nil,
		jwtExpiry:        15 * time.Minute,
		refreshExpiry:    7 * 24 * time.Hour,
		lockoutThreshold: threshold,
		lockoutWindow:    window,
		lockoutEnabled:   true,
	}
	// Use the NewService path to get proper logger and hasher
	real := NewService(nil, repo, tm, activity.NewNoopLogger(), nil, nil,
		15*time.Minute, 7*24*time.Hour, threshold, window, true)
	svc.hasher = real.hasher
	svc.logger = real.logger

	return &testHarness{
		svc:      svc,
		repo:     repo,
		tokenMgr: tm,
	}
}

// ============================================================================
// NewService Tests
// ============================================================================

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("with nil logger uses default", func(t *testing.T) {
		t.Parallel()
		svc := NewService(nil, nil, nil, nil, nil, nil,
			15*time.Minute, 7*24*time.Hour, 5, 15*time.Minute, true)
		require.NotNil(t, svc)
		assert.NotNil(t, svc.logger)
		assert.Equal(t, 15*time.Minute, svc.jwtExpiry)
		assert.Equal(t, 7*24*time.Hour, svc.refreshExpiry)
		assert.Equal(t, 5, svc.lockoutThreshold)
		assert.Equal(t, 15*time.Minute, svc.lockoutWindow)
		assert.True(t, svc.lockoutEnabled)
		assert.NotNil(t, svc.hasher)
	})

	t.Run("with provided logger", func(t *testing.T) {
		t.Parallel()
		logger := testLogger()
		svc := NewService(nil, nil, nil, nil, nil, logger,
			30*time.Minute, 14*24*time.Hour, 10, 30*time.Minute, false)
		require.NotNil(t, svc)
		assert.Equal(t, 30*time.Minute, svc.jwtExpiry)
		assert.Equal(t, 14*24*time.Hour, svc.refreshExpiry)
		assert.Equal(t, 10, svc.lockoutThreshold)
		assert.False(t, svc.lockoutEnabled)
	})
}

// ============================================================================
// ptrToString Tests
// ============================================================================

func TestPtrToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer returns empty string",
			input:    nil,
			expected: "",
		},
		{
			name:     "non-nil pointer returns value",
			input:    strPtr("hello"),
			expected: "hello",
		},
		{
			name:     "empty string pointer returns empty string",
			input:    strPtr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, ptrToString(tt.input))
		})
	}
}

// ============================================================================
// Login Tests - Lockout Scenarios
// ============================================================================

func TestLogin_AccountLocked(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	username := "lockeduser"
	ip := netip.MustParseAddr("192.168.1.1")

	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(5), nil).Once()

	resp, err := h.svc.Login(ctx, username, "password", &ip, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account locked")
	assert.Nil(t, resp)
}

func TestLogin_LockoutCheckError_ContinuesLogin(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	username := "testuser"

	// Lockout check fails
	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(0), fmt.Errorf("db error")).Once()

	// Login continues: user not found
	h.repo.On("GetUserByUsername", ctx, username).
		Return(nil, fmt.Errorf("not found")).Once()
	h.repo.On("GetUserByEmail", ctx, username).
		Return(nil, fmt.Errorf("not found")).Once()

	resp, err := h.svc.Login(ctx, username, "password", nil, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
}

func TestLogin_LockoutBelowThreshold_ContinuesLogin(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	username := "testuser"

	// Below threshold
	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(3), nil).Once()

	// User not found
	h.repo.On("GetUserByUsername", ctx, username).
		Return(nil, fmt.Errorf("not found")).Once()
	h.repo.On("GetUserByEmail", ctx, username).
		Return(nil, fmt.Errorf("not found")).Once()

	resp, err := h.svc.Login(ctx, username, "password", nil, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
}

func TestLogin_WrongPassword_RecordsFailedAttemptWithLockout(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	username := "testuser"
	ip := netip.MustParseAddr("10.0.0.1")

	// Below threshold
	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(0), nil).Once()

	// User found with known password hash
	user := &db.SharedUser{
		ID:           uuid.Must(uuid.NewV7()),
		Username:     username,
		PasswordHash: dummyPasswordHash, // will not match "wrongpassword"
	}
	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()

	// Failed login attempt recorded
	h.repo.On("RecordFailedLoginAttempt", ctx, username, ip.String()).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, username, "wrongpassword", &ip, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
}

func TestLogin_WrongPassword_RecordFailedAttemptError(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	username := "testuser"
	ip := netip.MustParseAddr("10.0.0.1")

	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(0), nil).Once()

	// User not found - use dummy hash
	h.repo.On("GetUserByUsername", ctx, username).
		Return(nil, fmt.Errorf("not found")).Once()
	h.repo.On("GetUserByEmail", ctx, username).
		Return(nil, fmt.Errorf("not found")).Once()

	// Record failed attempt fails - should not affect result
	h.repo.On("RecordFailedLoginAttempt", ctx, username, ip.String()).
		Return(fmt.Errorf("db error")).Once()

	resp, err := h.svc.Login(ctx, username, "wrongpassword", &ip, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
}

// ============================================================================
// Login Tests - Successful Login Path
// ============================================================================

func TestLogin_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	// Hash the password for verification
	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()

	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("access-token-123", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("refresh-token-456", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "refresh-token-456").
		Return("hashed-refresh").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "access-token-123", resp.AccessToken)
	assert.Equal(t, "refresh-token-456", resp.RefreshToken)
	assert.Equal(t, userID, resp.User.ID)
	assert.Equal(t, int64(h.svc.jwtExpiry.Seconds()), resp.ExpiresIn)
}

func TestLogin_SuccessWithIPAndUserAgent(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}
	ip := netip.MustParseAddr("192.168.1.100")
	ua := "Mozilla/5.0"
	device := "Chrome Desktop"
	fingerprint := "abc123"

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, username, password, &ip, &ua, &device, &fingerprint)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "at", resp.AccessToken)
}

func TestLogin_SuccessWithLockout_ClearsAttempts(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(0), nil).Once()
	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("ClearFailedLoginAttemptsByUsername", ctx, username).
		Return(nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestLogin_SuccessWithLockout_ClearAttemptsError(t *testing.T) {
	t.Parallel()
	h := newTestHarnessWithLockout(t, 5, 15*time.Minute)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("CountFailedLoginAttemptsByUsername", ctx, username, mock.AnythingOfType("time.Time")).
		Return(int64(0), nil).Once()
	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	// Clear fails but should not fail login
	h.repo.On("ClearFailedLoginAttemptsByUsername", ctx, username).
		Return(fmt.Errorf("db error")).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestLogin_FoundByEmail(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	email := "test@example.com"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     "testuser",
		Email:        email,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	// Not found by username, found by email
	h.repo.On("GetUserByUsername", ctx, email).
		Return(nil, fmt.Errorf("not found")).Once()
	h.repo.On("GetUserByEmail", ctx, email).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "testuser").
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, email, password, nil, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, userID, resp.User.ID)
}

func TestLogin_AccountDisabled(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "disableduser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := false
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "account is disabled", err.Error())
	assert.Nil(t, resp)
}

func TestLogin_GenerateAccessTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("", fmt.Errorf("signing error")).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate access token")
	assert.Nil(t, resp)
}

func TestLogin_GenerateRefreshTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("", fmt.Errorf("random gen error")).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate refresh token")
	assert.Nil(t, resp)
}

func TestLogin_StoreRefreshTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{}, fmt.Errorf("db error")).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store refresh token")
	assert.Nil(t, resp)
}

func TestLogin_UpdateLastLoginError_StillSucceeds(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	// Fails but should not fail login
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(fmt.Errorf("db error")).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestLogin_NilIsActive_TreatedAsActive(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	// IsActive is nil (default)
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     nil,
	}

	h.repo.On("GetUserByUsername", ctx, username).
		Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

// ============================================================================
// Logout Tests
// ============================================================================

func TestLogout_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	h.tokenMgr.On("HashRefreshToken", "token123").Return("hash123").Once()
	h.repo.On("RevokeAuthTokenByHash", ctx, "hash123").Return(nil).Once()

	err := h.svc.Logout(ctx, "token123")
	require.NoError(t, err)
}

func TestLogout_Error(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	h.tokenMgr.On("HashRefreshToken", "token123").Return("hash123").Once()
	h.repo.On("RevokeAuthTokenByHash", ctx, "hash123").Return(fmt.Errorf("db error")).Once()

	err := h.svc.Logout(ctx, "token123")
	require.Error(t, err)
}

// ============================================================================
// LogoutAll Tests
// ============================================================================

func TestLogoutAll_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	h.repo.On("RevokeAllUserAuthTokens", ctx, userID).Return(nil).Once()

	err := h.svc.LogoutAll(ctx, userID)
	require.NoError(t, err)
}

func TestLogoutAll_Error(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	h.repo.On("RevokeAllUserAuthTokens", ctx, userID).Return(fmt.Errorf("db error")).Once()

	err := h.svc.LogoutAll(ctx, userID)
	require.Error(t, err)
}

// ============================================================================
// RefreshToken Tests
// ============================================================================

func TestRefreshToken_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	tokenID := uuid.Must(uuid.NewV7())
	authToken := AuthToken{ID: tokenID, UserID: userID}
	user := &db.SharedUser{ID: userID, Username: "testuser"}

	h.tokenMgr.On("HashRefreshToken", "rt").Return("rth").Once()
	h.repo.On("GetAuthTokenByHash", ctx, "rth").Return(authToken, nil).Once()
	h.repo.On("GetUserByID", ctx, userID).Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "testuser").Return("new-at", nil).Once()
	h.repo.On("UpdateAuthTokenLastUsed", ctx, tokenID).Return(nil).Once()

	resp, err := h.svc.RefreshToken(ctx, "rt")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-at", resp.AccessToken)
	assert.Equal(t, "rt", resp.RefreshToken) // same refresh token returned
	assert.Equal(t, "testuser", resp.User.Username)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	h.tokenMgr.On("HashRefreshToken", "bad").Return("bad-hash").Once()
	h.repo.On("GetAuthTokenByHash", ctx, "bad-hash").Return(AuthToken{}, fmt.Errorf("not found")).Once()

	resp, err := h.svc.RefreshToken(ctx, "bad")

	require.Error(t, err)
	assert.Equal(t, "invalid or expired refresh token", err.Error())
	assert.Nil(t, resp)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	authToken := AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}

	h.tokenMgr.On("HashRefreshToken", "rt").Return("rth").Once()
	h.repo.On("GetAuthTokenByHash", ctx, "rth").Return(authToken, nil).Once()
	h.repo.On("GetUserByID", ctx, userID).Return(nil, fmt.Errorf("not found")).Once()

	resp, err := h.svc.RefreshToken(ctx, "rt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, resp)
}

func TestRefreshToken_GenerateAccessTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	authToken := AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}
	user := &db.SharedUser{ID: userID, Username: "testuser"}

	h.tokenMgr.On("HashRefreshToken", "rt").Return("rth").Once()
	h.repo.On("GetAuthTokenByHash", ctx, "rth").Return(authToken, nil).Once()
	h.repo.On("GetUserByID", ctx, userID).Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "testuser").Return("", fmt.Errorf("err")).Once()

	resp, err := h.svc.RefreshToken(ctx, "rt")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate access token")
	assert.Nil(t, resp)
}

func TestRefreshToken_UpdateLastUsedError_StillSucceeds(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	tokenID := uuid.Must(uuid.NewV7())
	authToken := AuthToken{ID: tokenID, UserID: userID}
	user := &db.SharedUser{ID: userID, Username: "testuser"}

	h.tokenMgr.On("HashRefreshToken", "rt").Return("rth").Once()
	h.repo.On("GetAuthTokenByHash", ctx, "rth").Return(authToken, nil).Once()
	h.repo.On("GetUserByID", ctx, userID).Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "testuser").Return("new-at", nil).Once()
	h.repo.On("UpdateAuthTokenLastUsed", ctx, tokenID).Return(fmt.Errorf("db err")).Once()

	resp, err := h.svc.RefreshToken(ctx, "rt")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-at", resp.AccessToken)
}

// ============================================================================
// ResendVerification Tests
// ============================================================================

func TestResendVerification_InvalidateError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	h.repo.On("InvalidateUserEmailVerificationTokens", ctx, userID).
		Return(fmt.Errorf("db error")).Once()

	err := h.svc.ResendVerification(ctx, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to invalidate old tokens")
}

func TestResendVerification_UserNotFound(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	h.repo.On("InvalidateUserEmailVerificationTokens", ctx, userID).
		Return(nil).Once()
	h.repo.On("GetUserByID", ctx, userID).
		Return(nil, fmt.Errorf("not found")).Once()

	err := h.svc.ResendVerification(ctx, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestResendVerification_CreateTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	user := &db.SharedUser{ID: userID, Email: "u@e.com", Username: "u"}

	h.repo.On("InvalidateUserEmailVerificationTokens", ctx, userID).
		Return(nil).Once()
	h.repo.On("GetUserByID", ctx, userID).
		Return(user, nil).Once()
	h.tokenMgr.On("HashRefreshToken", mock.AnythingOfType("string")).
		Return("hash").Once()
	h.repo.On("CreateEmailVerificationToken", ctx, mock.AnythingOfType("auth.CreateEmailVerificationTokenParams")).
		Return(EmailVerificationToken{}, fmt.Errorf("db error")).Once()

	err := h.svc.ResendVerification(ctx, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create verification token")
}

func TestResendVerification_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	user := &db.SharedUser{ID: userID, Email: "u@e.com", Username: "u"}

	h.repo.On("InvalidateUserEmailVerificationTokens", ctx, userID).
		Return(nil).Once()
	h.repo.On("GetUserByID", ctx, userID).
		Return(user, nil).Once()
	h.tokenMgr.On("HashRefreshToken", mock.AnythingOfType("string")).
		Return("hash").Once()
	h.repo.On("CreateEmailVerificationToken", ctx, mock.AnythingOfType("auth.CreateEmailVerificationTokenParams")).
		Return(EmailVerificationToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()

	err := h.svc.ResendVerification(ctx, userID)
	require.NoError(t, err)
}

// ============================================================================
// Register Tests - Pre-transaction error paths
// ============================================================================

func TestRegister_EmptyPassword_HashFails(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	req := RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "", // empty password should fail hasher
	}

	user, err := h.svc.Register(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to hash password")
	assert.Nil(t, user)
}

// ============================================================================
// VerifyEmail Tests
// ============================================================================

func TestVerifyEmail_InvalidToken(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	h.tokenMgr.On("HashRefreshToken", "bad-token").Return("bad-hash").Once()
	h.repo.On("GetEmailVerificationToken", ctx, "bad-hash").
		Return(EmailVerificationToken{}, fmt.Errorf("not found")).Once()

	err := h.svc.VerifyEmail(ctx, "bad-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired verification token")
}

// VerifyEmail success path requires transactions (pool.Begin) and is covered by integration tests.

// ============================================================================
// RegisterFromOIDC Tests
// ============================================================================

func TestRegisterFromOIDC_CreateUserError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	req := RegisterFromOIDCRequest{
		Username: "oidcuser",
		Email:    "oidc@example.com",
	}

	h.repo.On("CreateUser", ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(db.SharedUser{}, fmt.Errorf("unique constraint")).Once()

	user, err := h.svc.RegisterFromOIDC(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.Nil(t, user)
}

func TestRegisterFromOIDC_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	displayName := "OIDC User"
	req := RegisterFromOIDCRequest{
		Username:    "oidcuser",
		Email:       "oidc@example.com",
		DisplayName: &displayName,
	}

	createdUser := db.SharedUser{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	h.repo.On("CreateUser", ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(createdUser, nil).Once()
	h.repo.On("UpdateUserEmailVerified", ctx, userID, true).
		Return(nil).Once()

	user, err := h.svc.RegisterFromOIDC(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
}

func TestRegisterFromOIDC_EmailVerificationError_StillSucceeds(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	req := RegisterFromOIDCRequest{
		Username: "oidcuser",
		Email:    "oidc@example.com",
	}

	createdUser := db.SharedUser{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	h.repo.On("CreateUser", ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(createdUser, nil).Once()
	h.repo.On("UpdateUserEmailVerified", ctx, userID, true).
		Return(fmt.Errorf("db error")).Once()

	user, err := h.svc.RegisterFromOIDC(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, user)
}

// ============================================================================
// CreateSessionForUser Tests
// ============================================================================

func TestCreateSessionForUser_UserNotFound(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	h.repo.On("GetUserByID", ctx, userID).
		Return(nil, fmt.Errorf("not found")).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, resp)
}

func TestCreateSessionForUser_AccountDisabled(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := false

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.Error(t, err)
	assert.Equal(t, "account is disabled", err.Error())
	assert.Nil(t, resp)
}

func TestCreateSessionForUser_GenerateAccessTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := true

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "u").
		Return("", fmt.Errorf("err")).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate access token")
	assert.Nil(t, resp)
}

func TestCreateSessionForUser_GenerateRefreshTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := true

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "u").
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("", fmt.Errorf("err")).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate refresh token")
	assert.Nil(t, resp)
}

func TestCreateSessionForUser_StoreTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := true

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "u").
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{}, fmt.Errorf("db error")).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store refresh token")
	assert.Nil(t, resp)
}

func TestCreateSessionForUser_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := true

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "u").
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "at", resp.AccessToken)
	assert.Equal(t, "rt", resp.RefreshToken)
}

func TestCreateSessionForUser_WithIPAddress(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := true
	ip := netip.MustParseAddr("10.0.0.1")
	ua := "TestAgent"
	device := "TestDevice"

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "u").
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(nil).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, &ip, &ua, &device)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCreateSessionForUser_LastLoginError_StillSucceeds(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	isActive := true

	h.repo.On("GetUserByID", ctx, userID).
		Return(&db.SharedUser{ID: userID, Username: "u", IsActive: &isActive}, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, "u").
		Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").
		Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").
		Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).
		Return(fmt.Errorf("db err")).Once()

	resp, err := h.svc.CreateSessionForUser(ctx, userID, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// ============================================================================
// ChangePassword Tests
// ============================================================================

func TestChangePassword_UserNotFound(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	h.repo.On("GetUserByID", ctx, userID).
		Return(nil, fmt.Errorf("not found")).Once()

	err := h.svc.ChangePassword(ctx, userID, "old", "new")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	// Hash a known password
	correctHash, err := h.svc.hasher.HashPassword("correct-password")
	require.NoError(t, err)

	user := &db.SharedUser{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: correctHash,
	}
	h.repo.On("GetUserByID", ctx, userID).Return(user, nil).Once()

	err = h.svc.ChangePassword(ctx, userID, "wrong-password", "new-password")
	require.Error(t, err)
	assert.Equal(t, "invalid current password", err.Error())
}

func TestChangePassword_EmptyOldPassword(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	correctHash, err := h.svc.hasher.HashPassword("correct-password")
	require.NoError(t, err)

	user := &db.SharedUser{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: correctHash,
	}
	h.repo.On("GetUserByID", ctx, userID).Return(user, nil).Once()

	// Empty old password should fail hasher.VerifyPassword
	err = h.svc.ChangePassword(ctx, userID, "", "new-password")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password verification failed")
}

// ChangePassword success path requires transactions (pool.Begin) and is covered by integration tests.

// ============================================================================
// RequestPasswordReset Tests
// ============================================================================

func TestRequestPasswordReset_UserNotFound_SilentSuccess(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	h.repo.On("GetUserByEmail", ctx, "unknown@example.com").
		Return(nil, fmt.Errorf("not found")).Once()

	err := h.svc.RequestPasswordReset(ctx, "unknown@example.com", nil, nil)
	require.NoError(t, err) // should NOT reveal that email doesn't exist
}

func TestRequestPasswordReset_InvalidateOldTokensError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	user := &db.SharedUser{ID: userID, Email: "u@e.com"}

	h.repo.On("GetUserByEmail", ctx, "u@e.com").Return(user, nil).Once()
	h.repo.On("InvalidateUserPasswordResetTokens", ctx, userID).
		Return(fmt.Errorf("db error")).Once()

	err := h.svc.RequestPasswordReset(ctx, "u@e.com", nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to invalidate old tokens")
}

func TestRequestPasswordReset_CreateTokenError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	user := &db.SharedUser{ID: userID, Email: "u@e.com"}

	h.repo.On("GetUserByEmail", ctx, "u@e.com").Return(user, nil).Once()
	h.repo.On("InvalidateUserPasswordResetTokens", ctx, userID).Return(nil).Once()
	h.tokenMgr.On("HashRefreshToken", mock.AnythingOfType("string")).Return("hash").Once()
	h.repo.On("CreatePasswordResetToken", ctx, mock.AnythingOfType("auth.CreatePasswordResetTokenParams")).
		Return(PasswordResetToken{}, fmt.Errorf("db error")).Once()

	err := h.svc.RequestPasswordReset(ctx, "u@e.com", nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create reset token")
}

func TestRequestPasswordReset_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	user := &db.SharedUser{ID: userID, Email: "u@e.com", Username: "u"}

	h.repo.On("GetUserByEmail", ctx, "u@e.com").Return(user, nil).Once()
	h.repo.On("InvalidateUserPasswordResetTokens", ctx, userID).Return(nil).Once()
	h.tokenMgr.On("HashRefreshToken", mock.AnythingOfType("string")).Return("hash").Once()
	h.repo.On("CreatePasswordResetToken", ctx, mock.AnythingOfType("auth.CreatePasswordResetTokenParams")).
		Return(PasswordResetToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()

	err := h.svc.RequestPasswordReset(ctx, "u@e.com", nil, nil)
	require.NoError(t, err)
}

func TestRequestPasswordReset_WithIPAndUserAgent(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	user := &db.SharedUser{ID: userID, Email: "u@e.com", Username: "u"}
	ip := netip.MustParseAddr("10.0.0.1")
	ua := "TestAgent"

	h.repo.On("GetUserByEmail", ctx, "u@e.com").Return(user, nil).Once()
	h.repo.On("InvalidateUserPasswordResetTokens", ctx, userID).Return(nil).Once()
	h.tokenMgr.On("HashRefreshToken", mock.AnythingOfType("string")).Return("hash").Once()
	h.repo.On("CreatePasswordResetToken", ctx, mock.AnythingOfType("auth.CreatePasswordResetTokenParams")).
		Return(PasswordResetToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()

	err := h.svc.RequestPasswordReset(ctx, "u@e.com", &ip, &ua)
	require.NoError(t, err)
}

// ============================================================================
// ResetPassword Tests
// ============================================================================

func TestResetPassword_InvalidToken(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	h.tokenMgr.On("HashRefreshToken", "bad").Return("bad-hash").Once()
	h.repo.On("GetPasswordResetToken", ctx, "bad-hash").
		Return(PasswordResetToken{}, fmt.Errorf("not found")).Once()

	err := h.svc.ResetPassword(ctx, "bad", "newpass")
	require.Error(t, err)
	assert.Equal(t, "invalid or expired reset token", err.Error())
}

func TestResetPassword_EmptyNewPassword_HashFails(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	resetTokenID := uuid.Must(uuid.NewV7())
	resetToken := PasswordResetToken{
		ID:     resetTokenID,
		UserID: userID,
	}

	h.tokenMgr.On("HashRefreshToken", "valid-token").Return("valid-hash").Once()
	h.repo.On("GetPasswordResetToken", ctx, "valid-hash").Return(resetToken, nil).Once()

	// Empty password should fail the hasher
	err := h.svc.ResetPassword(ctx, "valid-token", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to hash password")
}

// ResetPassword success path requires transactions (pool.Begin) and is covered by integration tests.

// ============================================================================
// MFA Integration Tests
// ============================================================================

func TestNewMFAAuthenticator(t *testing.T) {
	t.Parallel()
	auth := NewMFAAuthenticator(nil)
	require.NotNil(t, auth)
	assert.Nil(t, auth.mfaManager)
}

func TestCompleteMFALogin_FailedVerification(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	result := &mfa.VerificationResult{
		Success: false,
		Method:  mfa.VerifyMethodTOTP,
	}

	err := h.svc.CompleteMFALogin(ctx, sessionID, result)
	require.Error(t, err)
	assert.Equal(t, ErrInvalidMFACode, err)
}

func TestCompleteMFALogin_Success(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	result := &mfa.VerificationResult{
		Success: true,
		Method:  mfa.VerifyMethodTOTP,
		UserID:  userID,
	}

	h.repo.On("MarkSessionMFAVerified", ctx, sessionID).Return(nil).Once()

	err := h.svc.CompleteMFALogin(ctx, sessionID, result)
	require.NoError(t, err)
}

func TestCompleteMFALogin_MarkSessionError(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	result := &mfa.VerificationResult{
		Success: true,
		Method:  mfa.VerifyMethodTOTP,
		UserID:  userID,
	}

	h.repo.On("MarkSessionMFAVerified", ctx, sessionID).
		Return(fmt.Errorf("db error")).Once()

	err := h.svc.CompleteMFALogin(ctx, sessionID, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark session as MFA verified")
}

func TestGetSessionMFAInfo_SessionNotFound(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	h.repo.On("GetSessionByID", ctx, sessionID).
		Return(nil, fmt.Errorf("not found")).Once()

	info, err := h.svc.GetSessionMFAInfo(ctx, sessionID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
	assert.Nil(t, info)
}

func TestGetSessionMFAInfo_NotVerified(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	session := &db.SharedSession{
		ID:          sessionID,
		MfaVerified: false,
		MfaVerifiedAt: pgtype.Timestamptz{
			Valid: false,
		},
	}

	h.repo.On("GetSessionByID", ctx, sessionID).Return(session, nil).Once()

	info, err := h.svc.GetSessionMFAInfo(ctx, sessionID)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.False(t, info.MFAVerified)
	assert.Nil(t, info.MFAVerifiedAt)
}

func TestGetSessionMFAInfo_Verified(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()
	sessionID := uuid.Must(uuid.NewV7())

	verifiedAt := time.Now().Add(-10 * time.Minute)
	session := &db.SharedSession{
		ID:          sessionID,
		MfaVerified: true,
		MfaVerifiedAt: pgtype.Timestamptz{
			Time:  verifiedAt,
			Valid: true,
		},
	}

	h.repo.On("GetSessionByID", ctx, sessionID).Return(session, nil).Once()

	info, err := h.svc.GetSessionMFAInfo(ctx, sessionID)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.True(t, info.MFAVerified)
	require.NotNil(t, info.MFAVerifiedAt)
	assert.WithinDuration(t, verifiedAt, *info.MFAVerifiedAt, time.Second)
}

// ============================================================================
// MFAAuthenticator.VerifyMFA Tests
// ============================================================================

func TestVerifyMFA_UnsupportedMethod(t *testing.T) {
	t.Parallel()
	auth := NewMFAAuthenticator(nil)

	_, err := auth.VerifyMFA(context.Background(), MFAVerifyRequest{
		Method: "sms",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported MFA method: sms")
}

func TestVerifyMFA_TOTPEmptyCode(t *testing.T) {
	t.Parallel()
	auth := NewMFAAuthenticator(nil)

	_, err := auth.VerifyMFA(context.Background(), MFAVerifyRequest{
		Method: "totp",
		Code:   "",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "totp code is required")
}

func TestVerifyMFA_BackupCodeEmptyCode(t *testing.T) {
	t.Parallel()
	auth := NewMFAAuthenticator(nil)

	_, err := auth.VerifyMFA(context.Background(), MFAVerifyRequest{
		Method: "backup_code",
		Code:   "",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "backup code is required")
}

func TestVerifyMFA_WebAuthn_NilAssertion(t *testing.T) {
	t.Parallel()
	auth := NewMFAAuthenticator(nil)

	_, err := auth.VerifyMFA(context.Background(), MFAVerifyRequest{
		Method: "webauthn",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "webauthn assertion data is required")
}

// ============================================================================
// LoginWithMFA Tests
// ============================================================================

// LoginWithMFA delegates to Login, which uses password hashing.
// We test that the nil MFA authenticator case works and returns login response.
func TestLoginWithMFA_NilAuthenticator_ReturnsLoginResponse(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	userID := uuid.Must(uuid.NewV7())
	username := "testuser"
	password := "SecureP@ssw0rd!"

	passwordHash, err := h.svc.hasher.HashPassword(password)
	require.NoError(t, err)

	isActive := true
	user := &db.SharedUser{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     &isActive,
	}

	h.repo.On("GetUserByUsername", ctx, username).Return(user, nil).Once()
	h.tokenMgr.On("GenerateAccessToken", userID, username).Return("at", nil).Once()
	h.tokenMgr.On("GenerateRefreshToken").Return("rt", nil).Once()
	h.tokenMgr.On("HashRefreshToken", "rt").Return("rth").Once()
	h.repo.On("CreateAuthToken", ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(AuthToken{ID: uuid.Must(uuid.NewV7()), UserID: userID}, nil).Once()
	h.repo.On("UpdateUserLastLogin", ctx, userID).Return(nil).Once()

	loginResp, mfaResp, err := h.svc.LoginWithMFA(ctx, username, password, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, loginResp)
	assert.Nil(t, mfaResp)
	assert.Equal(t, "at", loginResp.AccessToken)
}

func TestLoginWithMFA_LoginFails(t *testing.T) {
	t.Parallel()
	h := newTestHarness(t)
	ctx := context.Background()

	// User not found
	h.repo.On("GetUserByUsername", ctx, "nouser").
		Return(nil, fmt.Errorf("not found")).Once()
	h.repo.On("GetUserByEmail", ctx, "nouser").
		Return(nil, fmt.Errorf("not found")).Once()

	loginResp, mfaResp, err := h.svc.LoginWithMFA(ctx, "nouser", "pass", nil, nil, nil, nil, nil)
	require.Error(t, err)
	assert.Nil(t, loginResp)
	assert.Nil(t, mfaResp)
}

// ============================================================================
// MFA Error Sentinel Values
// ============================================================================

func TestMFAErrorValues(t *testing.T) {
	t.Parallel()

	assert.True(t, errors.Is(ErrMFARequired, ErrMFARequired))
	assert.True(t, errors.Is(ErrInvalidMFACode, ErrInvalidMFACode))
	assert.True(t, errors.Is(ErrMFANotEnabled, ErrMFANotEnabled))

	assert.Equal(t, "mfa verification required", ErrMFARequired.Error())
	assert.Equal(t, "invalid mfa code", ErrInvalidMFACode.Error())
	assert.Equal(t, "mfa not enabled for user", ErrMFANotEnabled.Error())
}

// ============================================================================
// Repository PG - Conversion Functions (no DB required)
// ============================================================================

func TestNewRepositoryPG(t *testing.T) {
	t.Parallel()
	repo := NewRepositoryPG(nil)
	require.NotNil(t, repo)
	assert.Nil(t, repo.queries)
}

func TestAuthTokenFromDB(t *testing.T) {
	t.Parallel()

	t.Run("all fields populated", func(t *testing.T) {
		t.Parallel()
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		ip := netip.MustParseAddr("10.0.0.1")
		now := time.Now()
		deviceName := "Test Device"
		fingerprint := "abc123"
		ua := "Mozilla/5.0"

		row := db.SharedAuthToken{
			ID:                id,
			UserID:            userID,
			TokenHash:         "hash123",
			TokenType:         "refresh",
			DeviceName:        &deviceName,
			DeviceFingerprint: &fingerprint,
			IpAddress:         ip,
			UserAgent:         &ua,
			ExpiresAt:         now.Add(7 * 24 * time.Hour),
			RevokedAt:         pgtype.Timestamptz{Time: now, Valid: true},
			LastUsedAt:        pgtype.Timestamptz{Time: now, Valid: true},
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		token := authTokenFromDB(row)

		assert.Equal(t, id, token.ID)
		assert.Equal(t, userID, token.UserID)
		assert.Equal(t, "hash123", token.TokenHash)
		assert.Equal(t, "refresh", token.TokenType)
		require.NotNil(t, token.DeviceName)
		assert.Equal(t, "Test Device", *token.DeviceName)
		require.NotNil(t, token.DeviceFingerprint)
		assert.Equal(t, "abc123", *token.DeviceFingerprint)
		require.NotNil(t, token.IPAddress)
		assert.Equal(t, ip, *token.IPAddress)
		require.NotNil(t, token.UserAgent)
		assert.Equal(t, "Mozilla/5.0", *token.UserAgent)
		require.NotNil(t, token.RevokedAt)
		require.NotNil(t, token.LastUsedAt)
		assert.Equal(t, now, token.CreatedAt)
		assert.Equal(t, now, token.UpdatedAt)
	})

	t.Run("optional fields nil", func(t *testing.T) {
		t.Parallel()
		row := db.SharedAuthToken{
			ID:         uuid.Must(uuid.NewV7()),
			UserID:     uuid.Must(uuid.NewV7()),
			TokenHash:  "hash",
			TokenType:  "refresh",
			IpAddress:  netip.Addr{}, // zero value - not valid
			RevokedAt:  pgtype.Timestamptz{Valid: false},
			LastUsedAt: pgtype.Timestamptz{Valid: false},
		}

		token := authTokenFromDB(row)

		assert.Nil(t, token.IPAddress)
		assert.Nil(t, token.RevokedAt)
		assert.Nil(t, token.LastUsedAt)
		assert.Nil(t, token.DeviceName)
		assert.Nil(t, token.DeviceFingerprint)
		assert.Nil(t, token.UserAgent)
	})
}

func TestPasswordResetTokenFromDB(t *testing.T) {
	t.Parallel()

	t.Run("all fields populated", func(t *testing.T) {
		t.Parallel()
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		ip := netip.MustParseAddr("192.168.1.1")
		now := time.Now()
		ua := "TestAgent"

		row := db.SharedPasswordResetToken{
			ID:        id,
			UserID:    userID,
			TokenHash: "reset-hash",
			IpAddress: ip,
			UserAgent: &ua,
			ExpiresAt: now.Add(1 * time.Hour),
			UsedAt:    pgtype.Timestamptz{Time: now, Valid: true},
			CreatedAt: now,
		}

		token := passwordResetTokenFromDB(row)

		assert.Equal(t, id, token.ID)
		assert.Equal(t, userID, token.UserID)
		assert.Equal(t, "reset-hash", token.TokenHash)
		require.NotNil(t, token.IPAddress)
		assert.Equal(t, ip, *token.IPAddress)
		require.NotNil(t, token.UserAgent)
		assert.Equal(t, "TestAgent", *token.UserAgent)
		require.NotNil(t, token.UsedAt)
		assert.Equal(t, now, token.CreatedAt)
	})

	t.Run("optional fields nil", func(t *testing.T) {
		t.Parallel()
		row := db.SharedPasswordResetToken{
			ID:        uuid.Must(uuid.NewV7()),
			UserID:    uuid.Must(uuid.NewV7()),
			TokenHash: "hash",
			IpAddress: netip.Addr{},
			UsedAt:    pgtype.Timestamptz{Valid: false},
		}

		token := passwordResetTokenFromDB(row)

		assert.Nil(t, token.IPAddress)
		assert.Nil(t, token.UsedAt)
		assert.Nil(t, token.UserAgent)
	})
}

func TestEmailVerificationTokenFromDB(t *testing.T) {
	t.Parallel()

	t.Run("all fields populated", func(t *testing.T) {
		t.Parallel()
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		ip := netip.MustParseAddr("10.0.0.1")
		now := time.Now()
		ua := "TestAgent"

		row := db.SharedEmailVerificationToken{
			ID:         id,
			UserID:     userID,
			TokenHash:  "verify-hash",
			Email:      "user@example.com",
			IpAddress:  ip,
			UserAgent:  &ua,
			ExpiresAt:  now.Add(24 * time.Hour),
			VerifiedAt: pgtype.Timestamptz{Time: now, Valid: true},
			CreatedAt:  now,
		}

		token := emailVerificationTokenFromDB(row)

		assert.Equal(t, id, token.ID)
		assert.Equal(t, userID, token.UserID)
		assert.Equal(t, "verify-hash", token.TokenHash)
		assert.Equal(t, "user@example.com", token.Email)
		require.NotNil(t, token.IPAddress)
		assert.Equal(t, ip, *token.IPAddress)
		require.NotNil(t, token.UserAgent)
		assert.Equal(t, "TestAgent", *token.UserAgent)
		require.NotNil(t, token.VerifiedAt)
		assert.Equal(t, now, token.CreatedAt)
	})

	t.Run("optional fields nil", func(t *testing.T) {
		t.Parallel()
		row := db.SharedEmailVerificationToken{
			ID:         uuid.Must(uuid.NewV7()),
			UserID:     uuid.Must(uuid.NewV7()),
			TokenHash:  "hash",
			Email:      "test@test.com",
			IpAddress:  netip.Addr{},
			VerifiedAt: pgtype.Timestamptz{Valid: false},
		}

		token := emailVerificationTokenFromDB(row)

		assert.Nil(t, token.IPAddress)
		assert.Nil(t, token.VerifiedAt)
		assert.Nil(t, token.UserAgent)
	})
}

// ============================================================================
// NewServiceForTestingWithEmail Tests
// ============================================================================

func TestNewServiceForTestingWithEmail(t *testing.T) {
	t.Parallel()
	repo := newMockRepo(t)
	tm := newMockTokenMgr(t)
	actLogger := activity.NewNoopLogger()

	svc := NewServiceForTestingWithEmail(
		nil, repo, tm, actLogger, nil,
		15*time.Minute, 7*24*time.Hour,
	)
	require.NotNil(t, svc)
	assert.Equal(t, 15*time.Minute, svc.jwtExpiry)
	assert.Equal(t, 7*24*time.Hour, svc.refreshExpiry)
	assert.False(t, svc.lockoutEnabled)
	assert.Nil(t, svc.emailService)
}

// ============================================================================
// Helpers
// ============================================================================

func strPtr(s string) *string {
	return &s
}

func testLogger() *slog.Logger {
	return slog.Default()
}
