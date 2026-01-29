// Package voyage provides adult scene domain models (QAR obfuscation: scenes → voyages).
package voyage

import (
	"go.uber.org/fx"
)

// Module provides voyage (adult scene) dependencies for fx.
var Module = fx.Module("qar.voyage",
	fx.Provide(
		NewSQLCRepository,
		NewService,
	),
)
