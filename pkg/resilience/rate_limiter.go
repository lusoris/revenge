package resilience

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Rate limiter errors.
var (
	ErrRateLimited = errors.New("rate limited")
)

// RateLimiterConfig configures the rate limiter.
type RateLimiterConfig struct {
	// Name for identification.
	Name string

	// Rate is requests per second.
	Rate float64

	// Burst is max tokens available.
	Burst int

	// WaitTimeout is max time to wait for a token.
	WaitTimeout time.Duration
}

// DefaultRateLimiterConfig returns sensible defaults.
func DefaultRateLimiterConfig(name string) RateLimiterConfig {
	return RateLimiterConfig{
		Name:        name,
		Rate:        100, // 100 req/s
		Burst:       200,
		WaitTimeout: time.Second,
	}
}

// TokenBucketLimiter implements token bucket rate limiting.
type TokenBucketLimiter struct {
	config RateLimiterConfig

	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
	refillRate float64 // tokens per nanosecond
	maxTokens  float64
}

// NewTokenBucketLimiter creates a new token bucket limiter.
func NewTokenBucketLimiter(config RateLimiterConfig) *TokenBucketLimiter {
	if config.Rate <= 0 {
		config.Rate = 100
	}
	if config.Burst <= 0 {
		config.Burst = int(config.Rate)
	}

	return &TokenBucketLimiter{
		config:     config,
		tokens:     float64(config.Burst),
		lastRefill: time.Now(),
		refillRate: config.Rate / float64(time.Second),
		maxTokens:  float64(config.Burst),
	}
}

// Allow checks if request is allowed (non-blocking).
func (l *TokenBucketLimiter) Allow() bool {
	return l.AllowN(1)
}

// AllowN checks if n requests are allowed.
func (l *TokenBucketLimiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.tokens >= float64(n) {
		l.tokens -= float64(n)
		return true
	}
	return false
}

// Wait blocks until a token is available or context is cancelled.
func (l *TokenBucketLimiter) Wait(ctx context.Context) error {
	return l.WaitN(ctx, 1)
}

// WaitN waits for n tokens.
func (l *TokenBucketLimiter) WaitN(ctx context.Context, n int) error {
	// Try immediately
	if l.AllowN(n) {
		return nil
	}

	// Calculate wait time
	l.mu.Lock()
	tokensNeeded := float64(n) - l.tokens
	waitDuration := time.Duration(tokensNeeded / l.refillRate)
	l.mu.Unlock()

	// Check timeout
	if l.config.WaitTimeout > 0 && waitDuration > l.config.WaitTimeout {
		return ErrRateLimited
	}

	// Wait
	timer := time.NewTimer(waitDuration)
	defer timer.Stop()

	select {
	case <-timer.C:
		// Try again after wait
		if l.AllowN(n) {
			return nil
		}
		return ErrRateLimited
	case <-ctx.Done():
		return ctx.Err()
	}
}

// refill adds tokens based on time passed (must hold lock).
func (l *TokenBucketLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill)
	l.lastRefill = now

	l.tokens += float64(elapsed) * l.refillRate
	if l.tokens > l.maxTokens {
		l.tokens = l.maxTokens
	}
}

// Tokens returns available tokens.
func (l *TokenBucketLimiter) Tokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	return l.tokens
}

// Stats returns rate limiter statistics.
func (l *TokenBucketLimiter) Stats() RateLimiterStats {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()

	return RateLimiterStats{
		Name:      l.config.Name,
		Rate:      l.config.Rate,
		Burst:     l.config.Burst,
		Available: l.tokens,
		MaxTokens: l.maxTokens,
	}
}

// RateLimiterStats contains rate limiter statistics.
type RateLimiterStats struct {
	Name      string
	Rate      float64
	Burst     int
	Available float64
	MaxTokens float64
}

// SlidingWindowLimiter implements sliding window rate limiting.
type SlidingWindowLimiter struct {
	config RateLimiterConfig

	mu       sync.Mutex
	requests []time.Time
	window   time.Duration
}

// NewSlidingWindowLimiter creates a new sliding window limiter.
func NewSlidingWindowLimiter(config RateLimiterConfig) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		config:   config,
		requests: make([]time.Time, 0, config.Burst),
		window:   time.Second,
	}
}

// Allow checks if request is allowed.
func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.window)

	// Remove old requests
	valid := l.requests[:0]
	for _, t := range l.requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	l.requests = valid

	// Check limit
	if float64(len(l.requests)) >= l.config.Rate {
		return false
	}

	l.requests = append(l.requests, now)
	return true
}

// RateLimiterRegistry manages multiple rate limiters.
type RateLimiterRegistry struct {
	mu       sync.RWMutex
	limiters map[string]*TokenBucketLimiter
	defaults RateLimiterConfig
}

// NewRateLimiterRegistry creates a new registry.
func NewRateLimiterRegistry(defaults RateLimiterConfig) *RateLimiterRegistry {
	return &RateLimiterRegistry{
		limiters: make(map[string]*TokenBucketLimiter),
		defaults: defaults,
	}
}

// Get returns or creates a rate limiter.
func (r *RateLimiterRegistry) Get(name string) *TokenBucketLimiter {
	r.mu.RLock()
	rl, ok := r.limiters[name]
	r.mu.RUnlock()

	if ok {
		return rl
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check
	if rl, ok := r.limiters[name]; ok {
		return rl
	}

	config := r.defaults
	config.Name = name
	rl = NewTokenBucketLimiter(config)
	r.limiters[name] = rl
	return rl
}

// Stats returns all rate limiter stats.
func (r *RateLimiterRegistry) Stats() []RateLimiterStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make([]RateLimiterStats, 0, len(r.limiters))
	for _, rl := range r.limiters {
		stats = append(stats, rl.Stats())
	}
	return stats
}

// PerKeyLimiter provides rate limiting per key (e.g., per user, per IP).
type PerKeyLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*TokenBucketLimiter
	config   RateLimiterConfig
	cleanup  time.Duration
}

// NewPerKeyLimiter creates a per-key rate limiter.
func NewPerKeyLimiter(config RateLimiterConfig, cleanupInterval time.Duration) *PerKeyLimiter {
	pkl := &PerKeyLimiter{
		limiters: make(map[string]*TokenBucketLimiter),
		config:   config,
		cleanup:  cleanupInterval,
	}

	if cleanupInterval > 0 {
		go pkl.cleanupLoop()
	}

	return pkl
}

// Allow checks if key is allowed.
func (pkl *PerKeyLimiter) Allow(key string) bool {
	return pkl.getLimiter(key).Allow()
}

// Wait waits for key to be allowed.
func (pkl *PerKeyLimiter) Wait(ctx context.Context, key string) error {
	return pkl.getLimiter(key).Wait(ctx)
}

func (pkl *PerKeyLimiter) getLimiter(key string) *TokenBucketLimiter {
	pkl.mu.RLock()
	rl, ok := pkl.limiters[key]
	pkl.mu.RUnlock()

	if ok {
		return rl
	}

	pkl.mu.Lock()
	defer pkl.mu.Unlock()

	if rl, ok := pkl.limiters[key]; ok {
		return rl
	}

	config := pkl.config
	config.Name = key
	rl = NewTokenBucketLimiter(config)
	pkl.limiters[key] = rl
	return rl
}

func (pkl *PerKeyLimiter) cleanupLoop() {
	ticker := time.NewTicker(pkl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		pkl.mu.Lock()
		for key, rl := range pkl.limiters {
			// Remove if at full capacity (unused)
			if rl.Tokens() >= rl.maxTokens {
				delete(pkl.limiters, key)
			}
		}
		pkl.mu.Unlock()
	}
}
