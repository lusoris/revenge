// Package movie provides movie-related business logic.
package movie

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// CachedService wraps the movie service with caching.
type CachedService struct {
	Service
	cache  *cache.Cache
	logger *zap.Logger
}

// NewCachedService creates a new cached movie service.
func NewCachedService(svc Service, cache *cache.Cache, logger *zap.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   cache,
		logger:  logger.Named("movie-cache"),
	}
}

// GetMovie retrieves a movie by ID with caching (5 min TTL).
func (s *CachedService) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
	if s.cache == nil {
		return s.Service.GetMovie(ctx, id)
	}

	cacheKey := cache.MovieKey(id.String())

	// Try cache first
	var movie Movie
	if err := s.cache.GetJSON(ctx, cacheKey, &movie); err == nil {
		s.logger.Debug("movie cache hit", zap.String("id", id.String()))
		return &movie, nil
	}

	s.logger.Debug("movie cache miss", zap.String("id", id.String()))

	// Cache miss - load from database
	result, err := s.Service.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieTTL); setErr != nil {
			s.logger.Warn("failed to cache movie", zap.Error(setErr))
		}
	}()

	return result, nil
}

// ListMovies returns a paginated list of movies with caching (1 min TTL).
func (s *CachedService) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	if s.cache == nil {
		return s.Service.ListMovies(ctx, filters)
	}

	// Generate cache key from filters
	cacheKey := s.listMoviesCacheKey(filters)

	// Try cache first
	var movies []Movie
	if err := s.cache.GetJSON(ctx, cacheKey, &movies); err == nil {
		s.logger.Debug("list movies cache hit", zap.String("key", cacheKey))
		return movies, nil
	}

	s.logger.Debug("list movies cache miss", zap.String("key", cacheKey))

	// Cache miss - load from database
	result, err := s.Service.ListMovies(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Cache the result async (shorter TTL for lists)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, time.Minute); setErr != nil {
			s.logger.Warn("failed to cache movie list", zap.Error(setErr))
		}
	}()

	return result, nil
}

// ListRecentlyAdded returns recently added movies with caching (2 min TTL).
func (s *CachedService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, error) {
	if s.cache == nil {
		return s.Service.ListRecentlyAdded(ctx, limit, offset)
	}

	cacheKey := fmt.Sprintf("%srecently-added:%d:%d", cache.KeyPrefixMovie, limit, offset)

	// Try cache first
	var movies []Movie
	if err := s.cache.GetJSON(ctx, cacheKey, &movies); err == nil {
		s.logger.Debug("recently added cache hit")
		return movies, nil
	}

	// Cache miss - load from database
	result, err := s.Service.ListRecentlyAdded(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, 2*time.Minute); setErr != nil {
			s.logger.Warn("failed to cache recently added", zap.Error(setErr))
		}
	}()

	return result, nil
}

// ListTopRated returns top-rated movies with caching (5 min TTL).
func (s *CachedService) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, error) {
	if s.cache == nil {
		return s.Service.ListTopRated(ctx, minVotes, limit, offset)
	}

	cacheKey := fmt.Sprintf("%stop-rated:%d:%d:%d", cache.KeyPrefixMovie, minVotes, limit, offset)

	// Try cache first
	var movies []Movie
	if err := s.cache.GetJSON(ctx, cacheKey, &movies); err == nil {
		s.logger.Debug("top rated cache hit")
		return movies, nil
	}

	// Cache miss - load from database
	result, err := s.Service.ListTopRated(ctx, minVotes, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieTTL); setErr != nil {
			s.logger.Warn("failed to cache top rated", zap.Error(setErr))
		}
	}()

	return result, nil
}

// GetMovieCast returns the cast for a movie with caching (10 min TTL).
func (s *CachedService) GetMovieCast(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	if s.cache == nil {
		return s.Service.GetMovieCast(ctx, movieID)
	}

	cacheKey := fmt.Sprintf("%s%s:cast", cache.KeyPrefixMovie, movieID.String())

	// Try cache first
	var cast []MovieCredit
	if err := s.cache.GetJSON(ctx, cacheKey, &cast); err == nil {
		s.logger.Debug("movie cast cache hit", zap.String("movie_id", movieID.String()))
		return cast, nil
	}

	// Cache miss - load from database
	result, err := s.Service.GetMovieCast(ctx, movieID)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache movie cast", zap.Error(setErr))
		}
	}()

	return result, nil
}

// GetMovieCrew returns the crew for a movie with caching (10 min TTL).
func (s *CachedService) GetMovieCrew(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	if s.cache == nil {
		return s.Service.GetMovieCrew(ctx, movieID)
	}

	cacheKey := fmt.Sprintf("%s%s:crew", cache.KeyPrefixMovie, movieID.String())

	// Try cache first
	var crew []MovieCredit
	if err := s.cache.GetJSON(ctx, cacheKey, &crew); err == nil {
		s.logger.Debug("movie crew cache hit", zap.String("movie_id", movieID.String()))
		return crew, nil
	}

	// Cache miss - load from database
	result, err := s.Service.GetMovieCrew(ctx, movieID)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache movie crew", zap.Error(setErr))
		}
	}()

	return result, nil
}

// GetMovieGenres returns genres for a movie with caching (10 min TTL).
func (s *CachedService) GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	if s.cache == nil {
		return s.Service.GetMovieGenres(ctx, movieID)
	}

	cacheKey := fmt.Sprintf("%s%s:genres", cache.KeyPrefixMovie, movieID.String())

	// Try cache first
	var genres []MovieGenre
	if err := s.cache.GetJSON(ctx, cacheKey, &genres); err == nil {
		s.logger.Debug("movie genres cache hit", zap.String("movie_id", movieID.String()))
		return genres, nil
	}

	// Cache miss - load from database
	result, err := s.Service.GetMovieGenres(ctx, movieID)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache movie genres", zap.Error(setErr))
		}
	}()

	return result, nil
}

// GetMovieCollection retrieves a collection by ID with caching (10 min TTL).
func (s *CachedService) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	if s.cache == nil {
		return s.Service.GetMovieCollection(ctx, id)
	}

	cacheKey := fmt.Sprintf("%scollection:%s", cache.KeyPrefixMovie, id.String())

	// Try cache first
	var collection MovieCollection
	if err := s.cache.GetJSON(ctx, cacheKey, &collection); err == nil {
		s.logger.Debug("collection cache hit", zap.String("id", id.String()))
		return &collection, nil
	}

	// Cache miss - load from database
	result, err := s.Service.GetMovieCollection(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache collection", zap.Error(setErr))
		}
	}()

	return result, nil
}

// GetContinueWatching returns movies the user is currently watching with caching (1 min TTL).
func (s *CachedService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	if s.cache == nil {
		return s.Service.GetContinueWatching(ctx, userID, limit)
	}

	cacheKey := fmt.Sprintf("%scontinue-watching:%s:%d", cache.KeyPrefixMovie, userID.String(), limit)

	// Try cache first
	var items []ContinueWatchingItem
	if err := s.cache.GetJSON(ctx, cacheKey, &items); err == nil {
		s.logger.Debug("continue watching cache hit", zap.String("user_id", userID.String()))
		return items, nil
	}

	// Cache miss - load from database
	result, err := s.Service.GetContinueWatching(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	// Cache the result async (short TTL since watch progress changes frequently)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, time.Minute); setErr != nil {
			s.logger.Warn("failed to cache continue watching", zap.Error(setErr))
		}
	}()

	return result, nil
}

// UpdateMovie updates an existing movie and invalidates cache.
func (s *CachedService) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	result, err := s.Service.UpdateMovie(ctx, params)
	if err != nil {
		return nil, err
	}

	// Invalidate movie cache
	if s.cache != nil {
		s.invalidateMovie(ctx, params.ID)
	}

	return result, nil
}

// DeleteMovie soft-deletes a movie and invalidates cache.
func (s *CachedService) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	if err := s.Service.DeleteMovie(ctx, id); err != nil {
		return err
	}

	// Invalidate movie cache
	if s.cache != nil {
		s.invalidateMovie(ctx, id)
	}

	return nil
}

// UpdateWatchProgress updates watch progress and invalidates continue watching cache.
func (s *CachedService) UpdateWatchProgress(ctx context.Context, userID, movieID uuid.UUID, progressSeconds, durationSeconds int32) (*MovieWatched, error) {
	result, err := s.Service.UpdateWatchProgress(ctx, userID, movieID, progressSeconds, durationSeconds)
	if err != nil {
		return nil, err
	}

	// Invalidate continue watching cache for user
	if s.cache != nil {
		pattern := fmt.Sprintf("%scontinue-watching:%s:*", cache.KeyPrefixMovie, userID.String())
		if invErr := s.cache.Invalidate(ctx, pattern); invErr != nil {
			s.logger.Warn("failed to invalidate continue watching cache", zap.Error(invErr))
		}
	}

	return result, nil
}

// MarkAsWatched marks a movie as watched and invalidates continue watching cache.
func (s *CachedService) MarkAsWatched(ctx context.Context, userID, movieID uuid.UUID) error {
	if err := s.Service.MarkAsWatched(ctx, userID, movieID); err != nil {
		return err
	}

	// Invalidate continue watching cache for user
	if s.cache != nil {
		pattern := fmt.Sprintf("%scontinue-watching:%s:*", cache.KeyPrefixMovie, userID.String())
		if invErr := s.cache.Invalidate(ctx, pattern); invErr != nil {
			s.logger.Warn("failed to invalidate continue watching cache", zap.Error(invErr))
		}
	}

	return nil
}

// invalidateMovie invalidates all cache entries for a movie.
func (s *CachedService) invalidateMovie(ctx context.Context, movieID uuid.UUID) {
	// Delete movie itself
	if err := s.cache.Delete(ctx, cache.MovieKey(movieID.String())); err != nil {
		s.logger.Warn("failed to invalidate movie cache", zap.Error(err))
	}

	// Delete movie metadata (cast, crew, genres)
	pattern := fmt.Sprintf("%s%s:*", cache.KeyPrefixMovie, movieID.String())
	if err := s.cache.Invalidate(ctx, pattern); err != nil {
		s.logger.Warn("failed to invalidate movie metadata cache", zap.Error(err))
	}

	// Invalidate list caches (recently-added, top-rated)
	// These will be refreshed on next request
	patterns := []string{
		cache.KeyPrefixMovie + "recently-added:*",
		cache.KeyPrefixMovie + "top-rated:*",
		cache.KeyPrefixMovie + "list:*",
	}
	for _, p := range patterns {
		if err := s.cache.Invalidate(ctx, p); err != nil {
			s.logger.Warn("failed to invalidate movie list cache", zap.String("pattern", p), zap.Error(err))
		}
	}
}

// InvalidateAllMovies invalidates all movie caches.
func (s *CachedService) InvalidateAllMovies(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Invalidate(ctx, cache.KeyPrefixMovie+"*")
}

// listMoviesCacheKey generates a cache key for movie list queries.
func (s *CachedService) listMoviesCacheKey(filters ListFilters) string {
	// Create a deterministic key from filters
	key := fmt.Sprintf("%s:%d:%d",
		filters.OrderBy,
		filters.Limit,
		filters.Offset,
	)

	// Hash the key to keep it short
	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixMovie + "list:" + hex.EncodeToString(hash[:8])
}
