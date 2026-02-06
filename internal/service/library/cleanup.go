package library

import (
	"context"
	"fmt"
	"time"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/raft"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// LibraryScanCleanupJobKind is the unique identifier for library scan cleanup jobs.
const LibraryScanCleanupJobKind = "library_scan_cleanup"

// LibraryScanCleanupArgs defines the arguments for library scan cleanup jobs.
type LibraryScanCleanupArgs struct {
	// RetentionDays specifies how many days of scan history to keep
	RetentionDays int `json:"retention_days"`

	// DryRun if true, only logs what would be deleted without actual deletion
	DryRun bool `json:"dry_run,omitempty"`
}

// Kind returns the job kind identifier.
func (LibraryScanCleanupArgs) Kind() string {
	return LibraryScanCleanupJobKind
}

// InsertOpts returns the default insert options for library scan cleanup jobs.
// Cleanup runs on the low-priority queue since it's maintenance work.
func (LibraryScanCleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueLow,
	}
}

// LibraryScanCleanupWorker performs periodic library scan history cleanup.
type LibraryScanCleanupWorker struct {
	river.WorkerDefaults[LibraryScanCleanupArgs]
	leaderElection *raft.LeaderElection
	repo           Repository
	logger         *zap.Logger
}

// NewLibraryScanCleanupWorker creates a new library scan cleanup worker.
func NewLibraryScanCleanupWorker(leaderElection *raft.LeaderElection, repo Repository, logger *zap.Logger) *LibraryScanCleanupWorker {
	return &LibraryScanCleanupWorker{
		leaderElection: leaderElection,
		repo:           repo,
		logger:         logger.Named("library-scan-cleanup"),
	}
}

// Timeout returns the maximum execution time for library scan cleanup jobs.
func (w *LibraryScanCleanupWorker) Timeout(job *river.Job[LibraryScanCleanupArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the library scan cleanup job.
func (w *LibraryScanCleanupWorker) Work(ctx context.Context, job *river.Job[LibraryScanCleanupArgs]) error {
	args := job.Args

	// Check if this node is the leader (only leader should run cleanup jobs)
	if w.leaderElection != nil && !w.leaderElection.IsLeader() {
		w.logger.Info("skipping library scan cleanup job: not the leader node",
			zap.Int64("job_id", job.ID),
			zap.String("leader", w.leaderElection.LeaderAddr()),
		)
		return nil
	}

	w.logger.Info("starting library scan cleanup job",
		zap.Int64("job_id", job.ID),
		zap.Int("retention_days", args.RetentionDays),
		zap.Bool("dry_run", args.DryRun),
		zap.Bool("is_leader", w.leaderElection == nil || w.leaderElection.IsLeader()),
	)

	// Validate arguments
	if args.RetentionDays <= 0 {
		args.RetentionDays = 30 // Default to 30 days for scan history
	}

	olderThan := time.Now().AddDate(0, 0, -args.RetentionDays)

	if args.DryRun {
		w.logger.Info("dry run mode: would delete library scan records",
			zap.Int64("job_id", job.ID),
			zap.Time("older_than", olderThan),
		)
		return nil
	}

	// Perform actual cleanup
	deleted, err := w.repo.DeleteOldScans(ctx, olderThan)
	if err != nil {
		w.logger.Error("failed to cleanup library scans",
			zap.Int64("job_id", job.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cleanup scans: %w", err)
	}

	w.logger.Info("library scan cleanup job completed",
		zap.Int64("job_id", job.ID),
		zap.Int64("deleted_count", deleted),
		zap.Time("older_than", olderThan),
	)

	return nil
}

// ScheduleLibraryScanCleanup creates a periodic cleanup job.
// This should be called during application startup to ensure cleanup runs regularly.
func ScheduleLibraryScanCleanup(client *river.Client[any], retentionDays int) error {
	_, err := client.Insert(context.Background(), LibraryScanCleanupArgs{
		RetentionDays: retentionDays,
		DryRun:        false,
	}, &river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 24 * time.Hour, // Run once per day
		},
	})
	return err
}
