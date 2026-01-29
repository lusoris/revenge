# HTTP Handler Patterns

> Instructions for implementing HTTP handlers, middleware, and API converters in Revenge.

## Handler Structure

Handlers implement ogen-generated interfaces:

```go
package handler

import (
    "context"

    "github.com/lusoris/revenge/internal/api"
    "github.com/lusoris/revenge/internal/service/movie"
)

type MovieHandler struct {
    service *movie.Service
}

func NewMovieHandler(service *movie.Service) *MovieHandler {
    return &MovieHandler{service: service}
}

// Implement ogen-generated interface
var _ api.MoviesHandler = (*MovieHandler)(nil)
```

---

## Response Patterns

### Success Responses

```go
func (h *MovieHandler) GetMovie(ctx context.Context, params api.GetMovieParams) (api.GetMovieRes, error) {
    movie, err := h.service.GetByID(ctx, params.ID)
    if err != nil {
        if errors.Is(err, movie.ErrNotFound) {
            return &api.GetMovieNotFound{}, nil
        }
        return nil, err // 500 Internal Server Error
    }

    return convertMovieToAPI(movie), nil
}
```

### Error Responses

ogen generates typed error responses. Use them:

```go
func (h *MovieHandler) CreateMovie(ctx context.Context, req *api.CreateMovieRequest) (api.CreateMovieRes, error) {
    movie, err := h.service.Create(ctx, convertAPIToMovie(req))
    if err != nil {
        switch {
        case errors.Is(err, movie.ErrDuplicate):
            return &api.CreateMovieConflict{
                Message: "Movie already exists",
            }, nil
        case errors.Is(err, movie.ErrInvalidInput):
            return &api.CreateMovieBadRequest{
                Message: err.Error(),
            }, nil
        default:
            return nil, err
        }
    }

    return convertMovieToAPI(movie), nil
}
```

---

## Middleware

Location: `internal/api/middleware/`

### Authentication Middleware

```go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/lusoris/revenge/internal/service/auth"
)

type AuthMiddleware struct {
    auth *auth.Service
}

func NewAuthMiddleware(auth *auth.Service) *AuthMiddleware {
    return &AuthMiddleware{auth: auth}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        if token == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        session, err := m.auth.ValidateToken(r.Context(), token)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), userContextKey, session.User)
        ctx = context.WithValue(ctx, sessionContextKey, session)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func extractToken(r *http.Request) string {
    // Check Authorization header
    auth := r.Header.Get("Authorization")
    if strings.HasPrefix(auth, "Bearer ") {
        return strings.TrimPrefix(auth, "Bearer ")
    }

    // Check X-Emby-Token header (Jellyfin compatibility)
    if token := r.Header.Get("X-Emby-Token"); token != "" {
        return token
    }

    // Check query parameter
    return r.URL.Query().Get("api_key")
}
```

### Context Helpers

```go
package middleware

type contextKey string

const (
    userContextKey    contextKey = "user"
    sessionContextKey contextKey = "session"
)

func UserFromContext(ctx context.Context) *user.User {
    u, _ := ctx.Value(userContextKey).(*user.User)
    return u
}

func SessionFromContext(ctx context.Context) *session.Session {
    s, _ := ctx.Value(sessionContextKey).(*session.Session)
    return s
}
```

### RBAC Middleware

```go
package middleware

type RBACMiddleware struct {
    rbac *rbac.Service
}

func (m *RBACMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := UserFromContext(r.Context())
            if user == nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            allowed, err := m.rbac.CheckPermission(r.Context(), user.ID, permission)
            if err != nil || !allowed {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### Request Logging

```go
package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

func RequestLogging(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            wrapped := &responseWriter{ResponseWriter: w, status: 200}

            next.ServeHTTP(wrapped, r)

            logger.Info("request",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.Int("status", wrapped.status),
                slog.Duration("duration", time.Since(start)),
            )
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (w *responseWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}
```

---

## API Converters

Location: `internal/api/converter/`

Convert between domain entities and API types:

```go
package converter

import (
    "github.com/lusoris/revenge/internal/api"
    "github.com/lusoris/revenge/internal/content/movie"
)

// MovieToAPI converts domain Movie to API response
func MovieToAPI(m *movie.Movie) *api.Movie {
    return &api.Movie{
        ID:              m.ID,
        Title:           m.Title,
        OriginalTitle:   api.NewOptString(m.OriginalTitle),
        Overview:        api.NewOptString(m.Overview),
        Tagline:         api.NewOptString(m.Tagline),
        ReleaseDate:     optTimeFromPtr(m.ReleaseDate),
        RuntimeTicks:    api.NewOptInt64(m.RuntimeTicks),
        CommunityRating: api.NewOptFloat64(m.CommunityRating),
        PosterPath:      api.NewOptString(m.PosterPath),
        BackdropPath:    api.NewOptString(m.BackdropPath),
        TmdbID:          optIntFromPtr(m.TmdbID),
        ImdbID:          api.NewOptString(m.ImdbID),
    }
}

// MoviesToAPI converts slice of domain Movies
func MoviesToAPI(movies []*movie.Movie) []api.Movie {
    result := make([]api.Movie, len(movies))
    for i, m := range movies {
        result[i] = *MovieToAPI(m)
    }
    return result
}

// MovieFromAPI converts API request to domain Movie
func MovieFromAPI(req *api.CreateMovieRequest) *movie.Movie {
    return &movie.Movie{
        Title:         req.Title,
        OriginalTitle: req.OriginalTitle.Or(""),
        Overview:      req.Overview.Or(""),
        ReleaseDate:   ptrFromOptTime(req.ReleaseDate),
    }
}

// Helper functions
func optTimeFromPtr(t *time.Time) api.OptDateTime {
    if t == nil {
        return api.OptDateTime{}
    }
    return api.NewOptDateTime(*t)
}

func ptrFromOptTime(opt api.OptDateTime) *time.Time {
    if !opt.Set {
        return nil
    }
    return &opt.Value
}

func optIntFromPtr(i *int) api.OptInt {
    if i == nil {
        return api.OptInt{}
    }
    return api.NewOptInt(*i)
}
```

### Pagination Helpers

```go
package converter

// PaginationFromAPI extracts pagination params
func PaginationFromAPI(limit, offset api.OptInt) (int, int) {
    l := limit.Or(20)  // default limit
    o := offset.Or(0)  // default offset

    // Enforce max limit
    if l > 100 {
        l = 100
    }

    return l, o
}

// PaginationToAPI creates pagination response
func PaginationToAPI(total, limit, offset int) api.Pagination {
    return api.Pagination{
        Total:  total,
        Limit:  limit,
        Offset: offset,
    }
}
```

---

## Security Context

### ogen Security Handler

```go
package api

import (
    "context"

    "github.com/lusoris/revenge/internal/service/auth"
)

type SecurityHandler struct {
    auth *auth.Service
}

func NewSecurityHandler(auth *auth.Service) *SecurityHandler {
    return &SecurityHandler{auth: auth}
}

// HandleBearerAuth implements ogen security handler
func (h *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
    session, err := h.auth.ValidateToken(ctx, t.Token)
    if err != nil {
        return ctx, err
    }

    ctx = context.WithValue(ctx, userContextKey, session.User)
    ctx = context.WithValue(ctx, sessionContextKey, session)
    return ctx, nil
}

// HandleApiKeyAuth implements ogen security handler for X-Emby-Token
func (h *SecurityHandler) HandleApiKeyAuth(ctx context.Context, operationName string, t api.ApiKeyAuth) (context.Context, error) {
    session, err := h.auth.ValidateToken(ctx, t.APIKey)
    if err != nil {
        return ctx, err
    }

    ctx = context.WithValue(ctx, userContextKey, session.User)
    return ctx, nil
}
```

---

## Handler Registration

```go
package handler

import "go.uber.org/fx"

var Module = fx.Module("handlers",
    fx.Provide(NewMovieHandler),
    fx.Provide(NewTVShowHandler),
    fx.Provide(NewLibraryHandler),
    fx.Provide(NewUserHandler),
    fx.Provide(NewAuthHandler),
    fx.Provide(NewSecurityHandler),
)
```

---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index
- [ogen-api.instructions.md](ogen-api.instructions.md) - ogen API generation
- [service-layer-patterns.instructions.md](service-layer-patterns.instructions.md) - Service patterns
- [rbac-casbin.instructions.md](rbac-casbin.instructions.md) - RBAC patterns
- [API Design](../../docs/dev/design/technical/API.md) - API documentation
