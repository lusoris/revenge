// Package auth provides authentication services for the Jellyfin Go server.
package auth

import (
	"time"

	"go.uber.org/fx"

	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/pkg/config"
)

// Module provides all auth-related services for dependency injection.
var Module = fx.Module("auth",
	fx.Provide(
		NewPasswordService,
		NewTokenService,
		NewService,
		// Bind interfaces to implementations
		AsPasswordService,
		AsTokenService,
		AsAuthService,
	),
)

// AsPasswordService binds PasswordService to domain.PasswordService interface.
func AsPasswordService(s *PasswordService) domain.PasswordService {
	return s
}

// AsTokenService binds TokenService to domain.TokenService interface.
func AsTokenService(s *TokenService) domain.TokenService {
	return s
}

// AsAuthService binds Service to domain.AuthService interface.
func AsAuthService(s *Service) domain.AuthService {
	return s
}

// PasswordServiceParams holds dependencies for PasswordService.
type PasswordServiceParams struct {
	fx.In
	Config *config.Config
}

// NewPasswordService creates a new PasswordService with fx integration.
func NewPasswordService(p PasswordServiceParams) *PasswordService {
	return newPasswordService(p.Config.Auth.BcryptCost)
}

// TokenServiceParams holds dependencies for TokenService.
type TokenServiceParams struct {
	fx.In
	Config *config.Config
}

// NewTokenService creates a new TokenService with fx integration.
func NewTokenService(p TokenServiceParams) *TokenService {
	accessDuration := parseDuration(p.Config.Auth.AccessTokenDuration, 15*time.Minute)
	return newTokenService(p.Config.Auth.JWTSecret, accessDuration)
}

// ServiceParams holds dependencies for the main auth Service.
type ServiceParams struct {
	fx.In
	Config    *config.Config
	Users     domain.UserRepository
	Sessions  domain.SessionRepository
	Passwords domain.PasswordService
	Tokens    domain.TokenService
}

// NewService creates a new auth Service with fx integration.
func NewService(p ServiceParams) *Service {
	accessDuration := parseDuration(p.Config.Auth.AccessTokenDuration, 15*time.Minute)
	refreshDuration := parseDuration(p.Config.Auth.RefreshTokenDuration, 7*24*time.Hour)

	return newService(
		p.Users,
		p.Sessions,
		p.Passwords,
		p.Tokens,
		p.Config.Auth.MaxSessionsPerUser,
		accessDuration,
		refreshDuration,
	)
}

// parseDuration parses a duration string, returning the default on error.
func parseDuration(s string, defaultVal time.Duration) time.Duration {
	if s == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultVal
	}
	return d
}
