package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/config"
	"github.com/stretchr/testify/assert"
)

func dummyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestCORSMiddleware_NoOrigin(t *testing.T) {
	cfg := config.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	handler := middleware.CORSMiddleware(cfg)(dummyHandler())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"), "no CORS headers without Origin")
}

func TestCORSMiddleware_WildcardOrigin(t *testing.T) {
	cfg := config.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	handler := middleware.CORSMiddleware(cfg)(dummyHandler())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Contains(t, rec.Header().Get("Access-Control-Expose-Headers"), "X-Request-ID")
	assert.Equal(t, "43200", rec.Header().Get("Access-Control-Max-Age"))
	assert.Contains(t, rec.Header().Get("Vary"), "Origin")
}

func TestCORSMiddleware_SpecificOrigins(t *testing.T) {
	cfg := config.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:5173", "https://app.example.com"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}

	handler := middleware.CORSMiddleware(cfg)(dummyHandler())

	t.Run("allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://evil.com")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestCORSMiddleware_Preflight(t *testing.T) {
	cfg := config.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	handler := middleware.CORSMiddleware(cfg)(dummyHandler())
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code, "preflight returns 204")
	assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

func TestCORSMiddleware_NoCredentials(t *testing.T) {
	cfg := config.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
	}

	handler := middleware.CORSMiddleware(cfg)(dummyHandler())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Credentials"), "no credentials header when disabled")
	assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
}
