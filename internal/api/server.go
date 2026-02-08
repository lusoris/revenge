// Package api provides the HTTP API server implementation using ogen-generated code.
package api

import (
	"context"
	"fmt"
	"net/http"

	openapidoc "github.com/lusoris/revenge/api/openapi"
	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/api/sse"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/image"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/observability"
	"github.com/lusoris/revenge/internal/integration/radarr"
	"github.com/lusoris/revenge/internal/integration/sonarr"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/hls"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/service/mfa"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/search"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"go.uber.org/fx"
	"log/slog"
)

// scalarHTML is the Scalar API reference UI, loaded from the embedded OpenAPI spec.
const scalarHTML = `<!doctype html>
<html>
<head>
  <title>Revenge API Reference</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
</head>
<body>
  <script id="api-reference" data-url="/api/openapi.yaml"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

// Server wraps the ogen-generated HTTP server with lifecycle management.
type Server struct {
	httpServer         *http.Server
	ogenServer         *ogen.Server
	logger             *slog.Logger
	authLimiter        *middleware.RateLimiter
	globalLimiter      *middleware.RateLimiter
	redisAuthLimiter   *middleware.RedisRateLimiter
	redisGlobalLimiter *middleware.RedisRateLimiter
}
// ServerParams defines the dependencies required to create the API server.
type ServerParams struct {
	fx.In

	Config          *config.Config
	Logger          *slog.Logger
	HealthService   *health.Service
	SettingsService settings.Service
	UserService     *user.Service
	AuthService     *auth.Service
	SessionService  *session.Service
	RBACService     *rbac.Service
	APIKeyService   *apikeys.Service
	OIDCService     *oidc.Service
	ActivityService *activity.Service
	LibraryService  *library.Service
	SearchService       *search.MovieSearchService   `optional:"true"`
	TVShowSearchService *search.TVShowSearchService `optional:"true"`
	TokenManager    auth.TokenManager
	// Cache client for Redis-based rate limiting
	CacheClient *cache.Client `optional:"true"`
	// MFA services
	TOTPService        *mfa.TOTPService
	BackupCodesService *mfa.BackupCodesService
	MFAManager         *mfa.MFAManager
	WebAuthnService    *mfa.WebAuthnService `optional:"true"`
	// Content modules
	MovieHandler    *movie.Handler
	MetadataService metadata.Service `optional:"true"`
	ImageService    *image.Service
	TVShowService   tvshow.Service         `optional:"true"`
	// Playback / HLS streaming (optional)
	PlaybackService *playback.Service      `optional:"true"`
	StreamHandler   *hls.StreamHandler     `optional:"true"`
	// SSE real-time events (optional)
	SSEHandler      *sse.Handler           `optional:"true"`
	// Integration services (optional)
	RadarrService *radarr.SyncService `optional:"true"`
	SonarrService *sonarr.SyncService `optional:"true"`
	RiverClient   *jobs.Client        `optional:"true"`
	Lifecycle     fx.Lifecycle
}

// NewServer creates a new HTTP API server with ogen-generated handlers.
func NewServer(p ServerParams) (*Server, error) {
	// Create MFA handler
	mfaHandler := NewMFAHandler(
		p.TOTPService,
		p.BackupCodesService,
		p.MFAManager,
		p.WebAuthnService,
		p.Logger.With("component", "mfa"),
	)

	// Create the handler implementation
	handler := &Handler{
		logger:          p.Logger.With("component", "api"),
		cfg:             p.Config,
		healthService:   p.HealthService,
		userService:     p.UserService,
		settingsService: p.SettingsService,
		authService:     p.AuthService,
		sessionService:  p.SessionService,
		rbacService:     p.RBACService,
		apikeyService:   p.APIKeyService,
		oidcService:     p.OIDCService,
		activityService: p.ActivityService,
		libraryService:  p.LibraryService,
		searchService:       p.SearchService,
		tvshowSearchService: p.TVShowSearchService,
		tokenManager:    p.TokenManager,
		mfaHandler:      mfaHandler,
		movieHandler:    p.MovieHandler,
		metadataService: p.MetadataService,
		imageService:    p.ImageService,
		tvshowService:   p.TVShowService,
	}

	// Wire up optional playback service
	if p.PlaybackService != nil {
		handler.playbackService = p.PlaybackService
	}
	// Wire up optional Radarr integration
	if p.RadarrService != nil {
		handler.radarrService = p.RadarrService
	}
	// Wire up optional Sonarr integration
	if p.SonarrService != nil {
		handler.sonarrService = p.SonarrService
	}
	if p.RiverClient != nil {
		handler.riverClient = p.RiverClient
	}

	// Create ogen server, optionally with rate limiting middleware
	var ogenServer *ogen.Server
	var authLimiter, globalLimiter *middleware.RateLimiter
	var redisAuthLimiter, redisGlobalLimiter *middleware.RedisRateLimiter
	var err error

	if p.Config.Server.RateLimit.Enabled {
		// Determine if we should use Redis backend
		useRedis := p.Config.Server.RateLimit.Backend == "redis" && p.CacheClient != nil && p.CacheClient.RueidisClient() != nil

		if useRedis {
			// Create Redis-based rate limiters
			redisAuthConfig := middleware.AuthRedisRateLimiterConfig()
			redisGlobalConfig := middleware.DefaultRedisRateLimiterConfig()

			// Override with custom config if provided
			if p.Config.Server.RateLimit.Auth.RequestsPerSecond > 0 {
				redisAuthConfig.RequestsPerSecond = p.Config.Server.RateLimit.Auth.RequestsPerSecond
			}
			if p.Config.Server.RateLimit.Auth.Burst > 0 {
				redisAuthConfig.Burst = p.Config.Server.RateLimit.Auth.Burst
			}
			if p.Config.Server.RateLimit.Global.RequestsPerSecond > 0 {
				redisGlobalConfig.RequestsPerSecond = p.Config.Server.RateLimit.Global.RequestsPerSecond
			}
			if p.Config.Server.RateLimit.Global.Burst > 0 {
				redisGlobalConfig.Burst = p.Config.Server.RateLimit.Global.Burst
			}

			rueidisClient := p.CacheClient.RueidisClient()
			redisAuthLimiter = middleware.NewRedisRateLimiter(redisAuthConfig, rueidisClient, p.Logger)
			redisGlobalLimiter = middleware.NewRedisRateLimiter(redisGlobalConfig, rueidisClient, p.Logger)

			p.Logger.Info("Using Redis-based rate limiting",
				slog.Float64("auth.rps", redisAuthConfig.RequestsPerSecond),
				slog.Int("auth.burst", redisAuthConfig.Burst),
				slog.Float64("global.rps", redisGlobalConfig.RequestsPerSecond),
				slog.Int("global.burst", redisGlobalConfig.Burst),
			)

			// Create ogen server with Redis rate limiting middleware
			ogenServer, err = ogen.NewServer(
				handler,
				handler,
				ogen.WithMiddleware(
					middleware.RequestIDMiddleware(),
					middleware.RequestMetadataMiddleware(),
					observability.HTTPMetricsMiddleware(),
					redisAuthLimiter.Middleware(),
					redisGlobalLimiter.Middleware(),
				),
				ogen.WithErrorHandler(middleware.ErrorHandler),
			)
		} else {
			// Create in-memory rate limiters
			authConfig := middleware.AuthRateLimitConfig()
			globalConfig := middleware.DefaultRateLimitConfig()

			// Override with custom config if provided
			if p.Config.Server.RateLimit.Auth.RequestsPerSecond > 0 {
				authConfig.RequestsPerSecond = p.Config.Server.RateLimit.Auth.RequestsPerSecond
			}
			if p.Config.Server.RateLimit.Auth.Burst > 0 {
				authConfig.Burst = p.Config.Server.RateLimit.Auth.Burst
			}
			if p.Config.Server.RateLimit.Global.RequestsPerSecond > 0 {
				globalConfig.RequestsPerSecond = p.Config.Server.RateLimit.Global.RequestsPerSecond
			}
			if p.Config.Server.RateLimit.Global.Burst > 0 {
				globalConfig.Burst = p.Config.Server.RateLimit.Global.Burst
			}

			// Auth rate limiter: strict limits for login/MFA endpoints
			authLimiter = middleware.NewRateLimiter(authConfig, p.Logger)
			// Global rate limiter: generous limits for all endpoints
			globalLimiter = middleware.NewRateLimiter(globalConfig, p.Logger)

			p.Logger.Info("Using in-memory rate limiting",
				slog.Float64("auth.rps", authConfig.RequestsPerSecond),
				slog.Int("auth.burst", authConfig.Burst),
				slog.Float64("global.rps", globalConfig.RequestsPerSecond),
				slog.Int("global.burst", globalConfig.Burst),
			)

			// Create ogen server with in-memory rate limiting middleware
			ogenServer, err = ogen.NewServer(
				handler,
				handler,
				ogen.WithMiddleware(
					middleware.RequestIDMiddleware(),
					middleware.RequestMetadataMiddleware(),
					observability.HTTPMetricsMiddleware(),
					authLimiter.Middleware(),
					globalLimiter.Middleware(),
				),
				ogen.WithErrorHandler(middleware.ErrorHandler),
			)
		}
	} else {
		// Create ogen server without rate limiting
		ogenServer, err = ogen.NewServer(
			handler,
			handler,
			ogen.WithMiddleware(
				middleware.RequestIDMiddleware(),
				middleware.RequestMetadataMiddleware(),
				observability.HTTPMetricsMiddleware(),
			),
			ogen.WithErrorHandler(middleware.ErrorHandler),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create ogen server: %w", err)
	}

	// Create HTTP server with RequestID wrapper
	// The wrapper adds X-Request-ID to response headers
	rootHandler := http.Handler(middleware.RequestIDHTTPWrapper(ogenServer))

	// Serve OpenAPI spec for frontend dev tools
	specHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Write(openapidoc.Spec) //nolint:errcheck
	})

	// Serve interactive API documentation (Scalar)
	docsHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(scalarHTML)) //nolint:errcheck
	})

	// Register non-ogen routes on a mux (spec, docs, HLS streaming)
	mux := http.NewServeMux()
	mux.Handle("GET /api/openapi.yaml", specHandler)
	mux.Handle("GET /api/docs", docsHandler)
	if p.StreamHandler != nil {
		mux.Handle("/api/v1/playback/stream/", p.StreamHandler)
	}
	if p.SSEHandler != nil {
		mux.Handle("GET /api/v1/events", p.SSEHandler)
	}
	mux.Handle("/", rootHandler)
	rootHandler = mux

	// Wrap with CORS middleware (outermost layer so all responses get CORS headers,
	// including preflight OPTIONS, error responses, and HLS endpoints).
	rootHandler = middleware.CORSMiddleware(p.Config.Server.CORS)(rootHandler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port),
		Handler:      rootHandler,
		ReadTimeout:  p.Config.Server.ReadTimeout,
		WriteTimeout: p.Config.Server.WriteTimeout,
		IdleTimeout:  p.Config.Server.IdleTimeout,
	}

	server := &Server{
		httpServer:         httpServer,
		ogenServer:         ogenServer,
		logger:             p.Logger.With("component", "server"),
		authLimiter:        authLimiter,
		globalLimiter:      globalLimiter,
		redisAuthLimiter:   redisAuthLimiter,
		redisGlobalLimiter: redisGlobalLimiter,
	}

	// Register lifecycle hooks
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start HTTP server in background
			go func() {
				addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port)
				server.logger.Info("Starting HTTP server",
					slog.String("address", addr),
					slog.Bool("ratelimit.enabled", p.Config.Server.RateLimit.Enabled),
					slog.String("ratelimit.backend", p.Config.Server.RateLimit.Backend),
				)
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					server.logger.Error("HTTP server error", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.logger.Info("Stopping HTTP server")
			// Stop rate limiters if they were created
			if server.authLimiter != nil {
				server.authLimiter.Stop()
			}
			if server.globalLimiter != nil {
				server.globalLimiter.Stop()
			}
			if server.redisAuthLimiter != nil {
				server.redisAuthLimiter.Stop()
			}
			if server.redisGlobalLimiter != nil {
				server.redisGlobalLimiter.Stop()
			}
			return httpServer.Shutdown(ctx)
		},
	})

	return server, nil
}
