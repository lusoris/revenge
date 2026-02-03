package movie

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
