package search

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// CachedSeasonSearchService wraps the season search service with caching.
type CachedSeasonSearchService struct {
	*SeasonSearchService
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedSeasonSearchService creates a new cached season search service.
func NewCachedSeasonSearchService(svc *SeasonSearchService, cache *cache.Cache, logger *slog.Logger) *CachedSeasonSearchService {
	return &CachedSeasonSearchService{
		SeasonSearchService: svc,
		cache:               cache,
		logger:              logger.With("component", "season-search-cache"),
	}
}

// SearchSeasons searches for seasons with caching (30 sec TTL).
func (s *CachedSeasonSearchService) SearchSeasons(ctx context.Context, params SeasonSearchParams) (*SeasonSearchResult, error) {
	if s.cache == nil {
		return s.SeasonSearchService.SearchSeasons(ctx, params)
	}

	cacheKey := s.searchCacheKey(params)

	var result SeasonSearchResult
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("season search cache hit", slog.String("query", params.Query))
		return &result, nil
	}

	s.logger.Debug("season search cache miss", slog.String("query", params.Query))

	searchResult, err := s.SeasonSearchService.SearchSeasons(ctx, params)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, searchResult, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache season search result", slog.Any("error", setErr))
		}
	}()

	return searchResult, nil
}

// AutocompleteSeasons provides search suggestions for season names with caching (30 sec TTL).
func (s *CachedSeasonSearchService) AutocompleteSeasons(ctx context.Context, query string, limit int) ([]string, error) {
	if s.cache == nil {
		return s.SeasonSearchService.AutocompleteSeasons(ctx, query, limit)
	}

	cacheKey := fmt.Sprintf("%sseasons:autocomplete:%s:%d", cache.KeyPrefixSearch, query, limit)

	var results []string
	if err := s.cache.GetJSON(ctx, cacheKey, &results); err == nil {
		s.logger.Debug("season autocomplete cache hit", slog.String("query", query))
		return results, nil
	}

	autocompleteResults, err := s.SeasonSearchService.AutocompleteSeasons(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, autocompleteResults, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache season autocomplete result", slog.Any("error", setErr))
		}
	}()

	return autocompleteResults, nil
}

// InvalidateSearchCache invalidates all season search caches.
func (s *CachedSeasonSearchService) InvalidateSearchCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Debug("invalidating season search cache")
	return s.cache.Invalidate(ctx, cache.KeyPrefixSearch+"seasons:*")
}

// searchCacheKey generates a cache key for season search queries.
func (s *CachedSeasonSearchService) searchCacheKey(params SeasonSearchParams) string {
	key := fmt.Sprintf("q:%s|f:%s|s:%s|p:%d|pp:%d|fb:%v",
		params.Query,
		params.FilterBy,
		params.SortBy,
		params.Page,
		params.PerPage,
		params.FacetBy,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixSearch + "seasons:" + hex.EncodeToString(hash[:8])
}
