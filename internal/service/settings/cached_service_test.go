package settings

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupCachedTestService(t *testing.T) (*CachedService, testutil.DB) {
	t.Helper()
	svc, testDB := setupTestService(t)

	// Create cache instance (L1-only for tests)
	testCache, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)

	cachedSvc := NewCachedService(svc, testCache, logging.NewTestLogger())
	return cachedSvc, testDB
}

func TestNewCachedService(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

	testCache, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)

	cached := NewCachedService(svc, testCache, logging.NewTestLogger())
	require.NotNil(t, cached)
	assert.NotNil(t, cached.Service)
	assert.NotNil(t, cached.cache)
}

// ============================================================================
// Server Settings Tests
// ============================================================================

func TestCachedService_GetServerSetting(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set a server setting
	created, err := svc.SetServerSetting(ctx, "test.key", "test value", userID)
	require.NoError(t, err)
	require.NotNil(t, created)

	// First call - cache miss
	setting, err := svc.GetServerSetting(ctx, "test.key")
	require.NoError(t, err)
	assert.Equal(t, "test.key", setting.Key)

	// Second call - cache hit
	setting, err = svc.GetServerSetting(ctx, "test.key")
	require.NoError(t, err)
	assert.Equal(t, "test.key", setting.Key)
}

func TestCachedService_GetServerSetting_NoCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	// Create cached service without cache (nil cache)
	cachedSvc := NewCachedService(svc, nil, logging.NewTestLogger())

	userID := createTestUser(t, testDB)

	// Set a server setting
	_, err := svc.SetServerSetting(ctx, "nocache.key", "value", userID)
	require.NoError(t, err)

	// Get should work without cache
	setting, err := cachedSvc.GetServerSetting(ctx, "nocache.key")
	require.NoError(t, err)
	assert.Equal(t, "nocache.key", setting.Key)
}

func TestCachedService_ListServerSettings(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set multiple server settings
	_, err := svc.SetServerSetting(ctx, "list.key1", "value1", userID)
	require.NoError(t, err)
	_, err = svc.SetServerSetting(ctx, "list.key2", "value2", userID)
	require.NoError(t, err)

	// First call - cache miss
	settings, err := svc.ListServerSettings(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(settings), 2)

	// Second call - cache hit
	settings, err = svc.ListServerSettings(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(settings), 2)
}

func TestCachedService_ListServerSettings_NoCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	cachedSvc := NewCachedService(svc, nil, logging.NewTestLogger())

	// Should work without cache
	settings, err := cachedSvc.ListServerSettings(ctx)
	require.NoError(t, err)
	assert.NotNil(t, settings)
}

func TestCachedService_ListPublicServerSettings(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set a public server setting
	_, err := svc.SetServerSetting(ctx, "public.key", "public value", userID)
	require.NoError(t, err)

	// First call - cache miss
	settings, err := svc.ListPublicServerSettings(ctx)
	require.NoError(t, err)
	assert.NotNil(t, settings)

	// Second call - cache hit
	settings, err = svc.ListPublicServerSettings(ctx)
	require.NoError(t, err)
	assert.NotNil(t, settings)
}

func TestCachedService_ListPublicServerSettings_NoCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	cachedSvc := NewCachedService(svc, nil, logging.NewTestLogger())

	// Should work without cache
	settings, err := cachedSvc.ListPublicServerSettings(ctx)
	require.NoError(t, err)
	assert.NotNil(t, settings)
}

func TestCachedService_SetServerSetting_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set and cache a setting
	_, err := svc.SetServerSetting(ctx, "invalidate.key", "old value", userID)
	require.NoError(t, err)

	// Get to cache it
	setting, err := svc.GetServerSetting(ctx, "invalidate.key")
	require.NoError(t, err)
	assert.Equal(t, "old value", setting.Value)

	// Update the setting - should invalidate cache
	_, err = svc.SetServerSetting(ctx, "invalidate.key", "new value", userID)
	require.NoError(t, err)

	// Get again - should get fresh data (not cached old value)
	setting, err = svc.GetServerSetting(ctx, "invalidate.key")
	require.NoError(t, err)
	assert.Equal(t, "new value", setting.Value)
}

func TestCachedService_DeleteServerSetting_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set and cache a setting
	_, err := svc.SetServerSetting(ctx, "delete.key", "value", userID)
	require.NoError(t, err)

	// Get to cache it
	_, err = svc.GetServerSetting(ctx, "delete.key")
	require.NoError(t, err)

	// Delete the setting - should invalidate cache
	err = svc.DeleteServerSetting(ctx, "delete.key")
	require.NoError(t, err)

	// Get again - should get error (not cached value)
	_, err = svc.GetServerSetting(ctx, "delete.key")
	assert.Error(t, err)
}

// ============================================================================
// User Settings Tests
// ============================================================================

func TestCachedService_GetUserSetting(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set a user setting
	created, err := svc.SetUserSetting(ctx, userID, "user.test.key", "user value")
	require.NoError(t, err)
	require.NotNil(t, created)

	// First call - cache miss
	setting, err := svc.GetUserSetting(ctx, userID, "user.test.key")
	require.NoError(t, err)
	assert.Equal(t, "user.test.key", setting.Key)

	// Second call - cache hit
	setting, err = svc.GetUserSetting(ctx, userID, "user.test.key")
	require.NoError(t, err)
	assert.Equal(t, "user.test.key", setting.Key)
}

func TestCachedService_GetUserSetting_NoCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	cachedSvc := NewCachedService(svc, nil, logging.NewTestLogger())

	userID := createTestUser(t, testDB)

	// Set a user setting
	_, err := svc.SetUserSetting(ctx, userID, "user.nocache.key", "value")
	require.NoError(t, err)

	// Get should work without cache
	setting, err := cachedSvc.GetUserSetting(ctx, userID, "user.nocache.key")
	require.NoError(t, err)
	assert.Equal(t, "user.nocache.key", setting.Key)
}

func TestCachedService_SetUserSetting_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set and cache a user setting
	_, err := svc.SetUserSetting(ctx, userID, "user.invalidate.key", "old value")
	require.NoError(t, err)

	// Get to cache it
	setting, err := svc.GetUserSetting(ctx, userID, "user.invalidate.key")
	require.NoError(t, err)
	assert.Equal(t, "old value", setting.Value)

	// Update the setting - should invalidate cache
	_, err = svc.SetUserSetting(ctx, userID, "user.invalidate.key", "new value")
	require.NoError(t, err)

	// Get again - should get fresh data
	setting, err = svc.GetUserSetting(ctx, userID, "user.invalidate.key")
	require.NoError(t, err)
	assert.Equal(t, "new value", setting.Value)
}

func TestCachedService_SetUserSettingsBulk_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set initial settings
	_, err := svc.SetUserSetting(ctx, userID, "bulk.key1", "value1")
	require.NoError(t, err)
	_, err = svc.SetUserSetting(ctx, userID, "bulk.key2", "value2")
	require.NoError(t, err)

	// Cache them
	_, err = svc.GetUserSetting(ctx, userID, "bulk.key1")
	require.NoError(t, err)

	// Bulk update - should invalidate cache
	err = svc.SetUserSettingsBulk(ctx, userID, map[string]interface{}{
		"bulk.key1": "new1",
		"bulk.key2": "new2",
		"bulk.key3": "new3",
	})
	require.NoError(t, err)

	// Verify fresh data is returned
	setting, err := svc.GetUserSetting(ctx, userID, "bulk.key1")
	require.NoError(t, err)
	assert.Equal(t, "new1", setting.Value)
}

func TestCachedService_DeleteUserSetting_InvalidatesCache(t *testing.T) {
	t.Parallel()
	svc, testDB := setupCachedTestService(t)
	ctx := context.Background()

	userID := createTestUser(t, testDB)

	// Set and cache a user setting
	_, err := svc.SetUserSetting(ctx, userID, "user.delete.key", "value")
	require.NoError(t, err)

	// Get to cache it
	_, err = svc.GetUserSetting(ctx, userID, "user.delete.key")
	require.NoError(t, err)

	// Delete the setting - should invalidate cache
	err = svc.DeleteUserSetting(ctx, userID, "user.delete.key")
	require.NoError(t, err)

	// Get again - should get error (not cached value)
	_, err = svc.GetUserSetting(ctx, userID, "user.delete.key")
	assert.Error(t, err)
}
