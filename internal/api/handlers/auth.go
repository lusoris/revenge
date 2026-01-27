// Package handlers provides HTTP handlers for the Revenge Go API.
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/netip"
	"strings"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/domain"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService domain.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes registers auth routes on the given mux.
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Public routes
	mux.HandleFunc("POST /Users/AuthenticateByName", h.Login)

	// Protected routes
	mux.Handle("POST /Sessions/Logout", auth.Required(http.HandlerFunc(h.Logout)))
	mux.Handle("POST /Auth/Logout", auth.Required(http.HandlerFunc(h.Logout)))
	mux.Handle("POST /Auth/Refresh", auth.Required(http.HandlerFunc(h.RefreshToken)))

	// Password management (requires auth)
	mux.Handle("POST /Users/{userId}/Password", auth.Required(http.HandlerFunc(h.ChangePassword)))
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"Username"`
	Pw       string `json:"Pw"`
}

// LoginResponse represents the login response body.
// Matches Revenge API AuthenticationResult.
type LoginResponse struct {
	User        UserDTO    `json:"User"`
	AccessToken string     `json:"AccessToken"`
	ServerID    string     `json:"ServerId"`
	SessionInfo SessionDTO `json:"SessionInfo,omitempty"`
}

// UserDTO represents a user in API responses.
// Matches Revenge API UserDto.
type UserDTO struct {
	ID                        string `json:"Id"`
	Name                      string `json:"Name"`
	ServerID                  string `json:"ServerId,omitempty"`
	HasPassword               bool   `json:"HasPassword"`
	HasConfiguredPassword     bool   `json:"HasConfiguredPassword"`
	HasConfiguredEasyPassword bool   `json:"HasConfiguredEasyPassword"`
	EnableAutoLogin           bool   `json:"EnableAutoLogin"`
}

// SessionDTO represents a session in API responses.
type SessionDTO struct {
	ID       string `json:"Id"`
	UserID   string `json:"UserId"`
	UserName string `json:"UserName"`
}

// Login handles POST /Users/AuthenticateByName
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if req.Username == "" {
		BadRequest(w, "Username is required")
		return
	}

	// Parse client info from MediaBrowser header
	deviceID, deviceName, clientName, clientVersion := parseClientInfo(r)

	// Parse IP address
	ipAddr := parseClientIP(r)

	result, err := h.authService.Login(r.Context(), domain.LoginParams{
		Username:      req.Username,
		Password:      req.Pw,
		DeviceID:      deviceID,
		DeviceName:    deviceName,
		ClientName:    clientName,
		ClientVersion: clientVersion,
		IPAddress:     ipAddr,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			Unauthorized(w, "Invalid username or password")
		case errors.Is(err, domain.ErrUserDisabled):
			Forbidden(w, "User account is disabled")
		default:
			slog.Error("login failed", slog.Any("error", err))
			InternalError(w, err)
		}
		return
	}

	// Build response matching Revenge API format
	resp := LoginResponse{
		User: UserDTO{
			ID:                    result.User.ID.String(),
			Name:                  result.User.Username,
			HasPassword:           result.User.PasswordHash != nil,
			HasConfiguredPassword: result.User.PasswordHash != nil,
		},
		AccessToken: result.AccessToken,
		SessionInfo: SessionDTO{
			ID:       result.SessionID.String(),
			UserID:   result.User.ID.String(),
			UserName: result.User.Username,
		},
	}

	OK(w, resp)
}

// Logout handles POST /Sessions/Logout and POST /Auth/Logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractTokenFromRequest(r)
	if token == "" {
		Unauthorized(w, "No token provided")
		return
	}

	if err := h.authService.Logout(r.Context(), token); err != nil {
		slog.Error("logout failed", slog.Any("error", err))
		InternalError(w, err)
		return
	}

	NoContent(w)
}

// RefreshTokenRequest represents the refresh token request body.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshToken handles POST /Auth/Refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		BadRequest(w, "Refresh token is required")
		return
	}

	ipAddr := parseClientIP(r)

	result, err := h.authService.RefreshToken(r.Context(), domain.RefreshParams{
		RefreshToken: req.RefreshToken,
		IPAddress:    ipAddr,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			Unauthorized(w, "Invalid refresh token")
		case errors.Is(err, domain.ErrSessionExpired):
			Unauthorized(w, "Refresh token expired")
		case errors.Is(err, domain.ErrUserDisabled):
			Forbidden(w, "User account is disabled")
		default:
			slog.Error("token refresh failed", slog.Any("error", err))
			InternalError(w, err)
		}
		return
	}

	OK(w, map[string]any{
		"AccessToken":  result.AccessToken,
		"RefreshToken": result.RefreshToken,
		"ExpiresAt":    result.ExpiresAt,
	})
}

// ChangePasswordRequest represents the change password request body.
type ChangePasswordRequest struct {
	CurrentPw string `json:"CurrentPw"`
	NewPw     string `json:"NewPw"`
	// ResetPassword is used by admins to reset without current password
	ResetPassword bool `json:"ResetPassword,omitempty"`
}

// ChangePassword handles POST /Users/{userId}/Password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		BadRequest(w, "User ID is required")
		return
	}

	targetUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	// Get current user from context
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		Unauthorized(w, "Authentication required")
		return
	}

	// Check authorization: users can change their own password, admins can change any
	isOwnPassword := claims.UserID == targetUserID
	if !isOwnPassword && !claims.IsAdmin {
		Forbidden(w, "Cannot change another user's password")
		return
	}

	// Admin reset (no current password required)
	if req.ResetPassword && claims.IsAdmin {
		if req.NewPw == "" {
			BadRequest(w, "New password is required")
			return
		}
		if err := h.authService.ResetPassword(r.Context(), targetUserID, req.NewPw); err != nil {
			slog.Error("password reset failed",
				slog.String("target_user_id", targetUserID.String()),
				slog.Any("error", err))
			InternalError(w, err)
			return
		}
		NoContent(w)
		return
	}

	// Normal password change (requires current password)
	if req.NewPw == "" {
		BadRequest(w, "New password is required")
		return
	}

	if err := h.authService.ChangePassword(r.Context(), targetUserID, req.CurrentPw, req.NewPw); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			Unauthorized(w, "Current password is incorrect")
		case errors.Is(err, domain.ErrUserNotFound):
			NotFound(w, "User not found")
		default:
			slog.Error("password change failed",
				slog.String("user_id", targetUserID.String()),
				slog.Any("error", err))
			InternalError(w, err)
		}
		return
	}

	NoContent(w)
}

// parseClientInfo extracts client info from the MediaBrowser authorization header.
func parseClientInfo(r *http.Request) (deviceID, deviceName, clientName, clientVersion *string) {
	params := parseMediaBrowserHeader(r)

	if v, ok := params["DeviceId"]; ok && v != "" {
		deviceID = &v
	}
	if v, ok := params["Device"]; ok && v != "" {
		deviceName = &v
	}
	if v, ok := params["Client"]; ok && v != "" {
		clientName = &v
	}
	if v, ok := params["Version"]; ok && v != "" {
		clientVersion = &v
	}

	return
}

// parseMediaBrowserHeader parses the MediaBrowser or X-Emby-Authorization header.
func parseMediaBrowserHeader(r *http.Request) map[string]string {
	params := make(map[string]string)

	header := r.Header.Get("Authorization")
	if header == "" {
		header = r.Header.Get("X-Emby-Authorization")
	}
	if header == "" {
		return params
	}

	// Remove "MediaBrowser " prefix
	header = strings.TrimPrefix(header, "MediaBrowser ")

	// Parse key="value" pairs
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "="); idx > 0 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])
			value = strings.Trim(value, "\"")
			params[key] = value
		}
	}

	return params
}

// parseClientIP extracts the client IP address from the request.
func parseClientIP(r *http.Request) *netip.Addr {
	// Check X-Forwarded-For header first (for reverse proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP (client IP)
		if idx := strings.Index(xff, ","); idx > 0 {
			xff = xff[:idx]
		}
		xff = strings.TrimSpace(xff)
		if addr, err := netip.ParseAddr(xff); err == nil {
			return &addr
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if addr, err := netip.ParseAddr(xri); err == nil {
			return &addr
		}
	}

	// Fall back to RemoteAddr
	host := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(host, ":"); idx > 0 {
		// Check if it's IPv6 (has brackets)
		if strings.HasPrefix(host, "[") {
			if bracketIdx := strings.Index(host, "]"); bracketIdx > 0 {
				host = host[1:bracketIdx]
			}
		} else {
			host = host[:idx]
		}
	}

	if addr, err := netip.ParseAddr(host); err == nil {
		return &addr
	}

	return nil
}

// extractTokenFromRequest extracts the JWT token from the request.
func extractTokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
		if strings.HasPrefix(auth, "MediaBrowser ") {
			params := parseMediaBrowserHeader(r)
			if token, ok := params["Token"]; ok {
				return token
			}
		}
	}

	embyAuth := r.Header.Get("X-Emby-Authorization")
	if embyAuth != "" {
		params := make(map[string]string)
		parts := strings.Split(strings.TrimPrefix(embyAuth, "MediaBrowser "), ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if idx := strings.Index(part, "="); idx > 0 {
				key := strings.TrimSpace(part[:idx])
				value := strings.Trim(strings.TrimSpace(part[idx+1:]), "\"")
				params[key] = value
			}
		}
		if token, ok := params["Token"]; ok {
			return token
		}
	}

	if token := r.URL.Query().Get("api_key"); token != "" {
		return token
	}

	return ""
}
