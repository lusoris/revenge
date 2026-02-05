package library

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanSummary(t *testing.T) {
	t.Run("AddError", func(t *testing.T) {
		summary := &ScanSummary{}

		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		summary.AddError(err1)
		summary.AddError(err2)

		assert.Len(t, summary.Errors, 2)
		assert.Equal(t, err1, summary.Errors[0])
		assert.Equal(t, err2, summary.Errors[1])
	})

	t.Run("HasErrors", func(t *testing.T) {
		summary := &ScanSummary{}
		assert.False(t, summary.HasErrors())

		summary.AddError(errors.New("test"))
		assert.True(t, summary.HasErrors())
	})
}

func TestMatchResult(t *testing.T) {
	type Movie struct {
		ID    string
		Title string
	}

	t.Run("IsMatched with content", func(t *testing.T) {
		movie := &Movie{ID: "1", Title: "Test"}
		result := MatchResult[Movie]{
			FilePath:   "/path/to/movie.mkv",
			Content:    movie,
			MatchType:  MatchTypeTitle,
			Confidence: 0.95,
		}

		assert.True(t, result.IsMatched())
	})

	t.Run("IsMatched without content", func(t *testing.T) {
		result := MatchResult[Movie]{
			FilePath:  "/path/to/movie.mkv",
			Content:   nil,
			MatchType: MatchTypeUnmatched,
		}

		assert.False(t, result.IsMatched())
	})

	t.Run("IsMatched with error", func(t *testing.T) {
		movie := &Movie{ID: "1", Title: "Test"}
		result := MatchResult[Movie]{
			FilePath:  "/path/to/movie.mkv",
			Content:   movie,
			MatchType: MatchTypeTitle,
			Error:     errors.New("some error"),
		}

		assert.False(t, result.IsMatched())
	})
}

func TestScanItem(t *testing.T) {
	t.Run("GetYear with year in metadata", func(t *testing.T) {
		item := ScanItem{
			FilePath:    "/path/to/movie.mkv",
			ParsedTitle: "Test Movie",
			Metadata:    map[string]any{"year": 1999},
		}

		year := item.GetYear()
		assert.NotNil(t, year)
		assert.Equal(t, 1999, *year)
	})

	t.Run("GetYear without metadata", func(t *testing.T) {
		item := ScanItem{
			FilePath:    "/path/to/movie.mkv",
			ParsedTitle: "Test Movie",
			Metadata:    nil,
		}

		year := item.GetYear()
		assert.Nil(t, year)
	})

	t.Run("GetYear with no year in metadata", func(t *testing.T) {
		item := ScanItem{
			FilePath:    "/path/to/movie.mkv",
			ParsedTitle: "Test Movie",
			Metadata:    map[string]any{"title": "Test"},
		}

		year := item.GetYear()
		assert.Nil(t, year)
	})

	t.Run("GetYear with wrong type", func(t *testing.T) {
		item := ScanItem{
			FilePath:    "/path/to/movie.mkv",
			ParsedTitle: "Test Movie",
			Metadata:    map[string]any{"year": "1999"}, // String instead of int
		}

		year := item.GetYear()
		assert.Nil(t, year)
	})
}

func TestMatchTypeConstants(t *testing.T) {
	assert.Equal(t, MatchType("exact"), MatchTypeExact)
	assert.Equal(t, MatchType("title"), MatchTypeTitle)
	assert.Equal(t, MatchType("fuzzy"), MatchTypeFuzzy)
	assert.Equal(t, MatchType("manual"), MatchTypeManual)
	assert.Equal(t, MatchType("unmatched"), MatchTypeUnmatched)
}

func TestDefaultRefreshOptions(t *testing.T) {
	opts := DefaultRefreshOptions()

	assert.True(t, opts.RefreshCredits)
	assert.True(t, opts.RefreshGenres)
	assert.True(t, opts.RefreshImages)
	assert.Equal(t, []string{"en-US"}, opts.Languages)
}

func TestMediaFileInfo(t *testing.T) {
	info := MediaFileInfo{
		Path:            "/path/to/movie.mkv",
		Size:            1024 * 1024 * 1000, // ~1GB
		Container:       "mkv",
		Resolution:      "1920x1080",
		ResolutionLabel: "1080p",
		VideoCodec:      "hevc",
		VideoProfile:    "Main 10",
		AudioCodec:      "dts",
		BitrateKbps:     8000,
		DurationSeconds: 7200,
		Framerate:       23.976,
		DynamicRange:    "HDR10",
		ColorSpace:      "bt2020",
		AudioChannels:   6,
		AudioLayout:     "5.1",
		Languages:       []string{"en", "de"},
		SubtitleLangs:   []string{"en", "de", "fr"},
	}

	// Just verify the struct can be populated correctly
	assert.Equal(t, "/path/to/movie.mkv", info.Path)
	assert.Equal(t, "mkv", info.Container)
	assert.Equal(t, "1080p", info.ResolutionLabel)
	assert.Equal(t, "hevc", info.VideoCodec)
	assert.Equal(t, "HDR10", info.DynamicRange)
	assert.Equal(t, 6, info.AudioChannels)
	assert.Len(t, info.Languages, 2)
	assert.Len(t, info.SubtitleLangs, 3)
}
