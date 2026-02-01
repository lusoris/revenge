package logging

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides logging dependencies.
var Module = fx.Module("logging",
	fx.Provide(
		ProvideSlogLogger,
		ProvideZapLogger,
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

// ProvideZapLogger creates a zap.Logger from configuration.
func ProvideZapLogger(cfg *config.Config) (*zap.Logger, error) {
	logConfig := Config{
		Level:       cfg.Logging.Level,
		Format:      cfg.Logging.Format,
		Development: cfg.Logging.Development,
		Output:      nil, // Use default (stdout)
	}

	logger, err := NewZapLogger(logConfig)
	if err != nil {
		return nil, err
	}

	// Set as global logger
	zap.ReplaceGlobals(logger)

	return logger, nil
}
