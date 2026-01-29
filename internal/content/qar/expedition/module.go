// Package expedition provides adult movie domain models (QAR obfuscation: movies → expeditions).
package expedition

import (
	"go.uber.org/fx"
)

// Module provides expedition (adult movie) dependencies for fx.
var Module = fx.Module("qar.expedition",
	fx.Provide(
		NewSQLCRepository,
		NewService,
	),
)
