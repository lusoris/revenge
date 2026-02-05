package oidc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func makeTestProvider(id uuid.UUID, name string, enabled bool) *oidc.Provider {
	now := time.Now()
	return &oidc.Provider{
		ID:              id,
		Name:            name,
		DisplayName:     name,
		ProviderType:    oidc.ProviderTypeGeneric,
		IssuerURL:       "https://auth.example.com",
		ClientID:        "test-client-id",
		Scopes:          []string{"openid", "profile", "email"},
		ClaimMappings:   oidc.ClaimMappings{Username: "preferred_username", Email: "email", Name: "name", Picture: "picture"},
		AutoCreateUsers: true,
		AllowLinking:    true,
		IsEnabled:       enabled,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func makeTestUserLink(id, userID, providerID uuid.UUID, subject string) *oidc.UserLink {
	now := time.Now()
	email := "user@example.com"
	return &oidc.UserLink{
		ID:         id,
		UserID:     userID,
		ProviderID: providerID,
		Subject:    subject,
		Email:      &email,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func setupOIDCService(repo oidc.Repository) *oidc.Service {
	logger := zap.NewNop()
	return oidc.NewService(repo, logger, "https://app.example.com/callback", nil)
}

func TestOIDCService_AddProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		created := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(nil, oidc.ErrProviderNotFound)
		mockRepo.On("CreateProvider", mock.Anything, mock.AnythingOfType("oidc.CreateProviderRequest")).
			Return(created, nil)

		req := oidc.CreateProviderRequest{
			Name:                  "test-provider",
			DisplayName:          "Test Provider",
			ProviderType:         oidc.ProviderTypeGeneric,
			IssuerURL:            "https://auth.example.com",
			ClientID:             "test-client-id",
			ClientSecretEncrypted: []byte("secret"),
		}

		provider, err := svc.AddProvider(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, "test-provider", provider.Name)
	})

	t.Run("invalid provider type", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		req := oidc.CreateProviderRequest{
			Name:         "test-provider",
			ProviderType: "invalid_type",
			IssuerURL:    "https://auth.example.com",
		}

		provider, err := svc.AddProvider(context.Background(), req)

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrInvalidProviderType)
	})

	t.Run("invalid issuer URL", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		req := oidc.CreateProviderRequest{
			Name:         "test-provider",
			ProviderType: oidc.ProviderTypeGeneric,
			IssuerURL:    "://invalid-url",
		}

		provider, err := svc.AddProvider(context.Background(), req)

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrInvalidIssuerURL)
	})

	t.Run("provider name exists", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		existing := makeTestProvider(uuid.New(), "test-provider", true)
		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(existing, nil)

		req := oidc.CreateProviderRequest{
			Name:         "test-provider",
			ProviderType: oidc.ProviderTypeGeneric,
			IssuerURL:    "https://auth.example.com",
		}

		provider, err := svc.AddProvider(context.Background(), req)

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrProviderNameExists)
	})
}

func TestOIDCService_GetProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		expected := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(expected, nil)

		provider, err := svc.GetProvider(context.Background(), providerID)

		require.NoError(t, err)
		assert.Equal(t, providerID, provider.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(nil, oidc.ErrProviderNotFound)

		provider, err := svc.GetProvider(context.Background(), providerID)

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})
}

func TestOIDCService_GetProviderByName_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		expected := makeTestProvider(uuid.New(), "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(expected, nil)

		provider, err := svc.GetProviderByName(context.Background(), "test-provider")

		require.NoError(t, err)
		assert.Equal(t, "test-provider", provider.Name)
	})
}

func TestOIDCService_GetDefaultProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		expected := makeTestProvider(uuid.New(), "default-provider", true)
		expected.IsDefault = true

		mockRepo.On("GetDefaultProvider", mock.Anything).Return(expected, nil)

		provider, err := svc.GetDefaultProvider(context.Background())

		require.NoError(t, err)
		assert.True(t, provider.IsDefault)
	})
}

func TestOIDCService_ListProviders_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providers := []oidc.Provider{
			*makeTestProvider(uuid.New(), "provider1", true),
			*makeTestProvider(uuid.New(), "provider2", false),
		}

		mockRepo.On("ListProviders", mock.Anything).Return(providers, nil)

		result, err := svc.ListProviders(context.Background())

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

func TestOIDCService_ListEnabledProviders_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providers := []oidc.Provider{
			*makeTestProvider(uuid.New(), "enabled-provider", true),
		}

		mockRepo.On("ListEnabledProviders", mock.Anything).Return(providers, nil)

		result, err := svc.ListEnabledProviders(context.Background())

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.True(t, result[0].IsEnabled)
	})
}

func TestOIDCService_UpdateProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		updated := makeTestProvider(providerID, "test-provider", true)
		updated.DisplayName = "Updated Provider"

		mockRepo.On("UpdateProvider", mock.Anything, providerID, mock.AnythingOfType("oidc.UpdateProviderRequest")).
			Return(updated, nil)

		newDisplayName := "Updated Provider"
		req := oidc.UpdateProviderRequest{DisplayName: &newDisplayName}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)

		require.NoError(t, err)
		assert.Equal(t, "Updated Provider", provider.DisplayName)
	})

	t.Run("invalid provider type", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		invalidType := "invalid_type"
		req := oidc.UpdateProviderRequest{ProviderType: &invalidType}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrInvalidProviderType)
	})

	t.Run("invalid issuer URL", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		invalidURL := "://invalid"
		req := oidc.UpdateProviderRequest{IssuerURL: &invalidURL}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrInvalidIssuerURL)
	})
}

func TestOIDCService_DeleteProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("DeleteStatesByProvider", mock.Anything, providerID).Return(nil)
		mockRepo.On("DeleteProvider", mock.Anything, providerID).Return(nil)

		err := svc.DeleteProvider(context.Background(), providerID)

		assert.NoError(t, err)
	})

	t.Run("states cleanup error ignored", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("DeleteStatesByProvider", mock.Anything, providerID).Return(errors.New("cleanup error"))
		mockRepo.On("DeleteProvider", mock.Anything, providerID).Return(nil)

		err := svc.DeleteProvider(context.Background(), providerID)

		assert.NoError(t, err) // Cleanup error is logged but not returned
	})
}

func TestOIDCService_EnableProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("EnableProvider", mock.Anything, providerID).Return(nil)

		err := svc.EnableProvider(context.Background(), providerID)

		assert.NoError(t, err)
	})
}

func TestOIDCService_DisableProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("DeleteStatesByProvider", mock.Anything, providerID).Return(nil)
		mockRepo.On("DisableProvider", mock.Anything, providerID).Return(nil)

		err := svc.DisableProvider(context.Background(), providerID)

		assert.NoError(t, err)
	})
}

func TestOIDCService_SetDefaultProvider_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("SetDefaultProvider", mock.Anything, providerID).Return(nil)

		err := svc.SetDefaultProvider(context.Background(), providerID)

		assert.NoError(t, err)
	})
}

func TestOIDCService_UnlinkUser_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()

		mockRepo.On("DeleteUserLinkByUserAndProvider", mock.Anything, userID, providerID).Return(nil)

		err := svc.UnlinkUser(context.Background(), userID, providerID)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()

		mockRepo.On("DeleteUserLinkByUserAndProvider", mock.Anything, userID, providerID).
			Return(oidc.ErrUserLinkNotFound)

		err := svc.UnlinkUser(context.Background(), userID, providerID)

		assert.ErrorIs(t, err, oidc.ErrUserLinkNotFound)
	})
}

func TestOIDCService_ListUserLinks_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		links := []oidc.UserLinkWithProvider{
			{
				UserLink:            *makeTestUserLink(uuid.New(), userID, uuid.New(), "subject1"),
				ProviderName:        "provider1",
				ProviderDisplayName: "Provider 1",
			},
		}

		mockRepo.On("ListUserLinks", mock.Anything, userID).Return(links, nil)

		result, err := svc.ListUserLinks(context.Background(), userID)

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestOIDCService_CleanupExpiredStates_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("DeleteExpiredStates", mock.Anything).Return(int64(10), nil)

		count, err := svc.CleanupExpiredStates(context.Background())

		require.NoError(t, err)
		assert.Equal(t, int64(10), count)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("DeleteExpiredStates", mock.Anything).Return(int64(0), errors.New("db error"))

		count, err := svc.CleanupExpiredStates(context.Background())

		assert.Equal(t, int64(0), count)
		assert.Error(t, err)
	})
}

func TestOIDCService_LinkUser_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("linking disabled", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)
		provider.AllowLinking = false

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(provider, nil)

		link, err := svc.LinkUser(context.Background(), userID, providerID, "subject", nil, nil)

		assert.Nil(t, link)
		assert.ErrorIs(t, err, oidc.ErrLinkingDisabled)
	})

	t.Run("user already linked", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)
		existingLink := makeTestUserLink(uuid.New(), userID, providerID, "subject")

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(provider, nil)
		mockRepo.On("GetUserLinkByUserAndProvider", mock.Anything, userID, providerID).
			Return(existingLink, nil)

		link, err := svc.LinkUser(context.Background(), userID, providerID, "subject", nil, nil)

		assert.Nil(t, link)
		assert.ErrorIs(t, err, oidc.ErrUserLinkExists)
	})

	t.Run("provider not found", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(nil, oidc.ErrProviderNotFound)

		link, err := svc.LinkUser(context.Background(), userID, providerID, "subject", nil, nil)

		assert.Nil(t, link)
		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})

	t.Run("success with user info", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()
		linkID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(provider, nil)
		mockRepo.On("GetUserLinkByUserAndProvider", mock.Anything, userID, providerID).
			Return(nil, oidc.ErrUserLinkNotFound)

		createdLink := makeTestUserLink(linkID, userID, providerID, "subject123")
		mockRepo.On("CreateUserLink", mock.Anything, mock.AnythingOfType("oidc.CreateUserLinkRequest")).
			Return(createdLink, nil)

		userInfo := &oidc.UserInfo{
			Subject:  "subject123",
			Email:    "user@example.com",
			Name:     "Test User",
			Username: "testuser",
		}

		link, err := svc.LinkUser(context.Background(), userID, providerID, "subject123", userInfo, nil)

		require.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, linkID, link.ID)
	})

	t.Run("check existing link error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(provider, nil)
		mockRepo.On("GetUserLinkByUserAndProvider", mock.Anything, userID, providerID).
			Return(nil, errors.New("database error"))

		link, err := svc.LinkUser(context.Background(), userID, providerID, "subject", nil, nil)

		assert.Nil(t, link)
		assert.Error(t, err)
	})
}

func TestOIDCService_AddProvider_Extended_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("GetProviderByName returns unexpected error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").
			Return(nil, errors.New("database connection error"))

		req := oidc.CreateProviderRequest{
			Name:                  "test-provider",
			ProviderType:         oidc.ProviderTypeGeneric,
			IssuerURL:            "https://auth.example.com",
			ClientID:             "test-client-id",
			ClientSecretEncrypted: []byte("secret"),
		}

		provider, err := svc.AddProvider(context.Background(), req)

		assert.Nil(t, provider)
		assert.Error(t, err)
		assert.NotErrorIs(t, err, oidc.ErrProviderNameExists)
	})

	t.Run("sets default scopes when empty", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		created := makeTestProvider(providerID, "test-provider", true)
		created.Scopes = []string{"openid", "profile", "email"}

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(nil, oidc.ErrProviderNotFound)
		mockRepo.On("CreateProvider", mock.Anything, mock.MatchedBy(func(req oidc.CreateProviderRequest) bool {
			return len(req.Scopes) == 3 && req.Scopes[0] == "openid"
		})).Return(created, nil)

		req := oidc.CreateProviderRequest{
			Name:                  "test-provider",
			ProviderType:         oidc.ProviderTypeGeneric,
			IssuerURL:            "https://auth.example.com",
			ClientID:             "test-client-id",
			ClientSecretEncrypted: []byte("secret"),
			Scopes:               []string{}, // Empty scopes
		}

		provider, err := svc.AddProvider(context.Background(), req)

		require.NoError(t, err)
		assert.NotNil(t, provider)
	})

	t.Run("sets default claim mappings", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		created := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(nil, oidc.ErrProviderNotFound)
		mockRepo.On("CreateProvider", mock.Anything, mock.MatchedBy(func(req oidc.CreateProviderRequest) bool {
			return req.ClaimMappings.Username == "preferred_username" &&
				req.ClaimMappings.Email == "email" &&
				req.ClaimMappings.Name == "name" &&
				req.ClaimMappings.Picture == "picture"
		})).Return(created, nil)

		req := oidc.CreateProviderRequest{
			Name:                  "test-provider",
			ProviderType:         oidc.ProviderTypeGeneric,
			IssuerURL:            "https://auth.example.com",
			ClientID:             "test-client-id",
			ClientSecretEncrypted: []byte("secret"),
			// No ClaimMappings set - should use defaults
		}

		provider, err := svc.AddProvider(context.Background(), req)

		require.NoError(t, err)
		assert.NotNil(t, provider)
	})

	t.Run("valid provider types", func(t *testing.T) {
		validTypes := []string{"oidc", "generic", "authentik", "keycloak"}

		for _, pt := range validTypes {
			t.Run(pt, func(t *testing.T) {
				mockRepo := NewMockOIDCRepository(t)
				svc := setupOIDCService(mockRepo)

				providerID := uuid.New()
				created := makeTestProvider(providerID, "test-provider", true)
				created.ProviderType = pt

				mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(nil, oidc.ErrProviderNotFound)
				mockRepo.On("CreateProvider", mock.Anything, mock.AnythingOfType("oidc.CreateProviderRequest")).
					Return(created, nil)

				req := oidc.CreateProviderRequest{
					Name:                  "test-provider",
					ProviderType:         pt,
					IssuerURL:            "https://auth.example.com",
					ClientID:             "test-client-id",
					ClientSecretEncrypted: []byte("secret"),
				}

				provider, err := svc.AddProvider(context.Background(), req)

				require.NoError(t, err)
				assert.Equal(t, pt, provider.ProviderType)
			})
		}
	})
}

func TestOIDCService_UpdateProvider_Extended_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("update with client secret", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		updated := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("UpdateProvider", mock.Anything, providerID, mock.AnythingOfType("oidc.UpdateProviderRequest")).
			Return(updated, nil)

		req := oidc.UpdateProviderRequest{
			ClientSecretEncrypted: []byte("new-secret"),
		}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)

		require.NoError(t, err)
		assert.NotNil(t, provider)
	})

	t.Run("update repo error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("UpdateProvider", mock.Anything, providerID, mock.AnythingOfType("oidc.UpdateProviderRequest")).
			Return(nil, errors.New("database error"))

		newDisplayName := "Updated Provider"
		req := oidc.UpdateProviderRequest{DisplayName: &newDisplayName}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)

		assert.Nil(t, provider)
		assert.Error(t, err)
	})
}

func TestOIDCService_DisableProvider_Extended_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("states cleanup error ignored", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("DeleteStatesByProvider", mock.Anything, providerID).Return(errors.New("cleanup error"))
		mockRepo.On("DisableProvider", mock.Anything, providerID).Return(nil)

		err := svc.DisableProvider(context.Background(), providerID)

		assert.NoError(t, err) // States cleanup error should be logged but not returned
	})

	t.Run("disable error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("DeleteStatesByProvider", mock.Anything, providerID).Return(nil)
		mockRepo.On("DisableProvider", mock.Anything, providerID).Return(errors.New("disable failed"))

		err := svc.DisableProvider(context.Background(), providerID)

		assert.Error(t, err)
	})
}

func TestOIDCService_DeleteProvider_Extended_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("delete repo error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("DeleteStatesByProvider", mock.Anything, providerID).Return(nil)
		mockRepo.On("DeleteProvider", mock.Anything, providerID).Return(oidc.ErrProviderNotFound)

		err := svc.DeleteProvider(context.Background(), providerID)

		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})
}

func TestOIDCService_GetAuthURL_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("provider disabled", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		provider := makeTestProvider(uuid.New(), "test-provider", false) // Disabled

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(provider, nil)

		result, err := svc.GetAuthURL(context.Background(), "test-provider", "https://app.example.com", nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, oidc.ErrProviderDisabled)
	})

	t.Run("provider not found", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("GetProviderByName", mock.Anything, "unknown-provider").
			Return(nil, oidc.ErrProviderNotFound)

		result, err := svc.GetAuthURL(context.Background(), "unknown-provider", "https://app.example.com", nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})

	t.Run("state creation error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		provider := makeTestProvider(uuid.New(), "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(provider, nil)
		mockRepo.On("CreateState", mock.Anything, mock.AnythingOfType("oidc.CreateStateRequest")).
			Return(nil, errors.New("database error"))

		result, err := svc.GetAuthURL(context.Background(), "test-provider", "https://app.example.com", nil)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(provider, nil)
		mockRepo.On("CreateState", mock.Anything, mock.AnythingOfType("oidc.CreateStateRequest")).
			Return(&oidc.State{
				ID:         uuid.New(),
				State:      "random-state",
				ProviderID: providerID,
			}, nil)

		result, err := svc.GetAuthURL(context.Background(), "test-provider", "https://app.example.com", nil)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.URL)
		assert.NotEmpty(t, result.State)
	})

	t.Run("success with user ID for linking", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		userID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(provider, nil)
		mockRepo.On("CreateState", mock.Anything, mock.MatchedBy(func(req oidc.CreateStateRequest) bool {
			return req.UserID != nil && *req.UserID == userID
		})).Return(&oidc.State{
			ID:         uuid.New(),
			State:      "random-state",
			ProviderID: providerID,
			UserID:     &userID,
		}, nil)

		result, err := svc.GetAuthURL(context.Background(), "test-provider", "https://app.example.com", &userID)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestOIDCService_WithEncryption_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	// Create a 32-byte encryption key for AES-256
	encryptKey := []byte("12345678901234567890123456789012")

	setupEncryptedService := func(repo oidc.Repository) *oidc.Service {
		logger := zap.NewNop()
		return oidc.NewService(repo, logger, "https://app.example.com/callback", encryptKey)
	}

	t.Run("AddProvider encrypts client secret", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupEncryptedService(mockRepo)

		providerID := uuid.New()
		created := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(nil, oidc.ErrProviderNotFound)
		mockRepo.On("CreateProvider", mock.Anything, mock.MatchedBy(func(req oidc.CreateProviderRequest) bool {
			// Secret should be encrypted (different from original)
			return len(req.ClientSecretEncrypted) > 0 &&
				string(req.ClientSecretEncrypted) != "my-secret"
		})).Return(created, nil)

		req := oidc.CreateProviderRequest{
			Name:                  "test-provider",
			ProviderType:         oidc.ProviderTypeGeneric,
			IssuerURL:            "https://auth.example.com",
			ClientID:             "test-client-id",
			ClientSecretEncrypted: []byte("my-secret"),
		}

		provider, err := svc.AddProvider(context.Background(), req)

		require.NoError(t, err)
		assert.NotNil(t, provider)
	})

	t.Run("UpdateProvider encrypts new client secret", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupEncryptedService(mockRepo)

		providerID := uuid.New()
		updated := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("UpdateProvider", mock.Anything, providerID, mock.MatchedBy(func(req oidc.UpdateProviderRequest) bool {
			// Secret should be encrypted
			return len(req.ClientSecretEncrypted) > 0 &&
				string(req.ClientSecretEncrypted) != "new-secret"
		})).Return(updated, nil)

		req := oidc.UpdateProviderRequest{
			ClientSecretEncrypted: []byte("new-secret"),
		}

		provider, err := svc.UpdateProvider(context.Background(), providerID, req)

		require.NoError(t, err)
		assert.NotNil(t, provider)
	})
}

func TestOIDCService_EnableProvider_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("enable error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("EnableProvider", mock.Anything, providerID).Return(oidc.ErrProviderNotFound)

		err := svc.EnableProvider(context.Background(), providerID)

		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})
}

func TestOIDCService_SetDefaultProvider_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("set default error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()

		mockRepo.On("SetDefaultProvider", mock.Anything, providerID).Return(oidc.ErrProviderNotFound)

		err := svc.SetDefaultProvider(context.Background(), providerID)

		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})
}

func TestOIDCService_GetProviderByName_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("not found", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("GetProviderByName", mock.Anything, "unknown").Return(nil, oidc.ErrProviderNotFound)

		provider, err := svc.GetProviderByName(context.Background(), "unknown")

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})
}

func TestOIDCService_GetDefaultProvider_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("no default set", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("GetDefaultProvider", mock.Anything).Return(nil, oidc.ErrProviderNotFound)

		provider, err := svc.GetDefaultProvider(context.Background())

		assert.Nil(t, provider)
		assert.ErrorIs(t, err, oidc.ErrProviderNotFound)
	})
}

func TestOIDCService_ListProviders_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("database error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("ListProviders", mock.Anything).Return(nil, errors.New("database error"))

		providers, err := svc.ListProviders(context.Background())

		assert.Nil(t, providers)
		assert.Error(t, err)
	})
}

func TestOIDCService_ListEnabledProviders_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("database error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		mockRepo.On("ListEnabledProviders", mock.Anything).Return(nil, errors.New("database error"))

		providers, err := svc.ListEnabledProviders(context.Background())

		assert.Nil(t, providers)
		assert.Error(t, err)
	})
}

func TestOIDCService_ListUserLinks_Error_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("database error", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()

		mockRepo.On("ListUserLinks", mock.Anything, userID).Return(nil, errors.New("database error"))

		links, err := svc.ListUserLinks(context.Background(), userID)

		assert.Nil(t, links)
		assert.Error(t, err)
	})
}

func TestOIDCService_LinkUser_WithToken_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("success with OAuth2 token", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		userID := uuid.New()
		providerID := uuid.New()
		linkID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)

		mockRepo.On("GetProvider", mock.Anything, providerID).Return(provider, nil)
		mockRepo.On("GetUserLinkByUserAndProvider", mock.Anything, userID, providerID).
			Return(nil, oidc.ErrUserLinkNotFound)

		createdLink := makeTestUserLink(linkID, userID, providerID, "subject123")
		mockRepo.On("CreateUserLink", mock.Anything, mock.MatchedBy(func(req oidc.CreateUserLinkRequest) bool {
			return req.UserID == userID && req.ProviderID == providerID && req.Subject == "subject123"
		})).Return(createdLink, nil)

		userInfo := &oidc.UserInfo{
			Subject:  "subject123",
			Email:    "user@example.com",
			Name:     "Test User",
			Username: "testuser",
		}

		// OAuth2 token - without encryption configured, tokens are stored as-is
		token := &oauth2.Token{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
			Expiry:       time.Now().Add(time.Hour),
		}

		link, err := svc.LinkUser(context.Background(), userID, providerID, "subject123", userInfo, token)

		require.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, linkID, link.ID)
	})
}

func TestOIDC_MarshalHelpers_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("MarshalClaimMappings", func(t *testing.T) {
		cm := oidc.ClaimMappings{
			Username: "preferred_username",
			Email:    "email",
			Name:     "name",
			Picture:  "picture",
			Roles:    "roles",
		}

		data, err := oidc.MarshalClaimMappings(cm)

		require.NoError(t, err)
		assert.Contains(t, string(data), "preferred_username")
	})

	t.Run("UnmarshalClaimMappings", func(t *testing.T) {
		data := []byte(`{"username":"sub","email":"mail","name":"displayName","picture":"photo"}`)

		cm, err := oidc.UnmarshalClaimMappings(data)

		require.NoError(t, err)
		assert.Equal(t, "sub", cm.Username)
		assert.Equal(t, "mail", cm.Email)
	})

	t.Run("UnmarshalClaimMappings invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid}`)

		_, err := oidc.UnmarshalClaimMappings(data)

		assert.Error(t, err)
	})

	t.Run("MarshalRoleMappings", func(t *testing.T) {
		rm := map[string]string{
			"admin": "administrator",
			"user":  "member",
		}

		data, err := oidc.MarshalRoleMappings(rm)

		require.NoError(t, err)
		assert.Contains(t, string(data), "admin")
	})

	t.Run("MarshalRoleMappings nil", func(t *testing.T) {
		data, err := oidc.MarshalRoleMappings(nil)

		require.NoError(t, err)
		assert.Equal(t, "{}", string(data))
	})

	t.Run("UnmarshalRoleMappings", func(t *testing.T) {
		data := []byte(`{"admin":"administrator","user":"member"}`)

		rm, err := oidc.UnmarshalRoleMappings(data)

		require.NoError(t, err)
		assert.Equal(t, "administrator", rm["admin"])
	})

	t.Run("UnmarshalRoleMappings invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid}`)

		_, err := oidc.UnmarshalRoleMappings(data)

		assert.Error(t, err)
	})
}

func TestOIDCService_GetAuthURL_CustomEndpoints_Short(t *testing.T) {
	if testing.Short() {
		t.Log("Running short test")
	}

	t.Run("provider with custom endpoints", func(t *testing.T) {
		mockRepo := NewMockOIDCRepository(t)
		svc := setupOIDCService(mockRepo)

		providerID := uuid.New()
		provider := makeTestProvider(providerID, "test-provider", true)
		// Set custom endpoints
		authEndpoint := "https://auth.example.com/custom/authorize"
		tokenEndpoint := "https://auth.example.com/custom/token"
		provider.AuthorizationEndpoint = &authEndpoint
		provider.TokenEndpoint = &tokenEndpoint

		mockRepo.On("GetProviderByName", mock.Anything, "test-provider").Return(provider, nil)
		mockRepo.On("CreateState", mock.Anything, mock.AnythingOfType("oidc.CreateStateRequest")).
			Return(&oidc.State{
				ID:         uuid.New(),
				State:      "random-state",
				ProviderID: providerID,
			}, nil)

		result, err := svc.GetAuthURL(context.Background(), "test-provider", "https://app.example.com", nil)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.URL, "custom/authorize")
	})
}
