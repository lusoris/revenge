package user

import (
	"context"
	"fmt"
	"os"
	"strings"
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

func setupTestRepo(t *testing.T) (Repository, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewPostgresRepository(queries)
	return repo, testDB
}

// contains checks if s contains any of the substrings
func contains(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ============================================================================
// User CRUD Tests
// ============================================================================

func TestPostgresRepository_CreateUser(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		params  CreateUserParams
		wantErr bool
	}{
		{
			name: "valid user",
			params: CreateUserParams{
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hashedpassword123",
				DisplayName:  new("Test User"),
				Timezone:     new("Europe/Berlin"),
				QarEnabled:   new(false),
				IsActive:     new(true),
				IsAdmin:      new(false),
			},
			wantErr: false,
		},
		{
			name: "minimal user",
			params: CreateUserParams{
				Username:     "minimaluser",
				Email:        "minimal@example.com",
				PasswordHash: "hash",
			},
			wantErr: false,
		},
		{
			name: "admin user",
			params: CreateUserParams{
				Username:     "adminuser",
				Email:        "admin@example.com",
				PasswordHash: "adminhash",
				IsAdmin:      new(true),
				IsActive:     new(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.CreateUser(ctx, tt.params)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.NotEqual(t, uuid.Nil, user.ID)
			assert.Equal(t, tt.params.Username, user.Username)
			assert.Equal(t, tt.params.Email, user.Email)
			assert.Equal(t, tt.params.PasswordHash, user.PasswordHash)
			assert.NotZero(t, user.CreatedAt)
			assert.NotZero(t, user.UpdatedAt)

			if tt.params.DisplayName != nil {
				assert.Equal(t, tt.params.DisplayName, user.DisplayName)
			}
			if tt.params.IsAdmin != nil && *tt.params.IsAdmin {
				assert.NotNil(t, user.IsAdmin)
				assert.True(t, *user.IsAdmin)
			}
		})
	}
}

func TestPostgresRepository_GetUserByID(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user first
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "getbyid_user",
		Email:        "getbyid@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	t.Run("existing user", func(t *testing.T) {
		user, err := repo.GetUserByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
		assert.Equal(t, created.Username, user.Username)
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := repo.GetUserByID(ctx, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
		assert.True(t, contains(err.Error(), "not found", "no rows"))
	})
}

func TestPostgresRepository_GetUserByUsername(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "findme_username",
		Email:        "findme@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	t.Run("existing username", func(t *testing.T) {
		user, err := repo.GetUserByUsername(ctx, "findme_username")
		require.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
	})

	t.Run("non-existent username", func(t *testing.T) {
		_, err := repo.GetUserByUsername(ctx, "nonexistent")
		require.Error(t, err)
		assert.True(t, contains(err.Error(), "not found", "no rows"))
	})
}

func TestPostgresRepository_GetUserByEmail(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "email_user",
		Email:        "unique_email@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	t.Run("existing email", func(t *testing.T) {
		user, err := repo.GetUserByEmail(ctx, "unique_email@example.com")
		require.NoError(t, err)
		assert.Equal(t, created.ID, user.ID)
	})

	t.Run("non-existent email", func(t *testing.T) {
		_, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")
		require.Error(t, err)
		assert.True(t, contains(err.Error(), "not found", "no rows"))
	})
}

func TestPostgresRepository_UpdateUser(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "update_user",
		Email:        "update@example.com",
		PasswordHash: "hash",
		IsActive:     new(true),
	})
	require.NoError(t, err)

	t.Run("update display name", func(t *testing.T) {
		updated, err := repo.UpdateUser(ctx, created.ID, UpdateUserParams{
			DisplayName: new("New Display Name"),
		})
		require.NoError(t, err)
		require.NotNil(t, updated.DisplayName)
		assert.Equal(t, "New Display Name", *updated.DisplayName)
	})

	t.Run("update multiple fields", func(t *testing.T) {
		updated, err := repo.UpdateUser(ctx, created.ID, UpdateUserParams{
			Email:      new("newemail@example.com"),
			Timezone:   new("America/New_York"),
			QarEnabled: new(true),
		})
		require.NoError(t, err)
		assert.Equal(t, "newemail@example.com", updated.Email)
		require.NotNil(t, updated.Timezone)
		assert.Equal(t, "America/New_York", *updated.Timezone)
		require.NotNil(t, updated.QarEnabled)
		assert.True(t, *updated.QarEnabled)
	})

	t.Run("update non-existent user", func(t *testing.T) {
		_, err := repo.UpdateUser(ctx, uuid.Must(uuid.NewV7()), UpdateUserParams{
			DisplayName: new("Test"),
		})
		require.Error(t, err)
		assert.True(t, contains(err.Error(), "not found", "no rows"))
	})
}

func TestPostgresRepository_UpdatePassword(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "password_user",
		Email:        "password@example.com",
		PasswordHash: "oldhash",
	})
	require.NoError(t, err)

	err = repo.UpdatePassword(ctx, created.ID, "newhash")
	require.NoError(t, err)

	// Verify password was updated
	user, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "newhash", user.PasswordHash)
}

func TestPostgresRepository_UpdateLastLogin(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "lastlogin_user",
		Email:        "lastlogin@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	assert.False(t, created.LastLoginAt.Valid)

	err = repo.UpdateLastLogin(ctx, created.ID)
	require.NoError(t, err)

	// Verify last_login_at was set
	user, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, user.LastLoginAt.Valid)
}

func TestPostgresRepository_VerifyEmail(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "verify_user",
		Email:        "verify@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	// Default should be false or nil
	if created.EmailVerified != nil {
		assert.False(t, *created.EmailVerified)
	}

	err = repo.VerifyEmail(ctx, created.ID)
	require.NoError(t, err)

	// Verify email_verified was set
	user, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, user.EmailVerified)
	assert.True(t, *user.EmailVerified)
	assert.True(t, user.EmailVerifiedAt.Valid)
}

func TestPostgresRepository_DeleteUser(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "delete_user",
		Email:        "delete@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Soft delete
	err = repo.DeleteUser(ctx, created.ID)
	require.NoError(t, err)

	// User should not be found via normal GetUserByID (which filters deleted_at IS NULL)
	_, err = repo.GetUserByID(ctx, created.ID)
	require.Error(t, err)
	assert.True(t, contains(err.Error(), "not found", "no rows"))
}

func TestPostgresRepository_HardDeleteUser(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	created, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "harddelete_user",
		Email:        "harddelete@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Hard delete
	err = repo.HardDeleteUser(ctx, created.ID)
	require.NoError(t, err)

	// User should not exist
	_, err = repo.GetUserByID(ctx, created.ID)
	require.Error(t, err)
	assert.True(t, contains(err.Error(), "not found", "no rows"))
}

func TestPostgresRepository_ListUsers(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create test users with different states
	_, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "inactive_regular",
		Email:        "inactive@example.com",
		PasswordHash: "hash",
		IsAdmin:      new(false),
		IsActive:     new(false),
	})
	require.NoError(t, err)

	_, err = repo.CreateUser(ctx, CreateUserParams{
		Username:     "active_regular",
		Email:        "active@example.com",
		PasswordHash: "hash",
		IsAdmin:      new(false),
		IsActive:     new(true),
	})
	require.NoError(t, err)

	_, err = repo.CreateUser(ctx, CreateUserParams{
		Username:     "active_admin",
		Email:        "admin@example.com",
		PasswordHash: "hash",
		IsAdmin:      new(true),
		IsActive:     new(true),
	})
	require.NoError(t, err)

	t.Run("list all users with no filters", func(t *testing.T) {
		// No filters should return ALL users
		users, count, err := repo.ListUsers(ctx, UserFilters{
			Limit:  100,
			Offset: 0,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3, "Should return all users when no filters")
		assert.GreaterOrEqual(t, count, int64(3))
	})

	t.Run("pagination", func(t *testing.T) {
		users, _, err := repo.ListUsers(ctx, UserFilters{
			Limit:  2,
			Offset: 0,
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(users), 2)
	})

	t.Run("filter by active", func(t *testing.T) {
		users, count, err := repo.ListUsers(ctx, UserFilters{
			IsActive: new(true),
			Limit:    100,
			Offset:   0,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 2, "Should return active users (regular + admin)")
		assert.GreaterOrEqual(t, count, int64(2))
		// All returned users should be active
		for _, u := range users {
			if u.IsActive != nil {
				assert.True(t, *u.IsActive)
			}
		}
	})

	t.Run("filter by inactive", func(t *testing.T) {
		users, count, err := repo.ListUsers(ctx, UserFilters{
			IsActive: new(false),
			Limit:    100,
			Offset:   0,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1, "Should return inactive user")
		assert.GreaterOrEqual(t, count, int64(1))
	})

	t.Run("filter by admin", func(t *testing.T) {
		users, count, err := repo.ListUsers(ctx, UserFilters{
			IsAdmin: new(true),
			Limit:   100,
			Offset:  0,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1, "Should return admin user")
		assert.GreaterOrEqual(t, count, int64(1))
		// All returned users should be admins
		for _, u := range users {
			if u.IsAdmin != nil {
				assert.True(t, *u.IsAdmin)
			}
		}
	})

	t.Run("filter by both active and admin", func(t *testing.T) {
		users, count, err := repo.ListUsers(ctx, UserFilters{
			IsActive: new(true),
			IsAdmin:  new(true),
			Limit:    100,
			Offset:   0,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1, "Should return active admin user")
		assert.GreaterOrEqual(t, count, int64(1))
	})
}

// ============================================================================
// User Preferences Tests
// ============================================================================

func TestPostgresRepository_UpsertUserPreferences(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "prefs_user",
		Email:        "prefs@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	t.Run("create preferences", func(t *testing.T) {
		prefs, err := repo.UpsertUserPreferences(ctx, UpsertPreferencesParams{
			UserID:            user.ID,
			Theme:             new("dark"),
			DisplayLanguage:   new("de"),
			ShowAdultContent:  new(false),
			AutoPlayVideos:    new(true),
			ProfileVisibility: new("private"),
		})
		require.NoError(t, err)
		assert.Equal(t, user.ID, prefs.UserID)
		require.NotNil(t, prefs.Theme)
		assert.Equal(t, "dark", *prefs.Theme)
		require.NotNil(t, prefs.DisplayLanguage)
		assert.Equal(t, "de", *prefs.DisplayLanguage)
	})

	t.Run("update preferences", func(t *testing.T) {
		prefs, err := repo.UpsertUserPreferences(ctx, UpsertPreferencesParams{
			UserID: user.ID,
			Theme:  new("light"),
		})
		require.NoError(t, err)
		require.NotNil(t, prefs.Theme)
		assert.Equal(t, "light", *prefs.Theme)
	})
}

func TestPostgresRepository_GetUserPreferences(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create user and preferences
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "getprefs_user",
		Email:        "getprefs@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	_, err = repo.UpsertUserPreferences(ctx, UpsertPreferencesParams{
		UserID: user.ID,
		Theme:  new("dark"),
	})
	require.NoError(t, err)

	t.Run("get existing preferences", func(t *testing.T) {
		prefs, err := repo.GetUserPreferences(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, prefs.UserID)
	})

	t.Run("get non-existent preferences", func(t *testing.T) {
		_, err := repo.GetUserPreferences(ctx, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
		assert.True(t, contains(err.Error(), "not found", "no rows"))
	})
}

func TestPostgresRepository_DeleteUserPreferences(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create user and preferences
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "delprefs_user",
		Email:        "delprefs@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	_, err = repo.UpsertUserPreferences(ctx, UpsertPreferencesParams{
		UserID: user.ID,
		Theme:  new("dark"),
	})
	require.NoError(t, err)

	err = repo.DeleteUserPreferences(ctx, user.ID)
	require.NoError(t, err)

	_, err = repo.GetUserPreferences(ctx, user.ID)
	require.Error(t, err)
}

// ============================================================================
// User Avatars Tests
// ============================================================================

func TestPostgresRepository_CreateAvatar(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "avatar_user",
		Email:        "avatar@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	t.Run("create avatar", func(t *testing.T) {
		avatar, err := repo.CreateAvatar(ctx, CreateAvatarParams{
			UserID:        user.ID,
			FilePath:      "/avatars/test.jpg",
			MimeType:      "image/jpeg",
			FileSizeBytes: 1024,
			Width:         100,
			Height:        100,
			Version:       1,
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, avatar.ID)
		assert.Equal(t, user.ID, avatar.UserID)
		assert.Equal(t, "/avatars/test.jpg", avatar.FilePath)
	})
}

func TestPostgresRepository_GetCurrentAvatar(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "currentavatar_user",
		Email:        "currentavatar@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Create avatar and set as current
	created, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:        user.ID,
		FilePath:      "/avatars/current.jpg",
		MimeType:      "image/jpeg",
		FileSizeBytes: 1024,
		Width:         100,
		Height:        100,
		Version:       1,
	})
	require.NoError(t, err)

	// Set as current
	err = repo.SetCurrentAvatar(ctx, created.ID)
	require.NoError(t, err)

	t.Run("get current avatar", func(t *testing.T) {
		avatar, err := repo.GetCurrentAvatar(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, avatar.ID)
		require.NotNil(t, avatar.IsCurrent)
		assert.True(t, *avatar.IsCurrent)
	})

	t.Run("no current avatar", func(t *testing.T) {
		_, err := repo.GetCurrentAvatar(ctx, uuid.Must(uuid.NewV7()))
		require.Error(t, err)
	})
}

func TestPostgresRepository_ListUserAvatars(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "listavatars_user",
		Email:        "listavatars@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Create multiple avatars (must unset current before creating new ones due to unique constraint)
	for i := 1; i <= 3; i++ {
		// Unset any existing current avatar before creating new one
		_ = repo.UnsetCurrentAvatars(ctx, user.ID)
		_, err := repo.CreateAvatar(ctx, CreateAvatarParams{
			UserID:        user.ID,
			FilePath:      fmt.Sprintf("/avatars/v%d.jpg", i),
			MimeType:      "image/jpeg",
			FileSizeBytes: int64(1024 * i),
			Width:         int32(100 * i),
			Height:        int32(100 * i),
			Version:       int32(i),
		})
		require.NoError(t, err)
	}

	avatars, err := repo.ListUserAvatars(ctx, user.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, avatars, 3)
}

func TestPostgresRepository_SetCurrentAvatar(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "setcurrent_user",
		Email:        "setcurrent@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Create two avatars
	avatar1, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:        user.ID,
		FilePath:      "/avatars/v1.jpg",
		MimeType:      "image/jpeg",
		FileSizeBytes: 1024,
		Version:       1,
	})
	require.NoError(t, err)

	// Unset current before creating second avatar (unique constraint)
	err = repo.UnsetCurrentAvatars(ctx, user.ID)
	require.NoError(t, err)

	avatar2, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:        user.ID,
		FilePath:      "/avatars/v2.jpg",
		MimeType:      "image/jpeg",
		FileSizeBytes: 2048,
		Version:       2,
	})
	require.NoError(t, err)

	// avatar2 is now current (from CreateAvatar), verify avatar1 is not
	a1, err := repo.GetAvatarByID(ctx, avatar1.ID)
	require.NoError(t, err)
	if a1.IsCurrent != nil {
		assert.False(t, *a1.IsCurrent)
	}

	// Test switching current avatar: unset all and set avatar1 as current
	err = repo.UnsetCurrentAvatars(ctx, user.ID)
	require.NoError(t, err)

	err = repo.SetCurrentAvatar(ctx, avatar1.ID)
	require.NoError(t, err)

	// Verify avatar1 is now current
	current, err := repo.GetCurrentAvatar(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, avatar1.ID, current.ID)

	// Verify avatar2 is no longer current
	a2, err := repo.GetAvatarByID(ctx, avatar2.ID)
	require.NoError(t, err)
	if a2.IsCurrent != nil {
		assert.False(t, *a2.IsCurrent)
	}
}

func TestPostgresRepository_DeleteAvatar(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user and avatar
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "deleteavatar_user",
		Email:        "deleteavatar@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	avatar, err := repo.CreateAvatar(ctx, CreateAvatarParams{
		UserID:        user.ID,
		FilePath:      "/avatars/delete.jpg",
		MimeType:      "image/jpeg",
		FileSizeBytes: 1024,
		Version:       1,
	})
	require.NoError(t, err)

	t.Run("soft delete", func(t *testing.T) {
		err := repo.DeleteAvatar(ctx, avatar.ID)
		require.NoError(t, err)

		// Avatar should not be found via normal GetAvatarByID (which filters deleted_at IS NULL)
		_, err = repo.GetAvatarByID(ctx, avatar.ID)
		require.Error(t, err)
		assert.True(t, contains(err.Error(), "not found", "no rows"))
	})
}

func TestPostgresRepository_GetLatestAvatarVersion(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	// Create a user
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username:     "latestversion_user",
		Email:        "latestversion@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	t.Run("no avatars", func(t *testing.T) {
		version, err := repo.GetLatestAvatarVersion(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int32(0), version)
	})

	// Create avatars - must unset current before each one due to unique constraint
	for i := 1; i <= 3; i++ {
		// Unset any current avatar before creating new one
		err := repo.UnsetCurrentAvatars(ctx, user.ID)
		require.NoError(t, err)

		_, err = repo.CreateAvatar(ctx, CreateAvatarParams{
			UserID:        user.ID,
			FilePath:      fmt.Sprintf("/avatars/v%d.jpg", i),
			MimeType:      "image/jpeg",
			FileSizeBytes: 1024,
			Version:       int32(i),
		})
		require.NoError(t, err)
	}

	t.Run("with avatars", func(t *testing.T) {
		version, err := repo.GetLatestAvatarVersion(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int32(3), version)
	})
}
