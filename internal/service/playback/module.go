package playback

import "go.uber.org/fx"

// Module provides playback service dependencies.
var Module = fx.Module("playback",
	fx.Provide(NewService),
)
