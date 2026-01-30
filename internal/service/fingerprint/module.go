// Package fingerprint provides video fingerprinting services.
package fingerprint

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides the fingerprint service to the application.
var Module = fx.Module("fingerprint",
	fx.Provide(
		NewServiceFromConfig,
	),
)

// NewServiceFromConfig creates a fingerprint service from application config.
func NewServiceFromConfig(cfg *config.Config, logger *slog.Logger) *Service {
	fpCfg := DefaultConfig()

	// Override from config if available
	if cfg.Adult.Fingerprint.FFProbePath != "" {
		fpCfg.FFProbePath = cfg.Adult.Fingerprint.FFProbePath
	}
	if cfg.Adult.Fingerprint.FFMpegPath != "" {
		fpCfg.FFMpegPath = cfg.Adult.Fingerprint.FFMpegPath
	}
	fpCfg.GeneratePHash = cfg.Adult.Fingerprint.GeneratePHash
	fpCfg.GenerateMD5 = cfg.Adult.Fingerprint.GenerateMD5
	if cfg.Adult.Fingerprint.Timeout > 0 {
		fpCfg.Timeout = cfg.Adult.Fingerprint.Timeout
	}

	return NewService(fpCfg, logger)
}
