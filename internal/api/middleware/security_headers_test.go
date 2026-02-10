package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	t.Parallel()

	handler := SecurityHeadersMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	tests := []struct {
		header string
		want   string
	}{
		{"Content-Security-Policy", "default-src 'self'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; img-src 'self' image.tmdb.org *.image.tmdb.org data: https://cdn.jsdelivr.net; media-src 'self' blob:; connect-src 'self' https://cdn.jsdelivr.net https://proxy.scalar.com; font-src 'self' https://cdn.jsdelivr.net https://fonts.scalar.com data: blob:; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"},
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{"Permissions-Policy", "camera=(), microphone=(), geolocation=()"},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			assert.Equal(t, tt.want, rec.Header().Get(tt.header))
		})
	}
}

func TestSecurityHeadersMiddleware_PassesThrough(t *testing.T) {
	t.Parallel()

	called := false
	handler := SecurityHeadersMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.Header().Set("X-Custom", "value")
		w.WriteHeader(http.StatusCreated)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/test", nil))

	assert.True(t, called, "next handler should be called")
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "value", rec.Header().Get("X-Custom"))
	// Security headers still present
	assert.NotEmpty(t, rec.Header().Get("Content-Security-Policy"))
}
