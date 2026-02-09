package jobs

import (
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"log/slog"

	"github.com/lusoris/revenge/internal/content/tvshow"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/service/search"
)

// Module provides the TV show jobs workers.
var Module = fx.Module("tvshowjobs",
	fx.Provide(
		provideLibraryScanWorker,
		provideMetadataRefreshWorker,
		provideFileMatchWorker,
		provideSearchIndexWorker,
		provideSeriesRefreshWorker,
	),
	fx.Invoke(RegisterWorkers),
)

// WorkerProviderParams holds common dependencies for worker providers.
type WorkerProviderParams struct {
	fx.In

	Service              tvshow.Service
	MetadataProvider     tvshow.MetadataProvider      `optional:"true"`
	SearchService        *search.TVShowSearchService  `optional:"true"`
	EpisodeSearchService *search.EpisodeSearchService `optional:"true"`
	JobClient            *infrajobs.Client
	Logger               *slog.Logger
}

// provideLibraryScanWorker creates a library scan worker with optional metadata provider.
func provideLibraryScanWorker(p WorkerProviderParams) *LibraryScanWorker {
	return NewLibraryScanWorker(p.Service, p.MetadataProvider, p.JobClient, p.Logger)
}

// provideMetadataRefreshWorker creates a metadata refresh worker.
func provideMetadataRefreshWorker(p WorkerProviderParams) *MetadataRefreshWorker {
	return NewMetadataRefreshWorker(p.Service, p.JobClient, p.Logger)
}

// provideFileMatchWorker creates a file match worker with optional metadata provider.
func provideFileMatchWorker(p WorkerProviderParams) *FileMatchWorker {
	return NewFileMatchWorker(p.Service, p.MetadataProvider, p.Logger)
}

// provideSearchIndexWorker creates a search index worker.
func provideSearchIndexWorker(p WorkerProviderParams) *SearchIndexWorker {
	return NewSearchIndexWorker(p.Service, p.SearchService, p.EpisodeSearchService, p.Logger)
}

// provideSeriesRefreshWorker creates a series refresh worker.
func provideSeriesRefreshWorker(p WorkerProviderParams) *SeriesRefreshWorker {
	return NewSeriesRefreshWorker(p.Service, p.JobClient, p.Logger)
}

// RegisterWorkers registers all TV show job workers with the River workers registry.
func RegisterWorkers(
	workers *river.Workers,
	libraryScanWorker *LibraryScanWorker,
	metadataRefreshWorker *MetadataRefreshWorker,
	fileMatchWorker *FileMatchWorker,
	searchIndexWorker *SearchIndexWorker,
	seriesRefreshWorker *SeriesRefreshWorker,
) error {
	river.AddWorker(workers, libraryScanWorker)
	river.AddWorker(workers, metadataRefreshWorker)
	river.AddWorker(workers, fileMatchWorker)
	river.AddWorker(workers, searchIndexWorker)
	river.AddWorker(workers, seriesRefreshWorker)
	return nil
}

// RegisterWorkersParams defines the parameters for RegisterWorkers.
type RegisterWorkersParams struct {
	fx.In

	Workers               *river.Workers
	LibraryScanWorker     *LibraryScanWorker
	MetadataRefreshWorker *MetadataRefreshWorker
	FileMatchWorker       *FileMatchWorker
	SearchIndexWorker     *SearchIndexWorker      `optional:"true"`
	SeriesRefreshWorker   *SeriesRefreshWorker
	TVShowService         tvshow.Service
	MetadataProvider      tvshow.MetadataProvider `optional:"true"`
	Logger                *slog.Logger
}
