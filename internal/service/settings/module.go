package settings

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

// Module provides the settings service and its dependencies.
var Module = fx.Module("settings",
	fx.Provide(
		NewPostgresRepository,
		NewService,
	),
)

// RepositoryParams defines the dependencies for creating a repository.
type RepositoryParams struct {
	fx.In
	Pool *pgxpool.Pool
}

// ServiceParams defines the dependencies for creating a service.
type ServiceParams struct {
	fx.In
	Repo Repository
}
