// Package crew provides adult performer domain models (QAR obfuscation: performers → crew).
package crew

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed crew repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		pool:   pool,
		logger: logger.With(slog.String("repository", "qar.crew")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) Create(ctx context.Context, crew *Crew) error {
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, crew *Crew) error {
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *SQLCRepository) GetByCharter(ctx context.Context, charter string) (*Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) GetByRegistry(ctx context.Context, registry string) (*Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) Search(ctx context.Context, query string, limit, offset int) ([]Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) ListNames(ctx context.Context, crewID uuid.UUID) ([]CrewName, error) {
	return nil, nil
}

func (r *SQLCRepository) AddName(ctx context.Context, crewID uuid.UUID, name string) error {
	return nil
}

func (r *SQLCRepository) RemoveName(ctx context.Context, crewID uuid.UUID, name string) error {
	return nil
}

func (r *SQLCRepository) ListPortraits(ctx context.Context, crewID uuid.UUID) ([]CrewPortrait, error) {
	return nil, nil
}

func (r *SQLCRepository) AddPortrait(ctx context.Context, portrait *CrewPortrait) error {
	return nil
}

func (r *SQLCRepository) SetPrimaryPortrait(ctx context.Context, crewID, portraitID uuid.UUID) error {
	return nil
}

func (r *SQLCRepository) ListExpeditionCrew(ctx context.Context, expeditionID uuid.UUID) ([]Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) ListVoyageCrew(ctx context.Context, voyageID uuid.UUID) ([]Crew, error) {
	return nil, nil
}

func (r *SQLCRepository) AddExpeditionCrew(ctx context.Context, expeditionID, crewID uuid.UUID, characterName string) error {
	return nil
}

func (r *SQLCRepository) AddVoyageCrew(ctx context.Context, voyageID, crewID uuid.UUID, role string) error {
	return nil
}

func (r *SQLCRepository) RemoveExpeditionCrew(ctx context.Context, expeditionID, crewID uuid.UUID) error {
	return nil
}

func (r *SQLCRepository) RemoveVoyageCrew(ctx context.Context, voyageID, crewID uuid.UUID) error {
	return nil
}
