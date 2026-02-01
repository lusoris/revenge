package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/app"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/logging"
)

func TestNewHTTPServer(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:            "localhost",
			Port:            8080,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
	}

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, nil)

	server := app.NewHTTPServer(cfg, logger, healthService)

	require.NotNil(t, server)
	assert.Equal(t, "localhost:8080", server.Addr)
	assert.Equal(t, 5*time.Second, server.ReadTimeout)
	assert.Equal(t, 5*time.Second, server.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.IdleTimeout)
}

func TestHealthEndpoints(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, nil)
	healthService.MarkStartupComplete()

	server := app.NewHTTPServer(cfg, logger, healthService)
	require.NotNil(t, server)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkBody      func(*testing.T, string)
	}{
		{
			name:           "liveness endpoint",
			path:           "/health/live",
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, `"status":"healthy"`)
			},
		},
		// Note: readiness and startup require database pool, tested in integration tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			server.Handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)

			tt.checkBody(t, string(body))
		})
	}
}

// TestReadinessAndStartup are tested in integration tests with real database
// Unit tests only cover liveness which doesn't require external dependencies
