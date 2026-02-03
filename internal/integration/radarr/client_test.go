package radarr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeJSON is a helper to write JSON responses in tests with proper Content-Type.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	//nolint:errcheck // test helper
	json.NewEncoder(w).Encode(v)
}

func TestNewClient(t *testing.T) {
	client := NewClient(Config{
		BaseURL:   "http://localhost:7878",
		APIKey:    "test-api-key",
		RateLimit: 10.0,
		CacheTTL:  5 * time.Minute,
		Timeout:   30 * time.Second,
	})

	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:7878", client.baseURL)
	assert.Equal(t, "test-api-key", client.apiKey)
}

func TestNewClient_Defaults(t *testing.T) {
	client := NewClient(Config{
		BaseURL: "http://localhost:7878",
		APIKey:  "test-api-key",
	})

	assert.NotNil(t, client)
	// Check that defaults are applied
	assert.Equal(t, 5*time.Minute, client.cacheTTL)
}

func TestClient_GetSystemStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/system/status", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-Api-Key"))

		status := SystemStatus{
			AppName: "Radarr",
			Version: "5.0.0.0",
			Branch:  "main",
		}
		writeJSON(w, status)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	status, err := client.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Radarr", status.AppName)
	assert.Equal(t, "5.0.0.0", status.Version)
}

func TestClient_GetAllMovies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/movie", r.URL.Path)

		movies := []Movie{
			{
				ID:      1,
				Title:   "Inception",
				Year:    2010,
				TMDbID:  27205,
				HasFile: true,
			},
			{
				ID:      2,
				Title:   "The Matrix",
				Year:    1999,
				TMDbID:  603,
				HasFile: true,
			},
		}
		writeJSON(w, movies)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	movies, err := client.GetAllMovies(context.Background())
	require.NoError(t, err)
	require.Len(t, movies, 2)
	assert.Equal(t, "Inception", movies[0].Title)
	assert.Equal(t, 2010, movies[0].Year)
}

func TestClient_GetMovie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/movie/1", r.URL.Path)

		movie := Movie{
			ID:       1,
			Title:    "Inception",
			Year:     2010,
			TMDbID:   27205,
			Overview: "A thief who steals corporate secrets...",
		}
		writeJSON(w, movie)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	movie, err := client.GetMovie(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "Inception", movie.Title)
	assert.Equal(t, 27205, movie.TMDbID)
}

func TestClient_GetMovie_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	movie, err := client.GetMovie(context.Background(), 999)
	assert.ErrorIs(t, err, ErrMovieNotFound)
	assert.Nil(t, movie)
}

func TestClient_GetMovieByTMDbID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/movie", r.URL.Path)
		assert.Equal(t, "27205", r.URL.Query().Get("tmdbId"))

		movies := []Movie{
			{
				ID:     1,
				Title:  "Inception",
				Year:   2010,
				TMDbID: 27205,
			},
		}
		writeJSON(w, movies)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	movie, err := client.GetMovieByTMDbID(context.Background(), 27205)
	require.NoError(t, err)
	assert.Equal(t, "Inception", movie.Title)
}

func TestClient_GetMovieByTMDbID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, []Movie{})
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	movie, err := client.GetMovieByTMDbID(context.Background(), 99999)
	assert.ErrorIs(t, err, ErrMovieNotFound)
	assert.Nil(t, movie)
}

func TestClient_GetQualityProfiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/qualityprofile", r.URL.Path)

		profiles := []QualityProfile{
			{
				ID:   1,
				Name: "HD-1080p",
			},
			{
				ID:   2,
				Name: "Ultra-HD",
			},
		}
		writeJSON(w, profiles)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	profiles, err := client.GetQualityProfiles(context.Background())
	require.NoError(t, err)
	require.Len(t, profiles, 2)
	assert.Equal(t, "HD-1080p", profiles[0].Name)
}

func TestClient_GetRootFolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/rootfolder", r.URL.Path)

		folders := []RootFolder{
			{
				ID:         1,
				Path:       "/movies",
				Accessible: true,
				FreeSpace:  1000000000000,
			},
		}
		writeJSON(w, folders)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	folders, err := client.GetRootFolders(context.Background())
	require.NoError(t, err)
	require.Len(t, folders, 1)
	assert.Equal(t, "/movies", folders[0].Path)
}

func TestClient_AddMovie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/movie", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		movie := Movie{
			ID:     1,
			Title:  "New Movie",
			Year:   2024,
			TMDbID: 12345,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		//nolint:errcheck
		json.NewEncoder(w).Encode(movie)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	movie, err := client.AddMovie(context.Background(), AddMovieRequest{
		Title:            "New Movie",
		QualityProfileID: 1,
		TMDbID:           12345,
		RootFolderPath:   "/movies",
		Monitored:        true,
	})
	require.NoError(t, err)
	assert.Equal(t, "New Movie", movie.Title)
}

func TestClient_DeleteMovie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/movie/1", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "true", r.URL.Query().Get("deleteFiles"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	err := client.DeleteMovie(context.Background(), 1, true, false)
	require.NoError(t, err)
}

func TestClient_RefreshMovie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/command", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		command := Command{
			ID:     1,
			Name:   "RefreshMovie",
			Status: "started",
		}
		writeJSON(w, command)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	})

	command, err := client.RefreshMovie(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "RefreshMovie", command.Name)
	assert.Equal(t, "started", command.Status)
}

func TestClient_IsHealthy(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			status := SystemStatus{Version: "5.0.0.0"}
			writeJSON(w, status)
		}))
		defer server.Close()

		client := NewClient(Config{
			BaseURL: server.URL,
			APIKey:  "test-api-key",
		})

		assert.True(t, client.IsHealthy(context.Background()))
	})

	t.Run("unhealthy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(Config{
			BaseURL: server.URL,
			APIKey:  "test-api-key",
		})

		assert.False(t, client.IsHealthy(context.Background()))
	})
}

func TestClient_Caching(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		status := SystemStatus{Version: "5.0.0.0"}
		writeJSON(w, status)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL:  server.URL,
		APIKey:   "test-api-key",
		CacheTTL: 1 * time.Hour,
	})

	// First call
	_, err := client.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Second call should use cache
	_, err = client.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, callCount) // Still 1 because of cache

	// Clear cache
	client.ClearCache()

	// Third call should hit server again
	_, err = client.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}
