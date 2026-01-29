package tmdb

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/pkg/lazy"
)

// ModuleParams contains dependencies for the TMDb module.
type ModuleParams struct {
	fx.In

	Config Config `optional:"true"`
	Logger *slog.Logger
}

// NewLazyClient creates a lazily-initialized TMDb client.
// The client is only created when first needed, following the offloading pattern.
func NewLazyClient(p ModuleParams) *lazy.Service[*Client] {
	return lazy.New(func() (*Client, error) {
		p.Logger.Info("initializing TMDb client (lazy)")
		return NewClient(p.Config, p.Logger), nil
	})
}

// NewLazyProvider creates a lazily-initialized TMDb provider.
func NewLazyProvider(lazyClient *lazy.Service[*Client], logger *slog.Logger) *lazy.Service[*Provider] {
	return lazy.New(func() (*Provider, error) {
		client, err := lazyClient.Get()
		if err != nil {
			return nil, err
		}
		return NewProvider(client, logger), nil
	})
}

// Module provides TMDb client and provider for dependency injection.
// Note: This module provides lazy.Service wrappers. The parent movie module
// is responsible for converting these to the MetadataProvider interface.
var Module = fx.Module("tmdb",
	fx.Provide(
		NewLazyClient,
		NewLazyProvider,
	),
)
