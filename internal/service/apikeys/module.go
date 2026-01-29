package apikeys

import "go.uber.org/fx"

// Module provides API keys service dependencies.
var Module = fx.Module("apikeys",
	fx.Provide(NewService),
)
