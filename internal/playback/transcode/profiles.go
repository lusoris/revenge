// Package transcode provides FFmpeg-based transcoding and remuxing for HLS streaming.
package transcode

// QualityProfile defines a target encoding quality for HLS streaming.
type QualityProfile struct {
	Name         string // "original", "1080p", "720p", "480p"
	MaxWidth     int    // 0 = no limit (original)
	MaxHeight    int    // 0 = no limit (original)
	VideoBitrate int    // kbps (0 = copy)
	AudioBitrate int    // kbps (0 = copy)
	VideoCodec   string // "copy" or "libx264"
	AudioCodec   string // "copy" or "aac"
}

// DefaultProfiles contains the standard quality profiles for HLS streaming.
var DefaultProfiles = map[string]QualityProfile{
	"original": {
		Name:         "original",
		MaxWidth:     0,
		MaxHeight:    0,
		VideoBitrate: 0,
		AudioBitrate: 0,
		VideoCodec:   "copy",
		AudioCodec:   "copy",
	},
	"1080p": {
		Name:         "1080p",
		MaxWidth:     1920,
		MaxHeight:    1080,
		VideoBitrate: 5000,
		AudioBitrate: 192,
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
	},
	"720p": {
		Name:         "720p",
		MaxWidth:     1280,
		MaxHeight:    720,
		VideoBitrate: 2800,
		AudioBitrate: 128,
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
	},
	"480p": {
		Name:         "480p",
		MaxWidth:     854,
		MaxHeight:    480,
		VideoBitrate: 1400,
		AudioBitrate: 96,
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
	},
}

// GetEnabledProfiles returns the quality profiles matching the given names.
// Unknown names are silently skipped.
func GetEnabledProfiles(names []string) []QualityProfile {
	profiles := make([]QualityProfile, 0, len(names))
	for _, name := range names {
		if p, ok := DefaultProfiles[name]; ok {
			profiles = append(profiles, p)
		}
	}
	return profiles
}

// EstimateBandwidth returns the estimated total bandwidth in bits/sec for a profile.
// Used in HLS master playlist BANDWIDTH attribute.
func (p QualityProfile) EstimateBandwidth(sourceVideoBitrate, sourceAudioBitrate int64) int {
	videoBps := p.VideoBitrate * 1000
	if videoBps == 0 && sourceVideoBitrate > 0 {
		videoBps = int(sourceVideoBitrate) * 1000
	}
	audioBps := p.AudioBitrate * 1000
	if audioBps == 0 && sourceAudioBitrate > 0 {
		audioBps = int(sourceAudioBitrate) * 1000
	}
	// Add 10% overhead for container/muxing
	return int(float64(videoBps+audioBps) * 1.1)
}
