// Package tmdb provides a TMDb (The Movie Database) provider implementation.
package tmdb

import (
	"time"
)

// API response types for TMDb.

// MovieResponse is the TMDb API response for movie details.
type MovieResponse struct {
	ID                  int                   `json:"id"`
	IMDbID              *string               `json:"imdb_id"`
	Title               string                `json:"title"`
	OriginalTitle       string                `json:"original_title"`
	OriginalLanguage    string                `json:"original_language"`
	Overview            *string               `json:"overview"`
	Tagline             *string               `json:"tagline"`
	ReleaseDate         string                `json:"release_date"`
	Runtime             *int                  `json:"runtime"`
	Budget              *int64                `json:"budget"`
	Revenue             *int64                `json:"revenue"`
	Status              string                `json:"status"`
	VoteAverage         float64               `json:"vote_average"`
	VoteCount           int                   `json:"vote_count"`
	Popularity          float64               `json:"popularity"`
	Adult               bool                  `json:"adult"`
	Video               bool                  `json:"video"`
	PosterPath          *string               `json:"poster_path"`
	BackdropPath        *string               `json:"backdrop_path"`
	Homepage            *string               `json:"homepage"`
	Genres              []GenreResponse       `json:"genres"`
	ProductionCompanies []CompanyResponse     `json:"production_companies"`
	ProductionCountries []CountryResponse     `json:"production_countries"`
	SpokenLanguages     []LanguageResponse    `json:"spoken_languages"`
	BelongsToCollection *CollectionRefResponse `json:"belongs_to_collection"`

	// Appended data (via append_to_response)
	Credits       *CreditsResponse      `json:"credits,omitempty"`
	Images        *ImagesResponse       `json:"images,omitempty"`
	ReleaseDates  *ReleaseDatesWrapper  `json:"release_dates,omitempty"`
	Translations  *TranslationsWrapper  `json:"translations,omitempty"`
	ExternalIDs   *ExternalIDsResponse  `json:"external_ids,omitempty"`
	Videos        *VideosResponse       `json:"videos,omitempty"`
}

// SearchResultsResponse is the TMDb API response for search results.
type SearchResultsResponse struct {
	Page         int                   `json:"page"`
	Results      []MovieSearchResponse `json:"results"`
	TotalPages   int                   `json:"total_pages"`
	TotalResults int                   `json:"total_results"`
}

// MovieSearchResponse is a single movie search result.
type MovieSearchResponse struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       *string `json:"poster_path"`
	BackdropPath     *string `json:"backdrop_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	Adult            bool    `json:"adult"`
	Video            bool    `json:"video"`
	GenreIDs         []int   `json:"genre_ids"`
}

// TVSearchResultsResponse is the TMDb API response for TV search results.
type TVSearchResultsResponse struct {
	Page         int                  `json:"page"`
	Results      []TVSearchResponse   `json:"results"`
	TotalPages   int                  `json:"total_pages"`
	TotalResults int                  `json:"total_results"`
}

// TVSearchResponse is a single TV show search result.
type TVSearchResponse struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
	OriginalLanguage string   `json:"original_language"`
	Overview         string   `json:"overview"`
	FirstAirDate     string   `json:"first_air_date"`
	PosterPath       *string  `json:"poster_path"`
	BackdropPath     *string  `json:"backdrop_path"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	Popularity       float64  `json:"popularity"`
	Adult            bool     `json:"adult"`
	GenreIDs         []int    `json:"genre_ids"`
	OriginCountry    []string `json:"origin_country"`
}

// TVResponse is the TMDb API response for TV show details.
type TVResponse struct {
	ID               int                   `json:"id"`
	Name             string                `json:"name"`
	OriginalName     string                `json:"original_name"`
	OriginalLanguage string                `json:"original_language"`
	Overview         *string               `json:"overview"`
	Tagline          *string               `json:"tagline"`
	Status           string                `json:"status"`
	Type             string                `json:"type"`
	FirstAirDate     string                `json:"first_air_date"`
	LastAirDate      string                `json:"last_air_date"`
	InProduction     bool                  `json:"in_production"`
	NumberOfSeasons  int                   `json:"number_of_seasons"`
	NumberOfEpisodes int                   `json:"number_of_episodes"`
	EpisodeRunTime   []int                 `json:"episode_run_time"`
	VoteAverage      float64               `json:"vote_average"`
	VoteCount        int                   `json:"vote_count"`
	Popularity       float64               `json:"popularity"`
	Adult            bool                  `json:"adult"`
	PosterPath       *string               `json:"poster_path"`
	BackdropPath     *string               `json:"backdrop_path"`
	Homepage         *string               `json:"homepage"`
	Genres           []GenreResponse       `json:"genres"`
	Networks         []NetworkResponse     `json:"networks"`
	CreatedBy        []CreatorResponse     `json:"created_by"`
	OriginCountry    []string              `json:"origin_country"`
	SpokenLanguages  []LanguageResponse    `json:"spoken_languages"`
	Seasons          []SeasonSummaryResponse `json:"seasons"`

	// Appended data
	Credits        *CreditsResponse         `json:"credits,omitempty"`
	Images         *ImagesResponse          `json:"images,omitempty"`
	ContentRatings *ContentRatingsWrapper   `json:"content_ratings,omitempty"`
	Translations   *TranslationsWrapper     `json:"translations,omitempty"`
	ExternalIDs    *ExternalIDsResponse     `json:"external_ids,omitempty"`
	Videos         *VideosResponse          `json:"videos,omitempty"`
}

// SeasonResponse is the TMDb API response for season details.
type SeasonResponse struct {
	ID           int                      `json:"id"`
	Name         string                   `json:"name"`
	Overview     *string                  `json:"overview"`
	SeasonNumber int                      `json:"season_number"`
	AirDate      string                   `json:"air_date"`
	PosterPath   *string                  `json:"poster_path"`
	VoteAverage  float64                  `json:"vote_average"`
	Episodes     []EpisodeSummaryResponse `json:"episodes"`

	// Appended data
	Credits      *CreditsResponse     `json:"credits,omitempty"`
	Images       *ImagesResponse      `json:"images,omitempty"`
	Translations *TranslationsWrapper `json:"translations,omitempty"`
}

// SeasonSummaryResponse is a season summary in TV show details.
type SeasonSummaryResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     *string `json:"overview"`
	SeasonNumber int     `json:"season_number"`
	AirDate      string  `json:"air_date"`
	PosterPath   *string `json:"poster_path"`
	VoteAverage  float64 `json:"vote_average"`
	EpisodeCount int     `json:"episode_count"`
}

// EpisodeResponse is the TMDb API response for episode details.
type EpisodeResponse struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Overview       *string          `json:"overview"`
	SeasonNumber   int              `json:"season_number"`
	EpisodeNumber  int              `json:"episode_number"`
	AirDate        string           `json:"air_date"`
	Runtime        *int             `json:"runtime"`
	StillPath      *string          `json:"still_path"`
	VoteAverage    float64          `json:"vote_average"`
	VoteCount      int              `json:"vote_count"`
	ProductionCode *string          `json:"production_code"`
	GuestStars     []CastResponse   `json:"guest_stars"`
	Crew           []CrewResponse   `json:"crew"`

	// Appended data
	Credits      *CreditsResponse     `json:"credits,omitempty"`
	Images       *ImagesResponse      `json:"images,omitempty"`
	Translations *TranslationsWrapper `json:"translations,omitempty"`
}

// EpisodeSummaryResponse is an episode summary in season details.
type EpisodeSummaryResponse struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       *string `json:"overview"`
	SeasonNumber   int     `json:"season_number"`
	EpisodeNumber  int     `json:"episode_number"`
	AirDate        string  `json:"air_date"`
	Runtime        *int    `json:"runtime"`
	StillPath      *string `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
	ProductionCode *string `json:"production_code"`
}

// PersonSearchResultsResponse is the TMDb API response for person search.
type PersonSearchResultsResponse struct {
	Page         int                    `json:"page"`
	Results      []PersonSearchResponse `json:"results"`
	TotalPages   int                    `json:"total_pages"`
	TotalResults int                    `json:"total_results"`
}

// PersonSearchResponse is a single person search result.
type PersonSearchResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	ProfilePath *string              `json:"profile_path"`
	Popularity  float64              `json:"popularity"`
	Adult       bool                 `json:"adult"`
	KnownFor    []KnownForResponse   `json:"known_for"`
}

// KnownForResponse represents a media item a person is known for.
type KnownForResponse struct {
	MediaType    string  `json:"media_type"`
	ID           int     `json:"id"`
	Title        string  `json:"title"`         // For movies
	Name         string  `json:"name"`          // For TV
	PosterPath   *string `json:"poster_path"`
	ReleaseDate  string  `json:"release_date"`  // For movies
	FirstAirDate string  `json:"first_air_date"` // For TV
}

// PersonResponse is the TMDb API response for person details.
type PersonResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	AlsoKnownAs  []string  `json:"also_known_as"`
	Biography    *string   `json:"biography"`
	Birthday     *string   `json:"birthday"`
	Deathday     *string   `json:"deathday"`
	Gender       int       `json:"gender"`
	PlaceOfBirth *string   `json:"place_of_birth"`
	ProfilePath  *string   `json:"profile_path"`
	Homepage     *string   `json:"homepage"`
	Popularity   float64   `json:"popularity"`
	Adult        bool      `json:"adult"`
	IMDbID       *string   `json:"imdb_id"`
	KnownForDept string    `json:"known_for_department"`

	// Appended data
	ExternalIDs *ExternalIDsResponse `json:"external_ids,omitempty"`
	Images      *PersonImagesResponse `json:"images,omitempty"`
}

// PersonCreditsResponse is the TMDb API response for person credits.
type PersonCreditsResponse struct {
	ID   int                     `json:"id"`
	Cast []PersonCastCredit      `json:"cast"`
	Crew []PersonCrewCredit      `json:"crew"`
}

// PersonCastCredit represents a cast credit in a person's filmography.
type PersonCastCredit struct {
	MediaType    string  `json:"media_type"`
	ID           int     `json:"id"`
	Title        string  `json:"title"`          // For movies
	Name         string  `json:"name"`           // For TV
	Character    *string `json:"character"`
	PosterPath   *string `json:"poster_path"`
	ReleaseDate  string  `json:"release_date"`   // For movies
	FirstAirDate string  `json:"first_air_date"` // For TV
	VoteAverage  float64 `json:"vote_average"`
	EpisodeCount *int    `json:"episode_count"`  // For TV
}

// PersonCrewCredit represents a crew credit in a person's filmography.
type PersonCrewCredit struct {
	MediaType    string  `json:"media_type"`
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Name         string  `json:"name"`
	Job          *string `json:"job"`
	Department   *string `json:"department"`
	PosterPath   *string `json:"poster_path"`
	ReleaseDate  string  `json:"release_date"`
	FirstAirDate string  `json:"first_air_date"`
	VoteAverage  float64 `json:"vote_average"`
	EpisodeCount *int    `json:"episode_count"`
}

// CreditsResponse is the TMDb API response for credits.
type CreditsResponse struct {
	ID   int            `json:"id"`
	Cast []CastResponse `json:"cast"`
	Crew []CrewResponse `json:"crew"`
}

// CastResponse is a single cast member.
type CastResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Character   string  `json:"character"`
	Order       int     `json:"order"`
	CreditID    string  `json:"credit_id"`
	Gender      int     `json:"gender"`
	ProfilePath *string `json:"profile_path"`
}

// CrewResponse is a single crew member.
type CrewResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Job         string  `json:"job"`
	Department  string  `json:"department"`
	CreditID    string  `json:"credit_id"`
	Gender      int     `json:"gender"`
	ProfilePath *string `json:"profile_path"`
}

// ImagesResponse is the TMDb API response for images.
type ImagesResponse struct {
	ID        int             `json:"id"`
	Backdrops []ImageResponse `json:"backdrops"`
	Posters   []ImageResponse `json:"posters"`
	Logos     []ImageResponse `json:"logos"`
	Stills    []ImageResponse `json:"stills"`
}

// PersonImagesResponse is the TMDb API response for person images.
type PersonImagesResponse struct {
	ID       int             `json:"id"`
	Profiles []ImageResponse `json:"profiles"`
}

// ImageResponse is a single image.
type ImageResponse struct {
	AspectRatio float64 `json:"aspect_ratio"`
	FilePath    string  `json:"file_path"`
	Height      int     `json:"height"`
	Width       int     `json:"width"`
	ISO639_1    *string `json:"iso_639_1"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// GenreResponse is a genre from TMDb.
type GenreResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CompanyResponse is a production company from TMDb.
type CompanyResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	LogoPath      *string `json:"logo_path"`
	OriginCountry string  `json:"origin_country"`
}

// CountryResponse is a production country from TMDb.
type CountryResponse struct {
	ISO3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}

// LanguageResponse is a spoken language from TMDb.
type LanguageResponse struct {
	ISO639_1    string `json:"iso_639_1"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

// NetworkResponse is a TV network from TMDb.
type NetworkResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	LogoPath      *string `json:"logo_path"`
	OriginCountry string  `json:"origin_country"`
}

// CreatorResponse is a TV show creator from TMDb.
type CreatorResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Gender      int     `json:"gender"`
	ProfilePath *string `json:"profile_path"`
	CreditID    string  `json:"credit_id"`
}

// CollectionRefResponse is a collection reference from TMDb.
type CollectionRefResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	PosterPath   *string `json:"poster_path"`
	BackdropPath *string `json:"backdrop_path"`
}

// CollectionResponse is the TMDb API response for collection details.
type CollectionResponse struct {
	ID           int                   `json:"id"`
	Name         string                `json:"name"`
	Overview     *string               `json:"overview"`
	PosterPath   *string               `json:"poster_path"`
	BackdropPath *string               `json:"backdrop_path"`
	Parts        []MovieSearchResponse `json:"parts"`
}

// ReleaseDatesWrapper wraps the release dates response.
type ReleaseDatesWrapper struct {
	ID      int                    `json:"id"`
	Results []CountryReleaseResponse `json:"results"`
}

// CountryReleaseResponse contains releases for a specific country.
type CountryReleaseResponse struct {
	ISO3166_1    string                `json:"iso_3166_1"`
	ReleaseDates []ReleaseDateResponse `json:"release_dates"`
}

// ReleaseDateResponse is a single release date entry.
type ReleaseDateResponse struct {
	Certification string `json:"certification"`
	ISO639_1      string `json:"iso_639_1"`
	ReleaseDate   string `json:"release_date"`
	Type          int    `json:"type"`
	Note          string `json:"note"`
}

// ContentRatingsWrapper wraps the content ratings response.
type ContentRatingsWrapper struct {
	ID      int                     `json:"id"`
	Results []ContentRatingResponse `json:"results"`
}

// ContentRatingResponse is a content rating for a specific country.
type ContentRatingResponse struct {
	ISO3166_1   string   `json:"iso_3166_1"`
	Rating      string   `json:"rating"`
	Descriptors []string `json:"descriptors"`
}

// TranslationsWrapper wraps the translations response.
type TranslationsWrapper struct {
	ID           int                   `json:"id"`
	Translations []TranslationResponse `json:"translations"`
}

// TranslationResponse is a single translation entry.
type TranslationResponse struct {
	ISO3166_1   string              `json:"iso_3166_1"`
	ISO639_1    string              `json:"iso_639_1"`
	Name        string              `json:"name"`
	EnglishName string              `json:"english_name"`
	Data        TranslationDataResponse `json:"data"`
}

// TranslationDataResponse contains translated content.
type TranslationDataResponse struct {
	Title    string `json:"title"`    // For movies
	Name     string `json:"name"`     // For TV
	Overview string `json:"overview"`
	Tagline  string `json:"tagline"`
	Homepage string `json:"homepage"`
	Runtime  *int   `json:"runtime"`
}

// ExternalIDsResponse is the TMDb API response for external IDs.
type ExternalIDsResponse struct {
	ID           int     `json:"id"`
	IMDbID       *string `json:"imdb_id"`
	TVDbID       *int    `json:"tvdb_id"`
	WikidataID   *string `json:"wikidata_id"`
	FacebookID   *string `json:"facebook_id"`
	InstagramID  *string `json:"instagram_id"`
	TwitterID    *string `json:"twitter_id"`
	TikTokID     *string `json:"tiktok_id"`
	YouTubeID    *string `json:"youtube_id"`
	FreebaseID   *string `json:"freebase_id"`
	FreebaseMID  *string `json:"freebase_mid"`
	TVRageID     *int    `json:"tvrage_id"`
}

// VideosResponse is the TMDb API response for videos.
type VideosResponse struct {
	ID      int             `json:"id"`
	Results []VideoResponse `json:"results"`
}

// VideoResponse is a single video entry.
type VideoResponse struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Site      string `json:"site"`
	Type      string `json:"type"`
	Official  bool   `json:"official"`
	Published string `json:"published_at"`
}

// ErrorResponse is the TMDb API error response.
type ErrorResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}

// CacheEntry wraps cached data with expiration.
type CacheEntry struct {
	Data      any
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired.
func (c *CacheEntry) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
