// Package health provides health check handlers for HTTP endpoints.
package health

import (
	"encoding/json"
	"net/http"
)

// Handler provides HTTP handlers for health checks.
// This is a separate layer from Service for clean separation of concerns.
type Handler struct {
	service *Service
}

// NewHandler creates a new health handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Response represents the JSON response for health endpoints.
type Response struct {
	Status  Status         `json:"status"`
	Message string         `json:"message,omitempty"`
	Checks  []CheckResult  `json:"checks,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

// HandleLiveness handles GET /health/live requests.
// Returns 200 if the service is alive, 503 otherwise.
func (h *Handler) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	result := h.service.Liveness(r.Context())
	h.writeResponse(w, result)
}

// HandleReadiness handles GET /health/ready requests.
// Returns 200 if the service is ready, 503 otherwise.
func (h *Handler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	result := h.service.Readiness(r.Context())
	h.writeResponse(w, result)
}

// HandleStartup handles GET /health/startup requests.
// Returns 200 if startup is complete, 503 otherwise.
func (h *Handler) HandleStartup(w http.ResponseWriter, r *http.Request) {
	result := h.service.Startup(r.Context())
	h.writeResponse(w, result)
}

// HandleFull handles GET /health requests with full dependency checks.
// Returns detailed health information for all dependencies.
func (h *Handler) HandleFull(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Run all checks (returns map[string]CheckResult)
	checksMap := h.service.FullCheck(ctx)

	// Convert map to slice for response
	checks := make([]CheckResult, 0, len(checksMap))
	for _, check := range checksMap {
		checks = append(checks, check)
	}

	// Determine overall status
	overallStatus := StatusHealthy
	for _, check := range checks {
		if check.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
			break
		}
		if check.Status == StatusDegraded && overallStatus != StatusUnhealthy {
			overallStatus = StatusDegraded
		}
	}

	resp := Response{
		Status: overallStatus,
		Checks: checks,
	}

	statusCode := http.StatusOK
	if overallStatus != StatusHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

// writeResponse writes a CheckResult as JSON response.
func (h *Handler) writeResponse(w http.ResponseWriter, result CheckResult) {
	statusCode := http.StatusOK
	if result.Status != StatusHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	resp := Response{
		Status:  result.Status,
		Message: result.Message,
		Details: result.Details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

// RegisterRoutes registers health check routes on a ServeMux.
// This is useful when using standard library routing.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.HandleFull)
	mux.HandleFunc("GET /health/live", h.HandleLiveness)
	mux.HandleFunc("GET /health/ready", h.HandleReadiness)
	mux.HandleFunc("GET /health/startup", h.HandleStartup)
}
