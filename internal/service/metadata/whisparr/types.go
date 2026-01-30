// Package whisparr provides a Whisparr v3 API client for adult content acquisition.
// Whisparr is a fork of Radarr specialized for adult content management.
package whisparr

import "time"

// SystemStatus represents Whisparr system status.
type SystemStatus struct {
	AppName                string `json:"appName"`
	InstanceName           string `json:"instanceName"`
	Version                string `json:"version"`
	BuildTime              string `json:"buildTime"`
	IsDebug                bool   `json:"isDebug"`
	IsProduction           bool   `json:"isProduction"`
	IsAdmin                bool   `json:"isAdmin"`
	IsUserInteractive      bool   `json:"isUserInteractive"`
	StartupPath            string `json:"startupPath"`
	AppData                string `json:"appData"`
	OsName                 string `json:"osName"`
	OsVersion              string `json:"osVersion"`
	IsMonoRuntime          bool   `json:"isMonoRuntime"`
	IsMono                 bool   `json:"isMono"`
	IsLinux                bool   `json:"isLinux"`
	IsOsx                  bool   `json:"isOsx"`
	IsWindows              bool   `json:"isWindows"`
	IsDocker               bool   `json:"isDocker"`
	Mode                   string `json:"mode"`
	Branch                 string `json:"branch"`
	DatabaseType           string `json:"databaseType"`
	DatabaseVersion        string `json:"databaseVersion"`
	Authentication         string `json:"authentication"`
	MigrationVersion       int    `json:"migrationVersion"`
	UrlBase                string `json:"urlBase"`
	RuntimeVersion         string `json:"runtimeVersion"`
	RuntimeName            string `json:"runtimeName"`
	StartTime              string `json:"startTime"`
	PackageVersion         string `json:"packageVersion"`
	PackageAuthor          string `json:"packageAuthor"`
	PackageUpdateMechanism string `json:"packageUpdateMechanism"`
}

// HealthCheck represents a health check result.
type HealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}

// Movie represents a Whisparr movie/scene item.
// Whisparr uses "movie" terminology for both full-length movies and scenes.
type Movie struct {
	ID                    int            `json:"id"`
	Title                 string         `json:"title"`
	OriginalTitle         string         `json:"originalTitle,omitempty"`
	OriginalLanguage      *Language      `json:"originalLanguage,omitempty"`
	SortTitle             string         `json:"sortTitle"`
	SizeOnDisk            int64          `json:"sizeOnDisk"`
	Overview              string         `json:"overview"`
	Status                string         `json:"status"`
	Images                []Image        `json:"images"`
	Downloaded            bool           `json:"downloaded"`
	Year                  int            `json:"year"`
	Path                  string         `json:"path"`
	QualityProfileID      int            `json:"qualityProfileId"`
	MovieFile             *MovieFile     `json:"movieFile,omitempty"`
	Studio                *Studio        `json:"studio,omitempty"`
	ForeignId             string         `json:"foreignId"` // External ID (StashDB/TPDB)
	StashID               string         `json:"stashId,omitempty"`
	Credits               []Credit       `json:"credits"`
	Genres                []string       `json:"genres"`
	Tags                  []int          `json:"tags"`
	Added                 time.Time      `json:"added"`
	ReleaseDate           string         `json:"releaseDate,omitempty"`
	Runtime               int            `json:"runtime"` // minutes
	Certification         string         `json:"certification,omitempty"`
	HasFile               bool           `json:"hasFile"`
	Monitored             bool           `json:"monitored"`
	IsAvailable           bool           `json:"isAvailable"`
	FolderName            string         `json:"folderName"`
	CleanTitle            string         `json:"cleanTitle"`
	TitleSlug             string         `json:"titleSlug"`
	RootFolderPath        string         `json:"rootFolderPath,omitempty"`
	MinimumAvailability   string         `json:"minimumAvailability"`
	AddOptions            *AddOptions    `json:"addOptions,omitempty"`
	Ratings               *Ratings       `json:"ratings,omitempty"`
	SecondaryYear         int            `json:"secondaryYear,omitempty"`
	SecondaryYearSourceID int            `json:"secondaryYearSourceId,omitempty"`
	ItemType              string         `json:"itemType"` // "movie" or "scene"
}

// Language represents a language reference.
type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Image represents a media image.
type Image struct {
	CoverType string `json:"coverType"` // poster, fanart, banner
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// MovieFile represents a downloaded file.
type MovieFile struct {
	ID                  int       `json:"id"`
	MovieID             int       `json:"movieId"`
	RelativePath        string    `json:"relativePath"`
	Path                string    `json:"path"`
	Size                int64     `json:"size"`
	DateAdded           time.Time `json:"dateAdded"`
	SceneName           string    `json:"sceneName,omitempty"`
	IndexerFlags        int       `json:"indexerFlags"`
	Quality             *Quality  `json:"quality"`
	MediaInfo           *MediaInfo `json:"mediaInfo,omitempty"`
	OriginalFilePath    string    `json:"originalFilePath,omitempty"`
	QualityCutoffNotMet bool      `json:"qualityCutoffNotMet"`
	Languages           []Language `json:"languages"`
	ReleaseGroup        string    `json:"releaseGroup,omitempty"`
	Edition             string    `json:"edition,omitempty"`
}

// Quality represents quality information.
type Quality struct {
	Quality  *QualityDetail `json:"quality"`
	Revision *Revision      `json:"revision"`
}

// QualityDetail contains quality specifics.
type QualityDetail struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution int    `json:"resolution"`
	Modifier   string `json:"modifier"`
}

// Revision represents quality revision info.
type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// MediaInfo contains technical media information.
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
	Resolution            string  `json:"resolution"`
	RunTime               string  `json:"runTime"`
	ScanType              string  `json:"scanType"`
	Subtitles             string  `json:"subtitles"`
	VideoDynamicRange     string  `json:"videoDynamicRange"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
}

// Studio represents an adult studio.
type Studio struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Name      string  `json:"name"`
	ForeignID string  `json:"foreignId"` // StashDB/TPDB studio ID
	Images    []Image `json:"images"`
}

// Credit represents a performer credit.
type Credit struct {
	PersonName    string  `json:"personName"`
	CreditID      int     `json:"creditId"`
	PersonID      int     `json:"personId"`
	MovieMetaID   int     `json:"movieMetadataId"`
	Character     string  `json:"character"`
	Department    string  `json:"department"`
	Job           string  `json:"job"`
	Type          string  `json:"type"` // "cast" or "crew"
	Images        []Image `json:"images"`
	Order         int     `json:"order"`
	PersonIMDbID  string  `json:"personImdbId,omitempty"`
	PersonTMDbID  int     `json:"personTmdbId,omitempty"`
	StashID       string  `json:"stashId,omitempty"` // StashDB performer ID
}

// AddOptions for adding new movies.
type AddOptions struct {
	SearchForMovie      bool `json:"searchForMovie"`
	IgnoreEpisodesWithFiles bool `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
	Monitor             string `json:"monitor"`
}

// Ratings contains rating information.
type Ratings struct {
	IMDB       *Rating `json:"imdb,omitempty"`
	TMDB       *Rating `json:"tmdb,omitempty"`
	Metacritic *Rating `json:"metacritic,omitempty"`
	Rottentomatoes *Rating `json:"rottenTomatoes,omitempty"`
}

// Rating represents a single rating source.
type Rating struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// RootFolder represents a Whisparr root folder.
type RootFolder struct {
	ID              int    `json:"id"`
	Path            string `json:"path"`
	Accessible      bool   `json:"accessible"`
	FreeSpace       int64  `json:"freeSpace"`
	UnmappedFolders []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"unmappedFolders"`
}

// QualityProfile represents a quality profile.
type QualityProfile struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	UpgradeAllowed    bool   `json:"upgradeAllowed"`
	Cutoff            int    `json:"cutoff"`
	MinFormatScore    int    `json:"minFormatScore"`
	CutoffFormatScore int    `json:"cutoffFormatScore"`
}

// Performer represents an adult performer from Whisparr.
// Maps to Crew in QAR terminology.
type Performer struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	ForeignID     string    `json:"foreignId"` // StashDB performer ID
	StashID       string    `json:"stashId,omitempty"`
	Overview      string    `json:"overview,omitempty"`
	Gender        string    `json:"gender,omitempty"`
	Images        []Image   `json:"images"`
	RootFolderPath string   `json:"rootFolderPath,omitempty"`
	Monitored     bool      `json:"monitored"`
	Added         time.Time `json:"added"`
	Statistics    *PerformerStats `json:"statistics,omitempty"`
}

// PerformerStats contains performer statistics.
type PerformerStats struct {
	MovieCount       int   `json:"movieCount"`
	MovieFileCount   int   `json:"movieFileCount"`
	SizeOnDisk       int64 `json:"sizeOnDisk"`
	PercentOfMovies  float64 `json:"percentOfMovies"`
}

// AddMovieOptions contains options for adding a movie.
type AddMovieOptions struct {
	Title               string      `json:"title"`
	ForeignID           string      `json:"foreignId"`
	QualityProfileID    int         `json:"qualityProfileId"`
	RootFolderPath      string      `json:"rootFolderPath"`
	Monitored           bool        `json:"monitored"`
	MinimumAvailability string      `json:"minimumAvailability"`
	Tags                []int       `json:"tags,omitempty"`
	AddOptions          *AddOptions `json:"addOptions,omitempty"`
}

// SearchMovieResult represents a movie search result.
type SearchMovieResult struct {
	ForeignID    string    `json:"foreignId"`
	Title        string    `json:"title"`
	Year         int       `json:"year"`
	Overview     string    `json:"overview,omitempty"`
	Studio       *Studio   `json:"studio,omitempty"`
	Runtime      int       `json:"runtime"`
	ReleaseDate  string    `json:"releaseDate,omitempty"`
	Images       []Image   `json:"images"`
	Genres       []string  `json:"genres,omitempty"`
	Ratings      *Ratings  `json:"ratings,omitempty"`
	ItemType     string    `json:"itemType"` // movie or scene
}

// Command represents a Whisparr command.
type Command struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Queued          time.Time `json:"queued"`
	Started         time.Time `json:"started,omitempty"`
	Ended           time.Time `json:"ended,omitempty"`
	Duration        string    `json:"duration,omitempty"`
	Exception       string    `json:"exception,omitempty"`
	Trigger         string    `json:"trigger"`
	StateChangeTime time.Time `json:"stateChangeTime,omitempty"`
	SendUpdatesToClient bool  `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool  `json:"updateScheduledTask"`
	LastExecutionTime   time.Time `json:"lastExecutionTime,omitempty"`
}
