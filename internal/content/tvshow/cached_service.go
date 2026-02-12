// Package tvshow provides TV show-related business logic.
package tvshow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/infra/cache"
)

// CachedService wraps the TV show service with caching.
type CachedService struct {
	Service
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedService creates a new cached TV show service.
func NewCachedService(svc Service, cache *cache.Cache, logger *slog.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   cache,
		logger:  logger.With("component", "tvshow-cache"),
	}
}

// =============================================================================
// Series Operations (Cached Reads)
// =============================================================================

// GetSeries retrieves a series by ID with caching (5 min TTL).
func (s *CachedService) GetSeries(ctx context.Context, id uuid.UUID) (*Series, error) {
	if s.cache == nil {
		return s.Service.GetSeries(ctx, id)
	}

	cacheKey := cache.TVShowKey(id.String())

	var series Series
	if err := s.cache.GetJSON(ctx, cacheKey, &series); err == nil {
		s.logger.Debug("series cache hit", slog.String("id", id.String()))
		return &series, nil
	}

	s.logger.Debug("series cache miss", slog.String("id", id.String()))

	result, err := s.Service.GetSeries(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.TVShowTTL); setErr != nil {
			s.logger.Warn("failed to cache series", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// ListSeries returns a paginated list of series with caching (1 min TTL).
func (s *CachedService) ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error) {
	if s.cache == nil {
		return s.Service.ListSeries(ctx, filters)
	}

	cacheKey := s.listSeriesCacheKey(filters)

	var series []Series
	if err := s.cache.GetJSON(ctx, cacheKey, &series); err == nil {
		s.logger.Debug("list series cache hit", slog.String("key", cacheKey))
		return series, nil
	}

	s.logger.Debug("list series cache miss", slog.String("key", cacheKey))

	result, err := s.Service.ListSeries(ctx, filters)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, time.Minute); setErr != nil {
			s.logger.Warn("failed to cache series list", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// ListRecentlyAdded returns recently added series with caching (2 min TTL).
// Cache key includes pagination params, caches both items and total.
func (s *CachedService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Series, int64, error) {
	if s.cache == nil {
		return s.Service.ListRecentlyAdded(ctx, limit, offset)
	}

	cacheKey := fmt.Sprintf("%srecently-added:%d:%d", cache.KeyPrefixTVShow, limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedSeries struct {
		Items []Series `json:"items"`
		Total int64    `json:"total"`
	}
	var cached cachedSeries
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("recently added series cache hit")
		return cached.Items, cached.Total, nil
	}

	// Cache miss - load from database
	items, total, err := s.Service.ListRecentlyAdded(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		toCache := cachedSeries{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, 2*time.Minute); setErr != nil {
			s.logger.Warn("failed to cache recently added series", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
}

// =============================================================================
// Season Operations (Cached Reads)
// =============================================================================

// ListSeasons returns seasons for a series with caching (5 min TTL).
func (s *CachedService) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]Season, error) {
	if s.cache == nil {
		return s.Service.ListSeasons(ctx, seriesID)
	}

	cacheKey := cache.TVShowSeasonsKey(seriesID.String())

	var seasons []Season
	if err := s.cache.GetJSON(ctx, cacheKey, &seasons); err == nil {
		s.logger.Debug("seasons cache hit", slog.String("series_id", seriesID.String()))
		return seasons, nil
	}

	result, err := s.Service.ListSeasons(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.TVShowSeasonTTL); setErr != nil {
			s.logger.Warn("failed to cache seasons", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// ListEpisodesBySeason returns episodes for a season with caching (5 min TTL).
func (s *CachedService) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error) {
	if s.cache == nil {
		return s.Service.ListEpisodesBySeason(ctx, seasonID)
	}

	cacheKey := cache.TVShowEpisodesKey(seasonID.String())

	var episodes []Episode
	if err := s.cache.GetJSON(ctx, cacheKey, &episodes); err == nil {
		s.logger.Debug("episodes cache hit", slog.String("season_id", seasonID.String()))
		return episodes, nil
	}

	result, err := s.Service.ListEpisodesBySeason(ctx, seasonID)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.TVShowSeasonTTL); setErr != nil {
			s.logger.Warn("failed to cache episodes", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// =============================================================================
// Credits / Genres / Networks (Cached Reads)
// =============================================================================

// GetSeriesCast returns the cast for a series with caching (10 min TTL).
// Cache key includes pagination params.
func (s *CachedService) GetSeriesCast(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, int64, error) {
	if s.cache == nil {
		return s.Service.GetSeriesCast(ctx, seriesID, limit, offset)
	}

	cacheKey := fmt.Sprintf("%s%s:%d:%d", cache.KeyPrefixTVShowCast, seriesID.String(), limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedCredits struct {
		Items []SeriesCredit `json:"items"`
		Total int64          `json:"total"`
	}
	var cached cachedCredits
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("series cast cache hit", slog.String("series_id", seriesID.String()))
		return cached.Items, cached.Total, nil
	}

	// Cache miss - load from database
	items, total, err := s.Service.GetSeriesCast(ctx, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		toCache := cachedCredits{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, cache.TVShowMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache series cast", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
}

// GetSeriesCrew returns the crew for a series with caching (10 min TTL).
// Cache key includes pagination params.
func (s *CachedService) GetSeriesCrew(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, int64, error) {
	if s.cache == nil {
		return s.Service.GetSeriesCrew(ctx, seriesID, limit, offset)
	}

	cacheKey := fmt.Sprintf("%s%s:%d:%d", cache.KeyPrefixTVShowCrew, seriesID.String(), limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedCredits struct {
		Items []SeriesCredit `json:"items"`
		Total int64          `json:"total"`
	}
	var cached cachedCredits
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("series crew cache hit", slog.String("series_id", seriesID.String()))
		return cached.Items, cached.Total, nil
	}

	// Cache miss - load from database
	items, total, err := s.Service.GetSeriesCrew(ctx, seriesID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		toCache := cachedCredits{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, cache.TVShowMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache series crew", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
}

// GetSeriesGenres returns genres for a series with caching (10 min TTL).
func (s *CachedService) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error) {
	if s.cache == nil {
		return s.Service.GetSeriesGenres(ctx, seriesID)
	}

	cacheKey := cache.TVShowGenresKey(seriesID.String())

	var genres []SeriesGenre
	if err := s.cache.GetJSON(ctx, cacheKey, &genres); err == nil {
		s.logger.Debug("series genres cache hit", slog.String("series_id", seriesID.String()))
		return genres, nil
	}

	result, err := s.Service.GetSeriesGenres(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.TVShowMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache series genres", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// GetSeriesNetworks returns networks for a series with caching (10 min TTL).
func (s *CachedService) GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]Network, error) {
	if s.cache == nil {
		return s.Service.GetSeriesNetworks(ctx, seriesID)
	}

	cacheKey := cache.TVShowNetworksKey(seriesID.String())

	var networks []Network
	if err := s.cache.GetJSON(ctx, cacheKey, &networks); err == nil {
		s.logger.Debug("series networks cache hit", slog.String("series_id", seriesID.String()))
		return networks, nil
	}

	result, err := s.Service.GetSeriesNetworks(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.TVShowMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache series networks", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// =============================================================================
// Watch Progress (Cached Reads)
// =============================================================================

// GetContinueWatching returns series the user is currently watching with caching (1 min TTL).
func (s *CachedService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	if s.cache == nil {
		return s.Service.GetContinueWatching(ctx, userID, limit)
	}

	cacheKey := fmt.Sprintf("%scontinue-watching:%s:%d", cache.KeyPrefixTVShow, userID.String(), limit)

	var items []ContinueWatchingItem
	if err := s.cache.GetJSON(ctx, cacheKey, &items); err == nil {
		s.logger.Debug("continue watching cache hit", slog.String("user_id", userID.String()))
		return items, nil
	}

	result, err := s.Service.GetContinueWatching(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, time.Minute); setErr != nil {
			s.logger.Warn("failed to cache continue watching", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// =============================================================================
// Write Operations (Cache Invalidation)
// =============================================================================

// UpdateSeries updates a series and invalidates cache.
func (s *CachedService) UpdateSeries(ctx context.Context, params UpdateSeriesParams) (*Series, error) {
	result, err := s.Service.UpdateSeries(ctx, params)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.invalidateSeries(ctx, params.ID)
	}

	return result, nil
}

// DeleteSeries deletes a series and invalidates cache.
func (s *CachedService) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	if err := s.Service.DeleteSeries(ctx, id); err != nil {
		return err
	}

	if s.cache != nil {
		s.invalidateSeries(ctx, id)
	}

	return nil
}

// UpdateEpisodeProgress updates watch progress and invalidates continue watching cache.
func (s *CachedService) UpdateEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID, progressSeconds, durationSeconds int32) (*EpisodeWatched, error) {
	result, err := s.Service.UpdateEpisodeProgress(ctx, userID, episodeID, progressSeconds, durationSeconds)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return result, nil
}

// MarkEpisodeWatched marks an episode as watched and invalidates continue watching cache.
func (s *CachedService) MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) error {
	if err := s.Service.MarkEpisodeWatched(ctx, userID, episodeID); err != nil {
		return err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return nil
}

// MarkEpisodesWatchedBulk marks multiple episodes as watched and invalidates continue watching cache.
func (s *CachedService) MarkEpisodesWatchedBulk(ctx context.Context, userID uuid.UUID, episodeIDs []uuid.UUID) (int64, error) {
	affected, err := s.Service.MarkEpisodesWatchedBulk(ctx, userID, episodeIDs)
	if err != nil {
		return 0, err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return affected, nil
}

// MarkSeasonWatched marks a season as watched and invalidates continue watching cache.
func (s *CachedService) MarkSeasonWatched(ctx context.Context, userID, seasonID uuid.UUID) error {
	if err := s.Service.MarkSeasonWatched(ctx, userID, seasonID); err != nil {
		return err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return nil
}

// MarkSeriesWatched marks a series as watched and invalidates continue watching cache.
func (s *CachedService) MarkSeriesWatched(ctx context.Context, userID, seriesID uuid.UUID) error {
	if err := s.Service.MarkSeriesWatched(ctx, userID, seriesID); err != nil {
		return err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return nil
}

// RemoveEpisodeProgress removes episode progress and invalidates continue watching cache.
func (s *CachedService) RemoveEpisodeProgress(ctx context.Context, userID, episodeID uuid.UUID) error {
	if err := s.Service.RemoveEpisodeProgress(ctx, userID, episodeID); err != nil {
		return err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return nil
}

// RemoveSeriesProgress removes series progress and invalidates continue watching cache.
func (s *CachedService) RemoveSeriesProgress(ctx context.Context, userID, seriesID uuid.UUID) error {
	if err := s.Service.RemoveSeriesProgress(ctx, userID, seriesID); err != nil {
		return err
	}

	if s.cache != nil {
		s.invalidateContinueWatching(ctx, userID)
	}

	return nil
}

// =============================================================================
// Cache Invalidation Helpers
// =============================================================================

// invalidateSeries invalidates all cache entries for a series.
func (s *CachedService) invalidateSeries(ctx context.Context, seriesID uuid.UUID) {
	id := seriesID.String()

	// Delete series itself
	if err := s.cache.Delete(ctx, cache.TVShowKey(id)); err != nil {
		s.logger.Warn("failed to invalidate series cache", slog.Any("error", err))
	}

	// Delete series metadata (cast, crew, genres, networks, seasons)
	keys := []string{
		cache.TVShowCastKey(id),
		cache.TVShowCrewKey(id),
		cache.TVShowGenresKey(id),
		cache.TVShowNetworksKey(id),
		cache.TVShowSeasonsKey(id),
	}
	for _, key := range keys {
		if err := s.cache.Delete(ctx, key); err != nil {
			s.logger.Warn("failed to invalidate series metadata cache", slog.String("key", key), slog.Any("error", err))
		}
	}

	// Invalidate episode caches for this series
	pattern := cache.KeyPrefixTVShowEpisodes + "*"
	if err := s.cache.Invalidate(ctx, pattern); err != nil {
		s.logger.Warn("failed to invalidate episode caches", slog.Any("error", err))
	}

	// Invalidate list caches (recently-added, etc.)
	patterns := []string{
		cache.KeyPrefixTVShowRecent + "*",
		cache.KeyPrefixTVShowList + "*",
	}
	for _, p := range patterns {
		if err := s.cache.Invalidate(ctx, p); err != nil {
			s.logger.Warn("failed to invalidate series list cache", slog.String("pattern", p), slog.Any("error", err))
		}
	}
}

// invalidateContinueWatching invalidates the continue watching cache for a user.
func (s *CachedService) invalidateContinueWatching(ctx context.Context, userID uuid.UUID) {
	pattern := fmt.Sprintf("%scontinue-watching:%s:*", cache.KeyPrefixTVShow, userID.String())
	if err := s.cache.Invalidate(ctx, pattern); err != nil {
		s.logger.Warn("failed to invalidate continue watching cache", slog.Any("error", err))
	}
}

// InvalidateAllTVShows invalidates all TV show caches.
func (s *CachedService) InvalidateAllTVShows(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Invalidate(ctx, cache.KeyPrefixTVShow+"*")
}

// listSeriesCacheKey generates a cache key for series list queries.
func (s *CachedService) listSeriesCacheKey(filters SeriesListFilters) string {
	key := fmt.Sprintf("%s:%d:%d",
		filters.OrderBy,
		filters.Limit,
		filters.Offset,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixTVShowList + hex.EncodeToString(hash[:8])
}

// ListDistinctGenres returns all distinct TV show genres with caching (10 min TTL).
// Genre lists are near-static and frequently requested for filter UIs.
func (s *CachedService) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	if s.cache == nil {
		return s.Service.ListDistinctGenres(ctx)
	}

	cacheKey := cache.KeyPrefixTVShowGenres + "distinct"

	// Try cache first
	var genres []content.GenreSummary
	if err := s.cache.GetJSON(ctx, cacheKey, &genres); err == nil {
		s.logger.Debug("tvshow distinct genres cache hit")
		return genres, nil
	}

	// Cache miss - load from database
	result, err := s.Service.ListDistinctGenres(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.TVShowMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache distinct tvshow genres", slog.Any("error", setErr))
		}
	}()

	return result, nil
}
