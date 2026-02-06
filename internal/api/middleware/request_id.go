package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
)

// requestIDKey is the context key for request ID
type requestIDKey struct{}

// RequestIDMiddleware extracts or generates X-Request-ID header for request correlation.
//
// Request ID handling:
// - If X-Request-ID header is present, use it
// - If not present, generate a new UUID
// - Store in context for access by handlers and logging
//
// Note: This middleware only handles context. Use RequestIDHTTPWrapper to add
// X-Request-ID to response headers.
//
// This enables:
// - Request tracing across services
// - Correlation of logs for a single request
// - Debugging distributed systems
func RequestIDMiddleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// Extract or generate request ID
		requestID := req.Raw.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.Must(uuid.NewV7()).String()
		}

		// Store in context
		ctx := context.WithValue(req.Context, requestIDKey{}, requestID)
		req.Context = ctx

		// Execute next middleware/handler
		return next(req)
	}
}

// RequestIDHTTPWrapper is an HTTP middleware that adds X-Request-ID to response headers.
// Should wrap the ogen server to ensure the response header is always set.
//
// Usage:
//
//	httpServer := &http.Server{
//	    Handler: middleware.RequestIDHTTPWrapper(ogenServer),
//	}
func RequestIDHTTPWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract request ID from context (set by RequestIDMiddleware)
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			// Fallback: generate if not in context (shouldn't happen if middleware is configured correctly)
			requestID = uuid.Must(uuid.NewV7()).String()
			r = r.WithContext(WithRequestID(r.Context(), requestID))
		}

		// Set response header
		w.Header().Set("X-Request-ID", requestID)

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from the context.
// Returns empty string if not found.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// WithRequestID stores a request ID in the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}
