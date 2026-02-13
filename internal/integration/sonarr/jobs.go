package sonarr

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/riverqueue/river"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
)

// SonarrSyncOperation defines the type of sync operation.
type SonarrSyncOperation string

const (
	// SonarrSyncOperationFull performs a full library sync.
	SonarrSyncOperationFull SonarrSyncOperation = "full"
	// SonarrSyncOperationSingle syncs a single series.
	SonarrSyncOperationSingle SonarrSyncOperation = "single"
)

// SonarrSyncJobArgs are the arguments for Sonarr sync jobs.
type SonarrSyncJobArgs struct {
	// Operation is the type of sync operation to perform.
	Operation SonarrSyncOperation `json:"operation"`
	// SonarrSeriesID is the Sonarr series ID for single sync (not used for full).
	SonarrSeriesID int `json:"sonarr_series_id,omitempty"`
}

// Kind returns the unique job kind for River.
func (SonarrSyncJobArgs) Kind() string {
	return "sonarr_sync"
}

// InsertOpts returns the default insert options for Sonarr sync jobs.
// Sonarr sync runs on the high-priority queue for responsive integration.
func (SonarrSyncJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueHigh,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 5 * time.Minute,
		},
	}
}

// SonarrSyncWorker handles Sonarr library sync operations.
type SonarrSyncWorker struct {
	river.WorkerDefaults[SonarrSyncJobArgs]
	syncService *SyncService
	logger      *slog.Logger
}

// NewSonarrSyncWorker creates a new Sonarr sync worker.
func NewSonarrSyncWorker(syncService *SyncService, logger *slog.Logger) *SonarrSyncWorker {
	return &SonarrSyncWorker{
		syncService: syncService,
		logger:      logger.With("component", "sonarr_sync_worker"),
	}
}

// Timeout returns the maximum execution time for Sonarr sync jobs.
func (w *SonarrSyncWorker) Timeout(job *river.Job[SonarrSyncJobArgs]) time.Duration {
	return 10 * time.Minute
}

// Work executes the Sonarr sync job.
func (w *SonarrSyncWorker) Work(ctx context.Context, job *river.Job[SonarrSyncJobArgs]) error {
	args := job.Args

	w.logger.Info("starting sonarr sync operation",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("operation", string(args.Operation)),
		slog.Int("sonarr_series_id", args.SonarrSeriesID),
	)

	// Check if sync service is available
	if w.syncService == nil {
		w.logger.Debug("sonarr sync service not available, skipping")
		return nil
	}

	switch args.Operation {
	case SonarrSyncOperationFull:
		return w.fullSync(ctx)
	case SonarrSyncOperationSingle:
		return w.singleSync(ctx, args.SonarrSeriesID)
	default:
		return fmt.Errorf("unknown operation: %s", args.Operation)
	}
}

// fullSync performs a full library sync from Sonarr.
func (w *SonarrSyncWorker) fullSync(ctx context.Context) error {
	result, err := w.syncService.SyncLibrary(ctx)
	if err != nil {
		w.logger.Error("full sync failed", slog.Any("error", err))
		return err
	}

	w.logger.Info("full sync completed",
		slog.Int("series_added", result.SeriesAdded),
		slog.Int("series_updated", result.SeriesUpdated),
		slog.Int("series_skipped", result.SeriesSkipped),
		slog.Int("episodes_added", result.EpisodesAdded),
		slog.Int("episodes_updated", result.EpisodesUpdated),
		slog.Int("errors", len(result.Errors)),
		slog.Duration("duration", result.Duration),
	)

	return nil
}

// singleSync syncs a single series from Sonarr.
func (w *SonarrSyncWorker) singleSync(ctx context.Context, sonarrSeriesID int) error {
	if err := w.syncService.SyncSeries(ctx, sonarrSeriesID); err != nil {
		w.logger.Error("single series sync failed",
			slog.Int("sonarr_series_id", sonarrSeriesID),
			slog.Any("error", err),
		)
		return err
	}

	w.logger.Info("single series sync completed",
		slog.Int("sonarr_series_id", sonarrSeriesID),
	)

	return nil
}

// SonarrWebhookJobArgs are the arguments for Sonarr webhook processing jobs.
type SonarrWebhookJobArgs struct {
	// Payload is the webhook payload from Sonarr.
	Payload WebhookPayload `json:"payload"`
}

// Kind returns the unique job kind for River.
func (SonarrWebhookJobArgs) Kind() string {
	return "sonarr_webhook"
}

// InsertOpts returns the default insert options for Sonarr webhook jobs.
// Webhooks run on the high-priority queue for responsive user experience.
func (SonarrWebhookJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueHigh,
		MaxAttempts: 5,
	}
}

// SonarrWebhookWorker handles Sonarr webhook processing.
type SonarrWebhookWorker struct {
	river.WorkerDefaults[SonarrWebhookJobArgs]
	webhookHandler *WebhookHandler
	logger         *slog.Logger
}

// NewSonarrWebhookWorker creates a new Sonarr webhook worker.
func NewSonarrWebhookWorker(webhookHandler *WebhookHandler, logger *slog.Logger) *SonarrWebhookWorker {
	return &SonarrWebhookWorker{
		webhookHandler: webhookHandler,
		logger:         logger.With("component", "sonarr_webhook_worker"),
	}
}

// Timeout returns the maximum execution time for Sonarr webhook jobs.
func (w *SonarrWebhookWorker) Timeout(job *river.Job[SonarrWebhookJobArgs]) time.Duration {
	return 1 * time.Minute
}

// Work executes the Sonarr webhook processing job.
func (w *SonarrWebhookWorker) Work(ctx context.Context, job *river.Job[SonarrWebhookJobArgs]) error {
	args := job.Args

	seriesID := 0
	if args.Payload.Series != nil {
		seriesID = args.Payload.Series.ID
	}

	w.logger.Info("processing sonarr webhook",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("event_type", args.Payload.EventType),
		slog.Int("series_id", seriesID),
	)

	// Check if webhook handler is available
	if w.webhookHandler == nil {
		w.logger.Debug("sonarr webhook handler not available, skipping")
		return nil
	}

	if err := w.webhookHandler.HandleWebhook(ctx, &args.Payload); err != nil {
		w.logger.Error("webhook processing failed",
			slog.String("event_type", args.Payload.EventType),
			slog.Any("error", err),
		)
		return err
	}

	w.logger.Info("webhook processed successfully",
		slog.String("event_type", args.Payload.EventType),
	)

	return nil
}
