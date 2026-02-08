package radarr

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
	"log/slog"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
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

// InsertOpts returns the default insert options for Radarr sync jobs.
// Radarr sync runs on the high-priority queue for responsive integration.
func (RadarrSyncJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueHigh,
	}
}

// RadarrSyncWorker handles Radarr library sync operations.
type RadarrSyncWorker struct {
	river.WorkerDefaults[RadarrSyncJobArgs]
	syncService *SyncService
	logger      *slog.Logger
}

// NewRadarrSyncWorker creates a new Radarr sync worker.
func NewRadarrSyncWorker(syncService *SyncService, logger *slog.Logger) *RadarrSyncWorker {
	return &RadarrSyncWorker{
		syncService: syncService,
		logger:      logger.With("component", "radarr_sync_worker"),
	}
}

// Timeout returns the maximum execution time for Radarr sync jobs.
func (w *RadarrSyncWorker) Timeout(job *river.Job[RadarrSyncJobArgs]) time.Duration {
	return 10 * time.Minute
}

// Work executes the Radarr sync job.
func (w *RadarrSyncWorker) Work(ctx context.Context, job *river.Job[RadarrSyncJobArgs]) error {
	args := job.Args

	w.logger.Info("starting radarr sync operation",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("operation", string(args.Operation)),
		slog.Int("radarr_movie_id", args.RadarrMovieID),
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
		w.logger.Error("full sync failed", slog.Any("error",err))
		return err
	}

	w.logger.Info("full sync completed",
		slog.Int("added", result.MoviesAdded),
		slog.Int("updated", result.MoviesUpdated),
		slog.Int("skipped", result.MoviesSkipped),
		slog.Int("errors", len(result.Errors)),
		slog.Duration("duration", result.Duration),
	)

	return nil
}

// singleSync syncs a single movie from Radarr.
func (w *RadarrSyncWorker) singleSync(ctx context.Context, radarrMovieID int) error {
	if err := w.syncService.SyncMovie(ctx, radarrMovieID); err != nil {
		w.logger.Error("single movie sync failed",
			slog.Int("radarr_movie_id", radarrMovieID),
			slog.Any("error",err),
		)
		return err
	}

	w.logger.Info("single movie sync completed",
		slog.Int("radarr_movie_id", radarrMovieID),
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

// InsertOpts returns the default insert options for Radarr webhook jobs.
// Webhooks run on the high-priority queue for responsive user experience.
func (RadarrWebhookJobArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: infrajobs.QueueHigh,
	}
}

// RadarrWebhookWorker handles Radarr webhook processing.
type RadarrWebhookWorker struct {
	river.WorkerDefaults[RadarrWebhookJobArgs]
	webhookHandler *WebhookHandler
	logger         *slog.Logger
}

// NewRadarrWebhookWorker creates a new Radarr webhook worker.
func NewRadarrWebhookWorker(webhookHandler *WebhookHandler, logger *slog.Logger) *RadarrWebhookWorker {
	return &RadarrWebhookWorker{
		webhookHandler: webhookHandler,
		logger:         logger.With("component", "radarr_webhook_worker"),
	}
}

// Timeout returns the maximum execution time for Radarr webhook jobs.
func (w *RadarrWebhookWorker) Timeout(job *river.Job[RadarrWebhookJobArgs]) time.Duration {
	return 1 * time.Minute
}

// Work executes the Radarr webhook processing job.
func (w *RadarrWebhookWorker) Work(ctx context.Context, job *river.Job[RadarrWebhookJobArgs]) error {
	args := job.Args

	var movieID int
	if args.Payload.Movie != nil {
		movieID = args.Payload.Movie.ID
	}
	w.logger.Info("processing radarr webhook",
		slog.String("job_id", fmt.Sprintf("%d", job.ID)),
		slog.String("event_type", args.Payload.EventType),
		slog.Int("movie_id", movieID),
	)

	// Check if webhook handler is available
	if w.webhookHandler == nil {
		w.logger.Debug("radarr webhook handler not available, skipping")
		return nil
	}

	if err := w.webhookHandler.HandleWebhook(ctx, &args.Payload); err != nil {
		w.logger.Error("webhook processing failed",
			slog.String("event_type", args.Payload.EventType),
			slog.Any("error",err),
		)
		return err
	}

	w.logger.Info("webhook processed successfully",
		slog.String("event_type", args.Payload.EventType),
	)

	return nil
}
