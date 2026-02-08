package kitsu

// ListResponse wraps a list of resources (JSON:API format).
type ListResponse[T any] struct {
	Data     []ResourceObject[T] `json:"data"`
	Included []IncludedResource  `json:"included,omitempty"`
	Meta     Meta                `json:"meta,omitempty"`
	Links    Links               `json:"links,omitempty"`
}

// SingleResponse wraps a single resource (JSON:API format).
type SingleResponse[T any] struct {
	Data     ResourceObject[T]  `json:"data"`
	Included []IncludedResource `json:"included,omitempty"`
}

// ResourceObject is a JSON:API resource.
type ResourceObject[T any] struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    T             `json:"attributes"`
	Relationships Relationships `json:"relationships,omitempty"`
	Links         Links         `json:"links,omitempty"`
}

// IncludedResource is a generic included resource (for sideloading).
type IncludedResource struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes"`
	Links      Links          `json:"links,omitempty"`
}

// Relationships maps relationship names to relationship data.
type Relationships map[string]Relationship

// Relationship is a JSON:API relationship.
type Relationship struct {
	Data  *RelationshipData `json:"data,omitempty"`
	Links Links             `json:"links,omitempty"`
}

// RelationshipData can be a single or list of resource identifiers.
type RelationshipData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Meta contains response metadata.
type Meta struct {
	Count int `json:"count"`
}

// Links contains pagination and resource links.
type Links struct {
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
	Self  string `json:"self,omitempty"`
}

// AnimeAttributes contains Kitsu anime attributes.
type AnimeAttributes struct {
	CreatedAt           string            `json:"createdAt"`
	UpdatedAt           string            `json:"updatedAt"`
	Slug                string            `json:"slug"`
	Synopsis            string            `json:"synopsis"`
	Description         string            `json:"description"`
	CoverImageTopOffset int               `json:"coverImageTopOffset"`
	Titles              map[string]string `json:"titles"`           // en, en_jp, ja_jp, etc.
	CanonicalTitle      string            `json:"canonicalTitle"`
	AbbreviatedTitles   []string          `json:"abbreviatedTitles"`
	AverageRating       *string           `json:"averageRating"`    // "82.47" percentage string
	RatingFrequencies   map[string]string `json:"ratingFrequencies"`
	UserCount           int               `json:"userCount"`
	FavoritesCount      int               `json:"favoritesCount"`
	StartDate           *string           `json:"startDate"`        // "1998-04-03"
	EndDate             *string           `json:"endDate"`
	NextRelease         *string           `json:"nextRelease"`
	PopularityRank      *int              `json:"popularityRank"`
	RatingRank          *int              `json:"ratingRank"`
	AgeRating           *string           `json:"ageRating"`        // G, PG, R, R18
	AgeRatingGuide      *string           `json:"ageRatingGuide"`   // "17+ (violence & profanity)"
	Subtype             string            `json:"subtype"`          // TV, movie, OVA, ONA, special, music
	Status              string            `json:"status"`           // current, finished, tba, unreleased, upcoming
	TBA                 *string           `json:"tba"`
	PosterImage         *ImageSet         `json:"posterImage"`
	CoverImage          *ImageSet         `json:"coverImage"`
	EpisodeCount        *int              `json:"episodeCount"`
	EpisodeLength       *int              `json:"episodeLength"`    // Length in minutes
	TotalLength         *int              `json:"totalLength"`      // Total runtime in minutes
	YoutubeVideoID      *string           `json:"youtubeVideoId"`
	ShowType            *string           `json:"showType"`         // TV, movie, etc.
	NSFW                bool              `json:"nsfw"`
}

// EpisodeAttributes contains Kitsu episode attributes.
type EpisodeAttributes struct {
	CreatedAt      string            `json:"createdAt"`
	UpdatedAt      string            `json:"updatedAt"`
	Titles         map[string]string `json:"titles"`
	CanonicalTitle string            `json:"canonicalTitle"`
	SeasonNumber   *int              `json:"seasonNumber"`
	Number         *int              `json:"number"`
	RelativeNumber *int              `json:"relativeNumber"`
	Synopsis       string            `json:"synopsis"`
	Airdate        *string           `json:"airdate"`  // "1998-04-01"
	Length         *int              `json:"length"`    // Minutes
	Thumbnail      *ImageSet         `json:"thumbnail"`
}

// CharacterAttributes contains Kitsu character attributes.
type CharacterAttributes struct {
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	MalID       *int      `json:"malId"`
	Description *string   `json:"description"`
	Image       *ImageSet `json:"image"`
}

// CategoryAttributes contains Kitsu category attributes.
type CategoryAttributes struct {
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TotalMediaCount int       `json:"totalMediaCount"`
	Slug            string    `json:"slug"`
	NSFW            bool      `json:"nsfw"`
	ChildCount      int       `json:"childCount"`
	Image           *ImageSet `json:"image"`
}

// MappingAttributes contains external ID mapping.
type MappingAttributes struct {
	ExternalSite string `json:"externalSite"` // myanimelist/anime, thetvdb/series, etc.
	ExternalID   string `json:"externalId"`
}

// CastingAttributes contains casting information.
type CastingAttributes struct {
	Role       string  `json:"role"`        // e.g., "producer"
	VoiceActor bool    `json:"voiceActor"`
	Featured   bool    `json:"featured"`
	Language   *string `json:"language"`    // e.g., "Japanese", "English"
}

// ImageSet contains Kitsu image URLs at various sizes.
type ImageSet struct {
	Tiny     *string       `json:"tiny"`
	Small    *string       `json:"small"`
	Medium   *string       `json:"medium"`
	Large    *string       `json:"large"`
	Original *string       `json:"original"`
	Meta     *ImageSetMeta `json:"meta,omitempty"`
}

// ImageSetMeta contains image dimensions.
type ImageSetMeta struct {
	Dimensions map[string]ImageDimension `json:"dimensions,omitempty"`
}

// ImageDimension contains width and height.
type ImageDimension struct {
	Width  *int `json:"width"`
	Height *int `json:"height"`
}
