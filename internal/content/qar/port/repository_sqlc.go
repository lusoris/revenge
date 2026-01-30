// Package port provides adult studio domain models (QAR obfuscation: studios → ports).
package port

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

// ErrPortNotFound is returned when a port cannot be found.
var ErrPortNotFound = errors.New("port not found")

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed port repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.port")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Port, error) {
	row, err := r.queries.GetPortByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPortNotFound
		}
		return nil, err
	}
	return r.rowToPort(&row), nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Port, error) {
	rows, err := r.queries.ListPorts(ctx, adultdb.ListPortsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToPorts(rows), nil
}

func (r *SQLCRepository) ListRoot(ctx context.Context, limit, offset int) ([]Port, error) {
	rows, err := r.queries.ListRootPorts(ctx, adultdb.ListRootPortsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToPorts(rows), nil
}

func (r *SQLCRepository) ListChildren(ctx context.Context, parentID uuid.UUID) ([]Port, error) {
	rows, err := r.queries.ListPortChildren(ctx, pgtype.UUID{Bytes: parentID, Valid: true})
	if err != nil {
		return nil, err
	}
	return r.rowsToPorts(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, port *Port) error {
	params := r.portToCreateParams(port)
	row, err := r.queries.CreatePort(ctx, params)
	if err != nil {
		return err
	}
	port.ID = row.ID
	port.CreatedAt = row.CreatedAt
	port.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, port *Port) error {
	params := r.portToUpdateParams(port)
	_, err := r.queries.UpdatePort(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPortNotFound
		}
		return err
	}
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePort(ctx, id)
}

func (r *SQLCRepository) GetByStashDBID(ctx context.Context, stashdbID string) (*Port, error) {
	row, err := r.queries.GetPortByStashDBID(ctx, &stashdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPortNotFound
		}
		return nil, err
	}
	return r.rowToPort(&row), nil
}

func (r *SQLCRepository) GetByTPDBID(ctx context.Context, tpdbID string) (*Port, error) {
	row, err := r.queries.GetPortByTPDBID(ctx, &tpdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPortNotFound
		}
		return nil, err
	}
	return r.rowToPort(&row), nil
}

func (r *SQLCRepository) Search(ctx context.Context, query string, limit, offset int) ([]Port, error) {
	rows, err := r.queries.SearchPorts(ctx, adultdb.SearchPortsParams{
		Column1: &query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToPorts(rows), nil
}

func (r *SQLCRepository) CountExpeditions(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.CountExpeditionsByPort(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *SQLCRepository) CountVoyages(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.CountVoyagesByPort(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

// rowToPort converts a database row to a domain entity.
func (r *SQLCRepository) rowToPort(row *adultdb.QarPort) *Port {
	if row == nil {
		return nil
	}

	p := &Port{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.ParentID.Valid {
		parentID := uuid.UUID(row.ParentID.Bytes)
		p.ParentID = &parentID
	}

	if row.StashdbID != nil {
		p.StashDBID = *row.StashdbID
	}

	if row.TpdbID != nil {
		p.TPDBID = *row.TpdbID
	}

	if row.Url != nil {
		p.URL = *row.Url
	}

	if row.LogoPath != nil {
		p.LogoPath = *row.LogoPath
	}

	return p
}

// rowsToPorts converts multiple database rows to domain entities.
func (r *SQLCRepository) rowsToPorts(rows []adultdb.QarPort) []Port {
	result := make([]Port, 0, len(rows))
	for i := range rows {
		if p := r.rowToPort(&rows[i]); p != nil {
			result = append(result, *p)
		}
	}
	return result
}

// portToCreateParams converts a domain entity to create parameters.
func (r *SQLCRepository) portToCreateParams(p *Port) adultdb.CreatePortParams {
	params := adultdb.CreatePortParams{
		Name: p.Name,
	}

	if p.ParentID != nil {
		params.ParentID = pgtype.UUID{Bytes: *p.ParentID, Valid: true}
	}

	if p.StashDBID != "" {
		params.StashdbID = &p.StashDBID
	}

	if p.TPDBID != "" {
		params.TpdbID = &p.TPDBID
	}

	if p.URL != "" {
		params.Url = &p.URL
	}

	if p.LogoPath != "" {
		params.LogoPath = &p.LogoPath
	}

	return params
}

// portToUpdateParams converts a domain entity to update parameters.
func (r *SQLCRepository) portToUpdateParams(p *Port) adultdb.UpdatePortParams {
	params := adultdb.UpdatePortParams{
		ID:   p.ID,
		Name: p.Name,
	}

	if p.ParentID != nil {
		params.ParentID = pgtype.UUID{Bytes: *p.ParentID, Valid: true}
	}

	if p.StashDBID != "" {
		params.StashdbID = &p.StashDBID
	}

	if p.TPDBID != "" {
		params.TpdbID = &p.TPDBID
	}

	if p.URL != "" {
		params.Url = &p.URL
	}

	if p.LogoPath != "" {
		params.LogoPath = &p.LogoPath
	}

	return params
}
