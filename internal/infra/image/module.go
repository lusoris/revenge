package image

import (
	"go.uber.org/fx"
	"log/slog"
)

// Module provides image service dependencies.
var Module = fx.Module("image",
	fx.Provide(
		NewImageService,
	),
)

// ImageServiceParams holds dependencies for creating the image service.
type ImageServiceParams struct {
	fx.In

	Config Config `name:"imageConfig"`
	Logger *slog.Logger
}

// NewImageService creates a new image service with FX dependencies.
func NewImageService(p ImageServiceParams) (*Service, error) {
	return NewService(p.Config, p.Logger)
}
