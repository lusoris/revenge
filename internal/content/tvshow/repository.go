package tvshow

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content"
)

// Repository defines database operations for TV shows
type Repository interface {
	// Series CRUD
	GetSeries(ctx context.Context, id uuid.UUID) (*Series, error)
	GetSeriesByTMDbID(ctx context.Context, tmdbID int32) (*Series, error)
	GetSeriesByTVDbID(ctx context.Context, tvdbID int32) (*Series, error)
	GetSeriesBySonarrID(ctx context.Context, sonarrID int32) (*Series, error)
	ListSeries(ctx context.Context, filters SeriesListFilters) ([]Series, error)
	CountSeries(ctx context.Context) (int64, error)
	SearchSeriesByTitle(ctx context.Context, query string, limit, offset int32) ([]Series, error)
	SearchSeriesByTitleAnyLanguage(ctx context.Context, query string, limit, offset int32) ([]Series, error)
	ListRecentlyAddedSeries(ctx context.Context, limit, offset int32) ([]Series, error)
	ListSeriesByGenre(ctx context.Context, tmdbGenreID int32, limit, offset int32) ([]Series, error)
	ListSeriesByNetwork(ctx context.Context, networkID uuid.UUID, limit, offset int32) ([]Series, error)
	ListSeriesByStatus(ctx context.Context, status string, limit, offset int32) ([]Series, error)
	CreateSeries(ctx context.Context, params CreateSeriesParams) (*Series, error)
	UpdateSeries(ctx context.Context, params UpdateSeriesParams) (*Series, error)
	UpdateSeriesStats(ctx context.Context, seriesID uuid.UUID) error
	DeleteSeries(ctx context.Context, id uuid.UUID) error

	// Seasons
	GetSeason(ctx context.Context, id uuid.UUID) (*Season, error)
	GetSeasonByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) (*Season, error)
	ListSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) ([]Season, error)
	ListSeasonsBySeriesWithEpisodeCount(ctx context.Context, seriesID uuid.UUID) ([]SeasonWithEpisodeCount, error)
	CreateSeason(ctx context.Context, params CreateSeasonParams) (*Season, error)
	UpsertSeason(ctx context.Context, params CreateSeasonParams) (*Season, error)
	UpdateSeason(ctx context.Context, params UpdateSeasonParams) (*Season, error)
	DeleteSeason(ctx context.Context, id uuid.UUID) error
	DeleteSeasonsBySeries(ctx context.Context, seriesID uuid.UUID) error

	// Episodes
	GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error)
	GetEpisodeByTMDbID(ctx context.Context, tmdbID int32) (*Episode, error)
	GetEpisodeByNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber, episodeNumber int32) (*Episode, error)
	ListEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) ([]Episode, error)
	ListEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) ([]Episode, error)
	ListEpisodesBySeasonNumber(ctx context.Context, seriesID uuid.UUID, seasonNumber int32) ([]Episode, error)
	ListRecentEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error)
	ListUpcomingEpisodes(ctx context.Context, limit, offset int32) ([]EpisodeWithSeriesInfo, error)
	CountEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) (int64, error)
	CountEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) (int64, error)
	CreateEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error)
	UpsertEpisode(ctx context.Context, params CreateEpisodeParams) (*Episode, error)
	UpdateEpisode(ctx context.Context, params UpdateEpisodeParams) (*Episode, error)
	DeleteEpisode(ctx context.Context, id uuid.UUID) error
	DeleteEpisodesBySeason(ctx context.Context, seasonID uuid.UUID) error
	DeleteEpisodesBySeries(ctx context.Context, seriesID uuid.UUID) error

	// Episode Files
	GetEpisodeFile(ctx context.Context, id uuid.UUID) (*EpisodeFile, error)
	GetEpisodeFileByPath(ctx context.Context, path string) (*EpisodeFile, error)
	GetEpisodeFileBySonarrID(ctx context.Context, sonarrFileID int32) (*EpisodeFile, error)
	ListEpisodeFilesByEpisode(ctx context.Context, episodeID uuid.UUID) ([]EpisodeFile, error)
	CreateEpisodeFile(ctx context.Context, params CreateEpisodeFileParams) (*EpisodeFile, error)
	UpdateEpisodeFile(ctx context.Context, params UpdateEpisodeFileParams) (*EpisodeFile, error)
	DeleteEpisodeFile(ctx context.Context, id uuid.UUID) error
	DeleteEpisodeFilesByEpisode(ctx context.Context, episodeID uuid.UUID) error

	// Series Credits
	CreateSeriesCredit(ctx context.Context, params CreateSeriesCreditParams) (*SeriesCredit, error)
	ListSeriesCast(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, error)
	ListSeriesCrew(ctx context.Context, seriesID uuid.UUID, limit, offset int32) ([]SeriesCredit, error)
	CountSeriesCast(ctx context.Context, seriesID uuid.UUID) (int64, error)
	CountSeriesCrew(ctx context.Context, seriesID uuid.UUID) (int64, error)
	DeleteSeriesCredits(ctx context.Context, seriesID uuid.UUID) error

	// Episode Credits
	CreateEpisodeCredit(ctx context.Context, params CreateEpisodeCreditParams) (*EpisodeCredit, error)
	ListEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error)
	ListEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]EpisodeCredit, error)
	DeleteEpisodeCredits(ctx context.Context, episodeID uuid.UUID) error

	// Genres
	AddSeriesGenre(ctx context.Context, seriesID uuid.UUID, tmdbGenreID int32, name string) error
	ListSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]SeriesGenre, error)
	ListDistinctSeriesGenres(ctx context.Context) ([]content.GenreSummary, error)
	DeleteSeriesGenres(ctx context.Context, seriesID uuid.UUID) error

	// Networks
	CreateNetwork(ctx context.Context, params CreateNetworkParams) (*Network, error)
	GetNetwork(ctx context.Context, id uuid.UUID) (*Network, error)
	GetNetworkByTMDbID(ctx context.Context, tmdbID int32) (*Network, error)
	ListNetworksBySeries(ctx context.Context, seriesID uuid.UUID) ([]Network, error)
	AddSeriesNetwork(ctx context.Context, seriesID, networkID uuid.UUID) error
	DeleteSeriesNetworks(ctx context.Context, seriesID uuid.UUID) error

	// Watch Progress
	CreateOrUpdateWatchProgress(ctx context.Context, params CreateWatchProgressParams) (*EpisodeWatched, error)
	MarkEpisodeWatched(ctx context.Context, userID, episodeID uuid.UUID, durationSeconds int32) (*EpisodeWatched, error)
	GetWatchProgress(ctx context.Context, userID, episodeID uuid.UUID) (*EpisodeWatched, error)
	DeleteWatchProgress(ctx context.Context, userID, episodeID uuid.UUID) error
	DeleteSeriesWatchProgress(ctx context.Context, userID, seriesID uuid.UUID) error
	ListContinueWatchingSeries(ctx context.Context, userID uuid.UUID, limit int32) ([]ContinueWatchingItem, error)
	ListWatchedEpisodesBySeries(ctx context.Context, userID, seriesID uuid.UUID) ([]WatchedEpisodeItem, error)
	ListWatchedEpisodesByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]WatchedEpisodeItem, error)
	GetSeriesWatchStats(ctx context.Context, userID, seriesID uuid.UUID) (*SeriesWatchStats, error)
	GetUserTVStats(ctx context.Context, userID uuid.UUID) (*UserTVStats, error)
	GetNextUnwatchedEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error)
}

// CreateSeriesParams contains parameters for creating a series
type CreateSeriesParams struct {
	TMDbID           *int32
	TVDbID           *int32
	IMDbID           *string
	SonarrID         *int32
	Title            string
	OriginalTitle    *string
	OriginalLanguage string
	Tagline          *string
	Overview         *string
	Status           *string
	Type             *string
	FirstAirDate     *string
	LastAirDate      *string
	VoteAverage      *string
	VoteCount        *int32
	Popularity       *string
	PosterPath       *string
	BackdropPath     *string
	TotalSeasons     int32
	TotalEpisodes    int32
	TrailerURL       *string
	Homepage         *string
	TitlesI18n       map[string]string
	TaglinesI18n     map[string]string
	OverviewsI18n    map[string]string
	AgeRatings       map[string]map[string]string
	ExternalRatings  []ExternalRating
	MetadataUpdatedAt *string
}

// UpdateSeriesParams contains parameters for updating a series
type UpdateSeriesParams struct {
	ID                uuid.UUID
	TMDbID            *int32
	TVDbID            *int32
	IMDbID            *string
	SonarrID          *int32
	Title             *string
	OriginalTitle     *string
	OriginalLanguage  *string
	Tagline           *string
	Overview          *string
	Status            *string
	Type              *string
	FirstAirDate      *string
	LastAirDate       *string
	VoteAverage       *string
	VoteCount         *int32
	Popularity        *string
	PosterPath        *string
	BackdropPath      *string
	TotalSeasons      *int32
	TotalEpisodes     *int32
	TrailerURL        *string
	Homepage          *string
	TitlesI18n        map[string]string
	TaglinesI18n      map[string]string
	OverviewsI18n     map[string]string
	AgeRatings        map[string]map[string]string
	ExternalRatings   []ExternalRating
	MetadataUpdatedAt *string
}

// CreateSeasonParams contains parameters for creating a season
type CreateSeasonParams struct {
	SeriesID      uuid.UUID
	TMDbID        *int32
	SeasonNumber  int32
	Name          string
	Overview      *string
	PosterPath    *string
	EpisodeCount  int32
	AirDate       *string
	VoteAverage   *string
	NamesI18n     map[string]string
	OverviewsI18n map[string]string
}

// UpdateSeasonParams contains parameters for updating a season
type UpdateSeasonParams struct {
	ID            uuid.UUID
	TMDbID        *int32
	SeasonNumber  *int32
	Name          *string
	Overview      *string
	PosterPath    *string
	EpisodeCount  *int32
	AirDate       *string
	VoteAverage   *string
	NamesI18n     map[string]string
	OverviewsI18n map[string]string
}

// CreateEpisodeParams contains parameters for creating an episode
type CreateEpisodeParams struct {
	SeriesID       uuid.UUID
	SeasonID       uuid.UUID
	TMDbID         *int32
	TVDbID         *int32
	IMDbID         *string
	SeasonNumber   int32
	EpisodeNumber  int32
	Title          string
	Overview       *string
	AirDate        *string
	Runtime        *int32
	VoteAverage    *string
	VoteCount      *int32
	StillPath      *string
	ProductionCode *string
	TitlesI18n     map[string]string
	OverviewsI18n  map[string]string
}

// UpdateEpisodeParams contains parameters for updating an episode
type UpdateEpisodeParams struct {
	ID             uuid.UUID
	TMDbID         *int32
	TVDbID         *int32
	IMDbID         *string
	SeasonNumber   *int32
	EpisodeNumber  *int32
	Title          *string
	Overview       *string
	AirDate        *string
	Runtime        *int32
	VoteAverage    *string
	VoteCount      *int32
	StillPath      *string
	ProductionCode *string
	TitlesI18n     map[string]string
	OverviewsI18n  map[string]string
}

// CreateEpisodeFileParams contains parameters for creating an episode file
type CreateEpisodeFileParams struct {
	EpisodeID         uuid.UUID
	FilePath          string
	FileName          string
	FileSize          int64
	Container         *string
	Resolution        *string
	QualityProfile    *string
	VideoCodec        *string
	AudioCodec        *string
	BitrateKbps       *int32
	DurationSeconds   *string
	AudioLanguages    []string
	SubtitleLanguages []string
	SonarrFileID      *int32
}

// UpdateEpisodeFileParams contains parameters for updating an episode file
type UpdateEpisodeFileParams struct {
	ID                uuid.UUID
	FilePath          *string
	FileName          *string
	FileSize          *int64
	Container         *string
	Resolution        *string
	QualityProfile    *string
	VideoCodec        *string
	AudioCodec        *string
	BitrateKbps       *int32
	DurationSeconds   *string
	AudioLanguages    []string
	SubtitleLanguages []string
	SonarrFileID      *int32
}

// CreateSeriesCreditParams contains parameters for creating a series credit
type CreateSeriesCreditParams struct {
	SeriesID     uuid.UUID
	TMDbPersonID int32
	Name         string
	CreditType   string
	Character    *string
	Job          *string
	Department   *string
	CastOrder    *int32
	ProfilePath  *string
}

// CreateEpisodeCreditParams contains parameters for creating an episode credit
type CreateEpisodeCreditParams struct {
	EpisodeID    uuid.UUID
	TMDbPersonID int32
	Name         string
	CreditType   string
	Character    *string
	Job          *string
	Department   *string
	CastOrder    *int32
	ProfilePath  *string
}

// CreateNetworkParams contains parameters for creating a network
type CreateNetworkParams struct {
	TMDbID        int32
	Name          string
	LogoPath      *string
	OriginCountry *string
}

// CreateWatchProgressParams contains parameters for creating/updating watch progress
type CreateWatchProgressParams struct {
	UserID          uuid.UUID
	EpisodeID       uuid.UUID
	ProgressSeconds int32
	DurationSeconds int32
	IsCompleted     bool
}

// SeasonWithEpisodeCount extends Season with episode count
type SeasonWithEpisodeCount struct {
	Season
	ActualEpisodeCount int64
}

// EpisodeWithSeriesInfo extends Episode with series information
type EpisodeWithSeriesInfo struct {
	Episode
	SeriesID          uuid.UUID
	SeriesTitle       string
	SeriesPosterPath  *string
}

// WatchedEpisodeItem represents a watched episode with additional info
type WatchedEpisodeItem struct {
	EpisodeWatched
	SeasonNumber     int32
	EpisodeNumber    int32
	EpisodeTitle     string
	EpisodeStillPath *string
	SeriesID         uuid.UUID
	SeriesTitle      string
	SeriesPosterPath *string
}

// GetProgressPercent returns the watch progress as a percentage.
func (c *ContinueWatchingItem) GetProgressPercent() float64 {
	if c.DurationSeconds == 0 {
		return 0
	}
	return float64(c.ProgressSeconds) / float64(c.DurationSeconds) * 100
}

// GetLastWatchedFormatted returns a human-readable last watched time
func (c *ContinueWatchingItem) GetLastWatchedFormatted() string {
	now := time.Now()
	diff := now.Sub(c.LastWatchedAt)

	switch {
	case diff < time.Minute:
		return "Just now"
	case diff < time.Hour:
		return "Less than an hour ago"
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return string(rune('0'+hours/10)) + string(rune('0'+hours%10)) + " hours ago"
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "Yesterday"
		}
		return string(rune('0'+days)) + " days ago"
	default:
		return c.LastWatchedAt.Format("Jan 2, 2006")
	}
}
