package logging

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
)

// Module provides logging dependencies.
var Module = fx.Module("logging",
	fx.Provide(
		ProvideSlogLogger,
	),
)

// ProvideSlogLogger creates an slog.Logger from configuration.
func ProvideSlogLogger(cfg *config.Config) *slog.Logger {
	logConfig := Config{
		Level:       cfg.Logging.Level,
		Format:      cfg.Logging.Format,
		Development: cfg.Logging.Development,
		Output:      nil, // Use default (stdout)
	}

	logger := NewLogger(logConfig)

	// Set as default logger
	slog.SetDefault(logger)

	return logger
}
