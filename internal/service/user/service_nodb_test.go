package user

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/storage"
)

// ============================================================================
// Inline Mock Repository (no mockery dependency)
// ============================================================================

// mockRepo is an inline mock implementation of the Repository interface.
// Each field is a function that can be set per-test to control behavior.
type mockRepo struct {
	getUserByIDFn           func(ctx context.Context, id uuid.UUID) (*db.SharedUser, error)
	getUserByUsernameFn     func(ctx context.Context, username string) (*db.SharedUser, error)
	getUserByEmailFn        func(ctx context.Context, email string) (*db.SharedUser, error)
	listUsersFn             func(ctx context.Context, filters UserFilters) ([]db.SharedUser, int64, error)
	createUserFn            func(ctx context.Context, params CreateUserParams) (*db.SharedUser, error)
	updateUserFn            func(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*db.SharedUser, error)
	updatePasswordFn        func(ctx context.Context, id uuid.UUID, passwordHash string) error
	updateLastLoginFn       func(ctx context.Context, id uuid.UUID) error
	verifyEmailFn           func(ctx context.Context, id uuid.UUID) error
	deleteUserFn            func(ctx context.Context, id uuid.UUID) error
	hardDeleteUserFn        func(ctx context.Context, id uuid.UUID) error
	getUserPreferencesFn    func(ctx context.Context, userID uuid.UUID) (*db.SharedUserPreference, error)
	upsertUserPreferencesFn func(ctx context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error)
	deleteUserPreferencesFn func(ctx context.Context, userID uuid.UUID) error
	getCurrentAvatarFn      func(ctx context.Context, userID uuid.UUID) (*db.SharedUserAvatar, error)
	getAvatarByIDFn         func(ctx context.Context, id uuid.UUID) (*db.SharedUserAvatar, error)
	listUserAvatarsFn       func(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.SharedUserAvatar, error)
	createAvatarFn          func(ctx context.Context, params CreateAvatarParams) (*db.SharedUserAvatar, error)
	unsetCurrentAvatarsFn   func(ctx context.Context, userID uuid.UUID) error
	setCurrentAvatarFn      func(ctx context.Context, id uuid.UUID) error
	deleteAvatarFn          func(ctx context.Context, id uuid.UUID) error
	hardDeleteAvatarFn      func(ctx context.Context, id uuid.UUID) error
	getLatestAvatarVerFn    func(ctx context.Context, userID uuid.UUID) (int32, error)
}

func (m *mockRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*db.SharedUser, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) GetUserByUsername(ctx context.Context, username string) (*db.SharedUser, error) {
	if m.getUserByUsernameFn != nil {
		return m.getUserByUsernameFn(ctx, username)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*db.SharedUser, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, email)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) ListUsers(ctx context.Context, filters UserFilters) ([]db.SharedUser, int64, error) {
	if m.listUsersFn != nil {
		return m.listUsersFn(ctx, filters)
	}
	return nil, 0, nil
}

func (m *mockRepo) CreateUser(ctx context.Context, params CreateUserParams) (*db.SharedUser, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, params)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) UpdateUser(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*db.SharedUser, error) {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, id, params)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, id, passwordHash)
	}
	return nil
}

func (m *mockRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	if m.updateLastLoginFn != nil {
		return m.updateLastLoginFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	if m.verifyEmailFn != nil {
		return m.verifyEmailFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.deleteUserFn != nil {
		return m.deleteUserFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) HardDeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.hardDeleteUserFn != nil {
		return m.hardDeleteUserFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*db.SharedUserPreference, error) {
	if m.getUserPreferencesFn != nil {
		return m.getUserPreferencesFn(ctx, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) UpsertUserPreferences(ctx context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
	if m.upsertUserPreferencesFn != nil {
		return m.upsertUserPreferencesFn(ctx, params)
	}
	return &db.SharedUserPreference{UserID: params.UserID}, nil
}

func (m *mockRepo) DeleteUserPreferences(ctx context.Context, userID uuid.UUID) error {
	if m.deleteUserPreferencesFn != nil {
		return m.deleteUserPreferencesFn(ctx, userID)
	}
	return nil
}

func (m *mockRepo) GetCurrentAvatar(ctx context.Context, userID uuid.UUID) (*db.SharedUserAvatar, error) {
	if m.getCurrentAvatarFn != nil {
		return m.getCurrentAvatarFn(ctx, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) GetAvatarByID(ctx context.Context, id uuid.UUID) (*db.SharedUserAvatar, error) {
	if m.getAvatarByIDFn != nil {
		return m.getAvatarByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) ListUserAvatars(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.SharedUserAvatar, error) {
	if m.listUserAvatarsFn != nil {
		return m.listUserAvatarsFn(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *mockRepo) CreateAvatar(ctx context.Context, params CreateAvatarParams) (*db.SharedUserAvatar, error) {
	if m.createAvatarFn != nil {
		return m.createAvatarFn(ctx, params)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) UnsetCurrentAvatars(ctx context.Context, userID uuid.UUID) error {
	if m.unsetCurrentAvatarsFn != nil {
		return m.unsetCurrentAvatarsFn(ctx, userID)
	}
	return nil
}

func (m *mockRepo) SetCurrentAvatar(ctx context.Context, id uuid.UUID) error {
	if m.setCurrentAvatarFn != nil {
		return m.setCurrentAvatarFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) DeleteAvatar(ctx context.Context, id uuid.UUID) error {
	if m.deleteAvatarFn != nil {
		return m.deleteAvatarFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) HardDeleteAvatar(ctx context.Context, id uuid.UUID) error {
	if m.hardDeleteAvatarFn != nil {
		return m.hardDeleteAvatarFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) GetLatestAvatarVersion(ctx context.Context, userID uuid.UUID) (int32, error) {
	if m.getLatestAvatarVerFn != nil {
		return m.getLatestAvatarVerFn(ctx, userID)
	}
	return 0, nil
}

// ============================================================================
// Test Helpers (no DB needed)
// ============================================================================

func newTestService(repo Repository) *Service {
	avatarCfg := config.AvatarConfig{
		StoragePath:  "/tmp/test-avatars",
		MaxSizeBytes: 5 * 1024 * 1024,
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
	}
	return NewService(nil, repo, activity.NewNoopLogger(), storage.NewMockStorage(), avatarCfg)
}

func makeUser(id uuid.UUID, username, email string) *db.SharedUser {
	now := time.Now()
	displayName := "Test User"
	isActive := true
	isAdmin := false
	return &db.SharedUser{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=4$c2FsdA$aGFzaA",
		DisplayName:  &displayName,
		IsActive:     &isActive,
		IsAdmin:      &isAdmin,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// ============================================================================
// Service Constructor Test
// ============================================================================

func TestNoDB_NewService(t *testing.T) {
	t.Parallel()

	repo := &mockRepo{}
	svc := newTestService(repo)

	require.NotNil(t, svc)
	assert.NotNil(t, svc.repo)
	assert.NotNil(t, svc.hasher)
	assert.NotNil(t, svc.activityLogger)
	assert.NotNil(t, svc.storage)
	assert.Nil(t, svc.pool, "pool should be nil in no-DB tests")
}

// ============================================================================
// User Management Tests
// ============================================================================

func TestNoDB_GetUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() *mockRepo
		wantErr bool
	}{
		{
			name: "success",
			setup: func() *mockRepo {
				return &mockRepo{
					getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
						return makeUser(id, "testuser", "test@example.com"), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "not found",
			setup: func() *mockRepo {
				return &mockRepo{
					getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
						return nil, errors.New("user not found")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := newTestService(tt.setup())
			userID := uuid.Must(uuid.NewV7())

			result, err := svc.GetUser(ctx, userID)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, userID, result.ID)
			}
		})
	}
}

func TestNoDB_GetUserByUsername(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByUsernameFn: func(_ context.Context, username string) (*db.SharedUser, error) {
				return makeUser(uuid.Must(uuid.NewV7()), username, "test@example.com"), nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.GetUserByUsername(ctx, "testuser")
		require.NoError(t, err)
		assert.Equal(t, "testuser", result.Username)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{}
		svc := newTestService(repo)
		_, err := svc.GetUserByUsername(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestNoDB_GetUserByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByEmailFn: func(_ context.Context, email string) (*db.SharedUser, error) {
				return makeUser(uuid.Must(uuid.NewV7()), "testuser", email), nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.GetUserByEmail(ctx, "test@example.com")
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{}
		svc := newTestService(repo)
		_, err := svc.GetUserByEmail(ctx, "nonexistent@example.com")
		require.Error(t, err)
	})
}

func TestNoDB_ListUsers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name      string
		filters   UserFilters
		wantCount int64
		wantLen   int
	}{
		{
			name:      "returns users",
			filters:   UserFilters{Limit: 10, Offset: 0},
			wantCount: 3,
			wantLen:   3,
		},
		{
			name:      "empty result",
			filters:   UserFilters{Limit: 10, Offset: 100},
			wantCount: 0,
			wantLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockRepo{
				listUsersFn: func(_ context.Context, _ UserFilters) ([]db.SharedUser, int64, error) {
					if tt.wantLen == 0 {
						return nil, 0, nil
					}
					users := make([]db.SharedUser, tt.wantLen)
					for i := range users {
						users[i] = *makeUser(uuid.Must(uuid.NewV7()), "user", "user@example.com")
					}
					return users, tt.wantCount, nil
				},
			}
			svc := newTestService(repo)
			result, count, err := svc.ListUsers(ctx, tt.filters)
			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, count)
			assert.Len(t, result, tt.wantLen)
		})
	}
}

func TestNoDB_CreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		created := makeUser(uuid.Must(uuid.NewV7()), "newuser", "new@example.com")
		repo := &mockRepo{
			// Username not found (good)
			getUserByUsernameFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			// Email not found (good)
			getUserByEmailFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			createUserFn: func(_ context.Context, params CreateUserParams) (*db.SharedUser, error) {
				created.PasswordHash = params.PasswordHash
				return created, nil
			},
			upsertUserPreferencesFn: func(_ context.Context, _ UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				return &db.SharedUserPreference{}, nil
			},
		}
		svc := newTestService(repo)

		result, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "newuser",
			Email:        "new@example.com",
			PasswordHash: "password123",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "newuser", result.Username)
		// Password should have been hashed (argon2id)
		assert.True(t, strings.HasPrefix(result.PasswordHash, "$argon2id$"))
	})

	t.Run("missing username", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Email:        "new@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("missing email", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "newuser",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("missing password", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username: "newuser",
			Email:    "new@example.com",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("username already exists", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByUsernameFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return makeUser(uuid.Must(uuid.NewV7()), "existing", "old@example.com"), nil
			},
		}
		svc := newTestService(repo)
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "existing",
			Email:        "new@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})

	t.Run("email already exists", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByUsernameFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			getUserByEmailFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return makeUser(uuid.Must(uuid.NewV7()), "other", "existing@example.com"), nil
			},
		}
		svc := newTestService(repo)
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "newuser",
			Email:        "existing@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})

	t.Run("repo create error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByUsernameFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			getUserByEmailFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			createUserFn: func(_ context.Context, _ CreateUserParams) (*db.SharedUser, error) {
				return nil, errors.New("db error")
			},
		}
		svc := newTestService(repo)
		_, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "newuser",
			Email:        "new@example.com",
			PasswordHash: "password123",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user")
	})

	t.Run("preferences upsert error does not fail creation", func(t *testing.T) {
		t.Parallel()
		created := makeUser(uuid.Must(uuid.NewV7()), "newuser", "new@example.com")
		repo := &mockRepo{
			getUserByUsernameFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			getUserByEmailFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			createUserFn: func(_ context.Context, params CreateUserParams) (*db.SharedUser, error) {
				created.PasswordHash = params.PasswordHash
				return created, nil
			},
			upsertUserPreferencesFn: func(_ context.Context, _ UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				return nil, errors.New("preferences error")
			},
		}
		svc := newTestService(repo)
		result, err := svc.CreateUser(ctx, CreateUserParams{
			Username:     "newuser",
			Email:        "new@example.com",
			PasswordHash: "password123",
		})
		// Should still succeed even if preferences fail
		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func TestNoDB_UpdateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success with display name change", func(t *testing.T) {
		t.Parallel()
		oldUser := makeUser(userID, "testuser", "test@example.com")
		updatedUser := makeUser(userID, "testuser", "test@example.com")
		newName := "New Name"
		updatedUser.DisplayName = &newName

		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return oldUser, nil
			},
			updateUserFn: func(_ context.Context, _ uuid.UUID, _ UpdateUserParams) (*db.SharedUser, error) {
				return updatedUser, nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.UpdateUser(ctx, userID, UpdateUserParams{DisplayName: &newName})
		require.NoError(t, err)
		assert.Equal(t, "New Name", *result.DisplayName)
	})

	t.Run("success with email change", func(t *testing.T) {
		t.Parallel()
		oldUser := makeUser(userID, "testuser", "old@example.com")
		newEmail := "new@example.com"
		updatedUser := makeUser(userID, "testuser", newEmail)

		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return oldUser, nil
			},
			updateUserFn: func(_ context.Context, _ uuid.UUID, _ UpdateUserParams) (*db.SharedUser, error) {
				return updatedUser, nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.UpdateUser(ctx, userID, UpdateUserParams{Email: &newEmail})
		require.NoError(t, err)
		assert.Equal(t, newEmail, result.Email)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
		}
		svc := newTestService(repo)
		_, err := svc.UpdateUser(ctx, userID, UpdateUserParams{DisplayName: new("name")})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repo update error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return makeUser(userID, "testuser", "test@example.com"), nil
			},
			updateUserFn: func(_ context.Context, _ uuid.UUID, _ UpdateUserParams) (*db.SharedUser, error) {
				return nil, errors.New("db error")
			},
		}
		svc := newTestService(repo)
		_, err := svc.UpdateUser(ctx, userID, UpdateUserParams{DisplayName: new("name")})
		require.Error(t, err)
	})
}

func TestNoDB_UpdatePassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		// Hash a password to set as the "old" password
		oldHash, err := svc.HashPassword("oldpassword")
		require.NoError(t, err)

		user := makeUser(userID, "testuser", "test@example.com")
		user.PasswordHash = oldHash

		var capturedHash string
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return user, nil
			},
			updatePasswordFn: func(_ context.Context, _ uuid.UUID, hash string) error {
				capturedHash = hash
				return nil
			},
		}
		svc = newTestService(repo)
		err = svc.UpdatePassword(ctx, userID, "oldpassword", "newpassword")
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(capturedHash, "$argon2id$"), "new password should be hashed")
	})

	t.Run("wrong old password", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		oldHash, err := svc.HashPassword("oldpassword")
		require.NoError(t, err)

		user := makeUser(userID, "testuser", "test@example.com")
		user.PasswordHash = oldHash

		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return user, nil
			},
		}
		svc = newTestService(repo)
		err = svc.UpdatePassword(ctx, userID, "wrongpassword", "newpassword")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid old password")
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
		}
		svc := newTestService(repo)
		err := svc.UpdatePassword(ctx, userID, "old", "new")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repo update error", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		oldHash, err := svc.HashPassword("oldpassword")
		require.NoError(t, err)

		user := makeUser(userID, "testuser", "test@example.com")
		user.PasswordHash = oldHash

		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return user, nil
			},
			updatePasswordFn: func(_ context.Context, _ uuid.UUID, _ string) error {
				return errors.New("db error")
			},
		}
		svc = newTestService(repo)
		err = svc.UpdatePassword(ctx, userID, "oldpassword", "newpassword")
		require.Error(t, err)
	})
}

func TestNoDB_DeleteUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return makeUser(userID, "testuser", "test@example.com"), nil
			},
			deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		}
		svc := newTestService(repo)
		err := svc.DeleteUser(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("success even when user not found for logging", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return nil, errors.New("not found")
			},
			deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		}
		svc := newTestService(repo)
		err := svc.DeleteUser(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("delete error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
				return makeUser(userID, "testuser", "test@example.com"), nil
			},
			deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		}
		svc := newTestService(repo)
		err := svc.DeleteUser(ctx, userID)
		require.Error(t, err)
	})
}

func TestNoDB_HardDeleteUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		deletedPrefs := false
		repo := &mockRepo{
			deleteUserPreferencesFn: func(_ context.Context, _ uuid.UUID) error {
				deletedPrefs = true
				return nil
			},
			hardDeleteUserFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		}
		svc := newTestService(repo)
		err := svc.HardDeleteUser(ctx, userID)
		require.NoError(t, err)
		assert.True(t, deletedPrefs, "preferences should be deleted")
	})

	t.Run("preferences error is ignored", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			deleteUserPreferencesFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("prefs error")
			},
			hardDeleteUserFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		}
		svc := newTestService(repo)
		err := svc.HardDeleteUser(ctx, userID)
		require.NoError(t, err, "preferences error should not prevent hard delete")
	})

	t.Run("hard delete error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			hardDeleteUserFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		}
		svc := newTestService(repo)
		err := svc.HardDeleteUser(ctx, userID)
		require.Error(t, err)
	})
}

func TestNoDB_VerifyEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{
			verifyEmailFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		})
		err := svc.VerifyEmail(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{
			verifyEmailFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		})
		err := svc.VerifyEmail(ctx, userID)
		require.Error(t, err)
	})
}

func TestNoDB_RecordLogin(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{
			updateLastLoginFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		})
		err := svc.RecordLogin(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{
			updateLastLoginFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		})
		err := svc.RecordLogin(ctx, userID)
		require.Error(t, err)
	})
}

// ============================================================================
// Password Management Tests
// ============================================================================

func TestNoDB_HashPassword(t *testing.T) {
	t.Parallel()
	svc := newTestService(&mockRepo{})

	t.Run("hashes password with argon2id", func(t *testing.T) {
		t.Parallel()
		hash, err := svc.HashPassword("testpassword")
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
		assert.NotEqual(t, "testpassword", hash)
	})

	t.Run("different inputs produce different hashes", func(t *testing.T) {
		t.Parallel()
		hash1, err := svc.HashPassword("password1")
		require.NoError(t, err)
		hash2, err := svc.HashPassword("password2")
		require.NoError(t, err)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("same input produces different hashes due to salt", func(t *testing.T) {
		t.Parallel()
		hash1, err := svc.HashPassword("samepassword")
		require.NoError(t, err)
		hash2, err := svc.HashPassword("samepassword")
		require.NoError(t, err)
		assert.NotEqual(t, hash1, hash2, "salted hashing should produce unique hashes")
	})
}

func TestNoDB_VerifyPassword(t *testing.T) {
	t.Parallel()
	svc := newTestService(&mockRepo{})

	hash, err := svc.HashPassword("testpassword")
	require.NoError(t, err)

	t.Run("correct password", func(t *testing.T) {
		t.Parallel()
		err := svc.VerifyPassword(hash, "testpassword")
		require.NoError(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		err := svc.VerifyPassword(hash, "wrongpassword")
		require.Error(t, err)
	})

	t.Run("empty password", func(t *testing.T) {
		t.Parallel()
		err := svc.VerifyPassword(hash, "")
		require.Error(t, err)
	})

	t.Run("invalid hash format", func(t *testing.T) {
		t.Parallel()
		err := svc.VerifyPassword("not-a-valid-hash", "testpassword")
		require.Error(t, err)
	})
}

// ============================================================================
// User Preferences Tests
// ============================================================================

func TestNoDB_GetUserPreferences(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("returns existing preferences", func(t *testing.T) {
		t.Parallel()
		theme := "dark"
		repo := &mockRepo{
			getUserPreferencesFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserPreference, error) {
				return &db.SharedUserPreference{UserID: userID, Theme: &theme}, nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.GetUserPreferences(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "dark", *result.Theme)
	})

	t.Run("creates defaults when not found", func(t *testing.T) {
		t.Parallel()
		var capturedParams UpsertPreferencesParams
		repo := &mockRepo{
			getUserPreferencesFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserPreference, error) {
				return nil, errors.New("not found")
			},
			upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				capturedParams = params
				return &db.SharedUserPreference{UserID: params.UserID}, nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.GetUserPreferences(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, userID, result.UserID)

		// Verify default preferences were passed
		assert.Equal(t, userID, capturedParams.UserID)
		require.NotNil(t, capturedParams.Theme)
		assert.Equal(t, "system", *capturedParams.Theme)
		require.NotNil(t, capturedParams.ProfileVisibility)
		assert.Equal(t, "private", *capturedParams.ProfileVisibility)
		require.NotNil(t, capturedParams.DisplayLanguage)
		assert.Equal(t, "en-US", *capturedParams.DisplayLanguage)
	})
}

func TestNoDB_UpdateUserPreferences(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("valid theme", func(t *testing.T) {
		t.Parallel()
		themes := []string{"light", "dark", "system"}
		for _, theme := range themes {
			t.Run(theme, func(t *testing.T) {
				t.Parallel()
				repo := &mockRepo{
					upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
						return &db.SharedUserPreference{UserID: params.UserID, Theme: params.Theme}, nil
					},
				}
				svc := newTestService(repo)
				result, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
					UserID: userID,
					Theme:  new(theme),
				})
				require.NoError(t, err)
				assert.Equal(t, theme, *result.Theme)
			})
		}
	})

	t.Run("invalid theme", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		_, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
			UserID: userID,
			Theme:  new("invalid"),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid theme")
	})

	t.Run("valid profile visibility", func(t *testing.T) {
		t.Parallel()
		visibilities := []string{"public", "friends", "private"}
		for _, vis := range visibilities {
			t.Run(vis, func(t *testing.T) {
				t.Parallel()
				repo := &mockRepo{
					upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
						return &db.SharedUserPreference{UserID: params.UserID, ProfileVisibility: params.ProfileVisibility}, nil
					},
				}
				svc := newTestService(repo)
				result, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
					UserID:            userID,
					ProfileVisibility: new(vis),
				})
				require.NoError(t, err)
				assert.Equal(t, vis, *result.ProfileVisibility)
			})
		}
	})

	t.Run("invalid profile visibility", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		_, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
			UserID:            userID,
			ProfileVisibility: new("invalid"),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid profile visibility")
	})

	t.Run("no validation fields passes", func(t *testing.T) {
		t.Parallel()
		showAdult := true
		repo := &mockRepo{
			upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				return &db.SharedUserPreference{UserID: params.UserID, ShowAdultContent: params.ShowAdultContent}, nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.UpdateUserPreferences(ctx, UpsertPreferencesParams{
			UserID:           userID,
			ShowAdultContent: &showAdult,
		})
		require.NoError(t, err)
		require.NotNil(t, result.ShowAdultContent)
		assert.True(t, *result.ShowAdultContent)
	})
}

func TestNoDB_UpdateNotificationPreferences(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("all notification types", func(t *testing.T) {
		t.Parallel()
		var capturedParams UpsertPreferencesParams
		repo := &mockRepo{
			upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				capturedParams = params
				return &db.SharedUserPreference{UserID: userID}, nil
			},
		}
		svc := newTestService(repo)

		emailSettings := &NotificationSettings{Enabled: true, Frequency: "daily"}
		pushSettings := &NotificationSettings{Enabled: false}
		digestSettings := &NotificationSettings{Enabled: true, Frequency: "weekly"}

		err := svc.UpdateNotificationPreferences(ctx, userID, emailSettings, pushSettings, digestSettings)
		require.NoError(t, err)

		// Verify all notifications were marshalled
		require.NotNil(t, capturedParams.EmailNotifications)
		require.NotNil(t, capturedParams.PushNotifications)
		require.NotNil(t, capturedParams.DigestNotifications)

		var email NotificationSettings
		err = json.Unmarshal(*capturedParams.EmailNotifications, &email)
		require.NoError(t, err)
		assert.True(t, email.Enabled)
		assert.Equal(t, "daily", email.Frequency)

		var push NotificationSettings
		err = json.Unmarshal(*capturedParams.PushNotifications, &push)
		require.NoError(t, err)
		assert.False(t, push.Enabled)

		var digest NotificationSettings
		err = json.Unmarshal(*capturedParams.DigestNotifications, &digest)
		require.NoError(t, err)
		assert.True(t, digest.Enabled)
		assert.Equal(t, "weekly", digest.Frequency)
	})

	t.Run("nil notification types are skipped", func(t *testing.T) {
		t.Parallel()
		var capturedParams UpsertPreferencesParams
		repo := &mockRepo{
			upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				capturedParams = params
				return &db.SharedUserPreference{UserID: userID}, nil
			},
		}
		svc := newTestService(repo)

		err := svc.UpdateNotificationPreferences(ctx, userID, nil, nil, nil)
		require.NoError(t, err)

		assert.Nil(t, capturedParams.EmailNotifications)
		assert.Nil(t, capturedParams.PushNotifications)
		assert.Nil(t, capturedParams.DigestNotifications)
	})

	t.Run("partial notification types", func(t *testing.T) {
		t.Parallel()
		var capturedParams UpsertPreferencesParams
		repo := &mockRepo{
			upsertUserPreferencesFn: func(_ context.Context, params UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				capturedParams = params
				return &db.SharedUserPreference{UserID: userID}, nil
			},
		}
		svc := newTestService(repo)

		emailSettings := &NotificationSettings{Enabled: true, Frequency: "instant"}
		err := svc.UpdateNotificationPreferences(ctx, userID, emailSettings, nil, nil)
		require.NoError(t, err)

		require.NotNil(t, capturedParams.EmailNotifications)
		assert.Nil(t, capturedParams.PushNotifications)
		assert.Nil(t, capturedParams.DigestNotifications)
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			upsertUserPreferencesFn: func(_ context.Context, _ UpsertPreferencesParams) (*db.SharedUserPreference, error) {
				return nil, errors.New("db error")
			},
		}
		svc := newTestService(repo)
		err := svc.UpdateNotificationPreferences(ctx, userID,
			&NotificationSettings{Enabled: true}, nil, nil)
		require.Error(t, err)
	})
}

// ============================================================================
// Avatar Management Tests
// ============================================================================

func TestNoDB_GetCurrentAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		avatarID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			getCurrentAvatarFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path/to/avatar.png"}, nil
			},
		}
		svc := newTestService(repo)
		result, err := svc.GetCurrentAvatar(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, avatarID, result.ID)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		_, err := svc.GetCurrentAvatar(ctx, userID)
		require.Error(t, err)
	})
}

func TestNoDB_ListUserAvatars(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	t.Run("default limit when 0", func(t *testing.T) {
		t.Parallel()
		var capturedLimit int32
		repo := &mockRepo{
			listUserAvatarsFn: func(_ context.Context, _ uuid.UUID, limit, _ int32) ([]db.SharedUserAvatar, error) {
				capturedLimit = limit
				return []db.SharedUserAvatar{}, nil
			},
		}
		svc := newTestService(repo)
		_, err := svc.ListUserAvatars(ctx, userID, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, int32(10), capturedLimit)
	})

	t.Run("negative limit defaults to 10", func(t *testing.T) {
		t.Parallel()
		var capturedLimit int32
		repo := &mockRepo{
			listUserAvatarsFn: func(_ context.Context, _ uuid.UUID, limit, _ int32) ([]db.SharedUserAvatar, error) {
				capturedLimit = limit
				return []db.SharedUserAvatar{}, nil
			},
		}
		svc := newTestService(repo)
		_, err := svc.ListUserAvatars(ctx, userID, -5, 0)
		require.NoError(t, err)
		assert.Equal(t, int32(10), capturedLimit)
	})

	t.Run("max limit capped at 100", func(t *testing.T) {
		t.Parallel()
		var capturedLimit int32
		repo := &mockRepo{
			listUserAvatarsFn: func(_ context.Context, _ uuid.UUID, limit, _ int32) ([]db.SharedUserAvatar, error) {
				capturedLimit = limit
				return []db.SharedUserAvatar{}, nil
			},
		}
		svc := newTestService(repo)
		_, err := svc.ListUserAvatars(ctx, userID, 500, 0)
		require.NoError(t, err)
		assert.Equal(t, int32(100), capturedLimit)
	})

	t.Run("valid limit passes through", func(t *testing.T) {
		t.Parallel()
		var capturedLimit int32
		repo := &mockRepo{
			listUserAvatarsFn: func(_ context.Context, _ uuid.UUID, limit, _ int32) ([]db.SharedUserAvatar, error) {
				capturedLimit = limit
				return []db.SharedUserAvatar{}, nil
			},
		}
		svc := newTestService(repo)
		_, err := svc.ListUserAvatars(ctx, userID, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, int32(50), capturedLimit)
	})
}

func TestNoDB_SetCurrentAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	avatarID := uuid.Must(uuid.NewV7())
	otherUserID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path.png"}, nil
			},
			unsetCurrentAvatarsFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
			setCurrentAvatarFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
			updateUserFn: func(_ context.Context, _ uuid.UUID, _ UpdateUserParams) (*db.SharedUser, error) {
				return makeUser(userID, "testuser", "test@example.com"), nil
			},
		}
		svc := newTestService(repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)
		require.NoError(t, err)
	})

	t.Run("avatar not found", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar not found")
	})

	t.Run("avatar belongs to other user", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: otherUserID}, nil
			},
		}
		svc := newTestService(repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar does not belong to user")
	})

	t.Run("unset error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path.png"}, nil
			},
			unsetCurrentAvatarsFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		}
		svc := newTestService(repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unset current avatars")
	})

	t.Run("set current error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: userID, FilePath: "/path.png"}, nil
			},
			unsetCurrentAvatarsFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
			setCurrentAvatarFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		}
		svc := newTestService(repo)
		err := svc.SetCurrentAvatar(ctx, userID, avatarID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to set current avatar")
	})
}

func TestNoDB_DeleteAvatar(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())
	avatarID := uuid.Must(uuid.NewV7())
	otherUserID := uuid.Must(uuid.NewV7())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: userID}, nil
			},
			deleteAvatarFn: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		}
		svc := newTestService(repo)
		err := svc.DeleteAvatar(ctx, userID, avatarID)
		require.NoError(t, err)
	})

	t.Run("avatar not found", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})
		err := svc.DeleteAvatar(ctx, userID, avatarID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar not found")
	})

	t.Run("avatar belongs to other user", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: otherUserID}, nil
			},
		}
		svc := newTestService(repo)
		err := svc.DeleteAvatar(ctx, userID, avatarID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "avatar does not belong to user")
	})

	t.Run("delete error", func(t *testing.T) {
		t.Parallel()
		repo := &mockRepo{
			getAvatarByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUserAvatar, error) {
				return &db.SharedUserAvatar{ID: avatarID, UserID: userID}, nil
			},
			deleteAvatarFn: func(_ context.Context, _ uuid.UUID) error {
				return errors.New("db error")
			},
		}
		svc := newTestService(repo)
		err := svc.DeleteAvatar(ctx, userID, avatarID)
		require.Error(t, err)
	})
}

// ============================================================================
// Validation Tests (unexported functions)
// ============================================================================

func TestNoDB_validatePreferences(t *testing.T) {
	t.Parallel()
	svc := newTestService(&mockRepo{})

	tests := []struct {
		name    string
		params  UpsertPreferencesParams
		wantErr string
	}{
		{name: "valid theme light", params: UpsertPreferencesParams{Theme: new("light")}},
		{name: "valid theme dark", params: UpsertPreferencesParams{Theme: new("dark")}},
		{name: "valid theme system", params: UpsertPreferencesParams{Theme: new("system")}},
		{name: "invalid theme", params: UpsertPreferencesParams{Theme: new("rainbow")}, wantErr: "invalid theme"},
		{name: "valid visibility public", params: UpsertPreferencesParams{ProfileVisibility: new("public")}},
		{name: "valid visibility friends", params: UpsertPreferencesParams{ProfileVisibility: new("friends")}},
		{name: "valid visibility private", params: UpsertPreferencesParams{ProfileVisibility: new("private")}},
		{name: "invalid visibility", params: UpsertPreferencesParams{ProfileVisibility: new("hidden")}, wantErr: "invalid profile visibility"},
		{name: "nil theme and visibility passes", params: UpsertPreferencesParams{}},
		{name: "both valid", params: UpsertPreferencesParams{Theme: new("dark"), ProfileVisibility: new("public")}},
		{name: "invalid theme overrides valid visibility", params: UpsertPreferencesParams{Theme: new("neon"), ProfileVisibility: new("public")}, wantErr: "invalid theme"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

func TestNoDB_validateAvatarMetadata(t *testing.T) {
	t.Parallel()
	svc := newTestService(&mockRepo{})

	tests := []struct {
		name     string
		metadata AvatarMetadata
		wantErr  string
	}{
		{
			name:     "valid jpeg",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/jpeg", Width: 100, Height: 100},
		},
		{
			name:     "valid png",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 100, Height: 100},
		},
		{
			name:     "valid gif",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/gif", Width: 100, Height: 100},
		},
		{
			name:     "valid webp",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/webp", Width: 100, Height: 100},
		},
		{
			name:     "boundary min width",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 16, Height: 100},
		},
		{
			name:     "boundary max width",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 4096, Height: 100},
		},
		{
			name:     "boundary min height",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 100, Height: 16},
		},
		{
			name:     "boundary max height",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 100, Height: 4096},
		},
		{
			name:     "boundary max file size",
			metadata: AvatarMetadata{FileSizeBytes: 5 * 1024 * 1024, MimeType: "image/png", Width: 100, Height: 100},
		},
		{
			name:     "file too large",
			metadata: AvatarMetadata{FileSizeBytes: 5*1024*1024 + 1, MimeType: "image/png", Width: 100, Height: 100},
			wantErr:  "exceeds maximum",
		},
		{
			name:     "invalid mime type bmp",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/bmp", Width: 100, Height: 100},
			wantErr:  "invalid MIME type",
		},
		{
			name:     "invalid mime type pdf",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "application/pdf", Width: 100, Height: 100},
			wantErr:  "invalid MIME type",
		},
		{
			name:     "width too small",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 15, Height: 100},
			wantErr:  "invalid width",
		},
		{
			name:     "width too large",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 4097, Height: 100},
			wantErr:  "invalid width",
		},
		{
			name:     "height too small",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 100, Height: 15},
			wantErr:  "invalid height",
		},
		{
			name:     "height too large",
			metadata: AvatarMetadata{FileSizeBytes: 1024, MimeType: "image/png", Width: 100, Height: 4097},
			wantErr:  "invalid height",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

func TestNoDB_getDefaultPreferences(t *testing.T) {
	t.Parallel()
	svc := newTestService(&mockRepo{})
	userID := uuid.Must(uuid.NewV7())

	defaults := svc.getDefaultPreferences(userID)

	assert.Equal(t, userID, defaults.UserID)

	// Validate all default values
	require.NotNil(t, defaults.ProfileVisibility)
	assert.Equal(t, "private", *defaults.ProfileVisibility)

	require.NotNil(t, defaults.ShowEmail)
	assert.False(t, *defaults.ShowEmail)

	require.NotNil(t, defaults.ShowActivity)
	assert.True(t, *defaults.ShowActivity)

	require.NotNil(t, defaults.Theme)
	assert.Equal(t, "system", *defaults.Theme)

	require.NotNil(t, defaults.DisplayLanguage)
	assert.Equal(t, "en-US", *defaults.DisplayLanguage)

	require.NotNil(t, defaults.ShowAdultContent)
	assert.False(t, *defaults.ShowAdultContent)

	require.NotNil(t, defaults.ShowSpoilers)
	assert.False(t, *defaults.ShowSpoilers)

	require.NotNil(t, defaults.AutoPlayVideos)
	assert.True(t, *defaults.AutoPlayVideos)

	// Validate notification defaults are valid JSON
	require.NotNil(t, defaults.EmailNotifications)
	var emailNotif map[string]any
	err := json.Unmarshal(*defaults.EmailNotifications, &emailNotif)
	require.NoError(t, err)
	assert.True(t, emailNotif["enabled"].(bool))
	assert.Equal(t, "instant", emailNotif["frequency"].(string))

	require.NotNil(t, defaults.PushNotifications)
	var pushNotif map[string]any
	err = json.Unmarshal(*defaults.PushNotifications, &pushNotif)
	require.NoError(t, err)
	assert.False(t, pushNotif["enabled"].(bool))

	require.NotNil(t, defaults.DigestNotifications)
	var digestNotif map[string]any
	err = json.Unmarshal(*defaults.DigestNotifications, &digestNotif)
	require.NoError(t, err)
	assert.True(t, digestNotif["enabled"].(bool))
	assert.Equal(t, "weekly", digestNotif["frequency"].(string))
}

// ============================================================================
// CachedService Tests (no DB needed)
// ============================================================================

func TestNoDB_CachedService_NewCachedService(t *testing.T) {
	t.Parallel()
	svc := newTestService(&mockRepo{})
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)

	cached := NewCachedService(svc, testCache, logging.NewTestLogger())
	require.NotNil(t, cached)
	assert.NotNil(t, cached.Service)
	assert.NotNil(t, cached.cache)
	assert.NotNil(t, cached.logger)
}

func TestNoDB_CachedService_GetUser_NilCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			return makeUser(id, "testuser", "test@example.com"), nil
		},
	}
	svc := newTestService(repo)
	cached := NewCachedService(svc, nil, logging.NewTestLogger())

	result, err := cached.GetUser(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
}

func TestNoDB_CachedService_GetUser_CacheMissAndPopulate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	callCount := 0
	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			callCount++
			return makeUser(id, "testuser", "test@example.com"), nil
		},
	}
	svc := newTestService(repo)
	// L1 TTL must be <= UserTTL (1min) so Set() actually stores in L1
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	// First call: cache miss, hits repo
	result, err := cached.GetUser(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, 1, callCount, "first call should hit repo")

	// Wait for async cache population
	time.Sleep(50 * time.Millisecond)

	// Second call: cache hit, does NOT hit repo
	result2, err := cached.GetUser(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, result2.ID)
	assert.Equal(t, 1, callCount, "second call should hit cache, not repo")
}

func TestNoDB_CachedService_GetUser_RepoError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
			return nil, errors.New("db unavailable")
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	_, err = cached.GetUser(ctx, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db unavailable")
}

func TestNoDB_CachedService_GetUserByUsername_NilCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := &mockRepo{
		getUserByUsernameFn: func(_ context.Context, username string) (*db.SharedUser, error) {
			return makeUser(uuid.Must(uuid.NewV7()), username, "test@example.com"), nil
		},
	}
	svc := newTestService(repo)
	cached := NewCachedService(svc, nil, logging.NewTestLogger())

	result, err := cached.GetUserByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
}

func TestNoDB_CachedService_GetUserByUsername_CacheMissAndPopulate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	callCount := 0
	repo := &mockRepo{
		getUserByUsernameFn: func(_ context.Context, username string) (*db.SharedUser, error) {
			callCount++
			return makeUser(uuid.Must(uuid.NewV7()), username, "test@example.com"), nil
		},
	}
	svc := newTestService(repo)
	// L1 TTL must be <= UserTTL (1min) so Set() actually stores in L1
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	// First call: cache miss
	result, err := cached.GetUserByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, 1, callCount)

	// Wait for async cache population
	time.Sleep(50 * time.Millisecond)

	// Second call: cache hit
	result2, err := cached.GetUserByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", result2.Username)
	assert.Equal(t, 1, callCount, "second call should hit cache")
}

func TestNoDB_CachedService_GetUserByUsername_RepoError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := &mockRepo{
		getUserByUsernameFn: func(_ context.Context, _ string) (*db.SharedUser, error) {
			return nil, errors.New("db unavailable")
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	_, err = cached.GetUserByUsername(ctx, "testuser")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db unavailable")
}

func TestNoDB_CachedService_UpdateUser_NilCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			return makeUser(id, "testuser", "test@example.com"), nil
		},
		updateUserFn: func(_ context.Context, id uuid.UUID, _ UpdateUserParams) (*db.SharedUser, error) {
			u := makeUser(id, "testuser", "test@example.com")
			newName := "Updated"
			u.DisplayName = &newName
			return u, nil
		},
	}
	svc := newTestService(repo)
	cached := NewCachedService(svc, nil, logging.NewTestLogger())

	newName := "Updated"
	result, err := cached.UpdateUser(ctx, userID, UpdateUserParams{DisplayName: &newName})
	require.NoError(t, err)
	assert.Equal(t, "Updated", *result.DisplayName)
}

func TestNoDB_CachedService_UpdateUser_WithCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	callCount := 0
	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			callCount++
			return makeUser(id, "testuser", "test@example.com"), nil
		},
		updateUserFn: func(_ context.Context, id uuid.UUID, _ UpdateUserParams) (*db.SharedUser, error) {
			u := makeUser(id, "testuser", "test@example.com")
			newName := "Updated"
			u.DisplayName = &newName
			return u, nil
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	// Populate cache via GetUser
	_, err = cached.GetUser(ctx, userID)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond) // wait for async cache set

	// Update user (should invalidate cache)
	newName := "Updated"
	result, err := cached.UpdateUser(ctx, userID, UpdateUserParams{DisplayName: &newName})
	require.NoError(t, err)
	assert.Equal(t, "Updated", *result.DisplayName)

	// Wait for async cache invalidation
	time.Sleep(50 * time.Millisecond)

	// Next GetUser should hit repo again (cache was invalidated)
	prevCallCount := callCount
	_, err = cached.GetUser(ctx, userID)
	require.NoError(t, err)
	assert.Greater(t, callCount, prevCallCount, "repo should be called again after cache invalidation")
}

func TestNoDB_CachedService_UpdateUser_RepoError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
			return nil, errors.New("not found")
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	_, err = cached.UpdateUser(ctx, userID, UpdateUserParams{DisplayName: new("name")})
	require.Error(t, err)
}

func TestNoDB_CachedService_DeleteUser_NilCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			return makeUser(id, "testuser", "test@example.com"), nil
		},
		deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	svc := newTestService(repo)
	cached := NewCachedService(svc, nil, logging.NewTestLogger())

	err := cached.DeleteUser(ctx, userID)
	require.NoError(t, err)
}

func TestNoDB_CachedService_DeleteUser_WithCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			return makeUser(id, "testuser", "test@example.com"), nil
		},
		deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	// Populate cache first
	_, err = cached.GetUser(ctx, userID)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond) // wait for async cache set

	// Delete user (should invalidate cache)
	err = cached.DeleteUser(ctx, userID)
	require.NoError(t, err)

	// Wait for async cache invalidation
	time.Sleep(50 * time.Millisecond)
}

func TestNoDB_CachedService_DeleteUser_RepoError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, id uuid.UUID) (*db.SharedUser, error) {
			return makeUser(id, "testuser", "test@example.com"), nil
		},
		deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
			return errors.New("db error")
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	err = cached.DeleteUser(ctx, userID)
	require.Error(t, err)
}

func TestNoDB_CachedService_DeleteUser_UserNotFoundForLogging(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, _ uuid.UUID) (*db.SharedUser, error) {
			return nil, errors.New("not found")
		},
		deleteUserFn: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	svc := newTestService(repo)
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	// Should succeed even if user lookup fails
	err = cached.DeleteUser(ctx, userID)
	require.NoError(t, err)

	// Wait for async cache invalidation (user nil branch)
	time.Sleep(50 * time.Millisecond)
}

func TestNoDB_CachedService_InvalidateUserCache_NilCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	svc := newTestService(&mockRepo{})
	cached := NewCachedService(svc, nil, logging.NewTestLogger())

	err := cached.InvalidateUserCache(ctx, userID)
	require.NoError(t, err)
}

func TestNoDB_CachedService_InvalidateUserCache_WithCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := uuid.Must(uuid.NewV7())

	svc := newTestService(&mockRepo{})
	testCache, err := cache.NewCache(nil, 1000, 30*time.Second)
	require.NoError(t, err)
	defer testCache.Close()
	cached := NewCachedService(svc, testCache, logging.NewTestLogger())

	err = cached.InvalidateUserCache(ctx, userID)
	require.NoError(t, err)
}

// ============================================================================
// Repository Interface Compliance Test
// ============================================================================

func TestNoDB_MockRepoImplementsRepository(t *testing.T) {
	t.Parallel()
	var _ Repository = (*mockRepo)(nil)
}

// ============================================================================
// AvatarMetadata Struct Test
// ============================================================================

func TestNoDB_AvatarMetadata(t *testing.T) {
	t.Parallel()
	isAnimated := true
	ip := "192.168.1.1"
	ua := "TestAgent/1.0"

	metadata := AvatarMetadata{
		FileName:              "avatar.png",
		FileSizeBytes:         1024,
		MimeType:              "image/png",
		Width:                 256,
		Height:                256,
		IsAnimated:            &isAnimated,
		UploadedFromIP:        &ip,
		UploadedFromUserAgent: &ua,
	}

	assert.Equal(t, "avatar.png", metadata.FileName)
	assert.Equal(t, int64(1024), metadata.FileSizeBytes)
	assert.Equal(t, "image/png", metadata.MimeType)
	assert.Equal(t, int32(256), metadata.Width)
	assert.Equal(t, int32(256), metadata.Height)
	require.NotNil(t, metadata.IsAnimated)
	assert.True(t, *metadata.IsAnimated)
	require.NotNil(t, metadata.UploadedFromIP)
	assert.Equal(t, "192.168.1.1", *metadata.UploadedFromIP)
	require.NotNil(t, metadata.UploadedFromUserAgent)
	assert.Equal(t, "TestAgent/1.0", *metadata.UploadedFromUserAgent)
}
