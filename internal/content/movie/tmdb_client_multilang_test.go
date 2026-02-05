package movie

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTMDbClient_GetMovieWithLanguage(t *testing.T) {
	t.Run("fetches movie in specific language", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/movie/550", r.URL.Path)
			assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
			assert.Equal(t, "de-DE", r.URL.Query().Get("language"))

			overview := "German overview..."
			movie := TMDbMovie{
				ID:               550,
				Title:            "Fight Club",
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
			}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieWithLanguage(context.Background(), 550, "de-DE")

		require.NoError(t, err)
		assert.Equal(t, 550, result.ID)
		assert.Equal(t, "Fight Club", result.Title)
		assert.NotNil(t, result.Overview)
		assert.Equal(t, "German overview...", *result.Overview)
	})

	t.Run("caches results by language", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			movie := TMDbMovie{ID: 550, Title: "Fight Club"}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		// First call
		_, err := client.GetMovieWithLanguage(context.Background(), 550, "en-US")
		require.NoError(t, err)

		// Second call should use cache
		_, err = client.GetMovieWithLanguage(context.Background(), 550, "en-US")
		require.NoError(t, err)

		assert.Equal(t, 1, callCount, "Should only call API once, second should use cache")
	})

	t.Run("different languages have separate cache entries", func(t *testing.T) {
		languagesRequested := []string{}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")
			languagesRequested = append(languagesRequested, lang)

			movie := TMDbMovie{ID: 550, Title: "Fight Club"}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		_, err := client.GetMovieWithLanguage(context.Background(), 550, "en-US")
		require.NoError(t, err)

		_, err = client.GetMovieWithLanguage(context.Background(), 550, "de-DE")
		require.NoError(t, err)

		_, err = client.GetMovieWithLanguage(context.Background(), 550, "fr-FR")
		require.NoError(t, err)

		assert.Equal(t, []string{"en-US", "de-DE", "fr-FR"}, languagesRequested)
	})
}

func TestTMDbClient_GetMovieMultiLanguage(t *testing.T) {
	t.Run("fetches movie in multiple languages", func(t *testing.T) {
		languagesRequested := []string{}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")
			languagesRequested = append(languagesRequested, lang)

			title := "Fight Club"
			overview := "English overview"

			switch lang {
			case "de-DE":
				title = "Fight Club (German)"
				overview = "German overview"
			case "fr-FR":
				title = "Fight Club (French)"
				overview = "French overview"
			}

			movie := TMDbMovie{
				ID:               550,
				Title:            title,
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
			}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieMultiLanguage(context.Background(), 550, []string{"en-US", "de-DE", "fr-FR"})

		require.NoError(t, err)
		assert.Len(t, result.Movies, 3)

		// Check language codes are mapped correctly (en-US → en, de-DE → de, etc.)
		assert.Contains(t, result.Movies, "en")
		assert.Contains(t, result.Movies, "de")
		assert.Contains(t, result.Movies, "fr")

		assert.Equal(t, "Fight Club", result.Movies["en"].Title)
		assert.Equal(t, "Fight Club (German)", result.Movies["de"].Title)
		assert.Equal(t, "Fight Club (French)", result.Movies["fr"].Title)

		assert.Equal(t, []string{"en-US", "de-DE", "fr-FR"}, languagesRequested)
	})

	t.Run("uses default language when none specified", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")
			assert.Equal(t, "en-US", lang)

			movie := TMDbMovie{ID: 550, Title: "Fight Club"}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieMultiLanguage(context.Background(), 550, nil)

		require.NoError(t, err)
		assert.Len(t, result.Movies, 1)
		assert.Contains(t, result.Movies, "en")
	})

	t.Run("continues fetching if one language fails", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			lang := r.URL.Query().Get("language")

			// Fail on German request
			if lang == "de-DE" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			movie := TMDbMovie{ID: 550, Title: "Fight Club"}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, movie)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieMultiLanguage(context.Background(), 550, []string{"en-US", "de-DE", "fr-FR"})

		require.NoError(t, err)
		assert.Len(t, result.Movies, 2, "Should have en and fr, but not de")
		assert.Contains(t, result.Movies, "en")
		assert.Contains(t, result.Movies, "fr")
		assert.NotContains(t, result.Movies, "de")
		assert.Equal(t, 3, callCount, "Should have tried all three languages")
	})

	t.Run("returns error if all languages fail", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		_, err := client.GetMovieMultiLanguage(context.Background(), 550, []string{"en-US", "de-DE"})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch movie in any language")
	})
}

func TestTMDbClient_GetMovieReleaseDates(t *testing.T) {
	t.Run("fetches release dates and age ratings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/movie/550/release_dates", r.URL.Path)
			assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))

			response := TMDbReleaseDatesResponse{
				ID: 550,
				Results: []TMDbCountryRelease{
					{
						ISO3166_1: "US",
						ReleaseDates: []TMDbReleaseDate{
							{
								Certification: "R",
								Type:          3, // Theatrical
								ReleaseDate:   "1999-10-15",
							},
						},
					},
					{
						ISO3166_1: "DE",
						ReleaseDates: []TMDbReleaseDate{
							{
								Certification: "16",
								Type:          3,
								ReleaseDate:   "1999-11-10",
							},
						},
					},
					{
						ISO3166_1: "GB",
						ReleaseDates: []TMDbReleaseDate{
							{
								Certification: "18",
								Type:          3,
								ReleaseDate:   "1999-11-12",
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieReleaseDates(context.Background(), 550)

		require.NoError(t, err)
		assert.Equal(t, 550, result.ID)
		assert.Len(t, result.Results, 3)

		// Check US rating
		usRelease := result.Results[0]
		assert.Equal(t, "US", usRelease.ISO3166_1)
		assert.Len(t, usRelease.ReleaseDates, 1)
		assert.Equal(t, "R", usRelease.ReleaseDates[0].Certification)
		assert.Equal(t, 3, usRelease.ReleaseDates[0].Type)

		// Check DE rating
		deRelease := result.Results[1]
		assert.Equal(t, "DE", deRelease.ISO3166_1)
		assert.Equal(t, "16", deRelease.ReleaseDates[0].Certification)

		// Check GB rating
		gbRelease := result.Results[2]
		assert.Equal(t, "GB", gbRelease.ISO3166_1)
		assert.Equal(t, "18", gbRelease.ReleaseDates[0].Certification)
	})

	t.Run("caches release dates", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			response := TMDbReleaseDatesResponse{
				ID:      550,
				Results: []TMDbCountryRelease{},
			}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		// First call
		_, err := client.GetMovieReleaseDates(context.Background(), 550)
		require.NoError(t, err)

		// Second call should use cache
		_, err = client.GetMovieReleaseDates(context.Background(), 550)
		require.NoError(t, err)

		assert.Equal(t, 1, callCount, "Should only call API once")
	})

	t.Run("handles empty release dates", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := TMDbReleaseDatesResponse{
				ID:      550,
				Results: []TMDbCountryRelease{},
			}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieReleaseDates(context.Background(), 550)

		require.NoError(t, err)
		assert.Equal(t, 550, result.ID)
		assert.Empty(t, result.Results)
	})

	t.Run("handles multiple release dates per country", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := TMDbReleaseDatesResponse{
				ID: 550,
				Results: []TMDbCountryRelease{
					{
						ISO3166_1: "US",
						ReleaseDates: []TMDbReleaseDate{
							{
								Certification: "",
								Type:          1, // Premiere
								ReleaseDate:   "1999-09-10",
							},
							{
								Certification: "R",
								Type:          3, // Theatrical
								ReleaseDate:   "1999-10-15",
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			writeJSON(w, response)
		}))
		defer server.Close()

		client := createTestClient(server.URL, "test-api-key")

		result, err := client.GetMovieReleaseDates(context.Background(), 550)

		require.NoError(t, err)
		assert.Len(t, result.Results[0].ReleaseDates, 2)
		assert.Equal(t, 1, result.Results[0].ReleaseDates[0].Type)
		assert.Equal(t, 3, result.Results[0].ReleaseDates[1].Type)
		assert.Equal(t, "R", result.Results[0].ReleaseDates[1].Certification)
	})
}
