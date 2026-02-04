package movie

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/util"
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

func (m *TMDbMapper) MapMovie(tmdbMovie *TMDbMovie) *Movie {
	mov := &Movie{
		ID:               uuid.New(),
		TMDbID:           tmdbParseOptionalInt32(tmdbMovie.ID),
		IMDbID:           tmdbParseOptionalString(tmdbMovie.IMDbID),
		Title:            tmdbMovie.Title,
		OriginalTitle:    tmdbParseOptionalString(&tmdbMovie.OriginalTitle),
		OriginalLanguage: tmdbParseOptionalString(&tmdbMovie.OriginalLanguage),
		Overview:         tmdbParseOptionalString(tmdbMovie.Overview),
		Tagline:          tmdbParseOptionalString(tmdbMovie.Tagline),
		ReleaseDate:      parseReleaseDate(tmdbMovie.ReleaseDate),
		Runtime:          tmdbParseOptionalInt32Ptr(tmdbMovie.Runtime),
		Budget:           parseOptionalInt64Ptr(tmdbMovie.Budget),
		Revenue:          parseOptionalInt64Ptr(tmdbMovie.Revenue),
		Status:           tmdbParseOptionalString(&tmdbMovie.Status),
		VoteAverage:      parseDecimal(tmdbMovie.VoteAverage),
		VoteCount:        tmdbParseOptionalInt32(tmdbMovie.VoteCount),
		Popularity:       parseDecimal(tmdbMovie.Popularity),
		PosterPath:       tmdbParseOptionalString(tmdbMovie.PosterPath),
		BackdropPath:     tmdbParseOptionalString(tmdbMovie.BackdropPath),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if tmdbMovie.ReleaseDate != "" {
		mov.Year = tmdbParseOptionalInt32Ptr(ExtractYear(tmdbMovie.ReleaseDate))
	}

	return mov
}

func (m *TMDbMapper) MapSearchResult(result *TMDbSearchResult) *Movie {
	mov := &Movie{
		ID:               uuid.New(),
		TMDbID:           tmdbParseOptionalInt32(result.ID),
		Title:            result.Title,
		OriginalTitle:    tmdbParseOptionalString(&result.OriginalTitle),
		OriginalLanguage: tmdbParseOptionalString(&result.OriginalLanguage),
		Overview:         tmdbParseOptionalString(&result.Overview),
		ReleaseDate:      parseReleaseDate(result.ReleaseDate),
		VoteAverage:      parseDecimal(result.VoteAverage),
		VoteCount:        tmdbParseOptionalInt32(result.VoteCount),
		Popularity:       parseDecimal(result.Popularity),
		PosterPath:       tmdbParseOptionalString(result.PosterPath),
		BackdropPath:     tmdbParseOptionalString(result.BackdropPath),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if result.ReleaseDate != "" {
		mov.Year = tmdbParseOptionalInt32Ptr(ExtractYear(result.ReleaseDate))
	}

	return mov
}

func (m *TMDbMapper) MapCredits(movieID uuid.UUID, credits *TMDbCredits) []MovieCredit {
	var result []MovieCredit

	for _, cast := range credits.Cast {
		result = append(result, MovieCredit{
			ID:           uuid.New(),
			MovieID:      movieID,
			TMDbPersonID: util.SafeIntToInt32(cast.ID),
			Name:         cast.Name,
			CreditType:   "cast",
			Character:    tmdbParseOptionalString(&cast.Character),
			Department:   nil,
			Job:          nil,
			CastOrder:    tmdbParseOptionalInt32Ptr(&cast.Order),
			ProfilePath:  tmdbParseOptionalString(cast.ProfilePath),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	for _, crew := range credits.Crew {
		result = append(result, MovieCredit{
			ID:           uuid.New(),
			MovieID:      movieID,
			TMDbPersonID: util.SafeIntToInt32(crew.ID),
			Name:         crew.Name,
			CreditType:   "crew",
			Character:    nil,
			Department:   tmdbParseOptionalString(&crew.Department),
			Job:          tmdbParseOptionalString(&crew.Job),
			CastOrder:    nil,
			ProfilePath:  tmdbParseOptionalString(crew.ProfilePath),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	return result
}

func (m *TMDbMapper) MapGenres(movieID uuid.UUID, genres []Genre) []MovieGenre {
	var result []MovieGenre

	for _, genre := range genres {
		result = append(result, MovieGenre{
			ID:          uuid.New(),
			MovieID:     movieID,
			TMDbGenreID: util.SafeIntToInt32(genre.ID),
			Name:        genre.Name,
			CreatedAt:   time.Now(),
		})
	}

	return result
}

func (m *TMDbMapper) MapCollection(collection *TMDbCollectionDetails) *MovieCollection {
	return &MovieCollection{
		ID:               uuid.New(),
		TMDbCollectionID: tmdbParseOptionalInt32(collection.ID),
		Name:             collection.Name,
		Overview:         tmdbParseOptionalString(&collection.Overview),
		PosterPath:       tmdbParseOptionalString(collection.PosterPath),
		BackdropPath:     tmdbParseOptionalString(collection.BackdropPath),
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

func tmdbParseOptionalString(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

func tmdbParseOptionalInt32(i int) *int32 {
	if i == 0 {
		return nil
	}
	val := util.SafeIntToInt32(i)
	return &val
}

func tmdbParseOptionalInt32Ptr(i *int) *int32 {
	if i == nil || *i == 0 {
		return nil
	}
	val := util.SafeIntToInt32(*i)
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
