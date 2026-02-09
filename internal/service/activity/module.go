package activity

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/infra/database/db"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

// Module provides activity service dependencies.
var Module = fx.Module("activity",
	fx.Provide(
		newRepository,
		NewService,
		provideLogger,
		NewActivityCleanupWorker,
		NewActivityLogWorker,
	),
	fx.Invoke(
		registerActivityLogWorker,
	),
)

func newRepository(queries *db.Queries) Repository {
	return NewRepositoryPg(queries)
}

// provideLogger returns an AsyncLogger backed by River jobs when the job
// client is available, falling back to the synchronous ServiceLogger.
func provideLogger(svc *Service, client *infrajobs.Client, logger *slog.Logger) Logger {
	if client != nil {
		return NewAsyncLogger(client, logger)
	}
	return NewLogger(svc)
}

// registerActivityLogWorker registers the activity log worker with River.
func registerActivityLogWorker(workers *river.Workers, worker *ActivityLogWorker) {
	river.AddWorker(workers, worker)
}
