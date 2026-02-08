package tvmaze

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// stripHTML removes HTML tags from a string.
func stripHTML(s string) string {
	return strings.TrimSpace(htmlTagRegex.ReplaceAllString(s, ""))
}

// mapShowToTVShowSearchResult converts a TVmaze Show to TVShowSearchResult.
func mapShowToTVShowSearchResult(show Show) metadata.TVShowSearchResult {
	result := metadata.TVShowSearchResult{
		ProviderID:       strconv.Itoa(show.ID),
		Provider:         metadata.ProviderTVmaze,
		Name:             show.Name,
		OriginalLanguage: show.Language,
		OriginCountries:  mapOriginCountries(show),
	}

	if show.Summary != nil {
		overview := stripHTML(*show.Summary)
		result.Overview = overview
	}

	if show.Rating.Average != nil {
		result.VoteAverage = *show.Rating.Average
	}

	result.FirstAirDate = parseDate(show.Premiered)
	if show.Premiered != nil {
		if year, err := strconv.Atoi((*show.Premiered)[:4]); err == nil {
			result.Year = &year
		}
	}

	if show.Image != nil && show.Image.Original != "" {
		result.PosterPath = &show.Image.Original
	}

	for _, g := range show.Genres {
		result.GenreIDs = append(result.GenreIDs, genreNameToID(g))
	}

	return result
}

// mapShowToTVShowMetadata converts a TVmaze Show to TVShowMetadata.
func mapShowToTVShowMetadata(show *Show) *metadata.TVShowMetadata {
	if show == nil {
		return nil
	}

	m := &metadata.TVShowMetadata{
		ProviderID:       strconv.Itoa(show.ID),
		Provider:         metadata.ProviderTVmaze,
		Name:             show.Name,
		OriginalLanguage: show.Language,
		Status:           show.Status,
		Type:             show.Type,
		OriginCountries:  mapOriginCountries(*show),
	}

	if show.Externals.IMDb != nil {
		m.IMDbID = show.Externals.IMDb
	}
	if show.Externals.TVDb != nil {
		tvdbID := int32(*show.Externals.TVDb)
		m.TVDbID = &tvdbID
	}

	if show.Summary != nil {
		overview := stripHTML(*show.Summary)
		m.Overview = &overview
	}

	m.FirstAirDate = parseDate(show.Premiered)
	m.LastAirDate = parseDate(show.Ended)
	m.InProduction = show.Status == "Running"

	if show.Runtime != nil {
		m.EpisodeRuntime = []int{*show.Runtime}
	} else if show.AverageRuntime != nil {
		m.EpisodeRuntime = []int{*show.AverageRuntime}
	}

	if show.Rating.Average != nil {
		m.VoteAverage = *show.Rating.Average
	}

	if show.Image != nil && show.Image.Original != "" {
		m.PosterPath = &show.Image.Original
	}

	if show.OfficialSite != nil {
		m.Homepage = show.OfficialSite
	}

	// Map genres
	for _, g := range show.Genres {
		m.Genres = append(m.Genres, metadata.Genre{
			ID:   genreNameToID(g),
			Name: g,
		})
	}

	// Map network
	if show.Network != nil {
		n := metadata.Network{
			ID:   show.Network.ID,
			Name: show.Network.Name,
		}
		if show.Network.Country != nil {
			n.OriginCountry = show.Network.Country.Code
		}
		m.Networks = []metadata.Network{n}
	} else if show.WebChannel != nil {
		n := metadata.Network{
			ID:   show.WebChannel.ID,
			Name: show.WebChannel.Name,
		}
		if show.WebChannel.Country != nil {
			n.OriginCountry = show.WebChannel.Country.Code
		}
		m.Networks = []metadata.Network{n}
	}

	return m
}

// mapSeasons converts TVmaze seasons to SeasonSummary.
func mapSeasons(seasons []Season) []metadata.SeasonSummary {
	if len(seasons) == 0 {
		return nil
	}

	result := make([]metadata.SeasonSummary, 0, len(seasons))
	for _, s := range seasons {
		ss := metadata.SeasonSummary{
			ProviderID:   strconv.Itoa(s.ID),
			SeasonNumber: s.Number,
			Name:         s.Name,
			AirDate:      parseDate(s.PremiereDate),
		}
		if s.EpisodeOrder != nil {
			ss.EpisodeCount = *s.EpisodeOrder
		}
		if s.Summary != nil {
			overview := stripHTML(*s.Summary)
			ss.Overview = &overview
		}
		if s.Image != nil && s.Image.Original != "" {
			ss.PosterPath = &s.Image.Original
		}
		result = append(result, ss)
	}
	return result
}

// mapEpisodes converts TVmaze episodes to EpisodeSummary for a specific season.
func mapEpisodes(episodes []Episode, seasonNum int) []metadata.EpisodeSummary {
	var result []metadata.EpisodeSummary
	for _, ep := range episodes {
		if ep.Season != seasonNum {
			continue
		}
		es := metadata.EpisodeSummary{
			ProviderID: strconv.Itoa(ep.ID),
			Name:       ep.Name,
			AirDate:    parseAirdate(ep.Airdate),
		}
		if ep.Number != nil {
			es.EpisodeNumber = *ep.Number
		}
		if ep.Runtime != nil {
			rt := int32(*ep.Runtime)
			es.Runtime = &rt
		}
		if ep.Rating.Average != nil {
			es.VoteAverage = *ep.Rating.Average
		}
		if ep.Summary != nil {
			overview := stripHTML(*ep.Summary)
			es.Overview = &overview
		}
		if ep.Image != nil && ep.Image.Original != "" {
			es.StillPath = &ep.Image.Original
		}
		result = append(result, es)
	}
	return result
}

// mapCast converts TVmaze cast to Credits.
func mapCast(cast []CastMember, crew []CrewMember) *metadata.Credits {
	credits := &metadata.Credits{}

	for i, cm := range cast {
		member := metadata.CastMember{
			ProviderID: strconv.Itoa(cm.Person.ID),
			Name:       cm.Person.Name,
			Character:  cm.Character.Name,
			Order:      i,
			Gender:     mapGender(cm.Person.Gender),
		}
		if cm.Person.Image != nil && cm.Person.Image.Original != "" {
			member.ProfilePath = &cm.Person.Image.Original
		}
		credits.Cast = append(credits.Cast, member)
	}

	for _, cm := range crew {
		member := metadata.CrewMember{
			ProviderID: strconv.Itoa(cm.Person.ID),
			Name:       cm.Person.Name,
			Job:        cm.Type,
			Department: mapDepartment(cm.Type),
			Gender:     mapGender(cm.Person.Gender),
		}
		if cm.Person.Image != nil && cm.Person.Image.Original != "" {
			member.ProfilePath = &cm.Person.Image.Original
		}
		credits.Crew = append(credits.Crew, member)
	}

	return credits
}

// mapImages converts TVmaze ShowImages to metadata.Images.
func mapImages(imgs []ShowImage) *metadata.Images {
	if len(imgs) == 0 {
		return nil
	}

	images := &metadata.Images{}
	for _, img := range imgs {
		if img.Resolutions.Original == nil {
			continue
		}
		mi := metadata.Image{
			FilePath: img.Resolutions.Original.URL,
			Width:    img.Resolutions.Original.Width,
			Height:   img.Resolutions.Original.Height,
		}
		if img.Resolutions.Original.Width > 0 && img.Resolutions.Original.Height > 0 {
			mi.AspectRatio = float64(img.Resolutions.Original.Width) / float64(img.Resolutions.Original.Height)
		}

		switch img.Type {
		case "poster":
			images.Posters = append(images.Posters, mi)
		case "background":
			images.Backdrops = append(images.Backdrops, mi)
		case "banner":
			images.Backdrops = append(images.Backdrops, mi)
		case "typography":
			images.Logos = append(images.Logos, mi)
		}
	}
	return images
}

// mapExternalIDs creates ExternalIDs from TVmaze externals.
func mapExternalIDs(show *Show) *metadata.ExternalIDs {
	if show == nil {
		return nil
	}
	ids := &metadata.ExternalIDs{
		IMDbID: show.Externals.IMDb,
	}
	if show.Externals.TVDb != nil {
		tvdbID := int32(*show.Externals.TVDb)
		ids.TVDbID = &tvdbID
	}
	if show.Externals.TVRage != nil {
		rageID := int32(*show.Externals.TVRage)
		ids.TVRageID = &rageID
	}
	return ids
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

func parseAirdate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
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
	default:
		return 0
	}
}

func mapDepartment(crewType string) string {
	switch crewType {
	case "Creator", "Developer":
		return "Writing"
	case "Executive Producer", "Producer", "Co-Executive Producer":
		return "Production"
	case "Director":
		return "Directing"
	default:
		return "Production"
	}
}

func mapOriginCountries(show Show) []string {
	if show.Network != nil && show.Network.Country != nil {
		return []string{show.Network.Country.Code}
	}
	if show.WebChannel != nil && show.WebChannel.Country != nil {
		return []string{show.WebChannel.Country.Code}
	}
	return nil
}

// genreNameToID maps genre names to stable IDs (matching TMDb convention where possible).
func genreNameToID(name string) int {
	switch name {
	case "Action":
		return 10759
	case "Adventure":
		return 10759
	case "Animation":
		return 16
	case "Comedy":
		return 35
	case "Crime":
		return 80
	case "Documentary":
		return 99
	case "Drama":
		return 18
	case "Family":
		return 10751
	case "Fantasy":
		return 10765
	case "History":
		return 36
	case "Horror":
		return 27
	case "Music":
		return 10402
	case "Mystery":
		return 9648
	case "Romance":
		return 10749
	case "Science-Fiction":
		return 10765
	case "Thriller":
		return 53
	case "War":
		return 10768
	case "Western":
		return 37
	default:
		return 0
	}
}
