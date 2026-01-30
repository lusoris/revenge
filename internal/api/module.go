package api

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/qar/crew"
	"github.com/lusoris/revenge/internal/content/qar/expedition"
	"github.com/lusoris/revenge/internal/content/qar/flag"
	"github.com/lusoris/revenge/internal/content/qar/fleet"
	"github.com/lusoris/revenge/internal/content/qar/port"
	"github.com/lusoris/revenge/internal/content/qar/voyage"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/rbac"
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

	Config         *config.Config
	AuthService    *auth.Service
	UserService    *user.Service
	SessionService *session.Service
	LibraryService *library.Service
	RBACService    *rbac.CasbinService
	MovieService   *movie.Service  `optional:"true"`
	TVShowService  *tvshow.Service `optional:"true"`
	RiverClient    *river.Client[pgx.Tx]
	HealthChecker  *health.Checker
	Logger         *slog.Logger
	BuildInfo      BuildInfo

	// QAR (adult content) services - all optional
	ExpeditionService *expedition.Service `optional:"true"`
	VoyageService     *voyage.Service     `optional:"true"`
	CrewService       *crew.Service       `optional:"true"`
	PortService       *port.Service       `optional:"true"`
	FlagService       *flag.Service       `optional:"true"`
	FleetService      *fleet.Service      `optional:"true"`
}

// provideHandlerParams creates HandlerParams from dependencies.
func provideHandlerParams(deps HandlerDeps) HandlerParams {
	return HandlerParams{
		AuthService:       deps.AuthService,
		UserService:       deps.UserService,
		SessionService:    deps.SessionService,
		LibraryService:    deps.LibraryService,
		RBACService:       deps.RBACService,
		MovieService:      deps.MovieService,
		TVShowService:     deps.TVShowService,
		RiverClient:       deps.RiverClient,
		HealthChecker:     deps.HealthChecker,
		Logger:            deps.Logger,
		AdultEnabled:      deps.Config.Adult.Enabled,
		Version:           deps.BuildInfo.Version,
		BuildTime:         deps.BuildInfo.BuildTime,
		GitCommit:         deps.BuildInfo.GitCommit,
		ExpeditionService: deps.ExpeditionService,
		VoyageService:     deps.VoyageService,
		CrewService:       deps.CrewService,
		PortService:       deps.PortService,
		FlagService:       deps.FlagService,
		FleetService:      deps.FleetService,
	}
}
