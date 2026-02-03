// Package radarr provides a client for the Radarr API v3.
// Radarr is a PRIMARY metadata provider for movies in the Revenge media system.
package radarr

import "time"

// Movie represents a movie in Radarr.
type Movie struct {
	ID                    int              `json:"id"`
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

// MediaInfo contains technical information about the media file.
type MediaInfo struct {
	AudioBitrate          int      `json:"audioBitrate,omitempty"`
	AudioChannels         float64  `json:"audioChannels,omitempty"`
	AudioCodec            string   `json:"audioCodec,omitempty"`
	AudioLanguages        string   `json:"audioLanguages,omitempty"`
	AudioStreamCount      int      `json:"audioStreamCount,omitempty"`
	VideoBitDepth         int      `json:"videoBitDepth,omitempty"`
	VideoBitrate          int      `json:"videoBitrate,omitempty"`
	VideoCodec            string   `json:"videoCodec,omitempty"`
	VideoFps              float64  `json:"videoFps,omitempty"`
	VideoDynamicRange     string   `json:"videoDynamicRange,omitempty"`
	VideoDynamicRangeType string   `json:"videoDynamicRangeType,omitempty"`
	Resolution            string   `json:"resolution,omitempty"`
	RunTime               string   `json:"runTime,omitempty"`
	ScanType              string   `json:"scanType,omitempty"`
	Subtitles             string   `json:"subtitles,omitempty"`
}

// Quality represents quality information.
type Quality struct {
	Quality  QualityInfo  `json:"quality,omitempty"`
	Revision QualityRevision `json:"revision,omitempty"`
}

// QualityInfo contains quality details.
type QualityInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Source     string `json:"source,omitempty"`
	Resolution int    `json:"resolution,omitempty"`
	Modifier   string `json:"modifier,omitempty"`
}

// QualityRevision contains quality revision info.
type QualityRevision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// Language represents a language.
type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// AlternateTitle represents an alternate title.
type AlternateTitle struct {
	SourceType      string `json:"sourceType,omitempty"`
	MovieMetadataID int    `json:"movieMetadataId,omitempty"`
	Title           string `json:"title,omitempty"`
	CleanTitle      string `json:"cleanTitle,omitempty"`
}

// Image represents a movie image.
type Image struct {
	CoverType string `json:"coverType,omitempty"`
	URL       string `json:"url,omitempty"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// Ratings contains rating information.
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

// QualityProfile represents a quality profile.
type QualityProfile struct {
	ID                    int              `json:"id"`
	Name                  string           `json:"name"`
	UpgradeAllowed        bool             `json:"upgradeAllowed"`
	Cutoff                int              `json:"cutoff"`
	Items                 []QualityItem    `json:"items,omitempty"`
	MinFormatScore        int              `json:"minFormatScore"`
	CutoffFormatScore     int              `json:"cutoffFormatScore"`
	FormatItems           []FormatItem     `json:"formatItems,omitempty"`
	Language              *Language        `json:"language,omitempty"`
}

// QualityItem represents a quality item in a profile.
type QualityItem struct {
	ID      int          `json:"id,omitempty"`
	Name    string       `json:"name,omitempty"`
	Quality *QualityInfo `json:"quality,omitempty"`
	Items   []QualityItem `json:"items,omitempty"`
	Allowed bool         `json:"allowed"`
}

// FormatItem represents a custom format item.
type FormatItem struct {
	Format int `json:"format"`
	Name   string `json:"name,omitempty"`
	Score  int    `json:"score"`
}

// RootFolder represents a root folder.
type RootFolder struct {
	ID              int    `json:"id"`
	Path            string `json:"path"`
	Accessible      bool   `json:"accessible"`
	FreeSpace       int64  `json:"freeSpace,omitempty"`
	UnmappedFolders []UnmappedFolder `json:"unmappedFolders,omitempty"`
}

// UnmappedFolder represents an unmapped folder.
type UnmappedFolder struct {
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	RelativePath string `json:"relativePath,omitempty"`
}

// SystemStatus represents Radarr system status.
type SystemStatus struct {
	AppName                string `json:"appName,omitempty"`
	InstanceName           string `json:"instanceName,omitempty"`
	Version                string `json:"version,omitempty"`
	BuildTime              string `json:"buildTime,omitempty"`
	IsDebug                bool   `json:"isDebug"`
	IsProduction           bool   `json:"isProduction"`
	IsAdmin                bool   `json:"isAdmin"`
	IsUserInteractive      bool   `json:"isUserInteractive"`
	StartupPath            string `json:"startupPath,omitempty"`
	AppData                string `json:"appData,omitempty"`
	OsName                 string `json:"osName,omitempty"`
	OsVersion              string `json:"osVersion,omitempty"`
	IsDocker               bool   `json:"isDocker"`
	IsLinux                bool   `json:"isLinux"`
	IsOsx                  bool   `json:"isOsx"`
	IsWindows              bool   `json:"isWindows"`
	IsNetCore              bool   `json:"isNetCore"`
	RuntimeVersion         string `json:"runtimeVersion,omitempty"`
	RuntimeName            string `json:"runtimeName,omitempty"`
	StartTime              string `json:"startTime,omitempty"`
	PackageVersion         string `json:"packageVersion,omitempty"`
	PackageAuthor          string `json:"packageAuthor,omitempty"`
	PackageUpdateMechanism string `json:"packageUpdateMechanism,omitempty"`
	Branch                 string `json:"branch,omitempty"`
	Authentication         string `json:"authentication,omitempty"`
	SqliteVersion          string `json:"sqliteVersion,omitempty"`
	MigrationVersion       int    `json:"migrationVersion,omitempty"`
	URLBase                string `json:"urlBase,omitempty"`
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

// Tag represents a tag.
type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
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

// WebhookRelease represents release info in a webhook.
type WebhookRelease struct {
	Quality          string   `json:"quality,omitempty"`
	QualityVersion   int      `json:"qualityVersion,omitempty"`
	ReleaseGroup     string   `json:"releaseGroup,omitempty"`
	ReleaseTitle     string   `json:"releaseTitle,omitempty"`
	Indexer          string   `json:"indexer,omitempty"`
	Size             int64    `json:"size,omitempty"`
	CustomFormatScore int     `json:"customFormatScore,omitempty"`
	CustomFormats    []string `json:"customFormats,omitempty"`
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
