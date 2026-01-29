// Package metadata provides a unified metadata service that orchestrates
// multiple providers with intelligent fallback and caching.
package metadata

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/service/metadata/radarr"
	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
)

// Errors returned by the metadata service.
var (
	ErrNotFound         = errors.New("metadata not found")
	ErrNoProviders      = errors.New("no metadata providers available")
	ErrAllProvidersFailed = errors.New("all metadata providers failed")
)

// MovieMetadata represents unified movie metadata from any provider.
type MovieMetadata struct {
	// IDs
	TMDbID  int
	IMDbID  string
	RadarrID int

	// Basic info
	Title          string
	OriginalTitle  string
	SortTitle      string
	Overview       string
	Tagline        string
	RuntimeMinutes int
	ReleaseDate    *time.Time
	Year           int
	Status         string

	// Ratings
	Rating    float64
	VoteCount int

	// Financial (from TMDb)
	Budget  int64
	Revenue int64

	// Media info
	Certification string
	Genres        []string
	Studios       []StudioInfo
	Collection    *CollectionInfo

	// Images
	PosterURL      string
	BackdropURL    string
	YouTubeTrailer string
	Images         []ImageInfo

	// Credits (from TMDb)
	Cast []CastInfo
	Crew []CrewInfo

	// File info (from Radarr)
	HasFile    bool
	Path       string
	SizeOnDisk int64
	Quality    string

	// Metadata source
	Source     string
	FetchedAt  time.Time
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

// MovieProvider is the interface for movie metadata providers.
type MovieProvider interface {
	Name() string
	Priority() int
	IsAvailable() bool
}

// Service is the central metadata orchestration service.
// It implements the Servarr-first strategy with intelligent fallback.
type Service struct {
	radarr *radarr.Provider
	tmdb   *tmdb.Provider
	local  *cache.LocalCache
	api    *cache.APICache
	logger *slog.Logger
	mu     sync.RWMutex
}

// NewService creates a new metadata service.
func NewService(
	radarrProvider *radarr.Provider,
	tmdbProvider *tmdb.Provider,
	localCache *cache.LocalCache,
	apiCache *cache.APICache,
	logger *slog.Logger,
) *Service {
	return &Service{
		radarr: radarrProvider,
		tmdb:   tmdbProvider,
		local:  localCache,
		api:    apiCache,
		logger: logger.With("service", "metadata"),
	}
}

// GetMovieMetadata fetches movie metadata using the provider hierarchy.
// Priority: Radarr (local) -> TMDb (fallback)
func (s *Service) GetMovieMetadata(ctx context.Context, tmdbID int) (*MovieMetadata, error) {
	cacheKey := fmt.Sprintf("metadata:movie:%d", tmdbID)

	// Check local cache first (Tier 1)
	if s.local != nil {
		var cached MovieMetadata
		if s.local.GetJSON(cacheKey, &cached) {
			s.logger.Debug("metadata cache hit", "key", cacheKey, "source", "local")
			return &cached, nil
		}
	}

	// Try providers in priority order
	metadata, err := s.fetchMovieFromProviders(ctx, tmdbID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if s.local != nil && metadata != nil {
		_ = s.local.SetJSONWithTTL(cacheKey, metadata, 5*time.Minute)
	}

	return metadata, nil
}

// fetchMovieFromProviders tries providers in priority order.
func (s *Service) fetchMovieFromProviders(ctx context.Context, tmdbID int) (*MovieMetadata, error) {
	var lastErr error

	// Priority 1: Radarr (local, curated metadata)
	if s.radarr != nil && s.radarr.IsAvailable() {
		s.logger.Debug("trying radarr provider", "tmdb_id", tmdbID)
		meta, err := s.radarr.GetMovieMetadata(ctx, tmdbID)
		if err == nil && meta != nil {
			result := s.convertRadarrMetadata(meta)
			// Enhance with TMDb data if available
			if s.tmdb != nil && s.tmdb.IsAvailable() {
				s.enhanceWithTMDb(ctx, result, tmdbID)
			}
			return result, nil
		}
		if !errors.Is(err, radarr.ErrNotFound) && !errors.Is(err, radarr.ErrUnavailable) {
			lastErr = fmt.Errorf("radarr: %w", err)
		}
		s.logger.Debug("radarr provider failed, trying fallback", "error", err)
	}

	// Priority 2: TMDb (fallback for missing items)
	if s.tmdb != nil && s.tmdb.IsAvailable() {
		s.logger.Debug("trying tmdb provider", "tmdb_id", tmdbID)
		meta, err := s.tmdb.GetMovieMetadata(ctx, tmdbID)
		if err == nil && meta != nil {
			return s.convertTMDbMetadata(meta), nil
		}
		if !errors.Is(err, tmdb.ErrNotFound) {
			lastErr = fmt.Errorf("tmdb: %w", err)
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, ErrNotFound
}

// SearchMovies searches for movies across providers.
func (s *Service) SearchMovies(ctx context.Context, query string, year int) ([]MovieSearchResult, error) {
	if s.tmdb != nil && s.tmdb.IsAvailable() {
		results, err := s.tmdb.SearchMovies(ctx, query, year)
		if err != nil {
			return nil, err
		}

		searchResults := make([]MovieSearchResult, len(results))
		for i, r := range results {
			searchResults[i] = MovieSearchResult{
				TMDbID:    r.TMDbID,
				Title:     r.Title,
				Year:      r.Year,
				Overview:  r.Overview,
				PosterURL: r.PosterURL,
				Score:     r.Score,
			}
		}
		return searchResults, nil
	}

	return nil, ErrNoProviders
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

// MatchMovie attempts to match a movie by title/year/IMDb ID.
func (s *Service) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*MovieMetadata, error) {
	// Try Radarr first (local library)
	if s.radarr != nil && s.radarr.IsAvailable() {
		meta, err := s.radarr.MatchMovie(ctx, title, year, imdbID)
		if err == nil && meta != nil {
			result := s.convertRadarrMetadata(meta)
			if s.tmdb != nil && s.tmdb.IsAvailable() && result.TMDbID > 0 {
				s.enhanceWithTMDb(ctx, result, result.TMDbID)
			}
			return result, nil
		}
	}

	// Fall back to TMDb
	if s.tmdb != nil && s.tmdb.IsAvailable() {
		meta, err := s.tmdb.MatchMovie(ctx, title, year, imdbID)
		if err == nil && meta != nil {
			return s.convertTMDbMetadata(meta), nil
		}
	}

	return nil, ErrNotFound
}

// GetProviderStatus returns status of all configured providers.
func (s *Service) GetProviderStatus(ctx context.Context) []ProviderStatus {
	var status []ProviderStatus

	if s.radarr != nil {
		ps := ProviderStatus{
			Name:      s.radarr.Name(),
			Priority:  s.radarr.Priority(),
			Available: s.radarr.IsAvailable(),
		}
		if ps.Available {
			if err := s.radarr.Ping(ctx); err != nil {
				ps.Available = false
				ps.Error = err.Error()
			}
		}
		status = append(status, ps)
	}

	if s.tmdb != nil {
		ps := ProviderStatus{
			Name:      s.tmdb.Name(),
			Priority:  s.tmdb.Priority(),
			Available: s.tmdb.IsAvailable(),
		}
		status = append(status, ps)
	}

	// Sort by priority
	sort.Slice(status, func(i, j int) bool {
		return status[i].Priority < status[j].Priority
	})

	return status
}

// ProviderStatus represents the status of a metadata provider.
type ProviderStatus struct {
	Name      string
	Priority  int
	Available bool
	Error     string
}

// convertRadarrMetadata converts Radarr metadata to unified format.
func (s *Service) convertRadarrMetadata(m *radarr.MovieMetadata) *MovieMetadata {
	meta := &MovieMetadata{
		TMDbID:         m.TMDbID,
		IMDbID:         m.IMDbID,
		RadarrID:       m.RadarrID,
		Title:          m.Title,
		OriginalTitle:  m.OriginalTitle,
		SortTitle:      m.SortTitle,
		Overview:       m.Overview,
		Tagline:        m.Tagline,
		RuntimeMinutes: m.RuntimeMinutes,
		Year:           m.Year,
		Status:         m.Status,
		Rating:         m.Rating,
		VoteCount:      m.VoteCount,
		Certification:  m.Certification,
		Genres:         m.Genres,
		PosterURL:      m.PosterURL,
		BackdropURL:    m.BackdropURL,
		YouTubeTrailer: m.YouTubeTrailer,
		HasFile:        m.HasFile,
		Path:           m.Path,
		SizeOnDisk:     m.SizeOnDisk,
		Quality:        m.Quality,
		Source:         "radarr",
		FetchedAt:      time.Now(),
	}

	if m.ReleaseDate != nil {
		meta.ReleaseDate = m.ReleaseDate
	}
	if m.InCinemas != nil {
		meta.ReleaseDate = m.InCinemas
	}

	return meta
}

// convertTMDbMetadata converts TMDb metadata to unified format.
func (s *Service) convertTMDbMetadata(m *tmdb.MovieMetadata) *MovieMetadata {
	meta := &MovieMetadata{
		TMDbID:         m.TMDbID,
		IMDbID:         m.IMDbID,
		Title:          m.Title,
		OriginalTitle:  m.OriginalTitle,
		Overview:       m.Overview,
		Tagline:        m.Tagline,
		RuntimeMinutes: m.RuntimeMinutes,
		Year:           m.ReleaseYear,
		Status:         m.Status,
		Rating:         m.Rating,
		VoteCount:      m.VoteCount,
		Budget:         m.Budget,
		Revenue:        m.Revenue,
		Genres:         m.Genres,
		PosterURL:      m.PosterURL,
		BackdropURL:    m.BackdropURL,
		Source:         "tmdb",
		FetchedAt:      time.Now(),
	}

	if !m.ReleaseDate.IsZero() {
		meta.ReleaseDate = &m.ReleaseDate
	}

	// Convert studios
	meta.Studios = make([]StudioInfo, len(m.Studios))
	for i, s := range m.Studios {
		meta.Studios[i] = StudioInfo{
			TMDbID:        s.TMDbID,
			Name:          s.Name,
			LogoURL:       s.LogoURL,
			OriginCountry: s.OriginCountry,
		}
	}

	// Convert collection
	if m.Collection != nil {
		meta.Collection = &CollectionInfo{
			TMDbID:      m.Collection.TMDbID,
			Name:        m.Collection.Name,
			Overview:    m.Collection.Overview,
			PosterURL:   m.Collection.PosterURL,
			BackdropURL: m.Collection.BackdropURL,
		}
	}

	// Convert cast
	meta.Cast = make([]CastInfo, len(m.Cast))
	for i, c := range m.Cast {
		meta.Cast[i] = CastInfo{
			TMDbID:     c.TMDbID,
			Name:       c.Name,
			Character:  c.Character,
			Order:      c.Order,
			ProfileURL: c.ProfileURL,
		}
	}

	// Convert crew
	meta.Crew = make([]CrewInfo, len(m.Crew))
	for i, c := range m.Crew {
		meta.Crew[i] = CrewInfo{
			TMDbID:     c.TMDbID,
			Name:       c.Name,
			Department: c.Department,
			Job:        c.Job,
			ProfileURL: c.ProfileURL,
		}
	}

	// Convert images
	meta.Images = make([]ImageInfo, len(m.Images))
	for i, img := range m.Images {
		meta.Images[i] = ImageInfo{
			Type:        img.Type,
			URL:         img.URL,
			Width:       img.Width,
			Height:      img.Height,
			AspectRatio: img.AspectRatio,
			Language:    img.Language,
			VoteAverage: img.VoteAverage,
			VoteCount:   img.VoteCount,
		}
	}

	// Find YouTube trailer
	for _, v := range m.Videos {
		if v.Site == "YouTube" && (v.Type == "Trailer" || v.Type == "Teaser") {
			meta.YouTubeTrailer = v.Key
			break
		}
	}

	return meta
}

// enhanceWithTMDb adds TMDb-specific data to Radarr metadata.
func (s *Service) enhanceWithTMDb(ctx context.Context, meta *MovieMetadata, tmdbID int) {
	tmdbMeta, err := s.tmdb.GetMovieMetadata(ctx, tmdbID)
	if err != nil {
		s.logger.Debug("failed to enhance with tmdb", "tmdb_id", tmdbID, "error", err)
		return
	}

	// Add data that Radarr doesn't provide
	if meta.Budget == 0 {
		meta.Budget = tmdbMeta.Budget
	}
	if meta.Revenue == 0 {
		meta.Revenue = tmdbMeta.Revenue
	}
	if meta.Tagline == "" {
		meta.Tagline = tmdbMeta.Tagline
	}
	if len(meta.Studios) == 0 {
		meta.Studios = make([]StudioInfo, len(tmdbMeta.Studios))
		for i, s := range tmdbMeta.Studios {
			meta.Studios[i] = StudioInfo{
				TMDbID:        s.TMDbID,
				Name:          s.Name,
				LogoURL:       s.LogoURL,
				OriginCountry: s.OriginCountry,
			}
		}
	}
	if meta.Collection == nil && tmdbMeta.Collection != nil {
		meta.Collection = &CollectionInfo{
			TMDbID:      tmdbMeta.Collection.TMDbID,
			Name:        tmdbMeta.Collection.Name,
			Overview:    tmdbMeta.Collection.Overview,
			PosterURL:   tmdbMeta.Collection.PosterURL,
			BackdropURL: tmdbMeta.Collection.BackdropURL,
		}
	}
	if len(meta.Cast) == 0 {
		meta.Cast = make([]CastInfo, len(tmdbMeta.Cast))
		for i, c := range tmdbMeta.Cast {
			meta.Cast[i] = CastInfo{
				TMDbID:     c.TMDbID,
				Name:       c.Name,
				Character:  c.Character,
				Order:      c.Order,
				ProfileURL: c.ProfileURL,
			}
		}
	}
	if len(meta.Crew) == 0 {
		meta.Crew = make([]CrewInfo, len(tmdbMeta.Crew))
		for i, c := range tmdbMeta.Crew {
			meta.Crew[i] = CrewInfo{
				TMDbID:     c.TMDbID,
				Name:       c.Name,
				Department: c.Department,
				Job:        c.Job,
				ProfileURL: c.ProfileURL,
			}
		}
	}
	if len(meta.Images) == 0 {
		meta.Images = make([]ImageInfo, len(tmdbMeta.Images))
		for i, img := range tmdbMeta.Images {
			meta.Images[i] = ImageInfo{
				Type:        img.Type,
				URL:         img.URL,
				Width:       img.Width,
				Height:      img.Height,
				AspectRatio: img.AspectRatio,
				Language:    img.Language,
				VoteAverage: img.VoteAverage,
				VoteCount:   img.VoteCount,
			}
		}
	}
	if meta.YouTubeTrailer == "" {
		for _, v := range tmdbMeta.Videos {
			if v.Site == "YouTube" && (v.Type == "Trailer" || v.Type == "Teaser") {
				meta.YouTubeTrailer = v.Key
				break
			}
		}
	}
}
