package moviejobs

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
	"log/slog"

	"github.com/lusoris/revenge/internal/content/movie"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/observability"
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

// InsertOpts returns the default insert options for movie library scan jobs.
// Library scans run on the bulk queue since they're resource-intensive batch operations.
func (MovieLibraryScanArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueBulk,
	}
}

// MovieLibraryScanWorker is a worker that scans movie libraries.
type MovieLibraryScanWorker struct {
	river.WorkerDefaults[MovieLibraryScanArgs]
	libraryService *movie.LibraryService
	logger         *slog.Logger
}

// NewMovieLibraryScanWorker creates a new movie library scan worker.
func NewMovieLibraryScanWorker(
	libraryService *movie.LibraryService,
	logger *slog.Logger,
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

// Timeout returns the maximum execution time for movie library scan jobs.
func (w *MovieLibraryScanWorker) Timeout(job *river.Job[MovieLibraryScanArgs]) time.Duration {
	return 30 * time.Minute
}

// Work performs the movie library scan job.
func (w *MovieLibraryScanWorker) Work(ctx context.Context, job *river.Job[MovieLibraryScanArgs]) error {
	args := job.Args
	scanStart := time.Now()

	w.logger.Info("starting movie library scan",
		slog.Any("paths", args.Paths),
		slog.Bool("force", args.Force),
	)

	// Call library service to scan the library
	summary, err := w.libraryService.ScanLibrary(ctx)
	if err != nil {
		w.logger.Error("library scan failed",
			slog.Any("error",err),
		)
		observability.LibraryScanErrorsTotal.WithLabelValues("movies", "fatal").Inc()
		return fmt.Errorf("library scan failed: %w", err)
	}

	observability.LibraryScanDuration.WithLabelValues("movies").Observe(time.Since(scanStart).Seconds())
	observability.LibraryFilesScanned.WithLabelValues("movies").Add(float64(summary.TotalFiles))
	if len(summary.Errors) > 0 {
		observability.LibraryScanErrorsTotal.WithLabelValues("movies", "scan").Add(float64(len(summary.Errors)))
	}

	// Log summary
	w.logger.Info("library scan completed",
		slog.Int("total_files", summary.TotalFiles),
		slog.Int("matched_files", summary.MatchedFiles),
		slog.Int("unmatched_files", summary.UnmatchedFiles),
		slog.Int("new_movies", summary.NewMovies),
		slog.Int("existing_movies", summary.ExistingMovies),
		slog.Int("errors", len(summary.Errors)),
	)

	// Log first 10 errors if any (avoid spam)
	for i, scanErr := range summary.Errors {
		if i >= 10 {
			w.logger.Warn("additional errors truncated",
				slog.Int("total_errors", len(summary.Errors)),
			)
			break
		}
		w.logger.Warn("scan error",
			slog.Any("error",scanErr),
		)
	}

	return nil
}
