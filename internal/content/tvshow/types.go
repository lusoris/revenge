// Package tvshow provides types, adapters, and services for TV show content management.
// It handles series, seasons, episodes, and episode files with full localization support.
package tvshow

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/lusoris/revenge/internal/content"
)

// Series represents a TV show with metadata from TMDb/TVDb/Sonarr.
type Series struct {
	ID               uuid.UUID
	TMDbID           *int32
	TVDbID           *int32
	IMDbID           *string
	SonarrID         *int32
	Title            string
	OriginalTitle    *string
	OriginalLanguage string
	Tagline          *string
	Overview         *string
	Status           *string            // "Returning Series", "Ended", "Canceled", etc.
	Type             *string            // "Scripted", "Reality", "Documentary", etc.
	FirstAirDate     *time.Time
	LastAirDate      *time.Time
	VoteAverage      *decimal.Decimal
	VoteCount        *int32
	Popularity       *decimal.Decimal
	PosterPath       *string
	BackdropPath     *string
	TotalSeasons     int32
	TotalEpisodes    int32
	TrailerURL       *string
	Homepage         *string
	MetadataUpdatedAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Multi-language support (Hybrid JSONB approach)
	TitlesI18n    map[string]string            // {"en": "Breaking Bad", "de": "Breaking Bad", "es": "Breaking Bad"}
	TaglinesI18n  map[string]string            // {"en": "All Hail the King", "de": "Heil dem König"}
	OverviewsI18n map[string]string            // {"en": "A chemistry teacher...", "de": "Ein Chemielehrer..."}
	AgeRatings    map[string]map[string]string // {"US": {"TV": "TV-MA"}, "DE": {"FSK": "16"}}

	// External ratings from various providers (IMDb, Rotten Tomatoes, Metacritic, TMDb, etc.)
	ExternalRatings []ExternalRating

	// Relations (populated by repository when requested)
	Seasons []Season
	Genres  []SeriesGenre
	Credits []SeriesCredit
	Networks []Network
}

// ExternalRating is an alias for the shared content.ExternalRating type.
type ExternalRating = content.ExternalRating

// GetTitle returns the series title in the preferred language with fallback chain:
// 1. Requested language from TitlesI18n
// 2. English from TitlesI18n
// 3. Original language from TitlesI18n
// 4. OriginalTitle field
// 5. Default Title field
func (s *Series) GetTitle(lang string) string {
	// Try requested language
	if s.TitlesI18n != nil {
		if title, ok := s.TitlesI18n[lang]; ok && title != "" {
			return title
		}
	}

	// Fallback to English
	if s.TitlesI18n != nil {
		if title, ok := s.TitlesI18n["en"]; ok && title != "" {
			return title
		}
	}

	// Fallback to original language
	if s.OriginalLanguage != "" && s.TitlesI18n != nil {
		if title, ok := s.TitlesI18n[s.OriginalLanguage]; ok && title != "" {
			return title
		}
	}

	// Fallback to OriginalTitle
	if s.OriginalTitle != nil && *s.OriginalTitle != "" {
		return *s.OriginalTitle
	}

	// Final fallback to default field
	return s.Title
}

// GetTagline returns the series tagline in the preferred language with fallback chain:
// 1. Requested language from TaglinesI18n
// 2. English from TaglinesI18n
// 3. Original language from TaglinesI18n
// 4. Default Tagline field
func (s *Series) GetTagline(lang string) string {
	// Try requested language
	if s.TaglinesI18n != nil {
		if tagline, ok := s.TaglinesI18n[lang]; ok && tagline != "" {
			return tagline
		}
	}

	// Fallback to English
	if s.TaglinesI18n != nil {
		if tagline, ok := s.TaglinesI18n["en"]; ok && tagline != "" {
			return tagline
		}
	}

	// Fallback to original language
	if s.OriginalLanguage != "" && s.TaglinesI18n != nil {
		if tagline, ok := s.TaglinesI18n[s.OriginalLanguage]; ok && tagline != "" {
			return tagline
		}
	}

	// Final fallback to default field
	if s.Tagline != nil {
		return *s.Tagline
	}

	return ""
}

// GetOverview returns the series overview in the preferred language with fallback chain:
// 1. Requested language from OverviewsI18n
// 2. English from OverviewsI18n
// 3. Original language from OverviewsI18n
// 4. Default Overview field
func (s *Series) GetOverview(lang string) string {
	// Try requested language
	if s.OverviewsI18n != nil {
		if overview, ok := s.OverviewsI18n[lang]; ok && overview != "" {
			return overview
		}
	}

	// Fallback to English
	if s.OverviewsI18n != nil {
		if overview, ok := s.OverviewsI18n["en"]; ok && overview != "" {
			return overview
		}
	}

	// Fallback to original language
	if s.OriginalLanguage != "" && s.OverviewsI18n != nil {
		if overview, ok := s.OverviewsI18n[s.OriginalLanguage]; ok && overview != "" {
			return overview
		}
	}

	// Final fallback to default field
	if s.Overview != nil {
		return *s.Overview
	}

	return ""
}

// GetAgeRating returns the age rating for a specific country and rating system.
// Returns empty string if not found.
//
// Examples:
//   - GetAgeRating("US", "TV") -> "TV-MA"
//   - GetAgeRating("DE", "FSK") -> "16"
//   - GetAgeRating("GB", "BBFC") -> "15"
func (s *Series) GetAgeRating(country, system string) string {
	if s.AgeRatings == nil {
		return ""
	}

	countryRatings, ok := s.AgeRatings[country]
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
// that have translations for this series.
func (s *Series) GetAvailableLanguages() []string {
	if s.TitlesI18n == nil {
		return []string{}
	}

	langs := make([]string, 0, len(s.TitlesI18n))
	for lang := range s.TitlesI18n {
		langs = append(langs, lang)
	}

	return langs
}

// GetAvailableAgeRatingCountries returns a list of all countries
// that have age ratings for this series.
func (s *Series) GetAvailableAgeRatingCountries() []string {
	if s.AgeRatings == nil {
		return []string{}
	}

	countries := make([]string, 0, len(s.AgeRatings))
	for country := range s.AgeRatings {
		countries = append(countries, country)
	}

	return countries
}

// IsEnded returns true if the series has ended (no more episodes coming).
func (s *Series) IsEnded() bool {
	if s.Status == nil {
		return false
	}
	status := *s.Status
	return status == "Ended" || status == "Canceled"
}

// Season represents a season within a TV series.
type Season struct {
	ID           uuid.UUID
	SeriesID     uuid.UUID
	TMDbID       *int32
	SeasonNumber int32
	Name         string
	Overview     *string
	PosterPath   *string
	EpisodeCount int32
	AirDate      *time.Time
	VoteAverage  *decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Multi-language support
	NamesI18n     map[string]string // {"en": "Season 1", "de": "Staffel 1"}
	OverviewsI18n map[string]string

	// Relations (populated by repository when requested)
	Episodes []Episode
}

// GetName returns the season name in the preferred language.
func (s *Season) GetName(lang string) string {
	if s.NamesI18n != nil {
		if name, ok := s.NamesI18n[lang]; ok && name != "" {
			return name
		}
		if name, ok := s.NamesI18n["en"]; ok && name != "" {
			return name
		}
	}
	return s.Name
}

// GetOverview returns the season overview in the preferred language.
func (s *Season) GetOverview(lang string) string {
	if s.OverviewsI18n != nil {
		if overview, ok := s.OverviewsI18n[lang]; ok && overview != "" {
			return overview
		}
		if overview, ok := s.OverviewsI18n["en"]; ok && overview != "" {
			return overview
		}
	}
	if s.Overview != nil {
		return *s.Overview
	}
	return ""
}

// IsSpecials returns true if this is a specials season (season 0).
func (s *Season) IsSpecials() bool {
	return s.SeasonNumber == 0
}

// Episode represents an episode within a season.
type Episode struct {
	ID             uuid.UUID
	SeriesID       uuid.UUID
	SeasonID       uuid.UUID
	TMDbID         *int32
	TVDbID         *int32
	IMDbID         *string
	SeasonNumber   int32
	EpisodeNumber  int32
	Title          string
	Overview       *string
	AirDate        *time.Time
	Runtime        *int32 // minutes
	VoteAverage    *decimal.Decimal
	VoteCount      *int32
	StillPath      *string
	ProductionCode *string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Multi-language support
	TitlesI18n    map[string]string
	OverviewsI18n map[string]string

	// Relations (populated by repository when requested)
	Files   []EpisodeFile
	Credits []EpisodeCredit
}

// GetTitle returns the episode title in the preferred language.
func (e *Episode) GetTitle(lang string) string {
	if e.TitlesI18n != nil {
		if title, ok := e.TitlesI18n[lang]; ok && title != "" {
			return title
		}
		if title, ok := e.TitlesI18n["en"]; ok && title != "" {
			return title
		}
	}
	return e.Title
}

// GetOverview returns the episode overview in the preferred language.
func (e *Episode) GetOverview(lang string) string {
	if e.OverviewsI18n != nil {
		if overview, ok := e.OverviewsI18n[lang]; ok && overview != "" {
			return overview
		}
		if overview, ok := e.OverviewsI18n["en"]; ok && overview != "" {
			return overview
		}
	}
	if e.Overview != nil {
		return *e.Overview
	}
	return ""
}

// EpisodeCode returns the episode code in SxxExx format.
func (e *Episode) EpisodeCode() string {
	return formatEpisodeCode(e.SeasonNumber, e.EpisodeNumber)
}

// HasAired returns true if the episode has already aired.
func (e *Episode) HasAired() bool {
	if e.AirDate == nil {
		return false
	}
	return e.AirDate.Before(time.Now())
}

// EpisodeFile represents a physical media file for an episode.
type EpisodeFile struct {
	ID                uuid.UUID
	EpisodeID         uuid.UUID
	FilePath          string
	FileName          string
	FileSize          int64
	Container         *string // mkv, mp4, avi, etc.
	Resolution        *string // 1920x1080
	QualityProfile    *string // HDTV-1080p, WEB-DL, Bluray
	VideoCodec        *string // h264, hevc
	AudioCodec        *string // aac, dts, truehd
	BitrateKbps       *int32
	DurationSeconds   *decimal.Decimal
	AudioLanguages    []string
	SubtitleLanguages []string
	SonarrFileID      *int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// SeriesCredit represents a cast or crew member for a series.
type SeriesCredit struct {
	ID           uuid.UUID
	SeriesID     uuid.UUID
	TMDbPersonID int32
	Name         string
	CreditType   string // "cast" or "crew"
	Character    *string
	Job          *string
	Department   *string
	CastOrder    *int32
	ProfilePath  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsCast returns true if this is a cast credit.
func (c *SeriesCredit) IsCast() bool {
	return c.CreditType == "cast"
}

// IsCrew returns true if this is a crew credit.
func (c *SeriesCredit) IsCrew() bool {
	return c.CreditType == "crew"
}

// EpisodeCredit represents a cast or crew member for a specific episode.
type EpisodeCredit struct {
	ID           uuid.UUID
	EpisodeID    uuid.UUID
	TMDbPersonID int32
	Name         string
	CreditType   string // "cast" or "crew"
	Character    *string
	Job          *string
	Department   *string
	CastOrder    *int32
	ProfilePath  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsCast returns true if this is a cast credit.
func (c *EpisodeCredit) IsCast() bool {
	return c.CreditType == "cast"
}

// IsCrew returns true if this is a crew credit.
func (c *EpisodeCredit) IsCrew() bool {
	return c.CreditType == "crew"
}

// SeriesGenre represents a genre associated with a series.
type SeriesGenre struct {
	ID        uuid.UUID
	SeriesID  uuid.UUID
	Slug      string
	Name      string
	CreatedAt time.Time
}

// Network represents a TV network or streaming service.
type Network struct {
	ID            uuid.UUID
	TMDbID        int32
	Name          string
	LogoPath      *string
	OriginCountry *string
	CreatedAt     time.Time
}

// EpisodeWatched represents watch progress for a user on an episode.
type EpisodeWatched struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	EpisodeID       uuid.UUID
	ProgressSeconds int32
	DurationSeconds int32
	IsCompleted     bool
	WatchCount      int32
	LastWatchedAt   time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ProgressPercent calculates the watch progress percentage.
func (w *EpisodeWatched) ProgressPercent() float64 {
	if w.DurationSeconds == 0 {
		return 0
	}
	return float64(w.ProgressSeconds) / float64(w.DurationSeconds) * 100
}

// SeriesWatchStats contains aggregated watch statistics for a series.
type SeriesWatchStats struct {
	WatchedCount    int64
	InProgressCount int64
	TotalWatches    int64
	TotalEpisodes   int64
}

// CompletionPercent calculates the series completion percentage.
func (s *SeriesWatchStats) CompletionPercent() float64 {
	if s.TotalEpisodes == 0 {
		return 0
	}
	return float64(s.WatchedCount) / float64(s.TotalEpisodes) * 100
}

// UserTVStats contains aggregated TV watching statistics for a user.
type UserTVStats struct {
	SeriesCount       int64
	EpisodesWatched   int64
	EpisodesInProgress int64
	TotalWatches      int64
}

// SeriesListFilters contains filters for listing series.
type SeriesListFilters struct {
	OrderBy string // "title", "first_air_date", "added", "rating", "popularity"
	Limit   int32
	Offset  int32
}

// EpisodeListFilters contains filters for listing episodes.
type EpisodeListFilters struct {
	SeriesID *uuid.UUID
	SeasonID *uuid.UUID
	Limit    int32
	Offset   int32
}

// ContinueWatchingItem represents a series in the user's continue watching list.
type ContinueWatchingItem struct {
	Series            *Series
	LastEpisodeID     uuid.UUID
	LastSeasonNumber  int32
	LastEpisodeNumber int32
	LastEpisodeTitle  string
	ProgressSeconds   int32
	DurationSeconds   int32
	LastWatchedAt     time.Time
}

// NextEpisode represents the next episode to watch in a series.
type NextEpisode struct {
	Episode       *Episode
	IsNewSeason   bool
	IsSeriesFinale bool
}

// formatEpisodeCode formats season and episode numbers as SxxExx.
func formatEpisodeCode(season, episode int32) string {
	return fmt.Sprintf("S%02dE%02d", season, episode)
}
