package notification

import (
	"go.uber.org/fx"
)

// Module provides notification service dependencies.
var Module = fx.Module("notification",
	fx.Provide(
		NewDispatcher,
		// Expose *Dispatcher as notification.Service for consumers that
		// depend on the interface (e.g. SSE module).
		func(d *Dispatcher) Service { return d },
	),
)
