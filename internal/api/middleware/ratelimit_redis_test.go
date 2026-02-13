package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/ogen-go/ogen/middleware"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewRedisRateLimiter(t *testing.T) {
	logger := logging.NewTestLogger()
	config := DefaultRedisRateLimiterConfig()

	// Test creation without Redis client (fallback mode)
	rl := NewRedisRateLimiter(config, nil, logger)
	defer rl.Stop()

	require.NotNil(t, rl)
	require.NotNil(t, rl.fallback)
	require.Nil(t, rl.client)
	assert.False(t, rl.isHealthy())
}

func TestRedisRateLimiter_DefaultConfig(t *testing.T) {
	config := DefaultRedisRateLimiterConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, float64(10), config.RequestsPerSecond)
	assert.Equal(t, 20, config.Burst)
	assert.Equal(t, time.Second, config.WindowSize)
	assert.Equal(t, "ratelimit:", config.KeyPrefix)
	assert.Empty(t, config.Operations)
}

func TestRedisRateLimiter_AuthConfig(t *testing.T) {
	config := AuthRedisRateLimiterConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, float64(1), config.RequestsPerSecond)
	assert.Equal(t, 5, config.Burst)
	assert.Equal(t, "ratelimit:auth:", config.KeyPrefix)
	assert.Contains(t, config.Operations, "Login")
	assert.Contains(t, config.Operations, "VerifyTOTP")
	assert.Contains(t, config.Operations, "ForgotPassword")
	assert.Contains(t, config.Operations, "ResetPassword")
	assert.Contains(t, config.Operations, "VerifyEmail")
	assert.Contains(t, config.Operations, "BeginWebAuthnLogin")
	assert.Contains(t, config.Operations, "FinishWebAuthnLogin")
}

func TestRedisRateLimiter_ShouldLimit(t *testing.T) {
	logger := logging.NewTestLogger()

	tests := []struct {
		name          string
		config        RedisRateLimiterConfig
		operationName string
		expected      bool
	}{
		{
			name: "disabled returns false",
			config: RedisRateLimiterConfig{
				Enabled:    false,
				Operations: nil,
			},
			operationName: "AnyOperation",
			expected:      false,
		},
		{
			name: "no operations limits all",
			config: RedisRateLimiterConfig{
				Enabled:    true,
				Operations: nil,
			},
			operationName: "AnyOperation",
			expected:      true,
		},
		{
			name: "operation in list is limited",
			config: RedisRateLimiterConfig{
				Enabled:    true,
				Operations: []string{"LoginUser", "RegisterUser"},
			},
			operationName: "LoginUser",
			expected:      true,
		},
		{
			name: "operation not in list is not limited",
			config: RedisRateLimiterConfig{
				Enabled:    true,
				Operations: []string{"LoginUser", "RegisterUser"},
			},
			operationName: "GetMovies",
			expected:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRedisRateLimiter(tc.config, nil, logger)
			defer rl.Stop()

			result := rl.shouldLimit(tc.operationName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRedisRateLimiter_Allow_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logging.NewTestLogger()
	mockClient := mock.NewClient(ctrl)

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             5,
		WindowSize:        time.Second,
		KeyPrefix:         "test:",
	}

	rl := &RedisRateLimiter{
		config:  config,
		client:  mockClient,
		logger:  logger.With("component", "ratelimit-redis"),
		healthy: true,
	}

	t.Run("allows when under limit", func(t *testing.T) {
		// Mock the EVAL command to return 1 (allowed)
		mockClient.EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Return(mock.Result(mock.RedisInt64(1)))

		allowed, err := rl.allow(context.Background(), "192.168.1.1")
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("denies when over limit", func(t *testing.T) {
		// Mock the EVAL command to return 0 (denied)
		mockClient.EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Return(mock.Result(mock.RedisInt64(0)))

		allowed, err := rl.allow(context.Background(), "192.168.1.2")
		require.NoError(t, err)
		assert.False(t, allowed)
	})
}

func TestRedisRateLimiter_Middleware_UseFallback(t *testing.T) {
	logger := logging.NewTestLogger()

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 100,
		Burst:             100,
		WindowSize:        time.Second,
		KeyPrefix:         "test:",
	}

	// Create limiter without Redis (will use fallback)
	rl := NewRedisRateLimiter(config, nil, logger)
	defer rl.Stop()

	// Create middleware
	mw := rl.Middleware()

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	middlewareReq := middleware.Request{
		Context:       context.Background(),
		OperationName: "TestOperation",
		Raw:           req,
	}

	// Next function that returns success
	next := func(req middleware.Request) (middleware.Response, error) {
		return middleware.Response{Type: nil}, nil
	}

	// Should use fallback and allow the request
	_, err := mw(middlewareReq, next)
	assert.NoError(t, err)
}

func TestRedisRateLimiter_Middleware_SkipsNonMatchingOperations(t *testing.T) {
	logger := logging.NewTestLogger()

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 1,
		Burst:             1,
		WindowSize:        time.Second,
		Operations:        []string{"LoginUser"}, // Only limit LoginUser
		KeyPrefix:         "test:",
	}

	rl := NewRedisRateLimiter(config, nil, logger)
	defer rl.Stop()

	mw := rl.Middleware()

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	middlewareReq := middleware.Request{
		Context:       context.Background(),
		OperationName: "GetMovies", // Not in the limited operations
		Raw:           req,
	}

	nextCalled := false
	next := func(req middleware.Request) (middleware.Response, error) {
		nextCalled = true
		return middleware.Response{Type: nil}, nil
	}

	_, err := mw(middlewareReq, next)
	assert.NoError(t, err)
	assert.True(t, nextCalled, "next should be called for non-matching operations")
}

func TestRedisRateLimiter_Middleware_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logging.NewTestLogger()
	mockClient := mock.NewClient(ctrl)

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             5,
		WindowSize:        time.Second,
		KeyPrefix:         "test:",
	}

	rl := &RedisRateLimiter{
		config:  config,
		client:  mockClient,
		logger:  logger.With("component", "ratelimit-redis"),
		healthy: true,
	}

	// Create fallback limiter
	fallbackConfig := RateLimitConfig{
		Enabled:           config.Enabled,
		RequestsPerSecond: config.RequestsPerSecond,
		Burst:             config.Burst,
		CleanupInterval:   5 * time.Minute,
		TTL:               10 * time.Minute,
	}
	rl.fallback = NewRateLimiter(fallbackConfig, logger)
	defer rl.fallback.Stop()

	t.Run("allows requests within limit", func(t *testing.T) {
		mockClient.EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Return(mock.Result(mock.RedisInt64(1)))

		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "10.0.0.1:12345"

		middlewareReq := middleware.Request{
			Context:       context.Background(),
			OperationName: "LoginUser",
			Raw:           req,
		}

		nextCalled := false
		next := func(req middleware.Request) (middleware.Response, error) {
			nextCalled = true
			return middleware.Response{Type: nil}, nil
		}

		mw := rl.Middleware()
		_, err := mw(middlewareReq, next)

		assert.NoError(t, err)
		assert.True(t, nextCalled)
	})

	t.Run("blocks requests exceeding limit", func(t *testing.T) {
		mockClient.EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Return(mock.Result(mock.RedisInt64(0)))

		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "10.0.0.2:12345"

		middlewareReq := middleware.Request{
			Context:       context.Background(),
			OperationName: "LoginUser",
			Raw:           req,
		}

		nextCalled := false
		next := func(req middleware.Request) (middleware.Response, error) {
			nextCalled = true
			return middleware.Response{Type: nil}, nil
		}

		mw := rl.Middleware()
		_, err := mw(middlewareReq, next)

		assert.Error(t, err)
		assert.False(t, nextCalled)

		// Check it's a rate limit error
		rateLimitErr, ok := err.(*RateLimitError)
		require.True(t, ok)
		assert.Equal(t, "10.0.0.2", rateLimitErr.IP)
		assert.Equal(t, "LoginUser", rateLimitErr.OperationName)
	})
}

func TestRedisRateLimiter_Stats(t *testing.T) {
	logger := logging.NewTestLogger()

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             20,
		WindowSize:        time.Second,
		KeyPrefix:         "test:",
	}

	rl := NewRedisRateLimiter(config, nil, logger)
	defer rl.Stop()

	stats := rl.Stats()

	assert.Equal(t, "redis", stats["backend"])
	assert.Equal(t, false, stats["healthy"])
	assert.Equal(t, float64(10), stats["requests_per_second"])
	assert.Equal(t, 20, stats["burst"])
	assert.Equal(t, "1s", stats["window_size"])
	assert.Equal(t, "test:", stats["key_prefix"])
}

func TestRedisRateLimiter_HealthyStateManagement(t *testing.T) {
	logger := logging.NewTestLogger()
	config := DefaultRedisRateLimiterConfig()

	t.Run("starts unhealthy without client", func(t *testing.T) {
		rl := NewRedisRateLimiter(config, nil, logger)
		defer rl.Stop()

		assert.False(t, rl.isHealthy())
	})

	t.Run("starts healthy with client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mock.NewClient(ctrl)

		// Create limiter directly to avoid starting health check goroutine
		rl := &RedisRateLimiter{
			config:  config,
			client:  mockClient,
			logger:  logger.With("component", "ratelimit-redis"),
			healthy: true,
		}

		assert.True(t, rl.isHealthy())
	})
}

func TestRedisRateLimiter_FallbackOnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logging.NewTestLogger()
	mockClient := mock.NewClient(ctrl)

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 100,
		Burst:             100,
		WindowSize:        time.Second,
		KeyPrefix:         "test:",
	}

	rl := &RedisRateLimiter{
		config:  config,
		client:  mockClient,
		logger:  logger.With("component", "ratelimit-redis"),
		healthy: true,
	}

	// Create fallback limiter
	fallbackConfig := RateLimitConfig{
		Enabled:           config.Enabled,
		RequestsPerSecond: config.RequestsPerSecond,
		Burst:             config.Burst,
		CleanupInterval:   5 * time.Minute,
		TTL:               10 * time.Minute,
	}
	rl.fallback = NewRateLimiter(fallbackConfig, logger)
	defer rl.fallback.Stop()

	// Simulate Redis error
	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.ErrorResult(rueidis.Nil))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	middlewareReq := middleware.Request{
		Context:       context.Background(),
		OperationName: "TestOperation",
		Raw:           req,
	}

	nextCalled := false
	next := func(req middleware.Request) (middleware.Response, error) {
		nextCalled = true
		return middleware.Response{Type: nil}, nil
	}

	mw := rl.Middleware()
	_, err := mw(middlewareReq, next)

	// Should fall back to in-memory and allow the request
	assert.NoError(t, err)
	assert.True(t, nextCalled, "should use fallback and allow request on Redis error")
}

func TestRedisRateLimiter_DifferentIPsHaveSeparateLimits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logging.NewTestLogger()
	mockClient := mock.NewClient(ctrl)

	config := RedisRateLimiterConfig{
		Enabled:           true,
		RequestsPerSecond: 1,
		Burst:             1,
		WindowSize:        time.Second,
		KeyPrefix:         "test:",
	}

	rl := &RedisRateLimiter{
		config:  config,
		client:  mockClient,
		logger:  logger.With("component", "ratelimit-redis"),
		healthy: true,
	}

	// Create fallback limiter
	fallbackConfig := RateLimitConfig{
		Enabled:           config.Enabled,
		RequestsPerSecond: config.RequestsPerSecond,
		Burst:             config.Burst,
		CleanupInterval:   5 * time.Minute,
		TTL:               10 * time.Minute,
	}
	rl.fallback = NewRateLimiter(fallbackConfig, logger)
	defer rl.fallback.Stop()

	mw := rl.Middleware()

	// First IP - allowed
	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.RedisInt64(1)))

	req1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"

	_, err := mw(middleware.Request{
		Context:       context.Background(),
		OperationName: "TestOp",
		Raw:           req1,
	}, func(req middleware.Request) (middleware.Response, error) {
		return middleware.Response{Type: nil}, nil
	})
	assert.NoError(t, err)

	// Second IP - also allowed (separate limit)
	mockClient.EXPECT().
		Do(gomock.Any(), gomock.Any()).
		Return(mock.Result(mock.RedisInt64(1)))

	req2 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"

	_, err = mw(middleware.Request{
		Context:       context.Background(),
		OperationName: "TestOp",
		Raw:           req2,
	}, func(req middleware.Request) (middleware.Response, error) {
		return middleware.Response{Type: nil}, nil
	})
	assert.NoError(t, err)
}
