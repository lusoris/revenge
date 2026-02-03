package auth

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Service implements auth business logic
type Service struct {
	repo          Repository
	tokenManager  TokenManager
	hasher        *crypto.PasswordHasher
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
}

// NewService creates a new auth service
func NewService(repo Repository, tokenManager TokenManager, jwtExpiry, refreshExpiry time.Duration) *Service {
	return &Service{
		repo:          repo,
		tokenManager:  tokenManager,
		hasher:        crypto.NewPasswordHasher(),
		jwtExpiry:     jwtExpiry,
		refreshExpiry: refreshExpiry,
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

	// Create user in database
	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate email verification token
	token, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store verification token in database
	tokenHash := s.tokenManager.HashRefreshToken(token)
	_, err = s.repo.CreateEmailVerificationToken(ctx, CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24h expiry per AUTH.md
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// TODO: Send verification email (requires email service)
	// For now, token is generated but not sent

	return &user, nil
}

// VerifyEmail verifies a user's email address
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	tokenHash := s.tokenManager.HashRefreshToken(token)

	// Retrieve token from database
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

	// TODO: Send verification email
	return nil
}

// ============================================================================
// Login & Logout
// ============================================================================

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, username, password string, ipAddress *netip.Addr, userAgent, deviceName, deviceFingerprint *string) (*LoginResponse, error) {
	// Retrieve user by username or email
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		// Try email if username not found
		user, err = s.repo.GetUserByEmail(ctx, username)
		if err != nil {
			return nil, errors.New("invalid username or password")
		}
	}

	// Verify password using Argon2id
	match, err := s.hasher.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("password verification failed: %w", err)
	}
	if !match {
		return nil, errors.New("invalid username or password")
	}

	// Check if account is active
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

	// Update last login time
	if err := s.repo.UpdateUserLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("failed to update last login: %v\n", err)
	}

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

	// Update password in database
	if err := s.repo.UpdateUserPassword(ctx, userID, newPasswordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens (force re-login on all devices for security)
	if err := s.repo.RevokeAllUserAuthTokens(ctx, userID); err != nil {
		// Log error but don't fail password change
		fmt.Printf("failed to revoke tokens: %v\n", err)
	}

	return nil
}

// RequestPasswordReset generates a password reset token
func (s *Service) RequestPasswordReset(ctx context.Context, email string, ipAddress *netip.Addr, userAgent *string) (string, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists (security)
		// Return success even if user not found
		return "", nil
	}

	// Invalidate old reset tokens
	if err := s.repo.InvalidateUserPasswordResetTokens(ctx, user.ID); err != nil {
		return "", fmt.Errorf("failed to invalidate old tokens: %w", err)
	}

	// Generate reset token
	token, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
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
		return "", fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send reset email
	return token, nil
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

	// Update password
	if err := s.repo.UpdateUserPassword(ctx, resetToken.UserID, newPasswordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.repo.MarkPasswordResetTokenUsed(ctx, resetToken.ID); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Revoke all refresh tokens (force re-login)
	if err := s.repo.RevokeAllUserAuthTokens(ctx, resetToken.UserID); err != nil {
		// Log error but don't fail reset
		fmt.Printf("failed to revoke tokens: %v\n", err)
	}

	return nil
}

// ============================================================================
// Helper Functions
// ============================================================================
