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

// CachedUserSearchService wraps the user search service with caching.
type CachedUserSearchService struct {
	*UserSearchService
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedUserSearchService creates a new cached user search service.
func NewCachedUserSearchService(svc *UserSearchService, cache *cache.Cache, logger *slog.Logger) *CachedUserSearchService {
	return &CachedUserSearchService{
		UserSearchService: svc,
		cache:             cache,
		logger:            logger.With("component", "user-search-cache"),
	}
}

// SearchUsers searches for users with caching (30 sec TTL).
func (s *CachedUserSearchService) SearchUsers(ctx context.Context, params UserSearchParams) (*UserSearchResult, error) {
	if s.cache == nil {
		return s.UserSearchService.SearchUsers(ctx, params)
	}

	cacheKey := s.searchCacheKey(params)

	var result UserSearchResult
	if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil {
		s.logger.Debug("user search cache hit", slog.String("query", params.Query))
		return &result, nil
	}

	s.logger.Debug("user search cache miss", slog.String("query", params.Query))

	searchResult, err := s.UserSearchService.SearchUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, searchResult, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache user search result", slog.Any("error", setErr))
		}
	}()

	return searchResult, nil
}

// AutocompleteUsers provides search suggestions for usernames with caching (30 sec TTL).
func (s *CachedUserSearchService) AutocompleteUsers(ctx context.Context, query string, limit int) ([]string, error) {
	if s.cache == nil {
		return s.UserSearchService.AutocompleteUsers(ctx, query, limit)
	}

	cacheKey := fmt.Sprintf("%susers:autocomplete:%s:%d", cache.KeyPrefixSearch, query, limit)

	var results []string
	if err := s.cache.GetJSON(ctx, cacheKey, &results); err == nil {
		s.logger.Debug("user autocomplete cache hit", slog.String("query", query))
		return results, nil
	}

	autocompleteResults, err := s.UserSearchService.AutocompleteUsers(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, autocompleteResults, cache.SearchResultsTTL); setErr != nil {
			s.logger.Warn("failed to cache user autocomplete result", slog.Any("error", setErr))
		}
	}()

	return autocompleteResults, nil
}

// InvalidateSearchCache invalidates all user search caches.
func (s *CachedUserSearchService) InvalidateSearchCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	s.logger.Debug("invalidating user search cache")
	return s.cache.Invalidate(ctx, cache.KeyPrefixSearch+"users:*")
}

// searchCacheKey generates a cache key for user search queries.
func (s *CachedUserSearchService) searchCacheKey(params UserSearchParams) string {
	key := fmt.Sprintf("q:%s|f:%s|s:%s|p:%d|pp:%d|fb:%v",
		params.Query,
		params.FilterBy,
		params.SortBy,
		params.Page,
		params.PerPage,
		params.FacetBy,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixSearch + "users:" + hex.EncodeToString(hash[:8])
}
