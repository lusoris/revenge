// Package user provides user management services for Jellyfin Go.
package user

import (
	"go.uber.org/fx"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// Module provides all user-related services for dependency injection.
var Module = fx.Module("user",
	fx.Provide(
		NewService,
	),
)

// ServiceParams holds dependencies for the user Service.
type ServiceParams struct {
	fx.In
	Users     domain.UserRepository
	Sessions  domain.SessionRepository
	Passwords domain.PasswordService
}

// NewService creates a new user Service with fx integration.
func NewService(p ServiceParams) *Service {
	return newService(p.Users, p.Sessions, p.Passwords)
}
