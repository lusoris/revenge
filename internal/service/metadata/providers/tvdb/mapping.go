package tvdb

import (
	"strconv"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapTVSearchResult converts a TVDb search result to metadata type.
func mapTVSearchResult(r *SearchResult) metadata.TVShowSearchResult {
	result := metadata.TVShowSearchResult{
		ProviderID: r.TVDbID,
		Provider:   metadata.ProviderTVDb,
		Name:       r.Name,
		Overview:   r.Overview,
	}

	if r.ImageURL != nil {
		result.PosterPath = r.ImageURL
	}

	if r.Year != "" {
		if year, err := strconv.Atoi(r.Year); err == nil {
			result.Year = &year
		}
	}

	if r.FirstAirTime != "" {
		result.FirstAirDate = parseDate(r.FirstAirTime)
	}

	return result
}

// mapTVShowMetadata converts a TVDb series response to metadata type.
func mapTVShowMetadata(r *SeriesResponse) *metadata.TVShowMetadata {
	tvdbID := int32(r.ID)
	result := &metadata.TVShowMetadata{
		ProviderID:       strconv.Itoa(r.ID),
		Provider:         metadata.ProviderTVDb,
		TVDbID:           &tvdbID,
		Name:             r.Name,
		OriginalLanguage: r.OriginalLanguage,
		Overview:         r.Overview,
		VoteAverage:      float64(r.Score) / 10.0, // TVDb score is 0-100
	}

	if r.Status != nil {
		result.Status = r.Status.Name
	}

	if r.FirstAired != "" {
		result.FirstAirDate = parseDate(r.FirstAired)
	}
	if r.LastAired != "" {
		result.LastAirDate = parseDate(r.LastAired)
	}

	if r.Status != nil && r.Status.KeepUpdated {
		result.InProduction = true
	}

	result.NumberOfSeasons = countSeasons(r.Seasons)
	result.NumberOfEpisodes = 0 // Would need to fetch episodes to count

	if r.AverageRuntime > 0 {
		result.EpisodeRuntime = []int{r.AverageRuntime}
	}

	// Find best image from artworks
	for _, art := range r.Artworks {
		switch art.Type {
		case ArtworkTypePoster:
			if result.PosterPath == nil {
				result.PosterPath = &art.Image
			}
		case ArtworkTypeBackground:
			if result.BackdropPath == nil {
				result.BackdropPath = &art.Image
			}
		}
	}

	// Fallback to main image
	if result.PosterPath == nil && r.Image != nil {
		result.PosterPath = r.Image
	}

	// Map genres
	result.Genres = make([]metadata.Genre, len(r.Genres))
	for i, g := range r.Genres {
		result.Genres[i] = metadata.Genre{ID: g.ID, Name: g.Name}
	}

	// Map networks
	result.Networks = make([]metadata.Network, len(r.Networks))
	for i, n := range r.Networks {
		result.Networks[i] = metadata.Network{
			ID:            n.ID,
			Name:          n.Name,
			OriginCountry: ptrString(n.Country),
		}
	}

	// Map seasons
	result.Seasons = make([]metadata.SeasonSummary, 0, len(r.Seasons))
	for _, s := range r.Seasons {
		// Only include default/official seasons
		if s.Type != nil && s.Type.Type != "default" && s.Type.Type != "official" {
			continue
		}
		season := metadata.SeasonSummary{
			ProviderID:   strconv.Itoa(s.ID),
			SeasonNumber: s.Number,
			Name:         s.Name,
			PosterPath:   s.Image,
		}
		if s.Year != "" {
			if t := parseDate(s.Year + "-01-01"); t != nil {
				season.AirDate = t
			}
		}
		result.Seasons = append(result.Seasons, season)
	}

	// Map remote IDs
	for _, remote := range r.RemoteIDs {
		switch remote.Type {
		case RemoteIDTypeIMDb:
			result.IMDbID = &remote.ID
		case RemoteIDTypeTMDb:
			if id, err := strconv.Atoi(remote.ID); err == nil {
				tmdbID := int32(id)
				result.TMDbID = &tmdbID
			}
		}
	}

	// Map trailers
	for _, t := range r.Trailers {
		if t.URL != "" {
			result.TrailerURL = &t.URL
			break
		}
	}

	// Map origin countries
	if r.OriginalCountry != "" {
		result.OriginCountries = []string{r.OriginalCountry}
	}

	return result
}

// mapSeasonMetadata converts a TVDb season response to metadata type.
func mapSeasonMetadata(r *SeasonResponse, showID string) *metadata.SeasonMetadata {
	tvdbID := int32(r.ID)
	result := &metadata.SeasonMetadata{
		ProviderID:   strconv.Itoa(r.ID),
		Provider:     metadata.ProviderTVDb,
		TVDbID:       &tvdbID,
		ShowID:       showID,
		SeasonNumber: r.Number,
		Name:         r.Name,
		Overview:     r.Overview,
		PosterPath:   r.Image,
	}

	if r.Year != "" {
		result.AirDate = parseDate(r.Year + "-01-01")
	}

	// Map localized data
	if len(r.Overviews) > 0 {
		result.Translations = make(map[string]*metadata.LocalizedSeasonData)
		for lang, overview := range r.Overviews {
			result.Translations[lang] = &metadata.LocalizedSeasonData{
				Language: lang,
				Overview: overview,
			}
		}
	}

	return result
}

// mapEpisodeMetadata converts a TVDb episode response to metadata type.
func mapEpisodeMetadata(r *EpisodeResponse, showID string) *metadata.EpisodeMetadata {
	tvdbID := int32(r.ID)
	result := &metadata.EpisodeMetadata{
		ProviderID:     strconv.Itoa(r.ID),
		Provider:       metadata.ProviderTVDb,
		TVDbID:         &tvdbID,
		ShowID:         showID,
		SeasonNumber:   r.SeasonNumber,
		EpisodeNumber:  r.Number,
		Name:           r.Name,
		Overview:       r.Overview,
		StillPath:      r.Image,
		ProductionCode: r.ProductionCode,
	}

	if r.Aired != "" {
		result.AirDate = parseDate(r.Aired)
	}

	if r.Runtime != nil && *r.Runtime > 0 {
		runtime := int32(*r.Runtime)
		result.Runtime = &runtime
	}

	// Map characters to guest stars and crew
	for _, c := range r.Characters {
		if c.Type == CharacterTypeActor {
			result.GuestStars = append(result.GuestStars, metadata.CastMember{
				ProviderID:  strconv.Itoa(c.ID),
				Name:        c.PersonName,
				Character:   c.Name,
				Order:       c.Sort,
				ProfilePath: c.PersonImgURL,
			})
		} else {
			job := characterTypeToJob(c.Type)
			result.Crew = append(result.Crew, metadata.CrewMember{
				ProviderID:  strconv.Itoa(c.ID),
				Name:        c.PersonName,
				Job:         job,
				Department:  characterTypeToDepartment(c.Type),
				ProfilePath: c.PersonImgURL,
			})
		}
	}

	// Map localized data
	if len(r.Overviews) > 0 {
		result.Translations = make(map[string]*metadata.LocalizedEpisodeData)
		for lang, overview := range r.Overviews {
			result.Translations[lang] = &metadata.LocalizedEpisodeData{
				Language: lang,
				Overview: overview,
			}
		}
	}

	return result
}

// mapPersonSearchResult converts a TVDb search result to metadata type.
func mapPersonSearchResult(r *SearchResult) metadata.PersonSearchResult {
	return metadata.PersonSearchResult{
		ProviderID:  r.TVDbID,
		Provider:    metadata.ProviderTVDb,
		Name:        r.Name,
		ProfilePath: r.ImageURL,
	}
}

// mapPersonMetadata converts a TVDb person response to metadata type.
func mapPersonMetadata(r *PersonResponse, lang string) *metadata.PersonMetadata {
	result := &metadata.PersonMetadata{
		ProviderID:  strconv.Itoa(r.ID),
		Provider:    metadata.ProviderTVDb,
		Name:        r.Name,
		ProfilePath: r.Image,
		Gender:      r.Gender,
		Popularity:  float64(r.Score),
	}

	// Initialize translations map
	result.Translations = make(map[string]*metadata.LocalizedPersonData)

	if r.Birth != nil && *r.Birth != "" {
		result.Birthday = parseDate(*r.Birth)
	}
	if r.Death != nil && *r.Death != "" {
		result.Deathday = parseDate(*r.Death)
	}
	if r.BirthPlace != nil {
		result.PlaceOfBirth = r.BirthPlace
	}

	// Map aliases
	for _, a := range r.Aliases {
		result.AlsoKnownAs = append(result.AlsoKnownAs, a.Name)
	}

	// Map biographies
	for _, b := range r.Biographies {
		bio := b.Biography
		result.Translations[b.Language] = &metadata.LocalizedPersonData{
			Language:  b.Language,
			Biography: bio,
		}
		if b.Language == lang || (result.Biography == nil && b.Language == "eng") {
			result.Biography = &bio
		}
	}

	// Map remote IDs
	for _, remote := range r.RemoteIDs {
		switch remote.Type {
		case RemoteIDTypeIMDb:
			result.IMDbID = &remote.ID
		case RemoteIDTypeTMDb:
			if id, err := strconv.Atoi(remote.ID); err == nil {
				tmdbID := int32(id)
				result.TMDbID = &tmdbID
			}
		}
	}

	return result
}

// mapPersonCredits converts a TVDb person response to person credits.
func mapPersonCredits(r *PersonResponse) *metadata.PersonCredits {
	result := &metadata.PersonCredits{
		ProviderID: strconv.Itoa(r.ID),
		Provider:   metadata.ProviderTVDb,
	}

	for _, c := range r.Characters {
		mediaType := "tv"
		if c.MovieID != nil {
			mediaType = "movie"
		}

		mediaID := ""
		if c.SeriesID != nil {
			mediaID = strconv.Itoa(*c.SeriesID)
		} else if c.MovieID != nil {
			mediaID = strconv.Itoa(*c.MovieID)
		}

		if c.Type == CharacterTypeActor {
			result.CastCredits = append(result.CastCredits, metadata.MediaCredit{
				MediaType:  mediaType,
				MediaID:    mediaID,
				Title:      c.Name,
				Character:  &c.Name,
				PosterPath: c.Image,
			})
		} else {
			job := characterTypeToJob(c.Type)
			result.CrewCredits = append(result.CrewCredits, metadata.MediaCredit{
				MediaType:  mediaType,
				MediaID:    mediaID,
				Title:      c.Name,
				Job:        &job,
				PosterPath: c.Image,
			})
		}
	}

	return result
}

// mapCharactersToCredits converts TVDb characters to metadata credits.
func mapCharactersToCredits(characters []CharacterResponse) *metadata.Credits {
	result := &metadata.Credits{}

	for _, c := range characters {
		if c.Type == CharacterTypeActor {
			result.Cast = append(result.Cast, metadata.CastMember{
				ProviderID:  strconv.Itoa(c.ID),
				Name:        c.PersonName,
				Character:   c.Name,
				Order:       c.Sort,
				ProfilePath: c.PersonImgURL,
			})
		} else {
			result.Crew = append(result.Crew, metadata.CrewMember{
				ProviderID:  strconv.Itoa(c.ID),
				Name:        c.PersonName,
				Job:         characterTypeToJob(c.Type),
				Department:  characterTypeToDepartment(c.Type),
				ProfilePath: c.PersonImgURL,
			})
		}
	}

	return result
}

// mapArtworksToImages converts TVDb artworks to metadata images.
func mapArtworksToImages(artworks []ArtworkResponse) *metadata.Images {
	result := &metadata.Images{}

	for _, art := range artworks {
		img := metadata.Image{
			FilePath:    art.Image,
			Width:       art.Width,
			Height:      art.Height,
			VoteAverage: float64(art.Score),
			Language:    art.Language,
		}
		if art.Width > 0 && art.Height > 0 {
			img.AspectRatio = float64(art.Width) / float64(art.Height)
		}

		switch art.Type {
		case ArtworkTypePoster:
			result.Posters = append(result.Posters, img)
		case ArtworkTypeBackground:
			result.Backdrops = append(result.Backdrops, img)
		case ArtworkTypeBanner:
			// Banner goes to logos for now
			result.Logos = append(result.Logos, img)
		case ArtworkTypeClearLogo, ArtworkTypeClearArt:
			result.Logos = append(result.Logos, img)
		}
	}

	return result
}

// mapContentRatings converts TVDb content ratings to metadata type.
func mapContentRatings(ratings []ContentRatingResponse) []metadata.ContentRating {
	results := make([]metadata.ContentRating, len(ratings))
	for i, r := range ratings {
		results[i] = metadata.ContentRating{
			CountryCode: r.Country,
			Rating:      r.Name,
		}
	}
	return results
}

// mapOverviewsToTranslations converts TVDb overviews to metadata translations.
func mapOverviewsToTranslations(overviews map[string]string, nameTranslations []string) []metadata.Translation {
	results := make([]metadata.Translation, 0, len(overviews))

	for lang, overview := range overviews {
		results = append(results, metadata.Translation{
			Language: lang,
			Data: &metadata.TranslationData{
				Overview: overview,
			},
		})
	}

	return results
}

// mapRemoteIDsToExternalIDs converts TVDb remote IDs to metadata external IDs.
func mapRemoteIDsToExternalIDs(remoteIDs []RemoteIDResponse, tvdbID int32) *metadata.ExternalIDs {
	result := &metadata.ExternalIDs{
		TVDbID: &tvdbID,
	}

	for _, remote := range remoteIDs {
		switch remote.Type {
		case RemoteIDTypeIMDb:
			result.IMDbID = &remote.ID
		case RemoteIDTypeTMDb:
			if id, err := strconv.Atoi(remote.ID); err == nil {
				tmdbID := int32(id)
				result.TMDbID = &tmdbID
			}
		}
	}

	return result
}

// parseDate parses a date string.
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	// Try common formats
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t
		}
	}

	return nil
}

// countSeasons counts the number of default/official seasons.
func countSeasons(seasons []SeasonSummaryResponse) int {
	count := 0
	for _, s := range seasons {
		if s.Type == nil || s.Type.Type == "default" || s.Type.Type == "official" {
			count++
		}
	}
	return count
}

// characterTypeToJob converts TVDb character type to job name.
func characterTypeToJob(t int) string {
	switch t {
	case CharacterTypeDirector:
		return "Director"
	case CharacterTypeWriter:
		return "Writer"
	case CharacterTypeProducer:
		return "Producer"
	default:
		return "Unknown"
	}
}

// characterTypeToDepartment converts TVDb character type to department.
func characterTypeToDepartment(t int) string {
	switch t {
	case CharacterTypeDirector:
		return "Directing"
	case CharacterTypeWriter:
		return "Writing"
	case CharacterTypeProducer:
		return "Production"
	default:
		return "Unknown"
	}
}

// ptrString returns the value of a pointer or empty string.
func ptrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
