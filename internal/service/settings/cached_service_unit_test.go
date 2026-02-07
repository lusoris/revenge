package settings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/service/settings"
)

// mockService implements settings.Service for unit tests without a database.
type mockService struct {
	getServerSettingFn            func(ctx context.Context, key string) (*settings.ServerSetting, error)
	listServerSettingsFn          func(ctx context.Context) ([]settings.ServerSetting, error)
	listServerSettingsByCategoryFn func(ctx context.Context, category string) ([]settings.ServerSetting, error)
	listPublicServerSettingsFn    func(ctx context.Context) ([]settings.ServerSetting, error)
	setServerSettingFn            func(ctx context.Context, key string, value interface{}, updatedBy uuid.UUID) (*settings.ServerSetting, error)
	deleteServerSettingFn         func(ctx context.Context, key string) error
	getUserSettingFn              func(ctx context.Context, userID uuid.UUID, key string) (*settings.UserSetting, error)
	listUserSettingsFn            func(ctx context.Context, userID uuid.UUID) ([]settings.UserSetting, error)
	listUserSettingsByCategoryFn  func(ctx context.Context, userID uuid.UUID, category string) ([]settings.UserSetting, error)
	setUserSettingFn              func(ctx context.Context, userID uuid.UUID, key string, value interface{}) (*settings.UserSetting, error)
	setUserSettingsBulkFn         func(ctx context.Context, userID uuid.UUID, s map[string]interface{}) error
	deleteUserSettingFn           func(ctx context.Context, userID uuid.UUID, key string) error
}

func (m *mockService) GetServerSetting(ctx context.Context, key string) (*settings.ServerSetting, error) {
	if m.getServerSettingFn != nil {
		return m.getServerSettingFn(ctx, key)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) ListServerSettings(ctx context.Context) ([]settings.ServerSetting, error) {
	if m.listServerSettingsFn != nil {
		return m.listServerSettingsFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) ListServerSettingsByCategory(ctx context.Context, category string) ([]settings.ServerSetting, error) {
	if m.listServerSettingsByCategoryFn != nil {
		return m.listServerSettingsByCategoryFn(ctx, category)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) ListPublicServerSettings(ctx context.Context) ([]settings.ServerSetting, error) {
	if m.listPublicServerSettingsFn != nil {
		return m.listPublicServerSettingsFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) SetServerSetting(ctx context.Context, key string, value interface{}, updatedBy uuid.UUID) (*settings.ServerSetting, error) {
	if m.setServerSettingFn != nil {
		return m.setServerSettingFn(ctx, key, value, updatedBy)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) DeleteServerSetting(ctx context.Context, key string) error {
	if m.deleteServerSettingFn != nil {
		return m.deleteServerSettingFn(ctx, key)
	}
	return errors.New("not implemented")
}

func (m *mockService) GetUserSetting(ctx context.Context, userID uuid.UUID, key string) (*settings.UserSetting, error) {
	if m.getUserSettingFn != nil {
		return m.getUserSettingFn(ctx, userID, key)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) ListUserSettings(ctx context.Context, userID uuid.UUID) ([]settings.UserSetting, error) {
	if m.listUserSettingsFn != nil {
		return m.listUserSettingsFn(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) ListUserSettingsByCategory(ctx context.Context, userID uuid.UUID, category string) ([]settings.UserSetting, error) {
	if m.listUserSettingsByCategoryFn != nil {
		return m.listUserSettingsByCategoryFn(ctx, userID, category)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) SetUserSetting(ctx context.Context, userID uuid.UUID, key string, value interface{}) (*settings.UserSetting, error) {
	if m.setUserSettingFn != nil {
		return m.setUserSettingFn(ctx, userID, key, value)
	}
	return nil, errors.New("not implemented")
}

func (m *mockService) SetUserSettingsBulk(ctx context.Context, userID uuid.UUID, s map[string]interface{}) error {
	if m.setUserSettingsBulkFn != nil {
		return m.setUserSettingsBulkFn(ctx, userID, s)
	}
	return errors.New("not implemented")
}

func (m *mockService) DeleteUserSetting(ctx context.Context, userID uuid.UUID, key string) error {
	if m.deleteUserSettingFn != nil {
		return m.deleteUserSettingFn(ctx, userID, key)
	}
	return errors.New("not implemented")
}

func newTestCache(t *testing.T) *cache.Cache {
	t.Helper()
	// L1 TTL must be <= the shortest setting TTL used (UserSettingsTTL = 2min)
	// so that Set() actually populates L1 (it skips L1 when ttl < l1TTL).
	c, err := cache.NewCache(nil, 1000, 1*time.Minute)
	require.NoError(t, err)
	return c
}

// ============================================================================
// CachedService Constructor Tests
// ============================================================================

func TestCachedService_NewCachedService_Unit(t *testing.T) {
	t.Parallel()

	t.Run("with cache", func(t *testing.T) {
		t.Parallel()
		svc := &mockService{}
		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())
		require.NotNil(t, cached)
	})

	t.Run("with nil cache", func(t *testing.T) {
		t.Parallel()
		svc := &mockService{}
		cached := settings.NewCachedService(svc, nil, zap.NewNop())
		require.NotNil(t, cached)
	})
}

// ============================================================================
// GetServerSetting Tests
// ============================================================================

func TestCachedService_GetServerSetting_Unit(t *testing.T) {
	t.Parallel()

	t.Run("no cache falls back to underlying service", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			getServerSettingFn: func(_ context.Context, key string) (*settings.ServerSetting, error) {
				return &settings.ServerSetting{Key: key, Value: "value", DataType: "string"}, nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		result, err := cached.GetServerSetting(context.Background(), "test.key")
		require.NoError(t, err)
		assert.Equal(t, "test.key", result.Key)
		assert.Equal(t, "value", result.Value)
	})

	t.Run("cache miss loads from service", func(t *testing.T) {
		t.Parallel()

		calls := 0
		svc := &mockService{
			getServerSettingFn: func(_ context.Context, key string) (*settings.ServerSetting, error) {
				calls++
				return &settings.ServerSetting{Key: key, Value: "from_db", DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		result, err := cached.GetServerSetting(context.Background(), "cache.miss.key")
		require.NoError(t, err)
		assert.Equal(t, "from_db", result.Value)
		assert.Equal(t, 1, calls)
	})

	t.Run("cache hit returns cached value", func(t *testing.T) {
		t.Parallel()

		calls := 0
		svc := &mockService{
			getServerSettingFn: func(_ context.Context, key string) (*settings.ServerSetting, error) {
				calls++
				return &settings.ServerSetting{Key: key, Value: "from_db", DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		// First call - cache miss
		_, err := cached.GetServerSetting(context.Background(), "cache.hit.key")
		require.NoError(t, err)

		// Wait for async cache write
		time.Sleep(50 * time.Millisecond)

		// Second call - should hit cache
		result, err := cached.GetServerSetting(context.Background(), "cache.hit.key")
		require.NoError(t, err)
		assert.Equal(t, "cache.hit.key", result.Key)
		assert.Equal(t, "from_db", result.Value)
		assert.Equal(t, 1, calls) // Only called once
	})

	t.Run("underlying service error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			getServerSettingFn: func(_ context.Context, _ string) (*settings.ServerSetting, error) {
				return nil, errors.New("database error")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		_, err := cached.GetServerSetting(context.Background(), "error.key")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

// ============================================================================
// ListServerSettings Tests
// ============================================================================

func TestCachedService_ListServerSettings_Unit(t *testing.T) {
	t.Parallel()

	t.Run("no cache falls back to underlying service", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			listServerSettingsFn: func(_ context.Context) ([]settings.ServerSetting, error) {
				return []settings.ServerSetting{
					{Key: "key1", Value: "val1", DataType: "string"},
					{Key: "key2", Value: "val2", DataType: "string"},
				}, nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		result, err := cached.ListServerSettings(context.Background())
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("cache miss then cache hit", func(t *testing.T) {
		t.Parallel()

		calls := 0
		svc := &mockService{
			listServerSettingsFn: func(_ context.Context) ([]settings.ServerSetting, error) {
				calls++
				return []settings.ServerSetting{
					{Key: "key1", Value: "val1", DataType: "string"},
				}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		// First call
		result, err := cached.ListServerSettings(context.Background())
		require.NoError(t, err)
		assert.Len(t, result, 1)

		// Wait for async cache write
		time.Sleep(50 * time.Millisecond)

		// Second call - cache hit
		result, err = cached.ListServerSettings(context.Background())
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, calls)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			listServerSettingsFn: func(_ context.Context) ([]settings.ServerSetting, error) {
				return nil, errors.New("list failed")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		_, err := cached.ListServerSettings(context.Background())
		require.Error(t, err)
	})
}

// ============================================================================
// ListPublicServerSettings Tests
// ============================================================================

func TestCachedService_ListPublicServerSettings_Unit(t *testing.T) {
	t.Parallel()

	t.Run("no cache falls back to underlying service", func(t *testing.T) {
		t.Parallel()

		isPublic := true
		svc := &mockService{
			listPublicServerSettingsFn: func(_ context.Context) ([]settings.ServerSetting, error) {
				return []settings.ServerSetting{
					{Key: "public.key", Value: "val", DataType: "string", IsPublic: &isPublic},
				}, nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		result, err := cached.ListPublicServerSettings(context.Background())
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("cache miss then cache hit", func(t *testing.T) {
		t.Parallel()

		calls := 0
		isPublic := true
		svc := &mockService{
			listPublicServerSettingsFn: func(_ context.Context) ([]settings.ServerSetting, error) {
				calls++
				return []settings.ServerSetting{
					{Key: "pub.key", Value: "pub_val", DataType: "string", IsPublic: &isPublic},
				}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		// First call
		_, err := cached.ListPublicServerSettings(context.Background())
		require.NoError(t, err)

		// Wait for async cache write
		time.Sleep(50 * time.Millisecond)

		// Second call
		result, err := cached.ListPublicServerSettings(context.Background())
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, calls)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			listPublicServerSettingsFn: func(_ context.Context) ([]settings.ServerSetting, error) {
				return nil, errors.New("db error")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		_, err := cached.ListPublicServerSettings(context.Background())
		require.Error(t, err)
	})
}

// ============================================================================
// SetServerSetting Tests
// ============================================================================

func TestCachedService_SetServerSetting_Unit(t *testing.T) {
	t.Parallel()

	adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("success with cache invalidation", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setServerSettingFn: func(_ context.Context, key string, value interface{}, _ uuid.UUID) (*settings.ServerSetting, error) {
				return &settings.ServerSetting{Key: key, Value: value, DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		result, err := cached.SetServerSetting(context.Background(), "new.key", "value", adminID)
		require.NoError(t, err)
		assert.Equal(t, "new.key", result.Key)
	})

	t.Run("success with nil cache", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setServerSettingFn: func(_ context.Context, key string, value interface{}, _ uuid.UUID) (*settings.ServerSetting, error) {
				return &settings.ServerSetting{Key: key, Value: value, DataType: "string"}, nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		result, err := cached.SetServerSetting(context.Background(), "new.key", "value", adminID)
		require.NoError(t, err)
		assert.Equal(t, "new.key", result.Key)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setServerSettingFn: func(_ context.Context, _ string, _ interface{}, _ uuid.UUID) (*settings.ServerSetting, error) {
				return nil, errors.New("write failed")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		_, err := cached.SetServerSetting(context.Background(), "fail.key", "value", adminID)
		require.Error(t, err)
	})

	t.Run("cache invalidation after set", func(t *testing.T) {
		t.Parallel()

		calls := 0
		svc := &mockService{
			getServerSettingFn: func(_ context.Context, key string) (*settings.ServerSetting, error) {
				calls++
				return &settings.ServerSetting{Key: key, Value: "value_" + string(rune('0'+calls)), DataType: "string"}, nil
			},
			setServerSettingFn: func(_ context.Context, key string, value interface{}, _ uuid.UUID) (*settings.ServerSetting, error) {
				return &settings.ServerSetting{Key: key, Value: value, DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		// Load into cache
		_, err := cached.GetServerSetting(context.Background(), "inv.key")
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond)
		assert.Equal(t, 1, calls)

		// Update setting - should invalidate cache
		_, err = cached.SetServerSetting(context.Background(), "inv.key", "new_value", adminID)
		require.NoError(t, err)

		// Next get should miss cache and call service again
		_, err = cached.GetServerSetting(context.Background(), "inv.key")
		require.NoError(t, err)
		assert.Equal(t, 2, calls) // Called again after invalidation
	})
}

// ============================================================================
// DeleteServerSetting Tests
// ============================================================================

func TestCachedService_DeleteServerSetting_Unit(t *testing.T) {
	t.Parallel()

	t.Run("success with cache invalidation", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			deleteServerSettingFn: func(_ context.Context, _ string) error {
				return nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		err := cached.DeleteServerSetting(context.Background(), "del.key")
		require.NoError(t, err)
	})

	t.Run("success with nil cache", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			deleteServerSettingFn: func(_ context.Context, _ string) error {
				return nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		err := cached.DeleteServerSetting(context.Background(), "del.key")
		require.NoError(t, err)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			deleteServerSettingFn: func(_ context.Context, _ string) error {
				return errors.New("delete failed")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		err := cached.DeleteServerSetting(context.Background(), "fail.key")
		require.Error(t, err)
	})
}

// ============================================================================
// GetUserSetting Tests
// ============================================================================

func TestCachedService_GetUserSetting_Unit(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("no cache falls back to underlying service", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			getUserSettingFn: func(_ context.Context, uid uuid.UUID, key string) (*settings.UserSetting, error) {
				return &settings.UserSetting{UserID: uid, Key: key, Value: "user_val", DataType: "string"}, nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		result, err := cached.GetUserSetting(context.Background(), userID, "user.key")
		require.NoError(t, err)
		assert.Equal(t, "user.key", result.Key)
		assert.Equal(t, userID, result.UserID)
	})

	t.Run("cache miss then hit", func(t *testing.T) {
		t.Parallel()

		calls := 0
		svc := &mockService{
			getUserSettingFn: func(_ context.Context, uid uuid.UUID, key string) (*settings.UserSetting, error) {
				calls++
				return &settings.UserSetting{UserID: uid, Key: key, Value: "cached_val", DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		// First call - miss
		_, err := cached.GetUserSetting(context.Background(), userID, "cache.user.key")
		require.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		// Second call - hit
		result, err := cached.GetUserSetting(context.Background(), userID, "cache.user.key")
		require.NoError(t, err)
		assert.Equal(t, "cached_val", result.Value)
		assert.Equal(t, 1, calls)
	})

	t.Run("error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			getUserSettingFn: func(_ context.Context, _ uuid.UUID, _ string) (*settings.UserSetting, error) {
				return nil, errors.New("not found")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		_, err := cached.GetUserSetting(context.Background(), userID, "missing.key")
		require.Error(t, err)
	})
}

// ============================================================================
// SetUserSetting Tests
// ============================================================================

func TestCachedService_SetUserSetting_Unit(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("success with cache invalidation", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setUserSettingFn: func(_ context.Context, uid uuid.UUID, key string, value interface{}) (*settings.UserSetting, error) {
				return &settings.UserSetting{UserID: uid, Key: key, Value: value, DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		result, err := cached.SetUserSetting(context.Background(), userID, "user.set.key", "new_val")
		require.NoError(t, err)
		assert.Equal(t, "user.set.key", result.Key)
	})

	t.Run("success with nil cache", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setUserSettingFn: func(_ context.Context, uid uuid.UUID, key string, value interface{}) (*settings.UserSetting, error) {
				return &settings.UserSetting{UserID: uid, Key: key, Value: value, DataType: "string"}, nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		result, err := cached.SetUserSetting(context.Background(), userID, "user.set.key", "new_val")
		require.NoError(t, err)
		assert.Equal(t, "user.set.key", result.Key)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setUserSettingFn: func(_ context.Context, _ uuid.UUID, _ string, _ interface{}) (*settings.UserSetting, error) {
				return nil, errors.New("upsert failed")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		_, err := cached.SetUserSetting(context.Background(), userID, "fail.key", "value")
		require.Error(t, err)
	})
}

// ============================================================================
// SetUserSettingsBulk Tests
// ============================================================================

func TestCachedService_SetUserSettingsBulk_Unit(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("success with cache invalidation", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setUserSettingsBulkFn: func(_ context.Context, _ uuid.UUID, _ map[string]interface{}) error {
				return nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		err := cached.SetUserSettingsBulk(context.Background(), userID, map[string]interface{}{
			"bulk.key1": "val1",
			"bulk.key2": "val2",
		})
		require.NoError(t, err)
	})

	t.Run("success with nil cache", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setUserSettingsBulkFn: func(_ context.Context, _ uuid.UUID, _ map[string]interface{}) error {
				return nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		err := cached.SetUserSettingsBulk(context.Background(), userID, map[string]interface{}{
			"bulk.key1": "val1",
		})
		require.NoError(t, err)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			setUserSettingsBulkFn: func(_ context.Context, _ uuid.UUID, _ map[string]interface{}) error {
				return errors.New("bulk failed")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		err := cached.SetUserSettingsBulk(context.Background(), userID, map[string]interface{}{
			"fail.key": "val",
		})
		require.Error(t, err)
	})
}

// ============================================================================
// DeleteUserSetting Tests
// ============================================================================

func TestCachedService_DeleteUserSetting_Unit(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("success with cache invalidation", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			deleteUserSettingFn: func(_ context.Context, _ uuid.UUID, _ string) error {
				return nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		err := cached.DeleteUserSetting(context.Background(), userID, "del.user.key")
		require.NoError(t, err)
	})

	t.Run("success with nil cache", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			deleteUserSettingFn: func(_ context.Context, _ uuid.UUID, _ string) error {
				return nil
			},
		}

		cached := settings.NewCachedService(svc, nil, zap.NewNop())

		err := cached.DeleteUserSetting(context.Background(), userID, "del.user.key")
		require.NoError(t, err)
	})

	t.Run("underlying error propagates", func(t *testing.T) {
		t.Parallel()

		svc := &mockService{
			deleteUserSettingFn: func(_ context.Context, _ uuid.UUID, _ string) error {
				return errors.New("delete failed")
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		err := cached.DeleteUserSetting(context.Background(), userID, "fail.key")
		require.Error(t, err)
	})
}

// ============================================================================
// Cache Invalidation Flow Tests
// ============================================================================

func TestCachedService_CacheInvalidation_Flows_Unit(t *testing.T) {
	t.Parallel()

	t.Run("set server setting invalidates get cache", func(t *testing.T) {
		t.Parallel()

		getCalls := 0
		adminID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

		svc := &mockService{
			getServerSettingFn: func(_ context.Context, key string) (*settings.ServerSetting, error) {
				getCalls++
				return &settings.ServerSetting{Key: key, Value: "val_" + string(rune('0'+getCalls)), DataType: "string"}, nil
			},
			setServerSettingFn: func(_ context.Context, key string, value interface{}, _ uuid.UUID) (*settings.ServerSetting, error) {
				return &settings.ServerSetting{Key: key, Value: value, DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())
		ctx := context.Background()

		// First get - populates cache
		r1, err := cached.GetServerSetting(ctx, "flow.key")
		require.NoError(t, err)
		assert.Equal(t, "val_1", r1.Value)
		time.Sleep(50 * time.Millisecond)

		// Set - invalidates cache
		_, err = cached.SetServerSetting(ctx, "flow.key", "updated", adminID)
		require.NoError(t, err)

		// Get again - should miss cache
		r2, err := cached.GetServerSetting(ctx, "flow.key")
		require.NoError(t, err)
		assert.Equal(t, "val_2", r2.Value)
		assert.Equal(t, 2, getCalls)
	})

	t.Run("delete server setting invalidates cache", func(t *testing.T) {
		t.Parallel()

		getCalls := 0
		svc := &mockService{
			getServerSettingFn: func(_ context.Context, key string) (*settings.ServerSetting, error) {
				getCalls++
				return &settings.ServerSetting{Key: key, Value: "cached", DataType: "string"}, nil
			},
			deleteServerSettingFn: func(_ context.Context, _ string) error {
				return nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())
		ctx := context.Background()

		// Get - populates cache
		_, err := cached.GetServerSetting(ctx, "del.flow.key")
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond)

		// Delete - invalidates cache
		err = cached.DeleteServerSetting(ctx, "del.flow.key")
		require.NoError(t, err)

		// Get again - should miss cache
		_, err = cached.GetServerSetting(ctx, "del.flow.key")
		require.NoError(t, err)
		assert.Equal(t, 2, getCalls)
	})

	t.Run("set user setting invalidates specific key", func(t *testing.T) {
		t.Parallel()

		userID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
		getCalls := 0

		svc := &mockService{
			getUserSettingFn: func(_ context.Context, uid uuid.UUID, key string) (*settings.UserSetting, error) {
				getCalls++
				return &settings.UserSetting{UserID: uid, Key: key, Value: "val_" + string(rune('0'+getCalls)), DataType: "string"}, nil
			},
			setUserSettingFn: func(_ context.Context, uid uuid.UUID, key string, value interface{}) (*settings.UserSetting, error) {
				return &settings.UserSetting{UserID: uid, Key: key, Value: value, DataType: "string"}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())
		ctx := context.Background()

		// Get - populates cache
		_, err := cached.GetUserSetting(ctx, userID, "user.flow.key")
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond)

		// Set - invalidates cache
		_, err = cached.SetUserSetting(ctx, userID, "user.flow.key", "updated")
		require.NoError(t, err)

		// Get again - should miss cache
		r, err := cached.GetUserSetting(ctx, userID, "user.flow.key")
		require.NoError(t, err)
		assert.Equal(t, "val_2", r.Value)
		assert.Equal(t, 2, getCalls)
	})
}

// ============================================================================
// ListServerSettingsByCategory delegation test
// ============================================================================

func TestCachedService_ListServerSettingsByCategory_Unit(t *testing.T) {
	t.Parallel()

	t.Run("delegates to underlying service", func(t *testing.T) {
		t.Parallel()

		category := "general"
		svc := &mockService{
			listServerSettingsByCategoryFn: func(_ context.Context, cat string) ([]settings.ServerSetting, error) {
				assert.Equal(t, category, cat)
				return []settings.ServerSetting{
					{Key: "general.theme", Value: "dark", DataType: "string", Category: &category},
				}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		// CachedService does not override ListServerSettingsByCategory, so it delegates
		result, err := cached.ListServerSettingsByCategory(context.Background(), "general")
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

// ============================================================================
// ListUserSettings delegation test
// ============================================================================

func TestCachedService_ListUserSettings_Unit(t *testing.T) {
	t.Parallel()

	t.Run("delegates to underlying service", func(t *testing.T) {
		t.Parallel()

		userID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
		svc := &mockService{
			listUserSettingsFn: func(_ context.Context, uid uuid.UUID) ([]settings.UserSetting, error) {
				assert.Equal(t, userID, uid)
				return []settings.UserSetting{
					{UserID: uid, Key: "pref1", Value: "v1", DataType: "string"},
				}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		result, err := cached.ListUserSettings(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

// ============================================================================
// ListUserSettingsByCategory delegation test
// ============================================================================

func TestCachedService_ListUserSettingsByCategory_Unit(t *testing.T) {
	t.Parallel()

	t.Run("delegates to underlying service", func(t *testing.T) {
		t.Parallel()

		userID := uuid.MustParse("55555555-5555-5555-5555-555555555555")
		svc := &mockService{
			listUserSettingsByCategoryFn: func(_ context.Context, uid uuid.UUID, cat string) ([]settings.UserSetting, error) {
				assert.Equal(t, userID, uid)
				assert.Equal(t, "display", cat)
				return []settings.UserSetting{
					{UserID: uid, Key: "display.lang", Value: "en", DataType: "string"},
				}, nil
			},
		}

		c := newTestCache(t)
		cached := settings.NewCachedService(svc, c, zap.NewNop())

		result, err := cached.ListUserSettingsByCategory(context.Background(), userID, "display")
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}
