package api

import (
	"context"
	"time"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/integration/sonarr"
	"go.uber.org/zap"
)

// sonarrService is an optional dependency for Sonarr integration.
// When nil, all Sonarr endpoints return 503 Service Unavailable.
type sonarrService interface {
	GetStatus() sonarr.SyncStatus
	IsHealthy(ctx context.Context) bool
	GetSystemStatus(ctx context.Context) (*sonarr.SystemStatus, error)
	GetQualityProfiles(ctx context.Context) ([]sonarr.QualityProfile, error)
	GetRootFolders(ctx context.Context) ([]sonarr.RootFolder, error)
	SyncLibrary(ctx context.Context) (*sonarr.SyncResult, error)
}

// AdminGetSonarrStatus returns the current Sonarr integration status.
// GET /api/v1/admin/integrations/sonarr/status
func (h *Handler) AdminGetSonarrStatus(ctx context.Context) (ogen.AdminGetSonarrStatusRes, error) {
	h.logger.Debug("AdminGetSonarrStatus called")

	// Check admin authorization
	if !h.isAdmin(ctx) {
		return &ogen.AdminGetSonarrStatusForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Check if Sonarr integration is configured
	if h.sonarrService == nil {
		return &ogen.AdminGetSonarrStatusServiceUnavailable{
			Code:    503,
			Message: "Sonarr integration not configured",
		}, nil
	}

	// Check connection health
	connected := h.sonarrService.IsHealthy(ctx)

	// Get sync status
	syncStatus := h.sonarrService.GetStatus()

	response := &ogen.SonarrStatus{
		Connected: connected,
		SyncStatus: ogen.SonarrSyncStatus{
			IsRunning:     syncStatus.IsRunning,
			SeriesAdded:   syncStatus.SeriesAdded,
			SeriesUpdated: syncStatus.SeriesUpdated,
			SeriesRemoved: syncStatus.SeriesRemoved,
			TotalSeries:   syncStatus.TotalSeries,
		},
	}

	// Add optional episode counts
	if syncStatus.EpisodesAdded > 0 {
		response.SyncStatus.EpisodesAdded = ogen.NewOptInt(syncStatus.EpisodesAdded)
	}
	if syncStatus.EpisodesUpdated > 0 {
		response.SyncStatus.EpisodesUpdated = ogen.NewOptInt(syncStatus.EpisodesUpdated)
	}

	// Add last sync time if available
	if !syncStatus.LastSync.IsZero() {
		response.SyncStatus.LastSync = ogen.NewOptDateTime(syncStatus.LastSync)
	}

	// Add last sync error if present
	if syncStatus.LastSyncError != "" {
		response.SyncStatus.LastSyncError = ogen.NewOptString(syncStatus.LastSyncError)
	}

	// Get system status from Sonarr for version info
	if connected {
		if status, err := h.sonarrService.GetSystemStatus(ctx); err == nil {
			response.Version = ogen.NewOptString(status.Version)
			response.InstanceName = ogen.NewOptString(status.InstanceName)
			if status.StartTime != "" {
				if t, err := time.Parse(time.RFC3339, status.StartTime); err == nil {
					response.StartTime = ogen.NewOptDateTime(t)
				}
			}
		}
	}

	h.logger.Info("Sonarr status retrieved",
		zap.Bool("connected", connected),
		zap.Bool("sync_running", syncStatus.IsRunning))

	return response, nil
}

// AdminTriggerSonarrSync triggers a full library sync from Sonarr.
// POST /api/v1/admin/integrations/sonarr/sync
func (h *Handler) AdminTriggerSonarrSync(ctx context.Context) (ogen.AdminTriggerSonarrSyncRes, error) {
	h.logger.Debug("AdminTriggerSonarrSync called")

	// Check admin authorization
	if !h.isAdmin(ctx) {
		return &ogen.AdminTriggerSonarrSyncForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Check if Sonarr integration is configured
	if h.sonarrService == nil {
		return &ogen.AdminTriggerSonarrSyncServiceUnavailable{
			Code:    503,
			Message: "Sonarr integration not configured",
		}, nil
	}

	// Check if sync is already in progress
	status := h.sonarrService.GetStatus()
	if status.IsRunning {
		return &ogen.AdminTriggerSonarrSyncConflict{
			Code:    409,
			Message: "Sync already in progress",
		}, nil
	}

	// If we have a River client, queue the job
	if h.riverClient != nil {
		_, err := h.riverClient.Insert(ctx, &sonarr.SonarrSyncJobArgs{
			Operation: sonarr.SonarrSyncOperationFull,
		}, nil)
		if err != nil {
			h.logger.Error("Failed to queue Sonarr sync job", zap.Error(err))
			return &ogen.AdminTriggerSonarrSyncServiceUnavailable{
				Code:    503,
				Message: "Failed to queue sync job",
			}, nil
		}

		h.logger.Info("Sonarr sync job queued")
		return &ogen.SonarrSyncResponse{
			Message: "Sync job queued",
			Status:  ogen.SonarrSyncResponseStatusQueued,
		}, nil
	}

	// No River client, start sync directly (blocking)
	go func() {
		// Use a new context with timeout since the request context will be done
		syncCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if _, err := h.sonarrService.SyncLibrary(syncCtx); err != nil {
			h.logger.Error("Sonarr sync failed", zap.Error(err))
		}
	}()

	h.logger.Info("Sonarr sync started directly")
	return &ogen.SonarrSyncResponse{
		Message: "Sync started",
		Status:  ogen.SonarrSyncResponseStatusStarted,
	}, nil
}

// AdminGetSonarrQualityProfiles returns all quality profiles from Sonarr.
// GET /api/v1/admin/integrations/sonarr/quality-profiles
func (h *Handler) AdminGetSonarrQualityProfiles(ctx context.Context) (ogen.AdminGetSonarrQualityProfilesRes, error) {
	h.logger.Debug("AdminGetSonarrQualityProfiles called")

	// Check admin authorization
	if !h.isAdmin(ctx) {
		return &ogen.AdminGetSonarrQualityProfilesForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Check if Sonarr integration is configured
	if h.sonarrService == nil {
		return &ogen.AdminGetSonarrQualityProfilesServiceUnavailable{
			Code:    503,
			Message: "Sonarr integration not configured",
		}, nil
	}

	profiles, err := h.sonarrService.GetQualityProfiles(ctx)
	if err != nil {
		h.logger.Error("Failed to get quality profiles from Sonarr", zap.Error(err))
		return &ogen.AdminGetSonarrQualityProfilesServiceUnavailable{
			Code:    503,
			Message: "Failed to connect to Sonarr",
		}, nil
	}

	ogenProfiles := make([]ogen.SonarrQualityProfile, 0, len(profiles))
	for _, p := range profiles {
		ogenProfiles = append(ogenProfiles, ogen.SonarrQualityProfile{
			ID:             p.ID,
			Name:           p.Name,
			UpgradeAllowed: ogen.NewOptBool(p.UpgradeAllowed),
			Cutoff:         ogen.NewOptInt(p.Cutoff),
			MinFormatScore: ogen.NewOptInt(p.MinFormatScore),
		})
	}

	h.logger.Info("Retrieved quality profiles from Sonarr", zap.Int("count", len(ogenProfiles)))
	return &ogen.SonarrQualityProfileList{
		Profiles: ogenProfiles,
	}, nil
}

// AdminGetSonarrRootFolders returns all root folders from Sonarr.
// GET /api/v1/admin/integrations/sonarr/root-folders
func (h *Handler) AdminGetSonarrRootFolders(ctx context.Context) (ogen.AdminGetSonarrRootFoldersRes, error) {
	h.logger.Debug("AdminGetSonarrRootFolders called")

	// Check admin authorization
	if !h.isAdmin(ctx) {
		return &ogen.AdminGetSonarrRootFoldersForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Check if Sonarr integration is configured
	if h.sonarrService == nil {
		return &ogen.AdminGetSonarrRootFoldersServiceUnavailable{
			Code:    503,
			Message: "Sonarr integration not configured",
		}, nil
	}

	folders, err := h.sonarrService.GetRootFolders(ctx)
	if err != nil {
		h.logger.Error("Failed to get root folders from Sonarr", zap.Error(err))
		return &ogen.AdminGetSonarrRootFoldersServiceUnavailable{
			Code:    503,
			Message: "Failed to connect to Sonarr",
		}, nil
	}

	ogenFolders := make([]ogen.SonarrRootFolder, 0, len(folders))
	for _, f := range folders {
		ogenFolders = append(ogenFolders, ogen.SonarrRootFolder{
			ID:         f.ID,
			Path:       f.Path,
			Accessible: f.Accessible,
			FreeSpace:  ogen.NewOptInt64(f.FreeSpace),
		})
	}

	h.logger.Info("Retrieved root folders from Sonarr", zap.Int("count", len(ogenFolders)))
	return &ogen.SonarrRootFolderList{
		Folders: ogenFolders,
	}, nil
}

// HandleSonarrWebhook handles incoming webhook events from Sonarr.
// POST /api/v1/webhooks/sonarr
func (h *Handler) HandleSonarrWebhook(ctx context.Context, req *ogen.SonarrWebhookPayload) (ogen.HandleSonarrWebhookRes, error) {
	h.logger.Debug("HandleSonarrWebhook called", zap.String("event_type", string(req.EventType)))

	// Convert ogen types to internal types
	payload := convertSonarrWebhookPayload(req)

	// If we have a River client, queue the job for async processing
	if h.riverClient != nil {
		_, err := h.riverClient.Insert(ctx, &sonarr.SonarrWebhookJobArgs{
			Payload: *payload,
		}, nil)
		if err != nil {
			h.logger.Error("Failed to queue Sonarr webhook job", zap.Error(err))
			return &ogen.Error{
				Code:    400,
				Message: "Failed to process webhook",
			}, nil
		}

		h.logger.Info("Sonarr webhook queued", zap.String("event", string(req.EventType)))
		return &ogen.HandleSonarrWebhookAccepted{}, nil
	}

	// No River client and no direct handler - just acknowledge
	h.logger.Warn("Sonarr webhook received but no handler configured")
	return &ogen.HandleSonarrWebhookAccepted{}, nil
}

// convertSonarrWebhookPayload converts ogen webhook payload to internal sonarr type.
func convertSonarrWebhookPayload(req *ogen.SonarrWebhookPayload) *sonarr.WebhookPayload {
	if req == nil {
		return nil
	}

	payload := &sonarr.WebhookPayload{
		EventType:          string(req.EventType),
		InstanceName:       req.InstanceName.Value,
		ApplicationURL:     req.ApplicationUrl.Value,
		DownloadClient:     req.DownloadClient.Value,
		DownloadClientType: req.DownloadClientType.Value,
		DownloadID:         req.DownloadId.Value,
		IsUpgrade:          req.IsUpgrade.Value,
	}

	// Convert series
	if series := req.Series; series.Set {
		payload.Series = &sonarr.WebhookSeries{
			ID:        series.Value.ID.Value,
			Title:     series.Value.Title.Value,
			TitleSlug: series.Value.TitleSlug.Value,
			Path:      series.Value.Path.Value,
			TVDbID:    series.Value.TvdbId.Value,
			TVMazeID:  series.Value.TvMazeId.Value,
			IMDbID:    series.Value.ImdbId.Value,
			Type:      series.Value.Type.Value,
		}
	}

	// Convert episodes
	if len(req.Episodes) > 0 {
		payload.Episodes = make([]sonarr.WebhookEpisode, 0, len(req.Episodes))
		for _, ep := range req.Episodes {
			episode := sonarr.WebhookEpisode{
				ID:            ep.ID.Value,
				EpisodeNumber: ep.EpisodeNumber.Value,
				SeasonNumber:  ep.SeasonNumber.Value,
				Title:         ep.Title.Value,
			}
			if ep.AirDate.Set {
				episode.AirDate = ep.AirDate.Value.String()
			}
			if ep.AirDateUtc.Set {
				episode.AirDateUtc = ep.AirDateUtc.Value.Format(time.RFC3339)
			}
			payload.Episodes = append(payload.Episodes, episode)
		}
	}

	// Convert episode file
	if file := req.EpisodeFile; file.Set {
		payload.EpisodeFile = &sonarr.WebhookEpisodeFile{
			ID:             file.Value.ID.Value,
			RelativePath:   file.Value.RelativePath.Value,
			Path:           file.Value.Path.Value,
			Quality:        file.Value.Quality.Value,
			QualityVersion: file.Value.QualityVersion.Value,
			Size:           file.Value.Size.Value,
		}
		if file.Value.DateAdded.Set {
			dateAdded := file.Value.DateAdded.Value
			payload.EpisodeFile.DateAdded = &dateAdded
		}
	}

	// Convert deleted files
	if len(req.DeletedFiles) > 0 {
		payload.DeletedFiles = make([]sonarr.WebhookEpisodeFile, 0, len(req.DeletedFiles))
		for _, f := range req.DeletedFiles {
			file := sonarr.WebhookEpisodeFile{
				ID:             f.ID.Value,
				RelativePath:   f.RelativePath.Value,
				Path:           f.Path.Value,
				Quality:        f.Quality.Value,
				QualityVersion: f.QualityVersion.Value,
				Size:           f.Size.Value,
			}
			if f.DateAdded.Set {
				dateAdded := f.DateAdded.Value
				file.DateAdded = &dateAdded
			}
			payload.DeletedFiles = append(payload.DeletedFiles, file)
		}
	}

	// Convert release
	if release := req.Release; release.Set {
		payload.Release = &sonarr.WebhookRelease{
			Quality:      release.Value.Quality.Value,
			ReleaseGroup: release.Value.ReleaseGroup.Value,
			ReleaseTitle: release.Value.ReleaseTitle.Value,
			Indexer:      release.Value.Indexer.Value,
			Size:         release.Value.Size.Value,
		}
	}

	return payload
}
