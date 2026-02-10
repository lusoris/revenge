package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/riverqueue/river"
)

// CleanupJobKind is the unique identifier for cleanup jobs.
const CleanupJobKind = "cleanup"

// Supported cleanup target types.
const (
	CleanupTargetExpiredTokens      = "expired_tokens"
	CleanupTargetRevokedTokens      = "revoked_tokens"
	CleanupTargetPasswordResets     = "password_resets"
	CleanupTargetEmailVerifications = "email_verifications"
	CleanupTargetFailedLogins       = "failed_logins"
	CleanupTargetAll                = "all"
)

// AuthCleanupRepository defines the cleanup methods needed from the auth repository.
type AuthCleanupRepository interface {
	DeleteExpiredAuthTokens(ctx context.Context) error
	DeleteRevokedAuthTokens(ctx context.Context) error
	DeleteExpiredPasswordResetTokens(ctx context.Context) error
	DeleteUsedPasswordResetTokens(ctx context.Context) error
	DeleteExpiredEmailVerificationTokens(ctx context.Context) error
	DeleteVerifiedEmailTokens(ctx context.Context) error
	DeleteOldFailedLoginAttempts(ctx context.Context) error
}

// CleanupArgs defines the arguments for cleanup jobs.
type CleanupArgs struct {
	// TargetType specifies what to clean up (e.g., "expired_tokens", "all")
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

// InsertOpts returns the default insert options for cleanup jobs.
// Cleanup runs on the low-priority queue since it's maintenance work.
func (CleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: QueueLow,
	}
}

// CleanupWorker performs periodic cleanup operations.
// Leader election is handled by River's built-in leader election via the
// river_leader table — periodic jobs only run on the River-elected leader.
type CleanupWorker struct {
	river.WorkerDefaults[CleanupArgs]
	authRepo AuthCleanupRepository
	logger   *slog.Logger
}

// NewCleanupWorker creates a new cleanup worker.
func NewCleanupWorker(authRepo AuthCleanupRepository, logger *slog.Logger) *CleanupWorker {
	if logger == nil {
		logger = slog.Default()
	}
	return &CleanupWorker{
		authRepo: authRepo,
		logger:   logger,
	}
}

// Timeout returns the maximum execution time for cleanup jobs.
func (w *CleanupWorker) Timeout(job *river.Job[CleanupArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the cleanup job.
// Leader election is handled by River's periodic job scheduler — this worker
// is only invoked on the River-elected leader node.
func (w *CleanupWorker) Work(ctx context.Context, job *river.Job[CleanupArgs]) error {
	args := job.Args

	w.logger.Info("starting cleanup job",
		"job_id", job.ID,
		"target_type", args.TargetType,
		"older_than", args.OlderThan,
		"dry_run", args.DryRun,
	)

	// Validate arguments
	if err := w.validateArgs(args); err != nil {
		w.logger.Error("invalid cleanup job arguments",
			"job_id", job.ID,
			"error", err,
		)
		return fmt.Errorf("invalid arguments: %w", err)
	}

	if args.DryRun {
		w.logger.Info("dry run mode: would perform cleanup",
			"job_id", job.ID,
			"target_type", args.TargetType,
		)
		return nil
	}

	var errs []error
	switch args.TargetType {
	case CleanupTargetExpiredTokens:
		errs = w.cleanupExpiredTokens(ctx)
	case CleanupTargetRevokedTokens:
		errs = w.cleanupRevokedTokens(ctx)
	case CleanupTargetPasswordResets:
		errs = w.cleanupPasswordResets(ctx)
	case CleanupTargetEmailVerifications:
		errs = w.cleanupEmailVerifications(ctx)
	case CleanupTargetFailedLogins:
		errs = w.cleanupFailedLogins(ctx)
	case CleanupTargetAll:
		errs = w.cleanupAll(ctx)
	default:
		return fmt.Errorf("unknown target_type: %s", args.TargetType)
	}

	if len(errs) > 0 {
		for _, err := range errs {
			w.logger.Error("cleanup error", "job_id", job.ID, "error", err)
		}
		return fmt.Errorf("cleanup completed with %d errors", len(errs))
	}

	w.logger.Info("cleanup job completed", "job_id", job.ID, "target_type", args.TargetType)
	return nil
}

func (w *CleanupWorker) cleanupExpiredTokens(ctx context.Context) []error {
	var errs []error
	if err := w.authRepo.DeleteExpiredAuthTokens(ctx); err != nil {
		errs = append(errs, fmt.Errorf("expired auth tokens: %w", err))
	}
	return errs
}

func (w *CleanupWorker) cleanupRevokedTokens(ctx context.Context) []error {
	var errs []error
	if err := w.authRepo.DeleteRevokedAuthTokens(ctx); err != nil {
		errs = append(errs, fmt.Errorf("revoked auth tokens: %w", err))
	}
	return errs
}

func (w *CleanupWorker) cleanupPasswordResets(ctx context.Context) []error {
	var errs []error
	if err := w.authRepo.DeleteExpiredPasswordResetTokens(ctx); err != nil {
		errs = append(errs, fmt.Errorf("expired password reset tokens: %w", err))
	}
	if err := w.authRepo.DeleteUsedPasswordResetTokens(ctx); err != nil {
		errs = append(errs, fmt.Errorf("used password reset tokens: %w", err))
	}
	return errs
}

func (w *CleanupWorker) cleanupEmailVerifications(ctx context.Context) []error {
	var errs []error
	if err := w.authRepo.DeleteExpiredEmailVerificationTokens(ctx); err != nil {
		errs = append(errs, fmt.Errorf("expired email verification tokens: %w", err))
	}
	if err := w.authRepo.DeleteVerifiedEmailTokens(ctx); err != nil {
		errs = append(errs, fmt.Errorf("verified email tokens: %w", err))
	}
	return errs
}

func (w *CleanupWorker) cleanupFailedLogins(ctx context.Context) []error {
	var errs []error
	if err := w.authRepo.DeleteOldFailedLoginAttempts(ctx); err != nil {
		errs = append(errs, fmt.Errorf("old failed login attempts: %w", err))
	}
	return errs
}

func (w *CleanupWorker) cleanupAll(ctx context.Context) []error {
	var errs []error
	errs = append(errs, w.cleanupExpiredTokens(ctx)...)
	errs = append(errs, w.cleanupRevokedTokens(ctx)...)
	errs = append(errs, w.cleanupPasswordResets(ctx)...)
	errs = append(errs, w.cleanupEmailVerifications(ctx)...)
	errs = append(errs, w.cleanupFailedLogins(ctx)...)
	return errs
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


