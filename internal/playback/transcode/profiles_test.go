package transcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnabledProfiles(t *testing.T) {
	profiles := GetEnabledProfiles([]string{"original", "720p", "480p"})

	require.Len(t, profiles, 3)
	assert.Equal(t, "original", profiles[0].Name)
	assert.Equal(t, "720p", profiles[1].Name)
	assert.Equal(t, "480p", profiles[2].Name)
}

func TestGetEnabledProfiles_UnknownSkipped(t *testing.T) {
	profiles := GetEnabledProfiles([]string{"original", "4k", "nonexistent"})

	require.Len(t, profiles, 1)
	assert.Equal(t, "original", profiles[0].Name)
}

func TestGetEnabledProfiles_Empty(t *testing.T) {
	profiles := GetEnabledProfiles([]string{})
	assert.Empty(t, profiles)
}

func TestQualityProfile_EstimateBandwidth(t *testing.T) {
	p := DefaultProfiles["1080p"]
	bw := p.EstimateBandwidth(0, 0)

	// 5000 video + 192 audio = 5192 kbps * 1000 = 5192000 bps * 1.1 = ~5711200
	expected := int(float64(5000*1000+192*1000) * 1.1)
	assert.Equal(t, expected, bw)
}

func TestQualityProfile_EstimateBandwidth_Original(t *testing.T) {
	p := DefaultProfiles["original"]
	bw := p.EstimateBandwidth(8000, 320)

	// Uses source bitrates: 8000 + 320 = 8320 kbps * 1000 * 1.1
	expected := int(float64(8000*1000+320*1000) * 1.1)
	assert.Equal(t, expected, bw)
}

func TestDefaultProfiles_Completeness(t *testing.T) {
	required := []string{"original", "1080p", "720p", "480p"}
	for _, name := range required {
		_, ok := DefaultProfiles[name]
		assert.True(t, ok, "default profiles should include %s", name)
	}
}

func TestDefaultProfiles_OriginalIsCopy(t *testing.T) {
	p := DefaultProfiles["original"]
	assert.Equal(t, "copy", p.VideoCodec)
	assert.Equal(t, "copy", p.AudioCodec)
	assert.Equal(t, 0, p.VideoBitrate)
	assert.Equal(t, 0, p.AudioBitrate)
}
