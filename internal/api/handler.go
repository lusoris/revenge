// Package api provides HTTP API handlers for the Revenge media server.
package api

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/pkg/health"
)

// ctxKey is a type for context keys.
type ctxKey string

const (
	// ctxKeyUser is the context key for the authenticated user.
	ctxKeyUser ctxKey = "user"
	// ctxKeySession is the context key for the current session.
	ctxKeySession ctxKey = "session"
)

// Handler implements the generated API handler interface.
type Handler struct {
	gen.UnimplementedHandler
	authService    *auth.Service
	userService    *user.Service
	sessionService *session.Service
	libraryService *library.Service
	movieService   *movie.Service
	tvshowService  *tvshow.Service
	riverClient    *river.Client[pgx.Tx]
	healthChecker  *health.Checker
	logger         *slog.Logger
	adultEnabled   bool
	version        string
	buildTime      string
	gitCommit      string
}

// BuildInfo contains build metadata for server info responses.
type BuildInfo struct {
	Version   string
	BuildTime string
	GitCommit string
}

// HandlerParams contains parameters for creating a new Handler.
type HandlerParams struct {
	AuthService    *auth.Service
	UserService    *user.Service
	SessionService *session.Service
	LibraryService *library.Service
	MovieService   *movie.Service
	TVShowService  *tvshow.Service
	RiverClient    *river.Client[pgx.Tx]
	HealthChecker  *health.Checker
	Logger         *slog.Logger
	AdultEnabled   bool
	Version        string
	BuildTime      string
	GitCommit      string
}

// NewHandler creates a new API handler.
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		authService:    params.AuthService,
		userService:    params.UserService,
		sessionService: params.SessionService,
		libraryService: params.LibraryService,
		movieService:   params.MovieService,
		tvshowService:  params.TVShowService,
		riverClient:    params.RiverClient,
		healthChecker:  params.HealthChecker,
		logger:         params.Logger.With(slog.String("component", "api")),
		adultEnabled:   params.AdultEnabled,
		version:        params.Version,
		buildTime:      params.BuildTime,
		gitCommit:      params.GitCommit,
	}
}

// UserFromContext retrieves the authenticated user from context.
func UserFromContext(ctx context.Context) (*db.User, bool) {
	u, ok := ctx.Value(ctxKeyUser).(*db.User)
	return u, ok
}

// SessionFromContext retrieves the current session from context.
func SessionFromContext(ctx context.Context) (*db.Session, bool) {
	s, ok := ctx.Value(ctxKeySession).(*db.Session)
	return s, ok
}

// contextWithUser adds a user to the context.
func contextWithUser(ctx context.Context, u *db.User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, u)
}

// contextWithSession adds a session to the context.
func contextWithSession(ctx context.Context, s *db.Session) context.Context {
	return context.WithValue(ctx, ctxKeySession, s)
}

// requireUser gets the user from context or returns an error.
func requireUser(ctx context.Context) (*db.User, error) {
	u, ok := UserFromContext(ctx)
	if !ok || u == nil {
		return nil, ErrUnauthorized
	}
	return u, nil
}

// requireSession gets the session from context or returns an error.
func requireSession(ctx context.Context) (*db.Session, error) {
	s, ok := SessionFromContext(ctx)
	if !ok || s == nil {
		return nil, ErrUnauthorized
	}
	return s, nil
}

// requireAdmin checks if the user is an admin.
func requireAdmin(ctx context.Context) (*db.User, error) {
	u, err := requireUser(ctx)
	if err != nil {
		return nil, err
	}
	if !u.IsAdmin {
		return nil, ErrForbidden
	}
	return u, nil
}

func (h *Handler) requireMovieService() (*movie.Service, error) {
	if h.movieService == nil {
		return nil, ErrModuleDisabled
	}
	return h.movieService, nil
}

// ptrString returns a pointer to a string, or nil if empty.
func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ptrInt32 returns a pointer to an int32.
func ptrInt32(i int) *int32 {
	i32 := int32(i)
	return &i32
}

// ptrBool returns a pointer to a bool.
func ptrBool(b bool) *bool {
	return &b
}
