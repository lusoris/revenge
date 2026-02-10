package movie

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/lusoris/revenge/internal/content"
)

// Movie represents a movie with metadata from TMDb/Radarr
type Movie struct {
	ID                uuid.UUID
	TMDbID            *int32
	IMDbID            *string
	Title             string
	OriginalTitle     *string
	Year              *int32
	ReleaseDate       *time.Time
	Runtime           *int32
	Overview          *string
	Tagline           *string
	Status            *string
	OriginalLanguage  *string
	PosterPath        *string
	BackdropPath      *string
	TrailerURL        *string
	VoteAverage       *decimal.Decimal
	VoteCount         *int32
	Popularity        *decimal.Decimal
	Budget            *int64
	Revenue           *int64
	LibraryAddedAt    time.Time
	MetadataUpdatedAt *time.Time
	RadarrID          *int32
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Multi-language support (Hybrid JSONB approach)
	// Default fields above are kept for simple queries and backwards compatibility
	TitlesI18n    map[string]string            // {"en": "The Shawshank Redemption", "de": "Die Verurteilten", "fr": "Les Évadés"}
	TaglinesI18n  map[string]string            // {"en": "Fear can hold you prisoner...", "de": "Angst kann dich gefangen halten..."}
	OverviewsI18n map[string]string            // {"en": "Imprisoned in the 1940s...", "de": "In den 1940er Jahren eingesperrt..."}
	AgeRatings    map[string]map[string]string // {"US": {"MPAA": "R"}, "DE": {"FSK": "12"}, "GB": {"BBFC": "15"}}

	// External ratings from various providers (IMDb, Rotten Tomatoes, Metacritic, TMDb, etc.)
	ExternalRatings []ExternalRating
}

// ExternalRating is an alias for the shared content.ExternalRating type.
type ExternalRating = content.ExternalRating

// GetTitle returns the movie title in the preferred language with fallback chain:
// 1. Requested language from TitlesI18n
// 2. English from TitlesI18n
// 3. Original language from TitlesI18n
// 4. OriginalTitle field
// 5. Default Title field
func (m *Movie) GetTitle(lang string) string {
	// Try requested language
	if m.TitlesI18n != nil {
		if title, ok := m.TitlesI18n[lang]; ok && title != "" {
			return title
		}
	}

	// Fallback to English
	if m.TitlesI18n != nil {
		if title, ok := m.TitlesI18n["en"]; ok && title != "" {
			return title
		}
	}

	// Fallback to original language
	if m.OriginalLanguage != nil && *m.OriginalLanguage != "" && m.TitlesI18n != nil {
		if title, ok := m.TitlesI18n[*m.OriginalLanguage]; ok && title != "" {
			return title
		}
	}

	// Fallback to OriginalTitle
	if m.OriginalTitle != nil && *m.OriginalTitle != "" {
		return *m.OriginalTitle
	}

	// Final fallback to default field
	return m.Title
}

// GetTagline returns the movie tagline in the preferred language with fallback chain:
// 1. Requested language from TaglinesI18n
// 2. English from TaglinesI18n
// 3. Original language from TaglinesI18n
// 4. Default Tagline field
func (m *Movie) GetTagline(lang string) string {
	// Try requested language
	if m.TaglinesI18n != nil {
		if tagline, ok := m.TaglinesI18n[lang]; ok && tagline != "" {
			return tagline
		}
	}

	// Fallback to English
	if m.TaglinesI18n != nil {
		if tagline, ok := m.TaglinesI18n["en"]; ok && tagline != "" {
			return tagline
		}
	}

	// Fallback to original language
	if m.OriginalLanguage != nil && *m.OriginalLanguage != "" && m.TaglinesI18n != nil {
		if tagline, ok := m.TaglinesI18n[*m.OriginalLanguage]; ok && tagline != "" {
			return tagline
		}
	}

	// Final fallback to default field
	if m.Tagline != nil {
		return *m.Tagline
	}

	return ""
}

// GetOverview returns the movie overview in the preferred language with fallback chain:
// 1. Requested language from OverviewsI18n
// 2. English from OverviewsI18n
// 3. Original language from OverviewsI18n
// 4. Default Overview field
func (m *Movie) GetOverview(lang string) string {
	// Try requested language
	if m.OverviewsI18n != nil {
		if overview, ok := m.OverviewsI18n[lang]; ok && overview != "" {
			return overview
		}
	}

	// Fallback to English
	if m.OverviewsI18n != nil {
		if overview, ok := m.OverviewsI18n["en"]; ok && overview != "" {
			return overview
		}
	}

	// Fallback to original language
	if m.OriginalLanguage != nil && *m.OriginalLanguage != "" && m.OverviewsI18n != nil {
		if overview, ok := m.OverviewsI18n[*m.OriginalLanguage]; ok && overview != "" {
			return overview
		}
	}

	// Final fallback to default field
	if m.Overview != nil {
		return *m.Overview
	}

	return ""
}

// GetAgeRating returns the age rating for a specific country and rating system.
// Returns empty string if not found.
//
// Examples:
//   - GetAgeRating("US", "MPAA") -> "R"
//   - GetAgeRating("DE", "FSK") -> "12"
//   - GetAgeRating("GB", "BBFC") -> "15"
func (m *Movie) GetAgeRating(country, system string) string {
	if m.AgeRatings == nil {
		return ""
	}

	countryRatings, ok := m.AgeRatings[country]
	if !ok {
		return ""
	}

	rating, ok := countryRatings[system]
	if !ok {
		return ""
	}

	return rating
}

// GetAvailableLanguages returns a list of all available language codes
// that have translations for this movie.
func (m *Movie) GetAvailableLanguages() []string {
	if m.TitlesI18n == nil {
		return []string{}
	}

	langs := make([]string, 0, len(m.TitlesI18n))
	for lang := range m.TitlesI18n {
		langs = append(langs, lang)
	}

	return langs
}

// GetAvailableAgeRatingCountries returns a list of all countries
// that have age ratings for this movie.
func (m *Movie) GetAvailableAgeRatingCountries() []string {
	if m.AgeRatings == nil {
		return []string{}
	}

	countries := make([]string, 0, len(m.AgeRatings))
	for country := range m.AgeRatings {
		countries = append(countries, country)
	}

	return countries
}

// MovieFile represents a physical media file for a movie
type MovieFile struct {
	ID                uuid.UUID
	MovieID           uuid.UUID
	FilePath          string
	FileSize          int64
	FileName          string
	Resolution        *string
	QualityProfile    *string
	VideoCodec        *string
	AudioCodec        *string
	Container         *string
	DurationSeconds   *int32
	BitrateKbps       *int32
	Framerate         *decimal.Decimal
	DynamicRange      *string
	ColorSpace        *string
	AudioChannels     *string
	AudioLanguages    []string
	SubtitleLanguages []string
	RadarrFileID      *int32
	LastScannedAt     *time.Time
	IsMonitored       *bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// MovieCredit represents cast or crew member for a movie
type MovieCredit struct {
	ID           uuid.UUID
	MovieID      uuid.UUID
	TMDbPersonID int32
	Name         string
	CreditType   string // 'cast' or 'crew'
	Character    *string
	Job          *string
	Department   *string
	CastOrder    *int32
	ProfilePath  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MovieCollection represents a collection of related movies
type MovieCollection struct {
	ID               uuid.UUID
	TMDbCollectionID *int32
	Name             string
	Overview         *string
	PosterPath       *string
	BackdropPath     *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// MovieGenre represents a genre associated with a movie
type MovieGenre struct {
	ID          uuid.UUID
	MovieID     uuid.UUID
	TMDbGenreID int32
	Name        string
	CreatedAt   time.Time
}

// MovieWatched represents watch progress for a user
type MovieWatched struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	MovieID         uuid.UUID
	ProgressSeconds int32
	DurationSeconds int32
	ProgressPercent *int32 // Generated column
	IsCompleted     bool
	WatchCount      int32
	LastWatchedAt   time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ListFilters contains filters for listing movies
type ListFilters struct {
	OrderBy string // "title", "year", "added", "rating"
	Limit   int32
	Offset  int32
}

// SearchFilters contains filters for searching movies
type SearchFilters struct {
	Limit  int32
	Offset int32
}
