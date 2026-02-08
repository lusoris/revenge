package anidb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapAnimeToTVShowSearchResult converts a title dump match to a TVShowSearchResult.
// Note: Search results are minimal since the title dump only has aid+title.
// Full details require a separate GetAnime call.
func mapTitleToTVShowSearchResult(entry TitleDumpEntry) metadata.TVShowSearchResult {
	return metadata.TVShowSearchResult{
		ProviderID:       strconv.Itoa(entry.AID),
		Provider:         metadata.ProviderAniDB,
		Name:             entry.Title,
		OriginalLanguage: "ja",
		OriginCountries:  []string{"JP"},
	}
}

// mapAnimeToTVShowMetadata converts an AniDB AnimeResponse to TVShowMetadata.
func mapAnimeToTVShowMetadata(a *AnimeResponse) *metadata.TVShowMetadata {
	if a == nil {
		return nil
	}

	m := &metadata.TVShowMetadata{
		ProviderID:       strconv.Itoa(a.ID),
		Provider:         metadata.ProviderAniDB,
		OriginalLanguage: "ja",
		OriginCountries:  []string{"JP"},
		Adult:            a.Restricted,
		NumberOfEpisodes: a.EpCount,
	}

	// Titles — find the best title
	m.Name = bestTitle(a.Titles, "en")
	m.OriginalName = bestTitle(a.Titles, "ja")
	if m.Name == "" {
		m.Name = bestTitle(a.Titles, "x-jat") // Romaji
	}

	// Description
	if a.Description != "" {
		desc := cleanDescription(a.Description)
		m.Overview = &desc
	}

	// Type
	m.Type = mapType(a.Type)
	m.Status = mapStatusFromDates(a.StartDate, a.EndDate)

	// Dates
	m.FirstAirDate = parseDate(a.StartDate)
	m.LastAirDate = parseDate(a.EndDate)

	// Single season for anime
	if a.EpCount > 0 {
		m.NumberOfSeasons = 1
	}

	// Episode runtime (from first regular episode)
	for _, ep := range a.Episodes.Episode {
		if ep.EpNo.Type == 1 && ep.Length > 0 {
			m.EpisodeRuntime = []int{ep.Length}
			break
		}
	}

	// Ratings
	if a.Ratings.Permanent.Count > 0 {
		m.VoteAverage = a.Ratings.Permanent.Value
		m.ExternalRatings = append(m.ExternalRatings, metadata.ExternalRating{
			Source: "AniDB (Permanent)",
			Value:  fmt.Sprintf("%.2f", a.Ratings.Permanent.Value),
			Score:  a.Ratings.Permanent.Value * 10, // 0-10 → 0-100
		})
	}
	if a.Ratings.Temporary.Count > 0 {
		m.ExternalRatings = append(m.ExternalRatings, metadata.ExternalRating{
			Source: "AniDB (Temporary)",
			Value:  fmt.Sprintf("%.2f", a.Ratings.Temporary.Value),
			Score:  a.Ratings.Temporary.Value * 10,
		})
	}

	// Cover image
	if a.Picture != "" {
		img := ImageBaseURL + a.Picture
		m.PosterPath = &img
	}

	// Homepage
	if a.URL != "" {
		m.Homepage = &a.URL
	}

	// Tags → Genres (use non-spoiler, high-weight tags)
	for _, tag := range a.Tags.Tag {
		if tag.GlobalSpoiler || tag.LocalSpoiler || tag.Weight < 200 {
			continue
		}
		m.Genres = append(m.Genres, metadata.Genre{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	// Studios from creators with "Animation Work" type
	for _, cr := range a.Creators.Name {
		if strings.Contains(strings.ToLower(cr.Type), "animation") ||
			strings.Contains(strings.ToLower(cr.Type), "work") {
			m.Networks = append(m.Networks, metadata.Network{
				ID:            cr.ID,
				Name:          cr.Text,
				OriginCountry: "JP",
			})
		}
	}

	return m
}

// mapCredits converts AniDB characters and creators to Credits.
func mapCredits(a *AnimeResponse) *metadata.Credits {
	if a == nil {
		return nil
	}

	credits := &metadata.Credits{}

	for i, ch := range a.Characters.Character {
		// Add voice actors as cast
		for _, va := range ch.Seiyuu {
			member := metadata.CastMember{
				ProviderID: strconv.Itoa(va.ID),
				Name:       va.Text,
				Character:  ch.Name,
				Order:      i,
			}
			if ch.Gender != "" {
				member.Gender = mapGender(ch.Gender)
			}
			if va.Picture != "" {
				img := ImageBaseURL + va.Picture
				member.ProfilePath = &img
			}
			credits.Cast = append(credits.Cast, member)
		}

		// If no seiyuu, add the character itself
		if len(ch.Seiyuu) == 0 && ch.Name != "" {
			member := metadata.CastMember{
				ProviderID: strconv.Itoa(ch.ID),
				Name:       ch.Name,
				Character:  ch.Type,
				Order:      i,
				Gender:     mapGender(ch.Gender),
			}
			if ch.Picture != "" {
				img := ImageBaseURL + ch.Picture
				member.ProfilePath = &img
			}
			credits.Cast = append(credits.Cast, member)
		}
	}

	// Creators as crew
	for _, cr := range a.Creators.Name {
		member := metadata.CrewMember{
			ProviderID: strconv.Itoa(cr.ID),
			Name:       cr.Text,
			Job:        cr.Type,
			Department: mapDepartment(cr.Type),
		}
		credits.Crew = append(credits.Crew, member)
	}

	return credits
}

// mapImages extracts images from an AniDB anime.
func mapImages(a *AnimeResponse) *metadata.Images {
	if a == nil || a.Picture == "" {
		return nil
	}

	img := ImageBaseURL + a.Picture
	return &metadata.Images{
		Posters: []metadata.Image{{FilePath: img}},
	}
}

// mapEpisodes converts AniDB episodes to metadata format.
func mapEpisodes(a *AnimeResponse, seasonNum int) []metadata.EpisodeSummary {
	if a == nil {
		return nil
	}

	var result []metadata.EpisodeSummary
	for _, ep := range a.Episodes.Episode {
		// Only regular episodes (type 1)
		if ep.EpNo.Type != 1 {
			continue
		}

		epNum, err := strconv.Atoi(ep.EpNo.Text)
		if err != nil {
			continue
		}

		es := metadata.EpisodeSummary{
			ProviderID:    strconv.Itoa(ep.ID),
			EpisodeNumber: epNum,
			AirDate:       parseDate(ep.Airdate),
		}

		// Find English title, fall back to romaji
		for _, t := range ep.Title {
			if t.Lang == "en" {
				es.Name = t.Text
				break
			}
		}
		if es.Name == "" {
			for _, t := range ep.Title {
				if t.Lang == "x-jat" {
					es.Name = t.Text
					break
				}
			}
		}
		if es.Name == "" && len(ep.Title) > 0 {
			es.Name = ep.Title[0].Text
		}

		if ep.Length > 0 {
			rt := int32(ep.Length)
			es.Runtime = &rt
		}

		if ep.Rating != nil {
			es.VoteAverage = ep.Rating.Value
		}

		result = append(result, es)
	}
	return result
}

// mapEpisodeToMetadata converts a single AniDB episode.
func mapEpisodeToMetadata(ep Episode, showID string) *metadata.EpisodeMetadata {
	epNum, err := strconv.Atoi(ep.EpNo.Text)
	if err != nil {
		return nil
	}

	em := &metadata.EpisodeMetadata{
		ProviderID:    strconv.Itoa(ep.ID),
		Provider:      metadata.ProviderAniDB,
		ShowID:        showID,
		SeasonNumber:  1,
		EpisodeNumber: epNum,
		AirDate:       parseDate(ep.Airdate),
	}

	// Find best title
	for _, t := range ep.Title {
		if t.Lang == "en" {
			em.Name = t.Text
			break
		}
	}
	if em.Name == "" {
		for _, t := range ep.Title {
			if t.Lang == "x-jat" {
				em.Name = t.Text
				break
			}
		}
	}
	if em.Name == "" && len(ep.Title) > 0 {
		em.Name = ep.Title[0].Text
	}

	if ep.Length > 0 {
		rt := int32(ep.Length)
		em.Runtime = &rt
	}

	if ep.Rating != nil {
		em.VoteAverage = ep.Rating.Value
	}

	return em
}

// findExternalIDs extracts external ID cross-references from AniDB resources.
func findExternalIDs(a *AnimeResponse) *metadata.ExternalIDs {
	if a == nil {
		return nil
	}

	ids := &metadata.ExternalIDs{}

	for _, res := range a.Resources.Resource {
		if len(res.ExternalID) == 0 || len(res.ExternalID[0].Identifier) == 0 {
			continue
		}
		extID := res.ExternalID[0].Identifier[0]

		switch res.Type {
		case 1: // ANN (Anime News Network)
			// Skip — no field for ANN
		case 2: // MAL
			// MAL ID — stored as metadata cross-reference
		case 4: // Official URL
			// Already have URL field
		case 6: // Wikipedia EN
			ids.WikidataID = &extID
		}
	}

	return ids
}

// Helper functions

func bestTitle(titles Titles, lang string) string {
	// First try official titles in the requested language
	for _, t := range titles.Title {
		if t.Lang == lang && t.Type == "official" {
			return t.Text
		}
	}
	// Then main title
	for _, t := range titles.Title {
		if t.Type == "main" {
			return t.Text
		}
	}
	// Then any title in the requested language
	for _, t := range titles.Title {
		if t.Lang == lang {
			return t.Text
		}
	}
	// Fallback to first title
	if len(titles.Title) > 0 {
		return titles.Title[0].Text
	}
	return ""
}

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", "2006-01", "2006"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return &t
		}
	}
	return nil
}

func cleanDescription(s string) string {
	// AniDB uses custom markup like [url=...], [b], etc. Strip it.
	s = strings.ReplaceAll(s, "[/b]", "")
	s = strings.ReplaceAll(s, "[/i]", "")
	s = strings.ReplaceAll(s, "[/u]", "")
	s = strings.ReplaceAll(s, "[b]", "")
	s = strings.ReplaceAll(s, "[i]", "")
	s = strings.ReplaceAll(s, "[u]", "")

	// Remove [url=...]...[/url] and similar
	for {
		start := strings.Index(s, "[url=")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "]")
		if end == -1 {
			break
		}
		s = s[:start] + s[start+end+1:]
	}
	s = strings.ReplaceAll(s, "[/url]", "")

	// Clean up source references like "Source: ANN"
	return strings.TrimSpace(s)
}

func mapType(t string) string {
	switch t {
	case "TV Series":
		return "Scripted"
	case "Movie":
		return "Movie"
	case "OVA":
		return "OVA"
	case "Web":
		return "ONA"
	case "TV Special":
		return "Special"
	case "Music Video":
		return "Music"
	default:
		return t
	}
}

func mapStatusFromDates(start, end string) string {
	if start == "" {
		return "Planned"
	}
	if end == "" {
		now := time.Now()
		startDate := parseDate(start)
		if startDate != nil && startDate.After(now) {
			return "Planned"
		}
		return "Returning Series"
	}
	return "Ended"
}

func mapGender(g string) int {
	switch strings.ToLower(g) {
	case "female":
		return 1
	case "male":
		return 2
	default:
		return 0
	}
}

func mapDepartment(creatorType string) string {
	lower := strings.ToLower(creatorType)
	switch {
	case strings.Contains(lower, "direction"):
		return "Directing"
	case strings.Contains(lower, "music"):
		return "Sound"
	case strings.Contains(lower, "character design"):
		return "Art"
	case strings.Contains(lower, "animation"):
		return "Art"
	case strings.Contains(lower, "series composition"), strings.Contains(lower, "original work"):
		return "Writing"
	default:
		return "Production"
	}
}
