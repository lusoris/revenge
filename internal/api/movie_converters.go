package api

import (
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content"
	"github.com/lusoris/revenge/internal/content/movie"
)

// externalRatingsToOgen converts domain external ratings to ogen format.
// Works with both movie.ExternalRating and tvshow.ExternalRating (both are
// type aliases for content.ExternalRating).
func externalRatingsToOgen(ratings []content.ExternalRating) []ogen.ExternalRating {
	if len(ratings) == 0 {
		return nil
	}
	out := make([]ogen.ExternalRating, len(ratings))
	for i, r := range ratings {
		out[i] = ogen.ExternalRating{
			Source: r.Source,
			Value:  r.Value,
			Score:  float32(r.Score),
		}
	}
	return out
}

// movieToOgen converts a movie domain type to ogen Movie.
func movieToOgen(m *movie.Movie) *ogen.Movie {
	o := &ogen.Movie{
		ID:             ogen.NewOptUUID(m.ID),
		Title:          ogen.NewOptString(m.Title),
		CreatedAt:      ogen.NewOptDateTime(m.CreatedAt),
		UpdatedAt:      ogen.NewOptDateTime(m.UpdatedAt),
		LibraryAddedAt: ogen.NewOptDateTime(m.LibraryAddedAt),
	}

	if m.TMDbID != nil {
		o.TmdbID.SetTo(int(*m.TMDbID))
	}
	setOpt(&o.ImdbID, m.IMDbID)
	setOpt(&o.OriginalTitle, m.OriginalTitle)
	setOptConv(&o.Year, m.Year, int32ToInt)
	setOpt(&o.ReleaseDate, m.ReleaseDate)
	setOptConv(&o.Runtime, m.Runtime, int32ToInt)
	setOpt(&o.Overview, m.Overview)
	setOpt(&o.Tagline, m.Tagline)
	setOpt(&o.Status, m.Status)
	setOpt(&o.OriginalLanguage, m.OriginalLanguage)
	setOpt(&o.PosterPath, m.PosterPath)
	setOpt(&o.BackdropPath, m.BackdropPath)
	setOpt(&o.TrailerURL, m.TrailerURL)
	setOptDecimalFloat32(&o.VoteAverage, m.VoteAverage)
	setOptConv(&o.VoteCount, m.VoteCount, int32ToInt)
	setOptDecimalFloat32(&o.Popularity, m.Popularity)
	if m.Budget != nil {
		o.Budget.SetTo(int64(*m.Budget))
	}
	if m.Revenue != nil {
		o.Revenue.SetTo(int64(*m.Revenue))
	}
	setOpt(&o.MetadataUpdatedAt, m.MetadataUpdatedAt)
	setOptConv(&o.RadarrID, m.RadarrID, int32ToInt)

	o.ExternalRatings = externalRatingsToOgen(m.ExternalRatings)

	return o
}

// movieFileToOgen converts a movie file domain type to ogen MovieFile.
func movieFileToOgen(f *movie.MovieFile) *ogen.MovieFile {
	o := &ogen.MovieFile{
		ID:        ogen.NewOptUUID(f.ID),
		MovieID:   ogen.NewOptUUID(f.MovieID),
		FilePath:  ogen.NewOptString(f.FilePath),
		FileSize:  ogen.NewOptInt64(f.FileSize),
		FileName:  ogen.NewOptString(f.FileName),
		CreatedAt: ogen.NewOptDateTime(f.CreatedAt),
		UpdatedAt: ogen.NewOptDateTime(f.UpdatedAt),
	}

	if f.Resolution != nil {
		o.Resolution.SetTo(*f.Resolution)
	}
	if f.QualityProfile != nil {
		o.QualityProfile.SetTo(*f.QualityProfile)
	}
	if f.VideoCodec != nil {
		o.VideoCodec.SetTo(*f.VideoCodec)
	}
	if f.AudioCodec != nil {
		o.AudioCodec.SetTo(*f.AudioCodec)
	}
	if f.Container != nil {
		o.Container.SetTo(*f.Container)
	}
	setOptConv(&o.DurationSeconds, f.DurationSeconds, int32ToInt)
	setOptConv(&o.BitrateKbps, f.BitrateKbps, int32ToInt)
	setOptDecimalFloat32(&o.Framerate, f.Framerate)
	setOpt(&o.DynamicRange, f.DynamicRange)
	setOpt(&o.ColorSpace, f.ColorSpace)
	setOpt(&o.AudioChannels, f.AudioChannels)
	if len(f.AudioLanguages) > 0 {
		o.AudioLanguages = f.AudioLanguages
	}
	if len(f.SubtitleLanguages) > 0 {
		o.SubtitleLanguages = f.SubtitleLanguages
	}
	setOptConv(&o.RadarrFileID, f.RadarrFileID, int32ToInt)
	setOpt(&o.LastScannedAt, f.LastScannedAt)
	setOpt(&o.IsMonitored, f.IsMonitored)

	return o
}

// movieCreditToOgen converts a movie credit domain type to ogen MovieCredit.
func movieCreditToOgen(c *movie.MovieCredit) *ogen.MovieCredit {
	o := &ogen.MovieCredit{
		ID:           ogen.NewOptUUID(c.ID),
		MovieID:      ogen.NewOptUUID(c.MovieID),
		TmdbPersonID: ogen.NewOptInt(int(c.TMDbPersonID)),
		Name:         ogen.NewOptString(c.Name),
		CreditType:   ogen.NewOptMovieCreditCreditType(ogen.MovieCreditCreditType(c.CreditType)),
		CreatedAt:    ogen.NewOptDateTime(c.CreatedAt),
		UpdatedAt:    ogen.NewOptDateTime(c.UpdatedAt),
	}

	setOpt(&o.Character, c.Character)
	setOpt(&o.Job, c.Job)
	setOpt(&o.Department, c.Department)
	setOptConv(&o.CastOrder, c.CastOrder, int32ToInt)
	setOpt(&o.ProfilePath, c.ProfilePath)

	return o
}

// movieCollectionToOgen converts a movie collection domain type to ogen MovieCollection.
func movieCollectionToOgen(c *movie.MovieCollection) *ogen.MovieCollection {
	o := &ogen.MovieCollection{
		ID:        ogen.NewOptUUID(c.ID),
		Name:      ogen.NewOptString(c.Name),
		CreatedAt: ogen.NewOptDateTime(c.CreatedAt),
		UpdatedAt: ogen.NewOptDateTime(c.UpdatedAt),
	}

	setOptConv(&o.TmdbCollectionID, c.TMDbCollectionID, int32ToInt)
	setOpt(&o.Overview, c.Overview)
	setOpt(&o.PosterPath, c.PosterPath)
	setOpt(&o.BackdropPath, c.BackdropPath)

	return o
}

// movieGenreToOgen converts a movie genre domain type to ogen MovieGenre.
func movieGenreToOgen(g *movie.MovieGenre) *ogen.MovieGenre {
	return &ogen.MovieGenre{
		ID:        ogen.NewOptUUID(g.ID),
		MovieID:   ogen.NewOptUUID(g.MovieID),
		Slug:      ogen.NewOptString(g.Slug),
		Name:      ogen.NewOptString(g.Name),
		CreatedAt: ogen.NewOptDateTime(g.CreatedAt),
	}
}

// movieWatchedToOgen converts a movie watched domain type to ogen MovieWatched.
func movieWatchedToOgen(w *movie.MovieWatched) *ogen.MovieWatched {
	o := &ogen.MovieWatched{
		ID:              ogen.NewOptUUID(w.ID),
		UserID:          ogen.NewOptUUID(w.UserID),
		MovieID:         ogen.NewOptUUID(w.MovieID),
		ProgressSeconds: ogen.NewOptInt(int(w.ProgressSeconds)),
		DurationSeconds: ogen.NewOptInt(int(w.DurationSeconds)),
		IsCompleted:     ogen.NewOptBool(w.IsCompleted),
		WatchCount:      ogen.NewOptInt(int(w.WatchCount)),
		LastWatchedAt:   ogen.NewOptDateTime(w.LastWatchedAt),
		CreatedAt:       ogen.NewOptDateTime(w.CreatedAt),
		UpdatedAt:       ogen.NewOptDateTime(w.UpdatedAt),
	}

	setOptConv(&o.ProgressPercent, w.ProgressPercent, int32ToInt)

	return o
}

// continueWatchingItemToOgen converts a continue watching item to ogen type.
func continueWatchingItemToOgen(item *movie.ContinueWatchingItem) ogen.ContinueWatchingItem {
	o := ogen.ContinueWatchingItem{
		ID:              ogen.NewOptUUID(item.ID),
		Title:           ogen.NewOptString(item.Title),
		CreatedAt:       ogen.NewOptDateTime(item.CreatedAt),
		UpdatedAt:       ogen.NewOptDateTime(item.UpdatedAt),
		LibraryAddedAt:  ogen.NewOptDateTime(item.LibraryAddedAt),
		ProgressSeconds: ogen.NewOptInt(int(item.ProgressSeconds)),
		DurationSeconds: ogen.NewOptInt(int(item.DurationSeconds)),
		LastWatchedAt:   ogen.NewOptDateTime(item.LastWatchedAt),
	}

	// Copy all movie fields
	setOptConv(&o.TmdbID, item.TMDbID, int32ToInt)
	setOpt(&o.ImdbID, item.IMDbID)
	setOpt(&o.OriginalTitle, item.OriginalTitle)
	setOptConv(&o.Year, item.Year, int32ToInt)
	setOpt(&o.ReleaseDate, item.ReleaseDate)
	setOptConv(&o.Runtime, item.Runtime, int32ToInt)
	setOpt(&o.Overview, item.Overview)
	setOpt(&o.Tagline, item.Tagline)
	setOpt(&o.Status, item.Status)
	setOpt(&o.OriginalLanguage, item.OriginalLanguage)
	setOpt(&o.PosterPath, item.PosterPath)
	setOpt(&o.BackdropPath, item.BackdropPath)
	setOpt(&o.TrailerURL, item.TrailerURL)
	setOptDecimalFloat32(&o.VoteAverage, item.VoteAverage)
	setOptConv(&o.VoteCount, item.VoteCount, int32ToInt)
	setOptDecimalFloat32(&o.Popularity, item.Popularity)
	if item.Budget != nil {
		o.Budget.SetTo(int64(*item.Budget))
	}
	if item.Revenue != nil {
		o.Revenue.SetTo(int64(*item.Revenue))
	}
	setOpt(&o.MetadataUpdatedAt, item.MetadataUpdatedAt)
	setOptConv(&o.RadarrID, item.RadarrID, int32ToInt)
	setOptConv(&o.ProgressPercent, item.ProgressPercent, int32ToInt)

	o.ExternalRatings = externalRatingsToOgen(item.ExternalRatings)

	return o
}

// watchedMovieItemToOgen converts a watched movie item to ogen type.
func watchedMovieItemToOgen(item *movie.WatchedMovieItem) ogen.WatchedMovieItem {
	o := ogen.WatchedMovieItem{
		ID:             ogen.NewOptUUID(item.ID),
		Title:          ogen.NewOptString(item.Title),
		CreatedAt:      ogen.NewOptDateTime(item.CreatedAt),
		UpdatedAt:      ogen.NewOptDateTime(item.UpdatedAt),
		LibraryAddedAt: ogen.NewOptDateTime(item.LibraryAddedAt),
		WatchCount:     ogen.NewOptInt(int(item.WatchCount)),
		LastWatchedAt:  ogen.NewOptDateTime(item.LastWatchedAt),
	}

	// Copy all movie fields
	setOptConv(&o.TmdbID, item.TMDbID, int32ToInt)
	setOpt(&o.ImdbID, item.IMDbID)
	setOpt(&o.OriginalTitle, item.OriginalTitle)
	setOptConv(&o.Year, item.Year, int32ToInt)
	setOpt(&o.ReleaseDate, item.ReleaseDate)
	setOptConv(&o.Runtime, item.Runtime, int32ToInt)
	setOpt(&o.Overview, item.Overview)
	setOpt(&o.Tagline, item.Tagline)
	setOpt(&o.Status, item.Status)
	setOpt(&o.OriginalLanguage, item.OriginalLanguage)
	setOpt(&o.PosterPath, item.PosterPath)
	setOpt(&o.BackdropPath, item.BackdropPath)
	setOpt(&o.TrailerURL, item.TrailerURL)
	setOptDecimalFloat32(&o.VoteAverage, item.VoteAverage)
	setOptConv(&o.VoteCount, item.VoteCount, int32ToInt)
	setOptDecimalFloat32(&o.Popularity, item.Popularity)
	if item.Budget != nil {
		o.Budget.SetTo(int64(*item.Budget))
	}
	if item.Revenue != nil {
		o.Revenue.SetTo(int64(*item.Revenue))
	}
	setOpt(&o.MetadataUpdatedAt, item.MetadataUpdatedAt)
	setOptConv(&o.RadarrID, item.RadarrID, int32ToInt)

	o.ExternalRatings = externalRatingsToOgen(item.ExternalRatings)

	return o
}

// userMovieStatsToOgen converts user movie stats to ogen type.
func userMovieStatsToOgen(stats *movie.UserMovieStats) *ogen.UserMovieStats {
	o := &ogen.UserMovieStats{
		WatchedCount:    ogen.NewOptInt64(stats.WatchedCount),
		InProgressCount: ogen.NewOptInt64(stats.InProgressCount),
	}

	setOpt(&o.TotalWatches, stats.TotalWatches)

	return o
}
