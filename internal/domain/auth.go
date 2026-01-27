// Package domain contains core business entities and repository interfaces.
package domain

import (
	"context"
	"net/netip"
	"time"

	"github.com/google/uuid"
)

// AuthResult represents the result of a successful authentication.
type AuthResult struct {
	User         *User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	SessionID    uuid.UUID
}

// LoginParams contains parameters for user login.
type LoginParams struct {
	Username      string
	Password      string
	DeviceID      *string
	DeviceName    *string
	ClientName    *string
	ClientVersion *string
	IPAddress     *netip.Addr
}

// RefreshParams contains parameters for token refresh.
type RefreshParams struct {
	RefreshToken string
	IPAddress    *netip.Addr
}

// TokenClaims represents the claims stored in a JWT token.
type TokenClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Username  string
	IsAdmin   bool
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// AuthService defines the interface for authentication operations.
type AuthService interface {
	// Login authenticates a user with username and password.
	// Returns tokens and session info on success.
	Login(ctx context.Context, params LoginParams) (*AuthResult, error)

	// Logout invalidates a session by its access token.
	Logout(ctx context.Context, accessToken string) error

	// LogoutAll invalidates all sessions for a user.
	LogoutAll(ctx context.Context, userID uuid.UUID) error

	// RefreshToken exchanges a refresh token for new access/refresh tokens.
	RefreshToken(ctx context.Context, params RefreshParams) (*AuthResult, error)

	// ValidateToken validates an access token and returns the claims.
	ValidateToken(ctx context.Context, accessToken string) (*TokenClaims, error)

	// GetSession retrieves session with user info by access token.
	GetSession(ctx context.Context, accessToken string) (*SessionWithUser, error)

	// ChangePassword changes a user's password (requires current password).
	ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error

	// ResetPassword sets a new password (admin operation, no current password required).
	ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

// PasswordService defines the interface for password operations.
type PasswordService interface {
	// Hash creates a bcrypt hash of the password.
	Hash(password string) (string, error)

	// Verify checks if a password matches a hash.
	Verify(password, hash string) error
}

// TokenService defines the interface for JWT token operations.
type TokenService interface {
	// GenerateAccessToken creates a new access token.
	GenerateAccessToken(claims TokenClaims) (string, error)

	// GenerateRefreshToken creates a new refresh token.
	GenerateRefreshToken() (string, error)

	// ValidateAccessToken validates an access token and extracts claims.
	ValidateAccessToken(token string) (*TokenClaims, error)

	// HashToken creates a SHA-256 hash of a token for storage.
	HashToken(token string) string
}
