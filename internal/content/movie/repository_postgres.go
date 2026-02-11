package movie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lusoris/revenge/internal/content"
	moviedb "github.com/lusoris/revenge/internal/content/movie/db"
	"github.com/lusoris/revenge/internal/util"
)

// postgresRepository implements the Repository interface using PostgreSQL
type postgresRepository struct {
	pool    *pgxpool.Pool
	queries *moviedb.Queries
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{
		pool:    pool,
		queries: moviedb.New(pool),
	}
}

// GetMovie retrieves a movie by ID
func (r *postgresRepository) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
	movie, err := r.queries.GetMovie(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// GetMovieByTMDbID retrieves a movie by TMDb ID
func (r *postgresRepository) GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error) {
	movie, err := r.queries.GetMovieByTMDbID(ctx, &tmdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("failed to get movie by TMDb ID: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// GetMovieByIMDbID retrieves a movie by IMDb ID
func (r *postgresRepository) GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	movie, err := r.queries.GetMovieByIMDbID(ctx, &imdbID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("failed to get movie by IMDb ID: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// GetMovieByRadarrID retrieves a movie by Radarr ID
func (r *postgresRepository) GetMovieByRadarrID(ctx context.Context, radarrID int32) (*Movie, error) {
	movie, err := r.queries.GetMovieByRadarrID(ctx, &radarrID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("failed to get movie by Radarr ID: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// dbMovieToMovie converts a database movie to a domain movie
func dbMovieToMovie(dbMovie moviedb.Movie) *Movie {
	return &Movie{
		ID:                dbMovie.ID,
		TMDbID:            dbMovie.TmdbID,
		IMDbID:            dbMovie.ImdbID,
		Title:             dbMovie.Title,
		OriginalTitle:     dbMovie.OriginalTitle,
		Year:              dbMovie.Year,
		ReleaseDate:       pgDateToTimePtr(dbMovie.ReleaseDate),
		Runtime:           dbMovie.Runtime,
		Overview:          dbMovie.Overview,
		Tagline:           dbMovie.Tagline,
		Status:            dbMovie.Status,
		OriginalLanguage:  dbMovie.OriginalLanguage,
		TitlesI18n:        unmarshalStringMap(dbMovie.TitlesI18n),
		TaglinesI18n:      unmarshalStringMap(dbMovie.TaglinesI18n),
		OverviewsI18n:     unmarshalStringMap(dbMovie.OverviewsI18n),
		AgeRatings:        unmarshalNestedStringMap(dbMovie.AgeRatings),
		ExternalRatings:   unmarshalExternalRatings(dbMovie.ExternalRatings),
		PosterPath:        dbMovie.PosterPath,
		BackdropPath:      dbMovie.BackdropPath,
		TrailerURL:        dbMovie.TrailerUrl,
		VoteAverage:       pgNumericToDecimalPtr(dbMovie.VoteAverage),
		VoteCount:         dbMovie.VoteCount,
		Popularity:        pgNumericToDecimalPtr(dbMovie.Popularity),
		Budget:            dbMovie.Budget,
		Revenue:           dbMovie.Revenue,
		LibraryAddedAt:    dbMovie.LibraryAddedAt,
		MetadataUpdatedAt: pgTimestamptzToTimePtr(dbMovie.MetadataUpdatedAt),
		RadarrID:          dbMovie.RadarrID,
		CreatedAt:         dbMovie.CreatedAt,
		UpdatedAt:         dbMovie.UpdatedAt,
	}
}

// Helper functions for pgtype conversions
func pgDateToTimePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func pgTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func pgNumericToDecimalPtr(n pgtype.Numeric) *decimal.Decimal {
	if !n.Valid {
		return nil
	}
	// Convert pgtype.Numeric to string then to decimal
	var s string
	_ = n.Scan(&s)
	d, _ := decimal.Parse(s)
	return &d
}

// stringToPgDate converts a string date (YYYY-MM-DD) to pgtype.Date
func stringToPgDate(s *string) pgtype.Date {
	if s == nil || *s == "" {
		return pgtype.Date{Valid: false}
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: t, Valid: true}
}

// stringToPgNumeric converts a string number to pgtype.Numeric
func stringToPgNumeric(s *string) pgtype.Numeric {
	if s == nil || *s == "" {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	if err := n.Scan(*s); err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return n
}

// stringToPgTimestamptz converts a string timestamp to pgtype.Timestamptz
func stringToPgTimestamptz(s *string) pgtype.Timestamptz {
	if s == nil || *s == "" {
		return pgtype.Timestamptz{Valid: false}
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// dbWatchedToWatched converts a database movie watched to a domain movie watched
func dbWatchedToWatched(dbWatched moviedb.MovieWatched) *MovieWatched {
	return &MovieWatched{
		ID:              dbWatched.ID,
		UserID:          dbWatched.UserID,
		MovieID:         dbWatched.MovieID,
		ProgressSeconds: dbWatched.ProgressSeconds,
		DurationSeconds: derefInt32(dbWatched.DurationSeconds),
		ProgressPercent: pgNumericToInt32Ptr(dbWatched.ProgressPercent),
		IsCompleted:     derefBool(dbWatched.IsCompleted),
		WatchCount:      derefInt32(dbWatched.WatchCount),
		LastWatchedAt:   dbWatched.LastWatchedAt,
		CreatedAt:       dbWatched.CreatedAt,
		UpdatedAt:       dbWatched.UpdatedAt,
	}
}

// dbContinueWatchingRowToMovie converts a ListContinueWatchingRow to a domain Movie
func dbContinueWatchingRowToMovie(row moviedb.ListContinueWatchingRow) *Movie {
	return &Movie{
		ID:                row.ID,
		TMDbID:            row.TmdbID,
		IMDbID:            row.ImdbID,
		Title:             row.Title,
		OriginalTitle:     row.OriginalTitle,
		Year:              row.Year,
		ReleaseDate:       pgDateToTimePtr(row.ReleaseDate),
		Runtime:           row.Runtime,
		Overview:          row.Overview,
		Tagline:           row.Tagline,
		Status:            row.Status,
		OriginalLanguage:  row.OriginalLanguage,
		TitlesI18n:        unmarshalStringMap(row.TitlesI18n),
		TaglinesI18n:      unmarshalStringMap(row.TaglinesI18n),
		OverviewsI18n:     unmarshalStringMap(row.OverviewsI18n),
		AgeRatings:        unmarshalNestedStringMap(row.AgeRatings),
		PosterPath:        row.PosterPath,
		BackdropPath:      row.BackdropPath,
		TrailerURL:        row.TrailerUrl,
		VoteAverage:       pgNumericToDecimalPtr(row.VoteAverage),
		VoteCount:         row.VoteCount,
		Popularity:        pgNumericToDecimalPtr(row.Popularity),
		Budget:            row.Budget,
		Revenue:           row.Revenue,
		LibraryAddedAt:    row.LibraryAddedAt,
		MetadataUpdatedAt: pgTimestamptzToTimePtr(row.MetadataUpdatedAt),
		RadarrID:          row.RadarrID,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}

// dbWatchedMovieRowToMovie converts a ListWatchedMoviesRow to a domain Movie
func dbWatchedMovieRowToMovie(row moviedb.ListWatchedMoviesRow) *Movie {
	return &Movie{
		ID:                row.ID,
		TMDbID:            row.TmdbID,
		IMDbID:            row.ImdbID,
		Title:             row.Title,
		OriginalTitle:     row.OriginalTitle,
		Year:              row.Year,
		ReleaseDate:       pgDateToTimePtr(row.ReleaseDate),
		Runtime:           row.Runtime,
		Overview:          row.Overview,
		Tagline:           row.Tagline,
		Status:            row.Status,
		OriginalLanguage:  row.OriginalLanguage,
		TitlesI18n:        unmarshalStringMap(row.TitlesI18n),
		TaglinesI18n:      unmarshalStringMap(row.TaglinesI18n),
		OverviewsI18n:     unmarshalStringMap(row.OverviewsI18n),
		AgeRatings:        unmarshalNestedStringMap(row.AgeRatings),
		PosterPath:        row.PosterPath,
		BackdropPath:      row.BackdropPath,
		TrailerURL:        row.TrailerUrl,
		VoteAverage:       pgNumericToDecimalPtr(row.VoteAverage),
		VoteCount:         row.VoteCount,
		Popularity:        pgNumericToDecimalPtr(row.Popularity),
		Budget:            row.Budget,
		Revenue:           row.Revenue,
		LibraryAddedAt:    row.LibraryAddedAt,
		MetadataUpdatedAt: pgTimestamptzToTimePtr(row.MetadataUpdatedAt),
		RadarrID:          row.RadarrID,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}

// pgNumericToInt32Ptr converts pgtype.Numeric to *int32
func pgNumericToInt32Ptr(n pgtype.Numeric) *int32 {
	if !n.Valid {
		return nil
	}
	// Convert to int64 first
	i64, err := n.Int64Value()
	if err != nil || !i64.Valid {
		return nil
	}
	i32 := util.SafeInt64ToInt32(i64.Int64)
	return &i32
}

// derefInt32 safely dereferences *int32, returning 0 if nil
func derefInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

// derefBool safely dereferences *bool, returning false if nil
func derefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// unmarshalStringMap unmarshals JSONB []byte to map[string]string
func unmarshalStringMap(data []byte) map[string]string {
	if len(data) == 0 {
		return nil
	}
	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

// unmarshalNestedStringMap unmarshals JSONB []byte to map[string]map[string]string
func unmarshalNestedStringMap(data []byte) map[string]map[string]string {
	if len(data) == 0 {
		return nil
	}
	var result map[string]map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

// marshalStringMap marshals map[string]string to JSONB []byte
func marshalStringMap(m map[string]string) []byte {
	if m == nil {
		return []byte("{}")
	}
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("{}")
	}
	return data
}

// marshalNestedStringMap marshals map[string]map[string]string to JSONB []byte
func marshalNestedStringMap(m map[string]map[string]string) []byte {
	if m == nil {
		return []byte("{}")
	}
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("{}")
	}
	return data
}

// marshalExternalRatings marshals []ExternalRating to JSONB []byte
func marshalExternalRatings(ratings []ExternalRating) []byte {
	if ratings == nil {
		return []byte("[]")
	}
	data, err := json.Marshal(ratings)
	if err != nil {
		return []byte("[]")
	}
	return data
}

// unmarshalExternalRatings unmarshals JSONB []byte to []ExternalRating
func unmarshalExternalRatings(data []byte) []ExternalRating {
	if len(data) == 0 {
		return nil
	}
	var result []ExternalRating
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

func (r *postgresRepository) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	dbMovies, err := r.queries.ListMovies(ctx, moviedb.ListMoviesParams{
		Limit:   filters.Limit,
		Offset:  filters.Offset,
		OrderBy: &filters.OrderBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list movies: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) CountMovies(ctx context.Context) (int64, error) {
	return r.queries.CountMovies(ctx)
}

func (r *postgresRepository) SearchMoviesByTitle(ctx context.Context, query string, limit, offset int32) ([]Movie, error) {
	dbMovies, err := r.queries.SearchMoviesByTitle(ctx, moviedb.SearchMoviesByTitleParams{
		Title:  query,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) SearchMoviesByTitleAnyLanguage(ctx context.Context, query string, limit, offset int32) ([]Movie, error) {
	dbMovies, err := r.queries.SearchMoviesByTitleAnyLanguage(ctx, moviedb.SearchMoviesByTitleAnyLanguageParams{
		Column1: &query,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search movies in any language: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) ListMoviesByYear(ctx context.Context, year int32, limit, offset int32) ([]Movie, error) {
	dbMovies, err := r.queries.ListMoviesByYear(ctx, moviedb.ListMoviesByYearParams{
		Year:   &year,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list movies by year: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, error) {
	dbMovies, err := r.queries.ListRecentlyAdded(ctx, moviedb.ListRecentlyAddedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recently added movies: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, error) {
	dbMovies, err := r.queries.ListTopRated(ctx, moviedb.ListTopRatedParams{
		VoteCount: &minVotes,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list top rated movies: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) CountTopRated(ctx context.Context, minVotes int32) (int64, error) {
	return r.queries.CountTopRated(ctx, &minVotes)
}

func (r *postgresRepository) CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error) {
	dbParams := moviedb.CreateMovieParams{
		TmdbID:            params.TMDbID,
		ImdbID:            params.IMDbID,
		Title:             params.Title,
		OriginalTitle:     params.OriginalTitle,
		Year:              params.Year,
		ReleaseDate:       stringToPgDate(params.ReleaseDate),
		Runtime:           params.Runtime,
		Overview:          params.Overview,
		Tagline:           params.Tagline,
		Status:            params.Status,
		OriginalLanguage:  params.OriginalLanguage,
		TitlesI18n:        marshalStringMap(params.TitlesI18n),
		TaglinesI18n:      marshalStringMap(params.TaglinesI18n),
		OverviewsI18n:     marshalStringMap(params.OverviewsI18n),
		AgeRatings:        marshalNestedStringMap(params.AgeRatings),
		ExternalRatings:   marshalExternalRatings(params.ExternalRatings),
		PosterPath:        params.PosterPath,
		BackdropPath:      params.BackdropPath,
		TrailerUrl:        params.TrailerURL,
		VoteAverage:       stringToPgNumeric(params.VoteAverage),
		VoteCount:         params.VoteCount,
		Popularity:        stringToPgNumeric(params.Popularity),
		Budget:            params.Budget,
		Revenue:           params.Revenue,
		RadarrID:          params.RadarrID,
		MetadataUpdatedAt: stringToPgTimestamptz(params.MetadataUpdatedAt),
	}
	movie, err := r.queries.CreateMovie(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

func (r *postgresRepository) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	// Convert maps to []byte for JSONB fields (only if not nil)
	var titlesI18n, taglinesI18n, overviewsI18n []byte
	var ageRatings []byte
	var externalRatings []byte

	if params.TitlesI18n != nil {
		titlesI18n = marshalStringMap(params.TitlesI18n)
	}
	if params.TaglinesI18n != nil {
		taglinesI18n = marshalStringMap(params.TaglinesI18n)
	}
	if params.OverviewsI18n != nil {
		overviewsI18n = marshalStringMap(params.OverviewsI18n)
	}
	if params.AgeRatings != nil {
		ageRatings = marshalNestedStringMap(params.AgeRatings)
	}
	if params.ExternalRatings != nil {
		externalRatings = marshalExternalRatings(params.ExternalRatings)
	}

	dbParams := moviedb.UpdateMovieParams{
		ID:                params.ID,
		TmdbID:            params.TMDbID,
		ImdbID:            params.IMDbID,
		Title:             params.Title,
		OriginalTitle:     params.OriginalTitle,
		Year:              params.Year,
		ReleaseDate:       stringToPgDate(params.ReleaseDate),
		Runtime:           params.Runtime,
		Overview:          params.Overview,
		Tagline:           params.Tagline,
		Status:            params.Status,
		OriginalLanguage:  params.OriginalLanguage,
		TitlesI18n:        titlesI18n,
		TaglinesI18n:      taglinesI18n,
		OverviewsI18n:     overviewsI18n,
		AgeRatings:        ageRatings,
		ExternalRatings:   externalRatings,
		PosterPath:        params.PosterPath,
		BackdropPath:      params.BackdropPath,
		TrailerUrl:        params.TrailerURL,
		VoteAverage:       stringToPgNumeric(params.VoteAverage),
		VoteCount:         params.VoteCount,
		Popularity:        stringToPgNumeric(params.Popularity),
		Budget:            params.Budget,
		Revenue:           params.Revenue,
		RadarrID:          params.RadarrID,
		MetadataUpdatedAt: stringToPgTimestamptz(params.MetadataUpdatedAt),
	}
	movie, err := r.queries.UpdateMovie(ctx, dbParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

func (r *postgresRepository) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteMovie(ctx, id)
}

func (r *postgresRepository) CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error) {
	file, err := r.queries.CreateMovieFile(ctx, moviedb.CreateMovieFileParams{
		MovieID:           params.MovieID,
		FilePath:          params.FilePath,
		FileSize:          params.FileSize,
		FileName:          params.FileName,
		Resolution:        params.Resolution,
		QualityProfile:    params.QualityProfile,
		VideoCodec:        params.VideoCodec,
		AudioCodec:        params.AudioCodec,
		Container:         params.Container,
		BitrateKbps:       params.BitrateKbps,
		AudioLanguages:    params.AudioLanguages,
		SubtitleLanguages: params.SubtitleLanguages,
		RadarrFileID:      params.RadarrFileID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create movie file: %w", err)
	}
	return dbMovieFileToMovieFile(file), nil
}

func (r *postgresRepository) GetMovieFile(ctx context.Context, id uuid.UUID) (*MovieFile, error) {
	file, err := r.queries.GetMovieFile(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieFileNotFound
		}
		return nil, fmt.Errorf("failed to get movie file: %w", err)
	}
	return dbMovieFileToMovieFile(file), nil
}

func (r *postgresRepository) GetMovieFileByPath(ctx context.Context, path string) (*MovieFile, error) {
	file, err := r.queries.GetMovieFileByPath(ctx, path)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieFileNotFound
		}
		return nil, fmt.Errorf("failed to get movie file by path: %w", err)
	}
	return dbMovieFileToMovieFile(file), nil
}

func (r *postgresRepository) GetMovieFileByRadarrID(ctx context.Context, radarrFileID int32) (*MovieFile, error) {
	file, err := r.queries.GetMovieFileByRadarrID(ctx, &radarrFileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieFileNotFound
		}
		return nil, fmt.Errorf("failed to get movie file by Radarr ID: %w", err)
	}
	return dbMovieFileToMovieFile(file), nil
}

func (r *postgresRepository) ListMovieFilesByMovieID(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	dbFiles, err := r.queries.ListMovieFilesByMovieID(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to list movie files: %w", err)
	}
	files := make([]MovieFile, len(dbFiles))
	for i, f := range dbFiles {
		files[i] = *dbMovieFileToMovieFile(f)
	}
	return files, nil
}

func (r *postgresRepository) UpdateMovieFile(ctx context.Context, params UpdateMovieFileParams) (*MovieFile, error) {
	file, err := r.queries.UpdateMovieFile(ctx, moviedb.UpdateMovieFileParams{
		ID:                params.ID,
		FilePath:          params.FilePath,
		FileSize:          params.FileSize,
		Resolution:        params.Resolution,
		QualityProfile:    params.QualityProfile,
		VideoCodec:        params.VideoCodec,
		AudioCodec:        params.AudioCodec,
		Container:         params.Container,
		BitrateKbps:       params.BitrateKbps,
		AudioLanguages:    params.AudioLanguages,
		SubtitleLanguages: params.SubtitleLanguages,
		RadarrFileID:      params.RadarrFileID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMovieFileNotFound
		}
		return nil, fmt.Errorf("failed to update movie file: %w", err)
	}
	return dbMovieFileToMovieFile(file), nil
}

func (r *postgresRepository) DeleteMovieFile(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteMovieFile(ctx, id)
}

func (r *postgresRepository) CreateMovieCredit(ctx context.Context, params CreateMovieCreditParams) (*MovieCredit, error) {
	credit, err := r.queries.CreateMovieCredit(ctx, moviedb.CreateMovieCreditParams{
		MovieID:      params.MovieID,
		TmdbPersonID: params.TMDbPersonID,
		Name:         params.Name,
		CreditType:   params.CreditType,
		Character:    params.Character,
		Job:          params.Job,
		Department:   params.Department,
		CastOrder:    params.CastOrder,
		ProfilePath:  params.ProfilePath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create movie credit: %w", err)
	}
	return dbCreditToCredit(credit), nil
}

func (r *postgresRepository) ListMovieCast(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, error) {
	dbCredits, err := r.queries.ListMovieCast(ctx, moviedb.ListMovieCastParams{
		MovieID: movieID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list movie cast: %w", err)
	}
	credits := make([]MovieCredit, len(dbCredits))
	for i, c := range dbCredits {
		credits[i] = *dbCreditToCredit(c)
	}
	return credits, nil
}

func (r *postgresRepository) CountMovieCast(ctx context.Context, movieID uuid.UUID) (int64, error) {
	return r.queries.CountMovieCast(ctx, movieID)
}

func (r *postgresRepository) ListMovieCrew(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, error) {
	dbCredits, err := r.queries.ListMovieCrew(ctx, moviedb.ListMovieCrewParams{
		MovieID: movieID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list movie crew: %w", err)
	}
	credits := make([]MovieCredit, len(dbCredits))
	for i, c := range dbCredits {
		credits[i] = *dbCreditToCredit(c)
	}
	return credits, nil
}

func (r *postgresRepository) CountMovieCrew(ctx context.Context, movieID uuid.UUID) (int64, error) {
	return r.queries.CountMovieCrew(ctx, movieID)
}

func (r *postgresRepository) DeleteMovieCredits(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.DeleteMovieCredits(ctx, movieID)
}

func (r *postgresRepository) CreateMovieCollection(ctx context.Context, params CreateMovieCollectionParams) (*MovieCollection, error) {
	coll, err := r.queries.CreateMovieCollection(ctx, moviedb.CreateMovieCollectionParams{
		TmdbCollectionID: params.TMDbCollectionID,
		Name:             params.Name,
		Overview:         params.Overview,
		PosterPath:       params.PosterPath,
		BackdropPath:     params.BackdropPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create movie collection: %w", err)
	}
	return dbCollectionToCollection(coll), nil
}

func (r *postgresRepository) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	coll, err := r.queries.GetMovieCollection(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCollectionNotFound
		}
		return nil, fmt.Errorf("failed to get movie collection: %w", err)
	}
	return dbCollectionToCollection(coll), nil
}

func (r *postgresRepository) GetMovieCollectionByTMDbID(ctx context.Context, tmdbCollectionID int32) (*MovieCollection, error) {
	coll, err := r.queries.GetMovieCollectionByTMDbID(ctx, &tmdbCollectionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCollectionNotFound
		}
		return nil, fmt.Errorf("failed to get movie collection by TMDb ID: %w", err)
	}
	return dbCollectionToCollection(coll), nil
}

func (r *postgresRepository) UpdateMovieCollection(ctx context.Context, params UpdateMovieCollectionParams) (*MovieCollection, error) {
	coll, err := r.queries.UpdateMovieCollection(ctx, moviedb.UpdateMovieCollectionParams{
		ID:               params.ID,
		TmdbCollectionID: params.TMDbCollectionID,
		Name:             params.Name,
		Overview:         params.Overview,
		PosterPath:       params.PosterPath,
		BackdropPath:     params.BackdropPath,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCollectionNotFound
		}
		return nil, fmt.Errorf("failed to update movie collection: %w", err)
	}
	return dbCollectionToCollection(coll), nil
}

func (r *postgresRepository) AddMovieToCollection(ctx context.Context, collectionID, movieID uuid.UUID, collectionOrder *int32) error {
	return r.queries.AddMovieToCollection(ctx, moviedb.AddMovieToCollectionParams{
		CollectionID:    collectionID,
		MovieID:         movieID,
		CollectionOrder: collectionOrder,
	})
}

func (r *postgresRepository) RemoveMovieFromCollection(ctx context.Context, collectionID, movieID uuid.UUID) error {
	return r.queries.RemoveMovieFromCollection(ctx, moviedb.RemoveMovieFromCollectionParams{
		CollectionID: collectionID,
		MovieID:      movieID,
	})
}

func (r *postgresRepository) ListMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error) {
	dbMovies, err := r.queries.ListMoviesByCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list movies by collection: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error) {
	coll, err := r.queries.GetCollectionForMovie(ctx, movieID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotInCollection
		}
		return nil, fmt.Errorf("failed to get collection for movie: %w", err)
	}
	return dbCollectionToCollection(coll), nil
}

// dbCollectionToCollection converts a database collection to a domain collection
func dbCollectionToCollection(dbColl moviedb.MovieCollection) *MovieCollection {
	return &MovieCollection{
		ID:               dbColl.ID,
		TMDbCollectionID: dbColl.TmdbCollectionID,
		Name:             dbColl.Name,
		Overview:         dbColl.Overview,
		PosterPath:       dbColl.PosterPath,
		BackdropPath:     dbColl.BackdropPath,
		CreatedAt:        dbColl.CreatedAt,
		UpdatedAt:        dbColl.UpdatedAt,
	}
}

// dbCreditToCredit converts a database movie credit to a domain movie credit
func dbCreditToCredit(dbCredit moviedb.MovieCredit) *MovieCredit {
	return &MovieCredit{
		ID:           dbCredit.ID,
		MovieID:      dbCredit.MovieID,
		TMDbPersonID: dbCredit.TmdbPersonID,
		Name:         dbCredit.Name,
		CreditType:   dbCredit.CreditType,
		Character:    dbCredit.Character,
		Job:          dbCredit.Job,
		Department:   dbCredit.Department,
		CastOrder:    dbCredit.CastOrder,
		ProfilePath:  dbCredit.ProfilePath,
		CreatedAt:    dbCredit.CreatedAt,
		UpdatedAt:    dbCredit.UpdatedAt,
	}
}

// dbGenreToGenre converts a database movie genre to a domain movie genre
func dbGenreToGenre(dbGenre moviedb.MovieGenre) *MovieGenre {
	return &MovieGenre{
		ID:        dbGenre.ID,
		MovieID:   dbGenre.MovieID,
		Slug:      dbGenre.Slug,
		Name:      dbGenre.Name,
		CreatedAt: dbGenre.CreatedAt,
	}
}

// dbMovieFileToMovieFile converts a database movie file to a domain movie file
func dbMovieFileToMovieFile(dbFile moviedb.MovieFile) *MovieFile {
	return &MovieFile{
		ID:                dbFile.ID,
		MovieID:           dbFile.MovieID,
		FilePath:          dbFile.FilePath,
		FileSize:          dbFile.FileSize,
		Resolution:        dbFile.Resolution,
		QualityProfile:    dbFile.QualityProfile,
		VideoCodec:        dbFile.VideoCodec,
		AudioCodec:        dbFile.AudioCodec,
		Container:         dbFile.Container,
		BitrateKbps:       dbFile.BitrateKbps,
		AudioLanguages:    dbFile.AudioLanguages,
		SubtitleLanguages: dbFile.SubtitleLanguages,
		RadarrFileID:      dbFile.RadarrFileID,
		CreatedAt:         dbFile.CreatedAt,
		UpdatedAt:         dbFile.UpdatedAt,
	}
}

func (r *postgresRepository) AddMovieGenre(ctx context.Context, movieID uuid.UUID, slug, name string) error {
	return r.queries.AddMovieGenre(ctx, moviedb.AddMovieGenreParams{
		MovieID: movieID,
		Slug:    slug,
		Name:    name,
	})
}

func (r *postgresRepository) ListMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	dbGenres, err := r.queries.ListMovieGenres(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to list movie genres: %w", err)
	}
	genres := make([]MovieGenre, len(dbGenres))
	for i, g := range dbGenres {
		genres[i] = *dbGenreToGenre(g)
	}
	return genres, nil
}

func (r *postgresRepository) ListDistinctMovieGenres(ctx context.Context) ([]content.GenreSummary, error) {
	rows, err := r.queries.ListDistinctMovieGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list distinct movie genres: %w", err)
	}
	genres := make([]content.GenreSummary, len(rows))
	for i, r := range rows {
		genres[i] = content.GenreSummary{
			Slug:      r.Slug,
			Name:      r.Name,
			ItemCount: r.ItemCount,
		}
	}
	return genres, nil
}

func (r *postgresRepository) DeleteMovieGenres(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.DeleteMovieGenres(ctx, movieID)
}

func (r *postgresRepository) ListMoviesByGenre(ctx context.Context, slug string, limit, offset int32) ([]Movie, error) {
	dbMovies, err := r.queries.ListMoviesByGenre(ctx, moviedb.ListMoviesByGenreParams{
		Slug:   slug,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list movies by genre: %w", err)
	}
	movies := make([]Movie, len(dbMovies))
	for i, m := range dbMovies {
		movies[i] = *dbMovieToMovie(m)
	}
	return movies, nil
}

func (r *postgresRepository) CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*MovieWatched, error) {
	watched, err := r.queries.CreateOrUpdateWatchProgress(ctx, moviedb.CreateOrUpdateWatchProgressParams{
		UserID:          params.UserID,
		MovieID:         params.MovieID,
		ProgressSeconds: params.ProgressSeconds,
		DurationSeconds: &params.DurationSeconds,
		IsCompleted:     &params.IsCompleted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create/update watch progress: %w", err)
	}
	return dbWatchedToWatched(watched), nil
}

func (r *postgresRepository) GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error) {
	watched, err := r.queries.GetWatchProgress(ctx, moviedb.GetWatchProgressParams{
		UserID:  userID,
		MovieID: movieID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProgressNotFound
		}
		return nil, fmt.Errorf("failed to get watch progress: %w", err)
	}
	return dbWatchedToWatched(watched), nil
}

func (r *postgresRepository) DeleteWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error {
	return r.queries.DeleteWatchProgress(ctx, moviedb.DeleteWatchProgressParams{
		UserID:  userID,
		MovieID: movieID,
	})
}

func (r *postgresRepository) ListContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	rows, err := r.queries.ListContinueWatching(ctx, moviedb.ListContinueWatchingParams{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list continue watching: %w", err)
	}
	items := make([]ContinueWatchingItem, len(rows))
	for i, row := range rows {
		items[i] = ContinueWatchingItem{
			Movie:           *dbContinueWatchingRowToMovie(row),
			ProgressSeconds: row.ProgressSeconds,
			DurationSeconds: derefInt32(row.DurationSeconds),
			ProgressPercent: pgNumericToInt32Ptr(row.ProgressPercent),
			LastWatchedAt:   row.LastWatchedAt,
		}
	}
	return items, nil
}

func (r *postgresRepository) ListWatchedMovies(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error) {
	rows, err := r.queries.ListWatchedMovies(ctx, moviedb.ListWatchedMoviesParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list watched movies: %w", err)
	}
	items := make([]WatchedMovieItem, len(rows))
	for i, row := range rows {
		items[i] = WatchedMovieItem{
			Movie:         *dbWatchedMovieRowToMovie(row),
			WatchCount:    derefInt32(row.WatchCount),
			LastWatchedAt: row.LastWatchedAt,
		}
	}
	return items, nil
}

func (r *postgresRepository) GetUserMovieStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error) {
	stats, err := r.queries.GetUserMovieStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user movie stats: %w", err)
	}
	return &UserMovieStats{
		WatchedCount:    stats.WatchedCount,
		InProgressCount: stats.InProgressCount,
		TotalWatches:    &stats.TotalWatches,
	}, nil
}
