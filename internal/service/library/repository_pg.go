package library

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// RepositoryPg implements Repository using PostgreSQL.
type RepositoryPg struct {
	queries *db.Queries
}

// NewRepositoryPg creates a new PostgreSQL repository.
func NewRepositoryPg(queries *db.Queries) *RepositoryPg {
	return &RepositoryPg{queries: queries}
}

// ============================================================================
// Library CRUD
// ============================================================================

// Create creates a new library.
func (r *RepositoryPg) Create(ctx context.Context, lib *Library) error {
	var scannerConfig []byte
	if lib.ScannerConfig != nil {
		var err error
		scannerConfig, err = json.Marshal(lib.ScannerConfig)
		if err != nil {
			return err
		}
	}

	result, err := r.queries.CreateLibrary(ctx, db.CreateLibraryParams{
		Name:               lib.Name,
		Type:               lib.Type,
		Paths:              lib.Paths,
		Enabled:            lib.Enabled,
		ScanOnStartup:      lib.ScanOnStartup,
		RealtimeMonitoring: lib.RealtimeMonitoring,
		MetadataProvider:   lib.MetadataProvider,
		PreferredLanguage:  lib.PreferredLanguage,
		ScannerConfig:      scannerConfig,
	})
	if err != nil {
		return err
	}

	lib.ID = result.ID
	lib.CreatedAt = result.CreatedAt
	lib.UpdatedAt = result.UpdatedAt
	return nil
}

// Get retrieves a library by ID.
func (r *RepositoryPg) Get(ctx context.Context, id uuid.UUID) (*Library, error) {
	result, err := r.queries.GetLibrary(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dbLibraryToLibrary(result), nil
}

// GetByName retrieves a library by name.
func (r *RepositoryPg) GetByName(ctx context.Context, name string) (*Library, error) {
	result, err := r.queries.GetLibraryByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dbLibraryToLibrary(result), nil
}

// List returns all libraries.
func (r *RepositoryPg) List(ctx context.Context) ([]Library, error) {
	results, err := r.queries.ListLibraries(ctx)
	if err != nil {
		return nil, err
	}
	libraries := make([]Library, len(results))
	for i, result := range results {
		libraries[i] = *dbLibraryToLibrary(result)
	}
	return libraries, nil
}

// ListEnabled returns all enabled libraries.
func (r *RepositoryPg) ListEnabled(ctx context.Context) ([]Library, error) {
	results, err := r.queries.ListEnabledLibraries(ctx)
	if err != nil {
		return nil, err
	}
	libraries := make([]Library, len(results))
	for i, result := range results {
		libraries[i] = *dbLibraryToLibrary(result)
	}
	return libraries, nil
}

// ListByType returns libraries by type.
func (r *RepositoryPg) ListByType(ctx context.Context, libType string) ([]Library, error) {
	results, err := r.queries.ListLibrariesByType(ctx, libType)
	if err != nil {
		return nil, err
	}
	libraries := make([]Library, len(results))
	for i, result := range results {
		libraries[i] = *dbLibraryToLibrary(result)
	}
	return libraries, nil
}

// Update updates a library.
func (r *RepositoryPg) Update(ctx context.Context, id uuid.UUID, update *LibraryUpdate) (*Library, error) {
	params := db.UpdateLibraryParams{
		ID:                 id,
		Name:               update.Name,
		Type:               update.Type,
		Paths:              update.Paths,
		Enabled:            update.Enabled,
		ScanOnStartup:      update.ScanOnStartup,
		RealtimeMonitoring: update.RealtimeMonitoring,
		MetadataProvider:   update.MetadataProvider,
		PreferredLanguage:  update.PreferredLanguage,
	}

	if update.ScannerConfig != nil {
		scannerConfig, err := json.Marshal(update.ScannerConfig)
		if err != nil {
			return nil, err
		}
		params.ScannerConfig = scannerConfig
	}

	result, err := r.queries.UpdateLibrary(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dbLibraryToLibrary(result), nil
}

// Delete deletes a library.
func (r *RepositoryPg) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteLibrary(ctx, id)
}

// Count returns total library count.
func (r *RepositoryPg) Count(ctx context.Context) (int64, error) {
	return r.queries.CountLibraries(ctx)
}

// CountByType returns library count by type.
func (r *RepositoryPg) CountByType(ctx context.Context, libType string) (int64, error) {
	return r.queries.CountLibrariesByType(ctx, libType)
}

// ============================================================================
// Library Scans
// ============================================================================

// CreateScan creates a new library scan.
func (r *RepositoryPg) CreateScan(ctx context.Context, scan *LibraryScan) error {
	result, err := r.queries.CreateLibraryScan(ctx, db.CreateLibraryScanParams{
		LibraryID: scan.LibraryID,
		ScanType:  scan.ScanType,
		Status:    scan.Status,
	})
	if err != nil {
		return err
	}

	scan.ID = result.ID
	scan.CreatedAt = result.CreatedAt
	return nil
}

// GetScan retrieves a scan by ID.
func (r *RepositoryPg) GetScan(ctx context.Context, id uuid.UUID) (*LibraryScan, error) {
	result, err := r.queries.GetLibraryScan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}
	return dbScanToScan(result), nil
}

// ListScans lists scans for a library.
func (r *RepositoryPg) ListScans(ctx context.Context, libraryID uuid.UUID, limit, offset int32) ([]LibraryScan, error) {
	results, err := r.queries.ListLibraryScans(ctx, db.ListLibraryScansParams{
		LibraryID: libraryID,
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, err
	}
	scans := make([]LibraryScan, len(results))
	for i, result := range results {
		scans[i] = *dbScanToScan(result)
	}
	return scans, nil
}

// CountScans counts scans for a library.
func (r *RepositoryPg) CountScans(ctx context.Context, libraryID uuid.UUID) (int64, error) {
	return r.queries.CountLibraryScans(ctx, libraryID)
}

// GetLatestScan gets the most recent scan for a library.
func (r *RepositoryPg) GetLatestScan(ctx context.Context, libraryID uuid.UUID) (*LibraryScan, error) {
	result, err := r.queries.GetLatestLibraryScan(ctx, libraryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}
	return dbScanToScan(result), nil
}

// GetRunningScans gets all running scans.
func (r *RepositoryPg) GetRunningScans(ctx context.Context) ([]LibraryScan, error) {
	results, err := r.queries.GetRunningScans(ctx)
	if err != nil {
		return nil, err
	}
	scans := make([]LibraryScan, len(results))
	for i, result := range results {
		scans[i] = *dbScanToScan(result)
	}
	return scans, nil
}

// UpdateScanStatus updates scan status.
func (r *RepositoryPg) UpdateScanStatus(ctx context.Context, id uuid.UUID, status *ScanStatusUpdate) (*LibraryScan, error) {
	params := db.UpdateLibraryScanStatusParams{
		ID:           id,
		Status:       status.Status,
		ErrorMessage: status.ErrorMessage,
	}

	if status.StartedAt != nil {
		params.StartedAt = pgtype.Timestamptz{Time: *status.StartedAt, Valid: true}
	}
	if status.CompletedAt != nil {
		params.CompletedAt = pgtype.Timestamptz{Time: *status.CompletedAt, Valid: true}
	}
	if status.DurationSeconds != nil {
		params.DurationSeconds = status.DurationSeconds
	}

	result, err := r.queries.UpdateLibraryScanStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}
	return dbScanToScan(result), nil
}

// UpdateScanProgress updates scan progress.
func (r *RepositoryPg) UpdateScanProgress(ctx context.Context, id uuid.UUID, progress *ScanProgress) (*LibraryScan, error) {
	result, err := r.queries.UpdateLibraryScanProgress(ctx, db.UpdateLibraryScanProgressParams{
		ID:           id,
		ItemsScanned: progress.ItemsScanned,
		ItemsAdded:   progress.ItemsAdded,
		ItemsUpdated: progress.ItemsUpdated,
		ItemsRemoved: progress.ItemsRemoved,
		ErrorsCount:  progress.ErrorsCount,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}
	return dbScanToScan(result), nil
}

// DeleteOldScans deletes scans older than the given time.
func (r *RepositoryPg) DeleteOldScans(ctx context.Context, olderThan time.Time) (int64, error) {
	return r.queries.DeleteOldLibraryScans(ctx, olderThan)
}

// ============================================================================
// Library Permissions
// ============================================================================

// GrantPermission grants a permission to a user.
func (r *RepositoryPg) GrantPermission(ctx context.Context, perm *Permission) error {
	result, err := r.queries.CreateLibraryPermission(ctx, db.CreateLibraryPermissionParams{
		LibraryID:  perm.LibraryID,
		UserID:     perm.UserID,
		Permission: perm.Permission,
	})
	if err != nil {
		return err
	}

	// Result may be empty due to ON CONFLICT DO NOTHING
	if result.ID != uuid.Nil {
		perm.ID = result.ID
		perm.CreatedAt = result.CreatedAt
	}
	return nil
}

// GetPermission retrieves a specific permission.
func (r *RepositoryPg) GetPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) (*Permission, error) {
	result, err := r.queries.GetLibraryPermission(ctx, db.GetLibraryPermissionParams{
		LibraryID:  libraryID,
		UserID:     userID,
		Permission: permission,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return dbPermToPermission(result), nil
}

// ListPermissions lists all permissions for a library.
func (r *RepositoryPg) ListPermissions(ctx context.Context, libraryID uuid.UUID) ([]Permission, error) {
	results, err := r.queries.ListLibraryPermissions(ctx, libraryID)
	if err != nil {
		return nil, err
	}
	perms := make([]Permission, len(results))
	for i, result := range results {
		perms[i] = *dbPermToPermission(result)
	}
	return perms, nil
}

// ListUserPermissions lists all library permissions for a user.
func (r *RepositoryPg) ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error) {
	results, err := r.queries.ListUserLibraryPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	perms := make([]Permission, len(results))
	for i, result := range results {
		perms[i] = *dbPermToPermission(result)
	}
	return perms, nil
}

// CheckPermission checks if a user has a specific permission.
func (r *RepositoryPg) CheckPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) (bool, error) {
	return r.queries.CheckLibraryPermission(ctx, db.CheckLibraryPermissionParams{
		LibraryID:  libraryID,
		UserID:     userID,
		Permission: permission,
	})
}

// GetUserAccessibleLibraries returns libraries a user can access.
func (r *RepositoryPg) GetUserAccessibleLibraries(ctx context.Context, userID uuid.UUID) ([]Library, error) {
	results, err := r.queries.GetUserAccessibleLibraries(ctx, userID)
	if err != nil {
		return nil, err
	}
	libraries := make([]Library, len(results))
	for i, result := range results {
		libraries[i] = *dbLibraryToLibrary(result)
	}
	return libraries, nil
}

// RevokePermission revokes a permission from a user.
func (r *RepositoryPg) RevokePermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) error {
	return r.queries.DeleteLibraryPermission(ctx, db.DeleteLibraryPermissionParams{
		LibraryID:  libraryID,
		UserID:     userID,
		Permission: permission,
	})
}

// RevokeAllPermissions revokes all permissions for a library.
func (r *RepositoryPg) RevokeAllPermissions(ctx context.Context, libraryID uuid.UUID) error {
	return r.queries.DeleteAllLibraryPermissions(ctx, libraryID)
}

// RevokeUserPermissions revokes all library permissions for a user.
func (r *RepositoryPg) RevokeUserPermissions(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteUserLibraryPermissions(ctx, userID)
}

// CountPermissions counts permissions for a library.
func (r *RepositoryPg) CountPermissions(ctx context.Context, libraryID uuid.UUID) (int64, error) {
	return r.queries.CountLibraryPermissions(ctx, libraryID)
}

// ============================================================================
// Helper Functions
// ============================================================================

func dbLibraryToLibrary(lib db.Library) *Library {
	result := &Library{
		ID:                 lib.ID,
		Name:               lib.Name,
		Type:               lib.Type,
		Paths:              lib.Paths,
		Enabled:            lib.Enabled,
		ScanOnStartup:      lib.ScanOnStartup,
		RealtimeMonitoring: lib.RealtimeMonitoring,
		MetadataProvider:   lib.MetadataProvider,
		PreferredLanguage:  lib.PreferredLanguage,
		CreatedAt:          lib.CreatedAt,
		UpdatedAt:          lib.UpdatedAt,
	}

	if len(lib.ScannerConfig) > 0 {
		var config map[string]interface{}
		if err := json.Unmarshal(lib.ScannerConfig, &config); err == nil {
			result.ScannerConfig = config
		}
	}

	return result
}

func dbScanToScan(scan db.LibraryScan) *LibraryScan {
	result := &LibraryScan{
		ID:              scan.ID,
		LibraryID:       scan.LibraryID,
		ScanType:        scan.ScanType,
		Status:          scan.Status,
		ItemsScanned:    scan.ItemsScanned,
		ItemsAdded:      scan.ItemsAdded,
		ItemsUpdated:    scan.ItemsUpdated,
		ItemsRemoved:    scan.ItemsRemoved,
		ErrorsCount:     scan.ErrorsCount,
		ErrorMessage:    scan.ErrorMessage,
		DurationSeconds: scan.DurationSeconds,
		CreatedAt:       scan.CreatedAt,
	}

	if scan.StartedAt.Valid {
		result.StartedAt = &scan.StartedAt.Time
	}
	if scan.CompletedAt.Valid {
		result.CompletedAt = &scan.CompletedAt.Time
	}

	return result
}

func dbPermToPermission(perm db.LibraryPermission) *Permission {
	return &Permission{
		ID:         perm.ID,
		LibraryID:  perm.LibraryID,
		UserID:     perm.UserID,
		Permission: perm.Permission,
		CreatedAt:  perm.CreatedAt,
	}
}
