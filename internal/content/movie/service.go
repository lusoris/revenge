package movie

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content"
)

// Service defines business logic for movies
type Service interface {
	// Movie operations
	GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error)
	GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error)
	GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error)
	ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error)
	CountMovies(ctx context.Context) (int64, error)
	SearchMovies(ctx context.Context, query string, filters SearchFilters) ([]Movie, error)
	ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, int64, error)
	ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, int64, error)
	CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error)
	UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error)
	DeleteMovie(ctx context.Context, id uuid.UUID) error

	// Movie files
	GetMovieFiles(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error)
	CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error)
	DeleteMovieFile(ctx context.Context, id uuid.UUID) error

	// Credits
	GetMovieCast(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error)
	GetMovieCrew(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error)

	// Collections
	GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error)
	GetMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error)
	GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error)

	// Genres
	GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error)
	GetMoviesByGenre(ctx context.Context, slug string, limit, offset int32) ([]Movie, error)
	ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error)

	// Watch progress
	UpdateWatchProgress(ctx context.Context, userID, movieID uuid.UUID, progressSeconds, durationSeconds int32) (*MovieWatched, error)
	GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error)
	MarkAsWatched(ctx context.Context, userID, movieID uuid.UUID) error
	RemoveWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error
	GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error)
	GetWatchHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error)

	// Metadata refresh
	RefreshMovieMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error
}

// movieService implements the Service interface
type movieService struct {
	repo             Repository
	metadataProvider MetadataProvider
}

// NewService creates a new movie service
func NewService(repo Repository, metadataProvider MetadataProvider) Service {
	return &movieService{
		repo:             repo,
		metadataProvider: metadataProvider,
	}
}

// GetMovie retrieves a movie by ID
func (s *movieService) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
	return s.repo.GetMovie(ctx, id)
}

// GetMovieByTMDbID retrieves a movie by TMDb ID
func (s *movieService) GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error) {
	return s.repo.GetMovieByTMDbID(ctx, tmdbID)
}

// GetMovieByIMDbID retrieves a movie by IMDb ID
func (s *movieService) GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	return s.repo.GetMovieByIMDbID(ctx, imdbID)
}

// ListMovies returns a paginated list of movies
func (s *movieService) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	return s.repo.ListMovies(ctx, filters)
}

// CountMovies returns the total number of movies
func (s *movieService) CountMovies(ctx context.Context) (int64, error) {
	return s.repo.CountMovies(ctx)
}

// SearchMovies searches for movies by title
func (s *movieService) SearchMovies(ctx context.Context, query string, filters SearchFilters) ([]Movie, error) {
	return s.repo.SearchMoviesByTitle(ctx, query, filters.Limit, filters.Offset)
}

// ListRecentlyAdded returns recently added movies with total count
func (s *movieService) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, int64, error) {
	movies, err := s.repo.ListRecentlyAdded(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.CountMovies(ctx)
	if err != nil {
		return nil, 0, err
	}
	return movies, count, nil
}

// ListTopRated returns top-rated movies with total count
func (s *movieService) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, int64, error) {
	movies, err := s.repo.ListTopRated(ctx, minVotes, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.CountTopRated(ctx, minVotes)
	if err != nil {
		return nil, 0, err
	}
	return movies, count, nil
}

// CreateMovie creates a new movie
func (s *movieService) CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error) {
	// Validate required fields
	if params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return s.repo.CreateMovie(ctx, params)
}

// UpdateMovie updates an existing movie
func (s *movieService) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	// Verify movie exists
	existing, err := s.repo.GetMovie(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Update only if something changed
	if params.Title != nil && *params.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	_ = existing // Use existing for validation if needed
	return s.repo.UpdateMovie(ctx, params)
}

// DeleteMovie soft-deletes a movie
func (s *movieService) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMovie(ctx, id)
}

// GetMovieFiles returns all files for a movie
func (s *movieService) GetMovieFiles(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	return s.repo.ListMovieFilesByMovieID(ctx, movieID)
}

// CreateMovieFile creates a new movie file
func (s *movieService) CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error) {
	// Verify movie exists
	_, err := s.repo.GetMovie(ctx, params.MovieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Check if file already exists
	existing, err := s.repo.GetMovieFileByPath(ctx, params.FilePath)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("file already exists at path: %s", params.FilePath)
	}

	return s.repo.CreateMovieFile(ctx, params)
}

// DeleteMovieFile deletes a movie file
func (s *movieService) DeleteMovieFile(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMovieFile(ctx, id)
}

// GetMovieCast returns the cast for a movie with total count
func (s *movieService) GetMovieCast(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error) {
	credits, err := s.repo.ListMovieCast(ctx, movieID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.CountMovieCast(ctx, movieID)
	if err != nil {
		return nil, 0, err
	}
	return credits, count, nil
}

// GetMovieCrew returns the crew for a movie with total count
func (s *movieService) GetMovieCrew(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, int64, error) {
	credits, err := s.repo.ListMovieCrew(ctx, movieID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.CountMovieCrew(ctx, movieID)
	if err != nil {
		return nil, 0, err
	}
	return credits, count, nil
}

// GetMovieCollection retrieves a collection by ID
func (s *movieService) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	return s.repo.GetMovieCollection(ctx, id)
}

// GetMoviesByCollection returns all movies in a collection
func (s *movieService) GetMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error) {
	return s.repo.ListMoviesByCollection(ctx, collectionID)
}

// GetCollectionForMovie returns the collection a movie belongs to
func (s *movieService) GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error) {
	return s.repo.GetCollectionForMovie(ctx, movieID)
}

// GetMovieGenres returns genres for a movie
func (s *movieService) GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	return s.repo.ListMovieGenres(ctx, movieID)
}

// GetMoviesByGenre returns movies filtered by genre
func (s *movieService) GetMoviesByGenre(ctx context.Context, slug string, limit, offset int32) ([]Movie, error) {
	return s.repo.ListMoviesByGenre(ctx, slug, limit, offset)
}

// ListDistinctGenres returns all distinct movie genres with item counts.
func (s *movieService) ListDistinctGenres(ctx context.Context) ([]content.GenreSummary, error) {
	return s.repo.ListDistinctMovieGenres(ctx)
}

// UpdateWatchProgress updates or creates watch progress for a user
func (s *movieService) UpdateWatchProgress(ctx context.Context, userID, movieID uuid.UUID, progressSeconds, durationSeconds int32) (*MovieWatched, error) {
	// Verify movie exists
	_, err := s.repo.GetMovie(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Calculate if completed (>90% watched)
	isCompleted := false
	if durationSeconds > 0 {
		progress := float64(progressSeconds) / float64(durationSeconds)
		isCompleted = progress > 0.90
	}

	params := CreateWatchProgressParams{
		UserID:          userID,
		MovieID:         movieID,
		ProgressSeconds: progressSeconds,
		DurationSeconds: durationSeconds,
		IsCompleted:     isCompleted,
	}

	return s.repo.CreateOrUpdateWatchProgress(ctx, params)
}

// GetWatchProgress retrieves watch progress for a user and movie
func (s *movieService) GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error) {
	return s.repo.GetWatchProgress(ctx, userID, movieID)
}

// MarkAsWatched marks a movie as watched (100% complete)
func (s *movieService) MarkAsWatched(ctx context.Context, userID, movieID uuid.UUID) error {
	// Get movie to get duration
	movie, err := s.repo.GetMovie(ctx, movieID)
	if err != nil {
		return fmt.Errorf("movie not found: %w", err)
	}

	// Use runtime if available, otherwise default to 7200 seconds (2 hours)
	durationSeconds := int32(7200)
	if movie.Runtime != nil && *movie.Runtime > 0 {
		durationSeconds = *movie.Runtime * 60 // Convert minutes to seconds
	}

	params := CreateWatchProgressParams{
		UserID:          userID,
		MovieID:         movieID,
		ProgressSeconds: durationSeconds,
		DurationSeconds: durationSeconds,
		IsCompleted:     true,
	}

	_, err = s.repo.CreateOrUpdateWatchProgress(ctx, params)
	return err
}

// RemoveWatchProgress removes watch progress for a user and movie
func (s *movieService) RemoveWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error {
	return s.repo.DeleteWatchProgress(ctx, userID, movieID)
}

// GetContinueWatching returns movies the user is currently watching
func (s *movieService) GetContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	return s.repo.ListContinueWatching(ctx, userID, limit)
}

// GetWatchHistory returns the user's watch history
func (s *movieService) GetWatchHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error) {
	return s.repo.ListWatchedMovies(ctx, userID, limit, offset)
}

// GetUserStats returns statistics about a user's movie watching
func (s *movieService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error) {
	return s.repo.GetUserMovieStats(ctx, userID)
}

// RefreshMovieMetadata triggers a metadata refresh for a movie.
// Options allow specifying force refresh and language overrides.
func (s *movieService) RefreshMovieMetadata(ctx context.Context, id uuid.UUID, opts ...MetadataRefreshOptions) error {
	// Check if metadata provider is available
	if s.metadataProvider == nil {
		return fmt.Errorf("metadata provider not configured")
	}

	// Get the movie
	mov, err := s.repo.GetMovie(ctx, id)
	if err != nil {
		return fmt.Errorf("get movie: %w", err)
	}

	// Enrich with latest metadata from TMDb, passing options through
	if err := s.metadataProvider.EnrichMovie(ctx, mov, opts...); err != nil {
		return fmt.Errorf("enrich movie: %w", err)
	}

	// Build update params from enriched movie
	params := movieToUpdateParams(mov)

	// Update the movie
	if _, err := s.repo.UpdateMovie(ctx, params); err != nil {
		return fmt.Errorf("update movie: %w", err)
	}

	// Update credits if TMDbID is available
	if mov.TMDbID != nil {
		credits, err := s.metadataProvider.GetMovieCredits(ctx, mov.ID, fmt.Sprintf("%d", *mov.TMDbID))
		if err == nil && len(credits) > 0 {
			// Delete existing credits and create new ones
			_ = s.repo.DeleteMovieCredits(ctx, mov.ID)
			for _, credit := range credits {
				_, _ = s.repo.CreateMovieCredit(ctx, CreateMovieCreditParams{
					MovieID:      credit.MovieID,
					TMDbPersonID: credit.TMDbPersonID,
					Name:         credit.Name,
					CreditType:   credit.CreditType,
					Character:    credit.Character,
					Job:          credit.Job,
					Department:   credit.Department,
					CastOrder:    credit.CastOrder,
					ProfilePath:  credit.ProfilePath,
				})
			}
		}

		// Update genres
		genres, err := s.metadataProvider.GetMovieGenres(ctx, mov.ID, fmt.Sprintf("%d", *mov.TMDbID))
		if err == nil && len(genres) > 0 {
			_ = s.repo.DeleteMovieGenres(ctx, mov.ID)
			for _, genre := range genres {
				_ = s.repo.AddMovieGenre(ctx, mov.ID, genre.Slug, genre.Name)
			}
		}
	}

	return nil
}

// movieToUpdateParams converts a Movie to UpdateMovieParams
func movieToUpdateParams(m *Movie) UpdateMovieParams {
	params := UpdateMovieParams{
		ID:               m.ID,
		TMDbID:           m.TMDbID,
		IMDbID:           m.IMDbID,
		Title:            &m.Title,
		OriginalTitle:    m.OriginalTitle,
		Year:             m.Year,
		Runtime:          m.Runtime,
		Overview:         m.Overview,
		Tagline:          m.Tagline,
		Status:           m.Status,
		OriginalLanguage: m.OriginalLanguage,
		TitlesI18n:       m.TitlesI18n,
		TaglinesI18n:     m.TaglinesI18n,
		OverviewsI18n:    m.OverviewsI18n,
		AgeRatings:       m.AgeRatings,
		ExternalRatings:  m.ExternalRatings,
		PosterPath:       m.PosterPath,
		BackdropPath:     m.BackdropPath,
		TrailerURL:       m.TrailerURL,
		VoteCount:        m.VoteCount,
		Budget:           m.Budget,
		Revenue:          m.Revenue,
		RadarrID:         m.RadarrID,
	}
	if m.ReleaseDate != nil {
		d := m.ReleaseDate.Format("2006-01-02")
		params.ReleaseDate = &d
	}
	if m.VoteAverage != nil {
		s := m.VoteAverage.String()
		params.VoteAverage = &s
	}
	if m.Popularity != nil {
		s := m.Popularity.String()
		params.Popularity = &s
	}
	return params
}
