package auth

import (
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
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
		func(repo Repository, tm TokenManager, activityLogger activity.Logger, cfg *config.Config) *Service {
			return NewService(repo, tm, activityLogger, cfg.Auth.JWTExpiry, cfg.Auth.RefreshExpiry)
		},
	),
)
