package library

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/service/activity"
	"go.uber.org/fx"
)

// Module provides library service dependencies.
var Module = fx.Module("library",
	fx.Provide(
		newRepository,
		newService,
		newCachedService,
		NewLibraryScanCleanupWorker,
		newPeriodicLibraryScanWorker,
	),
)

// newRepository creates a new library repository.
func newRepository(queries *db.Queries) Repository {
	return NewRepositoryPg(queries)
}

// newService creates a new library service with activity logger.
func newService(repo Repository, logger *slog.Logger, activityLogger activity.Logger) *Service {
	return NewService(repo, logger, activityLogger)
}

// newCachedService wraps the library service with caching.
// When cache is nil (disabled), CachedService passes through to the underlying Service.
func newCachedService(svc *Service, c *cache.Cache, logger *slog.Logger) *CachedService {
	return NewCachedService(svc, c, logger)
}

// newPeriodicLibraryScanWorker creates the periodic scan worker using the infra jobs client.
func newPeriodicLibraryScanWorker(repo Repository, client *infrajobs.Client, logger *slog.Logger) *PeriodicLibraryScanWorker {
	return NewPeriodicLibraryScanWorker(repo, client, logger)
}
