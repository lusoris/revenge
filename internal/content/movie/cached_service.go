// Package movie provides movie-related business logic.
package movie

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"log/slog"

	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/infra/cache"
)

// CachedService wraps the movie service with caching.
type CachedService struct {
	Service
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedService creates a new cached movie service.
func NewCachedService(svc Service, cache *cache.Cache, logger *slog.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   cache,
		logger:  logger.With("component", "movie-cache"),
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
		s.logger.Debug("movie cache hit", slog.String("id", id.String()))
		return &movie, nil
	}

	s.logger.Debug("movie cache miss", slog.String("id", id.String()))

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
			s.logger.Warn("failed to cache movie", slog.Any("error",setErr))
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
		s.logger.Debug("list movies cache hit", slog.String("key", cacheKey))
		return movies, nil
	}

	s.logger.Debug("list movies cache miss", slog.String("key", cacheKey))

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
			s.logger.Warn("failed to cache movie list", slog.Any("error",setErr))
		}
	}()

	return result, nil
}

// ListRecentlyAdded returns recently added movies with caching (2 min TTL).
// Cache key includes pagination params, caches both items and total.
func (s *CachedService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, int64, error) {
	if s.cache == nil {
		return s.Service.ListRecentlyAdded(ctx, limit, offset)
	}

	cacheKey := fmt.Sprintf("%srecently-added:%d:%d", cache.KeyPrefixMovie, limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedMovies struct {
		Items []Movie `json:"items"`
		Total int64   `json:"total"`
	}
	var cached cachedMovies
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("recently added cache hit")
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
		toCache := cachedMovies{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, 2*time.Minute); setErr != nil {
			s.logger.Warn("failed to cache recently added", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
}

// ListTopRated returns top-rated movies with caching (5 min TTL).
// Cache key includes pagination params, caches both items and total.
func (s *CachedService) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, int64, error) {
	if s.cache == nil {
		return s.Service.ListTopRated(ctx, minVotes, limit, offset)
	}

	cacheKey := fmt.Sprintf("%stop-rated:%d:%d:%d", cache.KeyPrefixMovie, minVotes, limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedMovies struct {
		Items []Movie `json:"items"`
		Total int64   `json:"total"`
	}
	var cached cachedMovies
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("top rated cache hit")
		return cached.Items, cached.Total, nil
	}

	// Cache miss - load from database
	items, total, err := s.Service.ListTopRated(ctx, minVotes, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		toCache := cachedMovies{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, cache.MovieTTL); setErr != nil {
			s.logger.Warn("failed to cache top rated", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
}

// GetMovieCast returns the cast for a movie with caching (10 min TTL).
// Cache key includes pagination params.
func (s *CachedService) GetMovieCast(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error) {
	if s.cache == nil {
		return s.Service.GetMovieCast(ctx, movieID, limit, offset)
	}

	cacheKey := fmt.Sprintf("%s%s:cast:%d:%d", cache.KeyPrefixMovie, movieID.String(), limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedCredits struct {
		Items []MovieCredit `json:"items"`
		Total int64         `json:"total"`
	}
	var cached cachedCredits
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("movie cast cache hit", slog.String("movie_id", movieID.String()))
		return cached.Items, cached.Total, nil
	}

	// Cache miss - load from database
	items, total, err := s.Service.GetMovieCast(ctx, movieID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		toCache := cachedCredits{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache movie cast", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
}

// GetMovieCrew returns the crew for a movie with caching (10 min TTL).
// Cache key includes pagination params.
func (s *CachedService) GetMovieCrew(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error) {
	if s.cache == nil {
		return s.Service.GetMovieCrew(ctx, movieID, limit, offset)
	}

	cacheKey := fmt.Sprintf("%s%s:crew:%d:%d", cache.KeyPrefixMovie, movieID.String(), limit, offset)

	// Try cache first - we cache both items and total count together
	type cachedCredits struct {
		Items []MovieCredit `json:"items"`
		Total int64         `json:"total"`
	}
	var cached cachedCredits
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("movie crew cache hit", slog.String("movie_id", movieID.String()))
		return cached.Items, cached.Total, nil
	}

	// Cache miss - load from database
	items, total, err := s.Service.GetMovieCrew(ctx, movieID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		toCache := cachedCredits{Items: items, Total: total}
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, toCache, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache movie crew", slog.Any("error", setErr))
		}
	}()

	return items, total, nil
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
		s.logger.Debug("movie genres cache hit", slog.String("movie_id", movieID.String()))
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
			s.logger.Warn("failed to cache movie genres", slog.Any("error",setErr))
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
		s.logger.Debug("collection cache hit", slog.String("id", id.String()))
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
			s.logger.Warn("failed to cache collection", slog.Any("error",setErr))
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
		s.logger.Debug("continue watching cache hit", slog.String("user_id", userID.String()))
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
			s.logger.Warn("failed to cache continue watching", slog.Any("error",setErr))
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
			s.logger.Warn("failed to invalidate continue watching cache", slog.Any("error",invErr))
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
			s.logger.Warn("failed to invalidate continue watching cache", slog.Any("error",invErr))
		}
	}

	return nil
}

// invalidateMovie invalidates all cache entries for a movie.
func (s *CachedService) invalidateMovie(ctx context.Context, movieID uuid.UUID) {
	// Delete movie itself
	if err := s.cache.Delete(ctx, cache.MovieKey(movieID.String())); err != nil {
		s.logger.Warn("failed to invalidate movie cache", slog.Any("error",err))
	}

	// Delete movie metadata (cast, crew, genres)
	pattern := fmt.Sprintf("%s%s:*", cache.KeyPrefixMovie, movieID.String())
	if err := s.cache.Invalidate(ctx, pattern); err != nil {
		s.logger.Warn("failed to invalidate movie metadata cache", slog.Any("error",err))
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
			s.logger.Warn("failed to invalidate movie list cache", slog.String("pattern", p), slog.Any("error",err))
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

// GetMovieFiles returns files for a movie with caching (5 min TTL).
// This is called on every playback start to locate the media file.
func (s *CachedService) GetMovieFiles(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	if s.cache == nil {
		return s.Service.GetMovieFiles(ctx, movieID)
	}

	cacheKey := cache.MovieFilesKey(movieID.String())

	// Try cache first
	var files []MovieFile
	if err := s.cache.GetJSON(ctx, cacheKey, &files); err == nil {
		s.logger.Debug("movie files cache hit", slog.String("movie_id", movieID.String()))
		return files, nil
	}

	// Cache miss - load from database
	result, err := s.Service.GetMovieFiles(ctx, movieID)
	if err != nil {
		return nil, err
	}

	// Cache the result async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieTTL); setErr != nil {
			s.logger.Warn("failed to cache movie files", slog.Any("error", setErr))
		}
	}()

	return result, nil
}

// ListDistinctGenres returns all distinct movie genres with caching (10 min TTL).
// Genre lists are near-static and frequently requested for filter UIs.
func (s *CachedService) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	if s.cache == nil {
		return s.Service.ListDistinctGenres(ctx)
	}

	cacheKey := cache.KeyPrefixMovieGenres + "distinct"

	// Try cache first
	var genres []content.GenreSummary
	if err := s.cache.GetJSON(ctx, cacheKey, &genres); err == nil {
		s.logger.Debug("movie distinct genres cache hit")
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
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.MovieMetaTTL); setErr != nil {
			s.logger.Warn("failed to cache distinct movie genres", slog.Any("error", setErr))
		}
	}()

	return result, nil
}
