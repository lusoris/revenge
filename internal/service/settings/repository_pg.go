package settings

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// pgRepository implements Repository using PostgreSQL (via sqlc).
type pgRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPostgresRepository creates a new PostgreSQL-backed settings repository.
func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// ============================================================================
// Server Settings
// ============================================================================

func (r *pgRepository) GetServerSetting(ctx context.Context, key string) (*db.SharedServerSetting, error) {
	setting, err := r.queries.GetServerSetting(ctx, key)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *pgRepository) ListServerSettings(ctx context.Context) ([]db.SharedServerSetting, error) {
	return r.queries.ListServerSettings(ctx)
}

func (r *pgRepository) ListServerSettingsByCategory(ctx context.Context, category string) ([]db.SharedServerSetting, error) {
	return r.queries.ListServerSettingsByCategory(ctx, &category)
}

func (r *pgRepository) ListPublicServerSettings(ctx context.Context) ([]db.SharedServerSetting, error) {
	return r.queries.ListPublicServerSettings(ctx)
}

func (r *pgRepository) UpsertServerSetting(ctx context.Context, params db.UpsertServerSettingParams) (*db.SharedServerSetting, error) {
	setting, err := r.queries.UpsertServerSetting(ctx, params)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *pgRepository) UpdateServerSetting(ctx context.Context, params db.UpdateServerSettingParams) (*db.SharedServerSetting, error) {
	setting, err := r.queries.UpdateServerSetting(ctx, params)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *pgRepository) DeleteServerSetting(ctx context.Context, key string) error {
	return r.queries.DeleteServerSetting(ctx, key)
}

// ============================================================================
// User Settings
// ============================================================================

func (r *pgRepository) GetUserSetting(ctx context.Context, userID uuid.UUID, key string) (*db.SharedUserSetting, error) {
	setting, err := r.queries.GetUserSetting(ctx, db.GetUserSettingParams{
		UserID: userID,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *pgRepository) ListUserSettings(ctx context.Context, userID uuid.UUID) ([]db.SharedUserSetting, error) {
	return r.queries.ListUserSettings(ctx, userID)
}

func (r *pgRepository) ListUserSettingsByCategory(ctx context.Context, userID uuid.UUID, category string) ([]db.SharedUserSetting, error) {
	return r.queries.ListUserSettingsByCategory(ctx, db.ListUserSettingsByCategoryParams{
		UserID:   userID,
		Category: &category,
	})
}

func (r *pgRepository) UpsertUserSetting(ctx context.Context, params db.UpsertUserSettingParams) (*db.SharedUserSetting, error) {
	setting, err := r.queries.UpsertUserSetting(ctx, params)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *pgRepository) UpdateUserSetting(ctx context.Context, params db.UpdateUserSettingParams) (*db.SharedUserSetting, error) {
	setting, err := r.queries.UpdateUserSetting(ctx, params)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *pgRepository) DeleteUserSetting(ctx context.Context, userID uuid.UUID, key string) error {
	return r.queries.DeleteUserSetting(ctx, db.DeleteUserSettingParams{
		UserID: userID,
		Key:    key,
	})
}

func (r *pgRepository) DeleteAllUserSettings(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteAllUserSettings(ctx, userID)
}
