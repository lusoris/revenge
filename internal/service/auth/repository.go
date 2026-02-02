package auth

import (
	"context"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Repository defines the data access interface for auth operations
type Repository interface {
	// User Operations (required for auth flows)
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.SharedUser, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*db.SharedUser, error)
	GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error)
	GetUserByEmail(ctx context.Context, email string) (*db.SharedUser, error)
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateUserEmailVerified(ctx context.Context, userID uuid.UUID, verified bool) error
	UpdateUserLastLogin(ctx context.Context, userID uuid.UUID) error

	// Auth Tokens (JWT refresh tokens)
	CreateAuthToken(ctx context.Context, params CreateAuthTokenParams) (AuthToken, error)
	GetAuthTokenByHash(ctx context.Context, tokenHash string) (AuthToken, error)
	GetAuthTokensByUserID(ctx context.Context, userID uuid.UUID) ([]AuthToken, error)
	GetAuthTokensByDeviceFingerprint(ctx context.Context, userID uuid.UUID, deviceFingerprint string) ([]AuthToken, error)
	UpdateAuthTokenLastUsed(ctx context.Context, id uuid.UUID) error
	RevokeAuthToken(ctx context.Context, id uuid.UUID) error
	RevokeAuthTokenByHash(ctx context.Context, tokenHash string) error
	RevokeAllUserAuthTokens(ctx context.Context, userID uuid.UUID) error
	RevokeAllUserAuthTokensExcept(ctx context.Context, userID uuid.UUID, exceptID uuid.UUID) error
	DeleteExpiredAuthTokens(ctx context.Context) error
	DeleteRevokedAuthTokens(ctx context.Context) error
	CountActiveAuthTokensByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	// Password Reset Tokens
	CreatePasswordResetToken(ctx context.Context, params CreatePasswordResetTokenParams) (PasswordResetToken, error)
	GetPasswordResetToken(ctx context.Context, tokenHash string) (PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID) error
	InvalidateUserPasswordResetTokens(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredPasswordResetTokens(ctx context.Context) error
	DeleteUsedPasswordResetTokens(ctx context.Context) error

	// Email Verification Tokens
	CreateEmailVerificationToken(ctx context.Context, params CreateEmailVerificationTokenParams) (EmailVerificationToken, error)
	GetEmailVerificationToken(ctx context.Context, tokenHash string) (EmailVerificationToken, error)
	MarkEmailVerificationTokenUsed(ctx context.Context, id uuid.UUID) error
	InvalidateUserEmailVerificationTokens(ctx context.Context, userID uuid.UUID) error
	InvalidateEmailVerificationTokensByEmail(ctx context.Context, email string) error
	DeleteExpiredEmailVerificationTokens(ctx context.Context) error
	DeleteVerifiedEmailTokens(ctx context.Context) error
}

// CreateAuthTokenParams parameters for creating an auth token
type CreateAuthTokenParams struct {
	UserID            uuid.UUID
	TokenHash         string
	TokenType         string
	DeviceName        *string
	DeviceFingerprint *string
	IPAddress         *netip.Addr
	UserAgent         *string
	ExpiresAt         time.Time
}

// AuthToken represents a stored refresh token
type AuthToken struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	TokenHash         string
	TokenType         string
	DeviceName        *string
	DeviceFingerprint *string
	IPAddress         *netip.Addr
	UserAgent         *string
	ExpiresAt         time.Time
	RevokedAt         *time.Time
	LastUsedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CreatePasswordResetTokenParams parameters for creating a password reset token
type CreatePasswordResetTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	IPAddress *netip.Addr
	UserAgent *string
	ExpiresAt time.Time
}

// PasswordResetToken represents a stored password reset token
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	IPAddress *netip.Addr
	UserAgent *string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// CreateEmailVerificationTokenParams parameters for creating an email verification token
type CreateEmailVerificationTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	Email     string
	IPAddress *netip.Addr
	UserAgent *string
	ExpiresAt time.Time
}

// EmailVerificationToken represents a stored email verification token
type EmailVerificationToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	Email      string
	IPAddress  *netip.Addr
	UserAgent  *string
	ExpiresAt  time.Time
	VerifiedAt *time.Time
	CreatedAt  time.Time
}
