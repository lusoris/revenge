package activity

import "go.uber.org/fx"

// Module provides activity logging service dependencies.
var Module = fx.Module("activity",
	fx.Provide(NewService),
)
