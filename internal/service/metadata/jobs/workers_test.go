package jobs

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/revenge/internal/infra/image"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestImageService creates an image.Service backed by the given HTTP server
// with a temp cache directory.
func newTestImageService(t *testing.T, serverURL string) *image.Service {
	t.Helper()

	svc, err := image.NewService(image.Config{
		BaseURL:  serverURL,
		CacheDir: t.TempDir(),
	}, slog.Default())
	require.NoError(t, err)

	return svc
}

func TestDownloadImageWorker_Work(t *testing.T) {
	t.Parallel()

	t.Run("downloads and caches image", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/w500/abc123.jpg", r.URL.Path)
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write([]byte("fake-jpeg-data"))
		}))
		defer srv.Close()

		worker := NewDownloadImageWorker(newTestImageService(t, srv.URL), slog.Default())

		err := worker.Work(context.Background(), &river.Job[DownloadImageArgs]{
			Args: DownloadImageArgs{
				ContentType: "movie",
				ContentID:   "550e8400-e29b-41d4-a716-446655440000",
				ImageType:   image.TypePoster,
				Path:        "/abc123.jpg",
				Size:        image.SizePosterLarge,
			},
		})

		assert.NoError(t, err)
	})

	t.Run("skips empty path", func(t *testing.T) {
		t.Parallel()

		// Server should never be called.
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("unexpected request to image server")
		}))
		defer srv.Close()

		worker := NewDownloadImageWorker(newTestImageService(t, srv.URL), slog.Default())

		err := worker.Work(context.Background(), &river.Job[DownloadImageArgs]{
			Args: DownloadImageArgs{
				ContentType: "movie",
				ContentID:   "550e8400-e29b-41d4-a716-446655440000",
				ImageType:   image.TypePoster,
				Path:        "",
				Size:        image.SizePosterLarge,
			},
		})

		assert.NoError(t, err)
	})

	t.Run("returns error on fetch failure", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		worker := NewDownloadImageWorker(newTestImageService(t, srv.URL), slog.Default())

		err := worker.Work(context.Background(), &river.Job[DownloadImageArgs]{
			Args: DownloadImageArgs{
				ContentType: "tvshow",
				ContentID:   "550e8400-e29b-41d4-a716-446655440001",
				ImageType:   image.TypeBackdrop,
				Path:        "/missing.jpg",
				Size:        image.SizeBackdropLarge,
			},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "download image")
		assert.Contains(t, err.Error(), "backdrop")
	})

	t.Run("returns error on invalid content type", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html>not an image</html>"))
		}))
		defer srv.Close()

		worker := NewDownloadImageWorker(newTestImageService(t, srv.URL), slog.Default())

		err := worker.Work(context.Background(), &river.Job[DownloadImageArgs]{
			Args: DownloadImageArgs{
				ContentType: "movie",
				ContentID:   "550e8400-e29b-41d4-a716-446655440002",
				ImageType:   image.TypePoster,
				Path:        "/poster.jpg",
				Size:        image.SizePosterMedium,
			},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "download image")
	})
}
