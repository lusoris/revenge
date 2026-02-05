package sonarr

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"go.uber.org/zap"
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

// SonarrSyncWorker handles Sonarr library sync operations.
type SonarrSyncWorker struct {
	river.WorkerDefaults[SonarrSyncJobArgs]
	syncService *SyncService
	logger      *zap.Logger
}

// NewSonarrSyncWorker creates a new Sonarr sync worker.
func NewSonarrSyncWorker(syncService *SyncService, logger *zap.Logger) *SonarrSyncWorker {
	return &SonarrSyncWorker{
		syncService: syncService,
		logger:      logger.Named("sonarr_sync_worker"),
	}
}

// Work executes the Sonarr sync job.
func (w *SonarrSyncWorker) Work(ctx context.Context, job *river.Job[SonarrSyncJobArgs]) error {
	args := job.Args

	w.logger.Info("starting sonarr sync operation",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("operation", string(args.Operation)),
		zap.Int("sonarr_series_id", args.SonarrSeriesID),
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
		w.logger.Error("full sync failed", zap.Error(err))
		return err
	}

	w.logger.Info("full sync completed",
		zap.Int("series_added", result.SeriesAdded),
		zap.Int("series_updated", result.SeriesUpdated),
		zap.Int("series_skipped", result.SeriesSkipped),
		zap.Int("episodes_added", result.EpisodesAdded),
		zap.Int("episodes_updated", result.EpisodesUpdated),
		zap.Int("errors", len(result.Errors)),
		zap.Duration("duration", result.Duration),
	)

	return nil
}

// singleSync syncs a single series from Sonarr.
func (w *SonarrSyncWorker) singleSync(ctx context.Context, sonarrSeriesID int) error {
	if err := w.syncService.SyncSeries(ctx, sonarrSeriesID); err != nil {
		w.logger.Error("single series sync failed",
			zap.Int("sonarr_series_id", sonarrSeriesID),
			zap.Error(err),
		)
		return err
	}

	w.logger.Info("single series sync completed",
		zap.Int("sonarr_series_id", sonarrSeriesID),
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

// SonarrWebhookWorker handles Sonarr webhook processing.
type SonarrWebhookWorker struct {
	river.WorkerDefaults[SonarrWebhookJobArgs]
	webhookHandler *WebhookHandler
	logger         *zap.Logger
}

// NewSonarrWebhookWorker creates a new Sonarr webhook worker.
func NewSonarrWebhookWorker(webhookHandler *WebhookHandler, logger *zap.Logger) *SonarrWebhookWorker {
	return &SonarrWebhookWorker{
		webhookHandler: webhookHandler,
		logger:         logger.Named("sonarr_webhook_worker"),
	}
}

// Work executes the Sonarr webhook processing job.
func (w *SonarrWebhookWorker) Work(ctx context.Context, job *river.Job[SonarrWebhookJobArgs]) error {
	args := job.Args

	seriesID := 0
	if args.Payload.Series != nil {
		seriesID = args.Payload.Series.ID
	}

	w.logger.Info("processing sonarr webhook",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("event_type", args.Payload.EventType),
		zap.Int("series_id", seriesID),
	)

	// Check if webhook handler is available
	if w.webhookHandler == nil {
		w.logger.Debug("sonarr webhook handler not available, skipping")
		return nil
	}

	if err := w.webhookHandler.HandleWebhook(ctx, &args.Payload); err != nil {
		w.logger.Error("webhook processing failed",
			zap.String("event_type", args.Payload.EventType),
			zap.Error(err),
		)
		return err
	}

	w.logger.Info("webhook processed successfully",
		zap.String("event_type", args.Payload.EventType),
	)

	return nil
}
