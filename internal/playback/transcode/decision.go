package transcode

import (
	"github.com/lusoris/revenge/internal/content/movie"
)

// HLS-compatible codecs that can be remuxed without transcoding.
var hlsCompatibleVideoCodecs = map[string]bool{
	"h264": true,
}

var hlsCompatibleAudioCodecs = map[string]bool{
	"aac":  true,
	"mp3":  true,
	"ac3":  true, // Dolby Digital 5.1
	"eac3": true, // Dolby Digital Plus / Atmos
}

// Decision describes the transcode/remux decision for a media file.
type Decision struct {
	CanRemux         bool
	SourceVideoCodec string
	SourceAudioCodec string
	SourceWidth      int
	SourceHeight     int
	Profiles         []ProfileDecision
}

// ProfileDecision describes the transcode/remux decision for a quality profile.
type ProfileDecision struct {
	Name           string
	Width          int
	Height         int
	VideoBitrate   int    // kbps (0 = copy)
	AudioBitrate   int    // kbps (0 = copy)
	NeedsTranscode bool
	VideoCodec     string // "copy" or "libx264"
	AudioCodec     string // "copy" or "aac"
}

// AnalyzeMedia examines probed MediaInfo and determines which profiles need
// transcoding and which can be remuxed directly to HLS.
func AnalyzeMedia(info *movie.MediaInfo, profiles []QualityProfile) Decision {
	videoCodec := info.VideoCodec
	audioCodec := ""
	if len(info.AudioStreams) > 0 {
		audioCodec = info.AudioStreams[0].Codec
	}

	canRemuxVideo := hlsCompatibleVideoCodecs[videoCodec]
	canRemuxAudio := hlsCompatibleAudioCodecs[audioCodec]
	canRemux := canRemuxVideo && canRemuxAudio

	d := Decision{
		CanRemux:         canRemux,
		SourceVideoCodec: videoCodec,
		SourceAudioCodec: audioCodec,
		SourceWidth:      info.Width,
		SourceHeight:     info.Height,
	}

	for _, p := range profiles {
		pd := analyzeProfile(p, info, canRemuxVideo, canRemuxAudio)
		if pd != nil {
			d.Profiles = append(d.Profiles, *pd)
		}
	}

	return d
}

func analyzeProfile(p QualityProfile, info *movie.MediaInfo, canRemuxVideo, canRemuxAudio bool) *ProfileDecision {
	isOriginal := p.MaxHeight == 0 && p.MaxWidth == 0

	pd := &ProfileDecision{
		Name: p.Name,
	}

	if isOriginal {
		// Original profile: use source dimensions
		pd.Width = info.Width
		pd.Height = info.Height

		if canRemuxVideo {
			pd.VideoCodec = "copy"
			pd.VideoBitrate = 0
		} else {
			pd.NeedsTranscode = true
			pd.VideoCodec = "libx264"
			pd.VideoBitrate = estimateOriginalBitrate(info)
		}

		if canRemuxAudio {
			pd.AudioCodec = "copy"
			pd.AudioBitrate = 0
		} else {
			pd.NeedsTranscode = true
			pd.AudioCodec = "aac"
			pd.AudioBitrate = 192
		}
	} else {
		// Sized profile: scale down if needed
		pd.Width = p.MaxWidth
		pd.Height = p.MaxHeight

		// If source is smaller or equal, use source dimensions
		if info.Height > 0 && info.Height <= p.MaxHeight {
			pd.Width = info.Width
			pd.Height = info.Height

			// Can potentially remux if codec is compatible
			if canRemuxVideo {
				pd.VideoCodec = "copy"
				pd.VideoBitrate = 0
			} else {
				pd.NeedsTranscode = true
				pd.VideoCodec = "libx264"
				pd.VideoBitrate = p.VideoBitrate
			}
		} else {
			// Must transcode to scale down
			pd.NeedsTranscode = true
			pd.VideoCodec = "libx264"
			pd.VideoBitrate = p.VideoBitrate
		}

		// Audio always needs transcode to AAC for sized profiles (consistent output)
		if canRemuxAudio && !pd.NeedsTranscode {
			pd.AudioCodec = "copy"
			pd.AudioBitrate = 0
		} else {
			pd.NeedsTranscode = true
			pd.AudioCodec = "aac"
			pd.AudioBitrate = p.AudioBitrate
		}
	}

	return pd
}

// estimateOriginalBitrate returns a reasonable bitrate for transcoding at original resolution.
func estimateOriginalBitrate(info *movie.MediaInfo) int {
	if info.VideoBitrateKbps > 0 {
		// Use source bitrate but cap at reasonable levels for H.264
		bitrate := int(info.VideoBitrateKbps)
		if bitrate > 20000 {
			bitrate = 20000
		}
		return bitrate
	}

	// Estimate from resolution
	switch {
	case info.Height >= 2160:
		return 15000
	case info.Height >= 1440:
		return 10000
	case info.Height >= 1080:
		return 5000
	case info.Height >= 720:
		return 2800
	default:
		return 1400
	}
}
