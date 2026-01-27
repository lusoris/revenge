package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/api/middleware"
	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/service/user"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRoutes registers user routes on the given mux.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Public route for getting current user (requires auth)
	mux.Handle("GET /Users/Me", auth.Required(http.HandlerFunc(h.GetCurrentUser)))

	// Protected routes - require authentication
	mux.Handle("GET /Users", auth.Required(http.HandlerFunc(h.ListUsers)))
	mux.Handle("GET /Users/{userId}", auth.Required(http.HandlerFunc(h.GetUser)))

	// Admin-only routes
	mux.Handle("POST /Users/New", auth.Required(auth.AdminRequired(http.HandlerFunc(h.CreateUser))))
	mux.Handle("POST /Users", auth.Required(auth.AdminRequired(http.HandlerFunc(h.CreateUser))))
	mux.Handle("POST /Users/{userId}", auth.Required(http.HandlerFunc(h.UpdateUser)))
	mux.Handle("DELETE /Users/{userId}", auth.Required(auth.AdminRequired(http.HandlerFunc(h.DeleteUser))))
}

// UserResponse represents a user in API responses.
// Matches Jellyfin API UserDto format.
type UserResponse struct {
	ID                        string     `json:"Id"`
	Name                      string     `json:"Name"`
	ServerID                  string     `json:"ServerId,omitempty"`
	HasPassword               bool       `json:"HasPassword"`
	HasConfiguredPassword     bool       `json:"HasConfiguredPassword"`
	HasConfiguredEasyPassword bool       `json:"HasConfiguredEasyPassword"`
	EnableAutoLogin           bool       `json:"EnableAutoLogin"`
	LastLoginDate             *string    `json:"LastLoginDate,omitempty"`
	LastActivityDate          *string    `json:"LastActivityDate,omitempty"`
	Policy                    UserPolicy `json:"Policy"`
}

// UserPolicy represents user permissions.
type UserPolicy struct {
	IsAdministrator bool `json:"IsAdministrator"`
	IsDisabled      bool `json:"IsDisabled"`
}

// userToResponse converts a domain.User to UserResponse.
func userToResponse(u *domain.User) UserResponse {
	resp := UserResponse{
		ID:                    u.ID.String(),
		Name:                  u.Username,
		HasPassword:           u.PasswordHash != nil,
		HasConfiguredPassword: u.PasswordHash != nil,
		Policy: UserPolicy{
			IsAdministrator: u.IsAdmin,
			IsDisabled:      u.IsDisabled,
		},
	}

	if u.LastLoginAt != nil {
		t := u.LastLoginAt.Format("2006-01-02T15:04:05.0000000Z")
		resp.LastLoginDate = &t
	}
	if u.LastActivityAt != nil {
		t := u.LastActivityAt.Format("2006-01-02T15:04:05.0000000Z")
		resp.LastActivityDate = &t
	}

	return resp
}

// GetCurrentUser handles GET /Users/Me
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		Unauthorized(w, "Authentication required")
		return
	}

	user, err := h.userService.GetByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			NotFound(w, "User not found")
			return
		}
		InternalError(w, err)
		return
	}

	OK(w, userToResponse(user))
}

// ListUsers handles GET /Users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("startIndex")

	limit := int32(50)
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil && l > 0 {
			limit = int32(l)
		}
	}

	offset := int32(0)
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && o >= 0 {
			offset = int32(o)
		}
	}

	users, err := h.userService.List(r.Context(), limit, offset)
	if err != nil {
		InternalError(w, err)
		return
	}

	// Convert to response format
	response := make([]UserResponse, len(users))
	for i, u := range users {
		response[i] = userToResponse(u)
	}

	OK(w, response)
}

// GetUser handles GET /Users/{userId}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		BadRequest(w, "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	// Check authorization: users can view themselves, admins can view anyone
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		Unauthorized(w, "Authentication required")
		return
	}

	if claims.UserID != userID && !claims.IsAdmin {
		Forbidden(w, "Cannot view other user's profile")
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			NotFound(w, "User not found")
			return
		}
		InternalError(w, err)
		return
	}

	OK(w, userToResponse(user))
}

// CreateUserRequest represents the create user request body.
type CreateUserRequest struct {
	Name     string `json:"Name"`
	Password string `json:"Password,omitempty"`
}

// CreateUser handles POST /Users/New and POST /Users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if req.Name == "" {
		BadRequest(w, "Name is required")
		return
	}

	newUser, err := h.userService.Create(r.Context(), user.CreateParams{
		Username: req.Name,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateUsername):
			BadRequest(w, "Username already exists")
		case errors.Is(err, domain.ErrDuplicateEmail):
			BadRequest(w, "Email already exists")
		default:
			slog.Error("failed to create user", slog.Any("error", err))
			InternalError(w, err)
		}
		return
	}

	Created(w, userToResponse(newUser))
}

// UpdateUserRequest represents the update user request body.
type UpdateUserRequest struct {
	Name   string      `json:"Name,omitempty"`
	Policy *UserPolicy `json:"Policy,omitempty"`
}

// UpdateUser handles POST /Users/{userId}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		BadRequest(w, "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	// Check authorization
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		Unauthorized(w, "Authentication required")
		return
	}

	// Users can update their own name, admins can update anyone's name and policy
	isOwnProfile := claims.UserID == userID
	if !isOwnProfile && !claims.IsAdmin {
		Forbidden(w, "Cannot update other user's profile")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	params := user.UpdateParams{ID: userID}

	if req.Name != "" {
		params.Username = &req.Name
	}

	// Only admins can change policy
	if req.Policy != nil && claims.IsAdmin {
		params.IsAdmin = &req.Policy.IsAdministrator
		params.IsDisabled = &req.Policy.IsDisabled
	}

	if err := h.userService.Update(r.Context(), params); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			NotFound(w, "User not found")
		case errors.Is(err, domain.ErrDuplicateUsername):
			BadRequest(w, "Username already exists")
		default:
			slog.Error("failed to update user",
				slog.String("user_id", userID.String()),
				slog.Any("error", err))
			InternalError(w, err)
		}
		return
	}

	NoContent(w)
}

// DeleteUser handles DELETE /Users/{userId}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	if userIDStr == "" {
		BadRequest(w, "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	// Prevent self-deletion
	claims := middleware.ClaimsFromContext(r.Context())
	if claims != nil && claims.UserID == userID {
		BadRequest(w, "Cannot delete your own account")
		return
	}

	if err := h.userService.Delete(r.Context(), userID); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			NotFound(w, "User not found")
		default:
			slog.Error("failed to delete user",
				slog.String("user_id", userID.String()),
				slog.Any("error", err))
			InternalError(w, err)
		}
		return
	}

	NoContent(w)
}
