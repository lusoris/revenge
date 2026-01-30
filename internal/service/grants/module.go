package grants

import "go.uber.org/fx"

// Module provides grants service dependencies.
var Module = fx.Module("grants",
	fx.Provide(NewService),
)
