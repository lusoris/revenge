// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/domain"
)

// Module provides all repository implementations for dependency injection.
var Module = fx.Options(
	fx.Provide(
		// User repository - provide concrete type and interface
		func(r *UserRepository) domain.UserRepository { return r },
		NewUserRepository,

		// Session repository
		func(r *SessionRepository) domain.SessionRepository { return r },
		NewSessionRepository,

		// OIDC provider repository
		func(r *OIDCProviderRepository) domain.OIDCProviderRepository { return r },
		NewOIDCProviderRepository,

		// OIDC user link repository
		func(r *OIDCUserLinkRepository) domain.OIDCUserLinkRepository { return r },
		NewOIDCUserLinkRepository,

		// Genre repository
		func(r *GenreRepository) domain.GenreRepository { return r },
		NewGenreRepository,
	),
)
