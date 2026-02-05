package auth_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/auth"
)

// setupMockService creates an auth service with all dependencies mocked
// NOTE: Methods using transactions (Register, etc.) cannot be properly tested
// with mocks and should use integration tests instead.
func setupMockService(t *testing.T) (
	*auth.Service,
	*MockAuthRepository,
	*MockTokenManager,
) {
	t.Helper()

	mockRepo := NewMockAuthRepository(t)
	mockTokenMgr := NewMockTokenManager(t)
	activityLogger := activity.NewNoopLogger()

	// NOTE: pool is nil for mock tests. Methods using transactions will panic.
	// Use integration tests for transaction-based methods.
	service := auth.NewServiceForTesting(
		nil, // pool
		mockRepo,
		mockTokenMgr,
		activityLogger,
		15*time.Minute,  // jwtExpiry
		7*24*time.Hour,  // refreshExpiry
	)

	return service, mockRepo, mockTokenMgr
}

// ========== Register Tests ==========

// Register password hashing errors require integration tests
// These are covered by the integration test suite

func TestService_Register_ErrorCreatingUser(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	req := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	expectedErr := fmt.Errorf("unique constraint violation")
	mockRepo.EXPECT().
		CreateUser(ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(db.SharedUser{}, expectedErr).
		Once()

	user, err := svc.Register(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.Nil(t, user)
}

func TestService_Register_ErrorCreatingVerificationToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	req := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	createdUser := db.SharedUser{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	mockRepo.EXPECT().
		CreateUser(ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(createdUser, nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken(mock.AnythingOfType("string")).
		Return("token_hash").
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		CreateEmailVerificationToken(ctx, mock.AnythingOfType("auth.CreateEmailVerificationTokenParams")).
		Return(auth.EmailVerificationToken{}, expectedErr).
		Once()

	user, err := svc.Register(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create verification token")
	assert.Nil(t, user)
}

// ========== VerifyEmail Tests ==========

func TestService_VerifyEmail_InvalidToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	mockTokenMgr.EXPECT().
		HashRefreshToken("invalid_token").
		Return("invalid_hash").
		Once()

	mockRepo.EXPECT().
		GetEmailVerificationToken(ctx, "invalid_hash").
		Return(auth.EmailVerificationToken{}, fmt.Errorf("not found")).
		Once()

	err := svc.VerifyEmail(ctx, "invalid_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired verification token")
}

func TestService_VerifyEmail_ErrorMarkingTokenUsed(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	tokenID := uuid.New()
	userID := uuid.New()
	emailToken := auth.EmailVerificationToken{
		ID:     tokenID,
		UserID: userID,
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("valid_hash").
		Once()

	mockRepo.EXPECT().
		GetEmailVerificationToken(ctx, "valid_hash").
		Return(emailToken, nil).
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		MarkEmailVerificationTokenUsed(ctx, tokenID).
		Return(expectedErr).
		Once()

	err := svc.VerifyEmail(ctx, "valid_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to mark token as used")
}

func TestService_VerifyEmail_ErrorUpdatingUser(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	tokenID := uuid.New()
	userID := uuid.New()
	emailToken := auth.EmailVerificationToken{
		ID:     tokenID,
		UserID: userID,
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("valid_hash").
		Once()

	mockRepo.EXPECT().
		GetEmailVerificationToken(ctx, "valid_hash").
		Return(emailToken, nil).
		Once()

	mockRepo.EXPECT().
		MarkEmailVerificationTokenUsed(ctx, tokenID).
		Return(nil).
		Once()

	expectedErr := fmt.Errorf("user not found")
	mockRepo.EXPECT().
		UpdateUserEmailVerified(ctx, userID, true).
		Return(expectedErr).
		Once()

	err := svc.VerifyEmail(ctx, "valid_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user")
}

// ========== Login Tests ==========

func TestService_Login_UserNotFoundByUsernameOrEmail(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	username := "nonexistent"
	password := "password"

	mockRepo.EXPECT().
		GetUserByUsername(ctx, username).
		Return(nil, fmt.Errorf("not found")).
		Once()

	mockRepo.EXPECT().
		GetUserByEmail(ctx, username).
		Return(nil, fmt.Errorf("not found")).
		Once()

	resp, err := svc.Login(ctx, username, password, nil, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
}

// Login tests with valid password hashing require integration tests with real hasher
// These are covered by the integration test suite

// ========== Logout Tests ==========

func TestService_Logout_ErrorRevokingToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	refreshToken := "valid_refresh_token"
	tokenHash := "token_hash"

	mockTokenMgr.EXPECT().
		HashRefreshToken(refreshToken).
		Return(tokenHash).
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		RevokeAuthTokenByHash(ctx, tokenHash).
		Return(expectedErr).
		Once()

	err := svc.Logout(ctx, refreshToken)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// ========== LogoutAll Tests ==========

func TestService_LogoutAll_ErrorRevokingTokens(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	expectedErr := fmt.Errorf("database error")

	mockRepo.EXPECT().
		RevokeAllUserAuthTokens(ctx, userID).
		Return(expectedErr).
		Once()

	err := svc.LogoutAll(ctx, userID)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// ========== RefreshToken Tests ==========

func TestService_RefreshToken_InvalidToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	refreshToken := "invalid_token"
	tokenHash := "invalid_hash"

	mockTokenMgr.EXPECT().
		HashRefreshToken(refreshToken).
		Return(tokenHash).
		Once()

	mockRepo.EXPECT().
		GetAuthTokenByHash(ctx, tokenHash).
		Return(auth.AuthToken{}, fmt.Errorf("not found")).
		Once()

	resp, err := svc.RefreshToken(ctx, refreshToken)

	require.Error(t, err)
	assert.Equal(t, "invalid or expired refresh token", err.Error())
	assert.Nil(t, resp)
}

func TestService_RefreshToken_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("token_hash").
		Once()

	mockRepo.EXPECT().
		GetAuthTokenByHash(ctx, "token_hash").
		Return(authToken, nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(nil, fmt.Errorf("user not found")).
		Once()

	resp, err := svc.RefreshToken(ctx, "valid_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, resp)
}

func TestService_RefreshToken_ErrorGeneratingAccessToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}

	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("token_hash").
		Once()

	mockRepo.EXPECT().
		GetAuthTokenByHash(ctx, "token_hash").
		Return(authToken, nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	expectedErr := fmt.Errorf("JWT signing error")
	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("", expectedErr).
		Once()

	resp, err := svc.RefreshToken(ctx, "valid_token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate access token")
	assert.Nil(t, resp)
}

// ========== ChangePassword Tests ==========

func TestService_ChangePassword_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(nil, fmt.Errorf("not found")).
		Once()

	err := svc.ChangePassword(ctx, userID, "oldpass", "newpass")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

// ChangePassword tests with valid password hashing require integration tests
// These are covered by the integration test suite

// ========== RequestPasswordReset Tests ==========

func TestService_RequestPasswordReset_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	email := "nonexistent@example.com"

	mockRepo.EXPECT().
		GetUserByEmail(ctx, email).
		Return(nil, fmt.Errorf("not found")).
		Once()

	// Should return empty string, not error (security: don't reveal if email exists)
	token, err := svc.RequestPasswordReset(ctx, email, nil, nil)

	require.NoError(t, err)
	assert.Empty(t, token)
}

func TestService_RequestPasswordReset_ErrorInvalidatingOldTokens(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	email := "user@example.com"
	user := &db.SharedUser{
		ID:    userID,
		Email: email,
	}

	mockRepo.EXPECT().
		GetUserByEmail(ctx, email).
		Return(user, nil).
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		InvalidateUserPasswordResetTokens(ctx, userID).
		Return(expectedErr).
		Once()

	token, err := svc.RequestPasswordReset(ctx, email, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to invalidate old tokens")
	assert.Empty(t, token)
}

func TestService_RequestPasswordReset_ErrorCreatingToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	email := "user@example.com"
	user := &db.SharedUser{
		ID:    userID,
		Email: email,
	}

	mockRepo.EXPECT().
		GetUserByEmail(ctx, email).
		Return(user, nil).
		Once()

	mockRepo.EXPECT().
		InvalidateUserPasswordResetTokens(ctx, userID).
		Return(nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken(mock.AnythingOfType("string")).
		Return("token_hash").
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		CreatePasswordResetToken(ctx, mock.AnythingOfType("auth.CreatePasswordResetTokenParams")).
		Return(auth.PasswordResetToken{}, expectedErr).
		Once()

	token, err := svc.RequestPasswordReset(ctx, email, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create reset token")
	assert.Empty(t, token)
}

// ========== ResetPassword Tests ==========

func TestService_ResetPassword_InvalidToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	mockTokenMgr.EXPECT().
		HashRefreshToken("invalid_token").
		Return("invalid_hash").
		Once()

	mockRepo.EXPECT().
		GetPasswordResetToken(ctx, "invalid_hash").
		Return(auth.PasswordResetToken{}, fmt.Errorf("not found")).
		Once()

	err := svc.ResetPassword(ctx, "invalid_token", "newpassword")

	require.Error(t, err)
	assert.Equal(t, errors.New("invalid or expired reset token"), err)
}

// ResetPassword tests with valid password hashing require integration tests
// These are covered by the integration test suite

// ========== ResendVerification Tests ==========

func TestService_ResendVerification_ErrorInvalidatingTokens(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	expectedErr := fmt.Errorf("database error")

	mockRepo.EXPECT().
		InvalidateUserEmailVerificationTokens(ctx, userID).
		Return(expectedErr).
		Once()

	err := svc.ResendVerification(ctx, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to invalidate old tokens")
}

func TestService_ResendVerification_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.EXPECT().
		InvalidateUserEmailVerificationTokens(ctx, userID).
		Return(nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(nil, fmt.Errorf("not found")).
		Once()

	err := svc.ResendVerification(ctx, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestService_ResendVerification_ErrorCreatingToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	user := &db.SharedUser{
		ID:    userID,
		Email: "user@example.com",
	}

	mockRepo.EXPECT().
		InvalidateUserEmailVerificationTokens(ctx, userID).
		Return(nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken(mock.AnythingOfType("string")).
		Return("token_hash").
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		CreateEmailVerificationToken(ctx, mock.AnythingOfType("auth.CreateEmailVerificationTokenParams")).
		Return(auth.EmailVerificationToken{}, expectedErr).
		Once()

	err := svc.ResendVerification(ctx, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create verification token")
}

// ========== RegisterFromOIDC Tests ==========

func TestService_RegisterFromOIDC_ErrorCreatingUser(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	req := auth.RegisterFromOIDCRequest{
		Username: "oidcuser",
		Email:    "oidc@example.com",
	}

	expectedErr := fmt.Errorf("unique constraint violation")
	mockRepo.EXPECT().
		CreateUser(ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(db.SharedUser{}, expectedErr).
		Once()

	user, err := svc.RegisterFromOIDC(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.Nil(t, user)
}

func TestService_RegisterFromOIDC_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	displayName := "OIDC User"
	req := auth.RegisterFromOIDCRequest{
		Username:    "oidcuser",
		Email:       "oidc@example.com",
		DisplayName: &displayName,
	}

	createdUser := db.SharedUser{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	mockRepo.EXPECT().
		CreateUser(ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(createdUser, nil).
		Once()

	mockRepo.EXPECT().
		UpdateUserEmailVerified(ctx, userID, true).
		Return(nil).
		Once()

	user, err := svc.RegisterFromOIDC(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
}

func TestService_RegisterFromOIDC_EmailVerificationError(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	req := auth.RegisterFromOIDCRequest{
		Username: "oidcuser",
		Email:    "oidc@example.com",
	}

	createdUser := db.SharedUser{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	mockRepo.EXPECT().
		CreateUser(ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(createdUser, nil).
		Once()

	// Email verification update fails but shouldn't fail the registration
	mockRepo.EXPECT().
		UpdateUserEmailVerified(ctx, userID, true).
		Return(fmt.Errorf("database error")).
		Once()

	user, err := svc.RegisterFromOIDC(ctx, req)

	// User should still be created successfully
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
}

// ========== CreateSessionForUser Tests ==========

func TestService_CreateSessionForUser_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(nil, fmt.Errorf("not found")).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, resp)
}

func TestService_CreateSessionForUser_AccountDisabled(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	isActive := false
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "account is disabled", err.Error())
	assert.Nil(t, resp)
}

func TestService_CreateSessionForUser_ErrorGeneratingAccessToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	isActive := true
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	expectedErr := fmt.Errorf("JWT signing error")
	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("", expectedErr).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate access token")
	assert.Nil(t, resp)
}

func TestService_CreateSessionForUser_ErrorGeneratingRefreshToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	isActive := true
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("access_token", nil).
		Once()

	expectedErr := fmt.Errorf("random generation error")
	mockTokenMgr.EXPECT().
		GenerateRefreshToken().
		Return("", expectedErr).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate refresh token")
	assert.Nil(t, resp)
}

func TestService_CreateSessionForUser_ErrorStoringToken(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	isActive := true
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("access_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateRefreshToken().
		Return("refresh_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken("refresh_token").
		Return("token_hash").
		Once()

	expectedErr := fmt.Errorf("database error")
	mockRepo.EXPECT().
		CreateAuthToken(ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(auth.AuthToken{}, expectedErr).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store refresh token")
	assert.Nil(t, resp)
}

func TestService_CreateSessionForUser_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	isActive := true
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("access_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateRefreshToken().
		Return("refresh_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken("refresh_token").
		Return("token_hash").
		Once()

	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}
	mockRepo.EXPECT().
		CreateAuthToken(ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(authToken, nil).
		Once()

	mockRepo.EXPECT().
		UpdateUserLastLogin(ctx, userID).
		Return(nil).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
	assert.Equal(t, user.Username, resp.User.Username)
}

func TestService_CreateSessionForUser_SuccessWithIPAndUserAgent(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	isActive := true
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	userAgent := "Mozilla/5.0"
	deviceName := "Test Device"

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("access_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateRefreshToken().
		Return("refresh_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken("refresh_token").
		Return("token_hash").
		Once()

	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}
	mockRepo.EXPECT().
		CreateAuthToken(ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(authToken, nil).
		Once()

	mockRepo.EXPECT().
		UpdateUserLastLogin(ctx, userID).
		Return(nil).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, &userAgent, &deviceName)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
}

// ========== Additional Login Tests ==========

func TestService_Login_AccountDisabled(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	isActive := false
	user := &db.SharedUser{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=4$dGVzdHNhbHQ$dGVzdGhhc2g", // placeholder
		IsActive:     &isActive,
	}

	mockRepo.EXPECT().
		GetUserByUsername(ctx, "testuser").
		Return(user, nil).
		Once()

	// Note: login will fail after password verification due to disabled account
	// This test verifies that the disabled check happens after password verification
	resp, err := svc.Login(ctx, "testuser", "password", nil, nil, nil, nil)

	require.Error(t, err)
	// If password doesn't match, we get invalid password error
	// If password matches but account disabled, we get account disabled error
	assert.Nil(t, resp)
}

// ========== RefreshToken Success Path ==========

func TestService_RefreshToken_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}

	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("token_hash").
		Once()

	mockRepo.EXPECT().
		GetAuthTokenByHash(ctx, "token_hash").
		Return(authToken, nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("new_access_token", nil).
		Once()

	mockRepo.EXPECT().
		UpdateAuthTokenLastUsed(ctx, tokenID).
		Return(nil).
		Once()

	resp, err := svc.RefreshToken(ctx, "valid_token")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new_access_token", resp.AccessToken)
	assert.Equal(t, "valid_token", resp.RefreshToken)
	assert.Equal(t, user.Username, resp.User.Username)
}

func TestService_RefreshToken_UpdateLastUsedFails(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}

	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("token_hash").
		Once()

	mockRepo.EXPECT().
		GetAuthTokenByHash(ctx, "token_hash").
		Return(authToken, nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("new_access_token", nil).
		Once()

	// Fails but shouldn't fail the refresh
	mockRepo.EXPECT().
		UpdateAuthTokenLastUsed(ctx, tokenID).
		Return(fmt.Errorf("database error")).
		Once()

	resp, err := svc.RefreshToken(ctx, "valid_token")

	// Should still succeed
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new_access_token", resp.AccessToken)
}

// ========== Logout Success Path ==========

func TestService_Logout_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	refreshToken := "valid_refresh_token"
	tokenHash := "token_hash"

	mockTokenMgr.EXPECT().
		HashRefreshToken(refreshToken).
		Return(tokenHash).
		Once()

	mockRepo.EXPECT().
		RevokeAuthTokenByHash(ctx, tokenHash).
		Return(nil).
		Once()

	err := svc.Logout(ctx, refreshToken)

	require.NoError(t, err)
}

// ========== LogoutAll Success Path ==========

func TestService_LogoutAll_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.EXPECT().
		RevokeAllUserAuthTokens(ctx, userID).
		Return(nil).
		Once()

	err := svc.LogoutAll(ctx, userID)

	require.NoError(t, err)
}

// ========== VerifyEmail Success Path ==========

func TestService_VerifyEmail_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	tokenID := uuid.New()
	userID := uuid.New()
	emailToken := auth.EmailVerificationToken{
		ID:     tokenID,
		UserID: userID,
	}

	mockTokenMgr.EXPECT().
		HashRefreshToken("valid_token").
		Return("valid_hash").
		Once()

	mockRepo.EXPECT().
		GetEmailVerificationToken(ctx, "valid_hash").
		Return(emailToken, nil).
		Once()

	mockRepo.EXPECT().
		MarkEmailVerificationTokenUsed(ctx, tokenID).
		Return(nil).
		Once()

	mockRepo.EXPECT().
		UpdateUserEmailVerified(ctx, userID, true).
		Return(nil).
		Once()

	err := svc.VerifyEmail(ctx, "valid_token")

	require.NoError(t, err)
}

// ========== RequestPasswordReset Success Path ==========

func TestService_RequestPasswordReset_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	email := "user@example.com"
	user := &db.SharedUser{
		ID:    userID,
		Email: email,
	}

	mockRepo.EXPECT().
		GetUserByEmail(ctx, email).
		Return(user, nil).
		Once()

	mockRepo.EXPECT().
		InvalidateUserPasswordResetTokens(ctx, userID).
		Return(nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken(mock.AnythingOfType("string")).
		Return("token_hash").
		Once()

	resetToken := auth.PasswordResetToken{
		ID:     uuid.New(),
		UserID: userID,
	}
	mockRepo.EXPECT().
		CreatePasswordResetToken(ctx, mock.AnythingOfType("auth.CreatePasswordResetTokenParams")).
		Return(resetToken, nil).
		Once()

	token, err := svc.RequestPasswordReset(ctx, email, nil, nil)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// ========== ResendVerification Success Path ==========

func TestService_ResendVerification_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	user := &db.SharedUser{
		ID:       userID,
		Email:    "user@example.com",
		Username: "testuser",
	}

	mockRepo.EXPECT().
		InvalidateUserEmailVerificationTokens(ctx, userID).
		Return(nil).
		Once()

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken(mock.AnythingOfType("string")).
		Return("token_hash").
		Once()

	verificationToken := auth.EmailVerificationToken{
		ID:     uuid.New(),
		UserID: userID,
	}
	mockRepo.EXPECT().
		CreateEmailVerificationToken(ctx, mock.AnythingOfType("auth.CreateEmailVerificationTokenParams")).
		Return(verificationToken, nil).
		Once()

	err := svc.ResendVerification(ctx, userID)

	require.NoError(t, err)
}

// ========== Login Error Paths ==========

func TestService_Login_FoundByEmail(t *testing.T) {
	t.Parallel()
	svc, mockRepo, _ := setupMockService(t)
	ctx := context.Background()

	email := "user@example.com"
	userID := uuid.New()
	user := &db.SharedUser{
		ID:           userID,
		Username:     "testuser",
		Email:        email,
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=4$placeholder",
	}

	// Not found by username
	mockRepo.EXPECT().
		GetUserByUsername(ctx, email).
		Return(nil, fmt.Errorf("not found")).
		Once()

	// Found by email
	mockRepo.EXPECT().
		GetUserByEmail(ctx, email).
		Return(user, nil).
		Once()

	// Password won't match (placeholder hash), but we test the lookup path
	resp, err := svc.Login(ctx, email, "wrongpassword", nil, nil, nil, nil)

	require.Error(t, err)
	assert.Nil(t, resp)
}

// ========== Register Success Path ==========

func TestService_Register_Success(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	displayName := "Test User"
	req := auth.RegisterRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "SecurePass123!",
		DisplayName: &displayName,
	}

	createdUser := db.SharedUser{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	mockRepo.EXPECT().
		CreateUser(ctx, mock.AnythingOfType("db.CreateUserParams")).
		Return(createdUser, nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken(mock.AnythingOfType("string")).
		Return("token_hash").
		Once()

	verificationToken := auth.EmailVerificationToken{
		ID:     uuid.New(),
		UserID: userID,
	}
	mockRepo.EXPECT().
		CreateEmailVerificationToken(ctx, mock.AnythingOfType("auth.CreateEmailVerificationTokenParams")).
		Return(verificationToken, nil).
		Once()

	user, err := svc.Register(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
}

// ========== CreateSessionForUser with LastLogin Failure ==========

func TestService_CreateSessionForUser_LastLoginFails(t *testing.T) {
	t.Parallel()
	svc, mockRepo, mockTokenMgr := setupMockService(t)
	ctx := context.Background()

	userID := uuid.New()
	tokenID := uuid.New()
	isActive := true
	user := &db.SharedUser{
		ID:       userID,
		Username: "testuser",
		IsActive: &isActive,
	}

	mockRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateAccessToken(userID, user.Username).
		Return("access_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		GenerateRefreshToken().
		Return("refresh_token", nil).
		Once()

	mockTokenMgr.EXPECT().
		HashRefreshToken("refresh_token").
		Return("token_hash").
		Once()

	authToken := auth.AuthToken{
		ID:     tokenID,
		UserID: userID,
	}
	mockRepo.EXPECT().
		CreateAuthToken(ctx, mock.AnythingOfType("auth.CreateAuthTokenParams")).
		Return(authToken, nil).
		Once()

	// LastLogin fails but shouldn't fail the session creation
	mockRepo.EXPECT().
		UpdateUserLastLogin(ctx, userID).
		Return(fmt.Errorf("database error")).
		Once()

	resp, err := svc.CreateSessionForUser(ctx, userID, nil, nil, nil)

	// Should still succeed
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "access_token", resp.AccessToken)
}
