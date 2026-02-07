package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/email"
)

// dummyPasswordHash is a precomputed Argon2id hash used for timing attack mitigation.
// When a user doesn't exist, we compare against this hash to ensure constant-time behavior.
// This prevents username enumeration via timing analysis.
// Hash of: "dummy-password-for-timing-attack-mitigation"
const dummyPasswordHash = "$argon2id$v=19$m=65536,t=1,p=24$tQMNjFt979tvL7ho1P6xXw$DXkAY76TwLxFcMyqpMQQowtoWwhHfcs5Da9lFIid0Bg"

// Service implements auth business logic
type Service struct {
	pool             *pgxpool.Pool
	repo             Repository
	tokenManager     TokenManager
	hasher           *crypto.PasswordHasher
	activityLogger   activity.Logger
	emailService     *email.Service
	jwtExpiry        time.Duration
	refreshExpiry    time.Duration
	lockoutThreshold int
	lockoutWindow    time.Duration
	lockoutEnabled   bool
}

// NewService creates a new auth service
func NewService(pool *pgxpool.Pool, repo Repository, tokenManager TokenManager, activityLogger activity.Logger, emailService *email.Service, jwtExpiry, refreshExpiry time.Duration, lockoutThreshold int, lockoutWindow time.Duration, lockoutEnabled bool) *Service {
	return &Service{
		pool:             pool,
		repo:             repo,
		tokenManager:     tokenManager,
		hasher:           crypto.NewPasswordHasher(),
		activityLogger:   activityLogger,
		emailService:     emailService,
		jwtExpiry:        jwtExpiry,
		refreshExpiry:    refreshExpiry,
		lockoutThreshold: lockoutThreshold,
		lockoutWindow:    lockoutWindow,
		lockoutEnabled:   lockoutEnabled,
	}
}

// RegisterRequest contains registration data
type RegisterRequest struct {
	Username    string
	Email       string
	Password    string
	DisplayName *string
}

// LoginResponse contains login result with tokens
type LoginResponse struct {
	User         *db.SharedUser
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds until access token expires
}

// ============================================================================
// Registration
// ============================================================================

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*db.SharedUser, error) {
	// Hash password using Argon2id (per AUTH.md)
	passwordHash, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate email verification token
	token, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}
	tokenHash := s.tokenManager.HashRefreshToken(token)

	// Begin transaction to ensure atomicity of user creation and token generation
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Create queries with transaction context
	txQueries := db.New(tx)

	// Create user in database (within transaction)
	user, err := txQueries.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Store verification token in database (within transaction)
	// Both operations succeed or both fail, preventing orphaned user accounts
	_, err = txQueries.CreateEmailVerificationToken(ctx, db.CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		Email:     user.Email,
		IpAddress: netip.Addr{}, // Empty value as IP not available in this context
		UserAgent: nil,           // Not available in service layer
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24h expiry per AUTH.md
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send verification email (outside transaction to avoid blocking)
	if s.emailService != nil {
		username := user.Username
		if err := s.emailService.SendVerificationEmail(ctx, user.Email, username, token); err != nil {
			// Log error but don't fail registration
			fmt.Printf("failed to send verification email: %v\n", err)
		}
	}

	return &user, nil
}

// VerifyEmail verifies a user's email address
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	tokenHash := s.tokenManager.HashRefreshToken(token)

	emailToken, err := s.repo.GetEmailVerificationToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token: %w", err)
	}

	// Mark token as used
	if err := s.repo.MarkEmailVerificationTokenUsed(ctx, emailToken.ID); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Update user's email_verified status
	if err := s.repo.UpdateUserEmailVerified(ctx, emailToken.UserID, true); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ResendVerification sends a new verification email
func (s *Service) ResendVerification(ctx context.Context, userID uuid.UUID) error {
	// Invalidate old verification tokens
	if err := s.repo.InvalidateUserEmailVerificationTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to invalidate old tokens: %w", err)
	}

	// Generate new token
	token, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Get user email
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Store new verification token
	tokenHash := s.tokenManager.HashRefreshToken(token)
	_, err = s.repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	// Send verification email
	if s.emailService != nil {
		if err := s.emailService.SendVerificationEmail(ctx, user.Email, user.Username, token); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}
	}

	return nil
}

// RegisterFromOIDCRequest contains data for OIDC user registration
type RegisterFromOIDCRequest struct {
	Username    string
	Email       string
	DisplayName *string
}

// RegisterFromOIDC creates a new user account from OIDC data.
// The user is created with a random unusable password and email is marked as verified
// (since OIDC provider already verified it).
func (s *Service) RegisterFromOIDC(ctx context.Context, req RegisterFromOIDCRequest) (*db.SharedUser, error) {
	// Generate a random unusable password
	// This password cannot be used to log in since it's random and not known to anyone
	randomPassword, err := crypto.GenerateSecureToken(64)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random password: %w", err)
	}

	// Hash the random password
	passwordHash, err := s.hasher.HashPassword(randomPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user in database with email already verified
	// OIDC providers verify email, so we don't need another verification step
	isActive := true
	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
		IsActive:     &isActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark email as verified since OIDC provider already verified it
	if err := s.repo.UpdateUserEmailVerified(ctx, user.ID, true); err != nil {
		// Log but don't fail - user is created, this is a secondary update
		fmt.Printf("failed to mark email as verified for OIDC user: %v\n", err)
	}

	return &user, nil
}

// ============================================================================
// Login & Logout
// ============================================================================

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, username, password string, ipAddress *netip.Addr, userAgent, deviceName, deviceFingerprint *string) (*LoginResponse, error) {
	// Helper to convert netip.Addr to net.IP for activity logging
	var activityIP net.IP
	if ipAddress != nil {
		activityIP = ipAddress.AsSlice()
	}

	// A7.5: Check account lockout (if enabled)
	if s.lockoutEnabled {
		since := time.Now().Add(-s.lockoutWindow)
		attemptCount, err := s.repo.CountFailedLoginAttemptsByUsername(ctx, username, since)
		if err != nil {
			// Log error but don't fail - continue with login attempt
			// This prevents lockout check failures from blocking legitimate logins
			fmt.Printf("failed to check lockout status: %v\n", err)
		} else if attemptCount >= int64(s.lockoutThreshold) {
			// Account is locked - log and return error
			_ = s.activityLogger.LogFailure(ctx, activity.LogFailureRequest{
				Username:     &username,
				Action:       activity.ActionUserLogin,
				ErrorMessage: "account locked due to too many failed attempts",
				IPAddress:    &activityIP,
				UserAgent:    userAgent,
			})
			return nil, fmt.Errorf("account locked due to too many failed login attempts. Please try again later")
		}
	}

	// Retrieve user by username or email
	// SECURITY: Always perform password hash comparison even if user not found
	// to prevent username enumeration via timing attacks
	user, err := s.repo.GetUserByUsername(ctx, username)
	userFound := (err == nil)
	if err != nil {
		// Try email if username not found
		user, err = s.repo.GetUserByEmail(ctx, username)
		userFound = (err == nil)
	}

	// Determine which hash to compare against
	// Use dummy hash if user not found to maintain constant-time behavior
	hashToCompare := dummyPasswordHash
	if userFound {
		hashToCompare = user.PasswordHash
	}

	// ALWAYS verify password (even if user not found) for timing attack mitigation
	// This ensures login timing is constant regardless of username validity
	match, err := s.hasher.VerifyPassword(password, hashToCompare)
	if err != nil {
		return nil, fmt.Errorf("password verification failed: %w", err)
	}

	// Check if user was found and password matched
	// Return same error message for both cases to prevent username enumeration
	if !userFound || !match {
		// A7.5: Record failed login attempt (if lockout enabled)
		if s.lockoutEnabled && ipAddress != nil {
			ipAddrStr := ipAddress.String()
			if err := s.repo.RecordFailedLoginAttempt(ctx, username, ipAddrStr); err != nil {
				// Log error but don't fail the login attempt
				fmt.Printf("failed to record failed login attempt: %v\n", err)
			}
		}

		// Log failed login attempt
		_ = s.activityLogger.LogFailure(ctx, activity.LogFailureRequest{
			Username:     &username,
			Action:       activity.ActionUserLogin,
			ErrorMessage: "invalid username or password",
			IPAddress:    &activityIP,
			UserAgent:    userAgent,
		})
		return nil, errors.New("invalid username or password")
	}

	// Check if account is active
	if user.IsActive != nil && !*user.IsActive {
		// Log failed login attempt
		_ = s.activityLogger.LogFailure(ctx, activity.LogFailureRequest{
			UserID:       &user.ID,
			Username:     &user.Username,
			Action:       activity.ActionUserLogin,
			ErrorMessage: "account is disabled",
			IPAddress:    &activityIP,
			UserAgent:    userAgent,
		})
		return nil, errors.New("account is disabled")
	}

	// Generate JWT access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (crypto/rand)
	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database (hashed with SHA-256)
	tokenHash := s.tokenManager.HashRefreshToken(refreshToken)
	_, err = s.repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:            user.ID,
		TokenHash:         tokenHash,
		TokenType:         "refresh",
		DeviceName:        deviceName,
		DeviceFingerprint: deviceFingerprint,
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		ExpiresAt:         time.Now().Add(s.refreshExpiry),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// A7.5: Clear failed login attempts on successful login (if lockout enabled)
	if s.lockoutEnabled {
		if err := s.repo.ClearFailedLoginAttemptsByUsername(ctx, username); err != nil {
			// Log error but don't fail login
			fmt.Printf("failed to clear failed login attempts: %v\n", err)
		}
	}

	// Update last login time
	if err := s.repo.UpdateUserLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	// Log successful login
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		UserID:       user.ID,
		Username:     user.Username,
		Action:       activity.ActionUserLogin,
		ResourceType: activity.ResourceTypeUser,
		ResourceID:   user.ID,
		IPAddress:    activityIP,
		UserAgent:    ptrToString(userAgent),
		Metadata: map[string]interface{}{
			"device_name": ptrToString(deviceName),
		},
	})

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}

// Logout invalidates a refresh token
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := s.tokenManager.HashRefreshToken(refreshToken)
	return s.repo.RevokeAuthTokenByHash(ctx, tokenHash)
}

// LogoutAll revokes all refresh tokens for a user
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.repo.RevokeAllUserAuthTokens(ctx, userID)
}

// RefreshToken generates a new access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	tokenHash := s.tokenManager.HashRefreshToken(refreshToken)

	// Retrieve token from database
	authToken, err := s.repo.GetAuthTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, authToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new JWT access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Update token last_used_at
	if err := s.repo.UpdateAuthTokenLastUsed(ctx, authToken.ID); err != nil {
		// Log error but don't fail refresh
		fmt.Printf("failed to update token last used: %v\n", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Return same refresh token
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}

// CreateSessionForUser creates a new session for an authenticated user.
// This is used for OIDC login where the user has already been verified by the OIDC provider.
func (s *Service) CreateSessionForUser(ctx context.Context, userID uuid.UUID, ipAddress *netip.Addr, userAgent, deviceName *string) (*LoginResponse, error) {
	// Helper to convert netip.Addr to net.IP for activity logging
	var activityIP net.IP
	if ipAddress != nil {
		activityIP = ipAddress.AsSlice()
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if user.IsActive != nil && !*user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Generate JWT access token
	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (crypto/rand)
	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash refresh token for storage
	tokenHash := s.tokenManager.HashRefreshToken(refreshToken)

	// Store refresh token in database
	var ipForStorage *netip.Addr
	if ipAddress != nil {
		ipForStorage = ipAddress
	}

	_, err = s.repo.CreateAuthToken(ctx, CreateAuthTokenParams{
		UserID:            user.ID,
		TokenHash:         tokenHash,
		ExpiresAt:         time.Now().Add(s.refreshExpiry),
		UserAgent:         userAgent,
		DeviceName:        deviceName,
		DeviceFingerprint: nil,
		IPAddress:         ipForStorage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Update last login time
	if err := s.repo.UpdateUserLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	// Log successful login
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		UserID:       user.ID,
		Username:     user.Username,
		Action:       activity.ActionUserLogin,
		ResourceType: activity.ResourceTypeUser,
		ResourceID:   user.ID,
		IPAddress:    activityIP,
		UserAgent:    ptrToString(userAgent),
		Metadata: map[string]interface{}{
			"device_name": ptrToString(deviceName),
			"oidc_login":  true,
		},
	})

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}

// ============================================================================
// Password Management
// ============================================================================

// ChangePassword updates a user's password
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	match, err := s.hasher.VerifyPassword(oldPassword, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("password verification failed: %w", err)
	}
	if !match {
		return errors.New("invalid current password")
	}

	// Hash new password
	newPasswordHash, err := s.hasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Begin transaction to ensure password update + token revocation are atomic
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txQueries := db.New(tx)

	// Update password in database
	if err := txQueries.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:           userID,
		PasswordHash: newPasswordHash,
	}); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens (force re-login on all devices for security)
	if err := txQueries.RevokeAllUserAuthTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log password change (outside transaction to avoid blocking)
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		UserID:       user.ID,
		Username:     user.Username,
		Action:       activity.ActionUserPasswordReset,
		ResourceType: activity.ResourceTypeUser,
		ResourceID:   user.ID,
	})

	return nil
}

// RequestPasswordReset generates a password reset token
// RequestPasswordReset initiates a password reset by sending a reset email.
// SECURITY: Never returns token to prevent information disclosure about email existence.
// Always returns success (nil error) to prevent email enumeration attacks.
func (s *Service) RequestPasswordReset(ctx context.Context, email string, ipAddress *netip.Addr, userAgent *string) error {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Silently succeed - don't reveal if email doesn't exist
		// This prevents email enumeration attacks
		return nil
	}

	// Invalidate old reset tokens
	if err := s.repo.InvalidateUserPasswordResetTokens(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to invalidate old tokens: %w", err)
	}

	// Generate reset token
	token, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Store reset token in database (hashed)
	tokenHash := s.tokenManager.HashRefreshToken(token)
	_, err = s.repo.CreatePasswordResetToken(ctx, CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1h expiry per AUTH.md
	})
	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send password reset email asynchronously (don't block request)
	// Email is only way to receive the token - never returned to API caller
	if s.emailService != nil {
		go func() {
			// Use background context with timeout for async operation
			emailCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := s.emailService.SendPasswordResetEmail(emailCtx, user.Email, user.Username, token); err != nil {
				// Log error for monitoring but don't expose to caller
				fmt.Printf("failed to send password reset email to %s: %v\n", user.Email, err)
			}
		}()
	}

	return nil
}

// ResetPassword resets a password using a reset token
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := s.tokenManager.HashRefreshToken(token)

	// Retrieve token from database
	resetToken, err := s.repo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	newPasswordHash, err := s.hasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Begin transaction to ensure password update + token marking + token revocation are atomic
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txQueries := db.New(tx)

	// Update password
	if err := txQueries.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:           resetToken.UserID,
		PasswordHash: newPasswordHash,
	}); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := txQueries.MarkPasswordResetTokenUsed(ctx, resetToken.ID); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Revoke all refresh tokens (force re-login)
	if err := txQueries.RevokeAllUserAuthTokens(ctx, resetToken.UserID); err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log password reset (outside transaction to avoid blocking)
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		UserID:       resetToken.UserID,
		Action:       activity.ActionUserPasswordReset,
		ResourceType: activity.ResourceTypeUser,
		ResourceID:   resetToken.UserID,
	})

	return nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// ptrToString safely dereferences a string pointer
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
