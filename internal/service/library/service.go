// Package library provides the library management service.
package library

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/domain"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Service handles library management operations.
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

// GetByID retrieves a library by its unique ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Library, error) {
	row, err := s.queries.GetLibraryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get library: %w", err)
	}

	return s.rowToLibrary(&row), nil
}

// GetByName retrieves a library by its name.
func (s *Service) GetByName(ctx context.Context, name string) (*domain.Library, error) {
	row, err := s.queries.GetLibraryByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get library by name: %w", err)
	}

	return s.getByNameRowToLibrary(&row), nil
}

// List retrieves all libraries.
func (s *Service) List(ctx context.Context) ([]*domain.Library, error) {
	rows, err := s.queries.ListLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list libraries: %w", err)
	}

	libraries := make([]*domain.Library, len(rows))
	for i := range rows {
		libraries[i] = s.listRowToLibrary(&rows[i])
	}

	return libraries, nil
}

// ListByType retrieves libraries of a specific type.
func (s *Service) ListByType(ctx context.Context, libType domain.LibraryType) ([]*domain.Library, error) {
	rows, err := s.queries.ListLibrariesByType(ctx, db.LibraryType(libType))
	if err != nil {
		return nil, fmt.Errorf("failed to list libraries by type: %w", err)
	}

	libraries := make([]*domain.Library, len(rows))
	for i := range rows {
		libraries[i] = s.listByTypeRowToLibrary(&rows[i])
	}

	return libraries, nil
}

// ListVisible retrieves only visible libraries.
func (s *Service) ListVisible(ctx context.Context) ([]*domain.Library, error) {
	rows, err := s.queries.ListVisibleLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list visible libraries: %w", err)
	}

	libraries := make([]*domain.Library, len(rows))
	for i := range rows {
		libraries[i] = s.listVisibleRowToLibrary(&rows[i])
	}

	return libraries, nil
}

// ListNonAdult retrieves libraries that are not adult content.
func (s *Service) ListNonAdult(ctx context.Context) ([]*domain.Library, error) {
	rows, err := s.queries.ListNonAdultLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list non-adult libraries: %w", err)
	}

	libraries := make([]*domain.Library, len(rows))
	for i := range rows {
		libraries[i] = s.listNonAdultRowToLibrary(&rows[i])
	}

	return libraries, nil
}

// ListForUser retrieves libraries accessible to a specific user.
func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Library, error) {
	rows, err := s.queries.ListLibrariesForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list libraries for user: %w", err)
	}

	libraries := make([]*domain.Library, len(rows))
	for i := range rows {
		libraries[i] = s.listForUserRowToLibrary(&rows[i])
	}

	return libraries, nil
}

// Create creates a new library.
func (s *Service) Create(ctx context.Context, params domain.CreateLibraryParams) (*domain.Library, error) {
	// Validate library type
	if !params.Type.IsValid() {
		return nil, fmt.Errorf("invalid library type: %s", params.Type)
	}

	// Check if name already exists
	exists, err := s.queries.LibraryNameExists(ctx, params.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check library name: %w", err)
	}
	if exists {
		return nil, domain.ErrAlreadyExists
	}

	// Determine if adult library based on type
	isAdult := params.Type.IsAdultType()

	// Convert settings to JSON
	settingsJSON, err := json.Marshal(params.Settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Convert scan interval
	var scanInterval pgtype.Int4
	if params.ScanIntervalHours != nil {
		scanInterval = pgtype.Int4{Int32: int32(*params.ScanIntervalHours), Valid: true}
	}

	row, err := s.queries.CreateLibrary(ctx, db.CreateLibraryParams{
		Name:              params.Name,
		Type:              db.LibraryType(params.Type),
		Paths:             params.Paths,
		Settings:          settingsJSON,
		IsVisible:         params.IsVisible,
		IsAdult:           isAdult,
		ScanIntervalHours: scanInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create library: %w", err)
	}

	s.logger.Info("library created",
		slog.String("library_id", row.ID.String()),
		slog.String("name", row.Name),
		slog.String("type", string(row.Type)),
	)

	return s.createRowToLibrary(&row), nil
}

// Update updates an existing library.
func (s *Service) Update(ctx context.Context, params domain.UpdateLibraryParams) error {
	// Check if name already exists (excluding current library)
	if params.Name != nil {
		exists, err := s.queries.LibraryNameExistsExcluding(ctx, db.LibraryNameExistsExcludingParams{
			Name: *params.Name,
			ID:   params.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to check library name: %w", err)
		}
		if exists {
			return domain.ErrAlreadyExists
		}
	}

	// Build update params with nullable types
	updateParams := db.UpdateLibraryParams{
		ID: params.ID,
	}

	if params.Name != nil {
		updateParams.Name = pgtype.Text{String: *params.Name, Valid: true}
	}

	if params.Paths != nil {
		updateParams.Paths = params.Paths
	}

	if params.Settings != nil {
		settingsJSON, err := json.Marshal(params.Settings)
		if err != nil {
			return fmt.Errorf("failed to marshal settings: %w", err)
		}
		updateParams.Settings = settingsJSON
	}

	if params.IsVisible != nil {
		updateParams.IsVisible = pgtype.Bool{Bool: *params.IsVisible, Valid: true}
	}

	if params.ScanIntervalHours != nil {
		updateParams.ScanIntervalHours = pgtype.Int4{Int32: int32(*params.ScanIntervalHours), Valid: true}
	}

	err := s.queries.UpdateLibrary(ctx, updateParams)
	if err != nil {
		return fmt.Errorf("failed to update library: %w", err)
	}

	s.logger.Info("library updated", slog.String("library_id", params.ID.String()))

	return nil
}

// Delete removes a library by its ID.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteLibrary(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete library: %w", err)
	}

	s.logger.Info("library deleted", slog.String("library_id", id.String()))

	return nil
}

// UpdateLastScan updates the library's last scan timestamp.
func (s *Service) UpdateLastScan(ctx context.Context, id uuid.UUID) error {
	err := s.queries.UpdateLibraryLastScan(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update last scan: %w", err)
	}

	return nil
}

// Count returns the total number of libraries.
func (s *Service) Count(ctx context.Context) (int64, error) {
	count, err := s.queries.CountLibraries(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count libraries: %w", err)
	}

	return count, nil
}

// Helper conversion functions for different row types

func (s *Service) rowToLibrary(row *db.GetLibraryByIDRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) listRowToLibrary(row *db.ListLibrariesRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) listByTypeRowToLibrary(row *db.ListLibrariesByTypeRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) listVisibleRowToLibrary(row *db.ListVisibleLibrariesRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) listNonAdultRowToLibrary(row *db.ListNonAdultLibrariesRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) listForUserRowToLibrary(row *db.ListLibrariesForUserRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) getByNameRowToLibrary(row *db.GetLibraryByNameRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) createRowToLibrary(row *db.CreateLibraryRow) *domain.Library {
	return s.convertLibrary(
		row.ID, row.Name, row.Type, row.Paths, row.Settings,
		row.IsVisible, row.IsAdult, row.ScanIntervalHours, row.LastScanAt,
		row.CreatedAt, row.UpdatedAt,
	)
}

func (s *Service) convertLibrary(
	id uuid.UUID,
	name string,
	libType db.LibraryType,
	paths []string,
	settingsRaw json.RawMessage,
	isVisible bool,
	isAdult bool,
	scanIntervalHours pgtype.Int4,
	lastScanAt pgtype.Timestamptz,
	createdAt, updatedAt time.Time,
) *domain.Library {
	var settings map[string]any
	if len(settingsRaw) > 0 {
		_ = json.Unmarshal(settingsRaw, &settings)
	}

	var scanInterval *int
	if scanIntervalHours.Valid {
		v := int(scanIntervalHours.Int32)
		scanInterval = &v
	}

	var lastScan *time.Time
	if lastScanAt.Valid {
		lastScan = &lastScanAt.Time
	}

	return &domain.Library{
		ID:                id,
		Name:              name,
		Type:              domain.LibraryType(libType),
		Paths:             paths,
		Settings:          settings,
		IsVisible:         isVisible,
		IsAdult:           isAdult,
		ScanIntervalHours: scanInterval,
		LastScanAt:        lastScan,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
