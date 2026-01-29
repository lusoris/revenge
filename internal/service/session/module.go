package session

import "go.uber.org/fx"

// Module provides session service dependencies.
var Module = fx.Module("session",
	fx.Provide(NewService),
)
