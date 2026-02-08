package anidb

import "encoding/xml"

// AnimeResponse is the top-level XML response for an anime lookup.
type AnimeResponse struct {
	XMLName    xml.Name       `xml:"anime"`
	ID         int            `xml:"id,attr"`
	Restricted bool           `xml:"restricted,attr"`
	Type       string         `xml:"type"`
	EpCount    int            `xml:"episodecount"`
	StartDate  string         `xml:"startdate"`
	EndDate    string         `xml:"enddate"`
	Titles     Titles         `xml:"titles"`
	Related    RelatedAnime   `xml:"relatedanime"`
	Similar    SimilarAnime   `xml:"similaranime"`
	URL        string         `xml:"url"`
	Creators   Creators       `xml:"creators"`
	Description string        `xml:"description"`
	Ratings    Ratings        `xml:"ratings"`
	Picture    string         `xml:"picture"`
	Resources  Resources      `xml:"resources"`
	Tags       Tags           `xml:"tags"`
	Characters Characters     `xml:"characters"`
	Episodes   Episodes       `xml:"episodes"`
}

// Titles contains all titles for an anime.
type Titles struct {
	Title []Title `xml:"title"`
}

// Title is a single anime title with language and type.
type Title struct {
	Lang string `xml:"lang,attr"`
	Type string `xml:"type,attr"` // main, official, synonym, short
	Text string `xml:",chardata"`
}

// RelatedAnime holds related anime entries.
type RelatedAnime struct {
	Anime []RelatedEntry `xml:"anime"`
}

// RelatedEntry is a reference to a related anime.
type RelatedEntry struct {
	ID   int    `xml:"id,attr"`
	Type string `xml:"type,attr"` // Sequel, Prequel, etc.
	Text string `xml:",chardata"`
}

// SimilarAnime holds similar anime entries.
type SimilarAnime struct {
	Anime []SimilarEntry `xml:"anime"`
}

// SimilarEntry is a reference to a similar anime with approval ratings.
type SimilarEntry struct {
	ID       int    `xml:"id,attr"`
	Approval int    `xml:"approval,attr"`
	Total    int    `xml:"total,attr"`
	Text     string `xml:",chardata"`
}

// Creators contains staff/creator information.
type Creators struct {
	Name []Creator `xml:"name"`
}

// Creator is a staff member (director, music, etc.).
type Creator struct {
	ID   int    `xml:"id,attr"`
	Type string `xml:"type,attr"` // Direction, Music, Animation Work, etc.
	Text string `xml:",chardata"`
}

// Ratings contains anime rating values.
type Ratings struct {
	Permanent RatingValue `xml:"permanent"`
	Temporary RatingValue `xml:"temporary"`
	Review    RatingValue `xml:"review"`
}

// RatingValue holds a rating with vote count.
type RatingValue struct {
	Count int     `xml:"count,attr"`
	Value float64 `xml:",chardata"`
}

// Resources contains external ID references.
type Resources struct {
	Resource []Resource `xml:"resource"`
}

// Resource is an external ID reference.
type Resource struct {
	Type       int                `xml:"type,attr"`
	ExternalID []ExternalEntity   `xml:"externalentity"`
}

// ExternalEntity contains identifiers for external resources.
type ExternalEntity struct {
	Identifier []string `xml:"identifier"`
}

// Tags contains anime tags.
type Tags struct {
	Tag []Tag `xml:"tag"`
}

// Tag is a content tag with weight.
type Tag struct {
	ID          int    `xml:"id,attr"`
	Weight      int    `xml:"weight,attr"`
	LocalSpoiler bool  `xml:"localspoiler,attr"`
	GlobalSpoiler bool `xml:"globalspoiler,attr"`
	Verified    bool   `xml:"verified,attr"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
}

// Characters contains character entries.
type Characters struct {
	Character []Character `xml:"character"`
}

// Character represents an anime character.
type Character struct {
	ID      int            `xml:"id,attr"`
	Type    string         `xml:"type,attr"` // main character in, secondary cast in, etc.
	Rating  *CharRating    `xml:"rating"`
	Name    string         `xml:"name"`
	Gender  string         `xml:"gender"`
	Picture string         `xml:"picture"`
	Seiyuu  []Seiyuu       `xml:"seiyuu"`
}

// CharRating is a character's rating.
type CharRating struct {
	Votes int     `xml:"votes,attr"`
	Value float64 `xml:",chardata"`
}

// Seiyuu is a voice actor for a character.
type Seiyuu struct {
	ID      int    `xml:"id,attr"`
	Picture string `xml:"picture,attr"`
	Text    string `xml:",chardata"`
}

// Episodes contains episode entries.
type Episodes struct {
	Episode []Episode `xml:"episode"`
}

// Episode represents a single episode.
type Episode struct {
	ID      int           `xml:"id,attr"`
	Update  string        `xml:"update,attr"`
	EpNo    EpNo          `xml:"epno"`
	Length  int           `xml:"length"`  // minutes
	Airdate string        `xml:"airdate"`
	Rating  *EpisodeRating `xml:"rating"`
	Title   []EpTitle     `xml:"title"`
}

// EpNo is the episode number with type.
type EpNo struct {
	Type int    `xml:"type,attr"` // 1=regular, 2=special, 3=credits, 4=trailer, 5=parody, 6=other
	Text string `xml:",chardata"`
}

// EpTitle is an episode title with language.
type EpTitle struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",chardata"`
}

// EpisodeRating is an episode's rating.
type EpisodeRating struct {
	Votes int     `xml:"votes,attr"`
	Value float64 `xml:",chardata"`
}

// ErrorResponse is the AniDB error format.
type ErrorResponse struct {
	XMLName xml.Name `xml:"error"`
	Text    string   `xml:",chardata"`
}

// TitleDumpEntry represents a parsed line from the title dump.
type TitleDumpEntry struct {
	AID    int
	Type   string // main, official, synonym, short
	Lang   string
	Title  string
}
