// Package crew provides adult performer domain models (QAR obfuscation: performers → crew).
package crew

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

// ErrCrewNotFound is returned when a crew member cannot be found.
var ErrCrewNotFound = errors.New("crew not found")

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed crew repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.crew")),
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*Crew, error) {
	row, err := r.queries.GetCrewByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCrewNotFound
		}
		return nil, err
	}
	return r.rowToCrew(&row), nil
}

func (r *SQLCRepository) List(ctx context.Context, limit, offset int) ([]Crew, error) {
	rows, err := r.queries.ListCrew(ctx, adultdb.ListCrewParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToCrews(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, crew *Crew) error {
	params := r.crewToCreateParams(crew)
	row, err := r.queries.CreateCrew(ctx, params)
	if err != nil {
		return err
	}
	crew.ID = row.ID
	crew.CreatedAt = row.CreatedAt
	crew.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *SQLCRepository) Update(ctx context.Context, crew *Crew) error {
	params := r.crewToUpdateParams(crew)
	_, err := r.queries.UpdateCrew(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCrewNotFound
		}
		return err
	}
	return nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCrew(ctx, id)
}

func (r *SQLCRepository) GetByCharter(ctx context.Context, charter string) (*Crew, error) {
	row, err := r.queries.GetCrewByCharter(ctx, &charter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCrewNotFound
		}
		return nil, err
	}
	return r.rowToCrew(&row), nil
}

func (r *SQLCRepository) GetByRegistry(ctx context.Context, registry string) (*Crew, error) {
	row, err := r.queries.GetCrewByRegistry(ctx, &registry)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCrewNotFound
		}
		return nil, err
	}
	return r.rowToCrew(&row), nil
}

func (r *SQLCRepository) Search(ctx context.Context, query string, limit, offset int) ([]Crew, error) {
	rows, err := r.queries.SearchCrew(ctx, adultdb.SearchCrewParams{
		Column1: &query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToCrews(rows), nil
}

func (r *SQLCRepository) ListNames(ctx context.Context, crewID uuid.UUID) ([]CrewName, error) {
	rows, err := r.queries.ListCrewNames(ctx, crewID)
	if err != nil {
		return nil, err
	}
	result := make([]CrewName, 0, len(rows))
	for _, row := range rows {
		result = append(result, CrewName{
			CrewID: row.CrewID,
			Name:   row.Name,
		})
	}
	return result, nil
}

func (r *SQLCRepository) AddName(ctx context.Context, crewID uuid.UUID, name string) error {
	return r.queries.AddCrewName(ctx, adultdb.AddCrewNameParams{
		CrewID: crewID,
		Name:   name,
	})
}

func (r *SQLCRepository) RemoveName(ctx context.Context, crewID uuid.UUID, name string) error {
	return r.queries.RemoveCrewName(ctx, adultdb.RemoveCrewNameParams{
		CrewID: crewID,
		Name:   name,
	})
}

func (r *SQLCRepository) ListPortraits(ctx context.Context, crewID uuid.UUID) ([]CrewPortrait, error) {
	rows, err := r.queries.ListCrewPortraits(ctx, pgtype.UUID{Bytes: crewID, Valid: true})
	if err != nil {
		return nil, err
	}
	result := make([]CrewPortrait, 0, len(rows))
	for _, row := range rows {
		p := CrewPortrait{
			ID:        row.ID,
			Path:      row.Path,
			CreatedAt: row.CreatedAt,
		}
		if row.CrewID.Valid {
			p.CrewID = uuid.UUID(row.CrewID.Bytes)
		}
		if row.Type != nil {
			p.Type = *row.Type
		}
		if row.Source != nil {
			p.Source = *row.Source
		}
		if row.PrimaryImage != nil {
			p.PrimaryImage = *row.PrimaryImage
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *SQLCRepository) AddPortrait(ctx context.Context, portrait *CrewPortrait) error {
	params := adultdb.AddCrewPortraitParams{
		CrewID: pgtype.UUID{Bytes: portrait.CrewID, Valid: true},
		Path:   portrait.Path,
	}
	if portrait.Type != "" {
		params.Type = &portrait.Type
	}
	if portrait.Source != "" {
		params.Source = &portrait.Source
	}
	params.PrimaryImage = &portrait.PrimaryImage

	row, err := r.queries.AddCrewPortrait(ctx, params)
	if err != nil {
		return err
	}
	portrait.ID = row.ID
	portrait.CreatedAt = row.CreatedAt
	return nil
}

func (r *SQLCRepository) SetPrimaryPortrait(ctx context.Context, crewID, portraitID uuid.UUID) error {
	return r.queries.SetPrimaryPortrait(ctx, adultdb.SetPrimaryPortraitParams{
		CrewID: pgtype.UUID{Bytes: crewID, Valid: true},
		ID:     portraitID,
	})
}

func (r *SQLCRepository) ListExpeditionCrew(ctx context.Context, expeditionID uuid.UUID) ([]Crew, error) {
	rows, err := r.queries.ListExpeditionCrew(ctx, expeditionID)
	if err != nil {
		return nil, err
	}
	return r.rowsToCrews(rows), nil
}

func (r *SQLCRepository) ListVoyageCrew(ctx context.Context, voyageID uuid.UUID) ([]Crew, error) {
	rows, err := r.queries.ListVoyageCrew(ctx, voyageID)
	if err != nil {
		return nil, err
	}
	return r.rowsToCrews(rows), nil
}

func (r *SQLCRepository) AddExpeditionCrew(ctx context.Context, expeditionID, crewID uuid.UUID, characterName string) error {
	params := adultdb.AddExpeditionCrewParams{
		ExpeditionID: expeditionID,
		CrewID:       crewID,
	}
	if characterName != "" {
		params.CharacterName = &characterName
	}
	return r.queries.AddExpeditionCrew(ctx, params)
}

func (r *SQLCRepository) AddVoyageCrew(ctx context.Context, voyageID, crewID uuid.UUID, role string) error {
	params := adultdb.AddVoyageCrewParams{
		VoyageID: voyageID,
		CrewID:   crewID,
	}
	if role != "" {
		params.Role = &role
	}
	return r.queries.AddVoyageCrew(ctx, params)
}

func (r *SQLCRepository) RemoveExpeditionCrew(ctx context.Context, expeditionID, crewID uuid.UUID) error {
	return r.queries.RemoveExpeditionCrew(ctx, adultdb.RemoveExpeditionCrewParams{
		ExpeditionID: expeditionID,
		CrewID:       crewID,
	})
}

func (r *SQLCRepository) RemoveVoyageCrew(ctx context.Context, voyageID, crewID uuid.UUID) error {
	return r.queries.RemoveVoyageCrew(ctx, adultdb.RemoveVoyageCrewParams{
		VoyageID: voyageID,
		CrewID:   crewID,
	})
}

// rowToCrew converts a database row to a domain entity.
func (r *SQLCRepository) rowToCrew(row *adultdb.QarCrew) *Crew {
	if row == nil {
		return nil
	}

	c := &Crew{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.Disambiguation != nil {
		c.Disambiguation = *row.Disambiguation
	}
	if row.Gender != nil {
		c.Gender = *row.Gender
	}
	if row.Christening.Valid {
		t := row.Christening.Time
		c.Christening = &t
	}
	if row.DeathDate.Valid {
		t := row.DeathDate.Time
		c.DeathDate = &t
	}
	if row.BirthCity != nil {
		c.BirthCity = *row.BirthCity
	}
	if row.Origin != nil {
		c.Origin = *row.Origin
	}
	if row.Nationality != nil {
		c.Nationality = *row.Nationality
	}
	if row.Rigging != nil {
		c.Rigging = *row.Rigging
	}
	if row.Compass != nil {
		c.Compass = *row.Compass
	}
	if row.HeightCm != nil {
		h := int(*row.HeightCm)
		c.HeightCM = &h
	}
	if row.WeightKg != nil {
		w := int(*row.WeightKg)
		c.WeightKG = &w
	}
	if row.Measurements != nil {
		c.Measurements = *row.Measurements
	}
	if row.CupSize != nil {
		c.CupSize = *row.CupSize
	}
	if row.BreastType != nil {
		c.BreastType = *row.BreastType
	}
	if row.Markings != nil {
		c.Markings = *row.Markings
	}
	if row.Anchors != nil {
		c.Anchors = *row.Anchors
	}
	if row.MaidenVoyage != nil {
		mv := int(*row.MaidenVoyage)
		c.MaidenVoyage = &mv
	}
	if row.LastPort != nil {
		lp := int(*row.LastPort)
		c.LastPort = &lp
	}
	if row.Bio != nil {
		c.Bio = *row.Bio
	}
	if row.StashID != nil {
		c.StashID = *row.StashID
	}
	if row.Charter != nil {
		c.Charter = *row.Charter
	}
	if row.Registry != nil {
		c.Registry = *row.Registry
	}
	if row.Manifest != nil {
		c.Manifest = *row.Manifest
	}
	if row.Twitter != nil {
		c.Twitter = *row.Twitter
	}
	if row.Instagram != nil {
		c.Instagram = *row.Instagram
	}
	if row.ImagePath != nil {
		c.ImagePath = *row.ImagePath
	}

	return c
}

// rowsToCrews converts multiple database rows to domain entities.
func (r *SQLCRepository) rowsToCrews(rows []adultdb.QarCrew) []Crew {
	result := make([]Crew, 0, len(rows))
	for i := range rows {
		if c := r.rowToCrew(&rows[i]); c != nil {
			result = append(result, *c)
		}
	}
	return result
}

// crewToCreateParams converts a domain entity to create parameters.
func (r *SQLCRepository) crewToCreateParams(c *Crew) adultdb.CreateCrewParams {
	params := adultdb.CreateCrewParams{
		Name: c.Name,
	}

	if c.Disambiguation != "" {
		params.Disambiguation = &c.Disambiguation
	}
	if c.Gender != "" {
		params.Gender = &c.Gender
	}
	if c.Christening != nil {
		params.Christening = pgtype.Date{Time: *c.Christening, Valid: true}
	}
	if c.DeathDate != nil {
		params.DeathDate = pgtype.Date{Time: *c.DeathDate, Valid: true}
	}
	if c.BirthCity != "" {
		params.BirthCity = &c.BirthCity
	}
	if c.Origin != "" {
		params.Origin = &c.Origin
	}
	if c.Nationality != "" {
		params.Nationality = &c.Nationality
	}
	if c.Rigging != "" {
		params.Rigging = &c.Rigging
	}
	if c.Compass != "" {
		params.Compass = &c.Compass
	}
	if c.HeightCM != nil {
		h := int32(*c.HeightCM)
		params.HeightCm = &h
	}
	if c.WeightKG != nil {
		w := int32(*c.WeightKG)
		params.WeightKg = &w
	}
	if c.Measurements != "" {
		params.Measurements = &c.Measurements
	}
	if c.CupSize != "" {
		params.CupSize = &c.CupSize
	}
	if c.BreastType != "" {
		params.BreastType = &c.BreastType
	}
	if c.Markings != "" {
		params.Markings = &c.Markings
	}
	if c.Anchors != "" {
		params.Anchors = &c.Anchors
	}
	if c.MaidenVoyage != nil {
		mv := int32(*c.MaidenVoyage)
		params.MaidenVoyage = &mv
	}
	if c.LastPort != nil {
		lp := int32(*c.LastPort)
		params.LastPort = &lp
	}
	if c.Bio != "" {
		params.Bio = &c.Bio
	}
	if c.StashID != "" {
		params.StashID = &c.StashID
	}
	if c.Charter != "" {
		params.Charter = &c.Charter
	}
	if c.Registry != "" {
		params.Registry = &c.Registry
	}
	if c.Manifest != "" {
		params.Manifest = &c.Manifest
	}
	if c.Twitter != "" {
		params.Twitter = &c.Twitter
	}
	if c.Instagram != "" {
		params.Instagram = &c.Instagram
	}
	if c.ImagePath != "" {
		params.ImagePath = &c.ImagePath
	}

	return params
}

// crewToUpdateParams converts a domain entity to update parameters.
func (r *SQLCRepository) crewToUpdateParams(c *Crew) adultdb.UpdateCrewParams {
	params := adultdb.UpdateCrewParams{
		ID:   c.ID,
		Name: c.Name,
	}

	if c.Disambiguation != "" {
		params.Disambiguation = &c.Disambiguation
	}
	if c.Gender != "" {
		params.Gender = &c.Gender
	}
	if c.Christening != nil {
		params.Christening = pgtype.Date{Time: *c.Christening, Valid: true}
	}
	if c.DeathDate != nil {
		params.DeathDate = pgtype.Date{Time: *c.DeathDate, Valid: true}
	}
	if c.BirthCity != "" {
		params.BirthCity = &c.BirthCity
	}
	if c.Origin != "" {
		params.Origin = &c.Origin
	}
	if c.Nationality != "" {
		params.Nationality = &c.Nationality
	}
	if c.Rigging != "" {
		params.Rigging = &c.Rigging
	}
	if c.Compass != "" {
		params.Compass = &c.Compass
	}
	if c.HeightCM != nil {
		h := int32(*c.HeightCM)
		params.HeightCm = &h
	}
	if c.WeightKG != nil {
		w := int32(*c.WeightKG)
		params.WeightKg = &w
	}
	if c.Measurements != "" {
		params.Measurements = &c.Measurements
	}
	if c.CupSize != "" {
		params.CupSize = &c.CupSize
	}
	if c.BreastType != "" {
		params.BreastType = &c.BreastType
	}
	if c.Markings != "" {
		params.Markings = &c.Markings
	}
	if c.Anchors != "" {
		params.Anchors = &c.Anchors
	}
	if c.MaidenVoyage != nil {
		mv := int32(*c.MaidenVoyage)
		params.MaidenVoyage = &mv
	}
	if c.LastPort != nil {
		lp := int32(*c.LastPort)
		params.LastPort = &lp
	}
	if c.Bio != "" {
		params.Bio = &c.Bio
	}
	if c.StashID != "" {
		params.StashID = &c.StashID
	}
	if c.Charter != "" {
		params.Charter = &c.Charter
	}
	if c.Registry != "" {
		params.Registry = &c.Registry
	}
	if c.Manifest != "" {
		params.Manifest = &c.Manifest
	}
	if c.Twitter != "" {
		params.Twitter = &c.Twitter
	}
	if c.Instagram != "" {
		params.Instagram = &c.Instagram
	}
	if c.ImagePath != "" {
		params.ImagePath = &c.ImagePath
	}

	return params
}
