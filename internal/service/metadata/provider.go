package metadata

import (
	"context"
)

// ProviderID identifies a metadata provider.
type ProviderID string

const (
	// ProviderTMDb is The Movie Database provider.
	ProviderTMDb ProviderID = "tmdb"
	// ProviderTVDb is TheTVDB provider.
	ProviderTVDb ProviderID = "tvdb"
	// ProviderFanartTV is Fanart.tv provider (for high-quality artwork).
	ProviderFanartTV ProviderID = "fanarttv"
	// ProviderOMDb is Open Movie Database provider.
	ProviderOMDb ProviderID = "omdb"
	// ProviderTVmaze is TVmaze provider (free TV metadata, no auth required).
	ProviderTVmaze ProviderID = "tvmaze"
)

// Provider represents an external metadata provider.
// Each provider (TMDb, TVDb, etc.) implements this interface.
type Provider interface {
	// ID returns the provider identifier.
	ID() ProviderID

	// Name returns the human-readable provider name.
	Name() string

	// Priority returns the provider priority (higher = preferred).
	// Used when aggregating results from multiple providers.
	Priority() int

	// SupportsMovies returns true if this provider supports movie metadata.
	SupportsMovies() bool

	// SupportsTVShows returns true if this provider supports TV show metadata.
	SupportsTVShows() bool

	// SupportsPeople returns true if this provider supports person metadata.
	SupportsPeople() bool

	// SupportsLanguage returns true if this provider supports the given language code.
	SupportsLanguage(lang string) bool

	// ClearCache clears any cached metadata for this provider.
	ClearCache()
}

// MovieProvider extends Provider with movie-specific methods.
type MovieProvider interface {
	Provider

	// SearchMovie searches for movies by title.
	SearchMovie(ctx context.Context, query string, opts SearchOptions) ([]MovieSearchResult, error)

	// GetMovie retrieves full movie details by provider ID.
	GetMovie(ctx context.Context, id string, lang string) (*MovieMetadata, error)

	// GetMovieCredits retrieves cast and crew for a movie.
	GetMovieCredits(ctx context.Context, id string) (*Credits, error)

	// GetMovieImages retrieves images (posters, backdrops) for a movie.
	GetMovieImages(ctx context.Context, id string) (*Images, error)

	// GetMovieReleaseDates retrieves release dates and certifications by country.
	GetMovieReleaseDates(ctx context.Context, id string) ([]ReleaseDate, error)

	// GetMovieTranslations retrieves available translations.
	GetMovieTranslations(ctx context.Context, id string) ([]Translation, error)

	// GetMovieExternalIDs retrieves external IDs (IMDb, etc.).
	GetMovieExternalIDs(ctx context.Context, id string) (*ExternalIDs, error)

	// GetSimilarMovies retrieves movies similar to the given movie.
	GetSimilarMovies(ctx context.Context, id string, opts SearchOptions) ([]MovieSearchResult, int, error)

	// GetMovieRecommendations retrieves recommended movies based on the given movie.
	GetMovieRecommendations(ctx context.Context, id string, opts SearchOptions) ([]MovieSearchResult, int, error)
}

// TVShowProvider extends Provider with TV show-specific methods.
type TVShowProvider interface {
	Provider

	// SearchTVShow searches for TV shows by title.
	SearchTVShow(ctx context.Context, query string, opts SearchOptions) ([]TVShowSearchResult, error)

	// GetTVShow retrieves full TV show details by provider ID.
	GetTVShow(ctx context.Context, id string, lang string) (*TVShowMetadata, error)

	// GetTVShowCredits retrieves cast and crew for a TV show.
	GetTVShowCredits(ctx context.Context, id string) (*Credits, error)

	// GetTVShowImages retrieves images for a TV show.
	GetTVShowImages(ctx context.Context, id string) (*Images, error)

	// GetTVShowContentRatings retrieves content ratings by country.
	GetTVShowContentRatings(ctx context.Context, id string) ([]ContentRating, error)

	// GetTVShowTranslations retrieves available translations.
	GetTVShowTranslations(ctx context.Context, id string) ([]Translation, error)

	// GetTVShowExternalIDs retrieves external IDs (TVDb, IMDb, etc.).
	GetTVShowExternalIDs(ctx context.Context, id string) (*ExternalIDs, error)

	// GetSeason retrieves season details.
	GetSeason(ctx context.Context, showID string, seasonNum int, lang string) (*SeasonMetadata, error)

	// GetSeasonCredits retrieves cast and crew for a season.
	GetSeasonCredits(ctx context.Context, showID string, seasonNum int) (*Credits, error)

	// GetSeasonImages retrieves images for a season.
	GetSeasonImages(ctx context.Context, showID string, seasonNum int) (*Images, error)

	// GetEpisode retrieves episode details.
	GetEpisode(ctx context.Context, showID string, seasonNum, episodeNum int, lang string) (*EpisodeMetadata, error)

	// GetEpisodeCredits retrieves cast and crew for an episode.
	GetEpisodeCredits(ctx context.Context, showID string, seasonNum, episodeNum int) (*Credits, error)

	// GetEpisodeImages retrieves images for an episode.
	GetEpisodeImages(ctx context.Context, showID string, seasonNum, episodeNum int) (*Images, error)
}

// PersonProvider extends Provider with person-specific methods.
type PersonProvider interface {
	Provider

	// SearchPerson searches for people by name.
	SearchPerson(ctx context.Context, query string, opts SearchOptions) ([]PersonSearchResult, error)

	// GetPerson retrieves full person details.
	GetPerson(ctx context.Context, id string, lang string) (*PersonMetadata, error)

	// GetPersonCredits retrieves the person's filmography.
	GetPersonCredits(ctx context.Context, id string) (*PersonCredits, error)

	// GetPersonImages retrieves images for a person.
	GetPersonImages(ctx context.Context, id string) (*Images, error)

	// GetPersonExternalIDs retrieves external IDs (IMDb, etc.).
	GetPersonExternalIDs(ctx context.Context, id string) (*ExternalIDs, error)
}

// ImageProvider extends Provider with image-specific methods.
type ImageProvider interface {
	Provider

	// GetImageURL constructs a full image URL from a path and size.
	GetImageURL(path string, size ImageSize) string

	// GetImageBaseURL returns the base URL for images.
	GetImageBaseURL() string

	// DownloadImage downloads an image by path and size.
	DownloadImage(ctx context.Context, path string, size ImageSize) ([]byte, error)
}

// CollectionProvider extends Provider with collection-specific methods.
type CollectionProvider interface {
	Provider

	// GetCollection retrieves movie collection details.
	GetCollection(ctx context.Context, id string, lang string) (*CollectionMetadata, error)

	// GetCollectionImages retrieves images for a collection.
	GetCollectionImages(ctx context.Context, id string) (*Images, error)
}

// SearchOptions configures metadata search queries.
type SearchOptions struct {
	// ProviderID selects a specific provider for the search.
	// When empty, the highest-priority provider is used with fallback.
	ProviderID ProviderID

	// Year filters results by release/air year.
	Year *int

	// Language specifies the response language (ISO 639-1, e.g., "en", "de").
	Language string

	// Page for pagination (1-indexed).
	Page int

	// IncludeAdult includes adult content in results.
	IncludeAdult bool

	// Region filters by region (ISO 3166-1, e.g., "US", "DE").
	Region string
}

// DefaultSearchOptions returns sensible defaults.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Language: "en",
		Page:     1,
	}
}

// ImageSize represents standard image sizes.
type ImageSize string

const (
	// Poster sizes
	ImageSizePosterW92       ImageSize = "w92"
	ImageSizePosterW154      ImageSize = "w154"
	ImageSizePosterW185      ImageSize = "w185"
	ImageSizePosterW342      ImageSize = "w342"
	ImageSizePosterW500      ImageSize = "w500"
	ImageSizePosterW780      ImageSize = "w780"
	ImageSizePosterOriginal  ImageSize = "original"

	// Backdrop sizes
	ImageSizeBackdropW300     ImageSize = "w300"
	ImageSizeBackdropW780     ImageSize = "w780"
	ImageSizeBackdropW1280    ImageSize = "w1280"
	ImageSizeBackdropOriginal ImageSize = "original"

	// Profile sizes
	ImageSizeProfileW45       ImageSize = "w45"
	ImageSizeProfileW185      ImageSize = "w185"
	ImageSizeProfileH632      ImageSize = "h632"
	ImageSizeProfileOriginal  ImageSize = "original"

	// Still (episode) sizes
	ImageSizeStillW92       ImageSize = "w92"
	ImageSizeStillW185      ImageSize = "w185"
	ImageSizeStillW300      ImageSize = "w300"
	ImageSizeStillOriginal  ImageSize = "original"

	// Logo sizes
	ImageSizeLogoW45       ImageSize = "w45"
	ImageSizeLogoW92       ImageSize = "w92"
	ImageSizeLogoW154      ImageSize = "w154"
	ImageSizeLogoW185      ImageSize = "w185"
	ImageSizeLogoW300      ImageSize = "w300"
	ImageSizeLogoW500      ImageSize = "w500"
	ImageSizeLogoOriginal  ImageSize = "original"
)
