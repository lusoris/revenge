package observability

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ogen-go/ogen/middleware"
)

// HTTPMetricsMiddleware returns an ogen middleware that records HTTP metrics.
func HTTPMetricsMiddleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		start := time.Now()

		// Track in-flight requests
		HTTPRequestsInFlight.Inc()
		defer HTTPRequestsInFlight.Dec()

		// Execute the request
		resp, err := next(req)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Normalize path for metrics (remove dynamic segments)
		path := normalizePath(req.Raw.URL.Path)
		method := req.Raw.Method

		// Determine status code
		statusCode := "500"
		if err == nil {
			// Try to get status from response type name
			statusCode = extractStatusFromResponse(resp)
		}

		// Record metrics
		HTTPRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
		HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

		return resp, err
	}
}

// normalizePath replaces dynamic path segments with placeholders
// to avoid high cardinality in metrics.
func normalizePath(path string) string {
	// Common patterns to normalize:
	// /api/v1/movies/123 -> /api/v1/movies/{id}
	// /api/v1/users/abc-def-ghi -> /api/v1/users/{id}
	// /api/v1/libraries/1/stats -> /api/v1/libraries/{id}/stats

	segments := strings.Split(path, "/")
	normalized := make([]string, len(segments))

	for i, seg := range segments {
		// Skip empty segments
		if seg == "" {
			normalized[i] = seg
			continue
		}

		// Check if this looks like an ID (UUID, numeric, or alphanumeric with hyphens)
		if isIDSegment(seg) {
			normalized[i] = "{id}"
		} else {
			normalized[i] = seg
		}
	}

	return strings.Join(normalized, "/")
}

// isIDSegment checks if a path segment looks like an ID.
func isIDSegment(seg string) bool {
	// Numeric ID
	if _, err := strconv.ParseInt(seg, 10, 64); err == nil {
		return true
	}

	// UUID pattern (8-4-4-4-12 hex chars with hyphens)
	if len(seg) == 36 && strings.Count(seg, "-") == 4 {
		return true
	}

	// Short UUID or other ID-like patterns (alphanumeric, 8+ chars)
	if len(seg) >= 8 && isAlphanumericWithHyphens(seg) && hasDigits(seg) {
		return true
	}

	return false
}

// isAlphanumericWithHyphens checks if string contains only alphanumeric chars and hyphens.
func isAlphanumericWithHyphens(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

// hasDigits checks if string contains at least one digit.
func hasDigits(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

// extractStatusFromResponse tries to extract status code from ogen response.
// This is a best-effort approach since ogen responses are type-safe.
func extractStatusFromResponse(resp middleware.Response) string {
	// ogen responses typically have a type name like "GetMovieOK" or "GetMovieNotFound"
	// We try to extract the status from common suffixes
	typeName := fmt.Sprintf("%T", resp.Type)

	switch {
	case strings.HasSuffix(typeName, "OK"):
		return "200"
	case strings.HasSuffix(typeName, "Created"):
		return "201"
	case strings.HasSuffix(typeName, "NoContent"):
		return "204"
	case strings.HasSuffix(typeName, "BadRequest"):
		return "400"
	case strings.HasSuffix(typeName, "Unauthorized"):
		return "401"
	case strings.HasSuffix(typeName, "Forbidden"):
		return "403"
	case strings.HasSuffix(typeName, "NotFound"):
		return "404"
	case strings.HasSuffix(typeName, "Conflict"):
		return "409"
	case strings.HasSuffix(typeName, "TooManyRequests"):
		return "429"
	case strings.HasSuffix(typeName, "InternalServerError"):
		return "500"
	default:
		return "200" // Default to 200 if we can't determine
	}
}

// StandardHTTPMetricsMiddleware returns a standard http.Handler middleware for non-ogen routes.
func StandardHTTPMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Track in-flight requests
		HTTPRequestsInFlight.Inc()
		defer HTTPRequestsInFlight.Dec()

		// Wrap response writer to capture status
		wrapped := &statusResponseWriter{ResponseWriter: w, status: 200}

		// Execute the request
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Normalize path
		path := normalizePath(r.URL.Path)
		method := r.Method
		status := strconv.Itoa(wrapped.status)

		// Record metrics
		HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	})
}

// statusResponseWriter wraps http.ResponseWriter to capture status code.
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
