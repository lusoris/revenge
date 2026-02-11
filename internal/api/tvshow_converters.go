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

	if s.TMDbID != nil {
		o.TmdbID.SetTo(int(*s.TMDbID))
	}
	if s.TVDbID != nil {
		o.TvdbID.SetTo(int(*s.TVDbID))
	}
	if s.IMDbID != nil {
		o.ImdbID.SetTo(*s.IMDbID)
	}
	if s.SonarrID != nil {
		o.SonarrID.SetTo(int(*s.SonarrID))
	}
	if s.OriginalTitle != nil {
		o.OriginalTitle.SetTo(*s.OriginalTitle)
	}
	if s.OriginalLanguage != "" {
		o.OriginalLanguage.SetTo(s.OriginalLanguage)
	}
	if s.Tagline != nil {
		o.Tagline.SetTo(*s.Tagline)
	}
	if s.Overview != nil {
		o.Overview.SetTo(*s.Overview)
	}
	if s.Status != nil {
		o.Status.SetTo(*s.Status)
	}
	if s.Type != nil {
		o.Type.SetTo(*s.Type)
	}
	if s.FirstAirDate != nil {
		o.FirstAirDate.SetTo(*s.FirstAirDate)
	}
	if s.LastAirDate != nil {
		o.LastAirDate.SetTo(*s.LastAirDate)
	}
	if s.VoteAverage != nil {
		f, _ := s.VoteAverage.Float64()
		o.VoteAverage.SetTo(float32(f))
	}
	if s.VoteCount != nil {
		o.VoteCount.SetTo(int(*s.VoteCount))
	}
	if s.Popularity != nil {
		f, _ := s.Popularity.Float64()
		o.Popularity.SetTo(float32(f))
	}
	if s.PosterPath != nil {
		o.PosterPath.SetTo(*s.PosterPath)
	}
	if s.BackdropPath != nil {
		o.BackdropPath.SetTo(*s.BackdropPath)
	}
	if s.TrailerURL != nil {
		o.TrailerURL.SetTo(*s.TrailerURL)
	}
	if s.Homepage != nil {
		o.Homepage.SetTo(*s.Homepage)
	}
	if s.MetadataUpdatedAt != nil {
		o.MetadataUpdatedAt.SetTo(*s.MetadataUpdatedAt)
	}

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

	if s.TMDbID != nil {
		o.TmdbID.SetTo(int(*s.TMDbID))
	}
	if s.Overview != nil {
		o.Overview.SetTo(*s.Overview)
	}
	if s.PosterPath != nil {
		o.PosterPath.SetTo(*s.PosterPath)
	}
	if s.AirDate != nil {
		o.AirDate.SetTo(*s.AirDate)
	}
	if s.VoteAverage != nil {
		f, _ := s.VoteAverage.Float64()
		o.VoteAverage.SetTo(float32(f))
	}

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

	if e.TMDbID != nil {
		o.TmdbID.SetTo(int(*e.TMDbID))
	}
	if e.TVDbID != nil {
		o.TvdbID.SetTo(int(*e.TVDbID))
	}
	if e.IMDbID != nil {
		o.ImdbID.SetTo(*e.IMDbID)
	}
	if e.Overview != nil {
		o.Overview.SetTo(*e.Overview)
	}
	if e.AirDate != nil {
		o.AirDate.SetTo(*e.AirDate)
	}
	if e.Runtime != nil {
		o.Runtime.SetTo(int(*e.Runtime))
	}
	if e.VoteAverage != nil {
		f, _ := e.VoteAverage.Float64()
		o.VoteAverage.SetTo(float32(f))
	}
	if e.VoteCount != nil {
		o.VoteCount.SetTo(int(*e.VoteCount))
	}
	if e.StillPath != nil {
		o.StillPath.SetTo(*e.StillPath)
	}
	if e.ProductionCode != nil {
		o.ProductionCode.SetTo(*e.ProductionCode)
	}

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

	if f.Container != nil {
		o.Container.SetTo(*f.Container)
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
	if f.BitrateKbps != nil {
		o.BitrateKbps.SetTo(int(*f.BitrateKbps))
	}
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
	if f.SonarrFileID != nil {
		o.SonarrFileID.SetTo(int(*f.SonarrFileID))
	}

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
		TmdbID:    ogen.NewOptInt(int(n.TMDbID)),
		Name:      ogen.NewOptString(n.Name),
		CreatedAt: ogen.NewOptDateTime(n.CreatedAt),
	}

	if n.LogoPath != nil {
		o.LogoPath.SetTo(*n.LogoPath)
	}
	if n.OriginCountry != nil {
		o.OriginCountry.SetTo(*n.OriginCountry)
	}

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
		if item.Series.PosterPath != nil {
			o.SeriesPosterPath.SetTo(*item.Series.PosterPath)
		}
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

	if e.SeriesPosterPath != nil {
		o.SeriesPosterPath.SetTo(*e.SeriesPosterPath)
	}
	if e.TMDbID != nil {
		o.TmdbID.SetTo(int(*e.TMDbID))
	}
	if e.TVDbID != nil {
		o.TvdbID.SetTo(int(*e.TVDbID))
	}
	if e.Overview != nil {
		o.Overview.SetTo(*e.Overview)
	}
	if e.AirDate != nil {
		o.AirDate.SetTo(*e.AirDate)
	}
	if e.Runtime != nil {
		o.Runtime.SetTo(int(*e.Runtime))
	}
	if e.StillPath != nil {
		o.StillPath.SetTo(*e.StillPath)
	}

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
