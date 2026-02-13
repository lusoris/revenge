package moviejobs

import (
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"log/slog"

	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/service/search"
)

// Module provides the movie jobs workers.
var Module = fx.Module("moviejobs",
	fx.Provide(
		NewMovieMetadataRefreshWorker,
		NewMovieLibraryScanWorker,
		NewMovieFileMatchWorker,
		NewMovieSearchIndexWorker,
	),
	fx.Invoke(RegisterWorkers),
)

// RegisterWorkers registers all movie job workers with the River workers registry.
func RegisterWorkers(
	workers *river.Workers,
	metadataRefreshWorker *MovieMetadataRefreshWorker,
	libraryScanWorker *MovieLibraryScanWorker,
	fileMatchWorker *MovieFileMatchWorker,
	searchIndexWorker *MovieSearchIndexWorker,
) error {
	river.AddWorker(workers, metadataRefreshWorker)
	river.AddWorker(workers, libraryScanWorker)
	river.AddWorker(workers, fileMatchWorker)
	river.AddWorker(workers, searchIndexWorker)
	return nil
}

// RegisterWorkersParams defines the parameters for RegisterWorkers.
type RegisterWorkersParams struct {
	fx.In

	Workers               *river.Workers
	MetadataRefreshWorker *MovieMetadataRefreshWorker
	LibraryScanWorker     *MovieLibraryScanWorker
	FileMatchWorker       *MovieFileMatchWorker
	SearchIndexWorker     *MovieSearchIndexWorker `optional:"true"`
	MovieService          movie.Service
	LibraryService        *movie.LibraryService
	SearchService         *search.MovieSearchService `optional:"true"`
	Logger                *slog.Logger
}
