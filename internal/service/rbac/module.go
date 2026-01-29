package rbac

import "go.uber.org/fx"

// Module provides RBAC service dependencies.
var Module = fx.Module("rbac",
	fx.Provide(NewService),
)
