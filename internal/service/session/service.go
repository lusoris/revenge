package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/errors"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/zap"
)

// Service handles session management operations
type Service struct {
	repo           Repository
	logger         *zap.Logger
	tokenLength    int
	expiry         time.Duration
	refreshExpiry  time.Duration
	maxPerUser     int
}

// DeviceInfo contains device metadata for session creation
type DeviceInfo struct {
	DeviceName *string
	UserAgent  *string
	IPAddress  *netip.Addr
}

// SessionInfo represents a session for API responses
type SessionInfo struct {
	ID             uuid.UUID
	DeviceName     *string
	IPAddress      *string
	UserAgent      *string
	CreatedAt      time.Time
	LastActivityAt time.Time
	ExpiresAt      time.Time
	IsActive       bool
	IsCurrent      bool
}

// CreateSession creates a new session for a user
func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID, deviceInfo DeviceInfo, scopes []string) (string, string, error) {
	// Check session limit
	count, err := s.repo.CountActiveUserSessions(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to count user sessions: %w", err)
	}

	if int(count) >= s.maxPerUser {
		s.logger.Warn("User has too many active sessions",
			zap.String("user_id", userID.String()),
			zap.Int64("count", count),
			zap.Int("max", s.maxPerUser))
		// Optionally revoke oldest session here
	}

	// Generate session token
	token, tokenHash, err := s.generateToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate session token: %w", err)
	}

	// Generate refresh token
	refreshToken, refreshTokenHash, err := s.generateToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	_, err = s.repo.CreateSession(ctx, CreateSessionParams{
		UserID:           userID,
		TokenHash:        tokenHash,
		RefreshTokenHash: &refreshTokenHash,
		DeviceName:       deviceInfo.DeviceName,
		UserAgent:        deviceInfo.UserAgent,
		IPAddress:        deviceInfo.IPAddress,
		Scopes:           scopes,
		ExpiresAt:        time.Now().Add(s.expiry),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("Session created",
		zap.String("user_id", userID.String()),
		zap.Int("active_sessions", int(count)+1))

	return token, refreshToken, nil
}

// ValidateSession validates a session token
func (s *Service) ValidateSession(ctx context.Context, token string) (*db.SharedSession, error) {
	tokenHash := s.hashToken(token)

	session, err := s.repo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, errors.ErrUnauthorized
	}

	// Update activity
	if err := s.repo.UpdateSessionActivity(ctx, session.ID); err != nil {
		s.logger.Warn("Failed to update session activity",
			zap.String("session_id", session.ID.String()),
			zap.Error(err))
	}

	return session, nil
}

// RefreshSession exchanges a refresh token for a new session token
func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenHash := s.hashToken(refreshToken)

	session, err := s.repo.GetSessionByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return "", fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	if session == nil {
		return "", errors.ErrUnauthorized
	}

	// Generate new session token
	newToken, newTokenHash, err := s.generateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate new session token: %w", err)
	}

	// Update session with new token hash
	// NOTE: This would require an UpdateSessionTokenHash query
	// For now, we'll create a new session and revoke the old one
	reason := "Refresh token rotation"
	if err := s.repo.RevokeSession(ctx, session.ID, &reason); err != nil {
		return "", fmt.Errorf("failed to revoke old session: %w", err)
	}

	// Create new session with same metadata
	var deviceName, userAgent *string
	var ipAddr *netip.Addr

	if session.DeviceName != nil {
		deviceName = session.DeviceName
	}
	if session.UserAgent != nil {
		userAgent = session.UserAgent
	}
	// IpAddress is netip.Addr (not nullable), convert to pointer
	if !session.IpAddress.IsUnspecified() {
		ipAddr = &session.IpAddress
	}

	_, err = s.repo.CreateSession(ctx, CreateSessionParams{
		UserID:           session.UserID,
		TokenHash:        newTokenHash,
		RefreshTokenHash: &refreshTokenHash, // Reuse refresh token
		DeviceName:       deviceName,
		UserAgent:        userAgent,
		IPAddress:        ipAddr,
		Scopes:           session.Scopes,
		ExpiresAt:        time.Now().Add(s.expiry),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create refreshed session: %w", err)
	}

	s.logger.Info("Session refreshed", zap.String("user_id", session.UserID.String()))

	return newToken, nil
}

// ListUserSessions lists all active sessions for a user
func (s *Service) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error) {
	sessions, err := s.repo.ListUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user sessions: %w", err)
	}

	result := make([]SessionInfo, len(sessions))
	for i, session := range sessions {
		result[i] = s.sessionToInfo(&session, false)
	}

	return result, nil
}

// RevokeSession revokes a specific session
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	reason := "User logout"
	if err := s.repo.RevokeSession(ctx, sessionID, &reason); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	s.logger.Info("Session revoked", zap.String("session_id", sessionID.String()))
	return nil
}

// RevokeAllUserSessions revokes all sessions for a user (logout everywhere)
func (s *Service) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	reason := "User logout all"
	if err := s.repo.RevokeAllUserSessions(ctx, userID, &reason); err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	s.logger.Info("All user sessions revoked", zap.String("user_id", userID.String()))
	return nil
}

// RevokeAllUserSessionsExcept revokes all sessions except the current one
func (s *Service) RevokeAllUserSessionsExcept(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) error {
	reason := "User logout all others"
	if err := s.repo.RevokeAllUserSessionsExcept(ctx, userID, currentSessionID, &reason); err != nil {
		return fmt.Errorf("failed to revoke other user sessions: %w", err)
	}

	s.logger.Info("Other user sessions revoked",
		zap.String("user_id", userID.String()),
		zap.String("kept_session_id", currentSessionID.String()))
	return nil
}

// CleanupExpiredSessions removes old expired and revoked sessions
func (s *Service) CleanupExpiredSessions(ctx context.Context) (int, error) {
	// Delete expired sessions
	if err := s.repo.DeleteExpiredSessions(ctx); err != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	// Delete revoked sessions
	if err := s.repo.DeleteRevokedSessions(ctx); err != nil {
		return 0, fmt.Errorf("failed to delete revoked sessions: %w", err)
	}

	s.logger.Info("Session cleanup completed")
	return 0, nil // TODO: Return actual count
}

// Helper methods

func (s *Service) generateToken() (string, string, error) {
	token := make([]byte, s.tokenLength)
	if _, err := rand.Read(token); err != nil {
		return "", "", err
	}

	tokenStr := hex.EncodeToString(token)
	tokenHash := s.hashToken(tokenStr)

	return tokenStr, tokenHash, nil
}

func (s *Service) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *Service) sessionToInfo(session *db.SharedSession, isCurrent bool) SessionInfo {
	info := SessionInfo{
		ID:             session.ID,
		CreatedAt:      session.CreatedAt,
		LastActivityAt: session.LastActivityAt,
		ExpiresAt:      session.ExpiresAt,
		IsActive:       !session.RevokedAt.Valid && session.ExpiresAt.After(time.Now()),
		IsCurrent:      isCurrent,
	}

	if session.DeviceName != nil {
		info.DeviceName = session.DeviceName
	}
	if session.UserAgent != nil {
		info.UserAgent = session.UserAgent
	}
	if !session.IpAddress.IsUnspecified() {
		ipStr := session.IpAddress.String()
		info.IPAddress = &ipStr
	}

	return info
}
