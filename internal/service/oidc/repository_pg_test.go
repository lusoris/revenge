package oidc

import (
	"context"
	"os"
	"testing"
	"time"

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

func setupTestRepository(t *testing.T) (*RepositoryPg, testutil.DB) {
	t.Helper()
	testDB := testutil.NewFastTestDB(t)
	queries := db.New(testDB.Pool())
	repo := NewRepositoryPg(queries)
	return repo, testDB
}

func createTestProvider(name string) CreateProviderRequest {
	return CreateProviderRequest{
		Name:                  name,
		DisplayName:           name + " Display",
		ProviderType:          "oidc",
		IssuerURL:             "https://example.com/" + name,
		ClientID:              "client-" + name,
		ClientSecretEncrypted: []byte("encrypted-secret"),
		Scopes:                []string{"openid", "profile", "email"},
		ClaimMappings: ClaimMappings{
			Username: "preferred_username",
			Email:    "email",
			Name:     "name",
		},
		AutoCreateUsers: true,
		UpdateUserInfo:  true,
		AllowLinking:    true,
		IsEnabled:       true,
		IsDefault:       false,
	}
}

// ============================================================================
// Provider CRUD Tests
// ============================================================================

func TestRepositoryPg_CreateProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("test")
	provider, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, provider.ID)
	assert.Equal(t, "test", provider.Name)
	assert.False(t, provider.CreatedAt.IsZero())
	assert.False(t, provider.UpdatedAt.IsZero())
}

func TestRepositoryPg_GetProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("get_test")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, provider.ID)
	assert.Equal(t, "get_test", provider.Name)
}

func TestRepositoryPg_GetProvider_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetProvider(ctx, uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
}

func TestRepositoryPg_GetProviderByName(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("by_name")
	_, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	provider, err := repo.GetProviderByName(ctx, "by_name")
	require.NoError(t, err)
	assert.Equal(t, "by_name", provider.Name)
}

func TestRepositoryPg_GetProviderByName_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetProviderByName(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestRepositoryPg_ListProviders(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Create multiple providers
	for i := range 3 {
		req := createTestProvider("list_" + string(rune('a'+i)))
		_, err := repo.CreateProvider(ctx, req)
		require.NoError(t, err)
	}

	providers, err := repo.ListProviders(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(providers), 3)
}

func TestRepositoryPg_ListEnabledProviders(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	// Create enabled provider
	req1 := createTestProvider("enabled")
	req1.IsEnabled = true
	_, err := repo.CreateProvider(ctx, req1)
	require.NoError(t, err)

	// Create disabled provider
	req2 := createTestProvider("disabled")
	req2.IsEnabled = false
	_, err = repo.CreateProvider(ctx, req2)
	require.NoError(t, err)

	providers, err := repo.ListEnabledProviders(ctx)
	require.NoError(t, err)

	// Check all are enabled
	for _, p := range providers {
		assert.True(t, p.IsEnabled)
	}
}

func TestRepositoryPg_UpdateProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("update")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	updateReq := UpdateProviderRequest{
		DisplayName: new("Updated Display"),
		IsEnabled:   new(false),
	}

	updated, err := repo.UpdateProvider(ctx, created.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "Updated Display", updated.DisplayName)
	assert.False(t, updated.IsEnabled)
	assert.True(t, updated.UpdatedAt.After(created.UpdatedAt))
}

func TestRepositoryPg_DeleteProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("delete")
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = repo.DeleteProvider(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetProvider(ctx, created.ID)
	assert.Error(t, err)
}

func TestRepositoryPg_EnableProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("enable")
	req.IsEnabled = false
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = repo.EnableProvider(ctx, created.ID)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, provider.IsEnabled)
}

func TestRepositoryPg_DisableProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("disable")
	req.IsEnabled = true
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = repo.DisableProvider(ctx, created.ID)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.False(t, provider.IsEnabled)
}

func TestRepositoryPg_SetDefaultProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("default")
	req.IsDefault = false
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	err = repo.SetDefaultProvider(ctx, created.ID)
	require.NoError(t, err)

	provider, err := repo.GetProvider(ctx, created.ID)
	require.NoError(t, err)
	assert.True(t, provider.IsDefault)
}

func TestRepositoryPg_GetDefaultProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	req := createTestProvider("default_provider")
	req.IsDefault = true
	created, err := repo.CreateProvider(ctx, req)
	require.NoError(t, err)

	provider, err := repo.GetDefaultProvider(ctx)
	require.NoError(t, err)
	assert.Equal(t, created.ID, provider.ID)
	assert.True(t, provider.IsDefault)
}

// ============================================================================
// User Link Tests
// ============================================================================

func TestRepositoryPg_CreateUserLink(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	// Create provider
	provider, err := repo.CreateProvider(ctx, createTestProvider("link_provider"))
	require.NoError(t, err)

	// Create user
	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "linkuser",
		Email:    "link@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:                user.ID,
		ProviderID:            provider.ID,
		Subject:               "sub-12345",
		Email:                 new("link@example.com"),
		AccessTokenEncrypted:  []byte("encrypted-access"),
		RefreshTokenEncrypted: []byte("encrypted-refresh"),
	}

	link, err := repo.CreateUserLink(ctx, req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, link.ID)
	assert.Equal(t, user.ID, link.UserID)
	assert.Equal(t, provider.ID, link.ProviderID)
	assert.Equal(t, "sub-12345", link.Subject)
}

func TestRepositoryPg_GetUserLink(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("get_link_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "getlinkuser",
		Email:    "getlink@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-get",
	}

	created, err := repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	link, err := repo.GetUserLink(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, link.ID)
	assert.Equal(t, "sub-get", link.Subject)
}

func TestRepositoryPg_GetUserLinkBySubject(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("subject_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "subjectuser",
		Email:    "subject@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-unique",
	}

	_, err = repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	link, err := repo.GetUserLinkBySubject(ctx, provider.ID, "sub-unique")
	require.NoError(t, err)
	assert.Equal(t, "sub-unique", link.Subject)
	assert.Equal(t, user.ID, link.UserID)
}

func TestRepositoryPg_GetUserLinkByUserAndProvider(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("user_prov_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "userprovuser",
		Email:    "userprov@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-userprov",
	}

	_, err = repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	link, err := repo.GetUserLinkByUserAndProvider(ctx, user.ID, provider.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, link.UserID)
	assert.Equal(t, provider.ID, link.ProviderID)
}

func TestRepositoryPg_ListUserLinks(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "multilink",
		Email:    "multilink@example.com",
	})

	// Create multiple providers and links
	for i := range 2 {
		provider, err := repo.CreateProvider(ctx, createTestProvider("list_link_"+string(rune('a'+i))))
		require.NoError(t, err)

		req := CreateUserLinkRequest{
			UserID:     user.ID,
			ProviderID: provider.ID,
			Subject:    "sub-list-" + string(rune('a'+i)),
		}
		_, err = repo.CreateUserLink(ctx, req)
		require.NoError(t, err)
	}

	links, err := repo.ListUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(links), 2)

	// Check provider info is included
	for _, link := range links {
		assert.NotEmpty(t, link.ProviderName)
		assert.NotEmpty(t, link.ProviderDisplayName)
	}
}

func TestRepositoryPg_UpdateUserLink(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("update_link_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "updatelinkuser",
		Email:    "updatelink@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-update",
	}

	created, err := repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	updateReq := UpdateUserLinkRequest{
		Email: new("updated@example.com"),
		Name:  new("Updated Name"),
	}

	updated, err := repo.UpdateUserLink(ctx, created.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "updated@example.com", *updated.Email)
	assert.Equal(t, "Updated Name", *updated.Name)
}

func TestRepositoryPg_UpdateUserLinkLastLogin(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("login_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "loginuser",
		Email:    "login@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-login",
	}

	created, err := repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	err = repo.UpdateUserLinkLastLogin(ctx, created.ID)
	require.NoError(t, err)

	link, err := repo.GetUserLink(ctx, created.ID)
	require.NoError(t, err)
	assert.NotNil(t, link.LastLoginAt)
}

func TestRepositoryPg_DeleteUserLink(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("delete_link_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "deletelinkuser",
		Email:    "deletelink@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-delete",
	}

	created, err := repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	err = repo.DeleteUserLink(ctx, created.ID)
	require.NoError(t, err)

	_, err = repo.GetUserLink(ctx, created.ID)
	assert.Error(t, err)
}

func TestRepositoryPg_DeleteUserLinkByUserAndProvider(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("del_by_up_provider"))
	require.NoError(t, err)

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "delbyupuser",
		Email:    "delbyup@example.com",
	})

	req := CreateUserLinkRequest{
		UserID:     user.ID,
		ProviderID: provider.ID,
		Subject:    "sub-delbyup",
	}

	_, err = repo.CreateUserLink(ctx, req)
	require.NoError(t, err)

	err = repo.DeleteUserLinkByUserAndProvider(ctx, user.ID, provider.ID)
	require.NoError(t, err)

	_, err = repo.GetUserLinkByUserAndProvider(ctx, user.ID, provider.ID)
	assert.Error(t, err)
}

func TestRepositoryPg_CountUserLinks(t *testing.T) {
	t.Parallel()
	repo, testDB := setupTestRepository(t)
	ctx := context.Background()

	user := testutil.CreateUser(t, testDB.Pool(), testutil.User{
		Username: "countlinkuser",
		Email:    "countlink@example.com",
	})

	initialCount, err := repo.CountUserLinks(ctx, user.ID)
	require.NoError(t, err)

	// Create links
	for i := range 2 {
		provider, err := repo.CreateProvider(ctx, createTestProvider("count_link_"+string(rune('a'+i))))
		require.NoError(t, err)

		req := CreateUserLinkRequest{
			UserID:     user.ID,
			ProviderID: provider.ID,
			Subject:    "sub-count-" + string(rune('a'+i)),
		}
		_, err = repo.CreateUserLink(ctx, req)
		require.NoError(t, err)
	}

	count, err := repo.CountUserLinks(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, initialCount+2, count)
}

// ============================================================================
// State Tests
// ============================================================================

func TestRepositoryPg_CreateState(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("state_provider"))
	require.NoError(t, err)

	req := CreateStateRequest{
		State:        "state-12345",
		ProviderID:   provider.ID,
		CodeVerifier: new("verifier-code"),
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}

	state, err := repo.CreateState(ctx, req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, state.ID)
	assert.Equal(t, "state-12345", state.State)
	assert.Equal(t, provider.ID, state.ProviderID)
}

func TestRepositoryPg_GetState(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("get_state_provider"))
	require.NoError(t, err)

	req := CreateStateRequest{
		State:      "state-get",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	}

	_, err = repo.CreateState(ctx, req)
	require.NoError(t, err)

	state, err := repo.GetState(ctx, "state-get")
	require.NoError(t, err)
	assert.Equal(t, "state-get", state.State)
	assert.Equal(t, provider.ID, state.ProviderID)
}

func TestRepositoryPg_GetState_NotFound(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetState(ctx, "nonexistent-state")
	assert.Error(t, err)
}

func TestRepositoryPg_DeleteState(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("del_state_provider"))
	require.NoError(t, err)

	req := CreateStateRequest{
		State:      "state-delete",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	}

	_, err = repo.CreateState(ctx, req)
	require.NoError(t, err)

	err = repo.DeleteState(ctx, "state-delete")
	require.NoError(t, err)

	_, err = repo.GetState(ctx, "state-delete")
	assert.Error(t, err)
}

func TestRepositoryPg_DeleteExpiredStates(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("expired_provider"))
	require.NoError(t, err)

	// Create expired state (we can't actually set created_at in the past, but test the path)
	req := CreateStateRequest{
		State:      "state-old",
		ProviderID: provider.ID,
		ExpiresAt:  time.Now().Add(-1 * time.Hour), // Expired
	}

	_, err = repo.CreateState(ctx, req)
	require.NoError(t, err)

	count, err := repo.DeleteExpiredStates(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

func TestRepositoryPg_DeleteStatesByProvider(t *testing.T) {
	t.Parallel()
	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	provider, err := repo.CreateProvider(ctx, createTestProvider("del_states_provider"))
	require.NoError(t, err)

	// Create multiple states
	for i := range 2 {
		req := CreateStateRequest{
			State:      "state-multi-" + string(rune('a'+i)),
			ProviderID: provider.ID,
			ExpiresAt:  time.Now().Add(10 * time.Minute),
		}
		_, err = repo.CreateState(ctx, req)
		require.NoError(t, err)
	}

	err = repo.DeleteStatesByProvider(ctx, provider.ID)
	require.NoError(t, err)

	// Check states are deleted
	_, err = repo.GetState(ctx, "state-multi-a")
	assert.Error(t, err)
}

// ============================================================================
// Helper Functions
// ============================================================================
