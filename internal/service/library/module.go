package library

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/content/shared"
)

// ModuleParams contains dependencies for the library aggregator module.
type ModuleParams struct {
	fx.In

	Providers []shared.LibraryProvider `group:"library_providers"`
	Logger    *slog.Logger
}

// ProvideService creates the library aggregator service.
func ProvideService(p ModuleParams) *Service {
	return NewService(p.Providers, p.Logger)
}

// Module provides the library aggregator service for fx.
var Module = fx.Module("library",
	fx.Provide(ProvideService),
)
