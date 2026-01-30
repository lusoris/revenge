// Package flag provides adult tag domain models (QAR obfuscation: tags → flags).
package flag

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	adultdb "github.com/lusoris/revenge/internal/content/qar/db"
)

// ErrFlagNotFound is returned when a flag cannot be found.
var ErrFlagNotFound = errors.New("flag not found")

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed flag repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.flag")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Flag, error) {
	row, err := r.queries.GetFlagByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFlagNotFound
		}
		return nil, err
	}
	return r.rowToFlag(&row), nil
}

func (r *SQLCRepository) GetByName(ctx context.Context, name string) (*Flag, error) {
	row, err := r.queries.GetFlagByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFlagNotFound
		}
		return nil, err
	}
	return r.rowToFlag(&row), nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Flag, error) {
	rows, err := r.queries.ListFlags(ctx, adultdb.ListFlagsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) ListRoot(ctx context.Context, limit, offset int) ([]Flag, error) {
	rows, err := r.queries.ListRootFlags(ctx, adultdb.ListRootFlagsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) ListChildren(ctx context.Context, parentID uuid.UUID) ([]Flag, error) {
	rows, err := r.queries.ListFlagChildren(ctx, pgtype.UUID{Bytes: parentID, Valid: true})
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) ListByWaters(ctx context.Context, waters string) ([]Flag, error) {
	rows, err := r.queries.ListFlagsByWaters(ctx, &waters)
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, flag *Flag) error {
	params := r.flagToCreateParams(flag)
	row, err := r.queries.CreateFlag(ctx, params)
	if err != nil {
		return err
	}
	flag.ID = row.ID
	flag.CreatedAt = row.CreatedAt
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, flag *Flag) error {
	params := r.flagToUpdateParams(flag)
	_, err := r.queries.UpdateFlag(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrFlagNotFound
		}
		return err
	}
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteFlag(ctx, id)
}

func (r *SQLCRepository) GetByStashDBID(ctx context.Context, stashdbID string) (*Flag, error) {
	row, err := r.queries.GetFlagByStashDBID(ctx, &stashdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFlagNotFound
		}
		return nil, err
	}
	return r.rowToFlag(&row), nil
}

func (r *SQLCRepository) Search(ctx context.Context, query string, limit, offset int) ([]Flag, error) {
	rows, err := r.queries.SearchFlags(ctx, adultdb.SearchFlagsParams{
		Column1: &query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) ListExpeditionFlags(ctx context.Context, expeditionID uuid.UUID) ([]Flag, error) {
	rows, err := r.queries.ListExpeditionFlags(ctx, expeditionID)
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) ListVoyageFlags(ctx context.Context, voyageID uuid.UUID) ([]Flag, error) {
	rows, err := r.queries.ListVoyageFlags(ctx, voyageID)
	if err != nil {
		return nil, err
	}
	return r.rowsToFlags(rows), nil
}

func (r *SQLCRepository) AddExpeditionFlag(ctx context.Context, expeditionID, flagID uuid.UUID) error {
	return r.queries.AddExpeditionFlag(ctx, adultdb.AddExpeditionFlagParams{
		ExpeditionID: expeditionID,
		FlagID:       flagID,
	})
}

func (r *SQLCRepository) AddVoyageFlag(ctx context.Context, voyageID, flagID uuid.UUID) error {
	return r.queries.AddVoyageFlag(ctx, adultdb.AddVoyageFlagParams{
		VoyageID: voyageID,
		FlagID:   flagID,
	})
}

func (r *SQLCRepository) RemoveExpeditionFlag(ctx context.Context, expeditionID, flagID uuid.UUID) error {
	return r.queries.RemoveExpeditionFlag(ctx, adultdb.RemoveExpeditionFlagParams{
		ExpeditionID: expeditionID,
		FlagID:       flagID,
	})
}

func (r *SQLCRepository) RemoveVoyageFlag(ctx context.Context, voyageID, flagID uuid.UUID) error {
	return r.queries.RemoveVoyageFlag(ctx, adultdb.RemoveVoyageFlagParams{
		VoyageID: voyageID,
		FlagID:   flagID,
	})
}

func (r *SQLCRepository) ClearExpeditionFlags(ctx context.Context, expeditionID uuid.UUID) error {
	return r.queries.ClearExpeditionFlags(ctx, expeditionID)
}

func (r *SQLCRepository) ClearVoyageFlags(ctx context.Context, voyageID uuid.UUID) error {
	return r.queries.ClearVoyageFlags(ctx, voyageID)
}

// rowToFlag converts a database row to a domain entity.
func (r *SQLCRepository) rowToFlag(row *adultdb.QarFlag) *Flag {
	if row == nil {
		return nil
	}

	f := &Flag{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
	}

	if row.Description != nil {
		f.Description = *row.Description
	}

	if row.ParentID.Valid {
		parentID := uuid.UUID(row.ParentID.Bytes)
		f.ParentID = &parentID
	}

	if row.StashdbID != nil {
		f.StashDBID = *row.StashdbID
	}

	if row.Waters != nil {
		f.Waters = *row.Waters
	}

	return f
}

// rowsToFlags converts multiple database rows to domain entities.
func (r *SQLCRepository) rowsToFlags(rows []adultdb.QarFlag) []Flag {
	result := make([]Flag, 0, len(rows))
	for i := range rows {
		if f := r.rowToFlag(&rows[i]); f != nil {
			result = append(result, *f)
		}
	}
	return result
}

// flagToCreateParams converts a domain entity to create parameters.
func (r *SQLCRepository) flagToCreateParams(f *Flag) adultdb.CreateFlagParams {
	params := adultdb.CreateFlagParams{
		Name: f.Name,
	}

	if f.Description != "" {
		params.Description = &f.Description
	}

	if f.ParentID != nil {
		params.ParentID = pgtype.UUID{Bytes: *f.ParentID, Valid: true}
	}

	if f.StashDBID != "" {
		params.StashdbID = &f.StashDBID
	}

	if f.Waters != "" {
		params.Waters = &f.Waters
	}

	return params
}

// flagToUpdateParams converts a domain entity to update parameters.
func (r *SQLCRepository) flagToUpdateParams(f *Flag) adultdb.UpdateFlagParams {
	params := adultdb.UpdateFlagParams{
		ID:   f.ID,
		Name: f.Name,
	}

	if f.Description != "" {
		params.Description = &f.Description
	}

	if f.ParentID != nil {
		params.ParentID = pgtype.UUID{Bytes: *f.ParentID, Valid: true}
	}

	if f.StashDBID != "" {
		params.StashdbID = &f.StashDBID
	}

	if f.Waters != "" {
		params.Waters = &f.Waters
	}

	return params
}
