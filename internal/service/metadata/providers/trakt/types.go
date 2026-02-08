package trakt

import "time"

// SearchResult is a single result from the Trakt search API.
type SearchResult struct {
	Type  string  `json:"type"` // movie, show, episode, person, list
	Score float64 `json:"score"`
	Movie *Movie  `json:"movie,omitempty"`
	Show  *Show   `json:"show,omitempty"`
}

// IDs contains cross-referenced IDs for a media item.
type IDs struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	IMDb  string `json:"imdb,omitempty"`
	TMDb  int    `json:"tmdb,omitempty"`
	TVDb  int    `json:"tvdb,omitempty"`
}

// Movie represents a Trakt movie.
type Movie struct {
	Title         string     `json:"title"`
	Year          int        `json:"year,omitempty"`
	IDs           IDs        `json:"ids"`
	Tagline       string     `json:"tagline,omitempty"`
	Overview      string     `json:"overview,omitempty"`
	Released      string     `json:"released,omitempty"`
	Runtime       int        `json:"runtime,omitempty"` // minutes
	Country       string     `json:"country,omitempty"`
	Trailer       string     `json:"trailer,omitempty"`
	Homepage      string     `json:"homepage,omitempty"`
	Status        string     `json:"status,omitempty"`
	Rating        float64    `json:"rating"`
	Votes         int        `json:"votes"`
	CommentCount  int        `json:"comment_count"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	Language      string     `json:"language,omitempty"`
	Languages     []string   `json:"available_translations,omitempty"`
	Genres        []string   `json:"genres,omitempty"`
	Certification string     `json:"certification,omitempty"`
}

// Show represents a Trakt TV show.
type Show struct {
	Title         string     `json:"title"`
	Year          int        `json:"year,omitempty"`
	IDs           IDs        `json:"ids"`
	Overview      string     `json:"overview,omitempty"`
	FirstAired    *time.Time `json:"first_aired,omitempty"`
	Runtime       int        `json:"runtime,omitempty"` // minutes
	Certification string     `json:"certification,omitempty"`
	Network       string     `json:"network,omitempty"`
	Country       string     `json:"country,omitempty"`
	Trailer       string     `json:"trailer,omitempty"`
	Homepage      string     `json:"homepage,omitempty"`
	Status        string     `json:"status,omitempty"`
	Rating        float64    `json:"rating"`
	Votes         int        `json:"votes"`
	CommentCount  int        `json:"comment_count"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	Language      string     `json:"language,omitempty"`
	Languages     []string   `json:"available_translations,omitempty"`
	Genres        []string   `json:"genres,omitempty"`
	AiredEpisodes int        `json:"aired_episodes,omitempty"`
}

// Season represents a Trakt season.
type Season struct {
	Number       int        `json:"number"`
	IDs          IDs        `json:"ids"`
	Rating       float64    `json:"rating"`
	Votes        int        `json:"votes"`
	EpisodeCount int        `json:"episode_count"`
	AiredEpisodes int       `json:"aired_episodes"`
	Title        string     `json:"title,omitempty"`
	Overview     string     `json:"overview,omitempty"`
	FirstAired   *time.Time `json:"first_aired,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	Network      string     `json:"network,omitempty"`
	Episodes     []Episode  `json:"episodes,omitempty"` // populated with ?extended=episodes
}

// Episode represents a Trakt episode.
type Episode struct {
	Season        int        `json:"season"`
	Number        int        `json:"number"`
	Title         string     `json:"title,omitempty"`
	IDs           IDs        `json:"ids"`
	Overview      string     `json:"overview,omitempty"`
	Rating        float64    `json:"rating"`
	Votes         int        `json:"votes"`
	CommentCount  int        `json:"comment_count"`
	FirstAired    *time.Time `json:"first_aired,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	Runtime       int        `json:"runtime,omitempty"`
}

// CastMember represents a person in a cast/crew listing.
type CastMember struct {
	Characters []string `json:"characters,omitempty"`
	Person     Person   `json:"person"`
}

// CrewMember represents a crew member.
type CrewMember struct {
	Jobs   []string `json:"jobs,omitempty"`
	Person Person   `json:"person"`
}

// Person represents a person (actor, director, etc.).
type Person struct {
	Name       string     `json:"name"`
	IDs        IDs        `json:"ids"`
	Biography  string     `json:"biography,omitempty"`
	Birthday   *time.Time `json:"birthday,omitempty"`
	Death      *time.Time `json:"death,omitempty"`
	Birthplace string     `json:"birthplace,omitempty"`
	Homepage   string     `json:"homepage,omitempty"`
	Gender     string     `json:"gender,omitempty"`
}

// Credits contains cast and crew for a movie or show.
type Credits struct {
	Cast []CastMember            `json:"cast,omitempty"`
	Crew map[string][]CrewMember `json:"crew,omitempty"` // keyed by department
}

// Translation is a localized version of a media item.
type Translation struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Language string `json:"language"`
	Tagline  string `json:"tagline,omitempty"`
	Country  string `json:"country,omitempty"`
}
