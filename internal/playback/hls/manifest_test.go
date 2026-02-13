package hls

import (
	"strings"
	"testing"

	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMasterPlaylist_Basic(t *testing.T) {
	profiles := []ProfileVariant{
		{Name: "original", Width: 1920, Height: 1080, Bandwidth: 8000000},
		{Name: "720p", Width: 1280, Height: 720, Bandwidth: 3000000},
		{Name: "480p", Width: 854, Height: 480, Bandwidth: 1500000},
	}

	playlist := GenerateMasterPlaylist(profiles, nil, nil)

	assert.Contains(t, playlist, "#EXTM3U")
	assert.Contains(t, playlist, "#EXT-X-VERSION:3")
	assert.Contains(t, playlist, `BANDWIDTH=8000000,RESOLUTION=1920x1080,NAME="original"`)
	assert.Contains(t, playlist, "original/index.m3u8")
	assert.Contains(t, playlist, `BANDWIDTH=3000000,RESOLUTION=1280x720,NAME="720p"`)
	assert.Contains(t, playlist, "720p/index.m3u8")
	assert.Contains(t, playlist, "480p/index.m3u8")

	// No subtitles attribute when no subs
	assert.NotContains(t, playlist, "SUBTITLES")
}

func TestGenerateMasterPlaylist_WithSubtitles(t *testing.T) {
	profiles := []ProfileVariant{
		{Name: "original", Width: 1920, Height: 1080, Bandwidth: 8000000},
	}
	subtitles := []SubtitleVariant{
		{Index: 0, Name: "English", Language: "en", IsDefault: true},
		{Index: 1, Name: "German", Language: "de", IsDefault: false},
	}

	playlist := GenerateMasterPlaylist(profiles, nil, subtitles)

	assert.Contains(t, playlist, `#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",NAME="English",DEFAULT=YES,LANGUAGE="en",URI="subs/0.vtt"`)
	assert.Contains(t, playlist, `#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",NAME="German",DEFAULT=NO,LANGUAGE="de",URI="subs/1.vtt"`)
	assert.Contains(t, playlist, `SUBTITLES="subs"`)
}

func TestGenerateMasterPlaylist_EmptyProfiles(t *testing.T) {
	playlist := GenerateMasterPlaylist(nil, nil, nil)

	assert.Contains(t, playlist, "#EXTM3U")
	// Should be valid but minimal
	lines := strings.Split(strings.TrimSpace(playlist), "\n")
	assert.GreaterOrEqual(t, len(lines), 2) // At least #EXTM3U and #EXT-X-VERSION
}

func TestGenerateMasterPlaylist_WithAudioTracks(t *testing.T) {
	profiles := []ProfileVariant{
		{Name: "original", Width: 1920, Height: 1080, Bandwidth: 8000000},
	}
	audio := []AudioVariant{
		{Index: 0, Name: "English 5.1", Language: "en", Channels: 6, IsDefault: true},
		{Index: 1, Name: "German", Language: "de", Channels: 2, IsDefault: false},
	}

	playlist := GenerateMasterPlaylist(profiles, audio, nil)

	assert.Contains(t, playlist, `#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="English 5.1",DEFAULT=YES,AUTOSELECT=YES,LANGUAGE="en",CHANNELS="6",URI="audio/0/index.m3u8"`)
	assert.Contains(t, playlist, `#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="German",DEFAULT=NO,AUTOSELECT=NO,LANGUAGE="de",CHANNELS="2",URI="audio/1/index.m3u8"`)
	assert.Contains(t, playlist, `AUDIO="audio"`)
	// No subtitles
	assert.NotContains(t, playlist, "SUBTITLES")
}

func TestEstimateBandwidth_WithSourceBitrates(t *testing.T) {
	from := estimateBandwidth(
		transcode.ProfileDecision{VideoBitrate: 5000, AudioBitrate: 192},
		0, 0,
	)

	// 5000*1000 + 192*1000 = 5192000 * 1.1 = 5711200
	expected := int(float64(5000*1000+192*1000) * 1.1)
	assert.Equal(t, expected, from)
}

func TestSegmentPath(t *testing.T) {
	path := SegmentPath("/tmp/segments", "720p", "seg-00001.ts")
	assert.Equal(t, "/tmp/segments/720p/seg-00001.ts", path)
}

func TestSubtitlePath(t *testing.T) {
	path := SubtitlePath("/tmp/segments", 2)
	assert.Equal(t, "/tmp/segments/subs/2.vtt", path)
}
