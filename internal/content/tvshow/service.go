package tvshow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service defines business logic for TV shows
type Service interface {
	// Series operations
	GetSeries(ctx context.Context, id uuid.UUID) (*Series, error)
	GetSeriesByTMDbID(ctx context.Context, tmdbID int32) (*Series, error)
	GetSeriesByTVDbID(ctx context.Context, tvdbID int32) (*Series, error)
	GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*Series, error)
	ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error)
	CountSeries(ctx context.Context) (int64, error)
	SearchSeries(ctx context.Context, query string, limit, offset int32) ([]Series, error)
	ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Series, error)
	ListByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Series, error)
	ListByNetwork(ctx context.Context, networkID uuid.UUID, limit, offset int32) ([]Series, error)
	ListByStatus(ctx context.Context, status string, limit, offset int32) ([]Series, error)
	CreateSeries(ctx context.Context, params CreateSeriesParams) (*Series, error)
	UpdateSeries(ctx context.Context, params UpdateSeriesParams) (*Series, error)
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	// Season operations
	GetSeason(ctx context.Context, id uuid.UUID) (*Season, error)
	GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) (*Season, error)
	ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]Season, error)
	ListSeasonsWithEpisodeCount(ctx context.Context, seriesID uuid.UUID) ([]SeasonWithEpisodeCount, error)
	CreateSeason(ctx context.Context, params CreateSeasonParams) (*Season, error)
	UpsertSeason(ctx context.Context, params CreateSeasonParams) (*Season, error)
	UpdateSeason(ctx context.Context, params UpdateSeasonParams) (*Season, error)
	DeleteSeason(ctx context.Context, id uuid.UUID) error

	// Episode operations
	GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error)
	GetEpisodeByTMDbID(ctx context.Context, tmdbID int32) (*Episode, error)
	GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int32) (*Episode, error)
	GetEpisodeByFile(ctx context.Context, filePath string) (*Episode, error)
	ListEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) ([]Episode, error)
	ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error)
	ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) ([]Episode, error)
	ListRecentEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error)
	ListUpcomingEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error)
	CreateEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error)
	UpsertEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error)
	UpdateEpisode(ctx context.Context, params UpdateEpisodeParams) (*Episode, error)
	DeleteEpisode(ctx context.Context, id uuid.UUID) error

	// Episode files
	GetEpisodeFile(ctx context.Context, id uuid.UUID) (*EpisodeFile, error)
	GetEpisodeFileByPath(ctx context.Context, filePath string) (*EpisodeFile, error)
	GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*EpisodeFile, error)
	ListEpisodeFiles(ctx context.Context, episodeID uuid.UUID) ([]EpisodeFile, error)
	CreateEpisodeFile(ctx context.Context, params CreateEpisodeFileParams) (*EpisodeFile, error)
	UpdateEpisodeFile(ctx context.Context, params UpdateEpisodeFileParams) (*EpisodeFile, error)
	DeleteEpisodeFile(ctx context.Context, id uuid.UUID) error

	// Credits
	GetSeriesCast(ctx context.Context, seriesID uuid.UUID) ([]SeriesCredit, error)
	GetSeriesCrew(ctx context.Context, seriesID uuid.UUID) ([]SeriesCredit, error)
	GetEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error)
	GetEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error)

	// Genres & Networks
	GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error)
	GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]Network, error)

	// Watch progress
	UpdateEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID, progressSeconds, durationSeconds int32) (*EpisodeWatched, error)
	GetEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatched, error)
	MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) error
	MarkSeasonWatched(ctx context.Context, userID, seasonID uuid.UUID) error
	MarkSeriesWatched(ctx context.Context, userID, seriesID uuid.UUID) error
	RemoveEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) error
	RemoveSeriesProgress(ctx context.Context, userID, seriesID uuid.UUID) error
	GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error)
	GetNextEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error)
	GetSeriesWatchStats(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchStats, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserTVStats, error)

	// Metadata refresh
	RefreshSeriesMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error
	RefreshSeasonMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error
	RefreshEpisodeMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error
}

// tvService implements the Service interface
type tvService struct {
	repo             Repository
	metadataProvider MetadataProvider
}

// NewService creates a new TV show service
func NewService(repo Repository, metadataProvider MetadataProvider) Service {
	return &tvService{
		repo:             repo,
		metadataProvider: metadataProvider,
	}
}

// =============================================================================
// Series Operations
// =============================================================================

func (s *tvService) GetSeries(ctx context.Context, id uuid.UUID) (*Series, error) {
	return s.repo.GetSeries(ctx, id)
}

func (s *tvService) GetSeriesByTMDbID(ctx context.Context, tmdbID int32) (*Series, error) {
	return s.repo.GetSeriesByTMDbID(ctx, tmdbID)
}

func (s *tvService) GetSeriesByTVDbID(ctx context.Context, tvdbID int32) (*Series, error) {
	return s.repo.GetSeriesByTVDbID(ctx, tvdbID)
}

func (s *tvService) GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*Series, error) {
	return s.repo.GetSeriesBySonarrID(ctx, sonarrID)
}

func (s *tvService) ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error) {
	return s.repo.ListSeries(ctx, filters)
}

func (s *tvService) CountSeries(ctx context.Context) (int64, error) {
	return s.repo.CountSeries(ctx)
}

func (s *tvService) SearchSeries(ctx context.Context, query string, limit, offset int32) ([]Series, error) {
	return s.repo.SearchSeriesByTitle(ctx, query, limit, offset)
}

func (s *tvService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Series, error) {
	return s.repo.ListRecentlyAddedSeries(ctx, limit, offset)
}

func (s *tvService) ListByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Series, error) {
	return s.repo.ListSeriesByGenre(ctx, tmdbGenreID, limit, offset)
}

func (s *tvService) ListByNetwork(ctx context.Context, networkID uuid.UUID, limit, offset int32) ([]Series, error) {
	return s.repo.ListSeriesByNetwork(ctx, networkID, limit, offset)
}

func (s *tvService) ListByStatus(ctx context.Context, status string, limit, offset int32) ([]Series, error) {
	return s.repo.ListSeriesByStatus(ctx, status, limit, offset)
}

func (s *tvService) CreateSeries(ctx context.Context, params CreateSeriesParams) (*Series, error) {
	// Validate required fields
	if params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Check if series already exists by TMDb ID
	if params.TMDbID != nil {
		existing, err := s.repo.GetSeriesByTMDbID(ctx, *params.TMDbID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("series with TMDb ID %d already exists", *params.TMDbID)
		}
	}

	return s.repo.CreateSeries(ctx, params)
}

func (s *tvService) UpdateSeries(ctx context.Context, params UpdateSeriesParams) (*Series, error) {
	// Verify series exists
	_, err := s.repo.GetSeries(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("series not found: %w", err)
	}

	// Validate title if provided
	if params.Title != nil && *params.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	return s.repo.UpdateSeries(ctx, params)
}

func (s *tvService) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSeries(ctx, id)
}

// =============================================================================
// Season Operations
// =============================================================================

func (s *tvService) GetSeason(ctx context.Context, id uuid.UUID) (*Season, error) {
	return s.repo.GetSeason(ctx, id)
}

func (s *tvService) GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) (*Season, error) {
	return s.repo.GetSeasonByNumber(ctx, seriesID, seasonNumber)
}

func (s *tvService) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]Season, error) {
	return s.repo.ListSeasonsBySeries(ctx, seriesID)
}

func (s *tvService) ListSeasonsWithEpisodeCount(ctx context.Context, seriesID uuid.UUID) ([]SeasonWithEpisodeCount, error) {
	return s.repo.ListSeasonsBySeriesWithEpisodeCount(ctx, seriesID)
}

func (s *tvService) CreateSeason(ctx context.Context, params CreateSeasonParams) (*Season, error) {
	// Verify series exists
	_, err := s.repo.GetSeries(ctx, params.SeriesID)
	if err != nil {
		return nil, fmt.Errorf("series not found: %w", err)
	}

	// Check if season already exists
	existing, err := s.repo.GetSeasonByNumber(ctx, params.SeriesID, params.SeasonNumber)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("season %d already exists for this series", params.SeasonNumber)
	}

	return s.repo.CreateSeason(ctx, params)
}

func (s *tvService) UpsertSeason(ctx context.Context, params CreateSeasonParams) (*Season, error) {
	// Verify series exists
	_, err := s.repo.GetSeries(ctx, params.SeriesID)
	if err != nil {
		return nil, fmt.Errorf("series not found: %w", err)
	}

	return s.repo.UpsertSeason(ctx, params)
}

func (s *tvService) UpdateSeason(ctx context.Context, params UpdateSeasonParams) (*Season, error) {
	// Verify season exists
	_, err := s.repo.GetSeason(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("season not found: %w", err)
	}

	return s.repo.UpdateSeason(ctx, params)
}

func (s *tvService) DeleteSeason(ctx context.Context, id uuid.UUID) error {
	// Delete episodes first
	if err := s.repo.DeleteEpisodesBySeason(ctx, id); err != nil {
		return fmt.Errorf("failed to delete episodes: %w", err)
	}
	return s.repo.DeleteSeason(ctx, id)
}

// =============================================================================
// Episode Operations
// =============================================================================

func (s *tvService) GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error) {
	return s.repo.GetEpisode(ctx, id)
}

func (s *tvService) GetEpisodeByTMDbID(ctx context.Context, tmdbID int32) (*Episode, error) {
	return s.repo.GetEpisodeByTMDbID(ctx, tmdbID)
}

func (s *tvService) GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int32) (*Episode, error) {
	return s.repo.GetEpisodeByNumber(ctx, seriesID, seasonNumber, episodeNumber)
}

func (s *tvService) GetEpisodeByFile(ctx context.Context, filePath string) (*Episode, error) {
	file, err := s.repo.GetEpisodeFileByPath(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}
	return s.repo.GetEpisode(ctx, file.EpisodeID)
}

func (s *tvService) ListEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) ([]Episode, error) {
	return s.repo.ListEpisodesBySeries(ctx, seriesID)
}

func (s *tvService) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error) {
	return s.repo.ListEpisodesBySeason(ctx, seasonID)
}

func (s *tvService) ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) ([]Episode, error) {
	return s.repo.ListEpisodesBySeasonNumber(ctx, seriesID, seasonNumber)
}

func (s *tvService) ListRecentEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error) {
	return s.repo.ListRecentEpisodes(ctx, limit, offset)
}

func (s *tvService) ListUpcomingEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error) {
	return s.repo.ListUpcomingEpisodes(ctx, limit, offset)
}

func (s *tvService) CreateEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error) {
	// Verify season exists
	_, err := s.repo.GetSeason(ctx, params.SeasonID)
	if err != nil {
		return nil, fmt.Errorf("season not found: %w", err)
	}

	// Check if episode already exists
	existing, err := s.repo.GetEpisodeByNumber(ctx, params.SeriesID, params.SeasonNumber, params.EpisodeNumber)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("episode S%02dE%02d already exists for this series", params.SeasonNumber, params.EpisodeNumber)
	}

	return s.repo.CreateEpisode(ctx, params)
}

func (s *tvService) UpsertEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error) {
	// Verify season exists
	_, err := s.repo.GetSeason(ctx, params.SeasonID)
	if err != nil {
		return nil, fmt.Errorf("season not found: %w", err)
	}

	return s.repo.UpsertEpisode(ctx, params)
}

func (s *tvService) UpdateEpisode(ctx context.Context, params UpdateEpisodeParams) (*Episode, error) {
	// Verify episode exists
	_, err := s.repo.GetEpisode(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("episode not found: %w", err)
	}

	return s.repo.UpdateEpisode(ctx, params)
}

func (s *tvService) DeleteEpisode(ctx context.Context, id uuid.UUID) error {
	// Delete associated files first
	if err := s.repo.DeleteEpisodeFilesByEpisode(ctx, id); err != nil {
		return fmt.Errorf("failed to delete episode files: %w", err)
	}
	return s.repo.DeleteEpisode(ctx, id)
}

// =============================================================================
// Episode File Operations
// =============================================================================

func (s *tvService) GetEpisodeFile(ctx context.Context, id uuid.UUID) (*EpisodeFile, error) {
	return s.repo.GetEpisodeFile(ctx, id)
}

func (s *tvService) GetEpisodeFileByPath(ctx context.Context, filePath string) (*EpisodeFile, error) {
	return s.repo.GetEpisodeFileByPath(ctx, filePath)
}

func (s *tvService) GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*EpisodeFile, error) {
	return s.repo.GetEpisodeFileBySonarrID(ctx, sonarrFileID)
}

func (s *tvService) ListEpisodeFiles(ctx context.Context, episodeID uuid.UUID) ([]EpisodeFile, error) {
	return s.repo.ListEpisodeFilesByEpisode(ctx, episodeID)
}

func (s *tvService) CreateEpisodeFile(ctx context.Context, params CreateEpisodeFileParams) (*EpisodeFile, error) {
	// Verify episode exists
	_, err := s.repo.GetEpisode(ctx, params.EpisodeID)
	if err != nil {
		return nil, fmt.Errorf("episode not found: %w", err)
	}

	// Check if file already exists
	existing, err := s.repo.GetEpisodeFileByPath(ctx, params.FilePath)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("file already exists at path: %s", params.FilePath)
	}

	return s.repo.CreateEpisodeFile(ctx, params)
}

func (s *tvService) UpdateEpisodeFile(ctx context.Context, params UpdateEpisodeFileParams) (*EpisodeFile, error) {
	// Verify file exists
	_, err := s.repo.GetEpisodeFile(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("episode file not found: %w", err)
	}

	return s.repo.UpdateEpisodeFile(ctx, params)
}

func (s *tvService) DeleteEpisodeFile(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEpisodeFile(ctx, id)
}

// =============================================================================
// Credits Operations
// =============================================================================

func (s *tvService) GetSeriesCast(ctx context.Context, seriesID uuid.UUID) ([]SeriesCredit, error) {
	return s.repo.ListSeriesCast(ctx, seriesID)
}

func (s *tvService) GetSeriesCrew(ctx context.Context, seriesID uuid.UUID) ([]SeriesCredit, error) {
	return s.repo.ListSeriesCrew(ctx, seriesID)
}

func (s *tvService) GetEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error) {
	return s.repo.ListEpisodeGuestStars(ctx, episodeID)
}

func (s *tvService) GetEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error) {
	return s.repo.ListEpisodeCrew(ctx, episodeID)
}

// =============================================================================
// Genres & Networks Operations
// =============================================================================

func (s *tvService) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error) {
	return s.repo.ListSeriesGenres(ctx, seriesID)
}

func (s *tvService) GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]Network, error) {
	return s.repo.ListNetworksBySeries(ctx, seriesID)
}

// =============================================================================
// Watch Progress Operations
// =============================================================================

func (s *tvService) UpdateEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID, progressSeconds, durationSeconds int32) (*EpisodeWatched, error) {
	// Verify episode exists
	_, err := s.repo.GetEpisode(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("episode not found: %w", err)
	}

	// Calculate if completed (>90% watched)
	isCompleted := false
	if durationSeconds > 0 {
		progress := float64(progressSeconds) / float64(durationSeconds)
		isCompleted = progress > 0.90
	}

	params := CreateWatchProgressParams{
		UserID:          userID,
		EpisodeID:       episodeID,
		ProgressSeconds: progressSeconds,
		DurationSeconds: durationSeconds,
		IsCompleted:     isCompleted,
	}

	return s.repo.CreateOrUpdateWatchProgress(ctx, params)
}

func (s *tvService) GetEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatched, error) {
	return s.repo.GetWatchProgress(ctx, userID, episodeID)
}

func (s *tvService) MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) error {
	// Get episode to get duration
	episode, err := s.repo.GetEpisode(ctx, episodeID)
	if err != nil {
		return fmt.Errorf("episode not found: %w", err)
	}

	// Use runtime if available, otherwise default to 2700 seconds (45 minutes)
	durationSeconds := int32(2700)
	if episode.Runtime != nil && *episode.Runtime > 0 {
		durationSeconds = *episode.Runtime * 60 // Convert minutes to seconds
	}

	_, err = s.repo.MarkEpisodeWatched(ctx, userID, episodeID, durationSeconds)
	return err
}

func (s *tvService) MarkSeasonWatched(ctx context.Context, userID, seasonID uuid.UUID) error {
	// Get all episodes in season
	episodes, err := s.repo.ListEpisodesBySeason(ctx, seasonID)
	if err != nil {
		return fmt.Errorf("failed to list episodes: %w", err)
	}

	// Mark each episode as watched
	for _, ep := range episodes {
		if err := s.MarkEpisodeWatched(ctx, userID, ep.ID); err != nil {
			return fmt.Errorf("failed to mark episode %d watched: %w", ep.EpisodeNumber, err)
		}
	}

	return nil
}

func (s *tvService) MarkSeriesWatched(ctx context.Context, userID, seriesID uuid.UUID) error {
	// Get all seasons
	seasons, err := s.repo.ListSeasonsBySeries(ctx, seriesID)
	if err != nil {
		return fmt.Errorf("failed to list seasons: %w", err)
	}

	// Mark each season as watched
	for _, season := range seasons {
		if err := s.MarkSeasonWatched(ctx, userID, season.ID); err != nil {
			return fmt.Errorf("failed to mark season %d watched: %w", season.SeasonNumber, err)
		}
	}

	return nil
}

func (s *tvService) RemoveEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) error {
	return s.repo.DeleteWatchProgress(ctx, userID, episodeID)
}

func (s *tvService) RemoveSeriesProgress(ctx context.Context, userID, seriesID uuid.UUID) error {
	return s.repo.DeleteSeriesWatchProgress(ctx, userID, seriesID)
}

func (s *tvService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	return s.repo.ListContinueWatchingSeries(ctx, userID, limit)
}

func (s *tvService) GetNextEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error) {
	return s.repo.GetNextUnwatchedEpisode(ctx, userID, seriesID)
}

func (s *tvService) GetSeriesWatchStats(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchStats, error) {
	return s.repo.GetSeriesWatchStats(ctx, userID, seriesID)
}

func (s *tvService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserTVStats, error) {
	return s.repo.GetUserTVStats(ctx, userID)
}

// =============================================================================
// Metadata Operations
// =============================================================================

func (s *tvService) RefreshSeriesMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error {
	// Check if metadata provider is available
	if s.metadataProvider == nil {
		return fmt.Errorf("metadata provider not configured")
	}

	// Get the series
	series, err := s.repo.GetSeries(ctx, id)
	if err != nil {
		return fmt.Errorf("get series: %w", err)
	}

	// Enrich with latest metadata from TMDb, passing options through
	if err := s.metadataProvider.EnrichSeries(ctx, series, opts...); err != nil {
		return fmt.Errorf("enrich series: %w", err)
	}

	// Build update params from enriched series
	params := seriesToUpdateParams(series)

	// Update the series
	if _, err := s.repo.UpdateSeries(ctx, params); err != nil {
		return fmt.Errorf("update series: %w", err)
	}

	// Update credits if TMDbID is available
	if series.TMDbID != nil {
		credits, err := s.metadataProvider.GetSeriesCredits(ctx, series.ID, int(*series.TMDbID))
		if err == nil && len(credits) > 0 {
			_ = s.repo.DeleteSeriesCredits(ctx, series.ID)
			for _, credit := range credits {
				_, _ = s.repo.CreateSeriesCredit(ctx, CreateSeriesCreditParams{
					SeriesID:     credit.SeriesID,
					TMDbPersonID: credit.TMDbPersonID,
					Name:         credit.Name,
					CreditType:   credit.CreditType,
					Character:    credit.Character,
					Job:          credit.Job,
					Department:   credit.Department,
					CastOrder:    credit.CastOrder,
					ProfilePath:  credit.ProfilePath,
				})
			}
		}

		// Update genres
		genres, err := s.metadataProvider.GetSeriesGenres(ctx, series.ID, int(*series.TMDbID))
		if err == nil && len(genres) > 0 {
			_ = s.repo.DeleteSeriesGenres(ctx, series.ID)
			for _, genre := range genres {
				_ = s.repo.AddSeriesGenre(ctx, series.ID, genre.TMDbGenreID, genre.Name)
			}
		}
	}

	return nil
}

func (s *tvService) RefreshSeasonMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error {
	// Check if metadata provider is available
	if s.metadataProvider == nil {
		return fmt.Errorf("metadata provider not configured")
	}

	// Get the season
	season, err := s.repo.GetSeason(ctx, id)
	if err != nil {
		return fmt.Errorf("get season: %w", err)
	}

	// Get the series to get TMDbID
	series, err := s.repo.GetSeries(ctx, season.SeriesID)
	if err != nil {
		return fmt.Errorf("get series: %w", err)
	}

	if series.TMDbID == nil {
		return fmt.Errorf("series has no TMDb ID for metadata refresh")
	}

	// Enrich with latest metadata from TMDb, passing options through
	if err := s.metadataProvider.EnrichSeason(ctx, season, *series.TMDbID, opts...); err != nil {
		return fmt.Errorf("enrich season: %w", err)
	}

	// Build update params from enriched season
	params := seasonToUpdateParams(season)

	// Update the season
	if _, err := s.repo.UpdateSeason(ctx, params); err != nil {
		return fmt.Errorf("update season: %w", err)
	}

	return nil
}

func (s *tvService) RefreshEpisodeMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error {
	// Check if metadata provider is available
	if s.metadataProvider == nil {
		return fmt.Errorf("metadata provider not configured")
	}

	// Get the episode
	episode, err := s.repo.GetEpisode(ctx, id)
	if err != nil {
		return fmt.Errorf("get episode: %w", err)
	}

	// Get the series to get TMDbID
	series, err := s.repo.GetSeries(ctx, episode.SeriesID)
	if err != nil {
		return fmt.Errorf("get series: %w", err)
	}

	if series.TMDbID == nil {
		return fmt.Errorf("series has no TMDb ID for metadata refresh")
	}

	// Enrich with latest metadata from TMDb, passing options through
	if err := s.metadataProvider.EnrichEpisode(ctx, episode, *series.TMDbID, opts...); err != nil {
		return fmt.Errorf("enrich episode: %w", err)
	}

	// Build update params from enriched episode
	params := episodeToUpdateParams(episode)

	// Update the episode
	if _, err := s.repo.UpdateEpisode(ctx, params); err != nil {
		return fmt.Errorf("update episode: %w", err)
	}

	return nil
}

// seriesToUpdateParams converts a Series to UpdateSeriesParams
func seriesToUpdateParams(s *Series) UpdateSeriesParams {
	params := UpdateSeriesParams{
		ID:               s.ID,
		TMDbID:           s.TMDbID,
		TVDbID:           s.TVDbID,
		IMDbID:           s.IMDbID,
		SonarrID:         s.SonarrID,
		Title:            &s.Title,
		OriginalTitle:    s.OriginalTitle,
		OriginalLanguage: &s.OriginalLanguage,
		Tagline:          s.Tagline,
		Overview:         s.Overview,
		Status:           s.Status,
		Type:             s.Type,
		VoteCount:        s.VoteCount,
		PosterPath:       s.PosterPath,
		BackdropPath:     s.BackdropPath,
		TotalSeasons:     &s.TotalSeasons,
		TotalEpisodes:    &s.TotalEpisodes,
		TrailerURL:       s.TrailerURL,
		Homepage:         s.Homepage,
		TitlesI18n:       s.TitlesI18n,
		TaglinesI18n:     s.TaglinesI18n,
		OverviewsI18n:    s.OverviewsI18n,
		AgeRatings:       s.AgeRatings,
		ExternalRatings:  s.ExternalRatings,
	}
	if s.FirstAirDate != nil {
		d := s.FirstAirDate.Format("2006-01-02")
		params.FirstAirDate = &d
	}
	if s.LastAirDate != nil {
		d := s.LastAirDate.Format("2006-01-02")
		params.LastAirDate = &d
	}
	if s.VoteAverage != nil {
		v := s.VoteAverage.String()
		params.VoteAverage = &v
	}
	if s.Popularity != nil {
		p := s.Popularity.String()
		params.Popularity = &p
	}
	return params
}

// seasonToUpdateParams converts a Season to UpdateSeasonParams
func seasonToUpdateParams(s *Season) UpdateSeasonParams {
	params := UpdateSeasonParams{
		ID:            s.ID,
		TMDbID:        s.TMDbID,
		SeasonNumber:  &s.SeasonNumber,
		Name:          &s.Name,
		Overview:      s.Overview,
		PosterPath:    s.PosterPath,
		EpisodeCount:  &s.EpisodeCount,
		NamesI18n:     s.NamesI18n,
		OverviewsI18n: s.OverviewsI18n,
	}
	if s.AirDate != nil {
		d := s.AirDate.Format("2006-01-02")
		params.AirDate = &d
	}
	if s.VoteAverage != nil {
		v := s.VoteAverage.String()
		params.VoteAverage = &v
	}
	return params
}

// episodeToUpdateParams converts an Episode to UpdateEpisodeParams
func episodeToUpdateParams(e *Episode) UpdateEpisodeParams {
	params := UpdateEpisodeParams{
		ID:             e.ID,
		TMDbID:         e.TMDbID,
		TVDbID:         e.TVDbID,
		IMDbID:         e.IMDbID,
		SeasonNumber:   &e.SeasonNumber,
		EpisodeNumber:  &e.EpisodeNumber,
		Title:          &e.Title,
		Overview:       e.Overview,
		Runtime:        e.Runtime,
		VoteCount:      e.VoteCount,
		StillPath:      e.StillPath,
		ProductionCode: e.ProductionCode,
		TitlesI18n:     e.TitlesI18n,
		OverviewsI18n:  e.OverviewsI18n,
	}
	if e.AirDate != nil {
		d := e.AirDate.Format("2006-01-02")
		params.AirDate = &d
	}
	if e.VoteAverage != nil {
		v := e.VoteAverage.String()
		params.VoteAverage = &v
	}
	return params
}
