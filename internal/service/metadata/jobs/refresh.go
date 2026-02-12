// Package jobs provides River job definitions for metadata operations.
package jobs

import (
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
)

// RefreshMovieArgs are the arguments for refreshing movie metadata.
// This job is handled by the movie module's worker.
type RefreshMovieArgs struct {
	MovieID   uuid.UUID `json:"movie_id"`
	Force     bool      `json:"force"`
	Languages []string  `json:"languages,omitempty"`
}

// Kind returns the unique job kind for River.
func (RefreshMovieArgs) Kind() string {
	return "metadata_refresh_movie"
}

// InsertOpts returns the default insert options for movie metadata refresh jobs.
// Metadata refreshes hit external APIs (TMDb) — limit retries to avoid hammering.
func (RefreshMovieArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

// RefreshTVShowArgs are the arguments for refreshing TV show metadata.
// This job is handled by the tvshow module's worker.
type RefreshTVShowArgs struct {
	SeriesID        uuid.UUID `json:"series_id"`
	Force           bool      `json:"force"`
	Languages       []string  `json:"languages,omitempty"`
	IncludeSeasons  bool      `json:"include_seasons"`
	IncludeEpisodes bool      `json:"include_episodes"`
}

// Kind returns the unique job kind for River.
func (RefreshTVShowArgs) Kind() string {
	return "metadata_refresh_tvshow"
}

// InsertOpts returns the default insert options for TV show metadata refresh jobs.
func (RefreshTVShowArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

// RefreshSeasonArgs are the arguments for refreshing season metadata.
type RefreshSeasonArgs struct {
	SeriesID        uuid.UUID `json:"series_id"`
	SeasonID        uuid.UUID `json:"season_id"`
	SeasonNumber    int       `json:"season_number"`
	Force           bool      `json:"force"`
	Languages       []string  `json:"languages,omitempty"`
	IncludeEpisodes bool      `json:"include_episodes"`
}

// Kind returns the unique job kind for River.
func (RefreshSeasonArgs) Kind() string {
	return "metadata_refresh_season"
}

// InsertOpts returns the default insert options for season metadata refresh jobs.
func (RefreshSeasonArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

// RefreshEpisodeArgs are the arguments for refreshing episode metadata.
type RefreshEpisodeArgs struct {
	SeriesID      uuid.UUID `json:"series_id"`
	SeasonID      uuid.UUID `json:"season_id"`
	EpisodeID     uuid.UUID `json:"episode_id"`
	SeasonNumber  int       `json:"season_number"`
	EpisodeNumber int       `json:"episode_number"`
	Force         bool      `json:"force"`
	Languages     []string  `json:"languages,omitempty"`
}

// Kind returns the unique job kind for River.
func (RefreshEpisodeArgs) Kind() string {
	return "metadata_refresh_episode"
}

// InsertOpts returns the default insert options for episode metadata refresh jobs.
func (RefreshEpisodeArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

// RefreshPersonArgs are the arguments for refreshing person metadata.
type RefreshPersonArgs struct {
	PersonID   uuid.UUID `json:"person_id"`
	ProviderID string    `json:"provider_id"`
	Force      bool      `json:"force"`
	Languages  []string  `json:"languages,omitempty"`
}

// Kind returns the unique job kind for River.
func (RefreshPersonArgs) Kind() string {
	return "metadata_refresh_person"
}

// InsertOpts returns the default insert options for person metadata refresh jobs.
func (RefreshPersonArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 30 * time.Minute,
		},
	}
}

// EnrichContentArgs are the arguments for enriching content with multiple providers.
type EnrichContentArgs struct {
	ContentType string    `json:"content_type"` // "movie" or "tvshow"
	ContentID   uuid.UUID `json:"content_id"`
	Providers   []string  `json:"providers,omitempty"` // Empty means all available
	Languages   []string  `json:"languages,omitempty"`
}

// Kind returns the unique job kind for River.
func (EnrichContentArgs) Kind() string {
	return "metadata_enrich_content"
}

// InsertOpts returns the default insert options for content enrichment jobs.
func (EnrichContentArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour,
		},
	}
}

// DownloadImageArgs are the arguments for downloading and caching an image.
type DownloadImageArgs struct {
	ContentType string `json:"content_type"` // "movie", "tvshow", "person", etc.
	ContentID   string `json:"content_id"`
	ImageType   string `json:"image_type"` // "poster", "backdrop", "profile", "still"
	Path        string `json:"path"`       // Provider image path
	Size        string `json:"size"`       // Image size
}

// Kind returns the unique job kind for River.
func (DownloadImageArgs) Kind() string {
	return "metadata_download_image"
}

// InsertOpts returns the default insert options for image download jobs.
func (DownloadImageArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       infrajobs.QueueDefault,
		MaxAttempts: 5,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: 1 * time.Hour,
		},
	}
}
