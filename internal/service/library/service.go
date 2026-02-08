package library

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/activity"
)

var (
	// ErrNotFound is returned when a library is not found.
	ErrNotFound = errors.New("library not found")
	// ErrScanNotFound is returned when a library scan is not found.
	ErrScanNotFound = errors.New("library scan not found")
	// ErrPermissionNotFound is returned when a permission is not found.
	ErrPermissionNotFound = errors.New("permission not found")
	// ErrInvalidLibraryType is returned when the library type is invalid.
	ErrInvalidLibraryType = errors.New("invalid library type")
	// ErrInvalidScanType is returned when the scan type is invalid.
	ErrInvalidScanType = errors.New("invalid scan type")
	// ErrInvalidPermission is returned when the permission is invalid.
	ErrInvalidPermission = errors.New("invalid permission")
	// ErrLibraryExists is returned when a library with the same name exists.
	ErrLibraryExists = errors.New("library with this name already exists")
	// ErrScanInProgress is returned when a scan is already running.
	ErrScanInProgress = errors.New("scan already in progress for this library")
	// ErrAccessDenied is returned when the user doesn't have permission.
	ErrAccessDenied = errors.New("access denied")
)

// Service provides library management functionality.
type Service struct {
	repo           Repository
	logger         *slog.Logger
	activityLogger activity.Logger
}

// NewService creates a new library service.
func NewService(repo Repository, logger *slog.Logger, activityLogger activity.Logger) *Service {
	return &Service{
		repo:           repo,
		logger:         logger.With("component", "library"),
		activityLogger: activityLogger,
	}
}

// ============================================================================
// Library CRUD
// ============================================================================

// CreateLibraryRequest represents a request to create a library.
type CreateLibraryRequest struct {
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`
	Paths              []string               `json:"paths"`
	Enabled            bool                   `json:"enabled"`
	ScanOnStartup      bool                   `json:"scan_on_startup"`
	RealtimeMonitoring bool                   `json:"realtime_monitoring"`
	MetadataProvider   *string                `json:"metadata_provider,omitempty"`
	PreferredLanguage  string                 `json:"preferred_language"`
	ScannerConfig      map[string]interface{} `json:"scanner_config,omitempty"`
}

// Create creates a new library.
func (s *Service) Create(ctx context.Context, req CreateLibraryRequest) (*Library, error) {
	// Validate library type
	if !IsValidLibraryType(req.Type) {
		return nil, ErrInvalidLibraryType
	}

	// Check if library with same name exists
	_, err := s.repo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, ErrLibraryExists
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	// Set defaults
	if req.PreferredLanguage == "" {
		req.PreferredLanguage = "en"
	}

	lib := &Library{
		Name:               req.Name,
		Type:               req.Type,
		Paths:              req.Paths,
		Enabled:            req.Enabled,
		ScanOnStartup:      req.ScanOnStartup,
		RealtimeMonitoring: req.RealtimeMonitoring,
		MetadataProvider:   req.MetadataProvider,
		PreferredLanguage:  req.PreferredLanguage,
		ScannerConfig:      req.ScannerConfig,
	}

	if err := s.repo.Create(ctx, lib); err != nil {
		s.logger.Error("failed to create library",
			slog.String("name", req.Name),
			slog.Any("error",err),
		)
		return nil, err
	}

	s.logger.Info("library created",
		slog.String("id", lib.ID.String()),
		slog.String("name", lib.Name),
		slog.String("type", lib.Type),
	)

	// Log library creation
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		Action:       activity.ActionLibraryCreate,
		ResourceType: activity.ResourceTypeLibrary,
		ResourceID:   lib.ID,
		Metadata: map[string]interface{}{
			"name":  lib.Name,
			"type":  lib.Type,
			"paths": lib.Paths,
		},
	})

	return lib, nil
}

// Get retrieves a library by ID.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Library, error) {
	return s.repo.Get(ctx, id)
}

// GetByName retrieves a library by name.
func (s *Service) GetByName(ctx context.Context, name string) (*Library, error) {
	return s.repo.GetByName(ctx, name)
}

// List returns all libraries.
func (s *Service) List(ctx context.Context) ([]Library, error) {
	return s.repo.List(ctx)
}

// ListEnabled returns all enabled libraries.
func (s *Service) ListEnabled(ctx context.Context) ([]Library, error) {
	return s.repo.ListEnabled(ctx)
}

// ListByType returns libraries by type.
func (s *Service) ListByType(ctx context.Context, libType string) ([]Library, error) {
	if !IsValidLibraryType(libType) {
		return nil, ErrInvalidLibraryType
	}
	return s.repo.ListByType(ctx, libType)
}

// ListAccessible returns libraries accessible to a user.
func (s *Service) ListAccessible(ctx context.Context, userID uuid.UUID) ([]Library, error) {
	return s.repo.GetUserAccessibleLibraries(ctx, userID)
}

// Update updates a library.
func (s *Service) Update(ctx context.Context, id uuid.UUID, update *LibraryUpdate) (*Library, error) {
	// Validate type if provided
	if update.Type != nil && !IsValidLibraryType(*update.Type) {
		return nil, ErrInvalidLibraryType
	}

	// Check if name is being changed to an existing name
	if update.Name != nil {
		existing, err := s.repo.GetByName(ctx, *update.Name)
		if err == nil && existing.ID != id {
			return nil, ErrLibraryExists
		}
		if err != nil && !errors.Is(err, ErrNotFound) {
			return nil, err
		}
	}

	lib, err := s.repo.Update(ctx, id, update)
	if err != nil {
		s.logger.Error("failed to update library",
			slog.String("id", id.String()),
			slog.Any("error",err),
		)
		return nil, err
	}

	s.logger.Info("library updated",
		slog.String("id", lib.ID.String()),
		slog.String("name", lib.Name),
	)

	// Log library update
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		Action:       activity.ActionLibraryUpdate,
		ResourceType: activity.ResourceTypeLibrary,
		ResourceID:   lib.ID,
		Metadata: map[string]interface{}{
			"name": lib.Name,
		},
	})

	return lib, nil
}

// Delete deletes a library.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Get library for logging
	lib, _ := s.repo.Get(ctx, id)

	// First revoke all permissions
	if err := s.repo.RevokeAllPermissions(ctx, id); err != nil {
		s.logger.Error("failed to revoke permissions for library",
			slog.String("id", id.String()),
			slog.Any("error",err),
		)
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete library",
			slog.String("id", id.String()),
			slog.Any("error",err),
		)
		return err
	}

	s.logger.Info("library deleted", slog.String("id", id.String()))

	// Log library deletion
	if lib != nil {
		_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
			Action:       activity.ActionLibraryDelete,
			ResourceType: activity.ResourceTypeLibrary,
			ResourceID:   id,
			Metadata: map[string]interface{}{
				"name": lib.Name,
			},
		})
	}

	return nil
}

// Count returns total library count.
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// ============================================================================
// Library Scans
// ============================================================================

// TriggerScan triggers a new library scan.
func (s *Service) TriggerScan(ctx context.Context, libraryID uuid.UUID, scanType string) (*LibraryScan, error) {
	// Validate scan type
	if !IsValidScanType(scanType) {
		return nil, ErrInvalidScanType
	}

	// Check if library exists
	lib, err := s.repo.Get(ctx, libraryID)
	if err != nil {
		return nil, err
	}

	// Check if there's already a running scan
	runningScans, err := s.repo.GetRunningScans(ctx)
	if err != nil {
		return nil, err
	}
	for _, scan := range runningScans {
		if scan.LibraryID == libraryID {
			return nil, ErrScanInProgress
		}
	}

	// Create scan record
	scan := &LibraryScan{
		LibraryID: libraryID,
		ScanType:  scanType,
		Status:    ScanStatusPending,
	}

	if err := s.repo.CreateScan(ctx, scan); err != nil {
		s.logger.Error("failed to create scan",
			slog.String("library_id", libraryID.String()),
			slog.String("scan_type", scanType),
			slog.Any("error",err),
		)
		return nil, err
	}

	s.logger.Info("scan triggered",
		slog.String("scan_id", scan.ID.String()),
		slog.String("library_id", libraryID.String()),
		slog.String("library_name", lib.Name),
		slog.String("scan_type", scanType),
	)

	return scan, nil
}

// GetScan retrieves a scan by ID.
func (s *Service) GetScan(ctx context.Context, id uuid.UUID) (*LibraryScan, error) {
	return s.repo.GetScan(ctx, id)
}

// ListScans lists scans for a library.
func (s *Service) ListScans(ctx context.Context, libraryID uuid.UUID, limit, offset int32) ([]LibraryScan, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	scans, err := s.repo.ListScans(ctx, libraryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.CountScans(ctx, libraryID)
	if err != nil {
		return nil, 0, err
	}

	return scans, count, nil
}

// GetLatestScan gets the most recent scan for a library.
func (s *Service) GetLatestScan(ctx context.Context, libraryID uuid.UUID) (*LibraryScan, error) {
	return s.repo.GetLatestScan(ctx, libraryID)
}

// GetRunningScans gets all running scans.
func (s *Service) GetRunningScans(ctx context.Context) ([]LibraryScan, error) {
	return s.repo.GetRunningScans(ctx)
}

// StartScan marks a scan as running.
func (s *Service) StartScan(ctx context.Context, scanID uuid.UUID) (*LibraryScan, error) {
	now := time.Now()
	return s.repo.UpdateScanStatus(ctx, scanID, &ScanStatusUpdate{
		Status:    ScanStatusRunning,
		StartedAt: &now,
	})
}

// CompleteScan marks a scan as completed.
func (s *Service) CompleteScan(ctx context.Context, scanID uuid.UUID, progress *ScanProgress) (*LibraryScan, error) {
	// Update progress first
	if progress != nil {
		if _, err := s.repo.UpdateScanProgress(ctx, scanID, progress); err != nil {
			return nil, err
		}
	}

	// Get scan to calculate duration
	scan, err := s.repo.GetScan(ctx, scanID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var duration int32
	if scan.StartedAt != nil {
		duration = int32(now.Sub(*scan.StartedAt).Seconds())
	}

	return s.repo.UpdateScanStatus(ctx, scanID, &ScanStatusUpdate{
		Status:          ScanStatusCompleted,
		CompletedAt:     &now,
		DurationSeconds: &duration,
	})
}

// FailScan marks a scan as failed.
func (s *Service) FailScan(ctx context.Context, scanID uuid.UUID, errorMsg string) (*LibraryScan, error) {
	// Get scan to calculate duration
	scan, err := s.repo.GetScan(ctx, scanID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var duration int32
	if scan.StartedAt != nil {
		duration = int32(now.Sub(*scan.StartedAt).Seconds())
	}

	return s.repo.UpdateScanStatus(ctx, scanID, &ScanStatusUpdate{
		Status:          ScanStatusFailed,
		CompletedAt:     &now,
		DurationSeconds: &duration,
		ErrorMessage:    &errorMsg,
	})
}

// CancelScan marks a scan as cancelled.
func (s *Service) CancelScan(ctx context.Context, scanID uuid.UUID) (*LibraryScan, error) {
	now := time.Now()
	return s.repo.UpdateScanStatus(ctx, scanID, &ScanStatusUpdate{
		Status:      ScanStatusCancelled,
		CompletedAt: &now,
	})
}

// UpdateScanProgress updates scan progress.
func (s *Service) UpdateScanProgress(ctx context.Context, scanID uuid.UUID, progress *ScanProgress) (*LibraryScan, error) {
	return s.repo.UpdateScanProgress(ctx, scanID, progress)
}

// ============================================================================
// Library Permissions
// ============================================================================

// GrantPermission grants a permission to a user for a library.
func (s *Service) GrantPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) error {
	// Validate permission
	if !IsValidPermission(permission) {
		return ErrInvalidPermission
	}

	// Check if library exists
	if _, err := s.repo.Get(ctx, libraryID); err != nil {
		return err
	}

	perm := &Permission{
		LibraryID:  libraryID,
		UserID:     userID,
		Permission: permission,
	}

	if err := s.repo.GrantPermission(ctx, perm); err != nil {
		s.logger.Error("failed to grant permission",
			slog.String("library_id", libraryID.String()),
			slog.String("user_id", userID.String()),
			slog.String("permission", permission),
			slog.Any("error",err),
		)
		return err
	}

	s.logger.Info("permission granted",
		slog.String("library_id", libraryID.String()),
		slog.String("user_id", userID.String()),
		slog.String("permission", permission),
	)

	return nil
}

// RevokePermission revokes a permission from a user.
func (s *Service) RevokePermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) error {
	if !IsValidPermission(permission) {
		return ErrInvalidPermission
	}

	if err := s.repo.RevokePermission(ctx, libraryID, userID, permission); err != nil {
		s.logger.Error("failed to revoke permission",
			slog.String("library_id", libraryID.String()),
			slog.String("user_id", userID.String()),
			slog.String("permission", permission),
			slog.Any("error",err),
		)
		return err
	}

	s.logger.Info("permission revoked",
		slog.String("library_id", libraryID.String()),
		slog.String("user_id", userID.String()),
		slog.String("permission", permission),
	)

	return nil
}

// CheckPermission checks if a user has a specific permission for a library.
func (s *Service) CheckPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) (bool, error) {
	if !IsValidPermission(permission) {
		return false, ErrInvalidPermission
	}
	return s.repo.CheckPermission(ctx, libraryID, userID, permission)
}

// ListPermissions lists all permissions for a library.
func (s *Service) ListPermissions(ctx context.Context, libraryID uuid.UUID) ([]Permission, error) {
	return s.repo.ListPermissions(ctx, libraryID)
}

// ListUserPermissions lists all library permissions for a user.
func (s *Service) ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error) {
	return s.repo.ListUserPermissions(ctx, userID)
}

// GetPermission gets a specific permission.
func (s *Service) GetPermission(ctx context.Context, libraryID, userID uuid.UUID, permission string) (*Permission, error) {
	return s.repo.GetPermission(ctx, libraryID, userID, permission)
}

// CanAccess checks if a user can access a library (has view permission or is admin).
func (s *Service) CanAccess(ctx context.Context, libraryID, userID uuid.UUID, isAdmin bool) (bool, error) {
	if isAdmin {
		return true, nil
	}
	return s.repo.CheckPermission(ctx, libraryID, userID, PermissionView)
}

// CanDownload checks if a user can download from a library.
func (s *Service) CanDownload(ctx context.Context, libraryID, userID uuid.UUID, isAdmin bool) (bool, error) {
	if isAdmin {
		return true, nil
	}
	return s.repo.CheckPermission(ctx, libraryID, userID, PermissionDownload)
}

// CanManage checks if a user can manage a library.
func (s *Service) CanManage(ctx context.Context, libraryID, userID uuid.UUID, isAdmin bool) (bool, error) {
	if isAdmin {
		return true, nil
	}
	return s.repo.CheckPermission(ctx, libraryID, userID, PermissionManage)
}
