// Package oidc provides OIDC/SSO authentication services.
package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/user"
)

var (
	// ErrProviderNotFound indicates the OIDC provider was not found.
	ErrProviderNotFound = errors.New("provider not found")
	// ErrProviderDisabled indicates the OIDC provider is disabled.
	ErrProviderDisabled = errors.New("provider disabled")
	// ErrLinkNotFound indicates the user link was not found.
	ErrLinkNotFound = errors.New("link not found")
	// ErrAutoProvisionDisabled indicates auto-provisioning is disabled.
	ErrAutoProvisionDisabled = errors.New("auto-provisioning disabled")
	// ErrSlugTaken indicates the provider slug is already in use.
	ErrSlugTaken = errors.New("provider slug already exists")
)

// ClaimMapping defines how OIDC claims map to user attributes.
type ClaimMapping struct {
	Username string `json:"username"` // Claim to use for username (default: preferred_username)
	Email    string `json:"email"`    // Claim to use for email (default: email)
	Name     string `json:"name"`     // Claim to use for name (default: name)
	Groups   string `json:"groups"`   // Claim to use for groups (default: groups)
}

// RoleMapping defines how OIDC groups map to roles.
type RoleMapping struct {
	AdminGroups []string `json:"admin_groups"` // Groups that grant admin access
}

// DefaultClaimMapping returns the default claim mapping.
func DefaultClaimMapping() ClaimMapping {
	return ClaimMapping{
		Username: "preferred_username",
		Email:    "email",
		Name:     "name",
		Groups:   "groups",
	}
}

// Service provides OIDC provider and user link management.
type Service struct {
	queries        *db.Queries
	userService    *user.Service
	sessionService *session.Service
	logger         *slog.Logger
}

// NewService creates a new OIDC service.
func NewService(
	queries *db.Queries,
	userService *user.Service,
	sessionService *session.Service,
	logger *slog.Logger,
) *Service {
	return &Service{
		queries:        queries,
		userService:    userService,
		sessionService: sessionService,
		logger:         logger.With(slog.String("service", "oidc")),
	}
}

// CreateProviderParams contains parameters for creating an OIDC provider.
type CreateProviderParams struct {
	Name            string
	Slug            string
	Enabled         bool
	IssuerURL       string
	ClientID        string
	ClientSecretEnc []byte // Encrypted client secret
	Scopes          []string
	ClaimMapping    *ClaimMapping
	RoleMapping     *RoleMapping
	AutoProvision   bool
	DefaultRole     string // user, admin
}

// CreateProvider creates a new OIDC provider.
func (s *Service) CreateProvider(ctx context.Context, params CreateProviderParams) (*db.OidcProvider, error) {
	// Check for existing slug
	_, err := s.queries.GetOIDCProviderBySlug(ctx, params.Slug)
	if err == nil {
		return nil, ErrSlugTaken
	}

	// Set default scopes
	scopes := params.Scopes
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"}
	}

	// Set default claim mapping
	claimMapping := DefaultClaimMapping()
	if params.ClaimMapping != nil {
		claimMapping = *params.ClaimMapping
	}
	claimJSON, _ := json.Marshal(claimMapping)

	// Set default role mapping
	roleMapping := RoleMapping{}
	if params.RoleMapping != nil {
		roleMapping = *params.RoleMapping
	}
	roleJSON, _ := json.Marshal(roleMapping)

	// Default role
	defaultRole := params.DefaultRole
	if defaultRole == "" {
		defaultRole = "user"
	}

	provider, err := s.queries.CreateOIDCProvider(ctx, db.CreateOIDCProviderParams{
		Name:            params.Name,
		Slug:            params.Slug,
		Enabled:         params.Enabled,
		IssuerUrl:       params.IssuerURL,
		ClientID:        params.ClientID,
		ClientSecretEnc: params.ClientSecretEnc,
		Scopes:          scopes,
		ClaimMapping:    claimJSON,
		RoleMapping:     roleJSON,
		AutoProvision:   params.AutoProvision,
		DefaultRole:     defaultRole,
	})
	if err != nil {
		return nil, fmt.Errorf("create provider: %w", err)
	}

	s.logger.Info("OIDC provider created",
		slog.String("provider_id", provider.ID.String()),
		slog.String("name", provider.Name),
		slog.String("slug", provider.Slug),
	)

	return &provider, nil
}

// GetProviderByID retrieves an OIDC provider by ID.
func (s *Service) GetProviderByID(ctx context.Context, id uuid.UUID) (*db.OidcProvider, error) {
	provider, err := s.queries.GetOIDCProviderByID(ctx, id)
	if err != nil {
		return nil, ErrProviderNotFound
	}
	return &provider, nil
}

// GetProviderBySlug retrieves an OIDC provider by slug.
func (s *Service) GetProviderBySlug(ctx context.Context, slug string) (*db.OidcProvider, error) {
	provider, err := s.queries.GetOIDCProviderBySlug(ctx, slug)
	if err != nil {
		return nil, ErrProviderNotFound
	}
	return &provider, nil
}

// GetEnabledProvider retrieves an OIDC provider by slug, checking if enabled.
func (s *Service) GetEnabledProvider(ctx context.Context, slug string) (*db.OidcProvider, error) {
	provider, err := s.GetProviderBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if !provider.Enabled {
		return nil, ErrProviderDisabled
	}
	return provider, nil
}

// ListProviders returns all OIDC providers (admin only).
func (s *Service) ListProviders(ctx context.Context) ([]db.OidcProvider, error) {
	providers, err := s.queries.ListOIDCProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("list providers: %w", err)
	}
	return providers, nil
}

// ListEnabledProviders returns all enabled OIDC providers (for login page).
func (s *Service) ListEnabledProviders(ctx context.Context) ([]db.OidcProvider, error) {
	providers, err := s.queries.ListEnabledOIDCProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("list enabled providers: %w", err)
	}
	return providers, nil
}

// DeleteProvider deletes an OIDC provider and all associated user links.
func (s *Service) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	// Delete all user links first
	if err := s.queries.DeleteOIDCLinksByProvider(ctx, id); err != nil {
		return fmt.Errorf("delete links: %w", err)
	}

	// Delete provider
	if err := s.queries.DeleteOIDCProvider(ctx, id); err != nil {
		return fmt.Errorf("delete provider: %w", err)
	}

	s.logger.Info("OIDC provider deleted",
		slog.String("provider_id", id.String()),
	)

	return nil
}

// GetUserLink retrieves an OIDC user link by provider and subject.
func (s *Service) GetUserLink(ctx context.Context, providerID uuid.UUID, subject string) (*db.OidcUserLink, error) {
	link, err := s.queries.GetOIDCLinkBySubject(ctx, db.GetOIDCLinkBySubjectParams{
		ProviderID: providerID,
		Subject:    subject,
	})
	if err != nil {
		return nil, ErrLinkNotFound
	}
	return &link, nil
}

// ListUserLinks returns all OIDC links for a user.
func (s *Service) ListUserLinks(ctx context.Context, userID uuid.UUID) ([]db.ListOIDCLinksByUserRow, error) {
	links, err := s.queries.ListOIDCLinksByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list links: %w", err)
	}
	return links, nil
}

// CreateLink creates a new OIDC user link.
func (s *Service) CreateLink(ctx context.Context, userID, providerID uuid.UUID, subject string, email, name *string, groups []string) (*db.OidcUserLink, error) {
	link, err := s.queries.CreateOIDCLink(ctx, db.CreateOIDCLinkParams{
		UserID:     userID,
		ProviderID: providerID,
		Subject:    subject,
		Email:      email,
		Name:       name,
		Groups:     groups,
	})
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	s.logger.Info("OIDC user link created",
		slog.String("user_id", userID.String()),
		slog.String("provider_id", providerID.String()),
		slog.String("subject", subject),
	)

	return &link, nil
}

// UpdateLinkLogin updates the link with latest login info.
func (s *Service) UpdateLinkLogin(ctx context.Context, linkID uuid.UUID, email, name *string, groups []string) error {
	if err := s.queries.UpdateOIDCLinkLogin(ctx, db.UpdateOIDCLinkLoginParams{
		ID:     linkID,
		Email:  email,
		Name:   name,
		Groups: groups,
	}); err != nil {
		return fmt.Errorf("update link login: %w", err)
	}
	return nil
}

// DeleteLink deletes an OIDC user link.
func (s *Service) DeleteLink(ctx context.Context, linkID uuid.UUID) error {
	if err := s.queries.DeleteOIDCLink(ctx, linkID); err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	s.logger.Info("OIDC user link deleted",
		slog.String("link_id", linkID.String()),
	)

	return nil
}

// DeleteUserLinks deletes all OIDC links for a user.
func (s *Service) DeleteUserLinks(ctx context.Context, userID uuid.UUID) error {
	if err := s.queries.DeleteOIDCLinksByUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user links: %w", err)
	}

	s.logger.Info("All OIDC links deleted for user",
		slog.String("user_id", userID.String()),
	)

	return nil
}

// ParseClaimMapping parses the claim mapping JSON from a provider.
func (s *Service) ParseClaimMapping(provider *db.OidcProvider) (ClaimMapping, error) {
	var mapping ClaimMapping
	if err := json.Unmarshal(provider.ClaimMapping, &mapping); err != nil {
		return DefaultClaimMapping(), nil
	}
	return mapping, nil
}

// ParseRoleMapping parses the role mapping JSON from a provider.
func (s *Service) ParseRoleMapping(provider *db.OidcProvider) (RoleMapping, error) {
	var mapping RoleMapping
	if err := json.Unmarshal(provider.RoleMapping, &mapping); err != nil {
		return RoleMapping{}, nil
	}
	return mapping, nil
}

// IsAdminFromGroups checks if any of the user's groups grant admin access.
func (s *Service) IsAdminFromGroups(roleMapping RoleMapping, groups []string) bool {
	for _, group := range groups {
		for _, adminGroup := range roleMapping.AdminGroups {
			if strings.EqualFold(group, adminGroup) {
				return true
			}
		}
	}
	return false
}

// ProviderInfo contains public info about a provider (for login page).
type ProviderInfo struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// GetPublicProviders returns public info for all enabled providers.
func (s *Service) GetPublicProviders(ctx context.Context) ([]ProviderInfo, error) {
	providers, err := s.ListEnabledProviders(ctx)
	if err != nil {
		return nil, err
	}

	infos := make([]ProviderInfo, len(providers))
	for i, p := range providers {
		infos[i] = ProviderInfo{
			Name: p.Name,
			Slug: p.Slug,
		}
	}
	return infos, nil
}
