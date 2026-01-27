package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/api/middleware"
	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/service/genre"
)

// GenreHandler handles genre-related HTTP requests.
type GenreHandler struct {
	genreService *genre.Service
}

// NewGenreHandler creates a new GenreHandler.
func NewGenreHandler(genreService *genre.Service) *GenreHandler {
	return &GenreHandler{genreService: genreService}
}

// RegisterRoutes registers genre routes on the given mux.
func (h *GenreHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Public routes (require auth)
	mux.Handle("GET /Genres", auth.Required(http.HandlerFunc(h.ListGenres)))
	mux.Handle("GET /Genres/{genreId}", auth.Required(http.HandlerFunc(h.GetGenre)))
	mux.Handle("GET /Genres/Search", auth.Required(http.HandlerFunc(h.SearchGenres)))
	mux.Handle("GET /Items/{itemId}/Genres", auth.Required(http.HandlerFunc(h.ListGenresForItem)))

	// Admin-only routes
	mux.Handle("POST /Genres", auth.Required(auth.AdminRequired(http.HandlerFunc(h.CreateGenre))))
	mux.Handle("PUT /Genres/{genreId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.UpdateGenre))))
	mux.Handle("DELETE /Genres/{genreId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.DeleteGenre))))
	mux.Handle("POST /Items/{itemId}/Genres/{genreId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.AssignGenreToItem))))
	mux.Handle("DELETE /Items/{itemId}/Genres/{genreId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.RemoveGenreFromItem))))
}

// GenreResponse represents a genre in API responses.
type GenreResponse struct {
	ID          string            `json:"Id"`
	Name        string            `json:"Name"`
	Slug        string            `json:"Slug"`
	Domain      string            `json:"Domain"`
	Description *string           `json:"Description,omitempty"`
	ParentID    *string           `json:"ParentId,omitempty"`
	ExternalIDs map[string]string `json:"ExternalIds,omitempty"`
	Children    []GenreResponse   `json:"Children,omitempty"`
}

// GenreListResponse wraps a list of genres.
type GenreListResponse struct {
	Items      []GenreResponse `json:"Items"`
	TotalCount int             `json:"TotalRecordCount"`
}

// CreateGenreRequest represents the request body for creating a genre.
type CreateGenreRequest struct {
	Name        string            `json:"Name"`
	Slug        string            `json:"Slug,omitempty"`
	Domain      string            `json:"Domain"`
	Description *string           `json:"Description,omitempty"`
	ParentID    *string           `json:"ParentId,omitempty"`
	ExternalIDs map[string]string `json:"ExternalIds,omitempty"`
}

// UpdateGenreRequest represents the request body for updating a genre.
type UpdateGenreRequest struct {
	Name        *string           `json:"Name,omitempty"`
	Slug        *string           `json:"Slug,omitempty"`
	Description *string           `json:"Description,omitempty"`
	ParentID    *string           `json:"ParentId,omitempty"`
	ExternalIDs map[string]string `json:"ExternalIds,omitempty"`
}

// genreToResponse converts a domain.Genre to GenreResponse.
func genreToResponse(g *domain.Genre) GenreResponse {
	resp := GenreResponse{
		ID:          g.ID.String(),
		Name:        g.Name,
		Slug:        g.Slug,
		Domain:      string(g.Domain),
		Description: g.Description,
		ExternalIDs: g.ExternalIDs,
	}

	if g.ParentID != nil {
		parentID := g.ParentID.String()
		resp.ParentID = &parentID
	}

	if len(g.Children) > 0 {
		resp.Children = make([]GenreResponse, len(g.Children))
		for i, child := range g.Children {
			resp.Children[i] = genreToResponse(child)
		}
	}

	return resp
}

// ListGenres handles GET /Genres
// Query params: domain (movie, tv, music, book, podcast, game), parentId, limit, offset
func (h *GenreHandler) ListGenres(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := domain.ListGenresParams{
		Limit:  50,
		Offset: 0,
	}

	// Parse domain filter
	if domainStr := r.URL.Query().Get("domain"); domainStr != "" {
		d := domain.GenreDomain(domainStr)
		if !d.IsValid() {
			BadRequest(w, "Invalid genre domain")
			return
		}
		params.Domain = &d
	}

	// Parse parent filter
	if parentIDStr := r.URL.Query().Get("parentId"); parentIDStr != "" {
		parentID, err := uuid.Parse(parentIDStr)
		if err != nil {
			BadRequest(w, "Invalid parent ID format")
			return
		}
		params.ParentID = &parentID
	}

	// Parse pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}

	// Check for hierarchy request
	if r.URL.Query().Get("hierarchy") == "true" && params.Domain != nil {
		genres, err := h.genreService.GetGenreHierarchy(ctx, *params.Domain)
		if err != nil {
			InternalError(w, err)
			return
		}

		items := make([]GenreResponse, len(genres))
		for i, g := range genres {
			items[i] = genreToResponse(g)
		}

		JSON(w, http.StatusOK, GenreListResponse{
			Items:      items,
			TotalCount: len(items),
		})
		return
	}

	genres, err := h.genreService.ListGenres(ctx, params)
	if err != nil {
		InternalError(w, err)
		return
	}

	items := make([]GenreResponse, len(genres))
	for i, g := range genres {
		items[i] = genreToResponse(g)
	}

	JSON(w, http.StatusOK, GenreListResponse{
		Items:      items,
		TotalCount: len(items),
	})
}

// GetGenre handles GET /Genres/{genreId}
func (h *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	genreIDStr := r.PathValue("genreId")
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		BadRequest(w, "Invalid genre ID format")
		return
	}

	g, err := h.genreService.GetGenre(ctx, genreID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Genre not found")
			return
		}
		InternalError(w, err)
		return
	}

	JSON(w, http.StatusOK, genreToResponse(g))
}

// SearchGenres handles GET /Genres/Search
// Query params: domain (required), query (required), limit
func (h *GenreHandler) SearchGenres(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	domainStr := r.URL.Query().Get("domain")
	if domainStr == "" {
		BadRequest(w, "Domain parameter is required")
		return
	}

	d := domain.GenreDomain(domainStr)
	if !d.IsValid() {
		BadRequest(w, "Invalid genre domain")
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		BadRequest(w, "Query parameter is required")
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	genres, err := h.genreService.SearchGenres(ctx, d, query, limit)
	if err != nil {
		InternalError(w, err)
		return
	}

	items := make([]GenreResponse, len(genres))
	for i, g := range genres {
		items[i] = genreToResponse(g)
	}

	JSON(w, http.StatusOK, GenreListResponse{
		Items:      items,
		TotalCount: len(items),
	})
}

// ListGenresForItem handles GET /Items/{itemId}/Genres
func (h *GenreHandler) ListGenresForItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemIDStr := r.PathValue("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		BadRequest(w, "Invalid item ID format")
		return
	}

	genres, err := h.genreService.ListGenresForMediaItem(ctx, itemID)
	if err != nil {
		InternalError(w, err)
		return
	}

	items := make([]GenreResponse, len(genres))
	for i, g := range genres {
		items[i] = genreToResponse(g)
	}

	JSON(w, http.StatusOK, GenreListResponse{
		Items:      items,
		TotalCount: len(items),
	})
}

// CreateGenre handles POST /Genres
func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	d := domain.GenreDomain(req.Domain)
	if !d.IsValid() {
		BadRequest(w, "Invalid genre domain")
		return
	}

	params := domain.CreateGenreParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Domain:      d,
		Description: req.Description,
		ExternalIDs: req.ExternalIDs,
	}

	if req.ParentID != nil {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			BadRequest(w, "Invalid parent ID format")
			return
		}
		params.ParentID = &parentID
	}

	g, err := h.genreService.CreateGenre(ctx, params)
	if err != nil {
		InternalError(w, err)
		return
	}

	Created(w, genreToResponse(g))
}

// UpdateGenre handles PUT /Genres/{genreId}
func (h *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	genreIDStr := r.PathValue("genreId")
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		BadRequest(w, "Invalid genre ID format")
		return
	}

	var req UpdateGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	params := domain.UpdateGenreParams{
		ID:          genreID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ExternalIDs: req.ExternalIDs,
	}

	if req.ParentID != nil {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			BadRequest(w, "Invalid parent ID format")
			return
		}
		params.ParentID = &parentID
	}

	g, err := h.genreService.UpdateGenre(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Genre not found")
			return
		}
		InternalError(w, err)
		return
	}

	JSON(w, http.StatusOK, genreToResponse(g))
}

// DeleteGenre handles DELETE /Genres/{genreId}
func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	genreIDStr := r.PathValue("genreId")
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		BadRequest(w, "Invalid genre ID format")
		return
	}

	if err := h.genreService.DeleteGenre(ctx, genreID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Genre not found")
			return
		}
		InternalError(w, err)
		return
	}

	NoContent(w)
}

// AssignGenreToItem handles POST /Items/{itemId}/Genres/{genreId}
func (h *GenreHandler) AssignGenreToItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemIDStr := r.PathValue("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		BadRequest(w, "Invalid item ID format")
		return
	}

	genreIDStr := r.PathValue("genreId")
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		BadRequest(w, "Invalid genre ID format")
		return
	}

	source := r.URL.Query().Get("source")
	if source == "" {
		source = domain.GenreSourceManual
	}

	if err := h.genreService.AssignGenreToMediaItem(ctx, itemID, genreID, source); err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}

// RemoveGenreFromItem handles DELETE /Items/{itemId}/Genres/{genreId}
func (h *GenreHandler) RemoveGenreFromItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemIDStr := r.PathValue("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		BadRequest(w, "Invalid item ID format")
		return
	}

	genreIDStr := r.PathValue("genreId")
	genreID, err := uuid.Parse(genreIDStr)
	if err != nil {
		BadRequest(w, "Invalid genre ID format")
		return
	}

	if err := h.genreService.RemoveGenreFromMediaItem(ctx, itemID, genreID); err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}
