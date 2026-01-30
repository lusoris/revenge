// Package voyage provides adult scene domain models (QAR obfuscation: scenes → voyages).
package voyage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	adultdb "github.com/lusoris/revenge/internal/content/qar/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// ErrVoyageNotFound is returned when a voyage cannot be found.
var ErrVoyageNotFound = errors.New("voyage not found")

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed voyage repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.voyage")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Voyage, error) {
	row, err := r.queries.GetVoyageByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVoyageNotFound
		}
		return nil, err
	}
	return r.rowToVoyage(&row), nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Voyage, error) {
	rows, err := r.queries.ListVoyages(ctx, adultdb.ListVoyagesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToVoyages(rows), nil
}

func (r *SQLCRepository) ListByFleet(ctx context.Context, fleetID uuid.UUID, limit, offset int) ([]Voyage, error) {
	rows, err := r.queries.ListVoyagesByFleet(ctx, adultdb.ListVoyagesByFleetParams{
		FleetID: fleetID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToVoyages(rows), nil
}

func (r *SQLCRepository) ListByPort(ctx context.Context, portID uuid.UUID, limit, offset int) ([]Voyage, error) {
	rows, err := r.queries.ListVoyagesByPort(ctx, adultdb.ListVoyagesByPortParams{
		PortID: pgtype.UUID{Bytes: portID, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToVoyages(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, voyage *Voyage) error {
	params := r.voyageToCreateParams(voyage)
	row, err := r.queries.CreateVoyage(ctx, params)
	if err != nil {
		return err
	}
	voyage.ID = row.ID
	voyage.CreatedAt = row.CreatedAt
	voyage.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, voyage *Voyage) error {
	params := r.voyageToUpdateParams(voyage)
	_, err := r.queries.UpdateVoyage(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrVoyageNotFound
		}
		return err
	}
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteVoyage(ctx, id)
}

func (r *SQLCRepository) GetByPath(ctx context.Context, path string) (*Voyage, error) {
	row, err := r.queries.GetVoyageByPath(ctx, path)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVoyageNotFound
		}
		return nil, err
	}
	return r.rowToVoyage(&row), nil
}

func (r *SQLCRepository) GetByOshash(ctx context.Context, oshash string) (*Voyage, error) {
	row, err := r.queries.GetVoyageByOshash(ctx, &oshash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVoyageNotFound
		}
		return nil, err
	}
	return r.rowToVoyage(&row), nil
}

func (r *SQLCRepository) GetByCoordinates(ctx context.Context, coordinates string) (*Voyage, error) {
	row, err := r.queries.GetVoyageByCoordinates(ctx, &coordinates)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVoyageNotFound
		}
		return nil, err
	}
	return r.rowToVoyage(&row), nil
}

func (r *SQLCRepository) GetByCharter(ctx context.Context, charter string) (*Voyage, error) {
	row, err := r.queries.GetVoyageByCharter(ctx, &charter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVoyageNotFound
		}
		return nil, err
	}
	return r.rowToVoyage(&row), nil
}

func (r *SQLCRepository) CountByFleet(ctx context.Context, fleetID uuid.UUID) (int64, error) {
	return r.queries.CountVoyagesByFleet(ctx, fleetID)
}

func (r *SQLCRepository) Search(ctx context.Context, query string, limit, offset int) ([]Voyage, error) {
	rows, err := r.queries.SearchVoyages(ctx, adultdb.SearchVoyagesParams{
		Column1: &query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToVoyages(rows), nil
}

// rowToVoyage converts a database row to a domain entity.
func (r *SQLCRepository) rowToVoyage(row *adultdb.QarVoyage) *Voyage {
	if row == nil {
		return nil
	}

	v := &Voyage{
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
		v.SortTitle = *row.SortTitle
	}

	if row.Overview != nil {
		v.Overview = *row.Overview
	}

	if row.LaunchDate.Valid {
		t := row.LaunchDate.Time
		v.LaunchDate = &t
	}

	if row.Distance != nil {
		v.Distance = int(*row.Distance)
	}

	if row.PortID.Valid {
		portID := uuid.UUID(row.PortID.Bytes)
		v.PortID = &portID
	}

	if row.Coordinates != nil {
		v.Coordinates = *row.Coordinates
	}

	if row.Oshash != nil {
		v.Oshash = *row.Oshash
	}

	if row.Md5 != nil {
		v.MD5 = *row.Md5
	}

	if row.CoverPath != nil {
		v.CoverPath = *row.CoverPath
	}

	if row.Charter != nil {
		v.Charter = *row.Charter
	}

	if row.Registry != nil {
		v.Registry = *row.Registry
	}

	if row.StashID != nil {
		v.StashID = *row.StashID
	}

	if row.WhisparrID != nil {
		id := int(*row.WhisparrID)
		v.WhisparrID = &id
	}

	return v
}

// rowsToVoyages converts multiple database rows to domain entities.
func (r *SQLCRepository) rowsToVoyages(rows []adultdb.QarVoyage) []Voyage {
	result := make([]Voyage, 0, len(rows))
	for i := range rows {
		if v := r.rowToVoyage(&rows[i]); v != nil {
			result = append(result, *v)
		}
	}
	return result
}

// voyageToCreateParams converts a domain entity to create parameters.
func (r *SQLCRepository) voyageToCreateParams(v *Voyage) adultdb.CreateVoyageParams {
	params := adultdb.CreateVoyageParams{
		FleetID: v.FleetID,
		Title:   v.Title,
		Path:    v.Path,
	}

	if v.SortTitle != "" {
		params.SortTitle = &v.SortTitle
	}

	if v.Overview != "" {
		params.Overview = &v.Overview
	}

	if v.LaunchDate != nil {
		params.LaunchDate = pgtype.Date{
			Time:  *v.LaunchDate,
			Valid: true,
		}
	}

	if v.Distance > 0 {
		dist := int32(v.Distance)
		params.Distance = &dist
	}

	if v.PortID != nil {
		params.PortID = pgtype.UUID{Bytes: *v.PortID, Valid: true}
	}

	if v.Coordinates != "" {
		params.Coordinates = &v.Coordinates
	}

	if v.Oshash != "" {
		params.Oshash = &v.Oshash
	}

	if v.MD5 != "" {
		params.Md5 = &v.MD5
	}

	if v.CoverPath != "" {
		params.CoverPath = &v.CoverPath
	}

	if v.Charter != "" {
		params.Charter = &v.Charter
	}

	if v.Registry != "" {
		params.Registry = &v.Registry
	}

	if v.StashID != "" {
		params.StashID = &v.StashID
	}

	if v.WhisparrID != nil {
		id := int32(*v.WhisparrID)
		params.WhisparrID = &id
	}

	return params
}

// voyageToUpdateParams converts a domain entity to update parameters.
func (r *SQLCRepository) voyageToUpdateParams(v *Voyage) adultdb.UpdateVoyageParams {
	params := adultdb.UpdateVoyageParams{
		ID:      v.ID,
		FleetID: v.FleetID,
		Title:   v.Title,
		Path:    v.Path,
	}

	if v.SortTitle != "" {
		params.SortTitle = &v.SortTitle
	}

	if v.Overview != "" {
		params.Overview = &v.Overview
	}

	if v.LaunchDate != nil {
		params.LaunchDate = pgtype.Date{
			Time:  *v.LaunchDate,
			Valid: true,
		}
	}

	if v.Distance > 0 {
		dist := int32(v.Distance)
		params.Distance = &dist
	}

	if v.PortID != nil {
		params.PortID = pgtype.UUID{Bytes: *v.PortID, Valid: true}
	}

	if v.Coordinates != "" {
		params.Coordinates = &v.Coordinates
	}

	if v.Oshash != "" {
		params.Oshash = &v.Oshash
	}

	if v.MD5 != "" {
		params.Md5 = &v.MD5
	}

	if v.CoverPath != "" {
		params.CoverPath = &v.CoverPath
	}

	if v.Charter != "" {
		params.Charter = &v.Charter
	}

	if v.Registry != "" {
		params.Registry = &v.Registry
	}

	if v.StashID != "" {
		params.StashID = &v.StashID
	}

	if v.WhisparrID != nil {
		id := int32(*v.WhisparrID)
		params.WhisparrID = &id
	}

	return params
}
