// Package graceful provides graceful shutdown utilities.
package graceful

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownConfig configures graceful shutdown.
type ShutdownConfig struct {
	// Timeout is max time for shutdown.
	Timeout time.Duration

	// DrainTimeout is time to wait for in-flight requests.
	DrainTimeout time.Duration

	// Signals to listen for.
	Signals []os.Signal
}

// DefaultShutdownConfig returns sensible defaults.
func DefaultShutdownConfig() ShutdownConfig {
	return ShutdownConfig{
		Timeout:      30 * time.Second,
		DrainTimeout: 5 * time.Second,
		Signals:      []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

// Shutdowner handles graceful shutdown.
type Shutdowner struct {
	config   ShutdownConfig
	logger   *slog.Logger
	hooks    []ShutdownHook
	mu       sync.Mutex
	done     chan struct{}
	started  bool
}

// ShutdownHook is called during shutdown.
type ShutdownHook struct {
	Name     string
	Priority int // Lower = earlier
	Fn       func(ctx context.Context) error
}

// NewShutdowner creates a new shutdowner.
func NewShutdowner(config ShutdownConfig, logger *slog.Logger) *Shutdowner {
	return &Shutdowner{
		config: config,
		logger: logger.With(slog.String("component", "shutdown")),
		hooks:  make([]ShutdownHook, 0),
		done:   make(chan struct{}),
	}
}

// Register adds a shutdown hook.
func (s *Shutdowner) Register(hook ShutdownHook) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hooks = append(s.hooks, hook)

	// Keep sorted by priority
	for i := len(s.hooks) - 1; i > 0; i-- {
		if s.hooks[i].Priority < s.hooks[i-1].Priority {
			s.hooks[i], s.hooks[i-1] = s.hooks[i-1], s.hooks[i]
		}
	}
}

// RegisterFunc is a convenience wrapper.
func (s *Shutdowner) RegisterFunc(name string, priority int, fn func(ctx context.Context) error) {
	s.Register(ShutdownHook{Name: name, Priority: priority, Fn: fn})
}

// Start begins listening for shutdown signals.
func (s *Shutdowner) Start() <-chan struct{} {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return s.done
	}
	s.started = true
	s.mu.Unlock()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, s.config.Signals...)

	go func() {
		sig := <-sigCh
		s.logger.Info("received shutdown signal", "signal", sig.String())

		// Ignore subsequent signals during shutdown
		signal.Stop(sigCh)

		s.shutdown()
	}()

	return s.done
}

// Trigger initiates shutdown programmatically.
func (s *Shutdowner) Trigger() {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	s.logger.Info("programmatic shutdown triggered")
	s.shutdown()
}

// shutdown runs all hooks.
func (s *Shutdowner) shutdown() {
	defer close(s.done)

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	s.mu.Lock()
	hooks := make([]ShutdownHook, len(s.hooks))
	copy(hooks, s.hooks)
	s.mu.Unlock()

	var errs []error

	for _, hook := range hooks {
		s.logger.Info("running shutdown hook", "name", hook.Name)

		start := time.Now()
		if err := hook.Fn(ctx); err != nil {
			s.logger.Error("shutdown hook failed",
				"name", hook.Name,
				"error", err,
				"duration", time.Since(start))
			errs = append(errs, err)
		} else {
			s.logger.Info("shutdown hook completed",
				"name", hook.Name,
				"duration", time.Since(start))
		}

		if ctx.Err() != nil {
			s.logger.Warn("shutdown timeout reached")
			break
		}
	}

	if len(errs) > 0 {
		s.logger.Error("shutdown completed with errors", "error_count", len(errs))
	} else {
		s.logger.Info("shutdown completed successfully")
	}
}

// Wait blocks until shutdown is complete.
func (s *Shutdowner) Wait() {
	<-s.done
}

// Done returns the done channel.
func (s *Shutdowner) Done() <-chan struct{} {
	return s.done
}

// DrainableServer wraps an HTTP server with request draining.
type DrainableServer interface {
	Shutdown(ctx context.Context) error
}

// DrainConnections creates a shutdown hook for draining connections.
func DrainConnections(name string, server DrainableServer, priority int) ShutdownHook {
	return ShutdownHook{
		Name:     name,
		Priority: priority,
		Fn: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	}
}

// MultiError combines multiple errors.
type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	if len(m.Errors) == 1 {
		return m.Errors[0].Error()
	}
	return errors.Join(m.Errors...).Error()
}

// WaitGroup with context awareness.
type WaitGroupContext struct {
	wg  sync.WaitGroup
	ctx context.Context
}

// NewWaitGroupContext creates a context-aware wait group.
func NewWaitGroupContext(ctx context.Context) *WaitGroupContext {
	return &WaitGroupContext{ctx: ctx}
}

// Go starts a goroutine tracked by the wait group.
func (wg *WaitGroupContext) Go(fn func(ctx context.Context)) {
	wg.wg.Add(1)
	go func() {
		defer wg.wg.Done()
		fn(wg.ctx)
	}()
}

// Wait waits for all goroutines or context cancellation.
func (wg *WaitGroupContext) Wait() error {
	done := make(chan struct{})
	go func() {
		wg.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-wg.ctx.Done():
		return wg.ctx.Err()
	}
}

// WaitWithTimeout waits with a timeout.
func (wg *WaitGroupContext) WaitWithTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		wg.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return context.DeadlineExceeded
	case <-wg.ctx.Done():
		return wg.ctx.Err()
	}
}
