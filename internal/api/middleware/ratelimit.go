// Package middleware provides API middleware components.
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/ogen-go/ogen/middleware"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimitConfig contains rate limiting configuration.
type RateLimitConfig struct {
	// Enabled controls whether rate limiting is active
	Enabled bool
	// RequestsPerSecond is the number of requests allowed per second per IP
	RequestsPerSecond float64
	// Burst is the maximum number of requests allowed in a burst
	Burst int
	// Operations is a list of operation names to apply rate limiting to
	// If empty, rate limiting is applied to all operations
	Operations []string
	// CleanupInterval is how often to clean up stale limiters
	CleanupInterval time.Duration
	// TTL is how long to keep a limiter after last use
	TTL time.Duration
}

// DefaultRateLimitConfig returns sensible defaults for rate limiting.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		Burst:             20,
		Operations:        nil, // All operations
		CleanupInterval:   5 * time.Minute,
		TTL:               10 * time.Minute,
	}
}

// AuthRateLimitConfig returns stricter rate limits for auth endpoints.
func AuthRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 1, // 1 request per second
		Burst:             5, // Allow burst of 5
		Operations: []string{
			"LoginUser",
			"VerifyMFA",
			"RequestPasswordReset",
			"ResetPassword",
			"VerifyEmail",
		},
		CleanupInterval: 5 * time.Minute,
		TTL:             10 * time.Minute,
	}
}

// ipLimiter tracks rate limiters per IP address.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter provides per-IP rate limiting for API endpoints.
type RateLimiter struct {
	config   RateLimitConfig
	limiters sync.Map
	logger   *zap.Logger
	stopCh   chan struct{}
	stopOnce sync.Once
}

// NewRateLimiter creates a new rate limiter with the given configuration.
func NewRateLimiter(config RateLimitConfig, logger *zap.Logger) *RateLimiter {
	// Apply defaults for zero values
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 5 * time.Minute
	}
	if config.TTL == 0 {
		config.TTL = 10 * time.Minute
	}

	rl := &RateLimiter{
		config: config,
		logger: logger.Named("ratelimit"),
		stopCh: make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Stop stops the rate limiter's cleanup goroutine.
func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopCh)
	})
}

// cleanup periodically removes stale IP limiters.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanupStale()
		case <-rl.stopCh:
			return
		}
	}
}

// cleanupStale removes limiters that haven't been used recently.
func (rl *RateLimiter) cleanupStale() {
	now := time.Now()
	cleaned := 0

	rl.limiters.Range(func(key, value any) bool {
		limiter := value.(*ipLimiter)
		if now.Sub(limiter.lastSeen) > rl.config.TTL {
			rl.limiters.Delete(key)
			cleaned++
		}
		return true
	})

	if cleaned > 0 {
		rl.logger.Debug("Cleaned up stale rate limiters",
			zap.Int("count", cleaned),
		)
	}
}

// getLimiter retrieves or creates a rate limiter for the given IP.
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	now := time.Now()

	if v, ok := rl.limiters.Load(ip); ok {
		il := v.(*ipLimiter)
		il.lastSeen = now
		return il.limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.Burst)
	il := &ipLimiter{
		limiter:  limiter,
		lastSeen: now,
	}

	// Use LoadOrStore to handle race conditions
	actual, _ := rl.limiters.LoadOrStore(ip, il)
	return actual.(*ipLimiter).limiter
}

// shouldLimit checks if the given operation should be rate limited.
func (rl *RateLimiter) shouldLimit(operationName string) bool {
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

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Return the first IP in the chain
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	// Remove port if present
	ip := r.RemoteAddr
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] == ':' {
			return ip[:i]
		}
	}

	return ip
}

// Middleware returns an ogen middleware that applies rate limiting.
func (rl *RateLimiter) Middleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// Check if this operation should be rate limited
		if !rl.shouldLimit(req.OperationName) {
			return next(req)
		}

		// Get client IP
		clientIP := getClientIP(req.Raw)

		// Get or create limiter for this IP
		limiter := rl.getLimiter(clientIP)

		// Check if request is allowed
		if !limiter.Allow() {
			rl.logger.Warn("Rate limit exceeded",
				zap.String("ip", clientIP),
				zap.String("operation", req.OperationName),
			)

			// Return rate limit error
			return middleware.Response{}, &RateLimitError{
				IP:            clientIP,
				OperationName: req.OperationName,
				RetryAfter:    time.Duration(1/rl.config.RequestsPerSecond) * time.Second,
			}
		}

		return next(req)
	}
}

// RateLimitError is returned when a client exceeds the rate limit.
type RateLimitError struct {
	IP            string
	OperationName string
	RetryAfter    time.Duration
}

// Error implements the error interface.
func (e *RateLimitError) Error() string {
	return "rate limit exceeded"
}

// StatusCode returns the HTTP status code for rate limiting.
func (e *RateLimitError) StatusCode() int {
	return http.StatusTooManyRequests
}

// ResponseHeaders returns headers to include in the error response.
func (e *RateLimitError) ResponseHeaders() http.Header {
	h := http.Header{}
	h.Set("Retry-After", e.RetryAfter.String())
	return h
}
