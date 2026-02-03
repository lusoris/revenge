package moviejobs

import (
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/content/movie"
)

// Module provides the movie jobs workers.
var Module = fx.Module("moviejobs",
	fx.Provide(
		NewMovieMetadataRefreshWorker,
		NewMovieLibraryScanWorker,
		NewMovieFileMatchWorker,
	),
)

// RegisterWorkers registers all movie job workers with the River workers registry.
func RegisterWorkers(
	workers *river.Workers,
	metadataRefreshWorker *MovieMetadataRefreshWorker,
	libraryScanWorker *MovieLibraryScanWorker,
	fileMatchWorker *MovieFileMatchWorker,
) error {
	river.AddWorker(workers, metadataRefreshWorker)
	river.AddWorker(workers, libraryScanWorker)
	river.AddWorker(workers, fileMatchWorker)
	return nil
}

// RegisterWorkersParams defines the parameters for RegisterWorkers.
type RegisterWorkersParams struct {
	fx.In

	Workers                  *river.Workers
	MetadataRefreshWorker    *MovieMetadataRefreshWorker
	LibraryScanWorker        *MovieLibraryScanWorker
	FileMatchWorker          *MovieFileMatchWorker
	MovieRepository          movie.Repository
	MetadataService          *movie.MetadataService
	LibraryService           *movie.LibraryService
	Logger                   *zap.Logger
}
