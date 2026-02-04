package health_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/health"
)

func newTestHandler(startupComplete bool) *health.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pool := &pgxpool.Pool{}
	service := health.NewService(logger, pool)
	if startupComplete {
		service.MarkStartupComplete()
	}
	return health.NewHandler(service)
}

func TestNewHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pool := &pgxpool.Pool{}
	service := health.NewService(logger, pool)

	handler := health.NewHandler(service)
	require.NotNil(t, handler)
}

func TestHandler_HandleLiveness(t *testing.T) {
	handler := newTestHandler(false)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()

	handler.HandleLiveness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Liveness should always return 200
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body health.Response
	err := json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	assert.Equal(t, health.StatusHealthy, body.Status)
	assert.NotEmpty(t, body.Message)
}

func TestHandler_HandleStartup(t *testing.T) {
	t.Run("before startup complete", func(t *testing.T) {
		handler := newTestHandler(false)

		req := httptest.NewRequest(http.MethodGet, "/health/startup", nil)
		w := httptest.NewRecorder()

		handler.HandleStartup(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should return 503 before startup is complete
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var body health.Response
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, health.StatusUnhealthy, body.Status)
	})

	t.Run("after startup complete", func(t *testing.T) {
		handler := newTestHandler(true)

		req := httptest.NewRequest(http.MethodGet, "/health/startup", nil)
		w := httptest.NewRecorder()

		handler.HandleStartup(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should return 200 after startup is complete
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body health.Response
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, health.StatusHealthy, body.Status)
	})
}

func TestHandler_HandleReadiness(t *testing.T) {
	t.Run("before startup complete", func(t *testing.T) {
		handler := newTestHandler(false)

		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		w := httptest.NewRecorder()

		handler.HandleReadiness(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should return 503 before startup is complete
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		var body health.Response
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, health.StatusUnhealthy, body.Status)
	})

	// Note: Testing with a real database requires an integration test
	// The "after startup complete" case would panic with a nil pool
}

func TestHandler_HandleFull(t *testing.T) {
	t.Run("before startup complete", func(t *testing.T) {
		handler := newTestHandler(false)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.HandleFull(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Should return 503 when not fully healthy
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var body health.Response
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		// Should have checks array
		assert.NotNil(t, body.Checks)
		assert.Greater(t, len(body.Checks), 0)

		// Overall status should be unhealthy
		assert.Equal(t, health.StatusUnhealthy, body.Status)
	})

	// Note: Testing with a real database requires an integration test
	// The "after startup complete" case would panic with a nil pool
}

func TestHandler_RegisterRoutes(t *testing.T) {
	// Use handler without startup complete to avoid database access
	handler := newTestHandler(false)
	mux := http.NewServeMux()

	// Should not panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(mux)
	})

	// Test that routes are registered
	server := httptest.NewServer(mux)
	defer server.Close()

	routes := []string{
		"/health",
		"/health/live",
		"/health/ready",
		"/health/startup",
	}

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			resp, err := http.Get(server.URL + route)
			require.NoError(t, err)
			defer resp.Body.Close()

			// All routes should return JSON
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			// Body should be valid JSON
			var body health.Response
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
		})
	}
}

func TestResponse_Structure(t *testing.T) {
	// Test that Response can be marshaled/unmarshaled correctly
	resp := health.Response{
		Status:  health.StatusHealthy,
		Message: "all checks passed",
		Checks: []health.CheckResult{
			{
				Name:    "database",
				Status:  health.StatusHealthy,
				Message: "connected",
			},
		},
		Details: map[string]interface{}{
			"version": "1.0.0",
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var unmarshaled health.Response
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, resp.Status, unmarshaled.Status)
	assert.Equal(t, resp.Message, unmarshaled.Message)
	assert.Equal(t, len(resp.Checks), len(unmarshaled.Checks))
	assert.Equal(t, "1.0.0", unmarshaled.Details["version"])
}
