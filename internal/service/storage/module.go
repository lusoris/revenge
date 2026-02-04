package storage

import (
	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides the storage service for fx dependency injection.
var Module = fx.Options(
	fx.Provide(provideStorage),
)

// provideStorage creates the storage service using avatar config.
func provideStorage(cfg *config.Config, logger *zap.Logger) (Storage, error) {
	return NewLocalStorage(cfg.Avatar, logger)
}
