package api

import (
	"log/slog"

	"go.uber.org/fx"

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
	HealthChecker  *health.Checker
	Logger         *slog.Logger
}

// provideHandlerParams creates HandlerParams from dependencies.
func provideHandlerParams(deps HandlerDeps) HandlerParams {
	return HandlerParams{
		AuthService:    deps.AuthService,
		UserService:    deps.UserService,
		SessionService: deps.SessionService,
		LibraryService: deps.LibraryService,
		HealthChecker:  deps.HealthChecker,
		Logger:         deps.Logger,
		Version:        "0.1.0",  // TODO: inject from build
		BuildTime:      "",       // TODO: inject from build
		GitCommit:      "",       // TODO: inject from build
	}
}
