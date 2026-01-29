// Package radarr provides a Radarr API v3 client and metadata provider.
package radarr

import (
	"context"
	"log/slog"
	"time"
)

// Config holds Radarr connection settings.
type Config struct {
	BaseURL string
	APIKey  string
}

// Provider is a Radarr metadata provider.
type Provider struct {
	client *Client
	cfg    Config
	logger *slog.Logger
}

// NewProvider creates a new Radarr provider.
func NewProvider(cfg Config, logger *slog.Logger) *Provider {
	p := &Provider{
		cfg:    cfg,
		logger: logger.With("provider", "radarr"),
	}

	// Only create client if configured
	if cfg.BaseURL != "" && cfg.APIKey != "" {
		client, err := NewClient(ClientConfig{
			BaseURL: cfg.BaseURL,
			APIKey:  cfg.APIKey,
		}, logger)
		if err == nil {
			p.client = client
		} else {
			logger.Warn("failed to create radarr client", "error", err)
		}
	}

	return p
}

// Name returns the provider name.
func (p *Provider) Name() string { return "radarr" }

// Priority returns the provider priority (lower = higher priority).
func (p *Provider) Priority() int { return 1 }

// IsAvailable reports whether Radarr is configured and reachable.
func (p *Provider) IsAvailable() bool {
	return p.client != nil
}

// Ping checks if Radarr is reachable.
func (p *Provider) Ping(ctx context.Context) error {
	if p.client == nil {
		return ErrUnavailable
	}
	return p.client.Ping(ctx)
}

// MovieMetadata represents metadata returned by Radarr.
type MovieMetadata struct {
	RadarrID       int
	TMDbID         int
	IMDbID         string
	Title          string
	OriginalTitle  string
	SortTitle      string
	Overview       string
	Tagline        string
	RuntimeMinutes int
	ReleaseDate    *time.Time
	InCinemas      *time.Time
	DigitalRelease *time.Time
	PhysicalRelease *time.Time
	Year           int
	Studio         string
	Certification  string
	Genres         []string
	Rating         float64
	VoteCount      int
	PosterURL      string
	BackdropURL    string
	YouTubeTrailer string
	Status         string
	HasFile        bool
	SizeOnDisk     int64
	Path           string
	Quality        string
}

// GetMovieMetadata fetches movie metadata from Radarr.
func (p *Provider) GetMovieMetadata(ctx context.Context, tmdbID int) (*MovieMetadata, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}

	movie, err := p.client.GetMovieByTMDbID(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	return p.convertMovie(movie), nil
}

// GetMovieByRadarrID fetches movie metadata by Radarr ID.
func (p *Provider) GetMovieByRadarrID(ctx context.Context, radarrID int) (*MovieMetadata, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}

	movie, err := p.client.GetMovie(ctx, radarrID)
	if err != nil {
		return nil, err
	}

	return p.convertMovie(movie), nil
}

// GetMovieByIMDbID fetches movie metadata by IMDb ID.
func (p *Provider) GetMovieByIMDbID(ctx context.Context, imdbID string) (*MovieMetadata, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}

	movie, err := p.client.GetMovieByIMDbID(ctx, imdbID)
	if err != nil {
		return nil, err
	}

	return p.convertMovie(movie), nil
}

// ListMovies lists all movies in Radarr.
func (p *Provider) ListMovies(ctx context.Context) ([]MovieMetadata, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}

	movies, err := p.client.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]MovieMetadata, 0, len(movies))
	for _, m := range movies {
		result = append(result, *p.convertMovie(&m))
	}

	return result, nil
}

// MatchMovie searches for a movie by title, year, or IMDb ID.
func (p *Provider) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*MovieMetadata, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}

	// Try IMDb ID first if available
	if imdbID != "" {
		meta, err := p.GetMovieByIMDbID(ctx, imdbID)
		if err == nil {
			return meta, nil
		}
	}

	// Fall back to listing and matching
	movies, err := p.client.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range movies {
		if m.Title == title && (year == 0 || m.Year == year) {
			return p.convertMovie(&m), nil
		}
		if m.OriginalTitle == title && (year == 0 || m.Year == year) {
			return p.convertMovie(&m), nil
		}
	}

	return nil, ErrNotFound
}

// GetSystemStatus returns Radarr system information.
func (p *Provider) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}
	return p.client.GetSystemStatus(ctx)
}

// GetHealth returns Radarr health check results.
func (p *Provider) GetHealth(ctx context.Context) ([]HealthCheck, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}
	return p.client.GetHealth(ctx)
}

// GetRootFolders returns Radarr root folder configurations.
func (p *Provider) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}
	return p.client.ListRootFolders(ctx)
}

// GetQualityProfiles returns Radarr quality profiles.
func (p *Provider) GetQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	if p.client == nil {
		return nil, ErrUnavailable
	}
	return p.client.ListQualityProfiles(ctx)
}

// convertMovie converts a Radarr Movie to MovieMetadata.
func (p *Provider) convertMovie(m *Movie) *MovieMetadata {
	meta := &MovieMetadata{
		RadarrID:       m.ID,
		TMDbID:         m.TMDbID,
		IMDbID:         m.IMDbID,
		Title:          m.Title,
		OriginalTitle:  m.OriginalTitle,
		SortTitle:      m.SortTitle,
		Overview:       m.Overview,
		RuntimeMinutes: m.Runtime,
		Year:           m.Year,
		Studio:         m.Studio,
		Certification:  m.Certification,
		Genres:         m.Genres,
		YouTubeTrailer: m.YouTubeTrailerID,
		Status:         m.Status,
		HasFile:        m.HasFile,
		SizeOnDisk:     m.SizeOnDisk,
		Path:           m.Path,
		InCinemas:      m.InCinemas,
		DigitalRelease: m.DigitalRelease,
		PhysicalRelease: m.PhysicalRelease,
	}

	// Extract rating
	if m.Ratings.TMDb != nil {
		meta.Rating = m.Ratings.TMDb.Value
		meta.VoteCount = m.Ratings.TMDb.Votes
	} else if m.Ratings.IMDb != nil {
		meta.Rating = m.Ratings.IMDb.Value
		meta.VoteCount = m.Ratings.IMDb.Votes
	}

	// Extract images
	for _, img := range m.Images {
		switch img.CoverType {
		case "poster":
			if meta.PosterURL == "" {
				meta.PosterURL = img.RemoteURL
				if meta.PosterURL == "" {
					meta.PosterURL = img.URL
				}
			}
		case "fanart":
			if meta.BackdropURL == "" {
				meta.BackdropURL = img.RemoteURL
				if meta.BackdropURL == "" {
					meta.BackdropURL = img.URL
				}
			}
		}
	}

	// Extract quality from movie file
	if m.MovieFile != nil {
		meta.Quality = m.MovieFile.Quality.Quality.Name
	}

	return meta
}

// Stats returns provider statistics.
func (p *Provider) Stats() ClientStats {
	if p.client == nil {
		return ClientStats{}
	}
	return p.client.Stats()
}

// Close releases provider resources.
func (p *Provider) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}
