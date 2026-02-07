package notification

import (
	"go.uber.org/fx"
)

// Module provides notification service dependencies.
var Module = fx.Module("notification",
	fx.Provide(
		NewDispatcher,
	),
)
