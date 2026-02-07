package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
)

// ============================================================================
// User OIDC Endpoints - Unauthorized (no user in context)
// ============================================================================

func TestHandler_OIDC_ListUserOIDCLinks_Unauthorized(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
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
		logger: zap.NewNop(),
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
		logger: zap.NewNop(),
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
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	result, err := handler.AdminListOIDCProviders(ctx)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminListOIDCProvidersForbidden)
	require.True(t, ok, "expected *ogen.AdminListOIDCProvidersForbidden, got %T", result)
}

func TestHandler_OIDC_AdminCreateOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
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

	_, ok := result.(*ogen.AdminCreateOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminCreateOIDCProviderForbidden, got %T", result)
}

func TestHandler_OIDC_AdminGetOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	params := ogen.AdminGetOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminGetOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminGetOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminGetOIDCProviderForbidden, got %T", result)
}

func TestHandler_OIDC_AdminUpdateOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	req := &ogen.UpdateOIDCProviderRequest{}
	params := ogen.AdminUpdateOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminUpdateOIDCProvider(ctx, req, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminUpdateOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminUpdateOIDCProviderForbidden, got %T", result)
}

func TestHandler_OIDC_AdminDeleteOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	params := ogen.AdminDeleteOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminDeleteOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminDeleteOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminDeleteOIDCProviderForbidden, got %T", result)
}

func TestHandler_OIDC_AdminEnableOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	params := ogen.AdminEnableOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminEnableOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminEnableOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminEnableOIDCProviderForbidden, got %T", result)
}

func TestHandler_OIDC_AdminDisableOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	params := ogen.AdminDisableOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminDisableOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminDisableOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminDisableOIDCProviderForbidden, got %T", result)
}

func TestHandler_OIDC_AdminSetDefaultOIDCProvider_Forbidden(t *testing.T) {
	t.Parallel()

	handler := &Handler{
		logger: zap.NewNop(),
	}

	ctx := context.Background()
	params := ogen.AdminSetDefaultOIDCProviderParams{ProviderId: uuid.New()}
	result, err := handler.AdminSetDefaultOIDCProvider(ctx, params)
	require.NoError(t, err)

	_, ok := result.(*ogen.AdminSetDefaultOIDCProviderForbidden)
	require.True(t, ok, "expected *ogen.AdminSetDefaultOIDCProviderForbidden, got %T", result)
}
