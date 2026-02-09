package movie

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content"
)

// Repository defines database operations for movies
type Repository interface {
	// Movie CRUD
	GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error)
	GetMovieByTMDbID(ctx context.Context, tmdbID int32) (*Movie, error)
	GetMovieByIMDbID(ctx context.Context, imdbID string) (*Movie, error)
	GetMovieByRadarrID(ctx context.Context, radarrID int32) (*Movie, error)
	ListMovies(ctx context.Context, filters ListFilters) ([]Movie, error)
	CountMovies(ctx context.Context) (int64, error)
	SearchMoviesByTitle(ctx context.Context, query string, limit, offset int32) ([]Movie, error)
	SearchMoviesByTitleAnyLanguage(ctx context.Context, query string, limit, offset int32) ([]Movie, error)
	ListMoviesByYear(ctx context.Context, year int32, limit, offset int32) ([]Movie, error)
	ListRecentlyAdded(ctx context.Context, limit, offset int32) ([]Movie, error)
	ListTopRated(ctx context.Context, minVotes int32, limit, offset int32) ([]Movie, error)
	CountTopRated(ctx context.Context, minVotes int32) (int64, error)
	CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error)
	UpdateMovie(ctx context.Context, params UpdateMovieParams) (*Movie, error)
	DeleteMovie(ctx context.Context, id uuid.UUID) error

	// Movie Files
	CreateMovieFile(ctx context.Context, params CreateMovieFileParams) (*MovieFile, error)
	GetMovieFile(ctx context.Context, id uuid.UUID) (*MovieFile, error)
	GetMovieFileByPath(ctx context.Context, path string) (*MovieFile, error)
	GetMovieFileByRadarrID(ctx context.Context, radarrFileID int32) (*MovieFile, error)
	ListMovieFilesByMovieID(ctx context.Context, movieID uuid.UUID) ([]MovieFile, error)
	UpdateMovieFile(ctx context.Context, params UpdateMovieFileParams) (*MovieFile, error)
	DeleteMovieFile(ctx context.Context, id uuid.UUID) error

	// Credits
	CreateMovieCredit(ctx context.Context, params CreateMovieCreditParams) (*MovieCredit, error)
	ListMovieCast(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, error)
	ListMovieCrew(ctx context.Context, movieID uuid.UUID, limit, offset int32) ([]MovieCredit, error)
	CountMovieCast(ctx context.Context, movieID uuid.UUID) (int64, error)
	CountMovieCrew(ctx context.Context, movieID uuid.UUID) (int64, error)
	DeleteMovieCredits(ctx context.Context, movieID uuid.UUID) error

	// Collections
	CreateMovieCollection(ctx context.Context, params CreateMovieCollectionParams) (*MovieCollection, error)
	GetMovieCollection(ctx context.Context, id uuid.UUID) (*MovieCollection, error)
	GetMovieCollectionByTMDbID(ctx context.Context, tmdbCollectionID int32) (*MovieCollection, error)
	UpdateMovieCollection(ctx context.Context, params UpdateMovieCollectionParams) (*MovieCollection, error)
	AddMovieToCollection(ctx context.Context, collectionID, movieID uuid.UUID, collectionOrder *int32) error
	RemoveMovieFromCollection(ctx context.Context, collectionID, movieID uuid.UUID) error
	ListMoviesByCollection(ctx context.Context, collectionID uuid.UUID) ([]Movie, error)
	GetCollectionForMovie(ctx context.Context, movieID uuid.UUID) (*MovieCollection, error)

	// Genres
	AddMovieGenre(ctx context.Context, movieID uuid.UUID, tmdbGenreID int32, name string) error
	ListMovieGenres(ctx context.Context, movieID uuid.UUID) ([]MovieGenre, error)
	ListDistinctMovieGenres(ctx context.Context) ([]content.GenreSummary, error)
	DeleteMovieGenres(ctx context.Context, movieID uuid.UUID) error
	ListMoviesByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Movie, error)

	// Watch Progress
	CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*MovieWatched, error)
	GetWatchProgress(ctx context.Context, userID, movieID uuid.UUID) (*MovieWatched, error)
	DeleteWatchProgress(ctx context.Context, userID, movieID uuid.UUID) error
	ListContinueWatching(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error)
	ListWatchedMovies(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedMovieItem, error)
	GetUserMovieStats(ctx context.Context, userID uuid.UUID) (*UserMovieStats, error)
}

// CreateMovieParams contains parameters for creating a movie
type CreateMovieParams struct {
	TMDbID            *int32
	IMDbID            *string
	Title             string
	OriginalTitle     *string
	Year              *int32
	ReleaseDate       *string
	Runtime           *int32
	Overview          *string
	Tagline           *string
	Status            *string
	OriginalLanguage  *string
	TitlesI18n        map[string]string
	TaglinesI18n      map[string]string
	OverviewsI18n     map[string]string
	AgeRatings        map[string]map[string]string
	ExternalRatings   []ExternalRating
	PosterPath        *string
	BackdropPath      *string
	TrailerURL        *string
	VoteAverage       *string
	VoteCount         *int32
	Popularity        *string
	Budget            *int64
	Revenue           *int64
	RadarrID          *int32
	MetadataUpdatedAt *string
}

// UpdateMovieParams contains parameters for updating a movie
type UpdateMovieParams struct {
	ID                uuid.UUID
	TMDbID            *int32
	IMDbID            *string
	Title             *string
	OriginalTitle     *string
	Year              *int32
	ReleaseDate       *string
	Runtime           *int32
	Overview          *string
	Tagline           *string
	Status            *string
	OriginalLanguage  *string
	TitlesI18n        map[string]string
	TaglinesI18n      map[string]string
	OverviewsI18n     map[string]string
	AgeRatings        map[string]map[string]string
	ExternalRatings   []ExternalRating
	PosterPath        *string
	BackdropPath      *string
	TrailerURL        *string
	VoteAverage       *string
	VoteCount         *int32
	Popularity        *string
	Budget            *int64
	Revenue           *int64
	RadarrID          *int32
	MetadataUpdatedAt *string
}

// CreateMovieFileParams contains parameters for creating a movie file
type CreateMovieFileParams struct {
	MovieID           uuid.UUID
	FilePath          string
	FileSize          int64
	Resolution        *string
	QualityProfile    *string
	VideoCodec        *string
	AudioCodec        *string
	Container         *string
	BitrateKbps       *int32
	AudioLanguages    []string
	SubtitleLanguages []string
	RadarrFileID      *int32
}

// UpdateMovieFileParams contains parameters for updating a movie file
type UpdateMovieFileParams struct {
	ID                uuid.UUID
	FilePath          *string
	FileSize          *int64
	Resolution        *string
	QualityProfile    *string
	VideoCodec        *string
	AudioCodec        *string
	Container         *string
	BitrateKbps       *int32
	AudioLanguages    []string
	SubtitleLanguages []string
	RadarrFileID      *int32
}

// CreateMovieCreditParams contains parameters for creating a movie credit
type CreateMovieCreditParams struct {
	MovieID      uuid.UUID
	TMDbPersonID int32
	Name         string
	CreditType   string
	Character    *string
	Job          *string
	Department   *string
	CastOrder    *int32
	ProfilePath  *string
}

// CreateMovieCollectionParams contains parameters for creating a movie collection
type CreateMovieCollectionParams struct {
	TMDbCollectionID *int32
	Name             string
	Overview         *string
	PosterPath       *string
	BackdropPath     *string
}

// UpdateMovieCollectionParams contains parameters for updating a movie collection
type UpdateMovieCollectionParams struct {
	ID               uuid.UUID
	TMDbCollectionID *int32
	Name             *string
	Overview         *string
	PosterPath       *string
	BackdropPath     *string
}

// CreateWatchProgressParams contains parameters for creating/updating watch progress
type CreateWatchProgressParams struct {
	UserID          uuid.UUID
	MovieID         uuid.UUID
	ProgressSeconds int32
	DurationSeconds int32
	IsCompleted     bool
}

// ContinueWatchingItem represents a movie with watch progress
type ContinueWatchingItem struct {
	Movie
	ProgressSeconds int32
	DurationSeconds int32
	ProgressPercent *int32
	LastWatchedAt   time.Time
}

// WatchedMovieItem represents a watched movie with statistics
type WatchedMovieItem struct {
	Movie
	WatchCount    int32
	LastWatchedAt time.Time
}

// UserMovieStats represents statistics for a user's movie watching
type UserMovieStats struct {
	WatchedCount    int64
	InProgressCount int64
	TotalWatches    *int64
}
