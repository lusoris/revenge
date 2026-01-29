package tvshow

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/content/shared"
	tvshowdb "github.com/lusoris/revenge/internal/content/tvshow/db"
	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
)

// ModuleParams contains dependencies for the tvshow module.
type ModuleParams struct {
	fx.In

	Pool          *pgxpool.Pool
	Logger        *slog.Logger
	ServiceConfig ServiceConfig `optional:"true"`
}

// ModuleResult contains outputs from the tvshow module.
type ModuleResult struct {
	fx.Out

	Repository      Repository
	Service         *Service
	LibraryService  *LibraryService
	LibraryProvider shared.LibraryProvider `group:"library_providers"`
}

// ProvideModule provides all tvshow module dependencies.
// Note: TV show data primarily comes from Servarr (Sonarr) API.
// Enrichment via TMDb is handled by background River jobs.
func ProvideModule(p ModuleParams) (ModuleResult, error) {
	repo := NewRepository(p.Pool)

	// Service doesn't need a metadata provider directly.
	// Metadata enrichment is done via River jobs that call the metadata service.
	service, err := NewService(repo, p.Logger, p.ServiceConfig)
	if err != nil {
		return ModuleResult{}, err
	}

	// Create library service
	queries := tvshowdb.New(p.Pool)
	libraryService := NewLibraryService(queries, p.Logger)

	return ModuleResult{
		Repository:      repo,
		Service:         service,
		LibraryService:  libraryService,
		LibraryProvider: libraryService,
	}, nil
}

// RegisterLifecycleHooks registers shutdown hooks for graceful cleanup.
func RegisterLifecycleHooks(lc fx.Lifecycle, service *Service, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down tvshow service")
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
	Client   *river.Client[pgx.Tx]
	Logger   *slog.Logger
	Provider *tmdb.Provider `optional:"true"` // Optional: for metadata enrichment
	Scanner  Scanner        `optional:"true"` // Optional: for library scanning
}

// RegisterRiverWorkers registers tvshow workers with River if available.
func RegisterRiverWorkers(p RiverParams) error {
	var provider MetadataProvider
	if p.Provider != nil {
		adapter := newTMDbAdapter(p.Provider)
		if adapter != nil && adapter.IsAvailable() {
			provider = adapter
		}
	}
	return RegisterWorkers(p.Workers, p.Service, p.Scanner, provider, p.Client, p.Logger)
}

// Module is the fx module for TV shows.
var Module = fx.Module("tvshow",
	fx.Provide(ProvideModule),
	fx.Invoke(RegisterLifecycleHooks),
)

// ModuleWithRiver provides the tvshow module with River job integration.
var ModuleWithRiver = fx.Module("tvshow-with-river",
	Module,
	fx.Invoke(RegisterRiverWorkers),
)
