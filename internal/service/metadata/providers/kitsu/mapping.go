package kitsu

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
	"github.com/lusoris/revenge/internal/util"
)

// mapAnimeToTVShowSearchResult converts a Kitsu anime to TVShowSearchResult.
func mapAnimeToTVShowSearchResult(res ResourceObject[AnimeAttributes]) metadata.TVShowSearchResult {
	a := res.Attributes
	result := metadata.TVShowSearchResult{
		ProviderID:       res.ID,
		Provider:         metadata.ProviderKitsu,
		Name:             a.CanonicalTitle,
		OriginalLanguage: "ja",
		OriginCountries:  []string{"JP"},
		Adult:            a.NSFW,
	}

	// Use English title if available
	if en, ok := a.Titles["en"]; ok && en != "" {
		result.Name = en
	}
	// Original name from Japanese title
	if jp, ok := a.Titles["ja_jp"]; ok && jp != "" {
		result.OriginalName = jp
	} else if enJP, ok := a.Titles["en_jp"]; ok && enJP != "" {
		result.OriginalName = enJP
	}

	result.Overview = a.Synopsis
	result.Popularity = float64(a.UserCount)

	if a.AverageRating != nil {
		if rating, err := strconv.ParseFloat(*a.AverageRating, 64); err == nil {
			result.VoteAverage = rating / 10.0 // Kitsu uses 0-100, we use 0-10
		}
	}

	result.VoteCount = a.UserCount

	result.FirstAirDate = parseDate(a.StartDate)
	if a.StartDate != nil {
		if year, err := strconv.Atoi((*a.StartDate)[:4]); err == nil && len(*a.StartDate) >= 4 {
			result.Year = &year
		}
	}

	if a.PosterImage != nil {
		if img := bestImage(a.PosterImage); img != "" {
			result.PosterPath = &img
		}
	}

	if a.CoverImage != nil {
		if img := bestImage(a.CoverImage); img != "" {
			result.BackdropPath = &img
		}
	}

	return result
}

// mapAnimeToTVShowMetadata converts Kitsu anime details to TVShowMetadata.
func mapAnimeToTVShowMetadata(res *SingleResponse[AnimeAttributes]) *metadata.TVShowMetadata {
	if res == nil {
		return nil
	}

	a := res.Data.Attributes
	m := &metadata.TVShowMetadata{
		ProviderID:       res.Data.ID,
		Provider:         metadata.ProviderKitsu,
		OriginalLanguage: "ja",
		OriginCountries:  []string{"JP"},
		Adult:            a.NSFW,
	}

	// Titles
	m.Name = a.CanonicalTitle
	if en, ok := a.Titles["en"]; ok && en != "" {
		m.Name = en
	}
	if jp, ok := a.Titles["ja_jp"]; ok && jp != "" {
		m.OriginalName = jp
	} else if enJP, ok := a.Titles["en_jp"]; ok && enJP != "" {
		m.OriginalName = enJP
	}

	if a.Synopsis != "" {
		m.Overview = &a.Synopsis
	}

	m.Status = mapStatus(a.Status)
	m.Type = mapSubtype(a.Subtype)
	m.InProduction = a.Status == "current"

	m.FirstAirDate = parseDate(a.StartDate)
	m.LastAirDate = parseDate(a.EndDate)

	if a.EpisodeCount != nil {
		m.NumberOfEpisodes = *a.EpisodeCount
	}
	if a.EpisodeLength != nil {
		m.EpisodeRuntime = []int{*a.EpisodeLength}
	}

	// Kitsu anime = single season typically
	if a.EpisodeCount != nil && *a.EpisodeCount > 0 {
		m.NumberOfSeasons = 1
	}

	if a.AverageRating != nil {
		if rating, err := strconv.ParseFloat(*a.AverageRating, 64); err == nil {
			m.VoteAverage = rating / 10.0
		}
	}
	m.Popularity = float64(a.UserCount)

	// External ratings
	if a.AverageRating != nil {
		if rating, err := strconv.ParseFloat(*a.AverageRating, 64); err == nil {
			m.ExternalRatings = append(m.ExternalRatings, metadata.ExternalRating{
				Source: "Kitsu",
				Value:  fmt.Sprintf("%.1f%%", rating),
				Score:  rating,
			})
		}
	}

	// Images
	if a.PosterImage != nil {
		if img := bestImage(a.PosterImage); img != "" {
			m.PosterPath = &img
		}
	}
	if a.CoverImage != nil {
		if img := bestImage(a.CoverImage); img != "" {
			m.BackdropPath = &img
		}
	}

	// Trailer
	if a.YoutubeVideoID != nil && *a.YoutubeVideoID != "" {
		url := "https://www.youtube.com/watch?v=" + *a.YoutubeVideoID
		m.TrailerURL = &url
	}

	// Age rating as content info
	if a.AgeRating != nil {
		siteURL := "https://kitsu.io/anime/" + a.Slug
		m.Homepage = &siteURL
	}

	// Categories from included resources
	for _, inc := range res.Included {
		if inc.Type == "categories" {
			title, _ := inc.Attributes["title"].(string)
			if title != "" {
				m.Genres = append(m.Genres, metadata.Genre{
					ID:   categoryToGenreID(title),
					Name: title,
				})
			}
		}
	}

	// Extract external IDs from included mappings
	for _, inc := range res.Included {
		if inc.Type == "mappings" {
			site, _ := inc.Attributes["externalSite"].(string)
			extID, _ := inc.Attributes["externalId"].(string)
			if site == "" || extID == "" {
				continue
			}
			switch {
			case strings.Contains(site, "myanimelist"):
				// MAL ID — no direct field
			case strings.Contains(site, "thetvdb"):
				if tvdbID, err := strconv.Atoi(extID); err == nil {
					id := util.SafeIntToInt32(tvdbID)
					m.TVDbID = &id
				}
			}
		}
	}

	return m
}

// mapEpisodes converts Kitsu episodes for a season.
func mapEpisodesToSummary(episodes *ListResponse[EpisodeAttributes], seasonNum int) []metadata.EpisodeSummary {
	if episodes == nil {
		return nil
	}

	var result []metadata.EpisodeSummary
	for _, ep := range episodes.Data {
		a := ep.Attributes

		// Filter by season if needed
		if seasonNum > 0 && a.SeasonNumber != nil && *a.SeasonNumber != seasonNum {
			continue
		}

		es := metadata.EpisodeSummary{
			ProviderID: ep.ID,
			Name:       a.CanonicalTitle,
			AirDate:    parseDate(a.Airdate),
		}

		if a.Number != nil {
			es.EpisodeNumber = *a.Number
		}
		if a.Length != nil {
			rt := util.SafeIntToInt32(*a.Length)
			es.Runtime = &rt
		}
		if a.Synopsis != "" {
			es.Overview = &a.Synopsis
		}
		if a.Thumbnail != nil {
			if img := bestImage(a.Thumbnail); img != "" {
				es.StillPath = &img
			}
		}

		result = append(result, es)
	}
	return result
}

// mapEpisodeToMetadata converts a single Kitsu episode.
func mapEpisodeToMetadata(ep ResourceObject[EpisodeAttributes], animeID string) *metadata.EpisodeMetadata {
	a := ep.Attributes
	em := &metadata.EpisodeMetadata{
		ProviderID: ep.ID,
		Provider:   metadata.ProviderKitsu,
		ShowID:     animeID,
		Name:       a.CanonicalTitle,
		AirDate:    parseDate(a.Airdate),
	}

	if a.SeasonNumber != nil {
		em.SeasonNumber = *a.SeasonNumber
	} else {
		em.SeasonNumber = 1
	}
	if a.Number != nil {
		em.EpisodeNumber = *a.Number
	}
	if a.Length != nil {
		rt := util.SafeIntToInt32(*a.Length)
		em.Runtime = &rt
	}
	if a.Synopsis != "" {
		em.Overview = &a.Synopsis
	}
	if a.Thumbnail != nil {
		if img := bestImage(a.Thumbnail); img != "" {
			em.StillPath = &img
		}
	}

	return em
}

// mapImages extracts images from Kitsu anime.
func mapImages(res *SingleResponse[AnimeAttributes]) *metadata.Images {
	if res == nil {
		return nil
	}

	a := res.Data.Attributes
	images := &metadata.Images{}

	if a.PosterImage != nil {
		for _, url := range imageURLs(a.PosterImage) {
			images.Posters = append(images.Posters, metadata.Image{
				FilePath: url,
			})
		}
	}

	if a.CoverImage != nil {
		for _, url := range imageURLs(a.CoverImage) {
			images.Backdrops = append(images.Backdrops, metadata.Image{
				FilePath: url,
			})
		}
	}

	if len(images.Posters) == 0 && len(images.Backdrops) == 0 {
		return nil
	}
	return images
}

// mapMappingsToExternalIDs converts Kitsu mappings to ExternalIDs.
func mapMappingsToExternalIDs(mappings *ListResponse[MappingAttributes]) *metadata.ExternalIDs {
	if mappings == nil || len(mappings.Data) == 0 {
		return nil
	}

	ids := &metadata.ExternalIDs{}
	for _, m := range mappings.Data {
		a := m.Attributes
		switch {
		case strings.Contains(a.ExternalSite, "thetvdb"):
			if tvdbID, err := strconv.Atoi(a.ExternalID); err == nil {
				id := util.SafeIntToInt32(tvdbID)
				ids.TVDbID = &id
			}
		case strings.Contains(a.ExternalSite, "anilist"):
			// AniList ID — useful for cross-referencing
		case strings.Contains(a.ExternalSite, "myanimelist"):
			// MAL ID
		case strings.Contains(a.ExternalSite, "anidb"):
			// AniDB ID
		}
	}

	return ids
}

// Helper functions

func bestImage(img *ImageSet) string {
	if img == nil {
		return ""
	}
	if img.Original != nil && *img.Original != "" {
		return *img.Original
	}
	if img.Large != nil && *img.Large != "" {
		return *img.Large
	}
	if img.Medium != nil && *img.Medium != "" {
		return *img.Medium
	}
	if img.Small != nil && *img.Small != "" {
		return *img.Small
	}
	return ""
}

func imageURLs(img *ImageSet) []string {
	if img == nil {
		return nil
	}
	var urls []string
	for _, u := range []*string{img.Original, img.Large, img.Medium, img.Small, img.Tiny} {
		if u != nil && *u != "" {
			urls = append(urls, *u)
		}
	}
	return urls
}

func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}

func mapStatus(s string) string {
	switch s {
	case "finished":
		return "Ended"
	case "current":
		return "Returning Series"
	case "upcoming", "tba":
		return "Planned"
	case "unreleased":
		return "Planned"
	default:
		return s
	}
}

func mapSubtype(s string) string {
	switch s {
	case "TV":
		return "Scripted"
	case "movie":
		return "Movie"
	case "OVA":
		return "OVA"
	case "ONA":
		return "ONA"
	case "special":
		return "Special"
	case "music":
		return "Music"
	default:
		return s
	}
}

// categoryToGenreID maps Kitsu category names to stable IDs.
func categoryToGenreID(name string) int {
	switch strings.ToLower(name) {
	case "action":
		return 10759
	case "adventure":
		return 10759
	case "comedy":
		return 35
	case "drama":
		return 18
	case "fantasy":
		return 10765
	case "horror":
		return 27
	case "music":
		return 10402
	case "mystery":
		return 9648
	case "romance":
		return 10749
	case "science fiction", "sci-fi":
		return 10765
	case "thriller":
		return 53
	case "slice of life":
		return 90005
	case "sports":
		return 90006
	case "supernatural":
		return 90007
	case "mecha":
		return 90003
	case "psychological":
		return 90004
	default:
		return 0
	}
}
