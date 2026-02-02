package apikeys

import (
	"context"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// RepositoryPg implements Repository using PostgreSQL with sqlc
type RepositoryPg struct {
	queries *db.Queries
}

// NewRepositoryPg creates a new PostgreSQL repository
func NewRepositoryPg(queries *db.Queries) Repository {
	return &RepositoryPg{queries: queries}
}

func (r *RepositoryPg) CreateAPIKey(ctx context.Context, params db.CreateAPIKeyParams) (db.SharedApiKey, error) {
	return r.queries.CreateAPIKey(ctx, params)
}

func (r *RepositoryPg) GetAPIKey(ctx context.Context, id uuid.UUID) (db.SharedApiKey, error) {
	return r.queries.GetAPIKey(ctx, id)
}

func (r *RepositoryPg) GetAPIKeyByHash(ctx context.Context, keyHash string) (db.SharedApiKey, error) {
	return r.queries.GetAPIKeyByHash(ctx, keyHash)
}

func (r *RepositoryPg) GetAPIKeyByPrefix(ctx context.Context, keyPrefix string) (db.SharedApiKey, error) {
	return r.queries.GetAPIKeyByPrefix(ctx, keyPrefix)
}

func (r *RepositoryPg) ListUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]db.SharedApiKey, error) {
	return r.queries.ListUserAPIKeys(ctx, userID)
}

func (r *RepositoryPg) ListActiveUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]db.SharedApiKey, error) {
	return r.queries.ListActiveUserAPIKeys(ctx, userID)
}

func (r *RepositoryPg) CountUserAPIKeys(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountUserAPIKeys(ctx, userID)
}

func (r *RepositoryPg) RevokeAPIKey(ctx context.Context, id uuid.UUID) error {
	return r.queries.RevokeAPIKey(ctx, id)
}

func (r *RepositoryPg) UpdateAPIKeyLastUsed(ctx context.Context, id uuid.UUID) error {
	return r.queries.UpdateAPIKeyLastUsed(ctx, id)
}

func (r *RepositoryPg) UpdateAPIKeyScopes(ctx context.Context, id uuid.UUID, scopes []string) error {
	return r.queries.UpdateAPIKeyScopes(ctx, db.UpdateAPIKeyScopesParams{
		ID:     id,
		Scopes: scopes,
	})
}

func (r *RepositoryPg) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAPIKey(ctx, id)
}

func (r *RepositoryPg) DeleteExpiredAPIKeys(ctx context.Context) error {
	return r.queries.DeleteExpiredAPIKeys(ctx)
}
