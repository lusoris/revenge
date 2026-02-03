package session

import (
	"context"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Repository defines the data access interface for session operations
type Repository interface {
	// Session CRUD
	CreateSession(ctx context.Context, params CreateSessionParams) (db.SharedSession, error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*db.SharedSession, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*db.SharedSession, error)
	GetSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*db.SharedSession, error)

	// Session Listing
	ListUserSessions(ctx context.Context, userID uuid.UUID) ([]db.SharedSession, error)
	ListAllUserSessions(ctx context.Context, userID uuid.UUID) ([]db.SharedSession, error)
	CountActiveUserSessions(ctx context.Context, userID uuid.UUID) (int64, error)

	// Session Updates
	UpdateSessionActivity(ctx context.Context, sessionID uuid.UUID) error
	UpdateSessionActivityByTokenHash(ctx context.Context, tokenHash string) error

	// Session Revocation
	RevokeSession(ctx context.Context, sessionID uuid.UUID, reason *string) error
	RevokeSessionByTokenHash(ctx context.Context, tokenHash string, reason *string) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID, reason *string) error
	RevokeAllUserSessionsExcept(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID, reason *string) error

	// Cleanup
	DeleteExpiredSessions(ctx context.Context) error
	DeleteRevokedSessions(ctx context.Context) error
	GetInactiveSessions(ctx context.Context, inactiveSince time.Time) ([]db.SharedSession, error)
	RevokeInactiveSessions(ctx context.Context, inactiveSince time.Time) error
}

// CreateSessionParams parameters for creating a session
type CreateSessionParams struct {
	UserID           uuid.UUID
	TokenHash        string
	RefreshTokenHash *string
	IPAddress        *netip.Addr
	UserAgent        *string
	DeviceName       *string
	Scopes           []string
	ExpiresAt        time.Time
}
