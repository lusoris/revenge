package auth

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
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
		func(pool *pgxpool.Pool, repo Repository, tm TokenManager, activityLogger activity.Logger, emailService *email.Service, logger *slog.Logger, cfg *config.Config) *Service {
			return NewService(
				pool,
				repo,
				tm,
				activityLogger,
				emailService,
				logger,
				cfg.Auth.JWTExpiry,
				cfg.Auth.RefreshExpiry,
				cfg.Auth.LockoutThreshold,
				cfg.Auth.LockoutWindow,
				cfg.Auth.LockoutEnabled,
			)
		},
		// Cleanup worker for expired/revoked auth tokens
		provideAuthCleanupWorker,
	),
)

// provideAuthCleanupWorker creates the auth token cleanup worker.
// The auth.Repository satisfies the jobs.AuthCleanupRepository interface.
func provideAuthCleanupWorker(repo Repository, logger *slog.Logger) *infrajobs.CleanupWorker {
	return infrajobs.NewCleanupWorker(repo, logger)
}
