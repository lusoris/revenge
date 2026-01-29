// Package tvshow provides TV show content management functionality.
package tvshow

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/lusoris/revenge/internal/content/tvshow/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// Series represents a TV series in the library.
type Series struct {
	shared.BaseEntity

	LibraryID uuid.UUID

	// Metadata
	Title         string
	SortTitle     string
	OriginalTitle string
	Tagline       string
	Overview      string

	// Airing info
	FirstAirDate *time.Time
	LastAirDate  *time.Time
	Year         int
	Status       string // Returning Series, Ended, Canceled
	Type         string // Scripted, Documentary, Reality

	// Ratings
	ContentRating   string
	RatingLevel     int
	CommunityRating float64
	VoteCount       int

	// Counts
	SeasonCount  int
	EpisodeCount int
	SpecialCount int

	// Images
	PosterPath       string
	PosterBlurhash   string
	BackdropPath     string
	BackdropBlurhash string
	LogoPath         string

	// External IDs
	TmdbID int
	ImdbID string
	TvdbID int

	// Primary network
	NetworkName     string
	NetworkLogoPath string

	// Stats
	DateAdded    time.Time
	LastPlayedAt *time.Time
	IsLocked     bool

	// Loaded on demand
	Seasons  []Season
	Cast     []CastMember
	Crew     []CrewMember
	Creators []CrewMember
	Networks []Network
	Genres   []Genre
	Images   []Image
	Videos   []Video
}

// Season represents a season of a TV series.
type Season struct {
	shared.BaseEntity

	SeriesID     uuid.UUID
	SeasonNumber int
	Name         string
	Overview     string

	// Air info
	AirDate *time.Time
	Year    int

	// Episode count
	EpisodeCount int

	// Images
	PosterPath     string
	PosterBlurhash string

	// External IDs
	TmdbID int
	TvdbID int

	// Loaded on demand
	Episodes []Episode
}

// Episode represents an episode of a TV series.
type Episode struct {
	shared.ContentEntity

	SeriesID uuid.UUID
	SeasonID uuid.UUID

	// File info
	Container    string
	SizeBytes    int64
	RuntimeTicks int64

	// Episode info
	SeasonNumber   int
	EpisodeNumber  int
	AbsoluteNumber *int

	// Metadata
	Overview       string
	ProductionCode string

	// Air info
	AirDate    *time.Time
	AirDateUTC *time.Time

	// Ratings
	CommunityRating float64
	VoteCount       int

	// Images
	StillPath     string
	StillBlurhash string

	// External IDs
	TmdbID int
	ImdbID string
	TvdbID int

	// Stats
	DateAdded    time.Time
	LastPlayedAt *time.Time
	PlayCount    int
	IsLocked     bool

	// Loaded on demand
	Cast       []CastMember
	GuestStars []CastMember
	Crew       []CrewMember
	Directors  []CrewMember
	Writers    []CrewMember
}

// Network represents a TV network.
type Network struct {
	ID            uuid.UUID
	Name          string
	LogoPath      string
	OriginCountry string
	TmdbID        int
	CreatedAt     time.Time
}

// Genre represents a TV show genre.
type Genre struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
}

// CastMember represents an actor's role in a series or episode.
type CastMember struct {
	PersonID             uuid.UUID
	Name                 string
	CharacterName        string
	BillingOrder         int
	IsGuest              bool
	PrimaryImageURL      string
	PrimaryImageBlurhash string
}

// CrewMember represents a crew member's role in a series or episode.
type CrewMember struct {
	PersonID             uuid.UUID
	Name                 string
	Role                 string
	Department           string
	Job                  string
	BillingOrder         int
	PrimaryImageURL      string
	PrimaryImageBlurhash string
}

// Image represents additional artwork for a series, season, or episode.
type Image struct {
	ID          uuid.UUID
	ImageType   string
	URL         string
	LocalPath   string
	Width       int
	Height      int
	AspectRatio float64
	Language    string
	VoteAverage float64
	VoteCount   int
	Blurhash    string
	Provider    string
	ProviderID  string
	IsPrimary   bool
	CreatedAt   time.Time
}

// Video represents a series trailer or clip.
type Video struct {
	ID         uuid.UUID
	VideoType  string
	Name       string
	Key        string
	Site       string
	Size       int
	Language   string
	IsOfficial bool
	TmdbID     string
	CreatedAt  time.Time
}

// SeriesUserRating represents a user's rating for a series.
type SeriesUserRating struct {
	UserID    uuid.UUID
	SeriesID  uuid.UUID
	Rating    float64
	Review    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EpisodeUserRating represents a user's rating for an episode.
type EpisodeUserRating struct {
	UserID    uuid.UUID
	EpisodeID uuid.UUID
	Rating    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EpisodeWatchHistory represents a user's watch progress for an episode.
type EpisodeWatchHistory struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	ProfileID        *uuid.UUID
	EpisodeID        uuid.UUID
	PositionTicks    int64
	DurationTicks    int64
	PlayedPercentage float64
	Completed        bool
	CompletedAt      *time.Time
	DeviceName       string
	DeviceType       string
	ClientName       string
	PlayMethod       string
	StartedAt        time.Time
	LastUpdatedAt    time.Time
}

// SeriesWatchProgress represents a user's overall progress through a series.
type SeriesWatchProgress struct {
	UserID            uuid.UUID
	SeriesID          uuid.UUID
	LastEpisodeID     *uuid.UUID
	LastSeasonNumber  int
	LastEpisodeNumber int
	TotalEpisodes     int
	WatchedEpisodes   int
	ProgressPercent   float64
	IsWatching        bool
	StartedAt         *time.Time
	LastWatchedAt     *time.Time
	CompletedAt       *time.Time
}

// FromDBSeries converts a database series to a domain series.
func FromDBSeries(s *db.Series) *Series {
	if s == nil {
		return nil
	}

	series := &Series{
		BaseEntity: shared.BaseEntity{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		LibraryID:    s.LibraryID,
		Title:        s.Title,
		DateAdded:    s.DateAdded,
		SeasonCount:  int(s.SeasonCount),
		EpisodeCount: int(s.EpisodeCount),
		SpecialCount: int(s.SpecialCount),
		IsLocked:     s.IsLocked,
	}

	// Optional fields
	if s.SortTitle != nil {
		series.SortTitle = *s.SortTitle
	}
	if s.OriginalTitle != nil {
		series.OriginalTitle = *s.OriginalTitle
	}
	if s.Tagline != nil {
		series.Tagline = *s.Tagline
	}
	if s.Overview != nil {
		series.Overview = *s.Overview
	}
	if s.FirstAirDate.Valid {
		t := s.FirstAirDate.Time
		series.FirstAirDate = &t
	}
	if s.LastAirDate.Valid {
		t := s.LastAirDate.Time
		series.LastAirDate = &t
	}
	if s.Year != nil {
		series.Year = int(*s.Year)
	}
	if s.Status != nil {
		series.Status = *s.Status
	}
	if s.Type != nil {
		series.Type = *s.Type
	}
	if s.ContentRating != nil {
		series.ContentRating = *s.ContentRating
	}
	if s.RatingLevel != nil {
		series.RatingLevel = int(*s.RatingLevel)
	}
	if s.CommunityRating.Valid {
		f, _ := s.CommunityRating.Float64Value()
		series.CommunityRating = f.Float64
	}
	if s.VoteCount != nil {
		series.VoteCount = int(*s.VoteCount)
	}
	if s.PosterPath != nil {
		series.PosterPath = *s.PosterPath
	}
	if s.PosterBlurhash != nil {
		series.PosterBlurhash = *s.PosterBlurhash
	}
	if s.BackdropPath != nil {
		series.BackdropPath = *s.BackdropPath
	}
	if s.BackdropBlurhash != nil {
		series.BackdropBlurhash = *s.BackdropBlurhash
	}
	if s.LogoPath != nil {
		series.LogoPath = *s.LogoPath
	}
	if s.TmdbID != nil {
		series.TmdbID = int(*s.TmdbID)
	}
	if s.ImdbID != nil {
		series.ImdbID = *s.ImdbID
	}
	if s.TvdbID != nil {
		series.TvdbID = int(*s.TvdbID)
	}
	if s.NetworkName != nil {
		series.NetworkName = *s.NetworkName
	}
	if s.NetworkLogoPath != nil {
		series.NetworkLogoPath = *s.NetworkLogoPath
	}
	if s.LastPlayedAt.Valid {
		t := s.LastPlayedAt.Time
		series.LastPlayedAt = &t
	}

	return series
}

// FromDBSeason converts a database season to a domain season.
func FromDBSeason(s *db.Season) *Season {
	if s == nil {
		return nil
	}

	season := &Season{
		BaseEntity: shared.BaseEntity{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		SeriesID:     s.SeriesID,
		SeasonNumber: int(s.SeasonNumber),
		EpisodeCount: int(s.EpisodeCount),
	}

	if s.Name != nil {
		season.Name = *s.Name
	}
	if s.Overview != nil {
		season.Overview = *s.Overview
	}
	if s.AirDate.Valid {
		t := s.AirDate.Time
		season.AirDate = &t
	}
	if s.Year != nil {
		season.Year = int(*s.Year)
	}
	if s.PosterPath != nil {
		season.PosterPath = *s.PosterPath
	}
	if s.PosterBlurhash != nil {
		season.PosterBlurhash = *s.PosterBlurhash
	}
	if s.TmdbID != nil {
		season.TmdbID = int(*s.TmdbID)
	}
	if s.TvdbID != nil {
		season.TvdbID = int(*s.TvdbID)
	}

	return season
}

// FromDBEpisode converts a database episode to a domain episode.
func FromDBEpisode(e *db.Episode) *Episode {
	if e == nil {
		return nil
	}

	episode := &Episode{
		ContentEntity: shared.ContentEntity{
			BaseEntity: shared.BaseEntity{
				ID:        e.ID,
				CreatedAt: e.CreatedAt,
				UpdatedAt: e.UpdatedAt,
			},
			LibraryID: uuid.UUID{}, // Episodes don't have direct library_id
			Path:      e.Path,
			Title:     e.Title,
		},
		SeriesID:      e.SeriesID,
		SeasonID:      e.SeasonID,
		SeasonNumber:  int(e.SeasonNumber),
		EpisodeNumber: int(e.EpisodeNumber),
		DateAdded:     e.DateAdded,
		PlayCount:     int(e.PlayCount),
		IsLocked:      e.IsLocked,
	}

	if e.Container != nil {
		episode.Container = *e.Container
	}
	if e.SizeBytes != nil {
		episode.SizeBytes = *e.SizeBytes
	}
	if e.RuntimeTicks != nil {
		episode.RuntimeTicks = *e.RuntimeTicks
	}
	if e.AbsoluteNumber != nil {
		n := int(*e.AbsoluteNumber)
		episode.AbsoluteNumber = &n
	}
	if e.Overview != nil {
		episode.Overview = *e.Overview
	}
	if e.ProductionCode != nil {
		episode.ProductionCode = *e.ProductionCode
	}
	if e.AirDate.Valid {
		t := e.AirDate.Time
		episode.AirDate = &t
	}
	if e.AirDateUtc.Valid {
		t := e.AirDateUtc.Time
		episode.AirDateUTC = &t
	}
	if e.CommunityRating.Valid {
		f, _ := e.CommunityRating.Float64Value()
		episode.CommunityRating = f.Float64
	}
	if e.VoteCount != nil {
		episode.VoteCount = int(*e.VoteCount)
	}
	if e.StillPath != nil {
		episode.StillPath = *e.StillPath
	}
	if e.StillBlurhash != nil {
		episode.StillBlurhash = *e.StillBlurhash
	}
	if e.TmdbID != nil {
		episode.TmdbID = int(*e.TmdbID)
	}
	if e.ImdbID != nil {
		episode.ImdbID = *e.ImdbID
	}
	if e.TvdbID != nil {
		episode.TvdbID = int(*e.TvdbID)
	}
	if e.LastPlayedAt.Valid {
		t := e.LastPlayedAt.Time
		episode.LastPlayedAt = &t
	}

	return episode
}

// ToDBCreateParams converts a domain series to database create params.
func (s *Series) ToDBCreateParams() db.CreateSeriesParams {
	params := db.CreateSeriesParams{
		LibraryID: s.LibraryID,
		Title:     s.Title,
	}

	if s.SortTitle != "" {
		params.SortTitle = &s.SortTitle
	}
	if s.OriginalTitle != "" {
		params.OriginalTitle = &s.OriginalTitle
	}
	if s.Tagline != "" {
		params.Tagline = &s.Tagline
	}
	if s.Overview != "" {
		params.Overview = &s.Overview
	}
	if s.FirstAirDate != nil {
		params.FirstAirDate = pgtype.Date{Time: *s.FirstAirDate, Valid: true}
	}
	if s.LastAirDate != nil {
		params.LastAirDate = pgtype.Date{Time: *s.LastAirDate, Valid: true}
	}
	if s.Year > 0 {
		y := int32(s.Year)
		params.Year = &y
	}
	if s.Status != "" {
		params.Status = &s.Status
	}
	if s.Type != "" {
		params.Type = &s.Type
	}
	if s.ContentRating != "" {
		params.ContentRating = &s.ContentRating
	}
	if s.RatingLevel > 0 {
		rl := int32(s.RatingLevel)
		params.RatingLevel = &rl
	}
	if s.CommunityRating > 0 {
		params.CommunityRating = numericFromFloat(s.CommunityRating)
	}
	if s.VoteCount > 0 {
		vc := int32(s.VoteCount)
		params.VoteCount = &vc
	}
	if s.PosterPath != "" {
		params.PosterPath = &s.PosterPath
	}
	if s.PosterBlurhash != "" {
		params.PosterBlurhash = &s.PosterBlurhash
	}
	if s.BackdropPath != "" {
		params.BackdropPath = &s.BackdropPath
	}
	if s.BackdropBlurhash != "" {
		params.BackdropBlurhash = &s.BackdropBlurhash
	}
	if s.LogoPath != "" {
		params.LogoPath = &s.LogoPath
	}
	if s.TmdbID > 0 {
		id := int32(s.TmdbID)
		params.TmdbID = &id
	}
	if s.ImdbID != "" {
		params.ImdbID = &s.ImdbID
	}
	if s.TvdbID > 0 {
		id := int32(s.TvdbID)
		params.TvdbID = &id
	}
	if s.NetworkName != "" {
		params.NetworkName = &s.NetworkName
	}
	if s.NetworkLogoPath != "" {
		params.NetworkLogoPath = &s.NetworkLogoPath
	}

	return params
}

// RuntimeDuration returns the episode runtime as time.Duration.
func (e *Episode) RuntimeDuration() time.Duration {
	return time.Duration(e.RuntimeTicks * 100)
}

// RuntimeMinutes returns the episode runtime in minutes.
func (e *Episode) RuntimeMinutes() int {
	if e.RuntimeTicks == 0 {
		return 0
	}
	return int(e.RuntimeTicks / 10_000_000 / 60)
}

// numericFromFloat converts a float64 to pgtype.Numeric.
func numericFromFloat(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(f)
	return n
}
