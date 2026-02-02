package api

import "go.uber.org/fx"

// Module provides the HTTP API server.
var Module = fx.Module(
	"api",
	fx.Provide(NewServer),
)
