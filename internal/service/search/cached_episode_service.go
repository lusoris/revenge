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

// CachedEpisodeSearchService wraps the episode search service with caching.
type CachedEpisodeSearchService struct {
	*EpisodeSearchService
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedEpisodeSearchService creates a new cached episode search service.
func NewCachedEpisodeSearchService(svc *EpisodeSearchService, cache *cache.Cache, logger *slog.Logger) *CachedEpisodeSearchService {
	return &CachedEpisodeSearchService{
		EpisodeSearchService: svc,
		cache:                cache,
		logger:               logger.With("component", "episode-search-cache"),
	}
}

// SearchEpisodes searches for episodes with caching (30 sec TTL).
func (s *CachedEpisodeSearchService) SearchEpisodes(ctx context.Context, params EpisodeSearchParams) (*EpisodeSearchResult, error) {
	if s.cache == nil {
		return s.EpisodeSearchService.SearchEpisodes(ctx, params)
	}

	cacheKey := s.searchCacheKey(params)

	var result EpisodeSearchResult
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("episode search cache hit", slog.String("query", params.Query))
		return &result, nil
	}

	s.logger.Debug("episode search cache miss", slog.String("query", params.Query))

	searchResult, err := s.EpisodeSearchService.SearchEpisodes(ctx, params)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, searchResult, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache episode search result", slog.Any("error", setErr))
		}
	}()

	return searchResult, nil
}

// AutocompleteEpisodes provides search suggestions for episode titles with caching (30 sec TTL).
func (s *CachedEpisodeSearchService) AutocompleteEpisodes(ctx context.Context, query string, limit int) ([]string, error) {
	if s.cache == nil {
		return s.EpisodeSearchService.AutocompleteEpisodes(ctx, query, limit)
	}

	cacheKey := fmt.Sprintf("%sepisodes:autocomplete:%s:%d", cache.KeyPrefixSearch, query, limit)

	var results []string
	if err := s.cache.GetJSON(ctx, cacheKey, &results); err == nil {
		s.logger.Debug("episode autocomplete cache hit", slog.String("query", query))
		return results, nil
	}

	autocompleteResults, err := s.EpisodeSearchService.AutocompleteEpisodes(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, autocompleteResults, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache episode autocomplete result", slog.Any("error", setErr))
		}
	}()

	return autocompleteResults, nil
}

// InvalidateSearchCache invalidates all episode search caches.
func (s *CachedEpisodeSearchService) InvalidateSearchCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Debug("invalidating episode search cache")
	return s.cache.Invalidate(ctx, cache.KeyPrefixSearch+"episodes:*")
}

// searchCacheKey generates a cache key for episode search queries.
func (s *CachedEpisodeSearchService) searchCacheKey(params EpisodeSearchParams) string {
	key := fmt.Sprintf("q:%s|f:%s|s:%s|p:%d|pp:%d|fb:%v",
		params.Query,
		params.FilterBy,
		params.SortBy,
		params.Page,
		params.PerPage,
		params.FacetBy,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixSearch + "episodes:" + hex.EncodeToString(hash[:8])
}
