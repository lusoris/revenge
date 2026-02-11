package simkl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// mapSearchResultToMovieSearchResult converts a Simkl SearchResult to MovieSearchResult.
func mapSearchResultToMovieSearchResult(r SearchResult) metadata.MovieSearchResult {
	result := metadata.MovieSearchResult{
		ProviderID: strconv.Itoa(r.IDs.Simkl),
		Provider:   metadata.ProviderSimkl,
		Title:      r.Title,
	}

	if r.Year > 0 {
		year := r.Year
		result.Year = &year
	}

	if r.Poster != "" {
		posterURL := ImageURL(r.Poster, "_ca")
		result.PosterPath = &posterURL
	}

	return result
}

// mapMovieToMetadata converts a Simkl Movie to MovieMetadata.
func mapMovieToMetadata(m *Movie) *metadata.MovieMetadata {
	if m == nil {
		return nil
	}

	md := &metadata.MovieMetadata{
		ProviderID:       strconv.Itoa(m.IDs.Simkl),
		Provider:         metadata.ProviderSimkl,
		Title:            m.Title,
		Status:           mapMovieStatus(m.Status),
	}

	if m.Overview != "" {
		md.Overview = &m.Overview
	}
	if m.Trailer != "" {
		md.TrailerURL = &m.Trailer
	}
	if m.Runtime > 0 {
		rt := util.SafeIntToInt32(m.Runtime)
		md.Runtime = &rt
	}

	// Certification (e.g. PG-13) is not mapped — MovieMetadata lacks a top-level
	// certification field. The data is available via Simkl's m.Certification.

	// Cross-referenced IDs
	if m.IDs.IMDb != "" {
		md.IMDbID = &m.IDs.IMDb
	}
	if m.IDs.TMDb > 0 {
		tmdbID := util.SafeIntToInt32(m.IDs.TMDb)
		md.TMDbID = &tmdbID
	}
	if m.IDs.TVDb > 0 {
		tvdbID := util.SafeIntToInt32(m.IDs.TVDb)
		md.TVDbID = &tvdbID
	}

	md.ReleaseDate = parseDate(m.ReleaseDate)

	if m.Ratings != nil && m.Ratings.Simkl != nil {
		md.VoteAverage = m.Ratings.Simkl.Rating
		if m.Ratings.Simkl.Votes > 0 {
			md.VoteCount = m.Ratings.Simkl.Votes
		}
	}

	// Add external ratings
	if m.Ratings != nil {
		// Simkl's own rating
		if m.Ratings.Simkl != nil && m.Ratings.Simkl.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "Simkl",
				Value:  fmt.Sprintf("%.1f/10", m.Ratings.Simkl.Rating),
				Score:  m.Ratings.Simkl.Rating * 10,
			})
		}
		if m.Ratings.IMDb != nil {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "IMDb",
				Value:  strconv.FormatFloat(m.Ratings.IMDb.Rating, 'f', 1, 64) + "/10",
				Score:  m.Ratings.IMDb.Rating * 10,
			})
		}
		if m.Ratings.MAL != nil {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "MAL",
				Value:  strconv.FormatFloat(m.Ratings.MAL.Rating, 'f', 1, 64) + "/10",
				Score:  m.Ratings.MAL.Rating * 10,
			})
		}
		if m.Ratings.Tmdb != nil && m.Ratings.Tmdb.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "TMDb",
				Value:  fmt.Sprintf("%.1f/10", m.Ratings.Tmdb.Rating),
				Score:  m.Ratings.Tmdb.Rating * 10,
			})
		}
		if m.Ratings.Trakt != nil && m.Ratings.Trakt.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "Trakt",
				Value:  fmt.Sprintf("%.1f/10", m.Ratings.Trakt.Rating),
				Score:  m.Ratings.Trakt.Rating * 10,
			})
		}
		if m.Ratings.Letterboxd != nil && m.Ratings.Letterboxd.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "Letterboxd",
				Value:  fmt.Sprintf("%.1f/5", m.Ratings.Letterboxd.Rating),
				Score:  m.Ratings.Letterboxd.Rating * 20,
			})
		}
	}

	for _, g := range m.Genres {
		md.Genres = append(md.Genres, metadata.Genre{
			ID:   genreNameToID(g),
			Name: g,
		})
	}

	if m.Country != "" {
		md.ProductionCountries = []metadata.ProductionCountry{
			{ISOCode: strings.ToUpper(m.Country)},
		}
	}

	// Images
	if m.Poster != "" {
		posterURL := ImageURL(m.Poster, "_m")
		md.PosterPath = &posterURL
	}
	if m.Fanart != "" {
		fanartURL := FanartURL(m.Fanart, "_medium")
		md.BackdropPath = &fanartURL
	}

	return md
}

// mapSearchResultToTVShowSearchResult converts a Simkl SearchResult to TVShowSearchResult.
func mapSearchResultToTVShowSearchResult(r SearchResult) metadata.TVShowSearchResult {
	result := metadata.TVShowSearchResult{
		ProviderID: strconv.Itoa(r.IDs.Simkl),
		Provider:   metadata.ProviderSimkl,
		Name:       r.Title,
	}

	if r.Year > 0 {
		year := r.Year
		result.Year = &year
	}

	if r.Poster != "" {
		posterURL := ImageURL(r.Poster, "_ca")
		result.PosterPath = &posterURL
	}

	return result
}

// mapShowToMetadata converts a Simkl Show to TVShowMetadata.
func mapShowToMetadata(s *Show) *metadata.TVShowMetadata {
	if s == nil {
		return nil
	}

	md := &metadata.TVShowMetadata{
		ProviderID: strconv.Itoa(s.IDs.Simkl),
		Provider:   metadata.ProviderSimkl,
		Name:       s.Title,
		Status:     mapShowStatus(s.Status),
	}

	if s.ENTitle != "" {
		md.OriginalName = s.Title
		md.Name = s.ENTitle
	}

	if s.Overview != "" {
		md.Overview = &s.Overview
	}
	if s.Trailer != "" {
		md.TrailerURL = &s.Trailer
	}
	if s.Runtime > 0 {
		md.EpisodeRuntime = []int{s.Runtime}
	}
	if s.TotalEpisodes > 0 {
		md.NumberOfEpisodes = s.TotalEpisodes
	}

	md.InProduction = strings.EqualFold(s.Status, "airing") || strings.EqualFold(s.Status, "ongoing")

	// Cross-referenced IDs
	if s.IDs.IMDb != "" {
		md.IMDbID = &s.IDs.IMDb
	}
	if s.IDs.TMDb > 0 {
		tmdbID := util.SafeIntToInt32(s.IDs.TMDb)
		md.TMDbID = &tmdbID
	}
	if s.IDs.TVDb > 0 {
		tvdbID := util.SafeIntToInt32(s.IDs.TVDb)
		md.TVDbID = &tvdbID
	}

	if s.Network != "" {
		md.Networks = []metadata.Network{
			{Name: s.Network},
		}
	}

	if s.Ratings != nil && s.Ratings.Simkl != nil {
		md.VoteAverage = s.Ratings.Simkl.Rating
		if s.Ratings.Simkl.Votes > 0 {
			md.VoteCount = s.Ratings.Simkl.Votes
		}
	}

	if s.Ratings != nil {
		// Simkl's own rating
		if s.Ratings.Simkl != nil && s.Ratings.Simkl.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "Simkl",
				Value:  fmt.Sprintf("%.1f/10", s.Ratings.Simkl.Rating),
				Score:  s.Ratings.Simkl.Rating * 10,
			})
		}
		if s.Ratings.IMDb != nil {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "IMDb",
				Value:  strconv.FormatFloat(s.Ratings.IMDb.Rating, 'f', 1, 64) + "/10",
				Score:  s.Ratings.IMDb.Rating * 10,
			})
		}
		if s.Ratings.MAL != nil {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "MAL",
				Value:  strconv.FormatFloat(s.Ratings.MAL.Rating, 'f', 1, 64) + "/10",
				Score:  s.Ratings.MAL.Rating * 10,
			})
		}
		if s.Ratings.Tmdb != nil && s.Ratings.Tmdb.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "TMDb",
				Value:  fmt.Sprintf("%.1f/10", s.Ratings.Tmdb.Rating),
				Score:  s.Ratings.Tmdb.Rating * 10,
			})
		}
		if s.Ratings.Trakt != nil && s.Ratings.Trakt.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "Trakt",
				Value:  fmt.Sprintf("%.1f/10", s.Ratings.Trakt.Rating),
				Score:  s.Ratings.Trakt.Rating * 10,
			})
		}
		if s.Ratings.Letterboxd != nil && s.Ratings.Letterboxd.Rating > 0 {
			md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
				Source: "Letterboxd",
				Value:  fmt.Sprintf("%.1f/5", s.Ratings.Letterboxd.Rating),
				Score:  s.Ratings.Letterboxd.Rating * 20,
			})
		}
	}

	for _, g := range s.Genres {
		md.Genres = append(md.Genres, metadata.Genre{
			ID:   genreNameToID(g),
			Name: g,
		})
	}

	if s.Country != "" {
		md.OriginCountries = []string{strings.ToUpper(s.Country)}
	}

	// Images
	if s.Poster != "" {
		posterURL := ImageURL(s.Poster, "_m")
		md.PosterPath = &posterURL
	}
	if s.Fanart != "" {
		fanartURL := FanartURL(s.Fanart, "_medium")
		md.BackdropPath = &fanartURL
	}

	return md
}

// mapEpisodesToSummaries converts Simkl episodes to EpisodeSummary for a given season.
func mapEpisodesToSummaries(episodes []Episode, seasonNum int) []metadata.EpisodeSummary {
	var result []metadata.EpisodeSummary
	for _, ep := range episodes {
		if ep.Season != seasonNum {
			continue
		}
		es := metadata.EpisodeSummary{
			ProviderID:    strconv.Itoa(ep.IDs.Simkl),
			EpisodeNumber: ep.Episode,
			Name:          ep.Title,
			AirDate:       ep.Date,
		}
		if ep.Img != "" {
			imgURL := fmt.Sprintf("%s/episodes/%s_w.jpg", imageBaseURL, ep.Img)
			es.StillPath = &imgURL
		}
		result = append(result, es)
	}
	return result
}

// mapExternalIDs creates ExternalIDs from Simkl IDs.
func mapExternalIDs(ids IDs) *metadata.ExternalIDs {
	ext := &metadata.ExternalIDs{}
	if ids.IMDb != "" {
		ext.IMDbID = &ids.IMDb
	}
	if ids.TMDb > 0 {
		tmdbID := util.SafeIntToInt32(ids.TMDb)
		ext.TMDbID = &tmdbID
	}
	if ids.TVDb > 0 {
		tvdbID := util.SafeIntToInt32(ids.TVDb)
		ext.TVDbID = &tvdbID
	}
	return ext
}

// --- Helpers ---

// parseDate parses a date string (YYYY-MM-DD) to *time.Time.
func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

// mapMovieStatus maps Simkl movie status to a standard string.
func mapMovieStatus(status string) string {
	switch strings.ToLower(status) {
	case "released":
		return "Released"
	case "upcoming":
		return "Upcoming"
	case "rumored":
		return "Rumored"
	case "planned":
		return "Planned"
	case "in production":
		return "In Production"
	case "post production":
		return "Post Production"
	case "cancelled":
		return "Canceled"
	default:
		return status
	}
}

// mapShowStatus maps Simkl show status to a standard string.
func mapShowStatus(status string) string {
	switch strings.ToLower(status) {
	case "airing", "returning series", "ongoing":
		return "Returning Series"
	case "ended":
		return "Ended"
	case "cancelled", "canceled":
		return "Canceled"
	case "tba", "upcoming", "planned", "in production":
		return "In Production"
	default:
		return status
	}
}

// genreNameToID maps a genre name to a numeric ID.
func genreNameToID(genre string) int {
	genreMap := map[string]int{
		"action":          28,
		"adventure":       12,
		"animation":       16,
		"anime":           7878,
		"comedy":          35,
		"crime":           80,
		"documentary":     99,
		"drama":           18,
		"family":          10751,
		"fantasy":         14,
		"history":         36,
		"horror":          27,
		"music":           10402,
		"mystery":         9648,
		"romance":         10749,
		"science-fiction":  878,
		"sci-fi":          878,
		"science fiction":  878,
		"thriller":        53,
		"war":             10752,
		"western":         37,
	}
	if id, ok := genreMap[strings.ToLower(genre)]; ok {
		return id
	}
	return 0
}
