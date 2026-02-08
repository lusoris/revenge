package tvmaze

import "time"

// ShowSearchResult is a single result from /search/shows.
type ShowSearchResult struct {
	Score float64 `json:"score"`
	Show  Show    `json:"show"`
}

// Show is the TVmaze show object.
type Show struct {
	ID           int        `json:"id"`
	URL          string     `json:"url"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`    // Scripted, Reality, etc.
	Language     string     `json:"language"`
	Genres       []string   `json:"genres"`
	Status       string     `json:"status"` // Running, Ended, etc.
	Runtime      *int       `json:"runtime"`
	AverageRuntime *int     `json:"averageRuntime"`
	Premiered    *string    `json:"premiered"` // "2008-01-20"
	Ended        *string    `json:"ended"`
	OfficialSite *string    `json:"officialSite"`
	Schedule     Schedule   `json:"schedule"`
	Rating       Rating     `json:"rating"`
	Weight       int        `json:"weight"`
	Network      *Network   `json:"network"`
	WebChannel   *Network   `json:"webChannel"`
	DVRCountry   *Country   `json:"dvr_country"`
	Externals    Externals  `json:"externals"`
	Image        *ImageSet  `json:"image"`
	Summary      *string    `json:"summary"` // HTML-formatted
	Updated      int64      `json:"updated"`
}

// Schedule is the airing schedule.
type Schedule struct {
	Time string   `json:"time"`
	Days []string `json:"days"`
}

// Rating contains the TVmaze rating.
type Rating struct {
	Average *float64 `json:"average"`
}

// Network is a TV network or streaming platform.
type Network struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Country      *Country `json:"country"`
	OfficialSite *string  `json:"officialSite"`
}

// Country represents a country.
type Country struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Timezone string `json:"timezone"`
}

// Externals contains external IDs.
type Externals struct {
	TVRage *int    `json:"tvrage"`
	TVDb   *int    `json:"thetvdb"`
	IMDb   *string `json:"imdb"`
}

// ImageSet contains medium and original image URLs.
type ImageSet struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

// Season is a TVmaze season object.
type Season struct {
	ID            int       `json:"id"`
	URL           string    `json:"url"`
	Number        int       `json:"number"`
	Name          string    `json:"name"`
	EpisodeOrder  *int      `json:"episodeOrder"`
	PremiereDate  *string   `json:"premiereDate"`
	EndDate       *string   `json:"endDate"`
	Network       *Network  `json:"network"`
	WebChannel    *Network  `json:"webChannel"`
	Image         *ImageSet `json:"image"`
	Summary       *string   `json:"summary"`
}

// Episode is a TVmaze episode object.
type Episode struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Season   int       `json:"season"`
	Number   *int      `json:"number"` // nil for specials
	Type     string    `json:"type"`   // "regular", "significant_special"
	Airdate  string    `json:"airdate"`
	Airtime  string    `json:"airtime"`
	Airstamp *time.Time `json:"airstamp"`
	Runtime  *int      `json:"runtime"`
	Rating   Rating    `json:"rating"`
	Image    *ImageSet `json:"image"`
	Summary  *string   `json:"summary"`
}

// CastMember is a TVmaze cast entry.
type CastMember struct {
	Person    Person    `json:"person"`
	Character Character `json:"character"`
	Self      bool      `json:"self"`
	Voice     bool      `json:"voice"`
}

// Person is a TVmaze person object.
type Person struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Country  *Country  `json:"country"`
	Birthday *string   `json:"birthday"`
	Deathday *string   `json:"deathday"`
	Gender   *string   `json:"gender"`
	Image    *ImageSet `json:"image"`
	Updated  int64     `json:"updated"`
}

// Character is a TVmaze character.
type Character struct {
	ID    int       `json:"id"`
	URL   string    `json:"url"`
	Name  string    `json:"name"`
	Image *ImageSet `json:"image"`
}

// CrewMember is a TVmaze crew entry.
type CrewMember struct {
	Type   string `json:"type"` // "Creator", "Executive Producer", etc.
	Person Person `json:"person"`
}

// ShowImage is a TVmaze image entry from /shows/:id/images.
type ShowImage struct {
	ID          int    `json:"id"`
	Type        string `json:"type"` // "poster", "background", "banner", "typography"
	Main        bool   `json:"main"`
	Resolutions struct {
		Original *ImageResolution `json:"original"`
		Medium   *ImageResolution `json:"medium"`
	} `json:"resolutions"`
}

// ImageResolution contains URL and dimensions for an image.
type ImageResolution struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
