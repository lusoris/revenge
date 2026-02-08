package api

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/validate"
	"log/slog"
)

// ============================================================================
// Library Endpoints
// ============================================================================

// ListLibraries returns libraries accessible to the authenticated user.
// Admins see all libraries, regular users see only libraries they have permission to access.
// GET /api/v1/libraries
func (h *Handler) ListLibraries(ctx context.Context) (ogen.ListLibrariesRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	var libs []library.Library
	var err error

	if h.isAdmin(ctx) {
		// Admins see all libraries
		libs, err = h.libraryService.List(ctx)
	} else {
		// Regular users see only libraries they have access to
		libs, err = h.libraryService.ListAccessible(ctx, userID)
	}

	if err != nil {
		h.logger.Error("failed to list libraries", slog.Any("error",err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to list libraries",
		}, nil
	}

	return &ogen.LibraryListResponse{
		Libraries: convertLibraries(libs),
		Total:     int64(len(libs)),
	}, nil
}

// CreateLibrary creates a new media library. Admin only.
// POST /api/v1/libraries
func (h *Handler) CreateLibrary(ctx context.Context, req *ogen.CreateLibraryRequest) (ogen.CreateLibraryRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.CreateLibraryForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Validate library type
	if !library.IsValidLibraryType(string(req.Type)) {
		return &ogen.CreateLibraryBadRequest{
			Code:    400,
			Message: "Invalid library type",
		}, nil
	}

	createReq := library.CreateLibraryRequest{
		Name:    req.Name,
		Type:    string(req.Type),
		Paths:   req.Paths,
		Enabled: true, // Default
	}

	// Apply optional fields
	if req.Enabled.IsSet() {
		createReq.Enabled = req.Enabled.Value
	}
	if req.ScanOnStartup.IsSet() {
		createReq.ScanOnStartup = req.ScanOnStartup.Value
	}
	if req.RealtimeMonitoring.IsSet() {
		createReq.RealtimeMonitoring = req.RealtimeMonitoring.Value
	}
	if req.MetadataProvider.IsSet() {
		provider := req.MetadataProvider.Value
		createReq.MetadataProvider = &provider
	}
	if req.PreferredLanguage.IsSet() {
		createReq.PreferredLanguage = req.PreferredLanguage.Value
	}
	if req.ScannerConfig.IsSet() {
		createReq.ScannerConfig = convertOgenScannerConfigToMap(req.ScannerConfig.Value)
	}

	created, err := h.libraryService.Create(ctx, createReq)
	if err != nil {
		if errors.Is(err, library.ErrLibraryExists) {
			return &ogen.CreateLibraryConflict{
				Code:    409,
				Message: "Library with this name already exists",
			}, nil
		}
		if errors.Is(err, library.ErrInvalidLibraryType) {
			return &ogen.CreateLibraryBadRequest{
				Code:    400,
				Message: "Invalid library type",
			}, nil
		}
		h.logger.Error("failed to create library", slog.Any("error",err))
		return &ogen.CreateLibraryBadRequest{
			Code:    500,
			Message: "Failed to create library",
		}, nil
	}

	return convertLibraryToOgen(created), nil
}

// GetLibrary returns details of a specific library.
// GET /api/v1/libraries/{libraryId}
func (h *Handler) GetLibrary(ctx context.Context, params ogen.GetLibraryParams) (ogen.GetLibraryRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.GetLibraryUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	lib, err := h.libraryService.Get(ctx, params.LibraryId)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.GetLibraryNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		h.logger.Error("failed to get library", slog.Any("error",err))
		return &ogen.GetLibraryNotFound{
			Code:    500,
			Message: "Failed to get library",
		}, nil
	}

	// Check access permission for non-admins
	isAdmin := h.isAdmin(ctx)
	if !isAdmin {
		canAccess, err := h.libraryService.CanAccess(ctx, params.LibraryId, userID, isAdmin)
		if err != nil {
			h.logger.Error("failed to check library access", slog.Any("error",err))
			return &ogen.GetLibraryNotFound{
				Code:    500,
				Message: "Failed to check library access",
			}, nil
		}
		if !canAccess {
			return &ogen.GetLibraryForbidden{
				Code:    403,
				Message: "Access denied to this library",
			}, nil
		}
	}

	return convertLibraryToOgen(lib), nil
}

// UpdateLibrary updates a library's settings. Admin only.
// PUT /api/v1/libraries/{libraryId}
func (h *Handler) UpdateLibrary(ctx context.Context, req *ogen.UpdateLibraryRequest, params ogen.UpdateLibraryParams) (ogen.UpdateLibraryRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.UpdateLibraryForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	update := &library.LibraryUpdate{}

	if req.Name.IsSet() {
		update.Name = &req.Name.Value
	}
	if len(req.Paths) > 0 {
		update.Paths = req.Paths
	}
	if req.Enabled.IsSet() {
		update.Enabled = &req.Enabled.Value
	}
	if req.ScanOnStartup.IsSet() {
		update.ScanOnStartup = &req.ScanOnStartup.Value
	}
	if req.RealtimeMonitoring.IsSet() {
		update.RealtimeMonitoring = &req.RealtimeMonitoring.Value
	}
	if req.MetadataProvider.IsSet() {
		update.MetadataProvider = &req.MetadataProvider.Value
	}
	if req.PreferredLanguage.IsSet() {
		update.PreferredLanguage = &req.PreferredLanguage.Value
	}
	if req.ScannerConfig.IsSet() {
		config := convertOgenUpdateScannerConfigToMap(req.ScannerConfig.Value)
		update.ScannerConfig = config
	}

	updated, err := h.libraryService.Update(ctx, params.LibraryId, update)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.UpdateLibraryNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		if errors.Is(err, library.ErrLibraryExists) {
			return &ogen.UpdateLibraryConflict{
				Code:    409,
				Message: "Library with this name already exists",
			}, nil
		}
		h.logger.Error("failed to update library", slog.Any("error",err))
		return &ogen.UpdateLibraryNotFound{
			Code:    500,
			Message: "Failed to update library",
		}, nil
	}

	return convertLibraryToOgen(updated), nil
}

// DeleteLibrary deletes a library and all its content. Admin only.
// DELETE /api/v1/libraries/{libraryId}
func (h *Handler) DeleteLibrary(ctx context.Context, params ogen.DeleteLibraryParams) (ogen.DeleteLibraryRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.DeleteLibraryForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	err := h.libraryService.Delete(ctx, params.LibraryId)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.DeleteLibraryNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		h.logger.Error("failed to delete library", slog.Any("error",err))
		return &ogen.DeleteLibraryNotFound{
			Code:    500,
			Message: "Failed to delete library",
		}, nil
	}

	return &ogen.DeleteLibraryNoContent{}, nil
}

// ============================================================================
// Library Scan Endpoints
// ============================================================================

// TriggerLibraryScan starts a library scan job. Admin only.
// POST /api/v1/libraries/{libraryId}/scan
func (h *Handler) TriggerLibraryScan(ctx context.Context, req *ogen.TriggerLibraryScanReq, params ogen.TriggerLibraryScanParams) (ogen.TriggerLibraryScanRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.TriggerLibraryScanForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	scanType := string(req.ScanType)
	if !library.IsValidScanType(scanType) {
		return &ogen.TriggerLibraryScanBadRequest{
			Code:    400,
			Message: "Invalid scan type",
		}, nil
	}

	scan, err := h.libraryService.TriggerScan(ctx, params.LibraryId, scanType)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.TriggerLibraryScanNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		if errors.Is(err, library.ErrScanInProgress) {
			return &ogen.TriggerLibraryScanConflict{
				Code:    409,
				Message: "A scan is already in progress for this library",
			}, nil
		}
		if errors.Is(err, library.ErrInvalidScanType) {
			return &ogen.TriggerLibraryScanBadRequest{
				Code:    400,
				Message: "Invalid scan type",
			}, nil
		}
		h.logger.Error("failed to trigger library scan", slog.Any("error",err))
		return &ogen.TriggerLibraryScanBadRequest{
			Code:    500,
			Message: "Failed to trigger library scan",
		}, nil
	}

	return convertLibraryScanToOgen(scan), nil
}

// ListLibraryScans returns the scan history for a library.
// GET /api/v1/libraries/{libraryId}/scans
func (h *Handler) ListLibraryScans(ctx context.Context, params ogen.ListLibraryScansParams) (ogen.ListLibraryScansRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.ListLibraryScansUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	// Check access permission for non-admins
	isAdmin := h.isAdmin(ctx)
	if !isAdmin {
		canAccess, err := h.libraryService.CanAccess(ctx, params.LibraryId, userID, isAdmin)
		if err != nil {
			h.logger.Error("failed to check library access", slog.Any("error",err))
			return &ogen.ListLibraryScansNotFound{
				Code:    500,
				Message: "Failed to check library access",
			}, nil
		}
		if !canAccess {
			return &ogen.ListLibraryScansForbidden{
				Code:    403,
				Message: "Access denied to this library",
			}, nil
		}
	}

	limit := int32(20)
	offset := int32(0)
	if params.Limit.IsSet() {
		l, err := validate.SafeInt32(params.Limit.Value)
		if err != nil {
			h.logger.Error("invalid limit value", slog.Any("error",err))
			return &ogen.ListLibraryScansForbidden{
				Code:    400,
				Message: "Invalid limit parameter",
			}, nil
		}
		limit = l
	}
	if params.Offset.IsSet() {
		o, err := validate.SafeInt32(params.Offset.Value)
		if err != nil {
			h.logger.Error("invalid offset value", slog.Any("error",err))
			return &ogen.ListLibraryScansForbidden{
				Code:    400,
				Message: "Invalid offset parameter",
			}, nil
		}
		offset = o
	}

	scans, total, err := h.libraryService.ListScans(ctx, params.LibraryId, limit, offset)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.ListLibraryScansNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		h.logger.Error("failed to list library scans", slog.Any("error",err))
		return &ogen.ListLibraryScansNotFound{
			Code:    500,
			Message: "Failed to list library scans",
		}, nil
	}

	return &ogen.LibraryScanListResponse{
		Scans: convertLibraryScans(scans),
		Total: total,
	}, nil
}

// ============================================================================
// Library Permission Endpoints
// ============================================================================

// ListLibraryPermissions returns all permissions for a library. Admin only.
// GET /api/v1/libraries/{libraryId}/permissions
func (h *Handler) ListLibraryPermissions(ctx context.Context, params ogen.ListLibraryPermissionsParams) (ogen.ListLibraryPermissionsRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.ListLibraryPermissionsForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Check if library exists
	_, err := h.libraryService.Get(ctx, params.LibraryId)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.ListLibraryPermissionsNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		h.logger.Error("failed to get library", slog.Any("error",err))
		return &ogen.ListLibraryPermissionsNotFound{
			Code:    500,
			Message: "Failed to get library",
		}, nil
	}

	perms, err := h.libraryService.ListPermissions(ctx, params.LibraryId)
	if err != nil {
		h.logger.Error("failed to list library permissions", slog.Any("error",err))
		return &ogen.ListLibraryPermissionsNotFound{
			Code:    500,
			Message: "Failed to list library permissions",
		}, nil
	}

	return &ogen.LibraryPermissionListResponse{
		Permissions: convertLibraryPermissions(perms),
	}, nil
}

// GrantLibraryPermission grants a user permission to access a library. Admin only.
// POST /api/v1/libraries/{libraryId}/permissions
func (h *Handler) GrantLibraryPermission(ctx context.Context, req *ogen.GrantLibraryPermissionReq, params ogen.GrantLibraryPermissionParams) (ogen.GrantLibraryPermissionRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.GrantLibraryPermissionForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	permission := string(req.Permission)
	if !library.IsValidPermission(permission) {
		return &ogen.GrantLibraryPermissionBadRequest{
			Code:    400,
			Message: "Invalid permission type",
		}, nil
	}

	err := h.libraryService.GrantPermission(ctx, params.LibraryId, req.UserId, permission)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.GrantLibraryPermissionNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		if errors.Is(err, library.ErrInvalidPermission) {
			return &ogen.GrantLibraryPermissionBadRequest{
				Code:    400,
				Message: "Invalid permission type",
			}, nil
		}
		h.logger.Error("failed to grant library permission", slog.Any("error",err))
		return &ogen.GrantLibraryPermissionBadRequest{
			Code:    500,
			Message: "Failed to grant library permission",
		}, nil
	}

	// Return the created permission - get it from the service
	perm, _ := h.libraryService.GetPermission(ctx, params.LibraryId, req.UserId, permission)
	if perm == nil {
		return &ogen.LibraryPermission{
			ID:         uuid.Must(uuid.NewV7()),
			LibraryId:  params.LibraryId,
			UserId:     req.UserId,
			Permission: ogen.LibraryPermissionPermission(permission),
			CreatedAt:  time.Now(),
		}, nil
	}
	return convertLibraryPermissionToOgen(perm), nil
}

// RevokeLibraryPermission revokes a user's permission for a library. Admin only.
// DELETE /api/v1/libraries/{libraryId}/permissions/{userId}
func (h *Handler) RevokeLibraryPermission(ctx context.Context, params ogen.RevokeLibraryPermissionParams) (ogen.RevokeLibraryPermissionRes, error) {
	if !h.isAdmin(ctx) {
		return &ogen.RevokeLibraryPermissionForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	permission := string(params.Permission)
	if !library.IsValidPermission(permission) {
		return &ogen.RevokeLibraryPermissionBadRequest{
			Code:    400,
			Message: "Invalid permission type",
		}, nil
	}

	err := h.libraryService.RevokePermission(ctx, params.LibraryId, params.UserId, permission)
	if err != nil {
		if errors.Is(err, library.ErrNotFound) {
			return &ogen.RevokeLibraryPermissionNotFound{
				Code:    404,
				Message: "Library not found",
			}, nil
		}
		if errors.Is(err, library.ErrPermissionNotFound) {
			return &ogen.RevokeLibraryPermissionNotFound{
				Code:    404,
				Message: "Permission not found",
			}, nil
		}
		h.logger.Error("failed to revoke library permission", slog.Any("error",err))
		return &ogen.RevokeLibraryPermissionNotFound{
			Code:    500,
			Message: "Failed to revoke library permission",
		}, nil
	}

	return &ogen.RevokeLibraryPermissionNoContent{}, nil
}

// ============================================================================
// Helper functions
// ============================================================================

func convertLibraries(libs []library.Library) []ogen.Library {
	result := make([]ogen.Library, len(libs))
	for i, lib := range libs {
		result[i] = ogen.Library{
			ID:                 lib.ID,
			Name:               lib.Name,
			Type:               ogen.LibraryType(lib.Type),
			Paths:              lib.Paths,
			Enabled:            lib.Enabled,
			ScanOnStartup:      lib.ScanOnStartup,
			RealtimeMonitoring: lib.RealtimeMonitoring,
			MetadataProvider:   optStringFromPtr(lib.MetadataProvider),
			PreferredLanguage:  optStringFromVal(lib.PreferredLanguage),
			ScannerConfig:      convertScannerConfigToOgen(lib.ScannerConfig),
			CreatedAt:          lib.CreatedAt,
			UpdatedAt:          lib.UpdatedAt,
		}
	}
	return result
}

func convertLibraryToOgen(lib *library.Library) *ogen.Library {
	return &ogen.Library{
		ID:                 lib.ID,
		Name:               lib.Name,
		Type:               ogen.LibraryType(lib.Type),
		Paths:              lib.Paths,
		Enabled:            lib.Enabled,
		ScanOnStartup:      lib.ScanOnStartup,
		RealtimeMonitoring: lib.RealtimeMonitoring,
		MetadataProvider:   optStringFromPtr(lib.MetadataProvider),
		PreferredLanguage:  optStringFromVal(lib.PreferredLanguage),
		ScannerConfig:      convertScannerConfigToOgen(lib.ScannerConfig),
		CreatedAt:          lib.CreatedAt,
		UpdatedAt:          lib.UpdatedAt,
	}
}

func convertLibraryScanToOgen(scan *library.LibraryScan) *ogen.LibraryScan {
	result := &ogen.LibraryScan{
		ID:           scan.ID,
		LibraryId:    scan.LibraryID,
		ScanType:     ogen.LibraryScanScanType(scan.ScanType),
		Status:       ogen.LibraryScanStatus(scan.Status),
		ItemsScanned: ogen.NewOptInt64(int64(scan.ItemsScanned)),
		ItemsAdded:   ogen.NewOptInt64(int64(scan.ItemsAdded)),
		ItemsUpdated: ogen.NewOptInt64(int64(scan.ItemsUpdated)),
		ItemsRemoved: ogen.NewOptInt64(int64(scan.ItemsRemoved)),
		ErrorCount:   ogen.NewOptInt64(int64(scan.ErrorsCount)),
		ErrorMessage: optStringFromPtr(scan.ErrorMessage),
		StartedAt:    optDateTimeFromPtr(scan.StartedAt),
		CompletedAt:  optDateTimeFromPtr(scan.CompletedAt),
		CreatedAt:    scan.CreatedAt,
	}
	return result
}

func convertLibraryScans(scans []library.LibraryScan) []ogen.LibraryScan {
	result := make([]ogen.LibraryScan, len(scans))
	for i, scan := range scans {
		result[i] = ogen.LibraryScan{
			ID:           scan.ID,
			LibraryId:    scan.LibraryID,
			ScanType:     ogen.LibraryScanScanType(scan.ScanType),
			Status:       ogen.LibraryScanStatus(scan.Status),
			ItemsScanned: ogen.NewOptInt64(int64(scan.ItemsScanned)),
			ItemsAdded:   ogen.NewOptInt64(int64(scan.ItemsAdded)),
			ItemsUpdated: ogen.NewOptInt64(int64(scan.ItemsUpdated)),
			ItemsRemoved: ogen.NewOptInt64(int64(scan.ItemsRemoved)),
			ErrorCount:   ogen.NewOptInt64(int64(scan.ErrorsCount)),
			ErrorMessage: optStringFromPtr(scan.ErrorMessage),
			StartedAt:    optDateTimeFromPtr(scan.StartedAt),
			CompletedAt:  optDateTimeFromPtr(scan.CompletedAt),
			CreatedAt:    scan.CreatedAt,
		}
	}
	return result
}

func convertLibraryPermissions(perms []library.Permission) []ogen.LibraryPermission {
	result := make([]ogen.LibraryPermission, len(perms))
	for i, perm := range perms {
		result[i] = ogen.LibraryPermission{
			ID:         perm.ID,
			LibraryId:  perm.LibraryID,
			UserId:     perm.UserID,
			Permission: ogen.LibraryPermissionPermission(perm.Permission),
			CreatedAt:  perm.CreatedAt,
		}
	}
	return result
}

func convertLibraryPermissionToOgen(perm *library.Permission) *ogen.LibraryPermission {
	return &ogen.LibraryPermission{
		ID:         perm.ID,
		LibraryId:  perm.LibraryID,
		UserId:     perm.UserID,
		Permission: ogen.LibraryPermissionPermission(perm.Permission),
		CreatedAt:  perm.CreatedAt,
	}
}

func convertScannerConfigToOgen(config map[string]any) ogen.OptLibraryScannerConfig {
	if config == nil {
		return ogen.OptLibraryScannerConfig{}
	}

	result := make(ogen.LibraryScannerConfig)
	for k, v := range config {
		data, err := json.Marshal(v)
		if err != nil {
			continue
		}
		result[k] = jx.Raw(data)
	}

	return ogen.NewOptLibraryScannerConfig(result)
}

func convertOgenScannerConfigToMap(config ogen.CreateLibraryRequestScannerConfig) map[string]any {
	if config == nil {
		return nil
	}
	result := make(map[string]any)
	for k, v := range config {
		var val any
		if err := json.Unmarshal(v, &val); err == nil {
			result[k] = val
		}
	}
	return result
}

func convertOgenUpdateScannerConfigToMap(config ogen.UpdateLibraryRequestScannerConfig) map[string]any {
	if config == nil {
		return nil
	}
	result := make(map[string]any)
	for k, v := range config {
		var val any
		if err := json.Unmarshal(v, &val); err == nil {
			result[k] = val
		}
	}
	return result
}

func optStringFromPtr(s *string) ogen.OptString {
	if s == nil || *s == "" {
		return ogen.OptString{}
	}
	return ogen.NewOptString(*s)
}

func optStringFromVal(s string) ogen.OptString {
	if s == "" {
		return ogen.OptString{}
	}
	return ogen.NewOptString(s)
}

func optDateTimeFromPtr(t *time.Time) ogen.OptDateTime {
	if t == nil {
		return ogen.OptDateTime{}
	}
	return ogen.NewOptDateTime(*t)
}
