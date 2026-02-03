// Package search provides search functionality using Typesense.
package search

import (
	"go.uber.org/fx"
)

// Module provides search service dependencies.
var Module = fx.Module("search_service",
	fx.Provide(NewMovieSearchService),
)
