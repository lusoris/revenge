// Package movie provides movie-related business logic.
package movie

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
	return cache.Get(ctx, s.cache, cache.MovieKey(id.String()), cache.MovieTTL, func(ctx context.Context) (*Movie, error) {
		return s.Service.GetMovie(ctx, id)
	})
}

// ListMovies returns a paginated list of movies with caching (1 min TTL).
func (s *CachedService) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	return cache.Get(ctx, s.cache, s.listMoviesCacheKey(filters), time.Minute, func(ctx context.Context) ([]Movie, error) {
		return s.Service.ListMovies(ctx, filters)
	})
}

// ListRecentlyAdded returns recently added movies with caching (2 min TTL).
func (s *CachedService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, int64, error) {
	key := fmt.Sprintf("%srecently-added:%d:%d", cache.KeyPrefixMovie, limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.RecentlyAddedTTL, func(ctx context.Context) (cache.Pair[[]Movie], error) {
		items, total, err := s.Service.ListRecentlyAdded(ctx, limit, offset)
		return cache.Pair[[]Movie]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// ListTopRated returns top-rated movies with caching (5 min TTL).
func (s *CachedService) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, int64, error) {
	key := fmt.Sprintf("%stop-rated:%d:%d:%d", cache.KeyPrefixMovie, minVotes, limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.TopRatedTTL, func(ctx context.Context) (cache.Pair[[]Movie], error) {
		items, total, err := s.Service.ListTopRated(ctx, minVotes, limit, offset)
		return cache.Pair[[]Movie]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// GetMovieCast returns the cast for a movie with caching (10 min TTL).
func (s *CachedService) GetMovieCast(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error) {
	key := fmt.Sprintf("%s%s:cast:%d:%d", cache.KeyPrefixMovie, movieID.String(), limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.MovieMetaTTL, func(ctx context.Context) (cache.Pair[[]MovieCredit], error) {
		items, total, err := s.Service.GetMovieCast(ctx, movieID, limit, offset)
		return cache.Pair[[]MovieCredit]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// GetMovieCrew returns the crew for a movie with caching (10 min TTL).
func (s *CachedService) GetMovieCrew(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error) {
	key := fmt.Sprintf("%s%s:crew:%d:%d", cache.KeyPrefixMovie, movieID.String(), limit, offset)
	result, err := cache.Get(ctx, s.cache, key, cache.MovieMetaTTL, func(ctx context.Context) (cache.Pair[[]MovieCredit], error) {
		items, total, err := s.Service.GetMovieCrew(ctx, movieID, limit, offset)
		return cache.Pair[[]MovieCredit]{Items: items, Total: total}, err
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Items, result.Total, nil
}

// GetMovieGenres returns genres for a movie with caching (10 min TTL).
func (s *CachedService) GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	key := fmt.Sprintf("%s%s:genres", cache.KeyPrefixMovie, movieID.String())
	return cache.Get(ctx, s.cache, key, cache.MovieMetaTTL, func(ctx context.Context) ([]MovieGenre, error) {
		return s.Service.GetMovieGenres(ctx, movieID)
	})
}

// GetMovieCollection retrieves a collection by ID with caching (10 min TTL).
func (s *CachedService) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	key := fmt.Sprintf("%scollection:%s", cache.KeyPrefixMovie, id.String())
	return cache.Get(ctx, s.cache, key, cache.MovieMetaTTL, func(ctx context.Context) (*MovieCollection, error) {
		return s.Service.GetMovieCollection(ctx, id)
	})
}

// GetContinueWatching returns movies the user is currently watching with caching (1 min TTL).
func (s *CachedService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	key := fmt.Sprintf("%scontinue-watching:%s:%d", cache.KeyPrefixMovie, userID.String(), limit)
	return cache.Get(ctx, s.cache, key, cache.ContinueWatchingTTL, func(ctx context.Context) ([]ContinueWatchingItem, error) {
		return s.Service.GetContinueWatching(ctx, userID, limit)
	})
}

// UpdateMovie updates an existing movie and invalidates cache.
func (s *CachedService) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	result, err := s.Service.UpdateMovie(ctx, params)
	if err != nil {
		return nil, err
	}

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

	if s.cache != nil {
		pattern := fmt.Sprintf("%scontinue-watching:%s:*", cache.KeyPrefixMovie, userID.String())
		if invErr := s.cache.Invalidate(ctx, pattern); invErr != nil {
			s.logger.Warn("failed to invalidate continue watching cache", slog.Any("error", invErr))
		}
	}

	return result, nil
}

// MarkAsWatched marks a movie as watched and invalidates continue watching cache.
func (s *CachedService) MarkAsWatched(ctx context.Context, userID, movieID uuid.UUID) error {
	if err := s.Service.MarkAsWatched(ctx, userID, movieID); err != nil {
		return err
	}

	if s.cache != nil {
		pattern := fmt.Sprintf("%scontinue-watching:%s:*", cache.KeyPrefixMovie, userID.String())
		if invErr := s.cache.Invalidate(ctx, pattern); invErr != nil {
			s.logger.Warn("failed to invalidate continue watching cache", slog.Any("error", invErr))
		}
	}

	return nil
}

// invalidateMovie invalidates all cache entries for a movie.
func (s *CachedService) invalidateMovie(ctx context.Context, movieID uuid.UUID) {
	if err := s.cache.Delete(ctx, cache.MovieKey(movieID.String())); err != nil {
		s.logger.Warn("failed to invalidate movie cache", slog.Any("error", err))
	}

	pattern := fmt.Sprintf("%s%s:*", cache.KeyPrefixMovie, movieID.String())
	if err := s.cache.Invalidate(ctx, pattern); err != nil {
		s.logger.Warn("failed to invalidate movie metadata cache", slog.Any("error", err))
	}

	patterns := []string{
		cache.KeyPrefixMovie + "recently-added:*",
		cache.KeyPrefixMovie + "top-rated:*",
		cache.KeyPrefixMovie + "list:*",
	}
	for _, p := range patterns {
		if err := s.cache.Invalidate(ctx, p); err != nil {
			s.logger.Warn("failed to invalidate movie list cache", slog.String("pattern", p), slog.Any("error", err))
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
	key := fmt.Sprintf("%s:%d:%d",
		filters.OrderBy,
		filters.Limit,
		filters.Offset,
	)

	hash := sha256.Sum256([]byte(key))
	return cache.KeyPrefixMovie + "list:" + hex.EncodeToString(hash[:8])
}

// GetMovieFiles returns files for a movie with caching (5 min TTL).
// This is called on every playback start to locate the media file.
func (s *CachedService) GetMovieFiles(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	return cache.Get(ctx, s.cache, cache.MovieFilesKey(movieID.String()), cache.MovieTTL, func(ctx context.Context) ([]MovieFile, error) {
		return s.Service.GetMovieFiles(ctx, movieID)
	})
}

// ListDistinctGenres returns all distinct movie genres with caching (10 min TTL).
// Genre lists are near-static and frequently requested for filter UIs.
func (s *CachedService) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	return cache.Get(ctx, s.cache, cache.KeyPrefixMovieGenres+"distinct", cache.MovieMetaTTL, func(ctx context.Context) ([]content.GenreSummary, error) {
		return s.Service.ListDistinctGenres(ctx)
	})
}
