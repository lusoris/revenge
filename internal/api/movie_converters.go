package api

import (
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/movie"
)

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
	if m.IMDbID != nil {
		o.ImdbID.SetTo(*m.IMDbID)
	}
	if m.OriginalTitle != nil {
		o.OriginalTitle.SetTo(*m.OriginalTitle)
	}
	if m.Year != nil {
		o.Year.SetTo(int(*m.Year))
	}
	if m.ReleaseDate != nil {
		o.ReleaseDate.SetTo(*m.ReleaseDate)
	}
	if m.Runtime != nil {
		o.Runtime.SetTo(int(*m.Runtime))
	}
	if m.Overview != nil {
		o.Overview.SetTo(*m.Overview)
	}
	if m.Tagline != nil {
		o.Tagline.SetTo(*m.Tagline)
	}
	if m.Status != nil {
		o.Status.SetTo(*m.Status)
	}
	if m.OriginalLanguage != nil {
		o.OriginalLanguage.SetTo(*m.OriginalLanguage)
	}
	if m.PosterPath != nil {
		o.PosterPath.SetTo(*m.PosterPath)
	}
	if m.BackdropPath != nil {
		o.BackdropPath.SetTo(*m.BackdropPath)
	}
	if m.TrailerURL != nil {
		o.TrailerURL.SetTo(*m.TrailerURL)
	}
	if m.VoteAverage != nil {
		f, _ := m.VoteAverage.Float64()
		o.VoteAverage.SetTo(float32(f))
	}
	if m.VoteCount != nil {
		o.VoteCount.SetTo(int(*m.VoteCount))
	}
	if m.Popularity != nil {
		f, _ := m.Popularity.Float64()
		o.Popularity.SetTo(float32(f))
	}
	if m.Budget != nil {
		o.Budget.SetTo(int64(*m.Budget))
	}
	if m.Revenue != nil {
		o.Revenue.SetTo(int64(*m.Revenue))
	}
	if m.MetadataUpdatedAt != nil {
		o.MetadataUpdatedAt.SetTo(*m.MetadataUpdatedAt)
	}
	if m.RadarrID != nil {
		o.RadarrID.SetTo(int(*m.RadarrID))
	}

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
	if f.DurationSeconds != nil {
		o.DurationSeconds.SetTo(int(*f.DurationSeconds))
	}
	if f.BitrateKbps != nil {
		o.BitrateKbps.SetTo(int(*f.BitrateKbps))
	}
	if f.Framerate != nil {
		fr, _ := f.Framerate.Float64()
		o.Framerate.SetTo(float32(fr))
	}
	if f.DynamicRange != nil {
		o.DynamicRange.SetTo(*f.DynamicRange)
	}
	if f.ColorSpace != nil {
		o.ColorSpace.SetTo(*f.ColorSpace)
	}
	if f.AudioChannels != nil {
		o.AudioChannels.SetTo(*f.AudioChannels)
	}
	if f.AudioLanguages != nil && len(f.AudioLanguages) > 0 {
		o.AudioLanguages = f.AudioLanguages
	}
	if f.SubtitleLanguages != nil && len(f.SubtitleLanguages) > 0 {
		o.SubtitleLanguages = f.SubtitleLanguages
	}
	if f.RadarrFileID != nil {
		o.RadarrFileID.SetTo(int(*f.RadarrFileID))
	}
	if f.LastScannedAt != nil {
		o.LastScannedAt.SetTo(*f.LastScannedAt)
	}
	if f.IsMonitored != nil {
		o.IsMonitored.SetTo(*f.IsMonitored)
	}

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

	if c.Character != nil {
		o.Character.SetTo(*c.Character)
	}
	if c.Job != nil {
		o.Job.SetTo(*c.Job)
	}
	if c.Department != nil {
		o.Department.SetTo(*c.Department)
	}
	if c.CastOrder != nil {
		o.CastOrder.SetTo(int(*c.CastOrder))
	}
	if c.ProfilePath != nil {
		o.ProfilePath.SetTo(*c.ProfilePath)
	}

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

	if c.TMDbCollectionID != nil {
		o.TmdbCollectionID.SetTo(int(*c.TMDbCollectionID))
	}
	if c.Overview != nil {
		o.Overview.SetTo(*c.Overview)
	}
	if c.PosterPath != nil {
		o.PosterPath.SetTo(*c.PosterPath)
	}
	if c.BackdropPath != nil {
		o.BackdropPath.SetTo(*c.BackdropPath)
	}

	return o
}

// movieGenreToOgen converts a movie genre domain type to ogen MovieGenre.
func movieGenreToOgen(g *movie.MovieGenre) *ogen.MovieGenre {
	return &ogen.MovieGenre{
		ID:          ogen.NewOptUUID(g.ID),
		MovieID:     ogen.NewOptUUID(g.MovieID),
		TmdbGenreID: ogen.NewOptInt(int(g.TMDbGenreID)),
		Name:        ogen.NewOptString(g.Name),
		CreatedAt:   ogen.NewOptDateTime(g.CreatedAt),
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

	if w.ProgressPercent != nil {
		o.ProgressPercent.SetTo(int(*w.ProgressPercent))
	}

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
	if item.TMDbID != nil {
		o.TmdbID.SetTo(int(*item.TMDbID))
	}
	if item.IMDbID != nil {
		o.ImdbID.SetTo(*item.IMDbID)
	}
	if item.OriginalTitle != nil {
		o.OriginalTitle.SetTo(*item.OriginalTitle)
	}
	if item.Year != nil {
		o.Year.SetTo(int(*item.Year))
	}
	if item.ReleaseDate != nil {
		o.ReleaseDate.SetTo(*item.ReleaseDate)
	}
	if item.Runtime != nil {
		o.Runtime.SetTo(int(*item.Runtime))
	}
	if item.Overview != nil {
		o.Overview.SetTo(*item.Overview)
	}
	if item.Tagline != nil {
		o.Tagline.SetTo(*item.Tagline)
	}
	if item.Status != nil {
		o.Status.SetTo(*item.Status)
	}
	if item.OriginalLanguage != nil {
		o.OriginalLanguage.SetTo(*item.OriginalLanguage)
	}
	if item.PosterPath != nil {
		o.PosterPath.SetTo(*item.PosterPath)
	}
	if item.BackdropPath != nil {
		o.BackdropPath.SetTo(*item.BackdropPath)
	}
	if item.TrailerURL != nil {
		o.TrailerURL.SetTo(*item.TrailerURL)
	}
	if item.VoteAverage != nil {
		f, _ := item.VoteAverage.Float64()
		o.VoteAverage.SetTo(float32(f))
	}
	if item.VoteCount != nil {
		o.VoteCount.SetTo(int(*item.VoteCount))
	}
	if item.Popularity != nil {
		f, _ := item.Popularity.Float64()
		o.Popularity.SetTo(float32(f))
	}
	if item.Budget != nil {
		o.Budget.SetTo(int64(*item.Budget))
	}
	if item.Revenue != nil {
		o.Revenue.SetTo(int64(*item.Revenue))
	}
	if item.MetadataUpdatedAt != nil {
		o.MetadataUpdatedAt.SetTo(*item.MetadataUpdatedAt)
	}
	if item.RadarrID != nil {
		o.RadarrID.SetTo(int(*item.RadarrID))
	}
	if item.ProgressPercent != nil {
		o.ProgressPercent.SetTo(int(*item.ProgressPercent))
	}

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
	if item.TMDbID != nil {
		o.TmdbID.SetTo(int(*item.TMDbID))
	}
	if item.IMDbID != nil {
		o.ImdbID.SetTo(*item.IMDbID)
	}
	if item.OriginalTitle != nil {
		o.OriginalTitle.SetTo(*item.OriginalTitle)
	}
	if item.Year != nil {
		o.Year.SetTo(int(*item.Year))
	}
	if item.ReleaseDate != nil {
		o.ReleaseDate.SetTo(*item.ReleaseDate)
	}
	if item.Runtime != nil {
		o.Runtime.SetTo(int(*item.Runtime))
	}
	if item.Overview != nil {
		o.Overview.SetTo(*item.Overview)
	}
	if item.Tagline != nil {
		o.Tagline.SetTo(*item.Tagline)
	}
	if item.Status != nil {
		o.Status.SetTo(*item.Status)
	}
	if item.OriginalLanguage != nil {
		o.OriginalLanguage.SetTo(*item.OriginalLanguage)
	}
	if item.PosterPath != nil {
		o.PosterPath.SetTo(*item.PosterPath)
	}
	if item.BackdropPath != nil {
		o.BackdropPath.SetTo(*item.BackdropPath)
	}
	if item.TrailerURL != nil {
		o.TrailerURL.SetTo(*item.TrailerURL)
	}
	if item.VoteAverage != nil {
		f, _ := item.VoteAverage.Float64()
		o.VoteAverage.SetTo(float32(f))
	}
	if item.VoteCount != nil {
		o.VoteCount.SetTo(int(*item.VoteCount))
	}
	if item.Popularity != nil {
		f, _ := item.Popularity.Float64()
		o.Popularity.SetTo(float32(f))
	}
	if item.Budget != nil {
		o.Budget.SetTo(int64(*item.Budget))
	}
	if item.Revenue != nil {
		o.Revenue.SetTo(int64(*item.Revenue))
	}
	if item.MetadataUpdatedAt != nil {
		o.MetadataUpdatedAt.SetTo(*item.MetadataUpdatedAt)
	}
	if item.RadarrID != nil {
		o.RadarrID.SetTo(int(*item.RadarrID))
	}

	return o
}

// userMovieStatsToOgen converts user movie stats to ogen type.
func userMovieStatsToOgen(stats *movie.UserMovieStats) *ogen.UserMovieStats {
	o := &ogen.UserMovieStats{
		WatchedCount:    ogen.NewOptInt64(stats.WatchedCount),
		InProgressCount: ogen.NewOptInt64(stats.InProgressCount),
	}

	if stats.TotalWatches != nil {
		o.TotalWatches.SetTo(*stats.TotalWatches)
	}

	return o
}
