package tvshow

import (
	"context"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
)

// SeriesMetadata holds normalized TV series metadata regardless of source.
type SeriesMetadata struct {
	TMDbID        int
	TvdbID        int
	IMDbID        string
	Title         string
	OriginalTitle string
	Overview      string
	FirstAirDate  time.Time
	Status        string
	Type          string
	Rating        float64
	VoteCount     int
	PosterURL     string
	BackdropURL   string
	Genres        []string
	Networks      []string
	Creators      []CreatorInfo
}

// SeasonMetadata holds normalized TV season metadata.
type SeasonMetadata struct {
	TMDbID       int
	SeasonNumber int
	Name         string
	Overview     string
	AirDate      time.Time
	PosterURL    string
	EpisodeCount int
}

// EpisodeMetadata holds normalized TV episode metadata.
type EpisodeMetadata struct {
	TMDbID         int
	TvdbID         int
	SeasonNumber   int
	EpisodeNumber  int
	Name           string
	Overview       string
	AirDate        time.Time
	Runtime        int
	StillURL       string
	Rating         float64
	VoteCount      int
	Directors      []CrewInfo
	Writers        []CrewInfo
	GuestStars     []CastInfo
}

// CreatorInfo represents series creator metadata.
type CreatorInfo struct {
	TMDbID     int
	Name       string
	ProfileURL string
}

// CastInfo represents cast member metadata.
type CastInfo struct {
	TMDbID     int
	Name       string
	Character  string
	Order      int
	ProfileURL string
}

// CrewInfo represents crew member metadata.
type CrewInfo struct {
	TMDbID     int
	Name       string
	Department string
	Job        string
	ProfileURL string
}

// MetadataProvider supplies TV show metadata in a provider-agnostic format.
type MetadataProvider interface {
	Name() string
	Priority() int
	IsAvailable() bool
	GetSeriesMetadata(ctx context.Context, tmdbID int) (*SeriesMetadata, error)
	GetSeasonMetadata(ctx context.Context, seriesTmdbID, seasonNumber int) (*SeasonMetadata, error)
	GetEpisodeMetadata(ctx context.Context, seriesTmdbID, seasonNumber, episodeNumber int) (*EpisodeMetadata, error)
	MatchSeries(ctx context.Context, title string, year int, tvdbID int, imdbID string) (*SeriesMetadata, error)
}

type tmdbAdapter struct {
	provider *tmdb.Provider
}

func newTMDbAdapter(provider *tmdb.Provider) MetadataProvider {
	if provider == nil {
		return nil
	}
	return &tmdbAdapter{provider: provider}
}

func (a *tmdbAdapter) Name() string      { return a.provider.Name() }
func (a *tmdbAdapter) Priority() int     { return a.provider.Priority() }
func (a *tmdbAdapter) IsAvailable() bool { return a.provider.IsAvailable() }

func (a *tmdbAdapter) GetSeriesMetadata(ctx context.Context, tmdbID int) (*SeriesMetadata, error) {
	// TODO: Implement when TMDb provider has TV series methods
	// For now, return error indicating not implemented
	return nil, ErrMetadataUnavailable
}

func (a *tmdbAdapter) GetSeasonMetadata(ctx context.Context, seriesTmdbID, seasonNumber int) (*SeasonMetadata, error) {
	// TODO: Implement when TMDb provider has TV season methods
	return nil, ErrMetadataUnavailable
}

func (a *tmdbAdapter) GetEpisodeMetadata(ctx context.Context, seriesTmdbID, seasonNumber, episodeNumber int) (*EpisodeMetadata, error) {
	// TODO: Implement when TMDb provider has TV episode methods
	return nil, ErrMetadataUnavailable
}

func (a *tmdbAdapter) MatchSeries(ctx context.Context, title string, year int, tvdbID int, imdbID string) (*SeriesMetadata, error) {
	// TODO: Implement when TMDb provider has TV search methods
	return nil, ErrMetadataUnavailable
}

// Future: Sonarr adapter
// type sonarrAdapter struct {
// 	provider *sonarr.Provider
// }
