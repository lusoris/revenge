// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/infra/database/db"
)

// SessionRepository implements domain.SessionRepository using PostgreSQL.
type SessionRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewSessionRepository creates a new PostgreSQL session repository.
func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// GetByID retrieves a session by its ID.
func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	session, err := r.queries.GetSessionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}
	return mapDBSessionToDomain(&session), nil
}

// GetByTokenHash retrieves a session by the hashed access token.
func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	session, err := r.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by token hash: %w", err)
	}

	// Check if session has expired
	if session.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrSessionExpired
	}

	return mapDBSessionToDomain(&session), nil
}

// GetByRefreshTokenHash retrieves a session by the hashed refresh token.
func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Session, error) {
	session, err := r.queries.GetSessionByRefreshTokenHash(ctx, pgtype.Text{String: refreshTokenHash, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by refresh token hash: %w", err)
	}
	return mapDBSessionToDomain(&session), nil
}

// GetWithUser retrieves a session with joined user information.
func (r *SessionRepository) GetWithUser(ctx context.Context, tokenHash string) (*domain.SessionWithUser, error) {
	row, err := r.queries.GetSessionWithUser(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session with user: %w", err)
	}

	result := &domain.SessionWithUser{
		Session: domain.Session{
			ID:        row.ID,
			UserID:    row.UserID,
			TokenHash: row.TokenHash,
			IPAddress: row.IpAddress,
			ExpiresAt: row.ExpiresAt,
			CreatedAt: row.CreatedAt,
		},
		Username:   row.Username,
		IsAdmin:    row.IsAdmin,
		IsDisabled: row.IsDisabled,
	}

	if row.RefreshTokenHash.Valid {
		result.RefreshTokenHash = &row.RefreshTokenHash.String
	}
	if row.DeviceID.Valid {
		result.DeviceID = &row.DeviceID.String
	}
	if row.DeviceName.Valid {
		result.DeviceName = &row.DeviceName.String
	}
	if row.ClientName.Valid {
		result.ClientName = &row.ClientName.String
	}
	if row.ClientVersion.Valid {
		result.ClientVersion = &row.ClientVersion.String
	}
	if row.RefreshExpiresAt.Valid {
		result.RefreshExpiresAt = &row.RefreshExpiresAt.Time
	}
	if row.Email.Valid {
		result.Email = &row.Email.String
	}
	if row.DisplayName.Valid {
		result.DisplayName = &row.DisplayName.String
	}

	return result, nil
}

// ListByUser retrieves all sessions for a specific user.
func (r *SessionRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	sessions, err := r.queries.ListUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user sessions: %w", err)
	}

	result := make([]*domain.Session, len(sessions))
	for i, s := range sessions {
		result[i] = mapDBSessionToDomain(&s)
	}
	return result, nil
}

// Create creates a new session.
func (r *SessionRepository) Create(ctx context.Context, params domain.CreateSessionParams) (*domain.Session, error) {
	dbParams := db.CreateSessionParams{
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		IpAddress: params.IPAddress,
		ExpiresAt: params.ExpiresAt,
	}

	if params.RefreshTokenHash != nil {
		dbParams.RefreshTokenHash = pgtype.Text{String: *params.RefreshTokenHash, Valid: true}
	}
	if params.DeviceID != nil {
		dbParams.DeviceID = pgtype.Text{String: *params.DeviceID, Valid: true}
	}
	if params.DeviceName != nil {
		dbParams.DeviceName = pgtype.Text{String: *params.DeviceName, Valid: true}
	}
	if params.ClientName != nil {
		dbParams.ClientName = pgtype.Text{String: *params.ClientName, Valid: true}
	}
	if params.ClientVersion != nil {
		dbParams.ClientVersion = pgtype.Text{String: *params.ClientVersion, Valid: true}
	}
	if params.RefreshExpiresAt != nil {
		dbParams.RefreshExpiresAt = pgtype.Timestamptz{Time: *params.RefreshExpiresAt, Valid: true}
	}

	session, err := r.queries.CreateSession(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return mapDBSessionToDomain(&session), nil
}

// UpdateRefreshToken updates the session's refresh token.
func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID, refreshTokenHash string, expiresAt time.Time) error {
	err := r.queries.UpdateSessionRefreshToken(ctx, db.UpdateSessionRefreshTokenParams{
		ID:               id,
		RefreshTokenHash: pgtype.Text{String: refreshTokenHash, Valid: true},
		RefreshExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}
	return nil
}

// Delete removes a session by its ID.
func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteByTokenHash removes a session by its token hash.
func (r *SessionRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	err := r.queries.DeleteSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session by token hash: %w", err)
	}
	return nil
}

// DeleteByUser removes all sessions for a user.
func (r *SessionRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.DeleteUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

// DeleteExpired removes all expired sessions and returns the count.
func (r *SessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	count, err := r.queries.DeleteExpiredSessions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return count, nil
}

// CountByUser returns the number of sessions for a user.
func (r *SessionRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.queries.CountUserSessions(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count user sessions: %w", err)
	}
	return count, nil
}

// Exists checks if a valid session exists for a token hash.
func (r *SessionRepository) Exists(ctx context.Context, tokenHash string) (bool, error) {
	exists, err := r.queries.SessionExists(ctx, tokenHash)
	if err != nil {
		return false, fmt.Errorf("failed to check session exists: %w", err)
	}
	return exists, nil
}

// mapDBSessionToDomain converts a database session to a domain session.
func mapDBSessionToDomain(s *db.Session) *domain.Session {
	session := &domain.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		TokenHash: s.TokenHash,
		IPAddress: s.IpAddress,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}

	if s.RefreshTokenHash.Valid {
		session.RefreshTokenHash = &s.RefreshTokenHash.String
	}
	if s.DeviceID.Valid {
		session.DeviceID = &s.DeviceID.String
	}
	if s.DeviceName.Valid {
		session.DeviceName = &s.DeviceName.String
	}
	if s.ClientName.Valid {
		session.ClientName = &s.ClientName.String
	}
	if s.ClientVersion.Valid {
		session.ClientVersion = &s.ClientVersion.String
	}
	if s.RefreshExpiresAt.Valid {
		session.RefreshExpiresAt = &s.RefreshExpiresAt.Time
	}

	return session
}

// Silence unused import warning
var _ = netip.Addr{}

// Ensure SessionRepository implements domain.SessionRepository.
var _ domain.SessionRepository = (*SessionRepository)(nil)
