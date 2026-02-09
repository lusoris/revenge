package trakt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapMovieToSearchResult converts a Trakt Movie to MovieSearchResult.
func mapMovieToSearchResult(m *Movie) metadata.MovieSearchResult {
	if m == nil {
		return metadata.MovieSearchResult{}
	}

	result := metadata.MovieSearchResult{
		ProviderID:       strconv.Itoa(m.IDs.Trakt),
		Provider:         metadata.ProviderTrakt,
		Title:            m.Title,
		Overview:         m.Overview,
		OriginalLanguage: m.Language,
		VoteAverage:      m.Rating,
		VoteCount:        m.Votes,
	}

	if m.Year > 0 {
		year := m.Year
		result.Year = &year
	}

	result.ReleaseDate = parseDate(m.Released)

	for _, g := range m.Genres {
		result.GenreIDs = append(result.GenreIDs, genreNameToID(g))
	}

	return result
}

// mapMovieToMetadata converts a Trakt Movie to MovieMetadata.
func mapMovieToMetadata(m *Movie) *metadata.MovieMetadata {
	if m == nil {
		return nil
	}

	md := &metadata.MovieMetadata{
		ProviderID:       strconv.Itoa(m.IDs.Trakt),
		Provider:         metadata.ProviderTrakt,
		Title:            m.Title,
		OriginalLanguage: m.Language,
		VoteAverage:      m.Rating,
		VoteCount:        m.Votes,
		Status:           mapStatus(m.Status),
	}

	if m.Tagline != "" {
		md.Tagline = &m.Tagline
	}
	if m.Overview != "" {
		md.Overview = &m.Overview
	}
	if m.Homepage != "" {
		md.Homepage = &m.Homepage
	}
	if m.Trailer != "" {
		md.TrailerURL = &m.Trailer
	}
	if m.Runtime > 0 {
		rt := int32(m.Runtime)
		md.Runtime = &rt
	}

	// Cross-reference IDs
	if m.IDs.IMDb != "" {
		md.IMDbID = &m.IDs.IMDb
	}
	if m.IDs.TMDb > 0 {
		tmdbID := int32(m.IDs.TMDb)
		md.TMDbID = &tmdbID
	}
	if m.IDs.TVDb > 0 {
		tvdbID := int32(m.IDs.TVDb)
		md.TVDbID = &tvdbID
	}

	md.ReleaseDate = parseDate(m.Released)

	// Add Trakt community rating as ExternalRating
	if m.Rating > 0 {
		md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
			Source: "Trakt",
			Value:  fmt.Sprintf("%.1f/10", m.Rating),
			Score:  m.Rating * 10,
		})
	}

	for _, g := range m.Genres {
		md.Genres = append(md.Genres, metadata.Genre{
			ID:   genreNameToID(g),
			Name: capitalizeGenre(g),
		})
	}

	if m.Country != "" {
		md.ProductionCountries = []metadata.ProductionCountry{
			{ISOCode: strings.ToUpper(m.Country)},
		}
	}

	// Certification is not an ExternalRating — it's a content rating (e.g. PG-13)
	// TODO: store in a dedicated Certification field when available

	return md
}

// mapShowToSearchResult converts a Trakt Show to TVShowSearchResult.
func mapShowToSearchResult(s *Show) metadata.TVShowSearchResult {
	if s == nil {
		return metadata.TVShowSearchResult{}
	}

	result := metadata.TVShowSearchResult{
		ProviderID:       strconv.Itoa(s.IDs.Trakt),
		Provider:         metadata.ProviderTrakt,
		Name:             s.Title,
		Overview:         s.Overview,
		OriginalLanguage: s.Language,
		VoteAverage:      s.Rating,
		VoteCount:        s.Votes,
	}

	if s.Year > 0 {
		year := s.Year
		result.Year = &year
	}

	if s.FirstAired != nil {
		result.FirstAirDate = s.FirstAired
	}

	if s.Country != "" {
		result.OriginCountries = []string{strings.ToUpper(s.Country)}
	}

	for _, g := range s.Genres {
		result.GenreIDs = append(result.GenreIDs, genreNameToID(g))
	}

	return result
}

// mapShowToMetadata converts a Trakt Show to TVShowMetadata.
func mapShowToMetadata(s *Show) *metadata.TVShowMetadata {
	if s == nil {
		return nil
	}

	md := &metadata.TVShowMetadata{
		ProviderID:       strconv.Itoa(s.IDs.Trakt),
		Provider:         metadata.ProviderTrakt,
		Name:             s.Title,
		OriginalLanguage: s.Language,
		VoteAverage:      s.Rating,
		VoteCount:        s.Votes,
		Status:           mapShowStatus(s.Status),
	}

	if s.Overview != "" {
		md.Overview = &s.Overview
	}
	if s.Homepage != "" {
		md.Homepage = &s.Homepage
	}
	if s.Trailer != "" {
		md.TrailerURL = &s.Trailer
	}
	if s.Runtime > 0 {
		md.EpisodeRuntime = []int{s.Runtime}
	}

	// Cross-reference IDs
	if s.IDs.IMDb != "" {
		md.IMDbID = &s.IDs.IMDb
	}
	if s.IDs.TMDb > 0 {
		tmdbID := int32(s.IDs.TMDb)
		md.TMDbID = &tmdbID
	}
	if s.IDs.TVDb > 0 {
		tvdbID := int32(s.IDs.TVDb)
		md.TVDbID = &tvdbID
	}

	if s.FirstAired != nil {
		md.FirstAirDate = s.FirstAired
	}

	md.InProduction = strings.EqualFold(s.Status, "returning series")

	if s.Country != "" {
		md.OriginCountries = []string{strings.ToUpper(s.Country)}
	}

	if s.Network != "" {
		md.Networks = []metadata.Network{
			{Name: s.Network},
		}
	}

	for _, g := range s.Genres {
		md.Genres = append(md.Genres, metadata.Genre{
			ID:   genreNameToID(g),
			Name: capitalizeGenre(g),
		})
	}

	// Add Trakt community rating as ExternalRating
	if s.Rating > 0 {
		md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
			Source: "Trakt",
			Value:  fmt.Sprintf("%.1f/10", s.Rating),
			Score:  s.Rating * 10,
		})
	}

	// Certification is not an ExternalRating — it's a content rating (e.g. TV-MA)
	// TODO: store in a dedicated Certification field when available

	return md
}

// mapSeasons converts Trakt seasons to SeasonSummary.
func mapSeasons(seasons []Season) []metadata.SeasonSummary {
	if len(seasons) == 0 {
		return nil
	}

	result := make([]metadata.SeasonSummary, 0, len(seasons))
	for _, s := range seasons {
		ss := metadata.SeasonSummary{
			ProviderID:   strconv.Itoa(s.IDs.Trakt),
			SeasonNumber: s.Number,
			Name:         s.Title,
			EpisodeCount: s.EpisodeCount,
			AirDate:      s.FirstAired,
			VoteAverage:  s.Rating,
		}
		if s.Overview != "" {
			ss.Overview = &s.Overview
		}
		result = append(result, ss)
	}
	return result
}

// mapEpisodesToSummaries converts Trakt episodes to EpisodeSummary.
func mapEpisodesToSummaries(episodes []Episode) []metadata.EpisodeSummary {
	if len(episodes) == 0 {
		return nil
	}

	result := make([]metadata.EpisodeSummary, 0, len(episodes))
	for _, ep := range episodes {
		es := metadata.EpisodeSummary{
			ProviderID:    strconv.Itoa(ep.IDs.Trakt),
			EpisodeNumber: ep.Number,
			Name:          ep.Title,
			AirDate:       ep.FirstAired,
			VoteAverage:   ep.Rating,
			VoteCount:     ep.Votes,
		}
		if ep.Overview != "" {
			es.Overview = &ep.Overview
		}
		if ep.Runtime > 0 {
			rt := int32(ep.Runtime)
			es.Runtime = &rt
		}
		result = append(result, es)
	}
	return result
}

// mapCredits converts Trakt Credits to metadata Credits.
func mapCredits(c *Credits) *metadata.Credits {
	if c == nil {
		return nil
	}

	credits := &metadata.Credits{}

	for i, cm := range c.Cast {
		character := ""
		if len(cm.Characters) > 0 {
			character = strings.Join(cm.Characters, ", ")
		}
		credits.Cast = append(credits.Cast, metadata.CastMember{
			ProviderID: strconv.Itoa(cm.Person.IDs.Trakt),
			Name:       cm.Person.Name,
			Character:  character,
			Order:      i,
		})
	}

	for dept, members := range c.Crew {
		for _, cm := range members {
			job := ""
			if len(cm.Jobs) > 0 {
				job = strings.Join(cm.Jobs, ", ")
			}
			credits.Crew = append(credits.Crew, metadata.CrewMember{
				ProviderID: strconv.Itoa(cm.Person.IDs.Trakt),
				Name:       cm.Person.Name,
				Job:        job,
				Department: capitalizeGenre(dept),
			})
		}
	}

	return credits
}

// mapTranslations converts Trakt translations to metadata translations.
func mapTranslations(translations []Translation) []metadata.Translation {
	if len(translations) == 0 {
		return nil
	}

	result := make([]metadata.Translation, 0, len(translations))
	for _, t := range translations {
		mt := metadata.Translation{
			Language: t.Language,
			ISOCode:  t.Country,
			Data: &metadata.TranslationData{
				Title:    t.Title,
				Overview: t.Overview,
				Tagline:  t.Tagline,
			},
		}
		result = append(result, mt)
	}
	return result
}

// mapExternalIDs creates ExternalIDs from Trakt IDs.
func mapExternalIDs(ids IDs) *metadata.ExternalIDs {
	ext := &metadata.ExternalIDs{}
	if ids.IMDb != "" {
		ext.IMDbID = &ids.IMDb
	}
	if ids.TMDb > 0 {
		tmdbID := int32(ids.TMDb)
		ext.TMDbID = &tmdbID
	}
	if ids.TVDb > 0 {
		tvdbID := int32(ids.TVDb)
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

// mapStatus maps Trakt movie status to a standard status string.
func mapStatus(status string) string {
	switch strings.ToLower(status) {
	case "released":
		return "Released"
	case "in production":
		return "In Production"
	case "post production":
		return "Post Production"
	case "planned":
		return "Planned"
	case "rumored":
		return "Rumored"
	case "canceled":
		return "Canceled"
	default:
		return status
	}
}

// mapShowStatus maps Trakt show status to a standard status string.
func mapShowStatus(status string) string {
	switch strings.ToLower(status) {
	case "returning series":
		return "Returning Series"
	case "ended":
		return "Ended"
	case "canceled":
		return "Canceled"
	case "in production":
		return "In Production"
	case "planned":
		return "Planned"
	default:
		return status
	}
}

// genreNameToID maps a Trakt genre slug to a numeric ID.
// IDs are arbitrary and unique per genre for internal use.
func genreNameToID(genre string) int {
	genreMap := map[string]int{
		"action":        28,
		"adventure":     12,
		"animation":     16,
		"anime":         7878,
		"comedy":        35,
		"crime":         80,
		"documentary":   99,
		"drama":         18,
		"family":        10751,
		"fantasy":       14,
		"history":       36,
		"holiday":       10770,
		"horror":        27,
		"music":         10402,
		"musical":       10402,
		"mystery":       9648,
		"romance":       10749,
		"science-fiction": 878,
		"short":         7777,
		"sports":        6666,
		"superhero":     5555,
		"suspense":      53,
		"thriller":      53,
		"war":           10752,
		"western":       37,
	}
	if id, ok := genreMap[strings.ToLower(genre)]; ok {
		return id
	}
	return 0
}

// capitalizeGenre capitalizes a genre slug.
func capitalizeGenre(s string) string {
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// AiredEpisodes is a field on TVShowMetadata that may not exist; fall back to NumEpisodes.
// We add a convenience accessor.
var _ = fmt.Sprintf // ensure fmt is used
