package session

import (
	"context"
	"database/sql"
	"errors"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// RepositoryPG implements Repository using PostgreSQL via sqlc
type RepositoryPG struct {
	queries *db.Queries
}

// CreateSession creates a new session
func (r *RepositoryPG) CreateSession(ctx context.Context, params CreateSessionParams) (db.SharedSession, error) {
	var ipAddr netip.Addr
	if params.IPAddress != nil {
		ipAddr = *params.IPAddress
	}

	return r.queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:           params.UserID,
		TokenHash:        params.TokenHash,
		RefreshTokenHash: params.RefreshTokenHash,
		IpAddress:        ipAddr,
		UserAgent:        params.UserAgent,
		DeviceName:       params.DeviceName,
		Scopes:           params.Scopes,
		ExpiresAt:        params.ExpiresAt,
	})
}

// GetSessionByTokenHash retrieves a session by token hash
func (r *RepositoryPG) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*db.SharedSession, error) {
	session, err := r.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// GetSessionByID retrieves a session by ID
func (r *RepositoryPG) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*db.SharedSession, error) {
	session, err := r.queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// GetSessionByRefreshTokenHash retrieves a session by refresh token hash
func (r *RepositoryPG) GetSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*db.SharedSession, error) {
	session, err := r.queries.GetSessionByRefreshTokenHash(ctx, &refreshTokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// ListUserSessions lists active sessions for a user
func (r *RepositoryPG) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]db.SharedSession, error) {
	return r.queries.ListUserSessions(ctx, userID)
}

// ListAllUserSessions lists all sessions (including expired) for a user
func (r *RepositoryPG) ListAllUserSessions(ctx context.Context, userID uuid.UUID) ([]db.SharedSession, error) {
	return r.queries.ListAllUserSessions(ctx, userID)
}

// CountActiveUserSessions counts active sessions for a user
func (r *RepositoryPG) CountActiveUserSessions(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountActiveUserSessions(ctx, userID)
}

// CountAllActiveSessions counts all active sessions across all users
func (r *RepositoryPG) CountAllActiveSessions(ctx context.Context) (int64, error) {
	return r.queries.CountAllActiveSessions(ctx)
}

// UpdateSessionActivity updates the last activity timestamp
func (r *RepositoryPG) UpdateSessionActivity(ctx context.Context, sessionID uuid.UUID) error {
	return r.queries.UpdateSessionActivity(ctx, sessionID)
}

// UpdateSessionActivityByTokenHash updates activity by token hash
func (r *RepositoryPG) UpdateSessionActivityByTokenHash(ctx context.Context, tokenHash string) error {
	return r.queries.UpdateSessionActivityByTokenHash(ctx, tokenHash)
}

// RevokeSession revokes a session by ID
func (r *RepositoryPG) RevokeSession(ctx context.Context, sessionID uuid.UUID, reason *string) error {
	return r.queries.RevokeSession(ctx, db.RevokeSessionParams{
		ID:     sessionID,
		Reason: reason,
	})
}

// RevokeSessionByTokenHash revokes a session by token hash
func (r *RepositoryPG) RevokeSessionByTokenHash(ctx context.Context, tokenHash string, reason *string) error {
	return r.queries.RevokeSessionByTokenHash(ctx, db.RevokeSessionByTokenHashParams{
		TokenHash: tokenHash,
		Reason:    reason,
	})
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *RepositoryPG) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID, reason *string) error {
	return r.queries.RevokeAllUserSessions(ctx, db.RevokeAllUserSessionsParams{
		UserID: userID,
		Reason: reason,
	})
}

// RevokeAllUserSessionsExcept revokes all sessions except one
func (r *RepositoryPG) RevokeAllUserSessionsExcept(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID, reason *string) error {
	return r.queries.RevokeAllUserSessionsExcept(ctx, db.RevokeAllUserSessionsExceptParams{
		UserID: userID,
		ID:     exceptID,
		Reason: reason,
	})
}

// DeleteExpiredSessions deletes old expired sessions and returns the count
func (r *RepositoryPG) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	return r.queries.DeleteExpiredSessions(ctx)
}

// DeleteRevokedSessions deletes old revoked sessions and returns the count
func (r *RepositoryPG) DeleteRevokedSessions(ctx context.Context) (int64, error) {
	return r.queries.DeleteRevokedSessions(ctx)
}

// GetInactiveSessions gets sessions inactive since a given time
func (r *RepositoryPG) GetInactiveSessions(ctx context.Context, inactiveSince time.Time) ([]db.SharedSession, error) {
	return r.queries.GetInactiveSessions(ctx, inactiveSince)
}

// RevokeInactiveSessions revokes sessions inactive since a given time
func (r *RepositoryPG) RevokeInactiveSessions(ctx context.Context, inactiveSince time.Time) error {
	return r.queries.RevokeInactiveSessions(ctx, inactiveSince)
}
