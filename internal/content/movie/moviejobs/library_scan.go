package moviejobs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"log/slog"

	"github.com/lusoris/revenge/internal/content/movie"
	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/observability"
	"github.com/lusoris/revenge/internal/service/library"
)

const MovieLibraryScanJobKind = "movie_library_scan"

// MovieLibraryScanArgs are the arguments for the movie library scan job.
type MovieLibraryScanArgs struct {
	ScanID    string   `json:"scan_id"`
	LibraryID string   `json:"library_id"`
	Paths     []string `json:"paths"`
	Force     bool     `json:"force"`
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
	libraryService     *movie.LibraryService
	scanStatusService  *library.Service
	logger             *slog.Logger
}

// NewMovieLibraryScanWorker creates a new movie library scan worker.
func NewMovieLibraryScanWorker(
	libraryService *movie.LibraryService,
	scanStatusService *library.Service,
	logger *slog.Logger,
) *MovieLibraryScanWorker {
	return &MovieLibraryScanWorker{
		libraryService:    libraryService,
		scanStatusService: scanStatusService,
		logger:            logger,
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
		slog.String("scan_id", args.ScanID),
		slog.String("library_id", args.LibraryID),
		slog.Any("paths", args.Paths),
		slog.Bool("force", args.Force),
	)

	// Mark scan as running if we have a scan ID and status service.
	var scanID uuid.UUID
	hasScanID := false
	if args.ScanID != "" && w.scanStatusService != nil {
		var err error
		scanID, err = uuid.Parse(args.ScanID)
		if err == nil {
			hasScanID = true
			if _, err := w.scanStatusService.StartScan(ctx, scanID); err != nil {
				w.logger.Warn("failed to mark scan as running",
					slog.String("scan_id", args.ScanID),
					slog.Any("error", err),
				)
			}
		}
	}

	// Use paths from the job args (from the library record), not from startup config.
	var summary *movie.ScanSummary
	var scanErr error
	if len(args.Paths) > 0 {
		summary, scanErr = w.libraryService.ScanLibraryWithPaths(ctx, args.Paths)
	} else {
		summary, scanErr = w.libraryService.ScanLibrary(ctx)
	}

	if scanErr != nil {
		w.logger.Error("library scan failed",
			slog.Any("error", scanErr),
		)
		observability.LibraryScanErrorsTotal.WithLabelValues("movies", "fatal").Inc()

		// Mark scan as failed.
		if hasScanID {
			if _, err := w.scanStatusService.FailScan(ctx, scanID, scanErr.Error()); err != nil {
				w.logger.Warn("failed to mark scan as failed",
					slog.String("scan_id", args.ScanID),
					slog.Any("error", err),
				)
			}
		}
		return fmt.Errorf("library scan failed: %w", scanErr)
	}

	observability.LibraryScanDuration.WithLabelValues("movies").Observe(time.Since(scanStart).Seconds())
	observability.LibraryFilesScanned.WithLabelValues("movies").Add(float64(summary.TotalFiles))
	if len(summary.Errors) > 0 {
		observability.LibraryScanErrorsTotal.WithLabelValues("movies", "scan").Add(float64(len(summary.Errors)))
	}

	// Mark scan as completed.
	if hasScanID {
		progress := &library.ScanProgress{
			ItemsScanned: int32(summary.TotalFiles),
			ItemsAdded:   int32(summary.NewMovies),
			ItemsUpdated: int32(summary.ExistingMovies),
			ErrorsCount:  int32(len(summary.Errors)),
		}
		if _, err := w.scanStatusService.CompleteScan(ctx, scanID, progress); err != nil {
			w.logger.Warn("failed to mark scan as completed",
				slog.String("scan_id", args.ScanID),
				slog.Any("error", err),
			)
		}
	}

	// Log summary
	w.logger.Info("library scan completed",
		slog.String("scan_id", args.ScanID),
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
			slog.Any("error", scanErr),
		)
	}

	return nil
}
