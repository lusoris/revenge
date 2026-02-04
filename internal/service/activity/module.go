package activity

import (
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/fx"
)

// Module provides activity service dependencies.
var Module = fx.Module("activity",
	fx.Provide(
		newRepository,
		NewService,
		NewLogger,
		NewActivityCleanupWorker,
	),
)

func newRepository(queries *db.Queries) Repository {
	return NewRepositoryPg(queries)
}
