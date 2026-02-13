package api

import (
	"context"
	"errors"
	"strings"

	"log/slog"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/oidc"
)

// ============================================================================
// Public OIDC Endpoints
// ============================================================================

// ListOIDCProviders returns enabled OIDC providers for login
func (h *Handler) ListOIDCProviders(ctx context.Context) (*ogen.OIDCProviderListResponse, error) {
	providers, err := h.oidcService.ListEnabledProviders(ctx)
	if err != nil {
		h.logger.Error("failed to list OIDC providers", slog.Any("error", err))
		return &ogen.OIDCProviderListResponse{
			Providers: []ogen.OIDCProviderInfo{},
			Total:     0,
		}, nil
	}

	infos := make([]ogen.OIDCProviderInfo, len(providers))
	for i, p := range providers {
		infos[i] = ogen.OIDCProviderInfo{
			Name:        p.Name,
			DisplayName: p.DisplayName,
			IsDefault:   ogen.NewOptBool(p.IsDefault),
		}
	}

	return &ogen.OIDCProviderListResponse{
		Providers: infos,
		Total:     int64(len(infos)),
	}, nil
}

// OidcAuthorize initiates the OIDC login flow.
// Returns the OIDC auth URL as JSON for client-side redirect (SPA pattern).
func (h *Handler) OidcAuthorize(ctx context.Context, params ogen.OidcAuthorizeParams) (ogen.OidcAuthorizeRes, error) {
	redirectURL := ""
	if params.RedirectURL.IsSet() {
		redirectURL = params.RedirectURL.Value
	}

	result, err := h.oidcService.GetAuthURL(ctx, params.Provider, redirectURL, nil)
	if err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.Error{
				Code:    404,
				Message: "Provider not found",
			}, nil
		}
		h.logger.Error("failed to get auth URL", slog.String("provider", params.Provider), slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to generate auth URL",
		}, nil
	}

	return &ogen.OIDCAuthURLResponse{
		AuthURL: result.URL,
	}, nil
}

// OidcCallback handles the OAuth2 callback
func (h *Handler) OidcCallback(ctx context.Context, params ogen.OidcCallbackParams) (ogen.OidcCallbackRes, error) {
	// Check for error from provider
	if params.Error.IsSet() {
		h.logger.Warn("OIDC provider returned error",
			slog.String("provider", params.Provider),
			slog.String("error", params.Error.Value),
			slog.String("description", params.ErrorDescription.Or("")),
		)
		return &ogen.Error{
			Code:    400,
			Message: params.ErrorDescription.Or("Authentication failed"),
		}, nil
	}

	result, err := h.oidcService.HandleCallback(ctx, params.State, params.Code)
	if err != nil {
		h.logger.Error("OIDC callback failed",
			slog.String("provider", params.Provider),
			slog.Any("error", err),
		)

		if errors.Is(err, oidc.ErrInvalidState) || errors.Is(err, oidc.ErrStateExpired) {
			return &ogen.Error{
				Code:    400,
				Message: "Invalid or expired authentication state",
			}, nil
		}
		if errors.Is(err, oidc.ErrAutoCreateDisabled) {
			return &ogen.Error{
				Code:    400,
				Message: "Account creation is disabled. Please register first or contact admin.",
			}, nil
		}

		return &ogen.Error{
			Code:    500,
			Message: "Authentication failed",
		}, nil
	}

	var userID uuid.UUID

	// Handle new user creation from OIDC
	if result.IsNewUser {
		// Generate username from OIDC user info
		username := result.UserInfo.Username
		if username == "" {
			// Fallback to email prefix if no username provided
			username = result.UserInfo.Email
			if idx := strings.Index(username, "@"); idx > 0 {
				username = username[:idx]
			}
		}

		// Create user from OIDC data
		displayName := result.UserInfo.Name
		user, err := h.authService.RegisterFromOIDC(ctx, auth.RegisterFromOIDCRequest{
			Username:    username,
			Email:       result.UserInfo.Email,
			DisplayName: &displayName,
		})
		if err != nil {
			h.logger.Error("failed to create user from OIDC",
				slog.String("provider", params.Provider),
				slog.String("email", result.UserInfo.Email),
				slog.Any("error", err),
			)
			return &ogen.Error{
				Code:    500,
				Message: "Failed to create user account",
			}, nil
		}

		userID = user.ID

		// Link user to OIDC provider
		_, err = h.oidcService.LinkUser(ctx, userID, result.ProviderID, result.UserInfo.Subject, result.UserInfo, nil)
		if err != nil {
			h.logger.Error("failed to link OIDC user",
				slog.String("provider", params.Provider),
				slog.String("user_id", userID.String()),
				slog.Any("error", err),
			)
			// User was created but link failed - log but continue
			// User can still use the account and link later
		}

		h.logger.Info("new user created from OIDC",
			slog.String("provider", params.Provider),
			slog.String("user_id", userID.String()),
			slog.String("username", username),
		)
	} else {
		userID = result.UserID
	}

	// Create session for the user
	loginResp, err := h.authService.CreateSessionForUser(ctx, userID, nil, nil, nil)
	if err != nil {
		h.logger.Error("failed to create session for OIDC user",
			slog.String("provider", params.Provider),
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to create session",
		}, nil
	}

	return &ogen.OIDCCallbackResponse{
		AccessToken:  loginResp.AccessToken,
		RefreshToken: ogen.NewOptString(loginResp.RefreshToken),
		TokenType:    "Bearer",
		ExpiresIn:    int(loginResp.ExpiresIn),
	}, nil
}

// ============================================================================
// User OIDC Link Endpoints
// ============================================================================

// ListUserOIDCLinks lists user's linked OIDC providers
func (h *Handler) ListUserOIDCLinks(ctx context.Context) (ogen.ListUserOIDCLinksRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Unauthorized",
		}, nil
	}

	links, err := h.oidcService.ListUserLinks(ctx, userID)
	if err != nil {
		h.logger.Error("failed to list user OIDC links", slog.String("user_id", userID.String()), slog.Any("error", err))
		return &ogen.Error{
			Code:    500,
			Message: "Failed to list OIDC links",
		}, nil
	}

	ogenLinks := make([]ogen.OIDCUserLink, len(links))
	for i, l := range links {
		link := ogen.OIDCUserLink{
			ID:                  l.ID,
			ProviderName:        l.ProviderName,
			ProviderDisplayName: l.ProviderDisplayName,
			LinkedAt:            l.CreatedAt,
		}
		if l.Email != nil {
			link.Email.SetTo(*l.Email)
		}
		if l.LastLoginAt != nil {
			link.LastLoginAt.SetTo(*l.LastLoginAt)
		}
		ogenLinks[i] = link
	}

	return &ogen.OIDCUserLinkListResponse{
		Links: ogenLinks,
		Total: int64(len(ogenLinks)),
	}, nil
}

// InitOIDCLink initiates linking an OIDC provider to user's account
func (h *Handler) InitOIDCLink(ctx context.Context, params ogen.InitOIDCLinkParams) (ogen.InitOIDCLinkRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.InitOIDCLinkUnauthorized{}, nil
	}

	result, err := h.oidcService.GetAuthURL(ctx, params.Provider, "", &userID)
	if err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.InitOIDCLinkNotFound{}, nil
		}
		if errors.Is(err, oidc.ErrLinkingDisabled) {
			return &ogen.InitOIDCLinkNotFound{}, nil
		}
		h.logger.Error("failed to init OIDC link",
			slog.String("user_id", userID.String()),
			slog.String("provider", params.Provider),
			slog.Any("error", err),
		)
		return &ogen.InitOIDCLinkNotFound{}, nil
	}

	return &ogen.OIDCAuthURLResponse{
		AuthURL: result.URL,
	}, nil
}

// UnlinkOIDCProvider unlinks an OIDC provider from user's account
func (h *Handler) UnlinkOIDCProvider(ctx context.Context, params ogen.UnlinkOIDCProviderParams) (ogen.UnlinkOIDCProviderRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.UnlinkOIDCProviderUnauthorized{}, nil
	}

	provider, err := h.oidcService.GetProviderByName(ctx, params.Provider)
	if err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.UnlinkOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to get provider",
			slog.String("provider", params.Provider),
			slog.Any("error", err),
		)
		return &ogen.UnlinkOIDCProviderNotFound{}, nil
	}

	if err := h.oidcService.UnlinkUser(ctx, userID, provider.ID); err != nil {
		h.logger.Error("failed to unlink OIDC provider",
			slog.String("user_id", userID.String()),
			slog.String("provider", params.Provider),
			slog.Any("error", err),
		)
		return &ogen.UnlinkOIDCProviderNotFound{}, nil
	}

	return &ogen.UnlinkOIDCProviderNoContent{}, nil
}

// ============================================================================
// Admin OIDC Endpoints
// ============================================================================

// AdminListOIDCProviders lists all OIDC providers (admin only)
func (h *Handler) AdminListOIDCProviders(ctx context.Context) (ogen.AdminListOIDCProvidersRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminListOIDCProvidersUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminListOIDCProvidersForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	providers, err := h.oidcService.ListProviders(ctx)
	if err != nil {
		h.logger.Error("failed to list OIDC providers", slog.Any("error", err))
		return &ogen.AdminListOIDCProvidersForbidden{}, nil
	}

	ogenProviders := make([]ogen.AdminOIDCProvider, len(providers))
	for i, p := range providers {
		ogenProviders[i] = providerToOgen(p)
	}

	return &ogen.AdminOIDCProviderListResponse{
		Providers: ogenProviders,
		Total:     int64(len(ogenProviders)),
	}, nil
}

// AdminCreateOIDCProvider creates a new OIDC provider (admin only)
func (h *Handler) AdminCreateOIDCProvider(ctx context.Context, req *ogen.CreateOIDCProviderRequest) (ogen.AdminCreateOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminCreateOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminCreateOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	createReq := oidc.CreateProviderRequest{
		Name:                  req.Name,
		DisplayName:           req.DisplayName,
		ProviderType:          string(req.ProviderType.Or(ogen.CreateOIDCProviderRequestProviderTypeGeneric)),
		IssuerURL:             req.IssuerURL,
		ClientID:              req.ClientID,
		ClientSecretEncrypted: []byte(req.ClientSecret),
		AutoCreateUsers:       req.AutoCreateUsers.Or(true),
		UpdateUserInfo:        req.UpdateUserInfo.Or(true),
		AllowLinking:          req.AllowLinking.Or(true),
		IsEnabled:             req.IsEnabled.Or(true),
		IsDefault:             req.IsDefault.Or(false),
	}

	if req.AuthorizationEndpoint.IsSet() {
		v := req.AuthorizationEndpoint.Value
		createReq.AuthorizationEndpoint = &v
	}
	if req.TokenEndpoint.IsSet() {
		v := req.TokenEndpoint.Value
		createReq.TokenEndpoint = &v
	}
	if req.UserInfoEndpoint.IsSet() {
		v := req.UserInfoEndpoint.Value
		createReq.UserInfoEndpoint = &v
	}
	if req.JwksURI.IsSet() {
		v := req.JwksURI.Value
		createReq.JWKSURI = &v
	}
	if req.EndSessionEndpoint.IsSet() {
		v := req.EndSessionEndpoint.Value
		createReq.EndSessionEndpoint = &v
	}

	if len(req.Scopes) > 0 {
		createReq.Scopes = req.Scopes
	}

	if req.ClaimMappings.IsSet() {
		createReq.ClaimMappings = claimMappingsFromOgen(req.ClaimMappings.Value)
	}

	if req.RoleMappings.IsSet() {
		createReq.RoleMappings = req.RoleMappings.Value
	}

	provider, err := h.oidcService.AddProvider(ctx, createReq)
	if err != nil {
		h.logger.Error("failed to create OIDC provider", slog.String("name", req.Name), slog.Any("error", err))

		if errors.Is(err, oidc.ErrProviderNameExists) {
			return &ogen.AdminCreateOIDCProviderConflict{
				Code:    409,
				Message: "Provider with this name already exists",
			}, nil
		}
		if errors.Is(err, oidc.ErrInvalidProviderType) {
			return &ogen.AdminCreateOIDCProviderBadRequest{
				Code:    400,
				Message: "Invalid provider type",
			}, nil
		}
		if errors.Is(err, oidc.ErrInvalidIssuerURL) {
			return &ogen.AdminCreateOIDCProviderBadRequest{
				Code:    400,
				Message: "Invalid issuer URL",
			}, nil
		}

		return &ogen.AdminCreateOIDCProviderBadRequest{
			Code:    500,
			Message: "Failed to create provider",
		}, nil
	}

	result := providerToOgen(*provider)
	return &result, nil
}

// AdminGetOIDCProvider gets a provider by ID (admin only)
func (h *Handler) AdminGetOIDCProvider(ctx context.Context, params ogen.AdminGetOIDCProviderParams) (ogen.AdminGetOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminGetOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminGetOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	provider, err := h.oidcService.GetProvider(ctx, params.ProviderId)
	if err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.AdminGetOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to get OIDC provider", slog.String("id", params.ProviderId.String()), slog.Any("error", err))
		return &ogen.AdminGetOIDCProviderNotFound{}, nil
	}

	result := providerToOgen(*provider)
	return &result, nil
}

// AdminUpdateOIDCProvider updates a provider (admin only)
func (h *Handler) AdminUpdateOIDCProvider(ctx context.Context, req *ogen.UpdateOIDCProviderRequest, params ogen.AdminUpdateOIDCProviderParams) (ogen.AdminUpdateOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminUpdateOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminUpdateOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	updateReq := oidc.UpdateProviderRequest{}

	if req.DisplayName.IsSet() {
		v := req.DisplayName.Value
		updateReq.DisplayName = &v
	}
	if req.ProviderType.IsSet() {
		v := string(req.ProviderType.Value)
		updateReq.ProviderType = &v
	}
	if req.IssuerURL.IsSet() {
		v := req.IssuerURL.Value
		updateReq.IssuerURL = &v
	}
	if req.ClientID.IsSet() {
		v := req.ClientID.Value
		updateReq.ClientID = &v
	}
	if req.ClientSecret.IsSet() {
		updateReq.ClientSecretEncrypted = []byte(req.ClientSecret.Value)
	}
	if req.AuthorizationEndpoint.IsSet() {
		v := req.AuthorizationEndpoint.Value
		updateReq.AuthorizationEndpoint = &v
	}
	if req.TokenEndpoint.IsSet() {
		v := req.TokenEndpoint.Value
		updateReq.TokenEndpoint = &v
	}
	if req.UserInfoEndpoint.IsSet() {
		v := req.UserInfoEndpoint.Value
		updateReq.UserInfoEndpoint = &v
	}
	if req.JwksURI.IsSet() {
		v := req.JwksURI.Value
		updateReq.JWKSURI = &v
	}
	if req.EndSessionEndpoint.IsSet() {
		v := req.EndSessionEndpoint.Value
		updateReq.EndSessionEndpoint = &v
	}
	if len(req.Scopes) > 0 {
		updateReq.Scopes = req.Scopes
	}
	if req.ClaimMappings.IsSet() {
		cm := claimMappingsFromOgen(req.ClaimMappings.Value)
		updateReq.ClaimMappings = &cm
	}
	if req.RoleMappings.IsSet() {
		updateReq.RoleMappings = req.RoleMappings.Value
	}
	if req.AutoCreateUsers.IsSet() {
		v := req.AutoCreateUsers.Value
		updateReq.AutoCreateUsers = &v
	}
	if req.UpdateUserInfo.IsSet() {
		v := req.UpdateUserInfo.Value
		updateReq.UpdateUserInfo = &v
	}
	if req.AllowLinking.IsSet() {
		v := req.AllowLinking.Value
		updateReq.AllowLinking = &v
	}
	if req.IsEnabled.IsSet() {
		v := req.IsEnabled.Value
		updateReq.IsEnabled = &v
	}
	if req.IsDefault.IsSet() {
		v := req.IsDefault.Value
		updateReq.IsDefault = &v
	}

	provider, err := h.oidcService.UpdateProvider(ctx, params.ProviderId, updateReq)
	if err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.AdminUpdateOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to update OIDC provider", slog.String("id", params.ProviderId.String()), slog.Any("error", err))
		return &ogen.AdminUpdateOIDCProviderBadRequest{
			Code:    500,
			Message: "Failed to update provider",
		}, nil
	}

	result := providerToOgen(*provider)
	return &result, nil
}

// AdminDeleteOIDCProvider deletes a provider (admin only)
func (h *Handler) AdminDeleteOIDCProvider(ctx context.Context, params ogen.AdminDeleteOIDCProviderParams) (ogen.AdminDeleteOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminDeleteOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminDeleteOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	if err := h.oidcService.DeleteProvider(ctx, params.ProviderId); err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.AdminDeleteOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to delete OIDC provider", slog.String("id", params.ProviderId.String()), slog.Any("error", err))
		return &ogen.AdminDeleteOIDCProviderNotFound{}, nil
	}

	return &ogen.AdminDeleteOIDCProviderNoContent{}, nil
}

// AdminEnableOIDCProvider enables a provider (admin only)
func (h *Handler) AdminEnableOIDCProvider(ctx context.Context, params ogen.AdminEnableOIDCProviderParams) (ogen.AdminEnableOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminEnableOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminEnableOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	if err := h.oidcService.EnableProvider(ctx, params.ProviderId); err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.AdminEnableOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to enable OIDC provider", slog.String("id", params.ProviderId.String()), slog.Any("error", err))
		return &ogen.AdminEnableOIDCProviderNotFound{}, nil
	}

	return &ogen.AdminEnableOIDCProviderNoContent{}, nil
}

// AdminDisableOIDCProvider disables a provider (admin only)
func (h *Handler) AdminDisableOIDCProvider(ctx context.Context, params ogen.AdminDisableOIDCProviderParams) (ogen.AdminDisableOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminDisableOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminDisableOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	if err := h.oidcService.DisableProvider(ctx, params.ProviderId); err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.AdminDisableOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to disable OIDC provider", slog.String("id", params.ProviderId.String()), slog.Any("error", err))
		return &ogen.AdminDisableOIDCProviderNotFound{}, nil
	}

	return &ogen.AdminDisableOIDCProviderNoContent{}, nil
}

// AdminSetDefaultOIDCProvider sets a provider as default (admin only)
func (h *Handler) AdminSetDefaultOIDCProvider(ctx context.Context, params ogen.AdminSetDefaultOIDCProviderParams) (ogen.AdminSetDefaultOIDCProviderRes, error) {
	if _, err := h.requireAdmin(ctx); err != nil {
		if errors.Is(err, errNotAuthenticated) {
			return &ogen.AdminSetDefaultOIDCProviderUnauthorized{Code: 401, Message: "Authentication required"}, nil
		}
		if errors.Is(err, errNotAdmin) {
			return &ogen.AdminSetDefaultOIDCProviderForbidden{Code: 403, Message: "Admin access required"}, nil
		}
		return nil, err
	}

	if err := h.oidcService.SetDefaultProvider(ctx, params.ProviderId); err != nil {
		if errors.Is(err, oidc.ErrProviderNotFound) {
			return &ogen.AdminSetDefaultOIDCProviderNotFound{}, nil
		}
		h.logger.Error("failed to set default OIDC provider", slog.String("id", params.ProviderId.String()), slog.Any("error", err))
		return &ogen.AdminSetDefaultOIDCProviderNotFound{}, nil
	}

	return &ogen.AdminSetDefaultOIDCProviderNoContent{}, nil
}

// ============================================================================
// Helpers
// ============================================================================

func providerToOgen(p oidc.Provider) ogen.AdminOIDCProvider {
	result := ogen.AdminOIDCProvider{
		ID:              p.ID,
		Name:            p.Name,
		DisplayName:     p.DisplayName,
		ProviderType:    ogen.AdminOIDCProviderProviderType(p.ProviderType),
		IssuerURL:       p.IssuerURL,
		ClientID:        p.ClientID,
		Scopes:          p.Scopes,
		AutoCreateUsers: p.AutoCreateUsers,
		UpdateUserInfo:  p.UpdateUserInfo,
		AllowLinking:    p.AllowLinking,
		IsEnabled:       p.IsEnabled,
		IsDefault:       p.IsDefault,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}

	if p.AuthorizationEndpoint != nil {
		result.AuthorizationEndpoint.SetTo(*p.AuthorizationEndpoint)
	}
	if p.TokenEndpoint != nil {
		result.TokenEndpoint.SetTo(*p.TokenEndpoint)
	}
	if p.UserInfoEndpoint != nil {
		result.UserInfoEndpoint.SetTo(*p.UserInfoEndpoint)
	}
	if p.JWKSURI != nil {
		result.JwksURI.SetTo(*p.JWKSURI)
	}
	if p.EndSessionEndpoint != nil {
		result.EndSessionEndpoint.SetTo(*p.EndSessionEndpoint)
	}

	result.ClaimMappings.SetTo(ogen.ClaimMappings{
		Username: ogen.NewOptString(p.ClaimMappings.Username),
		Email:    ogen.NewOptString(p.ClaimMappings.Email),
		Name:     ogen.NewOptString(p.ClaimMappings.Name),
		Picture:  ogen.NewOptString(p.ClaimMappings.Picture),
		Roles:    ogen.NewOptString(p.ClaimMappings.Roles),
	})

	if len(p.RoleMappings) > 0 {
		result.RoleMappings.SetTo(ogen.AdminOIDCProviderRoleMappings(p.RoleMappings))
	}

	return result
}

func claimMappingsFromOgen(cm ogen.ClaimMappings) oidc.ClaimMappings {
	return oidc.ClaimMappings{
		Username: cm.Username.Or("preferred_username"),
		Email:    cm.Email.Or("email"),
		Name:     cm.Name.Or("name"),
		Picture:  cm.Picture.Or("picture"),
		Roles:    cm.Roles.Or("groups"),
	}
}

// Ensure uuid is imported
var _ = uuid.UUID{}
