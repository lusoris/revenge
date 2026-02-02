package auth

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// RepositoryPG implements Repository using PostgreSQL via sqlc
type RepositoryPG struct {
	queries *db.Queries
}

// NewRepositoryPG creates a new PostgreSQL repository
func NewRepositoryPG(queries *db.Queries) *RepositoryPG {
	return &RepositoryPG{queries: queries}
}

// Auth Tokens

func (r *RepositoryPG) CreateAuthToken(ctx context.Context, params CreateAuthTokenParams) (AuthToken, error) {
	ipAddr := netip.Addr{}
	if params.IPAddress != nil {
		ipAddr = *params.IPAddress
	}

	row, err := r.queries.CreateAuthToken(ctx, db.CreateAuthTokenParams{
		UserID:            params.UserID,
		TokenHash:         params.TokenHash,
		TokenType:         params.TokenType,
		DeviceName:        params.DeviceName,
		DeviceFingerprint: params.DeviceFingerprint,
		IpAddress:         ipAddr,
		UserAgent:         params.UserAgent,
		ExpiresAt:         params.ExpiresAt,
	})
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed to create auth token: %w", err)
	}

	return authTokenFromDB(row), nil
}

func (r *RepositoryPG) GetAuthTokenByHash(ctx context.Context, tokenHash string) (AuthToken, error) {
	row, err := r.queries.GetAuthTokenByHash(ctx, tokenHash)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed to get auth token: %w", err)
	}
	return authTokenFromDB(row), nil
}

func (r *RepositoryPG) GetAuthTokensByUserID(ctx context.Context, userID uuid.UUID) ([]AuthToken, error) {
	rows, err := r.queries.GetAuthTokensByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth tokens by user: %w", err)
	}

	tokens := make([]AuthToken, len(rows))
	for i, row := range rows {
		tokens[i] = authTokenFromDB(row)
	}
	return tokens, nil
}

func (r *RepositoryPG) GetAuthTokensByDeviceFingerprint(ctx context.Context, userID uuid.UUID, deviceFingerprint string) ([]AuthToken, error) {
	rows, err := r.queries.GetAuthTokensByDeviceFingerprint(ctx, db.GetAuthTokensByDeviceFingerprintParams{
		UserID:            userID,
		DeviceFingerprint: &deviceFingerprint,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get auth tokens by device: %w", err)
	}

	tokens := make([]AuthToken, len(rows))
	for i, row := range rows {
		tokens[i] = authTokenFromDB(row)
	}
	return tokens, nil
}

func (r *RepositoryPG) UpdateAuthTokenLastUsed(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.UpdateAuthTokenLastUsed(ctx, id); err != nil {
		return fmt.Errorf("failed to update auth token last used: %w", err)
	}
	return nil
}

func (r *RepositoryPG) RevokeAuthToken(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.RevokeAuthToken(ctx, id); err != nil {
		return fmt.Errorf("failed to revoke auth token: %w", err)
	}
	return nil
}

func (r *RepositoryPG) RevokeAuthTokenByHash(ctx context.Context, tokenHash string) error {
	if err := r.queries.RevokeAuthTokenByHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to revoke auth token by hash: %w", err)
	}
	return nil
}

func (r *RepositoryPG) RevokeAllUserAuthTokens(ctx context.Context, userID uuid.UUID) error {
	if err := r.queries.RevokeAllUserAuthTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke all user auth tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) RevokeAllUserAuthTokensExcept(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID) error {
	if err := r.queries.RevokeAllUserAuthTokensExcept(ctx, db.RevokeAllUserAuthTokensExceptParams{
		UserID: userID,
		ID:     exceptID,
	}); err != nil {
		return fmt.Errorf("failed to revoke user auth tokens except one: %w", err)
	}
	return nil
}

func (r *RepositoryPG) DeleteExpiredAuthTokens(ctx context.Context) error {
	if err := r.queries.DeleteExpiredAuthTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete expired auth tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) DeleteRevokedAuthTokens(ctx context.Context) error {
	if err := r.queries.DeleteRevokedAuthTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete revoked auth tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) CountActiveAuthTokensByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.queries.CountActiveAuthTokensByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count active auth tokens: %w", err)
	}
	return count, nil
}

// Password Reset Tokens

func (r *RepositoryPG) CreatePasswordResetToken(ctx context.Context, params CreatePasswordResetTokenParams) (PasswordResetToken, error) {
	ipAddr := netip.Addr{}
	if params.IPAddress != nil {
		ipAddr = *params.IPAddress
	}

	row, err := r.queries.CreatePasswordResetToken(ctx, db.CreatePasswordResetTokenParams{
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		IpAddress: ipAddr,
		UserAgent: params.UserAgent,
		ExpiresAt: params.ExpiresAt,
	})
	if err != nil {
		return PasswordResetToken{}, fmt.Errorf("failed to create password reset token: %w", err)
	}

	return passwordResetTokenFromDB(row), nil
}

func (r *RepositoryPG) GetPasswordResetToken(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	row, err := r.queries.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return PasswordResetToken{}, fmt.Errorf("failed to get password reset token: %w", err)
	}
	return passwordResetTokenFromDB(row), nil
}

func (r *RepositoryPG) MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.MarkPasswordResetTokenUsed(ctx, id); err != nil {
		return fmt.Errorf("failed to mark password reset token used: %w", err)
	}
	return nil
}

func (r *RepositoryPG) InvalidateUserPasswordResetTokens(ctx context.Context, userID uuid.UUID) error {
	if err := r.queries.InvalidateUserPasswordResetTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to invalidate user password reset tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) DeleteExpiredPasswordResetTokens(ctx context.Context) error {
	if err := r.queries.DeleteExpiredPasswordResetTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete expired password reset tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) DeleteUsedPasswordResetTokens(ctx context.Context) error {
	if err := r.queries.DeleteUsedPasswordResetTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete used password reset tokens: %w", err)
	}
	return nil
}

// Email Verification Tokens

func (r *RepositoryPG) CreateEmailVerificationToken(ctx context.Context, params CreateEmailVerificationTokenParams) (EmailVerificationToken, error) {
	ipAddr := netip.Addr{}
	if params.IPAddress != nil {
		ipAddr = *params.IPAddress
	}

	row, err := r.queries.CreateEmailVerificationToken(ctx, db.CreateEmailVerificationTokenParams{
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		Email:     params.Email,
		IpAddress: ipAddr,
		UserAgent: params.UserAgent,
		ExpiresAt: params.ExpiresAt,
	})
	if err != nil {
		return EmailVerificationToken{}, fmt.Errorf("failed to create email verification token: %w", err)
	}

	return emailVerificationTokenFromDB(row), nil
}

func (r *RepositoryPG) GetEmailVerificationToken(ctx context.Context, tokenHash string) (EmailVerificationToken, error) {
	row, err := r.queries.GetEmailVerificationToken(ctx, tokenHash)
	if err != nil {
		return EmailVerificationToken{}, fmt.Errorf("failed to get email verification token: %w", err)
	}
	return emailVerificationTokenFromDB(row), nil
}

func (r *RepositoryPG) MarkEmailVerificationTokenUsed(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.MarkEmailVerificationTokenUsed(ctx, id); err != nil {
		return fmt.Errorf("failed to mark email verification token used: %w", err)
	}
	return nil
}

func (r *RepositoryPG) InvalidateUserEmailVerificationTokens(ctx context.Context, userID uuid.UUID) error {
	if err := r.queries.InvalidateUserEmailVerificationTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to invalidate user email verification tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) InvalidateEmailVerificationTokensByEmail(ctx context.Context, email string) error {
	if err := r.queries.InvalidateEmailVerificationTokensByEmail(ctx, email); err != nil {
		return fmt.Errorf("failed to invalidate email verification tokens by email: %w", err)
	}
	return nil
}

func (r *RepositoryPG) DeleteExpiredEmailVerificationTokens(ctx context.Context) error {
	if err := r.queries.DeleteExpiredEmailVerificationTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete expired email verification tokens: %w", err)
	}
	return nil
}

func (r *RepositoryPG) DeleteVerifiedEmailTokens(ctx context.Context) error {
	if err := r.queries.DeleteVerifiedEmailTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete verified email tokens: %w", err)
	}
	return nil
}

// Helper functions to convert from sqlc generated types

func authTokenFromDB(row db.SharedAuthToken) AuthToken {
	var ipAddr *netip.Addr
	if row.IpAddress.IsValid() {
		ipAddr = &row.IpAddress
	}

	var revokedAt *time.Time
	if row.RevokedAt.Valid {
		t := row.RevokedAt.Time
		revokedAt = &t
	}

	var lastUsedAt *time.Time
	if row.LastUsedAt.Valid {
		t := row.LastUsedAt.Time
		lastUsedAt = &t
	}

	return AuthToken{
		ID:                row.ID,
		UserID:            row.UserID,
		TokenHash:         row.TokenHash,
		TokenType:         row.TokenType,
		DeviceName:        row.DeviceName,
		DeviceFingerprint: row.DeviceFingerprint,
		IPAddress:         ipAddr,
		UserAgent:         row.UserAgent,
		ExpiresAt:         row.ExpiresAt,
		RevokedAt:         revokedAt,
		LastUsedAt:        lastUsedAt,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}

func passwordResetTokenFromDB(row db.SharedPasswordResetToken) PasswordResetToken {
	var ipAddr *netip.Addr
	if row.IpAddress.IsValid() {
		ipAddr = &row.IpAddress
	}

	var usedAt *time.Time
	if row.UsedAt.Valid {
		t := row.UsedAt.Time
		usedAt = &t
	}

	return PasswordResetToken{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		IPAddress: ipAddr,
		UserAgent: row.UserAgent,
		ExpiresAt: row.ExpiresAt,
		UsedAt:    usedAt,
		CreatedAt: row.CreatedAt,
	}
}

func emailVerificationTokenFromDB(row db.SharedEmailVerificationToken) EmailVerificationToken {
	var ipAddr *netip.Addr
	if row.IpAddress.IsValid() {
		ipAddr = &row.IpAddress
	}

	var verifiedAt *time.Time
	if row.VerifiedAt.Valid {
		t := row.VerifiedAt.Time
		verifiedAt = &t
	}

	return EmailVerificationToken{
		ID:         row.ID,
		UserID:     row.UserID,
		TokenHash:  row.TokenHash,
		Email:      row.Email,
		IPAddress:  ipAddr,
		UserAgent:  row.UserAgent,
		ExpiresAt:  row.ExpiresAt,
		VerifiedAt: verifiedAt,
		CreatedAt:  row.CreatedAt,
	}
}
