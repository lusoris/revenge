//go:build ogen

package api

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/pkg/health"
)

// Module provides API handler dependencies.
var Module = fx.Module("api",
	fx.Provide(
		NewHandler,
		NewSecurityHandler,
		provideHandlerParams,
	),
)

// HandlerDeps contains the dependencies for creating a Handler.
type HandlerDeps struct {
	fx.In

	AuthService    *auth.Service
	UserService    *user.Service
	SessionService *session.Service
	LibraryService *library.Service
	MovieService   *movie.Service `optional:"true"`
	RiverClient    *river.Client[pgx.Tx]
	HealthChecker  *health.Checker
	Logger         *slog.Logger
	BuildInfo      BuildInfo
	Config         *config.Config
}

// provideHandlerParams creates HandlerParams from dependencies.
func provideHandlerParams(deps HandlerDeps) HandlerParams {
	return HandlerParams{
		AuthService:    deps.AuthService,
		UserService:    deps.UserService,
		SessionService: deps.SessionService,
		LibraryService: deps.LibraryService,
		MovieService:   deps.MovieService,
		RiverClient:    deps.RiverClient,
		HealthChecker:  deps.HealthChecker,
		Logger:         deps.Logger,
		AdultEnabled:   deps.Config.Modules.Adult,
		Version:        deps.BuildInfo.Version,
		BuildTime:      deps.BuildInfo.BuildTime,
		GitCommit:      deps.BuildInfo.GitCommit,
	}
}
