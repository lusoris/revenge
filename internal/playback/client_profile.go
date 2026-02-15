package playback

import (
	"strings"
)

// ClientProfile describes the media capabilities of the requesting client.
// Sent by the frontend (which probes via MediaSource.isTypeSupported / canPlayType)
// and optionally supplemented by server-side User-Agent detection.
//
// This follows the same architecture as Jellyfin: the client self-reports what it
// can play, and the server uses this to decide whether to direct-stream, remux
// (with optional metadata stripping), or transcode.
type ClientProfile struct {
	// VideoCodecs the client can decode (e.g. ["h264","hevc","av1"]).
	// Empty = server decides based on User-Agent.
	VideoCodecs []string `json:"video_codecs,omitempty"`

	// AudioCodecs the client can decode (e.g. ["aac","mp3","opus","flac","ac3","eac3"]).
	// Empty = server decides based on User-Agent.
	AudioCodecs []string `json:"audio_codecs,omitempty"`

	// HDR capability flags — what dynamic range types the client can display.
	SupportsHDR10       bool `json:"supports_hdr10,omitempty"`
	SupportsHLG         bool `json:"supports_hlg,omitempty"`
	SupportsDolbyVision bool `json:"supports_dolby_vision,omitempty"`

	// MaxWidth and MaxHeight limit output resolution (0 = unlimited).
	MaxWidth  int `json:"max_width,omitempty"`
	MaxHeight int `json:"max_height,omitempty"`

	// MaxBitrateKbps limits the total bitrate (0 = unlimited).
	MaxBitrateKbps int `json:"max_bitrate_kbps,omitempty"`

	// MaxAudioChannels limits audio channel count (0 = unlimited).
	// E.g., stereo-only devices would set 2.
	MaxAudioChannels int `json:"max_audio_channels,omitempty"`
}

// CanDecodeVideo returns true if the client profile declares support for the given video codec.
func (cp *ClientProfile) CanDecodeVideo(codec string) bool {
	if len(cp.VideoCodecs) == 0 {
		return false
	}
	codec = strings.ToLower(codec)
	for _, c := range cp.VideoCodecs {
		if strings.ToLower(c) == codec {
			return true
		}
	}
	return false
}

// CanDecodeAudio returns true if the client profile declares support for the given audio codec.
func (cp *ClientProfile) CanDecodeAudio(codec string) bool {
	if len(cp.AudioCodecs) == 0 {
		return false
	}
	codec = strings.ToLower(codec)
	for _, c := range cp.AudioCodecs {
		if strings.ToLower(c) == codec {
			return true
		}
	}
	return false
}

// CanDisplayDynamicRange returns true if the client can display the given dynamic range type.
// For DV content where the client doesn't support DV but the content has an HDR10/SDR
// fallback layer, this returns false — the caller should strip DV metadata and fall back.
func (cp *ClientProfile) CanDisplayDynamicRange(dynamicRange string) bool {
	switch dynamicRange {
	case "SDR":
		return true // everyone supports SDR
	case "HDR10":
		return cp.SupportsHDR10
	case "HLG":
		return cp.SupportsHLG
	case "Dolby Vision":
		return cp.SupportsDolbyVision
	case "HDR":
		return cp.SupportsHDR10 // generic HDR ≈ HDR10
	default:
		return true // unknown → assume supported
	}
}

// DefaultBrowserProfile returns a conservative profile for generic web browsers.
// H.264 + AAC only, SDR only. Used when the client doesn't send capabilities
// and User-Agent detection doesn't match a known pattern.
func DefaultBrowserProfile() ClientProfile {
	return ClientProfile{
		VideoCodecs:     []string{"h264"},
		AudioCodecs:     []string{"aac", "mp3", "opus", "flac"},
		SupportsHDR10:   false,
		SupportsHLG:     false,
		MaxAudioChannels: 2,
	}
}

// ChromeDesktopProfile returns capabilities for Chrome on desktop.
// Chrome supports HEVC (hardware), AV1, H.264. Does NOT support DV.
// Supports HDR10 passthrough on compatible displays (Edge Chromium 121+, Chrome desktop).
func ChromeDesktopProfile() ClientProfile {
	return ClientProfile{
		VideoCodecs:     []string{"h264", "hevc", "av1"},
		AudioCodecs:     []string{"aac", "mp3", "opus", "flac"},
		SupportsHDR10:   true,
		SupportsHLG:     true,
		MaxAudioChannels: 6,
	}
}

// FirefoxDesktopProfile returns capabilities for Firefox on desktop.
// Firefox supports H.264, AV1. HEVC support varies (hardware-dependent on Windows/macOS).
func FirefoxDesktopProfile() ClientProfile {
	return ClientProfile{
		VideoCodecs:     []string{"h264", "av1"},
		AudioCodecs:     []string{"aac", "mp3", "opus", "flac"},
		SupportsHDR10:   false,
		SupportsHLG:     false,
		MaxAudioChannels: 6,
	}
}

// SafariProfile returns capabilities for Safari on macOS/iOS.
// Safari natively supports HEVC, H.264, and Dolby Vision (on compatible hardware).
func SafariProfile() ClientProfile {
	return ClientProfile{
		VideoCodecs:         []string{"h264", "hevc"},
		AudioCodecs:         []string{"aac", "mp3", "opus", "flac", "ac3", "eac3"},
		SupportsHDR10:       true,
		SupportsHLG:         true,
		SupportsDolbyVision: true,
		MaxAudioChannels:    8,
	}
}

// ProfileFromUserAgent returns a client profile based on User-Agent string parsing.
// This is a fallback — the preferred approach is for the client to probe and send
// its actual capabilities. Returns nil if no known pattern matches.
func ProfileFromUserAgent(ua string) *ClientProfile {
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") && !strings.Contains(ua, "chromium"):
		// Safari (macOS / iOS) — native HEVC + DV support
		p := SafariProfile()
		return &p

	case strings.Contains(ua, "firefox"):
		p := FirefoxDesktopProfile()
		return &p

	case strings.Contains(ua, "chrome") || strings.Contains(ua, "chromium") || strings.Contains(ua, "edg/"):
		// Chrome / Chromium / Edge — HEVC via hardware, no DV
		p := ChromeDesktopProfile()
		return &p

	default:
		return nil
	}
}
