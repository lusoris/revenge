package email

import (
	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides the email service for fx dependency injection.
var Module = fx.Options(
	fx.Provide(provideService),
)

// provideService creates the email service.
func provideService(cfg *config.Config, logger *zap.Logger) *Service {
	return NewService(cfg.Email, logger)
}
