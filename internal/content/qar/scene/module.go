package scene

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

// ModuleParams contains dependencies for the adult scene module.
type ModuleParams struct {
	fx.In

	Pool   *pgxpool.Pool
	Logger *slog.Logger
}

// ModuleResult contains outputs from the adult scene module.
type ModuleResult struct {
	fx.Out

	Repository Repository
	Service    *Service
}

// ProvideModule provides adult scene dependencies.
func ProvideModule(p ModuleParams) ModuleResult {
	repo := NewRepository(p.Pool)
	service := NewService(repo, p.Logger)

	return ModuleResult{
		Repository: repo,
		Service:    service,
	}
}

// Module provides adult scene services.
var Module = fx.Module("adult-scene",
	fx.Provide(ProvideModule),
)
