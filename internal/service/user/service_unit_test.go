package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/internal/testutil"
	"github.com/lusoris/revenge/internal/util/ptr"
)

// ============================================================================
// Test Helpers
// ============================================================================

func setupUnitTestService(t *testing.T, repo user.Repository) *user.Service {
	t.Helper()
	svc, _ := setupUnitTestServiceWithDB(t, repo)
	return svc
}

func setupUnitTestServiceWithDB(t *testing.T, repo user.Repository) (*user.Service, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	avatarCfg := config.AvatarConfig{
		StoragePath:  "/tmp/test-avatars",
		MaxSizeBytes: 5 * 1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
	}
	return user.NewService(testDB.Pool(), repo, activity.NewNoopLogger(), storage.NewMockStorage(), avatarCfg), testDB
}

// createDBUser creates a user directly in the test database for tests that need FK references.
// Returns the generated user ID.
func createDBUser(t *testing.T, testDB testutil.DB) uuid.UUID {
	t.Helper()
	queries := db.New(testDB.Pool())
	unique := uuid.Must(uuid.NewV7()).String()
	u, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     "testuser_" + unique,
		Email:        "test_" + unique + "@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)
	return u.ID
}

func makeTestUser(id uuid.UUID, username, email string) *db.SharedUser {
	now := time.Now()
	displayName := "Test User"
	isActive := true
	isAdmin := false
	return &db.SharedUser{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=4$test$hash",
		DisplayName:  &displayName,
		IsActive:     &isActive,
		IsAdmin:      &isAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// ============================================================================
// User Management Unit Tests
// ============================================================================

func TestUnit_GetUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		expectedUser := makeTestUser(userID, "testuser", "test@example.com")

		repo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.GetUser(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "testuser", result.Username)
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

		repo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("user not found"))

		svc := setupUnitTestService(t, repo)
		result, err := svc.GetUser(ctx, userID)

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestUnit_GetUserByUsername(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	expectedUser := makeTestUser(userID, "testuser", "test@example.com")

	repo.EXPECT().GetUserByUsername(ctx, "testuser").Return(expectedUser, nil)

	svc := setupUnitTestService(t, repo)
	result, err := svc.GetUserByUsername(ctx, "testuser")

	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
}

func TestUnit_GetUserByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	expectedUser := makeTestUser(userID, "testuser", "test@example.com")

	repo.EXPECT().GetUserByEmail(ctx, "test@example.com").Return(expectedUser, nil)

	svc := setupUnitTestService(t, repo)
	result, err := svc.GetUserByEmail(ctx, "test@example.com")

	require.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestUnit_ListUsers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	filters := user.UserFilters{Limit: 10, Offset: 0}
	users := []db.SharedUser{
		*makeTestUser(uuid.Must(uuid.NewV7()), "user1", "user1@example.com"),
		*makeTestUser(uuid.Must(uuid.NewV7()), "user2", "user2@example.com"),
	}

	repo.EXPECT().ListUsers(ctx, filters).Return(users, int64(2), nil)

	svc := setupUnitTestService(t, repo)
	result, count, err := svc.ListUsers(ctx, filters)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), count)
}

func TestUnit_CreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.CreateUserParams{
			Username:     "newuser",
			Email:        "new@example.com",
			PasswordHash: "password123",
		}
		createdUser := makeTestUser(uuid.Must(uuid.NewV7()), "newuser", "new@example.com")

		repo.EXPECT().GetUserByUsername(ctx, "newuser").Return(nil, errors.New("not found"))
		repo.EXPECT().GetUserByEmail(ctx, "new@example.com").Return(nil, errors.New("not found"))
		repo.EXPECT().CreateUser(ctx, mock.AnythingOfType("user.CreateUserParams")).Return(createdUser, nil)
		repo.EXPECT().UpsertUserPreferences(ctx, mock.AnythingOfType("user.UpsertPreferencesParams")).Return(&db.SharedUserPreference{}, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.CreateUser(ctx, params)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "newuser", result.Username)
	})

	t.Run("missing username", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.CreateUserParams{
			Email:        "new@example.com",
			PasswordHash: "password123",
		}

		svc := setupUnitTestService(t, repo)
		result, err := svc.CreateUser(ctx, params)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("missing email", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.CreateUserParams{
			Username:     "newuser",
			PasswordHash: "password123",
		}

		svc := setupUnitTestService(t, repo)
		result, err := svc.CreateUser(ctx, params)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("missing password", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.CreateUserParams{
			Username: "newuser",
			Email:    "new@example.com",
		}

		svc := setupUnitTestService(t, repo)
		result, err := svc.CreateUser(ctx, params)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("username already exists", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.CreateUserParams{
			Username:     "existing",
			Email:        "new@example.com",
			PasswordHash: "password123",
		}
		existingUser := makeTestUser(uuid.Must(uuid.NewV7()), "existing", "old@example.com")

		repo.EXPECT().GetUserByUsername(ctx, "existing").Return(existingUser, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.CreateUser(ctx, params)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "username already exists")
	})

	t.Run("email already exists", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.CreateUserParams{
			Username:     "newuser",
			Email:        "existing@example.com",
			PasswordHash: "password123",
		}
		existingUser := makeTestUser(uuid.Must(uuid.NewV7()), "otheruser", "existing@example.com")

		repo.EXPECT().GetUserByUsername(ctx, "newuser").Return(nil, errors.New("not found"))
		repo.EXPECT().GetUserByEmail(ctx, "existing@example.com").Return(existingUser, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.CreateUser(ctx, params)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "email already exists")
	})
}

func TestUnit_UpdateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("success", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.UpdateUserParams{
			DisplayName: ptr.To("New Name"),
		}
		oldUser := makeTestUser(userID, "testuser", "test@example.com")
		updatedUser := makeTestUser(userID, "testuser", "test@example.com")
		updatedUser.DisplayName = ptr.To("New Name")

		repo.EXPECT().GetUserByID(ctx, userID).Return(oldUser, nil)
		repo.EXPECT().UpdateUser(ctx, userID, params).Return(updatedUser, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.UpdateUser(ctx, userID, params)

		require.NoError(t, err)
		assert.Equal(t, "New Name", *result.DisplayName)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		params := user.UpdateUserParams{
			DisplayName: ptr.To("New Name"),
		}

		repo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("not found"))

		svc := setupUnitTestService(t, repo)
		result, err := svc.UpdateUser(ctx, userID, params)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUnit_DeleteUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)
	existingUser := makeTestUser(userID, "testuser", "test@example.com")

	repo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)
	repo.EXPECT().DeleteUser(ctx, userID).Return(nil)

	svc := setupUnitTestService(t, repo)
	err := svc.DeleteUser(ctx, userID)

	require.NoError(t, err)
}

func TestUnit_HardDeleteUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)

	repo.EXPECT().DeleteUserPreferences(ctx, userID).Return(nil)
	repo.EXPECT().HardDeleteUser(ctx, userID).Return(nil)

	svc := setupUnitTestService(t, repo)
	err := svc.HardDeleteUser(ctx, userID)

	require.NoError(t, err)
}

func TestUnit_VerifyEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)
	repo.EXPECT().VerifyEmail(ctx, userID).Return(nil)

	svc := setupUnitTestService(t, repo)
	err := svc.VerifyEmail(ctx, userID)

	require.NoError(t, err)
}

func TestUnit_RecordLogin(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)
	repo.EXPECT().UpdateLastLogin(ctx, userID).Return(nil)

	svc := setupUnitTestService(t, repo)
	err := svc.RecordLogin(ctx, userID)

	require.NoError(t, err)
}

// ============================================================================
// Password Unit Tests
// ============================================================================

func TestUnit_HashPassword(t *testing.T) {
	t.Parallel()
	repo := NewMockUserRepository(t)
	svc := setupUnitTestService(t, repo)

	hash, err := svc.HashPassword("testpassword")
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
	assert.NotEqual(t, "testpassword", hash)
}

func TestUnit_VerifyPassword(t *testing.T) {
	t.Parallel()
	repo := NewMockUserRepository(t)
	svc := setupUnitTestService(t, repo)

	hash, err := svc.HashPassword("testpassword")
	require.NoError(t, err)

	t.Run("correct password", func(t *testing.T) {
		err := svc.VerifyPassword(hash, "testpassword")
		require.NoError(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		err := svc.VerifyPassword(hash, "wrongpassword")
		require.Error(t, err)
	})
}

func TestUnit_UpdatePassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("success", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		// Hash the old password
		oldHash, _ := svc.HashPassword("oldpassword")
		existingUser := makeTestUser(userID, "testuser", "test@example.com")
		existingUser.PasswordHash = oldHash

		repo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)
		repo.EXPECT().UpdatePassword(ctx, userID, mock.AnythingOfType("string")).Return(nil)

		err := svc.UpdatePassword(ctx, userID, "oldpassword", "newpassword")
		require.NoError(t, err)
	})

	t.Run("wrong old password", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		oldHash, _ := svc.HashPassword("oldpassword")
		existingUser := makeTestUser(userID, "testuser", "test@example.com")
		existingUser.PasswordHash = oldHash

		repo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)

		err := svc.UpdatePassword(ctx, userID, "wrongpassword", "newpassword")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid old password")
	})

	t.Run("user not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)

		repo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("not found"))

		svc := setupUnitTestService(t, repo)
		err := svc.UpdatePassword(ctx, userID, "old", "new")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

// ============================================================================
// User Preferences Unit Tests
// ============================================================================

func TestUnit_GetUserPreferences(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("existing preferences", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		prefs := &db.SharedUserPreference{UserID: userID}

		repo.EXPECT().GetUserPreferences(ctx, userID).Return(prefs, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.GetUserPreferences(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, userID, result.UserID)
	})

	t.Run("create default when not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		prefs := &db.SharedUserPreference{UserID: userID}

		repo.EXPECT().GetUserPreferences(ctx, userID).Return(nil, errors.New("not found"))
		repo.EXPECT().UpsertUserPreferences(ctx, mock.AnythingOfType("user.UpsertPreferencesParams")).Return(prefs, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.GetUserPreferences(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, userID, result.UserID)
	})
}

func TestUnit_UpdateUserPreferences(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("valid theme", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		theme := "dark"
		prefs := &db.SharedUserPreference{UserID: userID, Theme: &theme}

		repo.EXPECT().UpsertUserPreferences(ctx, mock.AnythingOfType("user.UpsertPreferencesParams")).Return(prefs, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.UpdateUserPreferences(ctx, user.UpsertPreferencesParams{
			UserID: userID,
			Theme:  ptr.To("dark"),
		})

		require.NoError(t, err)
		assert.Equal(t, "dark", *result.Theme)
	})

	t.Run("invalid theme", func(t *testing.T) {
		repo := NewMockUserRepository(t)

		svc := setupUnitTestService(t, repo)
		result, err := svc.UpdateUserPreferences(ctx, user.UpsertPreferencesParams{
			UserID: userID,
			Theme:  ptr.To("invalid"),
		})

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid theme")
	})

	t.Run("invalid visibility", func(t *testing.T) {
		repo := NewMockUserRepository(t)

		svc := setupUnitTestService(t, repo)
		result, err := svc.UpdateUserPreferences(ctx, user.UpsertPreferencesParams{
			UserID:            userID,
			ProfileVisibility: ptr.To("invalid"),
		})

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid profile visibility")
	})
}

func TestUnit_UpdateNotificationPreferences(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)

	var capturedParams user.UpsertPreferencesParams
	repo.EXPECT().UpsertUserPreferences(ctx, mock.AnythingOfType("user.UpsertPreferencesParams")).
		Run(func(_ context.Context, params user.UpsertPreferencesParams) {
			capturedParams = params
		}).
		Return(&db.SharedUserPreference{UserID: userID}, nil)

	svc := setupUnitTestService(t, repo)

	// All three parameters are *NotificationSettings
	emailSettings := &user.NotificationSettings{Enabled: true, Frequency: "daily"}
	pushSettings := &user.NotificationSettings{Enabled: false}
	digestSettings := &user.NotificationSettings{Enabled: true, Frequency: "weekly"}

	err := svc.UpdateNotificationPreferences(ctx, userID, emailSettings, pushSettings, digestSettings)
	require.NoError(t, err)

	// Verify the captured params
	require.NotNil(t, capturedParams.EmailNotifications)
	require.NotNil(t, capturedParams.PushNotifications)
	require.NotNil(t, capturedParams.DigestNotifications)

	var email user.NotificationSettings
	err = json.Unmarshal(*capturedParams.EmailNotifications, &email)
	require.NoError(t, err)
	assert.True(t, email.Enabled)
	assert.Equal(t, "daily", email.Frequency)
}

// ============================================================================
// Avatar Unit Tests
// ============================================================================

func TestUnit_GetCurrentAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	avatarID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockUserRepository(t)
	avatar := &db.SharedUserAvatar{
		ID:       avatarID,
		UserID:   userID,
		FilePath: "/api/v1/files/avatars/test.png",
	}

	repo.EXPECT().GetCurrentAvatar(ctx, userID).Return(avatar, nil)

	svc := setupUnitTestService(t, repo)
	result, err := svc.GetCurrentAvatar(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, avatarID, result.ID)
}

func TestUnit_ListUserAvatars(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("default limit when 0", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		avatars := []db.SharedUserAvatar{{}, {}}

		repo.EXPECT().ListUserAvatars(ctx, userID, int32(10), int32(0)).Return(avatars, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.ListUserAvatars(ctx, userID, 0, 0)

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("max limit capped at 100", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		avatars := []db.SharedUserAvatar{}

		repo.EXPECT().ListUserAvatars(ctx, userID, int32(100), int32(0)).Return(avatars, nil)

		svc := setupUnitTestService(t, repo)
		result, err := svc.ListUserAvatars(ctx, userID, 200, 0)

		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestUnit_SetCurrentAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	avatarID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	otherUserID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	t.Run("success", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		avatar := &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path/to/avatar.png"}
		updatedUser := makeTestUser(userID, "testuser", "test@example.com")

		repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(avatar, nil)
		repo.EXPECT().UnsetCurrentAvatars(ctx, userID).Return(nil)
		repo.EXPECT().SetCurrentAvatar(ctx, avatarID).Return(nil)
		repo.EXPECT().UpdateUser(ctx, userID, mock.AnythingOfType("user.UpdateUserParams")).Return(updatedUser, nil)

		svc := setupUnitTestService(t, repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)

		require.NoError(t, err)
	})

	t.Run("avatar not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)

		repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(nil, errors.New("not found"))

		svc := setupUnitTestService(t, repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar not found")
	})

	t.Run("avatar belongs to other user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		avatar := &db.SharedUserAvatar{ID: avatarID, UserID: otherUserID}

		repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(avatar, nil)

		svc := setupUnitTestService(t, repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar does not belong to user")
	})
}

func TestUnit_DeleteAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	avatarID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	otherUserID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	t.Run("success", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		avatar := &db.SharedUserAvatar{ID: avatarID, UserID: userID}

		repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(avatar, nil)
		repo.EXPECT().DeleteAvatar(ctx, avatarID).Return(nil)

		svc := setupUnitTestService(t, repo)
		err := svc.DeleteAvatar(ctx, userID, avatarID)

		require.NoError(t, err)
	})

	t.Run("avatar not found", func(t *testing.T) {
		repo := NewMockUserRepository(t)

		repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(nil, errors.New("not found"))

		svc := setupUnitTestService(t, repo)
		err := svc.DeleteAvatar(ctx, userID, avatarID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar not found")
	})

	t.Run("avatar belongs to other user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		avatar := &db.SharedUserAvatar{ID: avatarID, UserID: otherUserID}

		repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(avatar, nil)

		svc := setupUnitTestService(t, repo)
		err := svc.DeleteAvatar(ctx, userID, avatarID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar does not belong to user")
	})
}

func TestUnit_UploadAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("invalid file size", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 10 * 1024 * 1024, // 10MB - exceeds 5MB limit
			MimeType:      "image/png",
			Width:         256,
			Height:        256,
		}

		result, err := svc.UploadAvatar(ctx, userID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "exceeds maximum")
	})

	t.Run("invalid mime type", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.bmp",
			FileSizeBytes: 1024,
			MimeType:      "image/bmp", // BMP not allowed
			Width:         256,
			Height:        256,
		}

		result, err := svc.UploadAvatar(ctx, userID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid MIME type")
	})

	t.Run("invalid width too small", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         8, // Too small (min 16)
			Height:        256,
		}

		result, err := svc.UploadAvatar(ctx, userID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid width")
	})

	t.Run("invalid width too large", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         5000, // Too large (max 4096)
			Height:        256,
		}

		result, err := svc.UploadAvatar(ctx, userID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid width")
	})

	t.Run("invalid height too small", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         256,
			Height:        8, // Too small (min 16)
		}

		result, err := svc.UploadAvatar(ctx, userID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid height")
	})

	t.Run("invalid height too large", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         256,
			Height:        5000, // Too large (max 4096)
		}

		result, err := svc.UploadAvatar(ctx, userID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid height")
	})

	// NOTE: UploadAvatar uses s.pool directly (not the repository mock) for transaction
	// support. Valid upload tests need a real user in the database for FK constraints.

	t.Run("valid jpeg upload", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc, testDB := setupUnitTestServiceWithDB(t, repo)
		dbUserID := createDBUser(t, testDB)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.jpg",
			FileSizeBytes: 1024,
			MimeType:      "image/jpeg",
			Width:         256,
			Height:        256,
		}

		result, err := svc.UploadAvatar(ctx, dbUserID, strings.NewReader("test"), metadata)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, dbUserID, result.UserID)
	})

	t.Run("valid png upload", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc, testDB := setupUnitTestServiceWithDB(t, repo)
		dbUserID := createDBUser(t, testDB)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 2048,
			MimeType:      "image/png",
			Width:         512,
			Height:        512,
		}

		result, err := svc.UploadAvatar(ctx, dbUserID, strings.NewReader("test image data"), metadata)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("valid gif upload with animation", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		isAnimated := true
		svc, testDB := setupUnitTestServiceWithDB(t, repo)
		dbUserID := createDBUser(t, testDB)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.gif",
			FileSizeBytes: 4096,
			MimeType:      "image/gif",
			Width:         128,
			Height:        128,
			IsAnimated:    &isAnimated,
		}

		result, err := svc.UploadAvatar(ctx, dbUserID, strings.NewReader("test gif data"), metadata)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("valid webp upload", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc, testDB := setupUnitTestServiceWithDB(t, repo)
		dbUserID := createDBUser(t, testDB)

		metadata := user.AvatarMetadata{
			FileName:      "avatar.webp",
			FileSizeBytes: 2000,
			MimeType:      "image/webp",
			Width:         400,
			Height:        400,
		}

		result, err := svc.UploadAvatar(ctx, dbUserID, strings.NewReader("test webp data"), metadata)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	// NOTE: Error injection subtests for UploadAvatar can't use mock expectations
	// because the method uses s.pool directly (not the repository) for transaction
	// support. Test FK constraint violation as proxy for DB error handling.
	t.Run("error on non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository(t)
		svc := setupUnitTestService(t, repo)
		nonExistentUserID := uuid.Must(uuid.NewV7())

		metadata := user.AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         256,
			Height:        256,
		}

		result, err := svc.UploadAvatar(ctx, nonExistentUserID, strings.NewReader("test"), metadata)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create avatar")
	})
}

func TestUnit_SetCurrentAvatar_UnsetError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	avatarID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockUserRepository(t)
	avatar := &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path/to/avatar.png"}

	repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(avatar, nil)
	repo.EXPECT().UnsetCurrentAvatars(ctx, userID).Return(errors.New("db error"))

	svc := setupUnitTestService(t, repo)
	err := svc.SetCurrentAvatar(ctx, userID, avatarID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unset current avatars")
}

func TestUnit_SetCurrentAvatar_SetError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	avatarID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	repo := NewMockUserRepository(t)
	avatar := &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path/to/avatar.png"}

	repo.EXPECT().GetAvatarByID(ctx, avatarID).Return(avatar, nil)
	repo.EXPECT().UnsetCurrentAvatars(ctx, userID).Return(nil)
	repo.EXPECT().SetCurrentAvatar(ctx, avatarID).Return(errors.New("db error"))

	svc := setupUnitTestService(t, repo)
	err := svc.SetCurrentAvatar(ctx, userID, avatarID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set current avatar")
}

func TestUnit_DeleteUser_NotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)

	repo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("not found"))
	repo.EXPECT().DeleteUser(ctx, userID).Return(nil)

	svc := setupUnitTestService(t, repo)
	err := svc.DeleteUser(ctx, userID)

	// Should succeed even if user not found for logging
	require.NoError(t, err)
}

func TestUnit_CreateUser_RepoError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	params := user.CreateUserParams{
		Username:     "newuser",
		Email:        "new@example.com",
		PasswordHash: "password123",
	}

	repo.EXPECT().GetUserByUsername(ctx, "newuser").Return(nil, errors.New("not found"))
	repo.EXPECT().GetUserByEmail(ctx, "new@example.com").Return(nil, errors.New("not found"))
	repo.EXPECT().CreateUser(ctx, mock.AnythingOfType("user.CreateUserParams")).Return(nil, errors.New("db error"))

	svc := setupUnitTestService(t, repo)
	result, err := svc.CreateUser(ctx, params)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestUnit_UpdateUser_RepoError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := NewMockUserRepository(t)
	params := user.UpdateUserParams{
		DisplayName: ptr.To("New Name"),
	}
	oldUser := makeTestUser(userID, "testuser", "test@example.com")

	repo.EXPECT().GetUserByID(ctx, userID).Return(oldUser, nil)
	repo.EXPECT().UpdateUser(ctx, userID, params).Return(nil, errors.New("db error"))

	svc := setupUnitTestService(t, repo)
	result, err := svc.UpdateUser(ctx, userID, params)

	require.Error(t, err)
	require.Nil(t, result)
}
