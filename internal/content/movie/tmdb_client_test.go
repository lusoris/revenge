package movie

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestNewTMDbClient(t *testing.T) {
	t.Run("With default settings", func(t *testing.T) {
		client := NewTMDbClient(TMDbConfig{
			APIKey: "test-api-key",
		})

		assert.NotNil(t, client)
		assert.Equal(t, "test-api-key", client.apiKey)
		assert.Equal(t, 24*time.Hour, client.cacheTTL)
	})

	t.Run("With custom settings", func(t *testing.T) {
		client := NewTMDbClient(TMDbConfig{
			APIKey:    "test-api-key",
			RateLimit: rate.Limit(10.0),
			CacheTTL:  1 * time.Hour,
		})

		assert.NotNil(t, client)
		assert.Equal(t, "test-api-key", client.apiKey)
		assert.Equal(t, 1*time.Hour, client.cacheTTL)
	})

	t.Run("With proxy URL", func(t *testing.T) {
		client := NewTMDbClient(TMDbConfig{
			APIKey:   "test-api-key",
			ProxyURL: "http://proxy.example.com:8080",
		})

		assert.NotNil(t, client)
	})
}

func TestTMDbClient_SearchMovies(t *testing.T) {
	t.Run("Successful search", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/movie", r.URL.Path)
			assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
			assert.Equal(t, "The Matrix", r.URL.Query().Get("query"))

			response := TMDbSearchResponse{
				Page:         1,
				TotalResults: 1,
				TotalPages:   1,
				Results: []TMDbSearchResult{
					{
						ID:          603,
						Title:       "The Matrix",
						ReleaseDate: "1999-03-31",
						Overview:    "A computer hacker learns...",
						Popularity:  50.5,
						VoteAverage: 8.7,
						VoteCount:   20000,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		result, err := client.SearchMovies(ctx, "The Matrix", nil)

		require.NoError(t, err)
		assert.Equal(t, 1, result.TotalResults)
		assert.Len(t, result.Results, 1)
		assert.Equal(t, 603, result.Results[0].ID)
		assert.Equal(t, "The Matrix", result.Results[0].Title)
	})

	t.Run("Search with year filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "1999", r.URL.Query().Get("year"))

			response := TMDbSearchResponse{
				Page:         1,
				TotalResults: 1,
				TotalPages:   1,
				Results: []TMDbSearchResult{
					{
						ID:          603,
						Title:       "The Matrix",
						ReleaseDate: "1999-03-31",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()
		year := 1999

		result, err := client.SearchMovies(ctx, "The Matrix", &year)

		require.NoError(t, err)
		assert.Len(t, result.Results, 1)
	})

	t.Run("API error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(TMDbError{
				StatusMessage: "Invalid API key",
				StatusCode:    7,
			})
		}))
		defer server.Close()

		client := createTestClient(server.URL, "invalid-key")
		ctx := context.Background()

		result, err := client.SearchMovies(ctx, "test", nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Invalid API key")
	})

	t.Run("Empty results", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := TMDbSearchResponse{
				Page:         1,
				TotalResults: 0,
				TotalPages:   0,
				Results:      []TMDbSearchResult{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		result, err := client.SearchMovies(ctx, "nonexistent movie xyz", nil)

		require.NoError(t, err)
		assert.Equal(t, 0, result.TotalResults)
		assert.Empty(t, result.Results)
	})
}

func TestTMDbClient_GetMovie(t *testing.T) {
	t.Run("Successful get movie", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/movie/603", r.URL.Path)
			assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))

			overview := "A computer hacker learns about the true nature of his reality."
			tagline := "Welcome to the Real World"
			posterPath := "/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg"
			backdropPath := "/fNG7i7RqMErkcqhohV2a6cV1Ehy.jpg"
			runtime := 136

			movie := TMDbMovie{
				ID:           603,
				Title:        "The Matrix",
				Overview:     &overview,
				ReleaseDate:  "1999-03-31",
				Runtime:      &runtime,
				VoteAverage:  8.7,
				VoteCount:    20000,
				Popularity:   50.5,
				Status:       "Released",
				Tagline:      &tagline,
				PosterPath:   &posterPath,
				BackdropPath: &backdropPath,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		result, err := client.GetMovie(ctx, 603)

		require.NoError(t, err)
		assert.Equal(t, 603, result.ID)
		assert.Equal(t, "The Matrix", result.Title)
		assert.NotNil(t, result.Runtime)
		assert.Equal(t, 136, *result.Runtime)
		assert.NotNil(t, result.Tagline)
		assert.Equal(t, "Welcome to the Real World", *result.Tagline)
	})

	t.Run("Movie not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(TMDbError{
				StatusMessage: "The resource you requested could not be found.",
				StatusCode:    34,
			})
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		result, err := client.GetMovie(ctx, 9999999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTMDbClient_GetMovieCredits(t *testing.T) {
	t.Run("Successful get credits", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/movie/603/credits", r.URL.Path)

			profilePath1 := "/4D0PpNI0kmP58hgrwGC3wCjxhnm.jpg"
			profilePath2 := "/7uvJSb6Y2efMG3L1QWfpfUJSxUB.jpg"
			profilePath3 := "/hb5y7FdNqJCYPVQLUUJWKF2YxNe.jpg"
			profilePath4 := "/rVFvYQtPqQNQHvGpg9FuIMB4kHs.jpg"

			credits := TMDbCredits{
				ID: 603,
				Cast: []CastMember{
					{
						ID:          6384,
						Name:        "Keanu Reeves",
						Character:   "Thomas A. Anderson / Neo",
						ProfilePath: &profilePath1,
						Order:       0,
					},
					{
						ID:          2975,
						Name:        "Laurence Fishburne",
						Character:   "Morpheus",
						ProfilePath: &profilePath2,
						Order:       1,
					},
				},
				Crew: []CrewMember{
					{
						ID:          905,
						Name:        "Lilly Wachowski",
						Job:         "Director",
						Department:  "Directing",
						ProfilePath: &profilePath3,
					},
					{
						ID:          9340,
						Name:        "Lana Wachowski",
						Job:         "Director",
						Department:  "Directing",
						ProfilePath: &profilePath4,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(credits)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		result, err := client.GetMovieCredits(ctx, 603)

		require.NoError(t, err)
		assert.Equal(t, 603, result.ID)
		assert.Len(t, result.Cast, 2)
		assert.Equal(t, "Keanu Reeves", result.Cast[0].Name)
		assert.Equal(t, "Thomas A. Anderson / Neo", result.Cast[0].Character)
		assert.Len(t, result.Crew, 2)
		assert.Equal(t, "Lilly Wachowski", result.Crew[0].Name)
	})
}

func TestTMDbClient_Cache(t *testing.T) {
	t.Run("Cache hit on second request", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			movie := TMDbMovie{
				ID:    603,
				Title: "The Matrix",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		// First request
		result1, err := client.GetMovie(ctx, 603)
		require.NoError(t, err)
		assert.Equal(t, 603, result1.ID)
		assert.Equal(t, 1, callCount)

		// Second request (should use cache)
		result2, err := client.GetMovie(ctx, 603)
		require.NoError(t, err)
		assert.Equal(t, 603, result2.ID)
		assert.Equal(t, 1, callCount) // Still 1, used cache
	})

	t.Run("ClearCache removes cached items", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			movie := TMDbMovie{
				ID:    603,
				Title: "The Matrix",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")
		ctx := context.Background()

		// First request
		_, err := client.GetMovie(ctx, 603)
		require.NoError(t, err)
		assert.Equal(t, 1, callCount)

		// Clear cache
		client.ClearCache()

		// Third request (should hit API again)
		_, err = client.GetMovie(ctx, 603)
		require.NoError(t, err)
		assert.Equal(t, 2, callCount)
	})
}

func TestTMDbClient_ImageURL(t *testing.T) {
	client := NewTMDbClient(TMDbConfig{APIKey: "test"})

	t.Run("Build poster URL", func(t *testing.T) {
		url := client.GetImageURL("/abc123.jpg", "w500")
		assert.Equal(t, "https://image.tmdb.org/t/p/w500/abc123.jpg", url)
	})

	t.Run("Build backdrop URL", func(t *testing.T) {
		url := client.GetImageURL("/backdrop.jpg", "original")
		assert.Equal(t, "https://image.tmdb.org/t/p/original/backdrop.jpg", url)
	})

	t.Run("Handle empty path", func(t *testing.T) {
		url := client.GetImageURL("", "w500")
		assert.Equal(t, "", url)
	})
}

// Helper function to create test client with custom base URL
func createTestClient(baseURL, apiKey string) *TMDbClient {
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(5 * time.Second)

	return &TMDbClient{
		client:      client,
		apiKey:      apiKey,
		rateLimiter: rate.NewLimiter(rate.Limit(100), 10), // High limit for tests
		cacheTTL:    1 * time.Hour,
	}
}
