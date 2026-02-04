package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expected   string
	}{
		{
			name:       "uses X-Forwarded-For first",
			headers:    map[string]string{"X-Forwarded-For": "10.0.0.1"},
			remoteAddr: "192.168.1.1:1234",
			expected:   "10.0.0.1",
		},
		{
			name:       "uses first IP from X-Forwarded-For chain",
			headers:    map[string]string{"X-Forwarded-For": "10.0.0.1, 10.0.0.2, 10.0.0.3"},
			remoteAddr: "192.168.1.1:1234",
			expected:   "10.0.0.1",
		},
		{
			name:       "trims whitespace from X-Forwarded-For",
			headers:    map[string]string{"X-Forwarded-For": "  10.0.0.1  ,  10.0.0.2  "},
			remoteAddr: "192.168.1.1:1234",
			expected:   "10.0.0.1",
		},
		{
			name:       "uses X-Real-IP if no X-Forwarded-For",
			headers:    map[string]string{"X-Real-IP": "10.0.0.5"},
			remoteAddr: "192.168.1.1:1234",
			expected:   "10.0.0.5",
		},
		{
			name:       "prefers X-Forwarded-For over X-Real-IP",
			headers:    map[string]string{"X-Forwarded-For": "10.0.0.1", "X-Real-IP": "10.0.0.5"},
			remoteAddr: "192.168.1.1:1234",
			expected:   "10.0.0.1",
		},
		{
			name:       "falls back to RemoteAddr without port",
			headers:    map[string]string{},
			remoteAddr: "192.168.1.1:1234",
			expected:   "192.168.1.1",
		},
		{
			name:       "handles RemoteAddr without port",
			headers:    map[string]string{},
			remoteAddr: "192.168.1.1",
			expected:   "192.168.1.1",
		},
		{
			name:       "handles IPv6 with port",
			headers:    map[string]string{},
			remoteAddr: "[::1]:8080",
			expected:   "::1",
		},
		{
			name:       "handles IPv6 without port",
			headers:    map[string]string{},
			remoteAddr: "::1",
			expected:   "::1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tc.remoteAddr
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			got := extractClientIP(req)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestRequestMetadataMiddleware(t *testing.T) {
	t.Run("extracts all metadata", func(t *testing.T) {
		mw := RequestMetadataMiddleware()

		httpReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		httpReq.RemoteAddr = "192.168.1.1:8080"
		httpReq.Header.Set("User-Agent", "TestClient/1.0")
		httpReq.Header.Set("Accept-Language", "en-US,en;q=0.9")

		var capturedMeta RequestMetadata
		next := func(req middleware.Request) (middleware.Response, error) {
			capturedMeta = GetRequestMetadata(req.Context)
			return middleware.Response{}, nil
		}

		req := middleware.Request{
			Context: httpReq.Context(),
			Raw:     httpReq,
		}

		_, err := mw(req, next)
		require.NoError(t, err)

		assert.Equal(t, "192.168.1.1", capturedMeta.IPAddress)
		assert.Equal(t, "TestClient/1.0", capturedMeta.UserAgent)
		assert.Equal(t, "en-US,en;q=0.9", capturedMeta.AcceptLanguage)
	})

	t.Run("handles missing headers gracefully", func(t *testing.T) {
		mw := RequestMetadataMiddleware()

		httpReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		httpReq.RemoteAddr = "10.0.0.1:1234"
		// No User-Agent or Accept-Language headers

		var capturedMeta RequestMetadata
		next := func(req middleware.Request) (middleware.Response, error) {
			capturedMeta = GetRequestMetadata(req.Context)
			return middleware.Response{}, nil
		}

		req := middleware.Request{
			Context: httpReq.Context(),
			Raw:     httpReq,
		}

		_, err := mw(req, next)
		require.NoError(t, err)

		assert.Equal(t, "10.0.0.1", capturedMeta.IPAddress)
		assert.Equal(t, "", capturedMeta.UserAgent)
		assert.Equal(t, "", capturedMeta.AcceptLanguage)
	})

	t.Run("uses X-Forwarded-For for IP", func(t *testing.T) {
		mw := RequestMetadataMiddleware()

		httpReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		httpReq.RemoteAddr = "192.168.1.1:8080"
		httpReq.Header.Set("X-Forwarded-For", "203.0.113.50, 70.41.3.18, 150.172.238.178")

		var capturedMeta RequestMetadata
		next := func(req middleware.Request) (middleware.Response, error) {
			capturedMeta = GetRequestMetadata(req.Context)
			return middleware.Response{}, nil
		}

		req := middleware.Request{
			Context: httpReq.Context(),
			Raw:     httpReq,
		}

		_, err := mw(req, next)
		require.NoError(t, err)

		assert.Equal(t, "203.0.113.50", capturedMeta.IPAddress)
	})
}

func TestContextHelpers(t *testing.T) {
	t.Run("GetIPAddress returns IP", func(t *testing.T) {
		ctx := WithRequestMetadata(httptest.NewRequest(http.MethodGet, "/", nil).Context(), RequestMetadata{
			IPAddress: "1.2.3.4",
		})
		assert.Equal(t, "1.2.3.4", GetIPAddress(ctx))
	})

	t.Run("GetUserAgent returns user agent", func(t *testing.T) {
		ctx := WithRequestMetadata(httptest.NewRequest(http.MethodGet, "/", nil).Context(), RequestMetadata{
			UserAgent: "Mozilla/5.0",
		})
		assert.Equal(t, "Mozilla/5.0", GetUserAgent(ctx))
	})

	t.Run("GetRequestMetadata returns empty struct for missing context", func(t *testing.T) {
		ctx := httptest.NewRequest(http.MethodGet, "/", nil).Context()
		meta := GetRequestMetadata(ctx)
		assert.Equal(t, "", meta.IPAddress)
		assert.Equal(t, "", meta.UserAgent)
		assert.Equal(t, "", meta.AcceptLanguage)
	})
}
