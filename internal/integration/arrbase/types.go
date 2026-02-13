package arrbase

// ============================================================================
// Shared types — identical between Radarr and Sonarr API v3 responses.
// These types are deserialized from JSON and MUST match the arr API contracts.
// ============================================================================

// MediaInfo contains technical information about a media file.
// Shared by Radarr MovieFile and Sonarr EpisodeFile.
type MediaInfo struct {
	AudioBitrate          int     `json:"audioBitrate,omitempty"`
	AudioChannels         float64 `json:"audioChannels,omitempty"`
	AudioCodec            string  `json:"audioCodec,omitempty"`
	AudioLanguages        string  `json:"audioLanguages,omitempty"`
	AudioStreamCount      int     `json:"audioStreamCount,omitempty"`
	VideoBitDepth         int     `json:"videoBitDepth,omitempty"`
	VideoBitrate          int     `json:"videoBitrate,omitempty"`
	VideoCodec            string  `json:"videoCodec,omitempty"`
	VideoFps              float64 `json:"videoFps,omitempty"`
	VideoDynamicRange     string  `json:"videoDynamicRange,omitempty"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType,omitempty"`
	Resolution            string  `json:"resolution,omitempty"`
	RunTime               string  `json:"runTime,omitempty"`
	ScanType              string  `json:"scanType,omitempty"`
	Subtitles             string  `json:"subtitles,omitempty"`
}

// Quality represents quality information (profile + revision).
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
	Modifier   string `json:"modifier,omitempty"` // Radarr uses this; Sonarr ignores it (zero value)
}

// QualityRevision contains quality revision info.
type QualityRevision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// Language represents a language reference.
type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// Image represents a media image (poster, banner, fanart).
type Image struct {
	CoverType string `json:"coverType,omitempty"`
	URL       string `json:"url,omitempty"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// QualityItem represents a quality item in a profile.
type QualityItem struct {
	ID      int           `json:"id,omitempty"`
	Name    string        `json:"name,omitempty"`
	Quality *QualityInfo  `json:"quality,omitempty"`
	Items   []QualityItem `json:"items,omitempty"`
	Allowed bool          `json:"allowed"`
}

// FormatItem represents a custom format item in a quality profile.
type FormatItem struct {
	Format int    `json:"format"`
	Name   string `json:"name,omitempty"`
	Score  int    `json:"score"`
}

// QualityProfile represents a quality profile configuration.
// NOTE: Radarr has an extra Language field; Sonarr omits it (nil).
type QualityProfile struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	UpgradeAllowed    bool          `json:"upgradeAllowed"`
	Cutoff            int           `json:"cutoff"`
	Items             []QualityItem `json:"items,omitempty"`
	MinFormatScore    int           `json:"minFormatScore"`
	CutoffFormatScore int           `json:"cutoffFormatScore"`
	FormatItems       []FormatItem  `json:"formatItems,omitempty"`
	Language          *Language     `json:"language,omitempty"` // Radarr only; nil for Sonarr
}

// RootFolder represents a root folder in the arr system.
type RootFolder struct {
	ID              int              `json:"id"`
	Path            string           `json:"path"`
	Accessible      bool             `json:"accessible"`
	FreeSpace       int64            `json:"freeSpace,omitempty"`
	UnmappedFolders []UnmappedFolder `json:"unmappedFolders,omitempty"`
}

// UnmappedFolder represents an unmapped folder within a root folder.
type UnmappedFolder struct {
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	RelativePath string `json:"relativePath,omitempty"`
}

// SystemStatus represents the system status from an arr instance.
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

// Tag represents a tag in the arr system.
type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// WebhookRelease represents release info in a webhook payload.
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

// Shared webhook event types common to all arr applications.
const (
	EventGrab                       = "Grab"
	EventDownload                   = "Download"
	EventRename                     = "Rename"
	EventHealth                     = "Health"
	EventHealthRestored             = "HealthRestored"
	EventApplicationUpdate          = "ApplicationUpdate"
	EventManualInteractionRequired  = "ManualInteractionRequired"
	EventTest                       = "Test"
)
