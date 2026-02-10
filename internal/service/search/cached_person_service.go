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

// CachedPersonSearchService wraps the person search service with caching.
type CachedPersonSearchService struct {
	*PersonSearchService
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedPersonSearchService creates a new cached person search service.
func NewCachedPersonSearchService(svc *PersonSearchService, cache *cache.Cache, logger *slog.Logger) *CachedPersonSearchService {
	return &CachedPersonSearchService{
		PersonSearchService: svc,
		cache:               cache,
		logger:              logger.With("component", "person-search-cache"),
	}
}

// SearchPersons searches for people with caching (30 sec TTL).
func (s *CachedPersonSearchService) SearchPersons(ctx context.Context, params PersonSearchParams) (*PersonSearchResult, error) {
	if s.cache == nil {
		return s.PersonSearchService.SearchPersons(ctx, params)
	}

	cacheKey := s.searchCacheKey(params)

	var result PersonSearchResult
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("person search cache hit", slog.String("query", params.Query))
		return &result, nil
	}

	s.logger.Debug("person search cache miss", slog.String("query", params.Query))

	searchResult, err := s.PersonSearchService.SearchPersons(ctx, params)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, searchResult, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache person search result", slog.Any("error", setErr))
		}
	}()

	return searchResult, nil
}

// AutocompletePersons provides search suggestions for person names with caching (30 sec TTL).
func (s *CachedPersonSearchService) AutocompletePersons(ctx context.Context, query string, limit int) ([]string, error) {
	if s.cache == nil {
		return s.PersonSearchService.AutocompletePersons(ctx, query, limit)
	}

	cacheKey := fmt.Sprintf("%speople:autocomplete:%s:%d", cache.KeyPrefixSearch, query, limit)

	var results []string
	if err := s.cache.GetJSON(ctx, cacheKey, &results); err == nil {
		s.logger.Debug("person autocomplete cache hit", slog.String("query", query))
		return results, nil
	}

	autocompleteResults, err := s.PersonSearchService.AutocompletePersons(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, autocompleteResults, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache person autocomplete result", slog.Any("error", setErr))
		}
	}()

	return autocompleteResults, nil
}

// InvalidateSearchCache invalidates all person search caches.
func (s *CachedPersonSearchService) InvalidateSearchCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Debug("invalidating person search cache")
	return s.cache.Invalidate(ctx, cache.KeyPrefixSearch+"people:*")
}

// searchCacheKey generates a cache key for person search queries.
func (s *CachedPersonSearchService) searchCacheKey(params PersonSearchParams) string {
	key := fmt.Sprintf("q:%s|f:%s|s:%s|p:%d|pp:%d|fb:%v",
		params.Query,
		params.FilterBy,
		params.SortBy,
		params.Page,
		params.PerPage,
		params.FacetBy,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixSearch + "people:" + hex.EncodeToString(hash[:8])
}
