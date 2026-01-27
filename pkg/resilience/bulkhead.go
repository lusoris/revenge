package resilience

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Bulkhead errors.
var (
	ErrBulkheadFull     = errors.New("bulkhead full")
	ErrBulkheadTimeout  = errors.New("bulkhead wait timeout")
	ErrBulkheadRejected = errors.New("bulkhead rejected")
)

// BulkheadConfig configures the bulkhead.
type BulkheadConfig struct {
	// Name for identification.
	Name string

	// MaxConcurrent is max simultaneous executions.
	MaxConcurrent int

	// MaxWait is how long to wait for a slot.
	MaxWait time.Duration

	// QueueSize is max waiting requests (0 = no queue, reject immediately).
	QueueSize int

	// OnReject callback when rejected.
	OnReject func(name string)
}

// DefaultBulkheadConfig returns sensible defaults.
func DefaultBulkheadConfig(name string) BulkheadConfig {
	return BulkheadConfig{
		Name:          name,
		MaxConcurrent: 10,
		MaxWait:       5 * time.Second,
		QueueSize:     100,
	}
}

// Bulkhead implements the bulkhead pattern to isolate failures.
type Bulkhead struct {
	config BulkheadConfig

	mu      sync.Mutex
	sem     chan struct{}
	waiting int
}

// NewBulkhead creates a new bulkhead.
func NewBulkhead(config BulkheadConfig) *Bulkhead {
	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = 10
	}
	if config.MaxWait == 0 {
		config.MaxWait = 5 * time.Second
	}

	return &Bulkhead{
		config: config,
		sem:    make(chan struct{}, config.MaxConcurrent),
	}
}

// Execute runs fn within bulkhead constraints.
func (b *Bulkhead) Execute(fn func() error) error {
	return b.ExecuteWithContext(context.Background(), func(_ context.Context) error {
		return fn()
	})
}

// ExecuteWithContext runs with context.
func (b *Bulkhead) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	// Check queue capacity
	if b.config.QueueSize > 0 {
		b.mu.Lock()
		if b.waiting >= b.config.QueueSize {
			b.mu.Unlock()
			if b.config.OnReject != nil {
				b.config.OnReject(b.config.Name)
			}
			return ErrBulkheadFull
		}
		b.waiting++
		b.mu.Unlock()
		defer func() {
			b.mu.Lock()
			b.waiting--
			b.mu.Unlock()
		}()
	}

	// Try to acquire slot
	select {
	case b.sem <- struct{}{}:
		// Got slot immediately
	case <-ctx.Done():
		return ctx.Err()
	default:
		// No slot available, wait if configured
		if b.config.MaxWait == 0 {
			if b.config.OnReject != nil {
				b.config.OnReject(b.config.Name)
			}
			return ErrBulkheadFull
		}

		timer := time.NewTimer(b.config.MaxWait)
		defer timer.Stop()

		select {
		case b.sem <- struct{}{}:
			// Got slot
		case <-timer.C:
			if b.config.OnReject != nil {
				b.config.OnReject(b.config.Name)
			}
			return ErrBulkheadTimeout
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Release slot when done
	defer func() { <-b.sem }()

	return fn(ctx)
}

// Stats returns bulkhead statistics.
func (b *Bulkhead) Stats() BulkheadStats {
	b.mu.Lock()
	defer b.mu.Unlock()

	return BulkheadStats{
		Name:          b.config.Name,
		MaxConcurrent: b.config.MaxConcurrent,
		Active:        len(b.sem),
		Waiting:       b.waiting,
		Available:     b.config.MaxConcurrent - len(b.sem),
	}
}

// BulkheadStats contains bulkhead statistics.
type BulkheadStats struct {
	Name          string
	MaxConcurrent int
	Active        int
	Waiting       int
	Available     int
}

// BulkheadRegistry manages multiple bulkheads.
type BulkheadRegistry struct {
	mu        sync.RWMutex
	bulkheads map[string]*Bulkhead
	defaults  BulkheadConfig
}

// NewBulkheadRegistry creates a new registry.
func NewBulkheadRegistry(defaults BulkheadConfig) *BulkheadRegistry {
	return &BulkheadRegistry{
		bulkheads: make(map[string]*Bulkhead),
		defaults:  defaults,
	}
}

// Get returns or creates a bulkhead.
func (r *BulkheadRegistry) Get(name string) *Bulkhead {
	r.mu.RLock()
	bh, ok := r.bulkheads[name]
	r.mu.RUnlock()

	if ok {
		return bh
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check
	if bh, ok := r.bulkheads[name]; ok {
		return bh
	}

	config := r.defaults
	config.Name = name
	bh = NewBulkhead(config)
	r.bulkheads[name] = bh
	return bh
}

// Stats returns all bulkhead stats.
func (r *BulkheadRegistry) Stats() []BulkheadStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make([]BulkheadStats, 0, len(r.bulkheads))
	for _, bh := range r.bulkheads {
		stats = append(stats, bh.Stats())
	}
	return stats
}
