// Package api provides the HTTP API server implementation using ogen-generated code.
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/mfa"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
	"github.com/lusoris/revenge/internal/service/user"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server wraps the ogen-generated HTTP server with lifecycle management.
type Server struct {
	httpServer *http.Server
	ogenServer *ogen.Server
	logger     *zap.Logger
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
	TokenManager    auth.TokenManager
	// MFA services
	TOTPService        *mfa.TOTPService
	BackupCodesService *mfa.BackupCodesService
	MFAManager         *mfa.MFAManager
	Lifecycle          fx.Lifecycle
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
		tokenManager:    p.TokenManager,
		mfaHandler:      mfaHandler,
	}

	// Create ogen server
	ogenServer, err := ogen.NewServer(handler, handler)
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
		httpServer: httpServer,
		ogenServer: ogenServer,
		logger:     p.Logger.Named("server"),
	}

	// Register lifecycle hooks
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start HTTP server in background
			go func() {
				addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port)
				server.logger.Info("Starting HTTP server",
					zap.String("address", addr),
				)
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					server.logger.Error("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.logger.Info("Stopping HTTP server")
			return httpServer.Shutdown(ctx)
		},
	})

	return server, nil
}
