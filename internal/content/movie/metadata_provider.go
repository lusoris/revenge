package movie

import (
	"context"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata/radarr"
	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
)

// Metadata holds normalized movie metadata regardless of source.
type Metadata struct {
	TMDbID         int
	IMDbID         string
	Title          string
	OriginalTitle  string
	Overview       string
	Tagline        string
	RuntimeMinutes int
	ReleaseDate    time.Time
	Budget         int64
	Revenue        int64
	Rating         float64
	VoteCount      int
	PosterURL      string
	BackdropURL    string
}

// MetadataProvider supplies movie metadata in a provider-agnostic format.
type MetadataProvider interface {
	Name() string
	Priority() int
	IsAvailable() bool
	GetMovieMetadata(ctx context.Context, tmdbID int) (*Metadata, error)
	MatchMovie(ctx context.Context, title string, year int, imdbID string) (*Metadata, error)
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

func (a *tmdbAdapter) GetMovieMetadata(ctx context.Context, tmdbID int) (*Metadata, error) {
	meta, err := a.provider.GetMovieMetadata(ctx, tmdbID)
	if err != nil {
		return nil, err
	}
	return normalizeTMDb(meta), nil
}

func (a *tmdbAdapter) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*Metadata, error) {
	meta, err := a.provider.MatchMovie(ctx, title, year, imdbID)
	if err != nil {
		return nil, err
	}
	return normalizeTMDb(meta), nil
}

func normalizeTMDb(meta *tmdb.MovieMetadata) *Metadata {
	if meta == nil {
		return nil
	}
	return &Metadata{
		TMDbID:         meta.TMDbID,
		IMDbID:         meta.IMDbID,
		Title:          meta.Title,
		OriginalTitle:  meta.OriginalTitle,
		Overview:       meta.Overview,
		Tagline:        meta.Tagline,
		RuntimeMinutes: meta.RuntimeMinutes,
		ReleaseDate:    meta.ReleaseDate,
		Budget:         meta.Budget,
		Revenue:        meta.Revenue,
		Rating:         meta.Rating,
		VoteCount:      meta.VoteCount,
		PosterURL:      meta.PosterURL,
		BackdropURL:    meta.BackdropURL,
	}
}

type radarrAdapter struct {
	provider *radarr.Provider
}

func newRadarrAdapter(provider *radarr.Provider) MetadataProvider {
	if provider == nil {
		return nil
	}
	return &radarrAdapter{provider: provider}
}

func (a *radarrAdapter) Name() string      { return a.provider.Name() }
func (a *radarrAdapter) Priority() int     { return a.provider.Priority() }
func (a *radarrAdapter) IsAvailable() bool { return a.provider.IsAvailable() }

func (a *radarrAdapter) GetMovieMetadata(ctx context.Context, tmdbID int) (*Metadata, error) {
	meta, err := a.provider.GetMovieMetadata(ctx, tmdbID)
	if err != nil {
		return nil, err
	}
	return normalizeRadarr(meta), nil
}

func (a *radarrAdapter) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*Metadata, error) {
	meta, err := a.provider.MatchMovie(ctx, title, year, imdbID)
	if err != nil {
		return nil, err
	}
	return normalizeRadarr(meta), nil
}

func normalizeRadarr(meta *radarr.MovieMetadata) *Metadata {
	if meta == nil {
		return nil
	}
	return &Metadata{
		TMDbID:         meta.TMDbID,
		IMDbID:         meta.IMDbID,
		Title:          meta.Title,
		OriginalTitle:  meta.OriginalTitle,
		Overview:       meta.Overview,
		Tagline:        meta.Tagline,
		RuntimeMinutes: meta.RuntimeMinutes,
		ReleaseDate:    meta.ReleaseDate,
		Budget:         meta.Budget,
		Revenue:        meta.Revenue,
		Rating:         meta.Rating,
		VoteCount:      meta.VoteCount,
		PosterURL:      meta.PosterURL,
		BackdropURL:    meta.BackdropURL,
	}
}
