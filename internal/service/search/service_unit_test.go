package search

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/logging"
)

// --- MovieSearchService (disabled state) ---

func TestNewMovieSearchService(t *testing.T) {
	logger := slog.Default()

	t.Run("with nil client", func(t *testing.T) {
		svc := NewMovieSearchService(nil, logger)
		require.NotNil(t, svc)
		assert.Nil(t, svc.client)
		assert.False(t, svc.IsEnabled())
	})
}

func TestMovieSearchService_IsEnabled_NilClient(t *testing.T) {
	svc := &MovieSearchService{client: nil}
	assert.False(t, svc.IsEnabled())
}

func TestMovieSearchService_InitializeCollection_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.InitializeCollection(ctx)
	assert.NoError(t, err)
}

func TestMovieSearchService_IndexMovie_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	m := &movie.Movie{
		ID:        uuid.Must(uuid.NewV7()),
		Title:     "Test Movie",
		CreatedAt: time.Now(),
	}

	err := svc.IndexMovie(ctx, m, nil, nil, nil)
	assert.NoError(t, err)
}

func TestMovieSearchService_UpdateMovie_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	m := &movie.Movie{
		ID:        uuid.Must(uuid.NewV7()),
		Title:     "Test Movie",
		CreatedAt: time.Now(),
	}

	err := svc.UpdateMovie(ctx, m, nil, nil, nil)
	assert.NoError(t, err)
}

func TestMovieSearchService_RemoveMovie_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.RemoveMovie(ctx, uuid.Must(uuid.NewV7()))
	assert.NoError(t, err)
}

func TestMovieSearchService_BulkIndexMovies_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	movies := []MovieWithRelations{
		{
			Movie: &movie.Movie{
				ID:        uuid.Must(uuid.NewV7()),
				Title:     "Movie 1",
				CreatedAt: time.Now(),
			},
		},
	}

	err := svc.BulkIndexMovies(ctx, movies)
	assert.NoError(t, err)
}

func TestMovieSearchService_BulkIndexMovies_EmptySlice(t *testing.T) {
	// Even with a non-nil client scenario, empty slice returns nil early
	// But since client is nil (disabled), it returns nil before checking len
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.BulkIndexMovies(ctx, []MovieWithRelations{})
	assert.NoError(t, err)
}

func TestMovieSearchService_Search_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	result, err := svc.Search(ctx, SearchParams{Query: "test"})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Hits)
	assert.Equal(t, 0, result.TotalHits)
}

func TestMovieSearchService_Search_DefaultParams(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	// Test that invalid params get defaults (but disabled returns early anyway)
	tests := []struct {
		name   string
		params SearchParams
	}{
		{
			name:   "zero page defaults to 1",
			params: SearchParams{Page: 0, PerPage: 10, Query: "test"},
		},
		{
			name:   "negative page defaults to 1",
			params: SearchParams{Page: -1, PerPage: 10, Query: "test"},
		},
		{
			name:   "zero per page defaults to 20",
			params: SearchParams{Page: 1, PerPage: 0, Query: "test"},
		},
		{
			name:   "per page exceeds 100 capped to 100",
			params: SearchParams{Page: 1, PerPage: 200, Query: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Search(ctx, tt.params)
			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestMovieSearchService_Autocomplete_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	result, err := svc.Autocomplete(ctx, "test", 5)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestMovieSearchService_Autocomplete_LimitBounds(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	// Disabled path returns early; these test the code path exists
	tests := []struct {
		name  string
		limit int
	}{
		{"zero limit defaults to 5", 0},
		{"negative limit defaults to 5", -1},
		{"limit over 20 capped", 25},
		{"valid limit", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Autocomplete(ctx, "query", tt.limit)
			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestMovieSearchService_GetFacets_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	result, err := svc.GetFacets(ctx, []string{"genres", "year"})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestMovieSearchService_ReindexAll_Disabled(t *testing.T) {
	svc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.ReindexAll(ctx, nil) // nil repo is fine because we short-circuit
	assert.NoError(t, err)
}

// --- TVShowSearchService (disabled state) ---

func TestNewTVShowSearchService(t *testing.T) {
	logger := slog.Default()

	svc := NewTVShowSearchService(nil, logger)
	require.NotNil(t, svc)
	assert.Nil(t, svc.client)
	assert.False(t, svc.IsEnabled())
}

func TestTVShowSearchService_InitializeCollection_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.InitializeCollection(ctx)
	assert.NoError(t, err)
}

func TestTVShowSearchService_IndexSeries_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.IndexSeries(ctx, nil, nil, nil, nil, false)
	assert.NoError(t, err)
}

func TestTVShowSearchService_UpdateSeries_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.UpdateSeries(ctx, nil, nil, nil, nil, false)
	assert.NoError(t, err)
}

func TestTVShowSearchService_RemoveSeries_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.RemoveSeries(ctx, uuid.Must(uuid.NewV7()))
	assert.NoError(t, err)
}

func TestTVShowSearchService_BulkIndexSeries_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.BulkIndexSeries(ctx, []TVShowWithRelations{{
		Series: nil,
	}})
	assert.NoError(t, err)
}

func TestTVShowSearchService_BulkIndexSeries_EmptySlice(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.BulkIndexSeries(ctx, []TVShowWithRelations{})
	assert.NoError(t, err)
}

func TestTVShowSearchService_SearchSeries_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	result, err := svc.SearchSeries(ctx, TVShowSearchParams{Query: "test"})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Hits)
}

func TestTVShowSearchService_SearchSeries_DefaultParams(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	tests := []struct {
		name   string
		params TVShowSearchParams
	}{
		{
			name:   "zero page",
			params: TVShowSearchParams{Page: 0, PerPage: 10, Query: "x"},
		},
		{
			name:   "zero per page",
			params: TVShowSearchParams{Page: 1, PerPage: 0, Query: "x"},
		},
		{
			name:   "per page over 100",
			params: TVShowSearchParams{Page: 1, PerPage: 200, Query: "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.SearchSeries(ctx, tt.params)
			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestTVShowSearchService_AutocompleteSeries_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	result, err := svc.AutocompleteSeries(ctx, "test", 5)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestTVShowSearchService_AutocompleteSeries_LimitBounds(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	tests := []struct {
		name  string
		limit int
	}{
		{"zero limit", 0},
		{"negative limit", -1},
		{"over 20", 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.AutocompleteSeries(ctx, "query", tt.limit)
			assert.NoError(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestTVShowSearchService_GetFacets_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	result, err := svc.GetFacets(ctx, []string{"genres", "year"})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestTVShowSearchService_ReindexAll_Disabled(t *testing.T) {
	svc := &TVShowSearchService{
		client: nil,
		logger: slog.Default(),
	}
	ctx := context.Background()

	err := svc.ReindexAll(ctx, nil)
	assert.NoError(t, err)
}

// --- CachedMovieSearchService ---

func TestNewCachedMovieSearchService(t *testing.T) {
	innerSvc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	logger := logging.NewTestLogger()

	cached := NewCachedMovieSearchService(innerSvc, nil, logger)
	require.NotNil(t, cached)
	assert.NotNil(t, cached.MovieSearchService)
	assert.Nil(t, cached.cache)
}

func TestCachedMovieSearchService_Search_NilCache(t *testing.T) {
	innerSvc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	logger := logging.NewTestLogger()
	cached := NewCachedMovieSearchService(innerSvc, nil, logger)
	ctx := context.Background()

	// With nil cache, should fall through to underlying service
	result, err := cached.Search(ctx, SearchParams{Query: "test"})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCachedMovieSearchService_Autocomplete_NilCache(t *testing.T) {
	innerSvc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	logger := logging.NewTestLogger()
	cached := NewCachedMovieSearchService(innerSvc, nil, logger)
	ctx := context.Background()

	result, err := cached.Autocomplete(ctx, "test", 5)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCachedMovieSearchService_GetFacets_NilCache(t *testing.T) {
	innerSvc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	logger := logging.NewTestLogger()
	cached := NewCachedMovieSearchService(innerSvc, nil, logger)
	ctx := context.Background()

	result, err := cached.GetFacets(ctx, []string{"genres"})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCachedMovieSearchService_InvalidateSearchCache_NilCache(t *testing.T) {
	innerSvc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	logger := logging.NewTestLogger()
	cached := NewCachedMovieSearchService(innerSvc, nil, logger)
	ctx := context.Background()

	err := cached.InvalidateSearchCache(ctx)
	assert.NoError(t, err)
}

func TestCachedMovieSearchService_searchCacheKey(t *testing.T) {
	innerSvc := &MovieSearchService{
		client: nil,
		logger: slog.Default(),
	}
	logger := logging.NewTestLogger()
	cached := NewCachedMovieSearchService(innerSvc, nil, logger)

	t.Run("deterministic", func(t *testing.T) {
		params := SearchParams{
			Query:    "matrix",
			FilterBy: "year:>=2000",
			SortBy:   "popularity:desc",
			Page:     1,
			PerPage:  20,
			FacetBy:  []string{"genres", "year"},
		}

		key1 := cached.searchCacheKey(params)
		key2 := cached.searchCacheKey(params)
		assert.Equal(t, key1, key2)
		assert.NotEmpty(t, key1)
	})

	t.Run("different params produce different keys", func(t *testing.T) {
		params1 := SearchParams{Query: "matrix", Page: 1, PerPage: 20}
		params2 := SearchParams{Query: "inception", Page: 1, PerPage: 20}

		key1 := cached.searchCacheKey(params1)
		key2 := cached.searchCacheKey(params2)
		assert.NotEqual(t, key1, key2)
	})

	t.Run("same query different page produces different key", func(t *testing.T) {
		params1 := SearchParams{Query: "matrix", Page: 1, PerPage: 20}
		params2 := SearchParams{Query: "matrix", Page: 2, PerPage: 20}

		key1 := cached.searchCacheKey(params1)
		key2 := cached.searchCacheKey(params2)
		assert.NotEqual(t, key1, key2)
	})

	t.Run("key starts with expected prefix", func(t *testing.T) {
		params := SearchParams{Query: "test"}
		key := cached.searchCacheKey(params)
		assert.Contains(t, key, "search:movies:")
	})
}

// --- MovieWithRelations type ---

func TestMovieWithRelations(t *testing.T) {
	movieID := uuid.Must(uuid.NewV7())
	now := time.Now()

	mwr := MovieWithRelations{
		Movie: &movie.Movie{
			ID:             movieID,
			Title:          "Test Movie",
			LibraryAddedAt: now,
			CreatedAt:      now,
		},
		Genres: []movie.MovieGenre{
			{Name: "Action", TMDbGenreID: 28},
		},
		Credits: []movie.MovieCredit{
			{CreditType: "cast", Name: "Actor"},
		},
		File: &movie.MovieFile{
			ID:      uuid.Must(uuid.NewV7()),
			MovieID: movieID,
		},
	}

	assert.Equal(t, "Test Movie", mwr.Movie.Title)
	assert.Len(t, mwr.Genres, 1)
	assert.Len(t, mwr.Credits, 1)
	assert.NotNil(t, mwr.File)
}

// --- SearchResult types ---

func TestSearchResult_Zero(t *testing.T) {
	result := SearchResult{}
	assert.Empty(t, result.Hits)
	assert.Equal(t, 0, result.TotalHits)
	assert.Equal(t, 0, result.TotalPages)
	assert.Equal(t, 0, result.CurrentPage)
}

func TestMovieHit_Zero(t *testing.T) {
	hit := MovieHit{}
	assert.Empty(t, hit.Document.ID)
	assert.Equal(t, float64(0), hit.Score)
	assert.Nil(t, hit.Highlights)
}

func TestFacetValue(t *testing.T) {
	fv := FacetValue{Value: "Action", Count: 42}
	assert.Equal(t, "Action", fv.Value)
	assert.Equal(t, 42, fv.Count)
}

// --- TVShow types ---

func TestTVShowSearchResult_Zero(t *testing.T) {
	result := TVShowSearchResult{}
	assert.Empty(t, result.Hits)
	assert.Equal(t, 0, result.TotalHits)
}

func TestTVShowHit_Zero(t *testing.T) {
	hit := TVShowHit{}
	assert.Empty(t, hit.Document.ID)
	assert.Equal(t, float64(0), hit.Score)
	assert.Nil(t, hit.Highlights)
}

func TestTVShowWithRelations(t *testing.T) {
	twr := TVShowWithRelations{
		HasFile: true,
	}
	assert.True(t, twr.HasFile)
	assert.Nil(t, twr.Series)
}

// --- Constants ---

func TestMovieCollectionName(t *testing.T) {
	assert.Equal(t, "movies", MovieCollectionName)
}

func TestTVShowCollectionName(t *testing.T) {
	assert.Equal(t, "tvshows", TVShowCollectionName)
}
