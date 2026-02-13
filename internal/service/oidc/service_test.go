package oidc

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupTestService(t *testing.T) (*Service, *RepositoryPg, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)

	logger := logging.NewTestLogger()
	encryptKey := make([]byte, 32) // 32 bytes for AES-256
	for i := range encryptKey {
		encryptKey[i] = byte(i)
	}

	service := NewService(repo, logger, "http://localhost:8080/callback", encryptKey)
	return service, repo, testDB
}

// ============================================================================
// Provider Management Tests
// ============================================================================

func TestService_AddProvider(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("add_svc")
	req.ClientSecretEncrypted = []byte("plaintext-secret")

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, provider.ID)
	assert.Equal(t, "add_svc", provider.Name)

	// Verify secret was encrypted (should not match plaintext)
	assert.NotEqual(t, []byte("plaintext-secret"), provider.ClientSecretEncrypted)
	assert.NotEmpty(t, provider.ClientSecretEncrypted)

	// Verify we can decrypt it back
	decrypted := svc.decryptSecret(provider.ClientSecretEncrypted)
	assert.Equal(t, []byte("plaintext-secret"), decrypted)

	// Verify default scopes
	assert.Contains(t, provider.Scopes, "openid")

	// Verify default claim mappings
	assert.Equal(t, "preferred_username", provider.ClaimMappings.Username)
	assert.Equal(t, "email", provider.ClaimMappings.Email)
	assert.Equal(t, "name", provider.ClaimMappings.Name)
}

func TestService_AddProvider_InvalidProviderType(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("invalid_type")
	req.ProviderType = "invalid-type"

	_, err := svc.AddProvider(ctx, req)
	assert.ErrorIs(t, err, ErrInvalidProviderType)
}

func TestService_AddProvider_InvalidIssuerURL(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("invalid_url")
	req.IssuerURL = "not a valid url ::"

	_, err := svc.AddProvider(ctx, req)
	assert.ErrorIs(t, err, ErrInvalidIssuerURL)
}

func TestService_AddProvider_DuplicateName(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("duplicate")
	req.ClientSecretEncrypted = []byte("secret")

	_, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Try to create again
	_, err = svc.AddProvider(ctx, req)
	assert.ErrorIs(t, err, ErrProviderNameExists)
}

func TestService_GetProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("get_svc")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	provider, err := svc.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, provider.ID)
}

func TestService_GetProviderByName(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("by_name_svc")
	_, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	provider, err := svc.GetProviderByName(ctx, "by_name_svc")
	require.NoError(t, err)
	assert.Equal(t, "by_name_svc", provider.Name)
}

func TestService_GetDefaultProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("default_svc")
	req.IsDefault = true
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	provider, err := svc.GetDefaultProvider(ctx)
	require.NoError(t, err)
	assert.Equal(t, created.ID, provider.ID)
}

func TestService_ListProviders(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	for i := range 2 {
		req := createTestProvider("list_svc_" + string(rune('a'+i)))
		_, err := repo.CreateProvider(ctx, req)
		require.NoError(t, err)
	}

	providers, err := svc.ListProviders(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(providers), 2)
}

func TestService_ListEnabledProviders(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req1 := createTestProvider("enabled_svc")
	req1.IsEnabled = true
	_, err := repo.CreateProvider(ctx, req1)
	require.NoError(t, err)

	req2 := createTestProvider("disabled_svc")
	req2.IsEnabled = false
	_, err = repo.CreateProvider(ctx, req2)
	require.NoError(t, err)

	providers, err := svc.ListEnabledProviders(ctx)
	require.NoError(t, err)

	for _, p := range providers {
		assert.True(t, p.IsEnabled)
	}
}

func TestService_UpdateProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("update_svc")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	updateReq := UpdateProviderRequest{
		DisplayName: new("Updated"),
		IsEnabled:   new(false),
	}

	updated, err := svc.UpdateProvider(ctx, created.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.DisplayName)
	assert.False(t, updated.IsEnabled)
}

func TestService_UpdateProvider_InvalidType(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("update_type_svc")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	invalidType := "invalid"
	updateReq := UpdateProviderRequest{
		ProviderType: &invalidType,
	}

	_, err = svc.UpdateProvider(ctx, created.ID, updateReq)
	assert.ErrorIs(t, err, ErrInvalidProviderType)
}

func TestService_UpdateProvider_InvalidIssuerURL(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("update_url_svc")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	invalidURL := "not a valid url ::"
	updateReq := UpdateProviderRequest{
		IssuerURL: &invalidURL,
	}

	_, err = svc.UpdateProvider(ctx, created.ID, updateReq)
	assert.ErrorIs(t, err, ErrInvalidIssuerURL)
}

func TestService_DeleteProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("delete_svc")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = svc.DeleteProvider(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetProvider(ctx, created.ID)
	assert.Error(t, err)
}

func TestService_EnableProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("enable_svc")
	req.IsEnabled = false
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = svc.EnableProvider(ctx, created.ID)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, provider.IsEnabled)
}

func TestService_DisableProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("disable_svc")
	req.IsEnabled = true
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = svc.DisableProvider(ctx, created.ID)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.False(t, provider.IsEnabled)
}

func TestService_SetDefaultProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("setdefault_svc")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = svc.SetDefaultProvider(ctx, created.ID)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, provider.IsDefault)
}

// ============================================================================
// OAuth Flow Tests
// ============================================================================

func TestService_GetAuthURL(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("auth_url")
	req.IsEnabled = true
	_, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	result, err := svc.GetAuthURL(ctx, "auth_url", "http://localhost/redirect", nil)
	require.NoError(t, err)
	assert.NotEmpty(t, result.URL)
	assert.NotEmpty(t, result.State)
	assert.Contains(t, result.URL, "code_challenge")
	assert.Contains(t, result.URL, "code_challenge_method=S256")
}

func TestService_GetAuthURL_DisabledProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	req := createTestProvider("disabled_auth")
	req.IsEnabled = false
	_, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	_, err = svc.GetAuthURL(ctx, "disabled_auth", "http://localhost/redirect", nil)
	assert.ErrorIs(t, err, ErrProviderDisabled)
}

func TestService_GetAuthURL_ProviderNotFound(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.GetAuthURL(ctx, "nonexistent", "http://localhost/redirect", nil)
	assert.Error(t, err)
}

// ============================================================================
// User Link Tests
// ============================================================================

func TestService_LinkUser(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupTestService(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("link_svc"))
	require.NoError(t, err)
	provider.AllowLinking = true
	updateReq := UpdateProviderRequest{AllowLinking: new(true)}
	provider, err = repo.UpdateProvider(ctx, provider.ID, updateReq)
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "linksvcuser",
		Email:    "linksvc@example.com",
	})

	userInfo := &UserInfo{
		Subject: "sub-svc",
		Email:   "linksvc@example.com",
		Name:    "Link User",
	}

	token := &oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	link, err := svc.LinkUser(ctx, user.ID, provider.ID, "sub-svc", userInfo, token)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, link.ID)
	assert.Equal(t, user.ID, link.UserID)
}

func TestService_ListUserLinks(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "getlinks",
		Email:    "getlinks@example.com",
	})

	provider, err := repo.CreateProvider(ctx, createTestProvider("getlinks_prov"))
	require.NoError(t, err)

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-getlinks",
	}
	_, err = repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	links, err := svc.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(links), 1)
}

func TestService_UnlinkUser(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupTestService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "unlink",
		Email:    "unlink@example.com",
	})

	provider, err := repo.CreateProvider(ctx, createTestProvider("unlink_prov"))
	require.NoError(t, err)

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-unlink",
	}
	_, err = repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	err = svc.UnlinkUser(ctx, user.ID, provider.ID)
	require.NoError(t, err)

	_, err = repo.GetUserLinkByUserAndProvider(ctx, user.ID, provider.ID)
	assert.Error(t, err)
}

func TestService_CleanupExpiredStates(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("cleanup_prov"))
	require.NoError(t, err)

	// Create expired state
	req := CreateStateRequest{
		State:      "state-expired",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
	}
	_, err = repo.CreateState(ctx, req)
	require.NoError(t, err)

	count, err := svc.CleanupExpiredStates(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestService_ValidateProviderType(t *testing.T) {
	tests := []struct {
		name  string
		ptype string
		valid bool
	}{
		{"generic", "generic", true},
		{"oidc", "oidc", true},
		{"authentik", "authentik", true},
		{"keycloak", "keycloak", true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, isValidProviderType(tt.ptype))
		})
	}
}

func TestService_GenerateRandomString(t *testing.T) {
	str1, err := generateRandomString(32)
	require.NoError(t, err)
	assert.Len(t, str1, 64) // hex encoding doubles length

	str2, err := generateRandomString(32)
	require.NoError(t, err)
	assert.NotEqual(t, str1, str2) // Should be random
}

func TestService_GenerateCodeChallenge(t *testing.T) {
	verifier := "test_verifier_string"
	challenge := generateCodeChallenge(verifier)

	assert.NotEmpty(t, challenge)
	assert.NotEqual(t, verifier, challenge)

	// Same verifier should produce same challenge
	challenge2 := generateCodeChallenge(verifier)
	assert.Equal(t, challenge, challenge2)
}

func TestService_ExtractUserInfo(t *testing.T) {
	svc, _, _ := setupTestService(t)

	provider := &Provider{
		ClaimMappings: ClaimMappings{
			Username: "preferred_username",
			Email:    "email",
			Name:     "name",
			Picture:  "picture",
		},
	}

	claims := map[string]any{
		"preferred_username": "testuser",
		"email":              "test@example.com",
		"name":               "Test User",
		"picture":            "http://example.com/pic.jpg",
	}

	userInfo := svc.extractUserInfo(provider, claims)
	assert.Equal(t, "testuser", userInfo.Username)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.Equal(t, "http://example.com/pic.jpg", userInfo.Picture)
}

func TestService_EncryptDecryptSecret(t *testing.T) {
	svc, _, _ := setupTestService(t)

	plaintext := []byte("my-secret-key")

	encrypted, err := svc.encryptSecret(plaintext)
	require.NoError(t, err)

	decrypted := svc.decryptSecret(encrypted)
	assert.Equal(t, plaintext, decrypted)
}

func TestService_BuildOAuth2Config(t *testing.T) {
	svc, _, _ := setupTestService(t)

	authURL := "https://example.com/auth"
	tokenURL := "https://example.com/token"
	provider := &Provider{
		Name:                  "test",
		ClientID:              "test-client",
		ClientSecretEncrypted: []byte("secret"),
		Scopes:                []string{"openid", "profile", "email"},
		AuthorizationEndpoint: &authURL,
		TokenEndpoint:         &tokenURL,
	}

	config := svc.buildOAuth2Config(provider, nil)
	assert.Equal(t, "test-client", config.ClientID)
	assert.Equal(t, "https://example.com/auth", config.Endpoint.AuthURL)
	assert.Equal(t, "https://example.com/token", config.Endpoint.TokenURL)
	assert.Contains(t, config.Scopes, "openid")
}

// ============================================================================
// HandleCallback Tests
// ============================================================================

func TestService_HandleCallback_InvalidState(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupTestService(t)
	ctx := context.Background()

	_, err := svc.HandleCallback(ctx, "invalid-state", "some-code")
	assert.ErrorIs(t, err, ErrInvalidState)
}

func TestService_HandleCallback_ExpiredState(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("callback_expired"))
	require.NoError(t, err)

	// Create expired state
	stateReq := CreateStateRequest{
		State:      "expired-state-123",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(-1 * time.Hour), // Already expired
	}
	_, err = repo.CreateState(ctx, stateReq)
	require.NoError(t, err)

	_, err = svc.HandleCallback(ctx, "expired-state-123", "some-code")
	assert.ErrorIs(t, err, ErrStateExpired)
}

func TestService_HandleCallback_DisabledProvider(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupTestService(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("callback_disabled"))
	require.NoError(t, err)

	// Disable the provider
	err = repo.DisableProvider(ctx, provider.ID)
	require.NoError(t, err)

	// Create valid state
	stateReq := CreateStateRequest{
		State:      "valid-state-disabled",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(5 * time.Minute),
	}
	_, err = repo.CreateState(ctx, stateReq)
	require.NoError(t, err)

	_, err = svc.HandleCallback(ctx, "valid-state-disabled", "some-code")
	assert.ErrorIs(t, err, ErrProviderDisabled)
}

// ============================================================================
// ExtractUserInfo Extended Tests
// ============================================================================

func TestService_ExtractUserInfo_MissingClaims(t *testing.T) {
	svc, _, _ := setupTestService(t)

	provider := &Provider{
		ClaimMappings: ClaimMappings{
			Username: "preferred_username",
			Email:    "email",
			Name:     "name",
		},
	}

	// Empty claims
	claims := map[string]any{}

	userInfo := svc.extractUserInfo(provider, claims)
	assert.Empty(t, userInfo.Username)
	assert.Empty(t, userInfo.Email)
	assert.Empty(t, userInfo.Name)
}

func TestService_ExtractUserInfo_WithRoles(t *testing.T) {
	svc, _, _ := setupTestService(t)

	provider := &Provider{
		ClaimMappings: ClaimMappings{
			Username: "preferred_username",
			Email:    "email",
			Roles:    "roles",
		},
	}

	claims := map[string]any{
		"preferred_username": "testuser",
		"email":              "test@example.com",
		"roles":              []any{"admin", "user"},
	}

	userInfo := svc.extractUserInfo(provider, claims)
	assert.Equal(t, "testuser", userInfo.Username)
	assert.Equal(t, "test@example.com", userInfo.Email)
}

// ============================================================================
// LinkUser Extended Tests
// ============================================================================

func TestService_LinkUser_LinkingNotAllowed(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupTestService(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("nolink_svc"))
	require.NoError(t, err)

	// Explicitly disable linking
	updateReq := UpdateProviderRequest{AllowLinking: new(false)}
	provider, err = repo.UpdateProvider(ctx, provider.ID, updateReq)
	require.NoError(t, err)
	assert.False(t, provider.AllowLinking)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "nolinksvcuser",
		Email:    "nolinksvc@example.com",
	})

	userInfo := &UserInfo{
		Subject: "sub-nolink",
		Email:   "nolinksvc@example.com",
		Name:    "No Link User",
	}

	token := &oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	_, err = svc.LinkUser(ctx, user.ID, provider.ID, "sub-nolink", userInfo, token)
	assert.ErrorIs(t, err, ErrLinkingDisabled)
}

func TestService_LinkUser_AlreadyLinked(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupTestService(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("already_linked"))
	require.NoError(t, err)
	updateReq := UpdateProviderRequest{AllowLinking: new(true)}
	provider, err = repo.UpdateProvider(ctx, provider.ID, updateReq)
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "alreadylinked",
		Email:    "alreadylinked@example.com",
	})

	// First link
	linkReq := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-already",
	}
	_, err = repo.CreateUserLink(ctx, linkReq)
	require.NoError(t, err)

	// Try to link again
	userInfo := &UserInfo{
		Subject: "sub-already",
		Email:   "alreadylinked@example.com",
		Name:    "Already Linked User",
	}

	token := &oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	_, err = svc.LinkUser(ctx, user.ID, provider.ID, "sub-already", userInfo, token)
	assert.ErrorIs(t, err, ErrUserLinkExists)
}
