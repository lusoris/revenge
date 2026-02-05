package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/testutil"
)

// Integration tests for password-related flows that require real password hashing

func TestService_Register_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	tests := []struct {
		name    string
		req     auth.RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful registration",
			req: auth.RegisterRequest{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "SecurePassword123!",
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			req: auth.RegisterRequest{
				Username: "newuser",
				Email:    "another@example.com",
				Password: "SecurePassword123!",
			},
			wantErr: true,
			errMsg:  "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Register(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.req.Username, user.Username)
				assert.Equal(t, tt.req.Email, user.Email)
				assert.NotEmpty(t, user.PasswordHash)
				// Verify password was hashed (not plaintext)
				assert.NotEqual(t, tt.req.Password, user.PasswordHash)
			}
		})
	}
}

func TestService_Login_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	logger := zap.NewNop()
	activitySvc := activity.NewService(activity.NewRepositoryPg(queries), logger)
	activityLogger := activity.NewLogger(activitySvc)

	svc := auth.NewServiceForTesting(
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register a user first
	password := "TestPassword123!"
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "logintest",
		Email:    "logintest@example.com",
		Password: password,
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid login with username",
			username: "logintest",
			password: password,
			wantErr:  false,
		},
		{
			name:     "valid login with email",
			username: "logintest@example.com",
			password: password,
			wantErr:  false,
		},
		{
			name:     "invalid password",
			username: "logintest",
			password: "WrongPassword123!",
			wantErr:  true,
			errMsg:   "invalid username or password",
		},
		{
			name:     "nonexistent user",
			username: "doesnotexist",
			password: password,
			wantErr:  true,
			errMsg:   "invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Login(ctx, tt.username, tt.password, nil, nil, nil, nil)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, user.ID, resp.User.ID)
				assert.Equal(t, user.Username, resp.User.Username)
			}
		})
	}
}

func TestService_ChangePassword_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register a user first
	oldPassword := "OldPassword123!"
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "changepasstest",
		Email:    "changepass@example.com",
		Password: oldPassword,
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uuid.UUID
		oldPassword string
		newPassword string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful password change",
			userID:      user.ID,
			oldPassword: oldPassword,
			newPassword: "NewPassword456!",
			wantErr:     false,
		},
		{
			name:        "invalid old password",
			userID:      user.ID,
			oldPassword: "WrongOldPassword!",
			newPassword: "NewPassword789!",
			wantErr:     true,
			errMsg:      "invalid current password",
		},
		{
			name:        "nonexistent user",
			userID:      uuid.New(),
			oldPassword: oldPassword,
			newPassword: "NewPassword999!",
			wantErr:     true,
			errMsg:      "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ChangePassword(ctx, tt.userID, tt.oldPassword, tt.newPassword)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)

				// Verify can login with new password
				resp, err := svc.Login(ctx, user.Username, tt.newPassword, nil, nil, nil, nil)
				require.NoError(t, err)
				assert.NotNil(t, resp)

				// Verify cannot login with old password
				_, err = svc.Login(ctx, user.Username, oldPassword, nil, nil, nil, nil)
				require.Error(t, err)
			}
		})
	}
}

func TestService_ResetPassword_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register a user first
	oldPassword := "OldPassword123!"
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "resetpasstest",
		Email:    "resetpass@example.com",
		Password: oldPassword,
	})
	require.NoError(t, err)

	// Request password reset
	resetToken, err := svc.RequestPasswordReset(ctx, user.Email, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, resetToken)

	tests := []struct {
		name        string
		token       string
		newPassword string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful password reset",
			token:       resetToken,
			newPassword: "NewResetPassword456!",
			wantErr:     false,
		},
		{
			name:        "invalid token",
			token:       "invalid-token-12345",
			newPassword: "NewPassword789!",
			wantErr:     true,
			errMsg:      "invalid or expired reset token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ResetPassword(ctx, tt.token, tt.newPassword)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)

				// Verify can login with new password
				resp, err := svc.Login(ctx, user.Username, tt.newPassword, nil, nil, nil, nil)
				require.NoError(t, err)
				assert.NotNil(t, resp)

				// Verify token cannot be reused
				err = svc.ResetPassword(ctx, tt.token, "AnotherPassword!")
				require.Error(t, err)
			}
		})
	}
}

func TestService_RefreshToken_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register and login
	password := "TestPassword123!"
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "refreshtest",
		Email:    "refreshtest@example.com",
		Password: password,
	})
	require.NoError(t, err)

	loginResp, err := svc.Login(ctx, user.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)

	t.Run("valid refresh token", func(t *testing.T) {
		resp, err := svc.RefreshToken(ctx, loginResp.RefreshToken)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken)
		// Access token should be different from original
		assert.NotEqual(t, loginResp.AccessToken, resp.AccessToken)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		resp, err := svc.RefreshToken(ctx, "invalid-token-12345")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid or expired refresh token")
		assert.Nil(t, resp)
	})

	t.Run("revoked refresh token", func(t *testing.T) {
		err := svc.Logout(ctx, loginResp.RefreshToken)
		require.NoError(t, err)

		resp, err := svc.RefreshToken(ctx, loginResp.RefreshToken)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}
