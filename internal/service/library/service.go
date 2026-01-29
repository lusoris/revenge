// Package library provides library management services.
package library

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	// ErrLibraryNotFound indicates the library was not found.
	ErrLibraryNotFound = errors.New("library not found")
	// ErrAccessDenied indicates the user doesn't have access to the library.
	ErrAccessDenied = errors.New("access denied")
	// ErrInvalidLibraryType indicates an invalid library type.
	ErrInvalidLibraryType = errors.New("invalid library type")
)

// Service provides library management operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new library service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "library")),
	}
}

// LibraryWithAccess contains a library and access info.
type LibraryWithAccess struct {
	Library   db.Library
	CanManage bool
}

// CreateParams contains parameters for creating a library.
type CreateParams struct {
	Name              string
	LibraryType       string // String type that gets converted to db.LibraryType
	Paths             []string
	ScanEnabled       bool
	ScanIntervalHours int32
	PreferredLanguage *string
	DownloadImages    bool
	DownloadNfo       bool
	GenerateChapters  bool
	IsPrivate         bool
	OwnerUserID       pgtype.UUID
	SortOrder         int32
	Icon              *string
}

// Create creates a new library.
func (s *Service) Create(ctx context.Context, params CreateParams) (*db.Library, error) {
	// Convert and validate library type
	libType := db.LibraryType(params.LibraryType)
	if !isValidLibraryType(libType) {
		return nil, ErrInvalidLibraryType
	}

	library, err := s.queries.CreateLibrary(ctx, db.CreateLibraryParams{
		Name:              params.Name,
		Type:              libType,
		Paths:             params.Paths,
		ScanEnabled:       params.ScanEnabled,
		ScanIntervalHours: params.ScanIntervalHours,
		PreferredLanguage: params.PreferredLanguage,
		DownloadImages:    params.DownloadImages,
		DownloadNfo:       params.DownloadNfo,
		GenerateChapters:  params.GenerateChapters,
		IsPrivate:         params.IsPrivate,
		OwnerUserID:       params.OwnerUserID,
		SortOrder:         params.SortOrder,
		Icon:              params.Icon,
	})
	if err != nil {
		return nil, fmt.Errorf("create library: %w", err)
	}

	s.logger.Info("Library created",
		slog.String("library_id", library.ID.String()),
		slog.String("name", library.Name),
		slog.String("type", string(library.Type)),
	)

	return &library, nil
}

// GetByID retrieves a library by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.Library, error) {
	library, err := s.queries.GetLibraryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get library: %w", err)
	}
	return &library, nil
}

// GetByIDWithAccess retrieves a library by ID, checking user access.
func (s *Service) GetByIDWithAccess(ctx context.Context, id, userID uuid.UUID, isAdmin bool) (*db.Library, error) {
	library, err := s.queries.GetLibraryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get library: %w", err)
	}

	// Admins have access to all libraries
	if isAdmin {
		return &library, nil
	}

	// Check access
	hasAccess, err := s.UserCanAccess(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrAccessDenied
	}

	return &library, nil
}

// List returns all libraries (admin only).
func (s *Service) List(ctx context.Context) ([]db.Library, error) {
	libraries, err := s.queries.ListLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("list libraries: %w", err)
	}
	return libraries, nil
}

// ListAccessible returns libraries accessible to a user.
func (s *Service) ListAccessible(ctx context.Context, userID uuid.UUID) ([]db.Library, error) {
	libraries, err := s.queries.ListAccessibleLibraries(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list accessible libraries: %w", err)
	}
	return libraries, nil
}

// ListAll returns all libraries with access info (for admins).
func (s *Service) ListAll(ctx context.Context) ([]LibraryWithAccess, error) {
	libraries, err := s.queries.ListLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("list libraries: %w", err)
	}

	result := make([]LibraryWithAccess, len(libraries))
	for i, lib := range libraries {
		result[i] = LibraryWithAccess{
			Library:   lib,
			CanManage: true, // Admins can manage all
		}
	}
	return result, nil
}

// ListForUser returns libraries accessible to a user with access info.
func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID) ([]LibraryWithAccess, error) {
	libraries, err := s.queries.ListAccessibleLibraries(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list accessible libraries: %w", err)
	}

	result := make([]LibraryWithAccess, len(libraries))
	for i, lib := range libraries {
		// Check if user can manage this library
		canManage := false
		if lib.OwnerUserID.Valid && lib.OwnerUserID.Bytes == userID {
			canManage = true
		}
		result[i] = LibraryWithAccess{
			Library:   lib,
			CanManage: canManage,
		}
	}
	return result, nil
}

// ListByType returns libraries of a specific type.
func (s *Service) ListByType(ctx context.Context, libraryType db.LibraryType) ([]db.Library, error) {
	libraries, err := s.queries.ListLibrariesByType(ctx, libraryType)
	if err != nil {
		return nil, fmt.Errorf("list libraries by type: %w", err)
	}
	return libraries, nil
}

// UpdateParams contains parameters for updating a library.
type UpdateParams struct {
	ID                uuid.UUID
	Name              *string
	Paths             []string
	ScanEnabled       *bool
	ScanIntervalHours *int32
	PreferredLanguage *string
	DownloadImages    *bool
	DownloadNfo       *bool
	GenerateChapters  *bool
	IsPrivate         *bool
	SortOrder         *int32
	Icon              *string
}

// Update updates a library.
func (s *Service) Update(ctx context.Context, params UpdateParams) (*db.Library, error) {
	library, err := s.queries.UpdateLibrary(ctx, db.UpdateLibraryParams{
		ID:                params.ID,
		Name:              params.Name,
		Paths:             params.Paths,
		ScanEnabled:       params.ScanEnabled,
		ScanIntervalHours: params.ScanIntervalHours,
		PreferredLanguage: params.PreferredLanguage,
		DownloadImages:    params.DownloadImages,
		DownloadNfo:       params.DownloadNfo,
		GenerateChapters:  params.GenerateChapters,
		IsPrivate:         params.IsPrivate,
		SortOrder:         params.SortOrder,
		Icon:              params.Icon,
	})
	if err != nil {
		return nil, fmt.Errorf("update library: %w", err)
	}

	s.logger.Info("Library updated",
		slog.String("library_id", library.ID.String()),
	)

	return &library, nil
}

// Delete deletes a library.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteLibrary(ctx, id); err != nil {
		return fmt.Errorf("delete library: %w", err)
	}

	s.logger.Info("Library deleted",
		slog.String("library_id", id.String()),
	)

	return nil
}

// GrantAccess grants a user access to a library.
func (s *Service) GrantAccess(ctx context.Context, libraryID, userID uuid.UUID, canManage bool) error {
	if err := s.queries.GrantLibraryAccess(ctx, db.GrantLibraryAccessParams{
		LibraryID: libraryID,
		UserID:    userID,
		CanManage: canManage,
	}); err != nil {
		return fmt.Errorf("grant access: %w", err)
	}

	s.logger.Info("Library access granted",
		slog.String("library_id", libraryID.String()),
		slog.String("user_id", userID.String()),
		slog.Bool("can_manage", canManage),
	)

	return nil
}

// RevokeAccess revokes a user's access to a library.
func (s *Service) RevokeAccess(ctx context.Context, libraryID, userID uuid.UUID) error {
	if err := s.queries.RevokeLibraryAccess(ctx, db.RevokeLibraryAccessParams{
		LibraryID: libraryID,
		UserID:    userID,
	}); err != nil {
		return fmt.Errorf("revoke access: %w", err)
	}

	s.logger.Info("Library access revoked",
		slog.String("library_id", libraryID.String()),
		slog.String("user_id", userID.String()),
	)

	return nil
}

// ListUsers returns all users with access to a library.
func (s *Service) ListUsers(ctx context.Context, libraryID uuid.UUID) ([]db.ListLibraryUsersRow, error) {
	users, err := s.queries.ListLibraryUsers(ctx, libraryID)
	if err != nil {
		return nil, fmt.Errorf("list library users: %w", err)
	}
	return users, nil
}

// UserCanAccess checks if a user can access a library.
func (s *Service) UserCanAccess(ctx context.Context, libraryID, userID uuid.UUID) (bool, error) {
	canAccess, err := s.queries.UserCanAccessLibrary(ctx, db.UserCanAccessLibraryParams{
		ID:          libraryID,
		OwnerUserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("check access: %w", err)
	}
	return canAccess, nil
}

// Count returns the total number of libraries.
func (s *Service) Count(ctx context.Context) (int64, error) {
	count, err := s.queries.CountLibraries(ctx)
	if err != nil {
		return 0, fmt.Errorf("count libraries: %w", err)
	}
	return count, nil
}

// isValidLibraryType validates a library type.
func isValidLibraryType(t db.LibraryType) bool {
	switch t {
	case db.LibraryTypeMovie,
		db.LibraryTypeTvshow,
		db.LibraryTypeMusic,
		db.LibraryTypeAudiobook,
		db.LibraryTypeBook,
		db.LibraryTypePodcast,
		db.LibraryTypePhoto,
		db.LibraryTypeLivetv,
		db.LibraryTypeComics,
		db.LibraryTypeAdultMovie,
		db.LibraryTypeAdultShow:
		return true
	default:
		return false
	}
}
