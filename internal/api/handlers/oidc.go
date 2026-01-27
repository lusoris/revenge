// Package handlers provides HTTP handlers for the Jellyfin Go API.
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/api/middleware"
	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/service/oidc"
)

// OIDCHandler handles OIDC-related HTTP requests.
type OIDCHandler struct {
	oidcService *oidc.Service
}

// NewOIDCHandler creates a new OIDCHandler.
func NewOIDCHandler(oidcService *oidc.Service) *OIDCHandler {
	return &OIDCHandler{oidcService: oidcService}
}

// RegisterRoutes registers OIDC routes on the given mux.
func (h *OIDCHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Public routes
	mux.HandleFunc("GET /Auth/OIDC/Providers", h.ListProviders)
	mux.HandleFunc("GET /Auth/OIDC/Authorize/{providerId}", h.Authorize)
	mux.HandleFunc("GET /Auth/OIDC/Callback", h.Callback)
	mux.HandleFunc("POST /Auth/OIDC/Callback", h.Callback)

	// Protected routes - manage user's OIDC links
	mux.Handle("GET /Users/{userId}/OIDC/Links", auth.Required(http.HandlerFunc(h.GetUserLinks)))
	mux.Handle("DELETE /Users/{userId}/OIDC/Links/{linkId}", auth.Required(http.HandlerFunc(h.UnlinkProvider)))
}

// OIDCProviderResponse represents an OIDC provider in API responses.
type OIDCProviderResponse struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
}

// OIDCProvidersResponse represents the list of OIDC providers.
type OIDCProvidersResponse struct {
	Providers []OIDCProviderResponse `json:"Providers"`
}

// ListProviders returns all enabled OIDC providers.
// GET /Auth/OIDC/Providers
func (h *OIDCHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.oidcService.GetProviders(r.Context())
	if err != nil {
		slog.Error("failed to list OIDC providers", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := OIDCProvidersResponse{
		Providers: make([]OIDCProviderResponse, len(providers)),
	}

	for i, p := range providers {
		response.Providers[i] = OIDCProviderResponse{
			ID:          p.ID.String(),
			Name:        p.Name,
			DisplayName: p.DisplayName,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode response", slog.Any("error", err))
	}
}

// AuthorizeResponse contains the authorization URL.
type AuthorizeResponse struct {
	AuthorizationURL string `json:"AuthorizationUrl"`
}

// Authorize initiates the OIDC authorization flow.
// GET /Auth/OIDC/Authorize/{providerId}?redirect_uri=...
func (h *OIDCHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	providerIDStr := r.PathValue("providerId")
	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		http.Error(w, "Invalid provider ID", http.StatusBadRequest)
		return
	}

	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		// Default to the callback endpoint
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		redirectURI = scheme + "://" + r.Host + "/Auth/OIDC/Callback"
	}

	authURL, err := h.oidcService.GetAuthorizationURL(r.Context(), providerID, redirectURI)
	if err != nil {
		if errors.Is(err, domain.ErrOIDCProviderNotFound) {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to generate authorization URL",
			slog.String("provider_id", providerIDStr),
			slog.Any("error", err))
		http.Error(w, "Failed to initiate OIDC flow", http.StatusInternalServerError)
		return
	}

	// Option 1: Return JSON with URL (for API clients)
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(AuthorizeResponse{AuthorizationURL: authURL}); err != nil {
			slog.Error("failed to encode response", slog.Any("error", err))
		}
		return
	}

	// Option 2: Redirect directly (for browser clients)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// CallbackRequest represents the OIDC callback parameters.
type CallbackRequest struct {
	Code          string  `json:"code"`
	State         string  `json:"state"`
	DeviceID      *string `json:"deviceId,omitempty"`
	DeviceName    *string `json:"deviceName,omitempty"`
	ClientName    *string `json:"clientName,omitempty"`
	ClientVersion *string `json:"clientVersion,omitempty"`
}

// Callback handles the OIDC callback from the identity provider.
// GET/POST /Auth/OIDC/Callback
func (h *OIDCHandler) Callback(w http.ResponseWriter, r *http.Request) {
	var code, state string

	if r.Method == http.MethodGet {
		code = r.URL.Query().Get("code")
		state = r.URL.Query().Get("state")
	} else {
		var req CallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		code = req.Code
		state = req.State
	}

	if code == "" || state == "" {
		// Check for error response from provider
		errorCode := r.URL.Query().Get("error")
		errorDesc := r.URL.Query().Get("error_description")
		if errorCode != "" {
			slog.Warn("OIDC provider returned error",
				slog.String("error", errorCode),
				slog.String("description", errorDesc))
			http.Error(w, "OIDC authentication failed: "+errorDesc, http.StatusBadRequest)
			return
		}
		http.Error(w, "Missing code or state parameter", http.StatusBadRequest)
		return
	}

	// Determine redirect URI (must match the one used in authorize)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	redirectURI := scheme + "://" + r.Host + "/Auth/OIDC/Callback"

	result, err := h.oidcService.HandleCallback(r.Context(), oidc.CallbackParams{
		Code:        code,
		State:       state,
		RedirectURI: redirectURI,
	})
	if err != nil {
		slog.Error("OIDC callback failed", slog.Any("error", err))
		http.Error(w, "OIDC authentication failed", http.StatusUnauthorized)
		return
	}

	// Return authentication result
	response := LoginResponse{
		User: UserDTO{
			ID:                    result.User.ID.String(),
			Name:                  result.User.Username,
			HasPassword:           result.User.PasswordHash != nil,
			HasConfiguredPassword: result.User.PasswordHash != nil,
		},
		AccessToken: result.AccessToken,
		ServerID:    "", // TODO: Get from config
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode response", slog.Any("error", err))
	}
}

// OIDCLinkResponse represents an OIDC link in API responses.
type OIDCLinkResponse struct {
	ID              string  `json:"Id"`
	ProviderID      string  `json:"ProviderId"`
	ProviderName    string  `json:"ProviderName"`
	ProviderDisplay string  `json:"ProviderDisplayName"`
	Subject         string  `json:"Subject"`
	Email           *string `json:"Email,omitempty"`
	LastLoginAt     *string `json:"LastLoginAt,omitempty"`
}

// OIDCLinksResponse represents the list of OIDC links for a user.
type OIDCLinksResponse struct {
	Links []OIDCLinkResponse `json:"Links"`
}

// GetUserLinks returns all OIDC links for a user.
// GET /Users/{userId}/OIDC/Links
func (h *OIDCHandler) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check authorization: user can only view their own links, or admin can view any
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims.UserID != userID && !claims.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	links, err := h.oidcService.GetUserLinks(r.Context(), userID)
	if err != nil {
		slog.Error("failed to get user OIDC links",
			slog.String("user_id", userIDStr),
			slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := OIDCLinksResponse{
		Links: make([]OIDCLinkResponse, len(links)),
	}

	for i, link := range links {
		var lastLogin *string
		if link.LastLoginAt != nil {
			t := link.LastLoginAt.Format("2006-01-02T15:04:05.0000000Z")
			lastLogin = &t
		}

		response.Links[i] = OIDCLinkResponse{
			ID:              link.ID.String(),
			ProviderID:      link.ProviderID.String(),
			ProviderName:    link.ProviderName,
			ProviderDisplay: link.ProviderDisplayName,
			Subject:         link.Subject,
			Email:           link.Email,
			LastLoginAt:     lastLogin,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode response", slog.Any("error", err))
	}
}

// UnlinkProvider removes an OIDC link from a user.
// DELETE /Users/{userId}/OIDC/Links/{linkId}
func (h *OIDCHandler) UnlinkProvider(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	linkIDStr := r.PathValue("linkId")
	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}

	// Check authorization
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims.UserID != userID && !claims.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.oidcService.UnlinkUser(r.Context(), userID, linkID); err != nil {
		if errors.Is(err, domain.ErrOIDCUserLinkNotFound) {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to unlink OIDC provider",
			slog.String("user_id", userIDStr),
			slog.String("link_id", linkIDStr),
			slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
