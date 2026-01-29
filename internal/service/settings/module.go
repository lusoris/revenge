package settings

import "go.uber.org/fx"

// Module provides server settings service dependencies.
var Module = fx.Module("settings",
	fx.Provide(NewService),
)
