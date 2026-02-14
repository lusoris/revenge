package auth_test

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
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
		testDB.Pool(),
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
	logger := logging.NewTestLogger()
	activitySvc := activity.NewService(activity.NewRepositoryPg(queries), logger)
	activityLogger := activity.NewLogger(activitySvc)

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
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
		testDB.Pool(),
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
			userID:      uuid.Must(uuid.NewV7()),
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
		testDB.Pool(),
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

	// Create password reset token directly in database for testing
	// (since RequestPasswordReset no longer returns the token for security reasons)
	plainToken := "test-reset-token-12345"
	tokenHash := tokenMgr.HashRefreshToken(plainToken)
	_, err = repo.CreatePasswordResetToken(ctx, auth.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		IPAddress: nil,
		UserAgent: nil,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	require.NoError(t, err)
	resetToken := plainToken

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
		testDB.Pool(),
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

// TestService_Register_TransactionAtomicity verifies that user registration is atomic:
// both user creation and email verification token creation succeed or both fail.
// This test ensures that transaction boundaries prevent orphaned user accounts.
// Ref: A7.1.1 in TODO_A7_SECURITY_FIXES.md
func TestService_Register_TransactionAtomicity(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	t.Run("successful registration creates both user and token atomically", func(t *testing.T) {
		req := auth.RegisterRequest{
			Username: "atomicuser",
			Email:    "atomic@example.com",
			Password: "SecurePassword123!",
		}

		// Register user
		user, err := svc.Register(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, req.Email, user.Email)

		// Verify user exists in database
		dbUser, err := queries.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, dbUser.ID)
		assert.Equal(t, user.Username, dbUser.Username)

		// Verify email verification token exists
		// Note: We can't easily check the token without the plain token value,
		// but the successful Register call ensures it was created atomically
		// within the same transaction as the user.
	})

	t.Run("failed registration prevents orphaned user records", func(t *testing.T) {
		// First registration succeeds
		req1 := auth.RegisterRequest{
			Username: "uniqueuser",
			Email:    "unique@example.com",
			Password: "SecurePassword123!",
		}
		user1, err := svc.Register(ctx, req1)
		require.NoError(t, err)
		require.NotNil(t, user1)

		// Second registration with same username should fail
		req2 := auth.RegisterRequest{
			Username: "uniqueuser", // Duplicate username
			Email:    "different@example.com",
			Password: "SecurePassword123!",
		}
		user2, err := svc.Register(ctx, req2)
		require.Error(t, err)
		assert.Nil(t, user2)
		assert.Contains(t, err.Error(), "failed to create user")

		// Verify no user exists with the second email
		_, err = queries.GetUserByEmail(ctx, req2.Email)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no rows")

		// This verifies transaction atomicity: if CreateUser fails,
		// CreateEmailVerificationToken is never executed, and if it were
		// executed and failed, the entire transaction would rollback.
	})
}

// ============================================================================
// VerifyEmail Integration Tests
// ============================================================================

func TestService_VerifyEmail_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	t.Run("successful email verification", func(t *testing.T) {
		// Register a user (creates user + email verification token in one transaction)
		user, err := svc.Register(ctx, auth.RegisterRequest{
			Username: "verifyuser",
			Email:    "verifyuser@example.com",
			Password: "SecurePassword123!",
		})
		require.NoError(t, err)
		require.NotNil(t, user)

		// Verify user email is not yet verified
		dbUser, err := queries.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		if dbUser.EmailVerified != nil {
			assert.False(t, *dbUser.EmailVerified)
		}

		// Create a verification token we know the plain value of
		plainToken := "known-verification-token-123"
		tokenHash := tokenMgr.HashRefreshToken(plainToken)
		_, err = repo.CreateEmailVerificationToken(ctx, auth.CreateEmailVerificationTokenParams{
			UserID:    user.ID,
			TokenHash: tokenHash,
			Email:     user.Email,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)

		// Verify email
		err = svc.VerifyEmail(ctx, plainToken)
		require.NoError(t, err)

		// Verify user's email is now verified
		dbUser, err = queries.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, dbUser.EmailVerified)
		assert.True(t, *dbUser.EmailVerified)
	})

	t.Run("invalid token fails", func(t *testing.T) {
		err := svc.VerifyEmail(ctx, "totally-invalid-token")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid or expired verification token")
	})

	t.Run("token cannot be reused", func(t *testing.T) {
		// Register another user
		user, err := svc.Register(ctx, auth.RegisterRequest{
			Username: "verifyreuse",
			Email:    "verifyreuse@example.com",
			Password: "SecurePassword123!",
		})
		require.NoError(t, err)

		// Create a verification token
		plainToken := "reuse-verification-token-456"
		tokenHash := tokenMgr.HashRefreshToken(plainToken)
		_, err = repo.CreateEmailVerificationToken(ctx, auth.CreateEmailVerificationTokenParams{
			UserID:    user.ID,
			TokenHash: tokenHash,
			Email:     user.Email,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		require.NoError(t, err)

		// First verification succeeds
		err = svc.VerifyEmail(ctx, plainToken)
		require.NoError(t, err)

		// Second verification with same token fails (token was marked used)
		err = svc.VerifyEmail(ctx, plainToken)
		require.Error(t, err)
	})
}

// TestService_VerifyEmail_TransactionAtomicity verifies that VerifyEmail is atomic:
// both token marking and user verification happen together or not at all.
func TestService_VerifyEmail_TransactionAtomicity(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register user
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "atomicverify",
		Email:    "atomicverify@example.com",
		Password: "SecurePassword123!",
	})
	require.NoError(t, err)

	// Create a verification token
	plainToken := "atomic-verify-token-789"
	tokenHash := tokenMgr.HashRefreshToken(plainToken)
	verifyToken, err := repo.CreateEmailVerificationToken(ctx, auth.CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	// Verify email
	err = svc.VerifyEmail(ctx, plainToken)
	require.NoError(t, err)

	// Both should have been updated atomically:
	// 1. Token should be marked as verified
	var verifiedAtCount int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE id = $1 AND verified_at IS NOT NULL",
		verifyToken.ID,
	).Scan(&verifiedAtCount)
	require.NoError(t, err)
	assert.Equal(t, 1, verifiedAtCount)

	// 2. User should have email_verified = true
	dbUser, err := queries.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, dbUser.EmailVerified)
	assert.True(t, *dbUser.EmailVerified)
}

// ============================================================================
// Register → VerifyEmail → Login Full Flow
// ============================================================================

func TestService_RegisterVerifyLogin_FullFlow(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()
	password := "FullFlowPassword123!"

	// Step 1: Register
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "fullflow",
		Email:    "fullflow@example.com",
		Password: password,
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	// Step 2: Verify email (create a known token for test)
	plainToken := "fullflow-verify-token"
	tokenHash := tokenMgr.HashRefreshToken(plainToken)
	_, err = repo.CreateEmailVerificationToken(ctx, auth.CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)

	err = svc.VerifyEmail(ctx, plainToken)
	require.NoError(t, err)

	// Confirm email is verified
	dbUser, err := queries.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, dbUser.EmailVerified)
	assert.True(t, *dbUser.EmailVerified)

	// Step 3: Login with verified account
	loginResp, err := svc.Login(ctx, user.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, loginResp)
	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
	assert.Equal(t, user.ID, loginResp.User.ID)

	// Step 4: Verify last login was updated
	dbUser, err = queries.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	assert.True(t, dbUser.LastLoginAt.Valid)
}

// ============================================================================
// Account Lockout Integration Tests
// ============================================================================

func TestService_Login_AccountLockout_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	logger := logging.NewTestLogger()
	activitySvc := activity.NewService(activity.NewRepositoryPg(queries), logger)
	activityLogger := activity.NewLogger(activitySvc)

	// Use lockout threshold of 3 with a 15-minute window
	svc := auth.NewServiceForTestingWithLockout(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
		3,              // lockout after 3 failed attempts
		15*time.Minute, // lockout window
	)

	ctx := context.Background()
	password := "LockoutTestPassword123!"
	ipAddr := netip.MustParseAddr("192.168.1.50")

	// Register a user
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "lockoutuser",
		Email:    "lockout@example.com",
		Password: password,
	})
	require.NoError(t, err)

	t.Run("failed attempts are recorded", func(t *testing.T) {
		// Try to login with wrong password 2 times
		for range 2 {
			_, err := svc.Login(ctx, user.Username, "WrongPassword!", &ipAddr, nil, nil, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid username or password")
		}

		// Verify attempts were recorded
		count, err := repo.CountFailedLoginAttemptsByUsername(ctx, user.Username, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// Login should still work (below threshold)
		resp, err := svc.Login(ctx, user.Username, password, &ipAddr, nil, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("account locks after threshold", func(t *testing.T) {
		// Clear previous attempts first (successful login above cleared them)
		// Record 3 failed attempts to reach lockout threshold
		for range 3 {
			_, err := svc.Login(ctx, user.Username, "WrongPassword!", &ipAddr, nil, nil, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid username or password")
		}

		// Now even a correct password should fail due to lockout
		_, err := svc.Login(ctx, user.Username, password, &ipAddr, nil, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "account locked")
	})

	t.Run("successful login clears failed attempts", func(t *testing.T) {
		// Clear the lockout by directly clearing failed attempts
		err := repo.ClearFailedLoginAttemptsByUsername(ctx, user.Username)
		require.NoError(t, err)

		// Login should succeed now
		resp, err := svc.Login(ctx, user.Username, password, &ipAddr, nil, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// After successful login, failed attempts should be cleared
		count, err := repo.CountFailedLoginAttemptsByUsername(ctx, user.Username, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestService_Login_LockoutWithNonexistentUser(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTestingWithLockout(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
		3,
		15*time.Minute,
	)

	ctx := context.Background()
	ipAddr := netip.MustParseAddr("10.0.0.1")

	// Try to login with a nonexistent user multiple times with lockout enabled
	// This should record failed attempts and eventually lock out
	for range 3 {
		_, err := svc.Login(ctx, "nonexistent", "password", &ipAddr, nil, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid username or password")
	}

	// After threshold, should get lockout error
	_, err := svc.Login(ctx, "nonexistent", "password", &ipAddr, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "account locked")
}

func TestService_Login_LockoutIPTracking(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTestingWithLockout(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
		3,
		15*time.Minute,
	)

	ctx := context.Background()
	ipAddr := netip.MustParseAddr("172.16.0.50")

	// Register a user
	_, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "iptrackuser",
		Email:    "iptrack@example.com",
		Password: "SecurePassword123!",
	})
	require.NoError(t, err)

	// Failed login attempts should be tracked by IP
	for range 2 {
		_, err := svc.Login(ctx, "iptrackuser", "WrongPassword!", &ipAddr, nil, nil, nil)
		require.Error(t, err)
	}

	// Verify IP-based tracking
	ipCount, err := repo.CountFailedLoginAttemptsByIP(ctx, ipAddr.String(), time.Now().Add(-1*time.Hour))
	require.NoError(t, err)
	assert.Equal(t, int64(2), ipCount)
}

// ============================================================================
// ResendVerification Integration Test
// ============================================================================

func TestService_ResendVerification_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register a user
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "resenduser",
		Email:    "resenduser@example.com",
		Password: "SecurePassword123!",
	})
	require.NoError(t, err)

	// Resend verification
	err = svc.ResendVerification(ctx, user.ID)
	require.NoError(t, err)

	// Verify new token was created (at least 2 tokens exist: one from register, one from resend)
	// The old ones should be invalidated
	var activeCount int
	err = testDB.Pool().QueryRow(ctx,
		"SELECT COUNT(*) FROM shared.email_verification_tokens WHERE user_id = $1 AND verified_at IS NULL",
		user.ID,
	).Scan(&activeCount)
	require.NoError(t, err)
	// After resend: old tokens invalidated, new token created
	assert.Equal(t, 1, activeCount)
}

// ============================================================================
// RequestPasswordReset Integration Test
// ============================================================================

func TestService_RequestPasswordReset_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Register a user
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "resetreqtest",
		Email:    "resetreq@example.com",
		Password: "SecurePassword123!",
	})
	require.NoError(t, err)

	t.Run("creates reset token for existing email", func(t *testing.T) {
		err := svc.RequestPasswordReset(ctx, user.Email, nil, nil)
		require.NoError(t, err)

		// Verify a reset token was created
		var tokenCount int
		err = testDB.Pool().QueryRow(ctx,
			"SELECT COUNT(*) FROM shared.password_reset_tokens WHERE user_id = $1 AND used_at IS NULL",
			user.ID,
		).Scan(&tokenCount)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, tokenCount, 1)
	})

	t.Run("silently succeeds for nonexistent email", func(t *testing.T) {
		// This should NOT error (prevents email enumeration)
		err := svc.RequestPasswordReset(ctx, "nonexistent@example.com", nil, nil)
		require.NoError(t, err)
	})

	t.Run("invalidates old reset tokens on new request", func(t *testing.T) {
		// First request
		err := svc.RequestPasswordReset(ctx, user.Email, nil, nil)
		require.NoError(t, err)

		// Second request should invalidate old tokens
		err = svc.RequestPasswordReset(ctx, user.Email, nil, nil)
		require.NoError(t, err)

		// Verify only 1 active (unused) token exists
		var activeTokens int
		err = testDB.Pool().QueryRow(ctx,
			"SELECT COUNT(*) FROM shared.password_reset_tokens WHERE user_id = $1 AND used_at IS NULL",
			user.ID,
		).Scan(&activeTokens)
		require.NoError(t, err)
		assert.Equal(t, 1, activeTokens)
	})
}

// ============================================================================
// Logout Integration Tests
// ============================================================================

func TestService_LogoutAll_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()
	password := "LogoutAllPassword123!"

	// Register and login multiple times to create multiple refresh tokens
	user, err := svc.Register(ctx, auth.RegisterRequest{
		Username: "logoutalluser",
		Email:    "logoutall@example.com",
		Password: password,
	})
	require.NoError(t, err)

	var refreshTokens []string
	for range 3 {
		resp, err := svc.Login(ctx, user.Username, password, nil, nil, nil, nil)
		require.NoError(t, err)
		refreshTokens = append(refreshTokens, resp.RefreshToken)
	}

	// LogoutAll should revoke all tokens
	err = svc.LogoutAll(ctx, user.ID)
	require.NoError(t, err)

	// All refresh tokens should now be invalid
	for _, rt := range refreshTokens {
		_, err := svc.RefreshToken(ctx, rt)
		require.Error(t, err)
	}

	// Verify no active auth tokens remain
	count, err := repo.CountActiveAuthTokensByUser(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// ============================================================================
// RegisterFromOIDC Integration Test
// ============================================================================

func TestService_RegisterFromOIDC_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	displayName := "OIDC Test User"
	user, err := svc.RegisterFromOIDC(ctx, auth.RegisterFromOIDCRequest{
		Username:    "oidcintegration",
		Email:       "oidcintegration@example.com",
		DisplayName: &displayName,
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	assert.Equal(t, "oidcintegration", user.Username)
	assert.Equal(t, "oidcintegration@example.com", user.Email)

	// OIDC users should have email verified
	dbUser, err := queries.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, dbUser.EmailVerified)
	assert.True(t, *dbUser.EmailVerified)
}

// ============================================================================
// CreateSessionForUser Integration Test
// ============================================================================

func TestService_CreateSessionForUser_Integration(t *testing.T) {
	t.Parallel()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := auth.NewRepositoryPG(queries)
	tokenMgr := auth.NewTokenManager("test-secret-key-at-least-32-characters-long", 15*time.Minute)
	activityLogger := activity.NewNoopLogger()

	svc := auth.NewServiceForTesting(
		testDB.Pool(),
		repo,
		tokenMgr,
		activityLogger,
		15*time.Minute,
		7*24*time.Hour,
	)

	ctx := context.Background()

	// Create a user via OIDC (so they have verified email + active account)
	user, err := svc.RegisterFromOIDC(ctx, auth.RegisterFromOIDCRequest{
		Username: "sessionuser",
		Email:    "sessionuser@example.com",
	})
	require.NoError(t, err)

	ipAddr := netip.MustParseAddr("10.0.0.1")
	userAgent := "TestAgent/1.0"
	deviceName := "Test Device"

	resp, err := svc.CreateSessionForUser(ctx, user.ID, uuid.Nil, &ipAddr, &userAgent, &deviceName)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, user.ID, resp.User.ID)

	// Verify refresh token works
	refreshResp, err := svc.RefreshToken(ctx, resp.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResp.AccessToken)
}
