package activity

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.StopSharedPostgres()
	os.Exit(code)
}

func setupTestRepository(t *testing.T) (*RepositoryPg, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	return repo, testDB
}

func createTestEntry(userID *uuid.UUID, action string, success bool) *Entry {
	username := "testuser"
	resourceType := "test_resource"
	resourceID := uuid.Must(uuid.NewV7())
	ip := net.ParseIP("192.168.1.1")
	userAgent := "test-agent"

	return &Entry{
		UserID:       userID,
		Username:     &username,
		Action:       action,
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
		Changes:      map[string]any{"field": "value"},
		Metadata:     map[string]any{"key": "value"},
		IPAddress:    &ip,
		UserAgent:    &userAgent,
		Success:      success,
	}
}

func TestRepositoryPg_Create(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Use nil userID since we don't have actual users in the database
	entry := createTestEntry(nil, "system.event", true)

	err := repo.Create(ctx, entry)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, entry.ID)
	assert.False(t, entry.CreatedAt.IsZero())
}

func TestRepositoryPg_Create_Minimal(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	entry := &Entry{
		Action:  "system.startup",
		Success: true,
	}

	err := repo.Create(ctx, entry)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, entry.ID)
}

func TestRepositoryPg_Get(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	entry := createTestEntry(nil, "system.login", true)
	require.NoError(t, repo.Create(ctx, entry))

	retrieved, err := repo.Get(ctx, entry.ID)
	require.NoError(t, err)
	assert.Equal(t, entry.ID, retrieved.ID)
	assert.Equal(t, entry.Action, retrieved.Action)
}

func TestRepositoryPg_Get_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, uuid.Must(uuid.NewV7()))
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRepositoryPg_List(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	for range 5 {
		entry := createTestEntry(nil, "test.action", true)
		require.NoError(t, repo.Create(ctx, entry))
	}

	entries, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(entries), 5)
}

func TestRepositoryPg_Count(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	countBefore, err := repo.Count(ctx)
	require.NoError(t, err)

	for range 3 {
		entry := createTestEntry(nil, "test.action", true)
		require.NoError(t, repo.Create(ctx, entry))
	}

	countAfter, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, countBefore+3, countAfter)
}

func TestRepositoryPg_Search_ByUserID(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Create actual users to test user_id filtering
	user1 := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "user1_" + uuid.Must(uuid.NewV7()).String()[:8],
		Email:    "user1_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
	})
	user2 := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "user2_" + uuid.Must(uuid.NewV7()).String()[:8],
		Email:    "user2_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
	})

	require.NoError(t, repo.Create(ctx, createTestEntry(&user1.ID, "user.login", true)))
	require.NoError(t, repo.Create(ctx, createTestEntry(&user1.ID, "user.logout", true)))
	require.NoError(t, repo.Create(ctx, createTestEntry(&user2.ID, "user.login", true)))

	filters := SearchFilters{
		UserID: &user1.ID,
		Limit:  10,
	}

	entries, count, err := repo.Search(ctx, filters)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
	assert.GreaterOrEqual(t, len(entries), 2)
	for _, entry := range entries {
		assert.Equal(t, user1.ID, *entry.UserID)
	}
}

func TestRepositoryPg_GetByUser(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Create actual users
	user1 := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "user1_" + uuid.Must(uuid.NewV7()).String()[:8],
		Email:    "user1_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
	})
	user2 := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "user2_" + uuid.Must(uuid.NewV7()).String()[:8],
		Email:    "user2_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
	})

	require.NoError(t, repo.Create(ctx, createTestEntry(&user1.ID, "action1", true)))
	require.NoError(t, repo.Create(ctx, createTestEntry(&user1.ID, "action2", true)))
	require.NoError(t, repo.Create(ctx, createTestEntry(&user2.ID, "action3", true)))

	entries, count, err := repo.GetByUser(ctx, user1.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
	assert.GreaterOrEqual(t, len(entries), 2)
}

func TestRepositoryPg_GetStats(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, createTestEntry(nil, "action1", true)))
	require.NoError(t, repo.Create(ctx, createTestEntry(nil, "action2", false)))

	stats, err := repo.GetStats(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.TotalCount, int64(2))
	assert.GreaterOrEqual(t, stats.SuccessCount, int64(1))
	assert.GreaterOrEqual(t, stats.FailedCount, int64(1))
	assert.NotNil(t, stats.OldestEntry)
	assert.NotNil(t, stats.NewestEntry)
}
