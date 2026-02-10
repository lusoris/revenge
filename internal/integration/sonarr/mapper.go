package sonarr

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/util"
	"github.com/lusoris/revenge/internal/util/ptr"
)

// Mapper converts Sonarr types to domain types.
type Mapper struct{}

// NewMapper creates a new Sonarr mapper.
func NewMapper() *Mapper {
	return &Mapper{}
}

// ToSeries converts a Sonarr series to a domain series.
func (m *Mapper) ToSeries(ss *Series) *tvshow.Series {
	result := &tvshow.Series{
		ID:               uuid.Must(uuid.NewV7()),
		TVDbID:           ptr.To(util.SafeIntToInt32(ss.TVDbID)),
		IMDbID:           ptrString(ss.IMDbID),
		SonarrID:         ptr.To(util.SafeIntToInt32(ss.ID)),
		Title:            ss.Title,
		OriginalLanguage: ss.OriginalLanguage.Name,
		Overview:         ptrString(ss.Overview),
		Status:           ptrString(ss.Status),
		FirstAirDate:     ss.FirstAired,
		LastAirDate:      ss.LastAired,
		TotalSeasons:     util.SafeIntToInt32(len(ss.Seasons)),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Calculate total episodes from statistics
	if ss.Statistics != nil {
		result.TotalEpisodes = util.SafeIntToInt32(ss.Statistics.TotalEpisodeCount)
	}

	// Set ratings
	if ss.Ratings.Value > 0 {
		d, _ := decimal.NewFromFloat64(ss.Ratings.Value)
		result.VoteAverage = &d
		result.VoteCount = ptr.To(util.SafeIntToInt32(ss.Ratings.Votes))
	}

	// Set images from Sonarr
	for _, img := range ss.Images {
		switch img.CoverType {
		case "poster":
			if img.RemoteURL != "" {
				result.PosterPath = ptrString(img.RemoteURL)
			} else if img.URL != "" {
				result.PosterPath = ptrString(img.URL)
			}
		case "fanart":
			if img.RemoteURL != "" {
				result.BackdropPath = ptrString(img.RemoteURL)
			} else if img.URL != "" {
				result.BackdropPath = ptrString(img.URL)
			}
		}
	}

	// Series type
	if ss.SeriesType != "" {
		result.Type = ptrString(mapSeriesType(ss.SeriesType))
	}

	return result
}

// ToSeason converts a Sonarr season to a domain season.
func (m *Mapper) ToSeason(si *SeasonInfo, seriesID uuid.UUID) *tvshow.Season {
	season := &tvshow.Season{
		ID:           uuid.Must(uuid.NewV7()),
		SeriesID:     seriesID,
		SeasonNumber: util.SafeIntToInt32(si.SeasonNumber),
		Name:         fmt.Sprintf("Season %d", si.SeasonNumber),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if si.SeasonNumber == 0 {
		season.Name = "Specials"
	}

	if si.Statistics != nil {
		season.EpisodeCount = util.SafeIntToInt32(si.Statistics.EpisodeCount)
	}

	return season
}

// ToEpisode converts a Sonarr episode to a domain episode.
func (m *Mapper) ToEpisode(se *Episode, seriesID, seasonID uuid.UUID) *tvshow.Episode {
	ep := &tvshow.Episode{
		ID:            uuid.Must(uuid.NewV7()),
		SeriesID:      seriesID,
		SeasonID:      seasonID,
		TVDbID:        ptr.To(util.SafeIntToInt32(se.TVDbID)),
		SeasonNumber:  util.SafeIntToInt32(se.SeasonNumber),
		EpisodeNumber: util.SafeIntToInt32(se.EpisodeNumber),
		Title:         se.Title,
		Overview:      ptrString(se.Overview),
		AirDate:       se.AirDateUtc,
		Runtime:       ptr.To(util.SafeIntToInt32(se.Runtime)),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set still image
	for _, img := range se.Images {
		if img.CoverType == "screenshot" {
			if img.RemoteURL != "" {
				ep.StillPath = ptrString(img.RemoteURL)
			} else if img.URL != "" {
				ep.StillPath = ptrString(img.URL)
			}
			break
		}
	}

	return ep
}

// ToEpisodeFile converts a Sonarr episode file to a domain episode file.
func (m *Mapper) ToEpisodeFile(sef *EpisodeFile, episodeID uuid.UUID) *tvshow.EpisodeFile {
	result := &tvshow.EpisodeFile{
		ID:           uuid.Must(uuid.NewV7()),
		EpisodeID:    episodeID,
		FilePath:     sef.Path,
		FileName:     sef.RelativePath,
		FileSize:     sef.Size,
		SonarrFileID: ptr.To(util.SafeIntToInt32(sef.ID)),
		CreatedAt:    sef.DateAdded,
		UpdatedAt:    time.Now(),
	}

	// Quality information
	if sef.Quality.Quality.Name != "" {
		result.QualityProfile = ptrString(sef.Quality.Quality.Name)
	}
	if sef.Quality.Quality.Resolution > 0 {
		result.Resolution = ptrString(resolutionToString(sef.Quality.Quality.Resolution))
	}

	// Container from path
	result.Container = extractContainer(sef.Path)

	// Media info
	if sef.MediaInfo != nil {
		mi := sef.MediaInfo
		result.VideoCodec = ptrString(mi.VideoCodec)
		result.AudioCodec = ptrString(mi.AudioCodec)
		result.BitrateKbps = ptr.To(util.SafeIntToInt32(mi.VideoBitrate / 1000))

		// Parse runtime duration
		if mi.RunTime != "" {
			if duration, err := time.ParseDuration(mi.RunTime); err == nil {
				d, _ := decimal.NewFromFloat64(duration.Seconds())
				result.DurationSeconds = &d
			}
		}

		// Audio languages
		if mi.AudioLanguages != "" {
			result.AudioLanguages = []string{mi.AudioLanguages}
		}

		// Subtitles
		if mi.Subtitles != "" {
			result.SubtitleLanguages = []string{mi.Subtitles}
		}
	}

	return result
}

// ToGenres extracts genres from a Sonarr series.
func (m *Mapper) ToGenres(ss *Series, seriesID uuid.UUID) []tvshow.SeriesGenre {
	genres := make([]tvshow.SeriesGenre, len(ss.Genres))
	for i, g := range ss.Genres {
		genres[i] = tvshow.SeriesGenre{
			ID:        uuid.Must(uuid.NewV7()),
			SeriesID:  seriesID,
			Name:      g,
			CreatedAt: time.Now(),
		}
	}
	return genres
}

// Helper functions

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func resolutionToString(res int) string {
	switch {
	case res >= 2160:
		return "4K"
	case res >= 1080:
		return "1080p"
	case res >= 720:
		return "720p"
	case res >= 480:
		return "480p"
	default:
		return "SD"
	}
}

func mapSeriesType(sonarrType string) string {
	switch sonarrType {
	case SeriesTypeStandard:
		return "Scripted"
	case SeriesTypeDaily:
		return "Talk Show"
	case SeriesTypeAnime:
		return "Animation"
	default:
		return "Scripted"
	}
}

func extractContainer(path string) *string {
	if len(path) < 4 {
		return nil
	}
	// Get file extension
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext := path[i+1:]
			return &ext
		}
	}
	return nil
}
