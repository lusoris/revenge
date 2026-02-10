package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	ogenmw "github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/middleware"
)

func TestCacheControlMiddleware_CatalogEndpoints(t *testing.T) {
	t.Parallel()

	cacheable := []string{
		"listMovies", "searchMovies", "getMovie", "getMovieCast",
		"listTVShows", "getTVShow", "getTVShowSeasons", "getTVShowEpisodes",
		"getMovieMetadata", "getPersonMetadata",
		"getSearchFacets", "listLibraries", "getLibrary",
		"listMetadataProviders", "getProxiedImage",
	}

	for _, op := range cacheable {
		t.Run(op, func(t *testing.T) {
			t.Parallel()
			cc := callCacheControlMiddleware(t, http.MethodGet, op)
			assert.Equal(t, "private, max-age=60", cc)
		})
	}
}

func TestCacheControlMiddleware_UserSpecificEndpoints(t *testing.T) {
	t.Parallel()

	userSpecific := []string{
		"getContinueWatching", "getWatchHistory", "getUserMovieStats",
		"getWatchProgress", "getCurrentUser", "listSessions",
		"listAPIKeys", "getMFAStatus", "getUserPreferences",
		"getTVContinueWatching", "getUserTVStats",
		"getTVShowWatchStats", "getTVShowNextEpisode",
		"getTVEpisodeProgress",
	}

	for _, op := range userSpecific {
		t.Run(op, func(t *testing.T) {
			t.Parallel()
			cc := callCacheControlMiddleware(t, http.MethodGet, op)
			assert.Equal(t, "private, no-store", cc)
		})
	}
}

func TestCacheControlMiddleware_HealthProbes(t *testing.T) {
	t.Parallel()

	for _, op := range []string{"getLiveness", "getReadiness", "getStartup"} {
		t.Run(op, func(t *testing.T) {
			t.Parallel()
			cc := callCacheControlMiddleware(t, http.MethodGet, op)
			assert.Equal(t, "no-cache", cc)
		})
	}
}

func TestCacheControlMiddleware_Mutations(t *testing.T) {
	t.Parallel()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, m := range methods {
		t.Run(m, func(t *testing.T) {
			t.Parallel()
			cc := callCacheControlMiddleware(t, m, "createLibrary")
			assert.Equal(t, "no-store", cc)
		})
	}
}

func TestCacheControlMiddleware_NoResponseWriter(t *testing.T) {
	t.Parallel()

	mw := middleware.CacheControlMiddleware()

	raw := httptest.NewRequest(http.MethodGet, "/test", nil)
	req := ogenmw.Request{
		Context:     raw.Context(), // no ResponseWriter in context
		OperationID: "listMovies",
		Raw:         raw,
	}

	called := false
	next := func(req ogenmw.Request) (ogenmw.Response, error) {
		called = true
		return ogenmw.Response{}, nil
	}

	_, err := mw(req, next)
	require.NoError(t, err)
	assert.True(t, called, "next should be called even without ResponseWriter")
}

// callCacheControlMiddleware uses the ResponseWriterMiddleware to properly
// inject the ResponseWriter into context, then runs the ogen middleware.
func callCacheControlMiddleware(t *testing.T, method, operationID string) string {
	t.Helper()

	mw := middleware.CacheControlMiddleware()
	rec := httptest.NewRecorder()

	// Use the real ResponseWriterMiddleware to inject the writer
	var capturedCtx context.Context
	inner := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	})
	middleware.ResponseWriterMiddleware(inner).ServeHTTP(rec, httptest.NewRequest(method, "/test", nil))

	raw := httptest.NewRequest(method, "/test", nil)
	raw = raw.WithContext(capturedCtx)

	req := ogenmw.Request{
		Context:     capturedCtx,
		OperationID: operationID,
		Raw:         raw,
	}

	next := func(req ogenmw.Request) (ogenmw.Response, error) {
		return ogenmw.Response{}, nil
	}

	_, err := mw(req, next)
	require.NoError(t, err)

	return rec.Header().Get("Cache-Control")
}
