// Package session provides session management services.
package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	// ErrSessionNotFound indicates the session was not found or is invalid.
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired indicates the session has expired.
	ErrSessionExpired = errors.New("session expired")
	// ErrSessionInactive indicates the session is no longer active.
	ErrSessionInactive = errors.New("session inactive")
	// ErrTooManySessions indicates the user has too many active sessions.
	ErrTooManySessions = errors.New("too many active sessions")
)

// Service provides session management operations.
type Service struct {
	queries            *db.Queries
	logger             *slog.Logger
	sessionDuration    time.Duration
	maxSessionsPerUser int // 0 = unlimited
}

// NewService creates a new session service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries:            queries,
		logger:             logger.With(slog.String("service", "session")),
		sessionDuration:    24 * time.Hour, // Default 24 hours
		maxSessionsPerUser: 0,              // Unlimited by default
	}
}

// SetSessionDuration sets the session duration.
func (s *Service) SetSessionDuration(d time.Duration) {
	s.sessionDuration = d
}

// SetMaxSessionsPerUser sets the maximum sessions per user (0 = unlimited).
func (s *Service) SetMaxSessionsPerUser(maxSessions int) {
	s.maxSessionsPerUser = maxSessions
}

// CreateParams contains parameters for creating a session.
type CreateParams struct {
	UserID        uuid.UUID
	ProfileID     *uuid.UUID
	DeviceName    *string
	DeviceType    *string
	ClientName    *string
	ClientVersion *string
	IPAddress     netip.Addr
	UserAgent     *string
}

// CreateResult contains the result of session creation.
type CreateResult struct {
	Session *db.Session
	Token   string // Raw token (not hashed) - only returned on creation
}

// Create creates a new session and returns the session with the raw token.
func (s *Service) Create(ctx context.Context, params CreateParams) (*CreateResult, error) {
	// Check session limit
	if s.maxSessionsPerUser > 0 {
		count, err := s.queries.CountActiveSessionsByUser(ctx, params.UserID)
		if err != nil {
			return nil, fmt.Errorf("count sessions: %w", err)
		}
		if int(count) >= s.maxSessionsPerUser {
			return nil, ErrTooManySessions
		}
	}

	// Generate secure token
	token, err := generateToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	// Hash token for storage
	tokenHash := hashToken(token)

	// Build profile ID
	var profileID pgtype.UUID
	if params.ProfileID != nil {
		profileID = pgtype.UUID{Bytes: *params.ProfileID, Valid: true}
	}

	// Create session
	session, err := s.queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:        params.UserID,
		ProfileID:     profileID,
		TokenHash:     tokenHash,
		DeviceName:    params.DeviceName,
		DeviceType:    params.DeviceType,
		ClientName:    params.ClientName,
		ClientVersion: params.ClientVersion,
		IpAddress:     params.IPAddress,
		UserAgent:     params.UserAgent,
		ExpiresAt:     time.Now().Add(s.sessionDuration),
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	s.logger.Info("Session created",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", session.UserID.String()),
	)

	return &CreateResult{
		Session: &session,
		Token:   token,
	}, nil
}

// ValidateToken validates a token and returns the associated session.
func (s *Service) ValidateToken(ctx context.Context, token string) (*db.Session, error) {
	tokenHash := hashToken(token)

	session, err := s.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if !session.IsActive {
		return nil, ErrSessionInactive
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// GetByID retrieves a session by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.Session, error) {
	session, err := s.queries.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return &session, nil
}

// ListByUser returns all sessions for a user.
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]db.Session, error) {
	sessions, err := s.queries.ListSessionsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	return sessions, nil
}

// UpdateActivity updates the last activity time for a session.
func (s *Service) UpdateActivity(ctx context.Context, sessionID uuid.UUID, profileID *uuid.UUID) error {
	var pid pgtype.UUID
	if profileID != nil {
		pid = pgtype.UUID{Bytes: *profileID, Valid: true}
	}

	if err := s.queries.UpdateSessionActivity(ctx, db.UpdateSessionActivityParams{
		ID:        sessionID,
		ProfileID: pid,
	}); err != nil {
		return fmt.Errorf("update activity: %w", err)
	}
	return nil
}

// Deactivate deactivates a session (soft delete).
func (s *Service) Deactivate(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeactivateSession(ctx, id); err != nil {
		return fmt.Errorf("deactivate session: %w", err)
	}

	s.logger.Info("Session deactivated",
		slog.String("session_id", id.String()),
	)

	return nil
}

// DeactivateAllForUser deactivates all sessions for a user.
func (s *Service) DeactivateAllForUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.queries.DeactivateUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("deactivate user sessions: %w", err)
	}

	s.logger.Info("All sessions deactivated for user",
		slog.String("user_id", userID.String()),
	)

	return nil
}

// Delete hard deletes a session.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteSession(ctx, id); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	s.logger.Info("Session deleted",
		slog.String("session_id", id.String()),
	)

	return nil
}

// CleanupExpired removes all expired sessions.
func (s *Service) CleanupExpired(ctx context.Context) error {
	if err := s.queries.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("cleanup expired: %w", err)
	}

	s.logger.Debug("Expired sessions cleaned up")

	return nil
}

// generateToken generates a cryptographically secure random token.
func generateToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashToken creates a SHA-256 hash of a token for storage.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(h[:])
}
