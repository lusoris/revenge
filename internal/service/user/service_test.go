// Package user provides user management services for Jellyfin Go.
package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// mockUserRepository implements domain.UserRepository for testing.
type mockUserRepository struct {
	users         map[uuid.UUID]*domain.User
	usernameIndex map[string]uuid.UUID
	emailIndex    map[string]uuid.UUID
	createErr     error
	updateErr     error
	deleteErr     error
	adminCount    int64
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:         make(map[uuid.UUID]*domain.User),
		usernameIndex: make(map[string]uuid.UUID),
		emailIndex:    make(map[string]uuid.UUID),
		adminCount:    1,
	}
}

func (m *mockUserRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	if id, ok := m.usernameIndex[username]; ok {
		return m.users[id], nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	if id, ok := m.emailIndex[email]; ok {
		return m.users[id], nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) List(_ context.Context, limit, offset int32) ([]*domain.User, error) {
	var result []*domain.User
	i := int32(0)
	for _, u := range m.users {
		if i >= offset && int32(len(result)) < limit {
			result = append(result, u)
		}
		i++
	}
	return result, nil
}

func (m *mockUserRepository) Create(_ context.Context, params domain.CreateUserParams) (*domain.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}

	user := &domain.User{
		ID:        uuid.New(),
		Username:  params.Username,
		Email:     params.Email,
		IsAdmin:   params.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if params.PasswordHash != nil {
		user.PasswordHash = params.PasswordHash
	}
	if params.DisplayName != nil {
		user.DisplayName = params.DisplayName
	}

	m.users[user.ID] = user
	m.usernameIndex[user.Username] = user.ID
	if user.Email != nil {
		m.emailIndex[*user.Email] = user.ID
	}
	return user, nil
}

func (m *mockUserRepository) Update(_ context.Context, params domain.UpdateUserParams) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	user, ok := m.users[params.ID]
	if !ok {
		return domain.ErrUserNotFound
	}

	if params.Username != nil {
		delete(m.usernameIndex, user.Username)
		user.Username = *params.Username
		m.usernameIndex[user.Username] = user.ID
	}
	if params.Email != nil {
		if user.Email != nil {
			delete(m.emailIndex, *user.Email)
		}
		user.Email = params.Email
		m.emailIndex[*params.Email] = user.ID
	}
	if params.DisplayName != nil {
		user.DisplayName = params.DisplayName
	}
	if params.IsAdmin != nil {
		user.IsAdmin = *params.IsAdmin
	}
	if params.IsDisabled != nil {
		user.IsDisabled = *params.IsDisabled
	}
	user.UpdatedAt = time.Now()

	return nil
}

func (m *mockUserRepository) Delete(_ context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	user, ok := m.users[id]
	if !ok {
		return domain.ErrUserNotFound
	}

	delete(m.usernameIndex, user.Username)
	if user.Email != nil {
		delete(m.emailIndex, *user.Email)
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) UpdateLastLogin(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (m *mockUserRepository) UpdateLastActivity(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (m *mockUserRepository) SetPassword(_ context.Context, id uuid.UUID, passwordHash string) error {
	if user, ok := m.users[id]; ok {
		user.PasswordHash = &passwordHash
		return nil
	}
	return domain.ErrUserNotFound
}

func (m *mockUserRepository) Count(_ context.Context) (int64, error) {
	return int64(len(m.users)), nil
}

func (m *mockUserRepository) CountAdmins(_ context.Context) (int64, error) {
	return m.adminCount, nil
}

func (m *mockUserRepository) UsernameExists(_ context.Context, username string) (bool, error) {
	_, ok := m.usernameIndex[username]
	return ok, nil
}

func (m *mockUserRepository) EmailExists(_ context.Context, email string) (bool, error) {
	_, ok := m.emailIndex[email]
	return ok, nil
}

// mockSessionRepository implements domain.SessionRepository for testing.
type mockSessionRepository struct {
	sessions map[uuid.UUID]*domain.Session
}

func newMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{
		sessions: make(map[uuid.UUID]*domain.Session),
	}
}

func (m *mockSessionRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Session, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetByTokenHash(_ context.Context, _ string) (*domain.Session, error) {
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetByRefreshTokenHash(_ context.Context, _ string) (*domain.Session, error) {
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetWithUser(_ context.Context, _ string) (*domain.SessionWithUser, error) {
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) ListByUser(_ context.Context, _ uuid.UUID) ([]*domain.Session, error) {
	return nil, nil
}

func (m *mockSessionRepository) Create(_ context.Context, params domain.CreateSessionParams) (*domain.Session, error) {
	s := &domain.Session{
		ID:        uuid.New(),
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: time.Now(),
	}
	m.sessions[s.ID] = s
	return s, nil
}

func (m *mockSessionRepository) UpdateRefreshToken(_ context.Context, _ uuid.UUID, _ string, _ time.Time) error {
	return nil
}

func (m *mockSessionRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.sessions, id)
	return nil
}

func (m *mockSessionRepository) DeleteByTokenHash(_ context.Context, _ string) error {
	return nil
}

func (m *mockSessionRepository) DeleteByUser(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (m *mockSessionRepository) DeleteExpired(_ context.Context) (int64, error) {
	return 0, nil
}

func (m *mockSessionRepository) CountByUser(_ context.Context, _ uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockSessionRepository) Exists(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// mockPasswordService implements domain.PasswordService for testing.
type mockPasswordService struct {
	hashErr   error
	verifyErr error
}

func (m *mockPasswordService) Hash(password string) (string, error) {
	if m.hashErr != nil {
		return "", m.hashErr
	}
	return "hashed_" + password, nil
}

func (m *mockPasswordService) Verify(_, _ string) error {
	return m.verifyErr
}

func TestService_GetByID(t *testing.T) {
	repo := newMockUserRepository()
	svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

	// Create a test user
	user := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.users[user.ID] = user
	repo.usernameIndex[user.Username] = user.ID

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr error
	}{
		{
			name:    "existing user",
			id:      user.ID,
			wantErr: nil,
		},
		{
			name:    "non-existing user",
			id:      uuid.New(),
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetByID(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && result.ID != tt.id {
				t.Errorf("GetByID() got ID = %v, want %v", result.ID, tt.id)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name     string
		params   CreateParams
		setup    func(*mockUserRepository)
		wantErr  bool
		errMatch string
	}{
		{
			name: "valid user",
			params: CreateParams{
				Username: "newuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty username",
			params: CreateParams{
				Username: "",
				Password: "password123",
			},
			wantErr:  true,
			errMatch: "username is required",
		},
		{
			name: "username too short",
			params: CreateParams{
				Username: "ab",
				Password: "password123",
			},
			wantErr:  true,
			errMatch: "at least 3 characters",
		},
		{
			name: "password too short",
			params: CreateParams{
				Username: "validuser",
				Password: "short",
			},
			wantErr:  true,
			errMatch: "at least 8 characters",
		},
		{
			name: "duplicate username",
			params: CreateParams{
				Username: "existinguser",
				Password: "password123",
			},
			setup: func(m *mockUserRepository) {
				m.usernameIndex["existinguser"] = uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

			result, err := svc.Create(context.Background(), tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errMatch != "" && err != nil {
				if !contains(err.Error(), tt.errMatch) {
					t.Errorf("Create() error = %v, should contain %v", err, tt.errMatch)
				}
			}
			if !tt.wantErr && result == nil {
				t.Error("Create() returned nil user")
			}
		})
	}
}

func TestService_List(t *testing.T) {
	repo := newMockUserRepository()
	svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

	// Create test users
	for i := 0; i < 5; i++ {
		user := &domain.User{
			ID:        uuid.New(),
			Username:  "user" + string(rune('0'+i)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repo.users[user.ID] = user
		repo.usernameIndex[user.Username] = user.ID
	}

	tests := []struct {
		name      string
		limit     int32
		offset    int32
		wantCount int
		wantLimit int32 // expected corrected limit
	}{
		{
			name:      "default limit",
			limit:     0,
			offset:    0,
			wantCount: 5,
			wantLimit: 50,
		},
		{
			name:      "max limit capped",
			limit:     200,
			offset:    0,
			wantCount: 5,
			wantLimit: 100,
		},
		{
			name:      "negative offset corrected",
			limit:     10,
			offset:    -5,
			wantCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.List(context.Background(), tt.limit, tt.offset)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}
			if len(result) > tt.wantCount {
				t.Errorf("List() got %d users, want at most %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*mockUserRepository) uuid.UUID
		wantErr     bool
		errContains string
	}{
		{
			name: "delete regular user",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "regularuser",
					IsAdmin:  false,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			wantErr: false,
		},
		{
			name: "cannot delete last admin",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "lastadmin",
					IsAdmin:  true,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				m.adminCount = 1
				return user.ID
			},
			wantErr:     true,
			errContains: "last administrator",
		},
		{
			name: "delete non-existing user",
			setup: func(_ *mockUserRepository) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			userID := tt.setup(repo)
			svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

			err := svc.Delete(context.Background(), userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Delete() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestService_SetPassword(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*mockUserRepository) uuid.UUID
		newPassword string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid password change",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "testuser",
				}
				m.users[user.ID] = user
				return user.ID
			},
			newPassword: "newpassword123",
			wantErr:     false,
		},
		{
			name: "password too short",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "testuser",
				}
				m.users[user.ID] = user
				return user.ID
			},
			newPassword: "short",
			wantErr:     true,
			errContains: "at least 8 characters",
		},
		{
			name: "user not found",
			setup: func(_ *mockUserRepository) uuid.UUID {
				return uuid.New()
			},
			newPassword: "validpassword",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			userID := tt.setup(repo)
			svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

			err := svc.SetPassword(context.Background(), userID, tt.newPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("SetPassword() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestService_GetByUsername(t *testing.T) {
	repo := newMockUserRepository()
	svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

	// Create a test user
	user := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.users[user.ID] = user
	repo.usernameIndex[user.Username] = user.ID

	tests := []struct {
		name     string
		username string
		wantErr  error
	}{
		{
			name:     "existing user",
			username: "testuser",
			wantErr:  nil,
		},
		{
			name:     "non-existing user",
			username: "nonexistent",
			wantErr:  domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetByUsername(context.Background(), tt.username)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && result.Username != tt.username {
				t.Errorf("GetByUsername() got Username = %v, want %v", result.Username, tt.username)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*mockUserRepository) uuid.UUID
		params      func(uuid.UUID) UpdateParams
		wantErr     bool
		errContains string
	}{
		{
			name: "update username",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "oldname",
					IsAdmin:  false,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				newName := "newname"
				return UpdateParams{ID: id, Username: &newName}
			},
			wantErr: false,
		},
		{
			name: "update to existing username",
			setup: func(m *mockUserRepository) uuid.UUID {
				// Create existing user with target username
				existing := &domain.User{
					ID:       uuid.New(),
					Username: "existingname",
				}
				m.users[existing.ID] = existing
				m.usernameIndex[existing.Username] = existing.ID

				// Create user to update
				user := &domain.User{
					ID:       uuid.New(),
					Username: "myname",
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				duplicate := "existingname"
				return UpdateParams{ID: id, Username: &duplicate}
			},
			wantErr: true,
		},
		{
			name: "update to empty username",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "validname",
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				empty := "   "
				return UpdateParams{ID: id, Username: &empty}
			},
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name: "update username too short",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "validname",
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				short := "ab"
				return UpdateParams{ID: id, Username: &short}
			},
			wantErr:     true,
			errContains: "at least 3 characters",
		},
		{
			name: "update email",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "testuser",
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				email := "newemail@test.com"
				return UpdateParams{ID: id, Email: &email}
			},
			wantErr: false,
		},
		{
			name: "update to duplicate email",
			setup: func(m *mockUserRepository) uuid.UUID {
				// Create existing user with target email
				existingEmail := "taken@test.com"
				existing := &domain.User{
					ID:       uuid.New(),
					Username: "other",
					Email:    &existingEmail,
				}
				m.users[existing.ID] = existing
				m.usernameIndex[existing.Username] = existing.ID
				m.emailIndex[existingEmail] = existing.ID

				// Create user to update
				user := &domain.User{
					ID:       uuid.New(),
					Username: "testuser",
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				taken := "taken@test.com"
				return UpdateParams{ID: id, Email: &taken}
			},
			wantErr: true,
		},
		{
			name: "update same email allowed",
			setup: func(m *mockUserRepository) uuid.UUID {
				email := "my@test.com"
				user := &domain.User{
					ID:       uuid.New(),
					Username: "testuser",
					Email:    &email,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				m.emailIndex[email] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				same := "my@test.com"
				return UpdateParams{ID: id, Email: &same}
			},
			wantErr: false,
		},
		{
			name: "remove admin from last admin",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "lastadmin",
					IsAdmin:  true,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				m.adminCount = 1
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				isAdmin := false
				return UpdateParams{ID: id, IsAdmin: &isAdmin}
			},
			wantErr:     true,
			errContains: "last administrator",
		},
		{
			name: "remove admin when multiple admins",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "admin1",
					IsAdmin:  true,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				m.adminCount = 2
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				isAdmin := false
				return UpdateParams{ID: id, IsAdmin: &isAdmin}
			},
			wantErr: false,
		},
		{
			name: "disable user",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:         uuid.New(),
					Username:   "activeuser",
					IsDisabled: false,
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				disabled := true
				return UpdateParams{ID: id, IsDisabled: &disabled}
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setup: func(_ *mockUserRepository) uuid.UUID {
				return uuid.New()
			},
			params: func(id uuid.UUID) UpdateParams {
				name := "newname"
				return UpdateParams{ID: id, Username: &name}
			},
			wantErr: true,
		},
		{
			name: "update same username allowed",
			setup: func(m *mockUserRepository) uuid.UUID {
				user := &domain.User{
					ID:       uuid.New(),
					Username: "sameuser",
				}
				m.users[user.ID] = user
				m.usernameIndex[user.Username] = user.ID
				return user.ID
			},
			params: func(id uuid.UUID) UpdateParams {
				same := "sameuser"
				return UpdateParams{ID: id, Username: &same}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			userID := tt.setup(repo)
			svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

			err := svc.Update(context.Background(), tt.params(userID))
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Update() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestService_Count(t *testing.T) {
	repo := newMockUserRepository()
	svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

	// Create test users
	for i := 0; i < 3; i++ {
		user := &domain.User{
			ID:       uuid.New(),
			Username: "user" + string(rune('0'+i)),
		}
		repo.users[user.ID] = user
		repo.usernameIndex[user.Username] = user.ID
	}

	count, err := svc.Count(context.Background())
	if err != nil {
		t.Errorf("Count() error = %v", err)
	}
	if count != 3 {
		t.Errorf("Count() = %d, want 3", count)
	}
}

func TestService_Create_WithEmail(t *testing.T) {
	tests := []struct {
		name     string
		params   CreateParams
		setup    func(*mockUserRepository)
		wantErr  bool
		errMatch string
	}{
		{
			name: "with valid email",
			params: CreateParams{
				Username: "newuser",
				Password: "password123",
				Email:    strPtr("test@example.com"),
			},
			wantErr: false,
		},
		{
			name: "with duplicate email",
			params: CreateParams{
				Username: "newuser",
				Password: "password123",
				Email:    strPtr("taken@example.com"),
			},
			setup: func(m *mockUserRepository) {
				m.emailIndex["taken@example.com"] = uuid.New()
			},
			wantErr: true,
		},
		{
			name: "with empty email string",
			params: CreateParams{
				Username: "newuser",
				Password: "password123",
				Email:    strPtr(""),
			},
			wantErr: false, // Empty string is treated as no email
		},
		{
			name: "without password",
			params: CreateParams{
				Username: "oidcuser",
				Password: "", // No password
			},
			wantErr: false, // Allowed for OIDC users
		},
		{
			name: "username too long",
			params: CreateParams{
				Username: "thisusernameiswaywaytoolongtobevalidandshouldfailvalidation1234567890",
				Password: "password123",
			},
			wantErr:  true,
			errMatch: "at most 64 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := newService(repo, newMockSessionRepository(), &mockPasswordService{})

			result, err := svc.Create(context.Background(), tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errMatch != "" && err != nil {
				if !contains(err.Error(), tt.errMatch) {
					t.Errorf("Create() error = %v, should contain %v", err, tt.errMatch)
				}
			}
			if !tt.wantErr && result == nil {
				t.Error("Create() returned nil user")
			}
		})
	}
}

func TestService_Create_HashError(t *testing.T) {
	repo := newMockUserRepository()
	svc := newService(repo, newMockSessionRepository(), &mockPasswordService{
		hashErr: errors.New("hash failed"),
	})

	_, err := svc.Create(context.Background(), CreateParams{
		Username: "newuser",
		Password: "password123",
	})
	if err == nil {
		t.Error("Create() expected error for hash failure")
	}
	if !contains(err.Error(), "hash password") {
		t.Errorf("Create() error = %v, should mention hash password", err)
	}
}

func strPtr(s string) *string {
	return &s
}
