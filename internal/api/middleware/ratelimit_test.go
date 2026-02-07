package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewRateLimiter(t *testing.T) {
	logger := zap.NewNop()
	config := DefaultRateLimitConfig()

	rl := NewRateLimiter(config, logger)
	require.NotNil(t, rl)

	// Stop the cleanup goroutine
	rl.Stop()
}

func TestRateLimiter_DefaultConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, float64(10), config.RequestsPerSecond)
	assert.Equal(t, 20, config.Burst)
	assert.Empty(t, config.Operations)
	assert.Equal(t, 5*time.Minute, config.CleanupInterval)
	assert.Equal(t, 10*time.Minute, config.TTL)
}

func TestRateLimiter_AuthConfig(t *testing.T) {
	config := AuthRateLimitConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, float64(1), config.RequestsPerSecond)
	assert.Equal(t, 5, config.Burst)
	assert.Contains(t, config.Operations, "LoginUser")
	assert.Contains(t, config.Operations, "VerifyMFA")
	assert.Contains(t, config.Operations, "RequestPasswordReset")
	assert.Contains(t, config.Operations, "ResetPassword")
	assert.Contains(t, config.Operations, "VerifyEmail")
}

func TestRateLimiter_ShouldLimit(t *testing.T) {
	tests := []struct {
		name          string
		config        RateLimitConfig
		operationName string
		expected      bool
	}{
		{
			name: "disabled returns false",
			config: RateLimitConfig{
				Enabled: false,
			},
			operationName: "LoginUser",
			expected:      false,
		},
		{
			name: "no operations limits all",
			config: RateLimitConfig{
				Enabled:    true,
				Operations: nil,
			},
			operationName: "AnyOperation",
			expected:      true,
		},
		{
			name: "operation in list is limited",
			config: RateLimitConfig{
				Enabled:    true,
				Operations: []string{"LoginUser", "VerifyMFA"},
			},
			operationName: "LoginUser",
			expected:      true,
		},
		{
			name: "operation not in list is not limited",
			config: RateLimitConfig{
				Enabled:    true,
				Operations: []string{"LoginUser", "VerifyMFA"},
			},
			operationName: "GetUser",
			expected:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := zap.NewNop()
			rl := NewRateLimiter(tc.config, logger)
			defer rl.Stop()

			result := rl.shouldLimit(tc.operationName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRateLimiter_GetLimiter(t *testing.T) {
	logger := zap.NewNop()
	config := DefaultRateLimitConfig()
	rl := NewRateLimiter(config, logger)
	defer rl.Stop()

	// Get limiter for first IP
	limiter1 := rl.getLimiter("192.168.1.1")
	require.NotNil(t, limiter1)

	// Get limiter again for same IP - should be the same
	limiter1Again := rl.getLimiter("192.168.1.1")
	assert.Same(t, limiter1, limiter1Again)

	// Get limiter for different IP - should be different
	limiter2 := rl.getLimiter("192.168.1.2")
	assert.NotSame(t, limiter1, limiter2)
}

func TestGetClientIP(t *testing.T) {
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tc.remoteAddr
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			result := getClientIP(req)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRateLimiter_Middleware(t *testing.T) {
	logger := zap.NewNop()

	t.Run("allows requests within limit", func(t *testing.T) {
		config := RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: 10,
			Burst:             2,
			Operations:        []string{"TestOperation"},
			CleanupInterval:   time.Hour,
			TTL:               time.Hour,
		}
		rl := NewRateLimiter(config, logger)
		defer rl.Stop()

		mw := rl.Middleware()

		// Create a mock next function that returns success
		nextCalled := 0
		next := func(req middleware.Request) (middleware.Response, error) {
			nextCalled++
			return middleware.Response{Type: "success"}, nil
		}

		// First request should succeed
		req := createTestRequest("TestOperation", "192.168.1.1:1234")
		resp, err := mw(req, next)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Type)
		assert.Equal(t, 1, nextCalled)

		// Second request should succeed (within burst)
		resp, err = mw(req, next)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Type)
		assert.Equal(t, 2, nextCalled)
	})

	t.Run("blocks requests exceeding limit", func(t *testing.T) {
		config := RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: 10,
			Burst:             1, // Only 1 request allowed
			Operations:        []string{"TestOperation"},
			CleanupInterval:   time.Hour,
			TTL:               time.Hour,
		}
		rl := NewRateLimiter(config, logger)
		defer rl.Stop()

		mw := rl.Middleware()

		next := func(req middleware.Request) (middleware.Response, error) {
			return middleware.Response{Type: "success"}, nil
		}

		req := createTestRequest("TestOperation", "192.168.1.1:1234")

		// First request succeeds
		resp, err := mw(req, next)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Type)

		// Second request immediately after should fail
		resp, err = mw(req, next)
		require.Error(t, err)
		assert.Nil(t, resp.Type) // Response.Type should be nil on error

		var rateLimitErr *RateLimitError
		require.ErrorAs(t, err, &rateLimitErr)
		assert.Equal(t, "192.168.1.1", rateLimitErr.IP)
		assert.Equal(t, "TestOperation", rateLimitErr.OperationName)
		assert.Equal(t, http.StatusTooManyRequests, rateLimitErr.StatusCode())
	})

	t.Run("skips non-matching operations", func(t *testing.T) {
		config := RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: 10,
			Burst:             1,
			Operations:        []string{"LoginUser"},
			CleanupInterval:   time.Hour,
			TTL:               time.Hour,
		}
		rl := NewRateLimiter(config, logger)
		defer rl.Stop()

		mw := rl.Middleware()

		nextCalled := 0
		next := func(req middleware.Request) (middleware.Response, error) {
			nextCalled++
			return middleware.Response{Type: "success"}, nil
		}

		// GetUser should not be limited even with many requests
		req := createTestRequest("GetUser", "192.168.1.1:1234")
		for i := 0; i < 10; i++ {
			resp, err := mw(req, next)
			require.NoError(t, err)
			assert.Equal(t, "success", resp.Type)
		}
		assert.Equal(t, 10, nextCalled)
	})

	t.Run("different IPs have separate limits", func(t *testing.T) {
		config := RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: 10,
			Burst:             1,
			Operations:        []string{"TestOperation"},
			CleanupInterval:   time.Hour,
			TTL:               time.Hour,
		}
		rl := NewRateLimiter(config, logger)
		defer rl.Stop()

		mw := rl.Middleware()

		next := func(req middleware.Request) (middleware.Response, error) {
			return middleware.Response{Type: "success"}, nil
		}

		// First IP uses its quota
		req1 := createTestRequest("TestOperation", "192.168.1.1:1234")
		_, err := mw(req1, next)
		require.NoError(t, err)

		// First IP is now limited
		_, err = mw(req1, next)
		require.Error(t, err)

		// Second IP should still work
		req2 := createTestRequest("TestOperation", "192.168.1.2:1234")
		resp, err := mw(req2, next)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Type)
	})
}

func TestRateLimitError(t *testing.T) {
	err := &RateLimitError{
		IP:            "192.168.1.1",
		OperationName: "LoginUser",
		RetryAfter:    time.Second,
	}

	assert.Equal(t, "rate limit exceeded", err.Error())
	assert.Equal(t, http.StatusTooManyRequests, err.StatusCode())

	headers := err.ResponseHeaders()
	assert.Equal(t, "1s", headers.Get("Retry-After"))
}

func TestRateLimiter_Cleanup(t *testing.T) {
	logger := zap.NewNop()
	config := RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             10,
		CleanupInterval:   50 * time.Millisecond,
		TTL:               100 * time.Millisecond,
	}
	rl := NewRateLimiter(config, logger)
	defer rl.Stop()

	// Create some limiters
	rl.getLimiter("192.168.1.1")
	rl.getLimiter("192.168.1.2")

	// Verify they exist
	_, ok1 := rl.limiters.Get("192.168.1.1")
	_, ok2 := rl.limiters.Get("192.168.1.2")
	assert.True(t, ok1)
	assert.True(t, ok2)

	// Wait for TTL-based eviction
	time.Sleep(200 * time.Millisecond)

	// Limiters should be evicted by TTL
	_, ok1 = rl.limiters.Get("192.168.1.1")
	_, ok2 = rl.limiters.Get("192.168.1.2")
	assert.False(t, ok1)
	assert.False(t, ok2)
}

// createTestRequest creates a middleware.Request for testing.
func createTestRequest(operationName, remoteAddr string) middleware.Request {
	httpReq := httptest.NewRequest(http.MethodPost, "/test", nil)
	httpReq.RemoteAddr = remoteAddr

	return middleware.Request{
		Context:       context.Background(),
		OperationName: operationName,
		Raw:           httpReq,
	}
}
