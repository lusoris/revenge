package user

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
)

// Module provides the user service and its dependencies
var Module = fx.Module("user",
	fx.Provide(
		// Repository
		func(queries *db.Queries) Repository {
			return NewPostgresRepository(queries)
		},
		// Service
		func(pool *pgxpool.Pool, repo Repository, activityLogger activity.Logger, store storage.Storage, cfg *config.Config) *Service {
			return NewService(pool, repo, activityLogger, store, cfg.Avatar)
		},
		// CachedService wraps the user service with caching.
		// When cache is nil (disabled), CachedService passes through to the underlying Service.
		func(svc *Service, c *cache.Cache, logger *slog.Logger) *CachedService {
			return NewCachedService(svc, c, logger)
		},
	),
)
