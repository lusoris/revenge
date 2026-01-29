package movie

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

// ModuleParams contains dependencies for the adult movie module.
type ModuleParams struct {
	fx.In

	Pool   *pgxpool.Pool
	Logger *slog.Logger
}

// ModuleResult contains outputs from the adult movie module.
type ModuleResult struct {
	fx.Out

	Repository Repository
	Service    *Service
}

// ProvideModule provides adult movie dependencies.
func ProvideModule(p ModuleParams) ModuleResult {
	repo := NewRepository(p.Pool)
	service := NewService(repo, p.Logger)

	return ModuleResult{
		Repository: repo,
		Service:    service,
	}
}

// Module provides adult movie services.
var Module = fx.Module("adult-movie",
	fx.Provide(ProvideModule),
)
