package hls

import (
	"testing"

	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bwCalc computes expected bandwidth the same way the production code does,
// avoiding Go constant-expression issues with float64-to-int.
func bwCalc(videoBps, audioBps int) int {
	return int(float64(videoBps+audioBps) * 1.1)
}

func TestProfileVariantsFromDecision(t *testing.T) {
	t.Run("empty profiles", func(t *testing.T) {
		variants := ProfileVariantsFromDecision(nil, 0, 0)
		assert.Empty(t, variants)
	})

	t.Run("single profile with explicit bitrate", func(t *testing.T) {
		profiles := []transcode.ProfileDecision{
			{Name: "720p", Width: 1280, Height: 720, VideoBitrate: 2800, AudioBitrate: 128},
		}
		variants := ProfileVariantsFromDecision(profiles, 0, 0)
		require.Len(t, variants, 1)
		assert.Equal(t, "720p", variants[0].Name)
		assert.Equal(t, 1280, variants[0].Width)
		assert.Equal(t, 720, variants[0].Height)
		expected := bwCalc(2800*1000, 128*1000)
		assert.Equal(t, expected, variants[0].Bandwidth)
	})

	t.Run("original profile uses source bitrates", func(t *testing.T) {
		profiles := []transcode.ProfileDecision{
			{Name: "original", Width: 1920, Height: 1080, VideoBitrate: 0, AudioBitrate: 0},
		}
		variants := ProfileVariantsFromDecision(profiles, 8000, 320)
		require.Len(t, variants, 1)
		expected := bwCalc(8000*1000, 320*1000)
		assert.Equal(t, expected, variants[0].Bandwidth)
	})

	t.Run("multiple profiles preserve order", func(t *testing.T) {
		profiles := []transcode.ProfileDecision{
			{Name: "original", Width: 1920, Height: 1080},
			{Name: "720p", Width: 1280, Height: 720, VideoBitrate: 2800, AudioBitrate: 128},
			{Name: "480p", Width: 854, Height: 480, VideoBitrate: 1400, AudioBitrate: 96},
		}
		variants := ProfileVariantsFromDecision(profiles, 5000, 192)
		require.Len(t, variants, 3)
		assert.Equal(t, "original", variants[0].Name)
		assert.Equal(t, "720p", variants[1].Name)
		assert.Equal(t, "480p", variants[2].Name)
	})
}

func TestEstimateBandwidth_AllBranches(t *testing.T) {
	t.Run("both bitrates from profile", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 5000, AudioBitrate: 192}
		bw := estimateBandwidth(pd, 0, 0)
		expected := bwCalc(5000*1000, 192*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("video from source audio from profile", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 0, AudioBitrate: 128}
		bw := estimateBandwidth(pd, 3000, 0)
		expected := bwCalc(3000*1000, 128*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("video from profile audio from source", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 2800, AudioBitrate: 0}
		bw := estimateBandwidth(pd, 0, 256)
		expected := bwCalc(2800*1000, 256*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("both from source", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 0, AudioBitrate: 0}
		bw := estimateBandwidth(pd, 8000, 320)
		expected := bwCalc(8000*1000, 320*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("both zero uses defaults", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 0, AudioBitrate: 0}
		bw := estimateBandwidth(pd, 0, 0)
		// Default: 2_000_000 video (Height=0 fallback) + 192000 audio
		expected := bwCalc(2_000_000, 192000)
		assert.Equal(t, expected, bw)
	})

	t.Run("video zero no source uses default video", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 0, AudioBitrate: 128}
		bw := estimateBandwidth(pd, 0, 0)
		// Default video 2_000_000 (Height=0) + 128*1000
		expected := bwCalc(2_000_000, 128*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("4K copy profile uses 40 Mbps default", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 0, AudioBitrate: 192, Height: 2160}
		bw := estimateBandwidth(pd, 0, 0)
		expected := bwCalc(40_000_000, 192*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("1080p copy profile uses 10 Mbps default", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 0, AudioBitrate: 192, Height: 1080}
		bw := estimateBandwidth(pd, 0, 0)
		expected := bwCalc(10_000_000, 192*1000)
		assert.Equal(t, expected, bw)
	})

	t.Run("audio zero no source uses default audio", func(t *testing.T) {
		pd := transcode.ProfileDecision{VideoBitrate: 2800, AudioBitrate: 0}
		bw := estimateBandwidth(pd, 0, 0)
		// 2800*1000 + default audio 192000
		expected := bwCalc(2800*1000, 192000)
		assert.Equal(t, expected, bw)
	})
}

func TestAudioRenditionSegmentPath(t *testing.T) {
	tests := []struct {
		name        string
		segmentDir  string
		trackIndex  int
		segmentFile string
		want        string
	}{
		{
			name:        "track 0 first segment",
			segmentDir:  "/tmp/segments",
			trackIndex:  0,
			segmentFile: "seg-00000.ts",
			want:        "/tmp/segments/audio/0/seg-00000.ts",
		},
		{
			name:        "track 3 later segment",
			segmentDir:  "/tmp/segments",
			trackIndex:  3,
			segmentFile: "seg-00042.ts",
			want:        "/tmp/segments/audio/3/seg-00042.ts",
		},
		{
			name:        "deep directory",
			segmentDir:  "/data/revenge/sessions/abc123",
			trackIndex:  1,
			segmentFile: "seg-00001.ts",
			want:        "/data/revenge/sessions/abc123/audio/1/seg-00001.ts",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := AudioRenditionSegmentPath(tc.segmentDir, tc.trackIndex, tc.segmentFile)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGenerateMasterPlaylist_WithAudioAndSubtitles(t *testing.T) {
	profiles := []ProfileVariant{
		{Name: "original", Width: 1920, Height: 1080, Bandwidth: 8000000, VideoCodec: "h264"},
		{Name: "720p", Width: 1280, Height: 720, Bandwidth: 3000000, VideoCodec: "libx264"},
	}
	audio := []AudioVariant{
		{Index: 0, Name: "English 5.1", Language: "en", Channels: 6, IsDefault: true, Codec: "aac"},
	}
	subtitles := []SubtitleVariant{
		{Index: 0, Name: "English", Language: "en", IsDefault: true},
	}

	playlist := GenerateMasterPlaylist(profiles, audio, subtitles)

	// Stream variants should have both AUDIO and SUBTITLES attributes
	assert.Contains(t, playlist, `AUDIO="audio"`)
	assert.Contains(t, playlist, `SUBTITLES="subs"`)

	// Verify CODECS are present with audio codec included
	assert.Contains(t, playlist, `CODECS="avc1.640028,mp4a.40.2"`)
	assert.Contains(t, playlist, `CODECS="avc1.64001f,mp4a.40.2"`)
	assert.Contains(t, playlist, `AUDIO="audio",SUBTITLES="subs"`)
}

func TestGenerateMasterPlaylist_NoAudioNoSubtitles(t *testing.T) {
	profiles := []ProfileVariant{
		{Name: "original", Width: 1920, Height: 1080, Bandwidth: 5000000, VideoCodec: "h264"},
	}

	playlist := GenerateMasterPlaylist(profiles, nil, nil)

	// Should not contain audio or subtitle group references
	assert.NotContains(t, playlist, `AUDIO=`)
	assert.NotContains(t, playlist, `SUBTITLES=`)
	assert.NotContains(t, playlist, "EXT-X-MEDIA")
}

func TestGenerateMasterPlaylist_MultipleAudioDefaultSelection(t *testing.T) {
	audio := []AudioVariant{
		{Index: 0, Name: "English", Language: "en", Channels: 2, IsDefault: false},
		{Index: 1, Name: "Spanish", Language: "es", Channels: 2, IsDefault: true},
		{Index: 2, Name: "French", Language: "fr", Channels: 6, IsDefault: false},
	}

	playlist := GenerateMasterPlaylist(nil, audio, nil)

	// Only the default track should have DEFAULT=YES
	assert.Contains(t, playlist, `NAME="English",DEFAULT=NO,AUTOSELECT=NO`)
	assert.Contains(t, playlist, `NAME="Spanish",DEFAULT=YES,AUTOSELECT=YES`)
	assert.Contains(t, playlist, `NAME="French",DEFAULT=NO,AUTOSELECT=NO`)
}

func TestCleanHEVCCodecString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "DV Main 10 L150 constraint 90 stripped",
			input:    "hvc1.2.4.L150.90",
			expected: "hvc1.2.4.L150",
		},
		{
			name:     "DV Main 10 H153 constraint 90 stripped",
			input:    "hvc1.2.4.H153.90",
			expected: "hvc1.2.4.H153",
		},
		{
			name:     "standard B0 stripped",
			input:    "hvc1.1.6.L120.B0",
			expected: "hvc1.1.6.L120",
		},
		{
			name:     "multi-byte constraints stripped",
			input:    "hvc1.2.4.L150.90.00.00",
			expected: "hvc1.2.4.L150",
		},
		{
			name:     "hev1 prefix also cleaned",
			input:    "hev1.2.4.L150.90",
			expected: "hev1.2.4.L150",
		},
		{
			name:     "no constraint part unchanged",
			input:    "hvc1.2.4.L150",
			expected: "hvc1.2.4.L150",
		},
		{
			name:     "non-HEVC string unchanged",
			input:    "avc1.640028",
			expected: "avc1.640028",
		},
		{
			name:     "empty string unchanged",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, cleanHEVCCodecString(tt.input))
		})
	}
}
