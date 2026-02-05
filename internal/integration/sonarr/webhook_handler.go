package sonarr

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lusoris/revenge/internal/util"
)

// WebhookHandler handles Sonarr webhook events.
type WebhookHandler struct {
	syncService *SyncService
	logger      *slog.Logger
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(syncService *SyncService, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		syncService: syncService,
		logger:      logger.With("handler", "sonarr_webhook"),
	}
}

// HandleWebhook processes a Sonarr webhook payload.
func (h *WebhookHandler) HandleWebhook(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("received webhook",
		"event_type", payload.EventType,
		"instance", payload.InstanceName,
	)

	switch payload.EventType {
	case EventTest:
		return h.handleTest(ctx, payload)
	case EventGrab:
		return h.handleGrab(ctx, payload)
	case EventDownload:
		return h.handleDownload(ctx, payload)
	case EventRename:
		return h.handleRename(ctx, payload)
	case EventSeriesAdd:
		return h.handleSeriesAdd(ctx, payload)
	case EventSeriesDelete:
		return h.handleSeriesDelete(ctx, payload)
	case EventEpisodeFileDelete:
		return h.handleEpisodeFileDelete(ctx, payload)
	case EventHealth:
		return h.handleHealth(ctx, payload)
	case EventHealthRestored:
		return h.handleHealthRestored(ctx, payload)
	case EventApplicationUpdate:
		return h.handleApplicationUpdate(ctx, payload)
	case EventManualInteractionRequired:
		return h.handleManualInteractionRequired(ctx, payload)
	default:
		h.logger.Warn("unknown webhook event type", "event_type", payload.EventType)
		return nil
	}
}

// handleTest handles the test webhook event.
func (h *WebhookHandler) handleTest(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("test webhook received", "instance", payload.InstanceName)
	return nil
}

// handleGrab handles the grab event (release was grabbed/sent to download client).
func (h *WebhookHandler) handleGrab(ctx context.Context, payload *WebhookPayload) error {
	if payload.Series == nil {
		return fmt.Errorf("grab event missing series info")
	}

	h.logger.Info("episode grab event",
		"series_id", payload.Series.ID,
		"series_title", payload.Series.Title,
		"episode_count", len(payload.Episodes),
		"release", payload.Release,
	)

	// We don't need to do much here - wait for the download event
	return nil
}

// handleDownload handles the download event (file was downloaded/imported).
func (h *WebhookHandler) handleDownload(ctx context.Context, payload *WebhookPayload) error {
	if payload.Series == nil {
		return fmt.Errorf("download event missing series info")
	}

	h.logger.Info("episode download event",
		"series_id", payload.Series.ID,
		"series_title", payload.Series.Title,
		"episode_count", len(payload.Episodes),
		"is_upgrade", payload.IsUpgrade,
	)

	// Sync the series to get the new file
	if err := h.syncService.SyncSeries(ctx, payload.Series.ID); err != nil {
		h.logger.Error("failed to sync series after download",
			"series_id", payload.Series.ID,
			"error", err,
		)
		return fmt.Errorf("sync series: %w", err)
	}

	return nil
}

// handleRename handles the rename event (files were renamed).
func (h *WebhookHandler) handleRename(ctx context.Context, payload *WebhookPayload) error {
	if payload.Series == nil {
		return fmt.Errorf("rename event missing series info")
	}

	h.logger.Info("episode rename event",
		"series_id", payload.Series.ID,
		"series_title", payload.Series.Title,
	)

	// Re-sync to get updated file paths
	if err := h.syncService.SyncSeries(ctx, payload.Series.ID); err != nil {
		h.logger.Error("failed to sync series after rename",
			"series_id", payload.Series.ID,
			"error", err,
		)
		return fmt.Errorf("sync series: %w", err)
	}

	return nil
}

// handleSeriesAdd handles the series add event.
func (h *WebhookHandler) handleSeriesAdd(ctx context.Context, payload *WebhookPayload) error {
	if payload.Series == nil {
		return fmt.Errorf("series add event missing series info")
	}

	h.logger.Info("series add event",
		"series_id", payload.Series.ID,
		"series_title", payload.Series.Title,
		"tvdb_id", payload.Series.TVDbID,
	)

	// Sync the new series
	if err := h.syncService.SyncSeries(ctx, payload.Series.ID); err != nil {
		h.logger.Error("failed to sync new series",
			"series_id", payload.Series.ID,
			"error", err,
		)
		return fmt.Errorf("sync new series: %w", err)
	}

	return nil
}

// handleSeriesDelete handles the series delete event.
func (h *WebhookHandler) handleSeriesDelete(ctx context.Context, payload *WebhookPayload) error {
	if payload.Series == nil {
		return fmt.Errorf("series delete event missing series info")
	}

	h.logger.Info("series delete event",
		"series_id", payload.Series.ID,
		"series_title", payload.Series.Title,
	)

	// We might want to mark the series as removed in our DB
	// For now, we just log it - actual deletion handled separately
	existing, err := h.syncService.tvshowRepo.GetSeriesBySonarrID(ctx, util.SafeIntToInt32(payload.Series.ID))
	if err != nil {
		h.logger.Warn("series not found in local db", "sonarr_id", payload.Series.ID)
		return nil
	}

	h.logger.Info("series marked for deletion",
		"series_id", existing.ID,
		"sonarr_id", payload.Series.ID,
	)

	// Note: We don't automatically delete - that's a policy decision
	// The series will be marked as removed on the next full sync

	return nil
}

// handleEpisodeFileDelete handles the episode file delete event.
func (h *WebhookHandler) handleEpisodeFileDelete(ctx context.Context, payload *WebhookPayload) error {
	if payload.Series == nil {
		return fmt.Errorf("episode file delete event missing series info")
	}

	h.logger.Info("episode file delete event",
		"series_id", payload.Series.ID,
		"series_title", payload.Series.Title,
		"deleted_files", len(payload.DeletedFiles),
	)

	// Re-sync to reflect the deleted files
	if err := h.syncService.SyncSeries(ctx, payload.Series.ID); err != nil {
		h.logger.Error("failed to sync series after file delete",
			"series_id", payload.Series.ID,
			"error", err,
		)
		return fmt.Errorf("sync series: %w", err)
	}

	return nil
}

// handleHealth handles health check failure events.
func (h *WebhookHandler) handleHealth(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Warn("Sonarr health check issue",
		"level", payload.Level,
		"message", payload.Message,
		"type", payload.Type,
		"wiki_url", payload.WikiURL,
	)
	return nil
}

// handleHealthRestored handles health restored events.
func (h *WebhookHandler) handleHealthRestored(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("Sonarr health restored",
		"message", payload.Message,
		"type", payload.Type,
	)
	return nil
}

// handleApplicationUpdate handles application update events.
func (h *WebhookHandler) handleApplicationUpdate(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Info("Sonarr application updated",
		"previous_version", payload.PreviousVersion,
		"new_version", payload.NewVersion,
	)
	return nil
}

// handleManualInteractionRequired handles manual interaction required events.
func (h *WebhookHandler) handleManualInteractionRequired(ctx context.Context, payload *WebhookPayload) error {
	h.logger.Warn("Sonarr requires manual interaction",
		"message", payload.Message,
		"download_client", payload.DownloadClient,
		"download_id", payload.DownloadID,
	)
	return nil
}
