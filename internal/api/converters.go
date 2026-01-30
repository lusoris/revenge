package api

import (
	"context"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/qar/crew"
	"github.com/lusoris/revenge/internal/content/qar/expedition"
	"github.com/lusoris/revenge/internal/content/qar/flag"
	"github.com/lusoris/revenge/internal/content/qar/port"
	"github.com/lusoris/revenge/internal/content/qar/voyage"
	"github.com/lusoris/revenge/internal/content/shared"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// userToAPI converts a db.User to a gen.User.
func userToAPI(u *db.User) gen.User {
	result := gen.User{
		ID:         u.ID,
		Username:   u.Username,
		IsAdmin:    u.IsAdmin,
		IsDisabled: u.IsDisabled,
		CreatedAt:  u.CreatedAt,
	}

	if u.Email != nil {
		result.Email = gen.NewOptString(*u.Email)
	}

	result.MaxRatingLevel = gen.NewOptInt(int(u.MaxRatingLevel))
	result.AdultEnabled = gen.NewOptBool(u.AdultEnabled)

	if u.PreferredLanguage != nil {
		result.PreferredLanguage = gen.NewOptString(*u.PreferredLanguage)
	}

	if u.LastLoginAt.Valid {
		result.LastLoginAt = gen.NewOptDateTime(u.LastLoginAt.Time)
	}

	result.UpdatedAt = gen.NewOptDateTime(u.UpdatedAt)

	return result
}

// sessionToAPI converts a db.Session to a gen.Session.
func sessionToAPI(s *db.Session) gen.Session {
	result := gen.Session{
		ID:           s.ID,
		UserId:       s.UserID,
		IsActive:     s.IsActive,
		LastActivity: s.LastActivity,
		CreatedAt:    s.CreatedAt,
	}

	if s.ProfileID.Valid {
		result.ProfileId = gen.NewOptUUID(s.ProfileID.Bytes)
	}

	if s.DeviceName != nil {
		result.DeviceName = gen.NewOptString(*s.DeviceName)
	}

	if s.DeviceType != nil {
		result.DeviceType = gen.NewOptString(*s.DeviceType)
	}

	if s.ClientName != nil {
		result.ClientName = gen.NewOptString(*s.ClientName)
	}

	if s.ClientVersion != nil {
		result.ClientVersion = gen.NewOptString(*s.ClientVersion)
	}

	// Convert IP address to string
	if s.IpAddress.IsValid() {
		result.IpAddress = gen.NewOptString(s.IpAddress.String())
	}

	result.ExpiresAt = gen.NewOptDateTime(s.ExpiresAt)

	return result
}

// libraryToAPI converts a shared.LibraryInfo to a gen.Library.
func libraryToAPI(l *shared.LibraryInfo) gen.Library {
	result := gen.Library{
		ID:    l.ID,
		Name:  l.Name,
		Type:  gen.LibraryType(l.Module), // Module name as type
		Paths: l.Paths,
	}

	// Adult content libraries are marked as private
	if l.IsAdult {
		result.IsPrivate = gen.NewOptBool(true)
	}

	// Module-specific settings are in l.Settings but we can't easily
	// map them to the generic Library API response.
	// TODO: Extract common settings from l.Settings if needed

	return result
}

// movieToAPI converts a movie.Movie to a gen.Movie.
func movieToAPI(m *movie.Movie, userData *gen.MovieUserData) gen.Movie {
	result := gen.Movie{
		ID:        m.ID,
		LibraryId: m.MovieLibraryID,
		Title:     m.Title,
	}

	if m.CollectionID != nil {
		result.CollectionId = gen.NewOptUUID(*m.CollectionID)
	}
	if m.OriginalTitle != "" {
		result.OriginalTitle = gen.NewOptString(m.OriginalTitle)
	}
	if m.SortTitle != "" {
		result.SortTitle = gen.NewOptString(m.SortTitle)
	}
	if m.Overview != "" {
		result.Overview = gen.NewOptString(m.Overview)
	}
	if m.Tagline != "" {
		result.Tagline = gen.NewOptString(m.Tagline)
	}
	if m.Year > 0 {
		result.Year = gen.NewOptInt(m.Year)
	}
	if m.ReleaseDate != nil {
		result.ReleaseDate = gen.NewOptDate(*m.ReleaseDate)
	}
	if m.RuntimeTicks > 0 {
		result.RuntimeTicks = gen.NewOptInt64(m.RuntimeTicks)
	}
	if m.CommunityRating > 0 {
		result.CommunityRating = gen.NewOptFloat64(m.CommunityRating)
	}
	if m.VoteCount > 0 {
		result.VoteCount = gen.NewOptInt(m.VoteCount)
	}
	if m.ContentRating != "" {
		result.ContentRating = gen.NewOptString(m.ContentRating)
	}
	if m.PosterPath != "" {
		result.PosterPath = gen.NewOptString(m.PosterPath)
	}
	if m.BackdropPath != "" {
		result.BackdropPath = gen.NewOptString(m.BackdropPath)
	}
	if m.TmdbID > 0 {
		result.TmdbId = gen.NewOptInt(m.TmdbID)
	}
	if m.ImdbID != "" {
		result.ImdbId = gen.NewOptString(m.ImdbID)
	}
	result.DateAdded = gen.NewOptDateTime(m.CreatedAt)

	if userData != nil {
		result.UserData = gen.NewOptMovieUserData(*userData)
	}

	return result
}

// movieFullToAPI converts a movie.Movie with relations to a gen.MovieFull.
func movieFullToAPI(m *movie.Movie, userData *gen.MovieUserData) gen.MovieFull {
	base := movieToAPI(m, userData)
	result := gen.MovieFull{
		ID:              base.ID,
		LibraryId:       base.LibraryId,
		CollectionId:    base.CollectionId,
		Title:           base.Title,
		OriginalTitle:   base.OriginalTitle,
		SortTitle:       base.SortTitle,
		Overview:        base.Overview,
		Tagline:         base.Tagline,
		Year:            base.Year,
		ReleaseDate:     base.ReleaseDate,
		RuntimeTicks:    base.RuntimeTicks,
		CommunityRating: base.CommunityRating,
		VoteCount:       base.VoteCount,
		ContentRating:   base.ContentRating,
		PosterPath:      base.PosterPath,
		BackdropPath:    base.BackdropPath,
		TmdbId:          base.TmdbId,
		ImdbId:          base.ImdbId,
		DateAdded:       base.DateAdded,
		UserData:        base.UserData,
	}

	if m.Budget > 0 {
		result.Budget = gen.NewOptInt64(m.Budget)
	}
	if m.Revenue > 0 {
		result.Revenue = gen.NewOptInt64(m.Revenue)
	}

	// Convert genres
	if len(m.Genres) > 0 {
		genres := make([]gen.Genre, len(m.Genres))
		for i, g := range m.Genres {
			genres[i] = gen.Genre{
				ID:   g.ID,
				Name: g.Name,
			}
		}
		result.Genres = genres
	}

	// Convert cast
	if len(m.Cast) > 0 {
		cast := make([]gen.CastMember, len(m.Cast))
		for i, c := range m.Cast {
			cast[i] = gen.CastMember{
				PersonId:  c.PersonID,
				Name:      c.Name,
				Character: c.CharacterName,
			}
			if c.BillingOrder > 0 {
				cast[i].Order = gen.NewOptInt(c.BillingOrder)
			}
			if c.PrimaryImageURL != "" {
				cast[i].ProfilePath = gen.NewOptString(c.PrimaryImageURL)
			}
		}
		result.Cast = cast
	}

	// Convert directors
	if len(m.Directors) > 0 {
		directors := make([]gen.CrewMember, len(m.Directors))
		for i, d := range m.Directors {
			directors[i] = crewMemberToAPI(&d)
		}
		result.Directors = directors
	}

	// Convert writers
	if len(m.Writers) > 0 {
		writers := make([]gen.CrewMember, len(m.Writers))
		for i, w := range m.Writers {
			writers[i] = crewMemberToAPI(&w)
		}
		result.Writers = writers
	}

	return result
}

// movieWithProgressToAPI converts a movie with watch history to API type.
func movieWithProgressToAPI(m *movie.Movie, history *movie.WatchHistory) gen.MovieWithProgress {
	base := movieToAPI(m, nil)
	result := gen.MovieWithProgress{
		ID:              base.ID,
		LibraryId:       base.LibraryId,
		CollectionId:    base.CollectionId,
		Title:           base.Title,
		OriginalTitle:   base.OriginalTitle,
		SortTitle:       base.SortTitle,
		Overview:        base.Overview,
		Tagline:         base.Tagline,
		Year:            base.Year,
		ReleaseDate:     base.ReleaseDate,
		RuntimeTicks:    base.RuntimeTicks,
		CommunityRating: base.CommunityRating,
		VoteCount:       base.VoteCount,
		ContentRating:   base.ContentRating,
		PosterPath:      base.PosterPath,
		BackdropPath:    base.BackdropPath,
		TmdbId:          base.TmdbId,
		ImdbId:          base.ImdbId,
		DateAdded:       base.DateAdded,
		UserData:        base.UserData,
	}

	if history != nil {
		result.Progress = gen.NewOptWatchProgress(gen.WatchProgress{
			PositionTicks:    history.PositionTicks,
			DurationTicks:    history.DurationTicks,
			PlayedPercentage: gen.NewOptFloat64(history.PlayedPercentage),
			LastPlayedAt:     gen.NewOptDateTime(history.LastUpdatedAt),
		})
	}

	return result
}

// collectionToAPI converts a movie.Collection to a gen.Collection.
func collectionToAPI(c *movie.Collection) gen.Collection {
	result := gen.Collection{
		ID:   c.ID,
		Name: c.Name,
	}

	if c.Overview != "" {
		result.Overview = gen.NewOptString(c.Overview)
	}
	if c.PosterPath != "" {
		result.PosterPath = gen.NewOptString(c.PosterPath)
	}
	if c.BackdropPath != "" {
		result.BackdropPath = gen.NewOptString(c.BackdropPath)
	}
	if c.TmdbID > 0 {
		result.TmdbId = gen.NewOptInt(c.TmdbID)
	}

	return result
}

// collectionFullToAPI converts a collection with movies to API type.
func collectionFullToAPI(c *movie.Collection, movies []*movie.Movie) gen.CollectionFull {
	base := collectionToAPI(c)
	result := gen.CollectionFull{
		ID:           base.ID,
		Name:         base.Name,
		Overview:     base.Overview,
		PosterPath:   base.PosterPath,
		BackdropPath: base.BackdropPath,
		TmdbId:       base.TmdbId,
	}

	if len(movies) > 0 {
		movieList := make([]gen.Movie, len(movies))
		for i, m := range movies {
			movieList[i] = movieToAPI(m, nil)
		}
		result.Movies = movieList
	}

	return result
}

// crewMemberToAPI converts a movie.CrewMember to a gen.CrewMember.
func crewMemberToAPI(c *movie.CrewMember) gen.CrewMember {
	result := gen.CrewMember{
		PersonId: c.PersonID,
		Name:     c.Name,
		Job:      c.Job,
	}

	if c.Department != "" {
		result.Department = gen.NewOptString(c.Department)
	}
	if c.PrimaryImageURL != "" {
		result.ProfilePath = gen.NewOptString(c.PrimaryImageURL)
	}

	return result
}

// =============================================================================
// TV Show Converters
// =============================================================================

// seriesToAPI converts a tvshow.Series to a gen.Series.
func seriesToAPI(s *tvshow.Series, userData *gen.SeriesUserData) gen.Series {
	result := gen.Series{
		ID:        s.ID,
		LibraryId: s.TvLibraryID,
		Title:     s.Title,
	}

	if s.OriginalTitle != "" {
		result.OriginalTitle = gen.NewOptString(s.OriginalTitle)
	}
	if s.SortTitle != "" {
		result.SortTitle = gen.NewOptString(s.SortTitle)
	}
	if s.Overview != "" {
		result.Overview = gen.NewOptString(s.Overview)
	}
	if s.Tagline != "" {
		result.Tagline = gen.NewOptString(s.Tagline)
	}
	if s.FirstAirDate != nil {
		result.FirstAirDate = gen.NewOptDate(*s.FirstAirDate)
	}
	if s.LastAirDate != nil {
		result.LastAirDate = gen.NewOptDate(*s.LastAirDate)
	}
	if s.Status != "" {
		result.Status = gen.NewOptSeriesStatus(gen.SeriesStatus(s.Status))
	}
	if s.SeasonCount > 0 {
		result.SeasonCount = gen.NewOptInt(s.SeasonCount)
	}
	if s.EpisodeCount > 0 {
		result.EpisodeCount = gen.NewOptInt(s.EpisodeCount)
	}
	if s.CommunityRating > 0 {
		result.CommunityRating = gen.NewOptFloat32(float32(s.CommunityRating))
	}
	if s.VoteCount > 0 {
		result.VoteCount = gen.NewOptInt(s.VoteCount)
	}
	if s.ContentRating != "" {
		result.ContentRating = gen.NewOptString(s.ContentRating)
	}
	if s.PosterPath != "" {
		result.PosterPath = gen.NewOptString(s.PosterPath)
	}
	if s.PosterBlurhash != "" {
		result.PosterBlurhash = gen.NewOptString(s.PosterBlurhash)
	}
	if s.BackdropPath != "" {
		result.BackdropPath = gen.NewOptString(s.BackdropPath)
	}
	if s.BackdropBlurhash != "" {
		result.BackdropBlurhash = gen.NewOptString(s.BackdropBlurhash)
	}
	if s.LogoPath != "" {
		result.LogoPath = gen.NewOptString(s.LogoPath)
	}
	if s.TmdbID > 0 {
		result.TmdbId = gen.NewOptInt(s.TmdbID)
	}
	if s.TvdbID > 0 {
		result.TvdbId = gen.NewOptInt(s.TvdbID)
	}
	if s.ImdbID != "" {
		result.ImdbId = gen.NewOptString(s.ImdbID)
	}
	result.DateAdded = gen.NewOptDateTime(s.DateAdded)

	if userData != nil {
		result.UserData = gen.NewOptSeriesUserData(*userData)
	}

	return result
}

// seriesFullToAPI converts a tvshow.Series with relations to a gen.SeriesFull.
func seriesFullToAPI(s *tvshow.Series, userData *gen.SeriesUserData) gen.SeriesFull {
	base := seriesToAPI(s, userData)
	result := gen.SeriesFull{
		ID:               base.ID,
		LibraryId:        base.LibraryId,
		Title:            base.Title,
		OriginalTitle:    base.OriginalTitle,
		SortTitle:        base.SortTitle,
		Overview:         base.Overview,
		Tagline:          base.Tagline,
		FirstAirDate:     base.FirstAirDate,
		LastAirDate:      base.LastAirDate,
		ContentRating:    base.ContentRating,
		SeasonCount:      base.SeasonCount,
		EpisodeCount:     base.EpisodeCount,
		CommunityRating:  base.CommunityRating,
		VoteCount:        base.VoteCount,
		PosterPath:       base.PosterPath,
		PosterBlurhash:   base.PosterBlurhash,
		BackdropPath:     base.BackdropPath,
		BackdropBlurhash: base.BackdropBlurhash,
		LogoPath:         base.LogoPath,
		TmdbId:           base.TmdbId,
		TvdbId:           base.TvdbId,
		ImdbId:           base.ImdbId,
		DateAdded:        base.DateAdded,
		UserData:         base.UserData,
	}

	// Convert status (different enum type)
	if s.Status != "" {
		result.Status = gen.NewOptSeriesFullStatus(gen.SeriesFullStatus(s.Status))
	}

	// Convert seasons
	if len(s.Seasons) > 0 {
		seasons := make([]gen.Season, len(s.Seasons))
		for i, season := range s.Seasons {
			seasons[i] = seasonToAPI(&season)
		}
		result.Seasons = seasons
	}

	// Convert genres
	if len(s.Genres) > 0 {
		genres := make([]gen.Genre, len(s.Genres))
		for i, g := range s.Genres {
			genres[i] = gen.Genre{
				ID:   g.ID,
				Name: g.Name,
			}
		}
		result.Genres = genres
	}

	// Convert cast
	if len(s.Cast) > 0 {
		cast := make([]gen.CastMember, len(s.Cast))
		for i, c := range s.Cast {
			cast[i] = gen.CastMember{
				PersonId:  c.PersonID,
				Name:      c.Name,
				Character: c.CharacterName,
			}
			if c.BillingOrder > 0 {
				cast[i].Order = gen.NewOptInt(c.BillingOrder)
			}
			if c.PrimaryImageURL != "" {
				cast[i].ProfilePath = gen.NewOptString(c.PrimaryImageURL)
			}
		}
		result.Cast = cast
	}

	return result
}

// seriesWithProgressToAPI converts a series with watch progress to API type.
func seriesWithProgressToAPI(s *tvshow.Series, progress *tvshow.SeriesWatchProgress) gen.SeriesWithProgress {
	base := seriesToAPI(s, nil)
	result := gen.SeriesWithProgress{
		ID:               base.ID,
		LibraryId:        base.LibraryId,
		Title:            base.Title,
		OriginalTitle:    base.OriginalTitle,
		SortTitle:        base.SortTitle,
		Tagline:          base.Tagline,
		Overview:         base.Overview,
		FirstAirDate:     base.FirstAirDate,
		LastAirDate:      base.LastAirDate,
		ContentRating:    base.ContentRating,
		SeasonCount:      base.SeasonCount,
		EpisodeCount:     base.EpisodeCount,
		CommunityRating:  base.CommunityRating,
		VoteCount:        base.VoteCount,
		PosterPath:       base.PosterPath,
		PosterBlurhash:   base.PosterBlurhash,
		BackdropPath:     base.BackdropPath,
		BackdropBlurhash: base.BackdropBlurhash,
		LogoPath:         base.LogoPath,
		TmdbId:           base.TmdbId,
		TvdbId:           base.TvdbId,
		ImdbId:           base.ImdbId,
		DateAdded:        base.DateAdded,
		UserData:         base.UserData,
	}

	// Convert status (different enum type)
	if s.Status != "" {
		result.Status = gen.NewOptSeriesWithProgressStatus(gen.SeriesWithProgressStatus(s.Status))
	}

	if progress != nil {
		result.Progress = gen.NewOptSeriesWatchProgress(gen.SeriesWatchProgress{
			EpisodesWatched: gen.NewOptInt(progress.WatchedEpisodes),
			EpisodesTotal:   gen.NewOptInt(progress.TotalEpisodes),
			PercentComplete: gen.NewOptFloat32(float32(progress.ProgressPercent)),
		})
		if progress.LastWatchedAt != nil {
			result.Progress.Value.LastWatchedAt = gen.NewOptDateTime(*progress.LastWatchedAt)
		}
	}

	return result
}

// seasonToAPI converts a tvshow.Season to a gen.Season.
func seasonToAPI(s *tvshow.Season) gen.Season {
	result := gen.Season{
		ID:           s.ID,
		SeriesId:     s.SeriesID,
		SeasonNumber: s.SeasonNumber,
	}

	if s.Name != "" {
		result.Name = gen.NewOptString(s.Name)
	}
	if s.Overview != "" {
		result.Overview = gen.NewOptString(s.Overview)
	}
	if s.AirDate != nil {
		result.AirDate = gen.NewOptDate(*s.AirDate)
	}
	if s.EpisodeCount > 0 {
		result.EpisodeCount = gen.NewOptInt(s.EpisodeCount)
	}
	if s.PosterPath != "" {
		result.PosterPath = gen.NewOptString(s.PosterPath)
	}
	if s.TmdbID > 0 {
		result.TmdbId = gen.NewOptInt(s.TmdbID)
	}
	if s.TvdbID > 0 {
		result.TvdbId = gen.NewOptInt(s.TvdbID)
	}

	return result
}

// seasonFullToAPI converts a tvshow.Season with episodes to a gen.SeasonFull.
func seasonFullToAPI(s *tvshow.Season) gen.SeasonFull {
	base := seasonToAPI(s)
	result := gen.SeasonFull{
		ID:           base.ID,
		SeriesId:     base.SeriesId,
		SeasonNumber: base.SeasonNumber,
		Name:         base.Name,
		Overview:     base.Overview,
		AirDate:      base.AirDate,
		EpisodeCount: base.EpisodeCount,
		PosterPath:   base.PosterPath,
		TmdbId:       base.TmdbId,
		TvdbId:       base.TvdbId,
	}

	if len(s.Episodes) > 0 {
		episodes := make([]gen.Episode, len(s.Episodes))
		for i, ep := range s.Episodes {
			episodes[i] = episodeToAPI(&ep)
		}
		result.Episodes = episodes
	}

	return result
}

// episodeToAPI converts a tvshow.Episode to a gen.Episode.
func episodeToAPI(e *tvshow.Episode) gen.Episode {
	result := gen.Episode{
		ID:            e.ID,
		SeriesId:      e.SeriesID,
		SeasonId:      e.SeasonID,
		SeasonNumber:  e.SeasonNumber,
		EpisodeNumber: e.EpisodeNumber,
		Title:         e.Title,
	}

	if e.Overview != "" {
		result.Overview = gen.NewOptString(e.Overview)
	}
	if e.AirDate != nil {
		result.AirDate = gen.NewOptDate(*e.AirDate)
	}
	if e.RuntimeTicks > 0 {
		result.RuntimeTicks = gen.NewOptInt64(e.RuntimeTicks)
	}
	if e.CommunityRating > 0 {
		result.CommunityRating = gen.NewOptFloat32(float32(e.CommunityRating))
	}
	if e.VoteCount > 0 {
		result.VoteCount = gen.NewOptInt(e.VoteCount)
	}
	if e.StillPath != "" {
		result.StillPath = gen.NewOptString(e.StillPath)
	}
	if e.StillBlurhash != "" {
		result.StillBlurhash = gen.NewOptString(e.StillBlurhash)
	}
	if e.TmdbID > 0 {
		result.TmdbId = gen.NewOptInt(e.TmdbID)
	}
	if e.TvdbID > 0 {
		result.TvdbId = gen.NewOptInt(e.TvdbID)
	}
	if e.ImdbID != "" {
		result.ImdbId = gen.NewOptString(e.ImdbID)
	}

	return result
}

// episodeFullToAPI converts a tvshow.Episode with relations to a gen.EpisodeFull.
func episodeFullToAPI(e *tvshow.Episode, userData *gen.EpisodeUserData) gen.EpisodeFull {
	base := episodeToAPI(e)
	result := gen.EpisodeFull{
		ID:              base.ID,
		SeriesId:        base.SeriesId,
		SeasonId:        base.SeasonId,
		SeasonNumber:    base.SeasonNumber,
		EpisodeNumber:   base.EpisodeNumber,
		Title:           base.Title,
		Overview:        base.Overview,
		AirDate:         base.AirDate,
		RuntimeTicks:    base.RuntimeTicks,
		CommunityRating: base.CommunityRating,
		VoteCount:       base.VoteCount,
		StillPath:       base.StillPath,
		StillBlurhash:   base.StillBlurhash,
		TmdbId:          base.TmdbId,
		TvdbId:          base.TvdbId,
		ImdbId:          base.ImdbId,
	}

	// Convert cast
	if len(e.Cast) > 0 {
		cast := make([]gen.CastMember, len(e.Cast))
		for i, c := range e.Cast {
			cast[i] = gen.CastMember{
				PersonId:  c.PersonID,
				Name:      c.Name,
				Character: c.CharacterName,
			}
			if c.BillingOrder > 0 {
				cast[i].Order = gen.NewOptInt(c.BillingOrder)
			}
			if c.PrimaryImageURL != "" {
				cast[i].ProfilePath = gen.NewOptString(c.PrimaryImageURL)
			}
		}
		result.Cast = cast
	}

	// Convert directors
	if len(e.Directors) > 0 {
		directors := make([]gen.CrewMember, len(e.Directors))
		for i, d := range e.Directors {
			directors[i] = tvshowCrewMemberToAPI(&d)
		}
		result.Directors = directors
	}

	// Convert writers
	if len(e.Writers) > 0 {
		writers := make([]gen.CrewMember, len(e.Writers))
		for i, w := range e.Writers {
			writers[i] = tvshowCrewMemberToAPI(&w)
		}
		result.Writers = writers
	}

	if userData != nil {
		result.UserData = gen.NewOptEpisodeUserData(*userData)
	}

	return result
}

// episodeWithSeriesToAPI converts an episode with series info for listings.
func episodeWithSeriesToAPI(e *tvshow.Episode, svc *tvshow.Service) gen.EpisodeWithSeries {
	base := episodeToAPI(e)
	result := gen.EpisodeWithSeries{
		ID:              base.ID,
		SeriesId:        base.SeriesId,
		SeasonId:        base.SeasonId,
		SeasonNumber:    base.SeasonNumber,
		EpisodeNumber:   base.EpisodeNumber,
		Title:           base.Title,
		Overview:        base.Overview,
		AirDate:         base.AirDate,
		RuntimeTicks:    base.RuntimeTicks,
		CommunityRating: base.CommunityRating,
		VoteCount:       base.VoteCount,
		StillPath:       base.StillPath,
		StillBlurhash:   base.StillBlurhash,
		TmdbId:          base.TmdbId,
		TvdbId:          base.TvdbId,
		ImdbId:          base.ImdbId,
	}

	// Fetch series info if available
	if svc != nil {
		if series, err := svc.GetSeries(context.Background(), e.SeriesID); err == nil {
			result.SeriesTitle = gen.NewOptString(series.Title)
			if series.PosterPath != "" {
				result.SeriesPosterPath = gen.NewOptString(series.PosterPath)
			}
		}
	}

	return result
}

// tvshowCrewMemberToAPI converts a tvshow.CrewMember to a gen.CrewMember.
func tvshowCrewMemberToAPI(c *tvshow.CrewMember) gen.CrewMember {
	result := gen.CrewMember{
		PersonId: c.PersonID,
		Name:     c.Name,
		Job:      c.Job,
	}

	if c.Department != "" {
		result.Department = gen.NewOptString(c.Department)
	}
	if c.PrimaryImageURL != "" {
		result.ProfilePath = gen.NewOptString(c.PrimaryImageURL)
	}

	return result
}

// =============================================================================
// QAR (Adult Content) Converters
// These convert between internal QAR entities and public API types.
// Internal: expedition/voyage/crew/port/flag (QAR obfuscation)
// API: AdultMovie/AdultScene/AdultPerformer/AdultStudio/AdultTag
// =============================================================================

// expeditionToAPI converts an expedition.Expedition to a gen.AdultMovie.
func expeditionToAPI(e *expedition.Expedition) gen.AdultMovie {
	result := gen.AdultMovie{
		ID:        e.ID,
		LibraryId: e.FleetID,
		Title:     e.Title,
		Path:      e.Path,
	}

	if e.SortTitle != "" {
		result.SortTitle = gen.NewOptString(e.SortTitle)
	}
	if e.Overview != "" {
		result.Overview = gen.NewOptString(e.Overview)
	}
	if e.LaunchDate != nil {
		result.ReleaseDate = gen.NewOptDate(*e.LaunchDate)
	}
	if e.RuntimeTicks > 0 {
		result.RuntimeTicks = gen.NewOptInt64(e.RuntimeTicks)
	}
	if e.PortID != nil {
		result.StudioId = gen.NewOptUUID(*e.PortID)
	}
	if e.Director != "" {
		result.Director = gen.NewOptString(e.Director)
	}
	if e.Series != "" {
		result.Series = gen.NewOptString(e.Series)
	}
	if e.Coordinates != "" {
		result.Phash = gen.NewOptString(e.Coordinates)
	}
	if e.Charter != "" {
		result.StashdbId = gen.NewOptString(e.Charter)
	}
	if e.Registry != "" {
		result.TpdbId = gen.NewOptString(e.Registry)
	}
	if e.WhisparrID != nil {
		result.WhisparrId = gen.NewOptInt(*e.WhisparrID)
	}
	result.CreatedAt = gen.NewOptDateTime(e.CreatedAt)
	result.UpdatedAt = gen.NewOptDateTime(e.UpdatedAt)

	return result
}

// voyageToAPI converts a voyage.Voyage to a gen.AdultScene.
func voyageToAPI(v *voyage.Voyage) gen.AdultScene {
	result := gen.AdultScene{
		ID:        v.ID,
		LibraryId: v.FleetID,
		Title:     v.Title,
		Path:      v.Path,
	}

	if v.SortTitle != "" {
		result.SortTitle = gen.NewOptString(v.SortTitle)
	}
	if v.Overview != "" {
		result.Overview = gen.NewOptString(v.Overview)
	}
	if v.LaunchDate != nil {
		result.ReleaseDate = gen.NewOptDate(*v.LaunchDate)
	}
	if v.Distance > 0 {
		result.RuntimeMinutes = gen.NewOptInt(v.Distance)
	}
	if v.PortID != nil {
		result.StudioId = gen.NewOptUUID(*v.PortID)
	}
	if v.Coordinates != "" {
		result.Phash = gen.NewOptString(v.Coordinates)
	}
	if v.Oshash != "" {
		result.Oshash = gen.NewOptString(v.Oshash)
	}
	if v.Charter != "" {
		result.StashdbId = gen.NewOptString(v.Charter)
	}
	if v.Registry != "" {
		result.TpdbId = gen.NewOptString(v.Registry)
	}
	if v.WhisparrID != nil {
		result.WhisparrId = gen.NewOptInt(*v.WhisparrID)
	}
	result.CreatedAt = gen.NewOptDateTime(v.CreatedAt)
	result.UpdatedAt = gen.NewOptDateTime(v.UpdatedAt)

	return result
}

// crewToAPI converts a crew.Crew to a gen.AdultPerformer.
func crewToAPI(c *crew.Crew) gen.AdultPerformer {
	result := gen.AdultPerformer{
		ID:   c.ID,
		Name: c.Name,
	}

	if c.Disambiguation != "" {
		result.Disambiguation = gen.NewOptString(c.Disambiguation)
	}
	if c.Gender != "" {
		result.Gender = gen.NewOptString(c.Gender)
	}
	if c.Christening != nil {
		result.Birthdate = gen.NewOptDate(*c.Christening)
	}
	if c.DeathDate != nil {
		result.DeathDate = gen.NewOptDate(*c.DeathDate)
	}
	if c.BirthCity != "" {
		result.BirthCity = gen.NewOptString(c.BirthCity)
	}
	if c.Origin != "" {
		result.Ethnicity = gen.NewOptString(c.Origin)
	}
	if c.Nationality != "" {
		result.Nationality = gen.NewOptString(c.Nationality)
	}
	if c.Rigging != "" {
		result.HairColor = gen.NewOptString(c.Rigging)
	}
	if c.Compass != "" {
		result.EyeColor = gen.NewOptString(c.Compass)
	}
	if c.HeightCM != nil {
		result.HeightCm = gen.NewOptInt(*c.HeightCM)
	}
	if c.WeightKG != nil {
		result.WeightKg = gen.NewOptInt(*c.WeightKG)
	}
	if c.Measurements != "" {
		result.Measurements = gen.NewOptString(c.Measurements)
	}
	if c.CupSize != "" {
		result.CupSize = gen.NewOptString(c.CupSize)
	}
	if c.BreastType != "" {
		result.BreastType = gen.NewOptString(c.BreastType)
	}
	if c.Markings != "" {
		result.Tattoos = gen.NewOptString(c.Markings)
	}
	if c.Anchors != "" {
		result.Piercings = gen.NewOptString(c.Anchors)
	}
	if c.MaidenVoyage != nil {
		result.CareerStart = gen.NewOptInt(*c.MaidenVoyage)
	}
	if c.LastPort != nil {
		result.CareerEnd = gen.NewOptInt(*c.LastPort)
	}
	if c.Bio != "" {
		result.Bio = gen.NewOptString(c.Bio)
	}
	if c.Charter != "" {
		result.StashdbId = gen.NewOptString(c.Charter)
	}
	if c.Registry != "" {
		result.TpdbId = gen.NewOptString(c.Registry)
	}
	if c.Manifest != "" {
		result.FreeonesId = gen.NewOptString(c.Manifest)
	}
	if c.Twitter != "" {
		result.Twitter = gen.NewOptString(c.Twitter)
	}
	if c.Instagram != "" {
		result.Instagram = gen.NewOptString(c.Instagram)
	}
	if c.ImagePath != "" {
		result.ImagePath = gen.NewOptString(c.ImagePath)
	}
	result.CreatedAt = gen.NewOptDateTime(c.CreatedAt)
	result.UpdatedAt = gen.NewOptDateTime(c.UpdatedAt)

	return result
}

// portToAPI converts a port.Port to a gen.AdultStudio.
func portToAPI(p *port.Port) gen.AdultStudio {
	result := gen.AdultStudio{
		ID:   p.ID,
		Name: p.Name,
	}

	if p.ParentID != nil {
		result.ParentId = gen.NewOptUUID(*p.ParentID)
	}
	if p.StashDBID != "" {
		result.StashdbId = gen.NewOptString(p.StashDBID)
	}
	if p.TPDBID != "" {
		result.TpdbId = gen.NewOptString(p.TPDBID)
	}
	if p.URL != "" {
		result.URL = gen.NewOptString(p.URL)
	}
	if p.LogoPath != "" {
		result.LogoPath = gen.NewOptString(p.LogoPath)
	}
	result.CreatedAt = gen.NewOptDateTime(p.CreatedAt)
	result.UpdatedAt = gen.NewOptDateTime(p.UpdatedAt)

	return result
}

// flagToAPI converts a flag.Flag to a gen.AdultTag.
func flagToAPI(f *flag.Flag) gen.AdultTag {
	result := gen.AdultTag{
		ID:   f.ID,
		Name: f.Name,
	}

	if f.Description != "" {
		result.Description = gen.NewOptString(f.Description)
	}
	if f.ParentID != nil {
		result.ParentId = gen.NewOptUUID(*f.ParentID)
	}
	if f.StashDBID != "" {
		result.StashdbId = gen.NewOptString(f.StashDBID)
	}
	// Note: f.Waters (category) is not exposed in the API schema

	return result
}
