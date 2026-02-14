package session

import (
	"context"
	"log/slog"
	"time"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/riverqueue/river"
)

// SessionMaintenanceJobKind is the unique identifier for session maintenance jobs.
const SessionMaintenanceJobKind = "session_maintenance"

// MaintenanceArgs defines the arguments for session maintenance jobs.
type MaintenanceArgs struct{}

// Kind returns the job kind identifier.
func (MaintenanceArgs) Kind() string {
	return SessionMaintenanceJobKind
}

// InsertOpts returns the default insert options.
func (MaintenanceArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueLow,
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByPeriod: 1 * time.Hour,
		},
	}
}

// MaintenanceWorker periodically cleans up expired/revoked sessions
// and reconciles the active sessions Prometheus gauge.
type MaintenanceWorker struct {
	river.WorkerDefaults[MaintenanceArgs]
	service *Service
	logger  *slog.Logger
}

// NewMaintenanceWorker creates a new session maintenance worker.
func NewMaintenanceWorker(service *Service, logger *slog.Logger) *MaintenanceWorker {
	return &MaintenanceWorker{
		service: service,
		logger:  logger.With("component", "session_maintenance"),
	}
}

// Timeout returns the maximum execution time for maintenance jobs.
func (w *MaintenanceWorker) Timeout(_ *river.Job[MaintenanceArgs]) time.Duration {
	return 2 * time.Minute
}

// Work executes the session maintenance job.
func (w *MaintenanceWorker) Work(ctx context.Context, _ *river.Job[MaintenanceArgs]) error {
	// 1. Clean up old expired and revoked sessions
	deleted, err := w.service.CleanupExpiredSessions(ctx)
	if err != nil {
		w.logger.Error("Session cleanup failed", slog.Any("error", err))
		return err
	}

	// 2. Reconcile the active sessions gauge with the real DB count
	w.service.ReconcileSessionGauge(ctx)

	w.logger.Info("Session maintenance completed",
		slog.Int("sessions_deleted", deleted))

	return nil
}
