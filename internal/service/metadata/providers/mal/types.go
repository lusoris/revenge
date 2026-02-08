package mal

// ListResponse wraps paginated search results.
type ListResponse struct {
	Data   []ListNode `json:"data"`
	Paging Paging     `json:"paging"`
}

// ListNode wraps a single item in a list response.
type ListNode struct {
	Node Anime `json:"node"`
}

// Paging contains cursor-based pagination links.
type Paging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// Anime represents a MAL anime entry.
type Anime struct {
	ID                     int               `json:"id"`
	Title                  string            `json:"title"`
	MainPicture            *Picture          `json:"main_picture"`
	AlternativeTitles      AlternativeTitles `json:"alternative_titles"`
	StartDate              string            `json:"start_date"`
	EndDate                string            `json:"end_date"`
	Synopsis               string            `json:"synopsis"`
	Mean                   *float64          `json:"mean"`
	Rank                   *int              `json:"rank"`
	Popularity             *int              `json:"popularity"`
	NumListUsers           int               `json:"num_list_users"`
	NumScoringUsers        int               `json:"num_scoring_users"`
	NSFW                   string            `json:"nsfw"`
	MediaType              string            `json:"media_type"`
	Status                 string            `json:"status"`
	Genres                 []Genre           `json:"genres"`
	NumEpisodes            int               `json:"num_episodes"`
	StartSeason            *Season           `json:"start_season"`
	Broadcast              *Broadcast        `json:"broadcast"`
	Source                 string            `json:"source"`
	AverageEpisodeDuration int               `json:"average_episode_duration"` // seconds
	Rating                 string            `json:"rating"`
	Pictures               []Picture         `json:"pictures"`
	Background             string            `json:"background"`
	RelatedAnime           []RelatedAnime    `json:"related_anime"`
	Recommendations        []Recommendation  `json:"recommendations"`
	Studios                []Studio          `json:"studios"`
	Statistics             *Statistics       `json:"statistics"`
	CreatedAt              string            `json:"created_at"`
	UpdatedAt              string            `json:"updated_at"`
}

// Picture contains image URLs.
type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// AlternativeTitles contains alternative names.
type AlternativeTitles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

// Genre is a genre tag.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Season represents an anime season.
type Season struct {
	Year   int    `json:"year"`
	Season string `json:"season"` // winter, spring, summer, fall
}

// Broadcast contains airing schedule info.
type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week"`
	StartTime    string `json:"start_time"`
}

// RelatedAnime contains a reference to a related anime.
type RelatedAnime struct {
	Node                  Anime  `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// Recommendation is a user recommendation.
type Recommendation struct {
	Node               Anime `json:"node"`
	NumRecommendations int   `json:"num_recommendations"`
}

// Studio represents an animation studio.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Statistics contains user statistics.
type Statistics struct {
	Status StatusStats `json:"status"`
}

// StatusStats contains per-status user counts.
type StatusStats struct {
	Watching    string `json:"watching"`
	Completed   string `json:"completed"`
	OnHold      string `json:"on_hold"`
	Dropped     string `json:"dropped"`
	PlanToWatch string `json:"plan_to_watch"`
}

// ErrorResponse is the MAL error format.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
