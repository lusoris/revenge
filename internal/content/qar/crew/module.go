// Package crew provides adult performer domain models (QAR obfuscation: performers → crew).
package crew

import (
	"go.uber.org/fx"
)

// Module provides crew (performer) dependencies for fx.
var Module = fx.Module("qar.crew",
	fx.Provide(
		NewSQLCRepository,
		NewService,
	),
)
