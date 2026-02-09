package mal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapAnimeToTVShowSearchResult converts a MAL Anime to TVShowSearchResult.
func mapAnimeToTVShowSearchResult(a Anime) metadata.TVShowSearchResult {
	result := metadata.TVShowSearchResult{
		ProviderID:       strconv.Itoa(a.ID),
		Provider:         metadata.ProviderMAL,
		Name:             a.Title,
		OriginalLanguage: "ja",
		OriginCountries:  []string{"JP"},
		Overview:         a.Synopsis,
		Adult:            a.NSFW == "black", // MAL NSFW: white (safe), gray (suggestive), black (explicit)
	}

	// Use English title if available
	if a.AlternativeTitles.En != "" {
		result.Name = a.AlternativeTitles.En
	}
	if a.AlternativeTitles.Ja != "" {
		result.OriginalName = a.AlternativeTitles.Ja
	}

	if a.Mean != nil {
		result.VoteAverage = *a.Mean
	}
	result.VoteCount = a.NumScoringUsers

	if a.Popularity != nil {
		result.Popularity = float64(*a.Popularity)
	} else {
		result.Popularity = float64(a.NumListUsers)
	}

	result.FirstAirDate = parseDate(a.StartDate)
	if a.StartSeason != nil {
		result.Year = &a.StartSeason.Year
	} else if year := parseYear(a.StartDate); year > 0 {
		result.Year = &year
	}

	if a.MainPicture != nil {
		img := a.MainPicture.Large
		if img == "" {
			img = a.MainPicture.Medium
		}
		if img != "" {
			result.PosterPath = &img
		}
	}

	for _, g := range a.Genres {
		result.GenreIDs = append(result.GenreIDs, malGenreToStandardID(g.ID))
	}

	return result
}

// mapAnimeToTVShowMetadata converts a MAL Anime to TVShowMetadata.
func mapAnimeToTVShowMetadata(a *Anime) *metadata.TVShowMetadata {
	if a == nil {
		return nil
	}

	m := &metadata.TVShowMetadata{
		ProviderID:       strconv.Itoa(a.ID),
		Provider:         metadata.ProviderMAL,
		Name:             a.Title,
		OriginalLanguage: "ja",
		OriginCountries:  []string{"JP"},
		Adult:            a.NSFW == "black",
	}

	// Titles
	if a.AlternativeTitles.En != "" {
		m.Name = a.AlternativeTitles.En
	}
	if a.AlternativeTitles.Ja != "" {
		m.OriginalName = a.AlternativeTitles.Ja
	}

	if a.Synopsis != "" {
		m.Overview = &a.Synopsis
	}

	m.Status = mapStatus(a.Status)
	m.Type = mapMediaType(a.MediaType)
	m.InProduction = a.Status == "currently_airing"

	m.FirstAirDate = parseDate(a.StartDate)
	m.LastAirDate = parseDate(a.EndDate)

	m.NumberOfEpisodes = a.NumEpisodes

	// Average episode duration (MAL returns seconds)
	if a.AverageEpisodeDuration > 0 {
		m.EpisodeRuntime = []int{a.AverageEpisodeDuration / 60} // seconds → minutes
	}

	// Single season for anime
	if a.NumEpisodes > 0 {
		m.NumberOfSeasons = 1
	}

	// MAL score is already on 0-10 scale
	if a.Mean != nil {
		m.VoteAverage = *a.Mean
	}

	if a.Popularity != nil {
		m.Popularity = float64(*a.Popularity)
	} else {
		m.Popularity = float64(a.NumListUsers)
	}

	// External ratings
	if a.Mean != nil {
		m.ExternalRatings = append(m.ExternalRatings, metadata.ExternalRating{
			Source: "MyAnimeList",
			Value:  fmt.Sprintf("%.2f", *a.Mean),
			Score:  *a.Mean * 10, // Normalize to 0-100
		})
	}

	// Images
	if a.MainPicture != nil {
		img := a.MainPicture.Large
		if img == "" {
			img = a.MainPicture.Medium
		}
		if img != "" {
			m.PosterPath = &img
		}
	}

	// Homepage
	homepage := fmt.Sprintf("https://myanimelist.net/anime/%d", a.ID)
	m.Homepage = &homepage

	// Genres
	for _, g := range a.Genres {
		m.Genres = append(m.Genres, metadata.Genre{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	// Studios → Networks
	for _, s := range a.Studios {
		m.Networks = append(m.Networks, metadata.Network{
			ID:            s.ID,
			Name:          s.Name,
			OriginCountry: "JP",
		})
	}

	return m
}

// mapImages extracts images from a MAL Anime.
func mapImages(a *Anime) *metadata.Images {
	if a == nil {
		return nil
	}

	images := &metadata.Images{}

	// Main picture
	if a.MainPicture != nil {
		for _, url := range []string{a.MainPicture.Large, a.MainPicture.Medium} {
			if url != "" {
				images.Posters = append(images.Posters, metadata.Image{FilePath: url})
			}
		}
	}

	// Additional pictures
	for _, pic := range a.Pictures {
		if pic.Large != "" {
			images.Posters = append(images.Posters, metadata.Image{FilePath: pic.Large})
		} else if pic.Medium != "" {
			images.Posters = append(images.Posters, metadata.Image{FilePath: pic.Medium})
		}
	}

	if len(images.Posters) == 0 {
		return nil
	}
	return images
}

// Helper functions

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	// MAL dates can be "2017-10-23", "2017-10", or "2017"
	for _, layout := range []string{"2006-01-02", "2006-01", "2006"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return &t
		}
	}
	return nil
}

func parseYear(s string) int {
	if len(s) < 4 {
		return 0
	}
	year, err := strconv.Atoi(s[:4])
	if err != nil {
		return 0
	}
	return year
}

func mapStatus(s string) string {
	switch s {
	case "finished_airing":
		return "Ended"
	case "currently_airing":
		return "Returning Series"
	case "not_yet_aired":
		return "Planned"
	default:
		return s
	}
}

func mapMediaType(s string) string {
	switch s {
	case "tv":
		return "Scripted"
	case "movie":
		return "Movie"
	case "ova":
		return "OVA"
	case "ona":
		return "ONA"
	case "special":
		return "Special"
	case "music":
		return "Music"
	default:
		return strings.ToUpper(s)
	}
}

func mapRating(s string) string {
	switch s {
	case "g":
		return "G"
	case "pg":
		return "PG"
	case "pg_13":
		return "PG-13"
	case "r":
		return "R"
	case "r+":
		return "R+"
	case "rx":
		return "Rx"
	default:
		return s
	}
}

// malGenreToStandardID maps MAL genre IDs to TMDb-compatible genre IDs where applicable.
func malGenreToStandardID(malID int) int {
	// MAL genre IDs are already stable, so we use them directly.
	// Some map to TMDb TV genre IDs for consistency.
	switch malID {
	case 1: // Action
		return 10759
	case 2: // Adventure
		return 10759
	case 4: // Comedy
		return 35
	case 8: // Drama
		return 18
	case 10: // Fantasy
		return 10765
	case 14: // Horror
		return 27
	case 7: // Mystery
		return 9648
	case 22: // Romance
		return 10749
	case 24: // Sci-Fi
		return 10765
	default:
		return malID // Use MAL IDs directly as stable identifiers
	}
}
