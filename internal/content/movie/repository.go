package movie

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Repository errors.
var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrStudioNotFound     = errors.New("studio not found")
)

// ListParams contains pagination and sorting parameters.
type ListParams struct {
	Limit     int
	Offset    int
	SortBy    string // title, date_added, release_date, rating
	SortOrder string // asc, desc
}

// DefaultListParams returns default list parameters.
func DefaultListParams() ListParams {
	return ListParams{
		Limit:     20,
		Offset:    0,
		SortBy:    "title",
		SortOrder: "asc",
	}
}

// Repository defines the interface for movie data access.
type Repository interface {
	// Core CRUD
	GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
	GetByPath(ctx context.Context, path string) (*Movie, error)
	GetByTmdbID(ctx context.Context, tmdbID int) (*Movie, error)
	GetByImdbID(ctx context.Context, imdbID string) (*Movie, error)
	List(ctx context.Context, params ListParams) ([]*Movie, error)
	ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, error)
	ListByCollection(ctx context.Context, collectionID uuid.UUID) ([]*Movie, error)
	ListRecentlyAdded(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Movie, error)
	ListRecentlyPlayed(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Movie, error)
	Search(ctx context.Context, query string, params ListParams) ([]*Movie, error)
	Count(ctx context.Context) (int64, error)
	CountByLibrary(ctx context.Context, libraryID uuid.UUID) (int64, error)
	Create(ctx context.Context, movie *Movie) error
	Update(ctx context.Context, movie *Movie) error
	UpdatePlaybackStats(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByLibrary(ctx context.Context, libraryID uuid.UUID) error
	ExistsByPath(ctx context.Context, path string) (bool, error)
	ExistsByTmdbID(ctx context.Context, tmdbID int) (bool, error)
	ListPaths(ctx context.Context, libraryID uuid.UUID) (map[uuid.UUID]string, error)

	// Collections
	GetCollectionByID(ctx context.Context, id uuid.UUID) (*Collection, error)
	GetCollectionByTmdbID(ctx context.Context, tmdbID int) (*Collection, error)
	ListCollections(ctx context.Context, params ListParams) ([]*Collection, error)
	CreateCollection(ctx context.Context, collection *Collection) error
	UpdateCollection(ctx context.Context, collection *Collection) error
	DeleteCollection(ctx context.Context, id uuid.UUID) error

	// Studios
	GetStudioByID(ctx context.Context, id uuid.UUID) (*Studio, error)
	GetStudioByTmdbID(ctx context.Context, tmdbID int) (*Studio, error)
	ListStudios(ctx context.Context, params ListParams) ([]*Studio, error)
	CreateStudio(ctx context.Context, studio *Studio) error
	LinkMovieStudio(ctx context.Context, movieID, studioID uuid.UUID, order int) error
	UnlinkMovieStudios(ctx context.Context, movieID uuid.UUID) error
	GetMovieStudios(ctx context.Context, movieID uuid.UUID) ([]Studio, error)

	// Genres
	GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]Genre, error)
	LinkMovieGenre(ctx context.Context, movieID, genreID uuid.UUID) error
	UnlinkMovieGenres(ctx context.Context, movieID uuid.UUID) error
	ListMoviesByGenre(ctx context.Context, genreID uuid.UUID, params ListParams) ([]*Movie, error)
	CountMoviesByGenre(ctx context.Context, genreID uuid.UUID) (int64, error)

	// Credits
	GetMovieCast(ctx context.Context, movieID uuid.UUID) ([]CastMember, error)
	GetMovieCrew(ctx context.Context, movieID uuid.UUID) ([]CrewMember, error)
	GetMovieDirectors(ctx context.Context, movieID uuid.UUID) ([]CrewMember, error)
	GetMovieWriters(ctx context.Context, movieID uuid.UUID) ([]CrewMember, error)
	CreateMovieCredit(ctx context.Context, movieID, personID uuid.UUID, role, character, department, job string, order int, isGuest bool, tmdbCreditID string) error
	DeleteMovieCredits(ctx context.Context, movieID uuid.UUID) error

	// Images
	GetMovieImages(ctx context.Context, movieID uuid.UUID) ([]Image, error)
	GetMovieImagesByType(ctx context.Context, movieID uuid.UUID, imageType string) ([]Image, error)
	CreateMovieImage(ctx context.Context, movieID uuid.UUID, img *Image) error
	DeleteMovieImages(ctx context.Context, movieID uuid.UUID) error

	// Videos
	GetMovieVideos(ctx context.Context, movieID uuid.UUID) ([]Video, error)
	CreateMovieVideo(ctx context.Context, movieID uuid.UUID, video *Video) error
	DeleteMovieVideos(ctx context.Context, movieID uuid.UUID) error

	// User Data
	GetUserRating(ctx context.Context, userID, movieID uuid.UUID) (*UserRating, error)
	SetUserRating(ctx context.Context, userID, movieID uuid.UUID, rating float64, review string) error
	DeleteUserRating(ctx context.Context, userID, movieID uuid.UUID) error
	IsFavorite(ctx context.Context, userID, movieID uuid.UUID) (bool, error)
	AddFavorite(ctx context.Context, userID, movieID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, movieID uuid.UUID) error
	ListFavorites(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Movie, error)
	CountFavorites(ctx context.Context, userID uuid.UUID) (int64, error)
	GetWatchHistory(ctx context.Context, userID, movieID uuid.UUID) (*WatchHistory, error)
	CreateWatchHistory(ctx context.Context, history *WatchHistory) error
	UpdateWatchHistory(ctx context.Context, id uuid.UUID, positionTicks int64, durationTicks *int64) error
	MarkWatchHistoryCompleted(ctx context.Context, id uuid.UUID) error
	ListResumeableMovies(ctx context.Context, userID uuid.UUID, limit int) ([]WatchHistory, error)
	IsWatched(ctx context.Context, userID, movieID uuid.UUID) (bool, error)
	IsInWatchlist(ctx context.Context, userID, movieID uuid.UUID) (bool, error)
	AddToWatchlist(ctx context.Context, userID, movieID uuid.UUID) error
	RemoveFromWatchlist(ctx context.Context, userID, movieID uuid.UUID) error
	ListWatchlist(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Movie, error)
	CountWatchlist(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteWatchHistory(ctx context.Context, id uuid.UUID) error
	CountCollections(ctx context.Context) (int64, error)
}
