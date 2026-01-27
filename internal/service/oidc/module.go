// Package oidc provides OIDC/SSO authentication services for Revenge Go.
package oidc

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/domain"
	"github.com/lusoris/revenge/pkg/config"
)

// Module provides OIDC-related services for dependency injection.
var Module = fx.Module("oidc",
	fx.Provide(
		NewService,
		AsOIDCService,
	),
)

// AsOIDCService binds Service to the Service interface.
func AsOIDCService(s *Service) AuthService {
	return s
}

// AuthService defines the interface for OIDC operations.
// This is used for dependency injection and mocking.
type AuthService interface {
	// GetProviders returns all enabled OIDC providers.
	GetProviders(ctx context.Context) ([]*domain.OIDCProvider, error)

	// GetAuthorizationURL generates the authorization URL for a provider.
	GetAuthorizationURL(ctx context.Context, providerID uuid.UUID, redirectURI string) (string, error)

	// HandleCallback processes the OIDC callback and returns auth result.
	HandleCallback(ctx context.Context, params CallbackParams) (*domain.AuthResult, error)

	// LinkUser links an existing user to an OIDC provider.
	LinkUser(ctx context.Context, userID, providerID uuid.UUID, subject string, email *string) error

	// UnlinkUser removes an OIDC link from a user.
	UnlinkUser(ctx context.Context, userID, linkID uuid.UUID) error

	// GetUserLinks returns all OIDC links for a user.
	GetUserLinks(ctx context.Context, userID uuid.UUID) ([]*domain.OIDCUserLinkWithProvider, error)
}

// Params holds dependencies for the OIDC service (fx injection).
type Params struct {
	fx.In
	Config    *config.Config
	Providers domain.OIDCProviderRepository
	UserLinks domain.OIDCUserLinkRepository
	Users     domain.UserRepository
	Sessions  domain.SessionRepository
	Tokens    domain.TokenService
	Passwords domain.PasswordService
}

// NewService creates a new OIDC service with fx integration.
func NewService(p Params) *Service {
	accessDuration := 15 * time.Minute
	refreshDuration := 7 * 24 * time.Hour

	// Parse durations from config
	if d, err := time.ParseDuration(p.Config.Auth.AccessTokenDuration); err == nil {
		accessDuration = d
	}
	if d, err := time.ParseDuration(p.Config.Auth.RefreshTokenDuration); err == nil {
		refreshDuration = d
	}

	return newService(ServiceParams{
		Providers:       p.Providers,
		UserLinks:       p.UserLinks,
		Users:           p.Users,
		Sessions:        p.Sessions,
		Tokens:          p.Tokens,
		Passwords:       p.Passwords,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	})
}

// newService creates a new OIDC service (internal constructor).
func newService(params ServiceParams) *Service {
	if params.AccessDuration <= 0 {
		params.AccessDuration = 15 * time.Minute
	}
	if params.RefreshDuration <= 0 {
		params.RefreshDuration = 7 * 24 * time.Hour
	}

	return &Service{
		providers:       params.Providers,
		userLinks:       params.UserLinks,
		users:           params.Users,
		sessions:        params.Sessions,
		tokens:          params.Tokens,
		passwords:       params.Passwords,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		stateStore:      make(map[string]*authState),
		discoveryCache:  make(map[string]*cachedDiscovery),
		stateExpiry:     10 * time.Minute,
		accessDuration:  params.AccessDuration,
		refreshDuration: params.RefreshDuration,
	}
}
