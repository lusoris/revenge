package movie

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

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

func (s *MetadataService) EnrichMovie(ctx context.Context, mov *Movie) error {
	if mov.TMDbID == nil {
		return fmt.Errorf("movie has no TMDb ID")
	}

	tmdbID := int(*mov.TMDbID)

	tmdbMovie, err := s.client.GetMovie(ctx, tmdbID)
	if err != nil {
		return fmt.Errorf("fetch movie metadata: %w", err)
	}

	enriched := s.mapper.MapMovie(tmdbMovie)

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
