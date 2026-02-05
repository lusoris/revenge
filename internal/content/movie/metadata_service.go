package movie

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MetadataProvider defines the interface for metadata operations.
type MetadataProvider interface {
	SearchMovies(ctx context.Context, query string, year *int) ([]*Movie, error)
	EnrichMovie(ctx context.Context, mov *Movie) error
	GetMovieCredits(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]MovieCredit, error)
	GetMovieGenres(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]MovieGenre, error)
	ClearCache()
}

// MetadataService implements MetadataProvider using TMDb
type MetadataService struct {
	client *TMDbClient
	mapper *TMDbMapper
}

func NewMetadataService(config TMDbConfig) *MetadataService {
	client := NewTMDbClient(config)
	mapper := NewTMDbMapper(client)

	return &MetadataService{
		client: client,
		mapper: mapper,
	}
}

func (s *MetadataService) SearchMovies(ctx context.Context, query string, year *int) ([]*Movie, error) {
	response, err := s.client.SearchMovies(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("search movies: %w", err)
	}

	movies := make([]*Movie, 0, len(response.Results))
	for i := range response.Results {
		mov := s.mapper.MapSearchResult(&response.Results[i])
		movies = append(movies, mov)
	}

	return movies, nil
}

func (s *MetadataService) GetMovieByTMDbID(ctx context.Context, tmdbID int) (*Movie, error) {
	tmdbMovie, err := s.client.GetMovie(ctx, tmdbID)
	if err != nil {
		return nil, fmt.Errorf("get movie: %w", err)
	}

	mov := s.mapper.MapMovie(tmdbMovie)
	return mov, nil
}

// GetMovieByTMDbIDMultiLanguage fetches movie metadata in multiple languages
func (s *MetadataService) GetMovieByTMDbIDMultiLanguage(ctx context.Context, tmdbID int, languages []string) (*Movie, error) {
	if len(languages) == 0 {
		languages = []string{"en-US", "de-DE", "fr-FR", "es-ES", "ja-JP"}
	}

	// Fetch movie in multiple languages
	multiLangResult, err := s.client.GetMovieMultiLanguage(ctx, tmdbID, languages)
	if err != nil {
		return nil, fmt.Errorf("get movie multi-language: %w", err)
	}

	// Fetch release dates for age ratings
	releaseDates, err := s.client.GetMovieReleaseDates(ctx, tmdbID)
	if err != nil {
		// Log warning but continue without age ratings
		releaseDates = nil
	}

	// Map to domain model with all languages
	mov := s.mapper.MapMultiLanguageMovie(multiLangResult, releaseDates)
	if mov == nil {
		return nil, fmt.Errorf("failed to map multi-language movie (English missing)")
	}

	return mov, nil
}

func (s *MetadataService) GetMovieCredits(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]MovieCredit, error) {
	credits, err := s.client.GetMovieCredits(ctx, tmdbID)
	if err != nil {
		return nil, fmt.Errorf("get movie credits: %w", err)
	}

	return s.mapper.MapCredits(movieID, credits), nil
}

func (s *MetadataService) GetMovieImages(ctx context.Context, tmdbID int) (*TMDbImages, error) {
	images, err := s.client.GetMovieImages(ctx, tmdbID)
	if err != nil {
		return nil, fmt.Errorf("get movie images: %w", err)
	}

	return images, nil
}

func (s *MetadataService) GetMovieGenres(ctx context.Context, movieID uuid.UUID, tmdbID int) ([]MovieGenre, error) {
	tmdbMovie, err := s.client.GetMovie(ctx, tmdbID)
	if err != nil {
		return nil, fmt.Errorf("get movie for genres: %w", err)
	}

	return s.mapper.MapGenres(movieID, tmdbMovie.Genres), nil
}

func (s *MetadataService) GetCollection(ctx context.Context, collectionID int) (*MovieCollection, error) {
	collection, err := s.client.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}

	return s.mapper.MapCollection(collection), nil
}

func (s *MetadataService) GetCollectionMovies(ctx context.Context, collectionID int) ([]*Movie, error) {
	collection, err := s.client.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collection movies: %w", err)
	}

	movies := make([]*Movie, 0, len(collection.Parts))
	for i := range collection.Parts {
		mov := s.mapper.MapSearchResult(&collection.Parts[i])
		movies = append(movies, mov)
	}

	return movies, nil
}

// GetSimilarMovies returns movies similar to the given movie.
func (s *MetadataService) GetSimilarMovies(ctx context.Context, tmdbID int) ([]*Movie, int, error) {
	response, err := s.client.GetSimilarMovies(ctx, tmdbID)
	if err != nil {
		return nil, 0, fmt.Errorf("get similar movies: %w", err)
	}

	movies := make([]*Movie, 0, len(response.Results))
	for i := range response.Results {
		mov := s.mapper.MapSearchResult(&response.Results[i])
		movies = append(movies, mov)
	}

	return movies, response.TotalResults, nil
}

// GetCollectionDetails returns full collection details from TMDb.
func (s *MetadataService) GetCollectionDetails(ctx context.Context, collectionID int) (*TMDbCollectionDetails, error) {
	return s.client.GetCollection(ctx, collectionID)
}

func (s *MetadataService) EnrichMovie(ctx context.Context, mov *Movie) error {
	return s.EnrichMovieWithLanguages(ctx, mov, nil)
}

// EnrichMovieWithLanguages enriches movie with metadata in multiple languages
func (s *MetadataService) EnrichMovieWithLanguages(ctx context.Context, mov *Movie, languages []string) error {
	if mov.TMDbID == nil {
		return fmt.Errorf("movie has no TMDb ID")
	}

	tmdbID := int(*mov.TMDbID)

	// Use default languages if not specified
	if len(languages) == 0 {
		languages = []string{"en-US", "de-DE", "fr-FR", "es-ES", "ja-JP"}
	}

	// Fetch movie in multiple languages
	multiLangResult, err := s.client.GetMovieMultiLanguage(ctx, tmdbID, languages)
	if err != nil {
		return fmt.Errorf("fetch multi-language metadata: %w", err)
	}

	// Fetch release dates for age ratings
	releaseDates, err := s.client.GetMovieReleaseDates(ctx, tmdbID)
	if err != nil {
		// Log warning but continue without age ratings
		releaseDates = nil
	}

	// Map to domain model with all languages
	enriched := s.mapper.MapMultiLanguageMovie(multiLangResult, releaseDates)
	if enriched == nil {
		return fmt.Errorf("failed to map multi-language movie (English missing)")
	}

	// Merge basic fields
	mov.Title = enriched.Title
	mov.IMDbID = enriched.IMDbID
	mov.OriginalTitle = enriched.OriginalTitle
	mov.OriginalLanguage = enriched.OriginalLanguage
	mov.Overview = enriched.Overview
	mov.Tagline = enriched.Tagline
	mov.ReleaseDate = enriched.ReleaseDate
	mov.Year = enriched.Year
	mov.Runtime = enriched.Runtime
	mov.Budget = enriched.Budget
	mov.Revenue = enriched.Revenue
	mov.Status = enriched.Status
	mov.VoteAverage = enriched.VoteAverage
	mov.VoteCount = enriched.VoteCount
	mov.Popularity = enriched.Popularity
	mov.PosterPath = enriched.PosterPath
	mov.BackdropPath = enriched.BackdropPath

	// Merge multi-language fields
	mov.TitlesI18n = enriched.TitlesI18n
	mov.TaglinesI18n = enriched.TaglinesI18n
	mov.OverviewsI18n = enriched.OverviewsI18n
	mov.AgeRatings = enriched.AgeRatings

	return nil
}

func (s *MetadataService) DownloadPoster(ctx context.Context, posterPath string, size string) ([]byte, error) {
	if size == "" {
		size = "w500"
	}

	data, err := s.client.DownloadImage(ctx, posterPath, size)
	if err != nil {
		return nil, fmt.Errorf("download poster: %w", err)
	}

	return data, nil
}

func (s *MetadataService) DownloadBackdrop(ctx context.Context, backdropPath string, size string) ([]byte, error) {
	if size == "" {
		size = "w1280"
	}

	data, err := s.client.DownloadImage(ctx, backdropPath, size)
	if err != nil {
		return nil, fmt.Errorf("download backdrop: %w", err)
	}

	return data, nil
}

func (s *MetadataService) GetPosterURL(posterPath *string, size string) *string {
	return s.mapper.GetPosterURL(posterPath, size)
}

func (s *MetadataService) GetBackdropURL(backdropPath *string, size string) *string {
	return s.mapper.GetBackdropURL(backdropPath, size)
}

func (s *MetadataService) ClearCache() {
	s.client.ClearCache()
}
