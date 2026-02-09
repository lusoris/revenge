package radarr

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/util"
	"github.com/lusoris/revenge/internal/util/ptr"
)

// Mapper converts Radarr types to domain types.
type Mapper struct{}

// NewMapper creates a new Radarr mapper.
func NewMapper() *Mapper {
	return &Mapper{}
}

// ToMovie converts a Radarr movie to a domain movie.
func (m *Mapper) ToMovie(rm *Movie) *movie.Movie {
	result := &movie.Movie{
		ID:               uuid.Must(uuid.NewV7()), // Generate new ID for local storage
		TMDbID:           ptr.To(util.SafeIntToInt32(rm.TMDbID)),
		IMDbID:           ptrString(rm.IMDbID),
		Title:            rm.Title,
		OriginalTitle:    ptrString(rm.OriginalTitle),
		Year:             ptr.To(util.SafeIntToInt32(rm.Year)),
		Runtime:          ptr.To(util.SafeIntToInt32(rm.Runtime)),
		Overview:         ptrString(rm.Overview),
		Status:           ptrString(rm.Status),
		OriginalLanguage: ptrString(rm.OriginalLanguage.Name),
		RadarrID:         ptr.To(util.SafeIntToInt32(rm.ID)),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Set release date from available sources
	if rm.DigitalRelease != nil {
		result.ReleaseDate = rm.DigitalRelease
	} else if rm.PhysicalRelease != nil {
		result.ReleaseDate = rm.PhysicalRelease
	} else if rm.InCinemas != nil {
		result.ReleaseDate = rm.InCinemas
	}

	// Set ratings from TMDb (preferred) or IMDb
	if rm.Ratings.TMDb != nil {
		d, _ := decimal.NewFromFloat64(rm.Ratings.TMDb.Value)
		result.VoteAverage = &d
		result.VoteCount = ptr.To(util.SafeIntToInt32(rm.Ratings.TMDb.Votes))
	} else if rm.Ratings.IMDb != nil {
		d, _ := decimal.NewFromFloat64(rm.Ratings.IMDb.Value)
		result.VoteAverage = &d
		result.VoteCount = ptr.To(util.SafeIntToInt32(rm.Ratings.IMDb.Votes))
	}

	// Set popularity
	if rm.Popularity > 0 {
		d, _ := decimal.NewFromFloat64(rm.Popularity)
		result.Popularity = &d
	}

	// Set library added time
	result.LibraryAddedAt = rm.Added

	// Set images from Radarr
	for _, img := range rm.Images {
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

	// YouTube trailer
	if rm.YouTubeTrailerID != "" {
		result.TrailerURL = ptrString("https://www.youtube.com/watch?v=" + rm.YouTubeTrailerID)
	}

	return result
}

// ToMovieFile converts a Radarr movie file to a domain movie file.
func (m *Mapper) ToMovieFile(rmf *MovieFile, movieID uuid.UUID) *movie.MovieFile {
	result := &movie.MovieFile{
		ID:           uuid.Must(uuid.NewV7()),
		MovieID:      movieID,
		FilePath:     rmf.Path,
		FileName:     rmf.RelativePath,
		FileSize:     rmf.Size,
		RadarrFileID: ptr.To(util.SafeIntToInt32(rmf.ID)),
		CreatedAt:    rmf.DateAdded,
		UpdatedAt:    time.Now(),
	}

	// Quality information
	if rmf.Quality.Quality.Name != "" {
		result.QualityProfile = ptrString(rmf.Quality.Quality.Name)
	}
	if rmf.Quality.Quality.Resolution > 0 {
		result.Resolution = ptrString(resolutionToString(rmf.Quality.Quality.Resolution))
	}

	// Media info
	if rmf.MediaInfo != nil {
		mi := rmf.MediaInfo
		result.VideoCodec = ptrString(mi.VideoCodec)
		result.AudioCodec = ptrString(mi.AudioCodec)
		result.DynamicRange = ptrString(mi.VideoDynamicRange)
		if mi.VideoFps > 0 {
			fps, _ := decimal.NewFromFloat64(mi.VideoFps)
			result.Framerate = &fps
		}
		result.BitrateKbps = ptr.To(util.SafeIntToInt32(mi.VideoBitrate / 1000))
		if mi.AudioChannels > 0 {
			result.AudioChannels = ptrString(formatAudioChannels(mi.AudioChannels))
		}

		// Parse runtime duration
		if mi.RunTime != "" {
			// RunTime is in format "HH:MM:SS" or similar
			if duration, err := time.ParseDuration(mi.RunTime); err == nil {
				result.DurationSeconds = ptr.To(util.SafeInt64ToInt32(int64(duration.Seconds())))
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

	// Release group - append to quality profile
	if rmf.ReleaseGroup != "" && result.QualityProfile != nil {
		q := *result.QualityProfile + " (" + rmf.ReleaseGroup + ")"
		result.QualityProfile = &q
	}

	return result
}

// ToMovieCollection converts a Radarr collection to a domain collection.
func (m *Mapper) ToMovieCollection(rc *Collection) *movie.MovieCollection {
	if rc == nil {
		return nil
	}

	result := &movie.MovieCollection{
		ID:               uuid.Must(uuid.NewV7()),
		TMDbCollectionID: ptr.To(util.SafeIntToInt32(rc.TMDbID)),
		Name:             rc.Name,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Set images
	for _, img := range rc.Images {
		switch img.CoverType {
		case "poster":
			if img.RemoteURL != "" {
				result.PosterPath = ptrString(img.RemoteURL)
			}
		case "fanart":
			if img.RemoteURL != "" {
				result.BackdropPath = ptrString(img.RemoteURL)
			}
		}
	}

	return result
}

// ToGenres extracts genres from a Radarr movie.
func (m *Mapper) ToGenres(rm *Movie, movieID uuid.UUID) []movie.MovieGenre {
	genres := make([]movie.MovieGenre, len(rm.Genres))
	for i, g := range rm.Genres {
		genres[i] = movie.MovieGenre{
			ID:      uuid.Must(uuid.NewV7()),
			MovieID: movieID,
			Name:    g,
			// TMDbGenreID would need to be looked up separately
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

func formatAudioChannels(channels float64) string {
	// Format audio channels like "5.1", "7.1", "2.0"
	if channels == float64(int(channels)) {
		return fmt.Sprintf("%.0f.0", channels)
	}
	return fmt.Sprintf("%.1f", channels)
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
