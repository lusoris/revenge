package sonarr

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/util"
)

// SyncService handles synchronization between Sonarr and Revenge.
// It implements the PRIMARY provider pattern from the metadata priority chain.
type SyncService struct {
	client      *Client
	mapper      *Mapper
	tvshowRepo  tvshow.Repository
	logger      *slog.Logger
	syncMu      sync.Mutex
	syncStatus  SyncStatus
}

// SyncStatus represents the current sync status.
type SyncStatus struct {
	IsRunning       bool      `json:"is_running"`
	LastSync        time.Time `json:"last_sync,omitempty"`
	LastSyncError   string    `json:"last_sync_error,omitempty"`
	SeriesAdded     int       `json:"series_added"`
	SeriesUpdated   int       `json:"series_updated"`
	SeriesRemoved   int       `json:"series_removed"`
	EpisodesAdded   int       `json:"episodes_added"`
	EpisodesUpdated int       `json:"episodes_updated"`
	TotalSeries     int       `json:"total_series"`
}

// SyncResult contains the result of a sync operation.
type SyncResult struct {
	SeriesAdded     int           `json:"series_added"`
	SeriesUpdated   int           `json:"series_updated"`
	SeriesRemoved   int           `json:"series_removed"`
	SeriesSkipped   int           `json:"series_skipped"`
	EpisodesAdded   int           `json:"episodes_added"`
	EpisodesUpdated int           `json:"episodes_updated"`
	Errors          []string      `json:"errors,omitempty"`
	Duration        time.Duration `json:"duration"`
}

// NewSyncService creates a new Sonarr sync service.
func NewSyncService(
	client *Client,
	mapper *Mapper,
	tvshowRepo tvshow.Repository,
	logger *slog.Logger,
) *SyncService {
	return &SyncService{
		client:     client,
		mapper:     mapper,
		tvshowRepo: tvshowRepo,
		logger:     logger.With("service", "sonarr_sync"),
	}
}

// GetStatus returns the current sync status.
func (s *SyncService) GetStatus() SyncStatus {
	s.syncMu.Lock()
	defer s.syncMu.Unlock()
	return s.syncStatus
}

// IsHealthy checks if Sonarr is reachable and healthy.
func (s *SyncService) IsHealthy(ctx context.Context) bool {
	return s.client.IsHealthy(ctx)
}

// GetSystemStatus returns Sonarr's system status.
func (s *SyncService) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	return s.client.GetSystemStatus(ctx)
}

// GetQualityProfiles returns all quality profiles from Sonarr.
func (s *SyncService) GetQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	return s.client.GetQualityProfiles(ctx)
}

// GetRootFolders returns all root folders from Sonarr.
func (s *SyncService) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	return s.client.GetRootFolders(ctx)
}

// LookupSeries searches for TV series via Sonarr's lookup API.
func (s *SyncService) LookupSeries(ctx context.Context, term string) ([]Series, error) {
	return s.client.LookupSeries(ctx, term)
}

// LookupSeriesByTVDbID looks up a series by TVDb ID via Sonarr.
func (s *SyncService) LookupSeriesByTVDbID(ctx context.Context, tvdbID int) (*Series, error) {
	return s.client.LookupSeriesByTVDbID(ctx, tvdbID)
}

// SyncLibrary performs a full library sync from Sonarr to Revenge.
func (s *SyncService) SyncLibrary(ctx context.Context) (*SyncResult, error) {
	s.syncMu.Lock()
	if s.syncStatus.IsRunning {
		s.syncMu.Unlock()
		return nil, fmt.Errorf("sync already in progress")
	}
	s.syncStatus.IsRunning = true
	s.syncMu.Unlock()

	defer func() {
		s.syncMu.Lock()
		s.syncStatus.IsRunning = false
		s.syncStatus.LastSync = time.Now()
		s.syncMu.Unlock()
	}()

	start := time.Now()
	s.logger.Info("starting library sync from Sonarr")

	result := &SyncResult{}

	// Get all series from Sonarr
	sonarrSeries, err := s.client.GetAllSeries(ctx)
	if err != nil {
		s.syncMu.Lock()
		s.syncStatus.LastSyncError = err.Error()
		s.syncMu.Unlock()
		return nil, fmt.Errorf("failed to get series from Sonarr: %w", err)
	}

	s.logger.Info("fetched series from Sonarr", "count", len(sonarrSeries))

	// Get existing series from Revenge by SonarrID
	existingSeries, err := s.getExistingSonarrSeries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing series: %w", err)
	}

	// Track which series we've seen for removal detection
	seenSonarrIDs := make(map[int]bool)

	// Process each series
	for _, ss := range sonarrSeries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		seenSonarrIDs[ss.ID] = true

		// Only sync series that have files
		if ss.Statistics == nil || ss.Statistics.EpisodeFileCount == 0 {
			result.SeriesSkipped++
			continue
		}

		if existingID, exists := existingSeries[ss.ID]; exists {
			// Update existing series
			if err := s.updateSeries(ctx, ss, existingID, result); err != nil {
				s.logger.Error("failed to update series", "sonarr_id", ss.ID, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("update series %d: %v", ss.ID, err))
				continue
			}
			result.SeriesUpdated++
		} else {
			// Add new series
			if err := s.addSeries(ctx, ss, result); err != nil {
				s.logger.Error("failed to add series", "sonarr_id", ss.ID, "error", err)
				result.Errors = append(result.Errors, fmt.Sprintf("add series %d: %v", ss.ID, err))
				continue
			}
			result.SeriesAdded++
		}
	}

	// Find and mark removed series
	for sonarrID := range existingSeries {
		if !seenSonarrIDs[sonarrID] {
			s.logger.Info("series no longer in Sonarr", "sonarr_id", sonarrID)
			result.SeriesRemoved++
		}
	}

	result.Duration = time.Since(start)

	s.syncMu.Lock()
	s.syncStatus.SeriesAdded = result.SeriesAdded
	s.syncStatus.SeriesUpdated = result.SeriesUpdated
	s.syncStatus.SeriesRemoved = result.SeriesRemoved
	s.syncStatus.EpisodesAdded = result.EpisodesAdded
	s.syncStatus.EpisodesUpdated = result.EpisodesUpdated
	s.syncStatus.TotalSeries = len(sonarrSeries)
	s.syncStatus.LastSyncError = ""
	s.syncMu.Unlock()

	s.logger.Info("library sync completed",
		"series_added", result.SeriesAdded,
		"series_updated", result.SeriesUpdated,
		"series_removed", result.SeriesRemoved,
		"series_skipped", result.SeriesSkipped,
		"episodes_added", result.EpisodesAdded,
		"episodes_updated", result.EpisodesUpdated,
		"duration", result.Duration,
	)

	return result, nil
}

// SyncSeries syncs a single series from Sonarr by its Sonarr ID.
func (s *SyncService) SyncSeries(ctx context.Context, sonarrID int) error {
	ss, err := s.client.GetSeries(ctx, sonarrID)
	if err != nil {
		return fmt.Errorf("failed to get series from Sonarr: %w", err)
	}

	// Get existing series by SonarrID
	existing, err := s.tvshowRepo.GetSeriesBySonarrID(ctx, util.SafeIntToInt32(sonarrID))
	if err != nil && err.Error() != "series not found" {
		return fmt.Errorf("failed to check existing series: %w", err)
	}

	result := &SyncResult{}
	if existing != nil {
		return s.updateSeries(ctx, *ss, existing.ID, result)
	}
	return s.addSeries(ctx, *ss, result)
}

// RefreshSeries triggers a refresh in Sonarr for a series by TVDb ID.
func (s *SyncService) RefreshSeries(ctx context.Context, tvdbID int) (*Command, error) {
	ss, err := s.client.GetSeriesByTVDbID(ctx, tvdbID)
	if err != nil {
		return nil, fmt.Errorf("series not found in Sonarr: %w", err)
	}
	return s.client.RefreshSeries(ctx, ss.ID)
}

// getExistingSonarrSeries returns a map of SonarrID -> SeriesID for all series with a SonarrID.
func (s *SyncService) getExistingSonarrSeries(ctx context.Context) (map[int]uuid.UUID, error) {
	series, err := s.tvshowRepo.ListSeries(ctx, tvshow.SeriesListFilters{
		Limit:  10000,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[int]uuid.UUID)
	for _, sr := range series {
		if sr.SonarrID != nil {
			result[int(*sr.SonarrID)] = sr.ID
		}
	}
	return result, nil
}

// addSeries creates a new series from Sonarr data.
func (s *SyncService) addSeries(ctx context.Context, ss Series, result *SyncResult) error {
	series := s.mapper.ToSeries(&ss)

	// Create the series
	createParams := s.seriesToCreateParams(series)
	created, err := s.tvshowRepo.CreateSeries(ctx, createParams)
	if err != nil {
		return fmt.Errorf("failed to create series: %w", err)
	}

	// Sync seasons and episodes
	if err := s.syncSeasonsAndEpisodes(ctx, ss.ID, created.ID, result); err != nil {
		s.logger.Warn("failed to sync seasons/episodes", "series_id", created.ID, "error", err)
	}

	s.logger.Debug("added series from Sonarr", "id", created.ID, "title", ss.Title)
	return nil
}

// updateSeries updates an existing series with Sonarr data.
func (s *SyncService) updateSeries(ctx context.Context, ss Series, existingID uuid.UUID, result *SyncResult) error {
	series := s.mapper.ToSeries(&ss)
	series.ID = existingID

	updateParams := s.seriesToUpdateParams(series)
	if _, err := s.tvshowRepo.UpdateSeries(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update series: %w", err)
	}

	// Sync seasons and episodes
	if err := s.syncSeasonsAndEpisodes(ctx, ss.ID, existingID, result); err != nil {
		s.logger.Warn("failed to sync seasons/episodes", "series_id", existingID, "error", err)
	}

	s.logger.Debug("updated series from Sonarr", "id", existingID, "title", ss.Title)
	return nil
}

// syncSeasonsAndEpisodes syncs seasons and episodes for a series.
func (s *SyncService) syncSeasonsAndEpisodes(ctx context.Context, sonarrSeriesID int, seriesID uuid.UUID, result *SyncResult) error {
	// Get all episodes from Sonarr
	sonarrEpisodes, err := s.client.GetEpisodes(ctx, sonarrSeriesID)
	if err != nil {
		return fmt.Errorf("failed to get episodes: %w", err)
	}

	// Group episodes by season
	episodesBySeason := make(map[int][]Episode)
	for _, ep := range sonarrEpisodes {
		episodesBySeason[ep.SeasonNumber] = append(episodesBySeason[ep.SeasonNumber], ep)
	}

	// Get existing seasons
	existingSeasons, err := s.tvshowRepo.ListSeasonsBySeries(ctx, seriesID)
	if err != nil {
		return fmt.Errorf("failed to get existing seasons: %w", err)
	}

	seasonMap := make(map[int32]uuid.UUID)
	for _, season := range existingSeasons {
		seasonMap[season.SeasonNumber] = season.ID
	}

	// Process each season
	for seasonNum, episodes := range episodesBySeason {
		var seasonID uuid.UUID

		if existingSeasonID, exists := seasonMap[util.SafeIntToInt32(seasonNum)]; exists {
			seasonID = existingSeasonID
		} else {
			// Create new season
			seasonParams := tvshow.CreateSeasonParams{
				SeriesID:     seriesID,
				SeasonNumber: util.SafeIntToInt32(seasonNum),
				Name:         fmt.Sprintf("Season %d", seasonNum),
				EpisodeCount: util.SafeIntToInt32(len(episodes)),
			}
			if seasonNum == 0 {
				seasonParams.Name = "Specials"
			}

			created, err := s.tvshowRepo.CreateSeason(ctx, seasonParams)
			if err != nil {
				s.logger.Warn("failed to create season", "series_id", seriesID, "season", seasonNum, "error", err)
				continue
			}
			seasonID = created.ID
		}

		// Sync episodes for this season
		if err := s.syncEpisodes(ctx, episodes, seriesID, seasonID, result); err != nil {
			s.logger.Warn("failed to sync episodes", "series_id", seriesID, "season", seasonNum, "error", err)
		}
	}

	return nil
}

// syncEpisodes syncs episodes for a season.
func (s *SyncService) syncEpisodes(ctx context.Context, sonarrEpisodes []Episode, seriesID, seasonID uuid.UUID, result *SyncResult) error {
	// Get existing episodes for this season
	existingEpisodes, err := s.tvshowRepo.ListEpisodesBySeason(ctx, seasonID)
	if err != nil {
		return fmt.Errorf("failed to get existing episodes: %w", err)
	}

	episodeMap := make(map[int32]uuid.UUID)
	for _, ep := range existingEpisodes {
		episodeMap[ep.EpisodeNumber] = ep.ID
	}

	for _, se := range sonarrEpisodes {
		// Only sync episodes that have files
		if !se.HasFile {
			continue
		}

		episode := s.mapper.ToEpisode(&se, seriesID, seasonID)

		if existingEpisodeID, exists := episodeMap[util.SafeIntToInt32(se.EpisodeNumber)]; exists {
			// Update existing episode
			episode.ID = existingEpisodeID
			updateParams := s.episodeToUpdateParams(episode)
			if _, err := s.tvshowRepo.UpdateEpisode(ctx, updateParams); err != nil {
				s.logger.Warn("failed to update episode", "episode_id", existingEpisodeID, "error", err)
				continue
			}
			result.EpisodesUpdated++

			// Sync episode files
			if se.EpisodeFile != nil {
				if err := s.syncEpisodeFile(ctx, se.EpisodeFile, existingEpisodeID); err != nil {
					s.logger.Warn("failed to sync episode file", "episode_id", existingEpisodeID, "error", err)
				}
			}
		} else {
			// Create new episode
			createParams := s.episodeToCreateParams(episode)
			created, err := s.tvshowRepo.CreateEpisode(ctx, createParams)
			if err != nil {
				s.logger.Warn("failed to create episode", "error", err)
				continue
			}
			result.EpisodesAdded++

			// Sync episode files
			if se.EpisodeFile != nil {
				if err := s.syncEpisodeFile(ctx, se.EpisodeFile, created.ID); err != nil {
					s.logger.Warn("failed to sync episode file", "episode_id", created.ID, "error", err)
				}
			}
		}
	}

	return nil
}

// syncEpisodeFile syncs an episode file.
func (s *SyncService) syncEpisodeFile(ctx context.Context, sef *EpisodeFile, episodeID uuid.UUID) error {
	file := s.mapper.ToEpisodeFile(sef, episodeID)

	// Check if file already exists by SonarrFileID
	existing, err := s.tvshowRepo.GetEpisodeFileBySonarrID(ctx, util.SafeIntToInt32(sef.ID))
	if err != nil && err.Error() != "episode file not found" {
		return err
	}

	if existing != nil {
		// Update existing file
		updateParams := s.episodeFileToUpdateParams(file, existing.ID)
		if _, err := s.tvshowRepo.UpdateEpisodeFile(ctx, updateParams); err != nil {
			return fmt.Errorf("failed to update episode file: %w", err)
		}
	} else {
		// Create new file
		createParams := s.episodeFileToCreateParams(file)
		if _, err := s.tvshowRepo.CreateEpisodeFile(ctx, createParams); err != nil {
			return fmt.Errorf("failed to create episode file: %w", err)
		}
	}

	return nil
}

// Conversion helpers

func (s *SyncService) seriesToCreateParams(sr *tvshow.Series) tvshow.CreateSeriesParams {
	params := tvshow.CreateSeriesParams{
		TVDbID:           sr.TVDbID,
		IMDbID:           sr.IMDbID,
		SonarrID:         sr.SonarrID,
		Title:            sr.Title,
		OriginalTitle:    sr.OriginalTitle,
		OriginalLanguage: sr.OriginalLanguage,
		Overview:         sr.Overview,
		Status:           sr.Status,
		Type:             sr.Type,
		TotalSeasons:     sr.TotalSeasons,
		TotalEpisodes:    sr.TotalEpisodes,
		PosterPath:       sr.PosterPath,
		BackdropPath:     sr.BackdropPath,
	}
	if sr.FirstAirDate != nil {
		d := sr.FirstAirDate.Format("2006-01-02")
		params.FirstAirDate = &d
	}
	if sr.LastAirDate != nil {
		d := sr.LastAirDate.Format("2006-01-02")
		params.LastAirDate = &d
	}
	return params
}

func (s *SyncService) seriesToUpdateParams(sr *tvshow.Series) tvshow.UpdateSeriesParams {
	params := tvshow.UpdateSeriesParams{
		ID:               sr.ID,
		TVDbID:           sr.TVDbID,
		IMDbID:           sr.IMDbID,
		SonarrID:         sr.SonarrID,
		Title:            &sr.Title,
		OriginalTitle:    sr.OriginalTitle,
		OriginalLanguage: &sr.OriginalLanguage,
		Overview:         sr.Overview,
		Status:           sr.Status,
		Type:             sr.Type,
		TotalSeasons:     &sr.TotalSeasons,
		TotalEpisodes:    &sr.TotalEpisodes,
		PosterPath:       sr.PosterPath,
		BackdropPath:     sr.BackdropPath,
	}
	if sr.FirstAirDate != nil {
		d := sr.FirstAirDate.Format("2006-01-02")
		params.FirstAirDate = &d
	}
	if sr.LastAirDate != nil {
		d := sr.LastAirDate.Format("2006-01-02")
		params.LastAirDate = &d
	}
	return params
}

func (s *SyncService) episodeToCreateParams(ep *tvshow.Episode) tvshow.CreateEpisodeParams {
	params := tvshow.CreateEpisodeParams{
		SeriesID:      ep.SeriesID,
		SeasonID:      ep.SeasonID,
		TVDbID:        ep.TVDbID,
		SeasonNumber:  ep.SeasonNumber,
		EpisodeNumber: ep.EpisodeNumber,
		Title:         ep.Title,
		Overview:      ep.Overview,
		Runtime:       ep.Runtime,
		StillPath:     ep.StillPath,
	}
	if ep.AirDate != nil {
		d := ep.AirDate.Format("2006-01-02")
		params.AirDate = &d
	}
	return params
}

func (s *SyncService) episodeToUpdateParams(ep *tvshow.Episode) tvshow.UpdateEpisodeParams {
	params := tvshow.UpdateEpisodeParams{
		ID:        ep.ID,
		TVDbID:    ep.TVDbID,
		Title:     &ep.Title,
		Overview:  ep.Overview,
		Runtime:   ep.Runtime,
		StillPath: ep.StillPath,
	}
	if ep.AirDate != nil {
		d := ep.AirDate.Format("2006-01-02")
		params.AirDate = &d
	}
	return params
}

func (s *SyncService) episodeFileToCreateParams(ef *tvshow.EpisodeFile) tvshow.CreateEpisodeFileParams {
	return tvshow.CreateEpisodeFileParams{
		EpisodeID:      ef.EpisodeID,
		FilePath:       ef.FilePath,
		FileName:       ef.FileName,
		FileSize:       ef.FileSize,
		Container:      ef.Container,
		Resolution:     ef.Resolution,
		QualityProfile: ef.QualityProfile,
		VideoCodec:     ef.VideoCodec,
		AudioCodec:     ef.AudioCodec,
		BitrateKbps:    ef.BitrateKbps,
		SonarrFileID:   ef.SonarrFileID,
	}
}

func (s *SyncService) episodeFileToUpdateParams(ef *tvshow.EpisodeFile, id uuid.UUID) tvshow.UpdateEpisodeFileParams {
	return tvshow.UpdateEpisodeFileParams{
		ID:             id,
		FilePath:       &ef.FilePath,
		FileName:       &ef.FileName,
		FileSize:       &ef.FileSize,
		Container:      ef.Container,
		Resolution:     ef.Resolution,
		QualityProfile: ef.QualityProfile,
		VideoCodec:     ef.VideoCodec,
		AudioCodec:     ef.AudioCodec,
		BitrateKbps:    ef.BitrateKbps,
		SonarrFileID:   ef.SonarrFileID,
	}
}
