// Package sonarr provides a client for the Sonarr API v3.
// Sonarr is a PRIMARY metadata provider for TV shows in the Revenge media system.
package sonarr

import "time"

// Series represents a TV series in Sonarr.
type Series struct {
	ID                    int              `json:"id"`
	Title                 string           `json:"title"`
	AlternateTitles       []AlternateTitle `json:"alternateTitles,omitempty"`
	SortTitle             string           `json:"sortTitle,omitempty"`
	Status                string           `json:"status,omitempty"`
	Ended                 bool             `json:"ended"`
	Overview              string           `json:"overview,omitempty"`
	PreviousAiring        *time.Time       `json:"previousAiring,omitempty"`
	NextAiring            *time.Time       `json:"nextAiring,omitempty"`
	Network               string           `json:"network,omitempty"`
	AirTime               string           `json:"airTime,omitempty"`
	Images                []Image          `json:"images,omitempty"`
	OriginalLanguage      Language         `json:"originalLanguage,omitempty"`
	Seasons               []SeasonInfo     `json:"seasons,omitempty"`
	Year                  int              `json:"year"`
	Path                  string           `json:"path,omitempty"`
	QualityProfileID      int              `json:"qualityProfileId"`
	SeasonFolder          bool             `json:"seasonFolder"`
	Monitored             bool             `json:"monitored"`
	MonitorNewItems       string           `json:"monitorNewItems,omitempty"`
	UseSceneNumbering     bool             `json:"useSceneNumbering"`
	Runtime               int              `json:"runtime,omitempty"`
	TVDbID                int              `json:"tvdbId"`
	TVRageID              int              `json:"tvRageId,omitempty"`
	TVMazeID              int              `json:"tvMazeId,omitempty"`
	IMDbID                string           `json:"imdbId,omitempty"`
	FirstAired            *time.Time       `json:"firstAired,omitempty"`
	LastAired             *time.Time       `json:"lastAired,omitempty"`
	SeriesType            string           `json:"seriesType,omitempty"`
	CleanTitle            string           `json:"cleanTitle,omitempty"`
	TitleSlug             string           `json:"titleSlug,omitempty"`
	RootFolderPath        string           `json:"rootFolderPath,omitempty"`
	Certification         string           `json:"certification,omitempty"`
	Genres                []string         `json:"genres,omitempty"`
	Tags                  []int            `json:"tags,omitempty"`
	Added                 time.Time        `json:"added"`
	Ratings               Ratings          `json:"ratings,omitempty"`
	Statistics            *Statistics      `json:"statistics,omitempty"`
	LanguageProfileID     int              `json:"languageProfileId,omitempty"`
}

// SeasonInfo represents season information within a series.
type SeasonInfo struct {
	SeasonNumber int         `json:"seasonNumber"`
	Monitored    bool        `json:"monitored"`
	Statistics   *Statistics `json:"statistics,omitempty"`
}

// Episode represents an episode in Sonarr.
type Episode struct {
	ID                       int        `json:"id"`
	SeriesID                 int        `json:"seriesId"`
	TVDbID                   int        `json:"tvdbId,omitempty"`
	EpisodeFileID            int        `json:"episodeFileId,omitempty"`
	SeasonNumber             int        `json:"seasonNumber"`
	EpisodeNumber            int        `json:"episodeNumber"`
	Title                    string     `json:"title,omitempty"`
	AirDate                  string     `json:"airDate,omitempty"`
	AirDateUtc               *time.Time `json:"airDateUtc,omitempty"`
	Overview                 string     `json:"overview,omitempty"`
	EpisodeFile              *EpisodeFile `json:"episodeFile,omitempty"`
	HasFile                  bool       `json:"hasFile"`
	Monitored                bool       `json:"monitored"`
	AbsoluteEpisodeNumber    int        `json:"absoluteEpisodeNumber,omitempty"`
	SceneAbsoluteEpisodeNumber int      `json:"sceneAbsoluteEpisodeNumber,omitempty"`
	SceneSeasonNumber        int        `json:"sceneSeasonNumber,omitempty"`
	SceneEpisodeNumber       int        `json:"sceneEpisodeNumber,omitempty"`
	UnverifiedSceneNumbering bool       `json:"unverifiedSceneNumbering,omitempty"`
	Runtime                  int        `json:"runtime,omitempty"`
	FinaleType               string     `json:"finaleType,omitempty"`
	Images                   []Image    `json:"images,omitempty"`
	Series                   *Series    `json:"series,omitempty"`
}

// EpisodeFile represents an episode file in Sonarr.
type EpisodeFile struct {
	ID                  int            `json:"id"`
	SeriesID            int            `json:"seriesId"`
	SeasonNumber        int            `json:"seasonNumber"`
	RelativePath        string         `json:"relativePath,omitempty"`
	Path                string         `json:"path,omitempty"`
	Size                int64          `json:"size"`
	DateAdded           time.Time      `json:"dateAdded"`
	SceneName           string         `json:"sceneName,omitempty"`
	ReleaseGroup        string         `json:"releaseGroup,omitempty"`
	Quality             Quality        `json:"quality,omitempty"`
	MediaInfo           *MediaInfo     `json:"mediaInfo,omitempty"`
	OriginalFilePath    string         `json:"originalFilePath,omitempty"`
	QualityCutoffNotMet bool           `json:"qualityCutoffNotMet"`
	Languages           []Language     `json:"languages,omitempty"`
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
	Quality  QualityInfo     `json:"quality,omitempty"`
	Revision QualityRevision `json:"revision,omitempty"`
}

// QualityInfo contains quality details.
type QualityInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Source     string `json:"source,omitempty"`
	Resolution int    `json:"resolution,omitempty"`
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
	Title           string `json:"title,omitempty"`
	SeasonNumber    int    `json:"seasonNumber,omitempty"`
	SceneSeasonNumber int  `json:"sceneSeasonNumber,omitempty"`
	SceneOrigin     string `json:"sceneOrigin,omitempty"`
	Comment         string `json:"comment,omitempty"`
}

// Image represents a series/episode image.
type Image struct {
	CoverType string `json:"coverType,omitempty"`
	URL       string `json:"url,omitempty"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// Ratings contains rating information.
type Ratings struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
}

// Statistics contains series/season statistics.
type Statistics struct {
	SeasonCount       int     `json:"seasonCount,omitempty"`
	EpisodeFileCount  int     `json:"episodeFileCount"`
	EpisodeCount      int     `json:"episodeCount"`
	TotalEpisodeCount int     `json:"totalEpisodeCount"`
	SizeOnDisk        int64   `json:"sizeOnDisk"`
	ReleaseGroups     []string `json:"releaseGroups,omitempty"`
	PercentOfEpisodes float64 `json:"percentOfEpisodes,omitempty"`
}

// QualityProfile represents a quality profile.
type QualityProfile struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	UpgradeAllowed    bool          `json:"upgradeAllowed"`
	Cutoff            int           `json:"cutoff"`
	Items             []QualityItem `json:"items,omitempty"`
	MinFormatScore    int           `json:"minFormatScore"`
	CutoffFormatScore int           `json:"cutoffFormatScore"`
	FormatItems       []FormatItem  `json:"formatItems,omitempty"`
}

// QualityItem represents a quality item in a profile.
type QualityItem struct {
	ID      int           `json:"id,omitempty"`
	Name    string        `json:"name,omitempty"`
	Quality *QualityInfo  `json:"quality,omitempty"`
	Items   []QualityItem `json:"items,omitempty"`
	Allowed bool          `json:"allowed"`
}

// FormatItem represents a custom format item.
type FormatItem struct {
	Format int    `json:"format"`
	Name   string `json:"name,omitempty"`
	Score  int    `json:"score"`
}

// RootFolder represents a root folder.
type RootFolder struct {
	ID              int              `json:"id"`
	Path            string           `json:"path"`
	Accessible      bool             `json:"accessible"`
	FreeSpace       int64            `json:"freeSpace,omitempty"`
	UnmappedFolders []UnmappedFolder `json:"unmappedFolders,omitempty"`
}

// UnmappedFolder represents an unmapped folder.
type UnmappedFolder struct {
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	RelativePath string `json:"relativePath,omitempty"`
}

// SystemStatus represents Sonarr system status.
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

// Command represents a command to execute in Sonarr.
type Command struct {
	ID                  int         `json:"id,omitempty"`
	Name                string      `json:"name"`
	CommandName         string      `json:"commandName,omitempty"`
	Message             string      `json:"message,omitempty"`
	Body                CommandBody `json:"body,omitempty"`
	Priority            string      `json:"priority,omitempty"`
	Status              string      `json:"status,omitempty"`
	Queued              time.Time   `json:"queued,omitempty"`
	Started             *time.Time  `json:"started,omitempty"`
	Ended               *time.Time  `json:"ended,omitempty"`
	Duration            string      `json:"duration,omitempty"`
	Trigger             string      `json:"trigger,omitempty"`
	ClientUserAgent     string      `json:"clientUserAgent,omitempty"`
	StateChangeTime     *time.Time  `json:"stateChangeTime,omitempty"`
	SendUpdatesToClient bool        `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool        `json:"updateScheduledTask"`
	LastExecutionTime   *time.Time  `json:"lastExecutionTime,omitempty"`
}

// CommandBody contains command body parameters.
type CommandBody struct {
	SeriesIDs           []int  `json:"seriesIds,omitempty"`
	SeriesID            int    `json:"seriesId,omitempty"`
	SeasonNumber        int    `json:"seasonNumber,omitempty"`
	EpisodeIDs          []int  `json:"episodeIds,omitempty"`
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

// CalendarEntry represents a calendar entry (upcoming episode).
type CalendarEntry struct {
	ID                    int        `json:"id"`
	SeriesID              int        `json:"seriesId"`
	TVDbID                int        `json:"tvdbId,omitempty"`
	EpisodeFileID         int        `json:"episodeFileId,omitempty"`
	SeasonNumber          int        `json:"seasonNumber"`
	EpisodeNumber         int        `json:"episodeNumber"`
	Title                 string     `json:"title,omitempty"`
	AirDate               string     `json:"airDate,omitempty"`
	AirDateUtc            *time.Time `json:"airDateUtc,omitempty"`
	Overview              string     `json:"overview,omitempty"`
	HasFile               bool       `json:"hasFile"`
	Monitored             bool       `json:"monitored"`
	AbsoluteEpisodeNumber int        `json:"absoluteEpisodeNumber,omitempty"`
	Series                *Series    `json:"series,omitempty"`
}

// HistoryRecord represents a history record.
type HistoryRecord struct {
	ID                  int        `json:"id"`
	EpisodeID           int        `json:"episodeId"`
	SeriesID            int        `json:"seriesId"`
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

// AddSeriesRequest represents a request to add a series.
type AddSeriesRequest struct {
	Title             string           `json:"title"`
	QualityProfileID  int              `json:"qualityProfileId"`
	TVDbID            int              `json:"tvdbId"`
	RootFolderPath    string           `json:"rootFolderPath"`
	Monitored         bool             `json:"monitored"`
	SeasonFolder      bool             `json:"seasonFolder"`
	SeriesType        string           `json:"seriesType,omitempty"`
	Tags              []int            `json:"tags,omitempty"`
	AddOptions        AddSeriesOptions `json:"addOptions,omitempty"`
	Seasons           []SeasonInfo     `json:"seasons,omitempty"`
	LanguageProfileID int              `json:"languageProfileId,omitempty"`
}

// AddSeriesOptions represents options when adding a series.
type AddSeriesOptions struct {
	IgnoreEpisodesWithFiles    bool   `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles bool   `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                    string `json:"monitor,omitempty"`
	SearchForMissingEpisodes   bool   `json:"searchForMissingEpisodes"`
	SearchForCutoffUnmetEpisodes bool `json:"searchForCutoffUnmetEpisodes,omitempty"`
}

// WebhookPayload represents a webhook payload from Sonarr.
type WebhookPayload struct {
	EventType          string               `json:"eventType"`
	InstanceName       string               `json:"instanceName,omitempty"`
	ApplicationURL     string               `json:"applicationUrl,omitempty"`
	Series             *WebhookSeries       `json:"series,omitempty"`
	Episodes           []WebhookEpisode     `json:"episodes,omitempty"`
	EpisodeFile        *WebhookEpisodeFile  `json:"episodeFile,omitempty"`
	DeletedFiles       []WebhookEpisodeFile `json:"deletedFiles,omitempty"`
	IsUpgrade          bool                 `json:"isUpgrade,omitempty"`
	DownloadClient     string               `json:"downloadClient,omitempty"`
	DownloadClientType string               `json:"downloadClientType,omitempty"`
	DownloadID         string               `json:"downloadId,omitempty"`
	Release            *WebhookRelease      `json:"release,omitempty"`
	// Health event fields
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
	WikiURL   string `json:"wikiUrl,omitempty"`
	// Application update fields
	PreviousVersion string `json:"previousVersion,omitempty"`
	NewVersion      string `json:"newVersion,omitempty"`
}

// WebhookSeries represents series info in a webhook.
type WebhookSeries struct {
	ID         int    `json:"id"`
	Title      string `json:"title,omitempty"`
	TitleSlug  string `json:"titleSlug,omitempty"`
	Path       string `json:"path,omitempty"`
	TVDbID     int    `json:"tvdbId,omitempty"`
	TVMazeID   int    `json:"tvMazeId,omitempty"`
	IMDbID     string `json:"imdbId,omitempty"`
	Type       string `json:"type,omitempty"`
	Year       int    `json:"year,omitempty"`
}

// WebhookEpisode represents episode info in a webhook.
type WebhookEpisode struct {
	ID            int    `json:"id"`
	EpisodeNumber int    `json:"episodeNumber"`
	SeasonNumber  int    `json:"seasonNumber"`
	Title         string `json:"title,omitempty"`
	AirDate       string `json:"airDate,omitempty"`
	AirDateUtc    string `json:"airDateUtc,omitempty"`
	Overview      string `json:"overview,omitempty"`
	SeriesID      int    `json:"seriesId,omitempty"`
}

// WebhookEpisodeFile represents episode file info in a webhook.
type WebhookEpisodeFile struct {
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
	Quality           string   `json:"quality,omitempty"`
	QualityVersion    int      `json:"qualityVersion,omitempty"`
	ReleaseGroup      string   `json:"releaseGroup,omitempty"`
	ReleaseTitle      string   `json:"releaseTitle,omitempty"`
	Indexer           string   `json:"indexer,omitempty"`
	Size              int64    `json:"size,omitempty"`
	CustomFormatScore int      `json:"customFormatScore,omitempty"`
	CustomFormats     []string `json:"customFormats,omitempty"`
}

// Webhook event types.
const (
	EventGrab             = "Grab"
	EventDownload         = "Download"
	EventRename           = "Rename"
	EventSeriesAdd        = "SeriesAdd"
	EventSeriesDelete     = "SeriesDelete"
	EventEpisodeFileDelete = "EpisodeFileDelete"
	EventHealth           = "Health"
	EventHealthRestored   = "HealthRestored"
	EventApplicationUpdate = "ApplicationUpdate"
	EventManualInteractionRequired = "ManualInteractionRequired"
	EventTest             = "Test"
)

// Series status values.
const (
	StatusContinuing = "continuing"
	StatusEnded      = "ended"
	StatusUpcoming   = "upcoming"
	StatusDeleted    = "deleted"
)

// Series type values.
const (
	SeriesTypeStandard = "standard"
	SeriesTypeDaily    = "daily"
	SeriesTypeAnime    = "anime"
)

// Monitor options for adding series.
const (
	MonitorAll          = "all"
	MonitorFuture       = "future"
	MonitorMissing      = "missing"
	MonitorExisting     = "existing"
	MonitorPilot        = "pilot"
	MonitorFirstSeason  = "firstSeason"
	MonitorLastSeason   = "lastSeason"
	MonitorMonitorSpecials = "monitorSpecials"
	MonitorUnmonitorSpecials = "unmonitorSpecials"
	MonitorNone         = "none"
)
