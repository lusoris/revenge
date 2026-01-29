package movie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/maypok86/otter"

	"github.com/lusoris/revenge/pkg/resilience"
)

// Service errors.
var (
	ErrMovieNotFoundInService = errors.New("movie not found")
	ErrLibraryNotFound        = errors.New("library not found")
	ErrMetadataUnavailable    = errors.New("metadata provider unavailable")
)

// ServiceConfig holds service configuration.
type ServiceConfig struct {
	CacheMaxEntries int           `koanf:"cache_max_entries"`
	CacheTTL        time.Duration `koanf:"cache_ttl"`
}

// DefaultServiceConfig returns sensible defaults.
var DefaultServiceConfig = ServiceConfig{
	CacheMaxEntries: 10_000,
	CacheTTL:        5 * time.Minute,
}

// Service provides movie business logic with caching and resilience.
// Note: Movie data primarily comes from Servarr (Radarr) API calls.
// Metadata enrichment (TMDb) is done via background River jobs.
type Service struct {
	repo   Repository
	cache  otter.Cache[string, []byte]
	logger *slog.Logger
	config ServiceConfig
}

// NewService creates a new movie service.
func NewService(
	repo Repository,
	_ interface{}, // deprecated: metadata provider no longer needed here
	logger *slog.Logger,
	config ServiceConfig,
) (*Service, error) {
	if config.CacheMaxEntries == 0 {
		config.CacheMaxEntries = DefaultServiceConfig.CacheMaxEntries
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = DefaultServiceConfig.CacheTTL
	}

	// Create local cache for hot movie data
	cache, err := otter.MustBuilder[string, []byte](config.CacheMaxEntries).
		CollectStats().
		Cost(func(key string, value []byte) uint32 {
			return uint32(len(value))
		}).
		WithTTL(config.CacheTTL).
		Build()
	if err != nil {
		return nil, fmt.Errorf("create cache: %w", err)
	}

	return &Service{
		repo:   repo,
		cache:  cache,
		logger: logger.With("service", "movie"),
		config: config,
	}, nil
}

// Close releases service resources.
func (s *Service) Close() {
	s.cache.Close()
}

// GetMovie retrieves a movie by ID with caching.
func (s *Service) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
	key := fmt.Sprintf("movie:%s", id)

	// Check local cache first
	if cached, found := s.cache.Get(key); found {
		var movie Movie
		if err := json.Unmarshal(cached, &movie); err == nil {
			return &movie, nil
		}
		// Cache corrupted, remove it
		s.cache.Delete(key)
	}

	// Fetch from repository
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return nil, ErrMovieNotFoundInService
		}
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(movie); err == nil {
		s.cache.Set(key, data)
	}

	return movie, nil
}

// GetMovieWithRelations retrieves a movie with all related data.
func (s *Service) GetMovieWithRelations(ctx context.Context, id uuid.UUID) (*Movie, error) {
	movie, err := s.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load relations in parallel for performance
	errCh := make(chan error, 4)

	go func() {
		genres, err := s.repo.GetMovieGenres(ctx, id)
		if err == nil {
			movie.Genres = genres
		}
		errCh <- err
	}()

	go func() {
		cast, err := s.repo.GetMovieCast(ctx, id)
		if err == nil {
			movie.Cast = cast
		}
		errCh <- err
	}()

	go func() {
		crew, err := s.repo.GetMovieCrew(ctx, id)
		if err == nil {
			movie.Crew = crew
			// Extract directors and writers
			for _, c := range crew {
				switch c.Job {
				case "Director":
					movie.Directors = append(movie.Directors, c)
				case "Screenplay", "Writer":
					movie.Writers = append(movie.Writers, c)
				}
			}
		}
		errCh <- err
	}()

	go func() {
		images, err := s.repo.GetMovieImages(ctx, id)
		if err == nil {
			movie.Images = images
		}
		errCh <- err
	}()

	// Wait for all goroutines
	for i := 0; i < 4; i++ {
		if err := <-errCh; err != nil {
			s.logger.Warn("failed to load movie relation", "movie_id", id, "error", err)
		}
	}

	return movie, nil
}

// ListMovies retrieves movies with pagination.
func (s *Service) ListMovies(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, int64, error) {
	movies, err := s.repo.ListByLibrary(ctx, libraryID, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByLibrary(ctx, libraryID)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

// ListAllMovies retrieves all movies with pagination.
func (s *Service) ListAllMovies(ctx context.Context, params ListParams) ([]*Movie, int64, error) {
	movies, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

// SearchMovies searches movies by query.
func (s *Service) SearchMovies(ctx context.Context, query string, params ListParams) ([]*Movie, error) {
	return s.repo.Search(ctx, query, params)
}

// CreateMovie creates a new movie from a file path.
func (s *Service) CreateMovie(ctx context.Context, movie *Movie) error {
	if err := s.repo.Create(ctx, movie); err != nil {
		return err
	}

	s.logger.Info("movie created",
		"id", movie.ID,
		"title", movie.Title,
		"path", movie.Path,
	)

	return nil
}

// UpdateMovie updates a movie.
func (s *Service) UpdateMovie(ctx context.Context, movie *Movie) error {
	if err := s.repo.Update(ctx, movie); err != nil {
		return err
	}

	// Invalidate cache
	key := fmt.Sprintf("movie:%s", movie.ID)
	s.cache.Delete(key)

	return nil
}

// DeleteMovie deletes a movie.
func (s *Service) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	key := fmt.Sprintf("movie:%s", id)
	s.cache.Delete(key)

	s.logger.Info("movie deleted", "id", id)
	return nil
}

// ApplyMetadata updates movie fields from external metadata.
// This is called by the enrichment River job, not directly by the service.
func (s *Service) ApplyMetadata(ctx context.Context, id uuid.UUID, updates MovieMetadataUpdate) error {
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if updates.Title != "" {
		movie.Title = updates.Title
	}
	if updates.OriginalTitle != "" {
		movie.OriginalTitle = updates.OriginalTitle
	}
	if updates.Overview != "" {
		movie.Overview = updates.Overview
	}
	if updates.Tagline != "" {
		movie.Tagline = updates.Tagline
	}
	if updates.RuntimeTicks > 0 {
		movie.RuntimeTicks = updates.RuntimeTicks
	}
	if updates.ReleaseDate != nil {
		movie.ReleaseDate = updates.ReleaseDate
		movie.Year = updates.ReleaseDate.Year()
	}
	if updates.Budget > 0 {
		movie.Budget = updates.Budget
	}
	if updates.Revenue > 0 {
		movie.Revenue = updates.Revenue
	}
	if updates.CommunityRating > 0 {
		movie.CommunityRating = updates.CommunityRating
	}
	if updates.VoteCount > 0 {
		movie.VoteCount = updates.VoteCount
	}
	if updates.PosterPath != "" {
		movie.PosterPath = updates.PosterPath
	}
	if updates.BackdropPath != "" {
		movie.BackdropPath = updates.BackdropPath
	}
	if updates.TmdbID > 0 {
		movie.TmdbID = updates.TmdbID
	}
	if updates.ImdbID != "" {
		movie.ImdbID = updates.ImdbID
	}

	if err := s.repo.Update(ctx, movie); err != nil {
		return err
	}

	// Invalidate cache
	key := fmt.Sprintf("movie:%s", id)
	s.cache.Delete(key)

	s.logger.Info("movie metadata updated",
		"id", id,
		"title", movie.Title,
		"tmdb_id", movie.TmdbID,
	)

	return nil
}

// MovieMetadataUpdate contains fields that can be updated from metadata providers.
type MovieMetadataUpdate struct {
	Title           string
	OriginalTitle   string
	Overview        string
	Tagline         string
	RuntimeTicks    int64
	ReleaseDate     *time.Time
	Budget          int64
	Revenue         int64
	CommunityRating float64
	VoteCount       int
	PosterPath      string
	BackdropPath    string
	TmdbID          int
	ImdbID          string
}

// User Data Operations

// GetUserRating returns a user's rating for a movie.
func (s *Service) GetUserRating(ctx context.Context, userID, movieID uuid.UUID) (*UserRating, error) {
	return s.repo.GetUserRating(ctx, userID, movieID)
}

// SetUserRating sets a user's rating for a movie.
func (s *Service) SetUserRating(ctx context.Context, userID, movieID uuid.UUID, rating float64, review string) error {
	return s.repo.SetUserRating(ctx, userID, movieID, rating, review)
}

// DeleteUserRating removes a user's rating for a movie.
func (s *Service) DeleteUserRating(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.repo.DeleteUserRating(ctx, userID, movieID)
}

// Favorites

// IsFavorite checks if a movie is in user's favorites.
func (s *Service) IsFavorite(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	return s.repo.IsFavorite(ctx, userID, movieID)
}

// AddFavorite adds a movie to user's favorites.
func (s *Service) AddFavorite(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.repo.AddFavorite(ctx, userID, movieID)
}

// RemoveFavorite removes a movie from user's favorites.
func (s *Service) RemoveFavorite(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.repo.RemoveFavorite(ctx, userID, movieID)
}

// ListFavorites returns user's favorite movies.
func (s *Service) ListFavorites(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Movie, int64, error) {
	movies, err := s.repo.ListFavorites(ctx, userID, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountFavorites(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

// Watchlist

// IsInWatchlist checks if a movie is in user's watchlist.
func (s *Service) IsInWatchlist(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	return s.repo.IsInWatchlist(ctx, userID, movieID)
}

// AddToWatchlist adds a movie to user's watchlist.
func (s *Service) AddToWatchlist(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.repo.AddToWatchlist(ctx, userID, movieID)
}

// RemoveFromWatchlist removes a movie from user's watchlist.
func (s *Service) RemoveFromWatchlist(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.repo.RemoveFromWatchlist(ctx, userID, movieID)
}

// ListWatchlist returns user's watchlist.
func (s *Service) ListWatchlist(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Movie, int64, error) {
	movies, err := s.repo.ListWatchlist(ctx, userID, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountWatchlist(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

// Watch History

// GetWatchHistory returns current watch progress for a movie.
func (s *Service) GetWatchHistory(ctx context.Context, userID, movieID uuid.UUID) (*WatchHistory, error) {
	return s.repo.GetWatchHistory(ctx, userID, movieID)
}

// UpdateWatchProgress updates the watch position for a movie.
func (s *Service) UpdateWatchProgress(ctx context.Context, history *WatchHistory) error {
	existing, err := s.repo.GetWatchHistory(ctx, history.UserID, history.MovieID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create new watch history
		return s.repo.CreateWatchHistory(ctx, history)
	}

	// Update existing
	return s.repo.UpdateWatchHistory(ctx, existing.ID, history.PositionTicks, &history.DurationTicks)
}

// MarkAsWatched marks a movie as completely watched.
func (s *Service) MarkAsWatched(ctx context.Context, userID, movieID uuid.UUID) error {
	existing, err := s.repo.GetWatchHistory(ctx, userID, movieID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create completed watch history
		history := &WatchHistory{
			UserID:    userID,
			MovieID:   movieID,
			Completed: true,
		}
		return s.repo.CreateWatchHistory(ctx, history)
	}

	return s.repo.MarkWatchHistoryCompleted(ctx, existing.ID)
}

// ListResumeableMovies returns movies the user can resume.
func (s *Service) ListResumeableMovies(ctx context.Context, userID uuid.UUID, limit int) ([]WatchHistory, error) {
	return s.repo.ListResumeableMovies(ctx, userID, limit)
}

// IsWatched checks if a user has watched a movie.
func (s *Service) IsWatched(ctx context.Context, userID, movieID uuid.UUID) (bool, error) {
	return s.repo.IsWatched(ctx, userID, movieID)
}

// MarkAsUnwatched marks a movie as unwatched for a user.
func (s *Service) MarkAsUnwatched(ctx context.Context, userID, movieID uuid.UUID) error {
	// Delete watch history for this movie
	history, err := s.repo.GetWatchHistory(ctx, userID, movieID)
	if err != nil || history == nil {
		return nil // Already unwatched
	}
	return s.repo.DeleteWatchHistory(ctx, history.ID)
}

// Collections

// ListCollections returns all collections with pagination.
func (s *Service) ListCollections(ctx context.Context, params ListParams) ([]*Collection, int64, error) {
	collections, err := s.repo.ListCollections(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountCollections(ctx)
	if err != nil {
		return nil, 0, err
	}

	return collections, total, nil
}

// GetCollection returns a collection by ID.
func (s *Service) GetCollection(ctx context.Context, id uuid.UUID) (*Collection, error) {
	return s.repo.GetCollectionByID(ctx, id)
}

// ListMoviesByCollection returns all movies in a collection.
func (s *Service) ListMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]*Movie, error) {
	return s.repo.ListByCollection(ctx, collectionID)
}

// ListRecentlyAdded returns recently added movies from the specified libraries.
func (s *Service) ListRecentlyAdded(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Movie, error) {
	return s.repo.ListRecentlyAdded(ctx, libraryIDs, limit)
}

// CacheStats returns local cache statistics.
func (s *Service) CacheStats() CacheStats {
	stats := s.cache.Stats()
	return CacheStats{
		Hits:      stats.Hits(),
		Misses:    stats.Misses(),
		Ratio:     stats.Ratio(),
		Evictions: stats.EvictedCount(),
	}
}

// CacheStats contains cache statistics.
type CacheStats struct {
	Hits      int64
	Misses    int64
	Ratio     float64
	Evictions int64
}

// ServiceWithRetry wraps service with retry capability.
type ServiceWithRetry struct {
	*Service
	retry resilience.Retry
}

// NewServiceWithRetry creates a service with retry wrapper.
func NewServiceWithRetry(service *Service) *ServiceWithRetry {
	return &ServiceWithRetry{
		Service: service,
		retry:   resilience.DefaultRetry(),
	}
}

// GetMovieWithRetry retrieves a movie with automatic retry.
func (s *ServiceWithRetry) GetMovieWithRetry(ctx context.Context, id uuid.UUID) (*Movie, error) {
	var movie *Movie
	err := s.retry.DoWithContext(ctx, func(ctx context.Context) error {
		var err error
		movie, err = s.Service.GetMovie(ctx, id)
		return err
	})
	return movie, err
}

// ApplyMetadataWithRetry applies metadata updates with automatic retry.
func (s *ServiceWithRetry) ApplyMetadataWithRetry(ctx context.Context, id uuid.UUID, updates MovieMetadataUpdate) error {
	return s.retry.DoWithContext(ctx, func(ctx context.Context) error {
		return s.Service.ApplyMetadata(ctx, id, updates)
	})
}
