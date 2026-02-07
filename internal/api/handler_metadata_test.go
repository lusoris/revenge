package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/metadata"
)

// mockMetadataService implements metadata.Service for testing.
type mockMetadataService struct {
	searchMovieResults     []metadata.MovieSearchResult
	searchMovieErr         error
	movieMetadata          *metadata.MovieMetadata
	movieMetadataErr       error
	collectionMetadata     *metadata.CollectionMetadata
	collectionMetadataErr  error
	searchTVShowResults    []metadata.TVShowSearchResult
	searchTVShowErr        error
	tvShowMetadata         *metadata.TVShowMetadata
	tvShowMetadataErr      error
	seasonMetadata         *metadata.SeasonMetadata
	seasonMetadataErr      error
	episodeMetadata        *metadata.EpisodeMetadata
	episodeMetadataErr     error
}

func (m *mockMetadataService) SearchMovie(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, error) {
	return m.searchMovieResults, m.searchMovieErr
}

func (m *mockMetadataService) GetMovieMetadata(_ context.Context, _ int32, _ []string) (*metadata.MovieMetadata, error) {
	return m.movieMetadata, m.movieMetadataErr
}

func (m *mockMetadataService) GetMovieCredits(_ context.Context, _ int32) (*metadata.Credits, error) {
	return nil, nil
}

func (m *mockMetadataService) GetMovieImages(_ context.Context, _ int32) (*metadata.Images, error) {
	return nil, nil
}

func (m *mockMetadataService) GetMovieReleaseDates(_ context.Context, _ int32) ([]metadata.ReleaseDate, error) {
	return nil, nil
}

func (m *mockMetadataService) GetMovieExternalIDs(_ context.Context, _ int32) (*metadata.ExternalIDs, error) {
	return nil, nil
}

func (m *mockMetadataService) GetSimilarMovies(_ context.Context, _ int32, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	return nil, 0, nil
}

func (m *mockMetadataService) GetMovieRecommendations(_ context.Context, _ int32, _ metadata.SearchOptions) ([]metadata.MovieSearchResult, int, error) {
	return nil, 0, nil
}

func (m *mockMetadataService) SearchTVShow(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.TVShowSearchResult, error) {
	return m.searchTVShowResults, m.searchTVShowErr
}

func (m *mockMetadataService) GetTVShowMetadata(_ context.Context, _ int32, _ []string) (*metadata.TVShowMetadata, error) {
	return m.tvShowMetadata, m.tvShowMetadataErr
}

func (m *mockMetadataService) GetTVShowCredits(_ context.Context, _ int32) (*metadata.Credits, error) {
	return nil, nil
}

func (m *mockMetadataService) GetTVShowImages(_ context.Context, _ int32) (*metadata.Images, error) {
	return nil, nil
}

func (m *mockMetadataService) GetTVShowContentRatings(_ context.Context, _ int32) ([]metadata.ContentRating, error) {
	return nil, nil
}

func (m *mockMetadataService) GetTVShowExternalIDs(_ context.Context, _ int32) (*metadata.ExternalIDs, error) {
	return nil, nil
}

func (m *mockMetadataService) GetSeasonMetadata(_ context.Context, _ int32, _ int, _ []string) (*metadata.SeasonMetadata, error) {
	return m.seasonMetadata, m.seasonMetadataErr
}

func (m *mockMetadataService) GetEpisodeMetadata(_ context.Context, _ int32, _, _ int, _ []string) (*metadata.EpisodeMetadata, error) {
	return m.episodeMetadata, m.episodeMetadataErr
}

func (m *mockMetadataService) SearchPerson(_ context.Context, _ string, _ metadata.SearchOptions) ([]metadata.PersonSearchResult, error) {
	return nil, nil
}

func (m *mockMetadataService) GetPersonMetadata(_ context.Context, _ int32, _ []string) (*metadata.PersonMetadata, error) {
	return nil, nil
}

func (m *mockMetadataService) GetPersonCredits(_ context.Context, _ int32) (*metadata.PersonCredits, error) {
	return nil, nil
}

func (m *mockMetadataService) GetPersonImages(_ context.Context, _ int32) (*metadata.Images, error) {
	return nil, nil
}

func (m *mockMetadataService) GetCollectionMetadata(_ context.Context, _ int32, _ []string) (*metadata.CollectionMetadata, error) {
	return m.collectionMetadata, m.collectionMetadataErr
}

func (m *mockMetadataService) GetImageURL(_ string, _ metadata.ImageSize) string {
	return ""
}

func (m *mockMetadataService) RefreshMovie(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (m *mockMetadataService) RefreshTVShow(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (m *mockMetadataService) ClearCache() {}

func (m *mockMetadataService) RegisterProvider(_ metadata.Provider) {}

func (m *mockMetadataService) GetProviders() []metadata.Provider {
	return nil
}

func newMetadataTestHandler(metaSvc metadata.Service) *Handler {
	return &Handler{
		logger:          zap.NewNop(),
		metadataService: metaSvc,
	}
}

// ============================================================================
// SearchMoviesMetadata Tests
// ============================================================================

func TestHandler_SearchMoviesMetadata_Success(t *testing.T) {
	t.Parallel()

	releaseDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	poster := "/poster.jpg"

	mock := &mockMetadataService{
		searchMovieResults: []metadata.MovieSearchResult{
			{
				ProviderID:    "12345",
				Title:         "Test Movie",
				OriginalTitle: "Test Original",
				Overview:      "A test movie",
				ReleaseDate:   &releaseDate,
				PosterPath:    &poster,
				VoteAverage:   7.5,
				VoteCount:     100,
				Popularity:    50.0,
			},
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.SearchMoviesMetadata(context.Background(), ogen.SearchMoviesMetadataParams{
		Q: "test",
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataSearchResults)
	require.True(t, ok)
	assert.Equal(t, 1, response.TotalResults.Value)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "Test Movie", response.Results[0].Title.Value)
	assert.Equal(t, 12345, response.Results[0].TmdbID.Value)
	assert.Equal(t, "A test movie", response.Results[0].Overview.Value)
}

func TestHandler_SearchMoviesMetadata_Empty(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		searchMovieResults: []metadata.MovieSearchResult{},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.SearchMoviesMetadata(context.Background(), ogen.SearchMoviesMetadataParams{
		Q: "nonexistent",
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataSearchResults)
	require.True(t, ok)
	assert.Equal(t, 0, response.TotalResults.Value)
	assert.Empty(t, response.Results)
}

func TestHandler_SearchMoviesMetadata_Error(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		searchMovieErr: errors.New("API error"),
	}
	handler := newMetadataTestHandler(mock)

	_, err := handler.SearchMoviesMetadata(context.Background(), ogen.SearchMoviesMetadataParams{
		Q: "test",
	})
	assert.Error(t, err)
}

func TestHandler_SearchMoviesMetadata_WithLimit(t *testing.T) {
	t.Parallel()

	results := make([]metadata.MovieSearchResult, 10)
	for i := range results {
		results[i] = metadata.MovieSearchResult{
			ProviderID: "1",
			Title:      "Movie",
		}
	}

	mock := &mockMetadataService{
		searchMovieResults: results,
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.SearchMoviesMetadata(context.Background(), ogen.SearchMoviesMetadataParams{
		Q:     "test",
		Limit: ogen.NewOptInt(3),
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataSearchResults)
	require.True(t, ok)
	assert.Len(t, response.Results, 3)
}

// ============================================================================
// GetMovieMetadata Tests
// ============================================================================

func TestHandler_GetMovieMetadata_Success(t *testing.T) {
	t.Parallel()

	tmdbID := int32(12345)
	imdbID := "tt1234567"
	overview := "A great movie"
	tagline := "Best movie ever"
	runtime := int32(120)

	mock := &mockMetadataService{
		movieMetadata: &metadata.MovieMetadata{
			TMDbID:        &tmdbID,
			IMDbID:        &imdbID,
			Title:         "Test Movie",
			OriginalTitle: "Test Original",
			Overview:      &overview,
			Tagline:       &tagline,
			Runtime:       &runtime,
			VoteAverage:   8.0,
			VoteCount:     200,
			Popularity:    75.0,
			Status:        "Released",
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetMovieMetadata(context.Background(), ogen.GetMovieMetadataParams{
		TmdbId: 12345,
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataMovie)
	require.True(t, ok)
	assert.Equal(t, "Test Movie", response.Title.Value)
	assert.Equal(t, 12345, response.TmdbID.Value)
	assert.Equal(t, "tt1234567", response.ImdbID.Value)
	assert.Equal(t, "A great movie", response.Overview.Value)
	assert.Equal(t, "Best movie ever", response.Tagline.Value)
	assert.Equal(t, 120, response.Runtime.Value)
	assert.Equal(t, "Released", response.Status.Value)
}

func TestHandler_GetMovieMetadata_NotFound(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		movieMetadata: nil,
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetMovieMetadata(context.Background(), ogen.GetMovieMetadataParams{
		TmdbId: 99999,
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.GetMovieMetadataNotFound)
	assert.True(t, ok)
}

func TestHandler_GetMovieMetadata_Error(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		movieMetadataErr: errors.New("API error"),
	}
	handler := newMetadataTestHandler(mock)

	_, err := handler.GetMovieMetadata(context.Background(), ogen.GetMovieMetadataParams{
		TmdbId: 12345,
	})
	assert.Error(t, err)
}

// ============================================================================
// GetCollectionMetadata Tests
// ============================================================================

func TestHandler_GetCollectionMetadata_Success(t *testing.T) {
	t.Parallel()

	overview := "A great collection"
	poster := "/collection_poster.jpg"

	mock := &mockMetadataService{
		collectionMetadata: &metadata.CollectionMetadata{
			ProviderID: "1000",
			Name:       "Test Collection",
			Overview:   &overview,
			PosterPath: &poster,
			Parts: []metadata.MovieSearchResult{
				{
					ProviderID: "101",
					Title:      "Part 1",
					VoteAverage: 7.0,
				},
				{
					ProviderID: "102",
					Title:      "Part 2",
					VoteAverage: 8.0,
				},
			},
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetCollectionMetadata(context.Background(), ogen.GetCollectionMetadataParams{
		TmdbId: 1000,
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataCollection)
	require.True(t, ok)
	assert.Equal(t, "Test Collection", response.Name.Value)
	assert.Equal(t, 1000, response.ID.Value)
	assert.Equal(t, "A great collection", response.Overview.Value)
	require.Len(t, response.Parts, 2)
	assert.Equal(t, "Part 1", response.Parts[0].Title.Value)
	assert.Equal(t, "Part 2", response.Parts[1].Title.Value)
}

func TestHandler_GetCollectionMetadata_NotFound(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		collectionMetadata: nil,
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetCollectionMetadata(context.Background(), ogen.GetCollectionMetadataParams{
		TmdbId: 99999,
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.GetCollectionMetadataNotFound)
	assert.True(t, ok)
}

// ============================================================================
// SearchTVShowsMetadata Tests
// ============================================================================

func TestHandler_SearchTVShowsMetadata_Success(t *testing.T) {
	t.Parallel()

	firstAir := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	poster := "/tv_poster.jpg"

	mock := &mockMetadataService{
		searchTVShowResults: []metadata.TVShowSearchResult{
			{
				ProviderID: "5000",
				Name:       "Test Show",
				OriginalName: "Original Show",
				Overview:   "A test show",
				FirstAirDate: &firstAir,
				PosterPath:   &poster,
				VoteAverage:  8.5,
				VoteCount:    500,
				Popularity:   90.0,
			},
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.SearchTVShowsMetadata(context.Background(), ogen.SearchTVShowsMetadataParams{
		Q: "test",
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataTVSearchResults)
	require.True(t, ok)
	assert.Equal(t, 1, response.TotalResults.Value)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "Test Show", response.Results[0].Name.Value)
	assert.Equal(t, 5000, response.Results[0].TmdbID.Value)
}

func TestHandler_SearchTVShowsMetadata_Error(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		searchTVShowErr: errors.New("API error"),
	}
	handler := newMetadataTestHandler(mock)

	_, err := handler.SearchTVShowsMetadata(context.Background(), ogen.SearchTVShowsMetadataParams{
		Q: "test",
	})
	assert.Error(t, err)
}

// ============================================================================
// GetTVShowMetadata Tests
// ============================================================================

func TestHandler_GetTVShowMetadata_Success(t *testing.T) {
	t.Parallel()

	tmdbID := int32(5000)
	tvdbID := int32(300000)
	overview := "Great show"

	mock := &mockMetadataService{
		tvShowMetadata: &metadata.TVShowMetadata{
			TMDbID:           &tmdbID,
			TVDbID:           &tvdbID,
			Name:             "Test Show",
			OriginalName:     "Original Show",
			Overview:         &overview,
			Status:           "Returning Series",
			Type:             "Scripted",
			NumberOfSeasons:  5,
			NumberOfEpisodes: 50,
			EpisodeRuntime:   []int{45, 60},
			VoteAverage:      8.5,
			VoteCount:        1000,
			Popularity:       90.0,
			Networks: []metadata.Network{
				{ID: 1, Name: "HBO"},
			},
			Genres: []metadata.Genre{
				{ID: 18, Name: "Drama"},
			},
			Seasons: []metadata.SeasonSummary{
				{ProviderID: "100", SeasonNumber: 1, Name: "Season 1", EpisodeCount: 10},
			},
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetTVShowMetadata(context.Background(), ogen.GetTVShowMetadataParams{
		TmdbId: 5000,
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataTVShow)
	require.True(t, ok)
	assert.Equal(t, "Test Show", response.Name.Value)
	assert.Equal(t, 5000, response.TmdbID.Value)
	assert.Equal(t, 300000, response.TvdbID.Value)
	assert.Equal(t, "Returning Series", response.Status.Value)
	assert.Equal(t, 5, response.NumberOfSeasons.Value)
	assert.Equal(t, 50, response.NumberOfEpisodes.Value)
	require.Len(t, response.Networks, 1)
	assert.Equal(t, "HBO", response.Networks[0].Name.Value)
	require.Len(t, response.Genres, 1)
	assert.Equal(t, "Drama", response.Genres[0].Name.Value)
	require.Len(t, response.Seasons, 1)
	assert.Equal(t, 1, response.Seasons[0].SeasonNumber.Value)
}

func TestHandler_GetTVShowMetadata_NotFound(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		tvShowMetadata: nil,
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetTVShowMetadata(context.Background(), ogen.GetTVShowMetadataParams{
		TmdbId: 99999,
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.GetTVShowMetadataNotFound)
	assert.True(t, ok)
}

// ============================================================================
// GetSeasonMetadata Tests
// ============================================================================

func TestHandler_GetSeasonMetadata_Success(t *testing.T) {
	t.Parallel()

	tmdbID := int32(9000)
	overview := "Season overview"
	airDate := time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)
	epOverview := "Episode overview"
	runtime := int32(45)

	mock := &mockMetadataService{
		seasonMetadata: &metadata.SeasonMetadata{
			TMDbID:       &tmdbID,
			SeasonNumber: 2,
			Name:         "Season 2",
			Overview:     &overview,
			AirDate:      &airDate,
			Episodes: []metadata.EpisodeSummary{
				{
					ProviderID:    "2001",
					EpisodeNumber: 1,
					Name:          "Pilot",
					Overview:      &epOverview,
					Runtime:       &runtime,
					VoteAverage:   9.0,
					VoteCount:     300,
				},
			},
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetSeasonMetadata(context.Background(), ogen.GetSeasonMetadataParams{
		TmdbId:       5000,
		SeasonNumber: 2,
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataSeason)
	require.True(t, ok)
	assert.Equal(t, "Season 2", response.Name.Value)
	assert.Equal(t, 2, response.SeasonNumber.Value)
	assert.Equal(t, "Season overview", response.Overview.Value)
	require.Len(t, response.Episodes, 1)
	assert.Equal(t, "Pilot", response.Episodes[0].Name.Value)
	assert.Equal(t, 1, response.Episodes[0].EpisodeNumber.Value)
}

func TestHandler_GetSeasonMetadata_NotFound(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		seasonMetadata: nil,
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetSeasonMetadata(context.Background(), ogen.GetSeasonMetadataParams{
		TmdbId:       5000,
		SeasonNumber: 99,
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.GetSeasonMetadataNotFound)
	assert.True(t, ok)
}

// ============================================================================
// GetEpisodeMetadata Tests
// ============================================================================

func TestHandler_GetEpisodeMetadata_Success(t *testing.T) {
	t.Parallel()

	tmdbID := int32(50000)
	overview := "Episode overview"
	runtime := int32(55)
	still := "/still.jpg"

	mock := &mockMetadataService{
		episodeMetadata: &metadata.EpisodeMetadata{
			TMDbID:        &tmdbID,
			SeasonNumber:  1,
			EpisodeNumber: 3,
			Name:          "Test Episode",
			Overview:      &overview,
			Runtime:       &runtime,
			StillPath:     &still,
			VoteAverage:   8.0,
			VoteCount:     150,
			Crew: []metadata.CrewMember{
				{
					ProviderID: "7001",
					Name:       "Director Name",
					Job:        "Director",
					Department: "Directing",
				},
			},
			GuestStars: []metadata.CastMember{
				{
					ProviderID: "8001",
					Name:       "Guest Actor",
					Character:  "Character Name",
					Order:      0,
				},
			},
		},
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetEpisodeMetadata(context.Background(), ogen.GetEpisodeMetadataParams{
		TmdbId:        5000,
		SeasonNumber:  1,
		EpisodeNumber: 3,
	})
	require.NoError(t, err)

	response, ok := result.(*ogen.MetadataEpisode)
	require.True(t, ok)
	assert.Equal(t, "Test Episode", response.Name.Value)
	assert.Equal(t, 3, response.EpisodeNumber.Value)
	assert.Equal(t, "Episode overview", response.Overview.Value)
	assert.Equal(t, 55, response.Runtime.Value)
	require.Len(t, response.Crew, 1)
	assert.Equal(t, "Director Name", response.Crew[0].Name.Value)
	assert.Equal(t, "Director", response.Crew[0].Job.Value)
	require.Len(t, response.GuestStars, 1)
	assert.Equal(t, "Guest Actor", response.GuestStars[0].Name.Value)
	assert.Equal(t, "Character Name", response.GuestStars[0].Character.Value)
}

func TestHandler_GetEpisodeMetadata_NotFound(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		episodeMetadata: nil,
	}
	handler := newMetadataTestHandler(mock)

	result, err := handler.GetEpisodeMetadata(context.Background(), ogen.GetEpisodeMetadataParams{
		TmdbId:        5000,
		SeasonNumber:  1,
		EpisodeNumber: 99,
	})
	require.NoError(t, err)

	_, ok := result.(*ogen.GetEpisodeMetadataNotFound)
	assert.True(t, ok)
}

func TestHandler_GetEpisodeMetadata_Error(t *testing.T) {
	t.Parallel()

	mock := &mockMetadataService{
		episodeMetadataErr: errors.New("API error"),
	}
	handler := newMetadataTestHandler(mock)

	_, err := handler.GetEpisodeMetadata(context.Background(), ogen.GetEpisodeMetadataParams{
		TmdbId:        5000,
		SeasonNumber:  1,
		EpisodeNumber: 1,
	})
	assert.Error(t, err)
}
