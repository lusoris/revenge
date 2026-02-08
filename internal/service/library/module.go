package library

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"go.uber.org/fx"
)

// Module provides library service dependencies.
var Module = fx.Module("library",
	fx.Provide(
		newRepository,
		newService,
		NewLibraryScanCleanupWorker,
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
