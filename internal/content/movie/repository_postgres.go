package movie

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	moviedb "github.com/lusoris/revenge/internal/content/movie/db"
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// GetMovieByTMDbID retrieves a movie by TMDb ID
func (r *postgresRepository) GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error) {
	movie, err := r.queries.GetMovieByTMDbID(ctx, &tmdbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get movie by TMDb ID: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// GetMovieByIMDbID retrieves a movie by IMDb ID
func (r *postgresRepository) GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error) {
	movie, err := r.queries.GetMovieByIMDbID(ctx, &imdbID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get movie by IMDb ID: %w", err)
	}
	return dbMovieToMovie(movie), nil
}

// GetMovieByRadarrID retrieves a movie by Radarr ID
func (r *postgresRepository) GetMovieByRadarrID(ctx context.Context, radarrID int32) (*Movie, error) {
	movie, err := r.queries.GetMovieByRadarrID(ctx, &radarrID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found: %w", err)
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
	d, _ := decimal.NewFromString(s)
	return &d
}

// Placeholder implementations for remaining methods
// TODO: Implement all repository methods

func (r *postgresRepository) ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) CountMovies(ctx context.Context) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (r *postgresRepository) SearchMoviesByTitle(ctx context.Context, query string, limit, offset int32) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMoviesByYear(ctx context.Context, year int32, limit, offset int32) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetMovieFile(ctx context.Context, id uuid.UUID) (*MovieFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetMovieFileByPath(ctx context.Context, path string) (*MovieFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMovieFilesByMovieID(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) UpdateMovieFile(ctx context.Context, params UpdateMovieFileParams) (*MovieFile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) DeleteMovieFile(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) CreateMovieCredit(ctx context.Context, params CreateMovieCreditParams) (*MovieCredit, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMovieCast(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMovieCrew(ctx context.Context, movieID uuid.UUID) ([]MovieCredit, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) DeleteMovieCredits(ctx context.Context, movieID uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) CreateMovieCollection(ctx context.Context, params CreateMovieCollectionParams) (*MovieCollection, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetMovieCollectionByTMDbID(ctx context.Context, tmdbCollectionID int32) (*MovieCollection, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) AddMovieToCollection(ctx context.Context, collectionID, movieID uuid.UUID, collectionOrder *int32) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) RemoveMovieFromCollection(ctx context.Context, collectionID, movieID uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) AddMovieGenre(ctx context.Context, movieID uuid.UUID, tmdbGenreID int32, name string) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) DeleteMovieGenres(ctx context.Context, movieID uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListMoviesByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Movie, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*MovieWatched, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) DeleteWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) ListWatchedMovies(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *postgresRepository) GetUserMovieStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error) {
	return nil, fmt.Errorf("not implemented")
}
