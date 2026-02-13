// Package middleware provides API middleware components.
package middleware

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"log/slog"

	"github.com/ogen-go/ogen/middleware"
	"github.com/redis/rueidis"

	"github.com/lusoris/revenge/internal/infra/observability"
)

// RedisRateLimiterConfig contains Redis-based rate limiting configuration.
type RedisRateLimiterConfig struct {
	// Name is the rate limiter category name used for metrics labels
	Name string
	// Enabled controls whether rate limiting is active
	Enabled bool
	// RequestsPerSecond is the number of requests allowed per second per IP
	RequestsPerSecond float64
	// Burst is the maximum number of requests allowed in a burst window
	Burst int
	// WindowSize is the sliding window size for rate limiting
	WindowSize time.Duration
	// Operations is a list of operation names to apply rate limiting to
	// If empty, rate limiting is applied to all operations
	Operations []string
	// KeyPrefix is the prefix for Redis keys
	KeyPrefix string
}

// DefaultRedisRateLimiterConfig returns sensible defaults for Redis rate limiting.
func DefaultRedisRateLimiterConfig() RedisRateLimiterConfig {
	return RedisRateLimiterConfig{
		Name:              "global",
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             20,
		WindowSize:        time.Second,
		Operations:        nil, // All operations
		KeyPrefix:         "ratelimit:",
	}
}

// AuthRedisRateLimiterConfig returns stricter rate limits for auth endpoints.
func AuthRedisRateLimiterConfig() RedisRateLimiterConfig {
	return RedisRateLimiterConfig{
		Name:              "auth",
		Enabled:           true,
		RequestsPerSecond: 1, // 1 request per second
		Burst:             5, // Allow burst of 5
		WindowSize:        time.Second,
		Operations: []string{
			"Login",
			"VerifyTOTP",
			"ForgotPassword",
			"ResetPassword",
			"VerifyEmail",
			"BeginWebAuthnLogin",
			"FinishWebAuthnLogin",
		},
		KeyPrefix: "ratelimit:auth:",
	}
}

// RedisRateLimiter provides per-IP rate limiting using Redis/Dragonfly.
// Uses a sliding window algorithm for accurate rate limiting across multiple instances.
type RedisRateLimiter struct {
	config   RedisRateLimiterConfig
	client   rueidis.Client
	logger   *slog.Logger
	fallback *RateLimiter // Fallback to in-memory limiter
	mu       sync.RWMutex
	healthy  bool
}

// NewRedisRateLimiter creates a new Redis-based rate limiter.
// If client is nil, uses in-memory fallback.
func NewRedisRateLimiter(config RedisRateLimiterConfig, client rueidis.Client, logger *slog.Logger) *RedisRateLimiter {
	rl := &RedisRateLimiter{
		config:  config,
		client:  client,
		logger:  logger.With("component", "ratelimit-redis"),
		healthy: client != nil,
	}

	// Create fallback in-memory limiter
	fallbackConfig := RateLimitConfig{
		Name:              config.Name,
		Enabled:           config.Enabled,
		RequestsPerSecond: config.RequestsPerSecond,
		Burst:             config.Burst,
		Operations:        config.Operations,
		CleanupInterval:   5 * time.Minute,
		TTL:               10 * time.Minute,
	}
	rl.fallback = NewRateLimiter(fallbackConfig, logger)

	// Start health check goroutine if we have a client
	if client != nil {
		go rl.healthCheck()
	}

	return rl
}

// Stop stops the rate limiter and its fallback.
func (rl *RedisRateLimiter) Stop() {
	if rl.fallback != nil {
		rl.fallback.Stop()
	}
}

// healthCheck periodically checks Redis connectivity.
func (rl *RedisRateLimiter) healthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := rl.client.Do(ctx, rl.client.B().Ping().Build()).Error()
		cancel()

		rl.mu.Lock()
		wasHealthy := rl.healthy
		rl.healthy = err == nil
		if !wasHealthy && rl.healthy {
			rl.logger.Info("Redis connection restored, using distributed rate limiting")
		} else if wasHealthy && !rl.healthy {
			rl.logger.Warn("Redis connection lost, falling back to in-memory rate limiting", slog.Any("error",err))
		}
		rl.mu.Unlock()
	}
}

// isHealthy returns whether Redis is currently available.
func (rl *RedisRateLimiter) isHealthy() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.healthy
}

// shouldLimit checks if the given operation should be rate limited.
func (rl *RedisRateLimiter) shouldLimit(operationName string) bool {
	if !rl.config.Enabled {
		return false
	}

	// If no operations specified, limit all
	if len(rl.config.Operations) == 0 {
		return true
	}

	// Check if operation is in the list
	for _, op := range rl.config.Operations {
		if op == operationName {
			return true
		}
	}

	return false
}

// slidingWindowScript is a Lua script for sliding window rate limiting.
// Uses a sorted set to track request timestamps.
// Wrapped with rueidis.NewLuaScript for EVALSHA optimization (sends SHA1 hash
// instead of full script text on repeated calls).
var slidingWindowScript = rueidis.NewLuaScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local expire = tonumber(ARGV[4])

-- Remove old entries outside the window
redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)

-- Count current entries
local count = redis.call('ZCARD', key)

if count >= limit then
    return 0
end

-- Add new entry with current timestamp as score
redis.call('ZADD', key, now, now .. '-' .. math.random())

-- Set expiry on the key
redis.call('EXPIRE', key, expire)

return 1
`)

// allow checks if a request is allowed using the sliding window algorithm.
func (rl *RedisRateLimiter) allow(ctx context.Context, ip string) (bool, error) {
	if rl.client == nil {
		return false, fmt.Errorf("redis client not available")
	}

	key := rl.config.KeyPrefix + ip
	now := time.Now().UnixMilli()
	windowMs := rl.config.WindowSize.Milliseconds()
	limit := rl.config.Burst
	expireSec := int64(rl.config.WindowSize.Seconds() * 2) // Double window for safety
	if expireSec < 1 {
		expireSec = 1
	}

	// Execute the Lua script via EVALSHA (falls back to EVAL on cache miss)
	result := slidingWindowScript.Exec(ctx, rl.client,
		[]string{key},
		[]string{
			strconv.FormatInt(now, 10),
			strconv.FormatInt(windowMs, 10),
			strconv.Itoa(limit),
			strconv.FormatInt(expireSec, 10),
		},
	)

	if err := result.Error(); err != nil {
		return false, fmt.Errorf("redis eval failed: %w", err)
	}

	allowed, err := result.AsInt64()
	if err != nil {
		return false, fmt.Errorf("failed to parse result: %w", err)
	}

	return allowed == 1, nil
}

// Middleware returns an ogen middleware that applies rate limiting.
func (rl *RedisRateLimiter) Middleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// Check if this operation should be rate limited
		if !rl.shouldLimit(req.OperationName) {
			return next(req)
		}

		// Get client IP
		clientIP := getClientIP(req.Raw)

		// Check if Redis is healthy
		if !rl.isHealthy() || rl.client == nil {
			// Use fallback in-memory limiter
			return rl.fallback.Middleware()(req, next)
		}

		// Use Redis rate limiting
		ctx, cancel := context.WithTimeout(req.Context, 100*time.Millisecond)
		allowed, err := rl.allow(ctx, clientIP)
		cancel()

		if err != nil {
			// On Redis error, log and use fallback
			rl.logger.Warn("Redis rate limit check failed, using fallback",
				slog.String("ip", clientIP),
				slog.Any("error",err),
			)
			return rl.fallback.Middleware()(req, next)
		}

		if !allowed {
			observability.RecordRateLimitHit(rl.config.Name, "blocked")
			rl.logger.Warn("Rate limit exceeded",
				slog.String("ip", clientIP),
				slog.String("operation", req.OperationName),
				slog.String("backend", "redis"),
			)

			return middleware.Response{}, &RateLimitError{
				IP:            clientIP,
				OperationName: req.OperationName,
				RetryAfter:    time.Duration(1/rl.config.RequestsPerSecond) * time.Second,
			}
		}

		observability.RecordRateLimitHit(rl.config.Name, "allowed")
		return next(req)
	}
}

// Stats returns current rate limiter statistics.
func (rl *RedisRateLimiter) Stats() map[string]any {
	rl.mu.RLock()
	healthy := rl.healthy
	rl.mu.RUnlock()

	return map[string]any{
		"backend":             "redis",
		"healthy":             healthy,
		"requests_per_second": rl.config.RequestsPerSecond,
		"burst":               rl.config.Burst,
		"window_size":         rl.config.WindowSize.String(),
		"key_prefix":          rl.config.KeyPrefix,
	}
}
