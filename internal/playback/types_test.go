package playback

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItoa(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want string
	}{
		{name: "zero", in: 0, want: "0"},
		{name: "single digit", in: 7, want: "7"},
		{name: "two digits", in: 42, want: "42"},
		{name: "three digits", in: 123, want: "123"},
		{name: "large number", in: 99999, want: "99999"},
		{name: "one", in: 1, want: "1"},
		{name: "ten", in: 10, want: "10"},
		{name: "hundred", in: 100, want: "100"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := itoa(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsBitmapSubtitle(t *testing.T) {
	tests := []struct {
		codec    string
		isBitmap bool
	}{
		// Bitmap formats
		{"hdmv_pgs_subtitle", true},
		{"dvd_subtitle", true},
		{"dvb_subtitle", true},
		{"pgssub", true},
		{"vobsub", true},

		// Text formats — should NOT be bitmap
		{"subrip", false},
		{"srt", false},
		{"ass", false},
		{"ssa", false},
		{"webvtt", false},
		{"mov_text", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.codec, func(t *testing.T) {
			assert.Equal(t, tc.isBitmap, isBitmapSubtitle(tc.codec))
		})
	}
}

func TestSubtitleURL(t *testing.T) {
	sessionID := uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef")

	tests := []struct {
		name       string
		trackIndex int
		want       string
	}{
		{
			name:       "track 0",
			trackIndex: 0,
			want:       "/api/v1/playback/stream/01234567-89ab-cdef-0123-456789abcdef/subs/0.vtt",
		},
		{
			name:       "track 3",
			trackIndex: 3,
			want:       "/api/v1/playback/stream/01234567-89ab-cdef-0123-456789abcdef/subs/3.vtt",
		},
		{
			name:       "track 12",
			trackIndex: 12,
			want:       "/api/v1/playback/stream/01234567-89ab-cdef-0123-456789abcdef/subs/12.vtt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := subtitleURL(sessionID, tc.trackIndex)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAudioTracksFromMediaInfo(t *testing.T) {
	t.Run("empty audio streams", func(t *testing.T) {
		info := &movie.MediaInfo{
			AudioStreams: []movie.AudioStreamInfo{},
		}
		tracks := AudioTracksFromMediaInfo(info)
		assert.Empty(t, tracks)
	})

	t.Run("single audio stream", func(t *testing.T) {
		info := &movie.MediaInfo{
			AudioStreams: []movie.AudioStreamInfo{
				{
					Index:     1,
					Codec:     "aac",
					Language:  "eng",
					Title:     "Stereo",
					Channels:  2,
					Layout:    "stereo",
					IsDefault: true,
				},
			},
		}
		tracks := AudioTracksFromMediaInfo(info)
		require.Len(t, tracks, 1)
		assert.Equal(t, 1, tracks[0].Index)
		assert.Equal(t, "aac", tracks[0].Codec)
		assert.Equal(t, "eng", tracks[0].Language)
		assert.Equal(t, "Stereo", tracks[0].Title)
		assert.Equal(t, 2, tracks[0].Channels)
		assert.Equal(t, "stereo", tracks[0].Layout)
		assert.True(t, tracks[0].IsDefault)
	})

	t.Run("multiple audio streams", func(t *testing.T) {
		info := &movie.MediaInfo{
			AudioStreams: []movie.AudioStreamInfo{
				{
					Index:     1,
					Codec:     "dts",
					Language:  "eng",
					Title:     "DTS-HD MA 7.1",
					Channels:  8,
					Layout:    "7.1",
					IsDefault: true,
				},
				{
					Index:     2,
					Codec:     "ac3",
					Language:  "eng",
					Title:     "Dolby Digital 5.1",
					Channels:  6,
					Layout:    "5.1",
					IsDefault: false,
				},
				{
					Index:     3,
					Codec:     "aac",
					Language:  "ger",
					Title:     "German Stereo",
					Channels:  2,
					Layout:    "stereo",
					IsDefault: false,
				},
			},
		}
		tracks := AudioTracksFromMediaInfo(info)
		require.Len(t, tracks, 3)

		// Verify ordering is preserved
		assert.Equal(t, 1, tracks[0].Index)
		assert.Equal(t, "dts", tracks[0].Codec)
		assert.True(t, tracks[0].IsDefault)

		assert.Equal(t, 2, tracks[1].Index)
		assert.Equal(t, "ac3", tracks[1].Codec)
		assert.False(t, tracks[1].IsDefault)

		assert.Equal(t, 3, tracks[2].Index)
		assert.Equal(t, "ger", tracks[2].Language)
	})
}

func TestSubtitleTracksFromMediaInfo(t *testing.T) {
	sessionID := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")

	t.Run("empty subtitle streams", func(t *testing.T) {
		info := &movie.MediaInfo{
			SubtitleStreams: []movie.SubtitleStreamInfo{},
		}
		tracks := SubtitleTracksFromMediaInfo(info, sessionID)
		assert.Empty(t, tracks)
	})

	t.Run("text subtitles only", func(t *testing.T) {
		info := &movie.MediaInfo{
			SubtitleStreams: []movie.SubtitleStreamInfo{
				{
					Index:    4,
					Codec:    "subrip",
					Language: "eng",
					Title:    "English",
					IsForced: false,
				},
				{
					Index:    5,
					Codec:    "ass",
					Language: "ger",
					Title:    "German",
					IsForced: false,
				},
			},
		}
		tracks := SubtitleTracksFromMediaInfo(info, sessionID)
		require.Len(t, tracks, 2)

		assert.Equal(t, 4, tracks[0].Index)
		assert.Equal(t, "subrip", tracks[0].Codec)
		assert.Equal(t, "eng", tracks[0].Language)
		assert.Equal(t, "English", tracks[0].Title)
		assert.False(t, tracks[0].IsForced)
		assert.Contains(t, tracks[0].URL, sessionID.String())
		assert.Contains(t, tracks[0].URL, "/subs/4.vtt")

		assert.Equal(t, 5, tracks[1].Index)
		assert.Equal(t, "ger", tracks[1].Language)
	})

	t.Run("bitmap subtitles are filtered out", func(t *testing.T) {
		info := &movie.MediaInfo{
			SubtitleStreams: []movie.SubtitleStreamInfo{
				{
					Index:    4,
					Codec:    "subrip",
					Language: "eng",
					Title:    "English",
				},
				{
					Index:    5,
					Codec:    "hdmv_pgs_subtitle",
					Language: "eng",
					Title:    "English PGS",
				},
				{
					Index:    6,
					Codec:    "dvd_subtitle",
					Language: "ger",
					Title:    "German VobSub",
				},
				{
					Index:    7,
					Codec:    "ass",
					Language: "jpn",
					Title:    "Japanese ASS",
				},
			},
		}
		tracks := SubtitleTracksFromMediaInfo(info, sessionID)
		require.Len(t, tracks, 2, "only text subtitles should be included")

		assert.Equal(t, "subrip", tracks[0].Codec)
		assert.Equal(t, "ass", tracks[1].Codec)
	})

	t.Run("all bitmap subtitles produces empty result", func(t *testing.T) {
		info := &movie.MediaInfo{
			SubtitleStreams: []movie.SubtitleStreamInfo{
				{Index: 4, Codec: "hdmv_pgs_subtitle", Language: "eng"},
				{Index: 5, Codec: "dvd_subtitle", Language: "ger"},
				{Index: 6, Codec: "dvb_subtitle", Language: "fra"},
				{Index: 7, Codec: "pgssub", Language: "jpn"},
				{Index: 8, Codec: "vobsub", Language: "spa"},
			},
		}
		tracks := SubtitleTracksFromMediaInfo(info, sessionID)
		assert.Empty(t, tracks)
	})

	t.Run("forced subtitle flag preserved", func(t *testing.T) {
		info := &movie.MediaInfo{
			SubtitleStreams: []movie.SubtitleStreamInfo{
				{
					Index:    4,
					Codec:    "subrip",
					Language: "eng",
					Title:    "English (Forced)",
					IsForced: true,
				},
			},
		}
		tracks := SubtitleTracksFromMediaInfo(info, sessionID)
		require.Len(t, tracks, 1)
		assert.True(t, tracks[0].IsForced)
	})
}
