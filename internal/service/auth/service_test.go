// Package auth provides authentication services for Jellyfin Go.
package auth

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

	m.users[user.ID] = user
	m.usernameIndex[user.Username] = user.ID
	if user.Email != nil {
		m.emailIndex[*user.Email] = user.ID
	}
	return user, nil
}

func (m *mockUserRepository) Update(_ context.Context, params domain.UpdateUserParams) error {
	user, ok := m.users[params.ID]
	if !ok {
		return domain.ErrUserNotFound
	}

	if params.Username != nil {
		delete(m.usernameIndex, user.Username)
		user.Username = *params.Username
		m.usernameIndex[user.Username] = user.ID
	}
	if params.IsDisabled != nil {
		user.IsDisabled = *params.IsDisabled
	}
	user.UpdatedAt = time.Now()
	return nil
}

func (m *mockUserRepository) Delete(_ context.Context, id uuid.UUID) error {
	user, ok := m.users[id]
	if !ok {
		return domain.ErrUserNotFound
	}
	delete(m.usernameIndex, user.Username)
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

func (m *mockUserRepository) addUser(user *domain.User) {
	m.users[user.ID] = user
	m.usernameIndex[user.Username] = user.ID
	if user.Email != nil {
		m.emailIndex[*user.Email] = user.ID
	}
}

// mockSessionRepository implements domain.SessionRepository for testing.
type mockSessionRepository struct {
	sessions         map[uuid.UUID]*domain.Session
	tokenHashIndex   map[string]uuid.UUID
	refreshHashIndex map[string]uuid.UUID
}

func newMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{
		sessions:         make(map[uuid.UUID]*domain.Session),
		tokenHashIndex:   make(map[string]uuid.UUID),
		refreshHashIndex: make(map[string]uuid.UUID),
	}
}

func (m *mockSessionRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Session, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetByTokenHash(_ context.Context, tokenHash string) (*domain.Session, error) {
	if id, ok := m.tokenHashIndex[tokenHash]; ok {
		return m.sessions[id], nil
	}
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetByRefreshTokenHash(_ context.Context, refreshHash string) (*domain.Session, error) {
	if id, ok := m.refreshHashIndex[refreshHash]; ok {
		return m.sessions[id], nil
	}
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) GetWithUser(_ context.Context, tokenHash string) (*domain.SessionWithUser, error) {
	if id, ok := m.tokenHashIndex[tokenHash]; ok {
		s := m.sessions[id]
		return &domain.SessionWithUser{
			Session:  *s,
			Username: "testuser",
			IsAdmin:  false,
		}, nil
	}
	return nil, domain.ErrSessionNotFound
}

func (m *mockSessionRepository) ListByUser(_ context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	var result []*domain.Session
	for _, s := range m.sessions {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockSessionRepository) Create(_ context.Context, params domain.CreateSessionParams) (*domain.Session, error) {
	s := &domain.Session{
		ID:        uuid.New(),
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: time.Now(),
	}
	if params.RefreshTokenHash != nil {
		s.RefreshTokenHash = params.RefreshTokenHash
		m.refreshHashIndex[*params.RefreshTokenHash] = s.ID
	}
	if params.RefreshExpiresAt != nil {
		s.RefreshExpiresAt = params.RefreshExpiresAt
	}
	m.sessions[s.ID] = s
	m.tokenHashIndex[params.TokenHash] = s.ID
	return s, nil
}

func (m *mockSessionRepository) UpdateRefreshToken(_ context.Context, id uuid.UUID, refreshHash string, expiresAt time.Time) error {
	if s, ok := m.sessions[id]; ok {
		if s.RefreshTokenHash != nil {
			delete(m.refreshHashIndex, *s.RefreshTokenHash)
		}
		s.RefreshTokenHash = &refreshHash
		s.RefreshExpiresAt = &expiresAt
		m.refreshHashIndex[refreshHash] = id
		return nil
	}
	return domain.ErrSessionNotFound
}

func (m *mockSessionRepository) Delete(_ context.Context, id uuid.UUID) error {
	if s, ok := m.sessions[id]; ok {
		delete(m.tokenHashIndex, s.TokenHash)
		if s.RefreshTokenHash != nil {
			delete(m.refreshHashIndex, *s.RefreshTokenHash)
		}
		delete(m.sessions, id)
		return nil
	}
	return nil // Not an error if not found
}

func (m *mockSessionRepository) DeleteByTokenHash(_ context.Context, tokenHash string) error {
	if id, ok := m.tokenHashIndex[tokenHash]; ok {
		return m.Delete(context.Background(), id)
	}
	return nil
}

func (m *mockSessionRepository) DeleteByUser(_ context.Context, userID uuid.UUID) error {
	for id, s := range m.sessions {
		if s.UserID == userID {
			delete(m.tokenHashIndex, s.TokenHash)
			if s.RefreshTokenHash != nil {
				delete(m.refreshHashIndex, *s.RefreshTokenHash)
			}
			delete(m.sessions, id)
		}
	}
	return nil
}

func (m *mockSessionRepository) DeleteExpired(_ context.Context) (int64, error) {
	return 0, nil
}

func (m *mockSessionRepository) CountByUser(_ context.Context, userID uuid.UUID) (int64, error) {
	count := int64(0)
	for _, s := range m.sessions {
		if s.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockSessionRepository) Exists(_ context.Context, tokenHash string) (bool, error) {
	_, ok := m.tokenHashIndex[tokenHash]
	return ok, nil
}

// mockPasswordServiceForAuth implements domain.PasswordService for testing.
type mockPasswordServiceForAuth struct {
	hashErr   error
	verifyErr error
}

func (m *mockPasswordServiceForAuth) Hash(password string) (string, error) {
	if m.hashErr != nil {
		return "", m.hashErr
	}
	return "hashed_" + password, nil
}

func (m *mockPasswordServiceForAuth) Verify(password, hash string) error {
	if m.verifyErr != nil {
		return m.verifyErr
	}
	if hash == "hashed_"+password {
		return nil
	}
	return domain.ErrInvalidCredentials
}

// mockTokenServiceForAuth implements domain.TokenService for testing.
type mockTokenServiceForAuth struct {
	accessToken string
	claims      *domain.TokenClaims
	validateErr error
}

func (m *mockTokenServiceForAuth) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	if m.accessToken != "" {
		return m.accessToken, nil
	}
	return "access_token_" + claims.UserID.String(), nil
}

func (m *mockTokenServiceForAuth) GenerateRefreshToken() (string, error) {
	return "refresh_token_" + uuid.New().String(), nil
}

func (m *mockTokenServiceForAuth) ValidateAccessToken(token string) (*domain.TokenClaims, error) {
	if m.validateErr != nil {
		return nil, m.validateErr
	}
	if m.claims != nil {
		return m.claims, nil
	}
	return &domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   false,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, nil
}

// errInvalidToken is a test sentinel for invalid tokens.
var errInvalidToken = errors.New("invalid token")

func (m *mockTokenServiceForAuth) HashToken(token string) string {
	return "hash_" + token
}

func (m *mockTokenServiceForAuth) AccessTokenDuration() time.Duration {
	return 15 * time.Minute
}

func TestService_Login(t *testing.T) {
	passwordHash := "hashed_password123"

	tests := []struct {
		name      string
		username  string
		password  string
		setupUser func(*mockUserRepository)
		verifyErr error
		wantErr   error
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "password123",
			setupUser: func(m *mockUserRepository) {
				m.addUser(&domain.User{
					ID:           uuid.New(),
					Username:     "testuser",
					PasswordHash: &passwordHash,
					IsDisabled:   false,
				})
			},
			wantErr: nil,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "password123",
			setupUser: func(_ *mockUserRepository) {
				// No user added
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name:     "disabled user",
			username: "disableduser",
			password: "password123",
			setupUser: func(m *mockUserRepository) {
				m.addUser(&domain.User{
					ID:           uuid.New(),
					Username:     "disableduser",
					PasswordHash: &passwordHash,
					IsDisabled:   true,
				})
			},
			wantErr: domain.ErrUserDisabled,
		},
		{
			name:     "invalid password",
			username: "testuser",
			password: "wrongpassword",
			setupUser: func(m *mockUserRepository) {
				m.addUser(&domain.User{
					ID:           uuid.New(),
					Username:     "testuser",
					PasswordHash: &passwordHash,
					IsDisabled:   false,
				})
			},
			verifyErr: domain.ErrInvalidCredentials,
			wantErr:   domain.ErrInvalidCredentials,
		},
		{
			name:     "user without password (OIDC only)",
			username: "oidcuser",
			password: "anypassword",
			setupUser: func(m *mockUserRepository) {
				m.addUser(&domain.User{
					ID:           uuid.New(),
					Username:     "oidcuser",
					PasswordHash: nil, // No password
					IsDisabled:   false,
				})
			},
			wantErr: domain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepository()
			sessionRepo := newMockSessionRepository()
			passwordSvc := &mockPasswordServiceForAuth{verifyErr: tt.verifyErr}
			tokenSvc := &mockTokenServiceForAuth{}

			if tt.setupUser != nil {
				tt.setupUser(userRepo)
			}

			svc := newService(userRepo, sessionRepo, passwordSvc, tokenSvc, 0, 15*time.Minute, 7*24*time.Hour)

			result, err := svc.Login(context.Background(), domain.LoginParams{
				Username: tt.username,
				Password: tt.password,
			})

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if result == nil {
					t.Error("Login() returned nil result")
					return
				}
				if result.AccessToken == "" {
					t.Error("Login() returned empty access token")
				}
				if result.RefreshToken == "" {
					t.Error("Login() returned empty refresh token")
				}
				if result.User == nil {
					t.Error("Login() returned nil user")
				}
			}
		})
	}
}

func TestService_Logout(t *testing.T) {
	tests := []struct {
		name         string
		accessToken  string
		setupSession func(*mockSessionRepository, *mockTokenServiceForAuth)
		wantErr      bool
	}{
		{
			name:        "successful logout",
			accessToken: "valid_token",
			setupSession: func(sr *mockSessionRepository, ts *mockTokenServiceForAuth) {
				sessionID := uuid.New()
				userID := uuid.New()
				tokenHash := "hash_valid_token"

				sr.sessions[sessionID] = &domain.Session{
					ID:        sessionID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(time.Hour),
				}
				sr.tokenHashIndex[tokenHash] = sessionID

				ts.claims = &domain.TokenClaims{
					UserID:    userID,
					SessionID: sessionID,
					Username:  "testuser",
				}
			},
			wantErr: false,
		},
		{
			name:        "logout with non-existent token succeeds silently",
			accessToken: "nonexistent_token",
			setupSession: func(_ *mockSessionRepository, _ *mockTokenServiceForAuth) {
				// No session - logout is idempotent, succeeds even if no session exists
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepository()
			sessionRepo := newMockSessionRepository()
			passwordSvc := &mockPasswordServiceForAuth{}
			tokenSvc := &mockTokenServiceForAuth{}

			if tt.setupSession != nil {
				tt.setupSession(sessionRepo, tokenSvc)
			}

			svc := newService(userRepo, sessionRepo, passwordSvc, tokenSvc, 0, 15*time.Minute, 7*24*time.Hour)

			err := svc.Logout(context.Background(), tt.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_LogoutAll(t *testing.T) {
	userID := uuid.New()

	userRepo := newMockUserRepository()
	sessionRepo := newMockSessionRepository()
	passwordSvc := &mockPasswordServiceForAuth{}
	tokenSvc := &mockTokenServiceForAuth{}

	// Create multiple sessions for the user
	for i := 0; i < 3; i++ {
		sessionID := uuid.New()
		tokenHash := "hash_" + uuid.New().String()
		sessionRepo.sessions[sessionID] = &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		sessionRepo.tokenHashIndex[tokenHash] = sessionID
	}

	svc := newService(userRepo, sessionRepo, passwordSvc, tokenSvc, 0, 15*time.Minute, 7*24*time.Hour)

	err := svc.LogoutAll(context.Background(), userID)
	if err != nil {
		t.Errorf("LogoutAll() error = %v", err)
	}

	// Verify all sessions for the user are deleted
	count, _ := sessionRepo.CountByUser(context.Background(), userID)
	if count != 0 {
		t.Errorf("LogoutAll() left %d sessions, want 0", count)
	}
}

func TestService_ValidateToken(t *testing.T) {
	tests := []struct {
		name        string
		accessToken string
		setupMock   func(*mockTokenServiceForAuth, *mockSessionRepository)
		wantErr     bool
	}{
		{
			name:        "valid token with existing session",
			accessToken: "valid_token",
			setupMock: func(ts *mockTokenServiceForAuth, sr *mockSessionRepository) {
				sessionID := uuid.New()
				tokenHash := "hash_valid_token"

				ts.claims = &domain.TokenClaims{
					UserID:    uuid.New(),
					SessionID: sessionID,
					Username:  "testuser",
					ExpiresAt: time.Now().Add(time.Hour),
				}

				sr.sessions[sessionID] = &domain.Session{
					ID:        sessionID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(time.Hour),
				}
				sr.tokenHashIndex[tokenHash] = sessionID
			},
			wantErr: false,
		},
		{
			name:        "invalid token",
			accessToken: "invalid_token",
			setupMock: func(ts *mockTokenServiceForAuth, _ *mockSessionRepository) {
				ts.validateErr = errInvalidToken
			},
			wantErr: true,
		},
		{
			name:        "valid token but session not found",
			accessToken: "orphan_token",
			setupMock: func(ts *mockTokenServiceForAuth, _ *mockSessionRepository) {
				ts.claims = &domain.TokenClaims{
					UserID:    uuid.New(),
					SessionID: uuid.New(),
					Username:  "testuser",
					ExpiresAt: time.Now().Add(time.Hour),
				}
				// No session in repository
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepository()
			sessionRepo := newMockSessionRepository()
			passwordSvc := &mockPasswordServiceForAuth{}
			tokenSvc := &mockTokenServiceForAuth{}

			if tt.setupMock != nil {
				tt.setupMock(tokenSvc, sessionRepo)
			}

			svc := newService(userRepo, sessionRepo, passwordSvc, tokenSvc, 0, 15*time.Minute, 7*24*time.Hour)

			claims, err := svc.ValidateToken(context.Background(), tt.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && claims == nil {
				t.Error("ValidateToken() returned nil claims")
			}
		})
	}
}

func TestService_ChangePassword(t *testing.T) {
	userID := uuid.New()
	currentHash := "hashed_currentpass"

	tests := []struct {
		name            string
		currentPassword string
		newPassword     string
		setupUser       func(*mockUserRepository)
		wantErr         bool
	}{
		{
			name:            "successful password change",
			currentPassword: "currentpass",
			newPassword:     "newpassword123",
			setupUser: func(m *mockUserRepository) {
				m.addUser(&domain.User{
					ID:           userID,
					Username:     "testuser",
					PasswordHash: &currentHash,
				})
			},
			wantErr: false,
		},
		{
			name:            "wrong current password",
			currentPassword: "wrongpass",
			newPassword:     "newpassword123",
			setupUser: func(m *mockUserRepository) {
				m.addUser(&domain.User{
					ID:           userID,
					Username:     "testuser",
					PasswordHash: &currentHash,
				})
			},
			wantErr: true,
		},
		{
			name:            "user not found",
			currentPassword: "currentpass",
			newPassword:     "newpassword123",
			setupUser:       func(_ *mockUserRepository) {},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepository()
			sessionRepo := newMockSessionRepository()
			passwordSvc := &mockPasswordServiceForAuth{}
			tokenSvc := &mockTokenServiceForAuth{}

			if tt.setupUser != nil {
				tt.setupUser(userRepo)
			}

			svc := newService(userRepo, sessionRepo, passwordSvc, tokenSvc, 0, 15*time.Minute, 7*24*time.Hour)

			err := svc.ChangePassword(context.Background(), userID, tt.currentPassword, tt.newPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
