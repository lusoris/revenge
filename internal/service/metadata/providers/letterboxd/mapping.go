package letterboxd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// mapFilmSummaryToSearchResult converts a Letterboxd FilmSummary to a MovieSearchResult.
func mapFilmSummaryToSearchResult(f *FilmSummary) metadata.MovieSearchResult {
	result := metadata.MovieSearchResult{
		ProviderID: f.ID,
		Provider:   metadata.ProviderLetterboxd,
		Title:      f.Name,
		Adult:      f.Adult,
	}

	if f.OriginalName != "" {
		result.OriginalTitle = f.OriginalName
	}

	if f.ReleaseYear > 0 {
		year := f.ReleaseYear
		result.Year = &year
		// Approximate release date from year
		t := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		result.ReleaseDate = &t
	}

	// Rating: Letterboxd uses 0.5-5.0, convert to 0-10 scale
	if f.Rating > 0 {
		result.VoteAverage = f.Rating * 2.0 // 5.0 → 10.0
	}

	// Poster: use the largest available size
	if f.Poster != nil {
		result.PosterPath = largestImageURL(f.Poster)
	}

	// Genre IDs
	for _, g := range f.Genres {
		if id := genreNameToID(g.Name); id > 0 {
			result.GenreIDs = append(result.GenreIDs, id)
		}
	}

	return result
}

// mapFilmToMetadata converts a full Letterboxd Film to MovieMetadata.
func mapFilmToMetadata(f *Film) metadata.MovieMetadata {
	m := metadata.MovieMetadata{
		ProviderID: f.ID,
		Provider:   metadata.ProviderLetterboxd,
		Title:      f.Name,
		Adult:      f.Adult,
	}

	if f.OriginalName != "" {
		m.OriginalTitle = f.OriginalName
	}

	if f.Tagline != "" {
		m.Tagline = &f.Tagline
	}

	if f.Description != "" {
		m.Overview = &f.Description
	}

	if f.ReleaseYear > 0 {
		t := time.Date(f.ReleaseYear, 1, 1, 0, 0, 0, 0, time.UTC)
		m.ReleaseDate = &t
	}

	if f.RunTime > 0 {
		runtime := util.SafeIntToInt32(f.RunTime)
		m.Runtime = &runtime
	}

	// Rating: convert 0.5-5.0 → 0-10 scale
	if f.Rating > 0 {
		m.VoteAverage = f.Rating * 2.0
		// Add Letterboxd as an ExternalRating (original 0.5-5.0 scale)
		m.ExternalRatings = append(m.ExternalRatings, metadata.ExternalRating{
			Source: "Letterboxd",
			Value:  fmt.Sprintf("%.1f/5", f.Rating),
			Score:  f.Rating * 20, // normalize to 0-100
		})
	}

	// Poster
	if f.Poster != nil {
		m.PosterPath = largestImageURL(f.Poster)
	}

	// Backdrop
	if f.Backdrop != nil {
		m.BackdropPath = largestImageURL(f.Backdrop)
	}

	// Trailer
	if f.Trailer != nil && f.Trailer.URL != "" {
		m.TrailerURL = &f.Trailer.URL
	}

	// Genres
	for _, g := range f.Genres {
		m.Genres = append(m.Genres, metadata.Genre{
			ID:   genreNameToID(g.Name),
			Name: g.Name,
		})
	}

	// Production countries
	for _, c := range f.Countries {
		m.ProductionCountries = append(m.ProductionCountries, metadata.ProductionCountry{
			ISOCode: c.Code,
			Name:    c.Name,
		})
	}

	// Spoken languages
	for _, l := range f.Languages {
		m.SpokenLanguages = append(m.SpokenLanguages, metadata.SpokenLanguage{
			ISOCode: l.Code,
			Name:    l.Name,
		})
	}

	// Original language
	if f.PrimaryLanguage != nil {
		m.OriginalLanguage = f.PrimaryLanguage.Code
	}

	// Extract external IDs from links
	tmdbID, imdbID := extractExternalIDs(f.Links)
	if imdbID != "" {
		m.IMDbID = &imdbID
	}
	if tmdbID > 0 {
		id32 := util.SafeIntToInt32(tmdbID)
		m.TMDbID = &id32
	}

	// Collection reference
	if f.FilmCollectionID != "" {
		m.Collection = &metadata.CollectionRef{
			Name: "", // We don't have the name from the film response
		}
	}

	// Homepage from links
	for _, link := range f.Links {
		if link.Type == "letterboxd" {
			m.Homepage = &link.URL
			break
		}
	}

	return m
}

// mapCredits extracts credits from the film's contributions.
func mapCredits(contributions []Contributions) metadata.Credits {
	credits := metadata.Credits{}

	for _, c := range contributions {
		for _, contributor := range c.Contributors {
			switch c.Type {
			case "Director", "CoDirector":
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        c.Type,
					Department: "Directing",
				})
			case "Actor":
				credits.Cast = append(credits.Cast, metadata.CastMember{
					Name:      contributor.Name,
					Character: contributor.CharacterName,
				})
			case "Writer", "OriginalWriter", "Story":
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        c.Type,
					Department: "Writing",
				})
			case "Producer", "ExecutiveProducer":
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        c.Type,
					Department: "Production",
				})
			case "Composer", "Songs":
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        c.Type,
					Department: "Sound",
				})
			case "Cinematography", "CameraOperator", "AdditionalPhotography":
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        c.Type,
					Department: "Camera",
				})
			case "Editor":
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        "Editor",
					Department: "Editing",
				})
			default:
				credits.Crew = append(credits.Crew, metadata.CrewMember{
					Name:       contributor.Name,
					Job:        c.Type,
					Department: mapDepartment(c.Type),
				})
			}
		}
	}

	return credits
}

// extractExternalIDs pulls TMDb and IMDb IDs from the links array.
func extractExternalIDs(links []Link) (tmdbID int, imdbID string) {
	for _, link := range links {
		switch link.Type {
		case "tmdb":
			if id, err := strconv.Atoi(link.ID); err == nil {
				tmdbID = id
			}
		case "imdb":
			imdbID = link.ID
		}
	}
	return
}

// largestImageURL returns a pointer to the URL of the largest image size.
func largestImageURL(img *Image) *string {
	if img == nil || len(img.Sizes) == 0 {
		return nil
	}

	largest := img.Sizes[0]
	for _, s := range img.Sizes[1:] {
		if s.Width*s.Height > largest.Width*largest.Height {
			largest = s
		}
	}

	if largest.URL == "" {
		return nil
	}

	return &largest.URL
}

// mapDepartment maps Letterboxd contribution types to department names.
func mapDepartment(contributionType string) string {
	switch contributionType {
	case "ProductionDesign", "ArtDirection", "SetDecoration":
		return "Art"
	case "Costumes", "MakeUp", "Hairstyling":
		return "Costume & Make-Up"
	case "VisualEffects", "SpecialEffects":
		return "Visual Effects"
	case "Lighting":
		return "Lighting"
	case "Sound":
		return "Sound"
	case "TitleDesign":
		return "Art"
	case "Stunts", "Choreography":
		return "Crew"
	case "Casting":
		return "Production"
	case "Studio":
		return "Production"
	default:
		return "Crew"
	}
}

// genreNameToID maps genre names to TMDb-compatible numeric IDs.
func genreNameToID(name string) int {
	genreMap := map[string]int{
		"Action":          28,
		"Adventure":       12,
		"Animation":       16,
		"Comedy":          35,
		"Crime":           80,
		"Documentary":     99,
		"Drama":           18,
		"Family":          10751,
		"Fantasy":         14,
		"History":         36,
		"Horror":          27,
		"Music":           10402,
		"Mystery":         9648,
		"Romance":         10749,
		"Science Fiction": 878,
		"TV Movie":        10770,
		"Thriller":        53,
		"War":             10752,
		"Western":         37,
	}

	if id, ok := genreMap[name]; ok {
		return id
	}

	// Case-insensitive fallback
	lower := strings.ToLower(name)
	for k, v := range genreMap {
		if strings.ToLower(k) == lower {
			return v
		}
	}

	return 0
}
