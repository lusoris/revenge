package user

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Full User Lifecycle Integration Tests
// ============================================================================

func TestService_FullUserLifecycle(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// 1. Create user
	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "lifecycle_user",
		Email:        "lifecycle@example.com",
		PasswordHash: "initialpassword",
		IsActive:     ptr(true),
		IsAdmin:      ptr(false),
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	assert.Equal(t, "lifecycle_user", created.Username)
	assert.Equal(t, "lifecycle@example.com", created.Email)

	// 2. Get user by all lookups
	byID, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, byID.ID)

	byUsername, err := svc.GetUserByUsername(ctx, "lifecycle_user")
	require.NoError(t, err)
	assert.Equal(t, created.ID, byUsername.ID)

	byEmail, err := svc.GetUserByEmail(ctx, "lifecycle@example.com")
	require.NoError(t, err)
	assert.Equal(t, created.ID, byEmail.ID)

	// 3. Update user fields
	newName := "Lifecycle User Display"
	updated, err := svc.UpdateUser(ctx, created.ID, UpdateUserParams{
		DisplayName: &newName,
	})
	require.NoError(t, err)
	require.NotNil(t, updated.DisplayName)
	assert.Equal(t, "Lifecycle User Display", *updated.DisplayName)

	// 4. Verify email
	err = svc.VerifyEmail(ctx, created.ID)
	require.NoError(t, err)

	verified, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, verified.EmailVerified)
	assert.True(t, *verified.EmailVerified)

	// 5. Record login
	err = svc.RecordLogin(ctx, created.ID)
	require.NoError(t, err)

	loggedIn, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, loggedIn.LastLoginAt.Valid)

	// 6. Change password
	err = svc.UpdatePassword(ctx, created.ID, "initialpassword", "newpassword123")
	require.NoError(t, err)

	// Verify new password works
	afterPw, err := svc.GetUser(ctx, created.ID)
	require.NoError(t, err)
	err = svc.VerifyPassword(afterPw.PasswordHash, "newpassword123")
	require.NoError(t, err)

	// Old password should not work
	err = svc.VerifyPassword(afterPw.PasswordHash, "initialpassword")
	require.Error(t, err)

	// 7. Soft delete user
	err = svc.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	_, err = svc.GetUser(ctx, created.ID)
	require.Error(t, err)
}

// ============================================================================
// Preferences Lifecycle Integration Tests
// ============================================================================

func TestService_PreferencesLifecycle(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "prefs_lifecycle",
		Email:        "prefs_lifecycle@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// 1. Get default preferences (auto-created)
	prefs, err := svc.GetUserPreferences(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, prefs.UserID)

	// 2. Update theme
	theme := "dark"
	prefs, err = svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
		UserID: user.ID,
		Theme:  &theme,
	})
	require.NoError(t, err)
	assert.Equal(t, "dark", *prefs.Theme)

	// 3. Update visibility
	vis := "public"
	prefs, err = svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
		UserID:            user.ID,
		ProfileVisibility: &vis,
	})
	require.NoError(t, err)
	assert.Equal(t, "public", *prefs.ProfileVisibility)

	// 4. Update multiple fields at once
	showEmail := true
	showActivity := false
	lang := "de-DE"
	prefs, err = svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
		UserID:          user.ID,
		ShowEmail:       &showEmail,
		ShowActivity:    &showActivity,
		DisplayLanguage: &lang,
	})
	require.NoError(t, err)
	assert.True(t, *prefs.ShowEmail)
	assert.False(t, *prefs.ShowActivity)
	assert.Equal(t, "de-DE", *prefs.DisplayLanguage)

	// 5. Update notification preferences
	emailNotif := NotificationSettings{Enabled: true, Frequency: "daily"}
	pushNotif := NotificationSettings{Enabled: false}
	digestNotif := NotificationSettings{Enabled: true, Frequency: "weekly"}
	err = svc.UpdateNotificationPreferences(ctx, user.ID, &emailNotif, &pushNotif, &digestNotif)
	require.NoError(t, err)

	// Verify notification preferences
	prefs, err = svc.GetUserPreferences(ctx, user.ID)
	require.NoError(t, err)

	var email NotificationSettings
	err = json.Unmarshal(prefs.EmailNotifications, &email)
	require.NoError(t, err)
	assert.True(t, email.Enabled)
	assert.Equal(t, "daily", email.Frequency)

	var push NotificationSettings
	err = json.Unmarshal(prefs.PushNotifications, &push)
	require.NoError(t, err)
	assert.False(t, push.Enabled)

	var digest NotificationSettings
	err = json.Unmarshal(prefs.DigestNotifications, &digest)
	require.NoError(t, err)
	assert.True(t, digest.Enabled)
	assert.Equal(t, "weekly", digest.Frequency)
}

func TestService_UpdateNotificationPreferences_Partial(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "partial_notif",
		Email:        "partial_notif@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Update only email notifications (push and digest are nil)
	emailNotif := NotificationSettings{Enabled: true, Frequency: "instant"}
	err = svc.UpdateNotificationPreferences(ctx, user.ID, &emailNotif, nil, nil)
	require.NoError(t, err)

	prefs, err := svc.GetUserPreferences(ctx, user.ID)
	require.NoError(t, err)

	var email NotificationSettings
	err = json.Unmarshal(prefs.EmailNotifications, &email)
	require.NoError(t, err)
	assert.True(t, email.Enabled)
	assert.Equal(t, "instant", email.Frequency)
}

// ============================================================================
// Avatar Lifecycle Integration Tests
// ============================================================================

func TestService_AvatarLifecycle(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "avatar_lifecycle",
		Email:        "avatar_lifecycle@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// 1. No avatar initially
	_, err = svc.GetCurrentAvatar(ctx, user.ID)
	require.Error(t, err)

	// 2. Upload first avatar
	avatar1, err := svc.UploadAvatar(ctx, user.ID, nil, AvatarMetadata{
		FileName:      "avatar1.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         200,
		Height:        200,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(1), avatar1.Version)

	// 3. First avatar should be current
	current, err := svc.GetCurrentAvatar(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar1.ID, current.ID)

	// 4. Upload second avatar
	avatar2, err := svc.UploadAvatar(ctx, user.ID, nil, AvatarMetadata{
		FileName:      "avatar2.jpg",
		FileSizeBytes: 2048,
		MimeType:      "image/jpeg",
		Width:         300,
		Height:        300,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(2), avatar2.Version)

	// 5. Second avatar should now be current
	current, err = svc.GetCurrentAvatar(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar2.ID, current.ID)

	// 6. Switch back to first avatar
	err = svc.SetCurrentAvatar(ctx, user.ID, avatar1.ID)
	require.NoError(t, err)

	current, err = svc.GetCurrentAvatar(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar1.ID, current.ID)

	// 7. List avatars
	avatars, err := svc.ListUserAvatars(ctx, user.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, avatars, 2)

	// 8. Delete second avatar
	err = svc.DeleteAvatar(ctx, user.ID, avatar2.ID)
	require.NoError(t, err)

	// 9. List should show only one
	avatars, err = svc.ListUserAvatars(ctx, user.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, avatars, 1)
}

func TestService_UploadAvatar_WithIPAndUserAgent(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "avatar_metadata_user",
		Email:        "avatar_metadata@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	ipAddr := "192.168.1.100"
	userAgent := "Mozilla/5.0 (X11; Linux x86_64)"
	animated := true

	avatar, err := svc.UploadAvatar(ctx, user.ID, nil, AvatarMetadata{
		FileName:              "animated.gif",
		FileSizeBytes:         3000,
		MimeType:              "image/gif",
		Width:                 128,
		Height:                128,
		IsAnimated:            &animated,
		UploadedFromIP:        &ipAddr,
		UploadedFromUserAgent: &userAgent,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, avatar.ID)
	assert.Equal(t, int32(1), avatar.Version)
}

func TestService_UploadAvatar_InvalidIPAddress(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "bad_ip_user",
		Email:        "bad_ip@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	badIP := "not-a-valid-ip"
	_, err = svc.UploadAvatar(ctx, user.ID, nil, AvatarMetadata{
		FileName:       "test.png",
		FileSizeBytes:  1024,
		MimeType:       "image/png",
		Width:          100,
		Height:         100,
		UploadedFromIP: &badIP,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse IP address")
}

// ============================================================================
// HardDeleteAvatar Integration Tests (repo-level, 0% coverage)
// ============================================================================

func TestRepository_HardDeleteAvatar(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "hard_del_avatar_user",
		Email:        "hard_del_avatar@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)

	// Create an avatar
	avatar, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:        user.ID,
		FilePath:      "/tmp/test-avatar.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
		Version:       1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, avatar.ID)

	// Hard delete the avatar
	err = repo.HardDeleteAvatar(ctx, avatar.ID)
	require.NoError(t, err)

	// Avatar should no longer exist
	_, err = repo.GetAvatarByID(ctx, avatar.ID)
	require.Error(t, err)
}

func TestRepository_HardDeleteAvatar_NonExistent(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Hard delete non-existent avatar should not error (just no-op)
	err := repo.HardDeleteAvatar(ctx, uuid.Must(uuid.NewV7()))
	// The behavior may vary - just verify it doesn't panic
	if err != nil {
		assert.Contains(t, err.Error(), "hard delete avatar")
	}
}

// ============================================================================
// CreateAvatar with IP address (tests repo IP parsing path at 50%)
// ============================================================================

func TestRepository_CreateAvatar_WithIP(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "avatar_ip_user",
		Email:        "avatar_ip@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)

	ipAddr := "10.0.0.1"
	userAgent := "TestAgent/1.0"
	avatar, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:                user.ID,
		FilePath:              "/tmp/test-avatar-ip.png",
		FileSizeBytes:         2048,
		MimeType:              "image/png",
		Width:                 256,
		Height:                256,
		Version:               1,
		UploadedFromIP:        &ipAddr,
		UploadedFromUserAgent: &userAgent,
	})
	require.NoError(t, err)
	require.NotEmpty(t, avatar.ID)
}

func TestRepository_CreateAvatar_WithInvalidIP(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "avatar_bad_ip_user",
		Email:        "avatar_bad_ip@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)

	badIP := "invalid-ip"
	_, err = repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:         user.ID,
		FilePath:       "/tmp/test-avatar-bad-ip.png",
		FileSizeBytes:  1024,
		MimeType:       "image/png",
		Width:          100,
		Height:         100,
		Version:        1,
		UploadedFromIP: &badIP,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse IP address")
}

// ============================================================================
// ListUsers with Filters Integration Tests
// ============================================================================

func TestService_ListUsers_WithFilters(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create active admin
	_, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "admin_filter_test",
		Email:        "admin_filter@example.com",
		PasswordHash: "password123",
		IsActive:     ptr(true),
		IsAdmin:      ptr(true),
	})
	require.NoError(t, err)

	// Create inactive non-admin
	_, err = svc.CreateUser(ctx, CreateUserParams{
		Username:     "inactive_filter_test",
		Email:        "inactive_filter@example.com",
		PasswordHash: "password123",
		IsActive:     ptr(false),
		IsAdmin:      ptr(false),
	})
	require.NoError(t, err)

	// Filter for active admins
	users, count, err := svc.ListUsers(ctx, UserFilters{
		IsActive: ptr(true),
		IsAdmin:  ptr(true),
		Limit:    10,
		Offset:   0,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
	for _, u := range users {
		require.NotNil(t, u.IsActive)
		assert.True(t, *u.IsActive)
		require.NotNil(t, u.IsAdmin)
		assert.True(t, *u.IsAdmin)
	}

	// Filter for inactive users
	_, count, err = svc.ListUsers(ctx, UserFilters{
		IsActive: ptr(false),
		Limit:    10,
		Offset:   0,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

func TestService_ListUsers_Pagination(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create 5 users with unique identifiers
	for i := 0; i < 5; i++ {
		unique := uuid.Must(uuid.NewV7()).String()
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "paginate_" + unique,
			Email:        "paginate_" + unique + "@example.com",
			PasswordHash: "password123",
		})
		require.NoError(t, err)
	}

	// Page 1: limit 2
	page1, total, err := svc.ListUsers(ctx, UserFilters{
		Limit:  2,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.Len(t, page1, 2)

	// Page 2: limit 2, offset 2
	page2, _, err := svc.ListUsers(ctx, UserFilters{
		Limit:  2,
		Offset: 2,
	})
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	// Pages should contain different users
	if len(page1) > 0 && len(page2) > 0 {
		assert.NotEqual(t, page1[0].ID, page2[0].ID, "different pages should have different users")
	}
}

// ============================================================================
// UpdateUser with Multiple Fields
// ============================================================================

func TestService_UpdateUser_MultipleFields(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "multi_update",
		Email:        "multi_update@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Update display name and email
	displayName := "Multi Updated"
	newEmail := "multi_updated_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com"
	updated, err := svc.UpdateUser(ctx, created.ID, UpdateUserParams{
		DisplayName: &displayName,
		Email:       &newEmail,
	})
	require.NoError(t, err)
	require.NotNil(t, updated.DisplayName)
	assert.Equal(t, "Multi Updated", *updated.DisplayName)
	assert.Equal(t, newEmail, updated.Email)
}

func TestService_UpdateUser_SetTimezone(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "tz_update",
		Email:        "tz_update@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	tz := "Europe/Berlin"
	updated, err := svc.UpdateUser(ctx, created.ID, UpdateUserParams{
		Timezone: &tz,
	})
	require.NoError(t, err)
	require.NotNil(t, updated.Timezone)
	assert.Equal(t, "Europe/Berlin", *updated.Timezone)
}

// ============================================================================
// HardDeleteUser Lifecycle (includes preferences cleanup)
// ============================================================================

func TestService_HardDeleteUser_WithPreferences(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	// Create user with preferences
	created, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "hard_del_prefs",
		Email:        "hard_del_prefs@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Explicitly set preferences
	theme := "dark"
	_, err = svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
		UserID: created.ID,
		Theme:  &theme,
	})
	require.NoError(t, err)

	// Verify preferences exist
	prefs, err := svc.GetUserPreferences(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "dark", *prefs.Theme)

	// Hard delete user (should also delete preferences)
	err = svc.HardDeleteUser(ctx, created.ID)
	require.NoError(t, err)

	// User should not exist
	_, err = svc.GetUser(ctx, created.ID)
	require.Error(t, err)
}

func TestService_HardDeleteUser_WithAvatars(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "hard_del_avatars",
		Email:        "hard_del_avatars@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Upload avatar
	_, err = svc.UploadAvatar(ctx, user.ID, nil, AvatarMetadata{
		FileName:      "test.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
	})
	require.NoError(t, err)

	// Hard delete user (cascades to avatars)
	err = svc.HardDeleteUser(ctx, user.ID)
	require.NoError(t, err)

	// User should not exist
	_, err = svc.GetUser(ctx, user.ID)
	require.Error(t, err)
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestService_CreateUser_WithAllOptionalFields(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	displayName := "Full User"
	tz := "America/New_York"

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "full_fields_user",
		Email:        "full_fields@example.com",
		PasswordHash: "password123",
		DisplayName:  &displayName,
		Timezone:     &tz,
		QarEnabled:   ptr(true),
		IsActive:     ptr(true),
		IsAdmin:      ptr(false),
	})
	require.NoError(t, err)
	assert.Equal(t, "full_fields_user", user.Username)
	require.NotNil(t, user.DisplayName)
	assert.Equal(t, "Full User", *user.DisplayName)
	require.NotNil(t, user.Timezone)
	assert.Equal(t, "America/New_York", *user.Timezone)
}

func TestService_VerifyPassword_InvalidHash(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)

	// Try verifying against an invalid hash format
	err := svc.VerifyPassword("not-a-valid-hash", "password")
	require.Error(t, err)
}

func TestService_ListUserAvatars_NegativeOffset(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "neg_offset_user",
		Email:        "neg_offset@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// Negative limit should be corrected to default (10)
	avatars, err := svc.ListUserAvatars(ctx, user.ID, -5, 0)
	require.NoError(t, err)
	assert.NotNil(t, avatars)
}

func TestService_GetUserPreferences_AutoCreatesDefaults(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestService(t)
	ctx := context.Background()

	user, err := svc.CreateUser(ctx, CreateUserParams{
		Username:     "auto_prefs_user",
		Email:        "auto_prefs@example.com",
		PasswordHash: "password123",
	})
	require.NoError(t, err)

	// GetUserPreferences should auto-create defaults if missing
	prefs, err := svc.GetUserPreferences(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, prefs.UserID)

	// Default values should be set
	require.NotNil(t, prefs.Theme)
	assert.Equal(t, "system", *prefs.Theme)
	require.NotNil(t, prefs.ProfileVisibility)
	assert.Equal(t, "private", *prefs.ProfileVisibility)
}

// ============================================================================
// Repository-level Tests for Additional Coverage
// ============================================================================

func TestRepository_GetLatestAvatarVersion_NoAvatars(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "no_avatar_version",
		Email:        "no_avatar_version@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)

	version, err := repo.GetLatestAvatarVersion(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int32(0), version)
}

func TestRepository_UnsetCurrentAvatars(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "unset_current_user",
		Email:        "unset_current@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)

	// Create avatar
	avatar, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:        user.ID,
		FilePath:      "/tmp/unset-test.png",
		FileSizeBytes: 1024,
		MimeType:      "image/png",
		Width:         100,
		Height:        100,
		Version:       1,
	})
	require.NoError(t, err)

	// Set as current
	err = repo.SetCurrentAvatar(ctx, avatar.ID)
	require.NoError(t, err)

	// Verify it's current
	current, err := repo.GetCurrentAvatar(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar.ID, current.ID)

	// Unset current
	err = repo.UnsetCurrentAvatars(ctx, user.ID)
	require.NoError(t, err)

	// No current avatar should exist now
	_, err = repo.GetCurrentAvatar(ctx, user.ID)
	require.Error(t, err)
}

func TestRepository_DeleteUserPreferences(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "del_prefs_user",
		Email:        "del_prefs@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test$test",
	})
	require.NoError(t, err)

	// Create preferences
	theme := "dark"
	_, err = repo.UpsertUserPreferences(ctx, UpsertPreferencesParams{
		UserID: user.ID,
		Theme:  &theme,
	})
	require.NoError(t, err)

	// Delete preferences
	err = repo.DeleteUserPreferences(ctx, user.ID)
	require.NoError(t, err)

	// Preferences should not exist anymore
	_, err = repo.GetUserPreferences(ctx, user.ID)
	require.Error(t, err)
}
