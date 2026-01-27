// Package auth provides authentication services for Jellyfin Go.
package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"time"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// Service implements domain.AuthService.
type Service struct {
	users              domain.UserRepository
	sessions           domain.SessionRepository
	passwords          domain.PasswordService
	tokens             domain.TokenService
	maxSessionsPerUser int
	accessDuration     time.Duration
	refreshDuration    time.Duration
}

// newService creates a new authentication service.
// Use NewService from module.go for fx integration.
func newService(
	users domain.UserRepository,
	sessions domain.SessionRepository,
	passwords domain.PasswordService,
	tokens domain.TokenService,
	maxSessionsPerUser int,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) *Service {
	if accessDuration <= 0 {
		accessDuration = 15 * time.Minute
	}
	if refreshDuration <= 0 {
		refreshDuration = 7 * 24 * time.Hour
	}

	return &Service{
		users:              users,
		sessions:           sessions,
		passwords:          passwords,
		tokens:             tokens,
		maxSessionsPerUser: maxSessionsPerUser,
		accessDuration:     accessDuration,
		refreshDuration:    refreshDuration,
	}
}

// Login authenticates a user with username and password.
func (s *Service) Login(ctx context.Context, params domain.LoginParams) (*domain.AuthResult, error) {
	// Find user by username
	user, err := s.users.GetByUsername(ctx, params.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Warn("login attempt for non-existent user",
				slog.String("username", params.Username))
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is disabled
	if user.IsDisabled {
		slog.Warn("login attempt for disabled user",
			slog.String("username", params.Username),
			slog.String("user_id", user.ID.String()))
		return nil, domain.ErrUserDisabled
	}

	// Verify password
	if user.PasswordHash == nil {
		slog.Warn("login attempt for user without password (OIDC-only)",
			slog.String("username", params.Username))
		return nil, domain.ErrInvalidCredentials
	}

	if err := s.passwords.Verify(params.Password, *user.PasswordHash); err != nil {
		slog.Warn("invalid password",
			slog.String("username", params.Username),
			slog.String("user_id", user.ID.String()))
		return nil, domain.ErrInvalidCredentials
	}

	// Check max sessions limit
	if s.maxSessionsPerUser > 0 {
		count, err := s.sessions.CountByUser(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to count sessions: %w", err)
		}
		if count >= int64(s.maxSessionsPerUser) {
			// Delete oldest sessions to make room
			sessions, err := s.sessions.ListByUser(ctx, user.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to list sessions: %w", err)
			}
			// Sessions are ordered by created_at DESC, so delete from the end
			for i := s.maxSessionsPerUser - 1; i < len(sessions); i++ {
				if err := s.sessions.Delete(ctx, sessions[i].ID); err != nil {
					slog.Warn("failed to delete old session",
						slog.String("session_id", sessions[i].ID.String()),
						slog.Any("error", err))
				}
			}
		}
	}

	// Create session and tokens
	result, err := s.createSession(ctx, user, params.DeviceID, params.DeviceName,
		params.ClientName, params.ClientVersion, params.IPAddress)
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := s.users.UpdateLastLogin(ctx, user.ID); err != nil {
		slog.Warn("failed to update last login",
			slog.String("user_id", user.ID.String()),
			slog.Any("error", err))
	}

	slog.Info("user logged in",
		slog.String("user_id", user.ID.String()),
		slog.String("username", user.Username),
		slog.String("session_id", result.SessionID.String()))

	return result, nil
}

// Logout invalidates a session by its access token.
func (s *Service) Logout(ctx context.Context, accessToken string) error {
	tokenHash := s.tokens.HashToken(accessToken)

	if err := s.sessions.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	slog.Info("session logged out")
	return nil
}

// LogoutAll invalidates all sessions for a user.
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	if err := s.sessions.DeleteByUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	slog.Info("all sessions logged out",
		slog.String("user_id", userID.String()))
	return nil
}

// RefreshToken exchanges a refresh token for new access/refresh tokens.
func (s *Service) RefreshToken(ctx context.Context, params domain.RefreshParams) (*domain.AuthResult, error) {
	refreshTokenHash := s.tokens.HashToken(params.RefreshToken)

	// Find session by refresh token
	session, err := s.sessions.GetByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check refresh token expiry
	if session.RefreshExpiresAt != nil && session.RefreshExpiresAt.Before(time.Now()) {
		// Clean up expired session - error ignored as this is best-effort cleanup
		_ = s.sessions.Delete(ctx, session.ID) //nolint:errcheck
		return nil, domain.ErrSessionExpired
	}

	// Get user
	user, err := s.users.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is disabled
	if user.IsDisabled {
		// Clean up session for disabled user - error ignored as this is best-effort cleanup
		_ = s.sessions.Delete(ctx, session.ID) //nolint:errcheck
		return nil, domain.ErrUserDisabled
	}

	// Generate new tokens
	now := time.Now()
	accessExpiry := now.Add(s.accessDuration)
	refreshExpiry := now.Add(s.refreshDuration)

	accessToken, err := s.tokens.GenerateAccessToken(domain.TokenClaims{
		UserID:    user.ID,
		SessionID: session.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		IssuedAt:  now,
		ExpiresAt: accessExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update session with new tokens
	newAccessTokenHash := s.tokens.HashToken(accessToken)
	newRefreshTokenHash := s.tokens.HashToken(newRefreshToken)

	// Delete old session and create new one (atomic token rotation)
	if err := s.sessions.Delete(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to delete old session: %w", err)
	}

	newSession, err := s.sessions.Create(ctx, domain.CreateSessionParams{
		UserID:           user.ID,
		TokenHash:        newAccessTokenHash,
		RefreshTokenHash: &newRefreshTokenHash,
		DeviceID:         session.DeviceID,
		DeviceName:       session.DeviceName,
		ClientName:       session.ClientName,
		ClientVersion:    session.ClientVersion,
		IPAddress:        params.IPAddress,
		ExpiresAt:        accessExpiry,
		RefreshExpiresAt: &refreshExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new session: %w", err)
	}

	slog.Info("token refreshed",
		slog.String("user_id", user.ID.String()),
		slog.String("session_id", newSession.ID.String()))

	return &domain.AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExpiry,
		SessionID:    newSession.ID,
	}, nil
}

// ValidateToken validates an access token and returns the claims.
func (s *Service) ValidateToken(ctx context.Context, accessToken string) (*domain.TokenClaims, error) {
	// Parse and validate JWT
	claims, err := s.tokens.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	// Verify session still exists and is valid
	tokenHash := s.tokens.HashToken(accessToken)
	exists, err := s.sessions.Exists(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to check session: %w", err)
	}
	if !exists {
		return nil, domain.ErrSessionNotFound
	}

	return claims, nil
}

// GetSession retrieves session with user info by access token.
func (s *Service) GetSession(ctx context.Context, accessToken string) (*domain.SessionWithUser, error) {
	tokenHash := s.tokens.HashToken(accessToken)

	session, err := s.sessions.GetWithUser(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ChangePassword changes a user's password (requires current password).
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if user.PasswordHash == nil {
		return errors.New("user has no password set")
	}

	if err := s.passwords.Verify(currentPassword, *user.PasswordHash); err != nil {
		return domain.ErrInvalidCredentials
	}

	// Hash and set new password
	newHash, err := s.passwords.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.users.SetPassword(ctx, userID, newHash); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	// Invalidate all sessions (security best practice)
	if err := s.sessions.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to delete sessions after password change",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
	}

	slog.Info("password changed",
		slog.String("user_id", userID.String()))

	return nil
}

// ResetPassword sets a new password (admin operation, no current password required).
func (s *Service) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Hash and set new password
	newHash, err := s.passwords.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.users.SetPassword(ctx, userID, newHash); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	// Invalidate all sessions
	if err := s.sessions.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to delete sessions after password reset",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
	}

	slog.Info("password reset",
		slog.String("user_id", userID.String()))

	return nil
}

// createSession creates a new session with tokens for a user.
func (s *Service) createSession(
	ctx context.Context,
	user *domain.User,
	deviceID, deviceName, clientName, clientVersion *string,
	ipAddress *netip.Addr,
) (*domain.AuthResult, error) {
	now := time.Now()
	accessExpiry := now.Add(s.accessDuration)
	refreshExpiry := now.Add(s.refreshDuration)

	// Create session first to get ID
	sessionID := uuid.New()

	// Generate tokens
	accessToken, err := s.tokens.GenerateAccessToken(domain.TokenClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		IssuedAt:  now,
		ExpiresAt: accessExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash tokens for storage
	accessTokenHash := s.tokens.HashToken(accessToken)
	refreshTokenHash := s.tokens.HashToken(refreshToken)

	// Create session
	session, err := s.sessions.Create(ctx, domain.CreateSessionParams{
		UserID:           user.ID,
		TokenHash:        accessTokenHash,
		RefreshTokenHash: &refreshTokenHash,
		DeviceID:         deviceID,
		DeviceName:       deviceName,
		ClientName:       clientName,
		ClientVersion:    clientVersion,
		IPAddress:        ipAddress,
		ExpiresAt:        accessExpiry,
		RefreshExpiresAt: &refreshExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
		SessionID:    session.ID,
	}, nil
}

// Ensure Service implements domain.AuthService.
var _ domain.AuthService = (*Service)(nil)
