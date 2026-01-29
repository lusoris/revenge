// Package health provides service health checking.
package health

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Status represents overall health status.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

// Category represents service criticality.
type Category string

const (
	CategoryCritical Category = "critical" // Must be healthy
	CategoryWarm     Category = "warm"     // Should be healthy
	CategoryCold     Category = "cold"     // Can be unhealthy
)

// Check defines a health check.
type Check struct {
	Name     string
	Category Category
	Check    func(ctx context.Context) error
	Timeout  time.Duration
}

// ServiceStatus represents a single service's health.
type ServiceStatus struct {
	Name        string        `json:"name"`
	Healthy     bool          `json:"healthy"`
	Initialized bool          `json:"initialized"`
	Category    Category      `json:"category"`
	Latency     time.Duration `json:"latency_ms"`
	Error       string        `json:"error,omitempty"`
}

// Report represents overall system health.
type Report struct {
	Status    Status                   `json:"status"`
	Services  map[string]ServiceStatus `json:"services"`
	CheckedAt time.Time                `json:"checked_at"`
}

// Checker manages health checks.
type Checker struct {
	mu     sync.RWMutex
	checks map[string]Check
	logger *slog.Logger

	// Cached status
	lastStatus    Report
	lastCheckTime time.Time
	cacheDuration time.Duration
}

// NewChecker creates a new health checker.
func NewChecker(logger *slog.Logger) *Checker {
	return &Checker{
		checks:        make(map[string]Check),
		logger:        logger.With(slog.String("component", "health")),
		cacheDuration: 5 * time.Second,
	}
}

// Register adds a health check.
func (c *Checker) Register(check Check) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if check.Timeout == 0 {
		check.Timeout = 5 * time.Second
	}

	c.checks[check.Name] = check
}

// RegisterFunc is a convenience method to register a simple check.
func (c *Checker) RegisterFunc(name string, category Category, check func(ctx context.Context) error) {
	c.Register(Check{
		Name:     name,
		Category: category,
		Check:    check,
		Timeout:  5 * time.Second,
	})
}

// Check performs all health checks.
func (c *Checker) Check(ctx context.Context) Report {
	c.mu.RLock()
	// Return cached if recent
	if time.Since(c.lastCheckTime) < c.cacheDuration {
		status := c.lastStatus
		c.mu.RUnlock()
		return status
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if time.Since(c.lastCheckTime) < c.cacheDuration {
		return c.lastStatus
	}

	status := Report{
		Status:    StatusHealthy,
		Services:  make(map[string]ServiceStatus),
		CheckedAt: time.Now(),
	}

	var wg sync.WaitGroup
	var statusMu sync.Mutex

	for name, check := range c.checks {
		wg.Add(1)
		go func(name string, check Check) {
			defer wg.Done()

			checkCtx, cancel := context.WithTimeout(ctx, check.Timeout)
			defer cancel()

			start := time.Now()
			err := check.Check(checkCtx)
			latency := time.Since(start)

			svcStatus := ServiceStatus{
				Name:     name,
				Healthy:  err == nil,
				Category: check.Category,
				Latency:  latency,
			}

			if err != nil {
				svcStatus.Error = err.Error()
			}

			statusMu.Lock()
			status.Services[name] = svcStatus
			statusMu.Unlock()
		}(name, check)
	}

	wg.Wait()

	// Determine overall status
	for _, svc := range status.Services {
		if !svc.Healthy {
			switch svc.Category {
			case CategoryCritical:
				status.Status = StatusUnhealthy
			case CategoryWarm:
				if status.Status != StatusUnhealthy {
					status.Status = StatusDegraded
				}
			}
		}
	}

	c.lastStatus = status
	c.lastCheckTime = time.Now()

	return status
}

// CheckCategory performs health checks for a specific category.
func (c *Checker) CheckCategory(ctx context.Context, category Category) Report {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := Report{
		Status:    StatusHealthy,
		Services:  make(map[string]ServiceStatus),
		CheckedAt: time.Now(),
	}

	for name, check := range c.checks {
		if check.Category != category {
			continue
		}

		checkCtx, cancel := context.WithTimeout(ctx, check.Timeout)
		start := time.Now()
		err := check.Check(checkCtx)
		latency := time.Since(start)
		cancel()

		svcStatus := ServiceStatus{
			Name:     name,
			Healthy:  err == nil,
			Category: check.Category,
			Latency:  latency,
		}

		if err != nil {
			svcStatus.Error = err.Error()
			if category == CategoryCritical {
				status.Status = StatusUnhealthy
			} else {
				status.Status = StatusDegraded
			}
		}

		status.Services[name] = svcStatus
	}

	return status
}

// IsHealthy returns true if all critical services are healthy.
func (c *Checker) IsHealthy(ctx context.Context) bool {
	status := c.Check(ctx)
	return status.Status != StatusUnhealthy
}

// IsReady returns true if all services are healthy (liveness).
func (c *Checker) IsReady(ctx context.Context) bool {
	status := c.Check(ctx)
	return status.Status == StatusHealthy
}
