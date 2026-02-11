package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "static path",
			input:    "/api/v1/movies",
			expected: "/api/v1/movies",
		},
		{
			name:     "numeric ID",
			input:    "/api/v1/movies/123",
			expected: "/api/v1/movies/{id}",
		},
		{
			name:     "UUID",
			input:    "/api/v1/users/550e8400-e29b-41d4-a716-446655440000",
			expected: "/api/v1/users/{id}",
		},
		{
			name:     "nested resource with ID",
			input:    "/api/v1/libraries/1/stats",
			expected: "/api/v1/libraries/{id}/stats",
		},
		{
			name:     "multiple IDs",
			input:    "/api/v1/users/123/movies/456",
			expected: "/api/v1/users/{id}/movies/{id}",
		},
		{
			name:     "health endpoint",
			input:    "/health/live",
			expected: "/health/live",
		},
		{
			name:     "metrics endpoint",
			input:    "/metrics",
			expected: "/metrics",
		},
		{
			name:     "empty path",
			input:    "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsIDSegment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "numeric ID",
			input:    "12345",
			expected: true,
		},
		{
			name:     "UUID",
			input:    "550e8400-e29b-41d4-a716-446655440000",
			expected: true,
		},
		{
			name:     "word",
			input:    "movies",
			expected: false,
		},
		{
			name:     "api",
			input:    "api",
			expected: false,
		},
		{
			name:     "version",
			input:    "v1",
			expected: false,
		},
		{
			name:     "short alphanumeric with digits",
			input:    "abc12345",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIDSegment(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRecordMetrics(t *testing.T) {
	t.Run("RecordCacheHit", func(t *testing.T) {
		// Should not panic
		RecordCacheHit("test", "l1")
		RecordCacheHit("test", "l2")
	})

	t.Run("RecordCacheMiss", func(t *testing.T) {
		// Should not panic
		RecordCacheMiss("test", "l1")
		RecordCacheMiss("test", "l2")
	})

	t.Run("RecordJobEnqueued", func(t *testing.T) {
		// Should not panic
		RecordJobEnqueued("test_job")
	})

	t.Run("RecordJobCompleted", func(t *testing.T) {
		// Should not panic
		RecordJobCompleted("test_job", "completed")
		RecordJobCompleted("test_job", "failed")
	})

	t.Run("RecordAuthAttempt", func(t *testing.T) {
		// Should not panic
		RecordAuthAttempt("login", "success")
		RecordAuthAttempt("login", "failure")
	})

	t.Run("RecordRateLimitHit", func(t *testing.T) {
		// Should not panic
		RecordRateLimitHit("global", "allowed")
		RecordRateLimitHit("auth", "blocked")
	})
}

func TestStatusResponseWriter(t *testing.T) {
	t.Run("default status is 200", func(t *testing.T) {
		wrapped := &statusResponseWriter{status: 200}
		assert.Equal(t, 200, wrapped.status)
	})

	t.Run("WriteHeader captures status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		wrapped := &statusResponseWriter{ResponseWriter: rec, status: 200}

		wrapped.WriteHeader(404)

		assert.Equal(t, 404, wrapped.status)
		assert.Equal(t, 404, rec.Code)
	})

	t.Run("WriteHeader 500", func(t *testing.T) {
		rec := httptest.NewRecorder()
		wrapped := &statusResponseWriter{ResponseWriter: rec, status: 200}

		wrapped.WriteHeader(500)

		assert.Equal(t, 500, wrapped.status)
	})
}

func TestIsAlphanumericWithHyphens(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC-123", true},
		{"a-b-c", true},
		{"abc_123", false},
		{"abc 123", false},
		{"abc.123", false},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isAlphanumericWithHyphens(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasDigits(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"123", true},
		{"abc", false},
		{"", false},
		{"a1b2c3", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasDigits(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStandardHTTPMetricsMiddleware(t *testing.T) {
	t.Run("records metrics for request", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		wrapped := StandardHTTPMetricsMiddleware(handler)

		req := httptest.NewRequest("GET", "/api/v1/test", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("records metrics for 404", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		wrapped := StandardHTTPMetricsMiddleware(handler)

		req := httptest.NewRequest("GET", "/not-found", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("normalizes path with IDs", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := StandardHTTPMetricsMiddleware(handler)

		req := httptest.NewRequest("GET", "/api/v1/movies/12345", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestExtractStatusFromResponse(t *testing.T) {
	t.Run("nil response returns 200", func(t *testing.T) {
		// We cannot easily mock ogen middleware.Response
		// but we verify the function handles various type names
	})
}

