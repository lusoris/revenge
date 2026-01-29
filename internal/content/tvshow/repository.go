package tvshow

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Repository errors.
var (
	ErrSeriesNotFound  = errors.New("series not found")
	ErrSeasonNotFound  = errors.New("season not found")
	ErrEpisodeNotFound = errors.New("episode not found")
	ErrNetworkNotFound = errors.New("network not found")
)

// ListParams contains pagination and sorting parameters.
type ListParams struct {
	Limit     int
	Offset    int
	SortBy    string // title, date_added, first_air_date, rating
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

// Repository defines the interface for TV show data access.
type Repository interface {
	// Series CRUD
	GetSeriesByID(ctx context.Context, id uuid.UUID) (*Series, error)
	GetSeriesByTmdbID(ctx context.Context, tmdbID int) (*Series, error)
	GetSeriesByImdbID(ctx context.Context, imdbID string) (*Series, error)
	GetSeriesByTvdbID(ctx context.Context, tvdbID int) (*Series, error)
	ListSeries(ctx context.Context, params ListParams) ([]*Series, error)
	ListSeriesByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Series, error)
	ListRecentlyAddedSeries(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Series, error)
	ListRecentlyPlayedSeries(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Series, error)
	ListCurrentlyAiringSeries(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Series, error)
	SearchSeries(ctx context.Context, query string, params ListParams) ([]*Series, error)
	CountSeries(ctx context.Context) (int64, error)
	CountSeriesByLibrary(ctx context.Context, libraryID uuid.UUID) (int64, error)
	CreateSeries(ctx context.Context, series *Series) error
	UpdateSeries(ctx context.Context, series *Series) error
	UpdateSeriesPlaybackStats(ctx context.Context, id uuid.UUID) error
	DeleteSeries(ctx context.Context, id uuid.UUID) error
	DeleteSeriesByLibrary(ctx context.Context, libraryID uuid.UUID) error
	SeriesExistsByTmdbID(ctx context.Context, tmdbID int) (bool, error)
	SeriesExistsByTvdbID(ctx context.Context, tvdbID int) (bool, error)

	// Seasons
	GetSeasonByID(ctx context.Context, id uuid.UUID) (*Season, error)
	GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int) (*Season, error)
	ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]*Season, error)
	ListSeasonsWithSpecials(ctx context.Context, seriesID uuid.UUID) ([]*Season, error)
	CountSeasons(ctx context.Context, seriesID uuid.UUID) (int64, error)
	CreateSeason(ctx context.Context, season *Season) error
	UpdateSeason(ctx context.Context, season *Season) error
	DeleteSeason(ctx context.Context, id uuid.UUID) error
	DeleteSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) error
	GetOrCreateSeason(ctx context.Context, seriesID uuid.UUID, seasonNumber int, name string) (*Season, error)

	// Episodes
	GetEpisodeByID(ctx context.Context, id uuid.UUID) (*Episode, error)
	GetEpisodeByPath(ctx context.Context, path string) (*Episode, error)
	GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error)
	GetEpisodeByAbsoluteNumber(ctx context.Context, seriesID uuid.UUID, absoluteNumber int) (*Episode, error)
	ListEpisodes(ctx context.Context, seriesID uuid.UUID) ([]*Episode, error)
	ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]*Episode, error)
	ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int) ([]*Episode, error)
	ListRecentlyAddedEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error)
	ListRecentlyAiredEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error)
	ListUpcomingEpisodes(ctx context.Context, libraryIDs []uuid.UUID, limit int) ([]*Episode, error)
	CountEpisodes(ctx context.Context, seriesID uuid.UUID) (int64, error)
	CountEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) (int64, error)
	CreateEpisode(ctx context.Context, episode *Episode) error
	UpdateEpisode(ctx context.Context, episode *Episode) error
	UpdateEpisodePlaybackStats(ctx context.Context, id uuid.UUID) error
	DeleteEpisode(ctx context.Context, id uuid.UUID) error
	DeleteEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) error
	DeleteEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) error
	EpisodeExistsByPath(ctx context.Context, path string) (bool, error)
	EpisodeExistsByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (bool, error)
	ListEpisodePaths(ctx context.Context, libraryID uuid.UUID) (map[uuid.UUID]string, error)
	GetNextEpisode(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error)
	GetPreviousEpisode(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int) (*Episode, error)

	// Networks
	GetNetworkByID(ctx context.Context, id uuid.UUID) (*Network, error)
	GetNetworkByTmdbID(ctx context.Context, tmdbID int) (*Network, error)
	ListNetworks(ctx context.Context, params ListParams) ([]*Network, error)
	CreateNetwork(ctx context.Context, network *Network) error
	GetOrCreateNetwork(ctx context.Context, name string, logoPath, originCountry string, tmdbID int) (*Network, error)
	GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]Network, error)
	LinkSeriesNetwork(ctx context.Context, seriesID, networkID uuid.UUID, order int) error
	UnlinkSeriesNetworks(ctx context.Context, seriesID uuid.UUID) error
	ListSeriesByNetwork(ctx context.Context, networkID uuid.UUID, params ListParams) ([]*Series, error)
	CountSeriesByNetwork(ctx context.Context, networkID uuid.UUID) (int64, error)

	// Genres
	GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]Genre, error)
	LinkSeriesGenre(ctx context.Context, seriesID, genreID uuid.UUID) error
	UnlinkSeriesGenres(ctx context.Context, seriesID uuid.UUID) error
	ListSeriesByGenre(ctx context.Context, genreID uuid.UUID, params ListParams) ([]*Series, error)
	CountSeriesByGenre(ctx context.Context, genreID uuid.UUID) (int64, error)
	GetOrCreateGenre(ctx context.Context, name string, tmdbID int) (*Genre, error)

	// Series Credits
	GetSeriesCast(ctx context.Context, seriesID uuid.UUID) ([]CastMember, error)
	GetSeriesCrew(ctx context.Context, seriesID uuid.UUID) ([]CrewMember, error)
	GetSeriesCreators(ctx context.Context, seriesID uuid.UUID) ([]CrewMember, error)
	CreateSeriesCredit(ctx context.Context, seriesID, personID uuid.UUID, role, character, department, job string, order int, tmdbCreditID string) error
	DeleteSeriesCredits(ctx context.Context, seriesID uuid.UUID) error

	// Episode Credits
	GetEpisodeCast(ctx context.Context, episodeID uuid.UUID) ([]CastMember, error)
	GetEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]CastMember, error)
	GetEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]CrewMember, error)
	GetEpisodeDirectors(ctx context.Context, episodeID uuid.UUID) ([]CrewMember, error)
	GetEpisodeWriters(ctx context.Context, episodeID uuid.UUID) ([]CrewMember, error)
	CreateEpisodeCredit(ctx context.Context, episodeID, personID uuid.UUID, role, character, department, job string, order int, isGuest bool, tmdbCreditID string) error
	DeleteEpisodeCredits(ctx context.Context, episodeID uuid.UUID) error

	// Images
	GetSeriesImages(ctx context.Context, seriesID uuid.UUID) ([]Image, error)
	GetSeriesImagesByType(ctx context.Context, seriesID uuid.UUID, imageType string) ([]Image, error)
	CreateSeriesImage(ctx context.Context, seriesID uuid.UUID, img *Image) error
	DeleteSeriesImages(ctx context.Context, seriesID uuid.UUID) error
	GetSeasonImages(ctx context.Context, seasonID uuid.UUID) ([]Image, error)
	CreateSeasonImage(ctx context.Context, seasonID uuid.UUID, img *Image) error
	DeleteSeasonImages(ctx context.Context, seasonID uuid.UUID) error
	GetEpisodeImages(ctx context.Context, episodeID uuid.UUID) ([]Image, error)
	CreateEpisodeImage(ctx context.Context, episodeID uuid.UUID, img *Image) error
	DeleteEpisodeImages(ctx context.Context, episodeID uuid.UUID) error

	// Videos
	GetSeriesVideos(ctx context.Context, seriesID uuid.UUID) ([]Video, error)
	CreateSeriesVideo(ctx context.Context, seriesID uuid.UUID, video *Video) error
	DeleteSeriesVideos(ctx context.Context, seriesID uuid.UUID) error

	// User Data - Series
	GetSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesUserRating, error)
	SetSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID, rating float64, review string) error
	DeleteSeriesUserRating(ctx context.Context, userID, seriesID uuid.UUID) error
	IsSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) (bool, error)
	AddSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) error
	RemoveSeriesFavorite(ctx context.Context, userID, seriesID uuid.UUID) error
	ListFavoriteSeries(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Series, error)
	CountFavoriteSeries(ctx context.Context, userID uuid.UUID) (int64, error)
	IsSeriesInWatchlist(ctx context.Context, userID, seriesID uuid.UUID) (bool, error)
	AddSeriesToWatchlist(ctx context.Context, userID, seriesID uuid.UUID) error
	RemoveSeriesFromWatchlist(ctx context.Context, userID, seriesID uuid.UUID) error
	ListSeriesWatchlist(ctx context.Context, userID uuid.UUID, params ListParams) ([]*Series, error)
	CountSeriesWatchlist(ctx context.Context, userID uuid.UUID) (int64, error)

	// User Data - Episodes
	GetEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeUserRating, error)
	SetEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID, rating float64) error
	DeleteEpisodeUserRating(ctx context.Context, userID, episodeID uuid.UUID) error
	GetEpisodeWatchHistory(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatchHistory, error)
	CreateEpisodeWatchHistory(ctx context.Context, history *EpisodeWatchHistory) error
	UpdateEpisodeWatchHistory(ctx context.Context, id uuid.UUID, positionTicks int64, durationTicks *int64) error
	MarkEpisodeWatchHistoryCompleted(ctx context.Context, id uuid.UUID) error
	DeleteEpisodeWatchHistory(ctx context.Context, id uuid.UUID) error
	ListResumeableEpisodes(ctx context.Context, userID uuid.UUID, limit int) ([]EpisodeWatchHistory, error)
	IsEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID) (bool, error)
	CountWatchedEpisodes(ctx context.Context, userID uuid.UUID) (int64, error)
	CountWatchedEpisodesBySeries(ctx context.Context, userID, seriesID uuid.UUID) (int64, error)

	// Series Watch Progress (Continue Watching)
	GetSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchProgress, error)
	ListContinueWatchingSeries(ctx context.Context, userID uuid.UUID, limit int) ([]*SeriesWatchProgress, error)
	ListCompletedSeries(ctx context.Context, userID uuid.UUID, params ListParams) ([]*SeriesWatchProgress, error)
	DeleteSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) error

	// External Ratings
	GetSeriesExternalRatings(ctx context.Context, seriesID uuid.UUID) (map[string]float64, error)
	UpsertSeriesExternalRating(ctx context.Context, seriesID uuid.UUID, source string, rating float64, voteCount int, certified bool) error
}
