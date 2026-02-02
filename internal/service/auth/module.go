package auth

import (
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Module provides the auth service and its dependencies
var Module = fx.Module("auth",
	fx.Provide(
		// TokenManager
		func(cfg *config.Config) TokenManager {
			return NewTokenManager(
				cfg.Auth.JWTSecret,
				cfg.Auth.JWTExpiry,
			)
		},
		// Repository
		func(queries *db.Queries) Repository {
			return NewRepositoryPG(queries)
		},
		// Service
		func(repo Repository, tm TokenManager, cfg *config.Config) *Service {
			return NewService(repo, tm, cfg.Auth.JWTExpiry, cfg.Auth.RefreshExpiry)
		},
	),
)
