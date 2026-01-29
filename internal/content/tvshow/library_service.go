package tvshow

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	tvshowdb "github.com/lusoris/revenge/internal/content/tvshow/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// LibraryService errors.
var (
	ErrLibraryNotFoundInService = errors.New("tv library not found")
	ErrLibraryAccessDenied      = errors.New("access denied to tv library")
)

// LibraryService provides TV library management implementing shared.LibraryProvider.
type LibraryService struct {
	queries *tvshowdb.Queries
	logger  *slog.Logger
}

// NewLibraryService creates a new TV library service.
func NewLibraryService(queries *tvshowdb.Queries, logger *slog.Logger) *LibraryService {
	return &LibraryService{
		queries: queries,
		logger:  logger.With("service", "tv_library"),
	}
}

// Ensure LibraryService implements LibraryProvider.
var _ shared.LibraryProvider = (*LibraryService)(nil)

// ModuleName returns the module identifier.
func (s *LibraryService) ModuleName() string {
	return "tvshow"
}

// ListLibraries returns all TV libraries accessible by the user.
func (s *LibraryService) ListLibraries(ctx context.Context, userID uuid.UUID) ([]shared.LibraryInfo, error) {
	libs, err := s.queries.ListAccessibleTVLibraries(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]shared.LibraryInfo, len(libs))
	for i, lib := range libs {
		// Get series count for this library
		count, err := s.queries.CountSeriesByTVLibrary(ctx, lib.ID)
		if err != nil {
			s.logger.Warn("failed to count series for library",
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
	lib, err := s.queries.GetTVLibraryByID(ctx, libraryID)
	if err != nil {
		return nil, ErrLibraryNotFoundInService
	}

	count, _ := s.queries.CountSeriesByTVLibrary(ctx, libraryID)
	info := s.toLibraryInfo(lib, count)
	return &info, nil
}

// CreateLibrary creates a new TV library.
func (s *LibraryService) CreateLibrary(ctx context.Context, req shared.CreateLibraryRequest) (*shared.LibraryInfo, error) {
	params := tvshowdb.CreateTVLibraryParams{
		Name:              req.Name,
		Paths:             req.Paths,
		ScanEnabled:       true,
		ScanIntervalHours: 24,
		TmdbEnabled:       true,
		TvdbEnabled:       true,
		DownloadBackdrops: true,
		DownloadNfo:       false,
		GenerateChapters:  false,
		AutoAddMissing:    false,
		IsPrivate:         false,
		SortOrder:         0,
	}

	// Apply module-specific settings if provided
	if settings, ok := req.Settings.(*TVLibrarySettings); ok && settings != nil {
		if settings.PreferredLanguage != "" {
			params.PreferredLanguage = &settings.PreferredLanguage
		}
		params.TmdbEnabled = settings.TmdbEnabled
		params.TvdbEnabled = settings.TvdbEnabled
		params.DownloadBackdrops = settings.DownloadBackdrops
		params.DownloadNfo = settings.DownloadNfo
		params.GenerateChapters = settings.GenerateChapters
		if settings.SeasonFolderFormat != "" {
			params.SeasonFolderFormat = &settings.SeasonFolderFormat
		}
		if settings.EpisodeNamingFormat != "" {
			params.EpisodeNamingFormat = &settings.EpisodeNamingFormat
		}
		params.AutoAddMissing = settings.AutoAddMissing
		params.IsPrivate = settings.IsPrivate
		if settings.OwnerUserID != uuid.Nil {
			params.OwnerUserID = pgtype.UUID{Bytes: settings.OwnerUserID, Valid: true}
		}
		params.SortOrder = int32(settings.SortOrder)
		if settings.Icon != "" {
			params.Icon = &settings.Icon
		}
	}

	lib, err := s.queries.CreateTVLibrary(ctx, params)
	if err != nil {
		return nil, err
	}

	s.logger.Info("tv library created",
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
	current, err := s.queries.GetTVLibraryByID(ctx, libraryID)
	if err != nil {
		return nil, ErrLibraryNotFoundInService
	}

	// Build update params
	params := tvshowdb.UpdateTVLibraryParams{
		ID:                  libraryID,
		Name:                current.Name,
		Paths:               current.Paths,
		ScanEnabled:         current.ScanEnabled,
		ScanIntervalHours:   current.ScanIntervalHours,
		PreferredLanguage:   current.PreferredLanguage,
		TmdbEnabled:         current.TmdbEnabled,
		TvdbEnabled:         current.TvdbEnabled,
		DownloadBackdrops:   current.DownloadBackdrops,
		DownloadNfo:         current.DownloadNfo,
		GenerateChapters:    current.GenerateChapters,
		SeasonFolderFormat:  current.SeasonFolderFormat,
		EpisodeNamingFormat: current.EpisodeNamingFormat,
		AutoAddMissing:      current.AutoAddMissing,
		IsPrivate:           current.IsPrivate,
		OwnerUserID:         current.OwnerUserID,
		SortOrder:           current.SortOrder,
		Icon:                current.Icon,
	}

	// Apply updates
	if req.Name != nil {
		params.Name = *req.Name
	}
	if len(req.Paths) > 0 {
		params.Paths = req.Paths
	}

	// Apply module-specific settings if provided
	if settings, ok := req.Settings.(*TVLibrarySettings); ok && settings != nil {
		if settings.PreferredLanguage != "" {
			params.PreferredLanguage = &settings.PreferredLanguage
		}
		params.TmdbEnabled = settings.TmdbEnabled
		params.TvdbEnabled = settings.TvdbEnabled
		params.DownloadBackdrops = settings.DownloadBackdrops
		params.DownloadNfo = settings.DownloadNfo
		params.GenerateChapters = settings.GenerateChapters
		if settings.SeasonFolderFormat != "" {
			params.SeasonFolderFormat = &settings.SeasonFolderFormat
		}
		if settings.EpisodeNamingFormat != "" {
			params.EpisodeNamingFormat = &settings.EpisodeNamingFormat
		}
		params.AutoAddMissing = settings.AutoAddMissing
		params.IsPrivate = settings.IsPrivate
		if settings.OwnerUserID != uuid.Nil {
			params.OwnerUserID = pgtype.UUID{Bytes: settings.OwnerUserID, Valid: true}
		}
		params.SortOrder = int32(settings.SortOrder)
		if settings.Icon != "" {
			params.Icon = &settings.Icon
		}
	}

	lib, err := s.queries.UpdateTVLibrary(ctx, params)
	if err != nil {
		return nil, err
	}

	s.logger.Info("tv library updated",
		"id", lib.ID,
		"name", lib.Name,
	)

	count, _ := s.queries.CountSeriesByTVLibrary(ctx, libraryID)
	info := s.toLibraryInfo(lib, count)
	return &info, nil
}

// DeleteLibrary removes a library and optionally its content.
func (s *LibraryService) DeleteLibrary(ctx context.Context, libraryID uuid.UUID, deleteContent bool) error {
	// Verify library exists
	lib, err := s.queries.GetTVLibraryByID(ctx, libraryID)
	if err != nil {
		return ErrLibraryNotFoundInService
	}

	if err := s.queries.DeleteTVLibrary(ctx, libraryID); err != nil {
		return err
	}

	s.logger.Info("tv library deleted",
		"id", lib.ID,
		"name", lib.Name,
		"delete_content", deleteContent,
	)

	return nil
}

// ScanLibrary triggers a library scan.
func (s *LibraryService) ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) error {
	// Verify library exists
	lib, err := s.queries.GetTVLibraryByID(ctx, libraryID)
	if err != nil {
		return ErrLibraryNotFoundInService
	}

	s.logger.Info("tv library scan triggered",
		"library_id", libraryID,
		"library_name", lib.Name,
		"full_scan", fullScan,
	)

	// TODO: Enqueue a River job for the actual scanning
	return nil
}

// toLibraryInfo converts a database TvLibrary to shared.LibraryInfo.
func (s *LibraryService) toLibraryInfo(lib tvshowdb.TvLibrary, itemCount int64) shared.LibraryInfo {
	var ownerID uuid.UUID
	if lib.OwnerUserID.Valid {
		ownerID = lib.OwnerUserID.Bytes
	}

	return shared.LibraryInfo{
		ID:        lib.ID,
		Module:    "tvshow",
		Name:      lib.Name,
		Paths:     lib.Paths,
		IsAdult:   false, // TV libraries are never adult content
		ItemCount: itemCount,
		Settings: &TVLibrarySettings{
			PreferredLanguage:   derefString(lib.PreferredLanguage),
			TmdbEnabled:         lib.TmdbEnabled,
			TvdbEnabled:         lib.TvdbEnabled,
			DownloadBackdrops:   lib.DownloadBackdrops,
			DownloadNfo:         lib.DownloadNfo,
			GenerateChapters:    lib.GenerateChapters,
			SeasonFolderFormat:  derefString(lib.SeasonFolderFormat),
			EpisodeNamingFormat: derefString(lib.EpisodeNamingFormat),
			AutoAddMissing:      lib.AutoAddMissing,
			IsPrivate:           lib.IsPrivate,
			OwnerUserID:         ownerID,
			SortOrder:           int(lib.SortOrder),
			Icon:                derefString(lib.Icon),
		},
	}
}

// TVLibrarySettings contains TV-specific library settings.
type TVLibrarySettings struct {
	PreferredLanguage   string    `json:"preferred_language,omitempty"`
	TmdbEnabled         bool      `json:"tmdb_enabled"`
	TvdbEnabled         bool      `json:"tvdb_enabled"`
	DownloadBackdrops   bool      `json:"download_backdrops"`
	DownloadNfo         bool      `json:"download_nfo"`
	GenerateChapters    bool      `json:"generate_chapters"`
	SeasonFolderFormat  string    `json:"season_folder_format,omitempty"`
	EpisodeNamingFormat string    `json:"episode_naming_format,omitempty"`
	AutoAddMissing      bool      `json:"auto_add_missing"`
	IsPrivate           bool      `json:"is_private"`
	OwnerUserID         uuid.UUID `json:"owner_user_id,omitempty"`
	SortOrder           int       `json:"sort_order"`
	Icon                string    `json:"icon,omitempty"`
}

// derefString safely dereferences a string pointer.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
