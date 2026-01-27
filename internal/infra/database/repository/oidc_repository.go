// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/revenge/internal/domain"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// OIDCProviderRepository implements domain.OIDCProviderRepository using PostgreSQL.
type OIDCProviderRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewOIDCProviderRepository creates a new PostgreSQL OIDC provider repository.
func NewOIDCProviderRepository(pool *pgxpool.Pool) *OIDCProviderRepository {
	return &OIDCProviderRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// GetByID retrieves a provider by its ID.
func (r *OIDCProviderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.OIDCProvider, error) {
	provider, err := r.queries.GetOIDCProviderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOIDCProviderNotFound
		}
		return nil, fmt.Errorf("failed to get OIDC provider by id: %w", err)
	}
	return mapDBOIDCProviderToDomain(&provider), nil
}

// GetByName retrieves a provider by its unique name.
func (r *OIDCProviderRepository) GetByName(ctx context.Context, name string) (*domain.OIDCProvider, error) {
	provider, err := r.queries.GetOIDCProviderByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOIDCProviderNotFound
		}
		return nil, fmt.Errorf("failed to get OIDC provider by name: %w", err)
	}
	return mapDBOIDCProviderToDomain(&provider), nil
}

// List retrieves all providers.
func (r *OIDCProviderRepository) List(ctx context.Context) ([]*domain.OIDCProvider, error) {
	providers, err := r.queries.ListOIDCProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list OIDC providers: %w", err)
	}

	result := make([]*domain.OIDCProvider, len(providers))
	for i, p := range providers {
		result[i] = mapDBOIDCProviderToDomain(&p)
	}
	return result, nil
}

// ListEnabled retrieves all enabled providers.
func (r *OIDCProviderRepository) ListEnabled(ctx context.Context) ([]*domain.OIDCProvider, error) {
	providers, err := r.queries.ListEnabledOIDCProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list enabled OIDC providers: %w", err)
	}

	result := make([]*domain.OIDCProvider, len(providers))
	for i, p := range providers {
		result[i] = mapDBOIDCProviderToDomain(&p)
	}
	return result, nil
}

// Create creates a new provider.
func (r *OIDCProviderRepository) Create(ctx context.Context, params domain.CreateOIDCProviderParams) (*domain.OIDCProvider, error) {
	claimMappings := params.ClaimMappings
	if claimMappings == nil {
		claimMappings = []byte("{}")
	}

	dbParams := db.CreateOIDCProviderParams{
		Name:                  params.Name,
		DisplayName:           params.DisplayName,
		IssuerUrl:             params.IssuerURL,
		ClientID:              params.ClientID,
		ClientSecretEncrypted: params.ClientSecretEncrypted,
		Scopes:                params.Scopes,
		Enabled:               params.Enabled,
		AutoCreateUsers:       params.AutoCreateUsers,
		DefaultAdmin:          params.DefaultAdmin,
		ClaimMappings:         claimMappings,
	}

	provider, err := r.queries.CreateOIDCProvider(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	return mapDBOIDCProviderToDomain(&provider), nil
}

// Update updates an existing provider.
func (r *OIDCProviderRepository) Update(ctx context.Context, params domain.UpdateOIDCProviderParams) (*domain.OIDCProvider, error) {
	// Get current provider to apply partial updates
	current, err := r.queries.GetOIDCProviderByID(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOIDCProviderNotFound
		}
		return nil, fmt.Errorf("failed to get OIDC provider for update: %w", err)
	}

	// Build update params - sqlc generates non-nullable types for COALESCE
	dbParams := db.UpdateOIDCProviderParams{
		ID: params.ID,
	}

	// Apply updates or keep current values
	if params.DisplayName != nil {
		dbParams.DisplayName = *params.DisplayName
	} else {
		dbParams.DisplayName = current.DisplayName
	}

	if params.IssuerURL != nil {
		dbParams.IssuerUrl = *params.IssuerURL
	} else {
		dbParams.IssuerUrl = current.IssuerUrl
	}

	if params.ClientID != nil {
		dbParams.ClientID = *params.ClientID
	} else {
		dbParams.ClientID = current.ClientID
	}

	if params.ClientSecretEncrypted != nil {
		dbParams.ClientSecretEncrypted = params.ClientSecretEncrypted
	} else {
		dbParams.ClientSecretEncrypted = current.ClientSecretEncrypted
	}

	if params.Scopes != nil {
		dbParams.Scopes = params.Scopes
	} else {
		dbParams.Scopes = current.Scopes
	}

	if params.Enabled != nil {
		dbParams.Enabled = *params.Enabled
	} else {
		dbParams.Enabled = current.Enabled
	}

	if params.AutoCreateUsers != nil {
		dbParams.AutoCreateUsers = *params.AutoCreateUsers
	} else {
		dbParams.AutoCreateUsers = current.AutoCreateUsers
	}

	if params.DefaultAdmin != nil {
		dbParams.DefaultAdmin = *params.DefaultAdmin
	} else {
		dbParams.DefaultAdmin = current.DefaultAdmin
	}

	if params.ClaimMappings != nil {
		dbParams.ClaimMappings = params.ClaimMappings
	} else {
		dbParams.ClaimMappings = current.ClaimMappings
	}

	provider, err := r.queries.UpdateOIDCProvider(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update OIDC provider: %w", err)
	}
	return mapDBOIDCProviderToDomain(&provider), nil
}

// SetEnabled enables or disables a provider.
func (r *OIDCProviderRepository) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	err := r.queries.UpdateOIDCProviderEnabled(ctx, db.UpdateOIDCProviderEnabledParams{
		ID:      id,
		Enabled: enabled,
	})
	if err != nil {
		return fmt.Errorf("failed to set OIDC provider enabled: %w", err)
	}
	return nil
}

// Delete removes a provider.
func (r *OIDCProviderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteOIDCProvider(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete OIDC provider: %w", err)
	}
	return nil
}

// mapDBOIDCProviderToDomain converts a database OIDC provider to a domain OIDC provider.
func mapDBOIDCProviderToDomain(p *db.OidcProvider) *domain.OIDCProvider {
	return &domain.OIDCProvider{
		ID:                    p.ID,
		Name:                  p.Name,
		DisplayName:           p.DisplayName,
		IssuerURL:             p.IssuerUrl,
		ClientID:              p.ClientID,
		ClientSecretEncrypted: p.ClientSecretEncrypted,
		Scopes:                p.Scopes,
		Enabled:               p.Enabled,
		AutoCreateUsers:       p.AutoCreateUsers,
		DefaultAdmin:          p.DefaultAdmin,
		ClaimMappings:         p.ClaimMappings,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}

// Ensure OIDCProviderRepository implements domain.OIDCProviderRepository.
var _ domain.OIDCProviderRepository = (*OIDCProviderRepository)(nil)

// OIDCUserLinkRepository implements domain.OIDCUserLinkRepository using PostgreSQL.
type OIDCUserLinkRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewOIDCUserLinkRepository creates a new PostgreSQL OIDC user link repository.
func NewOIDCUserLinkRepository(pool *pgxpool.Pool) *OIDCUserLinkRepository {
	return &OIDCUserLinkRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// Get retrieves a specific link by provider and subject.
func (r *OIDCUserLinkRepository) Get(ctx context.Context, providerID uuid.UUID, subject string) (*domain.OIDCUserLink, error) {
	link, err := r.queries.GetOIDCUserLink(ctx, db.GetOIDCUserLinkParams{
		ProviderID: providerID,
		Subject:    subject,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOIDCUserLinkNotFound
		}
		return nil, fmt.Errorf("failed to get OIDC user link: %w", err)
	}
	return mapDBOIDCUserLinkToDomain(&link), nil
}

// GetByUser retrieves all OIDC links for a user.
func (r *OIDCUserLinkRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.OIDCUserLinkWithProvider, error) {
	rows, err := r.queries.GetOIDCUserLinkByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OIDC user links by user: %w", err)
	}

	result := make([]*domain.OIDCUserLinkWithProvider, len(rows))
	for i, row := range rows {
		result[i] = &domain.OIDCUserLinkWithProvider{
			OIDCUserLink: domain.OIDCUserLink{
				ID:         row.ID,
				UserID:     row.UserID,
				ProviderID: row.ProviderID,
				Subject:    row.Subject,
				CreatedAt:  row.CreatedAt,
			},
			ProviderName:        row.ProviderName,
			ProviderDisplayName: row.ProviderDisplayName,
		}

		if row.Email.Valid {
			result[i].Email = &row.Email.String
		}
		if row.LastLoginAt.Valid {
			result[i].LastLoginAt = &row.LastLoginAt.Time
		}
	}
	return result, nil
}

// Create creates a new link.
func (r *OIDCUserLinkRepository) Create(ctx context.Context, params domain.CreateOIDCUserLinkParams) (*domain.OIDCUserLink, error) {
	dbParams := db.CreateOIDCUserLinkParams{
		UserID:     params.UserID,
		ProviderID: params.ProviderID,
		Subject:    params.Subject,
	}

	if params.Email != nil {
		dbParams.Email = pgtype.Text{String: *params.Email, Valid: true}
	}

	link, err := r.queries.CreateOIDCUserLink(ctx, dbParams)
	if err != nil {
		if isUniqueViolation(err, "oidc_user_links_provider_id_subject_key") {
			return nil, domain.ErrDuplicateOIDCLink
		}
		return nil, fmt.Errorf("failed to create OIDC user link: %w", err)
	}
	return mapDBOIDCUserLinkToDomain(&link), nil
}

// UpdateLastLogin updates the link's last login timestamp.
func (r *OIDCUserLinkRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	err := r.queries.UpdateOIDCUserLinkLastLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update OIDC user link last login: %w", err)
	}
	return nil
}

// Delete removes a link by its ID.
func (r *OIDCUserLinkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteOIDCUserLink(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete OIDC user link: %w", err)
	}
	return nil
}

// DeleteByUser removes all links for a user.
func (r *OIDCUserLinkRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.DeleteOIDCUserLinksByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete OIDC user links by user: %w", err)
	}
	return nil
}

// DeleteByProvider removes all links for a provider.
func (r *OIDCUserLinkRepository) DeleteByProvider(ctx context.Context, providerID uuid.UUID) error {
	err := r.queries.DeleteOIDCUserLinksByProvider(ctx, providerID)
	if err != nil {
		return fmt.Errorf("failed to delete OIDC user links by provider: %w", err)
	}
	return nil
}

// Exists checks if a link exists for a provider and subject.
func (r *OIDCUserLinkRepository) Exists(ctx context.Context, providerID uuid.UUID, subject string) (bool, error) {
	exists, err := r.queries.OIDCUserLinkExists(ctx, db.OIDCUserLinkExistsParams{
		ProviderID: providerID,
		Subject:    subject,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check OIDC user link exists: %w", err)
	}
	return exists, nil
}

// GetUserByOIDC retrieves a user by their OIDC identity.
func (r *OIDCUserLinkRepository) GetUserByOIDC(ctx context.Context, providerID uuid.UUID, subject string) (*domain.User, error) {
	user, err := r.queries.GetUserByOIDCLink(ctx, db.GetUserByOIDCLinkParams{
		ProviderID: providerID,
		Subject:    subject,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by OIDC link: %w", err)
	}
	return mapDBUserToDomain(&user), nil
}

// mapDBOIDCUserLinkToDomain converts a database OIDC user link to a domain OIDC user link.
func mapDBOIDCUserLinkToDomain(l *db.OidcUserLink) *domain.OIDCUserLink {
	link := &domain.OIDCUserLink{
		ID:         l.ID,
		UserID:     l.UserID,
		ProviderID: l.ProviderID,
		Subject:    l.Subject,
		CreatedAt:  l.CreatedAt,
	}

	if l.Email.Valid {
		link.Email = &l.Email.String
	}
	if l.LastLoginAt.Valid {
		link.LastLoginAt = &l.LastLoginAt.Time
	}

	return link
}

// Ensure OIDCUserLinkRepository implements domain.OIDCUserLinkRepository.
var _ domain.OIDCUserLinkRepository = (*OIDCUserLinkRepository)(nil)
