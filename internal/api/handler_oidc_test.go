package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/api/ogen"
)

// ============================================================================
// User OIDC Endpoints - Unauthorized (no user in context)
// ============================================================================

func TestHandler_OIDC_ListUserOIDCLinks_Unauthorized(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	result, err := handler.ListUserOIDCLinks(ctx)
	require.NoError(t, err)

	errRes, ok := result.(*ogen.Error)
	require.True(t, ok, "expected *ogen.Error, got %T", result)
	assert.Equal(t, 401, errRes.Code)
	assert.Equal(t, "Unauthorized", errRes.Message)
}

func TestHandler_OIDC_InitOIDCLink_Unauthorized(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.InitOIDCLinkParams{Provider: "google"}
	result, err := handler.InitOIDCLink(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.InitOIDCLinkUnauthorized)
	require.True(t, ok, "expected *ogen.InitOIDCLinkUnauthorized, got %T", result)
}

func TestHandler_OIDC_UnlinkOIDCProvider_Unauthorized(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.UnlinkOIDCProviderParams{Provider: "google"}
	result, err := handler.UnlinkOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.UnlinkOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.UnlinkOIDCProviderUnauthorized, got %T", result)
}

// ============================================================================
// Admin OIDC Endpoints - Forbidden (no user in context => isAdmin returns false)
// ============================================================================

func TestHandler_OIDC_AdminListOIDCProviders_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	result, err := handler.AdminListOIDCProviders(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminListOIDCProvidersUnauthorized)
	require.True(t, ok, "expected *ogen.AdminListOIDCProvidersUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminCreateOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.CreateOIDCProviderRequest{
		Name:         "test",
		DisplayName:  "Test",
		IssuerUrl:    "https://example.com",
		ClientId:     "client-id",
		ClientSecret: "secret",
	}
	result, err := handler.AdminCreateOIDCProvider(ctx, req)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminCreateOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminCreateOIDCProviderUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminGetOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.AdminGetOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminGetOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminGetOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminGetOIDCProviderUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminUpdateOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	req := &ogen.UpdateOIDCProviderRequest{}
	params := ogen.AdminUpdateOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminUpdateOIDCProvider(ctx, req, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminUpdateOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminUpdateOIDCProviderUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminDeleteOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.AdminDeleteOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminDeleteOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminDeleteOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminDeleteOIDCProviderUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminEnableOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.AdminEnableOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminEnableOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminEnableOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminEnableOIDCProviderUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminDisableOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.AdminDisableOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminDisableOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminDisableOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminDisableOIDCProviderUnauthorized, got %T", result)
}

func TestHandler_OIDC_AdminSetDefaultOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: logging.NewTestLogger(),
	}

	ctx := context.Background()
	params := ogen.AdminSetDefaultOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminSetDefaultOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminSetDefaultOIDCProviderUnauthorized)
	require.True(t, ok, "expected *ogen.AdminSetDefaultOIDCProviderUnauthorized, got %T", result)
}
