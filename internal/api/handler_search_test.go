package api

import (
	"context"
	"errors"
	"testing"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// ReindexSearch Tests
// ============================================================================

func TestHandler_ReindexSearch_NilRiverClient(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:      logging.NewTestLogger(),
		riverClient: nil,
	}

	result, err := handler.ReindexSearch(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "job queue not available")
}

func TestHandler_ReindexSearch_Success(t *testing.T) {
	t.Parallel()

	mockRiver := &mockRiverClient{}
	handler := &Handler{
		logger:      logging.NewTestLogger(),
		riverClient: mockRiver,
	}

	result, err := handler.ReindexSearch(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	accepted, ok := result.(*ogen.ReindexSearchAccepted)
	require.True(t, ok)
	assert.Equal(t, "Reindex job enqueued", accepted.Message.Value)
	assert.True(t, accepted.JobID.Set, "JobID should be set")

	// Verify job was inserted with correct args
	require.Len(t, mockRiver.insertedArgs, 1)
	indexArgs, ok := mockRiver.insertedArgs[0].(moviejobs.MovieSearchIndexArgs)
	require.True(t, ok)
	assert.Equal(t, moviejobs.SearchIndexOperationReindex, indexArgs.Operation)
}

func TestHandler_ReindexSearch_InsertError(t *testing.T) {
	t.Parallel()

	mockRiver := &mockRiverClient{
		insertError: errors.New("queue connection failed"),
	}
	handler := &Handler{
		logger:      logging.NewTestLogger(),
		riverClient: mockRiver,
	}

	result, err := handler.ReindexSearch(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to enqueue reindex job")
	assert.Contains(t, err.Error(), "queue connection failed")
}

// ============================================================================
// SearchLibraryMovies Tests
// ============================================================================

func TestHandler_SearchLibraryMovies_DisabledService(t *testing.T) {
	t.Parallel()

	// A MovieSearchService with nil client reports IsEnabled()=false and
	// returns empty results without contacting Typesense.
	logger := logging.NewTestLogger()
	svc := search.NewMovieSearchService(nil, logger)

	handler := &Handler{
		logger:        logger,
		searchService: svc,
	}

	params := ogen.SearchLibraryMoviesParams{
		Q: "test query",
	}

	result, err := handler.SearchLibraryMovies(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	searchResults, ok := result.(*ogen.SearchResults)
	require.True(t, ok, "expected *ogen.SearchResults")
	assert.Equal(t, 0, searchResults.TotalHits.Value)
	assert.Equal(t, 0, searchResults.TotalPages.Value)
	assert.Equal(t, 0, searchResults.SearchTimeMs.Value)
	assert.Empty(t, searchResults.Hits)
}

func TestHandler_SearchLibraryMovies_DefaultPagination(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewMovieSearchService(nil, logger)

	handler := &Handler{
		logger:        logger,
		searchService: svc,
	}

	// No page/perPage set — defaults should apply (page=1, perPage=20)
	params := ogen.SearchLibraryMoviesParams{
		Q: "action movies",
	}

	result, err := handler.SearchLibraryMovies(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	searchResults, ok := result.(*ogen.SearchResults)
	require.True(t, ok)
	// Disabled service returns an empty SearchResult; handler converts it
	assert.NotNil(t, searchResults)
	assert.Empty(t, searchResults.Hits)
}

func TestHandler_SearchLibraryMovies_CustomPagination(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewMovieSearchService(nil, logger)

	handler := &Handler{
		logger:        logger,
		searchService: svc,
	}

	params := ogen.SearchLibraryMoviesParams{
		Q:       "sci-fi",
		Page:    ogen.NewOptInt(3),
		PerPage: ogen.NewOptInt(10),
	}

	result, err := handler.SearchLibraryMovies(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	searchResults, ok := result.(*ogen.SearchResults)
	require.True(t, ok)
	assert.Empty(t, searchResults.Hits)
}

// ============================================================================
// AutocompleteMovies Tests
// ============================================================================

func TestHandler_AutocompleteMovies_DisabledService(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewMovieSearchService(nil, logger)

	handler := &Handler{
		logger:        logger,
		searchService: svc,
	}

	params := ogen.AutocompleteMoviesParams{
		Q: "bat",
	}

	result, err := handler.AutocompleteMovies(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	acResults, ok := result.(*ogen.AutocompleteResults)
	require.True(t, ok, "expected *ogen.AutocompleteResults")
	// Disabled service returns nil suggestions, handler passes them through
	assert.Nil(t, acResults.Suggestions)
}

func TestHandler_AutocompleteMovies_DisabledServiceWithLimit(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewMovieSearchService(nil, logger)

	handler := &Handler{
		logger:        logger,
		searchService: svc,
	}

	params := ogen.AutocompleteMoviesParams{
		Q:     "star",
		Limit: ogen.NewOptInt(10),
	}

	result, err := handler.AutocompleteMovies(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	acResults, ok := result.(*ogen.AutocompleteResults)
	require.True(t, ok)
	assert.Nil(t, acResults.Suggestions)
}

// ============================================================================
// GetSearchFacets Tests
// ============================================================================

func TestHandler_GetSearchFacets_DisabledService(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewMovieSearchService(nil, logger)

	handler := &Handler{
		logger:        logger,
		searchService: svc,
	}

	result, err := handler.GetSearchFacets(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	facets, ok := result.(*ogen.SearchFacets)
	require.True(t, ok, "expected *ogen.SearchFacets")
	// All facet slices should be nil/empty since the disabled service returns
	// a nil facets map.
	assert.Nil(t, facets.Genres)
	assert.Nil(t, facets.Years)
	assert.Nil(t, facets.Status)
	assert.Nil(t, facets.Directors)
	assert.Nil(t, facets.Resolution)
	assert.Nil(t, facets.HasFile)
}

// ============================================================================
// SearchLibraryTVShows Tests
// ============================================================================

func TestHandler_SearchLibraryTVShows_NilService(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:              logging.NewTestLogger(),
		tvshowSearchService: nil,
	}

	params := ogen.SearchLibraryTVShowsParams{
		Q: "breaking bad",
	}

	result, err := handler.SearchLibraryTVShows(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	tvResults, ok := result.(*ogen.TVShowSearchResults)
	require.True(t, ok, "expected *ogen.TVShowSearchResults")
	assert.Equal(t, 0, tvResults.TotalHits.Value)
	assert.Equal(t, 0, tvResults.TotalPages.Value)
	assert.Equal(t, 1, tvResults.CurrentPage.Value)
	assert.Equal(t, 0, tvResults.SearchTimeMs.Value)
	assert.Empty(t, tvResults.Hits)
}

func TestHandler_SearchLibraryTVShows_DisabledService(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewTVShowSearchService(nil, logger)

	handler := &Handler{
		logger:              logger,
		tvshowSearchService: svc,
	}

	params := ogen.SearchLibraryTVShowsParams{
		Q: "the wire",
	}

	result, err := handler.SearchLibraryTVShows(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	tvResults, ok := result.(*ogen.TVShowSearchResults)
	require.True(t, ok)
	assert.Equal(t, 0, tvResults.TotalHits.Value)
	assert.Equal(t, 0, tvResults.TotalPages.Value)
	assert.Empty(t, tvResults.Hits)
}

// ============================================================================
// AutocompleteTVShows Tests
// ============================================================================

func TestHandler_AutocompleteTVShows_NilService(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:              logging.NewTestLogger(),
		tvshowSearchService: nil,
	}

	params := ogen.AutocompleteTVShowsParams{
		Q: "break",
	}

	result, err := handler.AutocompleteTVShows(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	acResults, ok := result.(*ogen.AutocompleteResults)
	require.True(t, ok, "expected *ogen.AutocompleteResults")
	assert.Empty(t, acResults.Suggestions)
}

func TestHandler_AutocompleteTVShows_DisabledService(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewTVShowSearchService(nil, logger)

	handler := &Handler{
		logger:              logger,
		tvshowSearchService: svc,
	}

	params := ogen.AutocompleteTVShowsParams{
		Q:     "game",
		Limit: ogen.NewOptInt(8),
	}

	result, err := handler.AutocompleteTVShows(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	acResults, ok := result.(*ogen.AutocompleteResults)
	require.True(t, ok)
	// Disabled service returns nil,nil from AutocompleteSeries
	assert.Nil(t, acResults.Suggestions)
}

// ============================================================================
// GetTVShowSearchFacets Tests
// ============================================================================

func TestHandler_GetTVShowSearchFacets_NilService(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger:              logging.NewTestLogger(),
		tvshowSearchService: nil,
	}

	result, err := handler.GetTVShowSearchFacets(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	facets, ok := result.(*ogen.TVShowSearchFacets)
	require.True(t, ok, "expected *ogen.TVShowSearchFacets")
	// All facet slices should be nil when service is nil
	assert.Nil(t, facets.Genres)
	assert.Nil(t, facets.Years)
	assert.Nil(t, facets.Status)
	assert.Nil(t, facets.Type)
	assert.Nil(t, facets.Networks)
	assert.Nil(t, facets.HasFile)
}

func TestHandler_GetTVShowSearchFacets_DisabledService(t *testing.T) {
	t.Parallel()

	logger := logging.NewTestLogger()
	svc := search.NewTVShowSearchService(nil, logger)

	handler := &Handler{
		logger:              logger,
		tvshowSearchService: svc,
	}

	result, err := handler.GetTVShowSearchFacets(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	facets, ok := result.(*ogen.TVShowSearchFacets)
	require.True(t, ok)
	assert.Nil(t, facets.Genres)
	assert.Nil(t, facets.Years)
	assert.Nil(t, facets.Status)
	assert.Nil(t, facets.Type)
	assert.Nil(t, facets.Networks)
	assert.Nil(t, facets.HasFile)
}
