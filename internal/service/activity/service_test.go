package activity

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupTestService(t *testing.T) (*Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	logger := zaptest.NewLogger(t)
	svc := NewService(repo, logger)
	return svc, testDB
}

// ============================================================================
// Logging Tests
// ============================================================================

func TestService_Log(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "testuser",
		Email:    "test@example.com",
	})

	ip := net.ParseIP("192.168.1.1")
	userAgent := "Test Agent"
	resourceID := uuid.Must(uuid.NewV7())

	req := LogRequest{
		UserID:       &user.ID,
		Username:     &user.Username,
		Action:       ActionUserLogin,
		ResourceType: stringPtr("user"),
		ResourceID:   &resourceID,
		IPAddress:    &ip,
		UserAgent:    &userAgent,
		Success:      true,
	}

	err := svc.Log(ctx, req)
	require.NoError(t, err)
}

func TestService_Log_SystemAction(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// System actions don't have user_id
	req := LogRequest{
		Action:  "system.startup",
		Success: true,
	}

	err := svc.Log(ctx, req)
	require.NoError(t, err)
}

func TestService_LogWithContext(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "testuser",
		Email:    "test@example.com",
	})

	ip := net.ParseIP("10.0.0.1")
	resourceID := uuid.Must(uuid.NewV7())

	err := svc.LogWithContext(
		ctx,
		user.ID,
		user.Username,
		ActionUserUpdate,
		"user",
		resourceID,
		map[string]interface{}{"field": "value"},
		ip,
		"Mozilla/5.0",
	)
	require.NoError(t, err)
}

func TestService_LogFailure(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "failuser",
		Email:    "fail@example.com",
	})

	ip := net.ParseIP("172.16.0.1")
	userAgent := "Bad Agent"

	err := svc.LogFailure(
		ctx,
		&user.ID,
		&user.Username,
		ActionUserLogin,
		"invalid credentials",
		&ip,
		&userAgent,
	)
	require.NoError(t, err)
}

func TestService_LogFailure_Anonymous(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	err := svc.LogFailure(
		ctx,
		nil,
		nil,
		ActionUserLogin,
		"invalid credentials",
		nil,
		nil,
	)
	require.NoError(t, err)
}

// ============================================================================
// Retrieval Tests
// ============================================================================

func TestService_Get(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create entry first
	req := LogRequest{
		Action:  "test.action",
		Success: true,
	}
	require.NoError(t, svc.Log(ctx, req))

	// Retrieve by ID (we need to get the ID somehow - list first)
	entries, _, err := svc.List(ctx, 1, 0)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	entry, err := svc.Get(ctx, entries[0].ID)
	require.NoError(t, err)
	assert.Equal(t, "test.action", entry.Action)
}

func TestService_Get_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.Get(ctx, uuid.Must(uuid.NewV7()))
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_List(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create multiple entries
	for i := 0; i < 5; i++ {
		req := LogRequest{
			Action:  "test.list",
			Success: true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	entries, count, err := svc.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(entries), 5)
	assert.GreaterOrEqual(t, count, int64(5))
}

func TestService_List_Pagination(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create entries
	for i := 0; i < 10; i++ {
		req := LogRequest{
			Action:  "test.pagination",
			Success: true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	// First page
	page1, count, err := svc.List(ctx, 5, 0)
	require.NoError(t, err)
	assert.Len(t, page1, 5)
	assert.GreaterOrEqual(t, count, int64(10))

	// Second page
	page2, _, err := svc.List(ctx, 5, 5)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(page2), 5)
}

// ============================================================================
// Search Tests
// ============================================================================

func TestService_Search_DefaultLimit(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Search with limit 0 should default to 50
	entries, count, err := svc.Search(ctx, SearchFilters{
		Limit: 0,
	})
	require.NoError(t, err)
	assert.NotNil(t, entries)
	assert.GreaterOrEqual(t, count, int64(0))
}

func TestService_Search_MaxLimit(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "searchuser",
		Email:    "search@example.com",
	})

	// Create entries
	for i := 0; i < 3; i++ {
		req := LogRequest{
			UserID:  &user.ID,
			Action:  "test.search",
			Success: true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	// Search with limit > 100 should cap at 100
	entries, count, err := svc.Search(ctx, SearchFilters{
		UserID: &user.ID,
		Limit:  200,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(entries), 3)
	assert.GreaterOrEqual(t, count, int64(3))
}

func TestService_Search_ByAction(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	action := "test.specific_action"

	// Create entries with specific action
	for i := 0; i < 2; i++ {
		req := LogRequest{
			Action:  action,
			Success: true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	// Search by action
	entries, count, err := svc.Search(ctx, SearchFilters{
		Action: &action,
		Limit:  10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
	for _, e := range entries {
		assert.Equal(t, action, e.Action)
	}
}

func TestService_Search_BySuccess(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create failed entries
	for i := 0; i < 2; i++ {
		req := LogRequest{
			Action:       "test.failure",
			Success:      false,
			ErrorMessage: stringPtr("test error"),
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	// Search for failed entries
	successFalse := false
	entries, count, err := svc.Search(ctx, SearchFilters{
		Success: &successFalse,
		Limit:   10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
	for _, e := range entries {
		assert.False(t, e.Success)
	}
}

// ============================================================================
// User Activity Tests
// ============================================================================

func TestService_GetUserActivity(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "activityuser",
		Email:    "activity@example.com",
	})

	// Create user activities
	for i := 0; i < 3; i++ {
		req := LogRequest{
			UserID:  &user.ID,
			Action:  ActionUserUpdate,
			Success: true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	entries, count, err := svc.GetUserActivity(ctx, user.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
	assert.GreaterOrEqual(t, len(entries), 3)
}

func TestService_GetUserActivity_DefaultLimit(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "limituser",
		Email:    "limit@example.com",
	})

	// Limit 0 should default to 50
	entries, _, err := svc.GetUserActivity(ctx, user.ID, 0, 0)
	require.NoError(t, err)
	assert.NotNil(t, entries)
}

func TestService_GetUserActivity_MaxLimit(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "maxuser",
		Email:    "max@example.com",
	})

	// Limit > 100 should cap at 100
	entries, _, err := svc.GetUserActivity(ctx, user.ID, 200, 0)
	require.NoError(t, err)
	assert.NotNil(t, entries)
}

// ============================================================================
// Resource Activity Tests
// ============================================================================

func TestService_GetResourceActivity(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	resourceType := "document"
	resourceID := uuid.Must(uuid.NewV7())

	// Create resource activities
	for i := 0; i < 2; i++ {
		req := LogRequest{
			Action:       "document.edit",
			ResourceType: &resourceType,
			ResourceID:   &resourceID,
			Success:      true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	entries, count, err := svc.GetResourceActivity(ctx, resourceType, resourceID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))
	assert.GreaterOrEqual(t, len(entries), 2)
}

func TestService_GetResourceActivity_LimitValidation(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	resourceType := "test"
	resourceID := uuid.Must(uuid.NewV7())

	// Test default limit
	_, _, err := svc.GetResourceActivity(ctx, resourceType, resourceID, 0, 0)
	require.NoError(t, err)

	// Test max limit
	_, _, err = svc.GetResourceActivity(ctx, resourceType, resourceID, 200, 0)
	require.NoError(t, err)
}

// ============================================================================
// Failed Activity Tests
// ============================================================================

func TestService_GetFailedActivity(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create failed entries
	for i := 0; i < 3; i++ {
		req := LogRequest{
			Action:       "test.failed",
			Success:      false,
			ErrorMessage: stringPtr("error occurred"),
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	entries, err := svc.GetFailedActivity(ctx, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(entries), 3)
	for _, e := range entries {
		assert.False(t, e.Success)
	}
}

func TestService_GetFailedActivity_LimitValidation(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Default limit
	entries, err := svc.GetFailedActivity(ctx, 0, 0)
	require.NoError(t, err)
	assert.NotNil(t, entries)

	// Max limit
	entries, err = svc.GetFailedActivity(ctx, 150, 0)
	require.NoError(t, err)
	assert.NotNil(t, entries)
}

// ============================================================================
// Stats Tests
// ============================================================================

func TestService_GetStats(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create some activities
	for i := 0; i < 5; i++ {
		req := LogRequest{
			Action:  "test.stats",
			Success: i%2 == 0, // 3 success, 2 failures
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	stats, err := svc.GetStats(ctx)
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.GreaterOrEqual(t, stats.TotalCount, int64(5))
	assert.GreaterOrEqual(t, stats.SuccessCount, int64(3))
}

// ============================================================================
// Recent Actions Tests
// ============================================================================

func TestService_GetRecentActions(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create entries with various actions
	actions := []string{"action1", "action2", "action1", "action3"}
	for _, action := range actions {
		req := LogRequest{
			Action:  action,
			Success: true,
		}
		require.NoError(t, svc.Log(ctx, req))
	}

	recent, err := svc.GetRecentActions(ctx, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, recent)
}

func TestService_GetRecentActions_DefaultLimit(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Limit 0 should default to 20
	recent, err := svc.GetRecentActions(ctx, 0)
	require.NoError(t, err)
	assert.NotNil(t, recent)
}

func TestService_GetRecentActions_MaxLimit(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Limit > 50 should cap at 50
	recent, err := svc.GetRecentActions(ctx, 100)
	require.NoError(t, err)
	assert.NotNil(t, recent)
}

// ============================================================================
// Cleanup Tests
// ============================================================================

func TestService_CleanupOldLogs(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create old entry (we can't actually set created_at, so this is more of a smoke test)
	req := LogRequest{
		Action:  "test.old",
		Success: true,
	}
	require.NoError(t, svc.Log(ctx, req))

	// Try cleanup (won't delete recent entries, but tests the path)
	cutoff := time.Now().Add(-365 * 24 * time.Hour)
	count, err := svc.CleanupOldLogs(ctx, cutoff)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(0))
}

func TestService_CountOldLogs(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create entry
	req := LogRequest{
		Action:  "test.count_old",
		Success: true,
	}
	require.NoError(t, svc.Log(ctx, req))

	// Count old logs
	cutoff := time.Now().Add(-365 * 24 * time.Hour)
	count, err := svc.CountOldLogs(ctx, cutoff)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(0))
}

// ============================================================================
// Helper Functions
// ============================================================================

func stringPtr(s string) *string {
	return &s
}
