package api

import (
	"context"

	gen "github.com/lusoris/revenge/api/generated"
)

// GetServerInfo implements the getServerInfo operation.
func (h *Handler) GetServerInfo(ctx context.Context) (*gen.ServerInfo, error) {
	setupRequired, err := h.authService.IsSetupRequired(ctx)
	if err != nil {
		// If we can't check, assume setup is required
		setupRequired = true
	}

	result := &gen.ServerInfo{
		Version:       h.version,
		SetupRequired: setupRequired,
	}

	if h.buildTime != "" {
		result.BuildTime = gen.NewOptString(h.buildTime)
	}
	if h.gitCommit != "" {
		result.GitCommit = gen.NewOptString(h.gitCommit)
	}

	return result, nil
}

// GetHealth implements the getHealth operation.
func (h *Handler) GetHealth(ctx context.Context) (*gen.HealthStatus, error) {
	result := h.healthChecker.Check(ctx)

	checks := make(gen.HealthStatusChecks)
	for name, svc := range result.Services {
		if svc.Healthy {
			checks[name] = "healthy"
		} else {
			checks[name] = "unhealthy"
		}
	}

	status := gen.HealthStatusStatusHealthy
	if result.Status != "healthy" {
		status = gen.HealthStatusStatusUnhealthy
	}

	return &gen.HealthStatus{
		Status:  status,
		Version: h.version,
		Checks:  gen.NewOptHealthStatusChecks(checks),
	}, nil
}

// GetLiveness implements the getLiveness operation.
func (h *Handler) GetLiveness(ctx context.Context) error {
	// Always return OK - if we can respond, we're alive
	return nil
}

// GetReadiness implements the getReadiness operation.
func (h *Handler) GetReadiness(ctx context.Context) (gen.GetReadinessRes, error) {
	if h.healthChecker.IsReady(ctx) {
		return &gen.GetReadinessOK{}, nil
	}
	return &gen.GetReadinessServiceUnavailable{}, nil
}
