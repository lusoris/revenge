// Package radarr provides a Radarr API v3 client.
package radarr

import "time"

// Movie represents a Radarr movie.
type Movie struct {
	ID                    int           `json:"id"`
	Title                 string        `json:"title"`
	OriginalTitle         string        `json:"originalTitle"`
	OriginalLanguage      Language      `json:"originalLanguage"`
	SortTitle             string        `json:"sortTitle"`
	Overview              string        `json:"overview"`
	Year                  int           `json:"year"`
	Path                  string        `json:"path"`
	QualityProfileID      int           `json:"qualityProfileId"`
	Monitored             bool          `json:"monitored"`
	MinimumAvailability   string        `json:"minimumAvailability"`
	IsAvailable           bool          `json:"isAvailable"`
	FolderName            string        `json:"folderName"`
	Runtime               int           `json:"runtime"` // minutes
	CleanTitle            string        `json:"cleanTitle"`
	IMDbID                string        `json:"imdbId"`
	TMDbID                int           `json:"tmdbId"`
	TitleSlug             string        `json:"titleSlug"`
	Certification         string        `json:"certification"`
	Genres                []string      `json:"genres"`
	Tags                  []int         `json:"tags"`
	Added                 time.Time     `json:"added"`
	Ratings               Ratings       `json:"ratings"`
	MovieFile             *MovieFile    `json:"movieFile,omitempty"`
	Collection            *Collection   `json:"collection,omitempty"`
	HasFile               bool          `json:"hasFile"`
	Status                string        `json:"status"`
	Studio                string        `json:"studio"`
	Images                []Image       `json:"images"`
	Website               string        `json:"website"`
	YouTubeTrailerID      string        `json:"youTubeTrailerId"`
	InCinemas             *time.Time    `json:"inCinemas,omitempty"`
	PhysicalRelease       *time.Time    `json:"physicalRelease,omitempty"`
	DigitalRelease        *time.Time    `json:"digitalRelease,omitempty"`
	AlternativeTitles     []AltTitle    `json:"alternateTitles"`
	SecondaryYearSourceID int           `json:"secondaryYearSourceId"`
	SizeOnDisk            int64         `json:"sizeOnDisk"`
	RootFolderPath        string        `json:"rootFolderPath"`
	Popularity            float64       `json:"popularity"`
}

// Language represents a language in Radarr.
type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Ratings contains rating information.
type Ratings struct {
	IMDb           *Rating `json:"imdb,omitempty"`
	TMDb           *Rating `json:"tmdb,omitempty"`
	Metacritic     *Rating `json:"metacritic,omitempty"`
	RottenTomatoes *Rating `json:"rottenTomatoes,omitempty"`
}

// Rating represents a single rating.
type Rating struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// MovieFile represents a movie file on disk.
type MovieFile struct {
	ID                  int           `json:"id"`
	MovieID             int           `json:"movieId"`
	RelativePath        string        `json:"relativePath"`
	Path                string        `json:"path"`
	Size                int64         `json:"size"`
	DateAdded           time.Time     `json:"dateAdded"`
	SceneName           string        `json:"sceneName"`
	IndexerFlags        int           `json:"indexerFlags"`
	Quality             Quality       `json:"quality"`
	MediaInfo           *MediaInfo    `json:"mediaInfo,omitempty"`
	OriginalFilePath    string        `json:"originalFilePath"`
	QualityCutoffNotMet bool          `json:"qualityCutoffNotMet"`
	Languages           []Language    `json:"languages"`
	ReleaseGroup        string        `json:"releaseGroup"`
	Edition             string        `json:"edition"`
}

// Quality represents quality information.
type Quality struct {
	Quality  QualityDetail `json:"quality"`
	Revision Revision      `json:"revision"`
}

// QualityDetail contains quality details.
type QualityDetail struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution int    `json:"resolution"`
	Modifier   string `json:"modifier"`
}

// Revision represents a quality revision.
type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// MediaInfo contains media file details.
type MediaInfo struct {
	AudioBitrate          int     `json:"audioBitrate"`
	AudioChannels         float64 `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	AudioLanguages        string  `json:"audioLanguages"`
	AudioStreamCount      int     `json:"audioStreamCount"`
	VideoBitDepth         int     `json:"videoBitDepth"`
	VideoBitrate          int     `json:"videoBitrate"`
	VideoCodec            string  `json:"videoCodec"`
	VideoFps              float64 `json:"videoFps"`
	VideoDynamicRange     string  `json:"videoDynamicRange"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
	Resolution            string  `json:"resolution"`
	RunTime               string  `json:"runTime"`
	ScanType              string  `json:"scanType"`
	Subtitles             string  `json:"subtitles"`
}

// Collection represents a movie collection.
type Collection struct {
	Name   string  `json:"name"`
	TMDbID int     `json:"tmdbId"`
	Images []Image `json:"images"`
}

// Image represents an image.
type Image struct {
	CoverType string `json:"coverType"` // poster, fanart, banner
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

// AltTitle represents an alternative title.
type AltTitle struct {
	SourceType string `json:"sourceType"`
	MovieID    int    `json:"movieId"`
	Title      string `json:"title"`
	SourceID   int    `json:"sourceId"`
	Votes      int    `json:"votes"`
	VoteCount  int    `json:"voteCount"`
	Language   Language `json:"language"`
}

// RootFolder represents a root folder configuration.
type RootFolder struct {
	ID              int    `json:"id"`
	Path            string `json:"path"`
	Accessible      bool   `json:"accessible"`
	FreeSpace       int64  `json:"freeSpace"`
	UnmappedFolders []any  `json:"unmappedFolders"`
}

// QualityProfile represents a quality profile.
type QualityProfile struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	UpgradeAllowed    bool              `json:"upgradeAllowed"`
	Cutoff            int               `json:"cutoff"`
	Items             []QualityItem     `json:"items"`
	MinFormatScore    int               `json:"minFormatScore"`
	CutoffFormatScore int               `json:"cutoffFormatScore"`
	FormatItems       []any             `json:"formatItems"`
	Language          Language          `json:"language"`
}

// QualityItem represents a quality item in a profile.
type QualityItem struct {
	Quality *QualityDetail `json:"quality,omitempty"`
	Items   []QualityItem  `json:"items,omitempty"`
	Allowed bool           `json:"allowed"`
	Name    string         `json:"name,omitempty"`
	ID      int            `json:"id,omitempty"`
}

// SystemStatus represents Radarr system status.
type SystemStatus struct {
	AppName                string    `json:"appName"`
	InstanceName           string    `json:"instanceName"`
	Version                string    `json:"version"`
	BuildTime              time.Time `json:"buildTime"`
	IsDebug                bool      `json:"isDebug"`
	IsProduction           bool      `json:"isProduction"`
	IsAdmin                bool      `json:"isAdmin"`
	IsUserInteractive      bool      `json:"isUserInteractive"`
	StartupPath            string    `json:"startupPath"`
	AppData                string    `json:"appData"`
	OsName                 string    `json:"osName"`
	OsVersion              string    `json:"osVersion"`
	IsNetCore              bool      `json:"isNetCore"`
	IsLinux                bool      `json:"isLinux"`
	IsOsx                  bool      `json:"isOsx"`
	IsWindows              bool      `json:"isWindows"`
	IsDocker               bool      `json:"isDocker"`
	Mode                   string    `json:"mode"`
	Branch                 string    `json:"branch"`
	Authentication         string    `json:"authentication"`
	SqliteVersion          string    `json:"sqliteVersion"`
	MigrationVersion       int       `json:"migrationVersion"`
	URLBase                string    `json:"urlBase"`
	RuntimeVersion         string    `json:"runtimeVersion"`
	RuntimeName            string    `json:"runtimeName"`
	StartTime              time.Time `json:"startTime"`
	PackageVersion         string    `json:"packageVersion"`
	PackageAuthor          string    `json:"packageAuthor"`
	PackageUpdateMechanism string    `json:"packageUpdateMechanism"`
}

// HealthCheck represents a health check result.
type HealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}

// WebhookPayload represents an incoming webhook from Radarr.
type WebhookPayload struct {
	EventType                string        `json:"eventType"`
	InstanceName             string        `json:"instanceName"`
	ApplicationURL           string        `json:"applicationUrl"`
	Movie                    *WebhookMovie `json:"movie,omitempty"`
	RemoteMovie              *WebhookMovie `json:"remoteMovie,omitempty"`
	MovieFile                *MovieFile    `json:"movieFile,omitempty"`
	DeletedFiles             []MovieFile   `json:"deletedFiles,omitempty"`
	IsUpgrade                bool          `json:"isUpgrade"`
	DownloadClient           string        `json:"downloadClient,omitempty"`
	DownloadClientType       string        `json:"downloadClientType,omitempty"`
	DownloadID               string        `json:"downloadId,omitempty"`
	CustomFormatInfo         any           `json:"customFormatInfo,omitempty"`
	Release                  any           `json:"release,omitempty"`
	RenamedMovieFiles        []MovieFile   `json:"renamedMovieFiles,omitempty"`
	AddMethod                string        `json:"addMethod,omitempty"`
}

// WebhookMovie represents movie data in a webhook payload.
type WebhookMovie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	FolderPath  string   `json:"folderPath"`
	ReleaseDate string   `json:"releaseDate"`
	TMDbID      int      `json:"tmdbId"`
	IMDbID      string   `json:"imdbId"`
	Overview    string   `json:"overview"`
	Genres      []string `json:"genres"`
	Images      []Image  `json:"images"`
}

// AddMovieOptions represents options for adding a movie.
type AddMovieOptions struct {
	Title               string `json:"title"`
	TMDbID              int    `json:"tmdbId"`
	QualityProfileID    int    `json:"qualityProfileId"`
	RootFolderPath      string `json:"rootFolderPath"`
	MinimumAvailability string `json:"minimumAvailability"`
	Monitored           bool   `json:"monitored"`
	AddOptions          struct {
		SearchForMovie bool `json:"searchForMovie"`
	} `json:"addOptions"`
}

// APIError represents a Radarr API error.
type APIError struct {
	Message            string `json:"message"`
	Description        string `json:"description"`
	PropertyName       string `json:"propertyName"`
	AttemptedValue     any    `json:"attemptedValue"`
	Severity           string `json:"severity"`
	ErrorCode          string `json:"errorCode"`
	InfoLink           string `json:"infoLink"`
	IsWarning          bool   `json:"isWarning"`
}

// Error implements error interface.
func (e APIError) Error() string {
	if e.Description != "" {
		return e.Description
	}
	return e.Message
}
