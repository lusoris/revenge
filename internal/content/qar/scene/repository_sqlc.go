package scene

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	adultdb "github.com/lusoris/revenge/internal/content/qar/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

type pgRepository struct {
	queries *adultdb.Queries
}

// NewRepository creates a new adult scene repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{queries: adultdb.New(pool)}
}

func (r *pgRepository) GetByID(ctx context.Context, id uuid.UUID) (*Scene, error) {
	row, err := r.queries.GetAdultSceneByID(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrSceneNotFound
		}
		return nil, err
	}
	return sceneFromRow(row), nil
}

func (r *pgRepository) List(ctx context.Context, params ListParams) ([]*Scene, error) {
	rows, err := r.queries.ListAdultScenes(ctx, adultdb.ListAdultScenesParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}
	return scenesFromRows(rows), nil
}

func (r *pgRepository) ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Scene, error) {
	rows, err := r.queries.ListAdultScenesByLibrary(ctx, adultdb.ListAdultScenesByLibraryParams{
		LibraryID: libraryID,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}
	return scenesFromRows(rows), nil
}

func (r *pgRepository) Create(ctx context.Context, scene *Scene) error {
	_, err := r.queries.CreateAdultScene(ctx, adultdb.CreateAdultSceneParams{
		LibraryID:      scene.LibraryID,
		Title:          scene.Title,
		SortTitle:      stringPtr(scene.SortTitle),
		Overview:       stringPtr(scene.Overview),
		ReleaseDate:    pgDateFromTime(scene.ReleaseDate),
		RuntimeMinutes: int32Ptr(scene.RuntimeMinutes),
		StudioID:       pgUUIDFromPtr(scene.StudioID),
		WhisparrID:     nil,
		StashID:        nil,
		StashdbID:      nil,
		TpdbID:         nil,
		Path:           scene.Path,
		SizeBytes:      nil,
		VideoCodec:     nil,
		AudioCodec:     nil,
		Resolution:     nil,
		Oshash:         nil,
		Phash:          nil,
		Md5:            nil,
		CoverPath:      nil,
	})
	return err
}

func (r *pgRepository) Update(ctx context.Context, scene *Scene) error {
	_, err := r.queries.UpdateAdultScene(ctx, adultdb.UpdateAdultSceneParams{
		ID:             scene.ID,
		LibraryID:      scene.LibraryID,
		Title:          scene.Title,
		SortTitle:      stringPtr(scene.SortTitle),
		Overview:       stringPtr(scene.Overview),
		ReleaseDate:    pgDateFromTime(scene.ReleaseDate),
		RuntimeMinutes: int32Ptr(scene.RuntimeMinutes),
		StudioID:       pgUUIDFromPtr(scene.StudioID),
		WhisparrID:     nil,
		StashID:        nil,
		StashdbID:      nil,
		TpdbID:         nil,
		Path:           scene.Path,
		SizeBytes:      nil,
		VideoCodec:     nil,
		AudioCodec:     nil,
		Resolution:     nil,
		Oshash:         nil,
		Phash:          nil,
		Md5:            nil,
		CoverPath:      nil,
	})
	if err != nil {
		if isNoRows(err) {
			return ErrSceneNotFound
		}
		return err
	}
	return nil
}

func (r *pgRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteAdultScene(ctx, id); err != nil {
		if isNoRows(err) {
			return ErrSceneNotFound
		}
		return err
	}
	return nil
}

func scenesFromRows(rows []adultdb.QarScene) []*Scene {
	scenes := make([]*Scene, 0, len(rows))
	for _, row := range rows {
		scenes = append(scenes, sceneFromRow(row))
	}
	return scenes
}

func sceneFromRow(row adultdb.QarScene) *Scene {
	return &Scene{
		ContentEntity: shared.ContentEntity{
			BaseEntity: shared.BaseEntity{
				ID:        row.ID,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			LibraryID: row.LibraryID,
			Path:      row.Path,
			Title:     row.Title,
			SortTitle: stringOrEmpty(row.SortTitle),
		},
		ReleaseDate:    timeFromPgDate(row.ReleaseDate),
		RuntimeMinutes: intOrZero(row.RuntimeMinutes),
		Overview:       stringOrEmpty(row.Overview),
		StudioID:       uuidFromPg(row.StudioID),
	}
}

func timeFromPgDate(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

func pgDateFromTime(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func uuidFromPg(v pgtype.UUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	id := uuid.UUID(v.Bytes)
	return &id
}

func pgUUIDFromPtr(v *uuid.UUID) pgtype.UUID {
	if v == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *v, Valid: true}
}

func stringOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func int32Ptr(v int) *int32 {
	if v == 0 {
		return nil
	}
	val := int32(v)
	return &val
}

func intOrZero(v *int32) int {
	if v == nil {
		return 0
	}
	return int(*v)
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
