package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

// mockAuthService is a mock implementation of domain.AuthService for testing.
type mockAuthService struct {
	validateTokenFn func(ctx context.Context, token string) (*domain.TokenClaims, error)
	getSessionFn    func(ctx context.Context, token string) (*domain.SessionWithUser, error)
}

func (m *mockAuthService) Login(ctx context.Context, params domain.LoginParams) (*domain.AuthResult, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) Logout(ctx context.Context, token string) error {
	return errors.New("not implemented")
}

func (m *mockAuthService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return errors.New("not implemented")
}

func (m *mockAuthService) RefreshToken(ctx context.Context, params domain.RefreshParams) (*domain.AuthResult, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	if m.validateTokenFn != nil {
		return m.validateTokenFn(ctx, token)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) GetSession(ctx context.Context, token string) (*domain.SessionWithUser, error) {
	if m.getSessionFn != nil {
		return m.getSessionFn(ctx, token)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	return errors.New("not implemented")
}

func (m *mockAuthService) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return errors.New("not implemented")
}

func TestAuth_Required(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()
	validClaims := &domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Username:  "testuser",
		IsAdmin:   false,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	tests := []struct {
		name          string
		authHeader    string
		validateFn    func(ctx context.Context, token string) (*domain.TokenClaims, error)
		wantStatus    int
		wantClaimsSet bool
	}{
		{
			name:       "no auth header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid bearer token",
			authHeader: "Bearer valid-token",
			validateFn: func(ctx context.Context, token string) (*domain.TokenClaims, error) {
				if token == "valid-token" {
					return validClaims, nil
				}
				return nil, errors.New("invalid token")
			},
			wantStatus:    http.StatusOK,
			wantClaimsSet: true,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid-token",
			validateFn: func(ctx context.Context, token string) (*domain.TokenClaims, error) {
				return nil, errors.New("invalid token")
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "mediabrowser format",
			authHeader: `MediaBrowser Client="Test", Device="Test", DeviceId="123", Version="1.0", Token="mb-token"`,
			validateFn: func(ctx context.Context, token string) (*domain.TokenClaims, error) {
				if token == "mb-token" {
					return validClaims, nil
				}
				return nil, errors.New("invalid token")
			},
			wantStatus:    http.StatusOK,
			wantClaimsSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAuthService{validateTokenFn: tt.validateFn}
			auth := NewAuth(mock)

			var gotClaims *domain.TokenClaims
			handler := auth.Required(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotClaims = ClaimsFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantClaimsSet && gotClaims == nil {
				t.Error("expected claims to be set in context")
			}
			if !tt.wantClaimsSet && gotClaims != nil {
				t.Error("expected claims to NOT be set in context")
			}
		})
	}
}

func TestAuth_Optional(t *testing.T) {
	validClaims := &domain.TokenClaims{
		UserID:    uuid.New(),
		SessionID: uuid.New(),
		Username:  "testuser",
		IsAdmin:   false,
	}

	tests := []struct {
		name          string
		authHeader    string
		validateFn    func(ctx context.Context, token string) (*domain.TokenClaims, error)
		wantStatus    int
		wantClaimsSet bool
	}{
		{
			name:       "no auth header - request proceeds",
			authHeader: "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid token - claims set",
			authHeader: "Bearer valid-token",
			validateFn: func(ctx context.Context, token string) (*domain.TokenClaims, error) {
				return validClaims, nil
			},
			wantStatus:    http.StatusOK,
			wantClaimsSet: true,
		},
		{
			name:       "invalid token - request still proceeds",
			authHeader: "Bearer invalid-token",
			validateFn: func(ctx context.Context, token string) (*domain.TokenClaims, error) {
				return nil, errors.New("invalid")
			},
			wantStatus:    http.StatusOK,
			wantClaimsSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAuthService{validateTokenFn: tt.validateFn}
			auth := NewAuth(mock)

			var gotClaims *domain.TokenClaims
			handler := auth.Optional(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotClaims = ClaimsFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantClaimsSet && gotClaims == nil {
				t.Error("expected claims to be set")
			}
			if !tt.wantClaimsSet && gotClaims != nil {
				t.Error("expected claims to NOT be set")
			}
		})
	}
}

func TestAuth_AdminRequired(t *testing.T) {
	tests := []struct {
		name       string
		claims     *domain.TokenClaims
		wantStatus int
	}{
		{
			name:       "no claims in context",
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-admin user",
			claims: &domain.TokenClaims{
				UserID:  uuid.New(),
				IsAdmin: false,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "admin user",
			claims: &domain.TokenClaims{
				UserID:  uuid.New(),
				IsAdmin: true,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAuthService{}
			auth := NewAuth(mock)

			handler := auth.AdminRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.claims != nil {
				ctx := context.WithValue(req.Context(), ClaimsContextKey, tt.claims)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(r *http.Request)
		wantToken string
	}{
		{
			name:      "no token",
			setup:     func(r *http.Request) {},
			wantToken: "",
		},
		{
			name: "bearer token",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer my-token")
			},
			wantToken: "my-token",
		},
		{
			name: "mediabrowser token",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", `MediaBrowser Token="mb-token"`)
			},
			wantToken: "mb-token",
		},
		{
			name: "x-emby-authorization header",
			setup: func(r *http.Request) {
				r.Header.Set("X-Emby-Authorization", `MediaBrowser Token="emby-token"`)
			},
			wantToken: "emby-token",
		},
		{
			name: "query parameter",
			setup: func(r *http.Request) {
				q := r.URL.Query()
				q.Set("api_key", "query-token")
				r.URL.RawQuery = q.Encode()
			},
			wantToken: "query-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setup(req)

			got := extractToken(req)
			if got != tt.wantToken {
				t.Errorf("extractToken() = %q, want %q", got, tt.wantToken)
			}
		})
	}
}

func TestParseMediaBrowserAuth(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   map[string]string
	}{
		{
			name:   "empty header",
			header: "",
			want:   map[string]string{},
		},
		{
			name:   "full header",
			header: `MediaBrowser Client="Revenge Web", Device="Chrome", DeviceId="abc123", Version="10.8.0", Token="xyz789"`,
			want: map[string]string{
				"Client":   "Revenge Web",
				"Device":   "Chrome",
				"DeviceId": "abc123",
				"Version":  "10.8.0",
				"Token":    "xyz789",
			},
		},
		{
			name:   "token only",
			header: `MediaBrowser Token="only-token"`,
			want: map[string]string{
				"Token": "only-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMediaBrowserAuth(tt.header)

			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("parseMediaBrowserAuth()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
