package radarr

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

// ErrUnavailable indicates the Radarr provider is not configured or reachable.
var ErrUnavailable = errors.New("radarr provider unavailable")

// Config holds Radarr connection settings.
type Config struct {
	BaseURL string
	APIKey  string
}

// Provider is a Radarr metadata provider stub.
// It is intentionally minimal and only reports availability.
type Provider struct {
	cfg    Config
	logger *slog.Logger
}

// NewProvider creates a new Radarr provider.
func NewProvider(cfg Config, logger *slog.Logger) *Provider {
	return &Provider{
		cfg:    cfg,
		logger: logger.With("provider", "radarr"),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string { return "radarr" }

// Priority returns the provider priority (lower = higher priority).
func (p *Provider) Priority() int { return 1 }

// IsAvailable reports whether Radarr is configured.
func (p *Provider) IsAvailable() bool {
	return p != nil && p.cfg.BaseURL != "" && p.cfg.APIKey != ""
}

// MovieMetadata represents metadata returned by Radarr.
type MovieMetadata struct {
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

// GetMovieMetadata is a placeholder until Radarr integration is implemented.
func (p *Provider) GetMovieMetadata(ctx context.Context, tmdbID int) (*MovieMetadata, error) {
	return nil, ErrUnavailable
}

// MatchMovie is a placeholder until Radarr integration is implemented.
func (p *Provider) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*MovieMetadata, error) {
	return nil, ErrUnavailable
}
