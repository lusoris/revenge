package image

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	logger := logging.NewTestLogger()

	t.Run("With default config", func(t *testing.T) {
		svc, err := NewService(Config{}, logger)
		require.NoError(t, err)
		assert.NotNil(t, svc)
		assert.Equal(t, "https://image.tmdb.org/t/p", svc.config.BaseURL)
		assert.Equal(t, 7*24*time.Hour, svc.config.CacheTTL)
		assert.Equal(t, int64(10*1024*1024), svc.config.MaxSize)
	})

	t.Run("With custom config", func(t *testing.T) {
		cfg := Config{
			BaseURL:  "https://custom.example.com",
			CacheTTL: 1 * time.Hour,
			MaxSize:  5 * 1024 * 1024,
		}
		svc, err := NewService(cfg, logger)
		require.NoError(t, err)
		assert.Equal(t, "https://custom.example.com", svc.config.BaseURL)
		assert.Equal(t, 1*time.Hour, svc.config.CacheTTL)
	})

	t.Run("Creates cache directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		cacheDir := filepath.Join(tmpDir, "cache")

		cfg := Config{
			CacheDir: cacheDir,
		}
		_, err := NewService(cfg, logger)
		require.NoError(t, err)

		info, err := os.Stat(cacheDir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})
}

func TestService_GetImageURL(t *testing.T) {
	logger := logging.NewTestLogger()
	svc, _ := NewService(Config{}, logger)

	tests := []struct {
		name     string
		path     string
		size     string
		expected string
	}{
		{
			name:     "Poster with size",
			path:     "/abc123.jpg",
			size:     "w500",
			expected: "https://image.tmdb.org/t/p/w500/abc123.jpg",
		},
		{
			name:     "Backdrop original",
			path:     "/backdrop.jpg",
			size:     "original",
			expected: "https://image.tmdb.org/t/p/original/backdrop.jpg",
		},
		{
			name:     "Path without leading slash",
			path:     "profile.jpg",
			size:     "w185",
			expected: "https://image.tmdb.org/t/p/w185/profile.jpg",
		},
		{
			name:     "Empty path",
			path:     "",
			size:     "w500",
			expected: "",
		},
		{
			name:     "Full URL from fanart.tv",
			path:     "https://assets.fanart.tv/fanart/movies/550/movieposter/fight-club.jpg",
			size:     "w500",
			expected: "https://assets.fanart.tv/fanart/movies/550/movieposter/fight-club.jpg",
		},
		{
			name:     "Full URL from Radarr",
			path:     "https://artworks.thetvdb.com/banners/poster.jpg",
			size:     "w342",
			expected: "https://artworks.thetvdb.com/banners/poster.jpg",
		},
		{
			name:     "Full URL http",
			path:     "http://example.com/image.jpg",
			size:     "original",
			expected: "http://example.com/image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.GetImageURL(tt.path, tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestService_FetchImage(t *testing.T) {
	logger := logging.NewTestLogger()

	t.Run("Successful fetch", func(t *testing.T) {
		// Create test server
		imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(imageData)
		}))
		defer server.Close()

		svc, _ := NewService(Config{BaseURL: server.URL}, logger)

		data, contentType, err := svc.FetchImage(context.Background(), TypePoster, "/test.jpg", "w500")

		require.NoError(t, err)
		assert.Equal(t, imageData, data)
		assert.Equal(t, "image/jpeg", contentType)
	})

	t.Run("Empty path returns error", func(t *testing.T) {
		svc, _ := NewService(Config{}, logger)

		_, _, err := svc.FetchImage(context.Background(), TypePoster, "", "w500")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty image path")
	})

	t.Run("Server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		svc, _ := NewService(Config{BaseURL: server.URL}, logger)

		_, _, err := svc.FetchImage(context.Background(), TypePoster, "/notfound.jpg", "w500")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("Invalid content type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html>Not an image</html>"))
		}))
		defer server.Close()

		svc, _ := NewService(Config{BaseURL: server.URL}, logger)

		_, _, err := svc.FetchImage(context.Background(), TypePoster, "/test.html", "w500")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid content type")
	})

	t.Run("Image too large", func(t *testing.T) {
		largeData := make([]byte, 100)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(largeData)
		}))
		defer server.Close()

		svc, _ := NewService(Config{
			BaseURL: server.URL,
			MaxSize: 50, // Very small limit
		}, logger)

		_, _, err := svc.FetchImage(context.Background(), TypePoster, "/large.jpg", "w500")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image too large")
	})
}

func TestService_Cache(t *testing.T) {
	logger := logging.NewTestLogger()

	t.Run("Caches image on fetch", func(t *testing.T) {
		tmpDir := t.TempDir()
		callCount := 0

		imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(imageData)
		}))
		defer server.Close()

		svc, _ := NewService(Config{
			BaseURL:  server.URL,
			CacheDir: tmpDir,
		}, logger)

		// First fetch
		data1, _, err := svc.FetchImage(context.Background(), TypePoster, "/cached.jpg", "w500")
		require.NoError(t, err)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, imageData, data1)

		// Second fetch should use cache
		data2, _, err := svc.FetchImage(context.Background(), TypePoster, "/cached.jpg", "w500")
		require.NoError(t, err)
		assert.Equal(t, 1, callCount) // Still 1
		assert.Equal(t, imageData, data2)
	})
}

func TestService_ServeHTTP(t *testing.T) {
	logger := logging.NewTestLogger()

	t.Run("Valid image request", func(t *testing.T) {
		imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0}
		tmdbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(imageData)
		}))
		defer tmdbServer.Close()

		svc, _ := NewService(Config{BaseURL: tmdbServer.URL}, logger)

		req := httptest.NewRequest("GET", "/images/poster/w500/test.jpg", nil)
		rr := httptest.NewRecorder()

		svc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "image/jpeg", rr.Header().Get("Content-Type"))
		assert.Equal(t, imageData, rr.Body.Bytes())
	})

	t.Run("Invalid path format", func(t *testing.T) {
		svc, _ := NewService(Config{}, logger)

		req := httptest.NewRequest("GET", "/images/poster", nil)
		rr := httptest.NewRecorder()

		svc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Invalid image type", func(t *testing.T) {
		svc, _ := NewService(Config{}, logger)

		req := httptest.NewRequest("GET", "/images/invalid/w500/test.jpg", nil)
		rr := httptest.NewRecorder()

		svc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Invalid size", func(t *testing.T) {
		svc, _ := NewService(Config{}, logger)

		req := httptest.NewRequest("GET", "/images/poster/w9999/test.jpg", nil)
		rr := httptest.NewRecorder()

		svc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("isValidImageType", func(t *testing.T) {
		assert.True(t, isValidImageType("image/jpeg"))
		assert.True(t, isValidImageType("image/png"))
		assert.True(t, isValidImageType("image/gif"))
		assert.True(t, isValidImageType("image/webp"))
		assert.True(t, isValidImageType("image/svg+xml"))
		assert.False(t, isValidImageType("text/html"))
		assert.False(t, isValidImageType("application/json"))
	})

	t.Run("isValidType", func(t *testing.T) {
		assert.True(t, isValidType(TypePoster))
		assert.True(t, isValidType(TypeBackdrop))
		assert.True(t, isValidType(TypeProfile))
		assert.True(t, isValidType(TypeLogo))
		assert.False(t, isValidType("invalid"))
	})

	t.Run("isValidSize poster", func(t *testing.T) {
		assert.True(t, isValidSize(TypePoster, "w185"))
		assert.True(t, isValidSize(TypePoster, "w342"))
		assert.True(t, isValidSize(TypePoster, "w500"))
		assert.True(t, isValidSize(TypePoster, "original"))
		assert.False(t, isValidSize(TypePoster, "w9999"))
	})

	t.Run("isValidSize backdrop", func(t *testing.T) {
		assert.True(t, isValidSize(TypeBackdrop, "w300"))
		assert.True(t, isValidSize(TypeBackdrop, "w780"))
		assert.True(t, isValidSize(TypeBackdrop, "w1280"))
		assert.True(t, isValidSize(TypeBackdrop, "original"))
		assert.False(t, isValidSize(TypeBackdrop, "w185"))
	})

	t.Run("getContentType", func(t *testing.T) {
		assert.Equal(t, "image/jpeg", getContentType("test.jpg"))
		assert.Equal(t, "image/jpeg", getContentType("test.jpeg"))
		assert.Equal(t, "image/png", getContentType("test.png"))
		assert.Equal(t, "image/gif", getContentType("test.gif"))
		assert.Equal(t, "image/webp", getContentType("test.webp"))
		assert.Equal(t, "image/svg+xml", getContentType("test.svg"))
		assert.Equal(t, "application/octet-stream", getContentType("test.unknown"))
	})
}

func TestService_getDefaultSize(t *testing.T) {
	logger := logging.NewTestLogger()
	svc, _ := NewService(Config{}, logger)

	assert.Equal(t, SizePosterMedium, svc.getDefaultSize(TypePoster))
	assert.Equal(t, SizeBackdropMedium, svc.getDefaultSize(TypeBackdrop))
	assert.Equal(t, SizeProfileMedium, svc.getDefaultSize(TypeProfile))
	assert.Equal(t, SizePosterMedium, svc.getDefaultSize(TypeLogo))
	assert.Equal(t, "original", svc.getDefaultSize("unknown"))
}
