// Package middleware provides HTTP middleware for the Jellyfin Go API.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// UserContextKey is the context key for the authenticated user.
	UserContextKey contextKey = "user"
	// ClaimsContextKey is the context key for the token claims.
	ClaimsContextKey contextKey = "claims"
	// SessionContextKey is the context key for the session.
	SessionContextKey contextKey = "session"
)

// Auth is middleware that validates JWT tokens from the Authorization header.
type Auth struct {
	authService domain.AuthService
}

// NewAuth creates a new Auth middleware.
func NewAuth(authService domain.AuthService) *Auth {
	return &Auth{authService: authService}
}

// Required returns middleware that requires a valid JWT token.
// Returns 401 Unauthorized if the token is missing or invalid.
func (a *Auth) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			slog.Debug("missing authorization token",
				slog.String("path", r.URL.Path))
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		claims, err := a.authService.ValidateToken(r.Context(), token)
		if err != nil {
			slog.Debug("invalid token",
				slog.String("path", r.URL.Path),
				slog.Any("error", err))
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequiredWithSession returns middleware that requires a valid JWT token
// and loads the full session with user info.
func (a *Auth) RequiredWithSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			slog.Debug("missing authorization token",
				slog.String("path", r.URL.Path))
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		session, err := a.authService.GetSession(r.Context(), token)
		if err != nil {
			slog.Debug("invalid token or session",
				slog.String("path", r.URL.Path),
				slog.Any("error", err))
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add session to context (SessionWithUser contains user info via embedded fields)
		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Optional returns middleware that validates the token if present,
// but allows the request to proceed without authentication.
func (a *Auth) Optional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := a.authService.ValidateToken(r.Context(), token)
		if err != nil {
			// Token present but invalid - still allow request but log it
			slog.Debug("optional auth: invalid token",
				slog.String("path", r.URL.Path),
				slog.Any("error", err))
			next.ServeHTTP(w, r)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequired returns middleware that requires admin privileges.
// Must be used after Required or RequiredWithSession middleware.
func (a *Auth) AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := ClaimsFromContext(r.Context())
		if claims == nil {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		if !claims.IsAdmin {
			slog.Warn("non-admin access attempt to admin endpoint",
				slog.String("user_id", claims.UserID.String()),
				slog.String("path", r.URL.Path))
			http.Error(w, "Administrator privileges required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// extractToken extracts the JWT token from the Authorization header.
// Supports both "Bearer <token>" and "MediaBrowser Token=<token>" formats
// for Jellyfin client compatibility.
func extractToken(r *http.Request) string {
	// Check Authorization header
	auth := r.Header.Get("Authorization")
	if auth != "" {
		// Standard Bearer token
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}

		// Jellyfin MediaBrowser format: MediaBrowser Token="<token>", ...
		if strings.HasPrefix(auth, "MediaBrowser ") {
			params := parseMediaBrowserAuth(auth)
			if token, ok := params["Token"]; ok {
				return token
			}
		}
	}

	// Check X-Emby-Authorization header (alternative Jellyfin header)
	embyAuth := r.Header.Get("X-Emby-Authorization")
	if embyAuth != "" {
		params := parseMediaBrowserAuth(embyAuth)
		if token, ok := params["Token"]; ok {
			return token
		}
	}

	// Check query parameter (for WebSocket connections, etc.)
	if token := r.URL.Query().Get("api_key"); token != "" {
		return token
	}

	return ""
}

// parseMediaBrowserAuth parses the MediaBrowser authorization header format.
// Format: MediaBrowser Client="...", Device="...", DeviceId="...", Version="...", Token="..."
func parseMediaBrowserAuth(header string) map[string]string {
	params := make(map[string]string)

	// Remove "MediaBrowser " prefix if present
	header = strings.TrimPrefix(header, "MediaBrowser ")

	// Split by comma and parse key="value" pairs
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "="); idx > 0 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])
			// Remove quotes if present
			value = strings.Trim(value, "\"")
			params[key] = value
		}
	}

	return params
}

// ClaimsFromContext extracts TokenClaims from the request context.
func ClaimsFromContext(ctx context.Context) *domain.TokenClaims {
	claims, _ := ctx.Value(ClaimsContextKey).(*domain.TokenClaims) //nolint:errcheck // type assertion ok flag intentionally ignored
	return claims
}

// UserFromContext extracts the User from the request context.
func UserFromContext(ctx context.Context) *domain.User {
	user, _ := ctx.Value(UserContextKey).(*domain.User) //nolint:errcheck // type assertion ok flag intentionally ignored
	return user
}

// SessionFromContext extracts the SessionWithUser from the request context.
func SessionFromContext(ctx context.Context) *domain.SessionWithUser {
	session, _ := ctx.Value(SessionContextKey).(*domain.SessionWithUser) //nolint:errcheck // type assertion ok flag intentionally ignored
	return session
}
