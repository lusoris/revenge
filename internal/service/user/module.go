package user

import "go.uber.org/fx"

// Module provides user service dependencies.
var Module = fx.Module("user",
	fx.Provide(NewService),
)
