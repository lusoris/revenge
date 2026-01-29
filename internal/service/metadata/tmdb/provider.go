package tmdb

import (
	"context"
	"log/slog"
	"time"
)

// Provider implements TMDb metadata fetching for movies.
type Provider struct {
	client *Client
	logger *slog.Logger
}

// NewProvider creates a new TMDb metadata provider.
func NewProvider(client *Client, logger *slog.Logger) *Provider {
	return &Provider{
		client: client,
		logger: logger.With("provider", "tmdb"),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "tmdb"
}

// Priority returns the provider priority (lower = higher priority).
func (p *Provider) Priority() int {
	return 2 // TMDb is secondary to Servarr (Radarr)
}

// IsAvailable reports whether the provider is configured.
func (p *Provider) IsAvailable() bool {
	return p != nil && p.client != nil && p.client.config.APIKey != ""
}

// MovieSearchResult represents a movie search result.
type MovieSearchResult struct {
	TMDbID    int
	Title     string
	Year      int
	Overview  string
	PosterURL string
	Score     float64
}

// MovieMetadata represents full movie metadata from TMDb.
type MovieMetadata struct {
	TMDbID           int
	IMDbID           string
	Title            string
	OriginalTitle    string
	OriginalLanguage string
	Overview         string
	Tagline          string
	RuntimeMinutes   int
	ReleaseDate      time.Time
	ReleaseYear      int
	Budget           int64
	Revenue          int64
	Rating           float64
	VoteCount        int
	Popularity       float64
	Adult            bool
	Status           string
	Homepage         string
	PosterURL        string
	BackdropURL      string
	Genres           []string
	Studios          []StudioInfo
	Collection       *CollectionInfo
	Cast             []CastInfo
	Crew             []CrewInfo
	Images           []ImageInfo
	Videos           []VideoInfo
}

// StudioInfo represents studio metadata.
type StudioInfo struct {
	TMDbID        int
	Name          string
	LogoURL       string
	OriginCountry string
}

// CollectionInfo represents collection metadata.
type CollectionInfo struct {
	TMDbID      int
	Name        string
	Overview    string
	PosterURL   string
	BackdropURL string
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

// ImageInfo represents image metadata.
type ImageInfo struct {
	Type        string
	URL         string
	Width       int
	Height      int
	AspectRatio float64
	Language    string
	VoteAverage float64
	VoteCount   int
}

// VideoInfo represents video metadata.
type VideoInfo struct {
	Key         string
	Name        string
	Site        string
	Type        string
	Size        int
	Official    bool
	PublishedAt string
	Language    string
}

// SearchMovies searches for movies matching the query.
func (p *Provider) SearchMovies(ctx context.Context, query string, year int) ([]MovieSearchResult, error) {
	result, err := p.client.SearchMovies(ctx, query, year, 1)
	if err != nil {
		return nil, err
	}

	results := make([]MovieSearchResult, len(result.Results))
	for i, m := range result.Results {
		results[i] = MovieSearchResult{
			TMDbID:    m.ID,
			Title:     m.Title,
			Year:      m.ReleaseYear(),
			Overview:  m.Overview,
			PosterURL: p.client.PosterURL(m.PosterPath),
			Score:     m.VoteAverage,
		}
	}

	return results, nil
}

// GetMovieMetadata fetches full movie metadata from TMDb.
func (p *Provider) GetMovieMetadata(ctx context.Context, tmdbID int) (*MovieMetadata, error) {
	m, err := p.client.GetMovie(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	metadata := &MovieMetadata{
		TMDbID:           m.ID,
		IMDbID:           m.IMDbID,
		Title:            m.Title,
		OriginalTitle:    m.OriginalTitle,
		OriginalLanguage: m.OriginalLanguage,
		Overview:         m.Overview,
		Tagline:          m.Tagline,
		RuntimeMinutes:   m.Runtime,
		ReleaseYear:      m.ReleaseYear(),
		Budget:           m.Budget,
		Revenue:          m.Revenue,
		Rating:           m.VoteAverage,
		VoteCount:        m.VoteCount,
		Popularity:       m.Popularity,
		Adult:            m.Adult,
		Status:           m.Status,
		Homepage:         m.Homepage,
		PosterURL:        p.client.PosterURL(m.PosterPath),
		BackdropURL:      p.client.BackdropURL(m.BackdropPath),
	}

	// Parse release date
	if m.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", m.ReleaseDate); err == nil {
			metadata.ReleaseDate = t
		}
	}

	// Convert genres
	metadata.Genres = make([]string, len(m.Genres))
	for i, g := range m.Genres {
		metadata.Genres[i] = g.Name
	}

	// Convert production companies (studios)
	metadata.Studios = make([]StudioInfo, len(m.ProductionCompanies))
	for i, c := range m.ProductionCompanies {
		metadata.Studios[i] = StudioInfo{
			TMDbID:        c.ID,
			Name:          c.Name,
			LogoURL:       p.client.ImageURL(c.LogoPath, "w185"),
			OriginCountry: c.OriginCountry,
		}
	}

	// Convert collection
	if m.BelongsToCollection != nil {
		metadata.Collection = &CollectionInfo{
			TMDbID:      m.BelongsToCollection.ID,
			Name:        m.BelongsToCollection.Name,
			Overview:    m.BelongsToCollection.Overview,
			PosterURL:   p.client.PosterURL(m.BelongsToCollection.PosterPath),
			BackdropURL: p.client.BackdropURL(m.BelongsToCollection.BackdropPath),
		}
	}

	// Convert credits
	if m.Credits != nil {
		metadata.Cast = make([]CastInfo, len(m.Credits.Cast))
		for i, c := range m.Credits.Cast {
			metadata.Cast[i] = CastInfo{
				TMDbID:     c.ID,
				Name:       c.Name,
				Character:  c.Character,
				Order:      c.Order,
				ProfileURL: p.client.ProfileURL(c.ProfilePath),
			}
		}

		metadata.Crew = make([]CrewInfo, len(m.Credits.Crew))
		for i, c := range m.Credits.Crew {
			metadata.Crew[i] = CrewInfo{
				TMDbID:     c.ID,
				Name:       c.Name,
				Department: c.Department,
				Job:        c.Job,
				ProfileURL: p.client.ProfileURL(c.ProfilePath),
			}
		}
	}

	// Convert images
	if m.Images != nil {
		for _, img := range m.Images.Posters {
			metadata.Images = append(metadata.Images, ImageInfo{
				Type:        "poster",
				URL:         p.client.ImageURL(img.FilePath, "original"),
				Width:       img.Width,
				Height:      img.Height,
				AspectRatio: img.AspectRatio,
				Language:    img.ISO6391,
				VoteAverage: img.VoteAverage,
				VoteCount:   img.VoteCount,
			})
		}
		for _, img := range m.Images.Backdrops {
			metadata.Images = append(metadata.Images, ImageInfo{
				Type:        "backdrop",
				URL:         p.client.ImageURL(img.FilePath, "original"),
				Width:       img.Width,
				Height:      img.Height,
				AspectRatio: img.AspectRatio,
				Language:    img.ISO6391,
				VoteAverage: img.VoteAverage,
				VoteCount:   img.VoteCount,
			})
		}
		for _, img := range m.Images.Logos {
			metadata.Images = append(metadata.Images, ImageInfo{
				Type:        "logo",
				URL:         p.client.ImageURL(img.FilePath, "original"),
				Width:       img.Width,
				Height:      img.Height,
				AspectRatio: img.AspectRatio,
				Language:    img.ISO6391,
				VoteAverage: img.VoteAverage,
				VoteCount:   img.VoteCount,
			})
		}
	}

	// Convert videos
	if m.Videos != nil {
		for _, v := range m.Videos.Results {
			metadata.Videos = append(metadata.Videos, VideoInfo{
				Key:         v.Key,
				Name:        v.Name,
				Site:        v.Site,
				Type:        v.Type,
				Size:        v.Size,
				Official:    v.Official,
				PublishedAt: v.PublishedAt,
				Language:    v.ISO6391,
			})
		}
	}

	return metadata, nil
}

// FindByIMDbID finds a movie by its IMDb ID.
func (p *Provider) FindByIMDbID(ctx context.Context, imdbID string) (*MovieSearchResult, error) {
	result, err := p.client.FindByIMDbID(ctx, imdbID)
	if err != nil {
		return nil, err
	}

	return &MovieSearchResult{
		TMDbID:    result.ID,
		Title:     result.Title,
		Year:      result.ReleaseYear(),
		Overview:  result.Overview,
		PosterURL: p.client.PosterURL(result.PosterPath),
		Score:     result.VoteAverage,
	}, nil
}

// MatchMovie attempts to match a movie by title/year/imdb.
func (p *Provider) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*MovieMetadata, error) {
	// Try IMDb ID first (most reliable)
	if imdbID != "" {
		result, err := p.FindByIMDbID(ctx, imdbID)
		if err == nil && result != nil {
			return p.GetMovieMetadata(ctx, result.TMDbID)
		}
		p.logger.Debug("IMDb lookup failed, trying title search", "imdb_id", imdbID, "error", err)
	}

	// Fall back to title search
	results, err := p.SearchMovies(ctx, title, year)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	// Return first (best) match
	return p.GetMovieMetadata(ctx, results[0].TMDbID)
}
