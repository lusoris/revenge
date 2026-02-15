package transcode

import (
	"github.com/lusoris/revenge/internal/content/movie"
)

// ClientCapabilities is the subset of client profile data needed for transcode decisions.
// This avoids importing the playback package (which would create a circular dependency).
type ClientCapabilities struct {
	VideoCodecs         []string // codecs the client can decode
	AudioCodecs         []string // codecs the client can decode
	SupportsDolbyVision bool
	SupportsHDR10       bool
}

// hlsCompatibleVideoCodecs lists codecs that can be carried in fMP4 HLS segments.
// Whether the client can actually decode them is a player-side concern —
// HLS.js uses MediaSource.isTypeSupported() to skip unsupported levels.
var hlsCompatibleVideoCodecs = map[string]bool{
	"h264": true,
	"hevc": true, // H.265 — supported via fMP4 segments
	"h265": true, // alias for hevc
	"av1":  true, // AV1 — supported via fMP4 segments
}

// Audio codecs that browsers can decode via MediaSource Extensions (MSE).
// While fMP4 can carry AC-3/E-AC-3/TrueHD, Chrome and Firefox cannot decode
// them — only Safari and native apps support Dolby codecs. We must transcode
// unsupported codecs to AAC for reliable web playback.
var hlsCompatibleAudioCodecs = map[string]bool{
	"aac":  true,
	"mp3":  true,
	"opus": true,
	"flac": true,
}

// Decision describes the transcode/remux decision for a media file.
type Decision struct {
	CanRemux              bool
	SourceVideoCodec      string
	SourceVideoCodecString string // RFC 6381 CODECS string from actual extradata
	SourceAudioCodec      string
	SourceWidth           int
	SourceHeight          int
	SourceVideoBitrateKbps int64 // source video bitrate in kbps (0 = unknown)
	Profiles              []ProfileDecision
}

// ProfileDecision describes the transcode/remux decision for a quality profile.
type ProfileDecision struct {
	Name              string
	Width             int
	Height            int
	VideoBitrate      int // kbps (0 = copy)
	AudioBitrate      int // kbps (0 = copy)
	NeedsTranscode    bool
	VideoCodec        string // "copy" or "libx264"
	AudioCodec        string // "copy" or "aac"
	StripDolbyVision  bool   // strip DV metadata from HEVC (for clients that can't decode DV)
}

// AnalyzeMedia examines probed MediaInfo and determines which profiles need
// transcoding and which can be remuxed directly to HLS.
// clientCaps is optional — when nil, server-side defaults are used (conservative: no DV).
func AnalyzeMedia(info *movie.MediaInfo, profiles []QualityProfile, clientCaps *ClientCapabilities) Decision {
	videoCodec := info.VideoCodec
	audioCodec := ""
	if len(info.AudioStreams) > 0 {
		audioCodec = info.AudioStreams[0].Codec
	}

	canRemuxVideo := hlsCompatibleVideoCodecs[videoCodec]
	canRemuxAudio := hlsCompatibleAudioCodecs[audioCodec]
	canRemux := canRemuxVideo && canRemuxAudio

	// Determine if DV metadata needs to be stripped.
	// Strip when: content is DV AND client doesn't support DV.
	isDV := info.DynamicRange == "Dolby Vision"
	stripDV := false
	if isDV {
		if clientCaps == nil || !clientCaps.SupportsDolbyVision {
			stripDV = true
		}
	}

	d := Decision{
		CanRemux:               canRemux,
		SourceVideoCodec:       videoCodec,
		SourceVideoCodecString: info.VideoCodecString,
		SourceAudioCodec:       audioCodec,
		SourceWidth:            info.Width,
		SourceHeight:           info.Height,
		SourceVideoBitrateKbps: info.VideoBitrateKbps,
	}

	for _, p := range profiles {
		pd := analyzeProfile(p, info, canRemuxVideo, canRemuxAudio)
		if pd != nil {
			// Set DV stripping flag on remuxed HEVC profiles.
			// Transcoded profiles (H.264) don't carry DV metadata anyway.
			if stripDV && pd.VideoCodec == "copy" {
				pd.StripDolbyVision = true
			}
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
			// DTS, DTS-HD, DTS-X, etc. → transcode to AAC
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

			// Only copy H.264 — it's universally browser-compatible.
			// HEVC/AV1 may fail in some browsers' MSE, so sized profiles
			// always transcode to H.264 to serve as reliable fallbacks.
			if canRemuxVideo && info.VideoCodec == "h264" {
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

		// Audio: transcode to AAC for consistent output in sized profiles
		pd.NeedsTranscode = true
		pd.AudioCodec = "aac"
		pd.AudioBitrate = p.AudioBitrate
	}

	return pd
}

// estimateOriginalBitrate returns a reasonable bitrate for transcoding at original resolution.
func estimateOriginalBitrate(info *movie.MediaInfo) int {
	if info.VideoBitrateKbps > 0 {
		// Use source bitrate but cap at reasonable levels for H.264
		bitrate := min(int(info.VideoBitrateKbps), 20000)
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
