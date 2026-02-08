package apikeys

import (
	"time"

	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/fx"
)

// Module provides the API keys service
var Module = fx.Module("apikeys",
	fx.Provide(
		NewRepositoryPg,
		NewService,
		provideConfig,
	),
)

// provideConfig extracts API keys configuration
func provideConfig(cfg *config.Config) (int, time.Duration) {
	// For now, use defaults until we add config keys
	maxKeysPerUser := 10
	var defaultExpiry time.Duration = 0 // Never expire

	return maxKeysPerUser, defaultExpiry
}

// Params for the service constructor via fx
type Params struct {
	fx.In

	Queries        *db.Queries
	Logger         *slog.Logger
	MaxKeysPerUser int           `name:"apikeys_max_per_user"`
	DefaultExpiry  time.Duration `name:"apikeys_default_expiry"`
}

// Result for providing the service via fx
type Result struct {
	fx.Out

	Service *Service
}
