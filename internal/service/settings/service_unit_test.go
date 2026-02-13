package settings_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/settings"
)

// ============================================================================
// Test Helpers
// ============================================================================

func makeServerSetting(key, dataType string, value []byte) *db.SharedServerSetting {
	return &db.SharedServerSetting{
		Key:      key,
		Value:    value,
		DataType: dataType,
	}
}

func makeUserSetting(userID uuid.UUID, key, dataType string, value []byte) *db.SharedUserSetting {
	return &db.SharedUserSetting{
		UserID:   userID,
		Key:      key,
		Value:    value,
		DataType: dataType,
	}
}

// ============================================================================
// Server Settings Unit Tests
// ============================================================================

func TestUnit_GetServerSetting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)
		dbSetting := makeServerSetting("test.key", "string", []byte(`"test_value"`))

		repo.EXPECT().GetServerSetting(ctx, "test.key").Return(dbSetting, nil)

		svc := settings.NewService(repo)
		result, err := svc.GetServerSetting(ctx, "test.key")

		require.NoError(t, err)
		assert.Equal(t, "test.key", result.Key)
		assert.Equal(t, "test_value", result.Value)
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().GetServerSetting(ctx, "missing.key").Return(nil, pgx.ErrNoRows)

		svc := settings.NewService(repo)
		result, err := svc.GetServerSetting(ctx, "missing.key")

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().GetServerSetting(ctx, "error.key").Return(nil, errors.New("db error"))

		svc := settings.NewService(repo)
		result, err := svc.GetServerSetting(ctx, "error.key")

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestUnit_ListServerSettings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)
		dbSettings := []db.SharedServerSetting{
			*makeServerSetting("key1", "string", []byte(`"value1"`)),
			*makeServerSetting("key2", "number", []byte(`42`)),
		}

		repo.EXPECT().ListServerSettings(ctx).Return(dbSettings, nil)

		svc := settings.NewService(repo)
		result, err := svc.ListServerSettings(ctx)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "key1", result[0].Key)
		assert.Equal(t, "key2", result[1].Key)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().ListServerSettings(ctx).Return(nil, errors.New("db error"))

		svc := settings.NewService(repo)
		result, err := svc.ListServerSettings(ctx)

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestUnit_ListServerSettingsByCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockSettingsRepository(t)
	dbSettings := []db.SharedServerSetting{
		*makeServerSetting("cat1.key1", "string", []byte(`"value1"`)),
	}

	repo.EXPECT().ListServerSettingsByCategory(ctx, "cat1").Return(dbSettings, nil)

	svc := settings.NewService(repo)
	result, err := svc.ListServerSettingsByCategory(ctx, "cat1")

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestUnit_ListPublicServerSettings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockSettingsRepository(t)
	dbSettings := []db.SharedServerSetting{
		*makeServerSetting("public.key1", "string", []byte(`"value1"`)),
	}

	repo.EXPECT().ListPublicServerSettings(ctx).Return(dbSettings, nil)

	svc := settings.NewService(repo)
	result, err := svc.ListPublicServerSettings(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestUnit_SetServerSetting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("create new setting", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)
		dbSetting := makeServerSetting("new.key", "string", []byte(`"new_value"`))

		repo.EXPECT().GetServerSetting(ctx, "new.key").Return(nil, pgx.ErrNoRows)
		repo.EXPECT().UpsertServerSetting(ctx, mock.AnythingOfType("db.UpsertServerSettingParams")).Return(dbSetting, nil)

		svc := settings.NewService(repo)
		result, err := svc.SetServerSetting(ctx, "new.key", "new_value", adminID)

		require.NoError(t, err)
		assert.Equal(t, "new.key", result.Key)
	})

	t.Run("update existing setting", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)
		existingSetting := makeServerSetting("existing.key", "string", []byte(`"old_value"`))
		updatedSetting := makeServerSetting("existing.key", "string", []byte(`"new_value"`))

		repo.EXPECT().GetServerSetting(ctx, "existing.key").Return(existingSetting, nil)
		repo.EXPECT().UpdateServerSetting(ctx, mock.AnythingOfType("db.UpdateServerSettingParams")).Return(updatedSetting, nil)

		svc := settings.NewService(repo)
		result, err := svc.SetServerSetting(ctx, "existing.key", "new_value", adminID)

		require.NoError(t, err)
		assert.Equal(t, "existing.key", result.Key)
	})
}

func TestUnit_DeleteServerSetting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().DeleteServerSetting(ctx, "delete.key").Return(nil)

		svc := settings.NewService(repo)
		err := svc.DeleteServerSetting(ctx, "delete.key")

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().DeleteServerSetting(ctx, "error.key").Return(errors.New("db error"))

		svc := settings.NewService(repo)
		err := svc.DeleteServerSetting(ctx, "error.key")

		require.Error(t, err)
	})
}

// ============================================================================
// User Settings Unit Tests
// ============================================================================

func TestUnit_GetUserSetting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("success", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)
		dbSetting := makeUserSetting(userID, "user.key", "string", []byte(`"user_value"`))

		repo.EXPECT().GetUserSetting(ctx, userID, "user.key").Return(dbSetting, nil)

		svc := settings.NewService(repo)
		result, err := svc.GetUserSetting(ctx, userID, "user.key")

		require.NoError(t, err)
		assert.Equal(t, "user.key", result.Key)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "user_value", result.Value)
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().GetUserSetting(ctx, userID, "missing.key").Return(nil, pgx.ErrNoRows)

		svc := settings.NewService(repo)
		result, err := svc.GetUserSetting(ctx, userID, "missing.key")

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestUnit_ListUserSettings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	dbSettings := []db.SharedUserSetting{
		*makeUserSetting(userID, "user.key1", "string", []byte(`"value1"`)),
		*makeUserSetting(userID, "user.key2", "number", []byte(`42`)),
	}

	repo.EXPECT().ListUserSettings(ctx, userID).Return(dbSettings, nil)

	svc := settings.NewService(repo)
	result, err := svc.ListUserSettings(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUnit_ListUserSettingsByCategory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	dbSettings := []db.SharedUserSetting{
		*makeUserSetting(userID, "cat1.key1", "string", []byte(`"value1"`)),
	}

	repo.EXPECT().ListUserSettingsByCategory(ctx, userID, "cat1").Return(dbSettings, nil)

	svc := settings.NewService(repo)
	result, err := svc.ListUserSettingsByCategory(ctx, userID, "cat1")

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestUnit_SetUserSetting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	dbSetting := makeUserSetting(userID, "user.key", "string", []byte(`"new_value"`))

	repo.EXPECT().UpsertUserSetting(ctx, mock.AnythingOfType("db.UpsertUserSettingParams")).Return(dbSetting, nil)

	svc := settings.NewService(repo)
	result, err := svc.SetUserSetting(ctx, userID, "user.key", "new_value")

	require.NoError(t, err)
	assert.Equal(t, "user.key", result.Key)
	assert.Equal(t, userID, result.UserID)
}

func TestUnit_SetUserSettingsBulk(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("success", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		// Return properly formed settings
		repo.EXPECT().UpsertUserSetting(ctx, mock.AnythingOfType("db.UpsertUserSettingParams")).
			Return(makeUserSetting(userID, "bulk.key", "string", []byte(`"value"`)), nil).Times(3)

		svc := settings.NewService(repo)
		settingsMap := map[string]any{
			"bulk.key1": "value1",
			"bulk.key2": 42,
			"bulk.key3": true,
		}
		err := svc.SetUserSettingsBulk(ctx, userID, settingsMap)

		require.NoError(t, err)
	})

	t.Run("partial failure", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		// First call succeeds, second fails
		repo.EXPECT().UpsertUserSetting(ctx, mock.AnythingOfType("db.UpsertUserSettingParams")).
			Return(makeUserSetting(userID, "bulk.key1", "string", []byte(`"value1"`)), nil).Once()
		repo.EXPECT().UpsertUserSetting(ctx, mock.AnythingOfType("db.UpsertUserSettingParams")).
			Return(nil, errors.New("db error")).Once()

		svc := settings.NewService(repo)
		settingsMap := map[string]any{
			"bulk.key1": "value1",
			"bulk.key2": "value2",
		}
		err := svc.SetUserSettingsBulk(ctx, userID, settingsMap)

		require.Error(t, err)
	})
}

func TestUnit_DeleteUserSetting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("success", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().DeleteUserSetting(ctx, userID, "delete.key").Return(nil)

		svc := settings.NewService(repo)
		err := svc.DeleteUserSetting(ctx, userID, "delete.key")

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockSettingsRepository(t)

		repo.EXPECT().DeleteUserSetting(ctx, userID, "error.key").Return(errors.New("db error"))

		svc := settings.NewService(repo)
		err := svc.DeleteUserSetting(ctx, userID, "error.key")

		require.Error(t, err)
	})
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestUnit_MarshalValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{"string", "test"},
		{"number", 42},
		{"float", 3.14},
		{"bool", true},
		{"array", []string{"a", "b", "c"}},
		{"map", map[string]any{"key": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := settings.MarshalValue(tt.value)
			require.NoError(t, err)
			assert.NotEmpty(t, result)

			// Verify we can unmarshal it back
			var unmarshaled any
			err = json.Unmarshal(result, &unmarshaled)
			require.NoError(t, err)
		})
	}
}

func TestUnit_UnmarshalValue(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		data := json.RawMessage(`"test"`)
		var result string
		err := settings.UnmarshalValue(data, &result)
		require.NoError(t, err)
		assert.Equal(t, "test", result)
	})

	t.Run("number", func(t *testing.T) {
		data := json.RawMessage(`42`)
		var result int
		err := settings.UnmarshalValue(data, &result)
		require.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("bool", func(t *testing.T) {
		data := json.RawMessage(`true`)
		var result bool
		err := settings.UnmarshalValue(data, &result)
		require.NoError(t, err)
		assert.True(t, result)
	})
}

// ============================================================================
// Additional Error Path Tests
// ============================================================================

func TestUnit_GetUserSetting_RepositoryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().GetUserSetting(ctx, userID, "error.key").Return(nil, errors.New("database connection error"))

	svc := settings.NewService(repo)
	result, err := svc.GetUserSetting(ctx, userID, "error.key")

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user setting")
}

func TestUnit_ListUserSettings_RepositoryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().ListUserSettings(ctx, userID).Return(nil, errors.New("database error"))

	svc := settings.NewService(repo)
	result, err := svc.ListUserSettings(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list user settings")
}

func TestUnit_ListUserSettingsByCategory_RepositoryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().ListUserSettingsByCategory(ctx, userID, "cat1").Return(nil, errors.New("database error"))

	svc := settings.NewService(repo)
	result, err := svc.ListUserSettingsByCategory(ctx, userID, "cat1")

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list user settings by category")
}

func TestUnit_ListServerSettingsByCategory_RepositoryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().ListServerSettingsByCategory(ctx, "cat1").Return(nil, errors.New("database error"))

	svc := settings.NewService(repo)
	result, err := svc.ListServerSettingsByCategory(ctx, "cat1")

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list server settings by category")
}

func TestUnit_ListPublicServerSettings_RepositoryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().ListPublicServerSettings(ctx).Return(nil, errors.New("database error"))

	svc := settings.NewService(repo)
	result, err := svc.ListPublicServerSettings(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list public server settings")
}

func TestUnit_SetServerSetting_CheckExistingError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().GetServerSetting(ctx, "error.key").Return(nil, errors.New("connection refused"))

	svc := settings.NewService(repo)
	result, err := svc.SetServerSetting(ctx, "error.key", "value", adminID)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check existing setting")
}

func TestUnit_SetServerSetting_UpsertError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().GetServerSetting(ctx, "new.key").Return(nil, pgx.ErrNoRows)
	repo.EXPECT().UpsertServerSetting(ctx, mock.AnythingOfType("db.UpsertServerSettingParams")).
		Return(nil, errors.New("insert failed"))

	svc := settings.NewService(repo)
	result, err := svc.SetServerSetting(ctx, "new.key", "value", adminID)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to set server setting")
}

func TestUnit_SetServerSetting_UpdateError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockSettingsRepository(t)
	existingSetting := makeServerSetting("existing.key", "string", []byte(`"old_value"`))
	repo.EXPECT().GetServerSetting(ctx, "existing.key").Return(existingSetting, nil)
	repo.EXPECT().UpdateServerSetting(ctx, mock.AnythingOfType("db.UpdateServerSettingParams")).
		Return(nil, errors.New("update failed"))

	svc := settings.NewService(repo)
	result, err := svc.SetServerSetting(ctx, "existing.key", "new_value", adminID)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to set server setting")
}

func TestUnit_SetUserSetting_RepositoryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	repo.EXPECT().UpsertUserSetting(ctx, mock.AnythingOfType("db.UpsertUserSettingParams")).
		Return(nil, errors.New("insert failed"))

	svc := settings.NewService(repo)
	result, err := svc.SetUserSetting(ctx, userID, "user.key", "value")

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to set user setting")
}

func TestUnit_ListServerSettings_InvalidJSON(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockSettingsRepository(t)
	// Return a setting with invalid JSON value
	dbSettings := []db.SharedServerSetting{
		{
			Key:      "invalid.key",
			Value:    []byte(`{invalid json`), // Invalid JSON
			DataType: "string",
		},
	}

	repo.EXPECT().ListServerSettings(ctx).Return(dbSettings, nil)

	svc := settings.NewService(repo)
	result, err := svc.ListServerSettings(ctx)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

func TestUnit_ListUserSettings_InvalidJSON(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockSettingsRepository(t)
	// Return a setting with invalid JSON value
	dbSettings := []db.SharedUserSetting{
		{
			UserID:   userID,
			Key:      "invalid.key",
			Value:    []byte(`{invalid json`), // Invalid JSON
			DataType: "string",
		},
	}

	repo.EXPECT().ListUserSettings(ctx, userID).Return(dbSettings, nil)

	svc := settings.NewService(repo)
	result, err := svc.ListUserSettings(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}
