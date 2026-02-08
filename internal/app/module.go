// Package app provides the main application module that wires all dependencies together.
package app

import (
	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/api/sse"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	appcrypto "github.com/lusoris/revenge/internal/crypto"
	"github.com/lusoris/revenge/internal/content/movie/moviejobs"
	"github.com/lusoris/revenge/internal/content/tvshow"
	tvshowjobs "github.com/lusoris/revenge/internal/content/tvshow/jobs"
	"github.com/lusoris/revenge/internal/playback/playbackfx"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/image"
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
	"github.com/lusoris/revenge/internal/service/email"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/metadata/metadatafx"
	"github.com/lusoris/revenge/internal/service/mfa"
	"github.com/lusoris/revenge/internal/service/notification"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	searchsvc "github.com/lusoris/revenge/internal/service/search"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/storage"
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
	image.Module,
	appcrypto.Module,

	// Services
	settings.Module,
	user.Module,
	auth.Module,
	email.Module,
	session.Module,
	rbac.Module,
	apikeys.Module,
	mfa.Module,
	oidc.Module,
	activity.Module,
	notification.Module,
	storage.Module,
	library.Module,
	searchsvc.Module,

	// Content Modules
	movie.Module,
	tvshow.Module,

	// Playback / HLS Streaming
	playbackfx.Module,

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

	// SSE Real-Time Events
	sse.Module,

	// HTTP API Server (ogen-generated)
	api.Module,
)
