package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/domain"
	"github.com/lusoris/revenge/internal/service/library"
)

// LibraryHandler handles library-related HTTP requests.
type LibraryHandler struct {
	libraryService *library.Service
}

// NewLibraryHandler creates a new LibraryHandler.
func NewLibraryHandler(libraryService *library.Service) *LibraryHandler {
	return &LibraryHandler{libraryService: libraryService}
}

// RegisterRoutes registers library routes on the given mux.
func (h *LibraryHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Public routes (require auth)
	mux.Handle("GET /Library/VirtualFolders", auth.Required(http.HandlerFunc(h.ListLibraries)))
	mux.Handle("GET /Library/VirtualFolders/{libraryId}", auth.Required(http.HandlerFunc(h.GetLibrary)))

	// Admin-only routes
	mux.Handle("POST /Library/VirtualFolders", auth.Required(auth.AdminRequired(http.HandlerFunc(h.CreateLibrary))))
	mux.Handle("POST /Library/VirtualFolders/{libraryId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.UpdateLibrary))))
	mux.Handle("DELETE /Library/VirtualFolders/{libraryId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.DeleteLibrary))))
	mux.Handle("POST /Library/VirtualFolders/{libraryId}/Refresh", auth.Required(auth.AdminRequired(http.HandlerFunc(h.RefreshLibrary))))
}

// LibraryResponse represents a library in API responses.
// Matches Revenge API VirtualFolderInfo format.
type LibraryResponse struct {
	ID                 string             `json:"ItemId"`
	Name               string             `json:"Name"`
	CollectionType     string             `json:"CollectionType,omitempty"`
	Locations          []string           `json:"Locations"`
	LibraryOptions     LibraryOptions     `json:"LibraryOptions"`
	RefreshStatus      string             `json:"RefreshStatus,omitempty"`
	RefreshProgress    *float64           `json:"RefreshProgress,omitempty"`
	PrimaryImageItemID *string            `json:"PrimaryImageItemId,omitempty"`
}

// LibraryOptions represents library configuration options.
type LibraryOptions struct {
	EnableArchiveMediaFiles        bool     `json:"EnableArchiveMediaFiles"`
	EnablePhotos                   bool     `json:"EnablePhotos"`
	EnableRealtimeMonitor          bool     `json:"EnableRealtimeMonitor"`
	ExtractChapterImagesDuringLibraryScan bool `json:"ExtractChapterImagesDuringLibraryScan"`
	EnableChapterImageExtraction   bool     `json:"EnableChapterImageExtraction"`
	EnableInternetProviders        bool     `json:"EnableInternetProviders"`
	SaveLocalMetadata              bool     `json:"SaveLocalMetadata"`
	EnableAutomaticSeriesGrouping  bool     `json:"EnableAutomaticSeriesGrouping"`
	PreferredMetadataLanguage      string   `json:"PreferredMetadataLanguage,omitempty"`
	MetadataCountryCode            string   `json:"MetadataCountryCode,omitempty"`
	AutomaticRefreshIntervalDays   int      `json:"AutomaticRefreshIntervalDays"`
	PathInfos                      []PathInfo `json:"PathInfos"`
}

// PathInfo represents a library path.
type PathInfo struct {
	Path        string `json:"Path"`
	NetworkPath string `json:"NetworkPath,omitempty"`
}

// libraryToResponse converts a domain.Library to LibraryResponse.
func libraryToResponse(l *domain.Library) LibraryResponse {
	pathInfos := make([]PathInfo, len(l.Paths))
	for i, p := range l.Paths {
		pathInfos[i] = PathInfo{Path: p}
	}

	autoRefreshDays := 0
	if l.ScanIntervalHours != nil {
		autoRefreshDays = *l.ScanIntervalHours / 24
	}

	return LibraryResponse{
		ID:             l.ID.String(),
		Name:           l.Name,
		CollectionType: string(l.Type),
		Locations:      l.Paths,
		LibraryOptions: LibraryOptions{
			EnableRealtimeMonitor:        true,
			EnableInternetProviders:      true,
			AutomaticRefreshIntervalDays: autoRefreshDays,
			PathInfos:                    pathInfos,
		},
	}
}

// ListLibraries handles GET /Library/VirtualFolders
func (h *LibraryHandler) ListLibraries(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		Unauthorized(w, "Authentication required")
		return
	}

	// Non-admin users get filtered libraries (respects adult content settings)
	var libraries []*domain.Library
	var err error

	if claims.IsAdmin {
		libraries, err = h.libraryService.List(r.Context())
	} else {
		libraries, err = h.libraryService.ListForUser(r.Context(), claims.UserID)
	}

	if err != nil {
		InternalError(w, err)
		return
	}

	response := make([]LibraryResponse, len(libraries))
	for i, lib := range libraries {
		response[i] = libraryToResponse(lib)
	}

	OK(w, response)
}

// GetLibrary handles GET /Library/VirtualFolders/{libraryId}
func (h *LibraryHandler) GetLibrary(w http.ResponseWriter, r *http.Request) {
	libraryID, err := uuid.Parse(r.PathValue("libraryId"))
	if err != nil {
		BadRequest(w, "Invalid library ID")
		return
	}

	lib, err := h.libraryService.GetByID(r.Context(), libraryID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Library not found")
			return
		}
		InternalError(w, err)
		return
	}

	OK(w, libraryToResponse(lib))
}

// CreateLibraryRequest represents a request to create a library.
type CreateLibraryRequest struct {
	Name              string   `json:"Name"`
	CollectionType    string   `json:"CollectionType"`
	Paths             []string `json:"Paths"`
	RefreshLibrary    bool     `json:"RefreshLibrary"`
	LibraryOptions    *LibraryOptionsRequest `json:"LibraryOptions,omitempty"`
}

// LibraryOptionsRequest represents library options in requests.
type LibraryOptionsRequest struct {
	EnableRealtimeMonitor        bool `json:"EnableRealtimeMonitor"`
	EnableInternetProviders      bool `json:"EnableInternetProviders"`
	AutomaticRefreshIntervalDays int  `json:"AutomaticRefreshIntervalDays"`
}

// CreateLibrary handles POST /Library/VirtualFolders
func (h *LibraryHandler) CreateLibrary(w http.ResponseWriter, r *http.Request) {
	var req CreateLibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if req.Name == "" {
		BadRequest(w, "Name is required")
		return
	}

	if len(req.Paths) == 0 {
		BadRequest(w, "At least one path is required")
		return
	}

	libType := domain.LibraryType(req.CollectionType)
	if !libType.IsValid() {
		BadRequest(w, "Invalid collection type")
		return
	}

	var scanInterval *int
	if req.LibraryOptions != nil && req.LibraryOptions.AutomaticRefreshIntervalDays > 0 {
		hours := req.LibraryOptions.AutomaticRefreshIntervalDays * 24
		scanInterval = &hours
	}

	lib, err := h.libraryService.Create(r.Context(), domain.CreateLibraryParams{
		Name:              req.Name,
		Type:              libType,
		Paths:             req.Paths,
		IsVisible:         true,
		ScanIntervalHours: scanInterval,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			Error(w, http.StatusConflict, "Library with this name already exists")
			return
		}
		InternalError(w, err)
		return
	}

	Created(w, libraryToResponse(lib))
}

// UpdateLibraryRequest represents a request to update a library.
type UpdateLibraryRequest struct {
	Name           *string  `json:"Name,omitempty"`
	Paths          []string `json:"Paths,omitempty"`
	LibraryOptions *LibraryOptionsRequest `json:"LibraryOptions,omitempty"`
}

// UpdateLibrary handles POST /Library/VirtualFolders/{libraryId}
func (h *LibraryHandler) UpdateLibrary(w http.ResponseWriter, r *http.Request) {
	libraryID, err := uuid.Parse(r.PathValue("libraryId"))
	if err != nil {
		BadRequest(w, "Invalid library ID")
		return
	}

	var req UpdateLibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	params := domain.UpdateLibraryParams{ID: libraryID}

	if req.Name != nil {
		params.Name = req.Name
	}
	if req.Paths != nil {
		params.Paths = req.Paths
	}
	if req.LibraryOptions != nil && req.LibraryOptions.AutomaticRefreshIntervalDays > 0 {
		hours := req.LibraryOptions.AutomaticRefreshIntervalDays * 24
		params.ScanIntervalHours = &hours
	}

	if err := h.libraryService.Update(r.Context(), params); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			Error(w, http.StatusConflict, "Library with this name already exists")
			return
		}
		InternalError(w, err)
		return
	}

	NoContent(w)
}

// DeleteLibrary handles DELETE /Library/VirtualFolders/{libraryId}
func (h *LibraryHandler) DeleteLibrary(w http.ResponseWriter, r *http.Request) {
	libraryID, err := uuid.Parse(r.PathValue("libraryId"))
	if err != nil {
		BadRequest(w, "Invalid library ID")
		return
	}

	if err := h.libraryService.Delete(r.Context(), libraryID); err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}

// RefreshLibrary handles POST /Library/VirtualFolders/{libraryId}/Refresh
func (h *LibraryHandler) RefreshLibrary(w http.ResponseWriter, r *http.Request) {
	libraryID, err := uuid.Parse(r.PathValue("libraryId"))
	if err != nil {
		BadRequest(w, "Invalid library ID")
		return
	}

	// For now, just update the last scan timestamp
	// TODO: Implement actual library scanning
	if err := h.libraryService.UpdateLastScan(r.Context(), libraryID); err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}
