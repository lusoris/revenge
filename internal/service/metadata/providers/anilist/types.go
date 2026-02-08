package anilist

// GraphQLRequest is the request body for AniList GraphQL API.
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// GraphQLResponse wraps the top-level AniList response.
type GraphQLResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError is an AniList API error.
type GraphQLError struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Locations []struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"locations,omitempty"`
}

// PageData wraps paginated results.
type PageData struct {
	Page Page `json:"Page"`
}

// Page contains pagination info and results.
type Page struct {
	PageInfo PageInfo `json:"pageInfo"`
	Media    []Media  `json:"media,omitempty"`
	Staff    []Staff  `json:"staff,omitempty"`
}

// PageInfo contains pagination details.
type PageInfo struct {
	Total       int  `json:"total"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
	PerPage     int  `json:"perPage"`
}

// MediaData wraps a single media lookup.
type MediaData struct {
	Media *Media `json:"Media"`
}

// Media represents an AniList anime/manga entry.
type Media struct {
	ID                int                 `json:"id"`
	IDMal             *int                `json:"idMal"`
	Title             MediaTitle          `json:"title"`
	Type              string              `json:"type"`
	Format            string              `json:"format"`
	Status            string              `json:"status"`
	Description       *string             `json:"description"`
	StartDate         FuzzyDate           `json:"startDate"`
	EndDate           FuzzyDate           `json:"endDate"`
	Season            *string             `json:"season"`
	SeasonYear        *int                `json:"seasonYear"`
	Episodes          *int                `json:"episodes"`
	Duration          *int                `json:"duration"`
	CountryOfOrigin   *string             `json:"countryOfOrigin"`
	Source            *string             `json:"source"`
	Genres            []string            `json:"genres"`
	Tags              []MediaTag          `json:"tags"`
	AverageScore      *int                `json:"averageScore"`
	MeanScore         *int                `json:"meanScore"`
	Popularity        int                 `json:"popularity"`
	Favourites        int                 `json:"favourites"`
	IsAdult           bool                `json:"isAdult"`
	SiteURL           string              `json:"siteUrl"`
	CoverImage        CoverImage          `json:"coverImage"`
	BannerImage       *string             `json:"bannerImage"`
	Studios           StudioConnection    `json:"studios"`
	ExternalLinks     []ExternalLink      `json:"externalLinks"`
	StreamingEpisodes []StreamingEpisode  `json:"streamingEpisodes"`
	Trailer           *Trailer            `json:"trailer"`
	Characters        CharacterConnection `json:"characters"`
	Staff             StaffConnection     `json:"staff"`
	Relations         MediaConnection     `json:"relations"`
	Synonyms          []string            `json:"synonyms"`
}

// MediaTitle contains multi-language titles.
type MediaTitle struct {
	Romaji        *string `json:"romaji"`
	English       *string `json:"english"`
	Native        *string `json:"native"`
	UserPreferred *string `json:"userPreferred"`
}

// FuzzyDate represents an AniList partial date.
type FuzzyDate struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

// MediaTag represents a content tag.
type MediaTag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Rank     int    `json:"rank"`
	IsAdult  bool   `json:"isAdult"`
}

// CoverImage contains image URLs at different sizes.
type CoverImage struct {
	ExtraLarge *string `json:"extraLarge"`
	Large      *string `json:"large"`
	Medium     *string `json:"medium"`
	Color      *string `json:"color"`
}

// StudioConnection contains studio edges.
type StudioConnection struct {
	Edges []StudioEdge `json:"edges"`
}

// StudioEdge connects a studio to a media.
type StudioEdge struct {
	Node   Studio `json:"node"`
	IsMain bool   `json:"isMain"`
}

// Studio represents an animation studio.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ExternalLink is a link to an external site.
type ExternalLink struct {
	ID       int     `json:"id"`
	URL      *string `json:"url"`
	Site     string  `json:"site"`
	SiteID   *int    `json:"siteId"`
	Type     string  `json:"type"`
	Language *string `json:"language"`
}

// StreamingEpisode is a streaming source for an episode.
type StreamingEpisode struct {
	Title     *string `json:"title"`
	Thumbnail *string `json:"thumbnail"`
	URL       *string `json:"url"`
	Site      *string `json:"site"`
}

// Trailer is a media trailer reference.
type Trailer struct {
	ID        *string `json:"id"`
	Site      *string `json:"site"`
	Thumbnail *string `json:"thumbnail"`
}

// CharacterConnection contains character edges.
type CharacterConnection struct {
	Edges []CharacterEdge `json:"edges"`
}

// CharacterEdge connects a character to a media.
type CharacterEdge struct {
	Node        Character `json:"node"`
	Role        string    `json:"role"`
	VoiceActors []Staff   `json:"voiceActors"`
}

// Character represents an anime/manga character.
type Character struct {
	ID     int            `json:"id"`
	Name   CharacterName  `json:"name"`
	Image  CharacterImage `json:"image"`
	Gender *string        `json:"gender"`
}

// CharacterName contains character name variants.
type CharacterName struct {
	Full   *string `json:"full"`
	Native *string `json:"native"`
}

// CharacterImage contains character images.
type CharacterImage struct {
	Large  *string `json:"large"`
	Medium *string `json:"medium"`
}

// StaffConnection contains staff edges.
type StaffConnection struct {
	Edges []StaffEdge `json:"edges"`
}

// StaffEdge connects staff to a media.
type StaffEdge struct {
	Node Staff  `json:"node"`
	Role string `json:"role"`
}

// Staff represents a person (voice actor, director, etc.).
type Staff struct {
	ID                 int        `json:"id"`
	Name               StaffName  `json:"name"`
	LanguageV2         *string    `json:"languageV2"`
	Image              StaffImage `json:"image"`
	PrimaryOccupations []string   `json:"primaryOccupations"`
	Gender             *string    `json:"gender"`
}

// StaffName contains staff name variants.
type StaffName struct {
	Full   *string `json:"full"`
	Native *string `json:"native"`
}

// StaffImage contains staff images.
type StaffImage struct {
	Large  *string `json:"large"`
	Medium *string `json:"medium"`
}

// MediaConnection contains related media edges.
type MediaConnection struct {
	Edges []MediaEdge `json:"edges"`
}

// MediaEdge connects related media.
type MediaEdge struct {
	Node         Media  `json:"node"`
	RelationType string `json:"relationType"`
}
