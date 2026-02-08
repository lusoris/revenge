package letterboxd

// TokenResponse is the response from the OAuth2 token endpoint.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// SearchResponse is the response from the /search endpoint.
type SearchResponse struct {
	Next      string       `json:"next,omitempty"`
	Items     []SearchItem `json:"items"`
	ItemCount int          `json:"itemCount,omitempty"`
}

// SearchItem is a single item in a search response.
type SearchItem struct {
	Type  string       `json:"type"`
	Score float64      `json:"score"`
	Film  *FilmSummary `json:"film,omitempty"`
}

// FilmSummary is a lightweight film representation returned in lists/search.
type FilmSummary struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	OriginalName     string               `json:"originalName,omitempty"`
	SortingName      string               `json:"sortingName,omitempty"`
	AlternativeNames []string             `json:"alternativeNames,omitempty"`
	ReleaseYear      int                  `json:"releaseYear,omitempty"`
	RunTime          int                  `json:"runTime,omitempty"`
	Rating           float64              `json:"rating,omitempty"`
	Adult            bool                 `json:"adult"`
	FilmCollectionID string               `json:"filmCollectionId,omitempty"`
	Directors        []ContributorSummary `json:"directors"`
	Poster           *Image               `json:"poster,omitempty"`
	Genres           []Genre              `json:"genres"`
	Links            []Link               `json:"links"`
	Top250Position   int                  `json:"top250Position,omitempty"`
}

// Film is the full film representation from GET /film/{id}.
type Film struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	OriginalName       string               `json:"originalName,omitempty"`
	SortingName        string               `json:"sortingName,omitempty"`
	AlternativeNames   []string             `json:"alternativeNames,omitempty"`
	ReleaseYear        int                  `json:"releaseYear,omitempty"`
	RunTime            int                  `json:"runTime,omitempty"`
	Rating             float64              `json:"rating,omitempty"`
	Adult              bool                 `json:"adult"`
	FilmCollectionID   string               `json:"filmCollectionId,omitempty"`
	Top250Position     int                  `json:"top250Position,omitempty"`
	Tagline            string               `json:"tagline,omitempty"`
	Description        string               `json:"description,omitempty"`
	Directors          []ContributorSummary `json:"directors,omitempty"`
	Poster             *Image               `json:"poster,omitempty"`
	Backdrop           *Image               `json:"backdrop,omitempty"`
	BackdropFocalPoint float64              `json:"backdropFocalPoint,omitempty"`
	Trailer            *FilmTrailer         `json:"trailer,omitempty"`
	Genres             []Genre              `json:"genres"`
	Links              []Link               `json:"links"`
	Countries          []Country            `json:"countries,omitempty"`
	Languages          []Language           `json:"languages,omitempty"`
	PrimaryLanguage    *Language            `json:"primaryLanguage,omitempty"`
	Contributions      []Contributions      `json:"contributions"`
	Releases           []Release            `json:"releases,omitempty"`
}

// ContributorSummary is a lightweight contributor.
type ContributorSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CharacterName string `json:"characterName,omitempty"`
	TMDbID        string `json:"tmdbid,omitempty"`
}

// Contributions groups contributors by type.
type Contributions struct {
	Type         string               `json:"type"`
	Contributors []ContributorSummary `json:"contributors"`
}

// Image represents an image with multiple sizes.
type Image struct {
	Sizes []ImageSize `json:"sizes"`
}

// ImageSize is a specific size variant of an image.
type ImageSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// Genre is a film genre.
type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Link is a related URL.
type Link struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	URL  string `json:"url"`
}

// Country is a production country.
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Language represents a language.
type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// FilmTrailer is a film trailer (YouTube).
type FilmTrailer struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	URL  string `json:"url"`
}

// Release is film release information.
type Release struct {
	Type          string  `json:"type"`
	Country       Country `json:"country"`
	Certification string  `json:"certification,omitempty"`
	ReleaseDate   string  `json:"releaseDate"`
}

// FilmStatistics is statistical data about a film.
type FilmStatistics struct {
	Film   FilmIdentifier       `json:"film"`
	Counts FilmStatisticsCounts `json:"counts"`
	Rating float64              `json:"rating,omitempty"`
}

// FilmIdentifier contains just the LID of a film.
type FilmIdentifier struct {
	ID string `json:"id"`
}

// FilmStatisticsCounts contains counts for a film.
type FilmStatisticsCounts struct {
	Watches int `json:"watches"`
	Likes   int `json:"likes"`
	Ratings int `json:"ratings"`
	Fans    int `json:"fans"`
	Lists   int `json:"lists"`
	Reviews int `json:"reviews"`
}

// Contributor is the full contributor representation.
type Contributor struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	TMDbID     string           `json:"tmdbid,omitempty"`
	Bio        string           `json:"bio,omitempty"`
	Statistics ContributorStats `json:"statistics"`
	Links      []Link           `json:"links"`
}

// ContributorStats contains statistics for a contributor.
type ContributorStats struct {
	Contributions []ContributionStatistics `json:"contributions"`
}

// ContributionStatistics counts films per contribution type.
type ContributionStatistics struct {
	Type      string `json:"type"`
	FilmCount int    `json:"filmCount"`
}

// FilmCollection is a collection of related films.
type FilmCollection struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	FilmCount int           `json:"filmCount"`
	Films     []FilmSummary `json:"films"`
	Links     []Link        `json:"links"`
}
