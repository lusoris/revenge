package library

import (
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/fx"
)

// Module provides library service dependencies.
var Module = fx.Module("library",
	fx.Provide(
		newRepository,
		NewService,
	),
)

// newRepository creates a new library repository.
func newRepository(queries *db.Queries) Repository {
	return NewRepositoryPg(queries)
}
