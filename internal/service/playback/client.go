// Package playback provides client detection, bandwidth monitoring, and transcoding integration.
package playback

import (
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// ClientCapabilities represents what a client can play directly.
type ClientCapabilities struct {
	// Device identification
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"` // "tv", "mobile", "desktop", "browser"
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`

	// Video capabilities
	MaxVideoWidth       int      `json:"max_video_width"`
	MaxVideoHeight      int      `json:"max_video_height"`
	SupportedCodecs     []string `json:"supported_codecs"` // ["h264", "hevc", "av1", "vp9"]
	SupportsHDR         bool     `json:"supports_hdr"`
	SupportsHDR10       bool     `json:"supports_hdr10"`
	SupportsDolbyVision bool     `json:"supports_dolby_vision"`

	// Audio capabilities
	MaxAudioChannels     int      `json:"max_audio_channels"`     // 2, 6 (5.1), 8 (7.1)
	SupportedAudioCodecs []string `json:"supported_audio_codecs"` // ["aac", "ac3", "eac3", "dts", "truehd"]
	SupportsAtmos        bool     `json:"supports_atmos"`

	// Container support
	SupportedContainers []string `json:"supported_containers"` // ["mp4", "mkv", "webm", "hls", "dash"]

	// Subtitle support
	SupportedSubtitleFormats  []string `json:"supported_subtitle_formats"` // ["srt", "ass", "vtt", "pgs"]
	SupportsEmbeddedSubtitles bool     `json:"supports_embedded_subtitles"`
}

// ClientInfo contains detected client information.
type ClientInfo struct {
	SessionID    uuid.UUID
	UserID       uuid.UUID
	DeviceID     string
	IsExternal   bool // Outside local network
	IPAddress    string
	UserAgent    string
	Capabilities *ClientCapabilities
}

// ClientDetector detects client capabilities from requests.
type ClientDetector struct {
	localNetworks []*net.IPNet
}

// NewClientDetector creates a new client detector.
func NewClientDetector(localCIDRs []string) (*ClientDetector, error) {
	networks := make([]*net.IPNet, 0, len(localCIDRs))
	for _, cidr := range localCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		networks = append(networks, network)
	}

	// Add common private networks by default
	defaultCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
		"fe80::/10",
	}
	for _, cidr := range defaultCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		networks = append(networks, network)
	}

	return &ClientDetector{localNetworks: networks}, nil
}

// DetectClient extracts client information from an HTTP request.
func (d *ClientDetector) DetectClient(r *http.Request, userID uuid.UUID, deviceID string) *ClientInfo {
	ip := d.extractIP(r)

	return &ClientInfo{
		SessionID:  uuid.New(),
		UserID:     userID,
		DeviceID:   deviceID,
		IsExternal: !d.isLocalIP(ip),
		IPAddress:  ip,
		UserAgent:  r.UserAgent(),
	}
}

// extractIP gets the real client IP from the request.
func (d *ClientDetector) extractIP(r *http.Request) string {
	// Check X-Forwarded-For header (first IP is the client)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// isLocalIP checks if an IP is in any local network.
func (d *ClientDetector) isLocalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, network := range d.localNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// DetectDeviceType guesses device type from User-Agent.
func DetectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "tv") || strings.Contains(ua, "tizen") ||
		strings.Contains(ua, "webos") || strings.Contains(ua, "roku"):
		return "tv"
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "android") ||
		strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		return "mobile"
	case strings.Contains(ua, "electron") || strings.Contains(ua, "revenge"):
		return "desktop"
	default:
		return "browser"
	}
}

// DetectDeviceGroup determines the device group for profile selection.
func DetectDeviceGroup(userAgent string, capabilities *ClientCapabilities) DeviceGroup {
	ua := strings.ToLower(userAgent)

	// TV detection with 4K capability check
	if strings.Contains(ua, "tv") || strings.Contains(ua, "tizen") ||
		strings.Contains(ua, "webos") || strings.Contains(ua, "roku") ||
		strings.Contains(ua, "android tv") || strings.Contains(ua, "fire tv") {
		if capabilities != nil && capabilities.MaxVideoWidth >= 3840 {
			return DeviceGroupTV4K
		}
		return DeviceGroupTVHD
	}

	// iOS devices
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return DeviceGroupMobileiOS
	}

	// Android mobile
	if strings.Contains(ua, "android") &&
		(strings.Contains(ua, "mobile") || !strings.Contains(ua, "tv")) {
		return DeviceGroupMobileAndroid
	}

	// Desktop app (Electron, native)
	if strings.Contains(ua, "electron") || strings.Contains(ua, "revenge-desktop") {
		return DeviceGroupDesktopApp
	}

	// Browser detection
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return DeviceGroupBrowserLegacy // Safari has limited codec support
	}

	// Modern browsers (Chrome, Firefox, Edge)
	if strings.Contains(ua, "chrome") || strings.Contains(ua, "firefox") ||
		strings.Contains(ua, "edg/") {
		return DeviceGroupBrowserModern
	}

	// Fallback
	return DeviceGroupBrowserLegacy
}

// SelectProfile selects the best transcode profile for a client.
func SelectProfile(client *ClientInfo, bandwidthEstimate *BandwidthEstimate) *TranscodeProfile {
	// Determine device group
	group := DetectDeviceGroup(client.UserAgent, client.Capabilities)

	// Get base profile for device
	baseProfile := GetProfileForDeviceGroup(group)
	if baseProfile == nil {
		baseProfile = DefaultProfiles["h264_720p"]
	}

	// For external clients with bandwidth constraints, potentially downgrade
	if client.IsExternal && bandwidthEstimate != nil && bandwidthEstimate.IsReliable {
		bandwidthProfileID := GetProfileForBandwidth(bandwidthEstimate.RecommendedKbps)
		if bandwidthProfileID != "" {
			bandwidthProfile := GetProfile(bandwidthProfileID)
			// Use bandwidth profile if it has lower bitrate than device default
			if bandwidthProfile != nil && bandwidthProfile.MaxBitrate < baseProfile.MaxBitrate {
				return bandwidthProfile
			}
		}
	}

	return baseProfile
}

// DefaultCapabilitiesForDevice returns default capabilities for a device type.
func DefaultCapabilitiesForDevice(deviceType string) *ClientCapabilities {
	switch deviceType {
	case "tv":
		return &ClientCapabilities{
			DeviceType:           deviceType,
			MaxVideoWidth:        3840,
			MaxVideoHeight:       2160,
			SupportedCodecs:      []string{"h264", "hevc"},
			SupportsHDR:          true,
			MaxAudioChannels:     8,
			SupportedAudioCodecs: []string{"aac", "ac3", "eac3"},
			SupportedContainers:  []string{"mp4", "hls"},
		}
	case "mobile":
		return &ClientCapabilities{
			DeviceType:           deviceType,
			MaxVideoWidth:        1920,
			MaxVideoHeight:       1080,
			SupportedCodecs:      []string{"h264"},
			SupportsHDR:          false,
			MaxAudioChannels:     2,
			SupportedAudioCodecs: []string{"aac"},
			SupportedContainers:  []string{"mp4", "hls"},
		}
	case "desktop":
		return &ClientCapabilities{
			DeviceType:           deviceType,
			MaxVideoWidth:        3840,
			MaxVideoHeight:       2160,
			SupportedCodecs:      []string{"h264", "hevc", "av1"},
			SupportsHDR:          true,
			MaxAudioChannels:     8,
			SupportedAudioCodecs: []string{"aac", "ac3", "eac3", "dts"},
			SupportedContainers:  []string{"mp4", "mkv", "webm"},
		}
	default: // browser
		return &ClientCapabilities{
			DeviceType:           "browser",
			MaxVideoWidth:        1920,
			MaxVideoHeight:       1080,
			SupportedCodecs:      []string{"h264", "vp9"},
			SupportsHDR:          false,
			MaxAudioChannels:     2,
			SupportedAudioCodecs: []string{"aac", "opus"},
			SupportedContainers:  []string{"mp4", "webm", "hls"},
		}
	}
}
