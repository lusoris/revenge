package movie

import (
	"go.uber.org/fx"
)

// Module provides the movie content module
var Module = fx.Module("movie",
	fx.Provide(
		NewPostgresRepository,
		NewService,
		NewHandler,
	),
)
