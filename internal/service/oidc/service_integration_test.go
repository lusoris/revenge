package oidc

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

// ============================================================================
// Test Setup Helpers
// ============================================================================

var integrationEncryptKey = bytes.Repeat([]byte{0x42}, 32) // AES-256 key for tests

func setupIntegrationService(t *testing.T) (*Service, *RepositoryPg, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	svc := NewService(repo, zap.NewNop(), "http://localhost:8080/callback", integrationEncryptKey)
	return svc, repo, testDB
}

func setupIntegrationServiceNoEncryption(t *testing.T) (*Service, *RepositoryPg, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	svc := NewService(repo, zap.NewNop(), "http://localhost:8080/callback", nil)
	return svc, repo, testDB
}

// integrationProviderReq returns a CreateProviderRequest with sensible defaults for testing.
func integrationProviderReq(name string) CreateProviderRequest {
	return CreateProviderRequest{
		Name:                  name,
		DisplayName:           name + " Display",
		ProviderType:          ProviderTypeGeneric,
		IssuerURL:             "https://issuer.example.com/" + name,
		ClientID:              "client-" + name,
		ClientSecretEncrypted: []byte("client-secret-" + name),
		Scopes:                []string{"openid", "profile", "email"},
		ClaimMappings: ClaimMappings{
			Username: "preferred_username",
			Email:    "email",
			Name:     "name",
			Picture:  "picture",
		},
		AutoCreateUsers: true,
		UpdateUserInfo:  true,
		AllowLinking:    true,
		IsEnabled:       true,
		IsDefault:       false,
	}
}

// ============================================================================
// Provider CRUD Lifecycle Tests
// ============================================================================

func TestServiceIntegration_ProviderLifecycle(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	// AddProvider
	req := integrationProviderReq("lifecycle")
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, provider)
	assert.NotEqual(t, uuid.Nil, provider.ID)
	assert.Equal(t, "lifecycle", provider.Name)
	assert.Equal(t, "lifecycle Display", provider.DisplayName)
	assert.Equal(t, ProviderTypeGeneric, provider.ProviderType)
	assert.True(t, provider.IsEnabled)
	assert.True(t, provider.AutoCreateUsers)
	assert.True(t, provider.AllowLinking)

	// GetProvider by ID
	got, err := svc.GetProvider(ctx, provider.ID)
	require.NoError(t, err)
	assert.Equal(t, provider.ID, got.ID)
	assert.Equal(t, "lifecycle", got.Name)

	// GetProviderByName
	gotByName, err := svc.GetProviderByName(ctx, "lifecycle")
	require.NoError(t, err)
	assert.Equal(t, provider.ID, gotByName.ID)

	// UpdateProvider
	newDisplay := "Updated Display"
	newType := ProviderTypeKeycloak
	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		DisplayName:  &newDisplay,
		ProviderType: &newType,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Display", updated.DisplayName)
	assert.Equal(t, ProviderTypeKeycloak, updated.ProviderType)

	// ListProviders
	providers, err := svc.ListProviders(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(providers), 1)
	found := false
	for _, p := range providers {
		if p.ID == provider.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "created provider should appear in ListProviders")

	// DeleteProvider
	err = svc.DeleteProvider(ctx, provider.ID)
	require.NoError(t, err)

	// GetProvider after delete should fail
	_, err = svc.GetProvider(ctx, provider.ID)
	require.ErrorIs(t, err, ErrProviderNotFound)
}

// ============================================================================
// Provider Validation Tests
// ============================================================================

func TestServiceIntegration_AddProvider_InvalidType(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("invalid_type")
	req.ProviderType = "invalid_type_xyz"

	_, err := svc.AddProvider(ctx, req)
	require.ErrorIs(t, err, ErrInvalidProviderType)
}

func TestServiceIntegration_AddProvider_DuplicateName(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("duplicate")
	_, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Second add with same name should fail
	_, err = svc.AddProvider(ctx, req)
	require.ErrorIs(t, err, ErrProviderNameExists)
}

func TestServiceIntegration_AddProvider_DefaultScopes(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("default_scopes")
	req.Scopes = nil // Empty scopes should be set to defaults

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, []string{"openid", "profile", "email"}, provider.Scopes)
}

func TestServiceIntegration_AddProvider_DefaultClaimMappings(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("default_claims")
	req.ClaimMappings = ClaimMappings{} // All empty -- should get defaults

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "preferred_username", provider.ClaimMappings.Username)
	assert.Equal(t, "email", provider.ClaimMappings.Email)
	assert.Equal(t, "name", provider.ClaimMappings.Name)
	assert.Equal(t, "picture", provider.ClaimMappings.Picture)
}

func TestServiceIntegration_AddProvider_PreservesCustomClaimMappings(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("custom_claims")
	req.ClaimMappings = ClaimMappings{
		Username: "sub",
		Email:    "mail",
		// Name and Picture left empty -- should get defaults
	}

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "sub", provider.ClaimMappings.Username)
	assert.Equal(t, "mail", provider.ClaimMappings.Email)
	assert.Equal(t, "name", provider.ClaimMappings.Name)
	assert.Equal(t, "picture", provider.ClaimMappings.Picture)
}

func TestServiceIntegration_AddProvider_AllValidTypes(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	validTypes := []string{"oidc", ProviderTypeGeneric, ProviderTypeAuthentik, ProviderTypeKeycloak}

	for _, pt := range validTypes {
		t.Run(pt, func(t *testing.T) {
			req := integrationProviderReq("type_" + pt)
			req.ProviderType = pt

			provider, err := svc.AddProvider(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, pt, provider.ProviderType)
		})
	}
}

func TestServiceIntegration_UpdateProvider_InvalidType(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("update_invalid_type"))
	require.NoError(t, err)

	badType := "not_a_type"
	_, err = svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		ProviderType: &badType,
	})
	require.ErrorIs(t, err, ErrInvalidProviderType)
}

func TestServiceIntegration_UpdateProvider_NonExistent(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	newDisplay := "whatever"
	_, err := svc.UpdateProvider(ctx, uuid.Must(uuid.NewV7()), UpdateProviderRequest{
		DisplayName: &newDisplay,
	})
	require.ErrorIs(t, err, ErrProviderNotFound)
}

func TestServiceIntegration_GetProvider_NonExistent(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	_, err := svc.GetProvider(ctx, uuid.Must(uuid.NewV7()))
	require.ErrorIs(t, err, ErrProviderNotFound)
}

func TestServiceIntegration_GetProviderByName_NonExistent(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	_, err := svc.GetProviderByName(ctx, "nonexistent_provider_xyz")
	require.ErrorIs(t, err, ErrProviderNotFound)
}

func TestServiceIntegration_DeleteProvider_NonExistent(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Deleting a non-existent provider should not error (it's a no-op at SQL level)
	err := svc.DeleteProvider(ctx, uuid.Must(uuid.NewV7()))
	require.NoError(t, err)
}

// ============================================================================
// Encryption Roundtrip Tests (through real DB storage)
// ============================================================================

func TestServiceIntegration_EncryptionRoundtrip(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	clientSecret := []byte("my-super-secret-client-secret")

	req := integrationProviderReq("enc_roundtrip")
	req.ClientSecretEncrypted = clientSecret

	// AddProvider encrypts the secret before storage
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// The stored encrypted value should NOT be the plaintext
	assert.NotEqual(t, clientSecret, provider.ClientSecretEncrypted)

	// Retrieve via repo to get raw stored value
	stored, err := repo.GetProvider(ctx, provider.ID)
	require.NoError(t, err)
	assert.NotEqual(t, clientSecret, stored.ClientSecretEncrypted)

	// The service should be able to decrypt it
	decrypted := svc.decryptSecret(stored.ClientSecretEncrypted)
	assert.Equal(t, clientSecret, decrypted)
}

func TestServiceIntegration_EncryptionRoundtrip_UpdateSecret(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create provider with initial secret
	provider, err := svc.AddProvider(ctx, integrationProviderReq("enc_update"))
	require.NoError(t, err)

	// Update with new secret
	newSecret := []byte("new-super-secret")
	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		ClientSecretEncrypted: newSecret,
	})
	require.NoError(t, err)
	assert.NotEqual(t, newSecret, updated.ClientSecretEncrypted)

	// Verify decryption of updated secret
	stored, err := repo.GetProvider(ctx, updated.ID)
	require.NoError(t, err)
	decrypted := svc.decryptSecret(stored.ClientSecretEncrypted)
	assert.Equal(t, newSecret, decrypted)
}

func TestServiceIntegration_NoEncryptionKey_Passthrough(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationServiceNoEncryption(t)
	ctx := context.Background()

	clientSecret := []byte("plaintext-secret")
	req := integrationProviderReq("no_enc")
	req.ClientSecretEncrypted = clientSecret

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Without encryption key, the secret is stored as-is
	stored, err := repo.GetProvider(ctx, provider.ID)
	require.NoError(t, err)
	assert.Equal(t, clientSecret, stored.ClientSecretEncrypted)
}

// ============================================================================
// Enable / Disable Provider Tests
// ============================================================================

func TestServiceIntegration_EnableDisableProvider(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("enable_disable")
	req.IsEnabled = true
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	assert.True(t, provider.IsEnabled)

	// Disable
	err = svc.DisableProvider(ctx, provider.ID)
	require.NoError(t, err)

	got, err := svc.GetProvider(ctx, provider.ID)
	require.NoError(t, err)
	assert.False(t, got.IsEnabled)

	// Enable
	err = svc.EnableProvider(ctx, provider.ID)
	require.NoError(t, err)

	got, err = svc.GetProvider(ctx, provider.ID)
	require.NoError(t, err)
	assert.True(t, got.IsEnabled)
}

func TestServiceIntegration_ListEnabledProviders(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create enabled provider
	req1 := integrationProviderReq("list_enabled")
	req1.IsEnabled = true
	_, err := svc.AddProvider(ctx, req1)
	require.NoError(t, err)

	// Create disabled provider
	req2 := integrationProviderReq("list_disabled")
	req2.IsEnabled = false
	_, err = svc.AddProvider(ctx, req2)
	require.NoError(t, err)

	// ListEnabledProviders should only return enabled ones
	enabled, err := svc.ListEnabledProviders(ctx)
	require.NoError(t, err)
	for _, p := range enabled {
		assert.True(t, p.IsEnabled, "ListEnabledProviders should only return enabled providers")
	}
}

func TestServiceIntegration_DisableProvider_CleansUpStates(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("disable_cleanup"))
	require.NoError(t, err)

	// Create a state for this provider
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "state-to-cleanup",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	})
	require.NoError(t, err)

	// Verify state exists
	_, err = repo.GetState(ctx, "state-to-cleanup")
	require.NoError(t, err)

	// Disable provider should clean up states
	err = svc.DisableProvider(ctx, provider.ID)
	require.NoError(t, err)

	// State should be deleted
	_, err = repo.GetState(ctx, "state-to-cleanup")
	require.ErrorIs(t, err, ErrStateNotFound)
}

func TestServiceIntegration_DeleteProvider_CleansUpStates(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("delete_cleanup"))
	require.NoError(t, err)

	// Create a state for this provider
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "state-for-deleted-provider",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	})
	require.NoError(t, err)

	// Delete provider should clean up states
	err = svc.DeleteProvider(ctx, provider.ID)
	require.NoError(t, err)

	// State should be deleted
	_, err = repo.GetState(ctx, "state-for-deleted-provider")
	require.ErrorIs(t, err, ErrStateNotFound)
}

// ============================================================================
// Default Provider Tests
// ============================================================================

func TestServiceIntegration_SetDefaultProvider(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create two providers
	p1, err := svc.AddProvider(ctx, integrationProviderReq("default1"))
	require.NoError(t, err)

	p2, err := svc.AddProvider(ctx, integrationProviderReq("default2"))
	require.NoError(t, err)

	// Set p1 as default
	err = svc.SetDefaultProvider(ctx, p1.ID)
	require.NoError(t, err)

	defaultProv, err := svc.GetDefaultProvider(ctx)
	require.NoError(t, err)
	assert.Equal(t, p1.ID, defaultProv.ID)
	assert.True(t, defaultProv.IsDefault)

	// Set p2 as default (should unset p1)
	err = svc.SetDefaultProvider(ctx, p2.ID)
	require.NoError(t, err)

	defaultProv, err = svc.GetDefaultProvider(ctx)
	require.NoError(t, err)
	assert.Equal(t, p2.ID, defaultProv.ID)
	assert.True(t, defaultProv.IsDefault)

	// p1 should no longer be default
	p1Got, err := svc.GetProvider(ctx, p1.ID)
	require.NoError(t, err)
	assert.False(t, p1Got.IsDefault)
}

func TestServiceIntegration_GetDefaultProvider_NoneSet(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create a provider but don't set as default
	req := integrationProviderReq("no_default")
	req.IsDefault = false
	_, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	_, err = svc.GetDefaultProvider(ctx)
	require.ErrorIs(t, err, ErrProviderNotFound)
}

// ============================================================================
// Auth URL and State Tests
// ============================================================================

func TestServiceIntegration_GetAuthURL(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create an enabled provider
	provider, err := svc.AddProvider(ctx, integrationProviderReq("auth_url"))
	require.NoError(t, err)

	// GetAuthURL should create a state in the DB and return a URL
	result, err := svc.GetAuthURL(ctx, "auth_url", "https://redirect.example.com", nil)
	require.NoError(t, err)
	assert.NotEmpty(t, result.URL)
	assert.NotEmpty(t, result.State)

	// URL should contain PKCE challenge
	assert.Contains(t, result.URL, "code_challenge")
	assert.Contains(t, result.URL, "S256")
	assert.Contains(t, result.URL, "client-auth_url")

	// State should exist in DB
	state, err := repo.GetState(ctx, result.State)
	require.NoError(t, err)
	assert.Equal(t, provider.ID, state.ProviderID)
	assert.NotNil(t, state.CodeVerifier)
	assert.Nil(t, state.UserID) // No user ID provided
}

func TestServiceIntegration_GetAuthURL_WithUserID(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupIntegrationService(t)
	ctx := context.Background()

	_, err := svc.AddProvider(ctx, integrationProviderReq("auth_url_user"))
	require.NoError(t, err)

	// Create a real user (FK constraint on oidc_states.user_id)
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "authurluser",
		Email:    "authurl@example.com",
	})

	result, err := svc.GetAuthURL(ctx, "auth_url_user", "https://redirect.example.com", &user.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, result.URL)

	// State should have user ID
	state, err := repo.GetState(ctx, result.State)
	require.NoError(t, err)
	require.NotNil(t, state.UserID)
	assert.Equal(t, user.ID, *state.UserID)
}

func TestServiceIntegration_GetAuthURL_DisabledProvider(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("auth_url_disabled")
	req.IsEnabled = false
	_, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	_, err = svc.GetAuthURL(ctx, "auth_url_disabled", "https://redirect.example.com", nil)
	require.ErrorIs(t, err, ErrProviderDisabled)
}

func TestServiceIntegration_GetAuthURL_NonExistentProvider(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	_, err := svc.GetAuthURL(ctx, "nonexistent_xyz", "https://redirect.example.com", nil)
	require.ErrorIs(t, err, ErrProviderNotFound)
}

func TestServiceIntegration_GetAuthURL_GeneratesUniqueStates(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	_, err := svc.AddProvider(ctx, integrationProviderReq("unique_states"))
	require.NoError(t, err)

	result1, err := svc.GetAuthURL(ctx, "unique_states", "https://redirect.example.com", nil)
	require.NoError(t, err)

	result2, err := svc.GetAuthURL(ctx, "unique_states", "https://redirect.example.com", nil)
	require.NoError(t, err)

	assert.NotEqual(t, result1.State, result2.State)
	assert.NotEqual(t, result1.URL, result2.URL)
}

// ============================================================================
// Cleanup Expired States Tests
// ============================================================================

func TestServiceIntegration_CleanupExpiredStates(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("cleanup"))
	require.NoError(t, err)

	// Create an expired state
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "expired-state",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
	})
	require.NoError(t, err)

	// Create a valid state
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "valid-state",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	// Cleanup should remove expired states
	count, err := svc.CleanupExpiredStates(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// Expired state should be gone
	_, err = repo.GetState(ctx, "expired-state")
	require.ErrorIs(t, err, ErrStateNotFound)

	// Valid state should still exist
	_, err = repo.GetState(ctx, "valid-state")
	require.NoError(t, err)
}

func TestServiceIntegration_CleanupExpiredStates_NoneExpired(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("cleanup_none"))
	require.NoError(t, err)

	// Create only valid states
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "still-valid",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	})
	require.NoError(t, err)

	count, err := svc.CleanupExpiredStates(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// ============================================================================
// User Link Lifecycle Tests
// ============================================================================

func TestServiceIntegration_LinkUserLifecycle(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	// Create provider with linking enabled
	req := integrationProviderReq("link_lifecycle")
	req.AllowLinking = true
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Create a user
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "linklifecycleuser",
		Email:    "linklifecycle@example.com",
	})

	// LinkUser
	userInfo := &UserInfo{
		Subject:  "oidc-subject-123",
		Email:    "oidc@example.com",
		Name:     "OIDC User",
		Picture:  "https://example.com/pic.jpg",
		Username: "oidcuser",
	}

	token := &oauth2.Token{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
		Expiry:       time.Now().Add(time.Hour),
	}

	link, err := svc.LinkUser(ctx, user.ID, provider.ID, userInfo.Subject, userInfo, token)
	require.NoError(t, err)
	require.NotNil(t, link)
	assert.NotEqual(t, uuid.Nil, link.ID)
	assert.Equal(t, user.ID, link.UserID)
	assert.Equal(t, provider.ID, link.ProviderID)
	assert.Equal(t, "oidc-subject-123", link.Subject)

	// ListUserLinks
	links, err := svc.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, "link_lifecycle", links[0].ProviderName)
	assert.Equal(t, "link_lifecycle Display", links[0].ProviderDisplayName)

	// UnlinkUser
	err = svc.UnlinkUser(ctx, user.ID, provider.ID)
	require.NoError(t, err)

	// ListUserLinks after unlink should be empty
	links, err = svc.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, links, 0)
}

func TestServiceIntegration_LinkUser_WithEncryptedTokens(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("link_tokens"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "linktokenuser",
		Email:    "linktoken@example.com",
	})

	userInfo := &UserInfo{
		Subject: "subject-enc",
		Email:   "enc@example.com",
		Name:    "Enc User",
	}

	accessToken := "my-access-token"
	refreshToken := "my-refresh-token"
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(time.Hour),
	}

	link, err := svc.LinkUser(ctx, user.ID, provider.ID, "subject-enc", userInfo, token)
	require.NoError(t, err)

	// Retrieve from repo to verify encrypted tokens are stored
	stored, err := repo.GetUserLink(ctx, link.ID)
	require.NoError(t, err)

	// Tokens should be encrypted (not plaintext)
	assert.NotEqual(t, []byte(accessToken), stored.AccessTokenEncrypted)
	assert.NotEqual(t, []byte(refreshToken), stored.RefreshTokenEncrypted)

	// But they should decrypt correctly
	decryptedAccess := svc.decryptSecret(stored.AccessTokenEncrypted)
	assert.Equal(t, []byte(accessToken), decryptedAccess)

	decryptedRefresh := svc.decryptSecret(stored.RefreshTokenEncrypted)
	assert.Equal(t, []byte(refreshToken), decryptedRefresh)

	// Token expiry should be set
	assert.NotNil(t, stored.TokenExpiresAt)
}

func TestServiceIntegration_LinkUser_WithoutToken(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("link_no_token"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "linknotokenuser",
		Email:    "linknotoken@example.com",
	})

	userInfo := &UserInfo{
		Subject: "subject-no-token",
		Email:   "notoken@example.com",
	}

	link, err := svc.LinkUser(ctx, user.ID, provider.ID, "subject-no-token", userInfo, nil)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, link.ID)
}

func TestServiceIntegration_LinkUser_Duplicate(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("link_dup"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "linkdupuser",
		Email:    "linkdup@example.com",
	})

	userInfo := &UserInfo{Subject: "dup-subject"}

	// First link should succeed
	_, err = svc.LinkUser(ctx, user.ID, provider.ID, "dup-subject", userInfo, nil)
	require.NoError(t, err)

	// Second link to same provider should fail
	_, err = svc.LinkUser(ctx, user.ID, provider.ID, "dup-subject-2", userInfo, nil)
	require.ErrorIs(t, err, ErrUserLinkExists)
}

func TestServiceIntegration_LinkUser_LinkingDisabled(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("no_linking")
	req.AllowLinking = false
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "nolinkuser",
		Email:    "nolink@example.com",
	})

	userInfo := &UserInfo{Subject: "subject-nolink"}

	_, err = svc.LinkUser(ctx, user.ID, provider.ID, "subject-nolink", userInfo, nil)
	require.ErrorIs(t, err, ErrLinkingDisabled)
}

func TestServiceIntegration_LinkUser_ProviderNotFound(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "noprovuser",
		Email:    "noprov@example.com",
	})

	userInfo := &UserInfo{Subject: "subject-noprov"}

	_, err := svc.LinkUser(ctx, user.ID, uuid.Must(uuid.NewV7()), "subject-noprov", userInfo, nil)
	require.ErrorIs(t, err, ErrProviderNotFound)
}

func TestServiceIntegration_UnlinkUser_NoLink(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("unlink_nolink"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "unlinknouser",
		Email:    "unlinkno@example.com",
	})

	// Unlinking when no link exists should not error (SQL DELETE no-op)
	err = svc.UnlinkUser(ctx, user.ID, provider.ID)
	require.NoError(t, err)
}

func TestServiceIntegration_ListUserLinks_Empty(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "emptylinksuser",
		Email:    "emptylinks@example.com",
	})

	links, err := svc.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.Empty(t, links)
}

func TestServiceIntegration_LinkUser_MultipleProviders(t *testing.T) {
	t.Parallel()
	svc, _, testDB := setupIntegrationService(t)
	ctx := context.Background()

	// Create two providers
	p1, err := svc.AddProvider(ctx, integrationProviderReq("multi_p1"))
	require.NoError(t, err)

	p2, err := svc.AddProvider(ctx, integrationProviderReq("multi_p2"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "multiuser",
		Email:    "multi@example.com",
	})

	// Link to both providers
	_, err = svc.LinkUser(ctx, user.ID, p1.ID, "sub-p1", &UserInfo{Subject: "sub-p1"}, nil)
	require.NoError(t, err)

	_, err = svc.LinkUser(ctx, user.ID, p2.ID, "sub-p2", &UserInfo{Subject: "sub-p2"}, nil)
	require.NoError(t, err)

	// Should have 2 links
	links, err := svc.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, links, 2)

	// Unlink from p1 only
	err = svc.UnlinkUser(ctx, user.ID, p1.ID)
	require.NoError(t, err)

	links, err = svc.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, links, 1)
	assert.Equal(t, p2.ID, links[0].ProviderID)
}

// ============================================================================
// BuildOAuth2Config Integration Tests (with real encrypted secrets)
// ============================================================================

func TestServiceIntegration_BuildOAuth2Config_DecryptsSecret(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create provider with encrypted secret
	req := integrationProviderReq("oauth2_config")
	req.ClientSecretEncrypted = []byte("test-client-secret-value")
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Retrieve stored provider (with encrypted secret)
	stored, err := repo.GetProvider(ctx, provider.ID)
	require.NoError(t, err)

	// buildOAuth2Config should decrypt the secret
	config := svc.buildOAuth2Config(stored)
	assert.Equal(t, "test-client-secret-value", config.ClientSecret)
	assert.Equal(t, "client-oauth2_config", config.ClientID)
}

func TestServiceIntegration_BuildOAuth2Config_WithCustomEndpoints(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	authEndpoint := "https://auth.example.com/authorize"
	tokenEndpoint := "https://auth.example.com/token"

	req := integrationProviderReq("custom_endpoints")
	req.AuthorizationEndpoint = &authEndpoint
	req.TokenEndpoint = &tokenEndpoint

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	stored, err := repo.GetProvider(ctx, provider.ID)
	require.NoError(t, err)

	config := svc.buildOAuth2Config(stored)
	assert.Equal(t, authEndpoint, config.Endpoint.AuthURL)
	assert.Equal(t, tokenEndpoint, config.Endpoint.TokenURL)
}

func TestServiceIntegration_BuildOAuth2Config_DefaultEndpoints(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("default_endpoints")
	// No custom endpoints -- should use issuer URL defaults
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	stored, err := repo.GetProvider(ctx, provider.ID)
	require.NoError(t, err)

	config := svc.buildOAuth2Config(stored)
	assert.Equal(t, stored.IssuerURL+"/authorize", config.Endpoint.AuthURL)
	assert.Equal(t, stored.IssuerURL+"/token", config.Endpoint.TokenURL)
}

// ============================================================================
// ExtractUserInfo Integration Tests
// ============================================================================

func TestServiceIntegration_ExtractUserInfo_FullClaims(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("extract_info")
	req.ClaimMappings = ClaimMappings{
		Username: "preferred_username",
		Email:    "email",
		Name:     "name",
		Picture:  "picture",
		Roles:    "roles",
	}
	req.RoleMappings = map[string]string{
		"admin":  "system-admin",
		"editor": "content-editor",
	}

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	stored, err := repo.GetProvider(ctx, provider.ID)
	require.NoError(t, err)

	claims := map[string]any{
		"preferred_username": "johndoe",
		"email":              "john@example.com",
		"name":               "John Doe",
		"picture":            "https://example.com/photo.jpg",
		"roles":              []any{"admin", "editor", "viewer"},
	}

	info := svc.extractUserInfo(stored, claims)
	assert.Equal(t, "johndoe", info.Username)
	assert.Equal(t, "john@example.com", info.Email)
	assert.Equal(t, "John Doe", info.Name)
	assert.Equal(t, "https://example.com/photo.jpg", info.Picture)
	// Only mapped roles should be present
	assert.Equal(t, []string{"system-admin", "content-editor"}, info.Roles)
}

// ============================================================================
// Multiple Encryption Operations Test
// ============================================================================

func TestServiceIntegration_EncryptionDeterminism(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)

	plaintext := []byte("same-secret")

	// Encrypt twice should produce different ciphertexts (random nonce)
	enc1, err := svc.encryptSecret(plaintext)
	require.NoError(t, err)

	enc2, err := svc.encryptSecret(plaintext)
	require.NoError(t, err)

	assert.NotEqual(t, enc1, enc2, "same plaintext should produce different ciphertexts due to random nonce")

	// Both should decrypt to the same value
	dec1 := svc.decryptSecret(enc1)
	dec2 := svc.decryptSecret(enc2)
	assert.Equal(t, plaintext, dec1)
	assert.Equal(t, plaintext, dec2)
}

// ============================================================================
// Full Auth Flow State Management Test
// ============================================================================

func TestServiceIntegration_FullAuthFlowStateManagement(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create provider
	provider, err := svc.AddProvider(ctx, integrationProviderReq("full_flow"))
	require.NoError(t, err)

	// Step 1: Generate auth URL (creates state in DB)
	result, err := svc.GetAuthURL(ctx, "full_flow", "https://app.example.com/auth/callback", nil)
	require.NoError(t, err)

	// Step 2: Verify state exists in DB with correct properties
	state, err := repo.GetState(ctx, result.State)
	require.NoError(t, err)
	assert.Equal(t, provider.ID, state.ProviderID)
	assert.NotNil(t, state.CodeVerifier)
	redirectURL := "https://app.example.com/auth/callback"
	assert.Equal(t, &redirectURL, state.RedirectURL)
	assert.True(t, state.ExpiresAt.After(time.Now()))

	// Step 3: Cleanup expired states (current state should survive)
	count, err := svc.CleanupExpiredStates(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count) // Nothing expired yet

	// State should still exist
	_, err = repo.GetState(ctx, result.State)
	require.NoError(t, err)

	// Step 4: Manually delete the state (simulating used state)
	err = repo.DeleteState(ctx, result.State)
	require.NoError(t, err)

	_, err = repo.GetState(ctx, result.State)
	require.ErrorIs(t, err, ErrStateNotFound)
}

// ============================================================================
// Provider with Role Mappings Stored and Retrieved
// ============================================================================

func TestServiceIntegration_ProviderRoleMappings(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("role_mappings")
	req.RoleMappings = map[string]string{
		"oidc-admin":  "admin",
		"oidc-viewer": "viewer",
		"oidc-editor": "editor",
	}

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Retrieve and verify role mappings are stored correctly
	got, err := svc.GetProvider(ctx, provider.ID)
	require.NoError(t, err)
	assert.Len(t, got.RoleMappings, 3)
	assert.Equal(t, "admin", got.RoleMappings["oidc-admin"])
	assert.Equal(t, "viewer", got.RoleMappings["oidc-viewer"])
	assert.Equal(t, "editor", got.RoleMappings["oidc-editor"])
}

// ============================================================================
// Update Provider with new Scopes and Claim Mappings
// ============================================================================

func TestServiceIntegration_UpdateProvider_ScopesAndClaims(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("update_scopes"))
	require.NoError(t, err)

	newClaims := &ClaimMappings{
		Username: "sub",
		Email:    "user_email",
		Name:     "full_name",
		Picture:  "avatar_url",
		Roles:    "realm_access.roles",
	}

	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		Scopes:        []string{"openid", "groups", "offline_access"},
		ClaimMappings: newClaims,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"openid", "groups", "offline_access"}, updated.Scopes)
	assert.Equal(t, "sub", updated.ClaimMappings.Username)
	assert.Equal(t, "user_email", updated.ClaimMappings.Email)
	assert.Equal(t, "full_name", updated.ClaimMappings.Name)
	assert.Equal(t, "avatar_url", updated.ClaimMappings.Picture)
	assert.Equal(t, "realm_access.roles", updated.ClaimMappings.Roles)
}

// ============================================================================
// Update Provider Toggles
// ============================================================================

func TestServiceIntegration_UpdateProvider_Toggles(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("toggles")
	req.AutoCreateUsers = true
	req.AllowLinking = true
	req.UpdateUserInfo = true
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	falseVal := false
	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		AutoCreateUsers: &falseVal,
		AllowLinking:    &falseVal,
		UpdateUserInfo:  &falseVal,
	})
	require.NoError(t, err)
	assert.False(t, updated.AutoCreateUsers)
	assert.False(t, updated.AllowLinking)
	assert.False(t, updated.UpdateUserInfo)
}

// ============================================================================
// Update Provider Endpoints + All Optional Fields
// ============================================================================

func TestServiceIntegration_UpdateProvider_AllEndpoints(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("update_endpoints"))
	require.NoError(t, err)

	// Update ALL endpoint and config fields at once
	newDisplayName := "Updated Display"
	newProviderType := ProviderTypeGeneric
	newIssuerURL := "https://new-issuer.example.com"
	newClientID := "new-client-id"
	newSecret := []byte("new-secret")
	newAuthEndpoint := "https://new-issuer.example.com/authorize"
	newTokenEndpoint := "https://new-issuer.example.com/token"
	newUserInfoEndpoint := "https://new-issuer.example.com/userinfo"
	newJWKSURI := "https://new-issuer.example.com/.well-known/jwks.json"
	newEndSessionEndpoint := "https://new-issuer.example.com/logout"
	trueVal := true

	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		DisplayName:           &newDisplayName,
		ProviderType:          &newProviderType,
		IssuerURL:             &newIssuerURL,
		ClientID:              &newClientID,
		ClientSecretEncrypted: newSecret,
		AuthorizationEndpoint: &newAuthEndpoint,
		TokenEndpoint:         &newTokenEndpoint,
		UserInfoEndpoint:      &newUserInfoEndpoint,
		JWKSURI:               &newJWKSURI,
		EndSessionEndpoint:    &newEndSessionEndpoint,
		IsEnabled:             &trueVal,
		IsDefault:             &trueVal,
	})
	require.NoError(t, err)
	assert.Equal(t, newDisplayName, updated.DisplayName)
	assert.Equal(t, newIssuerURL, updated.IssuerURL)
	assert.Equal(t, newClientID, updated.ClientID)
	require.NotNil(t, updated.AuthorizationEndpoint)
	assert.Equal(t, newAuthEndpoint, *updated.AuthorizationEndpoint)
	require.NotNil(t, updated.TokenEndpoint)
	assert.Equal(t, newTokenEndpoint, *updated.TokenEndpoint)
	require.NotNil(t, updated.UserInfoEndpoint)
	assert.Equal(t, newUserInfoEndpoint, *updated.UserInfoEndpoint)
	require.NotNil(t, updated.JWKSURI)
	assert.Equal(t, newJWKSURI, *updated.JWKSURI)
	require.NotNil(t, updated.EndSessionEndpoint)
	assert.Equal(t, newEndSessionEndpoint, *updated.EndSessionEndpoint)
	assert.True(t, updated.IsDefault)
}

// ============================================================================
// Update Provider with RoleMappings
// ============================================================================

func TestServiceIntegration_UpdateProvider_RoleMappings(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("update_roles"))
	require.NoError(t, err)

	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		RoleMappings: map[string]string{
			"oidc-admin":  "admin",
			"oidc-viewer": "user",
		},
	})
	require.NoError(t, err)
	require.Len(t, updated.RoleMappings, 2)
	assert.Equal(t, "admin", updated.RoleMappings["oidc-admin"])
	assert.Equal(t, "user", updated.RoleMappings["oidc-viewer"])
}

// ============================================================================
// UpdateUserLink - direct repo test for all fields
// ============================================================================

func TestServiceIntegration_UpdateUserLink_AllFields(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupIntegrationService(t)
	ctx := context.Background()

	// Create provider
	provider, err := svc.AddProvider(ctx, integrationProviderReq("update_link"))
	require.NoError(t, err)

	// Create user
	queries := db.New(testDB.Pool())
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "link_update_user_" + uuid.Must(uuid.NewV7()).String()[:8],
		Email:        "link_update_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Link user
	link, err := svc.LinkUser(ctx, user.ID, provider.ID, "subject-update", &UserInfo{
		Subject: "subject-update",
		Email:   "original@example.com",
		Name:    "Original Name",
	}, &oauth2.Token{AccessToken: "initial-token"})
	require.NoError(t, err)

	// Update all fields
	newEmail := "updated@example.com"
	newName := "Updated Name"
	newPicture := "https://example.com/avatar.png"
	newAccessToken := []byte("encrypted-access-token")
	newRefreshToken := []byte("encrypted-refresh-token")
	tokenExpiry := time.Now().Add(1 * time.Hour)
	lastLogin := time.Now()

	updatedLink, err := repo.UpdateUserLink(ctx, link.ID, UpdateUserLinkRequest{
		Email:                 &newEmail,
		Name:                  &newName,
		PictureURL:            &newPicture,
		AccessTokenEncrypted:  newAccessToken,
		RefreshTokenEncrypted: newRefreshToken,
		TokenExpiresAt:        &tokenExpiry,
		LastLoginAt:           &lastLogin,
	})
	require.NoError(t, err)
	require.NotNil(t, updatedLink.Email)
	assert.Equal(t, newEmail, *updatedLink.Email)
	require.NotNil(t, updatedLink.Name)
	assert.Equal(t, newName, *updatedLink.Name)
	require.NotNil(t, updatedLink.PictureURL)
	assert.Equal(t, newPicture, *updatedLink.PictureURL)
	assert.NotNil(t, updatedLink.TokenExpiresAt)
	assert.NotNil(t, updatedLink.LastLoginAt)
}

// ============================================================================
// GetUserLinkBySubject - found and not-found paths
// ============================================================================

func TestServiceIntegration_GetUserLinkBySubject_Found(t *testing.T) {
	t.Parallel()
	svc, repo, testDB := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("subject_found"))
	require.NoError(t, err)

	queries := db.New(testDB.Pool())
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "subject_user_" + uuid.Must(uuid.NewV7()).String()[:8],
		Email:        "subject_" + uuid.Must(uuid.NewV7()).String()[:8] + "@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	// Link user with specific subject
	_, err = svc.LinkUser(ctx, user.ID, provider.ID, "unique-subject-123", &UserInfo{
		Subject: "unique-subject-123",
		Email:   "subject@example.com",
		Name:    "Subject User",
	}, nil)
	require.NoError(t, err)

	// Find by subject - success path
	link, err := repo.GetUserLinkBySubject(ctx, provider.ID, "unique-subject-123")
	require.NoError(t, err)
	assert.Equal(t, user.ID, link.UserID)
	assert.Equal(t, "unique-subject-123", link.Subject)

	// Not found path
	_, err = repo.GetUserLinkBySubject(ctx, provider.ID, "nonexistent-subject")
	assert.ErrorIs(t, err, ErrUserLinkNotFound)
}

// ============================================================================
// UpdateUserLink not found
// ============================================================================

func TestServiceIntegration_UpdateUserLink_NotFound(t *testing.T) {
	t.Parallel()
	_, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	newName := "nope"
	_, err := repo.UpdateUserLink(ctx, uuid.Must(uuid.NewV7()), UpdateUserLinkRequest{
		Name: &newName,
	})
	assert.ErrorIs(t, err, ErrUserLinkNotFound)
}

// ============================================================================
// CreateProvider with RoleMappings
// ============================================================================

func TestServiceIntegration_AddProvider_WithRoleMappings(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("with_roles")
	req.RoleMappings = map[string]string{
		"oidc-admin":  "admin",
		"oidc-editor": "user",
	}
	authEp := "https://issuer.example.com/authorize"
	tokenEp := "https://issuer.example.com/token"
	userInfoEp := "https://issuer.example.com/userinfo"
	jwksURI := "https://issuer.example.com/.well-known/jwks.json"
	endSessionEp := "https://issuer.example.com/logout"
	req.AuthorizationEndpoint = &authEp
	req.TokenEndpoint = &tokenEp
	req.UserInfoEndpoint = &userInfoEp
	req.JWKSURI = &jwksURI
	req.EndSessionEndpoint = &endSessionEp

	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)
	require.Len(t, provider.RoleMappings, 2)
	assert.Equal(t, "admin", provider.RoleMappings["oidc-admin"])
	require.NotNil(t, provider.AuthorizationEndpoint)
	assert.Equal(t, authEp, *provider.AuthorizationEndpoint)
	require.NotNil(t, provider.TokenEndpoint)
	assert.Equal(t, tokenEp, *provider.TokenEndpoint)
	require.NotNil(t, provider.UserInfoEndpoint)
	assert.Equal(t, userInfoEp, *provider.UserInfoEndpoint)
	require.NotNil(t, provider.JWKSURI)
	assert.Equal(t, jwksURI, *provider.JWKSURI)
	require.NotNil(t, provider.EndSessionEndpoint)
	assert.Equal(t, endSessionEp, *provider.EndSessionEndpoint)
}

// ============================================================================
// HandleCallback early error paths (no HTTP needed)
// ============================================================================

func TestServiceIntegration_HandleCallback_InvalidState(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	_, err := svc.HandleCallback(ctx, "nonexistent-state", "some-code")
	assert.ErrorIs(t, err, ErrInvalidState)
}

func TestServiceIntegration_HandleCallback_ExpiredState(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create a provider first
	provider, err := svc.AddProvider(ctx, integrationProviderReq("callback_expired"))
	require.NoError(t, err)

	// Create an already-expired state
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "expired-state-test",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
	})
	require.NoError(t, err)

	_, err = svc.HandleCallback(ctx, "expired-state-test", "some-code")
	assert.ErrorIs(t, err, ErrStateExpired)
}

func TestServiceIntegration_HandleCallback_ProviderDisabled(t *testing.T) {
	t.Parallel()
	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create a disabled provider
	req := integrationProviderReq("callback_disabled")
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Disable it
	err = svc.DisableProvider(ctx, provider.ID)
	require.NoError(t, err)

	// Create a valid state pointing to the disabled provider
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "disabled-provider-state",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	})
	require.NoError(t, err)

	_, err = svc.HandleCallback(ctx, "disabled-provider-state", "some-code")
	assert.ErrorIs(t, err, ErrProviderDisabled)
}

func TestServiceIntegration_HandleCallback_DiscoverySuccessTokenExchangeFail(t *testing.T) {
	t.Parallel()

	// Start fake OIDC discovery server
	mux := http.NewServeMux()
	var serverURL string
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"issuer": "%s",
			"authorization_endpoint": "%s/authorize",
			"token_endpoint": "%s/token",
			"jwks_uri": "%s/jwks",
			"userinfo_endpoint": "%s/userinfo",
			"id_token_signing_alg_values_supported": ["RS256"]
		}`, serverURL, serverURL, serverURL, serverURL, serverURL)
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"invalid_grant"}`)
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"keys":[]}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	serverURL = srv.URL

	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	// Create provider pointing to our test server
	req := integrationProviderReq("callback_discovery")
	req.IssuerURL = serverURL
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Create valid state
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:      "discovery-test-state",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	})
	require.NoError(t, err)

	// HandleCallback should succeed through discovery but fail at token exchange
	_, err = svc.HandleCallback(ctx, "discovery-test-state", "fake-auth-code")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTokenExchange)
}

func TestServiceIntegration_HandleCallback_WithCodeVerifier(t *testing.T) {
	t.Parallel()

	// Start fake OIDC server
	mux := http.NewServeMux()
	var serverURL string
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"issuer": "%s",
			"authorization_endpoint": "%s/authorize",
			"token_endpoint": "%s/token",
			"jwks_uri": "%s/jwks",
			"id_token_signing_alg_values_supported": ["RS256"]
		}`, serverURL, serverURL, serverURL, serverURL)
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"invalid_grant"}`)
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"keys":[]}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	serverURL = srv.URL

	svc, repo, _ := setupIntegrationService(t)
	ctx := context.Background()

	req := integrationProviderReq("callback_pkce")
	req.IssuerURL = serverURL
	provider, err := svc.AddProvider(ctx, req)
	require.NoError(t, err)

	// Create state with PKCE code verifier
	verifier := "test-code-verifier-1234567890abcdef"
	_, err = repo.CreateState(ctx, CreateStateRequest{
		State:        "pkce-test-state",
		ProviderID:   provider.ID,
		CodeVerifier: &verifier,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	})
	require.NoError(t, err)

	// Should exercise the code_verifier branch in HandleCallback
	_, err = svc.HandleCallback(ctx, "pkce-test-state", "fake-auth-code")
	assert.ErrorIs(t, err, ErrTokenExchange)
}

// ============================================================================
// UpdateProvider with invalid IssuerURL
// ============================================================================

func TestServiceIntegration_UpdateProvider_InvalidIssuerURL(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("invalid_issuer"))
	require.NoError(t, err)

	invalidURL := "://not-a-url"
	_, err = svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		IssuerURL: &invalidURL,
	})
	assert.ErrorIs(t, err, ErrInvalidIssuerURL)
}

// ============================================================================
// UpdateProvider with encrypted client secret
// ============================================================================

func TestServiceIntegration_UpdateProvider_ClientSecret(t *testing.T) {
	t.Parallel()
	svc, _, _ := setupIntegrationService(t)
	ctx := context.Background()

	provider, err := svc.AddProvider(ctx, integrationProviderReq("update_secret"))
	require.NoError(t, err)

	updated, err := svc.UpdateProvider(ctx, provider.ID, UpdateProviderRequest{
		ClientSecretEncrypted: []byte("new-secret-value"),
	})
	require.NoError(t, err)
	// Secret should be encrypted (different from plaintext)
	assert.NotEqual(t, "new-secret-value", string(updated.ClientSecretEncrypted))
	assert.NotEmpty(t, updated.ClientSecretEncrypted)
}
