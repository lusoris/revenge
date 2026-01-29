// Package middleware provides HTTP middleware for the application.
package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/rbac"
)

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const (
	// UserContextKey is the context key for the authenticated user.
	UserContextKey contextKey = "user"
	// SessionContextKey is the context key for the current session.
	SessionContextKey contextKey = "session"
)

// UserFromContext retrieves the user from the context.
func UserFromContext(ctx context.Context) (*db.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*db.User)
	return user, ok
}

// SessionFromContext retrieves the session from the context.
func SessionFromContext(ctx context.Context) (*db.Session, bool) {
	session, ok := ctx.Value(SessionContextKey).(*db.Session)
	return session, ok
}

// WithUser adds a user to the context.
func WithUser(ctx context.Context, user *db.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// WithSession adds a session to the context.
func WithSession(ctx context.Context, session *db.Session) context.Context {
	return context.WithValue(ctx, SessionContextKey, session)
}

// RBACMiddleware provides permission checking middleware.
type RBACMiddleware struct {
	rbacService *rbac.CasbinService
}

// NewRBACMiddleware creates a new RBAC middleware.
func NewRBACMiddleware(rbacService *rbac.CasbinService) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
	}
}

// RequirePermission returns middleware that checks for a specific permission.
func (m *RBACMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if err := m.rbacService.RequirePermission(r.Context(), user.ID, permission); err != nil {
				if errors.Is(err, rbac.ErrPermissionDenied) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission returns middleware that checks for any of the specified permissions.
func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if err := m.rbacService.RequireAnyPermission(r.Context(), user.ID, permissions); err != nil {
				if errors.Is(err, rbac.ErrPermissionDenied) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole returns middleware that checks for a specific role.
func (m *RBACMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user's role matches any of the required roles
			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// RequireAdmin returns middleware that checks if the user is an admin.
func (m *RBACMiddleware) RequireAdmin() func(http.Handler) http.Handler {
	return m.RequireRole("admin")
}

// RequireAdultAccess returns middleware that checks if user can access adult content.
func (m *RBACMiddleware) RequireAdultAccess() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			canAccess, err := m.rbacService.CanAccessAdultContent(r.Context(), user)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !canAccess {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CheckPermission is a helper function to check permission without middleware.
// Returns the user ID if authorized, or an error.
func CheckPermission(ctx context.Context, rbacService *rbac.CasbinService, permission string) (uuid.UUID, error) {
	user, ok := UserFromContext(ctx)
	if !ok {
		return uuid.Nil, errors.New("unauthorized")
	}

	if err := rbacService.RequirePermission(ctx, user.ID, permission); err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
