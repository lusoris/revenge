// Package supervisor provides process supervision and self-healing.
package supervisor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// Supervision strategies.
const (
	// StrategyOneForOne restarts only the failed service.
	StrategyOneForOne = iota
	// StrategyOneForAll restarts all services if one fails.
	StrategyOneForAll
	// StrategyRestForOne restarts failed and all services started after it.
	StrategyRestForOne
)

// Service states.
const (
	StateStarting = iota
	StateRunning
	StateStopping
	StateStopped
	StateFailed
)

var (
	// ErrMaxRestartsExceeded is returned when restart limit is hit.
	ErrMaxRestartsExceeded = errors.New("max restarts exceeded")
	// ErrServiceNotFound is returned when service doesn't exist.
	ErrServiceNotFound = errors.New("service not found")
)

// SupervisorConfig configures the supervisor.
type SupervisorConfig struct {
	// Name for logging.
	Name string

	// Strategy is the restart strategy.
	Strategy int

	// MaxRestarts is max restarts within MaxRestartWindow.
	MaxRestarts int

	// MaxRestartWindow is the time window for max restarts.
	MaxRestartWindow time.Duration

	// RestartDelay is initial delay before restart.
	RestartDelay time.Duration

	// MaxRestartDelay is max delay with exponential backoff.
	MaxRestartDelay time.Duration

	// ShutdownTimeout is max time to wait for graceful shutdown.
	ShutdownTimeout time.Duration
}

// DefaultSupervisorConfig returns sensible defaults.
func DefaultSupervisorConfig(name string) SupervisorConfig {
	return SupervisorConfig{
		Name:             name,
		Strategy:         StrategyOneForOne,
		MaxRestarts:      5,
		MaxRestartWindow: time.Minute,
		RestartDelay:     100 * time.Millisecond,
		MaxRestartDelay:  30 * time.Second,
		ShutdownTimeout:  30 * time.Second,
	}
}

// Service is a supervised service.
type Service interface {
	// Name returns the service name.
	Name() string
	// Start starts the service (blocking until stopped).
	Start(ctx context.Context) error
	// Stop gracefully stops the service.
	Stop(ctx context.Context) error
}

// ServiceFunc wraps a function as a Service.
type ServiceFunc struct {
	name    string
	startFn func(ctx context.Context) error
	stopFn  func(ctx context.Context) error
}

// NewServiceFunc creates a service from functions.
func NewServiceFunc(name string, start func(ctx context.Context) error, stop func(ctx context.Context) error) *ServiceFunc {
	return &ServiceFunc{name: name, startFn: start, stopFn: stop}
}

// Name returns the service name.
func (s *ServiceFunc) Name() string { return s.name }

// Start starts the service.
func (s *ServiceFunc) Start(ctx context.Context) error { return s.startFn(ctx) }

// Stop stops the service.
func (s *ServiceFunc) Stop(ctx context.Context) error {
	if s.stopFn != nil {
		return s.stopFn(ctx)
	}
	return nil
}

// supervisedService wraps a service with supervision state.
type supervisedService struct {
	service  Service
	state    atomic.Int32
	restarts []time.Time
	cancel   context.CancelFunc
	done     chan struct{}
	lastErr  error
	mu       sync.Mutex
}

// Supervisor manages and restarts services.
type Supervisor struct {
	config   SupervisorConfig
	logger   *slog.Logger
	services []*supervisedService
	mu       sync.RWMutex
	running  atomic.Bool
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewSupervisor creates a new supervisor.
func NewSupervisor(config SupervisorConfig, logger *slog.Logger) *Supervisor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Supervisor{
		config: config,
		logger: logger.With(slog.String("supervisor", config.Name)),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Add registers a service for supervision.
func (s *Supervisor) Add(svc Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.services = append(s.services, &supervisedService{
		service:  svc,
		restarts: make([]time.Time, 0),
		done:     make(chan struct{}),
	})
}

// Start starts all services and supervises them.
func (s *Supervisor) Start() error {
	s.mu.Lock()
	if s.running.Load() {
		s.mu.Unlock()
		return errors.New("supervisor already running")
	}
	s.running.Store(true)
	services := s.services
	s.mu.Unlock()

	for _, ss := range services {
		if err := s.startService(ss); err != nil {
			s.logger.Error("failed to start service", "service", ss.service.Name(), "error", err)
			// Stop already started services
			s.Stop()
			return err
		}
	}

	return nil
}

// startService starts a single service with supervision.
func (s *Supervisor) startService(ss *supervisedService) error {
	ss.mu.Lock()
	ctx, cancel := context.WithCancel(s.ctx)
	ss.cancel = cancel
	ss.state.Store(StateStarting)
	ss.done = make(chan struct{})
	ss.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runService(ss, ctx)
	}()

	// Wait briefly for startup
	time.Sleep(10 * time.Millisecond)
	return nil
}

// runService runs a service and handles restarts.
func (s *Supervisor) runService(ss *supervisedService, ctx context.Context) {
	restartDelay := s.config.RestartDelay

	for {
		ss.state.Store(StateRunning)
		s.logger.Info("starting service", "service", ss.service.Name())

		err := ss.service.Start(ctx)

		// Check if shutdown was requested
		if ctx.Err() != nil {
			ss.state.Store(StateStopped)
			close(ss.done)
			return
		}

		// Service exited unexpectedly
		ss.mu.Lock()
		ss.lastErr = err
		ss.mu.Unlock()

		if err != nil {
			s.logger.Error("service failed", "service", ss.service.Name(), "error", err)
		} else {
			s.logger.Warn("service exited unexpectedly", "service", ss.service.Name())
		}

		// Check restart limit
		if !s.canRestart(ss) {
			ss.state.Store(StateFailed)
			s.logger.Error("max restarts exceeded", "service", ss.service.Name())
			s.handleMaxRestartsExceeded(ss)
			return
		}

		// Wait before restart with exponential backoff
		ss.state.Store(StateStopping)
		s.logger.Info("restarting service",
			"service", ss.service.Name(),
			"delay", restartDelay)

		select {
		case <-ctx.Done():
			ss.state.Store(StateStopped)
			close(ss.done)
			return
		case <-time.After(restartDelay):
		}

		// Exponential backoff
		restartDelay *= 2
		if restartDelay > s.config.MaxRestartDelay {
			restartDelay = s.config.MaxRestartDelay
		}

		// Record restart
		ss.mu.Lock()
		ss.restarts = append(ss.restarts, time.Now())
		ss.mu.Unlock()
	}
}

// canRestart checks if service can be restarted.
func (s *Supervisor) canRestart(ss *supervisedService) bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Clean old restarts
	windowStart := time.Now().Add(-s.config.MaxRestartWindow)
	valid := ss.restarts[:0]
	for _, t := range ss.restarts {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	ss.restarts = valid

	return len(ss.restarts) < s.config.MaxRestarts
}

// handleMaxRestartsExceeded handles when a service exceeds restart limit.
func (s *Supervisor) handleMaxRestartsExceeded(ss *supervisedService) {
	switch s.config.Strategy {
	case StrategyOneForAll:
		s.logger.Warn("stopping all services due to max restarts exceeded",
			"failed_service", ss.service.Name())
		go s.Stop()

	case StrategyRestForOne:
		s.logger.Warn("stopping rest of services due to max restarts exceeded",
			"failed_service", ss.service.Name())
		s.mu.RLock()
		defer s.mu.RUnlock()

		found := false
		for _, other := range s.services {
			if other == ss {
				found = true
				continue
			}
			if found {
				s.stopService(other)
			}
		}

	default: // StrategyOneForOne
		// Just let this one fail
	}
}

// Stop gracefully stops all services.
func (s *Supervisor) Stop() error {
	if !s.running.Load() {
		return nil
	}

	s.logger.Info("stopping supervisor")
	s.cancel()

	s.mu.RLock()
	services := s.services
	s.mu.RUnlock()

	// Stop in reverse order
	for i := len(services) - 1; i >= 0; i-- {
		s.stopService(services[i])
	}

	// Wait for all to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("supervisor stopped")
	case <-time.After(s.config.ShutdownTimeout):
		s.logger.Warn("supervisor shutdown timed out")
	}

	s.running.Store(false)
	return nil
}

// stopService stops a single service.
func (s *Supervisor) stopService(ss *supervisedService) {
	ss.mu.Lock()
	if ss.cancel != nil {
		ss.cancel()
	}
	done := ss.done
	ss.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := ss.service.Stop(ctx); err != nil {
		s.logger.Error("error stopping service", "service", ss.service.Name(), "error", err)
	}

	// Wait for service goroutine to exit
	select {
	case <-done:
	case <-ctx.Done():
		s.logger.Warn("service stop timed out", "service", ss.service.Name())
	}
}

// RestartService manually restarts a service.
func (s *Supervisor) RestartService(name string) error {
	s.mu.RLock()
	var target *supervisedService
	for _, ss := range s.services {
		if ss.service.Name() == name {
			target = ss
			break
		}
	}
	s.mu.RUnlock()

	if target == nil {
		return ErrServiceNotFound
	}

	s.logger.Info("manual restart requested", "service", name)
	s.stopService(target)

	// Reset restart counter for manual restart
	target.mu.Lock()
	target.restarts = nil
	target.mu.Unlock()

	return s.startService(target)
}

// Status returns status of all services.
func (s *Supervisor) Status() []ServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statuses := make([]ServiceStatus, 0, len(s.services))
	for _, ss := range s.services {
		ss.mu.Lock()
		status := ServiceStatus{
			Name:         ss.service.Name(),
			State:        ss.state.Load(),
			RestartCount: len(ss.restarts),
		}
		if ss.lastErr != nil {
			status.LastError = ss.lastErr.Error()
		}
		ss.mu.Unlock()
		statuses = append(statuses, status)
	}

	return statuses
}

// ServiceStatus contains service status information.
type ServiceStatus struct {
	Name         string `json:"name"`
	State        int32  `json:"state"`
	RestartCount int    `json:"restart_count"`
	LastError    string `json:"last_error,omitempty"`
}

// StateString returns human-readable state.
func (s ServiceStatus) StateString() string {
	switch s.State {
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	case StateFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", s.State)
	}
}

// HealthCheck checks if all services are healthy.
func (s *Supervisor) HealthCheck() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ss := range s.services {
		state := ss.state.Load()
		if state == StateFailed {
			return fmt.Errorf("service %s has failed", ss.service.Name())
		}
		if state != StateRunning {
			return fmt.Errorf("service %s is not running (state: %d)", ss.service.Name(), state)
		}
	}

	return nil
}
