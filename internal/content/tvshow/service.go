package tvshow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/maypok86/otter"

	"github.com/lusoris/revenge/pkg/resilience"
)

// Service errors.
var (
	ErrSeriesNotFoundInService  = errors.New("series not found")
	ErrSeasonNotFoundInService  = errors.New("season not found")
	ErrEpisodeNotFoundInService = errors.New("episode not found")
	ErrLibraryNotFound          = errors.New("library not found")
	ErrMetadataUnavailable      = errors.New("metadata provider unavailable")
)

// ServiceConfig holds service configuration.
type ServiceConfig struct {
	CacheMaxEntries int           `koanf:"cache_max_entries"`
	CacheTTL        time.Duration `koanf:"cache_ttl"`
}

// DefaultServiceConfig returns sensible defaults.
var DefaultServiceConfig = ServiceConfig{
	CacheMaxEntries: 10_000,
	CacheTTL:        5 * time.Minute,
}

// Service provides TV show business logic with caching and resilience.
type Service struct {
	repo   Repository
	cache  otter.Cache[string, []byte]
	logger *slog.Logger
	config ServiceConfig
}

// NewService creates a new TV show service.
func NewService(
	repo Repository,
	logger *slog.Logger,
	config ServiceConfig,
) (*Service, error) {
	if config.CacheMaxEntries == 0 {
		config.CacheMaxEntries = DefaultServiceConfig.CacheMaxEntries
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = DefaultServiceConfig.CacheTTL
	}

	// Create local cache for hot TV show data
	cache, err := otter.MustBuilder[string, []byte](config.CacheMaxEntries).
		CollectStats().
		Cost(func(key string, value []byte) uint32 {
			return uint32(len(value))
		}).
		WithTTL(config.CacheTTL).
		Build()
	if err != nil {
		return nil, fmt.Errorf("create cache: %w", err)
	}

	return &Service{
		repo:   repo,
		cache:  cache,
		logger: logger.With("service", "tvshow"),
		config: config,
	}, nil
}

// Close releases service resources.
func (s *Service) Close() {
	s.cache.Close()
}

// =============================================================================
// Series Operations
// =============================================================================

// GetSeries retrieves a series by ID with caching.
func (s *Service) GetSeries(ctx context.Context, id uuid.UUID) (*Series, error) {
	key := fmt.Sprintf("series:%s", id)

	// Check local cache first
	if cached, found := s.cache.Get(key); found {
		var series Series
		if err := json.Unmarshal(cached, &series); err == nil {
			return &series, nil
		}
		s.cache.Delete(key)
	}

	// Fetch from repository
	series, err := s.repo.GetSeriesByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSeriesNotFound) {
			return nil, ErrSeriesNotFoundInService
		}
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(series); err == nil {
		s.cache.Set(key, data)
	}

	return series, nil
}

// GetSeriesWithRelations retrieves a series with all related data.
func (s *Service) GetSeriesWithRelations(ctx context.Context, id uuid.UUID) (*Series, error) {
	series, err := s.GetSeries(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load relations in parallel for performance
	errCh := make(chan error, 5)

	go func() {
		seasons, err := s.repo.ListSeasons(ctx, id)
		if err == nil {
			// Convert []*Season to []Season
			series.Seasons = make([]Season, len(seasons))
			for i, season := range seasons {
				series.Seasons[i] = *season
			}
		}
		errCh <- err
	}()

	go func() {
		genres, err := s.repo.GetSeriesGenres(ctx, id)
		if err == nil {
			series.Genres = genres
		}
		errCh <- err
	}()

	go func() {
		cast, err := s.repo.GetSeriesCast(ctx, id)
		if err == nil {
			series.Cast = cast
		}
		errCh <- err
	}()

	go func() {
		crew, err := s.repo.GetSeriesCrew(ctx, id)
		if err == nil {
			series.Crew = crew
			// Extract creators
			for _, c := range crew {
				if c.Role == "creator" || c.Role == "showrunner" {
					series.Creators = append(series.Creators, c)
				}
			}
		}
		errCh <- err
	}()

	go func() {
		images, err := s.repo.GetSeriesImages(ctx, id)
		if err == nil {
			series.Images = images
		}
		errCh <- err
	}()

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		if err := <-errCh; err != nil {
			s.logger.Warn("failed to load series relation", "series_id", id, "error", err)
		}
	}

	return series, nil
}

// ListSeries retrieves series with pagination.
func (s *Service) ListSeries(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Series, int64, error) {
	series, err := s.repo.ListSeriesByLibrary(ctx, libraryID, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountSeriesByLibrary(ctx, libraryID)
	if err != nil {
		return nil, 0, err
	}

	return series, total, nil
}

// ListAllSeries retrieves all series with pagination.
func (s *Service) ListAllSeries(ctx context.Context, params ListParams) ([]*Series, int64, error) {
	series, err := s.repo.ListSeries(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountSeries(ctx)
	if err != nil {
		return nil, 0, err
	}

	return series, total, nil
}

// SearchSeries searches series by query.
func (s *Service) SearchSeries(ctx context.Context, query string, params ListParams) ([]*Series, error) {
	return s.repo.SearchSeries(ctx, query, params)
}

// CreateSeries creates a new series.
func (s *Service) CreateSeries(ctx context.Context, series *Series) error {
	if err := s.repo.CreateSeries(ctx, series); err != nil {
		return err
	}

	s.logger.Info("series created",
		"id", series.ID,
		"title", series.Title,
	)

	return nil
}

// UpdateSeries updates a series.
func (s *Service) UpdateSeries(ctx context.Context, series *Series) error {
	if err := s.repo.UpdateSeries(ctx, series); err != nil {
		return err
	}

	// Invalidate cache
	key := fmt.Sprintf("series:%s", series.ID)
	s.cache.Delete(key)

	return nil
}

// DeleteSeries deletes a series.
func (s *Service) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteSeries(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	key := fmt.Sprintf("series:%s", id)
	s.cache.Delete(key)

	s.logger.Info("series deleted", "id", id)
	return nil
}

// =============================================================================
// Season Operations
// =============================================================================

// GetSeason retrieves a season by ID.
func (s *Service) GetSeason(ctx context.Context, id uuid.UUID) (*Season, error) {
	season, err := s.repo.GetSeasonByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSeasonNotFound) {
			return nil, ErrSeasonNotFoundInService
		}
		return nil, err
	}
	return season, nil
}

// GetSeasonByNumber retrieves a season by series ID and season number.
func (s *Service) GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int) (*Season, error) {
	season, err := s.repo.GetSeasonByNumber(ctx, seriesID, seasonNumber)
	if err != nil {
		if errors.Is(err, ErrSeasonNotFound) {
			return nil, ErrSeasonNotFoundInService
		}
		return nil, err
	}
	return season, nil
}

// GetSeasonWithEpisodes retrieves a season with its episodes.
func (s *Service) GetSeasonWithEpisodes(ctx context.Context, id uuid.UUID) (*Season, error) {
	season, err := s.GetSeason(ctx, id)
	if err != nil {
		return nil, err
	}

	episodes, err := s.repo.ListEpisodesBySeason(ctx, id)
	if err != nil {
		s.logger.Warn("failed to load season episodes", "season_id", id, "error", err)
	} else {
		// Convert []*Episode to []Episode
		season.Episodes = make([]Episode, len(episodes))
		for i, ep := range episodes {
			season.Episodes[i] = *ep
		}
	}

	return season, nil
}

// ListSeasons retrieves all seasons for a series.
func (s *Service) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]*Season, error) {
	return s.repo.ListSeasons(ctx, seriesID)
}

// CreateSeason creates a new season.
func (s *Service) CreateSeason(ctx context.Context, season *Season) error {
	if err := s.repo.CreateSeason(ctx, season); err != nil {
		return err
	}

	s.logger.Info("season created",
		"id", season.ID,
		"series_id", season.SeriesID,
		"season_number", season.SeasonNumber,
	)

	return nil
}

// UpdateSeason updates a season.
func (s *Service) UpdateSeason(ctx context.Context, season *Season) error {
	return s.repo.UpdateSeason(ctx, season)
}

// DeleteSeason deletes a season.
func (s *Service) DeleteSeason(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteSeason(ctx, id); err != nil {
		return err
	}

	s.logger.Info("season deleted", "id", id)
	return nil
}

// =============================================================================
// Episode Operations
// =============================================================================

// GetEpisode retrieves an episode by ID.
func (s *Service) GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error) {
	episode, err := s.repo.GetEpisodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrEpisodeNotFound) {
			return nil, ErrEpisodeNotFoundInService
		}
		return nil, err
	}
	return episode, nil
}

// GetEpisodeByNumber retrieves an episode by series ID, season number, and episode number.
func (s *Service) GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error) {
	episode, err := s.repo.GetEpisodeByNumber(ctx, seriesID, seasonNumber, episodeNumber)
	if err != nil {
		if errors.Is(err, ErrEpisodeNotFound) {
			return nil, ErrEpisodeNotFoundInService
		}
		return nil, err
	}
	return episode, nil
}

// GetEpisodeWithRelations retrieves an episode with all related data.
func (s *Service) GetEpisodeWithRelations(ctx context.Context, id uuid.UUID) (*Episode, error) {
	episode, err := s.GetEpisode(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load relations in parallel (cast and crew)
	errCh := make(chan error, 2)

	go func() {
		cast, err := s.repo.GetEpisodeCast(ctx, id)
		if err == nil {
			episode.Cast = cast
		}
		errCh <- err
	}()

	go func() {
		crew, err := s.repo.GetEpisodeCrew(ctx, id)
		if err == nil {
			episode.Crew = crew
			// Extract directors and writers
			for _, c := range crew {
				switch c.Role {
				case "director":
					episode.Directors = append(episode.Directors, c)
				case "writer":
					episode.Writers = append(episode.Writers, c)
				}
			}
		}
		errCh <- err
	}()

	// Wait for all goroutines
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			s.logger.Warn("failed to load episode relation", "episode_id", id, "error", err)
		}
	}

	return episode, nil
}

// ListEpisodes retrieves all episodes for a series.
func (s *Service) ListEpisodes(ctx context.Context, seriesID uuid.UUID) ([]*Episode, error) {
	return s.repo.ListEpisodes(ctx, seriesID)
}

// ListEpisodesBySeason retrieves all episodes for a specific season.
func (s *Service) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]*Episode, error) {
	return s.repo.ListEpisodesBySeason(ctx, seasonID)
}

// ListRecentlyAiredEpisodes returns recently aired episodes.
func (s *Service) ListRecentlyAiredEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error) {
	return s.repo.ListRecentlyAiredEpisodes(ctx, libraryIDs, limit)
}

// ListUpcomingEpisodes returns upcoming episodes.
func (s *Service) ListUpcomingEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error) {
	return s.repo.ListUpcomingEpisodes(ctx, libraryIDs, limit)
}

// GetNextEpisode returns the next episode to watch after the given episode.
func (s *Service) GetNextEpisode(ctx context.Context, episodeID uuid.UUID) (*Episode, error) {
	episode, err := s.GetEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetNextEpisode(ctx, episode.SeriesID, episode.SeasonNumber, episode.EpisodeNumber)
}

// GetPreviousEpisode returns the previous episode before the given episode.
func (s *Service) GetPreviousEpisode(ctx context.Context, episodeID uuid.UUID) (*Episode, error) {
	episode, err := s.GetEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetPreviousEpisode(ctx, episode.SeriesID, episode.SeasonNumber, episode.EpisodeNumber)
}

// CreateEpisode creates a new episode.
func (s *Service) CreateEpisode(ctx context.Context, episode *Episode) error {
	if err := s.repo.CreateEpisode(ctx, episode); err != nil {
		return err
	}

	s.logger.Info("episode created",
		"id", episode.ID,
		"season_id", episode.SeasonID,
		"episode_number", episode.EpisodeNumber,
	)

	return nil
}

// UpdateEpisode updates an episode.
func (s *Service) UpdateEpisode(ctx context.Context, episode *Episode) error {
	return s.repo.UpdateEpisode(ctx, episode)
}

// DeleteEpisode deletes an episode.
func (s *Service) DeleteEpisode(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteEpisode(ctx, id); err != nil {
		return err
	}

	s.logger.Info("episode deleted", "id", id)
	return nil
}

// =============================================================================
// User Data Operations - Series
// =============================================================================

// GetSeriesUserRating returns a user's rating for a series.
func (s *Service) GetSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesUserRating, error) {
	return s.repo.GetSeriesUserRating(ctx, userID, seriesID)
}

// SetSeriesUserRating sets a user's rating for a series.
func (s *Service) SetSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID, rating float64, review string) error {
	return s.repo.SetSeriesUserRating(ctx, userID, seriesID, rating, review)
}

// DeleteSeriesUserRating removes a user's rating for a series.
func (s *Service) DeleteSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID) error {
	return s.repo.DeleteSeriesUserRating(ctx, userID, seriesID)
}

// IsSeriesFavorite checks if a series is in user's favorites.
func (s *Service) IsSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) (bool, error) {
	return s.repo.IsSeriesFavorite(ctx, userID, seriesID)
}

// AddSeriesFavorite adds a series to user's favorites.
func (s *Service) AddSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) error {
	return s.repo.AddSeriesFavorite(ctx, userID, seriesID)
}

// RemoveSeriesFavorite removes a series from user's favorites.
func (s *Service) RemoveSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) error {
	return s.repo.RemoveSeriesFavorite(ctx, userID, seriesID)
}

// ListSeriesFavorites returns user's favorite series.
func (s *Service) ListSeriesFavorites(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Series, int64, error) {
	series, err := s.repo.ListFavoriteSeries(ctx, userID, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountFavoriteSeries(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return series, total, nil
}

// IsSeriesInWatchlist checks if a series is in user's watchlist.
func (s *Service) IsSeriesInWatchlist(ctx context.Context, userID, seriesID uuid.UUID) (bool, error) {
	return s.repo.IsSeriesInWatchlist(ctx, userID, seriesID)
}

// AddSeriesToWatchlist adds a series to user's watchlist.
func (s *Service) AddSeriesToWatchlist(ctx context.Context, userID, seriesID uuid.UUID) error {
	return s.repo.AddSeriesToWatchlist(ctx, userID, seriesID)
}

// RemoveSeriesFromWatchlist removes a series from user's watchlist.
func (s *Service) RemoveSeriesFromWatchlist(ctx context.Context, userID, seriesID uuid.UUID) error {
	return s.repo.RemoveSeriesFromWatchlist(ctx, userID, seriesID)
}

// ListSeriesWatchlist returns user's series watchlist.
func (s *Service) ListSeriesWatchlist(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Series, int64, error) {
	series, err := s.repo.ListSeriesWatchlist(ctx, userID, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountSeriesWatchlist(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return series, total, nil
}

// GetSeriesWatchProgress returns watch progress for a series.
func (s *Service) GetSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchProgress, error) {
	return s.repo.GetSeriesWatchProgress(ctx, userID, seriesID)
}

// =============================================================================
// User Data Operations - Episodes
// =============================================================================

// GetEpisodeUserRating returns a user's rating for an episode.
func (s *Service) GetEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeUserRating, error) {
	return s.repo.GetEpisodeUserRating(ctx, userID, episodeID)
}

// SetEpisodeUserRating sets a user's rating for an episode.
func (s *Service) SetEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID, rating float64) error {
	return s.repo.SetEpisodeUserRating(ctx, userID, episodeID, rating)
}

// DeleteEpisodeUserRating removes a user's rating for an episode.
func (s *Service) DeleteEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID) error {
	return s.repo.DeleteEpisodeUserRating(ctx, userID, episodeID)
}

// GetEpisodeWatchHistory returns watch history for an episode.
func (s *Service) GetEpisodeWatchHistory(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatchHistory, error) {
	return s.repo.GetEpisodeWatchHistory(ctx, userID, episodeID)
}

// UpdateEpisodeWatchProgress updates the watch position for an episode.
func (s *Service) UpdateEpisodeWatchProgress(ctx context.Context, history *EpisodeWatchHistory) error {
	existing, err := s.repo.GetEpisodeWatchHistory(ctx, history.UserID, history.EpisodeID)
	if err != nil {
		return err
	}

	if existing == nil {
		return s.repo.CreateEpisodeWatchHistory(ctx, history)
	}

	return s.repo.UpdateEpisodeWatchHistory(ctx, existing.ID, history.PositionTicks, &history.DurationTicks)
}

// MarkEpisodeAsWatched marks an episode as completely watched.
func (s *Service) MarkEpisodeAsWatched(ctx context.Context, userID, episodeID uuid.UUID) error {
	existing, err := s.repo.GetEpisodeWatchHistory(ctx, userID, episodeID)
	if err != nil {
		return err
	}

	if existing == nil {
		history := &EpisodeWatchHistory{
			UserID:    userID,
			EpisodeID: episodeID,
			Completed: true,
		}
		return s.repo.CreateEpisodeWatchHistory(ctx, history)
	}

	return s.repo.MarkEpisodeWatchHistoryCompleted(ctx, existing.ID)
}

// MarkEpisodeAsUnwatched marks an episode as unwatched.
func (s *Service) MarkEpisodeAsUnwatched(ctx context.Context, userID, episodeID uuid.UUID) error {
	history, err := s.repo.GetEpisodeWatchHistory(ctx, userID, episodeID)
	if err != nil || history == nil {
		return nil // Already unwatched
	}
	return s.repo.DeleteEpisodeWatchHistory(ctx, history.ID)
}

// IsEpisodeWatched checks if a user has watched an episode.
func (s *Service) IsEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) (bool, error) {
	return s.repo.IsEpisodeWatched(ctx, userID, episodeID)
}

// ListResumeableEpisodes returns episodes the user can resume.
func (s *Service) ListResumeableEpisodes(ctx context.Context, userID uuid.UUID, limit int) ([]EpisodeWatchHistory, error) {
	return s.repo.ListResumeableEpisodes(ctx, userID, limit)
}

// ListContinueWatchingSeries returns series the user is currently watching.
func (s *Service) ListContinueWatchingSeries(ctx context.Context, userID uuid.UUID, limit int) ([]*SeriesWatchProgress, error) {
	return s.repo.ListContinueWatchingSeries(ctx, userID, limit)
}

// =============================================================================
// Metadata Operations
// =============================================================================

// SeriesMetadataUpdate contains fields that can be updated from metadata providers.
type SeriesMetadataUpdate struct {
	Title           string
	OriginalTitle   string
	Overview        string
	FirstAirDate    *time.Time
	Status          string
	Type            string
	CommunityRating float64
	VoteCount       int
	PosterPath      string
	PosterBlurhash  string
	BackdropPath    string
	BackdropBlurhash string
	TmdbID          int
	TvdbID          int
	ImdbID          string
}

// ApplySeriesMetadata updates series fields from external metadata.
func (s *Service) ApplySeriesMetadata(ctx context.Context, id uuid.UUID, updates SeriesMetadataUpdate) error {
	series, err := s.repo.GetSeriesByID(ctx, id)
	if err != nil {
		return err
	}

	if updates.Title != "" {
		series.Title = updates.Title
	}
	if updates.OriginalTitle != "" {
		series.OriginalTitle = updates.OriginalTitle
	}
	if updates.Overview != "" {
		series.Overview = updates.Overview
	}
	if updates.FirstAirDate != nil {
		series.FirstAirDate = updates.FirstAirDate
		series.Year = updates.FirstAirDate.Year()
	}
	if updates.Status != "" {
		series.Status = updates.Status
	}
	if updates.Type != "" {
		series.Type = updates.Type
	}
	if updates.CommunityRating > 0 {
		series.CommunityRating = updates.CommunityRating
	}
	if updates.VoteCount > 0 {
		series.VoteCount = updates.VoteCount
	}
	if updates.PosterPath != "" {
		series.PosterPath = updates.PosterPath
	}
	if updates.PosterBlurhash != "" {
		series.PosterBlurhash = updates.PosterBlurhash
	}
	if updates.BackdropPath != "" {
		series.BackdropPath = updates.BackdropPath
	}
	if updates.BackdropBlurhash != "" {
		series.BackdropBlurhash = updates.BackdropBlurhash
	}
	if updates.TmdbID > 0 {
		series.TmdbID = updates.TmdbID
	}
	if updates.TvdbID > 0 {
		series.TvdbID = updates.TvdbID
	}
	if updates.ImdbID != "" {
		series.ImdbID = updates.ImdbID
	}

	if err := s.repo.UpdateSeries(ctx, series); err != nil {
		return err
	}

	// Invalidate cache
	key := fmt.Sprintf("series:%s", id)
	s.cache.Delete(key)

	s.logger.Info("series metadata updated",
		"id", id,
		"title", series.Title,
		"tmdb_id", series.TmdbID,
	)

	return nil
}

// =============================================================================
// Cache and Statistics
// =============================================================================

// CacheStats returns local cache statistics.
func (s *Service) CacheStats() CacheStats {
	stats := s.cache.Stats()
	return CacheStats{
		Hits:      stats.Hits(),
		Misses:    stats.Misses(),
		Ratio:     stats.Ratio(),
		Evictions: stats.EvictedCount(),
	}
}

// CacheStats contains cache statistics.
type CacheStats struct {
	Hits      int64
	Misses    int64
	Ratio     float64
	Evictions int64
}

// =============================================================================
// Retry Wrapper
// =============================================================================

// ServiceWithRetry wraps service with retry capability.
type ServiceWithRetry struct {
	*Service
	retry resilience.Retry
}

// NewServiceWithRetry creates a service with retry wrapper.
func NewServiceWithRetry(service *Service) *ServiceWithRetry {
	return &ServiceWithRetry{
		Service: service,
		retry:   resilience.DefaultRetry(),
	}
}

// GetSeriesWithRetry retrieves a series with automatic retry.
func (s *ServiceWithRetry) GetSeriesWithRetry(ctx context.Context, id uuid.UUID) (*Series, error) {
	var series *Series
	err := s.retry.DoWithContext(ctx, func(ctx context.Context) error {
		var err error
		series, err = s.Service.GetSeries(ctx, id)
		return err
	})
	return series, err
}

// ApplySeriesMetadataWithRetry applies metadata updates with automatic retry.
func (s *ServiceWithRetry) ApplySeriesMetadataWithRetry(ctx context.Context, id uuid.UUID, updates SeriesMetadataUpdate) error {
	return s.retry.DoWithContext(ctx, func(ctx context.Context) error {
		return s.Service.ApplySeriesMetadata(ctx, id, updates)
	})
}
