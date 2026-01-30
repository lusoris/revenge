// Package expedition provides adult movie domain models (QAR obfuscation: movies → expeditions).
package expedition

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/revenge/internal/content/qar/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// ErrExpeditionNotFound is returned when an expedition cannot be found.
var ErrExpeditionNotFound = errors.New("expedition not found")

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed expedition repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.expedition")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Expedition, error) {
	row, err := r.queries.GetExpeditionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExpeditionNotFound
		}
		return nil, err
	}
	return r.rowToExpedition(&row), nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Expedition, error) {
	rows, err := r.queries.ListExpeditions(ctx, adultdb.ListExpeditionsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToExpeditions(rows), nil
}

func (r *SQLCRepository) ListByFleet(ctx context.Context, fleetID uuid.UUID, limit, offset int) ([]Expedition, error) {
	rows, err := r.queries.ListExpeditionsByFleet(ctx, adultdb.ListExpeditionsByFleetParams{
		FleetID: fleetID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToExpeditions(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, expedition *Expedition) error {
	params := r.expeditionToCreateParams(expedition)
	row, err := r.queries.CreateExpedition(ctx, params)
	if err != nil {
		return err
	}
	expedition.ID = row.ID
	expedition.CreatedAt = row.CreatedAt
	expedition.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, expedition *Expedition) error {
	params := r.expeditionToUpdateParams(expedition)
	_, err := r.queries.UpdateExpedition(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrExpeditionNotFound
		}
		return err
	}
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteExpedition(ctx, id)
}

func (r *SQLCRepository) GetByPath(ctx context.Context, path string) (*Expedition, error) {
	row, err := r.queries.GetExpeditionByPath(ctx, path)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExpeditionNotFound
		}
		return nil, err
	}
	return r.rowToExpedition(&row), nil
}

func (r *SQLCRepository) GetByCoordinates(ctx context.Context, coordinates string) (*Expedition, error) {
	row, err := r.queries.GetExpeditionByCoordinates(ctx, &coordinates)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExpeditionNotFound
		}
		return nil, err
	}
	return r.rowToExpedition(&row), nil
}

func (r *SQLCRepository) GetByCharter(ctx context.Context, charter string) (*Expedition, error) {
	row, err := r.queries.GetExpeditionByCharter(ctx, &charter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExpeditionNotFound
		}
		return nil, err
	}
	return r.rowToExpedition(&row), nil
}

func (r *SQLCRepository) CountByFleet(ctx context.Context, fleetID uuid.UUID) (int64, error) {
	return r.queries.CountExpeditionsByFleet(ctx, fleetID)
}

func (r *SQLCRepository) Search(ctx context.Context, query string, limit, offset int) ([]Expedition, error) {
	rows, err := r.queries.SearchExpeditions(ctx, adultdb.SearchExpeditionsParams{
		Column1: &query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToExpeditions(rows), nil
}

// rowToExpedition converts a database row to a domain entity.
func (r *SQLCRepository) rowToExpedition(row *adultdb.QarExpedition) *Expedition {
	if row == nil {
		return nil
	}

	exp := &Expedition{
		ContentEntity: shared.ContentEntity{
			BaseEntity: shared.BaseEntity{
				ID:        row.ID,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			Path:  row.Path,
			Title: row.Title,
		},
		FleetID: row.FleetID,
	}

	if row.SortTitle != nil {
		exp.SortTitle = *row.SortTitle
	}

	if row.Overview != nil {
		exp.Overview = *row.Overview
	}

	if row.LaunchDate.Valid {
		t := row.LaunchDate.Time
		exp.LaunchDate = &t
	}

	if row.RuntimeTicks != nil {
		exp.RuntimeTicks = *row.RuntimeTicks
	}

	if row.PortID.Valid {
		portID := uuid.UUID(row.PortID.Bytes)
		exp.PortID = &portID
	}

	if row.Director != nil {
		exp.Director = *row.Director
	}

	if row.Series != nil {
		exp.Series = *row.Series
	}

	if row.Coordinates != nil {
		exp.Coordinates = *row.Coordinates
	}

	if row.Charter != nil {
		exp.Charter = *row.Charter
	}

	if row.Registry != nil {
		exp.Registry = *row.Registry
	}

	if row.WhisparrID != nil {
		id := int(*row.WhisparrID)
		exp.WhisparrID = &id
	}

	if row.HasFile != nil {
		exp.HasFile = *row.HasFile
	}

	if row.IsHdr != nil {
		exp.IsHDR = *row.IsHdr
	}

	if row.Is3d != nil {
		exp.Is3D = *row.Is3d
	}

	return exp
}

// rowsToExpeditions converts multiple database rows to domain entities.
func (r *SQLCRepository) rowsToExpeditions(rows []adultdb.QarExpedition) []Expedition {
	result := make([]Expedition, 0, len(rows))
	for i := range rows {
		if exp := r.rowToExpedition(&rows[i]); exp != nil {
			result = append(result, *exp)
		}
	}
	return result
}

// expeditionToCreateParams converts a domain entity to create parameters.
func (r *SQLCRepository) expeditionToCreateParams(exp *Expedition) adultdb.CreateExpeditionParams {
	params := adultdb.CreateExpeditionParams{
		FleetID: exp.FleetID,
		Title:   exp.Title,
		Path:    exp.Path,
	}

	if exp.SortTitle != "" {
		params.SortTitle = &exp.SortTitle
	}

	if exp.Overview != "" {
		params.Overview = &exp.Overview
	}

	if exp.LaunchDate != nil {
		params.LaunchDate = pgtype.Date{
			Time:  *exp.LaunchDate,
			Valid: true,
		}
	}

	if exp.RuntimeTicks > 0 {
		params.RuntimeTicks = &exp.RuntimeTicks
	}

	if exp.PortID != nil {
		params.PortID = pgtype.UUID{Bytes: *exp.PortID, Valid: true}
	}

	if exp.Director != "" {
		params.Director = &exp.Director
	}

	if exp.Series != "" {
		params.Series = &exp.Series
	}

	if exp.Coordinates != "" {
		params.Coordinates = &exp.Coordinates
	}

	if exp.Charter != "" {
		params.Charter = &exp.Charter
	}

	if exp.Registry != "" {
		params.Registry = &exp.Registry
	}

	if exp.WhisparrID != nil {
		id := int32(*exp.WhisparrID)
		params.WhisparrID = &id
	}

	params.HasFile = &exp.HasFile
	params.IsHdr = &exp.IsHDR
	params.Is3d = &exp.Is3D

	return params
}

// expeditionToUpdateParams converts a domain entity to update parameters.
func (r *SQLCRepository) expeditionToUpdateParams(exp *Expedition) adultdb.UpdateExpeditionParams {
	params := adultdb.UpdateExpeditionParams{
		ID:      exp.ID,
		FleetID: exp.FleetID,
		Title:   exp.Title,
		Path:    exp.Path,
	}

	if exp.SortTitle != "" {
		params.SortTitle = &exp.SortTitle
	}

	if exp.Overview != "" {
		params.Overview = &exp.Overview
	}

	if exp.LaunchDate != nil {
		params.LaunchDate = pgtype.Date{
			Time:  *exp.LaunchDate,
			Valid: true,
		}
	}

	if exp.RuntimeTicks > 0 {
		params.RuntimeTicks = &exp.RuntimeTicks
	}

	if exp.PortID != nil {
		params.PortID = pgtype.UUID{Bytes: *exp.PortID, Valid: true}
	}

	if exp.Director != "" {
		params.Director = &exp.Director
	}

	if exp.Series != "" {
		params.Series = &exp.Series
	}

	if exp.Coordinates != "" {
		params.Coordinates = &exp.Coordinates
	}

	if exp.Charter != "" {
		params.Charter = &exp.Charter
	}

	if exp.Registry != "" {
		params.Registry = &exp.Registry
	}

	if exp.WhisparrID != nil {
		id := int32(*exp.WhisparrID)
		params.WhisparrID = &id
	}

	params.HasFile = &exp.HasFile
	params.IsHdr = &exp.IsHDR
	params.Is3d = &exp.Is3D

	return params
}

// Helper to convert time.Time to pgtype.Date
func timeToPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}
