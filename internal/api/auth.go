package api

import (
	"context"
	"errors"
	"log/slog"
	"net/netip"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/user"
)

// Login implements the login operation.
func (h *Handler) Login(ctx context.Context, req *gen.LoginRequest) (gen.LoginRes, error) {
	params := auth.LoginParams{
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: netip.IPv4Unspecified(), // TODO: get from request context
	}

	if req.DeviceName.IsSet() {
		params.DeviceName = ptrString(req.DeviceName.Value)
	}
	if req.DeviceType.IsSet() {
		params.DeviceType = ptrString(req.DeviceType.Value)
	}
	if req.ClientName.IsSet() {
		params.ClientName = ptrString(req.ClientName.Value)
	}
	if req.ClientVersion.IsSet() {
		params.ClientVersion = ptrString(req.ClientVersion.Value)
	}

	result, err := h.authService.Login(ctx, params)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			return &gen.LoginUnauthorized{
				Code:    "invalid_credentials",
				Message: "Invalid username or password",
			}, nil
		}
		if errors.Is(err, user.ErrUserDisabled) {
			return &gen.LoginForbidden{
				Code:    "user_disabled",
				Message: "User account is disabled",
			}, nil
		}
		h.logger.Error("Login failed", slog.String("error", err.Error()))
		return &gen.LoginUnauthorized{
			Code:    "login_failed",
			Message: "Login failed",
		}, nil
	}

	return &gen.LoginResponse{
		Token:   result.Token,
		User:    userToAPI(result.User),
		Session: sessionToAPI(result.Session),
	}, nil
}

// Logout implements the logout operation.
func (h *Handler) Logout(ctx context.Context) (gen.LogoutRes, error) {
	sess, err := requireSession(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	if err := h.sessionService.Deactivate(ctx, sess.ID); err != nil {
		h.logger.Error("Logout failed", slog.String("error", err.Error()))
		return &gen.Error{
			Code:    "logout_failed",
			Message: "Failed to logout",
		}, nil
	}

	return &gen.LogoutNoContent{}, nil
}

// Register implements the register operation.
func (h *Handler) Register(ctx context.Context, req *gen.RegisterRequest) (gen.RegisterRes, error) {
	params := auth.RegisterParams{
		Username: req.Username,
		Password: req.Password,
	}

	if req.Email.IsSet() {
		params.Email = ptrString(req.Email.Value)
	}
	if req.PreferredLanguage.IsSet() {
		params.PreferredLanguage = ptrString(req.PreferredLanguage.Value)
	}

	usr, err := h.authService.Register(ctx, params)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			return &gen.ValidationError{
				Code:    "validation_error",
				Message: "User already exists",
				Errors: []gen.ValidationErrorErrorsItem{
					{Field: "username", Message: "Username is already taken"},
				},
			}, nil
		}
		h.logger.Error("Registration failed", slog.String("error", err.Error()))
		return &gen.Error{
			Code:    "registration_failed",
			Message: "Registration failed",
		}, nil
	}

	return &gen.User{
		ID:         usr.ID,
		Username:   usr.Username,
		IsAdmin:    usr.IsAdmin,
		IsDisabled: usr.IsDisabled,
		CreatedAt:  usr.CreatedAt,
	}, nil
}

// GetCurrentUser implements the getCurrentUser operation.
func (h *Handler) GetCurrentUser(ctx context.Context) (gen.GetCurrentUserRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	sess, err := requireSession(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	return &gen.CurrentUser{
		User:    userToAPI(usr),
		Session: sessionToAPI(sess),
	}, nil
}

// ChangePassword implements the changePassword operation.
func (h *Handler) ChangePassword(ctx context.Context, req *gen.ChangePasswordRequest) (gen.ChangePasswordRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	if err := h.authService.ChangePassword(ctx, usr.ID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			return &gen.Error{
				Code:    "invalid_password",
				Message: "Current password is incorrect",
			}, nil
		}
		h.logger.Error("Password change failed", slog.String("error", err.Error()))
		return &gen.Error{
			Code:    "change_failed",
			Message: "Failed to change password",
		}, nil
	}

	return &gen.ChangePasswordNoContent{}, nil
}

// ListSessions implements the listSessions operation.
func (h *Handler) ListSessions(ctx context.Context) (gen.ListSessionsRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	sessions, err := h.sessionService.ListByUser(ctx, usr.ID)
	if err != nil {
		h.logger.Error("List sessions failed", slog.String("error", err.Error()))
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list sessions",
		}, nil
	}

	result := make(gen.ListSessionsOKApplicationJSON, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, sessionToAPI(&s))
	}

	return &result, nil
}

// RevokeSession implements the revokeSession operation.
func (h *Handler) RevokeSession(ctx context.Context, params gen.RevokeSessionParams) (gen.RevokeSessionRes, error) {
	usr, err := requireUser(ctx)
	if err != nil {
		return &gen.RevokeSessionUnauthorized{
			Code:    "unauthorized",
			Message: "Not authenticated",
		}, nil
	}

	// Get the session to verify ownership
	sess, err := h.sessionService.GetByID(ctx, params.SessionId)
	if err != nil {
		return &gen.RevokeSessionNotFound{
			Code:    "not_found",
			Message: "Session not found",
		}, nil
	}

	// Verify the session belongs to the user (unless admin)
	if sess.UserID != usr.ID && !usr.IsAdmin {
		return &gen.RevokeSessionNotFound{
			Code:    "not_found",
			Message: "Session not found",
		}, nil
	}

	if err := h.sessionService.Deactivate(ctx, params.SessionId); err != nil {
		h.logger.Error("Revoke session failed", slog.String("error", err.Error()))
		return &gen.RevokeSessionNotFound{
			Code:    "revoke_failed",
			Message: "Failed to revoke session",
		}, nil
	}

	return &gen.RevokeSessionNoContent{}, nil
}
