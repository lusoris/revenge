package activity

import (
	"context"
	"fmt"
	"time"

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

// ActivityCleanupWorker performs periodic activity log cleanup.
type ActivityCleanupWorker struct {
	river.WorkerDefaults[ActivityCleanupArgs]
	service *Service
	logger  *zap.Logger
}

// NewActivityCleanupWorker creates a new activity cleanup worker.
func NewActivityCleanupWorker(service *Service, logger *zap.Logger) *ActivityCleanupWorker {
	return &ActivityCleanupWorker{
		service: service,
		logger:  logger.Named("activity-cleanup"),
	}
}

// Work executes the activity cleanup job.
func (w *ActivityCleanupWorker) Work(ctx context.Context, job *river.Job[ActivityCleanupArgs]) error {
	args := job.Args

	w.logger.Info("starting activity cleanup job",
		zap.Int64("job_id", job.ID),
		zap.Int("retention_days", args.RetentionDays),
		zap.Bool("dry_run", args.DryRun),
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
