package auth

import "go.uber.org/fx"

// Module provides auth service dependencies.
var Module = fx.Module("auth",
	fx.Provide(NewService),
)
