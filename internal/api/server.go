// Package api provides the HTTP API server implementation using ogen-generated code.
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/image"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/integration/radarr"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/mfa"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/search"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server wraps the ogen-generated HTTP server with lifecycle management.
type Server struct {
	httpServer    *http.Server
	ogenServer    *ogen.Server
	logger        *zap.Logger
	authLimiter   *middleware.RateLimiter
	globalLimiter *middleware.RateLimiter
}
// ServerParams defines the dependencies required to create the API server.
type ServerParams struct {
	fx.In

	Config          *config.Config
	Logger          *zap.Logger
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
	SearchService   *search.MovieSearchService `optional:"true"`
	TokenManager    auth.TokenManager
	// MFA services
	TOTPService        *mfa.TOTPService
	BackupCodesService *mfa.BackupCodesService
	MFAManager         *mfa.MFAManager
	// Content modules
	MovieHandler    *movie.Handler
	MetadataService *movie.MetadataService `optional:"true"`
	ImageService    *image.Service         `optional:"true"`
	// Integration services (optional)
	RadarrService *radarr.SyncService `optional:"true"`
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
		p.Logger.Named("mfa"),
	)

	// Create the handler implementation
	handler := &Handler{
		logger:          p.Logger.Named("api"),
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
		searchService:   p.SearchService,
		tokenManager:    p.TokenManager,
		mfaHandler:      mfaHandler,
		movieHandler:    p.MovieHandler,
		metadataService: p.MetadataService,
		imageService:    p.ImageService,
	}

	// Wire up optional Radarr integration
	if p.RadarrService != nil {
		handler.radarrService = p.RadarrService
	}
	if p.RiverClient != nil {
		handler.riverClient = p.RiverClient
	}

	// Create ogen server, optionally with rate limiting middleware
	var ogenServer *ogen.Server
	var authLimiter, globalLimiter *middleware.RateLimiter
	var err error

	if p.Config.Server.RateLimit.Enabled {
		// Create rate limiters from config (or use defaults)
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

		// Create ogen server with rate limiting middleware
		ogenServer, err = ogen.NewServer(
			handler,
			handler,
			ogen.WithMiddleware(
				authLimiter.Middleware(),
				globalLimiter.Middleware(),
			),
			ogen.WithErrorHandler(middleware.ErrorHandler),
		)
	} else {
		// Create ogen server without rate limiting
		ogenServer, err = ogen.NewServer(
			handler,
			handler,
			ogen.WithErrorHandler(middleware.ErrorHandler),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create ogen server: %w", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port),
		Handler:      ogenServer,
		ReadTimeout:  p.Config.Server.ReadTimeout,
		WriteTimeout: p.Config.Server.WriteTimeout,
		IdleTimeout:  p.Config.Server.IdleTimeout,
	}

	server := &Server{
		httpServer:    httpServer,
		ogenServer:    ogenServer,
		logger:        p.Logger.Named("server"),
		authLimiter:   authLimiter,
		globalLimiter: globalLimiter,
	}

	// Register lifecycle hooks
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start HTTP server in background
			go func() {
				addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port)
				if p.Config.Server.RateLimit.Enabled {
					server.logger.Info("Starting HTTP server",
						zap.String("address", addr),
						zap.String("ratelimit.auth", "1 req/s, burst 5"),
						zap.String("ratelimit.global", "10 req/s, burst 20"),
					)
				} else {
					server.logger.Info("Starting HTTP server",
						zap.String("address", addr),
						zap.Bool("ratelimit.enabled", false),
					)
				}
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					server.logger.Error("HTTP server error", zap.Error(err))
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
			return httpServer.Shutdown(ctx)
		},
	})

	return server, nil
}
