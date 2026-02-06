package apikeys_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeTestAPIKey(id, userID uuid.UUID, name string, scopes []string, isActive bool) db.SharedApiKey {
	now := time.Now()
	return db.SharedApiKey{
		ID:        id,
		UserID:    userID,
		Name:      name,
		KeyHash:   "testhash123",
		KeyPrefix: "rv_12345",
		Scopes:    scopes,
		IsActive:  isActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func setupAPIKeysService(repo apikeys.Repository) *apikeys.Service {
	logger := zap.NewNop()
	return apikeys.NewService(repo, logger, 10, 0)
}

func TestService_CreateKey_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		userID := uuid.Must(uuid.NewV7())
		keyID := uuid.Must(uuid.NewV7())

		mockRepo.On("CountUserAPIKeys", mock.Anything, userID).Return(int64(0), nil)
		mockRepo.On("CreateAPIKey", mock.Anything, mock.AnythingOfType("db.CreateAPIKeyParams")).
			Return(makeTestAPIKey(keyID, userID, "Test Key", []string{"read"}, true), nil)

		req := apikeys.CreateKeyRequest{
			Name:   "Test Key",
			Scopes: []string{"read"},
		}

		resp, err := svc.CreateKey(context.Background(), userID, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Test Key", resp.Key.Name)
		assert.NotEmpty(t, resp.RawKey)
	})

	t.Run("max keys exceeded", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CountUserAPIKeys", mock.Anything, userID).Return(int64(10), nil)

		req := apikeys.CreateKeyRequest{
			Name:   "Test Key",
			Scopes: []string{"read"},
		}

		resp, err := svc.CreateKey(context.Background(), userID, req)

		assert.Nil(t, resp)
		assert.ErrorIs(t, err, apikeys.ErrMaxKeysExceeded)
	})

	t.Run("invalid scope", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CountUserAPIKeys", mock.Anything, userID).Return(int64(0), nil)

		req := apikeys.CreateKeyRequest{
			Name:   "Test Key",
			Scopes: []string{"invalid_scope"},
		}

		resp, err := svc.CreateKey(context.Background(), userID, req)

		assert.Nil(t, resp)
		assert.ErrorIs(t, err, apikeys.ErrInvalidScope)
	})

	t.Run("count error", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("CountUserAPIKeys", mock.Anything, userID).Return(int64(0), errors.New("db error"))

		req := apikeys.CreateKeyRequest{
			Name:   "Test Key",
			Scopes: []string{"read"},
		}

		resp, err := svc.CreateKey(context.Background(), userID, req)

		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}

func TestService_GetKey_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetAPIKey", mock.Anything, keyID).
			Return(makeTestAPIKey(keyID, userID, "Test Key", []string{"read"}, true), nil)

		key, err := svc.GetKey(context.Background(), keyID)

		require.NoError(t, err)
		assert.Equal(t, keyID, key.ID)
		assert.Equal(t, "Test Key", key.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetAPIKey", mock.Anything, keyID).Return(db.SharedApiKey{}, errors.New("not found"))

		key, err := svc.GetKey(context.Background(), keyID)

		assert.Nil(t, key)
		assert.ErrorIs(t, err, apikeys.ErrKeyNotFound)
	})
}

func TestService_ListUserKeys_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		userID := uuid.Must(uuid.NewV7())

		keys := []db.SharedApiKey{
			makeTestAPIKey(uuid.Must(uuid.NewV7()), userID, "Key 1", []string{"read"}, true),
			makeTestAPIKey(uuid.Must(uuid.NewV7()), userID, "Key 2", []string{"write"}, true),
		}

		mockRepo.On("ListUserAPIKeys", mock.Anything, userID).Return(keys, nil)

		result, err := svc.ListUserKeys(context.Background(), userID)

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		userID := uuid.Must(uuid.NewV7())

		mockRepo.On("ListUserAPIKeys", mock.Anything, userID).Return(nil, errors.New("db error"))

		result, err := svc.ListUserKeys(context.Background(), userID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestService_ValidateKey_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("invalid format - too short", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		key, err := svc.ValidateKey(context.Background(), "rv_")

		assert.Nil(t, key)
		assert.ErrorIs(t, err, apikeys.ErrInvalidKeyFormat)
	})

	t.Run("invalid format - wrong prefix", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		key, err := svc.ValidateKey(context.Background(), "xx_"+string(make([]byte, 64)))

		assert.Nil(t, key)
		assert.ErrorIs(t, err, apikeys.ErrInvalidKeyFormat)
	})

	t.Run("key not found", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		// Valid format: rv_ + 64 hex chars
		rawKey := "rv_" + "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

		mockRepo.On("GetAPIKeyByHash", mock.Anything, mock.AnythingOfType("string")).
			Return(db.SharedApiKey{}, errors.New("not found"))

		key, err := svc.ValidateKey(context.Background(), rawKey)

		assert.Nil(t, key)
		assert.ErrorIs(t, err, apikeys.ErrKeyNotFound)
	})

	t.Run("key inactive", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		rawKey := "rv_" + "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

		inactiveKey := makeTestAPIKey(keyID, userID, "Inactive Key", []string{"read"}, false)
		mockRepo.On("GetAPIKeyByHash", mock.Anything, mock.AnythingOfType("string")).
			Return(inactiveKey, nil)

		key, err := svc.ValidateKey(context.Background(), rawKey)

		assert.Nil(t, key)
		assert.ErrorIs(t, err, apikeys.ErrKeyInactive)
	})

	t.Run("key expired", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		rawKey := "rv_" + "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

		expiredKey := makeTestAPIKey(keyID, userID, "Expired Key", []string{"read"}, true)
		pastTime := time.Now().Add(-24 * time.Hour)
		expiredKey.ExpiresAt = pgtype.Timestamptz{Time: pastTime, Valid: true}

		mockRepo.On("GetAPIKeyByHash", mock.Anything, mock.AnythingOfType("string")).
			Return(expiredKey, nil)

		key, err := svc.ValidateKey(context.Background(), rawKey)

		assert.Nil(t, key)
		assert.ErrorIs(t, err, apikeys.ErrKeyExpired)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		rawKey := "rv_" + "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

		validKey := makeTestAPIKey(keyID, userID, "Valid Key", []string{"read"}, true)
		mockRepo.On("GetAPIKeyByHash", mock.Anything, mock.AnythingOfType("string")).
			Return(validKey, nil)
		mockRepo.On("UpdateAPIKeyLastUsed", mock.Anything, keyID).Return(nil).Maybe()

		key, err := svc.ValidateKey(context.Background(), rawKey)

		require.NoError(t, err)
		assert.Equal(t, keyID, key.ID)
		// Give time for async update
		time.Sleep(10 * time.Millisecond)
	})
}

func TestService_RevokeKey_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())

		mockRepo.On("RevokeAPIKey", mock.Anything, keyID).Return(nil)

		err := svc.RevokeKey(context.Background(), keyID)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())

		mockRepo.On("RevokeAPIKey", mock.Anything, keyID).Return(errors.New("db error"))

		err := svc.RevokeKey(context.Background(), keyID)

		assert.Error(t, err)
	})
}

func TestService_CheckScope_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("has scope", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		key := makeTestAPIKey(keyID, userID, "Test Key", []string{"read", "write"}, true)
		mockRepo.On("GetAPIKey", mock.Anything, keyID).Return(key, nil)

		hasScope, err := svc.CheckScope(context.Background(), keyID, "read")

		require.NoError(t, err)
		assert.True(t, hasScope)
	})

	t.Run("has admin scope (grants all)", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		key := makeTestAPIKey(keyID, userID, "Admin Key", []string{"admin"}, true)
		mockRepo.On("GetAPIKey", mock.Anything, keyID).Return(key, nil)

		hasScope, err := svc.CheckScope(context.Background(), keyID, "read")

		require.NoError(t, err)
		assert.True(t, hasScope)
	})

	t.Run("missing scope", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		key := makeTestAPIKey(keyID, userID, "Test Key", []string{"read"}, true)
		mockRepo.On("GetAPIKey", mock.Anything, keyID).Return(key, nil)

		hasScope, err := svc.CheckScope(context.Background(), keyID, "write")

		require.NoError(t, err)
		assert.False(t, hasScope)
	})

	t.Run("key not found", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())

		mockRepo.On("GetAPIKey", mock.Anything, keyID).Return(db.SharedApiKey{}, errors.New("not found"))

		hasScope, err := svc.CheckScope(context.Background(), keyID, "read")

		assert.False(t, hasScope)
		assert.ErrorIs(t, err, apikeys.ErrKeyNotFound)
	})
}

func TestService_UpdateScopes_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		newScopes := []string{"read", "write"}

		mockRepo.On("UpdateAPIKeyScopes", mock.Anything, keyID, newScopes).Return(nil)

		err := svc.UpdateScopes(context.Background(), keyID, newScopes)

		assert.NoError(t, err)
	})

	t.Run("invalid scope", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		keyID := uuid.Must(uuid.NewV7())
		newScopes := []string{"invalid_scope"}

		err := svc.UpdateScopes(context.Background(), keyID, newScopes)

		assert.ErrorIs(t, err, apikeys.ErrInvalidScope)
	})
}

func TestService_CleanupExpiredKeys_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		mockRepo.On("DeleteExpiredAPIKeys", mock.Anything).Return(nil)

		err := svc.CleanupExpiredKeys(context.Background())

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockAPIKeysRepository(t)
		svc := setupAPIKeysService(mockRepo)

		mockRepo.On("DeleteExpiredAPIKeys", mock.Anything).Return(errors.New("db error"))

		err := svc.CleanupExpiredKeys(context.Background())

		assert.Error(t, err)
	})
}
