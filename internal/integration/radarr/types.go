// Package radarr provides a client for the Radarr API v3.
// Radarr is a PRIMARY metadata provider for movies in the Revenge media system.
package radarr

import (
	"time"

	"github.com/lusoris/revenge/internal/integration/arrbase"
)

// Type aliases for shared arr types.
// These allow existing code to reference radarr.MediaInfo etc. transparently.
type (
	MediaInfo       = arrbase.MediaInfo
	Quality         = arrbase.Quality
	QualityInfo     = arrbase.QualityInfo
	QualityRevision = arrbase.QualityRevision
	Language        = arrbase.Language
	Image           = arrbase.Image
	QualityItem     = arrbase.QualityItem
	FormatItem      = arrbase.FormatItem
	QualityProfile  = arrbase.QualityProfile
	RootFolder      = arrbase.RootFolder
	UnmappedFolder  = arrbase.UnmappedFolder
	SystemStatus    = arrbase.SystemStatus
	Tag             = arrbase.Tag
	WebhookRelease  = arrbase.WebhookRelease
)

// Movie represents a movie in Radarr.
type Movie struct {	ID                    int              `json:"id"`
	Title                 string           `json:"title"`
	OriginalTitle         string           `json:"originalTitle,omitempty"`
	OriginalLanguage      Language         `json:"originalLanguage,omitempty"`
	AlternateTitles       []AlternateTitle `json:"alternateTitles,omitempty"`
	SortTitle             string           `json:"sortTitle,omitempty"`
	SizeOnDisk            int64            `json:"sizeOnDisk,omitempty"`
	Status                string           `json:"status,omitempty"`
	Overview              string           `json:"overview,omitempty"`
	InCinemas             *time.Time       `json:"inCinemas,omitempty"`
	PhysicalRelease       *time.Time       `json:"physicalRelease,omitempty"`
	DigitalRelease        *time.Time       `json:"digitalRelease,omitempty"`
	Images                []Image          `json:"images,omitempty"`
	Website               string           `json:"website,omitempty"`
	Year                  int              `json:"year"`
	HasFile               bool             `json:"hasFile"`
	YouTubeTrailerID      string           `json:"youTubeTrailerId,omitempty"`
	Studio                string           `json:"studio,omitempty"`
	Path                  string           `json:"path,omitempty"`
	QualityProfileID      int              `json:"qualityProfileId"`
	Monitored             bool             `json:"monitored"`
	MinimumAvailability   string           `json:"minimumAvailability,omitempty"`
	IsAvailable           bool             `json:"isAvailable"`
	FolderName            string           `json:"folderName,omitempty"`
	Runtime               int              `json:"runtime,omitempty"`
	CleanTitle            string           `json:"cleanTitle,omitempty"`
	IMDbID                string           `json:"imdbId,omitempty"`
	TMDbID                int              `json:"tmdbId"`
	TitleSlug             string           `json:"titleSlug,omitempty"`
	RootFolderPath        string           `json:"rootFolderPath,omitempty"`
	Certification         string           `json:"certification,omitempty"`
	Genres                []string         `json:"genres,omitempty"`
	Tags                  []int            `json:"tags,omitempty"`
	Added                 time.Time        `json:"added"`
	Ratings               Ratings          `json:"ratings,omitempty"`
	MovieFile             *MovieFile       `json:"movieFile,omitempty"`
	Collection            *Collection      `json:"collection,omitempty"`
	Popularity            float64          `json:"popularity,omitempty"`
	Statistics            *Statistics      `json:"statistics,omitempty"`
}

// MovieFile represents a movie file in Radarr.
type MovieFile struct {
	ID                  int            `json:"id"`
	MovieID             int            `json:"movieId"`
	RelativePath        string         `json:"relativePath,omitempty"`
	Path                string         `json:"path,omitempty"`
	Size                int64          `json:"size"`
	DateAdded           time.Time      `json:"dateAdded"`
	SceneName           string         `json:"sceneName,omitempty"`
	IndexerFlags        int            `json:"indexerFlags,omitempty"`
	Quality             Quality        `json:"quality,omitempty"`
	MediaInfo           *MediaInfo     `json:"mediaInfo,omitempty"`
	OriginalFilePath    string         `json:"originalFilePath,omitempty"`
	QualityCutoffNotMet bool           `json:"qualityCutoffNotMet"`
	Languages           []Language     `json:"languages,omitempty"`
	ReleaseGroup        string         `json:"releaseGroup,omitempty"`
	Edition             string         `json:"edition,omitempty"`
}

// AlternateTitle represents an alternate title (Radarr-specific fields).
type AlternateTitle struct {
	SourceType      string `json:"sourceType,omitempty"`
	MovieMetadataID int    `json:"movieMetadataId,omitempty"`
	Title           string `json:"title,omitempty"`
	CleanTitle      string `json:"cleanTitle,omitempty"`
}

// Ratings contains rating information (Radarr-specific).
type Ratings struct {
	IMDb           *Rating `json:"imdb,omitempty"`
	TMDb           *Rating `json:"tmdb,omitempty"`
	Metacritic     *Rating `json:"metacritic,omitempty"`
	RottenTomatoes *Rating `json:"rottenTomatoes,omitempty"`
}

// Rating represents a single rating source.
type Rating struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
	Type  string  `json:"type,omitempty"`
}

// Collection represents a movie collection.
type Collection struct {
	Name   string  `json:"name,omitempty"`
	TMDbID int     `json:"tmdbId,omitempty"`
	Images []Image `json:"images,omitempty"`
}

// Statistics contains movie statistics.
type Statistics struct {
	MovieFileCount int   `json:"movieFileCount"`
	SizeOnDisk     int64 `json:"sizeOnDisk"`
	ReleaseGroups  []string `json:"releaseGroups,omitempty"`
}



// Command represents a command to execute in Radarr.
type Command struct {
	ID                  int       `json:"id,omitempty"`
	Name                string    `json:"name"`
	CommandName         string    `json:"commandName,omitempty"`
	Message             string    `json:"message,omitempty"`
	Body                CommandBody `json:"body,omitempty"`
	Priority            string    `json:"priority,omitempty"`
	Status              string    `json:"status,omitempty"`
	Queued              time.Time `json:"queued,omitempty"`
	Started             *time.Time `json:"started,omitempty"`
	Ended               *time.Time `json:"ended,omitempty"`
	Duration            string    `json:"duration,omitempty"`
	Trigger             string    `json:"trigger,omitempty"`
	ClientUserAgent     string    `json:"clientUserAgent,omitempty"`
	StateChangeTime     *time.Time `json:"stateChangeTime,omitempty"`
	SendUpdatesToClient bool      `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool      `json:"updateScheduledTask"`
	LastExecutionTime   *time.Time `json:"lastExecutionTime,omitempty"`
}

// CommandBody contains command body parameters.
type CommandBody struct {
	MovieIDs            []int  `json:"movieIds,omitempty"`
	MovieID             int    `json:"movieId,omitempty"`
	SendUpdatesToClient bool   `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool   `json:"updateScheduledTask"`
	CompletionMessage   string `json:"completionMessage,omitempty"`
	RequireDiskAccess   bool   `json:"requiresDiskAccess"`
	IsExclusive         bool   `json:"isExclusive"`
	IsTypeExclusive     bool   `json:"isTypeExclusive"`
	IsLongRunning       bool   `json:"isLongRunning"`
	Name                string `json:"name,omitempty"`
	Trigger             string `json:"trigger,omitempty"`
	SuppressMessages    bool   `json:"suppressMessages"`
}

// CalendarEntry represents a calendar entry.
type CalendarEntry struct {
	ID                  int        `json:"id"`
	Title               string     `json:"title"`
	SortTitle           string     `json:"sortTitle,omitempty"`
	Year                int        `json:"year"`
	TMDbID              int        `json:"tmdbId"`
	IMDbID              string     `json:"imdbId,omitempty"`
	InCinemas           *time.Time `json:"inCinemas,omitempty"`
	PhysicalRelease     *time.Time `json:"physicalRelease,omitempty"`
	DigitalRelease      *time.Time `json:"digitalRelease,omitempty"`
	Status              string     `json:"status,omitempty"`
	Overview            string     `json:"overview,omitempty"`
	Images              []Image    `json:"images,omitempty"`
	Monitored           bool       `json:"monitored"`
	MinimumAvailability string     `json:"minimumAvailability,omitempty"`
	HasFile             bool       `json:"hasFile"`
}

// HistoryRecord represents a history record.
type HistoryRecord struct {
	ID                  int        `json:"id"`
	MovieID             int        `json:"movieId"`
	SourceTitle         string     `json:"sourceTitle,omitempty"`
	Languages           []Language `json:"languages,omitempty"`
	Quality             Quality    `json:"quality,omitempty"`
	CustomFormats       []any      `json:"customFormats,omitempty"`
	CustomFormatScore   int        `json:"customFormatScore,omitempty"`
	QualityCutoffNotMet bool       `json:"qualityCutoffNotMet"`
	Date                time.Time  `json:"date"`
	DownloadID          string     `json:"downloadId,omitempty"`
	EventType           string     `json:"eventType,omitempty"`
	Data                map[string]any `json:"data,omitempty"`
}

// HistoryResponse represents the paginated history response.
type HistoryResponse struct {
	Page          int             `json:"page"`
	PageSize      int             `json:"pageSize"`
	SortKey       string          `json:"sortKey,omitempty"`
	SortDirection string          `json:"sortDirection,omitempty"`
	TotalRecords  int             `json:"totalRecords"`
	Records       []HistoryRecord `json:"records,omitempty"`
}

// AddMovieRequest represents a request to add a movie.
type AddMovieRequest struct {
	Title               string   `json:"title"`
	QualityProfileID    int      `json:"qualityProfileId"`
	TMDbID              int      `json:"tmdbId"`
	RootFolderPath      string   `json:"rootFolderPath"`
	Monitored           bool     `json:"monitored"`
	MinimumAvailability string   `json:"minimumAvailability,omitempty"`
	Tags                []int    `json:"tags,omitempty"`
	AddOptions          AddOptions `json:"addOptions,omitempty"`
}

// AddOptions represents options when adding a movie.
type AddOptions struct {
	IgnoreEpisodesWithFiles    bool   `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles bool   `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                    string `json:"monitor,omitempty"`
	SearchForMovie             bool   `json:"searchForMovie"`
	AddMethod                  string `json:"addMethod,omitempty"`
}

// WebhookPayload represents a webhook payload from Radarr.
type WebhookPayload struct {
	EventType                string         `json:"eventType"`
	InstanceName             string         `json:"instanceName,omitempty"`
	ApplicationURL           string         `json:"applicationUrl,omitempty"`
	Movie                    *WebhookMovie  `json:"movie,omitempty"`
	RemoteMovie              *WebhookRemoteMovie `json:"remoteMovie,omitempty"`
	MovieFile                *WebhookMovieFile   `json:"movieFile,omitempty"`
	DeletedFiles             []WebhookMovieFile  `json:"deletedFiles,omitempty"`
	IsUpgrade                bool           `json:"isUpgrade,omitempty"`
	DownloadClient           string         `json:"downloadClient,omitempty"`
	DownloadClientType       string         `json:"downloadClientType,omitempty"`
	DownloadID               string         `json:"downloadId,omitempty"`
	Release                  *WebhookRelease     `json:"release,omitempty"`
	// Health event fields
	Level                    string         `json:"level,omitempty"`
	Message                  string         `json:"message,omitempty"`
	Type                     string         `json:"type,omitempty"`
	WikiURL                  string         `json:"wikiUrl,omitempty"`
	// Application update fields
	PreviousVersion          string         `json:"previousVersion,omitempty"`
	NewVersion               string         `json:"newVersion,omitempty"`
}

// WebhookMovie represents movie info in a webhook.
type WebhookMovie struct {
	ID           int    `json:"id"`
	Title        string `json:"title,omitempty"`
	Year         int    `json:"year,omitempty"`
	FilePath     string `json:"filePath,omitempty"`
	ReleaseDate  string `json:"releaseDate,omitempty"`
	FolderPath   string `json:"folderPath,omitempty"`
	TMDbID       int    `json:"tmdbId,omitempty"`
	IMDbID       string `json:"imdbId,omitempty"`
	Overview     string `json:"overview,omitempty"`
}

// WebhookRemoteMovie represents remote movie info.
type WebhookRemoteMovie struct {
	TMDbID int    `json:"tmdbId,omitempty"`
	IMDbID string `json:"imdbId,omitempty"`
	Title  string `json:"title,omitempty"`
	Year   int    `json:"year,omitempty"`
}

// WebhookMovieFile represents movie file info in a webhook.
type WebhookMovieFile struct {
	ID             int        `json:"id,omitempty"`
	RelativePath   string     `json:"relativePath,omitempty"`
	Path           string     `json:"path,omitempty"`
	Quality        string     `json:"quality,omitempty"`
	QualityVersion int        `json:"qualityVersion,omitempty"`
	ReleaseGroup   string     `json:"releaseGroup,omitempty"`
	SceneName      string     `json:"sceneName,omitempty"`
	Size           int64      `json:"size,omitempty"`
	DateAdded      *time.Time `json:"dateAdded,omitempty"`
	MediaInfo      *MediaInfo `json:"mediaInfo,omitempty"`
}

// Webhook event types.
const (
	EventGrab            = "Grab"
	EventDownload        = "Download"
	EventRename          = "Rename"
	EventMovieDelete     = "MovieDelete"
	EventMovieFileDelete = "MovieFileDelete"
	EventHealth          = "Health"
	EventHealthRestored  = "HealthRestored"
	EventApplicationUpdate = "ApplicationUpdate"
	EventManualInteractionRequired = "ManualInteractionRequired"
	EventTest            = "Test"
)
