package movie

import (
	"time"
)

type TMDbMovie struct {
	ID                  int                 `json:"id"`
	IMDbID              *string             `json:"imdb_id"`
	Title               string              `json:"title"`
	OriginalTitle       string              `json:"original_title"`
	OriginalLanguage    string              `json:"original_language"`
	Overview            *string             `json:"overview"`
	Tagline             *string             `json:"tagline"`
	ReleaseDate         string              `json:"release_date"`
	Runtime             *int                `json:"runtime"`
	Budget              *int64              `json:"budget"`
	Revenue             *int64              `json:"revenue"`
	Status              string              `json:"status"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
	Popularity          float64             `json:"popularity"`
	Adult               bool                `json:"adult"`
	Video               bool                `json:"video"`
	PosterPath          *string             `json:"poster_path"`
	BackdropPath        *string             `json:"backdrop_path"`
	Homepage            *string             `json:"homepage"`
	Genres              []Genre             `json:"genres"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	BelongsToCollection *Collection         `json:"belongs_to_collection"`
}

type TMDbSearchResponse struct {
	Page         int                `json:"page"`
	Results      []TMDbSearchResult `json:"results"`
	TotalPages   int                `json:"total_pages"`
	TotalResults int                `json:"total_results"`
}

type TMDbSearchResult struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       *string `json:"poster_path"`
	BackdropPath     *string `json:"backdrop_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
	Adult            bool    `json:"adult"`
	Video            bool    `json:"video"`
	GenreIDs         []int   `json:"genre_ids"`
}

type TMDbCredits struct {
	ID   int          `json:"id"`
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

type CastMember struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Character   string  `json:"character"`
	Order       int     `json:"order"`
	CreditID    string  `json:"credit_id"`
	Gender      *int    `json:"gender"`
	ProfilePath *string `json:"profile_path"`
}

type CrewMember struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Job         string  `json:"job"`
	Department  string  `json:"department"`
	CreditID    string  `json:"credit_id"`
	Gender      *int    `json:"gender"`
	ProfilePath *string `json:"profile_path"`
}

type TMDbImages struct {
	ID        int     `json:"id"`
	Backdrops []Image `json:"backdrops"`
	Posters   []Image `json:"posters"`
	Logos     []Image `json:"logos"`
}

type Image struct {
	AspectRatio float64 `json:"aspect_ratio"`
	FilePath    string  `json:"file_path"`
	Height      int     `json:"height"`
	Width       int     `json:"width"`
	ISO639_1    *string `json:"iso_639_1"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

type TMDbCollectionDetails struct {
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	Overview     string             `json:"overview"`
	PosterPath   *string            `json:"poster_path"`
	BackdropPath *string            `json:"backdrop_path"`
	Parts        []TMDbSearchResult `json:"parts"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	LogoPath      *string `json:"logo_path"`
	OriginCountry string  `json:"origin_country"`
}

type ProductionCountry struct {
	ISO3166_1 string `json:"iso_3166_1"`
	Name      string `json:"name"`
}

type SpokenLanguage struct {
	ISO639_1    string `json:"iso_639_1"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

type Collection struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	PosterPath   *string `json:"poster_path"`
	BackdropPath *string `json:"backdrop_path"`
}

type TMDbError struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}

type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

func (c *CacheEntry) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
