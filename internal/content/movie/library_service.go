package movie

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lusoris/revenge/internal/content/movie/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// LibraryService errors.
var (
	ErrLibraryNotFoundInService = errors.New("movie library not found")
	ErrLibraryAccessDenied      = errors.New("access denied to movie library")
)

// LibraryService provides movie library management implementing shared.LibraryProvider.
type LibraryService struct {
	queries *moviedb.Queries
	logger  *slog.Logger
}

// NewLibraryService creates a new movie library service.
func NewLibraryService(queries *moviedb.Queries, logger *slog.Logger) *LibraryService {
	return &LibraryService{
		queries: queries,
		logger:  logger.With("service", "movie_library"),
	}
}

// Ensure LibraryService implements LibraryProvider.
var _ shared.LibraryProvider = (*LibraryService)(nil)

// ModuleName returns the module identifier.
func (s *LibraryService) ModuleName() string {
	return "movie"
}

// ListLibraries returns all movie libraries accessible by the user.
func (s *LibraryService) ListLibraries(ctx context.Context, userID uuid.UUID) ([]shared.LibraryInfo, error) {
	libs, err := s.queries.ListAccessibleMovieLibraries(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]shared.LibraryInfo, len(libs))
	for i, lib := range libs {
		// Get movie count for this library
		count, err := s.queries.CountMoviesByMovieLibrary(ctx, lib.ID)
		if err != nil {
			s.logger.Warn("failed to count movies for library",
				"library_id", lib.ID,
				"error", err,
			)
			count = 0
		}

		result[i] = s.toLibraryInfo(lib, count)
	}

	return result, nil
}

// GetLibrary returns a specific library by ID.
func (s *LibraryService) GetLibrary(ctx context.Context, libraryID uuid.UUID) (*shared.LibraryInfo, error) {
	lib, err := s.queries.GetMovieLibraryByID(ctx, libraryID)
	if err != nil {
		return nil, ErrLibraryNotFoundInService
	}

	count, _ := s.queries.CountMoviesByMovieLibrary(ctx, libraryID)
	info := s.toLibraryInfo(lib, count)
	return &info, nil
}

// CreateLibrary creates a new movie library.
func (s *LibraryService) CreateLibrary(ctx context.Context, req shared.CreateLibraryRequest) (*shared.LibraryInfo, error) {
	params := moviedb.CreateMovieLibraryParams{
		Name:              req.Name,
		Paths:             req.Paths,
		ScanEnabled:       true,
		ScanIntervalHours: 24,
		TmdbEnabled:       true,
		ImdbEnabled:       true,
		DownloadTrailers:  false,
		DownloadBackdrops: true,
		DownloadNfo:       false,
		GenerateChapters:  false,
		IsPrivate:         false,
		SortOrder:         0,
	}

	// Apply module-specific settings if provided
	if settings, ok := req.Settings.(*MovieLibrarySettings); ok && settings != nil {
		if settings.PreferredLanguage != "" {
			params.PreferredLanguage = &settings.PreferredLanguage
		}
		params.TmdbEnabled = settings.TmdbEnabled
		params.ImdbEnabled = settings.ImdbEnabled
		params.DownloadTrailers = settings.DownloadTrailers
		params.DownloadBackdrops = settings.DownloadBackdrops
		params.DownloadNfo = settings.DownloadNfo
		params.GenerateChapters = settings.GenerateChapters
		params.IsPrivate = settings.IsPrivate
		if settings.OwnerUserID != uuid.Nil {
			params.OwnerUserID = pgtype.UUID{Bytes: settings.OwnerUserID, Valid: true}
		}
		params.SortOrder = int32(settings.SortOrder)
		if settings.Icon != "" {
			params.Icon = &settings.Icon
		}
	}

	lib, err := s.queries.CreateMovieLibrary(ctx, params)
	if err != nil {
		return nil, err
	}

	s.logger.Info("movie library created",
		"id", lib.ID,
		"name", lib.Name,
		"paths", lib.Paths,
	)

	info := s.toLibraryInfo(lib, 0)
	return &info, nil
}

// UpdateLibrary updates library settings.
func (s *LibraryService) UpdateLibrary(ctx context.Context, libraryID uuid.UUID, req shared.UpdateLibraryRequest) (*shared.LibraryInfo, error) {
	// Get current library
	current, err := s.queries.GetMovieLibraryByID(ctx, libraryID)
	if err != nil {
		return nil, ErrLibraryNotFoundInService
	}

	// Build update params
	params := moviedb.UpdateMovieLibraryParams{
		ID:                libraryID,
		Name:              current.Name,
		Paths:             current.Paths,
		ScanEnabled:       current.ScanEnabled,
		ScanIntervalHours: current.ScanIntervalHours,
		PreferredLanguage: current.PreferredLanguage,
		TmdbEnabled:       current.TmdbEnabled,
		ImdbEnabled:       current.ImdbEnabled,
		DownloadTrailers:  current.DownloadTrailers,
		DownloadBackdrops: current.DownloadBackdrops,
		DownloadNfo:       current.DownloadNfo,
		GenerateChapters:  current.GenerateChapters,
		IsPrivate:         current.IsPrivate,
		OwnerUserID:       current.OwnerUserID,
		SortOrder:         current.SortOrder,
		Icon:              current.Icon,
	}

	// Apply updates
	if req.Name != nil {
		params.Name = *req.Name
	}
	if len(req.Paths) > 0 {
		params.Paths = req.Paths
	}

	// Apply module-specific settings if provided
	if settings, ok := req.Settings.(*MovieLibrarySettings); ok && settings != nil {
		if settings.PreferredLanguage != "" {
			params.PreferredLanguage = &settings.PreferredLanguage
		}
		params.TmdbEnabled = settings.TmdbEnabled
		params.ImdbEnabled = settings.ImdbEnabled
		params.DownloadTrailers = settings.DownloadTrailers
		params.DownloadBackdrops = settings.DownloadBackdrops
		params.DownloadNfo = settings.DownloadNfo
		params.GenerateChapters = settings.GenerateChapters
		params.IsPrivate = settings.IsPrivate
		if settings.OwnerUserID != uuid.Nil {
			params.OwnerUserID = pgtype.UUID{Bytes: settings.OwnerUserID, Valid: true}
		}
		params.SortOrder = int32(settings.SortOrder)
		if settings.Icon != "" {
			params.Icon = &settings.Icon
		}
	}

	lib, err := s.queries.UpdateMovieLibrary(ctx, params)
	if err != nil {
		return nil, err
	}

	s.logger.Info("movie library updated",
		"id", lib.ID,
		"name", lib.Name,
	)

	count, _ := s.queries.CountMoviesByMovieLibrary(ctx, libraryID)
	info := s.toLibraryInfo(lib, count)
	return &info, nil
}

// DeleteLibrary removes a library and optionally its content.
func (s *LibraryService) DeleteLibrary(ctx context.Context, libraryID uuid.UUID, deleteContent bool) error {
	// Verify library exists
	lib, err := s.queries.GetMovieLibraryByID(ctx, libraryID)
	if err != nil {
		return ErrLibraryNotFoundInService
	}

	// Note: If deleteContent is true, we would delete all movies in this library
	// This should be handled via a cascade or explicit delete job
	// For now, we'll just delete the library and let FK cascade handle movies

	if err := s.queries.DeleteMovieLibrary(ctx, libraryID); err != nil {
		return err
	}

	s.logger.Info("movie library deleted",
		"id", lib.ID,
		"name", lib.Name,
		"delete_content", deleteContent,
	)

	return nil
}

// ScanLibrary triggers a library scan.
func (s *LibraryService) ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) error {
	// Verify library exists
	lib, err := s.queries.GetMovieLibraryByID(ctx, libraryID)
	if err != nil {
		return ErrLibraryNotFoundInService
	}

	s.logger.Info("library scan triggered",
		"library_id", libraryID,
		"library_name", lib.Name,
		"full_scan", fullScan,
	)

	// TODO: Enqueue a River job for the actual scanning
	// This will be implemented when we add the scan job worker
	// For now, just log and return

	return nil
}

// toLibraryInfo converts a database MovieLibrary to shared.LibraryInfo.
func (s *LibraryService) toLibraryInfo(lib moviedb.MovieLibrary, itemCount int64) shared.LibraryInfo {
	var ownerID uuid.UUID
	if lib.OwnerUserID.Valid {
		ownerID = lib.OwnerUserID.Bytes
	}

	return shared.LibraryInfo{
		ID:        lib.ID,
		Module:    "movie",
		Name:      lib.Name,
		Paths:     lib.Paths,
		IsAdult:   false, // Movie libraries are never adult content
		ItemCount: itemCount,
		Settings: &MovieLibrarySettings{
			PreferredLanguage: derefString(lib.PreferredLanguage),
			TmdbEnabled:       lib.TmdbEnabled,
			ImdbEnabled:       lib.ImdbEnabled,
			DownloadTrailers:  lib.DownloadTrailers,
			DownloadBackdrops: lib.DownloadBackdrops,
			DownloadNfo:       lib.DownloadNfo,
			GenerateChapters:  lib.GenerateChapters,
			IsPrivate:         lib.IsPrivate,
			OwnerUserID:       ownerID,
			SortOrder:         int(lib.SortOrder),
			Icon:              derefString(lib.Icon),
		},
	}
}

// MovieLibrarySettings contains movie-specific library settings.
type MovieLibrarySettings struct {
	PreferredLanguage string    `json:"preferred_language,omitempty"`
	TmdbEnabled       bool      `json:"tmdb_enabled"`
	ImdbEnabled       bool      `json:"imdb_enabled"`
	DownloadTrailers  bool      `json:"download_trailers"`
	DownloadBackdrops bool      `json:"download_backdrops"`
	DownloadNfo       bool      `json:"download_nfo"`
	GenerateChapters  bool      `json:"generate_chapters"`
	IsPrivate         bool      `json:"is_private"`
	OwnerUserID       uuid.UUID `json:"owner_user_id,omitempty"`
	SortOrder         int       `json:"sort_order"`
	Icon              string    `json:"icon,omitempty"`
}

// derefString safely dereferences a string pointer.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
