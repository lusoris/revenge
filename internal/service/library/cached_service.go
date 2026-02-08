package library

import (
	"context"

	"log/slog"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// CachedService wraps the library service with caching.
type CachedService struct {
	*Service
	cache  *cache.Cache
	logger *slog.Logger
}

// NewCachedService creates a new cached library service.
func NewCachedService(svc *Service, c *cache.Cache, logger *slog.Logger) *CachedService {
	return &CachedService{
		Service: svc,
		cache:   c,
		logger:  logger.With("component", "library-cache"),
	}
}

// Get retrieves a library by ID with caching.
func (s *CachedService) Get(ctx context.Context, id uuid.UUID) (*Library, error) {
	if s.cache == nil {
		return s.Service.Get(ctx, id)
	}

	cacheKey := cache.LibraryKey(id.String())

	var lib Library
	if err := s.cache.GetJSON(ctx, cacheKey, &lib); err == nil {
		s.logger.Debug("library cache hit", slog.String("id", id.String()))
		return &lib, nil
	}

	result, err := s.Service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache async
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.LibraryStatsTTL); setErr != nil {
			s.logger.Warn("failed to cache library", slog.Any("error",setErr))
		}
	}()

	return result, nil
}

// List retrieves all libraries with caching.
func (s *CachedService) List(ctx context.Context) ([]Library, error) {
	if s.cache == nil {
		return s.Service.List(ctx)
	}

	cacheKey := cache.KeyPrefixLibrary + "list"

	var libs []Library
	if err := s.cache.GetJSON(ctx, cacheKey, &libs); err == nil {
		s.logger.Debug("library list cache hit")
		return libs, nil
	}

	result, err := s.Service.List(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.LibraryStatsTTL); setErr != nil {
			s.logger.Warn("failed to cache library list", slog.Any("error",setErr))
		}
	}()

	return result, nil
}

// Count retrieves the total library count with caching.
func (s *CachedService) Count(ctx context.Context) (int64, error) {
	if s.cache == nil {
		return s.Service.Count(ctx)
	}

	cacheKey := cache.KeyPrefixLibrary + "count"

	var count int64
	if err := s.cache.GetJSON(ctx, cacheKey, &count); err == nil {
		s.logger.Debug("library count cache hit")
		return count, nil
	}

	result, err := s.Service.Count(ctx)
	if err != nil {
		return 0, err
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
		defer cancel()
		if setErr := s.cache.SetJSON(cacheCtx, cacheKey, result, cache.LibraryStatsTTL); setErr != nil {
			s.logger.Warn("failed to cache library count", slog.Any("error",setErr))
		}
	}()

	return result, nil
}

// Write operations - invalidate cache

// Create creates a library and invalidates list cache.
func (s *CachedService) Create(ctx context.Context, req CreateLibraryRequest) (*Library, error) {
	result, err := s.Service.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
			defer cancel()
			s.invalidateLibraryLists(cacheCtx)
		}()
	}

	return result, nil
}

// Update updates a library and invalidates cache.
func (s *CachedService) Update(ctx context.Context, id uuid.UUID, update *LibraryUpdate) (*Library, error) {
	result, err := s.Service.Update(ctx, id, update)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
			defer cancel()
			if err := s.cache.InvalidateLibrary(cacheCtx, id.String()); err != nil {
				s.logger.Warn("failed to invalidate library cache", slog.Any("error",err))
			}
			s.invalidateLibraryLists(cacheCtx)
		}()
	}

	return result, nil
}

// Delete deletes a library and invalidates cache.
func (s *CachedService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.Service.Delete(ctx, id)
	if err != nil {
		return err
	}

	if s.cache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
			defer cancel()
			if err := s.cache.InvalidateLibrary(cacheCtx, id.String()); err != nil {
				s.logger.Warn("failed to invalidate library cache", slog.Any("error",err))
			}
			s.invalidateLibraryLists(cacheCtx)
		}()
	}

	return nil
}

// CompleteScan completes a scan and invalidates library stats cache.
func (s *CachedService) CompleteScan(ctx context.Context, scanID uuid.UUID, progress *ScanProgress) (*LibraryScan, error) {
	result, err := s.Service.CompleteScan(ctx, scanID, progress)
	if err != nil {
		return nil, err
	}

	// Invalidate library stats after scan completes
	if s.cache != nil && result != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), cache.SessionTTL)
			defer cancel()
			if err := s.cache.InvalidateLibrary(cacheCtx, result.LibraryID.String()); err != nil {
				s.logger.Warn("failed to invalidate library cache after scan", slog.Any("error",err))
			}
		}()
	}

	return result, nil
}

// invalidateLibraryLists invalidates list and count caches.
func (s *CachedService) invalidateLibraryLists(ctx context.Context) {
	patterns := []string{
		cache.KeyPrefixLibrary + "list",
		cache.KeyPrefixLibrary + "count",
	}
	for _, key := range patterns {
		if err := s.cache.Delete(ctx, key); err != nil {
			s.logger.Warn("failed to invalidate library list cache", slog.String("key", key), slog.Any("error",err))
		}
	}
}

// InvalidateLibraryCache invalidates all cache entries for a library.
func (s *CachedService) InvalidateLibraryCache(ctx context.Context, libraryID uuid.UUID) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.InvalidateLibrary(ctx, libraryID.String())
}
