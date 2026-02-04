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
func setupMockService(t *testing.T) (
	*auth.Service,
	*MockAuthRepository,
	*MockTokenManager,
) {
	t.Helper()

	mockRepo := NewMockAuthRepository(t)
	mockTokenMgr := NewMockTokenManager(t)
	activityLogger := activity.NewNoopLogger()

	service := auth.NewServiceForTesting(
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
