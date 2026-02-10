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

// CachedTVShowSearchService wraps the TV show search service with caching.
type CachedTVShowSearchService struct {
	*TVShowSearchService
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedTVShowSearchService creates a new cached TV show search service.
func NewCachedTVShowSearchService(svc *TVShowSearchService, cache *cache.Cache, logger *slog.Logger) *CachedTVShowSearchService {
	return &CachedTVShowSearchService{
		TVShowSearchService: svc,
		cache:               cache,
		logger:              logger.With("component", "tvshow-search-cache"),
	}
}

// SearchSeries searches for TV shows with caching (30 sec TTL).
func (s *CachedTVShowSearchService) SearchSeries(ctx context.Context, params TVShowSearchParams) (*TVShowSearchResult, error) {
	if s.cache == nil {
		return s.TVShowSearchService.SearchSeries(ctx, params)
	}

	cacheKey := s.searchCacheKey(params)

	var result TVShowSearchResult
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("tvshow search cache hit", slog.String("query", params.Query))
		return &result, nil
	}

	s.logger.Debug("tvshow search cache miss", slog.String("query", params.Query))

	searchResult, err := s.TVShowSearchService.SearchSeries(ctx, params)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, searchResult, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache tvshow search result", slog.Any("error", setErr))
		}
	}()

	return searchResult, nil
}

// AutocompleteSeries provides search suggestions for series titles with caching (30 sec TTL).
func (s *CachedTVShowSearchService) AutocompleteSeries(ctx context.Context, query string, limit int) ([]string, error) {
	if s.cache == nil {
		return s.TVShowSearchService.AutocompleteSeries(ctx, query, limit)
	}

	cacheKey := fmt.Sprintf("%sautocomplete:%s:%d", cache.KeyPrefixSearchTVShows, query, limit)

	var results []string
	if err := s.cache.GetJSON(ctx, cacheKey, &results); err == nil {
		s.logger.Debug("tvshow autocomplete cache hit", slog.String("query", query))
		return results, nil
	}

	autocompleteResults, err := s.TVShowSearchService.AutocompleteSeries(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, autocompleteResults, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache tvshow autocomplete result", slog.Any("error", setErr))
		}
	}()

	return autocompleteResults, nil
}

// GetFacets returns available facets with caching (2 min TTL).
func (s *CachedTVShowSearchService) GetFacets(ctx context.Context, facetNames []string) (map[string][]FacetValue, error) {
	if s.cache == nil {
		return s.TVShowSearchService.GetFacets(ctx, facetNames)
	}

	cacheKey := cache.KeyPrefixSearchTVShows + "facets"

	var facets map[string][]FacetValue
	if err := s.cache.GetJSON(ctx, cacheKey, &facets); err == nil {
		s.logger.Debug("tvshow facets cache hit")
		return facets, nil
	}

	facetResults, err := s.TVShowSearchService.GetFacets(ctx, facetNames)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, facetResults, 2*time.Minute); setErr != nil {
			s.logger.Warn("failed to cache tvshow facets", slog.Any("error", setErr))
		}
	}()

	return facetResults, nil
}

// InvalidateSearchCache invalidates all TV show search caches.
func (s *CachedTVShowSearchService) InvalidateSearchCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Debug("invalidating tvshow search cache")
	return s.cache.Invalidate(ctx, cache.KeyPrefixSearchTVShows+"*")
}

// searchCacheKey generates a cache key for TV show search queries.
func (s *CachedTVShowSearchService) searchCacheKey(params TVShowSearchParams) string {
	key := fmt.Sprintf("q:%s|f:%s|s:%s|p:%d|pp:%d|fb:%v",
		params.Query,
		params.FilterBy,
		params.SortBy,
		params.Page,
		params.PerPage,
		params.FacetBy,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixSearchTVShows + hex.EncodeToString(hash[:8])
}
