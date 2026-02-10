package activity

import (
	"context"
	"log/slog"
	"strings"
	"time"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
)

// ActivityLogJobKind is the unique identifier for activity log jobs.
const ActivityLogJobKind = "activity_log"

// ActivityLogArgs defines the arguments for activity log jobs.
// The embedded LogRequest is JSON-serializable and contains
// all data needed to persist the activity entry.
type ActivityLogArgs struct {
	LogRequest `json:"log_request"`
}

// Kind returns the job kind identifier.
func (ActivityLogArgs) Kind() string {
	return ActivityLogJobKind
}

// InsertOpts returns the default insert options for activity log jobs.
// By default, activity logs are low-priority maintenance work.
// Security-critical events are overridden to QueueCritical at insertion time.
func (ActivityLogArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueLow,
		MaxAttempts: 3,
	}
}

// securityActions are auth/security actions that must be processed
// on the critical queue for timely audit trail and brute-force detection.
var securityActions = map[string]bool{
	ActionUserLogin:         true,
	ActionUserLogout:        true,
	ActionUserCreate:        true,
	ActionUserDelete:        true,
	ActionUserPasswordReset: true,
	ActionSessionCreate:     true,
	ActionSessionRevoke:     true,
	ActionOIDCLogin:         true,
	ActionOIDCLink:          true,
	ActionOIDCUnlink:        true,
	ActionAdminRoleAssign:   true,
	ActionAdminRoleRevoke:   true,
	ActionAdminUserBan:      true,
	ActionAdminUserUnban:    true,
}

// isSecurityAction returns true if the action is security-critical
// and should be processed on the critical queue.
func isSecurityAction(action string) bool {
	if securityActions[action] {
		return true
	}
	// All failed actions are security-relevant (brute-force, abuse)
	if strings.HasSuffix(action, ".failed") {
		return true
	}
	return false
}

// ActivityLogWorker processes activity log jobs.
type ActivityLogWorker struct {
	river.WorkerDefaults[ActivityLogArgs]
	service *Service
	logger  *slog.Logger
}

// NewActivityLogWorker creates a new activity log worker.
func NewActivityLogWorker(service *Service, logger *slog.Logger) *ActivityLogWorker {
	return &ActivityLogWorker{
		service: service,
		logger:  logger.With("component", "activity-log-worker"),
	}
}

// Timeout returns the maximum execution time for activity log jobs.
func (w *ActivityLogWorker) Timeout(_ *river.Job[ActivityLogArgs]) time.Duration {
	return 10 * time.Second
}

// Work processes an activity log job by persisting the entry to the database.
func (w *ActivityLogWorker) Work(ctx context.Context, job *river.Job[ActivityLogArgs]) error {
	if err := w.service.Log(ctx, job.Args.LogRequest); err != nil {
		w.logger.Error("failed to persist activity log",
			slog.String("action", job.Args.LogRequest.Action),
			slog.Any("error", err),
		)
		return err
	}
	return nil
}
