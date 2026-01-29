// Package tmdb provides a TMDb (The Movie Database) API client.
package tmdb

import "time"

// Movie represents a TMDb movie response.
type Movie struct {
	ID               int     `json:"id"`
	IMDbID           string  `json:"imdb_id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Tagline          string  `json:"tagline"`
	ReleaseDate      string  `json:"release_date"`
	Runtime          int     `json:"runtime"` // minutes
	Budget           int64   `json:"budget"`
	Revenue          int64   `json:"revenue"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	Adult            bool    `json:"adult"`
	Video            bool    `json:"video"`
	Status           string  `json:"status"`
	Homepage         string  `json:"homepage"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`

	// Included via append_to_response
	Genres              []Genre      `json:"genres"`
	ProductionCompanies []Company    `json:"production_companies"`
	ProductionCountries []Country    `json:"production_countries"`
	SpokenLanguages     []Language   `json:"spoken_languages"`
	BelongsToCollection *Collection  `json:"belongs_to_collection"`
	Credits             *Credits     `json:"credits,omitempty"`
	Images              *Images      `json:"images,omitempty"`
	Videos              *Videos      `json:"videos,omitempty"`
	ExternalIDs         *ExternalIDs `json:"external_ids,omitempty"`
}

// ReleaseYear extracts year from release_date.
func (m *Movie) ReleaseYear() int {
	if len(m.ReleaseDate) < 4 {
		return 0
	}
	t, err := time.Parse("2006-01-02", m.ReleaseDate)
	if err != nil {
		return 0
	}
	return t.Year()
}

// Genre represents a TMDb genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Company represents a production company.
type Company struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

// Country represents a production country.
type Country struct {
	ISO31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

// Language represents a spoken language.
type Language struct {
	ISO6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

// Collection represents a movie collection.
type Collection struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

// Credits contains cast and crew.
type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// CastMember represents an actor.
type CastMember struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Character    string  `json:"character"`
	Order        int     `json:"order"`
	Gender       int     `json:"gender"`
	ProfilePath  string  `json:"profile_path"`
	Popularity   float64 `json:"popularity"`
	KnownFor     string  `json:"known_for_department"`
	Adult        bool    `json:"adult"`
	CreditID     string  `json:"credit_id"`
	CastID       int     `json:"cast_id"`
}

// CrewMember represents a crew member.
type CrewMember struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Department   string  `json:"department"`
	Job          string  `json:"job"`
	Gender       int     `json:"gender"`
	ProfilePath  string  `json:"profile_path"`
	Popularity   float64 `json:"popularity"`
	KnownFor     string  `json:"known_for_department"`
	Adult        bool    `json:"adult"`
	CreditID     string  `json:"credit_id"`
}

// Images contains movie images.
type Images struct {
	Backdrops []Image `json:"backdrops"`
	Posters   []Image `json:"posters"`
	Logos     []Image `json:"logos"`
}

// Image represents a single image.
type Image struct {
	FilePath    string  `json:"file_path"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	AspectRatio float64 `json:"aspect_ratio"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
	ISO6391     string  `json:"iso_639_1"`
}

// Videos contains movie videos.
type Videos struct {
	Results []Video `json:"results"`
}

// Video represents a single video.
type Video struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Site        string `json:"site"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Official    bool   `json:"official"`
	PublishedAt string `json:"published_at"`
	ISO6391     string `json:"iso_639_1"`
	ISO31661    string `json:"iso_3166_1"`
}

// ExternalIDs contains external provider IDs.
type ExternalIDs struct {
	IMDbID      string `json:"imdb_id"`
	WikidataID  string `json:"wikidata_id"`
	FacebookID  string `json:"facebook_id"`
	InstagramID string `json:"instagram_id"`
	TwitterID   string `json:"twitter_id"`
}

// SearchResult represents a movie search response.
type SearchResult struct {
	Page         int           `json:"page"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
	Results      []MovieResult `json:"results"`
}

// MovieResult represents a single search result.
type MovieResult struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	Adult            bool    `json:"adult"`
	Video            bool    `json:"video"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `json:"genre_ids"`
}

// ReleaseYear extracts year from release_date.
func (m *MovieResult) ReleaseYear() int {
	if len(m.ReleaseDate) < 4 {
		return 0
	}
	t, err := time.Parse("2006-01-02", m.ReleaseDate)
	if err != nil {
		return 0
	}
	return t.Year()
}

// FindResult represents a find by external ID response.
type FindResult struct {
	MovieResults     []MovieResult `json:"movie_results"`
	TVResults        []any         `json:"tv_results"`
	PersonResults    []any         `json:"person_results"`
	TVEpisodeResults []any         `json:"tv_episode_results"`
	TVSeasonResults  []any         `json:"tv_season_results"`
}

// Configuration represents TMDb API configuration.
type Configuration struct {
	Images     ImageConfig `json:"images"`
	ChangeKeys []string    `json:"change_keys"`
}

// ImageConfig contains image URL configuration.
type ImageConfig struct {
	BaseURL       string   `json:"base_url"`
	SecureBaseURL string   `json:"secure_base_url"`
	BackdropSizes []string `json:"backdrop_sizes"`
	LogoSizes     []string `json:"logo_sizes"`
	PosterSizes   []string `json:"poster_sizes"`
	ProfileSizes  []string `json:"profile_sizes"`
	StillSizes    []string `json:"still_sizes"`
}

// APIError represents a TMDb API error response.
type APIError struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}

// Error implements error interface.
func (e APIError) Error() string {
	return e.StatusMessage
}
