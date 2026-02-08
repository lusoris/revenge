package email

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
)

// Module provides the email service for fx dependency injection.
var Module = fx.Options(
	fx.Provide(provideService),
)

// provideService creates the email service.
func provideService(cfg *config.Config, logger *slog.Logger) *Service {
	return NewService(cfg.Email, logger)
}
