package library

import "go.uber.org/fx"

// Module provides library service dependencies.
var Module = fx.Module("library",
	fx.Provide(NewService),
)
