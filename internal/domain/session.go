// Package domain contains core business entities and repository interfaces.
package domain

import (
	"context"
	"net/netip"
	"time"

	"github.com/google/uuid"
)

// Session represents an active user session.
type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	TokenHash        string
	RefreshTokenHash *string
	DeviceID         *string
	DeviceName       *string
	ClientName       *string
	ClientVersion    *string
	IPAddress        *netip.Addr
	ExpiresAt        time.Time
	RefreshExpiresAt *time.Time
	CreatedAt        time.Time
}

// SessionWithUser includes session data joined with user information.
type SessionWithUser struct {
	Session
	Username    string
	Email       *string
	DisplayName *string
	IsAdmin     bool
	IsDisabled  bool
}

// CreateSessionParams contains parameters for creating a new session.
type CreateSessionParams struct {
	UserID           uuid.UUID
	TokenHash        string
	RefreshTokenHash *string
	DeviceID         *string
	DeviceName       *string
	ClientName       *string
	ClientVersion    *string
	IPAddress        *netip.Addr
	ExpiresAt        time.Time
	RefreshExpiresAt *time.Time
}

// SessionRepository defines the interface for session data access.
type SessionRepository interface {
	// GetByID retrieves a session by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Session, error)

	// GetByTokenHash retrieves a session by the hashed access token.
	GetByTokenHash(ctx context.Context, tokenHash string) (*Session, error)

	// GetByRefreshTokenHash retrieves a session by the hashed refresh token.
	GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*Session, error)

	// GetWithUser retrieves a session with joined user information.
	GetWithUser(ctx context.Context, tokenHash string) (*SessionWithUser, error)

	// ListByUser retrieves all sessions for a specific user.
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Session, error)

	// Create creates a new session.
	Create(ctx context.Context, params CreateSessionParams) (*Session, error)

	// UpdateRefreshToken updates the session's refresh token.
	UpdateRefreshToken(ctx context.Context, id uuid.UUID, refreshTokenHash string, expiresAt time.Time) error

	// Delete removes a session by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByTokenHash removes a session by its token hash.
	DeleteByTokenHash(ctx context.Context, tokenHash string) error

	// DeleteByUser removes all sessions for a user.
	DeleteByUser(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired removes all expired sessions and returns the count.
	DeleteExpired(ctx context.Context) (int64, error)

	// CountByUser returns the number of sessions for a user.
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	// Exists checks if a valid session exists for a token hash.
	Exists(ctx context.Context, tokenHash string) (bool, error)
}
