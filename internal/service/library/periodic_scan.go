package library

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/google/uuid"
	rivertype "github.com/riverqueue/river/rivertype"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
)

// PeriodicLibraryScanJobKind is the unique identifier for periodic library scan jobs.
const PeriodicLibraryScanJobKind = "periodic_library_scan"

// PeriodicLibraryScanArgs defines arguments for periodic library scan orchestration jobs.
// This job lists all enabled libraries and enqueues individual scan jobs for each.
type PeriodicLibraryScanArgs struct{}

// Kind returns the job kind identifier.
func (PeriodicLibraryScanArgs) Kind() string {
	return PeriodicLibraryScanJobKind
}

// InsertOpts returns the default insert options for periodic library scan jobs.
func (PeriodicLibraryScanArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByPeriod: 5 * time.Minute,
		},
	}
}

// jobInserter is a minimal interface for inserting River jobs.
type jobInserter interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// PeriodicLibraryScanWorker lists enabled libraries and enqueues individual scan jobs.
type PeriodicLibraryScanWorker struct {
	river.WorkerDefaults[PeriodicLibraryScanArgs]
	repo      Repository
	jobClient jobInserter
	logger    *slog.Logger
}

// NewPeriodicLibraryScanWorker creates a new periodic library scan worker.
func NewPeriodicLibraryScanWorker(repo Repository, jobClient jobInserter, logger *slog.Logger) *PeriodicLibraryScanWorker {
	return &PeriodicLibraryScanWorker{
		repo:      repo,
		jobClient: jobClient,
		logger:    logger.With("component", "periodic-library-scan"),
	}
}

// Timeout returns the maximum execution time for periodic library scan jobs.
func (w *PeriodicLibraryScanWorker) Timeout(job *river.Job[PeriodicLibraryScanArgs]) time.Duration {
	return 1 * time.Minute
}

// Work lists all enabled libraries and enqueues scan jobs for each.
func (w *PeriodicLibraryScanWorker) Work(ctx context.Context, job *river.Job[PeriodicLibraryScanArgs]) error {
	w.logger.Info("starting periodic library scan", slog.Int64("job_id", job.ID))

	libraries, err := w.repo.ListEnabled(ctx)
	if err != nil {
		return fmt.Errorf("failed to list enabled libraries: %w", err)
	}

	if len(libraries) == 0 {
		w.logger.Info("no enabled libraries found, skipping periodic scan")
		return nil
	}

	var enqueued, failed int
	for _, lib := range libraries {
		if err := w.enqueueScan(ctx, lib); err != nil {
			w.logger.Warn("failed to enqueue scan for library",
				slog.String("library_id", lib.ID.String()),
				slog.String("library_name", lib.Name),
				slog.Any("error", err),
			)
			failed++
		} else {
			enqueued++
		}
	}

	w.logger.Info("periodic library scan orchestration complete",
		slog.Int64("job_id", job.ID),
		slog.Int("total_libraries", len(libraries)),
		slog.Int("enqueued", enqueued),
		slog.Int("failed", failed),
	)

	if failed > 0 && enqueued == 0 {
		return fmt.Errorf("failed to enqueue any library scans (%d failures)", failed)
	}

	return nil
}

// enqueueScan enqueues the appropriate scan job for a library based on its type.
func (w *PeriodicLibraryScanWorker) enqueueScan(ctx context.Context, lib Library) error {
	switch lib.Type {
	case LibraryTypeTVShow:
		libID := lib.ID
		_, err := w.jobClient.Insert(ctx, &tvshowScanArgs{
			Paths:      lib.Paths,
			Force:      false,
			LibraryID:  &libID,
			AutoCreate: true,
		}, nil)
		return err
	case LibraryTypeMovie:
		_, err := w.jobClient.Insert(ctx, &movieScanArgs{
			ScanID:    "",
			LibraryID: lib.ID.String(),
			Paths:     lib.Paths,
			Force:     false,
		}, nil)
		return err
	default:
		w.logger.Debug("skipping unsupported library type for periodic scan",
			slog.String("type", lib.Type),
			slog.String("library_id", lib.ID.String()),
		)
		return nil
	}
}

// ---------------------------------------------------------------------------
// Mirror args types to avoid import cycles (moviejobs → library → moviejobs).
// The Kind() values MUST match the constants in the actual worker packages.
// ---------------------------------------------------------------------------

// tvshowScanArgs mirrors tvshowjobs.LibraryScanArgs.
type tvshowScanArgs struct {
	Paths      []string   `json:"paths"`
	Force      bool       `json:"force"`
	LibraryID  *uuid.UUID `json:"library_id,omitempty"`
	AutoCreate bool       `json:"auto_create"`
}

func (tvshowScanArgs) Kind() string { return "tvshow_library_scan" }

// movieScanArgs mirrors moviejobs.MovieLibraryScanArgs.
type movieScanArgs struct {
	ScanID    string   `json:"scan_id"`
	LibraryID string   `json:"library_id"`
	Paths     []string `json:"paths"`
	Force     bool     `json:"force"`
}

func (movieScanArgs) Kind() string { return "movie_library_scan" }
