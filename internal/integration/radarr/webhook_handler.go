package radarr

import (
	"context"
	"log/slog"
)

// WebhookHandler processes Radarr webhook events.
type WebhookHandler struct {
	syncService *SyncService
	logger      *slog.Logger
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(syncService *SyncService, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		syncService: syncService,
		logger:      logger.With("handler", "radarr_webhook"),
	}
}

// HandleWebhook processes a Radarr webhook payload.
// It handles events like grab, download, rename, movie delete, etc.
func (h *WebhookHandler) HandleWebhook(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("received radarr webhook",
		"event_type", payload.EventType,
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
	)

	switch payload.EventType {
	case EventGrab:
		return h.handleGrab(ctx, payload)
	case EventDownload:
		return h.handleDownload(ctx, payload)
	case EventRename:
		return h.handleRename(ctx, payload)
	case EventMovieDelete:
		return h.handleMovieDelete(ctx, payload)
	case EventMovieFileDelete:
		return h.handleMovieFileDelete(ctx, payload)
	case EventHealth:
		return h.handleHealthIssue(ctx, payload)
	case EventApplicationUpdate:
		return h.handleApplicationUpdate(ctx, payload)
	case EventManualInteractionRequired:
		return h.handleManualInteraction(ctx, payload)
	case EventTest:
		return h.handleTest(ctx, payload)
	default:
		h.logger.Debug("ignoring unknown event type", "event_type", payload.EventType)
		return nil
	}
}

// handleGrab handles the Grab event (movie was grabbed for download).
func (h *WebhookHandler) handleGrab(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("movie grabbed",
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
		"release_title", payload.Release.ReleaseTitle,
	)

	// We don't need to do anything for grabs - wait for download complete
	return nil
}

// handleDownload handles the Download event (movie was downloaded and imported).
func (h *WebhookHandler) handleDownload(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("movie downloaded",
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
		"is_upgrade", payload.IsUpgrade,
	)

	// Sync the movie from Radarr to update our database
	if err := h.syncService.SyncMovie(ctx, payload.Movie.ID); err != nil {
		h.logger.Error("failed to sync movie after download",
			"movie_id", payload.Movie.ID,
			"error", err,
		)
		return err
	}

	h.logger.Debug("synced movie after download", "movie_id", payload.Movie.ID)
	return nil
}

// handleRename handles the Rename event (movie file was renamed).
func (h *WebhookHandler) handleRename(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("movie renamed",
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
	)

	// Re-sync to update file paths
	if err := h.syncService.SyncMovie(ctx, payload.Movie.ID); err != nil {
		h.logger.Error("failed to sync movie after rename",
			"movie_id", payload.Movie.ID,
			"error", err,
		)
		return err
	}

	return nil
}

// handleMovieDelete handles the MovieDelete event (movie was deleted from Radarr).
func (h *WebhookHandler) handleMovieDelete(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("movie deleted from radarr",
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
		"delete_files", payload.DeletedFiles,
	)

	// For now, we just log this. The movie remains in Revenge but is marked as no longer in Radarr.
	// A future enhancement could offer to delete the movie from Revenge as well.
	return nil
}

// handleMovieFileDelete handles the MovieFileDelete event (movie file was deleted).
func (h *WebhookHandler) handleMovieFileDelete(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("movie file deleted",
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
	)

	// Re-sync to remove the file reference
	if err := h.syncService.SyncMovie(ctx, payload.Movie.ID); err != nil {
		h.logger.Error("failed to sync movie after file delete",
			"movie_id", payload.Movie.ID,
			"error", err,
		)
		return err
	}

	return nil
}

// handleHealthIssue handles the HealthIssue event.
func (h *WebhookHandler) handleHealthIssue(_ context.Context, payload *WebhookPayload) error {
	h.logger.Warn("radarr health issue",
		"level", payload.Level,
		"message", payload.Message,
		"wiki_url", payload.WikiURL,
	)
	return nil
}

// handleApplicationUpdate handles the ApplicationUpdate event.
func (h *WebhookHandler) handleApplicationUpdate(_ context.Context, payload *WebhookPayload) error {
	h.logger.Info("radarr application update",
		"previous_version", payload.PreviousVersion,
		"new_version", payload.NewVersion,
	)
	return nil
}

// handleManualInteraction handles the ManualInteractionRequired event.
func (h *WebhookHandler) handleManualInteraction(_ context.Context, payload *WebhookPayload) error {
	h.logger.Warn("manual interaction required",
		"movie_id", payload.Movie.ID,
		"title", payload.Movie.Title,
	)
	return nil
}

// handleTest handles the Test event (webhook test from Radarr).
func (h *WebhookHandler) handleTest(_ context.Context, _ *WebhookPayload) error {
	h.logger.Info("received test webhook from radarr")
	return nil
}
