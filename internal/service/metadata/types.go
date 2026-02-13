package metadata

import (
	"time"
)

// MovieSearchResult represents a movie search result from any provider.
type MovieSearchResult struct {
	// ProviderID is the ID from the source provider.
	ProviderID string
	// Provider identifies which provider returned this result.
	Provider ProviderID

	Title            string
	OriginalTitle    string
	OriginalLanguage string
	Overview         string
	ReleaseDate      *time.Time
	Year             *int
	PosterPath       *string
	BackdropPath     *string
	VoteAverage      float64
	VoteCount        int
	Popularity       float64
	Adult            bool
	GenreIDs         []int
}

// MovieMetadata contains full movie details from a provider.
type MovieMetadata struct {
	// Provider identifiers
	ProviderID string
	Provider   ProviderID
	IMDbID     *string
	TMDbID     *int32
	TVDbID     *int32

	// Basic info
	Title            string
	OriginalTitle    string
	OriginalLanguage string
	Tagline          *string
	Overview         *string
	Status           string

	// Dates and runtime
	ReleaseDate *time.Time
	Runtime     *int32

	// Financial
	Budget  *int64
	Revenue *int64

	// Ratings
	VoteAverage     float64
	VoteCount       int
	Popularity      float64
	Adult           bool
	ExternalRatings []ExternalRating // IMDb, Rotten Tomatoes, Metacritic (from OMDb)

	// Images
	PosterPath   *string
	BackdropPath *string
	LogoPath     *string
	Homepage     *string
	TrailerURL   *string

	// Related data
	Genres              []Genre
	ProductionCompanies []ProductionCompany
	ProductionCountries []ProductionCountry
	SpokenLanguages     []SpokenLanguage
	Collection          *CollectionRef

	// Localized data (populated when fetching multi-language)
	Translations map[string]*LocalizedMovieData
}

// LocalizedMovieData contains language-specific movie fields.
type LocalizedMovieData struct {
	Language string
	Title    string
	Overview string
	Tagline  string
	Homepage string
	Runtime  *int32
}

// TVShowSearchResult represents a TV show search result.
type TVShowSearchResult struct {
	ProviderID string
	Provider   ProviderID

	Name             string
	OriginalName     string
	OriginalLanguage string
	Overview         string
	FirstAirDate     *time.Time
	Year             *int
	PosterPath       *string
	BackdropPath     *string
	VoteAverage      float64
	VoteCount        int
	Popularity       float64
	Adult            bool
	GenreIDs         []int
	OriginCountries  []string
}

// TVShowMetadata contains full TV show details.
type TVShowMetadata struct {
	// Provider identifiers
	ProviderID string
	Provider   ProviderID
	IMDbID     *string
	TMDbID     *int32
	TVDbID     *int32

	// Basic info
	Name             string
	OriginalName     string
	OriginalLanguage string
	Tagline          *string
	Overview         *string
	Status           string
	Type             string // Scripted, Reality, Documentary, etc.

	// Dates
	FirstAirDate *time.Time
	LastAirDate  *time.Time
	InProduction bool

	// Counts
	NumberOfSeasons  int
	NumberOfEpisodes int
	EpisodeRuntime   []int

	// Ratings
	VoteAverage     float64
	VoteCount       int
	Popularity      float64
	Adult           bool
	ExternalRatings []ExternalRating // IMDb, Rotten Tomatoes, Metacritic (from OMDb)

	// Images
	PosterPath   *string
	BackdropPath *string
	Homepage     *string
	TrailerURL   *string

	// Related data
	Genres          []Genre
	Networks        []Network
	CreatedBy       []Creator
	OriginCountries []string
	SpokenLanguages []SpokenLanguage

	// Seasons (populated if requested)
	Seasons []SeasonSummary

	// Localized data
	Translations map[string]*LocalizedTVShowData
}

// LocalizedTVShowData contains language-specific TV show fields.
type LocalizedTVShowData struct {
	Language string
	Name     string
	Overview string
	Tagline  string
	Homepage string
}

// SeasonSummary contains basic season info (as part of TVShowMetadata).
type SeasonSummary struct {
	ProviderID   string
	SeasonNumber int
	Name         string
	Overview     *string
	PosterPath   *string
	AirDate      *time.Time
	EpisodeCount int
	VoteAverage  float64
}

// SeasonMetadata contains full season details.
type SeasonMetadata struct {
	ProviderID string
	Provider   ProviderID
	TMDbID     *int32
	TVDbID     *int32

	ShowID       string // Provider ID of the parent show
	SeasonNumber int
	Name         string
	Overview     *string
	PosterPath   *string
	AirDate      *time.Time
	VoteAverage  float64

	Episodes []EpisodeSummary

	// Localized data
	Translations map[string]*LocalizedSeasonData
}

// LocalizedSeasonData contains language-specific season fields.
type LocalizedSeasonData struct {
	Language string
	Name     string
	Overview string
}

// EpisodeSummary contains basic episode info (as part of SeasonMetadata).
type EpisodeSummary struct {
	ProviderID     string
	EpisodeNumber  int
	Name           string
	Overview       *string
	AirDate        *time.Time
	Runtime        *int32
	StillPath      *string
	VoteAverage    float64
	VoteCount      int
	ProductionCode *string
}

// EpisodeMetadata contains full episode details.
type EpisodeMetadata struct {
	ProviderID string
	Provider   ProviderID
	TMDbID     *int32
	TVDbID     *int32
	IMDbID     *string

	ShowID         string // Provider ID of the parent show
	SeasonNumber   int
	EpisodeNumber  int
	Name           string
	Overview       *string
	AirDate        *time.Time
	Runtime        *int32
	StillPath      *string
	VoteAverage    float64
	VoteCount      int
	ProductionCode *string

	// Guest stars and crew
	GuestStars []CastMember
	Crew       []CrewMember

	// Localized data
	Translations map[string]*LocalizedEpisodeData
}

// LocalizedEpisodeData contains language-specific episode fields.
type LocalizedEpisodeData struct {
	Language string
	Name     string
	Overview string
}

// PersonSearchResult represents a person search result.
type PersonSearchResult struct {
	ProviderID string
	Provider   ProviderID

	Name        string
	ProfilePath *string
	Popularity  float64
	Adult       bool
	KnownFor    []MediaReference // Movies/TV shows they're known for
}

// PersonMetadata contains full person details.
type PersonMetadata struct {
	ProviderID string
	Provider   ProviderID
	IMDbID     *string
	TMDbID     *int32

	Name         string
	AlsoKnownAs  []string
	Biography    *string
	Birthday     *time.Time
	Deathday     *time.Time
	Gender       int // 0=not specified, 1=female, 2=male, 3=non-binary
	PlaceOfBirth *string
	ProfilePath  *string
	Homepage     *string
	Popularity   float64
	Adult        bool
	KnownForDept string // Acting, Directing, etc.

	// Localized data
	Translations map[string]*LocalizedPersonData
}

// LocalizedPersonData contains language-specific person fields.
type LocalizedPersonData struct {
	Language  string
	Biography string
}

// PersonCredits contains a person's filmography.
type PersonCredits struct {
	ProviderID string
	Provider   ProviderID

	CastCredits []MediaCredit // Roles as actor
	CrewCredits []MediaCredit // Roles as crew
}

// MediaCredit represents a credit in a person's filmography.
type MediaCredit struct {
	MediaType    string // "movie" or "tv"
	MediaID      string // Provider ID
	Title        string
	Character    *string // For cast
	Job          *string // For crew
	Department   *string // For crew
	ReleaseDate  *time.Time
	PosterPath   *string
	VoteAverage  float64
	EpisodeCount *int // For TV credits
}

// MediaReference is a lightweight reference to a movie or TV show.
type MediaReference struct {
	MediaType  string
	ID         string
	Title      string
	PosterPath *string
}

// CollectionRef is a lightweight reference to a collection.
type CollectionRef struct {
	ID           int
	Name         string
	PosterPath   *string
	BackdropPath *string
}

// CollectionMetadata contains full collection details.
type CollectionMetadata struct {
	ProviderID string
	Provider   ProviderID

	Name         string
	Overview     *string
	PosterPath   *string
	BackdropPath *string
	Parts        []MovieSearchResult
}

// Credits contains cast and crew information.
type Credits struct {
	Cast []CastMember
	Crew []CrewMember
}

// CastMember represents an actor/performer.
type CastMember struct {
	ProviderID  string
	Name        string
	Character   string
	Order       int
	CreditID    *string
	Gender      int
	ProfilePath *string
}

// CrewMember represents a crew member.
type CrewMember struct {
	ProviderID  string
	Name        string
	Job         string
	Department  string
	CreditID    *string
	Gender      int
	ProfilePath *string
}

// Images contains categorized images.
type Images struct {
	Posters   []Image
	Backdrops []Image
	Logos     []Image
	Stills    []Image // For episodes
	Profiles  []Image // For people
}

// Image represents a single image.
type Image struct {
	FilePath    string
	AspectRatio float64
	Width       int
	Height      int
	VoteAverage float64
	VoteCount   int
	Language    *string // ISO 639-1
}

// Genre represents a content genre.
type Genre struct {
	ID   int
	Name string
}

// ProductionCompany represents a production company.
type ProductionCompany struct {
	ID            int
	Name          string
	LogoPath      *string
	OriginCountry string
}

// ProductionCountry represents a production country.
type ProductionCountry struct {
	ISOCode string // ISO 3166-1 alpha-2
	Name    string
}

// SpokenLanguage represents a spoken language.
type SpokenLanguage struct {
	ISOCode     string // ISO 639-1
	Name        string
	EnglishName string
}

// Network represents a TV network.
type Network struct {
	ID            int
	Name          string
	LogoPath      *string
	OriginCountry string
}

// Creator represents a TV show creator.
type Creator struct {
	ID          int
	Name        string
	Gender      int
	ProfilePath *string
	CreditID    *string
}

// ReleaseDate contains release information for a specific country.
type ReleaseDate struct {
	CountryCode   string // ISO 3166-1 alpha-2
	Certification string // Age rating (R, PG-13, FSK 12, etc.)
	ReleaseDate   *time.Time
	ReleaseType   int    // 1=Premiere, 2=Theatrical (limited), 3=Theatrical, 4=Digital, 5=Physical, 6=TV
	Language      string // ISO 639-1
	Note          string
}

// ContentRating contains a content rating for a specific country (TV shows).
type ContentRating struct {
	CountryCode string   // ISO 3166-1 alpha-2
	Rating      string   // Content rating
	Descriptors []string // Content descriptors (violence, language, etc.)
}

// Translation contains translation availability info.
type Translation struct {
	ISOCode     string // ISO 3166-1 alpha-2
	Language    string // ISO 639-1
	Name        string // Localized country name
	EnglishName string
	Data        *TranslationData
}

// TranslationData contains translated content.
type TranslationData struct {
	Title    string
	Overview string
	Tagline  string
	Homepage string
	Runtime  *int32
}

// ExternalRating represents a rating from an external source (IMDb, Rotten Tomatoes, etc.).
type ExternalRating struct {
	Source string  // "Internet Movie Database", "Rotten Tomatoes", "Metacritic"
	Value  string  // "8.8/10", "96%", "90/100"
	Score  float64 // Normalized 0-100 scale
}

// ExternalIDs contains external identifiers from various sources.
type ExternalIDs struct {
	IMDbID      *string
	TVDbID      *int32
	TMDbID      *int32
	WikidataID  *string
	FacebookID  *string
	InstagramID *string
	TwitterID   *string
	TikTokID    *string
	YouTubeID   *string
	FreebaseID  *string
	FreebaseMID *string
	TVRageID    *int32
}
