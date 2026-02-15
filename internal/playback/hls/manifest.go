// Package hls provides HLS manifest generation and segment serving for playback.
package hls

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/playback/transcode"
)

// GenerateMasterPlaylist creates an HLS master playlist (.m3u8) for a playback session.
// Uses EXT-X-VERSION:7 for fMP4 segment support (required for HEVC/AV1 passthrough).
// It references media playlists for each quality profile, audio renditions, and subtitle tracks.
// All audio tracks are muxed into segments so HLS.js can switch instantly without stream restart.
func GenerateMasterPlaylist(profiles []ProfileVariant, audioTracks []AudioVariant, subtitles []SubtitleVariant) string {
	var b strings.Builder

	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:7\n")
	b.WriteString("#EXT-X-INDEPENDENT-SEGMENTS\n")
	b.WriteString("\n")

	// Audio rendition entries — each track has its own segmented stream.
	// HLS.js downloads only the selected track's segments, preserving original
	// quality and saving bandwidth. Switching loads new audio segments instantly.
	for _, at := range audioTracks {
		defaultStr := "NO"
		autoSelect := "NO"
		if at.IsDefault {
			defaultStr = "YES"
			autoSelect = "YES"
		}
		fmt.Fprintf(&b, "#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio\",NAME=\"%s\",DEFAULT=%s,AUTOSELECT=%s,LANGUAGE=\"%s\",CHANNELS=\"%d\",URI=\"audio/%d/index.m3u8\"\n",
			at.Name, defaultStr, autoSelect, at.Language, at.Channels, at.Index)
	}

	if len(audioTracks) > 0 {
		b.WriteString("\n")
	}

	// Subtitle media entries
	for _, sub := range subtitles {
		defaultStr := "NO"
		if sub.IsDefault {
			defaultStr = "YES"
		}
		fmt.Fprintf(&b, "#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"subs\",NAME=\"%s\",DEFAULT=%s,LANGUAGE=\"%s\",URI=\"subs/%d.vtt\"\n",
			sub.Name, defaultStr, sub.Language, sub.Index)
	}

	if len(subtitles) > 0 {
		b.WriteString("\n")
	}

	// Determine the default audio codec string for CODECS attribute.
	// Used when audio is a separate rendition (not muxed into video variant).
	defaultAudioCodec := "mp4a.40.2" // AAC-LC fallback
	if len(audioTracks) > 0 {
		defaultAudioCodec = audioCodecString(audioTracks[0].Codec)
	}

	// Stream variants
	for _, p := range profiles {
		extraAttrs := ""
		if len(audioTracks) > 0 {
			extraAttrs += ",AUDIO=\"audio\""
		}
		if len(subtitles) > 0 {
			extraAttrs += ",SUBTITLES=\"subs\""
		}

		// Build RFC 6381 CODECS string (required for fMP4/HEVC/AV1)
		// Use the pre-built codec string from actual extradata when available.
		// Fall back to the resolution-based heuristic for transcoded profiles.
		codecs := p.VideoCodecString
		if codecs == "" {
			codecs = videoCodecString(p.VideoCodec, p.Height)
		}
		if len(audioTracks) > 0 {
			// Audio in separate rendition — include audio codec for compatibility
			codecs += "," + defaultAudioCodec
		}

		fmt.Fprintf(&b, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d,CODECS=\"%s\",NAME=\"%s\"%s\n",
			p.Bandwidth, p.Width, p.Height, codecs, p.Name, extraAttrs)
		fmt.Fprintf(&b, "%s/index.m3u8\n", p.Name)
	}

	return b.String()
}

// ProfileVariant describes a quality variant in the master playlist.
type ProfileVariant struct {
	Name            string
	Width           int
	Height          int
	Bandwidth       int    // bits per second
	VideoCodec      string // source codec: "h264", "hevc", "av1", "libx264"
	VideoCodecString string // pre-built RFC 6381 string from extradata (e.g. "hvc1.2.4.L150.90")
}

// AudioVariant describes an audio rendition in the master playlist.
// All audio tracks are muxed into segments, enabling instant track switching.
type AudioVariant struct {
	Index     int
	Name      string
	Language  string
	Channels  int
	IsDefault bool
	Codec     string // source codec: "aac", "ac3", "eac3", "opus", "flac", etc.
}

// SubtitleVariant describes a subtitle track in the master playlist.
type SubtitleVariant struct {
	Index     int
	Name      string
	Language  string
	IsDefault bool
}

// ProfileVariantsFromDecision creates playlist variants from transcode decisions.
func ProfileVariantsFromDecision(profiles []transcode.ProfileDecision, sourceVideoBitrate, sourceAudioBitrate int64) []ProfileVariant {
	variants := make([]ProfileVariant, 0, len(profiles))
	for _, pd := range profiles {
		bw := estimateBandwidth(pd, sourceVideoBitrate, sourceAudioBitrate)
		variants = append(variants, ProfileVariant{
			Name:      pd.Name,
			Width:     pd.Width,
			Height:    pd.Height,
			Bandwidth: bw,
		})
	}
	return variants
}

func estimateBandwidth(pd transcode.ProfileDecision, sourceVideoBps, sourceAudioBps int64) int {
	videoBps := pd.VideoBitrate * 1000
	if videoBps == 0 && sourceVideoBps > 0 {
		videoBps = int(sourceVideoBps) * 1000
	}
	audioBps := pd.AudioBitrate * 1000
	if audioBps == 0 && sourceAudioBps > 0 {
		audioBps = int(sourceAudioBps) * 1000
	}
	// Fallback estimates based on resolution when bitrate is unknown.
	// A flat 5 Mbps default is far too low for 4K content and causes the
	// original profile to sort below transcode profiles in the master playlist.
	if videoBps == 0 {
		switch {
		case pd.Height >= 2160:
			videoBps = 40_000_000 // 40 Mbps for 4K
		case pd.Height >= 1440:
			videoBps = 20_000_000 // 20 Mbps for 1440p
		case pd.Height >= 1080:
			videoBps = 10_000_000 // 10 Mbps for 1080p
		case pd.Height >= 720:
			videoBps = 5_000_000 // 5 Mbps for 720p
		default:
			videoBps = 2_000_000 // 2 Mbps fallback
		}
	}
	if audioBps == 0 {
		audioBps = 192000 // 192 kbps default
	}
	return int(float64(videoBps+audioBps) * 1.1)
}

// ReadMediaPlaylist reads the FFmpeg-generated media playlist from disk.
// If the file doesn't exist yet (FFmpeg still starting), it polls with retries.
func ReadMediaPlaylist(segmentDir, profile string) (string, error) {
	// Validate profile to prevent path traversal (CWE-22)
	for part := range strings.SplitSeq(profile, "/") {
		if part == "" || part == "." || part == ".." || strings.ContainsAny(part, "/\\") {
			return "", fmt.Errorf("invalid profile: %q", profile)
		}
	}

	playlistPath := filepath.Join(segmentDir, profile, "index.m3u8")

	// Ensure resolved path stays within segmentDir
	absPath, err := filepath.Abs(playlistPath)
	if err != nil {
		return "", fmt.Errorf("invalid playlist path: %w", err)
	}
	absDir, err := filepath.Abs(segmentDir)
	if err != nil {
		return "", fmt.Errorf("invalid segment dir: %w", err)
	}
	if !strings.HasPrefix(absPath, absDir+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal detected in profile: %q", profile)
	}

	// Poll for the playlist file to appear (FFmpeg may still be starting)
	// Use the validated absolute path for file operations.
	// On-demand transcodes need time before the first segment appears.
	// Poll for 15s — long enough for on-demand fallback profiles to produce
	// an init segment + first media segment. Returns 503 if still not ready.
	var content []byte
	for range 150 {
		content, err = os.ReadFile(absPath) //nolint:gosec // path validated above with traversal check
		if err == nil && len(content) > 0 {
			return string(content), nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err != nil {
		return "", fmt.Errorf("media playlist not available at %s: %w", absPath, err)
	}
	return "", fmt.Errorf("media playlist empty at %s", absPath)
}

// SegmentPath returns the filesystem path for a segment file.
// It validates that profile and segmentFile do not contain path traversal sequences.
func SegmentPath(segmentDir, profile, segmentFile string) string {
	// Caller should validate, but defense-in-depth: clean components
	profile = filepath.Base(profile)
	segmentFile = filepath.Base(segmentFile)
	return filepath.Join(segmentDir, profile, segmentFile)
}

// AudioRenditionSegmentPath returns the filesystem path for an audio rendition segment.
func AudioRenditionSegmentPath(segmentDir string, trackIndex int, segmentFile string) string {
	return filepath.Join(segmentDir, "audio", fmt.Sprintf("%d", trackIndex), segmentFile)
}

// SubtitlePath returns the filesystem path for a subtitle WebVTT file.
func SubtitlePath(segmentDir string, trackIndex int) string {
	return filepath.Join(segmentDir, "subs", fmt.Sprintf("%d.vtt", trackIndex))
}

// cleanHEVCCodecString strips the constraint indicator bytes from an HEVC
// CODECS string entirely. Dolby Vision content has DV-specific constraint
// bytes (e.g., ".90") and even standard ".B0" suffixes that Chrome/Firefox
// MSE reject via MediaSource.isTypeSupported().
//
// Jellyfin-web probes HEVC capability with bare strings like:
//   canPlayType('video/mp4; codecs="hvc1.2.4.L153"')
// — no constraint bytes at all. Chrome only recognises the 4-part form:
//   hvc1.{profile}.{compat}.{tier+level}
//
// For non-HEVC strings, returns the input unchanged.
func cleanHEVCCodecString(codecStr string) string {
	if !strings.HasPrefix(codecStr, "hvc1.") && !strings.HasPrefix(codecStr, "hev1.") {
		return codecStr
	}
	parts := strings.Split(codecStr, ".")
	// Format: hvc1.{profile}.{compat}.{tier+level}[.{constraints}...]
	// Keep only the first 4 parts — strip all constraint bytes.
	if len(parts) < 4 {
		return codecStr
	}
	return strings.Join(parts[:4], ".")
}

// videoCodecString returns an RFC 6381 codec string for the CODECS attribute.
// These are required for HLS.js to correctly identify fMP4 variant capabilities.
func videoCodecString(codec string, height int) string {
	switch codec {
	case "hevc", "h265":
		// hvc1.2.4.L{level} — Main 10 profile, Main tier, no constraint bytes.
		// Jellyfin-web probes with canPlayType('video/mp4; codecs="hvc1.2.4.L153"')
		// — Chrome/Firefox only accept the 4-part form without constraint suffixes.
		level := hevcLevel(height)
		return fmt.Sprintf("hvc1.2.4.L%d", level)
	case "av1":
		// av01.P.LLM.DD — Main profile, level from resolution, Main tier, 8/10-bit
		level := av1Level(height)
		return fmt.Sprintf("av01.0.%02dM.10", level)
	case "h264", "libx264", "":
		// avc1.PPCCLL — High profile, level from resolution
		level := avcLevel(height)
		return fmt.Sprintf("avc1.6400%02x", level)
	default:
		// Generic H.264 High 4.0 as safe fallback
		return "avc1.640028"
	}
}

// audioCodecString returns an RFC 6381 codec string for audio.
func audioCodecString(codec string) string {
	switch codec {
	case "aac":
		return "mp4a.40.2" // AAC-LC
	case "ac3":
		return "ac-3"
	case "eac3":
		return "ec-3"
	case "opus":
		return "Opus"
	case "flac":
		return "fLaC"
	case "mp3":
		return "mp4a.40.34"
	case "truehd":
		return "mlpa" // Dolby TrueHD
	case "dts":
		return "dtsc" // DTS Core
	case "dts_hd", "dtshd":
		return "dtsh" // DTS-HD
	case "dts_hd_ma", "dtshd_ma":
		return "dtsl" // DTS-HD Master Audio (lossless)
	case "dts_x", "dtsx":
		return "dtsx" // DTS:X
	default:
		return "mp4a.40.2" // AAC-LC fallback
	}
}

// hevcLevel returns the HEVC level_idc for the CODECS string based on resolution.
func hevcLevel(height int) int {
	switch {
	case height >= 2160:
		return 150 // Level 5.0
	case height >= 1440:
		return 120 // Level 4.0
	case height >= 1080:
		return 120 // Level 4.0
	case height >= 720:
		return 93 // Level 3.1
	default:
		return 90 // Level 3.0
	}
}

// av1Level returns the AV1 seq_level_idx for the CODECS string based on resolution.
func av1Level(height int) int {
	switch {
	case height >= 2160:
		return 13 // Level 5.1
	case height >= 1440:
		return 12 // Level 5.0
	case height >= 1080:
		return 9 // Level 4.1
	case height >= 720:
		return 8 // Level 4.0
	default:
		return 4 // Level 3.0
	}
}

// avcLevel returns the H.264 level (decimal) for the CODECS string based on resolution.
func avcLevel(height int) int {
	switch {
	case height >= 2160:
		return 51 // Level 5.1 (0x33)
	case height >= 1440:
		return 50 // Level 5.0 (0x32)
	case height >= 1080:
		return 40 // Level 4.0 (0x28)
	case height >= 720:
		return 31 // Level 3.1 (0x1F)
	default:
		return 30 // Level 3.0 (0x1E)
	}
}
