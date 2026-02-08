package metadata

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/google/uuid"
)

// Service provides a unified interface for metadata operations.
// It aggregates multiple providers and handles fallback logic.
type Service interface {
	// Movie operations
	SearchMovie(ctx context.Context, query string, opts SearchOptions) ([]MovieSearchResult, error)
	GetMovieMetadata(ctx context.Context, tmdbID int32, languages []string) (*MovieMetadata, error)
	GetMovieCredits(ctx context.Context, tmdbID int32) (*Credits, error)
	GetMovieImages(ctx context.Context, tmdbID int32) (*Images, error)
	GetMovieReleaseDates(ctx context.Context, tmdbID int32) ([]ReleaseDate, error)
	GetMovieExternalIDs(ctx context.Context, tmdbID int32) (*ExternalIDs, error)
	GetSimilarMovies(ctx context.Context, tmdbID int32, opts SearchOptions) ([]MovieSearchResult, int, error)
	GetMovieRecommendations(ctx context.Context, tmdbID int32, opts SearchOptions) ([]MovieSearchResult, int, error)

	// TV Show operations
	SearchTVShow(ctx context.Context, query string, opts SearchOptions) ([]TVShowSearchResult, error)
	GetTVShowMetadata(ctx context.Context, tmdbID int32, languages []string) (*TVShowMetadata, error)
	GetTVShowCredits(ctx context.Context, tmdbID int32) (*Credits, error)
	GetTVShowImages(ctx context.Context, tmdbID int32) (*Images, error)
	GetTVShowContentRatings(ctx context.Context, tmdbID int32) ([]ContentRating, error)
	GetTVShowExternalIDs(ctx context.Context, tmdbID int32) (*ExternalIDs, error)
	GetSeasonMetadata(ctx context.Context, tmdbID int32, seasonNum int, languages []string) (*SeasonMetadata, error)
	GetSeasonImages(ctx context.Context, tmdbID int32, seasonNum int) (*Images, error)
	GetEpisodeMetadata(ctx context.Context, tmdbID int32, seasonNum, episodeNum int, languages []string) (*EpisodeMetadata, error)
	GetEpisodeImages(ctx context.Context, tmdbID int32, seasonNum, episodeNum int) (*Images, error)

	// Person operations
	SearchPerson(ctx context.Context, query string, opts SearchOptions) ([]PersonSearchResult, error)
	GetPersonMetadata(ctx context.Context, tmdbID int32, languages []string) (*PersonMetadata, error)
	GetPersonCredits(ctx context.Context, tmdbID int32) (*PersonCredits, error)
	GetPersonImages(ctx context.Context, tmdbID int32) (*Images, error)

	// Collection operations
	GetCollectionMetadata(ctx context.Context, tmdbID int32, languages []string) (*CollectionMetadata, error)

	// Image operations
	GetImageURL(path string, size ImageSize) string

	// Refresh operations (triggers async jobs)
	RefreshMovie(ctx context.Context, movieID uuid.UUID) error
	RefreshTVShow(ctx context.Context, seriesID uuid.UUID) error

	// Cache management
	ClearCache()

	// Provider management
	RegisterProvider(provider Provider)
	GetProviders() []Provider
}

// ServiceConfig configures the metadata service.
type ServiceConfig struct {
	// DefaultLanguages are fetched if no specific languages are requested.
	DefaultLanguages []string

	// EnableProviderFallback enables trying secondary providers on failure.
	EnableProviderFallback bool

	// EnableEnrichment enables merging data from multiple providers.
	EnableEnrichment bool
}

// DefaultServiceConfig returns a config with sensible defaults.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		DefaultLanguages:       []string{"en"},
		EnableProviderFallback: true,
		EnableEnrichment:       false,
	}
}

// service implements the Service interface.
type service struct {
	config ServiceConfig

	// Providers by type
	movieProviders  []MovieProvider
	tvProviders     []TVShowProvider
	personProviders []PersonProvider
	imageProviders  []ImageProvider
	collectionProviders []CollectionProvider
	allProviders    []Provider

	mu sync.RWMutex

	// Job queue interface (will be set by fx module)
	jobQueue JobQueue
}

// JobQueue is the interface for submitting metadata refresh jobs.
type JobQueue interface {
	EnqueueRefreshMovie(ctx context.Context, movieID uuid.UUID, force bool, languages []string) error
	EnqueueRefreshTVShow(ctx context.Context, seriesID uuid.UUID, force bool, languages []string) error
}

// NewService creates a new metadata service.
func NewService(config ServiceConfig) *service {
	return &service{
		config: config,
	}
}

// SetJobQueue sets the job queue for async operations.
func (s *service) SetJobQueue(queue JobQueue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobQueue = queue
}

// RegisterProvider adds a provider to the service.
func (s *service) RegisterProvider(provider Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.allProviders = append(s.allProviders, provider)

	// Sort by priority (highest first)
	sort.Slice(s.allProviders, func(i, j int) bool {
		return s.allProviders[i].Priority() > s.allProviders[j].Priority()
	})

	// Register by capability
	if mp, ok := provider.(MovieProvider); ok && provider.SupportsMovies() {
		s.movieProviders = append(s.movieProviders, mp)
		sort.Slice(s.movieProviders, func(i, j int) bool {
			return s.movieProviders[i].Priority() > s.movieProviders[j].Priority()
		})
	}

	if tp, ok := provider.(TVShowProvider); ok && provider.SupportsTVShows() {
		s.tvProviders = append(s.tvProviders, tp)
		sort.Slice(s.tvProviders, func(i, j int) bool {
			return s.tvProviders[i].Priority() > s.tvProviders[j].Priority()
		})
	}

	if pp, ok := provider.(PersonProvider); ok && provider.SupportsPeople() {
		s.personProviders = append(s.personProviders, pp)
		sort.Slice(s.personProviders, func(i, j int) bool {
			return s.personProviders[i].Priority() > s.personProviders[j].Priority()
		})
	}

	if ip, ok := provider.(ImageProvider); ok {
		s.imageProviders = append(s.imageProviders, ip)
		sort.Slice(s.imageProviders, func(i, j int) bool {
			return s.imageProviders[i].Priority() > s.imageProviders[j].Priority()
		})
	}

	if cp, ok := provider.(CollectionProvider); ok {
		s.collectionProviders = append(s.collectionProviders, cp)
		sort.Slice(s.collectionProviders, func(i, j int) bool {
			return s.collectionProviders[i].Priority() > s.collectionProviders[j].Priority()
		})
	}
}

// GetProviders returns all registered providers.
func (s *service) GetProviders() []Provider {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Provider, len(s.allProviders))
	copy(result, s.allProviders)
	return result
}

// SearchMovie searches for movies using the primary provider.
func (s *service) SearchMovie(ctx context.Context, query string, opts SearchOptions) ([]MovieSearchResult, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if opts.Language == "" {
		opts.Language = s.getDefaultLanguage()
	}

	// Route to specific provider if requested
	if opts.ProviderID != "" {
		for _, p := range providers {
			if p.ID() == opts.ProviderID {
				return p.SearchMovie(ctx, query, opts)
			}
		}
		return nil, fmt.Errorf("provider %q does not support movie search: %w", opts.ProviderID, ErrNoProviders)
	}

	// Use primary provider with fallback
	results, err := providers[0].SearchMovie(ctx, query, opts)
	if err != nil && s.config.EnableProviderFallback && len(providers) > 1 {
		for _, p := range providers[1:] {
			results, err = p.SearchMovie(ctx, query, opts)
			if err == nil {
				break
			}
		}
	}

	return results, err
}

// GetMovieMetadata retrieves movie metadata.
func (s *service) GetMovieMetadata(ctx context.Context, tmdbID int32, languages []string) (*MovieMetadata, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if len(languages) == 0 {
		languages = s.config.DefaultLanguages
	}

	id := fmt.Sprintf("%d", tmdbID)
	var result *MovieMetadata
	var aggErr AggregateError

	// Get from primary provider
	for _, lang := range languages {
		metadata, err := providers[0].GetMovie(ctx, id, lang)
		if err != nil {
			aggErr.Add(err)
			continue
		}

		if result == nil {
			result = metadata
		} else {
			// Merge translations
			if result.Translations == nil {
				result.Translations = make(map[string]*LocalizedMovieData)
			}
			result.Translations[lang] = &LocalizedMovieData{
				Language: lang,
				Title:    metadata.Title,
				Overview: ptrToString(metadata.Overview),
				Tagline:  ptrToString(metadata.Tagline),
			}
		}
	}

	if result == nil {
		return nil, aggErr.First()
	}

	return result, nil
}

// GetMovieCredits retrieves movie credits.
func (s *service) GetMovieCredits(ctx context.Context, tmdbID int32) (*Credits, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetMovieCredits(ctx, id)
}

// GetMovieImages retrieves movie images.
func (s *service) GetMovieImages(ctx context.Context, tmdbID int32) (*Images, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetMovieImages(ctx, id)
}

// GetMovieReleaseDates retrieves movie release dates.
func (s *service) GetMovieReleaseDates(ctx context.Context, tmdbID int32) ([]ReleaseDate, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetMovieReleaseDates(ctx, id)
}

// GetMovieExternalIDs retrieves movie external IDs.
func (s *service) GetMovieExternalIDs(ctx context.Context, tmdbID int32) (*ExternalIDs, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetMovieExternalIDs(ctx, id)
}

// GetSimilarMovies retrieves movies similar to the given movie.
func (s *service) GetSimilarMovies(ctx context.Context, tmdbID int32, opts SearchOptions) ([]MovieSearchResult, int, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, 0, ErrNoProviders
	}

	if opts.Language == "" {
		opts.Language = s.getDefaultLanguage()
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetSimilarMovies(ctx, id, opts)
}

// GetMovieRecommendations retrieves recommended movies based on the given movie.
func (s *service) GetMovieRecommendations(ctx context.Context, tmdbID int32, opts SearchOptions) ([]MovieSearchResult, int, error) {
	s.mu.RLock()
	providers := s.movieProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, 0, ErrNoProviders
	}

	if opts.Language == "" {
		opts.Language = s.getDefaultLanguage()
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetMovieRecommendations(ctx, id, opts)
}

// SearchTVShow searches for TV shows.
func (s *service) SearchTVShow(ctx context.Context, query string, opts SearchOptions) ([]TVShowSearchResult, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if opts.Language == "" {
		opts.Language = s.getDefaultLanguage()
	}

	// Route to specific provider if requested
	if opts.ProviderID != "" {
		for _, p := range providers {
			if p.ID() == opts.ProviderID {
				return p.SearchTVShow(ctx, query, opts)
			}
		}
		return nil, fmt.Errorf("provider %q does not support TV show search: %w", opts.ProviderID, ErrNoProviders)
	}

	// Use primary provider with fallback
	results, err := providers[0].SearchTVShow(ctx, query, opts)
	if err != nil && s.config.EnableProviderFallback && len(providers) > 1 {
		for _, p := range providers[1:] {
			results, err = p.SearchTVShow(ctx, query, opts)
			if err == nil {
				break
			}
		}
	}

	return results, err
}

// GetTVShowMetadata retrieves TV show metadata.
func (s *service) GetTVShowMetadata(ctx context.Context, tmdbID int32, languages []string) (*TVShowMetadata, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if len(languages) == 0 {
		languages = s.config.DefaultLanguages
	}

	id := fmt.Sprintf("%d", tmdbID)
	var result *TVShowMetadata
	var aggErr AggregateError

	for _, lang := range languages {
		metadata, err := providers[0].GetTVShow(ctx, id, lang)
		if err != nil {
			aggErr.Add(err)
			continue
		}

		if result == nil {
			result = metadata
		} else {
			if result.Translations == nil {
				result.Translations = make(map[string]*LocalizedTVShowData)
			}
			result.Translations[lang] = &LocalizedTVShowData{
				Language: lang,
				Name:     metadata.Name,
				Overview: ptrToString(metadata.Overview),
				Tagline:  ptrToString(metadata.Tagline),
			}
		}
	}

	if result == nil {
		return nil, aggErr.First()
	}

	return result, nil
}

// GetTVShowCredits retrieves TV show credits.
func (s *service) GetTVShowCredits(ctx context.Context, tmdbID int32) (*Credits, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetTVShowCredits(ctx, id)
}

// GetTVShowImages retrieves TV show images.
func (s *service) GetTVShowImages(ctx context.Context, tmdbID int32) (*Images, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetTVShowImages(ctx, id)
}

// GetTVShowContentRatings retrieves TV show content ratings.
func (s *service) GetTVShowContentRatings(ctx context.Context, tmdbID int32) ([]ContentRating, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetTVShowContentRatings(ctx, id)
}

// GetTVShowExternalIDs retrieves TV show external IDs.
func (s *service) GetTVShowExternalIDs(ctx context.Context, tmdbID int32) (*ExternalIDs, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetTVShowExternalIDs(ctx, id)
}

// GetSeasonMetadata retrieves season metadata.
func (s *service) GetSeasonMetadata(ctx context.Context, tmdbID int32, seasonNum int, languages []string) (*SeasonMetadata, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if len(languages) == 0 {
		languages = s.config.DefaultLanguages
	}

	id := fmt.Sprintf("%d", tmdbID)
	var result *SeasonMetadata
	var aggErr AggregateError

	for _, lang := range languages {
		metadata, err := providers[0].GetSeason(ctx, id, seasonNum, lang)
		if err != nil {
			aggErr.Add(err)
			continue
		}

		if result == nil {
			result = metadata
		} else {
			if result.Translations == nil {
				result.Translations = make(map[string]*LocalizedSeasonData)
			}
			result.Translations[lang] = &LocalizedSeasonData{
				Language: lang,
				Name:     metadata.Name,
				Overview: ptrToString(metadata.Overview),
			}
		}
	}

	if result == nil {
		return nil, aggErr.First()
	}

	return result, nil
}

// GetSeasonImages retrieves images for a season.
func (s *service) GetSeasonImages(ctx context.Context, tmdbID int32, seasonNum int) (*Images, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	for _, p := range providers {
		images, err := p.GetSeasonImages(ctx, id, seasonNum)
		if err != nil {
			continue
		}
		return images, nil
	}

	return nil, ErrNotFound
}

// GetEpisodeMetadata retrieves episode metadata.
func (s *service) GetEpisodeMetadata(ctx context.Context, tmdbID int32, seasonNum, episodeNum int, languages []string) (*EpisodeMetadata, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if len(languages) == 0 {
		languages = s.config.DefaultLanguages
	}

	id := fmt.Sprintf("%d", tmdbID)
	var result *EpisodeMetadata
	var aggErr AggregateError

	for _, lang := range languages {
		metadata, err := providers[0].GetEpisode(ctx, id, seasonNum, episodeNum, lang)
		if err != nil {
			aggErr.Add(err)
			continue
		}

		if result == nil {
			result = metadata
		} else {
			if result.Translations == nil {
				result.Translations = make(map[string]*LocalizedEpisodeData)
			}
			result.Translations[lang] = &LocalizedEpisodeData{
				Language: lang,
				Name:     metadata.Name,
				Overview: ptrToString(metadata.Overview),
			}
		}
	}

	if result == nil {
		return nil, aggErr.First()
	}

	return result, nil
}

// GetEpisodeImages retrieves images for an episode.
func (s *service) GetEpisodeImages(ctx context.Context, tmdbID int32, seasonNum, episodeNum int) (*Images, error) {
	s.mu.RLock()
	providers := s.tvProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	for _, p := range providers {
		images, err := p.GetEpisodeImages(ctx, id, seasonNum, episodeNum)
		if err != nil {
			continue
		}
		return images, nil
	}

	return nil, ErrNotFound
}

// SearchPerson searches for people.
func (s *service) SearchPerson(ctx context.Context, query string, opts SearchOptions) ([]PersonSearchResult, error) {
	s.mu.RLock()
	providers := s.personProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if opts.Language == "" {
		opts.Language = s.getDefaultLanguage()
	}

	// Route to specific provider if requested
	if opts.ProviderID != "" {
		for _, p := range providers {
			if p.ID() == opts.ProviderID {
				return p.SearchPerson(ctx, query, opts)
			}
		}
		return nil, fmt.Errorf("provider %q does not support person search: %w", opts.ProviderID, ErrNoProviders)
	}

	// Use primary provider with fallback
	results, err := providers[0].SearchPerson(ctx, query, opts)
	if err != nil && s.config.EnableProviderFallback && len(providers) > 1 {
		for _, p := range providers[1:] {
			results, err = p.SearchPerson(ctx, query, opts)
			if err == nil {
				break
			}
		}
	}

	return results, err
}

// GetPersonMetadata retrieves person metadata.
func (s *service) GetPersonMetadata(ctx context.Context, tmdbID int32, languages []string) (*PersonMetadata, error) {
	s.mu.RLock()
	providers := s.personProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if len(languages) == 0 {
		languages = s.config.DefaultLanguages
	}

	id := fmt.Sprintf("%d", tmdbID)
	var result *PersonMetadata
	var aggErr AggregateError

	for _, lang := range languages {
		metadata, err := providers[0].GetPerson(ctx, id, lang)
		if err != nil {
			aggErr.Add(err)
			continue
		}

		if result == nil {
			result = metadata
		} else {
			if result.Translations == nil {
				result.Translations = make(map[string]*LocalizedPersonData)
			}
			result.Translations[lang] = &LocalizedPersonData{
				Language:  lang,
				Biography: ptrToString(metadata.Biography),
			}
		}
	}

	if result == nil {
		return nil, aggErr.First()
	}

	return result, nil
}

// GetPersonCredits retrieves person credits.
func (s *service) GetPersonCredits(ctx context.Context, tmdbID int32) (*PersonCredits, error) {
	s.mu.RLock()
	providers := s.personProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetPersonCredits(ctx, id)
}

// GetPersonImages retrieves person images.
func (s *service) GetPersonImages(ctx context.Context, tmdbID int32) (*Images, error) {
	s.mu.RLock()
	providers := s.personProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetPersonImages(ctx, id)
}

// GetCollectionMetadata retrieves collection metadata.
func (s *service) GetCollectionMetadata(ctx context.Context, tmdbID int32, languages []string) (*CollectionMetadata, error) {
	s.mu.RLock()
	providers := s.collectionProviders
	s.mu.RUnlock()

	if len(providers) == 0 {
		return nil, ErrNoProviders
	}

	if len(languages) == 0 {
		languages = s.config.DefaultLanguages
	}

	id := fmt.Sprintf("%d", tmdbID)
	return providers[0].GetCollection(ctx, id, languages[0])
}

// GetImageURL constructs an image URL using the primary image provider.
func (s *service) GetImageURL(path string, size ImageSize) string {
	s.mu.RLock()
	providers := s.imageProviders
	s.mu.RUnlock()

	if len(providers) == 0 || path == "" {
		return ""
	}

	return providers[0].GetImageURL(path, size)
}

// RefreshMovie triggers an async movie metadata refresh.
func (s *service) RefreshMovie(ctx context.Context, movieID uuid.UUID) error {
	s.mu.RLock()
	queue := s.jobQueue
	s.mu.RUnlock()

	if queue == nil {
		return fmt.Errorf("metadata: job queue not configured")
	}

	return queue.EnqueueRefreshMovie(ctx, movieID, false, s.config.DefaultLanguages)
}

// RefreshTVShow triggers an async TV show metadata refresh.
func (s *service) RefreshTVShow(ctx context.Context, seriesID uuid.UUID) error {
	s.mu.RLock()
	queue := s.jobQueue
	s.mu.RUnlock()

	if queue == nil {
		return fmt.Errorf("metadata: job queue not configured")
	}

	return queue.EnqueueRefreshTVShow(ctx, seriesID, false, s.config.DefaultLanguages)
}

// ClearCache clears cached metadata across all registered providers.
func (s *service) ClearCache() {
	s.mu.RLock()
	providers := s.allProviders
	s.mu.RUnlock()

	for _, p := range providers {
		p.ClearCache()
	}
}

// getDefaultLanguage returns the first default language.
func (s *service) getDefaultLanguage() string {
	if len(s.config.DefaultLanguages) > 0 {
		return s.config.DefaultLanguages[0]
	}
	return "en"
}

// ptrToString returns the value of a string pointer or empty string.
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Ensure service implements Service interface.
var _ Service = (*service)(nil)
