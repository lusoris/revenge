package activity

import (
	"context"
	"fmt"
	"time"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/raft"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// ActivityCleanupJobKind is the unique identifier for activity cleanup jobs.
const ActivityCleanupJobKind = "activity_cleanup"

// ActivityCleanupArgs defines the arguments for activity cleanup jobs.
type ActivityCleanupArgs struct {
	// RetentionDays specifies how many days of logs to keep
	RetentionDays int `json:"retention_days"`

	// DryRun if true, only logs what would be deleted without actual deletion
	DryRun bool `json:"dry_run,omitempty"`
}

// Kind returns the job kind identifier.
func (ActivityCleanupArgs) Kind() string {
	return ActivityCleanupJobKind
}

// InsertOpts returns the default insert options for activity cleanup jobs.
// Cleanup runs on the low-priority queue since it's maintenance work.
func (ActivityCleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueLow,
	}
}

// ActivityCleanupWorker performs periodic activity log cleanup.
type ActivityCleanupWorker struct {
	river.WorkerDefaults[ActivityCleanupArgs]
	leaderElection *raft.LeaderElection
	service        *Service
	logger         *zap.Logger
}

// NewActivityCleanupWorker creates a new activity cleanup worker.
func NewActivityCleanupWorker(leaderElection *raft.LeaderElection, service *Service, logger *zap.Logger) *ActivityCleanupWorker {
	return &ActivityCleanupWorker{
		leaderElection: leaderElection,
		service:        service,
		logger:         logger.Named("activity-cleanup"),
	}
}

// Timeout returns the maximum execution time for activity cleanup jobs.
func (w *ActivityCleanupWorker) Timeout(job *river.Job[ActivityCleanupArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the activity cleanup job.
func (w *ActivityCleanupWorker) Work(ctx context.Context, job *river.Job[ActivityCleanupArgs]) error {
	args := job.Args

	// Check if this node is the leader (only leader should run cleanup jobs)
	if w.leaderElection != nil && !w.leaderElection.IsLeader() {
		w.logger.Info("skipping activity cleanup job: not the leader node",
			zap.Int64("job_id", job.ID),
			zap.String("leader", w.leaderElection.LeaderAddr()),
		)
		return nil
	}

	w.logger.Info("starting activity cleanup job",
		zap.Int64("job_id", job.ID),
		zap.Int("retention_days", args.RetentionDays),
		zap.Bool("dry_run", args.DryRun),
		zap.Bool("is_leader", w.leaderElection == nil || w.leaderElection.IsLeader()),
	)

	// Validate arguments
	if args.RetentionDays <= 0 {
		args.RetentionDays = 90 // Default to 90 days
	}

	olderThan := time.Now().AddDate(0, 0, -args.RetentionDays)

	if args.DryRun {
		// Count how many logs would be deleted
		count, err := w.service.CountOldLogs(ctx, olderThan)
		if err != nil {
			w.logger.Error("failed to count old activity logs",
				zap.Int64("job_id", job.ID),
				zap.Error(err),
			)
			return fmt.Errorf("failed to count old logs: %w", err)
		}

		w.logger.Info("dry run: would delete activity logs",
			zap.Int64("job_id", job.ID),
			zap.Int64("count", count),
			zap.Time("older_than", olderThan),
		)

		return nil
	}

	// Perform actual cleanup
	deleted, err := w.service.CleanupOldLogs(ctx, olderThan)
	if err != nil {
		w.logger.Error("failed to cleanup activity logs",
			zap.Int64("job_id", job.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cleanup logs: %w", err)
	}

	w.logger.Info("activity cleanup job completed",
		zap.Int64("job_id", job.ID),
		zap.Int64("deleted_count", deleted),
		zap.Time("older_than", olderThan),
	)

	return nil
}

// ScheduleActivityCleanup creates a periodic cleanup job.
// This should be called during application startup to ensure cleanup runs regularly.
func ScheduleActivityCleanup(client *river.Client[any], retentionDays int) error {
	_, err := client.Insert(context.Background(), ActivityCleanupArgs{
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
