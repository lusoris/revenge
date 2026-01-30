// Package fleet provides adult library domain models (QAR obfuscation: libraries → fleets).
package fleet

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

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed fleet repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.fleet")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Fleet, error) {
	row, err := r.queries.GetFleetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFleetNotFound
		}
		return nil, err
	}
	return r.rowToFleet(&row), nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Fleet, error) {
	rows, err := r.queries.ListFleets(ctx, adultdb.ListFleetsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToFleets(rows), nil
}

func (r *SQLCRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]Fleet, error) {
	rows, err := r.queries.ListFleetsByOwner(ctx, pgtype.UUID{Bytes: ownerID, Valid: true})
	if err != nil {
		return nil, err
	}
	return r.rowsToFleets(rows), nil
}

func (r *SQLCRepository) ListByType(ctx context.Context, fleetType FleetType, limit, offset int) ([]Fleet, error) {
	rows, err := r.queries.ListFleetsByType(ctx, adultdb.ListFleetsByTypeParams{
		FleetType: string(fleetType),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToFleets(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, fleet *Fleet) error {
	params := r.fleetToCreateParams(fleet)
	row, err := r.queries.CreateFleet(ctx, params)
	if err != nil {
		return err
	}
	fleet.ID = row.ID
	fleet.CreatedAt = row.CreatedAt
	fleet.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, fleet *Fleet) error {
	params := r.fleetToUpdateParams(fleet)
	_, err := r.queries.UpdateFleet(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrFleetNotFound
		}
		return err
	}
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteFleet(ctx, id)
}

func (r *SQLCRepository) GetStats(ctx context.Context, id uuid.UUID) (*FleetStats, error) {
	row, err := r.queries.GetFleetStats(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFleetNotFound
		}
		return nil, err
	}
	return &FleetStats{
		ExpeditionCount: row.ExpeditionCount,
		VoyageCount:     row.VoyageCount,
		TotalSizeBytes:  int64(row.TotalSizeBytes),
	}, nil
}

func (r *SQLCRepository) CountExpeditions(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.CountFleetExpeditions(ctx, id)
}

func (r *SQLCRepository) CountVoyages(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.CountFleetVoyages(ctx, id)
}

// rowToFleet converts a database row to a domain entity.
func (r *SQLCRepository) rowToFleet(row *adultdb.QarFleet) *Fleet {
	if row == nil {
		return nil
	}

	f := &Fleet{
		ID:                row.ID,
		Name:              row.Name,
		FleetType:         FleetType(row.FleetType),
		Paths:             row.Paths,
		TPDBEnabled:       row.TpdbEnabled,
		WhisparrSync:      row.WhisparrSync,
		AutoTagCrew:       row.AutoTagCrew,
		FingerprintOnScan: row.FingerprintOnScan,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}

	if row.StashdbEndpoint != nil {
		f.StashDBEndpoint = *row.StashdbEndpoint
	}

	if row.OwnerUserID.Valid {
		ownerID := uuid.UUID(row.OwnerUserID.Bytes)
		f.OwnerUserID = &ownerID
	}

	return f
}

// rowsToFleets converts multiple database rows to domain entities.
func (r *SQLCRepository) rowsToFleets(rows []adultdb.QarFleet) []Fleet {
	result := make([]Fleet, 0, len(rows))
	for i := range rows {
		if f := r.rowToFleet(&rows[i]); f != nil {
			result = append(result, *f)
		}
	}
	return result
}

// fleetToCreateParams converts a domain entity to create parameters.
func (r *SQLCRepository) fleetToCreateParams(f *Fleet) adultdb.CreateFleetParams {
	params := adultdb.CreateFleetParams{
		Name:              f.Name,
		FleetType:         string(f.FleetType),
		Paths:             f.Paths,
		TpdbEnabled:       f.TPDBEnabled,
		WhisparrSync:      f.WhisparrSync,
		AutoTagCrew:       f.AutoTagCrew,
		FingerprintOnScan: f.FingerprintOnScan,
	}

	if f.StashDBEndpoint != "" {
		params.StashdbEndpoint = &f.StashDBEndpoint
	}

	if f.OwnerUserID != nil {
		params.OwnerUserID = pgtype.UUID{Bytes: *f.OwnerUserID, Valid: true}
	}

	return params
}

// fleetToUpdateParams converts a domain entity to update parameters.
func (r *SQLCRepository) fleetToUpdateParams(f *Fleet) adultdb.UpdateFleetParams {
	params := adultdb.UpdateFleetParams{
		ID:                f.ID,
		Name:              f.Name,
		FleetType:         string(f.FleetType),
		Paths:             f.Paths,
		TpdbEnabled:       f.TPDBEnabled,
		WhisparrSync:      f.WhisparrSync,
		AutoTagCrew:       f.AutoTagCrew,
		FingerprintOnScan: f.FingerprintOnScan,
	}

	if f.StashDBEndpoint != "" {
		params.StashdbEndpoint = &f.StashDBEndpoint
	}

	if f.OwnerUserID != nil {
		params.OwnerUserID = pgtype.UUID{Bytes: *f.OwnerUserID, Valid: true}
	}

	return params
}
