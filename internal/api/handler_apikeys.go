package api

import (
	"context"
	"errors"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"log/slog"
)

// ListAPIKeys lists all API keys for the authenticated user
func (h *Handler) ListAPIKeys(ctx context.Context) (ogen.ListAPIKeysRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	keys, err := h.apikeyService.ListUserKeys(ctx, userID)
	if err != nil {
		h.logger.Error("failed to list API keys",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to list API keys",
		}, nil
	}

	keyInfos := make([]ogen.APIKeyInfo, len(keys))
	for i, key := range keys {
		keyInfos[i] = apiKeyToOgen(key)
	}

	return &ogen.APIKeyListResponse{
		Keys:  keyInfos,
		Total: int64(len(keyInfos)),
	}, nil
}

// CreateAPIKey creates a new API key for the authenticated user
func (h *Handler) CreateAPIKey(ctx context.Context, req *ogen.CreateAPIKeyRequest) (ogen.CreateAPIKeyRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.CreateAPIKeyUnauthorized{}, nil
	}

	// Validate scopes
	if len(req.Scopes) == 0 {
		return &ogen.CreateAPIKeyBadRequest{
			Code:    400,
			Message: "At least one scope is required",
		}, nil
	}

	// Convert ogen scopes to string slice
	scopes := make([]string, len(req.Scopes))
	for i, s := range req.Scopes {
		scopes[i] = string(s)
	}

	// Convert ogen request to service request
	createReq := apikeys.CreateKeyRequest{
		Name:   req.Name,
		Scopes: scopes,
	}

	if req.Description.IsSet() {
		desc := req.Description.Value
		createReq.Description = &desc
	}

	if req.ExpiresAt.IsSet() {
		expires := req.ExpiresAt.Value
		createReq.ExpiresAt = &expires
	}

	// Create key
	resp, err := h.apikeyService.CreateKey(ctx, userID, createReq)
	if err != nil {
		h.logger.Error("failed to create API key",
			slog.String("user_id", userID.String()),
			slog.String("name", req.Name),
			slog.Any("error",err),
		)

		// Check for specific errors
		if errors.Is(err, apikeys.ErrMaxKeysExceeded) {
			return &ogen.CreateAPIKeyBadRequest{
				Code:    400,
				Message: "Maximum number of API keys exceeded",
			}, nil
		}

		if errors.Is(err, apikeys.ErrInvalidScope) {
			return &ogen.CreateAPIKeyBadRequest{
				Code:    400,
				Message: err.Error(),
			}, nil
		}

		return &ogen.CreateAPIKeyBadRequest{
			Code:    500,
			Message: "Failed to create API key",
		}, nil
	}

	return &ogen.CreateAPIKeyResponse{
		ID:        resp.Key.ID,
		Name:      resp.Key.Name,
		KeyPrefix: resp.Key.KeyPrefix,
		Scopes:    resp.Key.Scopes,
		CreatedAt: resp.Key.CreatedAt,
		APIKey:    resp.RawKey,
		Message:   "Store this key securely. It won't be shown again.",
	}, nil
}

// GetAPIKey gets details of a specific API key
func (h *Handler) GetAPIKey(ctx context.Context, params ogen.GetAPIKeyParams) (ogen.GetAPIKeyRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.GetAPIKeyUnauthorized{}, nil
	}

	key, err := h.apikeyService.GetKey(ctx, params.KeyId)
	if err != nil {
		if errors.Is(err, apikeys.ErrKeyNotFound) {
			return &ogen.GetAPIKeyNotFound{}, nil
		}

		h.logger.Error("failed to get API key",
			slog.String("user_id", userID.String()),
			slog.String("key_id", params.KeyId.String()),
			slog.Any("error",err),
		)
		return &ogen.GetAPIKeyNotFound{
			Code:    500,
			Message: "Failed to get API key",
		}, nil
	}

	// Verify ownership
	if key.UserID != userID {
		return &ogen.GetAPIKeyNotFound{}, nil
	}

	return apiKeyToOgenPtr(key), nil
}

// RevokeAPIKey revokes (deactivates) an API key
func (h *Handler) RevokeAPIKey(ctx context.Context, params ogen.RevokeAPIKeyParams) (ogen.RevokeAPIKeyRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.RevokeAPIKeyUnauthorized{}, nil
	}

	// Verify ownership first
	key, err := h.apikeyService.GetKey(ctx, params.KeyId)
	if err != nil {
		if errors.Is(err, apikeys.ErrKeyNotFound) {
			return &ogen.RevokeAPIKeyNotFound{}, nil
		}

		h.logger.Error("failed to get API key for revocation",
			slog.String("user_id", userID.String()),
			slog.String("key_id", params.KeyId.String()),
			slog.Any("error",err),
		)
		return &ogen.RevokeAPIKeyNotFound{
			Code:    500,
			Message: "Failed to revoke API key",
		}, nil
	}

	// Check ownership
	if key.UserID != userID {
		return &ogen.RevokeAPIKeyNotFound{}, nil
	}

	// Revoke key
	if err := h.apikeyService.RevokeKey(ctx, params.KeyId); err != nil {
		h.logger.Error("failed to revoke API key",
			slog.String("user_id", userID.String()),
			slog.String("key_id", params.KeyId.String()),
			slog.Any("error",err),
		)
		return &ogen.RevokeAPIKeyNotFound{
			Code:    500,
			Message: "Failed to revoke API key",
		}, nil
	}

	return &ogen.RevokeAPIKeyNoContent{}, nil
}

// ============================================================================
// Helpers
// ============================================================================

// apiKeyToOgen converts service APIKey to ogen APIKeyInfo
func apiKeyToOgen(key apikeys.APIKey) ogen.APIKeyInfo {
	// Convert string scopes to ogen scope items
	scopes := make([]ogen.APIKeyInfoScopesItem, len(key.Scopes))
	for i, s := range key.Scopes {
		scopes[i] = ogen.APIKeyInfoScopesItem(s)
	}

	info := ogen.APIKeyInfo{
		ID:        key.ID,
		UserID:    key.UserID,
		Name:      key.Name,
		KeyPrefix: key.KeyPrefix,
		Scopes:    scopes,
		IsActive:  key.IsActive,
		CreatedAt: key.CreatedAt,
		UpdatedAt: key.UpdatedAt,
	}

	if key.Description != nil {
		info.Description.SetTo(*key.Description)
	}

	if key.ExpiresAt != nil {
		info.ExpiresAt.SetTo(*key.ExpiresAt)
	}

	if key.LastUsedAt != nil {
		info.LastUsedAt.SetTo(*key.LastUsedAt)
	}

	return info
}

// apiKeyToOgenPtr converts service APIKey pointer to ogen APIKeyInfo pointer
func apiKeyToOgenPtr(key *apikeys.APIKey) *ogen.APIKeyInfo {
	if key == nil {
		return nil
	}
	info := apiKeyToOgen(*key)
	return &info
}
