package auth

import (
	"go.uber.org/fx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/email"
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
		func(pool *pgxpool.Pool, repo Repository, tm TokenManager, activityLogger activity.Logger, emailService *email.Service, cfg *config.Config) *Service {
			return NewService(pool, repo, tm, activityLogger, emailService, cfg.Auth.JWTExpiry, cfg.Auth.RefreshExpiry)
		},
	),
)
