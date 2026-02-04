package user

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/testutil"
)

// TestMain is in repository_pg_test.go

func setupTestService(t *testing.T) (*Service, *testutil.TestDB) {
	t.Helper()
	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewPostgresRepository(queries)
	activityLogger := activity.NewNoopLogger()
	svc := NewService(repo, activityLogger)
	return svc, testDB
}

// ============================================================================
// User Management Tests
// ============================================================================

func TestService_CreateUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	t.Run("valid user", func(t *testing.T) {
		user, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "password123",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		// Password should be hashed with Argon2id
		assert.NotEqual(t, "password123", user.PasswordHash)
		assert.True(t, strings.HasPrefix(user.PasswordHash, "$argon2id$"))
	})

	t.Run("missing username", func(t *testing.T) {
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Email:        "test2@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("missing email", func(t *testing.T) {
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "testuser2",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("missing password", func(t *testing.T) {
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username: "testuser3",
			Email:    "test3@example.com",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("duplicate username", func(t *testing.T) {
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "dupuser",
			Email:        "dup1@example.com",
			PasswordHash: "password123",
		})
		require.NoError(t, err)

		_, err = svc.CreateUser(ctx, CreateUserParams{
			Username:     "dupuser",
			Email:        "dup2@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})

	t.Run("duplicate email", func(t *testing.T) {
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "emaildup1",
			Email:        "duplicate@example.com",
			PasswordHash: "password123",
		})
		require.NoError(t, err)

		_, err = svc.CreateUser(ctx, CreateUserParams{
			Username:     "emaildup2",
			Email:        "duplicate@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})
}

func TestService_GetUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create user
	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "getuser",
		Email:        "getuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("existing user", func(t *testing.T) {
		user, err := svc.GetUser(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, "getuser", user.Username)
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := svc.GetUser(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestService_GetUserByUsername(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "byusername",
		Email:        "byusername@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("existing username", func(t *testing.T) {
		user, err := svc.GetUserByUsername(ctx, "byusername")
		require.NoError(t, err)
		assert.Equal(t, "byusername", user.Username)
	})

	t.Run("non-existent username", func(t *testing.T) {
		_, err := svc.GetUserByUsername(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestService_GetUserByEmail(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "byemail",
		Email:        "byemail@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("existing email", func(t *testing.T) {
		user, err := svc.GetUserByEmail(ctx, "byemail@example.com")
		require.NoError(t, err)
		assert.Equal(t, "byemail@example.com", user.Email)
	})

	t.Run("non-existent email", func(t *testing.T) {
		_, err := svc.GetUserByEmail(ctx, "nonexistent@example.com")
		require.Error(t, err)
	})
}

func TestService_ListUsers(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create users with is_active=false, is_admin=false (default filter)
	for i := 0; i < 5; i++ {
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "listuser" + string(rune('0'+i)),
			Email:        "listuser" + string(rune('0'+i)) + "@example.com",
			PasswordHash: "password123",
			IsActive:     ptr(false),
			IsAdmin:      ptr(false),
		})
		require.NoError(t, err)
	}

	users, count, err := svc.ListUsers(ctx, UserFilters{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 5)
	assert.GreaterOrEqual(t, count, int64(5))
}

func TestService_UpdateUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "updateuser",
		Email:        "updateuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("update display name", func(t *testing.T) {
		displayName := "Updated Name"
		updated, err := svc.UpdateUser(ctx, created.ID, UpdateUserParams{
			DisplayName: &displayName,
		})
		require.NoError(t, err)
		require.NotNil(t, updated.DisplayName)
		assert.Equal(t, "Updated Name", *updated.DisplayName)
	})

	t.Run("update non-existent user", func(t *testing.T) {
		displayName := "Name"
		_, err := svc.UpdateUser(ctx, uuid.New(), UpdateUserParams{
			DisplayName: &displayName,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestService_UpdatePassword(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "pwduser",
		Email:        "pwduser@example.com",
		PasswordHash: "oldpassword",
	})
	require.NoError(t, err)

	t.Run("valid password change", func(t *testing.T) {
		err := svc.UpdatePassword(ctx, created.ID, "oldpassword", "newpassword")
		require.NoError(t, err)

		// Verify new password works
		user, err := svc.GetUser(ctx, created.ID)
		require.NoError(t, err)
		err = svc.VerifyPassword(user.PasswordHash, "newpassword")
		require.NoError(t, err)
	})

	t.Run("wrong old password", func(t *testing.T) {
		err := svc.UpdatePassword(ctx, created.ID, "wrongpassword", "newpassword2")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid old password")
	})

	t.Run("non-existent user", func(t *testing.T) {
		err := svc.UpdatePassword(ctx, uuid.New(), "old", "new")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestService_DeleteUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "deleteuser",
		Email:        "deleteuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	err = svc.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	// User should not be found (soft deleted)
	_, err = svc.GetUser(ctx, created.ID)
	require.Error(t, err)
}

func TestService_HardDeleteUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "harddeleteuser",
		Email:        "harddeleteuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	err = svc.HardDeleteUser(ctx, created.ID)
	require.NoError(t, err)

	_, err = svc.GetUser(ctx, created.ID)
	require.Error(t, err)
}

func TestService_VerifyEmail(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "verifyemail",
		Email:        "verifyemail@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	err = svc.VerifyEmail(ctx, created.ID)
	require.NoError(t, err)

	user, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, user.EmailVerified)
	assert.True(t, *user.EmailVerified)
}

func TestService_RecordLogin(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "loginuser",
		Email:        "loginuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)
	assert.False(t, created.LastLoginAt.Valid)

	err = svc.RecordLogin(ctx, created.ID)
	require.NoError(t, err)

	user, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, user.LastLoginAt.Valid)
}

// ============================================================================
// Password Tests
// ============================================================================

func TestService_HashPassword(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

	hash, err := svc.HashPassword("testpassword")
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
	assert.NotEqual(t, "testpassword", hash)
}

func TestService_VerifyPassword(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

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

// ============================================================================
// User Preferences Tests
// ============================================================================

func TestService_GetUserPreferences(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "prefsuser",
		Email:        "prefsuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Should return defaults if not explicitly set
	prefs, err := svc.GetUserPreferences(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, prefs.UserID)
}

func TestService_UpdateUserPreferences(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "updateprefs",
		Email:        "updateprefs@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("valid preferences", func(t *testing.T) {
		theme := "dark"
		prefs, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
			UserID: created.ID,
			Theme:  &theme,
		})
		require.NoError(t, err)
		assert.Equal(t, "dark", *prefs.Theme)
	})

	t.Run("invalid theme", func(t *testing.T) {
		theme := "invalid"
		_, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
			UserID: created.ID,
			Theme:  &theme,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid theme")
	})

	t.Run("invalid visibility", func(t *testing.T) {
		vis := "invalid"
		_, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
			UserID:            created.ID,
			ProfileVisibility: &vis,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid profile visibility")
	})
}

func TestService_UpdateNotificationPreferences(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "notifprefs",
		Email:        "notifprefs@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	emailSettings := NotificationSettings{Enabled: true, Frequency: "daily"}
	pushSettings := NotificationSettings{Enabled: false}

	err = svc.UpdateNotificationPreferences(ctx, created.ID, &emailSettings, &pushSettings, nil)
	require.NoError(t, err)

	prefs, err := svc.GetUserPreferences(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, prefs.EmailNotifications)

	var email NotificationSettings
	err = json.Unmarshal(prefs.EmailNotifications, &email)
	require.NoError(t, err)
	assert.True(t, email.Enabled)
	assert.Equal(t, "daily", email.Frequency)
}

// ============================================================================
// Avatar Tests
// ============================================================================

func TestService_GetCurrentAvatar(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "avataruser",
		Email:        "avataruser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// No avatar initially
	_, err = svc.GetCurrentAvatar(ctx, created.ID)
	require.Error(t, err)

	// Upload avatar
	_, err = svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
		FileName:      "avatar.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	// Now should have current avatar
	avatar, err := svc.GetCurrentAvatar(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, avatar.UserID)
}

func TestService_ListUserAvatars(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "listavatars",
		Email:        "listavatars@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Create 3 avatars
	for i := 0; i < 3; i++ {
		_, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         100,
			Height:        100,
		})
		require.NoError(t, err)
	}

	avatars, err := svc.ListUserAvatars(ctx, created.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, avatars, 3)
}

func TestService_ListUserAvatars_Limits(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "limittestuser",
		Email:        "limittestuser@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Test default limit when 0 is passed
	avatars, err := svc.ListUserAvatars(ctx, created.ID, 0, 0)
	require.NoError(t, err)
	assert.NotNil(t, avatars)

	// Test max limit cap (should be capped at 100)
	avatars, err = svc.ListUserAvatars(ctx, created.ID, 500, 0)
	require.NoError(t, err)
	assert.NotNil(t, avatars)
}

func TestService_UploadAvatar(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "uploadavatar",
		Email:        "uploadavatar@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	t.Run("valid upload", func(t *testing.T) {
		avatar, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "avatar.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         100,
			Height:        100,
		})
		require.NoError(t, err)
		assert.Equal(t, int32(1), avatar.Version)
	})

	t.Run("version increments", func(t *testing.T) {
		avatar, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "avatar2.png",
			FileSizeBytes: 2048,
			MimeType:      "image/png",
			Width:         200,
			Height:        200,
		})
		require.NoError(t, err)
		assert.Equal(t, int32(2), avatar.Version)
	})

	t.Run("file too large", func(t *testing.T) {
		_, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "big.png",
			FileSizeBytes: 10 * 1024 * 1024, // 10MB
			MimeType:      "image/png",
			Width:         100,
			Height:        100,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum")
	})

	t.Run("invalid mime type", func(t *testing.T) {
		_, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "doc.pdf",
			FileSizeBytes: 1024,
			MimeType:      "application/pdf",
			Width:         100,
			Height:        100,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid MIME type")
	})

	t.Run("dimensions too small", func(t *testing.T) {
		_, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "tiny.png",
			FileSizeBytes: 100,
			MimeType:      "image/png",
			Width:         5,
			Height:        5,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid width")
	})

	t.Run("dimensions too large", func(t *testing.T) {
		_, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
			FileName:      "huge.png",
			FileSizeBytes: 1024,
			MimeType:      "image/png",
			Width:         10000,
			Height:        100,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid width")
	})
}

func TestService_SetCurrentAvatar(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "setcurrent",
		Email:        "setcurrent@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	avatar1, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
		FileName:      "a1.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	avatar2, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
		FileName:      "a2.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	// avatar2 should be current now
	current, err := svc.GetCurrentAvatar(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar2.ID, current.ID)

	// Switch to avatar1
	err = svc.SetCurrentAvatar(ctx, created.ID, avatar1.ID)
	require.NoError(t, err)

	current, err = svc.GetCurrentAvatar(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar1.ID, current.ID)
}

func TestService_SetCurrentAvatar_WrongUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user1, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	user2, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "user2",
		Email:        "user2@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	avatar1, err := svc.UploadAvatar(ctx, user1.ID, nil, AvatarMetadata{
		FileName:      "a1.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	// user2 tries to set user1's avatar as current
	err = svc.SetCurrentAvatar(ctx, user2.ID, avatar1.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to user")
}

func TestService_DeleteAvatar(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "delavatar",
		Email:        "delavatar@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	avatar, err := svc.UploadAvatar(ctx, created.ID, nil, AvatarMetadata{
		FileName:      "del.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	err = svc.DeleteAvatar(ctx, created.ID, avatar.ID)
	require.NoError(t, err)
}

func TestService_DeleteAvatar_WrongUser(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user1, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "deluser1",
		Email:        "deluser1@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	user2, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "deluser2",
		Email:        "deluser2@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	avatar1, err := svc.UploadAvatar(ctx, user1.ID, nil, AvatarMetadata{
		FileName:      "a1.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	// user2 tries to delete user1's avatar
	err = svc.DeleteAvatar(ctx, user2.ID, avatar1.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to user")
}

func TestService_DeleteAvatar_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "delnotfound",
		Email:        "delnotfound@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	err = svc.DeleteAvatar(ctx, created.ID, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "avatar not found")
}

// ============================================================================
// Validation Tests
// ============================================================================

func TestService_ValidatePreferences(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

	tests := []struct {
		name    string
		params  UpsertPreferencesParams
		wantErr string
	}{
		{
			name:    "valid theme light",
			params:  UpsertPreferencesParams{Theme: ptr("light")},
			wantErr: "",
		},
		{
			name:    "valid theme dark",
			params:  UpsertPreferencesParams{Theme: ptr("dark")},
			wantErr: "",
		},
		{
			name:    "valid theme system",
			params:  UpsertPreferencesParams{Theme: ptr("system")},
			wantErr: "",
		},
		{
			name:    "invalid theme",
			params:  UpsertPreferencesParams{Theme: ptr("rainbow")},
			wantErr: "invalid theme",
		},
		{
			name:    "valid visibility public",
			params:  UpsertPreferencesParams{ProfileVisibility: ptr("public")},
			wantErr: "",
		},
		{
			name:    "valid visibility friends",
			params:  UpsertPreferencesParams{ProfileVisibility: ptr("friends")},
			wantErr: "",
		},
		{
			name:    "valid visibility private",
			params:  UpsertPreferencesParams{ProfileVisibility: ptr("private")},
			wantErr: "",
		},
		{
			name:    "invalid visibility",
			params:  UpsertPreferencesParams{ProfileVisibility: ptr("hidden")},
			wantErr: "invalid profile visibility",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validatePreferences(tt.params)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestService_ValidateAvatarMetadata(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

	tests := []struct {
		name     string
		metadata AvatarMetadata
		wantErr  string
	}{
		{
			name: "valid jpeg",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/jpeg",
				Width:         100,
				Height:        100,
			},
			wantErr: "",
		},
		{
			name: "valid png",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/png",
				Width:         100,
				Height:        100,
			},
			wantErr: "",
		},
		{
			name: "valid gif",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/gif",
				Width:         100,
				Height:        100,
			},
			wantErr: "",
		},
		{
			name: "valid webp",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/webp",
				Width:         100,
				Height:        100,
			},
			wantErr: "",
		},
		{
			name: "too large",
			metadata: AvatarMetadata{
				FileSizeBytes: 10 * 1024 * 1024,
				MimeType:      "image/png",
				Width:         100,
				Height:        100,
			},
			wantErr: "exceeds maximum",
		},
		{
			name: "invalid mime",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "video/mp4",
				Width:         100,
				Height:        100,
			},
			wantErr: "invalid MIME type",
		},
		{
			name: "width too small",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/png",
				Width:         10,
				Height:        100,
			},
			wantErr: "invalid width",
		},
		{
			name: "width too large",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/png",
				Width:         5000,
				Height:        100,
			},
			wantErr: "invalid width",
		},
		{
			name: "height too small",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/png",
				Width:         100,
				Height:        10,
			},
			wantErr: "invalid height",
		},
		{
			name: "height too large",
			metadata: AvatarMetadata{
				FileSizeBytes: 1024,
				MimeType:      "image/png",
				Width:         100,
				Height:        5000,
			},
			wantErr: "invalid height",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validateAvatarMetadata(tt.metadata)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
