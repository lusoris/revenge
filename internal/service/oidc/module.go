package oidc

import "go.uber.org/fx"

// Module provides OIDC service dependencies.
var Module = fx.Module("oidc",
	fx.Provide(NewService),
)
