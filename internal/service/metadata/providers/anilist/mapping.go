package anilist

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapMediaToTVShowSearchResult converts an AniList Media to TVShowSearchResult.
func mapMediaToTVShowSearchResult(m Media) metadata.TVShowSearchResult {
	result := metadata.TVShowSearchResult{
		ProviderID: strconv.Itoa(m.ID),
		Provider:   metadata.ProviderAniList,
		Adult:      m.IsAdult,
	}

	// Use English title if available, fall back to romaji, then native.
	result.Name = bestTitle(m.Title)
	if m.Title.Romaji != nil && m.Title.English != nil && *m.Title.Romaji != *m.Title.English {
		result.OriginalName = *m.Title.Romaji
	} else if m.Title.Native != nil {
		result.OriginalName = *m.Title.Native
	}

	if m.CountryOfOrigin != nil {
		result.OriginalLanguage = strings.ToLower(*m.CountryOfOrigin)
		result.OriginCountries = []string{*m.CountryOfOrigin}
	} else {
		result.OriginalLanguage = "ja"
		result.OriginCountries = []string{"JP"}
	}

	if m.Description != nil {
		result.Overview = *m.Description
	}

	if m.AverageScore != nil {
		result.VoteAverage = float64(*m.AverageScore) / 10.0 // AniList uses 0-100, we use 0-10
	}

	result.Popularity = float64(m.Popularity)
	result.FirstAirDate = fuzzyDateToTime(m.StartDate)

	if m.StartDate.Year != nil {
		year := *m.StartDate.Year
		result.Year = &year
	}

	// Use best available cover image
	if img := bestCoverImage(m.CoverImage); img != "" {
		result.PosterPath = &img
	}

	if m.BannerImage != nil {
		result.BackdropPath = m.BannerImage
	}

	for _, g := range m.Genres {
		result.GenreIDs = append(result.GenreIDs, animeGenreToID(g))
	}

	return result
}

// mapMediaToTVShowMetadata converts an AniList Media to TVShowMetadata.
func mapMediaToTVShowMetadata(m *Media) *metadata.TVShowMetadata {
	if m == nil {
		return nil
	}

	md := &metadata.TVShowMetadata{
		ProviderID: strconv.Itoa(m.ID),
		Provider:   metadata.ProviderAniList,
		Adult:      m.IsAdult,
	}

	// Titles
	md.Name = bestTitle(m.Title)
	if m.Title.Native != nil {
		md.OriginalName = *m.Title.Native
	} else if m.Title.Romaji != nil {
		md.OriginalName = *m.Title.Romaji
	}

	if m.CountryOfOrigin != nil {
		md.OriginalLanguage = strings.ToLower(*m.CountryOfOrigin)
		md.OriginCountries = []string{*m.CountryOfOrigin}
	} else {
		md.OriginalLanguage = "ja"
		md.OriginCountries = []string{"JP"}
	}

	if m.Description != nil {
		md.Overview = m.Description
	}

	// Map AniList status to standard status
	md.Status = mapStatus(m.Status)
	md.Type = mapFormat(m.Format)
	md.InProduction = m.Status == "RELEASING"

	// Dates
	md.FirstAirDate = fuzzyDateToTime(m.StartDate)
	md.LastAirDate = fuzzyDateToTime(m.EndDate)

	// Episode info
	if m.Episodes != nil {
		md.NumberOfEpisodes = *m.Episodes
	}
	if m.Duration != nil {
		md.EpisodeRuntime = []int{*m.Duration}
	}

	// AniList counts anime by cour, single "season" for most anime.
	// Set season count to 1 for standard anime, could be refined with relations.
	if m.Episodes != nil && *m.Episodes > 0 {
		md.NumberOfSeasons = 1
	}

	// Ratings (AniList 0-100 → 0-10 scale)
	if m.AverageScore != nil {
		md.VoteAverage = float64(*m.AverageScore) / 10.0
	}
	md.Popularity = float64(m.Popularity)

	// External ratings
	if m.MeanScore != nil {
		md.ExternalRatings = append(md.ExternalRatings, metadata.ExternalRating{
			Source: "AniList",
			Value:  fmt.Sprintf("%d%%", *m.MeanScore),
			Score:  float64(*m.MeanScore),
		})
	}

	// Images
	if img := bestCoverImage(m.CoverImage); img != "" {
		md.PosterPath = &img
	}
	if m.BannerImage != nil {
		md.BackdropPath = m.BannerImage
	}

	// Homepage
	if m.SiteURL != "" {
		md.Homepage = &m.SiteURL
	}

	// Trailer
	if m.Trailer != nil && m.Trailer.ID != nil && m.Trailer.Site != nil {
		if *m.Trailer.Site == "youtube" {
			url := "https://www.youtube.com/watch?v=" + *m.Trailer.ID
			md.TrailerURL = &url
		}
	}

	// Genres
	for _, g := range m.Genres {
		md.Genres = append(md.Genres, metadata.Genre{
			ID:   animeGenreToID(g),
			Name: g,
		})
	}

	// Studios → Networks
	for _, edge := range m.Studios.Edges {
		n := metadata.Network{
			ID:            edge.Node.ID,
			Name:          edge.Node.Name,
			OriginCountry: "JP",
		}
		md.Networks = append(md.Networks, n)
	}

	// Map MAL ID as external ID
	if m.IDMal != nil {
		malID := int32(*m.IDMal)
		externalIDs := findExternalIDs(m)
		externalIDs.TMDbID = nil // We don't have TMDb mapping from AniList
		_ = malID               // stored via ExternalIDs in provider interface
	}

	return md
}

// mapCredits converts AniList characters and staff to Credits.
func mapCredits(m *Media) *metadata.Credits {
	if m == nil {
		return nil
	}

	credits := &metadata.Credits{}

	for i, edge := range m.Characters.Edges {
		var charName string
		if edge.Node.Name.Full != nil {
			charName = *edge.Node.Name.Full
		}

		// Main voice actor as the "cast" member
		for _, va := range edge.VoiceActors {
			var vaName string
			if va.Name.Full != nil {
				vaName = *va.Name.Full
			}

			member := metadata.CastMember{
				ProviderID: strconv.Itoa(va.ID),
				Name:       vaName,
				Character:  charName,
				Order:      i,
				Gender:     mapGender(va.Gender),
			}
			if va.Image.Large != nil {
				member.ProfilePath = va.Image.Large
			}
			credits.Cast = append(credits.Cast, member)
		}

		// If no voice actors, still record the character
		if len(edge.VoiceActors) == 0 && charName != "" {
			member := metadata.CastMember{
				ProviderID: strconv.Itoa(edge.Node.ID),
				Name:       charName,
				Character:  string(edge.Role),
				Order:      i,
				Gender:     mapGender(edge.Node.Gender),
			}
			if edge.Node.Image.Large != nil {
				member.ProfilePath = edge.Node.Image.Large
			}
			credits.Cast = append(credits.Cast, member)
		}
	}

	for _, edge := range m.Staff.Edges {
		var name string
		if edge.Node.Name.Full != nil {
			name = *edge.Node.Name.Full
		}

		member := metadata.CrewMember{
			ProviderID: strconv.Itoa(edge.Node.ID),
			Name:       name,
			Job:        edge.Role,
			Department: mapDepartment(edge.Role),
			Gender:     mapGender(edge.Node.Gender),
		}
		if edge.Node.Image.Large != nil {
			member.ProfilePath = edge.Node.Image.Large
		}
		credits.Crew = append(credits.Crew, member)
	}

	return credits
}

// mapImages extracts images from an AniList Media.
func mapImages(m *Media) *metadata.Images {
	if m == nil {
		return nil
	}

	images := &metadata.Images{}

	// Cover images as posters
	for _, url := range []string{
		safeStr(m.CoverImage.ExtraLarge),
		safeStr(m.CoverImage.Large),
		safeStr(m.CoverImage.Medium),
	} {
		if url != "" {
			images.Posters = append(images.Posters, metadata.Image{
				FilePath: url,
			})
		}
	}

	// Banner as backdrop
	if m.BannerImage != nil && *m.BannerImage != "" {
		images.Backdrops = append(images.Backdrops, metadata.Image{
			FilePath: *m.BannerImage,
		})
	}

	if len(images.Posters) == 0 && len(images.Backdrops) == 0 {
		return nil
	}

	return images
}

// findExternalIDs extracts external IDs from an AniList Media.
func findExternalIDs(m *Media) *metadata.ExternalIDs {
	if m == nil {
		return nil
	}

	ids := &metadata.ExternalIDs{}

	// MAL ID
	if m.IDMal != nil {
		// Store MAL ID — no direct field in ExternalIDs, but we can use it for cross-referencing.
		// We could store via WikidataID or similar, but for now we skip MAL-specific field.
	}

	// Parse external links for known sites
	for _, link := range m.ExternalLinks {
		switch link.Site {
		case "IMDb":
			if link.URL != nil {
				// Extract IMDb ID from URL
				imdbID := extractIMDbID(*link.URL)
				if imdbID != "" {
					ids.IMDbID = &imdbID
				}
			}
		case "Twitter", "X":
			if link.URL != nil {
				ids.TwitterID = link.URL
			}
		case "YouTube":
			if link.URL != nil {
				ids.YouTubeID = link.URL
			}
		}
	}

	return ids
}

// Helper functions

func bestTitle(t MediaTitle) string {
	if t.English != nil && *t.English != "" {
		return *t.English
	}
	if t.Romaji != nil && *t.Romaji != "" {
		return *t.Romaji
	}
	if t.UserPreferred != nil && *t.UserPreferred != "" {
		return *t.UserPreferred
	}
	if t.Native != nil {
		return *t.Native
	}
	return ""
}

func bestCoverImage(ci CoverImage) string {
	if ci.ExtraLarge != nil && *ci.ExtraLarge != "" {
		return *ci.ExtraLarge
	}
	if ci.Large != nil && *ci.Large != "" {
		return *ci.Large
	}
	if ci.Medium != nil && *ci.Medium != "" {
		return *ci.Medium
	}
	return ""
}

func fuzzyDateToTime(fd FuzzyDate) *time.Time {
	if fd.Year == nil {
		return nil
	}
	year := *fd.Year
	month := time.January
	day := 1
	if fd.Month != nil {
		month = time.Month(*fd.Month)
	}
	if fd.Day != nil {
		day = *fd.Day
	}
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &t
}

func mapStatus(s string) string {
	switch s {
	case "FINISHED":
		return "Ended"
	case "RELEASING":
		return "Returning Series"
	case "NOT_YET_RELEASED":
		return "Planned"
	case "CANCELLED":
		return "Canceled"
	case "HIATUS":
		return "Returning Series"
	default:
		return s
	}
}

func mapFormat(f string) string {
	switch f {
	case "TV":
		return "Scripted"
	case "TV_SHORT":
		return "Scripted"
	case "MOVIE":
		return "Movie"
	case "SPECIAL":
		return "Special"
	case "OVA":
		return "OVA"
	case "ONA":
		return "ONA"
	case "MUSIC":
		return "Music"
	default:
		return f
	}
}

func mapGender(g *string) int {
	if g == nil {
		return 0
	}
	switch *g {
	case "Female":
		return 1
	case "Male":
		return 2
	case "Non-binary":
		return 3
	default:
		return 0
	}
}

func mapDepartment(role string) string {
	lower := strings.ToLower(role)
	switch {
	case strings.Contains(lower, "director"):
		return "Directing"
	case strings.Contains(lower, "producer"):
		return "Production"
	case strings.Contains(lower, "writer"), strings.Contains(lower, "script"),
		strings.Contains(lower, "creator"), strings.Contains(lower, "story"),
		strings.Contains(lower, "composition"):
		return "Writing"
	case strings.Contains(lower, "music"), strings.Contains(lower, "sound"):
		return "Sound"
	case strings.Contains(lower, "art"), strings.Contains(lower, "design"),
		strings.Contains(lower, "animation"), strings.Contains(lower, "character design"):
		return "Art"
	default:
		return "Production"
	}
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func extractIMDbID(url string) string {
	// URLs like https://www.imdb.com/title/tt1234567/
	parts := strings.Split(url, "/")
	for _, p := range parts {
		if strings.HasPrefix(p, "tt") {
			return p
		}
	}
	return ""
}

// animeGenreToID maps AniList genres to stable IDs (using TMDb TV genre convention where possible).
func animeGenreToID(name string) int {
	switch name {
	case "Action":
		return 10759
	case "Adventure":
		return 10759
	case "Comedy":
		return 35
	case "Drama":
		return 18
	case "Ecchi":
		return 90001 // Custom
	case "Fantasy":
		return 10765
	case "Horror":
		return 27
	case "Mahou Shoujo":
		return 90002 // Custom
	case "Mecha":
		return 90003 // Custom
	case "Music":
		return 10402
	case "Mystery":
		return 9648
	case "Psychological":
		return 90004 // Custom
	case "Romance":
		return 10749
	case "Sci-Fi":
		return 10765
	case "Slice of Life":
		return 90005 // Custom
	case "Sports":
		return 90006 // Custom
	case "Supernatural":
		return 90007 // Custom
	case "Thriller":
		return 53
	default:
		return 0
	}
}
