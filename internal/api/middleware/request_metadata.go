// Package middleware provides HTTP middleware for the API server.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ogen-go/ogen/middleware"
)

// requestMetaKey is the context key for request metadata
type requestMetaKey struct{}

// RequestMetadata holds metadata extracted from HTTP requests.
// Used for session creation, activity logging, and security auditing.
type RequestMetadata struct {
	// IPAddress is the client IP (with X-Forwarded-For support)
	IPAddress string
	// UserAgent is the client's User-Agent header
	UserAgent string
	// AcceptLanguage is the client's Accept-Language header
	AcceptLanguage string
}

// WithRequestMetadata stores request metadata in the context
func WithRequestMetadata(ctx context.Context, meta RequestMetadata) context.Context {
	return context.WithValue(ctx, requestMetaKey{}, meta)
}

// GetRequestMetadata retrieves request metadata from the context.
// Returns empty metadata if not found (never returns error for convenience).
func GetRequestMetadata(ctx context.Context) RequestMetadata {
	meta, ok := ctx.Value(requestMetaKey{}).(RequestMetadata)
	if !ok {
		return RequestMetadata{}
	}
	return meta
}

// GetIPAddress is a convenience function to get just the IP address from context.
func GetIPAddress(ctx context.Context) string {
	return GetRequestMetadata(ctx).IPAddress
}

// GetUserAgent is a convenience function to get just the user agent from context.
func GetUserAgent(ctx context.Context) string {
	return GetRequestMetadata(ctx).UserAgent
}

// RequestMetadataMiddleware extracts client metadata from HTTP requests
// and stores it in the context for use by handlers.
//
// Extracted metadata:
// - IP address (with X-Forwarded-For, X-Real-IP support for proxies)
// - User-Agent header
// - Accept-Language header
func RequestMetadataMiddleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		meta := RequestMetadata{
			IPAddress:      extractClientIP(req.Raw),
			UserAgent:      req.Raw.Header.Get("User-Agent"),
			AcceptLanguage: req.Raw.Header.Get("Accept-Language"),
		}

		// Store in context
		ctx := WithRequestMetadata(req.Context, meta)
		req.Context = ctx

		return next(req)
	}
}

// extractClientIP extracts the client IP address from the request,
// supporting common proxy headers.
//
// Priority:
// 1. X-Forwarded-For (first IP in chain, leftmost is original client)
// 2. X-Real-IP
// 3. RemoteAddr
func extractClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (standard proxy header)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// The first IP is the original client
		if idx := strings.IndexByte(xff, ','); idx > 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header (nginx default)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr (remove port if present)
	return stripPort(r.RemoteAddr)
}

// stripPort removes the port from an address, handling both IPv4 and IPv6.
// Examples:
//   - "192.168.1.1:8080" -> "192.168.1.1"
//   - "[::1]:8080" -> "::1"
//   - "::1" -> "::1" (no port, unchanged)
//   - "192.168.1.1" -> "192.168.1.1" (no port, unchanged)
func stripPort(addr string) string {
	// Handle IPv6 with brackets: [::1]:8080
	if strings.HasPrefix(addr, "[") {
		if idx := strings.LastIndexByte(addr, ']'); idx > 0 {
			return addr[1:idx]
		}
		return addr // malformed, return as-is
	}

	// Count colons to distinguish IPv6 without brackets from IPv4 with port
	colonCount := strings.Count(addr, ":")

	// IPv6 without brackets (no port): multiple colons like "::1" or "2001:db8::1"
	if colonCount > 1 {
		return addr
	}

	// IPv4 with port: exactly one colon like "192.168.1.1:8080"
	if colonCount == 1 {
		if idx := strings.LastIndexByte(addr, ':'); idx > 0 {
			return addr[:idx]
		}
	}

	// No colons (IPv4 without port) or malformed
	return addr
}
