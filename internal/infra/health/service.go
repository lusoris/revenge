// Package health provides health check endpoints for Kubernetes and monitoring.
package health

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/database"
)

// Status represents the health status.
type Status string

const (
	// StatusHealthy indicates the service is healthy.
	StatusHealthy Status = "healthy"

	// StatusUnhealthy indicates the service is unhealthy.
	StatusUnhealthy Status = "unhealthy"

	// StatusDegraded indicates the service is degraded.
	StatusDegraded Status = "degraded"
)

// CheckResult represents the result of a health check.
type CheckResult struct {
	Name    string                 `json:"name"`
	Status  Status                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Service provides health check functionality.
type Service struct {
	logger *slog.Logger
	pool   *pgxpool.Pool

	// Startup state
	startupComplete bool
	startupMu       sync.RWMutex
}

// NewService creates a new health service.
func NewService(logger *slog.Logger, pool *pgxpool.Pool) *Service {
	return &Service{
		logger:          logger,
		pool:            pool,
		startupComplete: false,
	}
}

// MarkStartupComplete marks the startup process as complete.
// This allows the startup probe to succeed.
func (s *Service) MarkStartupComplete() {
	s.startupMu.Lock()
	defer s.startupMu.Unlock()
	s.startupComplete = true
	s.logger.Info("startup complete")
}

// Liveness checks if the service is alive.
// This should always return healthy unless the process is deadlocked.
func (s *Service) Liveness(ctx context.Context) CheckResult {
	return CheckResult{
		Name:    "liveness",
		Status:  StatusHealthy,
		Message: "service is alive",
	}
}

// Readiness checks if the service is ready to accept traffic.
// This checks if all dependencies are available.
func (s *Service) Readiness(ctx context.Context) CheckResult {
	// Check if startup is complete
	s.startupMu.RLock()
	if !s.startupComplete {
		s.startupMu.RUnlock()
		return CheckResult{
			Name:    "readiness",
			Status:  StatusUnhealthy,
			Message: "startup not complete",
		}
	}
	s.startupMu.RUnlock()

	// Check database
	dbCheck := s.checkDatabase(ctx)
	if dbCheck.Status != StatusHealthy {
		return CheckResult{
			Name:    "readiness",
			Status:  StatusUnhealthy,
			Message: "database not ready",
			Details: map[string]interface{}{
				"database": dbCheck,
			},
		}
	}

	return CheckResult{
		Name:    "readiness",
		Status:  StatusHealthy,
		Message: "service is ready",
		Details: map[string]interface{}{
			"database": dbCheck,
		},
	}
}

// Startup checks if the service has completed initialization.
// This is used by Kubernetes startup probes.
func (s *Service) Startup(ctx context.Context) CheckResult {
	s.startupMu.RLock()
	complete := s.startupComplete
	s.startupMu.RUnlock()

	if !complete {
		return CheckResult{
			Name:    "startup",
			Status:  StatusUnhealthy,
			Message: "initialization in progress",
		}
	}

	return CheckResult{
		Name:    "startup",
		Status:  StatusHealthy,
		Message: "startup complete",
	}
}

// checkDatabase checks if the database is healthy.
func (s *Service) checkDatabase(ctx context.Context) CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := database.Health(ctx, s.pool); err != nil {
		return CheckResult{
			Name:    "database",
			Status:  StatusUnhealthy,
			Message: err.Error(),
		}
	}

	// Get pool stats
	stats := database.Stats(s.pool)

	return CheckResult{
		Name:    "database",
		Status:  StatusHealthy,
		Message: "database is healthy",
		Details: stats,
	}
}

// FullCheck runs all health checks and returns a summary.
func (s *Service) FullCheck(ctx context.Context) map[string]CheckResult {
	return map[string]CheckResult{
		"liveness":  s.Liveness(ctx),
		"readiness": s.Readiness(ctx),
		"startup":   s.Startup(ctx),
	}
}
