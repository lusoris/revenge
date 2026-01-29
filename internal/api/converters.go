package api

import (
	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/content/movie"
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

// libraryToAPI converts a db.Library to a gen.Library.
func libraryToAPI(l *db.Library) gen.Library {
	result := gen.Library{
		ID:    l.ID,
		Name:  l.Name,
		Type:  gen.LibraryType(l.Type),
		Paths: l.Paths,
	}

	result.ScanEnabled = gen.NewOptBool(l.ScanEnabled)
	result.ScanIntervalHours = gen.NewOptInt(int(l.ScanIntervalHours))

	if l.LastScanAt.Valid {
		result.LastScanAt = gen.NewOptDateTime(l.LastScanAt.Time)
	}

	if l.PreferredLanguage != nil {
		result.PreferredLanguage = gen.NewOptString(*l.PreferredLanguage)
	}

	result.DownloadImages = gen.NewOptBool(l.DownloadImages)
	result.DownloadNfo = gen.NewOptBool(l.DownloadNfo)
	result.GenerateChapters = gen.NewOptBool(l.GenerateChapters)
	result.IsPrivate = gen.NewOptBool(l.IsPrivate)

	if l.OwnerUserID.Valid {
		result.OwnerUserId = gen.NewOptUUID(l.OwnerUserID.Bytes)
	}

	result.SortOrder = gen.NewOptInt(int(l.SortOrder))

	if l.Icon != nil {
		result.Icon = gen.NewOptString(*l.Icon)
	}

	result.CreatedAt = gen.NewOptDateTime(l.CreatedAt)
	result.UpdatedAt = gen.NewOptDateTime(l.UpdatedAt)

	return result
}

// movieToAPI converts a movie.Movie to a gen.Movie.
func movieToAPI(m *movie.Movie, userData *gen.MovieUserData) gen.Movie {
	result := gen.Movie{
		ID:        m.ID,
		LibraryId: m.LibraryID,
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
