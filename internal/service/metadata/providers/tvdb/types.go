// Package tvdb provides a TVDb (TheTVDB) provider implementation.
package tvdb

import (
	"time"
)

// API response types for TVDb v4 API.

// LoginRequest is the TVDb API login request.
type LoginRequest struct {
	APIKey string `json:"apikey"`
	PIN    string `json:"pin,omitempty"`
}

// LoginResponse is the TVDb API login response.
type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

// BaseResponse wraps all TVDb API responses.
type BaseResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

// ListResponse wraps paginated list responses.
type ListResponse[T any] struct {
	Status string `json:"status"`
	Data   []T    `json:"data"`
	Links  struct {
		Prev  *string `json:"prev"`
		Self  string  `json:"self"`
		Next  *string `json:"next"`
		Total int     `json:"total_items"`
	} `json:"links"`
}

// SearchResponse is the TVDb API search response.
type SearchResponse struct {
	Status string         `json:"status"`
	Data   []SearchResult `json:"data"`
}

// SearchResult is a single search result.
type SearchResult struct {
	ObjectID     string            `json:"objectID"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	Type         string            `json:"type"` // "series", "movie", "person"
	Year         string            `json:"year"`
	Network      string            `json:"network"`
	Overview     string            `json:"overview"`
	Status       string            `json:"status"`
	PrimaryType  string            `json:"primary_type"`
	ImageURL     *string           `json:"image_url"`
	Thumbnail    *string           `json:"thumbnail"`
	FirstAirTime string            `json:"first_air_time"`
	TVDbID       string            `json:"tvdb_id"`
	RemoteIDs    []string          `json:"remote_ids"`
	Aliases      []string          `json:"aliases"`
	Overviews    map[string]string `json:"overviews"`
	Translations map[string]string `json:"translations"`
}

// SeriesResponse is the TVDb API series details response.
type SeriesResponse struct {
	ObjectID             string                  `json:"objectID"`
	ID                   int                     `json:"id"`
	Name                 string                  `json:"name"`
	Slug                 string                  `json:"slug"`
	Year                 string                  `json:"year"`
	FirstAired           string                  `json:"firstAired"`
	LastAired            string                  `json:"lastAired"`
	NextAired            *string                 `json:"nextAired"`
	Status               *StatusResponse         `json:"status"`
	DefaultSeasonType    int                     `json:"defaultSeasonType"`
	IsOrderRandomized    bool                    `json:"isOrderRandomized"`
	LastUpdated          string                  `json:"lastUpdated"`
	AverageRuntime       int                     `json:"averageRuntime"`
	Score                int                     `json:"score"`
	OriginalCountry      string                  `json:"originalCountry"`
	OriginalLanguage     string                  `json:"originalLanguage"`
	OriginalNetwork      *NetworkResponse        `json:"originalNetwork"`
	Overview             *string                 `json:"overview"`
	Overviews            map[string]string       `json:"overviews"`
	NameTranslations     []string                `json:"nameTranslations"`
	OverviewTranslations []string                `json:"overviewTranslations"`
	Aliases              []AliasResponse         `json:"aliases"`
	Image                *string                 `json:"image"`
	Artworks             []ArtworkResponse       `json:"artworks"`
	Genres               []GenreResponse         `json:"genres"`
	Characters           []CharacterResponse     `json:"characters"`
	RemoteIDs            []RemoteIDResponse      `json:"remoteIds"`
	Seasons              []SeasonSummaryResponse `json:"seasons"`
	Networks             []NetworkResponse       `json:"networks"`
	ContentRatings       []ContentRatingResponse `json:"contentRatings"`
	Companies            []CompanyResponse       `json:"companies"`
	Tags                 []TagResponse           `json:"tags"`
	Trailers             []TrailerResponse       `json:"trailers"`
}

// SeriesExtendedResponse includes additional data.
type SeriesExtendedResponse struct {
	SeriesResponse
	Episodes []EpisodeResponse `json:"episodes"`
}

// SeasonResponse is the TVDb API season details response.
type SeasonResponse struct {
	ID               int                 `json:"id"`
	SeriesID         int                 `json:"seriesId"`
	Type             *SeasonTypeResponse `json:"type"`
	Name             string              `json:"name"`
	Number           int                 `json:"number"`
	Year             string              `json:"year"`
	LastUpdated      string              `json:"lastUpdated"`
	Image            *string             `json:"image"`
	Overview         *string             `json:"overview"`
	Overviews        map[string]string   `json:"overviews"`
	NameTranslations []string            `json:"nameTranslations"`
	Companies        []CompanyResponse   `json:"companies"`
}

// SeasonSummaryResponse is a season summary in series details.
type SeasonSummaryResponse struct {
	ID          int                 `json:"id"`
	SeriesID    int                 `json:"seriesId"`
	Type        *SeasonTypeResponse `json:"type"`
	Name        string              `json:"name"`
	Number      int                 `json:"number"`
	Year        string              `json:"year"`
	Image       *string             `json:"image"`
	LastUpdated string              `json:"lastUpdated"`
}

// SeasonTypeResponse describes season ordering type.
type SeasonTypeResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	AlternateName *string `json:"alternateName"`
}

// EpisodeResponse is the TVDb API episode details response.
type EpisodeResponse struct {
	ID               int                     `json:"id"`
	SeriesID         int                     `json:"seriesId"`
	Name             string                  `json:"name"`
	Aired            string                  `json:"aired"`
	Runtime          *int                    `json:"runtime"`
	SeasonNumber     int                     `json:"seasonNumber"`
	Number           int                     `json:"number"`
	AbsoluteNumber   *int                    `json:"absoluteNumber"`
	LastUpdated      string                  `json:"lastUpdated"`
	FinaleType       *string                 `json:"finaleType"`
	Year             string                  `json:"year"`
	Image            *string                 `json:"image"`
	Overview         *string                 `json:"overview"`
	Overviews        map[string]string       `json:"overviews"`
	NameTranslations []string                `json:"nameTranslations"`
	ProductionCode   *string                 `json:"productionCode"`
	IsMovie          bool                    `json:"isMovie"`
	Seasons          []SeasonSummaryResponse `json:"seasons"`
	Characters       []CharacterResponse     `json:"characters"`
	ContentRatings   []ContentRatingResponse `json:"contentRatings"`
}

// MovieResponse is the TVDb API movie details response.
type MovieResponse struct {
	ID                   int                     `json:"id"`
	Name                 string                  `json:"name"`
	Slug                 string                  `json:"slug"`
	Image                *string                 `json:"image"`
	Year                 string                  `json:"year"`
	Runtime              *int                    `json:"runtime"`
	LastUpdated          string                  `json:"lastUpdated"`
	Score                int                     `json:"score"`
	Status               *StatusResponse         `json:"status"`
	OriginalCountry      string                  `json:"originalCountry"`
	OriginalLanguage     string                  `json:"originalLanguage"`
	Overview             *string                 `json:"overview"`
	Overviews            map[string]string       `json:"overviews"`
	NameTranslations     []string                `json:"nameTranslations"`
	OverviewTranslations []string                `json:"overviewTranslations"`
	Aliases              []AliasResponse         `json:"aliases"`
	Artworks             []ArtworkResponse       `json:"artworks"`
	Genres               []GenreResponse         `json:"genres"`
	Characters           []CharacterResponse     `json:"characters"`
	RemoteIDs            []RemoteIDResponse      `json:"remoteIds"`
	ContentRatings       []ContentRatingResponse `json:"contentRatings"`
	Companies            []CompanyResponse       `json:"companies"`
	Trailers             []TrailerResponse       `json:"trailers"`
	Budget               *string                 `json:"budget"`
	BoxOffice            *string                 `json:"boxOffice"`
	Releases             []ReleaseResponse       `json:"releases"`
}

// PersonResponse is the TVDb API person details response.
type PersonResponse struct {
	ID                   int                 `json:"id"`
	Name                 string              `json:"name"`
	Slug                 string              `json:"slug"`
	Image                *string             `json:"image"`
	Birth                *string             `json:"birth"`
	Death                *string             `json:"death"`
	BirthPlace           *string             `json:"birthPlace"`
	Gender               int                 `json:"gender"`
	Score                int                 `json:"score"`
	LastUpdated          string              `json:"lastUpdated"`
	Aliases              []AliasResponse     `json:"aliases"`
	NameTranslations     []string            `json:"nameTranslations"`
	OverviewTranslations []string            `json:"overviewTranslations"`
	Overviews            map[string]string   `json:"overviews"`
	RemoteIDs            []RemoteIDResponse  `json:"remoteIds"`
	Characters           []CharacterResponse `json:"characters"`
	Biographies          []BiographyResponse `json:"biographies"`
}

// CharacterResponse is a character/credit entry.
type CharacterResponse struct {
	ID                   int             `json:"id"`
	Name                 string          `json:"name"`
	PeopleID             *int            `json:"peopleId"`
	SeriesID             *int            `json:"seriesId"`
	MovieID              *int            `json:"movieId"`
	EpisodeID            *int            `json:"episodeId"`
	Type                 int             `json:"type"` // 1=Actor, 2=Director, etc.
	Image                *string         `json:"image"`
	Sort                 int             `json:"sort"`
	IsFeatured           bool            `json:"isFeatured"`
	URL                  *string         `json:"url"`
	NameTranslations     []string        `json:"nameTranslations"`
	OverviewTranslations []string        `json:"overviewTranslations"`
	Aliases              []AliasResponse `json:"aliases"`
	PeopleType           string          `json:"peopleType"`
	PersonName           string          `json:"personName"`
	TagOptions           *TagResponse    `json:"tagOptions"`
	PersonImgURL         *string         `json:"personImgURL"`
}

// ArtworkResponse is an artwork entry.
type ArtworkResponse struct {
	ID              int     `json:"id"`
	Image           string  `json:"image"`
	Thumbnail       string  `json:"thumbnail"`
	Language        *string `json:"language"`
	Type            int     `json:"type"` // 1=Banner, 2=Poster, 3=Background, etc.
	Score           int     `json:"score"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	IncludesText    bool    `json:"includesText"`
	ThumbnailWidth  int     `json:"thumbnailWidth"`
	ThumbnailHeight int     `json:"thumbnailHeight"`
}

// GenreResponse is a genre entry.
type GenreResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// NetworkResponse is a TV network entry.
type NetworkResponse struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Slug               string  `json:"slug"`
	Abbreviation       *string `json:"abbreviation"`
	Country            *string `json:"country"`
	PrimaryCompanyType int     `json:"primaryCompanyType"`
}

// StatusResponse is a content status entry.
type StatusResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	RecordType  string `json:"recordType"`
	KeepUpdated bool   `json:"keepUpdated"`
}

// CompanyResponse is a production company entry.
type CompanyResponse struct {
	ID                   int             `json:"id"`
	Name                 string          `json:"name"`
	Slug                 string          `json:"slug"`
	Country              *string         `json:"country"`
	PrimaryCompanyType   int             `json:"primaryCompanyType"`
	ActiveDate           *string         `json:"activeDate"`
	InactiveDate         *string         `json:"inactiveDate"`
	NameTranslations     []string        `json:"nameTranslations"`
	OverviewTranslations []string        `json:"overviewTranslations"`
	Aliases              []AliasResponse `json:"aliases"`
}

// ContentRatingResponse is a content rating entry.
type ContentRatingResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	ContentType string `json:"contentType"`
	Order       int    `json:"order"`
	FullName    string `json:"fullname"`
	Description string `json:"description"`
}

// RemoteIDResponse is an external ID reference.
type RemoteIDResponse struct {
	ID         string `json:"id"`
	Type       int    `json:"type"` // 2=IMDb, 4=TMDb, etc.
	SourceName string `json:"sourceName"`
}

// AliasResponse is an alternative name.
type AliasResponse struct {
	Language *string `json:"language"`
	Name     string  `json:"name"`
}

// TrailerResponse is a video trailer entry.
type TrailerResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Language string `json:"language"`
	Runtime  int    `json:"runtime"`
}

// TagResponse is a tag entry.
type TagResponse struct {
	ID       int    `json:"id"`
	Tag      int    `json:"tag"`
	TagName  string `json:"tagName"`
	Name     string `json:"name"`
	HelpText string `json:"helpText"`
}

// ReleaseResponse is a movie release entry.
type ReleaseResponse struct {
	Country string `json:"country"`
	Date    string `json:"date"`
	Detail  string `json:"detail"`
}

// BiographyResponse is a person biography entry.
type BiographyResponse struct {
	Biography string `json:"biography"`
	Language  string `json:"language"`
}

// ErrorResponse is the TVDb API error response.
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RemoteIDType constants for TVDb.
const (
	RemoteIDTypeIMDb int = 2
	RemoteIDTypeTMDb int = 4
)

// ArtworkType constants for TVDb.
const (
	ArtworkTypeBanner     int = 1
	ArtworkTypePoster     int = 2
	ArtworkTypeBackground int = 3
	ArtworkTypeIcon       int = 4
	ArtworkTypeClearArt   int = 22
	ArtworkTypeClearLogo  int = 23
)

// CharacterType constants for TVDb.
const (
	CharacterTypeActor    int = 3
	CharacterTypeDirector int = 1
	CharacterTypeWriter   int = 2
	CharacterTypeProducer int = 4
)

// CacheEntry wraps cached data with expiration.
type CacheEntry struct {
	Data      any
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired.
func (c *CacheEntry) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
