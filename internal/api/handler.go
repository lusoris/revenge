package api

import (
	"context"
	"fmt"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/infra/health"
	"go.uber.org/zap"
)

// Handler implements the ogen.Handler interface for health check endpoints.
type Handler struct {
	logger        *zap.Logger
	healthService *health.Service
}

// GetLiveness implements the liveness probe endpoint.
// This always returns healthy unless the process is deadlocked.
func (h *Handler) GetLiveness(ctx context.Context) (*ogen.HealthCheck, error) {
	h.logger.Debug("Liveness check requested")

	return &ogen.HealthCheck{
		Name:    "liveness",
		Status:  ogen.HealthCheckStatusHealthy,
		Message: ogen.NewOptString("Service is alive"),
	}, nil
}

// GetReadiness implements the readiness probe endpoint.
// Returns healthy only if all dependencies are available.
func (h *Handler) GetReadiness(ctx context.Context) (ogen.GetReadinessRes, error) {
	h.logger.Debug("Readiness check requested")

	// Check if service is ready
	result := h.healthService.Readiness(ctx)

	healthCheck := &ogen.HealthCheck{
		Name: result.Name,
	}

	if result.Status == health.StatusHealthy {
		healthCheck.Status = ogen.HealthCheckStatusHealthy
		healthCheck.Message = ogen.NewOptString(result.Message)
		return (*ogen.GetReadinessOK)(healthCheck), nil
	}

	healthCheck.Status = ogen.HealthCheckStatusUnhealthy
	healthCheck.Message = ogen.NewOptString(result.Message)
	return (*ogen.GetReadinessServiceUnavailable)(healthCheck), nil
}

// GetStartup implements the startup probe endpoint.
// Returns healthy only after initialization is complete.
func (h *Handler) GetStartup(ctx context.Context) (ogen.GetStartupRes, error) {
	h.logger.Debug("Startup check requested")

	// Check if service has started
	result := h.healthService.Startup(ctx)

	healthCheck := &ogen.HealthCheck{
		Name: result.Name,
	}

	if result.Status == health.StatusHealthy {
		healthCheck.Status = ogen.HealthCheckStatusHealthy
		healthCheck.Message = ogen.NewOptString(result.Message)
		return (*ogen.GetStartupOK)(healthCheck), nil
	}

	healthCheck.Status = ogen.HealthCheckStatusUnhealthy
	healthCheck.Message = ogen.NewOptString(result.Message)
	return (*ogen.GetStartupServiceUnavailable)(healthCheck), nil
}

// NewError creates an error response for failed requests.
func (h *Handler) NewError(ctx context.Context, err error) *ogen.ErrorStatusCode {
	h.logger.Error("Request error", zap.Error(err))

	return &ogen.ErrorStatusCode{
		StatusCode: 500,
		Response: ogen.Error{
			Code:    500,
			Message: fmt.Sprintf("Internal server error: %v", err),
		},
	}
}
