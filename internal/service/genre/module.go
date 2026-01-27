// Package genre provides fx module for genre service.
package genre

import (
	"go.uber.org/fx"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// Module provides the genre service dependencies.
var Module = fx.Module("genre",
	fx.Provide(
		NewService,
		func(s *Service) domain.GenreService { return s },
	),
)
