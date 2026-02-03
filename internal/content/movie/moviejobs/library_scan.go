package moviejobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/content/movie"
)

const MovieLibraryScanJobKind = "movie_library_scan"

// MovieLibraryScanArgs are the arguments for the movie library scan job.
type MovieLibraryScanArgs struct {
	Paths []string `json:"paths"`
	Force bool     `json:"force"`
}

// Kind returns the job kind for the movie library scan job.
func (MovieLibraryScanArgs) Kind() string {
	return MovieLibraryScanJobKind
}

// MovieLibraryScanWorker is a worker that scans movie libraries.
type MovieLibraryScanWorker struct {
	river.WorkerDefaults[MovieLibraryScanArgs]
	libraryService *movie.LibraryService
	logger         *zap.Logger
}

// NewMovieLibraryScanWorker creates a new movie library scan worker.
func NewMovieLibraryScanWorker(
	libraryService *movie.LibraryService,
	logger *zap.Logger,
) *MovieLibraryScanWorker {
	return &MovieLibraryScanWorker{
		libraryService: libraryService,
		logger:         logger,
	}
}

// Kind returns the job kind.
func (w *MovieLibraryScanWorker) Kind() string {
	return MovieLibraryScanJobKind
}

// Work performs the movie library scan job.
func (w *MovieLibraryScanWorker) Work(ctx context.Context, job *river.Job[MovieLibraryScanArgs]) error {
	args := job.Args

	w.logger.Info("starting movie library scan",
		zap.Strings("paths", args.Paths),
		zap.Bool("force", args.Force),
	)

	// Call library service to scan the library
	summary, err := w.libraryService.ScanLibrary(ctx)
	if err != nil {
		w.logger.Error("library scan failed",
			zap.Error(err),
		)
		return fmt.Errorf("library scan failed: %w", err)
	}

	// Log summary
	w.logger.Info("library scan completed",
		zap.Int("total_files", summary.TotalFiles),
		zap.Int("matched_files", summary.MatchedFiles),
		zap.Int("unmatched_files", summary.UnmatchedFiles),
		zap.Int("new_movies", summary.NewMovies),
		zap.Int("existing_movies", summary.ExistingMovies),
		zap.Int("errors", len(summary.Errors)),
	)

	// Log first 10 errors if any (avoid spam)
	for i, scanErr := range summary.Errors {
		if i >= 10 {
			w.logger.Warn("additional errors truncated",
				zap.Int("total_errors", len(summary.Errors)),
			)
			break
		}
		w.logger.Warn("scan error",
			zap.Error(scanErr),
		)
	}

	return nil
}
