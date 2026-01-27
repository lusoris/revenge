// Package resilience provides fault-tolerance patterns.
package resilience

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Circuit breaker states.
const (
	StateClosed   = iota // Normal operation
	StateOpen            // Failing, reject requests
	StateHalfOpen        // Testing if recovered
)

var (
	// ErrCircuitOpen is returned when circuit is open.
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTooManyRequests is returned when half-open limit exceeded.
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// CircuitBreakerConfig configures the circuit breaker.
type CircuitBreakerConfig struct {
	// Name identifies this breaker in logs/metrics.
	Name string

	// MaxFailures before opening circuit.
	MaxFailures int

	// FailureRatio threshold (0-1). If set, uses ratio instead of count.
	FailureRatio float64

	// MinRequests before ratio is evaluated.
	MinRequests int

	// Timeout in open state before trying half-open.
	Timeout time.Duration

	// MaxHalfOpenRequests allowed in half-open state.
	MaxHalfOpenRequests int

	// OnStateChange callback.
	OnStateChange func(name string, from, to int)

	// IsSuccessful determines if error is a failure.
	// Default: any non-nil error is a failure.
	IsSuccessful func(err error) bool
}

// DefaultCircuitBreakerConfig returns sensible defaults.
func DefaultCircuitBreakerConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:                name,
		MaxFailures:         5,
		FailureRatio:        0,
		MinRequests:         10,
		Timeout:             30 * time.Second,
		MaxHalfOpenRequests: 3,
		OnStateChange:       nil,
		IsSuccessful:        func(err error) bool { return err == nil },
	}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	config CircuitBreakerConfig

	mu              sync.Mutex
	state           int
	failures        int
	successes       int
	requests        int
	halfOpenCount   int
	lastStateChange time.Time
	expiry          time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.MaxFailures == 0 {
		config.MaxFailures = 5
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxHalfOpenRequests == 0 {
		config.MaxHalfOpenRequests = 3
	}
	if config.IsSuccessful == nil {
		config.IsSuccessful = func(err error) bool { return err == nil }
	}

	return &CircuitBreaker{
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Execute runs the given function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	err := fn()
	cb.afterRequest(err)
	return err
}

// ExecuteWithContext runs with context support.
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	err := fn(ctx)
	cb.afterRequest(err)
	return err
}

// beforeRequest checks if request is allowed.
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		if now.After(cb.expiry) {
			cb.toHalfOpen()
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		if cb.halfOpenCount >= cb.config.MaxHalfOpenRequests {
			return ErrTooManyRequests
		}
		cb.halfOpenCount++
		return nil
	}

	return nil
}

// afterRequest records result.
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.requests++
	if cb.config.IsSuccessful(err) {
		cb.onSuccess()
	} else {
		cb.onFailure()
	}
}

func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		cb.successes++
		// Reset failure count on success in closed state
		cb.failures = 0

	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.MaxHalfOpenRequests {
			cb.toClosed()
		}
	}
}

func (cb *CircuitBreaker) onFailure() {
	switch cb.state {
	case StateClosed:
		cb.failures++
		if cb.shouldTrip() {
			cb.toOpen()
		}

	case StateHalfOpen:
		cb.toOpen()
	}
}

func (cb *CircuitBreaker) shouldTrip() bool {
	// Use ratio if configured
	if cb.config.FailureRatio > 0 && cb.requests >= cb.config.MinRequests {
		ratio := float64(cb.failures) / float64(cb.requests)
		return ratio >= cb.config.FailureRatio
	}

	// Otherwise use count
	return cb.failures >= cb.config.MaxFailures
}

func (cb *CircuitBreaker) toOpen() {
	if cb.state == StateOpen {
		return
	}

	from := cb.state
	cb.state = StateOpen
	cb.expiry = time.Now().Add(cb.config.Timeout)
	cb.lastStateChange = time.Now()

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, from, StateOpen)
	}
}

func (cb *CircuitBreaker) toHalfOpen() {
	from := cb.state
	cb.state = StateHalfOpen
	cb.halfOpenCount = 0
	cb.successes = 0
	cb.failures = 0
	cb.lastStateChange = time.Now()

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, from, StateHalfOpen)
	}
}

func (cb *CircuitBreaker) toClosed() {
	from := cb.state
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.requests = 0
	cb.halfOpenCount = 0
	cb.lastStateChange = time.Now()

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(cb.config.Name, from, StateClosed)
	}
}

// State returns current state.
func (cb *CircuitBreaker) State() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if should transition to half-open
	if cb.state == StateOpen && time.Now().After(cb.expiry) {
		cb.toHalfOpen()
	}

	return cb.state
}

// Stats returns circuit breaker statistics.
func (cb *CircuitBreaker) Stats() CircuitBreakerStats {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return CircuitBreakerStats{
		Name:            cb.config.Name,
		State:           cb.state,
		Failures:        cb.failures,
		Successes:       cb.successes,
		Requests:        cb.requests,
		LastStateChange: cb.lastStateChange,
	}
}

// CircuitBreakerStats contains breaker statistics.
type CircuitBreakerStats struct {
	Name            string
	State           int
	Failures        int
	Successes       int
	Requests        int
	LastStateChange time.Time
}

// Reset forces the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.toClosed()
}

// CircuitBreakerRegistry manages multiple circuit breakers.
type CircuitBreakerRegistry struct {
	mu       sync.RWMutex
	breakers map[string]*CircuitBreaker
	defaults CircuitBreakerConfig
}

// NewCircuitBreakerRegistry creates a new registry.
func NewCircuitBreakerRegistry(defaults CircuitBreakerConfig) *CircuitBreakerRegistry {
	return &CircuitBreakerRegistry{
		breakers: make(map[string]*CircuitBreaker),
		defaults: defaults,
	}
}

// Get returns or creates a circuit breaker.
func (r *CircuitBreakerRegistry) Get(name string) *CircuitBreaker {
	r.mu.RLock()
	cb, ok := r.breakers[name]
	r.mu.RUnlock()

	if ok {
		return cb
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check
	if cb, ok := r.breakers[name]; ok {
		return cb
	}

	config := r.defaults
	config.Name = name
	cb = NewCircuitBreaker(config)
	r.breakers[name] = cb
	return cb
}

// Stats returns all circuit breaker stats.
func (r *CircuitBreakerRegistry) Stats() []CircuitBreakerStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make([]CircuitBreakerStats, 0, len(r.breakers))
	for _, cb := range r.breakers {
		stats = append(stats, cb.Stats())
	}
	return stats
}

// Retry implements retry with exponential backoff.
type Retry struct {
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
	Multiplier  float64
	Jitter      float64 // 0-1, adds randomness
}

// DefaultRetry returns sensible defaults.
func DefaultRetry() Retry {
	return Retry{
		MaxAttempts: 3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     10 * time.Second,
		Multiplier:  2.0,
		Jitter:      0.1,
	}
}

// Do executes fn with retry.
func (r Retry) Do(fn func() error) error {
	return r.DoWithContext(context.Background(), func(_ context.Context) error {
		return fn()
	})
}

// DoWithContext executes with context.
func (r Retry) DoWithContext(ctx context.Context, fn func(context.Context) error) error {
	var lastErr error
	wait := r.InitialWait

	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't wait after last attempt
		if attempt == r.MaxAttempts {
			break
		}

		// Calculate wait with jitter
		jitteredWait := wait
		if r.Jitter > 0 {
			// Simple jitter: ±jitter%
			jitteredWait = time.Duration(float64(wait) * (1 + (r.Jitter * (2*randFloat() - 1))))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(jitteredWait):
		}

		// Exponential backoff
		wait = time.Duration(float64(wait) * r.Multiplier)
		if wait > r.MaxWait {
			wait = r.MaxWait
		}
	}

	return lastErr
}

// Simple random float using atomic counter for jitter.
var randCounter atomic.Uint64

func randFloat() float64 {
	v := randCounter.Add(1)
	return float64(v%1000) / 1000.0
}

// WithCircuitBreaker combines retry with circuit breaker.
func (r Retry) WithCircuitBreaker(cb *CircuitBreaker, fn func() error) error {
	return r.Do(func() error {
		return cb.Execute(fn)
	})
}
