package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
)

// ListSessions lists all active sessions for the authenticated user.
func (h *Handler) ListSessions(ctx context.Context) (ogen.ListSessionsRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	sessions, err := h.sessionService.ListUserSessions(ctx, userID)
	if err != nil {
		h.logger.Error("failed to list sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to list sessions",
		}, nil
	}

	sessionInfos := make([]ogen.SessionInfo, 0, len(sessions))
	for _, s := range sessions {
		sessionInfos = append(sessionInfos, ogen.SessionInfo{
			ID:             s.ID,
			DeviceName:     ogen.NewOptString(stringPtrToValue(s.DeviceName)),
			IPAddress:      stringPtrToValue(s.IPAddress),
			UserAgent:      ogen.NewOptString(stringPtrToValue(s.UserAgent)),
			CreatedAt:      s.CreatedAt,
			LastActivityAt: s.LastActivityAt,
			ExpiresAt:      s.ExpiresAt,
			IsActive:       s.IsActive,
			IsCurrent:      s.IsCurrent,
		})
	}

	return &ogen.SessionListResponse{
		Sessions: sessionInfos,
	}, nil
}

// GetCurrentSession gets information about the current session.
func (h *Handler) GetCurrentSession(ctx context.Context) (ogen.GetCurrentSessionRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.GetCurrentSessionUnauthorized{}, nil
	}

	// Get session ID from context (set by auth middleware)
	sessionID, ok := h.getSessionID(ctx)
	if !ok {
		h.logger.Warn("session ID not found in context", zap.String("user_id", userID.String()))
		return &ogen.GetCurrentSessionUnauthorized{}, nil
	}

	sessions, err := h.sessionService.ListUserSessions(ctx, userID)
	if err != nil {
		h.logger.Error("failed to list sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return &ogen.GetCurrentSessionNotFound{
			Code:    500,
			Message: "Failed to get current session",
		}, nil
	}

	// Find current session
	for _, s := range sessions {
		if s.ID == sessionID {
			return &ogen.SessionInfo{
				ID:             s.ID,
				DeviceName:     ogen.NewOptString(stringPtrToValue(s.DeviceName)),
			IPAddress:      stringPtrToValue(s.IPAddress),
				UserAgent:      ogen.NewOptString(stringPtrToValue(s.UserAgent)),
				CreatedAt:      s.CreatedAt,
				LastActivityAt: s.LastActivityAt,
				ExpiresAt:      s.ExpiresAt,
				IsActive:       s.IsActive,
				IsCurrent:      true,
			}, nil
		}
	}

	return &ogen.GetCurrentSessionNotFound{}, nil
}

// LogoutCurrent revokes the current session (logout).
func (h *Handler) LogoutCurrent(ctx context.Context) (ogen.LogoutCurrentRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	sessionID, ok := h.getSessionID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	if err := h.sessionService.RevokeSession(ctx, sessionID); err != nil {
		h.logger.Error("failed to revoke session",
			zap.String("user_id", userID.String()),
			zap.String("session_id", sessionID.String()),
			zap.Error(err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to revoke session",
		}, nil
	}

	return &ogen.LogoutCurrentNoContent{}, nil
}

// LogoutAll revokes all sessions for the authenticated user (logout everywhere).
func (h *Handler) LogoutAll(ctx context.Context) (ogen.LogoutAllRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	if err := h.sessionService.RevokeAllUserSessions(ctx, userID); err != nil {
		h.logger.Error("failed to revoke all sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to logout all sessions",
		}, nil
	}

	return &ogen.LogoutAllNoContent{}, nil
}

// RefreshSession refreshes access token using refresh token.
func (h *Handler) RefreshSession(ctx context.Context, req *ogen.RefreshSessionRequest) (ogen.RefreshSessionRes, error) {
	accessToken, refreshToken, err := h.sessionService.RefreshSession(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Warn("failed to refresh session",
			zap.Error(err),
		)
		return &ogen.Error{
			Code:    401,
			Message: "Invalid refresh token",
		}, nil
	}

	// Get expiry from config (same as auth service)
	expiresIn := int(h.cfg.Auth.JWTExpiry.Seconds())

	return &ogen.RefreshSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RevokeSession revokes a specific session by ID.
func (h *Handler) RevokeSession(ctx context.Context, params ogen.RevokeSessionParams) (ogen.RevokeSessionRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.RevokeSessionUnauthorized{}, nil
	}

	sessionID := params.SessionId

	// Verify session belongs to user
	sessions, err := h.sessionService.ListUserSessions(ctx, userID)
	if err != nil {
		h.logger.Error("failed to list sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		return &ogen.RevokeSessionBadRequest{
			Code:    500,
			Message: "Failed to verify session ownership",
		}, nil
	}

	found := false
	for _, s := range sessions {
		if s.ID == sessionID {
			found = true
			break
		}
	}

	if !found {
		return &ogen.RevokeSessionNotFound{}, nil
	}

	if err := h.sessionService.RevokeSession(ctx, sessionID); err != nil {
		h.logger.Error("failed to revoke session",
			zap.String("user_id", userID.String()),
			zap.String("session_id", sessionID.String()),
			zap.Error(err),
		)
		return &ogen.RevokeSessionBadRequest{
			Code:    500,
			Message: "Failed to revoke session",
		}, nil
	}

	return &ogen.RevokeSessionNoContent{}, nil
}

// stringPtrToValue converts *string to string (empty if nil).
func stringPtrToValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
