//go:build integration
// +build integration

package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jwtSecret = "test-secret-key-for-integration-tests-only"

func setupAuthService(t *testing.T) (*auth.Service, *user.Service, *pgxpool.Pool, func()) {
	ctx := context.Background()

	// Create database pool
	pool, err := pgxpool.New(ctx, testDatabaseURL)
	require.NoError(t, err)

	// Create queries
	queries := db.New(pool)

	// Create repositories
	authRepo := auth.NewRepositoryPG(queries)
	userRepo := user.NewPostgresRepository(queries)

	// Create token manager
	tokenManager := auth.NewTokenManager(jwtSecret, 15*time.Minute)

	// Create services
	authSvc := auth.NewService(authRepo, tokenManager, 15*time.Minute, 7*24*time.Hour)
	userSvc := user.NewService(userRepo)

	cleanup := func() {
		pool.Close()
	}

	return authSvc, userSvc, pool, cleanup
}

func TestAuthService_Register(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	req := auth.RegisterRequest{
		Username: fmt.Sprintf("newuser_%d", timestamp),
		Email:    fmt.Sprintf("new_%d@example.com", timestamp),
		Password: "SecurePassword123!",
	}

	user, err := authSvc.Register(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, req.Password, user.PasswordHash, "password should be hashed")
	assert.Contains(t, user.PasswordHash, "$argon2id$", "should use Argon2id")
}

func TestAuthService_Login(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register user first
	password := "SecurePassword123!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("loginuser_%d", timestamp),
		Email:    fmt.Sprintf("login_%d@example.com", timestamp),
		Password: password,
	}

	registeredUser, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Login with username
	resp, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, registeredUser.ID, resp.User.ID)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Greater(t, resp.ExpiresIn, int64(0))
}

func TestAuthService_LoginWithEmail(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register user
	password := "SecurePassword123!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("emaillogin_%d", timestamp),
		Email:    fmt.Sprintf("emaillogin_%d@example.com", timestamp),
		Password: password,
	}

	registeredUser, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Login with email instead of username
	resp, err := authSvc.Login(ctx, req.Email, password, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, registeredUser.ID, resp.User.ID)
}

func TestAuthService_LoginWrongPassword(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register user
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("wrongpwd_%d", timestamp),
		Email:    fmt.Sprintf("wrongpwd_%d@example.com", timestamp),
		Password: "CorrectPassword123!",
	}

	_, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Try login with wrong password
	_, err = authSvc.Login(ctx, req.Username, "WrongPassword123!", nil, nil, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username or password")
}

func TestAuthService_RefreshToken(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register and login
	password := "SecurePassword123!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("refreshuser_%d", timestamp),
		Email:    fmt.Sprintf("refresh_%d@example.com", timestamp),
		Password: password,
	}

	_, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	loginResp, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)

	// Refresh token should generate NEW access token (different timestamp)
	refreshResp, err := authSvc.RefreshToken(ctx, loginResp.RefreshToken)
	require.NoError(t, err)
	assert.NotNil(t, refreshResp)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEqual(t, loginResp.AccessToken, refreshResp.AccessToken, "should get new access token with different timestamp")
}

func TestAuthService_Logout(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register and login
	password := "SecurePassword123!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("logoutuser_%d", timestamp),
		Email:    fmt.Sprintf("logout_%d@example.com", timestamp),
		Password: password,
	}

	_, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	loginResp, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)

	// Logout
	err = authSvc.Logout(ctx, loginResp.RefreshToken)
	require.NoError(t, err)

	// Try to refresh with revoked token
	_, err = authSvc.RefreshToken(ctx, loginResp.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired refresh token")
}

func TestAuthService_ChangePassword(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register user
	oldPassword := "OldPassword123!"
	newPassword := "NewPassword456!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("changepwd_%d", timestamp),
		Email:    fmt.Sprintf("changepwd_%d@example.com", timestamp),
		Password: oldPassword,
	}

	user, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Change password
	err = authSvc.ChangePassword(ctx, user.ID, oldPassword, newPassword)
	require.NoError(t, err)

	// Login with new password
	resp, err := authSvc.Login(ctx, req.Username, newPassword, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Old password should not work
	_, err = authSvc.Login(ctx, req.Username, oldPassword, nil, nil, nil, nil)
	assert.Error(t, err)
}

func TestAuthService_ChangePasswordWrongOldPassword(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register user
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("wrongoldpwd_%d", timestamp),
		Email:    fmt.Sprintf("wrongoldpwd_%d@example.com", timestamp),
		Password: "CorrectPassword123!",
	}

	user, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Try to change with wrong old password
	err = authSvc.ChangePassword(ctx, user.ID, "WrongOldPassword!", "NewPassword456!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid current password")
}

// CRITICAL TEST: Cross-service password compatibility
func TestPasswordCompatibility_UserServiceToAuthService(t *testing.T) {
	authSvc, userSvc, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	password := "CrossServicePassword123!"

	// Create user via User Service
	userParams := user.CreateUserParams{
		Username:     fmt.Sprintf("crosssvc_%d", timestamp),
		Email:        fmt.Sprintf("crosssvc_%d@example.com", timestamp),
		PasswordHash: password, // Will be hashed by User Service
	}

	createdUser, err := userSvc.CreateUser(ctx, userParams)
	require.NoError(t, err)

	// Login via Auth Service with same password
	loginResp, err := authSvc.Login(ctx, userParams.Username, password, nil, nil, nil, nil)
	require.NoError(t, err, "User created by User Service should be able to login via Auth Service")
	assert.NotNil(t, loginResp)
	assert.Equal(t, createdUser.ID, loginResp.User.ID)
	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
}

// Test the opposite direction
func TestPasswordCompatibility_AuthServiceToUserService(t *testing.T) {
	authSvc, userSvc, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	password := "RegisteredPassword123!"

	// Register user via Auth Service
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("authtouse_%d", timestamp),
		Email:    fmt.Sprintf("authtouse_%d@example.com", timestamp),
		Password: password,
	}

	registeredUser, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Verify password via User Service
	user, err := userSvc.GetUser(ctx, registeredUser.ID)
	require.NoError(t, err)

	err = userSvc.VerifyPassword(user.PasswordHash, password)
	require.NoError(t, err, "User Service should be able to verify password created by Auth Service")
}

func TestAuthService_MultipleDeviceLogin(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register user
	password := "MultiDevice123!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("multidev_%d", timestamp),
		Email:    fmt.Sprintf("multidev_%d@example.com", timestamp),
		Password: password,
	}

	_, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Login from device 1
	device1 := "device-fingerprint-1"
	resp1, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, &device1)
	require.NoError(t, err)

	// Login from device 2
	device2 := "device-fingerprint-2"
	resp2, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, &device2)
	require.NoError(t, err)

	// Both tokens should be different
	assert.NotEqual(t, resp1.RefreshToken, resp2.RefreshToken)

	// Both should be able to refresh
	_, err = authSvc.RefreshToken(ctx, resp1.RefreshToken)
	require.NoError(t, err)

	_, err = authSvc.RefreshToken(ctx, resp2.RefreshToken)
	require.NoError(t, err)
}

func TestAuthService_LogoutAll(t *testing.T) {
	authSvc, _, _, cleanup := setupAuthService(t)
	defer cleanup()

	ctx := context.Background()
	timestamp := time.Now().UnixNano()

	// Register and create multiple sessions
	password := "LogoutAll123!"
	req := auth.RegisterRequest{
		Username: fmt.Sprintf("logoutall_%d", timestamp),
		Email:    fmt.Sprintf("logoutall_%d@example.com", timestamp),
		Password: password,
	}

	user, err := authSvc.Register(ctx, req)
	require.NoError(t, err)

	// Login from 3 devices
	resp1, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)

	resp2, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)

	resp3, err := authSvc.Login(ctx, req.Username, password, nil, nil, nil, nil)
	require.NoError(t, err)

	// Logout all devices
	err = authSvc.LogoutAll(ctx, user.ID)
	require.NoError(t, err)

	// All tokens should be revoked
	_, err = authSvc.RefreshToken(ctx, resp1.RefreshToken)
	assert.Error(t, err)

	_, err = authSvc.RefreshToken(ctx, resp2.RefreshToken)
	assert.Error(t, err)

	_, err = authSvc.RefreshToken(ctx, resp3.RefreshToken)
	assert.Error(t, err)
}
