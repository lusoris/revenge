// Package domain contains core business entities and repository interfaces.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OIDCProvider represents an OpenID Connect provider configuration.
type OIDCProvider struct {
	ID                    uuid.UUID
	Name                  string // Internal name (keycloak, authentik, etc.)
	DisplayName           string // UI display name
	IssuerURL             string // OIDC issuer URL
	ClientID              string
	ClientSecretEncrypted []byte   // Encrypted with server key
	Scopes                []string // Default: openid, profile, email
	Enabled               bool
	AutoCreateUsers       bool   // Create users on first login
	DefaultAdmin          bool   // New users are admins
	ClaimMappings         []byte // JSON object for custom claim mappings
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// OIDCUserLink represents a link between a local user and an OIDC identity.
type OIDCUserLink struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ProviderID  uuid.UUID
	Subject     string  // OIDC subject claim (unique identifier from provider)
	Email       *string // Email from OIDC (may differ from user email)
	CreatedAt   time.Time
	LastLoginAt *time.Time
}

// OIDCUserLinkWithProvider includes link data with provider information.
type OIDCUserLinkWithProvider struct {
	OIDCUserLink
	ProviderName        string
	ProviderDisplayName string
}

// CreateOIDCProviderParams contains parameters for creating an OIDC provider.
type CreateOIDCProviderParams struct {
	Name                  string
	DisplayName           string
	IssuerURL             string
	ClientID              string
	ClientSecretEncrypted []byte
	Scopes                []string
	Enabled               bool
	AutoCreateUsers       bool
	DefaultAdmin          bool
	ClaimMappings         []byte
}

// UpdateOIDCProviderParams contains parameters for updating an OIDC provider.
type UpdateOIDCProviderParams struct {
	ID                    uuid.UUID
	DisplayName           *string
	IssuerURL             *string
	ClientID              *string
	ClientSecretEncrypted []byte
	Scopes                []string
	Enabled               *bool
	AutoCreateUsers       *bool
	DefaultAdmin          *bool
	ClaimMappings         []byte
}

// CreateOIDCUserLinkParams contains parameters for linking a user to an OIDC identity.
type CreateOIDCUserLinkParams struct {
	UserID     uuid.UUID
	ProviderID uuid.UUID
	Subject    string
	Email      *string
}

// OIDCProviderRepository defines the interface for OIDC provider data access.
type OIDCProviderRepository interface {
	// GetByID retrieves a provider by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*OIDCProvider, error)

	// GetByName retrieves a provider by its unique name.
	GetByName(ctx context.Context, name string) (*OIDCProvider, error)

	// List retrieves all providers.
	List(ctx context.Context) ([]*OIDCProvider, error)

	// ListEnabled retrieves all enabled providers.
	ListEnabled(ctx context.Context) ([]*OIDCProvider, error)

	// Create creates a new provider.
	Create(ctx context.Context, params CreateOIDCProviderParams) (*OIDCProvider, error)

	// Update updates an existing provider.
	Update(ctx context.Context, params UpdateOIDCProviderParams) (*OIDCProvider, error)

	// SetEnabled enables or disables a provider.
	SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) error

	// Delete removes a provider.
	Delete(ctx context.Context, id uuid.UUID) error
}

// OIDCUserLinkRepository defines the interface for OIDC user link data access.
type OIDCUserLinkRepository interface {
	// Get retrieves a specific link by provider and subject.
	Get(ctx context.Context, providerID uuid.UUID, subject string) (*OIDCUserLink, error)

	// GetByUser retrieves all OIDC links for a user.
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*OIDCUserLinkWithProvider, error)

	// Create creates a new link.
	Create(ctx context.Context, params CreateOIDCUserLinkParams) (*OIDCUserLink, error)

	// UpdateLastLogin updates the link's last login timestamp.
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error

	// Delete removes a link by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByUser removes all links for a user.
	DeleteByUser(ctx context.Context, userID uuid.UUID) error

	// DeleteByProvider removes all links for a provider.
	DeleteByProvider(ctx context.Context, providerID uuid.UUID) error

	// Exists checks if a link exists for a provider and subject.
	Exists(ctx context.Context, providerID uuid.UUID, subject string) (bool, error)

	// GetUserByOIDC retrieves a user by their OIDC identity.
	GetUserByOIDC(ctx context.Context, providerID uuid.UUID, subject string) (*User, error)
}
