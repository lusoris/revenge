package whisparr

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides the Whisparr client to the application.
var Module = fx.Module("whisparr",
	fx.Provide(
		NewClientFromConfig,
	),
)

// NewClientFromConfig creates a Whisparr client from application config.
// Returns nil if Whisparr is not configured (optional dependency).
func NewClientFromConfig(cfg *config.Config, logger *slog.Logger) *Client {
	// Only create client if adult module is enabled and whisparr is configured
	if !cfg.Adult.Enabled || cfg.Adult.Whisparr.URL == "" || cfg.Adult.Whisparr.APIKey == "" {
		logger.Debug("whisparr client not configured, skipping")
		return nil
	}

	client, err := NewClient(ClientConfig{
		BaseURL:    cfg.Adult.Whisparr.URL,
		APIKey:     cfg.Adult.Whisparr.APIKey,
		Timeout:    30,
		RetryCount: 3,
	}, logger)

	if err != nil {
		logger.Warn("failed to create whisparr client", "error", err)
		return nil
	}

	return client
}
