package search

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// CachedMovieSearchService wraps the search service with caching.
type CachedMovieSearchService struct {
	*MovieSearchService
	cache  *cache.Cache
	logger *zap.Logger
}

// NewCachedMovieSearchService creates a new cached search service.
func NewCachedMovieSearchService(svc *MovieSearchService, cache *cache.Cache, logger *zap.Logger) *CachedMovieSearchService {
	return &CachedMovieSearchService{
		MovieSearchService: svc,
		cache:              cache,
		logger:             logger.Named("search-cache"),
	}
}

// Search searches for movies with caching (30 sec TTL).
func (s *CachedMovieSearchService) Search(ctx context.Context, params SearchParams) (*SearchResult, error) {
	if s.cache == nil {
		return s.MovieSearchService.Search(ctx, params)
	}

	// Generate cache key from search params
	cacheKey := s.searchCacheKey(params)

	// Try cache first
	var result SearchResult
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("search cache hit", zap.String("query", params.Query))
		return &result, nil
	}

	s.logger.Debug("search cache miss", zap.String("query", params.Query))

	// Cache miss - execute search
	searchResult, err := s.MovieSearchService.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, searchResult, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache search result", zap.Error(setErr))
		}
	}()

	return searchResult, nil
}

// Autocomplete performs autocomplete with caching (30 sec TTL).
func (s *CachedMovieSearchService) Autocomplete(ctx context.Context, query string, limit int) ([]string, error) {
	if s.cache == nil {
		return s.MovieSearchService.Autocomplete(ctx, query, limit)
	}

	cacheKey := fmt.Sprintf("%sautocomplete:%s:%d", cache.KeyPrefixSearchMovies, query, limit)

	// Try cache first
	var results []string
	if err := s.cache.GetJSON(ctx, cacheKey, &results); err == nil {
		s.logger.Debug("autocomplete cache hit", zap.String("query", query))
		return results, nil
	}

	// Cache miss - execute autocomplete
	autocompleteResults, err := s.MovieSearchService.Autocomplete(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, autocompleteResults, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache autocomplete result", zap.Error(setErr))
		}
	}()

	return autocompleteResults, nil
}

// GetFacets returns available facets with caching (2 min TTL).
func (s *CachedMovieSearchService) GetFacets(ctx context.Context, facetFields []string) (map[string][]FacetValue, error) {
	if s.cache == nil {
		return s.MovieSearchService.GetFacets(ctx, facetFields)
	}

	cacheKey := cache.KeyPrefixSearchMovies + "facets"

	// Try cache first
	var facets map[string][]FacetValue
	if err := s.cache.GetJSON(ctx, cacheKey, &facets); err == nil {
		s.logger.Debug("facets cache hit")
		return facets, nil
	}

	// Cache miss - get facets
	facetResults, err := s.MovieSearchService.GetFacets(ctx, facetFields)
	if err != nil {
		return nil, err
	}

	// Cache the result async (longer TTL since facets change less frequently)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, facetResults, 2*time.Minute); setErr != nil {
			s.logger.Warn("failed to cache facets", zap.Error(setErr))
		}
	}()

	return facetResults, nil
}

// InvalidateSearchCache invalidates all search caches.
// Call this after indexing or removing movies.
func (s *CachedMovieSearchService) InvalidateSearchCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Debug("invalidating search cache")
	return s.cache.Invalidate(ctx, cache.KeyPrefixSearch+"*")
}

// searchCacheKey generates a cache key for search queries.
func (s *CachedMovieSearchService) searchCacheKey(params SearchParams) string {
	// Create a deterministic key from params
	key := fmt.Sprintf("q:%s|f:%s|s:%s|p:%d|pp:%d|fb:%v",
		params.Query,
		params.FilterBy,
		params.SortBy,
		params.Page,
		params.PerPage,
		params.FacetBy,
	)

	// Hash the key to keep it short
	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixSearchMovies + hex.EncodeToString(hash[:8])
}
