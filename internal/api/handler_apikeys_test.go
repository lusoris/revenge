package api

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupAPIKeysTestHandler(t *testing.T) (*Handler, testutil.DB, uuid.UUID) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())

	// Set up API keys service
	repo := apikeys.NewRepositoryPg(queries)
	apikeyService := apikeys.NewService(repo, zap.NewNop(), 10, 0)

	// Create test user
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "testuser",
		Email:    "testuser@example.com",
	})

	handler := &Handler{
		logger:        zap.NewNop(),
		apikeyService: apikeyService,
	}

	return handler, testDB, user.ID
}

// ListAPIKeys tests

func TestHandler_ListAPIKeys_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupAPIKeysTestHandler(t)

	ctx := context.Background()

	result, err := handler.ListAPIKeys(ctx)
	require.NoError(t, err)

	errResponse, ok := result.(*ogen.Error)
	require.True(t, ok)
	assert.Equal(t, 401, errResponse.Code)
}

func TestHandler_ListAPIKeys_Empty(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.ListAPIKeys(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.APIKeyListResponse)
	require.True(t, ok)
	assert.Empty(t, response.Keys)
}

func TestHandler_ListAPIKeys_WithKeys(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	// Create test keys
	_, err := handler.apikeyService.CreateKey(context.Background(), userID, apikeys.CreateKeyRequest{
		Name:   "Test Key 1",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	_, err = handler.apikeyService.CreateKey(context.Background(), userID, apikeys.CreateKeyRequest{
		Name:   "Test Key 2",
		Scopes: []string{"write"},
	})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), userID)

	result, err := handler.ListAPIKeys(ctx)
	require.NoError(t, err)

	response, ok := result.(*ogen.APIKeyListResponse)
	require.True(t, ok)
	assert.Len(t, response.Keys, 2)
}

// CreateAPIKey tests

func TestHandler_CreateAPIKey_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupAPIKeysTestHandler(t)

	ctx := context.Background()
	req := &ogen.CreateAPIKeyRequest{
		Name:   "Test Key",
		Scopes: []ogen.CreateAPIKeyRequestScopesItem{"read"},
	}

	result, err := handler.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.CreateAPIKeyUnauthorized)
	require.True(t, ok)
}

func TestHandler_CreateAPIKey_NoScopes(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	req := &ogen.CreateAPIKeyRequest{
		Name:   "Test Key",
		Scopes: []ogen.CreateAPIKeyRequestScopesItem{},
	}

	result, err := handler.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	badRequest, ok := result.(*ogen.CreateAPIKeyBadRequest)
	require.True(t, ok)
	assert.Equal(t, 400, badRequest.Code)
	assert.Contains(t, badRequest.Message, "scope")
}

func TestHandler_CreateAPIKey_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	req := &ogen.CreateAPIKeyRequest{
		Name:   "Test Key",
		Scopes: []ogen.CreateAPIKeyRequestScopesItem{"read"},
	}
	req.Description.SetTo("Test Description")

	result, err := handler.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	response, ok := result.(*ogen.CreateAPIKeyResponse)
	require.True(t, ok)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "Test Key", response.Name)
	assert.NotEmpty(t, response.APIKey)
	assert.True(t, len(response.APIKey) > 10) // Should be rv_<64 hex chars>
	assert.NotEmpty(t, response.KeyPrefix)
}

func TestHandler_CreateAPIKey_WithExpiry(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	expiresAt := time.Now().Add(24 * time.Hour)

	ctx := contextWithUserID(context.Background(), userID)
	req := &ogen.CreateAPIKeyRequest{
		Name:   "Expiring Key",
		Scopes: []ogen.CreateAPIKeyRequestScopesItem{"read"},
	}
	req.ExpiresAt.SetTo(expiresAt)

	result, err := handler.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	response, ok := result.(*ogen.CreateAPIKeyResponse)
	require.True(t, ok)
	assert.NotEmpty(t, response.ID)
}

func TestHandler_CreateAPIKey_InvalidScope(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	req := &ogen.CreateAPIKeyRequest{
		Name:   "Test Key",
		Scopes: []ogen.CreateAPIKeyRequestScopesItem{"invalid_scope"},
	}

	result, err := handler.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	badRequest, ok := result.(*ogen.CreateAPIKeyBadRequest)
	require.True(t, ok)
	assert.Equal(t, 400, badRequest.Code)
}

// GetAPIKey tests

func TestHandler_GetAPIKey_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupAPIKeysTestHandler(t)

	ctx := context.Background()
	params := ogen.GetAPIKeyParams{
		KeyId: uuid.New(),
	}

	result, err := handler.GetAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetAPIKeyUnauthorized)
	require.True(t, ok)
}

func TestHandler_GetAPIKey_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.GetAPIKeyParams{
		KeyId: uuid.New(),
	}

	result, err := handler.GetAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetAPIKeyNotFound)
	require.True(t, ok)
}

func TestHandler_GetAPIKey_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	// Create a key first
	createResp, err := handler.apikeyService.CreateKey(context.Background(), userID, apikeys.CreateKeyRequest{
		Name:   "Test Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.GetAPIKeyParams{
		KeyId: createResp.Key.ID,
	}

	result, err := handler.GetAPIKey(ctx, params)
	require.NoError(t, err)

	response, ok := result.(*ogen.APIKeyInfo)
	require.True(t, ok)
	assert.Equal(t, createResp.Key.ID, response.ID)
	assert.Equal(t, "Test Key", response.Name)
}

func TestHandler_GetAPIKey_NotOwner(t *testing.T) {
	t.Parallel()
	handler, testDB, userID := setupAPIKeysTestHandler(t)

	// Create another user
	otherUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "otheruser",
		Email:    "other@example.com",
	})

	// Create a key for the other user
	createResp, err := handler.apikeyService.CreateKey(context.Background(), otherUser.ID, apikeys.CreateKeyRequest{
		Name:   "Other's Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Try to get it as first user
	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.GetAPIKeyParams{
		KeyId: createResp.Key.ID,
	}

	result, err := handler.GetAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.GetAPIKeyNotFound)
	require.True(t, ok, "Should not be able to access another user's key")
}

// RevokeAPIKey tests

func TestHandler_RevokeAPIKey_NoAuth(t *testing.T) {
	t.Parallel()
	handler, _, _ := setupAPIKeysTestHandler(t)

	ctx := context.Background()
	params := ogen.RevokeAPIKeyParams{
		KeyId: uuid.New(),
	}

	result, err := handler.RevokeAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeAPIKeyUnauthorized)
	require.True(t, ok)
}

func TestHandler_RevokeAPIKey_NotFound(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.RevokeAPIKeyParams{
		KeyId: uuid.New(),
	}

	result, err := handler.RevokeAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeAPIKeyNotFound)
	require.True(t, ok)
}

func TestHandler_RevokeAPIKey_Success(t *testing.T) {
	t.Parallel()
	handler, _, userID := setupAPIKeysTestHandler(t)

	// Create a key first
	createResp, err := handler.apikeyService.CreateKey(context.Background(), userID, apikeys.CreateKeyRequest{
		Name:   "Test Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.RevokeAPIKeyParams{
		KeyId: createResp.Key.ID,
	}

	result, err := handler.RevokeAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeAPIKeyNoContent)
	require.True(t, ok)

	// Verify key is revoked
	key, err := handler.apikeyService.GetKey(context.Background(), createResp.Key.ID)
	require.NoError(t, err)
	assert.False(t, key.IsActive, "Key should be inactive after revocation")
}

func TestHandler_RevokeAPIKey_NotOwner(t *testing.T) {
	t.Parallel()
	handler, testDB, userID := setupAPIKeysTestHandler(t)

	// Create another user
	otherUser := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "otheruser",
		Email:    "other@example.com",
	})

	// Create a key for the other user
	createResp, err := handler.apikeyService.CreateKey(context.Background(), otherUser.ID, apikeys.CreateKeyRequest{
		Name:   "Other's Key",
		Scopes: []string{"read"},
	})
	require.NoError(t, err)

	// Try to revoke it as first user
	ctx := contextWithUserID(context.Background(), userID)
	params := ogen.RevokeAPIKeyParams{
		KeyId: createResp.Key.ID,
	}

	result, err := handler.RevokeAPIKey(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.RevokeAPIKeyNotFound)
	require.True(t, ok, "Should not be able to revoke another user's key")

	// Verify key is still active
	key, err := handler.apikeyService.GetKey(context.Background(), createResp.Key.ID)
	require.NoError(t, err)
	assert.True(t, key.IsActive, "Key should still be active")
}
