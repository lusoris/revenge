package health

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

// Module provides health check dependencies.
var Module = fx.Module("health",
	fx.Provide(NewService),
	fx.Invoke(registerHooks),
)

// registerHooks registers lifecycle hooks for health service.
func registerHooks(lc fx.Lifecycle, service *Service, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("health service started")
			// Mark startup as complete after all modules have started
			// This is done in OnStart to ensure all dependencies are ready
			service.MarkStartupComplete()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("health service stopped")
			return nil
		},
	})
}
