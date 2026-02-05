package movie

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataService_GetMovieByTMDbIDMultiLanguage(t *testing.T) {
	t.Run("fetches movie with multiple languages and age ratings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")

			if r.URL.Path == "/movie/550/release_dates" {
				// Return age ratings
				response := TMDbReleaseDatesResponse{
					ID: 550,
					Results: []TMDbCountryRelease{
						{
							ISO3166_1: "US",
							ReleaseDates: []TMDbReleaseDate{
								{Certification: "R", Type: 3},
							},
						},
						{
							ISO3166_1: "DE",
							ReleaseDates: []TMDbReleaseDate{
								{Certification: "16", Type: 3},
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			// Return movie based on language
			title := "Fight Club"
			overview := "A depressed man suffering from insomnia..."
			tagline := "Mischief. Mayhem. Soap."

			switch lang {
			case "de-DE":
				overview = "Ein depressiver Mann, der unter Schlaflosigkeit leidet..."
				tagline = "Unfug. Chaos. Seife."
			case "fr-FR":
				overview = "Un homme déprimé souffrant d'insomnie..."
				tagline = "Espièglerie. Chaos. Savon."
			}

			movie := TMDbMovie{
				ID:               550,
				Title:            title,
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
				Tagline:          &tagline,
				ReleaseDate:      "1999-10-15",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		config := TMDbConfig{
			APIKey: "test-key",
		}
		service := NewMetadataService(config)
		service.client.client.SetBaseURL(server.URL)

		movie, err := service.GetMovieByTMDbIDMultiLanguage(context.Background(), 550, []string{"en-US", "de-DE", "fr-FR"})

		require.NoError(t, err)
		require.NotNil(t, movie)

		// Check basic fields
		assert.Equal(t, "Fight Club", movie.Title)
		assert.Equal(t, int32(550), *movie.TMDbID)

		// Check multi-language titles
		assert.Len(t, movie.TitlesI18n, 3)
		assert.Equal(t, "Fight Club", movie.TitlesI18n["en"])
		assert.Equal(t, "Fight Club", movie.TitlesI18n["de"])
		assert.Equal(t, "Fight Club", movie.TitlesI18n["fr"])

		// Check multi-language overviews
		assert.Len(t, movie.OverviewsI18n, 3)
		assert.Contains(t, movie.OverviewsI18n["en"], "insomnia")
		assert.Contains(t, movie.OverviewsI18n["de"], "Schlaflosigkeit")
		assert.Contains(t, movie.OverviewsI18n["fr"], "insomnie")

		// Check multi-language taglines
		assert.Len(t, movie.TaglinesI18n, 3)
		assert.Equal(t, "Mischief. Mayhem. Soap.", movie.TaglinesI18n["en"])
		assert.Equal(t, "Unfug. Chaos. Seife.", movie.TaglinesI18n["de"])
		assert.Equal(t, "Espièglerie. Chaos. Savon.", movie.TaglinesI18n["fr"])

		// Check age ratings
		assert.Len(t, movie.AgeRatings, 2)
		assert.Equal(t, "R", movie.AgeRatings["US"]["MPAA"])
		assert.Equal(t, "16", movie.AgeRatings["DE"]["FSK"])
	})

	t.Run("uses default languages when none specified", func(t *testing.T) {
		languagesRequested := []string{}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")
			if lang != "" {
				languagesRequested = append(languagesRequested, lang)
			}

			if r.URL.Path == "/movie/550/release_dates" {
				response := TMDbReleaseDatesResponse{ID: 550, Results: []TMDbCountryRelease{}}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			overview := "Overview"
			movie := TMDbMovie{
				ID:               550,
				Title:            "Fight Club",
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		config := TMDbConfig{APIKey: "test-key"}
		service := NewMetadataService(config)
		service.client.client.SetBaseURL(server.URL)

		movie, err := service.GetMovieByTMDbIDMultiLanguage(context.Background(), 550, nil)

		require.NoError(t, err)
		require.NotNil(t, movie)

		// Should have requested default languages
		assert.Contains(t, languagesRequested, "en-US")
		assert.Contains(t, languagesRequested, "de-DE")
		assert.Contains(t, languagesRequested, "fr-FR")
		assert.Contains(t, languagesRequested, "es-ES")
		assert.Contains(t, languagesRequested, "ja-JP")
	})

	t.Run("continues without age ratings if fetch fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/movie/550/release_dates" {
				// Fail age ratings request
				w.WriteHeader(http.StatusNotFound)
				return
			}

			overview := "Overview"
			movie := TMDbMovie{
				ID:               550,
				Title:            "Fight Club",
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		config := TMDbConfig{APIKey: "test-key"}
		service := NewMetadataService(config)
		service.client.client.SetBaseURL(server.URL)

		movie, err := service.GetMovieByTMDbIDMultiLanguage(context.Background(), 550, []string{"en-US"})

		// Should succeed without age ratings
		require.NoError(t, err)
		require.NotNil(t, movie)
		assert.Nil(t, movie.AgeRatings)
	})

	t.Run("returns error when English is not available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")

			// Only return 404 for English
			if lang == "en-US" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			overview := "Overview"
			movie := TMDbMovie{
				ID:               550,
				Title:            "Fight Club",
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		config := TMDbConfig{APIKey: "test-key"}
		service := NewMetadataService(config)
		service.client.client.SetBaseURL(server.URL)

		movie, err := service.GetMovieByTMDbIDMultiLanguage(context.Background(), 550, []string{"en-US", "de-DE"})

		require.Error(t, err)
		assert.Nil(t, movie)
		assert.Contains(t, err.Error(), "English missing")
	})
}

func TestMetadataService_EnrichMovieWithLanguages(t *testing.T) {
	t.Run("enriches movie with multi-language metadata", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.URL.Query().Get("language")

			if r.URL.Path == "/movie/550/release_dates" {
				response := TMDbReleaseDatesResponse{
					ID: 550,
					Results: []TMDbCountryRelease{
						{
							ISO3166_1: "US",
							ReleaseDates: []TMDbReleaseDate{
								{Certification: "R", Type: 3},
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			overview := "English overview"
			tagline := "English tagline"

			switch lang {
			case "de-DE":
				overview = "German overview"
				tagline = "German tagline"
			}

			movie := TMDbMovie{
				ID:               550,
				Title:            "Fight Club",
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
				Tagline:          &tagline,
				ReleaseDate:      "1999-10-15",
				Runtime:          intPtr(139),
				VoteAverage:      8.4,
				VoteCount:        26000,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		config := TMDbConfig{APIKey: "test-key"}
		service := NewMetadataService(config)
		service.client.client.SetBaseURL(server.URL)

		tmdbID := int32(550)
		movie := &Movie{
			TMDbID: &tmdbID,
			Title:  "Original Title",
		}

		err := service.EnrichMovieWithLanguages(context.Background(), movie, []string{"en-US", "de-DE"})

		require.NoError(t, err)

		// Check basic fields updated
		assert.Equal(t, "Fight Club", movie.Title)
		assert.NotNil(t, movie.Runtime)
		assert.Equal(t, int32(139), *movie.Runtime)

		// Check multi-language fields populated
		assert.Len(t, movie.TitlesI18n, 2)
		assert.Len(t, movie.OverviewsI18n, 2)
		assert.Len(t, movie.TaglinesI18n, 2)

		assert.Equal(t, "English overview", movie.OverviewsI18n["en"])
		assert.Equal(t, "German overview", movie.OverviewsI18n["de"])

		assert.Equal(t, "English tagline", movie.TaglinesI18n["en"])
		assert.Equal(t, "German tagline", movie.TaglinesI18n["de"])

		// Check age ratings populated
		assert.Len(t, movie.AgeRatings, 1)
		assert.Equal(t, "R", movie.AgeRatings["US"]["MPAA"])
	})

	t.Run("returns error when movie has no TMDb ID", func(t *testing.T) {
		config := TMDbConfig{APIKey: "test-key"}
		service := NewMetadataService(config)

		movie := &Movie{
			Title: "No TMDb ID",
		}

		err := service.EnrichMovieWithLanguages(context.Background(), movie, nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no TMDb ID")
	})

	t.Run("EnrichMovie uses default behavior with multi-language", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/movie/550/release_dates" {
				response := TMDbReleaseDatesResponse{ID: 550, Results: []TMDbCountryRelease{}}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}

			overview := "Overview"
			movie := TMDbMovie{
				ID:               550,
				Title:            "Fight Club",
				OriginalTitle:    "Fight Club",
				OriginalLanguage: "en",
				Overview:         &overview,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}))
		defer server.Close()

		config := TMDbConfig{APIKey: "test-key"}
		service := NewMetadataService(config)
		service.client.client.SetBaseURL(server.URL)

		tmdbID := int32(550)
		movie := &Movie{
			TMDbID: &tmdbID,
			Title:  "Original Title",
		}

		err := service.EnrichMovie(context.Background(), movie)

		require.NoError(t, err)
		assert.Equal(t, "Fight Club", movie.Title)
		// Should have i18n fields populated from default languages
		assert.NotEmpty(t, movie.TitlesI18n)
	})
}
