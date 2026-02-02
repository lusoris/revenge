package oidc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// RepositoryPg implements Repository using PostgreSQL
type RepositoryPg struct {
	q *db.Queries
}

// NewRepositoryPg creates a new PostgreSQL repository
func NewRepositoryPg(q *db.Queries) *RepositoryPg {
	return &RepositoryPg{q: q}
}

// ============================================================================
// Provider Methods
// ============================================================================

func (r *RepositoryPg) CreateProvider(ctx context.Context, req CreateProviderRequest) (*Provider, error) {
	claimMappingsJSON, err := MarshalClaimMappings(req.ClaimMappings)
	if err != nil {
		return nil, err
	}
	roleMappingsJSON, err := MarshalRoleMappings(req.RoleMappings)
	if err != nil {
		return nil, err
	}

	dbProvider, err := r.q.CreateOIDCProvider(ctx, db.CreateOIDCProviderParams{
		Name:                  req.Name,
		DisplayName:           req.DisplayName,
		ProviderType:          req.ProviderType,
		IssuerUrl:             req.IssuerURL,
		ClientID:              req.ClientID,
		ClientSecretEncrypted: req.ClientSecretEncrypted,
		AuthorizationEndpoint: req.AuthorizationEndpoint,
		TokenEndpoint:         req.TokenEndpoint,
		UserinfoEndpoint:      req.UserInfoEndpoint,
		JwksUri:               req.JWKSURI,
		EndSessionEndpoint:    req.EndSessionEndpoint,
		Scopes:                req.Scopes,
		ClaimMappings:         claimMappingsJSON,
		RoleMappings:          roleMappingsJSON,
		AutoCreateUsers:       req.AutoCreateUsers,
		UpdateUserInfo:        req.UpdateUserInfo,
		AllowLinking:          req.AllowLinking,
		IsEnabled:             req.IsEnabled,
		IsDefault:             req.IsDefault,
	})
	if err != nil {
		return nil, err
	}

	return dbProviderToProvider(dbProvider)
}

func (r *RepositoryPg) GetProvider(ctx context.Context, id uuid.UUID) (*Provider, error) {
	dbProvider, err := r.q.GetOIDCProvider(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}
	return dbProviderToProvider(dbProvider)
}

func (r *RepositoryPg) GetProviderByName(ctx context.Context, name string) (*Provider, error) {
	dbProvider, err := r.q.GetOIDCProviderByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}
	return dbProviderToProvider(dbProvider)
}

func (r *RepositoryPg) GetDefaultProvider(ctx context.Context) (*Provider, error) {
	dbProvider, err := r.q.GetDefaultOIDCProvider(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}
	return dbProviderToProvider(dbProvider)
}

func (r *RepositoryPg) ListProviders(ctx context.Context) ([]Provider, error) {
	dbProviders, err := r.q.ListOIDCProviders(ctx)
	if err != nil {
		return nil, err
	}

	providers := make([]Provider, len(dbProviders))
	for i, dbp := range dbProviders {
		p, err := dbProviderToProvider(dbp)
		if err != nil {
			return nil, err
		}
		providers[i] = *p
	}
	return providers, nil
}

func (r *RepositoryPg) ListEnabledProviders(ctx context.Context) ([]Provider, error) {
	dbProviders, err := r.q.ListEnabledOIDCProviders(ctx)
	if err != nil {
		return nil, err
	}

	providers := make([]Provider, len(dbProviders))
	for i, dbp := range dbProviders {
		p, err := dbProviderToProvider(dbp)
		if err != nil {
			return nil, err
		}
		providers[i] = *p
	}
	return providers, nil
}

func (r *RepositoryPg) UpdateProvider(ctx context.Context, id uuid.UUID, req UpdateProviderRequest) (*Provider, error) {
	params := db.UpdateOIDCProviderParams{ID: id}

	if req.DisplayName != nil {
		params.DisplayName = req.DisplayName
	}
	if req.ProviderType != nil {
		params.ProviderType = req.ProviderType
	}
	if req.IssuerURL != nil {
		params.IssuerUrl = req.IssuerURL
	}
	if req.ClientID != nil {
		params.ClientID = req.ClientID
	}
	if req.ClientSecretEncrypted != nil {
		params.ClientSecretEncrypted = req.ClientSecretEncrypted
	}
	if req.AuthorizationEndpoint != nil {
		params.AuthorizationEndpoint = req.AuthorizationEndpoint
	}
	if req.TokenEndpoint != nil {
		params.TokenEndpoint = req.TokenEndpoint
	}
	if req.UserInfoEndpoint != nil {
		params.UserinfoEndpoint = req.UserInfoEndpoint
	}
	if req.JWKSURI != nil {
		params.JwksUri = req.JWKSURI
	}
	if req.EndSessionEndpoint != nil {
		params.EndSessionEndpoint = req.EndSessionEndpoint
	}
	if req.Scopes != nil {
		params.Scopes = req.Scopes
	}
	if req.ClaimMappings != nil {
		cm, err := MarshalClaimMappings(*req.ClaimMappings)
		if err != nil {
			return nil, err
		}
		params.ClaimMappings = cm
	}
	if req.RoleMappings != nil {
		rm, err := MarshalRoleMappings(req.RoleMappings)
		if err != nil {
			return nil, err
		}
		params.RoleMappings = rm
	}
	if req.AutoCreateUsers != nil {
		params.AutoCreateUsers = req.AutoCreateUsers
	}
	if req.UpdateUserInfo != nil {
		params.UpdateUserInfo = req.UpdateUserInfo
	}
	if req.AllowLinking != nil {
		params.AllowLinking = req.AllowLinking
	}
	if req.IsEnabled != nil {
		params.IsEnabled = req.IsEnabled
	}
	if req.IsDefault != nil {
		params.IsDefault = req.IsDefault
	}

	dbProvider, err := r.q.UpdateOIDCProvider(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}
	return dbProviderToProvider(dbProvider)
}

func (r *RepositoryPg) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteOIDCProvider(ctx, id)
}

func (r *RepositoryPg) EnableProvider(ctx context.Context, id uuid.UUID) error {
	return r.q.EnableOIDCProvider(ctx, id)
}

func (r *RepositoryPg) DisableProvider(ctx context.Context, id uuid.UUID) error {
	return r.q.DisableOIDCProvider(ctx, id)
}

func (r *RepositoryPg) SetDefaultProvider(ctx context.Context, id uuid.UUID) error {
	// Clear existing default
	if err := r.q.SetDefaultOIDCProvider(ctx); err != nil {
		return err
	}
	// Set new default
	isDefault := true
	_, err := r.q.UpdateOIDCProvider(ctx, db.UpdateOIDCProviderParams{
		ID:        id,
		IsDefault: &isDefault,
	})
	return err
}

// ============================================================================
// User Link Methods
// ============================================================================

func (r *RepositoryPg) CreateUserLink(ctx context.Context, req CreateUserLinkRequest) (*UserLink, error) {
	var tokenExpiresAt pgtype.Timestamptz
	if req.TokenExpiresAt != nil {
		tokenExpiresAt = pgtype.Timestamptz{Time: *req.TokenExpiresAt, Valid: true}
	}

	dbLink, err := r.q.CreateOIDCUserLink(ctx, db.CreateOIDCUserLinkParams{
		UserID:                req.UserID,
		ProviderID:            req.ProviderID,
		Subject:               req.Subject,
		Email:                 req.Email,
		Name:                  req.Name,
		PictureUrl:            req.PictureURL,
		AccessTokenEncrypted:  req.AccessTokenEncrypted,
		RefreshTokenEncrypted: req.RefreshTokenEncrypted,
		TokenExpiresAt:        tokenExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return dbUserLinkToUserLink(dbLink), nil
}

func (r *RepositoryPg) GetUserLink(ctx context.Context, id uuid.UUID) (*UserLink, error) {
	dbLink, err := r.q.GetOIDCUserLink(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserLinkNotFound
		}
		return nil, err
	}
	return dbUserLinkToUserLink(dbLink), nil
}

func (r *RepositoryPg) GetUserLinkBySubject(ctx context.Context, providerID uuid.UUID, subject string) (*UserLink, error) {
	dbLink, err := r.q.GetOIDCUserLinkBySubject(ctx, db.GetOIDCUserLinkBySubjectParams{
		ProviderID: providerID,
		Subject:    subject,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserLinkNotFound
		}
		return nil, err
	}
	return dbUserLinkToUserLink(dbLink), nil
}

func (r *RepositoryPg) GetUserLinkByUserAndProvider(ctx context.Context, userID, providerID uuid.UUID) (*UserLink, error) {
	dbLink, err := r.q.GetOIDCUserLinkByUserAndProvider(ctx, db.GetOIDCUserLinkByUserAndProviderParams{
		UserID:     userID,
		ProviderID: providerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserLinkNotFound
		}
		return nil, err
	}
	return dbUserLinkToUserLink(dbLink), nil
}

func (r *RepositoryPg) ListUserLinks(ctx context.Context, userID uuid.UUID) ([]UserLinkWithProvider, error) {
	dbLinks, err := r.q.ListUserOIDCLinks(ctx, userID)
	if err != nil {
		return nil, err
	}

	links := make([]UserLinkWithProvider, len(dbLinks))
	for i, dbl := range dbLinks {
		link := dbUserLinkRowToUserLinkWithProvider(dbl)
		links[i] = link
	}
	return links, nil
}

func (r *RepositoryPg) UpdateUserLink(ctx context.Context, id uuid.UUID, req UpdateUserLinkRequest) (*UserLink, error) {
	params := db.UpdateOIDCUserLinkParams{ID: id}

	if req.Email != nil {
		params.Email = req.Email
	}
	if req.Name != nil {
		params.Name = req.Name
	}
	if req.PictureURL != nil {
		params.PictureUrl = req.PictureURL
	}
	if req.AccessTokenEncrypted != nil {
		params.AccessTokenEncrypted = req.AccessTokenEncrypted
	}
	if req.RefreshTokenEncrypted != nil {
		params.RefreshTokenEncrypted = req.RefreshTokenEncrypted
	}
	if req.TokenExpiresAt != nil {
		params.TokenExpiresAt = pgtype.Timestamptz{Time: *req.TokenExpiresAt, Valid: true}
	}
	if req.LastLoginAt != nil {
		params.LastLoginAt = pgtype.Timestamptz{Time: *req.LastLoginAt, Valid: true}
	}

	dbLink, err := r.q.UpdateOIDCUserLink(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserLinkNotFound
		}
		return nil, err
	}
	return dbUserLinkToUserLink(dbLink), nil
}

func (r *RepositoryPg) UpdateUserLinkLastLogin(ctx context.Context, id uuid.UUID) error {
	return r.q.UpdateOIDCUserLinkLastLogin(ctx, id)
}

func (r *RepositoryPg) DeleteUserLink(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteOIDCUserLink(ctx, id)
}

func (r *RepositoryPg) DeleteUserLinkByUserAndProvider(ctx context.Context, userID, providerID uuid.UUID) error {
	return r.q.DeleteOIDCUserLinkByUserAndProvider(ctx, db.DeleteOIDCUserLinkByUserAndProviderParams{
		UserID:     userID,
		ProviderID: providerID,
	})
}

func (r *RepositoryPg) CountUserLinks(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.q.CountUserOIDCLinks(ctx, userID)
}

// ============================================================================
// State Methods
// ============================================================================

func (r *RepositoryPg) CreateState(ctx context.Context, req CreateStateRequest) (*State, error) {
	var userID pgtype.UUID
	if req.UserID != nil {
		userID = pgtype.UUID{Bytes: *req.UserID, Valid: true}
	}

	dbState, err := r.q.CreateOIDCState(ctx, db.CreateOIDCStateParams{
		State:        req.State,
		CodeVerifier: req.CodeVerifier,
		ProviderID:   req.ProviderID,
		UserID:       userID,
		RedirectUrl:  req.RedirectURL,
		ExpiresAt:    req.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return dbStateToState(dbState), nil
}

func (r *RepositoryPg) GetState(ctx context.Context, state string) (*State, error) {
	dbState, err := r.q.GetOIDCState(ctx, state)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStateNotFound
		}
		return nil, err
	}
	return dbStateToState(dbState), nil
}

func (r *RepositoryPg) DeleteState(ctx context.Context, state string) error {
	return r.q.DeleteOIDCState(ctx, state)
}

func (r *RepositoryPg) DeleteExpiredStates(ctx context.Context) (int64, error) {
	return r.q.DeleteExpiredOIDCStates(ctx)
}

func (r *RepositoryPg) DeleteStatesByProvider(ctx context.Context, providerID uuid.UUID) error {
	return r.q.DeleteOIDCStatesByProvider(ctx, providerID)
}

// ============================================================================
// Conversion Helpers
// ============================================================================

func dbProviderToProvider(dbp db.SharedOidcProvider) (*Provider, error) {
	claimMappings, err := UnmarshalClaimMappings(dbp.ClaimMappings)
	if err != nil {
		return nil, err
	}
	roleMappings, err := UnmarshalRoleMappings(dbp.RoleMappings)
	if err != nil {
		return nil, err
	}

	return &Provider{
		ID:                    dbp.ID,
		Name:                  dbp.Name,
		DisplayName:           dbp.DisplayName,
		ProviderType:          dbp.ProviderType,
		IssuerURL:             dbp.IssuerUrl,
		ClientID:              dbp.ClientID,
		ClientSecretEncrypted: dbp.ClientSecretEncrypted,
		AuthorizationEndpoint: dbp.AuthorizationEndpoint,
		TokenEndpoint:         dbp.TokenEndpoint,
		UserInfoEndpoint:      dbp.UserinfoEndpoint,
		JWKSURI:               dbp.JwksUri,
		EndSessionEndpoint:    dbp.EndSessionEndpoint,
		Scopes:                dbp.Scopes,
		ClaimMappings:         claimMappings,
		RoleMappings:          roleMappings,
		AutoCreateUsers:       dbp.AutoCreateUsers,
		UpdateUserInfo:        dbp.UpdateUserInfo,
		AllowLinking:          dbp.AllowLinking,
		IsEnabled:             dbp.IsEnabled,
		IsDefault:             dbp.IsDefault,
		CreatedAt:             dbp.CreatedAt,
		UpdatedAt:             dbp.UpdatedAt,
	}, nil
}

func dbUserLinkToUserLink(dbl db.SharedOidcUserLink) *UserLink {
	link := &UserLink{
		ID:                    dbl.ID,
		UserID:                dbl.UserID,
		ProviderID:            dbl.ProviderID,
		Subject:               dbl.Subject,
		Email:                 dbl.Email,
		Name:                  dbl.Name,
		PictureURL:            dbl.PictureUrl,
		AccessTokenEncrypted:  dbl.AccessTokenEncrypted,
		RefreshTokenEncrypted: dbl.RefreshTokenEncrypted,
		CreatedAt:             dbl.CreatedAt,
		UpdatedAt:             dbl.UpdatedAt,
	}
	if dbl.TokenExpiresAt.Valid {
		link.TokenExpiresAt = &dbl.TokenExpiresAt.Time
	}
	if dbl.LastLoginAt.Valid {
		link.LastLoginAt = &dbl.LastLoginAt.Time
	}
	return link
}

func dbUserLinkRowToUserLinkWithProvider(dbl db.ListUserOIDCLinksRow) UserLinkWithProvider {
	link := UserLinkWithProvider{
		UserLink: UserLink{
			ID:                    dbl.ID,
			UserID:                dbl.UserID,
			ProviderID:            dbl.ProviderID,
			Subject:               dbl.Subject,
			Email:                 dbl.Email,
			Name:                  dbl.Name,
			PictureURL:            dbl.PictureUrl,
			AccessTokenEncrypted:  dbl.AccessTokenEncrypted,
			RefreshTokenEncrypted: dbl.RefreshTokenEncrypted,
			CreatedAt:             dbl.CreatedAt,
			UpdatedAt:             dbl.UpdatedAt,
		},
		ProviderName:        dbl.ProviderName,
		ProviderDisplayName: dbl.ProviderDisplayName,
	}
	if dbl.TokenExpiresAt.Valid {
		link.TokenExpiresAt = &dbl.TokenExpiresAt.Time
	}
	if dbl.LastLoginAt.Valid {
		link.LastLoginAt = &dbl.LastLoginAt.Time
	}
	return link
}

func dbStateToState(dbs db.SharedOidcState) *State {
	state := &State{
		ID:           dbs.ID,
		State:        dbs.State,
		CodeVerifier: dbs.CodeVerifier,
		ProviderID:   dbs.ProviderID,
		RedirectURL:  dbs.RedirectUrl,
		ExpiresAt:    dbs.ExpiresAt,
		CreatedAt:    dbs.CreatedAt,
	}
	if dbs.UserID.Valid {
		userID := uuid.UUID(dbs.UserID.Bytes)
		state.UserID = &userID
	}
	return state
}

// Ensure interface compliance
var _ Repository = (*RepositoryPg)(nil)

// Errors for ErrUserLinkNotFound, ErrProviderNotFound, and ErrStateNotFound
// are defined in service.go
// We need json import for MarshalClaimMappings and MarshalRoleMappings
var _ = json.Marshal
