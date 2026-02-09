package api

import (
	"context"
	"errors"
	"strings"

	"log/slog"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/internal/validate"
)

// ============================================================================
// Admin User Management Endpoints (Admin only)
// ============================================================================

// AdminListUsers lists and searches users with optional filters.
// GET /api/v1/admin/users
func (h *Handler) AdminListUsers(ctx context.Context, params ogen.AdminListUsersParams) (ogen.AdminListUsersRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminListUsersUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminListUsersForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	filters := user.UserFilters{
		Limit:  50,
		Offset: 0,
	}

	if params.Query.IsSet() {
		q := strings.TrimSpace(params.Query.Value)
		if q != "" {
			filters.Query = &q
		}
	}
	if params.IsActive.IsSet() {
		v := params.IsActive.Value
		filters.IsActive = &v
	}
	if params.IsAdmin.IsSet() {
		v := params.IsAdmin.Value
		filters.IsAdmin = &v
	}
	if params.Limit.IsSet() {
		limit, err := validate.SafeInt32(params.Limit.Value)
		if err != nil || limit < 1 || limit > 100 {
			return &ogen.AdminListUsersForbidden{Code: 400, Message: "limit must be 1-100"}, nil
		}
		filters.Limit = limit
	}
	if params.Offset.IsSet() {
		offset, err := validate.SafeInt32(params.Offset.Value)
		if err != nil || offset < 0 {
			return &ogen.AdminListUsersForbidden{Code: 400, Message: "offset must be >= 0"}, nil
		}
		filters.Offset = offset
	}

	users, total, err := h.userService.ListUsers(ctx, filters)
	if err != nil {
		h.logger.Error("failed to list users", slog.Any("error", err))
		return &ogen.AdminListUsersForbidden{Code: 500, Message: "Failed to list users"}, nil
	}

	items := make([]ogen.User, 0, len(users))
	for _, u := range users {
		items = append(items, ogen.User{
			ID:            u.ID,
			Username:      u.Username,
			Email:         u.Email,
			DisplayName:   ogen.NewOptString(stringPtrToString(u.DisplayName)),
			AvatarURL:     ogen.NewOptString(stringPtrToString(u.AvatarUrl)),
			Locale:        ogen.NewOptString(stringPtrToString(u.Locale)),
			Timezone:      ogen.NewOptString(stringPtrToString(u.Timezone)),
			QarEnabled:    ogen.NewOptBool(boolPtrToBool(u.QarEnabled)),
			IsActive:      boolPtrToBool(u.IsActive),
			IsAdmin:       ogen.NewOptBool(boolPtrToBool(u.IsAdmin)),
			EmailVerified: ogen.NewOptBool(boolPtrToBool(u.EmailVerified)),
			CreatedAt:     u.CreatedAt,
			LastLoginAt:   ogen.NewOptDateTime(u.LastLoginAt.Time),
		})
	}

	return &ogen.AdminUserListResponse{
		Users: items,
		Total: total,
	}, nil
}

// AdminDeleteUser soft-deletes a user account.
// DELETE /api/v1/admin/users/{userId}
func (h *Handler) AdminDeleteUser(ctx context.Context, params ogen.AdminDeleteUserParams) (ogen.AdminDeleteUserRes, error) {
	adminID, err := h.requireAdmin(ctx)
	if err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminDeleteUserUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminDeleteUserForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	// Prevent self-deletion
	if adminID == params.UserId {
		return &ogen.AdminDeleteUserForbidden{Code: 400, Message: "Cannot delete your own account"}, nil
	}

	// Verify user exists
	_, err = h.userService.GetUser(ctx, params.UserId)
	if err != nil {
		return &ogen.AdminDeleteUserNotFound{Code: 404, Message: "User not found"}, nil
	}

	if err := h.userService.DeleteUser(ctx, params.UserId); err != nil {
		h.logger.Error("failed to delete user",
			slog.String("user_id", params.UserId.String()),
			slog.Any("error", err),
		)
		return &ogen.AdminDeleteUserForbidden{Code: 500, Message: "Failed to delete user"}, nil
	}

	return &ogen.AdminDeleteUserNoContent{}, nil
}
