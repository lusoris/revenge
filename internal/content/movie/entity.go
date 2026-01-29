// Package movie provides movie content management functionality.
package movie

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/lusoris/revenge/internal/content/movie/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// Movie represents a movie in the library.
type Movie struct {
	shared.ContentEntity

	// File info
	Container    string
	SizeBytes    int64
	RuntimeTicks int64

	// Metadata
	OriginalTitle string
	Tagline       string
	Overview      string
	ReleaseDate   *time.Time
	Year          int
	ContentRating string
	RatingLevel   int

	// Financial
	Budget  int64
	Revenue int64

	// Ratings
	CommunityRating float64
	VoteCount       int
	CriticRating    float64
	CriticCount     int

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

	// Collection
	CollectionID    *uuid.UUID
	CollectionOrder int

	// Stats
	DateAdded    time.Time
	LastPlayedAt *time.Time
	PlayCount    int
	IsLocked     bool

	// Loaded on demand
	Collection *Collection
	Cast       []CastMember
	Crew       []CrewMember
	Directors  []CrewMember
	Writers    []CrewMember
	Studios    []Studio
	Genres     []Genre
	Images     []Image
	Videos     []Video
}

// Collection represents a movie collection (e.g., "The Dark Knight Trilogy").
type Collection struct {
	ID               uuid.UUID
	Name             string
	SortName         string
	Overview         string
	PosterPath       string
	PosterBlurhash   string
	BackdropPath     string
	BackdropBlurhash string
	TmdbID           int
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Loaded on demand
	Movies []*Movie
}

// Studio represents a movie production studio.
type Studio struct {
	ID        uuid.UUID
	Name      string
	LogoPath  string
	TmdbID    int
	CreatedAt time.Time
}

// Genre represents a movie genre.
type Genre struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
}

// CastMember represents an actor's role in a movie.
type CastMember struct {
	PersonID             uuid.UUID
	Name                 string
	CharacterName        string
	BillingOrder         int
	IsGuest              bool
	PrimaryImageURL      string
	PrimaryImageBlurhash string
}

// CrewMember represents a crew member's role in a movie.
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

// Image represents additional artwork for a movie.
type Image struct {
	ID          uuid.UUID
	ImageType   string
	Path        string
	Blurhash    string
	Width       int
	Height      int
	AspectRatio float64
	Language    string
	VoteAverage float64
	VoteCount   int
	IsPrimary   bool
	Source      string
	CreatedAt   time.Time
}

// Video represents a movie trailer or clip.
type Video struct {
	ID        uuid.UUID
	VideoType string
	Site      string
	Key       string
	Name      string
	Language  string
	Size      int
	CreatedAt time.Time
}

// UserRating represents a user's rating for a movie.
type UserRating struct {
	UserID    uuid.UUID
	MovieID   uuid.UUID
	Rating    float64
	Review    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WatchHistory represents a user's watch progress for a movie.
type WatchHistory struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	ProfileID        *uuid.UUID
	MovieID          uuid.UUID
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

// FromDBMovie converts a database movie to a domain movie.
func FromDBMovie(m *db.Movie) *Movie {
	if m == nil {
		return nil
	}

	movie := &Movie{
		ContentEntity: shared.ContentEntity{
			BaseEntity: shared.BaseEntity{
				ID:        m.ID,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
			LibraryID: m.LibraryID,
			Path:      m.Path,
			Title:     m.Title,
		},
		DateAdded: m.DateAdded,
		PlayCount: int(m.PlayCount),
		IsLocked:  m.IsLocked,
	}

	// Optional fields
	if m.SortTitle != nil {
		movie.SortTitle = *m.SortTitle
	}
	if m.Container != nil {
		movie.Container = *m.Container
	}
	if m.SizeBytes != nil {
		movie.SizeBytes = *m.SizeBytes
	}
	if m.RuntimeTicks != nil {
		movie.RuntimeTicks = *m.RuntimeTicks
	}
	if m.OriginalTitle != nil {
		movie.OriginalTitle = *m.OriginalTitle
	}
	if m.Tagline != nil {
		movie.Tagline = *m.Tagline
	}
	if m.Overview != nil {
		movie.Overview = *m.Overview
	}
	if m.ReleaseDate.Valid {
		t := m.ReleaseDate.Time
		movie.ReleaseDate = &t
	}
	if m.Year != nil {
		movie.Year = int(*m.Year)
	}
	if m.ContentRating != nil {
		movie.ContentRating = *m.ContentRating
	}
	if m.RatingLevel != nil {
		movie.RatingLevel = int(*m.RatingLevel)
	}
	if m.Budget != nil {
		movie.Budget = *m.Budget
	}
	if m.Revenue != nil {
		movie.Revenue = *m.Revenue
	}
	if m.CommunityRating.Valid {
		f, _ := m.CommunityRating.Float64Value()
		movie.CommunityRating = f.Float64
	}
	if m.VoteCount != nil {
		movie.VoteCount = int(*m.VoteCount)
	}
	if m.CriticRating.Valid {
		f, _ := m.CriticRating.Float64Value()
		movie.CriticRating = f.Float64
	}
	if m.CriticCount != nil {
		movie.CriticCount = int(*m.CriticCount)
	}
	if m.PosterPath != nil {
		movie.PosterPath = *m.PosterPath
	}
	if m.PosterBlurhash != nil {
		movie.PosterBlurhash = *m.PosterBlurhash
	}
	if m.BackdropPath != nil {
		movie.BackdropPath = *m.BackdropPath
	}
	if m.BackdropBlurhash != nil {
		movie.BackdropBlurhash = *m.BackdropBlurhash
	}
	if m.LogoPath != nil {
		movie.LogoPath = *m.LogoPath
	}
	if m.TmdbID != nil {
		movie.TmdbID = int(*m.TmdbID)
	}
	if m.ImdbID != nil {
		movie.ImdbID = *m.ImdbID
	}
	if m.TvdbID != nil {
		movie.TvdbID = int(*m.TvdbID)
	}
	if m.CollectionID.Valid {
		id := uuid.UUID(m.CollectionID.Bytes)
		movie.CollectionID = &id
	}
	if m.CollectionOrder != nil {
		movie.CollectionOrder = int(*m.CollectionOrder)
	}
	if m.LastPlayedAt.Valid {
		t := m.LastPlayedAt.Time
		movie.LastPlayedAt = &t
	}

	return movie
}

// ToDBCreateParams converts a domain movie to database create params.
func (m *Movie) ToDBCreateParams() db.CreateMovieParams {
	params := db.CreateMovieParams{
		LibraryID: m.LibraryID,
		Path:      m.Path,
		Title:     m.Title,
	}

	if m.SortTitle != "" {
		params.SortTitle = &m.SortTitle
	}
	if m.Container != "" {
		params.Container = &m.Container
	}
	if m.SizeBytes > 0 {
		params.SizeBytes = &m.SizeBytes
	}
	if m.RuntimeTicks > 0 {
		params.RuntimeTicks = &m.RuntimeTicks
	}
	if m.OriginalTitle != "" {
		params.OriginalTitle = &m.OriginalTitle
	}
	if m.Tagline != "" {
		params.Tagline = &m.Tagline
	}
	if m.Overview != "" {
		params.Overview = &m.Overview
	}
	if m.ReleaseDate != nil {
		params.ReleaseDate = pgtype.Date{Time: *m.ReleaseDate, Valid: true}
	}
	if m.Year > 0 {
		y := int32(m.Year)
		params.Year = &y
	}
	if m.ContentRating != "" {
		params.ContentRating = &m.ContentRating
	}
	if m.RatingLevel > 0 {
		rl := int32(m.RatingLevel)
		params.RatingLevel = &rl
	}
	if m.Budget > 0 {
		params.Budget = &m.Budget
	}
	if m.Revenue > 0 {
		params.Revenue = &m.Revenue
	}
	if m.CommunityRating > 0 {
		params.CommunityRating = numericFromFloat(m.CommunityRating)
	}
	if m.VoteCount > 0 {
		vc := int32(m.VoteCount)
		params.VoteCount = &vc
	}
	if m.PosterPath != "" {
		params.PosterPath = &m.PosterPath
	}
	if m.PosterBlurhash != "" {
		params.PosterBlurhash = &m.PosterBlurhash
	}
	if m.BackdropPath != "" {
		params.BackdropPath = &m.BackdropPath
	}
	if m.BackdropBlurhash != "" {
		params.BackdropBlurhash = &m.BackdropBlurhash
	}
	if m.LogoPath != "" {
		params.LogoPath = &m.LogoPath
	}
	if m.TmdbID > 0 {
		id := int32(m.TmdbID)
		params.TmdbID = &id
	}
	if m.ImdbID != "" {
		params.ImdbID = &m.ImdbID
	}
	if m.TvdbID > 0 {
		id := int32(m.TvdbID)
		params.TvdbID = &id
	}
	if m.CollectionID != nil {
		params.CollectionID = pgtype.UUID{Bytes: *m.CollectionID, Valid: true}
	}
	if m.CollectionOrder > 0 {
		co := int32(m.CollectionOrder)
		params.CollectionOrder = &co
	}

	return params
}

// RuntimeDuration returns the runtime as a time.Duration.
func (m *Movie) RuntimeDuration() time.Duration {
	return time.Duration(m.RuntimeTicks * 100)
}

// RuntimeMinutes returns the runtime in minutes.
func (m *Movie) RuntimeMinutes() int {
	if m.RuntimeTicks == 0 {
		return 0
	}
	return int(m.RuntimeTicks / 10_000_000 / 60)
}
