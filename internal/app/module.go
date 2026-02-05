// Package app provides the main application module that wires all dependencies together.
package app

import (
	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
	"github.com/lusoris/revenge/internal/content/tvshow"
	tvshowjobs "github.com/lusoris/revenge/internal/content/tvshow/jobs"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/infra/observability"
	"github.com/lusoris/revenge/internal/infra/raft"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/integration/radarr"
	"github.com/lusoris/revenge/internal/integration/sonarr"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/metadata/metadatafx"
	"github.com/lusoris/revenge/internal/service/mfa"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"go.uber.org/fx"
)

// Module is the main application module that includes all sub-modules.
var Module = fx.Module("app",
	// Configuration
	config.Module,

	// Infrastructure
	logging.Module,
	database.Module,
	cache.Module,
	search.Module,
	jobs.Module,
	raft.Module,
	health.Module,

	// Services
	settings.Module,
	user.Module,
	auth.Module,
	session.Module,
	rbac.Module,
	apikeys.Module,
	mfa.Module,
	oidc.Module,
	activity.Module,
	library.Module,

	// Content Modules
	movie.Module,
	tvshow.Module,

	// Job Workers
	moviejobs.Module,
	tvshowjobs.Module,

	// Integrations
	radarr.Module,
	sonarr.Module,

	// Metadata Service
	metadatafx.Module,

	// Observability (metrics, pprof)
	observability.Module,

	// HTTP API Server (ogen-generated)
	api.Module,
)
