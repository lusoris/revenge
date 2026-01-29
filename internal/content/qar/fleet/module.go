// Package fleet provides adult library domain models (QAR obfuscation: libraries → fleets).
package fleet

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	adultdb "github.com/lusoris/revenge/internal/content/qar/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// ModuleResult contains outputs from the fleet module.
type ModuleResult struct {
	fx.Out

	Repository      Repository
	Service         *Service
	LibraryService  *LibraryService
	LibraryProvider shared.LibraryProvider `group:"library_providers"`
}

// ProvideModule provides all fleet module dependencies.
func ProvideModule(pool *pgxpool.Pool, logger *slog.Logger) ModuleResult {
	queries := adultdb.New(pool)
	repo := NewSQLCRepository(pool, logger)
	service := NewService(repo, logger)
	libraryService := NewLibraryService(queries, logger)

	return ModuleResult{
		Repository:      repo,
		Service:         service,
		LibraryService:  libraryService,
		LibraryProvider: libraryService,
	}
}

// Module provides fleet (adult library) dependencies for fx.
var Module = fx.Module("qar.fleet",
	fx.Provide(ProvideModule),
)
