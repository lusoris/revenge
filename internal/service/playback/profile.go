// Package playback provides transcode profile definitions.
package playback

// TranscodeProfile defines a pre-configured transcoding profile for Blackbeard.
type TranscodeProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Video settings
	VideoCodec   string `json:"video_codec"` // h264, hevc, av1, vp9
	MaxWidth     int    `json:"max_width"`
	MaxHeight    int    `json:"max_height"`
	MaxBitrate   int    `json:"max_bitrate_kbps"`
	MaxFramerate int    `json:"max_framerate"`

	// Audio settings
	AudioCodec    string `json:"audio_codec"` // aac, ac3, eac3, opus
	AudioChannels int    `json:"audio_channels"`
	AudioBitrate  int    `json:"audio_bitrate_kbps"`

	// Container and streaming
	Container    string `json:"container"`     // mp4, webm, ts, mkv
	StreamFormat string `json:"stream_format"` // hls, dash, progressive

	// Feature flags
	AllowHDR         bool `json:"allow_hdr"`
	AllowHDR10       bool `json:"allow_hdr10"`
	AllowDolbyVision bool `json:"allow_dolby_vision"`
	HardwareDecode   bool `json:"hardware_decode"`
	HardwareEncode   bool `json:"hardware_encode"`
}

// DeviceGroup represents a category of client devices.
type DeviceGroup string

const (
	DeviceGroupTV4K          DeviceGroup = "tv_4k"
	DeviceGroupTVHD          DeviceGroup = "tv_hd"
	DeviceGroupMobileiOS     DeviceGroup = "mobile_ios"
	DeviceGroupMobileAndroid DeviceGroup = "mobile_android"
	DeviceGroupDesktopApp    DeviceGroup = "desktop_app"
	DeviceGroupBrowserModern DeviceGroup = "browser_modern"
	DeviceGroupBrowserLegacy DeviceGroup = "browser_legacy"
	DeviceGroupLowBandwidth  DeviceGroup = "low_bandwidth"
)

// DefaultProfiles contains pre-configured transcode profiles.
var DefaultProfiles = map[string]*TranscodeProfile{
	// 4K HDR profiles
	"hevc_4k_hdr": {
		ID:             "hevc_4k_hdr",
		Name:           "4K HDR (HEVC)",
		VideoCodec:     "hevc",
		MaxWidth:       3840,
		MaxHeight:      2160,
		MaxBitrate:     40000,
		MaxFramerate:   60,
		AudioCodec:     "eac3",
		AudioChannels:  8,
		AudioBitrate:   640,
		Container:      "mp4",
		StreamFormat:   "hls",
		AllowHDR:       true,
		AllowHDR10:     true,
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"hevc_4k": {
		ID:             "hevc_4k",
		Name:           "4K SDR (HEVC)",
		VideoCodec:     "hevc",
		MaxWidth:       3840,
		MaxHeight:      2160,
		MaxBitrate:     25000,
		MaxFramerate:   60,
		AudioCodec:     "aac",
		AudioChannels:  6,
		AudioBitrate:   384,
		Container:      "mp4",
		StreamFormat:   "hls",
		AllowHDR:       false,
		HardwareDecode: true,
		HardwareEncode: true,
	},

	// 1080p profiles
	"h264_1080p": {
		ID:             "h264_1080p",
		Name:           "1080p (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       1920,
		MaxHeight:      1080,
		MaxBitrate:     8000,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  6,
		AudioBitrate:   256,
		Container:      "mp4",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"h264_1080p_hls": {
		ID:             "h264_1080p_hls",
		Name:           "1080p HLS (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       1920,
		MaxHeight:      1080,
		MaxBitrate:     8000,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  2,
		AudioBitrate:   192,
		Container:      "ts",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"hevc_1080p": {
		ID:             "hevc_1080p",
		Name:           "1080p (HEVC)",
		VideoCodec:     "hevc",
		MaxWidth:       1920,
		MaxHeight:      1080,
		MaxBitrate:     6000,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  6,
		AudioBitrate:   256,
		Container:      "mp4",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"vp9_1080p_dash": {
		ID:             "vp9_1080p_dash",
		Name:           "1080p DASH (VP9)",
		VideoCodec:     "vp9",
		MaxWidth:       1920,
		MaxHeight:      1080,
		MaxBitrate:     6000,
		MaxFramerate:   30,
		AudioCodec:     "opus",
		AudioChannels:  2,
		AudioBitrate:   128,
		Container:      "webm",
		StreamFormat:   "dash",
		HardwareDecode: true,
		HardwareEncode: false, // VP9 encode usually software
	},

	// 720p profiles
	"h264_720p": {
		ID:             "h264_720p",
		Name:           "720p (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       1280,
		MaxHeight:      720,
		MaxBitrate:     4000,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  2,
		AudioBitrate:   128,
		Container:      "mp4",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"h264_720p_hls": {
		ID:             "h264_720p_hls",
		Name:           "720p HLS (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       1280,
		MaxHeight:      720,
		MaxBitrate:     4000,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  2,
		AudioBitrate:   128,
		Container:      "ts",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},

	// Low bandwidth profiles
	"h264_480p": {
		ID:             "h264_480p",
		Name:           "480p (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       854,
		MaxHeight:      480,
		MaxBitrate:     1500,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  2,
		AudioBitrate:   96,
		Container:      "ts",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"h264_480p_hls": {
		ID:             "h264_480p_hls",
		Name:           "480p HLS (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       854,
		MaxHeight:      480,
		MaxBitrate:     1500,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  2,
		AudioBitrate:   96,
		Container:      "ts",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
	"h264_360p": {
		ID:             "h264_360p",
		Name:           "360p (H.264)",
		VideoCodec:     "h264",
		MaxWidth:       640,
		MaxHeight:      360,
		MaxBitrate:     800,
		MaxFramerate:   30,
		AudioCodec:     "aac",
		AudioChannels:  2,
		AudioBitrate:   64,
		Container:      "ts",
		StreamFormat:   "hls",
		HardwareDecode: true,
		HardwareEncode: true,
	},
}

// DeviceGroupProfiles maps device groups to their default profile.
var DeviceGroupProfiles = map[DeviceGroup]string{
	DeviceGroupTV4K:          "hevc_4k_hdr",
	DeviceGroupTVHD:          "h264_1080p",
	DeviceGroupMobileiOS:     "h264_1080p_hls",
	DeviceGroupMobileAndroid: "h264_1080p_hls",
	DeviceGroupDesktopApp:    "hevc_4k",
	DeviceGroupBrowserModern: "vp9_1080p_dash",
	DeviceGroupBrowserLegacy: "h264_720p_hls",
	DeviceGroupLowBandwidth:  "h264_480p_hls",
}

// BandwidthProfiles maps bandwidth ranges to profile overrides.
// Key is minimum bandwidth in kbps.
var BandwidthProfiles = []struct {
	MinKbps   int
	ProfileID string
}{
	{25000, ""}, // Use device default
	{15000, "hevc_1080p"},
	{8000, "h264_1080p"},
	{3000, "h264_720p"},
	{1500, "h264_480p"},
	{0, "h264_360p"},
}

// GetProfileForBandwidth returns the appropriate profile for given bandwidth.
func GetProfileForBandwidth(bandwidthKbps int) string {
	for _, bp := range BandwidthProfiles {
		if bandwidthKbps >= bp.MinKbps {
			return bp.ProfileID
		}
	}
	return "h264_360p"
}

// GetProfile returns a profile by ID.
func GetProfile(id string) *TranscodeProfile {
	return DefaultProfiles[id]
}

// GetProfileForDeviceGroup returns the default profile for a device group.
func GetProfileForDeviceGroup(group DeviceGroup) *TranscodeProfile {
	profileID := DeviceGroupProfiles[group]
	return DefaultProfiles[profileID]
}
