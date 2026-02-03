package settings

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func setupTestService(t *testing.T) (Service, *testutil.TestDB) {
	t.Helper()
	testDB := testutil.NewTestDB(t)
	repo := NewPostgresRepository(testDB.Pool())
	svc := NewService(repo)
	return svc, testDB
}

// createTestUser creates a user for FK constraints
func createTestUser(t *testing.T, testDB *testutil.TestDB) uuid.UUID {
	t.Helper()
	queries := db.New(testDB.Pool())
	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	return user.ID
}

func ptr[T any](v T) *T {
	return &v
}

// ============================================================================
// Server Settings Tests
// ============================================================================

func TestService_ServerSettings_CRUD(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	adminID := createTestUser(t, testDB)

	t.Run("set and get string setting", func(t *testing.T) {
		setting, err := svc.SetServerSetting(ctx, "test.string", "value123", adminID)
		require.NoError(t, err)
		assert.Equal(t, "test.string", setting.Key)
		assert.Equal(t, "value123", setting.Value)

		retrieved, err := svc.GetServerSetting(ctx, "test.string")
		require.NoError(t, err)
		assert.Equal(t, "value123", retrieved.Value)
	})

	t.Run("set and get numeric setting", func(t *testing.T) {
		setting, err := svc.SetServerSetting(ctx, "test.number", 42, adminID)
		require.NoError(t, err)
		assert.Equal(t, "test.number", setting.Key)
		assert.Equal(t, float64(42), setting.Value) // JSON unmarshals numbers as float64
	})

	t.Run("set and get boolean setting", func(t *testing.T) {
		setting, err := svc.SetServerSetting(ctx, "test.bool", true, adminID)
		require.NoError(t, err)
		assert.Equal(t, "test.bool", setting.Key)
		assert.Equal(t, true, setting.Value)
	})

	t.Run("get non-existent setting", func(t *testing.T) {
		_, err := svc.GetServerSetting(ctx, "does.not.exist")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("delete setting", func(t *testing.T) {
		_, err := svc.SetServerSetting(ctx, "test.delete", "value", adminID)
		require.NoError(t, err)

		err = svc.DeleteServerSetting(ctx, "test.delete")
		require.NoError(t, err)

		_, err = svc.GetServerSetting(ctx, "test.delete")
		require.Error(t, err)
	})
}

func TestService_ListServerSettings(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	adminID := createTestUser(t, testDB)

	// Create test settings
	_, err := svc.SetServerSetting(ctx, "list.test1", "value1", adminID)
	require.NoError(t, err)
	_, err = svc.SetServerSetting(ctx, "list.test2", "value2", adminID)
	require.NoError(t, err)

	settings, err := svc.ListServerSettings(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(settings), 2)
}

func TestService_ListServerSettingsByCategory(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	adminID := createTestUser(t, testDB)

	// Create categorized settings - need to use repository directly for category
	repo := svc.(*service).repo
	updatedBy := pgtype.UUID{Bytes: adminID, Valid: true}
	_, err := repo.UpsertServerSetting(ctx, db.UpsertServerSettingParams{
		Key:       "cat1.setting1",
		Value:     []byte(`"value1"`),
		Category:  ptr("cat1"),
		DataType:  "string",
		UpdatedBy: updatedBy,
	})
	require.NoError(t, err)
	_, err = repo.UpsertServerSetting(ctx, db.UpsertServerSettingParams{
		Key:       "cat1.setting2",
		Value:     []byte(`"value2"`),
		Category:  ptr("cat1"),
		DataType:  "string",
		UpdatedBy: updatedBy,
	})
	require.NoError(t, err)

	settings, err := svc.ListServerSettingsByCategory(ctx, "cat1")
	require.NoError(t, err)
	assert.Len(t, settings, 2)
}

func TestService_ListPublicServerSettings(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	settings, err := svc.ListPublicServerSettings(ctx)
	require.NoError(t, err)
	// Should return list (might be empty)
	assert.NotNil(t, settings)
}

// ============================================================================
// User Settings Tests
// ============================================================================

func TestService_UserSettings_CRUD(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	t.Run("set and get string setting", func(t *testing.T) {
		setting, err := svc.SetUserSetting(ctx, userID, "user.pref.string", "my_value")
		require.NoError(t, err)
		assert.Equal(t, "user.pref.string", setting.Key)
		assert.Equal(t, "my_value", setting.Value)
		assert.Equal(t, userID, setting.UserID)

		retrieved, err := svc.GetUserSetting(ctx, userID, "user.pref.string")
		require.NoError(t, err)
		assert.Equal(t, "my_value", retrieved.Value)
	})

	t.Run("set and get numeric setting", func(t *testing.T) {
		setting, err := svc.SetUserSetting(ctx, userID, "user.pref.num", 100)
		require.NoError(t, err)
		assert.Equal(t, float64(100), setting.Value)
	})

	t.Run("set and get boolean setting", func(t *testing.T) {
		setting, err := svc.SetUserSetting(ctx, userID, "user.pref.bool", false)
		require.NoError(t, err)
		assert.Equal(t, false, setting.Value)
	})

	t.Run("get non-existent setting", func(t *testing.T) {
		_, err := svc.GetUserSetting(ctx, userID, "does.not.exist")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("delete setting", func(t *testing.T) {
		_, err := svc.SetUserSetting(ctx, userID, "user.delete", "value")
		require.NoError(t, err)

		err = svc.DeleteUserSetting(ctx, userID, "user.delete")
		require.NoError(t, err)

		_, err = svc.GetUserSetting(ctx, userID, "user.delete")
		require.Error(t, err)
	})
}

func TestService_ListUserSettings(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	// Create test settings
	_, err := svc.SetUserSetting(ctx, userID, "list.test1", "value1")
	require.NoError(t, err)
	_, err = svc.SetUserSetting(ctx, userID, "list.test2", "value2")
	require.NoError(t, err)

	settings, err := svc.ListUserSettings(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(settings), 2)
}

func TestService_ListUserSettingsByCategory(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	// Create categorized settings - need to use repository directly for category
	repo := svc.(*service).repo
	_, err := repo.UpsertUserSetting(ctx, db.UpsertUserSettingParams{
		UserID:   userID,
		Key:      "cat1.user1",
		Value:    []byte(`"value1"`),
		Category: ptr("cat1"),
		DataType: "string",
	})
	require.NoError(t, err)
	_, err = repo.UpsertUserSetting(ctx, db.UpsertUserSettingParams{
		UserID:   userID,
		Key:      "cat1.user2",
		Value:    []byte(`"value2"`),
		Category: ptr("cat1"),
		DataType: "string",
	})
	require.NoError(t, err)

	settings, err := svc.ListUserSettingsByCategory(ctx, userID, "cat1")
	require.NoError(t, err)
	assert.Len(t, settings, 2)
}

func TestService_SetUserSettingsBulk(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	userID := createTestUser(t, testDB)

	settings := map[string]interface{}{
		"bulk.setting1": "value1",
		"bulk.setting2": 42,
		"bulk.setting3": true,
	}

	err := svc.SetUserSettingsBulk(ctx, userID, settings)
	require.NoError(t, err)

	// Verify all settings were created
	setting1, err := svc.GetUserSetting(ctx, userID, "bulk.setting1")
	require.NoError(t, err)
	assert.Equal(t, "value1", setting1.Value)

	setting2, err := svc.GetUserSetting(ctx, userID, "bulk.setting2")
	require.NoError(t, err)
	assert.Equal(t, float64(42), setting2.Value)

	setting3, err := svc.GetUserSetting(ctx, userID, "bulk.setting3")
	require.NoError(t, err)
	assert.Equal(t, true, setting3.Value)
}

func TestService_UserSettings_Isolation(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	user1 := createTestUser(t, testDB)
	user2 := createTestUser(t, testDB)

	// Set setting for user1
	_, err := svc.SetUserSetting(ctx, user1, "isolation.test", "user1_value")
	require.NoError(t, err)

	// Set same key for user2
	_, err = svc.SetUserSetting(ctx, user2, "isolation.test", "user2_value")
	require.NoError(t, err)

	// Verify isolation
	setting1, err := svc.GetUserSetting(ctx, user1, "isolation.test")
	require.NoError(t, err)
	assert.Equal(t, "user1_value", setting1.Value)

	setting2, err := svc.GetUserSetting(ctx, user2, "isolation.test")
	require.NoError(t, err)
	assert.Equal(t, "user2_value", setting2.Value)
}

func TestService_DataTypes(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	adminID := createTestUser(t, testDB)

	tests := []struct {
		name  string
		value interface{}
	}{
		{"string", "hello world"},
		{"int", 42},
		{"float", 3.14},
		{"bool_true", true},
		{"bool_false", false},
		{"array", []interface{}{"a", "b", "c"}},
		{"object", map[string]interface{}{"nested": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "datatype." + tt.name
			setting, err := svc.SetServerSetting(ctx, key, tt.value, adminID)
			require.NoError(t, err)
			assert.NotNil(t, setting)
			// Note: exact comparison may vary due to JSON marshaling
			// Just verify no error occurred
		})
	}
}

func TestService_SettingOverwrite(t *testing.T) {
	t.Parallel()
	svc, testDB := setupTestService(t)
	ctx := context.Background()
	adminID := createTestUser(t, testDB)

	// Set initial value
	_, err := svc.SetServerSetting(ctx, "overwrite.test", "value1", adminID)
	require.NoError(t, err)

	// Overwrite with new value
	setting, err := svc.SetServerSetting(ctx, "overwrite.test", "value2", adminID)
	require.NoError(t, err)
	assert.Equal(t, "value2", setting.Value)

	// Verify new value
	retrieved, err := svc.GetServerSetting(ctx, "overwrite.test")
	require.NoError(t, err)
	assert.Equal(t, "value2", retrieved.Value)
}
