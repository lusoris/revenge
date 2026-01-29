//go:build ogen

package api

import (
	"context"
	"errors"
	"log/slog"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/service/user"
)

// ListUsers implements the listUsers operation.
func (h *Handler) ListUsers(ctx context.Context, params gen.ListUsersParams) (gen.ListUsersRes, error) {
	_, err := requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.ListUsersUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.ListUsersForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	// Get pagination params with defaults
	limit := int32(20)
	offset := int32(0)
	if params.Limit.IsSet() {
		limit = int32(params.Limit.Value)
	}
	if params.Offset.IsSet() {
		offset = int32(params.Offset.Value)
	}

	users, err := h.userService.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("List users failed", slog.String("error", err.Error()))
		return &gen.ListUsersForbidden{
			Code:    "list_failed",
			Message: "Failed to list users",
		}, nil
	}

	total, err := h.userService.Count(ctx)
	if err != nil {
		h.logger.Error("Count users failed", slog.String("error", err.Error()))
		total = int64(len(users))
	}

	apiUsers := make([]gen.User, 0, len(users))
	for _, u := range users {
		apiUsers = append(apiUsers, userToAPI(&u))
	}

	return &gen.UserListResponse{
		Users: apiUsers,
		Pagination: gen.PaginationMeta{
			Total:  int(total),
			Limit:  int(limit),
			Offset: int(offset),
		},
	}, nil
}

// CreateUser implements the createUser operation.
func (h *Handler) CreateUser(ctx context.Context, req *gen.UserCreate) (gen.CreateUserRes, error) {
	_, err := requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.CreateUserUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.CreateUserForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	params := user.CreateParams{
		Username:       req.Username,
		Password:       req.Password,
		IsAdmin:        req.IsAdmin.Or(false),
		MaxRatingLevel: int32(req.MaxRatingLevel.Or(100)),
		AdultEnabled:   req.AdultEnabled.Or(false),
	}

	if req.Email.IsSet() {
		params.Email = ptrString(req.Email.Value)
	}
	if req.PreferredLanguage.IsSet() {
		params.PreferredLanguage = ptrString(req.PreferredLanguage.Value)
	}

	usr, err := h.userService.Create(ctx, params)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			return &gen.CreateUserConflict{
				Code:    "user_exists",
				Message: "User already exists",
			}, nil
		}
		h.logger.Error("Create user failed", slog.String("error", err.Error()))
		return &gen.ValidationError{
			Code:    "create_failed",
			Message: "Failed to create user",
		}, nil
	}

	result := userToAPI(usr)
	return &result, nil
}

// GetUser implements the getUser operation.
func (h *Handler) GetUser(ctx context.Context, params gen.GetUserParams) (gen.GetUserRes, error) {
	currentUser, err := requireUser(ctx)
	if err != nil {
		return &gen.GetUserUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Users can only get their own info unless admin
	if params.UserId != currentUser.ID && !currentUser.IsAdmin {
		return &gen.GetUserNotFound{
			Code:    "not_found",
			Message: "User not found",
		}, nil
	}

	usr, err := h.userService.GetByID(ctx, params.UserId)
	if err != nil {
		return &gen.GetUserNotFound{
			Code:    "not_found",
			Message: "User not found",
		}, nil
	}

	result := userToAPI(usr)
	return &result, nil
}

// UpdateUser implements the updateUser operation.
func (h *Handler) UpdateUser(ctx context.Context, req *gen.UserUpdate, params gen.UpdateUserParams) (gen.UpdateUserRes, error) {
	currentUser, err := requireUser(ctx)
	if err != nil {
		return &gen.UpdateUserUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Users can only update their own info (limited fields) unless admin
	isSelf := params.UserId == currentUser.ID
	if !isSelf && !currentUser.IsAdmin {
		return &gen.UpdateUserNotFound{
			Code:    "not_found",
			Message: "User not found",
		}, nil
	}

	updateParams := user.UpdateParams{
		ID: params.UserId,
	}

	// Non-admin users can only update certain fields
	if isSelf && !currentUser.IsAdmin {
		// Regular users can update: email, preferred language
		if req.Email.IsSet() {
			updateParams.Email = ptrString(req.Email.Value)
		}
		if req.PreferredLanguage.IsSet() {
			updateParams.PreferredLanguage = ptrString(req.PreferredLanguage.Value)
		}
	} else {
		// Admins can update all fields
		if req.Username.IsSet() {
			updateParams.Username = ptrString(req.Username.Value)
		}
		if req.Email.IsSet() {
			updateParams.Email = ptrString(req.Email.Value)
		}
		if req.IsAdmin.IsSet() {
			updateParams.IsAdmin = ptrBool(req.IsAdmin.Value)
		}
		if req.IsDisabled.IsSet() {
			updateParams.IsDisabled = ptrBool(req.IsDisabled.Value)
		}
		if req.MaxRatingLevel.IsSet() {
			updateParams.MaxRatingLevel = ptrInt32(req.MaxRatingLevel.Value)
		}
		if req.AdultEnabled.IsSet() {
			updateParams.AdultEnabled = ptrBool(req.AdultEnabled.Value)
		}
		if req.PreferredLanguage.IsSet() {
			updateParams.PreferredLanguage = ptrString(req.PreferredLanguage.Value)
		}
	}

	usr, err := h.userService.Update(ctx, updateParams)
	if err != nil {
		h.logger.Error("Update user failed", slog.String("error", err.Error()))
		return &gen.ValidationError{
			Code:    "update_failed",
			Message: "Failed to update user",
		}, nil
	}

	result := userToAPI(usr)
	return &result, nil
}

// DeleteUser implements the deleteUser operation.
func (h *Handler) DeleteUser(ctx context.Context, params gen.DeleteUserParams) (gen.DeleteUserRes, error) {
	currentUser, err := requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteUserUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteUserForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}

	// Prevent self-deletion
	if params.UserId == currentUser.ID {
		return &gen.DeleteUserForbidden{
			Code:    "forbidden",
			Message: "Cannot delete your own account",
		}, nil
	}

	// Verify user exists
	_, err = h.userService.GetByID(ctx, params.UserId)
	if err != nil {
		return &gen.DeleteUserNotFound{
			Code:    "not_found",
			Message: "User not found",
		}, nil
	}

	if err := h.userService.Delete(ctx, params.UserId); err != nil {
		h.logger.Error("Delete user failed", slog.String("error", err.Error()))
		return &gen.DeleteUserForbidden{
			Code:    "delete_failed",
			Message: "Failed to delete user",
		}, nil
	}

	return &gen.DeleteUserNoContent{}, nil
}
