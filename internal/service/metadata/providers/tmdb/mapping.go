package tmdb

import (
	"strconv"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapMovieSearchResult converts a TMDb search result to metadata type.
func mapMovieSearchResult(r *MovieSearchResponse) metadata.MovieSearchResult {
	result := metadata.MovieSearchResult{
		ProviderID:       strconv.Itoa(r.ID),
		Provider:         metadata.ProviderTMDb,
		Title:            r.Title,
		OriginalTitle:    r.OriginalTitle,
		OriginalLanguage: r.OriginalLanguage,
		Overview:         r.Overview,
		PosterPath:       r.PosterPath,
		BackdropPath:     r.BackdropPath,
		VoteAverage:      r.VoteAverage,
		VoteCount:        r.VoteCount,
		Popularity:       r.Popularity,
		Adult:            r.Adult,
		GenreIDs:         r.GenreIDs,
	}

	if r.ReleaseDate != "" {
		result.ReleaseDate = parseDate(r.ReleaseDate)
		result.Year = extractYear(r.ReleaseDate)
	}

	return result
}

// mapMovieMetadata converts a TMDb movie response to metadata type.
func mapMovieMetadata(r *MovieResponse) *metadata.MovieMetadata {
	tmdbID := int32(r.ID)
	result := &metadata.MovieMetadata{
		ProviderID:       strconv.Itoa(r.ID),
		Provider:         metadata.ProviderTMDb,
		TMDbID:           &tmdbID,
		IMDbID:           r.IMDbID,
		Title:            r.Title,
		OriginalTitle:    r.OriginalTitle,
		OriginalLanguage: r.OriginalLanguage,
		Tagline:          r.Tagline,
		Overview:         r.Overview,
		Status:           r.Status,
		VoteAverage:      r.VoteAverage,
		VoteCount:        r.VoteCount,
		Popularity:       r.Popularity,
		Adult:            r.Adult,
		PosterPath:       r.PosterPath,
		BackdropPath:     r.BackdropPath,
		Homepage:         r.Homepage,
	}

	if r.ReleaseDate != "" {
		result.ReleaseDate = parseDate(r.ReleaseDate)
	}

	if r.Runtime != nil && *r.Runtime > 0 {
		runtime := int32(*r.Runtime)
		result.Runtime = &runtime
	}

	if r.Budget != nil && *r.Budget > 0 {
		result.Budget = r.Budget
	}

	if r.Revenue != nil && *r.Revenue > 0 {
		result.Revenue = r.Revenue
	}

	// Map genres
	result.Genres = make([]metadata.Genre, len(r.Genres))
	for i, g := range r.Genres {
		result.Genres[i] = metadata.Genre{ID: g.ID, Name: g.Name}
	}

	// Map production companies
	result.ProductionCompanies = make([]metadata.ProductionCompany, len(r.ProductionCompanies))
	for i, c := range r.ProductionCompanies {
		result.ProductionCompanies[i] = metadata.ProductionCompany{
			ID:            c.ID,
			Name:          c.Name,
			LogoPath:      c.LogoPath,
			OriginCountry: c.OriginCountry,
		}
	}

	// Map production countries
	result.ProductionCountries = make([]metadata.ProductionCountry, len(r.ProductionCountries))
	for i, c := range r.ProductionCountries {
		result.ProductionCountries[i] = metadata.ProductionCountry{
			ISOCode: c.ISO3166_1,
			Name:    c.Name,
		}
	}

	// Map spoken languages
	result.SpokenLanguages = make([]metadata.SpokenLanguage, len(r.SpokenLanguages))
	for i, l := range r.SpokenLanguages {
		result.SpokenLanguages[i] = metadata.SpokenLanguage{
			ISOCode:     l.ISO639_1,
			Name:        l.Name,
			EnglishName: l.EnglishName,
		}
	}

	// Map collection
	if r.BelongsToCollection != nil {
		result.Collection = &metadata.CollectionRef{
			ID:           r.BelongsToCollection.ID,
			Name:         r.BelongsToCollection.Name,
			PosterPath:   r.BelongsToCollection.PosterPath,
			BackdropPath: r.BelongsToCollection.BackdropPath,
		}
	}

	// Map videos (find trailer)
	if r.Videos != nil {
		for _, v := range r.Videos.Results {
			if v.Site == "YouTube" && (v.Type == "Trailer" || v.Type == "Teaser") {
				url := "https://www.youtube.com/watch?v=" + v.Key
				result.TrailerURL = &url
				break
			}
		}
	}

	return result
}

// mapTVSearchResult converts a TMDb TV search result to metadata type.
func mapTVSearchResult(r *TVSearchResponse) metadata.TVShowSearchResult {
	result := metadata.TVShowSearchResult{
		ProviderID:       strconv.Itoa(r.ID),
		Provider:         metadata.ProviderTMDb,
		Name:             r.Name,
		OriginalName:     r.OriginalName,
		OriginalLanguage: r.OriginalLanguage,
		Overview:         r.Overview,
		PosterPath:       r.PosterPath,
		BackdropPath:     r.BackdropPath,
		VoteAverage:      r.VoteAverage,
		VoteCount:        r.VoteCount,
		Popularity:       r.Popularity,
		Adult:            r.Adult,
		GenreIDs:         r.GenreIDs,
		OriginCountries:  r.OriginCountry,
	}

	if r.FirstAirDate != "" {
		result.FirstAirDate = parseDate(r.FirstAirDate)
		result.Year = extractYear(r.FirstAirDate)
	}

	return result
}

// mapTVShowMetadata converts a TMDb TV response to metadata type.
func mapTVShowMetadata(r *TVResponse) *metadata.TVShowMetadata {
	tmdbID := int32(r.ID)
	result := &metadata.TVShowMetadata{
		ProviderID:       strconv.Itoa(r.ID),
		Provider:         metadata.ProviderTMDb,
		TMDbID:           &tmdbID,
		Name:             r.Name,
		OriginalName:     r.OriginalName,
		OriginalLanguage: r.OriginalLanguage,
		Tagline:          r.Tagline,
		Overview:         r.Overview,
		Status:           r.Status,
		Type:             r.Type,
		InProduction:     r.InProduction,
		NumberOfSeasons:  r.NumberOfSeasons,
		NumberOfEpisodes: r.NumberOfEpisodes,
		EpisodeRuntime:   r.EpisodeRunTime,
		VoteAverage:      r.VoteAverage,
		VoteCount:        r.VoteCount,
		Popularity:       r.Popularity,
		Adult:            r.Adult,
		PosterPath:       r.PosterPath,
		BackdropPath:     r.BackdropPath,
		Homepage:         r.Homepage,
		OriginCountries:  r.OriginCountry,
	}

	if r.FirstAirDate != "" {
		result.FirstAirDate = parseDate(r.FirstAirDate)
	}
	if r.LastAirDate != "" {
		result.LastAirDate = parseDate(r.LastAirDate)
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
			LogoPath:      n.LogoPath,
			OriginCountry: n.OriginCountry,
		}
	}

	// Map creators
	result.CreatedBy = make([]metadata.Creator, len(r.CreatedBy))
	for i, c := range r.CreatedBy {
		result.CreatedBy[i] = metadata.Creator{
			ID:          c.ID,
			Name:        c.Name,
			Gender:      c.Gender,
			ProfilePath: c.ProfilePath,
			CreditID:    &c.CreditID,
		}
	}

	// Map spoken languages
	result.SpokenLanguages = make([]metadata.SpokenLanguage, len(r.SpokenLanguages))
	for i, l := range r.SpokenLanguages {
		result.SpokenLanguages[i] = metadata.SpokenLanguage{
			ISOCode:     l.ISO639_1,
			Name:        l.Name,
			EnglishName: l.EnglishName,
		}
	}

	// Map seasons
	result.Seasons = make([]metadata.SeasonSummary, len(r.Seasons))
	for i, s := range r.Seasons {
		result.Seasons[i] = metadata.SeasonSummary{
			ProviderID:   strconv.Itoa(s.ID),
			SeasonNumber: s.SeasonNumber,
			Name:         s.Name,
			Overview:     s.Overview,
			PosterPath:   s.PosterPath,
			EpisodeCount: s.EpisodeCount,
			VoteAverage:  s.VoteAverage,
		}
		if s.AirDate != "" {
			result.Seasons[i].AirDate = parseDate(s.AirDate)
		}
	}

	// Map videos (find trailer)
	if r.Videos != nil {
		for _, v := range r.Videos.Results {
			if v.Site == "YouTube" && (v.Type == "Trailer" || v.Type == "Teaser") {
				url := "https://www.youtube.com/watch?v=" + v.Key
				result.TrailerURL = &url
				break
			}
		}
	}

	return result
}

// mapSeasonMetadata converts a TMDb season response to metadata type.
func mapSeasonMetadata(r *SeasonResponse, showID string) *metadata.SeasonMetadata {
	tmdbID := int32(r.ID)
	result := &metadata.SeasonMetadata{
		ProviderID:   strconv.Itoa(r.ID),
		Provider:     metadata.ProviderTMDb,
		TMDbID:       &tmdbID,
		ShowID:       showID,
		SeasonNumber: r.SeasonNumber,
		Name:         r.Name,
		Overview:     r.Overview,
		PosterPath:   r.PosterPath,
		VoteAverage:  r.VoteAverage,
	}

	if r.AirDate != "" {
		result.AirDate = parseDate(r.AirDate)
	}

	// Map episodes
	result.Episodes = make([]metadata.EpisodeSummary, len(r.Episodes))
	for i, e := range r.Episodes {
		result.Episodes[i] = metadata.EpisodeSummary{
			ProviderID:     strconv.Itoa(e.ID),
			EpisodeNumber:  e.EpisodeNumber,
			Name:           e.Name,
			Overview:       e.Overview,
			StillPath:      e.StillPath,
			VoteAverage:    e.VoteAverage,
			VoteCount:      e.VoteCount,
			ProductionCode: e.ProductionCode,
		}
		if e.AirDate != "" {
			result.Episodes[i].AirDate = parseDate(e.AirDate)
		}
		if e.Runtime != nil && *e.Runtime > 0 {
			runtime := int32(*e.Runtime)
			result.Episodes[i].Runtime = &runtime
		}
	}

	return result
}

// mapEpisodeMetadata converts a TMDb episode response to metadata type.
func mapEpisodeMetadata(r *EpisodeResponse, showID string) *metadata.EpisodeMetadata {
	tmdbID := int32(r.ID)
	result := &metadata.EpisodeMetadata{
		ProviderID:     strconv.Itoa(r.ID),
		Provider:       metadata.ProviderTMDb,
		TMDbID:         &tmdbID,
		ShowID:         showID,
		SeasonNumber:   r.SeasonNumber,
		EpisodeNumber:  r.EpisodeNumber,
		Name:           r.Name,
		Overview:       r.Overview,
		StillPath:      r.StillPath,
		VoteAverage:    r.VoteAverage,
		VoteCount:      r.VoteCount,
		ProductionCode: r.ProductionCode,
	}

	if r.AirDate != "" {
		result.AirDate = parseDate(r.AirDate)
	}
	if r.Runtime != nil && *r.Runtime > 0 {
		runtime := int32(*r.Runtime)
		result.Runtime = &runtime
	}

	// Map guest stars
	result.GuestStars = make([]metadata.CastMember, len(r.GuestStars))
	for i, g := range r.GuestStars {
		result.GuestStars[i] = metadata.CastMember{
			ProviderID:  strconv.Itoa(g.ID),
			Name:        g.Name,
			Character:   g.Character,
			Order:       g.Order,
			CreditID:    &g.CreditID,
			Gender:      g.Gender,
			ProfilePath: g.ProfilePath,
		}
	}

	// Map crew
	result.Crew = make([]metadata.CrewMember, len(r.Crew))
	for i, c := range r.Crew {
		result.Crew[i] = metadata.CrewMember{
			ProviderID:  strconv.Itoa(c.ID),
			Name:        c.Name,
			Job:         c.Job,
			Department:  c.Department,
			CreditID:    &c.CreditID,
			Gender:      c.Gender,
			ProfilePath: c.ProfilePath,
		}
	}

	return result
}

// mapPersonSearchResult converts a TMDb person search result to metadata type.
func mapPersonSearchResult(r *PersonSearchResponse) metadata.PersonSearchResult {
	result := metadata.PersonSearchResult{
		ProviderID:  strconv.Itoa(r.ID),
		Provider:    metadata.ProviderTMDb,
		Name:        r.Name,
		ProfilePath: r.ProfilePath,
		Popularity:  r.Popularity,
		Adult:       r.Adult,
	}

	// Map known for
	result.KnownFor = make([]metadata.MediaReference, len(r.KnownFor))
	for i, kf := range r.KnownFor {
		title := kf.Title
		if kf.MediaType == "tv" {
			title = kf.Name
		}
		result.KnownFor[i] = metadata.MediaReference{
			MediaType:  kf.MediaType,
			ID:         strconv.Itoa(kf.ID),
			Title:      title,
			PosterPath: kf.PosterPath,
		}
	}

	return result
}

// mapPersonMetadata converts a TMDb person response to metadata type.
func mapPersonMetadata(r *PersonResponse) *metadata.PersonMetadata {
	tmdbID := int32(r.ID)
	result := &metadata.PersonMetadata{
		ProviderID:   strconv.Itoa(r.ID),
		Provider:     metadata.ProviderTMDb,
		TMDbID:       &tmdbID,
		IMDbID:       r.IMDbID,
		Name:         r.Name,
		AlsoKnownAs:  r.AlsoKnownAs,
		Biography:    r.Biography,
		Gender:       r.Gender,
		PlaceOfBirth: r.PlaceOfBirth,
		ProfilePath:  r.ProfilePath,
		Homepage:     r.Homepage,
		Popularity:   r.Popularity,
		Adult:        r.Adult,
		KnownForDept: r.KnownForDept,
	}

	if r.Birthday != nil && *r.Birthday != "" {
		result.Birthday = parseDate(*r.Birthday)
	}
	if r.Deathday != nil && *r.Deathday != "" {
		result.Deathday = parseDate(*r.Deathday)
	}

	return result
}

// mapPersonCredits converts a TMDb person credits response to metadata type.
func mapPersonCredits(r *PersonCreditsResponse) *metadata.PersonCredits {
	result := &metadata.PersonCredits{
		ProviderID: strconv.Itoa(r.ID),
		Provider:   metadata.ProviderTMDb,
	}

	// Map cast credits
	result.CastCredits = make([]metadata.MediaCredit, len(r.Cast))
	for i, c := range r.Cast {
		title := c.Title
		if c.MediaType == "tv" {
			title = c.Name
		}
		result.CastCredits[i] = metadata.MediaCredit{
			MediaType:    c.MediaType,
			MediaID:      strconv.Itoa(c.ID),
			Title:        title,
			Character:    c.Character,
			PosterPath:   c.PosterPath,
			VoteAverage:  c.VoteAverage,
			EpisodeCount: c.EpisodeCount,
		}
		if c.ReleaseDate != "" {
			result.CastCredits[i].ReleaseDate = parseDate(c.ReleaseDate)
		} else if c.FirstAirDate != "" {
			result.CastCredits[i].ReleaseDate = parseDate(c.FirstAirDate)
		}
	}

	// Map crew credits
	result.CrewCredits = make([]metadata.MediaCredit, len(r.Crew))
	for i, c := range r.Crew {
		title := c.Title
		if c.MediaType == "tv" {
			title = c.Name
		}
		result.CrewCredits[i] = metadata.MediaCredit{
			MediaType:    c.MediaType,
			MediaID:      strconv.Itoa(c.ID),
			Title:        title,
			Job:          c.Job,
			Department:   c.Department,
			PosterPath:   c.PosterPath,
			VoteAverage:  c.VoteAverage,
			EpisodeCount: c.EpisodeCount,
		}
		if c.ReleaseDate != "" {
			result.CrewCredits[i].ReleaseDate = parseDate(c.ReleaseDate)
		} else if c.FirstAirDate != "" {
			result.CrewCredits[i].ReleaseDate = parseDate(c.FirstAirDate)
		}
	}

	return result
}

// mapCredits converts a TMDb credits response to metadata type.
func mapCredits(r *CreditsResponse) *metadata.Credits {
	result := &metadata.Credits{
		Cast: make([]metadata.CastMember, len(r.Cast)),
		Crew: make([]metadata.CrewMember, len(r.Crew)),
	}

	for i, c := range r.Cast {
		result.Cast[i] = metadata.CastMember{
			ProviderID:  strconv.Itoa(c.ID),
			Name:        c.Name,
			Character:   c.Character,
			Order:       c.Order,
			CreditID:    &c.CreditID,
			Gender:      c.Gender,
			ProfilePath: c.ProfilePath,
		}
	}

	for i, c := range r.Crew {
		result.Crew[i] = metadata.CrewMember{
			ProviderID:  strconv.Itoa(c.ID),
			Name:        c.Name,
			Job:         c.Job,
			Department:  c.Department,
			CreditID:    &c.CreditID,
			Gender:      c.Gender,
			ProfilePath: c.ProfilePath,
		}
	}

	return result
}

// mapImages converts a TMDb images response to metadata type.
func mapImages(r *ImagesResponse) *metadata.Images {
	result := &metadata.Images{
		Posters:   make([]metadata.Image, len(r.Posters)),
		Backdrops: make([]metadata.Image, len(r.Backdrops)),
		Logos:     make([]metadata.Image, len(r.Logos)),
		Stills:    make([]metadata.Image, len(r.Stills)),
	}

	for i, img := range r.Posters {
		result.Posters[i] = mapImage(&img)
	}
	for i, img := range r.Backdrops {
		result.Backdrops[i] = mapImage(&img)
	}
	for i, img := range r.Logos {
		result.Logos[i] = mapImage(&img)
	}
	for i, img := range r.Stills {
		result.Stills[i] = mapImage(&img)
	}

	return result
}

// mapPersonImages converts a TMDb person images response to metadata type.
func mapPersonImages(r *PersonImagesResponse) *metadata.Images {
	result := &metadata.Images{
		Profiles: make([]metadata.Image, len(r.Profiles)),
	}

	for i, img := range r.Profiles {
		result.Profiles[i] = mapImage(&img)
	}

	return result
}

// mapImage converts a single TMDb image to metadata type.
func mapImage(r *ImageResponse) metadata.Image {
	return metadata.Image{
		FilePath:    r.FilePath,
		AspectRatio: r.AspectRatio,
		Width:       r.Width,
		Height:      r.Height,
		VoteAverage: r.VoteAverage,
		VoteCount:   r.VoteCount,
		Language:    r.ISO639_1,
	}
}

// mapReleaseDates converts a TMDb release dates response to metadata type.
func mapReleaseDates(r *ReleaseDatesWrapper) []metadata.ReleaseDate {
	var results []metadata.ReleaseDate

	for _, country := range r.Results {
		for _, rd := range country.ReleaseDates {
			result := metadata.ReleaseDate{
				CountryCode:   country.ISO3166_1,
				Certification: rd.Certification,
				ReleaseType:   rd.Type,
				Language:      rd.ISO639_1,
				Note:          rd.Note,
			}
			if rd.ReleaseDate != "" {
				result.ReleaseDate = parseDate(rd.ReleaseDate)
			}
			results = append(results, result)
		}
	}

	return results
}

// mapContentRatings converts a TMDb content ratings response to metadata type.
func mapContentRatings(r *ContentRatingsWrapper) []metadata.ContentRating {
	results := make([]metadata.ContentRating, len(r.Results))

	for i, cr := range r.Results {
		results[i] = metadata.ContentRating{
			CountryCode: cr.ISO3166_1,
			Rating:      cr.Rating,
			Descriptors: cr.Descriptors,
		}
	}

	return results
}

// mapTranslations converts a TMDb translations response to metadata type.
func mapTranslations(r *TranslationsWrapper) []metadata.Translation {
	results := make([]metadata.Translation, len(r.Translations))

	for i, t := range r.Translations {
		result := metadata.Translation{
			ISOCode:     t.ISO3166_1,
			Language:    t.ISO639_1,
			Name:        t.Name,
			EnglishName: t.EnglishName,
		}

		// Map translation data
		title := t.Data.Title
		if title == "" {
			title = t.Data.Name
		}
		result.Data = &metadata.TranslationData{
			Title:    title,
			Overview: t.Data.Overview,
			Tagline:  t.Data.Tagline,
			Homepage: t.Data.Homepage,
		}
		if t.Data.Runtime != nil {
			runtime := int32(*t.Data.Runtime)
			result.Data.Runtime = &runtime
		}

		results[i] = result
	}

	return results
}

// mapExternalIDs converts a TMDb external IDs response to metadata type.
func mapExternalIDs(r *ExternalIDsResponse, tmdbID int32) *metadata.ExternalIDs {
	result := &metadata.ExternalIDs{
		TMDbID: &tmdbID,
		IMDbID: r.IMDbID,
	}

	if r.TVDbID != nil && *r.TVDbID > 0 {
		tvdbID := int32(*r.TVDbID)
		result.TVDbID = &tvdbID
	}
	if r.TVRageID != nil && *r.TVRageID > 0 {
		rageID := int32(*r.TVRageID)
		result.TVRageID = &rageID
	}

	result.WikidataID = r.WikidataID
	result.FacebookID = r.FacebookID
	result.InstagramID = r.InstagramID
	result.TwitterID = r.TwitterID
	result.TikTokID = r.TikTokID
	result.YouTubeID = r.YouTubeID
	result.FreebaseID = r.FreebaseID
	result.FreebaseMID = r.FreebaseMID

	return result
}

// mapCollectionMetadata converts a TMDb collection response to metadata type.
func mapCollectionMetadata(r *CollectionResponse) *metadata.CollectionMetadata {
	result := &metadata.CollectionMetadata{
		ProviderID:   strconv.Itoa(r.ID),
		Provider:     metadata.ProviderTMDb,
		Name:         r.Name,
		Overview:     r.Overview,
		PosterPath:   r.PosterPath,
		BackdropPath: r.BackdropPath,
	}

	// Map parts (movies in collection)
	result.Parts = make([]metadata.MovieSearchResult, len(r.Parts))
	for i, p := range r.Parts {
		result.Parts[i] = mapMovieSearchResult(&p)
	}

	return result
}

// parseDate parses a date string in ISO format (YYYY-MM-DD or with time).
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	// Try full ISO 8601 format first
	if strings.Contains(dateStr, "T") {
		t, err := time.Parse(time.RFC3339, dateStr)
		if err == nil {
			return &t
		}
	}

	// Try simple date format
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

// extractYear extracts the year from a date string.
func extractYear(dateStr string) *int {
	if dateStr == "" || len(dateStr) < 4 {
		return nil
	}

	yearStr := dateStr[:4]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil
	}
	return &year
}
