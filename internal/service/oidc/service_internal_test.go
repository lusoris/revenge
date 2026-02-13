package oidc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/lusoris/revenge/internal/infra/logging"
)

// mockRepo implements Repository for same-package unit tests.
type mockRepo struct {
	// Provider methods
	createProviderFn           func(ctx context.Context, req CreateProviderRequest) (*Provider, error)
	getProviderFn              func(ctx context.Context, id uuid.UUID) (*Provider, error)
	getProviderByNameFn        func(ctx context.Context, name string) (*Provider, error)
	getDefaultProviderFn       func(ctx context.Context) (*Provider, error)
	listProvidersFn            func(ctx context.Context) ([]Provider, error)
	listEnabledProvidersFn     func(ctx context.Context) ([]Provider, error)
	updateProviderFn           func(ctx context.Context, id uuid.UUID, req UpdateProviderRequest) (*Provider, error)
	deleteProviderFn           func(ctx context.Context, id uuid.UUID) error
	enableProviderFn           func(ctx context.Context, id uuid.UUID) error
	disableProviderFn          func(ctx context.Context, id uuid.UUID) error
	setDefaultProviderFn       func(ctx context.Context, id uuid.UUID) error
	createUserLinkFn           func(ctx context.Context, req CreateUserLinkRequest) (*UserLink, error)
	getUserLinkFn              func(ctx context.Context, id uuid.UUID) (*UserLink, error)
	getUserLinkBySubjectFn     func(ctx context.Context, providerID uuid.UUID, subject string) (*UserLink, error)
	getUserLinkByUserProvFn    func(ctx context.Context, userID, providerID uuid.UUID) (*UserLink, error)
	listUserLinksFn            func(ctx context.Context, userID uuid.UUID) ([]UserLinkWithProvider, error)
	updateUserLinkFn           func(ctx context.Context, id uuid.UUID, req UpdateUserLinkRequest) (*UserLink, error)
	updateUserLinkLastLoginFn  func(ctx context.Context, id uuid.UUID) error
	deleteUserLinkFn           func(ctx context.Context, id uuid.UUID) error
	deleteUserLinkByUserProvFn func(ctx context.Context, userID, providerID uuid.UUID) error
	countUserLinksFn           func(ctx context.Context, userID uuid.UUID) (int64, error)
	createStateFn              func(ctx context.Context, req CreateStateRequest) (*State, error)
	getStateFn                 func(ctx context.Context, state string) (*State, error)
	deleteStateFn              func(ctx context.Context, state string) error
	deleteExpiredStatesFn      func(ctx context.Context) (int64, error)
	deleteStatesByProviderFn   func(ctx context.Context, providerID uuid.UUID) error
}

func (m *mockRepo) CreateProvider(ctx context.Context, req CreateProviderRequest) (*Provider, error) {
	if m.createProviderFn != nil {
		return m.createProviderFn(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetProvider(ctx context.Context, id uuid.UUID) (*Provider, error) {
	if m.getProviderFn != nil {
		return m.getProviderFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetProviderByName(ctx context.Context, name string) (*Provider, error) {
	if m.getProviderByNameFn != nil {
		return m.getProviderByNameFn(ctx, name)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetDefaultProvider(ctx context.Context) (*Provider, error) {
	if m.getDefaultProviderFn != nil {
		return m.getDefaultProviderFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) ListProviders(ctx context.Context) ([]Provider, error) {
	if m.listProvidersFn != nil {
		return m.listProvidersFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) ListEnabledProviders(ctx context.Context) ([]Provider, error) {
	if m.listEnabledProvidersFn != nil {
		return m.listEnabledProvidersFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) UpdateProvider(ctx context.Context, id uuid.UUID, req UpdateProviderRequest) (*Provider, error) {
	if m.updateProviderFn != nil {
		return m.updateProviderFn(ctx, id, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	if m.deleteProviderFn != nil {
		return m.deleteProviderFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) EnableProvider(ctx context.Context, id uuid.UUID) error {
	if m.enableProviderFn != nil {
		return m.enableProviderFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) DisableProvider(ctx context.Context, id uuid.UUID) error {
	if m.disableProviderFn != nil {
		return m.disableProviderFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) SetDefaultProvider(ctx context.Context, id uuid.UUID) error {
	if m.setDefaultProviderFn != nil {
		return m.setDefaultProviderFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) CreateUserLink(ctx context.Context, req CreateUserLinkRequest) (*UserLink, error) {
	if m.createUserLinkFn != nil {
		return m.createUserLinkFn(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetUserLink(ctx context.Context, id uuid.UUID) (*UserLink, error) {
	if m.getUserLinkFn != nil {
		return m.getUserLinkFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetUserLinkBySubject(ctx context.Context, providerID uuid.UUID, subject string) (*UserLink, error) {
	if m.getUserLinkBySubjectFn != nil {
		return m.getUserLinkBySubjectFn(ctx, providerID, subject)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetUserLinkByUserAndProvider(ctx context.Context, userID, providerID uuid.UUID) (*UserLink, error) {
	if m.getUserLinkByUserProvFn != nil {
		return m.getUserLinkByUserProvFn(ctx, userID, providerID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) ListUserLinks(ctx context.Context, userID uuid.UUID) ([]UserLinkWithProvider, error) {
	if m.listUserLinksFn != nil {
		return m.listUserLinksFn(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) UpdateUserLink(ctx context.Context, id uuid.UUID, req UpdateUserLinkRequest) (*UserLink, error) {
	if m.updateUserLinkFn != nil {
		return m.updateUserLinkFn(ctx, id, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) UpdateUserLinkLastLogin(ctx context.Context, id uuid.UUID) error {
	if m.updateUserLinkLastLoginFn != nil {
		return m.updateUserLinkLastLoginFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) DeleteUserLink(ctx context.Context, id uuid.UUID) error {
	if m.deleteUserLinkFn != nil {
		return m.deleteUserLinkFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) DeleteUserLinkByUserAndProvider(ctx context.Context, userID, providerID uuid.UUID) error {
	if m.deleteUserLinkByUserProvFn != nil {
		return m.deleteUserLinkByUserProvFn(ctx, userID, providerID)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) CountUserLinks(ctx context.Context, userID uuid.UUID) (int64, error) {
	if m.countUserLinksFn != nil {
		return m.countUserLinksFn(ctx, userID)
	}
	return 0, errors.New("not implemented")
}

func (m *mockRepo) CreateState(ctx context.Context, req CreateStateRequest) (*State, error) {
	if m.createStateFn != nil {
		return m.createStateFn(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetState(ctx context.Context, state string) (*State, error) {
	if m.getStateFn != nil {
		return m.getStateFn(ctx, state)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) DeleteState(ctx context.Context, state string) error {
	if m.deleteStateFn != nil {
		return m.deleteStateFn(ctx, state)
	}
	return errors.New("not implemented")
}

func (m *mockRepo) DeleteExpiredStates(ctx context.Context) (int64, error) {
	if m.deleteExpiredStatesFn != nil {
		return m.deleteExpiredStatesFn(ctx)
	}
	return 0, errors.New("not implemented")
}

func (m *mockRepo) DeleteStatesByProvider(ctx context.Context, providerID uuid.UUID) error {
	if m.deleteStatesByProviderFn != nil {
		return m.deleteStatesByProviderFn(ctx, providerID)
	}
	return errors.New("not implemented")
}

// newTestService creates a service for unit tests using the manual mock.
func newTestService(repo *mockRepo) *Service {
	return NewService(repo, logging.NewTestLogger(), "http://localhost:8080/callback", nil)
}

func newTestServiceWithKey(repo *mockRepo, key []byte) *Service {
	return NewService(repo, logging.NewTestLogger(), "http://localhost:8080/callback", key)
}

// ============================================================================
// getClaim Tests
// ============================================================================

func Test_getClaim(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		claims   map[string]any
		path     string
		expected any
	}{
		{
			name:     "simple key",
			claims:   map[string]any{"email": "test@example.com"},
			path:     "email",
			expected: "test@example.com",
		},
		{
			name: "nested key",
			claims: map[string]any{
				"realm_access": map[string]any{
					"roles": []string{"admin", "user"},
				},
			},
			path:     "realm_access.roles",
			expected: []string{"admin", "user"},
		},
		{
			name: "deeply nested key",
			claims: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": "deep_value",
					},
				},
			},
			path:     "a.b.c",
			expected: "deep_value",
		},
		{
			name:     "missing key",
			claims:   map[string]any{"email": "test@example.com"},
			path:     "name",
			expected: nil,
		},
		{
			name: "missing nested key",
			claims: map[string]any{
				"a": map[string]any{
					"b": "value",
				},
			},
			path:     "a.c",
			expected: nil,
		},
		{
			name:     "non-map in path",
			claims:   map[string]any{"email": "test@example.com"},
			path:     "email.sub",
			expected: nil,
		},
		{
			name:     "empty path",
			claims:   map[string]any{"": "empty_key_value"},
			path:     "",
			expected: "empty_key_value",
		},
		{
			name:     "numeric value",
			claims:   map[string]any{"count": float64(42)},
			path:     "count",
			expected: float64(42),
		},
		{
			name:     "boolean value",
			claims:   map[string]any{"active": true},
			path:     "active",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := getClaim(tt.claims, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// extractUserInfo Tests (internal)
// ============================================================================

func Test_extractUserInfo(t *testing.T) {
	t.Parallel()

	svc := newTestService(&mockRepo{})

	t.Run("all claims present", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Username: "preferred_username",
				Email:    "email",
				Name:     "name",
				Picture:  "picture",
			},
		}

		claims := map[string]any{
			"preferred_username": "johndoe",
			"email":              "john@example.com",
			"name":               "John Doe",
			"picture":            "https://example.com/photo.jpg",
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Equal(t, "johndoe", info.Username)
		assert.Equal(t, "john@example.com", info.Email)
		assert.Equal(t, "John Doe", info.Name)
		assert.Equal(t, "https://example.com/photo.jpg", info.Picture)
		assert.Equal(t, claims, info.Claims)
	})

	t.Run("missing claims return empty strings", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Username: "preferred_username",
				Email:    "email",
				Name:     "name",
				Picture:  "picture",
			},
		}

		claims := map[string]any{}

		info := svc.extractUserInfo(provider, claims)

		assert.Empty(t, info.Username)
		assert.Empty(t, info.Email)
		assert.Empty(t, info.Name)
		assert.Empty(t, info.Picture)
	})

	t.Run("roles from []any", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Roles: "roles",
			},
		}

		claims := map[string]any{
			"roles": []any{"admin", "user", 123}, // 123 is not a string, should be skipped
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Equal(t, []string{"admin", "user"}, info.Roles)
	})

	t.Run("roles from []string", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Roles: "realm_access.roles",
			},
		}

		claims := map[string]any{
			"realm_access": map[string]any{
				"roles": []string{"admin", "editor"},
			},
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Equal(t, []string{"admin", "editor"}, info.Roles)
	})

	t.Run("role mappings", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Roles: "roles",
			},
			RoleMappings: map[string]string{
				"oidc_admin":  "admin",
				"oidc_editor": "editor",
			},
		}

		claims := map[string]any{
			"roles": []any{"oidc_admin", "oidc_editor", "unknown_role"},
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Equal(t, []string{"admin", "editor"}, info.Roles)
	})

	t.Run("role mappings with no matching roles", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Roles: "roles",
			},
			RoleMappings: map[string]string{
				"oidc_admin": "admin",
			},
		}

		claims := map[string]any{
			"roles": []any{"unknown_role"},
		}

		info := svc.extractUserInfo(provider, claims)

		// No roles should match
		assert.Empty(t, info.Roles)
	})

	t.Run("empty roles claim mapping skips role extraction", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Roles: "", // Empty - no role extraction
			},
		}

		claims := map[string]any{
			"roles": []any{"admin"},
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Nil(t, info.Roles)
	})

	t.Run("non-string claim values are skipped", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Username: "preferred_username",
				Email:    "email",
				Name:     "name",
				Picture:  "picture",
			},
		}

		claims := map[string]any{
			"preferred_username": 12345,   // Not a string
			"email":              true,    // Not a string
			"name":               nil,     // nil
			"picture":            []int{}, // Not a string
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Empty(t, info.Username)
		assert.Empty(t, info.Email)
		assert.Empty(t, info.Name)
		assert.Empty(t, info.Picture)
	})

	t.Run("nested claim paths", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Email: "profile.email",
				Name:  "profile.display_name",
			},
		}

		claims := map[string]any{
			"profile": map[string]any{
				"email":        "nested@example.com",
				"display_name": "Nested User",
			},
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Equal(t, "nested@example.com", info.Email)
		assert.Equal(t, "Nested User", info.Name)
	})

	t.Run("roles claim is not an array", func(t *testing.T) {
		t.Parallel()
		provider := &Provider{
			ClaimMappings: ClaimMappings{
				Roles: "roles",
			},
		}

		claims := map[string]any{
			"roles": "single_role", // string, not array - should not be parsed as roles
		}

		info := svc.extractUserInfo(provider, claims)

		assert.Nil(t, info.Roles)
	})
}

// ============================================================================
// encryptSecret / decryptSecret Tests (internal)
// ============================================================================

func Test_encryptDecryptSecret(t *testing.T) {
	t.Parallel()

	// AES-256 requires 32-byte key
	key := []byte("01234567890123456789012345678901")

	t.Run("round trip", func(t *testing.T) {
		t.Parallel()
		svc := newTestServiceWithKey(&mockRepo{}, key)

		plaintext := []byte("super-secret-value")

		encrypted, err := svc.encryptSecret(plaintext)
		require.NoError(t, err)
		assert.NotEqual(t, plaintext, encrypted)
		assert.Greater(t, len(encrypted), len(plaintext)) // GCM adds nonce + tag

		decrypted := svc.decryptSecret(encrypted)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("no encryption key returns plaintext", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{}) // nil key

		plaintext := []byte("unencrypted-value")

		encrypted, err := svc.encryptSecret(plaintext)
		require.NoError(t, err)
		assert.Equal(t, plaintext, encrypted)
	})

	t.Run("no encryption key decrypts as passthrough", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{}) // nil key

		ciphertext := []byte("some-value")
		decrypted := svc.decryptSecret(ciphertext)
		assert.Equal(t, ciphertext, decrypted)
	})

	t.Run("decrypt with invalid key returns ciphertext fallback", func(t *testing.T) {
		t.Parallel()
		svc := newTestServiceWithKey(&mockRepo{}, key)

		// Encrypt with one key
		plaintext := []byte("secret-data")
		encrypted, err := svc.encryptSecret(plaintext)
		require.NoError(t, err)

		// Try to decrypt with a different key
		differentKey := []byte("different_key_different_key_1234")
		svcBadKey := newTestServiceWithKey(&mockRepo{}, differentKey)
		decrypted := svcBadKey.decryptSecret(encrypted)

		// Should return ciphertext as fallback (minus nonce) or ciphertext
		// The important thing is it doesn't panic
		assert.NotNil(t, decrypted)
	})

	t.Run("decrypt with short ciphertext returns ciphertext", func(t *testing.T) {
		t.Parallel()
		svc := newTestServiceWithKey(&mockRepo{}, key)

		// Ciphertext shorter than nonce size
		shortData := []byte("abc")
		decrypted := svc.decryptSecret(shortData)
		assert.Equal(t, shortData, decrypted)
	})

	t.Run("encrypt with invalid key length returns error", func(t *testing.T) {
		t.Parallel()
		badKey := []byte("short") // Not 16, 24, or 32 bytes
		svc := newTestServiceWithKey(&mockRepo{}, badKey)

		_, err := svc.encryptSecret([]byte("test"))
		require.Error(t, err)
	})

	t.Run("decrypt with invalid key length returns ciphertext", func(t *testing.T) {
		t.Parallel()
		badKey := []byte("short")
		svc := newTestServiceWithKey(&mockRepo{}, badKey)

		data := []byte("some-ciphertext-that-wont-decrypt")
		decrypted := svc.decryptSecret(data)
		assert.Equal(t, data, decrypted) // Falls back to returning as-is
	})

	t.Run("encrypt produces different ciphertexts for same plaintext", func(t *testing.T) {
		t.Parallel()
		svc := newTestServiceWithKey(&mockRepo{}, key)

		plaintext := []byte("same-data")

		enc1, err := svc.encryptSecret(plaintext)
		require.NoError(t, err)

		enc2, err := svc.encryptSecret(plaintext)
		require.NoError(t, err)

		// Due to random nonce, ciphertexts should differ
		assert.NotEqual(t, enc1, enc2)
	})

	t.Run("empty plaintext encrypts and decrypts", func(t *testing.T) {
		t.Parallel()
		svc := newTestServiceWithKey(&mockRepo{}, key)

		plaintext := []byte("")

		encrypted, err := svc.encryptSecret(plaintext)
		require.NoError(t, err)

		decrypted := svc.decryptSecret(encrypted)
		assert.Empty(t, decrypted)
	})
}

// ============================================================================
// buildOAuth2Config Tests (internal)
// ============================================================================

func Test_buildOAuth2Config(t *testing.T) {
	t.Parallel()

	t.Run("with custom endpoints", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})

		authURL := "https://auth.example.com/authorize"
		tokenURL := "https://auth.example.com/token"
		provider := &Provider{
			Name:                  "test-provider",
			ClientID:              "client-123",
			ClientSecretEncrypted: []byte("secret-456"),
			Scopes:                []string{"openid", "profile"},
			AuthorizationEndpoint: &authURL,
			TokenEndpoint:         &tokenURL,
		}

		config := svc.buildOAuth2Config(provider, nil)

		assert.Equal(t, "client-123", config.ClientID)
		assert.Equal(t, "secret-456", config.ClientSecret) // No encryption key, passthrough
		assert.Equal(t, authURL, config.Endpoint.AuthURL)
		assert.Equal(t, tokenURL, config.Endpoint.TokenURL)
		assert.Equal(t, []string{"openid", "profile"}, config.Scopes)
		assert.Contains(t, config.RedirectURL, "test-provider")
	})

	t.Run("without custom endpoints uses issuer URL defaults", func(t *testing.T) {
		t.Parallel()
		svc := newTestService(&mockRepo{})

		provider := &Provider{
			Name:                  "generic-provider",
			ClientID:              "client-abc",
			ClientSecretEncrypted: []byte("secret"),
			IssuerURL:             "https://issuer.example.com",
			Scopes:                []string{"openid"},
		}

		config := svc.buildOAuth2Config(provider, nil)

		assert.Equal(t, "https://issuer.example.com/authorize", config.Endpoint.AuthURL)
		assert.Equal(t, "https://issuer.example.com/token", config.Endpoint.TokenURL)
	})

	t.Run("callback URL with existing /callback/ path", func(t *testing.T) {
		svc := NewService(&mockRepo{}, logging.NewTestLogger(), "http://localhost:8080/callback/existing", nil)

		provider := &Provider{
			Name:                  "test",
			ClientID:              "client",
			ClientSecretEncrypted: []byte("secret"),
			Scopes:                []string{"openid"},
		}

		config := svc.buildOAuth2Config(provider, nil)

		// Should NOT append provider name since URL already contains /callback/
		assert.Equal(t, "http://localhost:8080/callback/existing", config.RedirectURL)
	})

	t.Run("callback URL without /callback/ appends provider name", func(t *testing.T) {
		svc := NewService(&mockRepo{}, logging.NewTestLogger(), "http://localhost:8080/auth", nil)

		provider := &Provider{
			Name:                  "keycloak",
			ClientID:              "client",
			ClientSecretEncrypted: []byte("secret"),
			Scopes:                []string{"openid"},
		}

		config := svc.buildOAuth2Config(provider, nil)

		assert.Equal(t, "http://localhost:8080/auth/keycloak", config.RedirectURL)
	})

	t.Run("callback URL with trailing slash", func(t *testing.T) {
		svc := NewService(&mockRepo{}, logging.NewTestLogger(), "http://localhost:8080/auth/", nil)

		provider := &Provider{
			Name:                  "authentik",
			ClientID:              "client",
			ClientSecretEncrypted: []byte("secret"),
			Scopes:                []string{"openid"},
		}

		config := svc.buildOAuth2Config(provider, nil)

		assert.Equal(t, "http://localhost:8080/auth/authentik", config.RedirectURL)
	})
}

// ============================================================================
// createUserLink Tests (internal)
// ============================================================================

func Test_createUserLink(t *testing.T) {
	t.Parallel()

	t.Run("success with token", func(t *testing.T) {
		t.Parallel()

		linkID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		providerID := uuid.Must(uuid.NewV7())

		repo := &mockRepo{
			createUserLinkFn: func(_ context.Context, req CreateUserLinkRequest) (*UserLink, error) {
				assert.Equal(t, userID, req.UserID)
				assert.Equal(t, providerID, req.ProviderID)
				assert.Equal(t, "subject-123", req.Subject)
				assert.NotNil(t, req.TokenExpiresAt)
				return &UserLink{
					ID:         linkID,
					UserID:     userID,
					ProviderID: providerID,
					Subject:    "subject-123",
				}, nil
			},
		}

		svc := newTestService(repo)

		expiry := time.Now().Add(time.Hour)
		token := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       expiry,
		}

		provider := &Provider{ID: providerID}
		userInfo := &UserInfo{
			Subject: "subject-123",
			Email:   "test@example.com",
			Name:    "Test User",
			Picture: "https://example.com/pic.jpg",
		}

		link, err := svc.createUserLink(context.Background(), provider, userID, userInfo, token)

		require.NoError(t, err)
		assert.Equal(t, linkID, link.ID)
	})

	t.Run("success without token", func(t *testing.T) {
		t.Parallel()

		linkID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		providerID := uuid.Must(uuid.NewV7())

		repo := &mockRepo{
			createUserLinkFn: func(_ context.Context, req CreateUserLinkRequest) (*UserLink, error) {
				// No token fields should be set
				assert.Nil(t, req.AccessTokenEncrypted)
				assert.Nil(t, req.RefreshTokenEncrypted)
				assert.Nil(t, req.TokenExpiresAt)
				return &UserLink{
					ID:         linkID,
					UserID:     userID,
					ProviderID: providerID,
					Subject:    "subject-456",
				}, nil
			},
		}

		svc := newTestService(repo)

		provider := &Provider{ID: providerID}
		userInfo := &UserInfo{
			Subject: "subject-456",
			Email:   "user@example.com",
			Name:    "User",
		}

		link, err := svc.createUserLink(context.Background(), provider, userID, userInfo, nil)

		require.NoError(t, err)
		assert.Equal(t, linkID, link.ID)
	})

	t.Run("success with token zero expiry", func(t *testing.T) {
		t.Parallel()

		userID := uuid.Must(uuid.NewV7())
		providerID := uuid.Must(uuid.NewV7())

		repo := &mockRepo{
			createUserLinkFn: func(_ context.Context, req CreateUserLinkRequest) (*UserLink, error) {
				// Zero expiry should not set TokenExpiresAt
				assert.Nil(t, req.TokenExpiresAt)
				return &UserLink{
					ID:         uuid.Must(uuid.NewV7()),
					UserID:     userID,
					ProviderID: providerID,
				}, nil
			},
		}

		svc := newTestService(repo)

		token := &oauth2.Token{
			AccessToken:  "at",
			RefreshToken: "rt",
			// Expiry is zero
		}

		provider := &Provider{ID: providerID}
		userInfo := &UserInfo{Subject: "sub"}

		_, err := svc.createUserLink(context.Background(), provider, userID, userInfo, token)
		require.NoError(t, err)
	})

	t.Run("success with token empty access and refresh", func(t *testing.T) {
		t.Parallel()

		userID := uuid.Must(uuid.NewV7())
		providerID := uuid.Must(uuid.NewV7())

		repo := &mockRepo{
			createUserLinkFn: func(_ context.Context, req CreateUserLinkRequest) (*UserLink, error) {
				assert.Nil(t, req.AccessTokenEncrypted)
				assert.Nil(t, req.RefreshTokenEncrypted)
				return &UserLink{
					ID:         uuid.Must(uuid.NewV7()),
					UserID:     userID,
					ProviderID: providerID,
				}, nil
			},
		}

		svc := newTestService(repo)

		token := &oauth2.Token{
			// AccessToken and RefreshToken are empty strings
		}

		provider := &Provider{ID: providerID}
		userInfo := &UserInfo{Subject: "sub"}

		_, err := svc.createUserLink(context.Background(), provider, userID, userInfo, token)
		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepo{
			createUserLinkFn: func(_ context.Context, _ CreateUserLinkRequest) (*UserLink, error) {
				return nil, errors.New("database error")
			},
		}

		svc := newTestService(repo)

		provider := &Provider{ID: uuid.Must(uuid.NewV7())}
		userInfo := &UserInfo{Subject: "sub"}

		_, err := svc.createUserLink(context.Background(), provider, uuid.Must(uuid.NewV7()), userInfo, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

// ============================================================================
// HandleCallback Tests (internal, without external OIDC discovery)
// ============================================================================

func Test_HandleCallback_StateLookup(t *testing.T) {
	t.Parallel()

	t.Run("state not found returns ErrInvalidState", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepo{
			getStateFn: func(_ context.Context, _ string) (*State, error) {
				return nil, ErrStateNotFound
			},
		}

		svc := newTestService(repo)

		_, err := svc.HandleCallback(context.Background(), "bad-state", "code")
		require.ErrorIs(t, err, ErrInvalidState)
	})

	t.Run("expired state returns ErrStateExpired", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			getStateFn: func(_ context.Context, _ string) (*State, error) {
				return &State{
					ID:         uuid.Must(uuid.NewV7()),
					State:      "expired-state",
					ProviderID: providerID,
					ExpiresAt:  time.Now().Add(-1 * time.Hour), // Expired
				}, nil
			},
			deleteStateFn: func(_ context.Context, _ string) error {
				return nil
			},
		}

		svc := newTestService(repo)

		_, err := svc.HandleCallback(context.Background(), "expired-state", "code")
		require.ErrorIs(t, err, ErrStateExpired)
	})

	t.Run("provider not found after state lookup", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			getStateFn: func(_ context.Context, _ string) (*State, error) {
				return &State{
					ID:         uuid.Must(uuid.NewV7()),
					State:      "valid-state",
					ProviderID: providerID,
					ExpiresAt:  time.Now().Add(10 * time.Minute),
				}, nil
			},
			deleteStateFn: func(_ context.Context, _ string) error {
				return nil
			},
			getProviderFn: func(_ context.Context, _ uuid.UUID) (*Provider, error) {
				return nil, ErrProviderNotFound
			},
		}

		svc := newTestService(repo)

		_, err := svc.HandleCallback(context.Background(), "valid-state", "code")
		require.ErrorIs(t, err, ErrProviderNotFound)
	})

	t.Run("disabled provider after state lookup", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			getStateFn: func(_ context.Context, _ string) (*State, error) {
				return &State{
					ID:         uuid.Must(uuid.NewV7()),
					State:      "valid-state",
					ProviderID: providerID,
					ExpiresAt:  time.Now().Add(10 * time.Minute),
				}, nil
			},
			deleteStateFn: func(_ context.Context, _ string) error {
				return nil
			},
			getProviderFn: func(_ context.Context, _ uuid.UUID) (*Provider, error) {
				return &Provider{
					ID:        providerID,
					IsEnabled: false, // Disabled
				}, nil
			},
		}

		svc := newTestService(repo)

		_, err := svc.HandleCallback(context.Background(), "valid-state", "code")
		require.ErrorIs(t, err, ErrProviderDisabled)
	})

	t.Run("OIDC discovery failure", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			getStateFn: func(_ context.Context, _ string) (*State, error) {
				return &State{
					ID:         uuid.Must(uuid.NewV7()),
					State:      "valid-state",
					ProviderID: providerID,
					ExpiresAt:  time.Now().Add(10 * time.Minute),
				}, nil
			},
			deleteStateFn: func(_ context.Context, _ string) error {
				return nil
			},
			getProviderFn: func(_ context.Context, _ uuid.UUID) (*Provider, error) {
				return &Provider{
					ID:        providerID,
					IssuerURL: "https://invalid-issuer-that-will-fail.example.com",
					IsEnabled: true,
				}, nil
			},
		}

		svc := newTestService(repo)

		_, err := svc.HandleCallback(context.Background(), "valid-state", "code")
		// Should fail at OIDC discovery
		require.ErrorIs(t, err, ErrDiscoveryFailed)
	})
}

// ============================================================================
// AddProvider validation Tests (internal)
// ============================================================================

func Test_AddProvider_Validation(t *testing.T) {
	t.Parallel()

	t.Run("preserves custom scopes", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			getProviderByNameFn: func(_ context.Context, _ string) (*Provider, error) {
				return nil, ErrProviderNotFound
			},
			createProviderFn: func(_ context.Context, req CreateProviderRequest) (*Provider, error) {
				assert.Equal(t, []string{"openid", "custom_scope"}, req.Scopes)
				return &Provider{
					ID:     providerID,
					Name:   req.Name,
					Scopes: req.Scopes,
				}, nil
			},
		}

		svc := newTestService(repo)

		req := CreateProviderRequest{
			Name:                  "test",
			ProviderType:          ProviderTypeGeneric,
			IssuerURL:             "https://issuer.example.com",
			ClientID:              "client",
			ClientSecretEncrypted: []byte("secret"),
			Scopes:                []string{"openid", "custom_scope"},
		}

		provider, err := svc.AddProvider(context.Background(), req)
		require.NoError(t, err)
		assert.NotNil(t, provider)
	})

	t.Run("preserves custom claim mappings", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepo{
			getProviderByNameFn: func(_ context.Context, _ string) (*Provider, error) {
				return nil, ErrProviderNotFound
			},
			createProviderFn: func(_ context.Context, req CreateProviderRequest) (*Provider, error) {
				assert.Equal(t, "sub", req.ClaimMappings.Username)
				assert.Equal(t, "mail", req.ClaimMappings.Email)
				assert.Equal(t, "name", req.ClaimMappings.Name)       // default
				assert.Equal(t, "picture", req.ClaimMappings.Picture) // default
				return &Provider{ID: uuid.Must(uuid.NewV7()), Name: req.Name}, nil
			},
		}

		svc := newTestService(repo)

		req := CreateProviderRequest{
			Name:                  "test",
			ProviderType:          ProviderTypeGeneric,
			IssuerURL:             "https://issuer.example.com",
			ClientID:              "client",
			ClientSecretEncrypted: []byte("secret"),
			ClaimMappings: ClaimMappings{
				Username: "sub",
				Email:    "mail",
				// Name and Picture are empty - should be filled with defaults
			},
		}

		_, err := svc.AddProvider(context.Background(), req)
		require.NoError(t, err)
	})

	t.Run("CreateProvider repo error", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepo{
			getProviderByNameFn: func(_ context.Context, _ string) (*Provider, error) {
				return nil, ErrProviderNotFound
			},
			createProviderFn: func(_ context.Context, _ CreateProviderRequest) (*Provider, error) {
				return nil, errors.New("insert failed")
			},
		}

		svc := newTestService(repo)

		req := CreateProviderRequest{
			Name:                  "test",
			ProviderType:          ProviderTypeGeneric,
			IssuerURL:             "https://issuer.example.com",
			ClientID:              "client",
			ClientSecretEncrypted: []byte("secret"),
		}

		_, err := svc.AddProvider(context.Background(), req)
		require.Error(t, err)
	})
}

// ============================================================================
// UpdateProvider encryption Tests (internal)
// ============================================================================

func Test_UpdateProvider_Encryption(t *testing.T) {
	t.Parallel()

	key := []byte("01234567890123456789012345678901")

	t.Run("encrypts client secret when provided", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			updateProviderFn: func(_ context.Context, _ uuid.UUID, req UpdateProviderRequest) (*Provider, error) {
				// The secret should be encrypted (different from original)
				assert.NotEqual(t, []byte("new-secret"), req.ClientSecretEncrypted)
				return &Provider{ID: providerID, Name: "test"}, nil
			},
		}

		svc := newTestServiceWithKey(repo, key)

		req := UpdateProviderRequest{
			ClientSecretEncrypted: []byte("new-secret"),
		}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)
		require.NoError(t, err)
		assert.NotNil(t, provider)
	})

	t.Run("no encryption when secret not provided", func(t *testing.T) {
		t.Parallel()

		providerID := uuid.Must(uuid.NewV7())
		repo := &mockRepo{
			updateProviderFn: func(_ context.Context, _ uuid.UUID, req UpdateProviderRequest) (*Provider, error) {
				assert.Nil(t, req.ClientSecretEncrypted)
				return &Provider{ID: providerID}, nil
			},
		}

		svc := newTestServiceWithKey(repo, key)

		newDisplayName := "Updated"
		req := UpdateProviderRequest{
			DisplayName: &newDisplayName,
		}

		_, err := svc.UpdateProvider(context.Background(), providerID, req)
		require.NoError(t, err)
	})
}

// ============================================================================
// generateRandomString Tests
// ============================================================================

func Test_generateRandomString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		length       int
		expectHexLen int
	}{
		{"16 bytes", 16, 32},
		{"32 bytes", 32, 64},
		{"64 bytes", 64, 128},
		{"1 byte", 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := generateRandomString(tt.length)
			require.NoError(t, err)
			assert.Len(t, s, tt.expectHexLen)
		})
	}

	t.Run("different calls produce different values", func(t *testing.T) {
		t.Parallel()
		s1, err := generateRandomString(32)
		require.NoError(t, err)
		s2, err := generateRandomString(32)
		require.NoError(t, err)
		assert.NotEqual(t, s1, s2)
	})
}

// ============================================================================
// generateCodeChallenge Tests
// ============================================================================

func Test_generateCodeChallenge(t *testing.T) {
	t.Parallel()

	t.Run("deterministic for same input", func(t *testing.T) {
		t.Parallel()
		c1 := generateCodeChallenge("test-verifier")
		c2 := generateCodeChallenge("test-verifier")
		assert.Equal(t, c1, c2)
	})

	t.Run("different for different input", func(t *testing.T) {
		t.Parallel()
		c1 := generateCodeChallenge("verifier-1")
		c2 := generateCodeChallenge("verifier-2")
		assert.NotEqual(t, c1, c2)
	})

	t.Run("not empty", func(t *testing.T) {
		t.Parallel()
		c := generateCodeChallenge("any-verifier")
		assert.NotEmpty(t, c)
	})
}

// ============================================================================
// isValidProviderType Tests
// ============================================================================

func Test_isValidProviderType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		providerType string
		valid        bool
	}{
		{"oidc", true},
		{"generic", true},
		{"authentik", true},
		{"keycloak", true},
		{"invalid", false},
		{"", false},
		{"OIDC", false},    // Case-sensitive
		{"Generic", false}, // Case-sensitive
	}

	for _, tt := range tests {
		t.Run(tt.providerType, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.valid, isValidProviderType(tt.providerType))
		})
	}
}

// ============================================================================
// GetAuthURL additional edge cases (internal)
// ============================================================================

func Test_GetAuthURL_WithUserID(t *testing.T) {
	t.Parallel()

	providerID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	repo := &mockRepo{
		getProviderByNameFn: func(_ context.Context, name string) (*Provider, error) {
			return &Provider{
				ID:        providerID,
				Name:      name,
				IsEnabled: true,
				ClientID:  "client-id",
				IssuerURL: "https://issuer.example.com",
				Scopes:    []string{"openid"},
			}, nil
		},
		createStateFn: func(_ context.Context, req CreateStateRequest) (*State, error) {
			assert.NotNil(t, req.UserID)
			assert.Equal(t, userID, *req.UserID)
			assert.NotEmpty(t, req.State)
			assert.NotNil(t, req.CodeVerifier)
			return &State{
				ID:         uuid.Must(uuid.NewV7()),
				State:      req.State,
				ProviderID: providerID,
				UserID:     req.UserID,
			}, nil
		},
	}

	svc := newTestService(repo)

	result, err := svc.GetAuthURL(context.Background(), "test-provider", "https://redirect.example.com", &userID)
	require.NoError(t, err)
	assert.NotEmpty(t, result.URL)
	assert.NotEmpty(t, result.State)
	assert.Contains(t, result.URL, "code_challenge")
	assert.Contains(t, result.URL, "S256")
}
