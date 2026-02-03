package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
)

// ErrorResponse represents the JSON error response structure.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ErrorHandler creates a custom error handler that properly handles
// rate limit errors and other custom errors.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	// Check for rate limit error
	if rateLimitErr, ok := err.(*RateLimitError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", rateLimitErr.RetryAfter.String())
		w.WriteHeader(http.StatusTooManyRequests)

		resp := ErrorResponse{
			Error:   "rate_limit_exceeded",
			Message: "Too many requests. Please try again later.",
			Code:    "RATE_LIMIT_EXCEEDED",
		}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// Fall back to default error handler
	ogenerrors.DefaultErrorHandler(ctx, w, r, err)
}
