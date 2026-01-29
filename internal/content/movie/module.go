package movie

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
)

// ModuleParams contains dependencies for the movie module.
type ModuleParams struct {
	fx.In

	Pool          *pgxpool.Pool
	Logger        *slog.Logger
	ServiceConfig ServiceConfig `optional:"true"`
}

// ModuleResult contains outputs from the movie module.
type ModuleResult struct {
	fx.Out

	Repository Repository
	Service    *Service
}

// ProvideModule provides all movie module dependencies.
// Note: Movie data primarily comes from Servarr (Radarr) API.
// Enrichment via TMDb is handled by background River jobs.
func ProvideModule(p ModuleParams) (ModuleResult, error) {
	repo := NewRepository(p.Pool)

	// Service doesn't need a metadata provider directly.
	// Metadata enrichment is done via River jobs that call the metadata service.
	service, err := NewService(repo, nil, p.Logger, p.ServiceConfig)
	if err != nil {
		return ModuleResult{}, err
	}

	return ModuleResult{
		Repository: repo,
		Service:    service,
	}, nil
}

// RegisterLifecycleHooks registers shutdown hooks for graceful cleanup.
func RegisterLifecycleHooks(lc fx.Lifecycle, service *Service, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down movie service")
			service.Close()
			return nil
		},
	})
}

// RiverParams contains optional dependencies for River worker registration.
type RiverParams struct {
	fx.In

	Workers  *river.Workers
	Service  *Service
	Client   *river.Client[any]
	Logger   *slog.Logger
	Provider *tmdb.Provider `optional:"true"` // Optional: for metadata enrichment
	Scanner  Scanner        `optional:"true"` // Optional: for library scanning
}

// RegisterRiverWorkers registers movie workers with River if available.
func RegisterRiverWorkers(p RiverParams) error {
	return RegisterWorkers(p.Workers, p.Service, p.Scanner, p.Provider, p.Client, p.Logger)
}

// Module is the fx module for movies.
var Module = fx.Module("movie",
	fx.Provide(ProvideModule),
	fx.Invoke(RegisterLifecycleHooks),
)

// ModuleWithRiver provides the movie module with River job integration.
var ModuleWithRiver = fx.Module("movie-with-river",
	Module,
	fx.Invoke(RegisterRiverWorkers),
)
