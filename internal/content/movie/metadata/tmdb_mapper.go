package metadata

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/shopspring/decimal"
)

type TMDbMapper struct {
	client *TMDbClient
}

func NewTMDbMapper(client *TMDbClient) *TMDbMapper {
	return &TMDbMapper{
		client: client,
	}
}

func (m *TMDbMapper) MapMovie(tmdbMovie *TMDbMovie) *movie.Movie {
	mov := &movie.Movie{
		ID:               uuid.New(),
		TMDbID:           parseOptionalInt32(tmdbMovie.ID),
		IMDbID:           parseOptionalString(tmdbMovie.IMDbID),
		Title:            tmdbMovie.Title,
		OriginalTitle:    parseOptionalString(&tmdbMovie.OriginalTitle),
		OriginalLanguage: parseOptionalString(&tmdbMovie.OriginalLanguage),
		Overview:         parseOptionalString(tmdbMovie.Overview),
		Tagline:          parseOptionalString(tmdbMovie.Tagline),
		ReleaseDate:      parseReleaseDate(tmdbMovie.ReleaseDate),
		Runtime:          parseOptionalInt32Ptr(tmdbMovie.Runtime),
		Budget:           parseOptionalInt64Ptr(tmdbMovie.Budget),
		Revenue:          parseOptionalInt64Ptr(tmdbMovie.Revenue),
		Status:           parseOptionalString(&tmdbMovie.Status),
		VoteAverage:      parseDecimal(tmdbMovie.VoteAverage),
		VoteCount:        parseOptionalInt32(tmdbMovie.VoteCount),
		Popularity:       parseDecimal(tmdbMovie.Popularity),
		PosterPath:       parseOptionalString(tmdbMovie.PosterPath),
		BackdropPath:     parseOptionalString(tmdbMovie.BackdropPath),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if tmdbMovie.ReleaseDate != "" {
		mov.Year = parseOptionalInt32Ptr(ExtractYear(tmdbMovie.ReleaseDate))
	}

	return mov
}

func (m *TMDbMapper) MapSearchResult(result *TMDbSearchResult) *movie.Movie {
	mov := &movie.Movie{
		ID:               uuid.New(),
		TMDbID:           parseOptionalInt32(result.ID),
		Title:            result.Title,
		OriginalTitle:    parseOptionalString(&result.OriginalTitle),
		OriginalLanguage: parseOptionalString(&result.OriginalLanguage),
		Overview:         parseOptionalString(&result.Overview),
		ReleaseDate:      parseReleaseDate(result.ReleaseDate),
		VoteAverage:      parseDecimal(result.VoteAverage),
		VoteCount:        parseOptionalInt32(result.VoteCount),
		Popularity:       parseDecimal(result.Popularity),
		PosterPath:       parseOptionalString(result.PosterPath),
		BackdropPath:     parseOptionalString(result.BackdropPath),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if result.ReleaseDate != "" {
		mov.Year = parseOptionalInt32Ptr(ExtractYear(result.ReleaseDate))
	}

	return mov
}

func (m *TMDbMapper) MapCredits(movieID uuid.UUID, credits *TMDbCredits) []movie.MovieCredit {
	var result []movie.MovieCredit

	for _, cast := range credits.Cast {
		result = append(result, movie.MovieCredit{
			ID:           uuid.New(),
			MovieID:      movieID,
			TMDbPersonID: int32(cast.ID),
			Name:         cast.Name,
			CreditType:   "cast",
			Character:    parseOptionalString(&cast.Character),
			Department:   nil,
			Job:          nil,
			CastOrder:    parseOptionalInt32Ptr(&cast.Order),
			ProfilePath:  parseOptionalString(cast.ProfilePath),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	for _, crew := range credits.Crew {
		result = append(result, movie.MovieCredit{
			ID:           uuid.New(),
			MovieID:      movieID,
			TMDbPersonID: int32(crew.ID),
			Name:         crew.Name,
			CreditType:   "crew",
			Character:    nil,
			Department:   parseOptionalString(&crew.Department),
			Job:          parseOptionalString(&crew.Job),
			CastOrder:    nil,
			ProfilePath:  parseOptionalString(crew.ProfilePath),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	return result
}

func (m *TMDbMapper) MapGenres(movieID uuid.UUID, genres []Genre) []movie.MovieGenre {
	var result []movie.MovieGenre

	for _, genre := range genres {
		result = append(result, movie.MovieGenre{
			ID:          uuid.New(),
			MovieID:     movieID,
			TMDbGenreID: int32(genre.ID),
			Name:        genre.Name,
			CreatedAt:   time.Now(),
		})
	}

	return result
}

func (m *TMDbMapper) MapCollection(collection *TMDbCollectionDetails) *movie.MovieCollection {
	return &movie.MovieCollection{
		ID:               uuid.New(),
		TMDbCollectionID: parseOptionalInt32(collection.ID),
		Name:             collection.Name,
		Overview:         parseOptionalString(&collection.Overview),
		PosterPath:       parseOptionalString(collection.PosterPath),
		BackdropPath:     parseOptionalString(collection.BackdropPath),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (m *TMDbMapper) GetPosterURL(posterPath *string, size string) *string {
	if posterPath == nil || *posterPath == "" {
		return nil
	}

	if size == "" {
		size = "w500"
	}

	url := m.client.GetImageURL(*posterPath, size)
	return &url
}

func (m *TMDbMapper) GetBackdropURL(backdropPath *string, size string) *string {
	if backdropPath == nil || *backdropPath == "" {
		return nil
	}

	if size == "" {
		size = "w1280"
	}

	url := m.client.GetImageURL(*backdropPath, size)
	return &url
}

func parseOptionalString(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func parseOptionalInt32(i int) *int32 {
	if i == 0 {
		return nil
	}
	val := int32(i)
	return &val
}

func parseOptionalInt32Ptr(i *int) *int32 {
	if i == nil || *i == 0 {
		return nil
	}
	val := int32(*i)
	return &val
}

func parseOptionalInt64Ptr(i *int64) *int64 {
	if i == nil || *i == 0 {
		return nil
	}
	return i
}

func parseDecimal(f float64) *decimal.Decimal {
	if f == 0 {
		return nil
	}
	d := decimal.NewFromFloat(f)
	return &d
}

func parseReleaseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}

	return &t
}

func ExtractYear(releaseDate string) *int {
	if releaseDate == "" {
		return nil
	}

	parts := strings.Split(releaseDate, "-")
	if len(parts) == 0 {
		return nil
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}

	return &year
}
