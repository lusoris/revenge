package observability

import "go.uber.org/fx"

// Module provides observability components for the application.
var Module = fx.Options(
	fx.Invoke(NewServer),
	fx.Invoke(StartCollector),
)
