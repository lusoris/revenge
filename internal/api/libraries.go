package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/service/library"
)

// ListLibraries implements the listLibraries operation.
func (h *Handler) ListLibraries(ctx context.Context) (gen.ListLibrariesRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	var libs []library.WithAccess
	var libErr error

	if usr.IsAdmin {
		// Admins can see all libraries
		libs, libErr = h.libraryService.ListAll(ctx)
	} else {
		// Regular users see only accessible libraries
		libs, libErr = h.libraryService.ListForUser(ctx, usr.ID)
	}

	if libErr != nil {
		h.logger.Error("List libraries failed", slog.String("error", libErr.Error()))
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list libraries",
		}, nil
	}

	result := make(gen.ListLibrariesOKApplicationJSON, 0, len(libs))
	for _, l := range libs {
		result = append(result, libraryToAPI(&l.Library))
	}

	return &result, nil
}

// CreateLibrary implements the createLibrary operation.
func (h *Handler) CreateLibrary(ctx context.Context, req *gen.LibraryCreate) (gen.CreateLibraryRes, error) {
	usr, err := requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.CreateLibraryUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.CreateLibraryForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	// Convert owner user ID
	var ownerUserID pgtype.UUID
	if req.IsPrivate.Or(false) {
		ownerUserID = pgtype.UUID{Bytes: usr.ID, Valid: true}
	}

	params := library.CreateParams{
		Name:              req.Name,
		LibraryType:       string(req.Type),
		Paths:             req.Paths,
		ScanEnabled:       req.ScanEnabled.Or(true),
		ScanIntervalHours: int32(req.ScanIntervalHours.Or(24)),
		DownloadImages:    req.DownloadImages.Or(true),
		DownloadNfo:       req.DownloadNfo.Or(false),
		GenerateChapters:  req.GenerateChapters.Or(false),
		IsPrivate:         req.IsPrivate.Or(false),
		OwnerUserID:       ownerUserID,
		SortOrder:         int32(req.SortOrder.Or(0)),
	}

	if req.PreferredLanguage.IsSet() {
		params.PreferredLanguage = ptrString(req.PreferredLanguage.Value)
	}
	if req.Icon.IsSet() {
		params.Icon = ptrString(req.Icon.Value)
	}

	lib, err := h.libraryService.Create(ctx, params)
	if err != nil {
		if errors.Is(err, library.ErrInvalidLibraryType) {
			return &gen.ValidationError{
				Code:    "validation_error",
				Message: "Invalid library type",
				Errors: []gen.ValidationErrorErrorsItem{
					{Field: "type", Message: "Invalid library type"},
				},
			}, nil
		}
		h.logger.Error("Create library failed", slog.String("error", err.Error()))
		return &gen.ValidationError{
			Code:    "create_failed",
			Message: "Failed to create library",
		}, nil
	}

	result := libraryToAPI(lib)
	return &result, nil
}

// GetLibrary implements the getLibrary operation.
func (h *Handler) GetLibrary(ctx context.Context, params gen.GetLibraryParams) (gen.GetLibraryRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.GetLibraryUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	lib, err := h.libraryService.GetByID(ctx, params.LibraryId)
	if err != nil {
		return &gen.GetLibraryNotFound{
			Code:    "not_found",
			Message: "Library not found",
		}, nil
	}

	// Check access for non-admin users
	if !usr.IsAdmin {
		hasAccess, err := h.libraryService.UserCanAccess(ctx, usr.ID, params.LibraryId)
		if err != nil || !hasAccess {
			return &gen.GetLibraryForbidden{
				Code:    "forbidden",
				Message: "No access to this library",
			}, nil
		}
	}

	result := libraryToAPI(lib)
	return &result, nil
}

// UpdateLibrary implements the updateLibrary operation.
func (h *Handler) UpdateLibrary(ctx context.Context, req *gen.LibraryUpdate, params gen.UpdateLibraryParams) (gen.UpdateLibraryRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.UpdateLibraryUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Get library to check ownership
	lib, err := h.libraryService.GetByID(ctx, params.LibraryId)
	if err != nil {
		return &gen.UpdateLibraryNotFound{
			Code:    "not_found",
			Message: "Library not found",
		}, nil
	}

	// Check if user can update (admin or owner of private library)
	isOwner := lib.OwnerUserID.Valid && lib.OwnerUserID.Bytes == usr.ID
	if !usr.IsAdmin && !isOwner {
		return &gen.UpdateLibraryForbidden{
			Code:    "forbidden",
			Message: "No permission to update this library",
		}, nil
	}

	updateParams := library.UpdateParams{
		ID: params.LibraryId,
	}

	if req.Name.IsSet() {
		updateParams.Name = ptrString(req.Name.Value)
	}
	if len(req.Paths) > 0 {
		updateParams.Paths = req.Paths
	}
	if req.ScanEnabled.IsSet() {
		updateParams.ScanEnabled = ptrBool(req.ScanEnabled.Value)
	}
	if req.ScanIntervalHours.IsSet() {
		updateParams.ScanIntervalHours = ptrInt32(req.ScanIntervalHours.Value)
	}
	if req.PreferredLanguage.IsSet() {
		updateParams.PreferredLanguage = ptrString(req.PreferredLanguage.Value)
	}
	if req.DownloadImages.IsSet() {
		updateParams.DownloadImages = ptrBool(req.DownloadImages.Value)
	}
	if req.DownloadNfo.IsSet() {
		updateParams.DownloadNfo = ptrBool(req.DownloadNfo.Value)
	}
	if req.GenerateChapters.IsSet() {
		updateParams.GenerateChapters = ptrBool(req.GenerateChapters.Value)
	}
	if req.IsPrivate.IsSet() {
		updateParams.IsPrivate = ptrBool(req.IsPrivate.Value)
	}
	if req.SortOrder.IsSet() {
		updateParams.SortOrder = ptrInt32(req.SortOrder.Value)
	}
	if req.Icon.IsSet() {
		updateParams.Icon = ptrString(req.Icon.Value)
	}

	updatedLib, err := h.libraryService.Update(ctx, updateParams)
	if err != nil {
		h.logger.Error("Update library failed", slog.String("error", err.Error()))
		return &gen.ValidationError{
			Code:    "update_failed",
			Message: "Failed to update library",
		}, nil
	}

	result := libraryToAPI(updatedLib)
	return &result, nil
}

// DeleteLibrary implements the deleteLibrary operation.
func (h *Handler) DeleteLibrary(ctx context.Context, params gen.DeleteLibraryParams) (gen.DeleteLibraryRes, error) {
	_, err := requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteLibraryUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteLibraryForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	// Verify library exists
	_, err = h.libraryService.GetByID(ctx, params.LibraryId)
	if err != nil {
		return &gen.DeleteLibraryNotFound{
			Code:    "not_found",
			Message: "Library not found",
		}, nil
	}

	if err := h.libraryService.Delete(ctx, params.LibraryId); err != nil {
		h.logger.Error("Delete library failed", slog.String("error", err.Error()))
		return &gen.DeleteLibraryForbidden{
			Code:    "delete_failed",
			Message: "Failed to delete library",
		}, nil
	}

	return &gen.DeleteLibraryNoContent{}, nil
}

// ScanLibrary implements the scanLibrary operation.
func (h *Handler) ScanLibrary(ctx context.Context, req gen.OptScanRequest, params gen.ScanLibraryParams) (gen.ScanLibraryRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.ScanLibraryUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Verify library exists and user has access
	lib, err := h.libraryService.GetByID(ctx, params.LibraryId)
	if err != nil {
		return &gen.ScanLibraryNotFound{
			Code:    "not_found",
			Message: "Library not found",
		}, nil
	}

	// Check access for non-admin users
	if !usr.IsAdmin {
		hasAccess, err := h.libraryService.UserCanAccess(ctx, usr.ID, params.LibraryId)
		if err != nil || !hasAccess {
			return &gen.ScanLibraryForbidden{
				Code:    "forbidden",
				Message: "No access to this library",
			}, nil
		}
	}

	// TODO: Actually queue the scan job when the job system is wired up
	// For now, return a placeholder response
	_ = lib

	return &gen.ScanResponse{
		JobId:   uuid.New(),
		Status:  gen.ScanResponseStatusQueued,
		Message: gen.NewOptString("Scan queued"),
	}, nil
}
