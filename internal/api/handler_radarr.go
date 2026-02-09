package api

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/integration/radarr"
)

// radarrService is an optional dependency for Radarr integration.
// When nil, all Radarr endpoints return 503 Service Unavailable.
type radarrService interface {
	GetStatus() radarr.SyncStatus
	IsHealthy(ctx context.Context) bool
	GetSystemStatus(ctx context.Context) (*radarr.SystemStatus, error)
	GetQualityProfiles(ctx context.Context) ([]radarr.QualityProfile, error)
	GetRootFolders(ctx context.Context) ([]radarr.RootFolder, error)
	SyncLibrary(ctx context.Context) (*radarr.SyncResult, error)
	LookupMovie(ctx context.Context, term string) ([]radarr.Movie, error)
}

// AdminGetRadarrStatus returns the current Radarr integration status.
// GET /api/v1/admin/integrations/radarr/status
func (h *Handler) AdminGetRadarrStatus(ctx context.Context) (ogen.AdminGetRadarrStatusRes, error) {
	h.logger.Debug("AdminGetRadarrStatus called")

	// Check admin authorization
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminGetRadarrStatusUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminGetRadarrStatusForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	// Check if Radarr integration is configured
	if h.radarrService == nil {
		return &ogen.AdminGetRadarrStatusServiceUnavailable{
			Code:    503,
			Message: "Radarr integration not configured",
		}, nil
	}

	// Check connection health
	connected := h.radarrService.IsHealthy(ctx)

	// Get sync status
	syncStatus := h.radarrService.GetStatus()

	response := &ogen.RadarrStatus{
		Connected: connected,
		SyncStatus: ogen.RadarrSyncStatus{
			IsRunning:     syncStatus.IsRunning,
			MoviesAdded:   syncStatus.MoviesAdded,
			MoviesUpdated: syncStatus.MoviesUpdated,
			MoviesRemoved: syncStatus.MoviesRemoved,
			TotalMovies:   syncStatus.TotalMovies,
		},
	}

	// Add last sync time if available
	if !syncStatus.LastSync.IsZero() {
		response.SyncStatus.LastSync = ogen.NewOptDateTime(syncStatus.LastSync)
	}

	// Add last sync error if present
	if syncStatus.LastSyncError != "" {
		response.SyncStatus.LastSyncError = ogen.NewOptString(syncStatus.LastSyncError)
	}

	// Get system status from Radarr for version info
	if connected {
		if status, err := h.radarrService.GetSystemStatus(ctx); err == nil {
			response.Version = ogen.NewOptString(status.Version)
			response.InstanceName = ogen.NewOptString(status.InstanceName)
			if status.StartTime != "" {
				if t, err := time.Parse(time.RFC3339, status.StartTime); err == nil {
					response.StartTime = ogen.NewOptDateTime(t)
				}
			}
		}
	}

	h.logger.Info("Radarr status retrieved",
		slog.Bool("connected", connected),
		slog.Bool("sync_running", syncStatus.IsRunning))

	return response, nil
}

// AdminTriggerRadarrSync triggers a full library sync from Radarr.
// POST /api/v1/admin/integrations/radarr/sync
func (h *Handler) AdminTriggerRadarrSync(ctx context.Context) (ogen.AdminTriggerRadarrSyncRes, error) {
	h.logger.Debug("AdminTriggerRadarrSync called")

	// Check admin authorization
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminTriggerRadarrSyncUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminTriggerRadarrSyncForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	// Check if Radarr integration is configured
	if h.radarrService == nil {
		return &ogen.AdminTriggerRadarrSyncServiceUnavailable{
			Code:    503,
			Message: "Radarr integration not configured",
		}, nil
	}

	// Check if sync is already in progress
	status := h.radarrService.GetStatus()
	if status.IsRunning {
		return &ogen.AdminTriggerRadarrSyncConflict{
			Code:    409,
			Message: "Sync already in progress",
		}, nil
	}

	// If we have a River client, queue the job
	if h.riverClient != nil {
		_, err := h.riverClient.Insert(ctx, &radarr.RadarrSyncJobArgs{
			Operation: radarr.RadarrSyncOperationFull,
		}, nil)
		if err != nil {
			h.logger.Error("Failed to queue Radarr sync job", slog.Any("error",err))
			return &ogen.AdminTriggerRadarrSyncServiceUnavailable{
				Code:    503,
				Message: "Failed to queue sync job",
			}, nil
		}

		h.logger.Info("Radarr sync job queued")
		return &ogen.RadarrSyncResponse{
			Message: "Sync job queued",
			Status:  ogen.RadarrSyncResponseStatusQueued,
		}, nil
	}

	// No River client, start sync directly (blocking)
	go func() {
		// Use a new context with timeout since the request context will be done
		syncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if _, err := h.radarrService.SyncLibrary(syncCtx); err != nil {
			h.logger.Error("Radarr sync failed", slog.Any("error",err))
		}
	}()

	h.logger.Info("Radarr sync started directly")
	return &ogen.RadarrSyncResponse{
		Message: "Sync started",
		Status:  ogen.RadarrSyncResponseStatusStarted,
	}, nil
}

// AdminGetRadarrQualityProfiles returns all quality profiles from Radarr.
// GET /api/v1/admin/integrations/radarr/quality-profiles
func (h *Handler) AdminGetRadarrQualityProfiles(ctx context.Context) (ogen.AdminGetRadarrQualityProfilesRes, error) {
	h.logger.Debug("AdminGetRadarrQualityProfiles called")

	// Check admin authorization
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminGetRadarrQualityProfilesUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminGetRadarrQualityProfilesForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	// Check if Radarr integration is configured
	if h.radarrService == nil {
		return &ogen.AdminGetRadarrQualityProfilesServiceUnavailable{
			Code:    503,
			Message: "Radarr integration not configured",
		}, nil
	}

	profiles, err := h.radarrService.GetQualityProfiles(ctx)
	if err != nil {
		h.logger.Error("Failed to get quality profiles from Radarr", slog.Any("error",err))
		return &ogen.AdminGetRadarrQualityProfilesServiceUnavailable{
			Code:    503,
			Message: "Failed to connect to Radarr",
		}, nil
	}

	ogenProfiles := make([]ogen.RadarrQualityProfile, 0, len(profiles))
	for _, p := range profiles {
		ogenProfiles = append(ogenProfiles, ogen.RadarrQualityProfile{
			ID:             p.ID,
			Name:           p.Name,
			UpgradeAllowed: ogen.NewOptBool(p.UpgradeAllowed),
			Cutoff:         ogen.NewOptInt(p.Cutoff),
			MinFormatScore: ogen.NewOptInt(p.MinFormatScore),
		})
	}

	h.logger.Info("Retrieved quality profiles from Radarr", slog.Int("count", len(ogenProfiles)))
	return &ogen.RadarrQualityProfileList{
		Profiles: ogenProfiles,
	}, nil
}

// AdminGetRadarrRootFolders returns all root folders from Radarr.
// GET /api/v1/admin/integrations/radarr/root-folders
func (h *Handler) AdminGetRadarrRootFolders(ctx context.Context) (ogen.AdminGetRadarrRootFoldersRes, error) {
	h.logger.Debug("AdminGetRadarrRootFolders called")

	// Check admin authorization
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminGetRadarrRootFoldersUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminGetRadarrRootFoldersForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	// Check if Radarr integration is configured
	if h.radarrService == nil {
		return &ogen.AdminGetRadarrRootFoldersServiceUnavailable{
			Code:    503,
			Message: "Radarr integration not configured",
		}, nil
	}

	folders, err := h.radarrService.GetRootFolders(ctx)
	if err != nil {
		h.logger.Error("Failed to get root folders from Radarr", slog.Any("error",err))
		return &ogen.AdminGetRadarrRootFoldersServiceUnavailable{
			Code:    503,
			Message: "Failed to connect to Radarr",
		}, nil
	}

	ogenFolders := make([]ogen.RadarrRootFolder, 0, len(folders))
	for _, f := range folders {
		ogenFolders = append(ogenFolders, ogen.RadarrRootFolder{
			ID:         f.ID,
			Path:       f.Path,
			Accessible: f.Accessible,
			FreeSpace:  ogen.NewOptInt64(f.FreeSpace),
		})
	}

	h.logger.Info("Retrieved root folders from Radarr", slog.Int("count", len(ogenFolders)))
	return &ogen.RadarrRootFolderList{
		Folders: ogenFolders,
	}, nil
}

// HandleRadarrWebhook handles incoming webhook events from Radarr.
// POST /api/v1/webhooks/radarr
func (h *Handler) HandleRadarrWebhook(ctx context.Context, req *ogen.RadarrWebhookPayload) (ogen.HandleRadarrWebhookRes, error) {
	h.logger.Debug("HandleRadarrWebhook called", slog.String("event_type", string(req.EventType)))

	// Convert ogen types to internal types
	payload := convertWebhookPayload(req)

	// If we have a River client, queue the job for async processing
	if h.riverClient != nil {
		_, err := h.riverClient.Insert(ctx, &radarr.RadarrWebhookJobArgs{
			Payload: *payload,
		}, nil)
		if err != nil {
			h.logger.Error("Failed to queue webhook job", slog.Any("error",err))
			return &ogen.Error{
				Code:    400,
				Message: "Failed to process webhook",
			}, nil
		}

		h.logger.Info("Radarr webhook queued", slog.String("event", string(req.EventType)))
		return &ogen.HandleRadarrWebhookAccepted{}, nil
	}

	// No River client and no direct handler - just acknowledge
	h.logger.Warn("Radarr webhook received but no handler configured")
	return &ogen.HandleRadarrWebhookAccepted{}, nil
}

// convertWebhookPayload converts ogen webhook payload to internal radarr type.
func convertWebhookPayload(req *ogen.RadarrWebhookPayload) *radarr.WebhookPayload {
	if req == nil {
		return nil
	}

	payload := &radarr.WebhookPayload{
		EventType:       string(req.EventType),
		InstanceName:    req.InstanceName.Value,
		ApplicationURL:  req.ApplicationUrl.Value,
		DownloadClient:  req.DownloadClient.Value,
		DownloadID:      req.DownloadId.Value,
		IsUpgrade:       req.IsUpgrade.Value,
	}

	// Convert movie
	if movie := req.Movie; movie.Set {
		payload.Movie = &radarr.WebhookMovie{
			ID:         int(movie.Value.ID.Value),
			Title:      movie.Value.Title.Value,
			Year:       int(movie.Value.Year.Value),
			TMDbID:     int(movie.Value.TmdbId.Value),
			IMDbID:     movie.Value.ImdbId.Value,
			FolderPath: movie.Value.FolderPath.Value,
		}
	}

	// Convert movie file
	if file := req.MovieFile; file.Set {
		payload.MovieFile = &radarr.WebhookMovieFile{
			ID:           int(file.Value.ID.Value),
			RelativePath: file.Value.RelativePath.Value,
			Path:         file.Value.Path.Value,
			Quality:      file.Value.Quality.Value,
			Size:         file.Value.Size.Value,
		}
	}

	// Convert release
	if release := req.Release; release.Set {
		payload.Release = &radarr.WebhookRelease{
			Quality:      release.Value.Quality.Value,
			ReleaseGroup: release.Value.ReleaseGroup.Value,
			ReleaseTitle: release.Value.ReleaseTitle.Value,
			Indexer:      release.Value.Indexer.Value,
			Size:         release.Value.Size.Value,
		}
	}

	return payload
}
