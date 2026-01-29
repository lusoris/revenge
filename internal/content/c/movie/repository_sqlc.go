package movie

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	adultdb "github.com/lusoris/revenge/internal/content/c/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

type pgRepository struct {
	queries *adultdb.Queries
}

// NewRepository creates a new adult movie repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgRepository{queries: adultdb.New(pool)}
}

func (r *pgRepository) GetByID(ctx context.Context, id uuid.UUID) (*Movie, error) {
	row, err := r.queries.GetAdultMovieByID(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movieFromRow(row), nil
}

func (r *pgRepository) List(ctx context.Context, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListAdultMovies(ctx, adultdb.ListAdultMoviesParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}
	return moviesFromRows(rows), nil
}

func (r *pgRepository) ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListAdultMoviesByLibrary(ctx, adultdb.ListAdultMoviesByLibraryParams{
		LibraryID: libraryID,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}
	return moviesFromRows(rows), nil
}

func (r *pgRepository) Create(ctx context.Context, movie *Movie) error {
	_, err := r.queries.CreateAdultMovie(ctx, adultdb.CreateAdultMovieParams{
		LibraryID:     movie.LibraryID,
		Title:         movie.Title,
		SortTitle:     stringPtr(movie.SortTitle),
		OriginalTitle: nil,
		Overview:      stringPtr(movie.Overview),
		ReleaseDate:   pgDateFromTime(movie.ReleaseDate),
		RuntimeTicks:  int64Ptr(movie.RuntimeTicks),
		StudioID:      pgUUIDFromPtr(movie.StudioID),
		Director:      stringPtr(movie.Director),
		Series:        stringPtr(movie.Series),
		Path:          movie.Path,
		SizeBytes:     nil,
		Container:     nil,
		VideoCodec:    nil,
		AudioCodec:    nil,
		Resolution:    nil,
		Phash:         nil,
		Oshash:        nil,
		WhisparrID:    nil,
		StashdbID:     nil,
		TpdbID:        nil,
		HasFile:       boolPtr(true),
		IsHdr:         boolPtr(false),
		Is3d:          boolPtr(false),
	})
	return err
}

func (r *pgRepository) Update(ctx context.Context, movie *Movie) error {
	_, err := r.queries.UpdateAdultMovie(ctx, adultdb.UpdateAdultMovieParams{
		ID:            movie.ID,
		LibraryID:     movie.LibraryID,
		Title:         movie.Title,
		SortTitle:     stringPtr(movie.SortTitle),
		OriginalTitle: nil,
		Overview:      stringPtr(movie.Overview),
		ReleaseDate:   pgDateFromTime(movie.ReleaseDate),
		RuntimeTicks:  int64Ptr(movie.RuntimeTicks),
		StudioID:      pgUUIDFromPtr(movie.StudioID),
		Director:      stringPtr(movie.Director),
		Series:        stringPtr(movie.Series),
		Path:          movie.Path,
		SizeBytes:     nil,
		Container:     nil,
		VideoCodec:    nil,
		AudioCodec:    nil,
		Resolution:    nil,
		Phash:         nil,
		Oshash:        nil,
		WhisparrID:    nil,
		StashdbID:     nil,
		TpdbID:        nil,
		HasFile:       boolPtr(true),
		IsHdr:         boolPtr(false),
		Is3d:          boolPtr(false),
	})
	if err != nil {
		if isNoRows(err) {
			return ErrMovieNotFound
		}
		return err
	}
	return nil
}

func (r *pgRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteAdultMovie(ctx, id); err != nil {
		if isNoRows(err) {
			return ErrMovieNotFound
		}
		return err
	}
	return nil
}

func moviesFromRows(rows []adultdb.CMovie) []*Movie {
	movies := make([]*Movie, 0, len(rows))
	for _, row := range rows {
		movies = append(movies, movieFromRow(row))
	}
	return movies
}

func movieFromRow(row adultdb.CMovie) *Movie {
	return &Movie{
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
		ReleaseDate:  timeFromPgDate(row.ReleaseDate),
		RuntimeTicks: int64OrZero(row.RuntimeTicks),
		Overview:     stringOrEmpty(row.Overview),
		StudioID:     uuidFromPg(row.StudioID),
		Director:     stringOrEmpty(row.Director),
		Series:       stringOrEmpty(row.Series),
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

func int64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func int64OrZero(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func boolPtr(v bool) *bool {
	return &v
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
