// Package lazy provides lazy initialization utilities.
package lazy

import (
	"sync"
	"sync/atomic"
	"time"
)

// Service wraps a lazily-initialized service.
// The service is only created when Get() is first called.
type Service[T any] struct {
	init     func() (T, error)
	instance atomic.Pointer[T]
	once     sync.Once
	err      error
	initTime time.Duration
}

// New creates a new lazy service wrapper.
func New[T any](init func() (T, error)) *Service[T] {
	return &Service[T]{init: init}
}

// Get returns the service instance, initializing it on first call.
// Subsequent calls return the cached instance.
func (s *Service[T]) Get() (T, error) {
	s.once.Do(func() {
		start := time.Now()
		instance, err := s.init()
		s.initTime = time.Since(start)
		if err != nil {
			s.err = err
			return
		}
		s.instance.Store(&instance)
	})

	if s.err != nil {
		var zero T
		return zero, s.err
	}

	return *s.instance.Load(), nil
}

// MustGet returns the service instance or panics on error.
// Use only when initialization cannot fail.
func (s *Service[T]) MustGet() T {
	instance, err := s.Get()
	if err != nil {
		panic("lazy service initialization failed: " + err.Error())
	}
	return instance
}

// IsInitialized returns whether the service has been initialized.
func (s *Service[T]) IsInitialized() bool {
	return s.instance.Load() != nil
}

// InitTime returns the time taken to initialize the service.
// Returns 0 if not yet initialized.
func (s *Service[T]) InitTime() time.Duration {
	return s.initTime
}

// Error returns the initialization error, if any.
func (s *Service[T]) Error() error {
	return s.err
}

// Reset resets the lazy service to uninitialized state.
// This is primarily useful for testing.
// Warning: This is not thread-safe with concurrent Get() calls.
func (s *Service[T]) Reset() {
	s.once = sync.Once{}
	s.instance.Store(nil)
	s.err = nil
	s.initTime = 0
}

// ServiceWithCleanup wraps a lazily-initialized service with cleanup.
type ServiceWithCleanup[T any] struct {
	Service[T]
	cleanup func(T) error
}

// NewWithCleanup creates a lazy service with a cleanup function.
func NewWithCleanup[T any](init func() (T, error), cleanup func(T) error) *ServiceWithCleanup[T] {
	return &ServiceWithCleanup[T]{
		Service: Service[T]{init: init},
		cleanup: cleanup,
	}
}

// Close calls the cleanup function if the service was initialized.
func (s *ServiceWithCleanup[T]) Close() error {
	if !s.IsInitialized() {
		return nil
	}

	instance, err := s.Get()
	if err != nil {
		return err
	}

	return s.cleanup(instance)
}

// Pool manages a pool of lazy services.
type Pool[T any] struct {
	services []*Service[T]
	mu       sync.Mutex
	index    int
}

// NewPool creates a pool of lazy services.
func NewPool[T any](size int, init func() (T, error)) *Pool[T] {
	services := make([]*Service[T], size)
	for i := range size {
		services[i] = New(init)
	}
	return &Pool[T]{services: services}
}

// Get returns the next service from the pool (round-robin).
func (p *Pool[T]) Get() (T, error) {
	p.mu.Lock()
	service := p.services[p.index]
	p.index = (p.index + 1) % len(p.services)
	p.mu.Unlock()

	return service.Get()
}

// InitializedCount returns how many services in the pool are initialized.
func (p *Pool[T]) InitializedCount() int {
	count := 0
	for _, s := range p.services {
		if s.IsInitialized() {
			count++
		}
	}
	return count
}
