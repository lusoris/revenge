package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lusoris/revenge/internal/infra/raft"
	"github.com/riverqueue/river"
)

// CleanupJobKind is the unique identifier for cleanup jobs.
const CleanupJobKind = "cleanup"

// CleanupArgs defines the arguments for cleanup jobs.
type CleanupArgs struct {
	// TargetType specifies what to clean up (e.g., "sessions", "jobs", "logs")
	TargetType string `json:"target_type"`

	// OlderThan specifies the age threshold for cleanup
	OlderThan time.Duration `json:"older_than"`

	// BatchSize limits the number of records to delete per batch
	BatchSize int `json:"batch_size,omitempty"`

	// DryRun if true, only logs what would be deleted without actual deletion
	DryRun bool `json:"dry_run,omitempty"`
}

// Kind returns the job kind identifier.
func (CleanupArgs) Kind() string {
	return CleanupJobKind
}

// CleanupWorker performs periodic cleanup operations.
type CleanupWorker struct {
	river.WorkerDefaults[CleanupArgs]
	leaderElection *raft.LeaderElection
	logger         *slog.Logger
}

// NewCleanupWorker creates a new cleanup worker.
func NewCleanupWorker(leaderElection *raft.LeaderElection, logger *slog.Logger) *CleanupWorker {
	if logger == nil {
		logger = slog.Default()
	}
	return &CleanupWorker{
		leaderElection: leaderElection,
		logger:         logger,
	}
}

// Work executes the cleanup job.
func (w *CleanupWorker) Work(ctx context.Context, job *river.Job[CleanupArgs]) error {
	args := job.Args

	// Check if this node is the leader (only leader should run cleanup jobs)
	if w.leaderElection != nil && !w.leaderElection.IsLeader() {
		w.logger.Info("skipping cleanup job: not the leader node",
			"job_id", job.ID,
			"target_type", args.TargetType,
			"leader", w.leaderElection.LeaderAddr(),
		)
		return nil
	}

	w.logger.Info("starting cleanup job",
		"job_id", job.ID,
		"target_type", args.TargetType,
		"older_than", args.OlderThan,
		"batch_size", args.BatchSize,
		"dry_run", args.DryRun,
		"is_leader", w.leaderElection == nil || w.leaderElection.IsLeader(),
	)

	// Validate arguments
	if err := w.validateArgs(args); err != nil {
		w.logger.Error("invalid cleanup job arguments",
			"job_id", job.ID,
			"error", err,
		)
		return fmt.Errorf("invalid arguments: %w", err)
	}

	// Simulate cleanup work
	if args.DryRun {
		w.logger.Info("dry run mode: would delete records",
			"job_id", job.ID,
			"target_type", args.TargetType,
		)
	} else {
		w.logger.Info("performing cleanup",
			"job_id", job.ID,
			"target_type", args.TargetType,
		)
		// Actual cleanup logic would go here (database operations)
		// For now, this is a stub that simulates work
	}

	w.logger.Info("cleanup job completed",
		"job_id", job.ID,
		"target_type", args.TargetType,
	)

	return nil
}

// validateArgs validates cleanup job arguments.
func (w *CleanupWorker) validateArgs(args CleanupArgs) error {
	if args.TargetType == "" {
		return fmt.Errorf("target_type is required")
	}

	if args.OlderThan <= 0 {
		return fmt.Errorf("older_than must be positive")
	}

	if args.BatchSize < 0 {
		return fmt.Errorf("batch_size cannot be negative")
	}

	return nil
}
