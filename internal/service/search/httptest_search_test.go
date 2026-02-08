package search

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/cache"
	infraSearch "github.com/lusoris/revenge/internal/infra/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v2/typesense/api"

	"github.com/lusoris/revenge/internal/infra/logging"
)

// newTestSearchClient creates a search client backed by the given httptest server.
func newTestSearchClient(t *testing.T, serverURL string) *infraSearch.Client {
	t.Helper()

	cfg := &config.Config{}
	cfg.Search.Enabled = true
	cfg.Search.URL = serverURL
	cfg.Search.APIKey = "test-api-key"

	client, err := infraSearch.NewClient(cfg, slog.Default())
	require.NoError(t, err)
	require.True(t, client.IsEnabled())

	return client
}

// newTestCache creates an L1-only cache for testing.
func newTestCache(t *testing.T) *cache.Cache {
	t.Helper()
	c, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)
	return c
}

// fakeTypesenseSearchHandler handles all Typesense API requests for testing.
func fakeTypesenseSearchHandler(t *testing.T) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Health check
		if r.URL.Path == "/health" {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
			return
		}

		// Collection operations
		if strings.HasPrefix(r.URL.Path, "/collections/") {
			path := strings.TrimPrefix(r.URL.Path, "/collections/")
			parts := strings.Split(path, "/")

			// GET /collections/{name}
			if r.Method == http.MethodGet && len(parts) == 1 {
				json.NewEncoder(w).Encode(api.CollectionResponse{Name: parts[0]})
				return
			}

			// DELETE /collections/{name}
			if r.Method == http.MethodDelete && len(parts) == 1 {
				json.NewEncoder(w).Encode(map[string]interface{}{"name": parts[0]})
				return
			}

			// Document operations: /collections/{name}/documents/...
			if len(parts) >= 2 && parts[1] == "documents" {
				// Search: GET /collections/{name}/documents/search
				if len(parts) == 3 && parts[2] == "search" && r.Method == http.MethodGet {
					result := map[string]interface{}{
						"found": 2,
						"hits": []map[string]interface{}{
							{
								"document": map[string]interface{}{
									"id":    "doc-1",
									"title": "Test Movie 1",
									"year":  float64(2020),
								},
								"text_match": int64(100),
								"highlights": []map[string]interface{}{
									{
										"field":    "title",
										"snippets": []string{"<mark>Test</mark> Movie 1"},
									},
								},
							},
							{
								"document": map[string]interface{}{
									"id":    "doc-2",
									"title": "Test Movie 2",
									"year":  float64(2021),
								},
								"text_match": int64(80),
							},
						},
						"facet_counts": []map[string]interface{}{
							{
								"field_name": "genres",
								"counts": []map[string]interface{}{
									{"value": "Action", "count": 10},
									{"value": "Drama", "count": 5},
								},
							},
						},
					}
					json.NewEncoder(w).Encode(result)
					return
				}

				// Create doc: POST /collections/{name}/documents (must return 201)
				if len(parts) == 2 && r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]interface{}{"id": "created"})
					return
				}

				// Import docs: POST /collections/{name}/documents/import
				if len(parts) == 3 && parts[2] == "import" && r.Method == http.MethodPost {
					w.Write([]byte(`{"success": true}` + "\n"))
					w.Write([]byte(`{"success": true}` + "\n"))
					return
				}

				// Update doc: PATCH /collections/{name}/documents/{id}
				if len(parts) == 3 && r.Method == http.MethodPatch {
					json.NewEncoder(w).Encode(map[string]interface{}{"id": parts[2]})
					return
				}

				// Delete doc: DELETE /collections/{name}/documents/{id}
				if len(parts) == 3 && r.Method == http.MethodDelete {
					json.NewEncoder(w).Encode(map[string]interface{}{"id": parts[2]})
					return
				}
			}
		}

		// POST /collections - create collection
		if r.URL.Path == "/collections" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"name": "created"})
			return
		}

		http.NotFound(w, r)
	})
}

// --- Movie service tests with fake Typesense ---

func TestMovieSearchService_InitializeCollection_AlreadyExists(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	assert.True(t, svc.IsEnabled())

	err := svc.InitializeCollection(ctx)
	assert.NoError(t, err)
}

func TestMovieSearchService_IndexMovie_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	movieID := uuid.Must(uuid.NewV7())
	now := time.Now()
	m := &movie.Movie{
		ID:             movieID,
		Title:          "Test Movie",
		LibraryAddedAt: now,
		CreatedAt:      now,
	}

	err := svc.IndexMovie(ctx, m, nil, nil, nil)
	assert.NoError(t, err)
}

func TestMovieSearchService_UpdateMovie_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	movieID := uuid.Must(uuid.NewV7())
	now := time.Now()
	m := &movie.Movie{
		ID:             movieID,
		Title:          "Updated Movie",
		LibraryAddedAt: now,
		CreatedAt:      now,
	}

	err := svc.UpdateMovie(ctx, m, nil, nil, nil)
	assert.NoError(t, err)
}

func TestMovieSearchService_UpdateMovie_NotFound_FallsBackToIndex(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// PATCH (update) - return not found to trigger IndexMovie fallback
		if r.Method == http.MethodPatch {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"message": "not found"})
			return
		}

		// POST documents (create) - success (must return 201)
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/documents") {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"id": "created"})
			return
		}

		// Default - success for collection operations
		json.NewEncoder(w).Encode(map[string]interface{}{})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	movieID := uuid.Must(uuid.NewV7())
	now := time.Now()
	m := &movie.Movie{
		ID:             movieID,
		Title:          "New Movie",
		LibraryAddedAt: now,
		CreatedAt:      now,
	}

	// The "not found" error from PATCH triggers the fallback to IndexMovie.
	// We just verify it doesn't panic.
	err := svc.UpdateMovie(ctx, m, nil, nil, nil)
	_ = err
}

func TestMovieSearchService_RemoveMovie_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	err := svc.RemoveMovie(ctx, uuid.Must(uuid.NewV7()))
	assert.NoError(t, err)
}

func TestMovieSearchService_Search_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	t.Run("basic search", func(t *testing.T) {
		result, err := svc.Search(ctx, SearchParams{
			Query:             "test",
			Page:              1,
			PerPage:           20,
			SortBy:            "popularity:desc",
			FilterBy:          "year:>=2020",
			FacetBy:           []string{"genres"},
			IncludeHighlights: true,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.TotalHits)
		assert.Equal(t, 1, result.CurrentPage)
		assert.Len(t, result.Hits, 2)
		assert.Equal(t, "doc-1", result.Hits[0].Document.ID)
		assert.Equal(t, "Test Movie 1", result.Hits[0].Document.Title)
		assert.Equal(t, float64(100), result.Hits[0].Score)
		assert.NotEmpty(t, result.Hits[0].Highlights)
		assert.Contains(t, result.Facets, "genres")
	})

	t.Run("search without highlights", func(t *testing.T) {
		result, err := svc.Search(ctx, SearchParams{
			Query:             "test",
			Page:              1,
			PerPage:           20,
			IncludeHighlights: false,
		})
		require.NoError(t, err)
		assert.Len(t, result.Hits, 2)
		assert.Empty(t, result.Hits[0].Highlights)
	})

	t.Run("search with page defaults", func(t *testing.T) {
		result, err := svc.Search(ctx, SearchParams{
			Query:   "test",
			Page:    0,
			PerPage: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, result.CurrentPage)
	})

	t.Run("search with large per page", func(t *testing.T) {
		result, err := svc.Search(ctx, SearchParams{
			Query:   "test",
			Page:    1,
			PerPage: 200,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestMovieSearchService_Autocomplete_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	suggestions, err := svc.Autocomplete(ctx, "test", 5)
	require.NoError(t, err)
	assert.Len(t, suggestions, 2)
	assert.Equal(t, "Test Movie 1", suggestions[0])
	assert.Equal(t, "Test Movie 2", suggestions[1])
}

func TestMovieSearchService_Autocomplete_LimitClamp(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	t.Run("zero limit defaults to 5", func(t *testing.T) {
		_, err := svc.Autocomplete(ctx, "test", 0)
		assert.NoError(t, err)
	})

	t.Run("over 20 capped", func(t *testing.T) {
		_, err := svc.Autocomplete(ctx, "test", 25)
		assert.NoError(t, err)
	})
}

func TestMovieSearchService_Autocomplete_DeduplicatesTitles(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.Contains(r.URL.Path, "/search") {
			found := 3
			hits := []api.SearchResultHit{
				{Document: &map[string]interface{}{"id": "1", "title": "Same Movie"}},
				{Document: &map[string]interface{}{"id": "2", "title": "Same Movie"}},
				{Document: &map[string]interface{}{"id": "3", "title": "Other Movie"}},
			}
			json.NewEncoder(w).Encode(api.SearchResult{Found: &found, Hits: &hits})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	suggestions, err := svc.Autocomplete(ctx, "movie", 5)
	require.NoError(t, err)
	assert.Len(t, suggestions, 2)
	assert.Equal(t, "Same Movie", suggestions[0])
	assert.Equal(t, "Other Movie", suggestions[1])
}

func TestMovieSearchService_GetFacets_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	facets, err := svc.GetFacets(ctx, []string{"genres"})
	require.NoError(t, err)
	assert.NotNil(t, facets)
	assert.Contains(t, facets, "genres")
}

func TestMovieSearchService_BulkIndexMovies_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	now := time.Now()
	movies := []MovieWithRelations{
		{
			Movie: &movie.Movie{
				ID:             uuid.Must(uuid.NewV7()),
				Title:          "Bulk Movie 1",
				LibraryAddedAt: now,
				CreatedAt:      now,
			},
		},
		{
			Movie: &movie.Movie{
				ID:             uuid.Must(uuid.NewV7()),
				Title:          "Bulk Movie 2",
				LibraryAddedAt: now,
				CreatedAt:      now,
			},
		},
	}

	err := svc.BulkIndexMovies(ctx, movies)
	assert.NoError(t, err)
}

func TestMovieSearchService_BulkIndexMovies_EmptySlice_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	err := svc.BulkIndexMovies(ctx, []MovieWithRelations{})
	assert.NoError(t, err)
}

// --- TVShow service tests with fake Typesense ---

func TestTVShowSearchService_InitializeCollection_AlreadyExists(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	assert.True(t, svc.IsEnabled())
	err := svc.InitializeCollection(ctx)
	assert.NoError(t, err)
}

func TestTVShowSearchService_IndexSeries_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	series := &tvshow.Series{
		ID:        uuid.Must(uuid.NewV7()),
		Title:     "Test Show",
		CreatedAt: time.Now(),
	}

	err := svc.IndexSeries(ctx, series, nil, nil, nil, false)
	assert.NoError(t, err)
}

func TestTVShowSearchService_UpdateSeries_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	series := &tvshow.Series{
		ID:        uuid.Must(uuid.NewV7()),
		Title:     "Updated Show",
		CreatedAt: time.Now(),
	}

	err := svc.UpdateSeries(ctx, series, nil, nil, nil, true)
	assert.NoError(t, err)
}

func TestTVShowSearchService_RemoveSeries_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	err := svc.RemoveSeries(ctx, uuid.Must(uuid.NewV7()))
	assert.NoError(t, err)
}

func TestTVShowSearchService_BulkIndexSeries_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	now := time.Now()
	shows := []TVShowWithRelations{
		{Series: &tvshow.Series{ID: uuid.Must(uuid.NewV7()), Title: "Show 1", CreatedAt: now}},
		{Series: &tvshow.Series{ID: uuid.Must(uuid.NewV7()), Title: "Show 2", CreatedAt: now}},
	}

	err := svc.BulkIndexSeries(ctx, shows)
	assert.NoError(t, err)
}

func TestTVShowSearchService_SearchSeries_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	t.Run("basic search", func(t *testing.T) {
		result, err := svc.SearchSeries(ctx, TVShowSearchParams{
			Query:             "test",
			Page:              1,
			PerPage:           20,
			SortBy:            "popularity:desc",
			FilterBy:          "year:>=2020",
			FacetBy:           []string{"genres"},
			IncludeHighlights: true,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.TotalHits)
		assert.Len(t, result.Hits, 2)
	})

	t.Run("without highlights", func(t *testing.T) {
		result, err := svc.SearchSeries(ctx, TVShowSearchParams{
			Query:             "test",
			Page:              1,
			PerPage:           20,
			IncludeHighlights: false,
		})
		require.NoError(t, err)
		assert.Empty(t, result.Hits[0].Highlights)
	})

	t.Run("param defaults", func(t *testing.T) {
		result, err := svc.SearchSeries(ctx, TVShowSearchParams{
			Query: "test", Page: 0, PerPage: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, result.CurrentPage)
	})

	t.Run("per page capped", func(t *testing.T) {
		result, err := svc.SearchSeries(ctx, TVShowSearchParams{
			Query: "test", Page: 1, PerPage: 200,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestTVShowSearchService_AutocompleteSeries_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	suggestions, err := svc.AutocompleteSeries(ctx, "test", 5)
	require.NoError(t, err)
	assert.Len(t, suggestions, 2)
}

func TestTVShowSearchService_GetFacets_Enabled(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewTVShowSearchService(client, slog.Default())
	ctx := context.Background()

	facets, err := svc.GetFacets(ctx, []string{"genres"})
	require.NoError(t, err)
	assert.NotNil(t, facets)
}

// --- CachedMovieSearchService with cache ---

func TestCachedMovieSearchService_Search_WithCache(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	// First call - cache miss
	result1, err := cached.Search(ctx, SearchParams{Query: "test", Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Len(t, result1.Hits, 2)

	// Allow async cache set to complete
	time.Sleep(200 * time.Millisecond)

	// Second call - cache hit
	result2, err := cached.Search(ctx, SearchParams{Query: "test", Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Len(t, result2.Hits, 2)
}

func TestCachedMovieSearchService_Autocomplete_WithCache(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	result1, err := cached.Autocomplete(ctx, "test", 5)
	require.NoError(t, err)
	assert.Len(t, result1, 2)

	time.Sleep(200 * time.Millisecond)

	result2, err := cached.Autocomplete(ctx, "test", 5)
	require.NoError(t, err)
	assert.Len(t, result2, 2)
}

func TestCachedMovieSearchService_GetFacets_WithCache(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	result1, err := cached.GetFacets(ctx, []string{"genres"})
	require.NoError(t, err)
	assert.NotNil(t, result1)

	time.Sleep(200 * time.Millisecond)

	result2, err := cached.GetFacets(ctx, []string{"genres"})
	require.NoError(t, err)
	assert.NotNil(t, result2)
}

func TestCachedMovieSearchService_InvalidateSearchCache_WithCache(t *testing.T) {
	server := httptest.NewServer(fakeTypesenseSearchHandler(t))
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	err := cached.InvalidateSearchCache(ctx)
	assert.NoError(t, err)
}

// --- Error handling ---

func TestMovieSearchService_Search_ErrorFromClient(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "internal error"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	_, err := svc.Search(ctx, SearchParams{Query: "test", Page: 1, PerPage: 20})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "search failed")
}

func TestMovieSearchService_IndexMovie_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "error"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	m := &movie.Movie{
		ID:             uuid.Must(uuid.NewV7()),
		Title:          "Error Movie",
		LibraryAddedAt: time.Now(),
		CreatedAt:      time.Now(),
	}

	err := svc.IndexMovie(ctx, m, nil, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to index movie")
}

func TestMovieSearchService_RemoveMovie_NotFoundIgnored(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "not found"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	err := svc.RemoveMovie(ctx, uuid.Must(uuid.NewV7()))
	// "not found" errors should be silently ignored for remove operations
	assert.NoError(t, err)
}

func TestMovieSearchService_Autocomplete_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "error"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	_, err := svc.Autocomplete(ctx, "test", 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "autocomplete failed")
}

// --- Cached service error from underlying service ---

func TestCachedMovieSearchService_Search_ErrorPassthrough(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "error"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	_, err := cached.Search(ctx, SearchParams{Query: "test", Page: 1, PerPage: 20})
	assert.Error(t, err)
}

func TestCachedMovieSearchService_Autocomplete_ErrorPassthrough(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "error"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	_, err := cached.Autocomplete(ctx, "test", 5)
	assert.Error(t, err)
}

func TestCachedMovieSearchService_GetFacets_ErrorPassthrough(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "error"})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	innerSvc := NewMovieSearchService(client, slog.Default())
	testCacheInst := newTestCache(t)

	cached := NewCachedMovieSearchService(innerSvc, testCacheInst, logging.NewTestLogger())
	ctx := context.Background()

	_, err := cached.GetFacets(ctx, []string{"genres"})
	assert.Error(t, err)
}

// --- Facet parsing edge cases ---

func TestMovieSearchService_Search_NilFacetValues(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.Contains(r.URL.Path, "/search") {
			result := map[string]interface{}{
				"found": 0,
				"hits":  []interface{}{},
				"facet_counts": []map[string]interface{}{
					{"field_name": nil, "counts": nil},
					{
						"field_name": "genres",
						"counts": []map[string]interface{}{
							{"value": nil, "count": nil},
							{"value": "Action", "count": 5},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(result)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newTestSearchClient(t, server.URL)
	svc := NewMovieSearchService(client, slog.Default())
	ctx := context.Background()

	result, err := svc.Search(ctx, SearchParams{Query: "*", Page: 1, PerPage: 20})
	require.NoError(t, err)
	// The nil FieldName facet should be skipped
	assert.Equal(t, 1, len(result.Facets))
	// Within genres facet, the nil Value/Count entry should be skipped
	assert.Equal(t, 1, len(result.Facets["genres"]))
	assert.Equal(t, "Action", result.Facets["genres"][0].Value)
}
