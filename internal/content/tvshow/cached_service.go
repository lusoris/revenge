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
	return cache.Get(ctx, s.cache, cache.TVShowKey(id.String()), cache.TVShowTTL, func(ctx context.Context) (*Series, error) {
		return s.Service.GetSeries(ctx, id)
	})
}

// ListSeries returns a paginated list of series with caching (1 min TTL).
func (s *CachedService) ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error) {
	return cache.Get(ctx, s.cache, s.listSeriesCacheKey(filters), time.Minute, func(ctx context.Context) ([]Series, error) {
		return s.Service.ListSeries(ctx, filters)
	})
}

// ListRecentlyAdded returns recently added series with caching (2 min TTL).
func (s *CachedService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Series, int64, error) {
	key := fmt.Sprintf("%srecently-added:%d:%d", cache.KeyPrefixTVShow, limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.RecentlyAddedTTL, func(ctx context.Context) (cache.Pair[[]Series], error) {
		items, total, err := s.Service.ListRecentlyAdded(ctx, limit, offset)
		return cache.Pair[[]Series]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// =============================================================================
// Season Operations (Cached Reads)
// =============================================================================

// ListSeasons returns seasons for a series with caching (5 min TTL).
func (s *CachedService) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]Season, error) {
	return cache.Get(ctx, s.cache, cache.TVShowSeasonsKey(seriesID.String()), cache.TVShowSeasonTTL, func(ctx context.Context) ([]Season, error) {
		return s.Service.ListSeasons(ctx, seriesID)
	})
}

// ListEpisodesBySeason returns episodes for a season with caching (5 min TTL).
func (s *CachedService) ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error) {
	return cache.Get(ctx, s.cache, cache.TVShowEpisodesKey(seasonID.String()), cache.TVShowSeasonTTL, func(ctx context.Context) ([]Episode, error) {
		return s.Service.ListEpisodesBySeason(ctx, seasonID)
	})
}

// =============================================================================
// Credits / Genres / Networks (Cached Reads)
// =============================================================================

// GetSeriesCast returns the cast for a series with caching (10 min TTL).
func (s *CachedService) GetSeriesCast(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, int64, error) {
	key := fmt.Sprintf("%s%s:%d:%d", cache.KeyPrefixTVShowCast, seriesID.String(), limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.TVShowMetaTTL, func(ctx context.Context) (cache.Pair[[]SeriesCredit], error) {
		items, total, err := s.Service.GetSeriesCast(ctx, seriesID, limit, offset)
		return cache.Pair[[]SeriesCredit]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// GetSeriesCrew returns the crew for a series with caching (10 min TTL).
func (s *CachedService) GetSeriesCrew(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, int64, error) {
	key := fmt.Sprintf("%s%s:%d:%d", cache.KeyPrefixTVShowCrew, seriesID.String(), limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.TVShowMetaTTL, func(ctx context.Context) (cache.Pair[[]SeriesCredit], error) {
		items, total, err := s.Service.GetSeriesCrew(ctx, seriesID, limit, offset)
		return cache.Pair[[]SeriesCredit]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// GetSeriesGenres returns genres for a series with caching (10 min TTL).
func (s *CachedService) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error) {
	return cache.Get(ctx, s.cache, cache.TVShowGenresKey(seriesID.String()), cache.TVShowMetaTTL, func(ctx context.Context) ([]SeriesGenre, error) {
		return s.Service.GetSeriesGenres(ctx, seriesID)
	})
}

// GetSeriesNetworks returns networks for a series with caching (10 min TTL).
func (s *CachedService) GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]Network, error) {
	return cache.Get(ctx, s.cache, cache.TVShowNetworksKey(seriesID.String()), cache.TVShowMetaTTL, func(ctx context.Context) ([]Network, error) {
		return s.Service.GetSeriesNetworks(ctx, seriesID)
	})
}

// =============================================================================
// Watch Progress (Cached Reads)
// =============================================================================

// GetContinueWatching returns series the user is currently watching with caching (1 min TTL).
func (s *CachedService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	key := fmt.Sprintf("%scontinue-watching:%s:%d", cache.KeyPrefixTVShow, userID.String(), limit)
	return cache.Get(ctx, s.cache, key, cache.ContinueWatchingTTL, func(ctx context.Context) ([]ContinueWatchingItem, error) {
		return s.Service.GetContinueWatching(ctx, userID, limit)
	})
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

	if err := s.cache.Delete(ctx, cache.TVShowKey(id)); err != nil {
		s.logger.Warn("failed to invalidate series cache", slog.Any("error", err))
	}

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

	pattern := cache.KeyPrefixTVShowEpisodes + "*"
	if err := s.cache.Invalidate(ctx, pattern); err != nil {
		s.logger.Warn("failed to invalidate episode caches", slog.Any("error", err))
	}

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
func (s *CachedService) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	return cache.Get(ctx, s.cache, cache.KeyPrefixTVShowGenres+"distinct", cache.TVShowMetaTTL, func(ctx context.Context) ([]content.GenreSummary, error) {
		return s.Service.ListDistinctGenres(ctx)
	})
}
