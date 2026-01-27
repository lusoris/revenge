package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/domain"
	"github.com/lusoris/revenge/internal/service/rating"
)

// RatingHandler handles content rating HTTP requests.
type RatingHandler struct {
	ratingService *rating.Service
}

// NewRatingHandler creates a new RatingHandler.
func NewRatingHandler(ratingService *rating.Service) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

// RegisterRoutes registers rating routes on the given mux.
func (h *RatingHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Public routes for rating system info (require auth)
	mux.Handle("GET /Ratings/Systems", auth.Required(http.HandlerFunc(h.ListRatingSystems)))
	mux.Handle("GET /Ratings/Systems/{code}", auth.Required(http.HandlerFunc(h.GetRatingSystem)))
	mux.Handle("GET /Ratings/Systems/{systemId}/Ratings", auth.Required(http.HandlerFunc(h.ListRatingsBySystem)))

	// Content rating routes
	mux.Handle("GET /Items/{itemId}/Ratings", auth.Required(http.HandlerFunc(h.GetContentRatings)))
	mux.Handle("GET /Items/{itemId}/Rating", auth.Required(http.HandlerFunc(h.GetContentDisplayRating)))

	// Admin-only routes for managing content ratings
	mux.Handle("POST /Items/{itemId}/Ratings", auth.Required(auth.AdminRequired(http.HandlerFunc(h.AddContentRating))))
	mux.Handle("DELETE /Items/{itemId}/Ratings/{ratingId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.RemoveContentRating))))
}

// RatingSystemResponse represents a rating system in API responses.
type RatingSystemResponse struct {
	ID           string   `json:"Id"`
	Code         string   `json:"Code"`
	Name         string   `json:"Name"`
	CountryCodes []string `json:"CountryCodes"`
	IsActive     bool     `json:"IsActive"`
}

// RatingResponse represents a rating in API responses.
type RatingResponse struct {
	ID              string  `json:"Id"`
	Code            string  `json:"Code"`
	Name            string  `json:"Name"`
	Description     *string `json:"Description,omitempty"`
	MinAge          *int    `json:"MinAge,omitempty"`
	NormalizedLevel int     `json:"NormalizedLevel"`
	IsAdult         bool    `json:"IsAdult"`
	IconURL         *string `json:"IconUrl,omitempty"`
	SystemCode      string  `json:"SystemCode"`
	SystemName      string  `json:"SystemName"`
}

// ContentRatingResponse represents a content rating in API responses.
type ContentRatingResponse struct {
	ID          string          `json:"Id"`
	ContentID   string          `json:"ContentId"`
	ContentType string          `json:"ContentType"`
	Source      *string         `json:"Source,omitempty"`
	Rating      *RatingResponse `json:"Rating,omitempty"`
}

// ratingSystemToResponse converts a domain.RatingSystem to RatingSystemResponse.
func ratingSystemToResponse(rs *domain.RatingSystem) RatingSystemResponse {
	return RatingSystemResponse{
		ID:           rs.ID.String(),
		Code:         rs.Code,
		Name:         rs.Name,
		CountryCodes: rs.CountryCodes,
		IsActive:     rs.IsActive,
	}
}

// ratingToResponse converts a domain.Rating to RatingResponse.
func ratingToResponse(r *domain.Rating) RatingResponse {
	resp := RatingResponse{
		ID:              r.ID.String(),
		Code:            r.Code,
		Name:            r.Name,
		Description:     r.Description,
		MinAge:          r.MinAge,
		NormalizedLevel: r.NormalizedLevel,
		IsAdult:         r.IsAdult,
		IconURL:         r.IconURL,
	}

	if r.System != nil {
		resp.SystemCode = r.System.Code
		resp.SystemName = r.System.Name
	}

	return resp
}

// contentRatingToResponse converts a domain.ContentRating to ContentRatingResponse.
func contentRatingToResponse(cr *domain.ContentRating) ContentRatingResponse {
	resp := ContentRatingResponse{
		ID:          cr.ID.String(),
		ContentID:   cr.ContentID.String(),
		ContentType: cr.ContentType,
		Source:      cr.Source,
	}

	if cr.Rating != nil {
		r := ratingToResponse(cr.Rating)
		resp.Rating = &r
	}

	return resp
}

// ListRatingSystems handles GET /Ratings/Systems
func (h *RatingHandler) ListRatingSystems(w http.ResponseWriter, r *http.Request) {
	// Optional country filter
	country := r.URL.Query().Get("country")

	var systems []*domain.RatingSystem
	var err error

	if country != "" {
		systems, err = h.ratingService.ListRatingSystemsByCountry(r.Context(), country)
	} else {
		systems, err = h.ratingService.ListRatingSystems(r.Context())
	}

	if err != nil {
		InternalError(w, err)
		return
	}

	response := make([]RatingSystemResponse, len(systems))
	for i, sys := range systems {
		response[i] = ratingSystemToResponse(sys)
	}

	OK(w, response)
}

// GetRatingSystem handles GET /Ratings/Systems/{code}
func (h *RatingHandler) GetRatingSystem(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		BadRequest(w, "Rating system code is required")
		return
	}

	system, err := h.ratingService.GetRatingSystemByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Rating system not found")
			return
		}
		InternalError(w, err)
		return
	}

	OK(w, ratingSystemToResponse(system))
}

// ListRatingsBySystem handles GET /Ratings/Systems/{systemId}/Ratings
func (h *RatingHandler) ListRatingsBySystem(w http.ResponseWriter, r *http.Request) {
	systemID, err := uuid.Parse(r.PathValue("systemId"))
	if err != nil {
		BadRequest(w, "Invalid system ID")
		return
	}

	ratings, err := h.ratingService.ListRatingsBySystem(r.Context(), systemID)
	if err != nil {
		InternalError(w, err)
		return
	}

	response := make([]RatingResponse, len(ratings))
	for i, rating := range ratings {
		response[i] = ratingToResponse(rating)
	}

	OK(w, response)
}

// GetContentRatings handles GET /Items/{itemId}/Ratings
func (h *RatingHandler) GetContentRatings(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	contentType := r.URL.Query().Get("contentType")
	if contentType == "" {
		contentType = "media_item"
	}

	ratings, err := h.ratingService.GetContentRatings(r.Context(), itemID, contentType)
	if err != nil {
		InternalError(w, err)
		return
	}

	response := make([]ContentRatingResponse, len(ratings))
	for i, cr := range ratings {
		response[i] = contentRatingToResponse(cr)
	}

	OK(w, response)
}

// GetContentDisplayRating handles GET /Items/{itemId}/Rating
// Returns the single best rating to display for the content.
func (h *RatingHandler) GetContentDisplayRating(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	contentType := r.URL.Query().Get("contentType")
	if contentType == "" {
		contentType = "media_item"
	}

	// Get user's preferred rating system
	preferredSystem := r.URL.Query().Get("preferredSystem")
	if preferredSystem == "" {
		// Could get from user settings, default to MPAA for now
		preferredSystem = "mpaa"
	}

	rating, err := h.ratingService.GetDisplayRating(r.Context(), itemID, contentType, preferredSystem)
	if err != nil {
		InternalError(w, err)
		return
	}

	if rating == nil {
		// No rating found, return empty response
		OK(w, nil)
		return
	}

	OK(w, contentRatingToResponse(rating))
}

// AddContentRatingRequest represents a request to add a content rating.
type AddContentRatingRequest struct {
	RatingID    string  `json:"RatingId"`
	ContentType string  `json:"ContentType,omitempty"`
	Source      *string `json:"Source,omitempty"`
}

// AddContentRating handles POST /Items/{itemId}/Ratings
func (h *RatingHandler) AddContentRating(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	var req AddContentRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	ratingID, err := uuid.Parse(req.RatingID)
	if err != nil {
		BadRequest(w, "Invalid rating ID")
		return
	}

	contentType := req.ContentType
	if contentType == "" {
		contentType = "media_item"
	}

	cr, err := h.ratingService.AddContentRating(r.Context(), itemID, contentType, ratingID, req.Source)
	if err != nil {
		InternalError(w, err)
		return
	}

	Created(w, contentRatingToResponse(cr))
}

// RemoveContentRating handles DELETE /Items/{itemId}/Ratings/{ratingId}
func (h *RatingHandler) RemoveContentRating(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	ratingID, err := uuid.Parse(r.PathValue("ratingId"))
	if err != nil {
		BadRequest(w, "Invalid rating ID")
		return
	}

	if err := h.ratingService.RemoveContentRating(r.Context(), itemID, ratingID); err != nil {
		InternalError(w, err)
		return
	}

	NoContent(w)
}
