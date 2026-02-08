package omdb

import (
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapMovieMetadata converts an OMDb response to MovieMetadata.
func mapMovieMetadata(resp *Response) *metadata.MovieMetadata {
	if resp == nil || resp.Type != "movie" {
		return nil
	}

	m := &metadata.MovieMetadata{
		ProviderID: resp.IMDbID,
		Provider:   metadata.ProviderOMDb,
		Title:      resp.Title,
	}

	if resp.IMDbID != "" {
		m.IMDbID = &resp.IMDbID
	}

	if resp.Plot != "" && resp.Plot != "N/A" {
		m.Overview = &resp.Plot
	}

	m.Runtime = parseRuntime(resp.Runtime)
	m.ReleaseDate = parseDate(resp.Released)
	m.Genres = parseGenres(resp.Genre)

	if resp.Poster != "" && resp.Poster != "N/A" {
		m.PosterPath = &resp.Poster
	}

	// Parse IMDb rating as VoteAverage
	if rating, err := strconv.ParseFloat(resp.IMDbRating, 64); err == nil {
		m.VoteAverage = rating
	}
	if votes, err := strconv.Atoi(strings.ReplaceAll(resp.IMDbVotes, ",", "")); err == nil {
		m.VoteCount = votes
	}

	// External ratings (the main value of OMDb)
	m.ExternalRatings = mapExternalRatings(resp)

	return m
}

// mapTVShowMetadata converts an OMDb response to TVShowMetadata.
func mapTVShowMetadata(resp *Response) *metadata.TVShowMetadata {
	if resp == nil || resp.Type != "series" {
		return nil
	}

	m := &metadata.TVShowMetadata{
		ProviderID: resp.IMDbID,
		Provider:   metadata.ProviderOMDb,
		Name:       resp.Title,
	}

	if resp.IMDbID != "" {
		m.IMDbID = &resp.IMDbID
	}

	if resp.Plot != "" && resp.Plot != "N/A" {
		m.Overview = &resp.Plot
	}

	m.FirstAirDate = parseDate(resp.Released)
	m.Genres = parseGenres(resp.Genre)

	if resp.Poster != "" && resp.Poster != "N/A" {
		m.PosterPath = &resp.Poster
	}

	if seasons, err := strconv.Atoi(resp.TotalSeasons); err == nil {
		m.NumberOfSeasons = seasons
	}

	// Parse IMDb rating
	if rating, err := strconv.ParseFloat(resp.IMDbRating, 64); err == nil {
		m.VoteAverage = rating
	}
	if votes, err := strconv.Atoi(strings.ReplaceAll(resp.IMDbVotes, ",", "")); err == nil {
		m.VoteCount = votes
	}

	// External ratings
	m.ExternalRatings = mapExternalRatings(resp)

	return m
}

// mapMovieSearchResults converts an OMDb search response to MovieSearchResults.
func mapMovieSearchResults(resp *SearchResponse) []metadata.MovieSearchResult {
	if resp == nil {
		return nil
	}

	var results []metadata.MovieSearchResult
	for _, r := range resp.Search {
		if r.Type != "movie" {
			continue
		}
		result := metadata.MovieSearchResult{
			ProviderID: r.IMDbID,
			Provider:   metadata.ProviderOMDb,
			Title:      r.Title,
		}
		if year, err := strconv.Atoi(r.Year); err == nil {
			result.Year = &year
		}
		if r.Poster != "" && r.Poster != "N/A" {
			result.PosterPath = &r.Poster
		}
		results = append(results, result)
	}
	return results
}

// mapTVShowSearchResults converts an OMDb search response to TVShowSearchResults.
func mapTVShowSearchResults(resp *SearchResponse) []metadata.TVShowSearchResult {
	if resp == nil {
		return nil
	}

	var results []metadata.TVShowSearchResult
	for _, r := range resp.Search {
		if r.Type != "series" {
			continue
		}
		result := metadata.TVShowSearchResult{
			ProviderID: r.IMDbID,
			Provider:   metadata.ProviderOMDb,
			Name:       r.Title,
		}
		if year, err := strconv.Atoi(r.Year); err == nil {
			result.Year = &year
		}
		if r.Poster != "" && r.Poster != "N/A" {
			result.PosterPath = &r.Poster
		}
		results = append(results, result)
	}
	return results
}

// mapExternalRatings extracts all ratings from an OMDb response.
func mapExternalRatings(resp *Response) []metadata.ExternalRating {
	if resp == nil {
		return nil
	}

	var ratings []metadata.ExternalRating
	for _, r := range resp.Ratings {
		rating := metadata.ExternalRating{
			Source: r.Source,
			Value:  r.Value,
			Score:  normalizeScore(r.Source, r.Value),
		}
		ratings = append(ratings, rating)
	}
	return ratings
}

// normalizeScore converts various rating formats to a 0-100 scale.
func normalizeScore(source, value string) float64 {
	switch source {
	case "Internet Movie Database":
		// "8.8/10" → 88.0
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			if score, err := strconv.ParseFloat(parts[0], 64); err == nil {
				return score * 10
			}
		}
	case "Rotten Tomatoes":
		// "96%" → 96.0
		value = strings.TrimSuffix(value, "%")
		if score, err := strconv.ParseFloat(value, 64); err == nil {
			return score
		}
	case "Metacritic":
		// "90/100" → 90.0
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			if score, err := strconv.ParseFloat(parts[0], 64); err == nil {
				return score
			}
		}
	}
	return 0
}

// parseRuntime parses "142 min" to int32 minutes.
func parseRuntime(s string) *int32 {
	s = strings.TrimSuffix(s, " min")
	if s == "" || s == "N/A" {
		return nil
	}
	if mins, err := strconv.Atoi(s); err == nil {
		v := int32(mins)
		return &v
	}
	return nil
}

// parseDate parses "01 Jan 2020" to time.Time.
func parseDate(s string) *time.Time {
	if s == "" || s == "N/A" {
		return nil
	}
	t, err := time.Parse("02 Jan 2006", s)
	if err != nil {
		return nil
	}
	return &t
}

// parseGenres splits "Action, Drama, Thriller" into Genre slices.
func parseGenres(s string) []metadata.Genre {
	if s == "" || s == "N/A" {
		return nil
	}
	parts := strings.Split(s, ", ")
	genres := make([]metadata.Genre, 0, len(parts))
	for _, p := range parts {
		genres = append(genres, metadata.Genre{Name: strings.TrimSpace(p)})
	}
	return genres
}
