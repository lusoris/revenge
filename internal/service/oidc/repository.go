// Package oidc provides OpenID Connect authentication service
package oidc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for OIDC data access
type Repository interface {
	// Provider management
	CreateProvider(ctx context.Context, req CreateProviderRequest) (*Provider, error)
	GetProvider(ctx context.Context, id uuid.UUID) (*Provider, error)
	GetProviderByName(ctx context.Context, name string) (*Provider, error)
	GetDefaultProvider(ctx context.Context) (*Provider, error)
	ListProviders(ctx context.Context) ([]Provider, error)
	ListEnabledProviders(ctx context.Context) ([]Provider, error)
	UpdateProvider(ctx context.Context, id uuid.UUID, req UpdateProviderRequest) (*Provider, error)
	DeleteProvider(ctx context.Context, id uuid.UUID) error
	EnableProvider(ctx context.Context, id uuid.UUID) error
	DisableProvider(ctx context.Context, id uuid.UUID) error
	SetDefaultProvider(ctx context.Context, id uuid.UUID) error

	// User links
	CreateUserLink(ctx context.Context, req CreateUserLinkRequest) (*UserLink, error)
	GetUserLink(ctx context.Context, id uuid.UUID) (*UserLink, error)
	GetUserLinkBySubject(ctx context.Context, providerID uuid.UUID, subject string) (*UserLink, error)
	GetUserLinkByUserAndProvider(ctx context.Context, userID, providerID uuid.UUID) (*UserLink, error)
	ListUserLinks(ctx context.Context, userID uuid.UUID) ([]UserLinkWithProvider, error)
	UpdateUserLink(ctx context.Context, id uuid.UUID, req UpdateUserLinkRequest) (*UserLink, error)
	UpdateUserLinkLastLogin(ctx context.Context, id uuid.UUID) error
	DeleteUserLink(ctx context.Context, id uuid.UUID) error
	DeleteUserLinkByUserAndProvider(ctx context.Context, userID, providerID uuid.UUID) error
	CountUserLinks(ctx context.Context, userID uuid.UUID) (int64, error)

	// OAuth state management
	CreateState(ctx context.Context, req CreateStateRequest) (*State, error)
	GetState(ctx context.Context, state string) (*State, error)
	DeleteState(ctx context.Context, state string) error
	DeleteExpiredStates(ctx context.Context) (int64, error)
	DeleteStatesByProvider(ctx context.Context, providerID uuid.UUID) error
}

// ============================================================================
// Domain Types
// ============================================================================

// ClaimMappings defines how to map OIDC claims to user fields
type ClaimMappings struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Roles    string `json:"roles"`
}

// Provider represents an OIDC provider configuration
type Provider struct {
	ID                    uuid.UUID         `json:"id"`
	Name                  string            `json:"name"`
	DisplayName           string            `json:"displayName"`
	ProviderType          string            `json:"providerType"`
	IssuerURL             string            `json:"issuerUrl"`
	ClientID              string            `json:"clientId"`
	ClientSecretEncrypted []byte            `json:"-"` // Never expose
	AuthorizationEndpoint *string           `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         *string           `json:"tokenEndpoint,omitempty"`
	UserInfoEndpoint      *string           `json:"userInfoEndpoint,omitempty"`
	JWKSURI               *string           `json:"jwksUri,omitempty"`
	EndSessionEndpoint    *string           `json:"endSessionEndpoint,omitempty"`
	Scopes                []string          `json:"scopes"`
	ClaimMappings         ClaimMappings     `json:"claimMappings"`
	RoleMappings          map[string]string `json:"roleMappings"`
	AutoCreateUsers       bool              `json:"autoCreateUsers"`
	UpdateUserInfo        bool              `json:"updateUserInfo"`
	AllowLinking          bool              `json:"allowLinking"`
	IsEnabled             bool              `json:"isEnabled"`
	IsDefault             bool              `json:"isDefault"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
}

// UserLink represents a link between a user and an OIDC provider
type UserLink struct {
	ID                    uuid.UUID  `json:"id"`
	UserID                uuid.UUID  `json:"userId"`
	ProviderID            uuid.UUID  `json:"providerId"`
	Subject               string     `json:"subject"`
	Email                 *string    `json:"email,omitempty"`
	Name                  *string    `json:"name,omitempty"`
	PictureURL            *string    `json:"pictureUrl,omitempty"`
	AccessTokenEncrypted  []byte     `json:"-"` // Never expose
	RefreshTokenEncrypted []byte     `json:"-"` // Never expose
	TokenExpiresAt        *time.Time `json:"tokenExpiresAt,omitempty"`
	LastLoginAt           *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

// UserLinkWithProvider includes provider info
type UserLinkWithProvider struct {
	UserLink
	ProviderName        string `json:"providerName"`
	ProviderDisplayName string `json:"providerDisplayName"`
}

// State represents an OAuth2 state for the auth flow
type State struct {
	ID           uuid.UUID  `json:"id"`
	State        string     `json:"state"`
	CodeVerifier *string    `json:"-"` // PKCE verifier - never expose
	Nonce        *string    `json:"-"` // OIDC nonce - never expose
	ProviderID   uuid.UUID  `json:"providerId"`
	UserID       *uuid.UUID `json:"userId,omitempty"` // For linking flow
	RedirectURL  *string    `json:"redirectUrl,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// ============================================================================
// Request Types
// ============================================================================

// CreateProviderRequest contains data for creating a provider
type CreateProviderRequest struct {
	Name                  string
	DisplayName           string
	ProviderType          string
	IssuerURL             string
	ClientID              string
	ClientSecretEncrypted []byte
	AuthorizationEndpoint *string
	TokenEndpoint         *string
	UserInfoEndpoint      *string
	JWKSURI               *string
	EndSessionEndpoint    *string
	Scopes                []string
	ClaimMappings         ClaimMappings
	RoleMappings          map[string]string
	AutoCreateUsers       bool
	UpdateUserInfo        bool
	AllowLinking          bool
	IsEnabled             bool
	IsDefault             bool
}

// UpdateProviderRequest contains data for updating a provider
type UpdateProviderRequest struct {
	DisplayName           *string
	ProviderType          *string
	IssuerURL             *string
	ClientID              *string
	ClientSecretEncrypted []byte
	AuthorizationEndpoint *string
	TokenEndpoint         *string
	UserInfoEndpoint      *string
	JWKSURI               *string
	EndSessionEndpoint    *string
	Scopes                []string
	ClaimMappings         *ClaimMappings
	RoleMappings          map[string]string
	AutoCreateUsers       *bool
	UpdateUserInfo        *bool
	AllowLinking          *bool
	IsEnabled             *bool
	IsDefault             *bool
}

// CreateUserLinkRequest contains data for creating a user link
type CreateUserLinkRequest struct {
	UserID                uuid.UUID
	ProviderID            uuid.UUID
	Subject               string
	Email                 *string
	Name                  *string
	PictureURL            *string
	AccessTokenEncrypted  []byte
	RefreshTokenEncrypted []byte
	TokenExpiresAt        *time.Time
}

// UpdateUserLinkRequest contains data for updating a user link
type UpdateUserLinkRequest struct {
	Email                 *string
	Name                  *string
	PictureURL            *string
	AccessTokenEncrypted  []byte
	RefreshTokenEncrypted []byte
	TokenExpiresAt        *time.Time
	LastLoginAt           *time.Time
}

// CreateStateRequest contains data for creating an OAuth state
type CreateStateRequest struct {
	State        string
	CodeVerifier *string
	Nonce        *string
	ProviderID   uuid.UUID
	UserID       *uuid.UUID
	RedirectURL  *string
	ExpiresAt    time.Time
}

// ============================================================================
// JSON Helpers
// ============================================================================

// MarshalClaimMappings converts ClaimMappings to JSON
func MarshalClaimMappings(cm ClaimMappings) (json.RawMessage, error) {
	return json.Marshal(cm)
}

// UnmarshalClaimMappings converts JSON to ClaimMappings
func UnmarshalClaimMappings(data json.RawMessage) (ClaimMappings, error) {
	var cm ClaimMappings
	if err := json.Unmarshal(data, &cm); err != nil {
		return ClaimMappings{}, err
	}
	return cm, nil
}

// MarshalRoleMappings converts role mappings to JSON
func MarshalRoleMappings(rm map[string]string) (json.RawMessage, error) {
	if rm == nil {
		return json.RawMessage("{}"), nil
	}
	return json.Marshal(rm)
}

// UnmarshalRoleMappings converts JSON to role mappings
func UnmarshalRoleMappings(data json.RawMessage) (map[string]string, error) {
	var rm map[string]string
	if err := json.Unmarshal(data, &rm); err != nil {
		return nil, err
	}
	return rm, nil
}
