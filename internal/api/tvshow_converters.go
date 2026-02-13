package api

import (
	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/content/tvshow"
)

// seriesToOgen converts a series domain type to ogen TVSeries.
func seriesToOgen(s *tvshow.Series) *ogen.TVSeries {
	o := &ogen.TVSeries{
		ID:           ogen.NewOptUUID(s.ID),
		Title:        ogen.NewOptString(s.Title),
		TotalSeasons: ogen.NewOptInt(int(s.TotalSeasons)),
		TotalEpisodes: ogen.NewOptInt(int(s.TotalEpisodes)),
		CreatedAt:    ogen.NewOptDateTime(s.CreatedAt),
		UpdatedAt:    ogen.NewOptDateTime(s.UpdatedAt),
	}

	setOptConv(&o.TmdbID, s.TMDbID, int32ToInt)
	setOptConv(&o.TvdbID, s.TVDbID, int32ToInt)
	setOpt(&o.ImdbID, s.IMDbID)
	setOptConv(&o.SonarrID, s.SonarrID, int32ToInt)
	setOpt(&o.OriginalTitle, s.OriginalTitle)
	if s.OriginalLanguage != "" {
		o.OriginalLanguage.SetTo(s.OriginalLanguage)
	}
	setOpt(&o.Tagline, s.Tagline)
	setOpt(&o.Overview, s.Overview)
	setOpt(&o.Status, s.Status)
	setOpt(&o.Type, s.Type)
	setOpt(&o.FirstAirDate, s.FirstAirDate)
	setOpt(&o.LastAirDate, s.LastAirDate)
	setOptDecimalFloat32(&o.VoteAverage, s.VoteAverage)
	setOptConv(&o.VoteCount, s.VoteCount, int32ToInt)
	setOptDecimalFloat32(&o.Popularity, s.Popularity)
	setOpt(&o.PosterPath, s.PosterPath)
	setOpt(&o.BackdropPath, s.BackdropPath)
	setOpt(&o.TrailerURL, s.TrailerURL)
	setOpt(&o.Homepage, s.Homepage)
	setOpt(&o.MetadataUpdatedAt, s.MetadataUpdatedAt)

	o.ExternalRatings = externalRatingsToOgen(s.ExternalRatings)

	return o
}

// seasonToOgen converts a season domain type to ogen TVSeason.
func seasonToOgen(s *tvshow.Season) *ogen.TVSeason {
	o := &ogen.TVSeason{
		ID:           ogen.NewOptUUID(s.ID),
		SeriesID:     ogen.NewOptUUID(s.SeriesID),
		SeasonNumber: ogen.NewOptInt(int(s.SeasonNumber)),
		Name:         ogen.NewOptString(s.Name),
		EpisodeCount: ogen.NewOptInt(int(s.EpisodeCount)),
		CreatedAt:    ogen.NewOptDateTime(s.CreatedAt),
		UpdatedAt:    ogen.NewOptDateTime(s.UpdatedAt),
	}

	setOptConv(&o.TmdbID, s.TMDbID, int32ToInt)
	setOpt(&o.Overview, s.Overview)
	setOpt(&o.PosterPath, s.PosterPath)
	setOpt(&o.AirDate, s.AirDate)
	setOptDecimalFloat32(&o.VoteAverage, s.VoteAverage)

	return o
}

// episodeToOgen converts an episode domain type to ogen TVEpisode.
func episodeToOgen(e *tvshow.Episode) *ogen.TVEpisode {
	o := &ogen.TVEpisode{
		ID:            ogen.NewOptUUID(e.ID),
		SeriesID:      ogen.NewOptUUID(e.SeriesID),
		SeasonID:      ogen.NewOptUUID(e.SeasonID),
		SeasonNumber:  ogen.NewOptInt(int(e.SeasonNumber)),
		EpisodeNumber: ogen.NewOptInt(int(e.EpisodeNumber)),
		Title:         ogen.NewOptString(e.Title),
		CreatedAt:     ogen.NewOptDateTime(e.CreatedAt),
		UpdatedAt:     ogen.NewOptDateTime(e.UpdatedAt),
	}

	setOptConv(&o.TmdbID, e.TMDbID, int32ToInt)
	setOptConv(&o.TvdbID, e.TVDbID, int32ToInt)
	setOpt(&o.ImdbID, e.IMDbID)
	setOpt(&o.Overview, e.Overview)
	setOpt(&o.AirDate, e.AirDate)
	setOptConv(&o.Runtime, e.Runtime, int32ToInt)
	setOptDecimalFloat32(&o.VoteAverage, e.VoteAverage)
	setOptConv(&o.VoteCount, e.VoteCount, int32ToInt)
	setOpt(&o.StillPath, e.StillPath)
	setOpt(&o.ProductionCode, e.ProductionCode)

	return o
}

// episodeFileToOgen converts an episode file domain type to ogen TVEpisodeFile.
func episodeFileToOgen(f *tvshow.EpisodeFile) *ogen.TVEpisodeFile {
	o := &ogen.TVEpisodeFile{
		ID:        ogen.NewOptUUID(f.ID),
		EpisodeID: ogen.NewOptUUID(f.EpisodeID),
		FilePath:  ogen.NewOptString(f.FilePath),
		FileName:  ogen.NewOptString(f.FileName),
		FileSize:  ogen.NewOptInt64(f.FileSize),
		CreatedAt: ogen.NewOptDateTime(f.CreatedAt),
		UpdatedAt: ogen.NewOptDateTime(f.UpdatedAt),
	}

	setOpt(&o.Container, f.Container)
	setOpt(&o.Resolution, f.Resolution)
	setOpt(&o.QualityProfile, f.QualityProfile)
	setOpt(&o.VideoCodec, f.VideoCodec)
	setOpt(&o.AudioCodec, f.AudioCodec)
	setOptConv(&o.BitrateKbps, f.BitrateKbps, int32ToInt)
	if f.DurationSeconds != nil {
		d, _ := f.DurationSeconds.Float64()
		o.DurationSeconds.SetTo(int(d))
	}
	if len(f.AudioLanguages) > 0 {
		o.AudioLanguages = f.AudioLanguages
	}
	if len(f.SubtitleLanguages) > 0 {
		o.SubtitleLanguages = f.SubtitleLanguages
	}
	setOptConv(&o.SonarrFileID, f.SonarrFileID, int32ToInt)

	return o
}

// seriesCreditToOgen converts a series credit domain type to ogen TVSeriesCredit.
func seriesCreditToOgen(c *tvshow.SeriesCredit) *ogen.TVSeriesCredit {
	o := &ogen.TVSeriesCredit{
		ID:           ogen.NewOptUUID(c.ID),
		SeriesID:     ogen.NewOptUUID(c.SeriesID),
		TmdbPersonID: ogen.NewOptInt(int(c.TMDbPersonID)),
		Name:         ogen.NewOptString(c.Name),
		CreditType:   ogen.NewOptTVSeriesCreditCreditType(ogen.TVSeriesCreditCreditType(c.CreditType)),
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

// seriesGenreToOgen converts a series genre domain type to ogen TVGenre.
func seriesGenreToOgen(g *tvshow.SeriesGenre) *ogen.TVGenre {
	return &ogen.TVGenre{
		ID:        ogen.NewOptUUID(g.ID),
		SeriesID:  ogen.NewOptUUID(g.SeriesID),
		Slug:      ogen.NewOptString(g.Slug),
		Name:      ogen.NewOptString(g.Name),
		CreatedAt: ogen.NewOptDateTime(g.CreatedAt),
	}
}

// networkToOgen converts a network domain type to ogen TVNetwork.
func networkToOgen(n *tvshow.Network) *ogen.TVNetwork {
	o := &ogen.TVNetwork{
		ID:        ogen.NewOptUUID(n.ID),
		Name:      ogen.NewOptString(n.Name),
		CreatedAt: ogen.NewOptDateTime(n.CreatedAt),
	}

	setOptConv(&o.TmdbID, n.TMDbID, int32ToInt)
	setOpt(&o.LogoPath, n.LogoPath)
	setOpt(&o.OriginCountry, n.OriginCountry)

	return o
}

// episodeWatchProgressToOgen converts episode watch progress domain type to ogen EpisodeWatchProgress.
func episodeWatchProgressToOgen(w *tvshow.EpisodeWatched) *ogen.EpisodeWatchProgress {
	return &ogen.EpisodeWatchProgress{
		ID:              ogen.NewOptUUID(w.ID),
		UserID:          ogen.NewOptUUID(w.UserID),
		EpisodeID:       ogen.NewOptUUID(w.EpisodeID),
		ProgressSeconds: ogen.NewOptInt(int(w.ProgressSeconds)),
		DurationSeconds: ogen.NewOptInt(int(w.DurationSeconds)),
		IsCompleted:     ogen.NewOptBool(w.IsCompleted),
		LastWatchedAt:   ogen.NewOptDateTime(w.LastWatchedAt),
		CreatedAt:       ogen.NewOptDateTime(w.CreatedAt),
		UpdatedAt:       ogen.NewOptDateTime(w.UpdatedAt),
	}
}

// tvContinueWatchingItemToOgen converts a TV continue watching item to ogen TVContinueWatchingItem.
func tvContinueWatchingItemToOgen(item *tvshow.ContinueWatchingItem) *ogen.TVContinueWatchingItem {
	o := &ogen.TVContinueWatchingItem{
		EpisodeID:       ogen.NewOptUUID(item.LastEpisodeID),
		EpisodeTitle:    ogen.NewOptString(item.LastEpisodeTitle),
		SeasonNumber:    ogen.NewOptInt(int(item.LastSeasonNumber)),
		EpisodeNumber:   ogen.NewOptInt(int(item.LastEpisodeNumber)),
		ProgressSeconds: ogen.NewOptInt(int(item.ProgressSeconds)),
		DurationSeconds: ogen.NewOptInt(int(item.DurationSeconds)),
		LastWatchedAt:   ogen.NewOptDateTime(item.LastWatchedAt),
	}

	if item.Series != nil {
		o.SeriesID.SetTo(item.Series.ID)
		o.SeriesTitle.SetTo(item.Series.Title)
		setOpt(&o.SeriesPosterPath, item.Series.PosterPath)
	}

	// Calculate progress percent
	if item.DurationSeconds > 0 {
		percent := float64(item.ProgressSeconds) / float64(item.DurationSeconds) * 100
		o.ProgressPercent.SetTo(float32(percent))
	}

	return o
}

// episodeWithSeriesInfoToOgen converts an episode with series info to ogen EpisodeWithSeriesInfo.
func episodeWithSeriesInfoToOgen(e *tvshow.EpisodeWithSeriesInfo) *ogen.EpisodeWithSeriesInfo {
	o := &ogen.EpisodeWithSeriesInfo{
		ID:            ogen.NewOptUUID(e.ID),
		SeriesID:      ogen.NewOptUUID(e.SeriesID),
		SeriesTitle:   ogen.NewOptString(e.SeriesTitle),
		SeasonID:      ogen.NewOptUUID(e.SeasonID),
		SeasonNumber:  ogen.NewOptInt(int(e.SeasonNumber)),
		EpisodeNumber: ogen.NewOptInt(int(e.EpisodeNumber)),
		Title:         ogen.NewOptString(e.Title),
		CreatedAt:     ogen.NewOptDateTime(e.CreatedAt),
	}

	setOpt(&o.SeriesPosterPath, e.SeriesPosterPath)
	setOptConv(&o.TmdbID, e.TMDbID, int32ToInt)
	setOptConv(&o.TvdbID, e.TVDbID, int32ToInt)
	setOpt(&o.Overview, e.Overview)
	setOpt(&o.AirDate, e.AirDate)
	setOptConv(&o.Runtime, e.Runtime, int32ToInt)
	setOpt(&o.StillPath, e.StillPath)

	return o
}

// seriesWatchStatsToOgen converts series watch stats to ogen SeriesWatchStats.
func seriesWatchStatsToOgen(stats *tvshow.SeriesWatchStats, seriesID string) *ogen.SeriesWatchStats {
	// Parse series ID if provided
	o := &ogen.SeriesWatchStats{
		TotalEpisodes:   ogen.NewOptInt(int(stats.TotalEpisodes)),
		WatchedEpisodes: ogen.NewOptInt(int(stats.WatchedCount)),
	}

	// Calculate completion percent
	if stats.TotalEpisodes > 0 {
		percent := float64(stats.WatchedCount) / float64(stats.TotalEpisodes) * 100
		o.CompletionPercent.SetTo(float32(percent))
	}

	return o
}

// userTVStatsToOgen converts user TV stats to ogen UserTVStats.
func userTVStatsToOgen(stats *tvshow.UserTVStats, userID string) *ogen.UserTVStats {
	return &ogen.UserTVStats{
		TotalSeries:          ogen.NewOptInt(int(stats.SeriesCount)),
		TotalEpisodesWatched: ogen.NewOptInt(int(stats.EpisodesWatched)),
		CurrentSeries:        ogen.NewOptInt(int(stats.EpisodesInProgress)),
	}
}
