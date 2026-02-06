package user

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupCachedService(t *testing.T) (*CachedService, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewPostgresRepository(queries)
	activityLogger := activity.NewNoopLogger()
	mockStorage := storage.NewMockStorage()
	avatarCfg := config.AvatarConfig{
		StoragePath:  "/tmp/test-avatars",
		MaxSizeBytes: 5 * 1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
	}
	baseSvc := NewService(testDB.Pool(), repo, activityLogger, mockStorage, avatarCfg)

	// Create cache instance (L1-only for tests)
	testCache, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)

	cachedSvc := NewCachedService(baseSvc, testCache, zap.NewNop())
	return cachedSvc, testDB
}

// Test NewCachedService
func TestNewCachedService(t *testing.T) {
	t.Parallel()
	baseSvc, _ := setupTestService(t)

	testCache, err := cache.NewCache(nil, 1000, 15*time.Minute)
	require.NoError(t, err)

	cached := NewCachedService(baseSvc, testCache, zap.NewNop())
	require.NotNil(t, cached)
	assert.NotNil(t, cached.Service)
	assert.NotNil(t, cached.cache)
}

// Test GetUser with caching
func TestCachedService_GetUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedService(t)
	ctx := context.Background()

	// Create a user
	created, err := svc.Service.CreateUser(ctx, CreateUserParams{
		Username:     "cacheduser",
		Email:        "cacheduser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("first call - cache miss", func(t *testing.T) {
		user, err := svc.GetUser(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, "cacheduser", user.Username)
	})

	t.Run("second call - cache hit", func(t *testing.T) {
		user, err := svc.GetUser(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, "cacheduser", user.Username)
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := svc.GetUser(ctx, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
	})
}

// Test GetUserByUsername with caching
func TestCachedService_GetUserByUsername(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedService(t)
	ctx := context.Background()

	// Create a user
	_, err := svc.Service.CreateUser(ctx, CreateUserParams{
		Username:     "cachedbyname",
		Email:        "cachedbyname@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("first call - cache miss", func(t *testing.T) {
		user, err := svc.GetUserByUsername(ctx, "cachedbyname")
		require.NoError(t, err)
		assert.Equal(t, "cachedbyname", user.Username)
	})

	t.Run("second call - cache hit", func(t *testing.T) {
		user, err := svc.GetUserByUsername(ctx, "cachedbyname")
		require.NoError(t, err)
		assert.Equal(t, "cachedbyname", user.Username)
	})

	t.Run("non-existent username", func(t *testing.T) {
		_, err := svc.GetUserByUsername(ctx, "doesnotexist")
		require.Error(t, err)
	})
}

// Test UpdateUser invalidates cache
func TestCachedService_UpdateUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedService(t)
	ctx := context.Background()

	// Create and cache a user
	created, err := svc.Service.CreateUser(ctx, CreateUserParams{
		Username:     "updatecached",
		Email:        "updatecached@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Get user to cache it
	_, err = svc.GetUser(ctx, created.ID)
	require.NoError(t, err)

	// Update user - should invalidate cache
	displayName := "Updated Name"
	updated, err := svc.UpdateUser(ctx, created.ID, UpdateUserParams{
		DisplayName: &displayName,
	})
	require.NoError(t, err)
	require.NotNil(t, updated.DisplayName)
	assert.Equal(t, "Updated Name", *updated.DisplayName)

	// Get user again - should fetch fresh data, not cached
	user, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, user.DisplayName)
	assert.Equal(t, "Updated Name", *user.DisplayName)
}

// Test DeleteUser invalidates cache
func TestCachedService_DeleteUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedService(t)
	ctx := context.Background()

	// Create and cache a user
	created, err := svc.Service.CreateUser(ctx, CreateUserParams{
		Username:     "deletecached",
		Email:        "deletecached@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Get user to cache it
	_, err = svc.GetUser(ctx, created.ID)
	require.NoError(t, err)

	// Delete user - should invalidate cache
	err = svc.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	// Get user again - should get error (user deleted)
	_, err = svc.GetUser(ctx, created.ID)
	require.Error(t, err)
}

// Test InvalidateUserCache
func TestCachedService_InvalidateUserCache(t *testing.T) {
	t.Parallel()
	svc, _ := setupCachedService(t)
	ctx := context.Background()

	// Create and cache a user
	created, err := svc.Service.CreateUser(ctx, CreateUserParams{
		Username:     "invalidatecache",
		Email:        "invalidatecache@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Get user to cache it
	user1, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)

	// Manually invalidate cache
	err = svc.InvalidateUserCache(ctx, created.ID)
	require.NoError(t, err)

	// Update user directly via Service (bypassing cached service)
	displayName := "Direct Update"
	_, err = svc.Service.UpdateUser(ctx, created.ID, UpdateUserParams{
		DisplayName: &displayName,
	})
	require.NoError(t, err)

	// Get user through cached service - should get fresh data after cache was invalidated
	user2, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, user2.DisplayName)
	assert.Equal(t, "Direct Update", *user2.DisplayName)

	// Should be different from originally cached user
	assert.NotEqual(t, user1.DisplayName, user2.DisplayName)
}
