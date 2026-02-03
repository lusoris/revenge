package radarr

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// RadarrSyncOperation defines the type of sync operation.
type RadarrSyncOperation string

const (
	// RadarrSyncOperationFull performs a full library sync.
	RadarrSyncOperationFull RadarrSyncOperation = "full"
	// RadarrSyncOperationSingle syncs a single movie.
	RadarrSyncOperationSingle RadarrSyncOperation = "single"
)

// RadarrSyncJobArgs are the arguments for Radarr sync jobs.
type RadarrSyncJobArgs struct {
	// Operation is the type of sync operation to perform.
	Operation RadarrSyncOperation `json:"operation"`
	// RadarrMovieID is the Radarr movie ID for single sync (not used for full).
	RadarrMovieID int `json:"radarr_movie_id,omitempty"`
}

// Kind returns the unique job kind for River.
func (RadarrSyncJobArgs) Kind() string {
	return "radarr_sync"
}

// RadarrSyncWorker handles Radarr library sync operations.
type RadarrSyncWorker struct {
	river.WorkerDefaults[RadarrSyncJobArgs]
	syncService *SyncService
	logger      *zap.Logger
}

// NewRadarrSyncWorker creates a new Radarr sync worker.
func NewRadarrSyncWorker(syncService *SyncService, logger *zap.Logger) *RadarrSyncWorker {
	return &RadarrSyncWorker{
		syncService: syncService,
		logger:      logger.Named("radarr_sync_worker"),
	}
}

// Work executes the Radarr sync job.
func (w *RadarrSyncWorker) Work(ctx context.Context, job *river.Job[RadarrSyncJobArgs]) error {
	args := job.Args

	w.logger.Info("starting radarr sync operation",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("operation", string(args.Operation)),
		zap.Int("radarr_movie_id", args.RadarrMovieID),
	)

	// Check if sync service is available
	if w.syncService == nil {
		w.logger.Debug("radarr sync service not available, skipping")
		return nil
	}

	switch args.Operation {
	case RadarrSyncOperationFull:
		return w.fullSync(ctx)
	case RadarrSyncOperationSingle:
		return w.singleSync(ctx, args.RadarrMovieID)
	default:
		return fmt.Errorf("unknown operation: %s", args.Operation)
	}
}

// fullSync performs a full library sync from Radarr.
func (w *RadarrSyncWorker) fullSync(ctx context.Context) error {
	result, err := w.syncService.SyncLibrary(ctx)
	if err != nil {
		w.logger.Error("full sync failed", zap.Error(err))
		return err
	}

	w.logger.Info("full sync completed",
		zap.Int("added", result.MoviesAdded),
		zap.Int("updated", result.MoviesUpdated),
		zap.Int("skipped", result.MoviesSkipped),
		zap.Int("errors", len(result.Errors)),
		zap.Duration("duration", result.Duration),
	)

	return nil
}

// singleSync syncs a single movie from Radarr.
func (w *RadarrSyncWorker) singleSync(ctx context.Context, radarrMovieID int) error {
	if err := w.syncService.SyncMovie(ctx, radarrMovieID); err != nil {
		w.logger.Error("single movie sync failed",
			zap.Int("radarr_movie_id", radarrMovieID),
			zap.Error(err),
		)
		return err
	}

	w.logger.Info("single movie sync completed",
		zap.Int("radarr_movie_id", radarrMovieID),
	)

	return nil
}

// RadarrWebhookJobArgs are the arguments for Radarr webhook processing jobs.
type RadarrWebhookJobArgs struct {
	// Payload is the webhook payload from Radarr.
	Payload WebhookPayload `json:"payload"`
}

// Kind returns the unique job kind for River.
func (RadarrWebhookJobArgs) Kind() string {
	return "radarr_webhook"
}

// RadarrWebhookWorker handles Radarr webhook processing.
type RadarrWebhookWorker struct {
	river.WorkerDefaults[RadarrWebhookJobArgs]
	webhookHandler *WebhookHandler
	logger         *zap.Logger
}

// NewRadarrWebhookWorker creates a new Radarr webhook worker.
func NewRadarrWebhookWorker(webhookHandler *WebhookHandler, logger *zap.Logger) *RadarrWebhookWorker {
	return &RadarrWebhookWorker{
		webhookHandler: webhookHandler,
		logger:         logger.Named("radarr_webhook_worker"),
	}
}

// Work executes the Radarr webhook processing job.
func (w *RadarrWebhookWorker) Work(ctx context.Context, job *river.Job[RadarrWebhookJobArgs]) error {
	args := job.Args

	w.logger.Info("processing radarr webhook",
		zap.String("job_id", fmt.Sprintf("%d", job.ID)),
		zap.String("event_type", args.Payload.EventType),
		zap.Int("movie_id", args.Payload.Movie.ID),
	)

	// Check if webhook handler is available
	if w.webhookHandler == nil {
		w.logger.Debug("radarr webhook handler not available, skipping")
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
