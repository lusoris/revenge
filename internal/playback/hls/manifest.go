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
// It references media playlists for each quality profile, audio renditions, and subtitle tracks.
// All audio tracks are muxed into segments so HLS.js can switch instantly without stream restart.
func GenerateMasterPlaylist(profiles []ProfileVariant, audioTracks []AudioVariant, subtitles []SubtitleVariant) string {
	var b strings.Builder

	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:3\n")
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

	// Stream variants
	for _, p := range profiles {
		extraAttrs := ""
		if len(audioTracks) > 0 {
			extraAttrs += ",AUDIO=\"audio\""
		}
		if len(subtitles) > 0 {
			extraAttrs += ",SUBTITLES=\"subs\""
		}

		fmt.Fprintf(&b, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d,NAME=\"%s\"%s\n",
			p.Bandwidth, p.Width, p.Height, p.Name, extraAttrs)
		fmt.Fprintf(&b, "%s/index.m3u8\n", p.Name)
	}

	return b.String()
}

// ProfileVariant describes a quality variant in the master playlist.
type ProfileVariant struct {
	Name      string
	Width     int
	Height    int
	Bandwidth int // bits per second
}

// AudioVariant describes an audio rendition in the master playlist.
// All audio tracks are muxed into segments, enabling instant track switching.
type AudioVariant struct {
	Index     int
	Name      string
	Language  string
	Channels  int
	IsDefault bool
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
	// Default estimates if both are zero
	if videoBps == 0 {
		videoBps = 5000000 // 5 Mbps default
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
	for _, part := range strings.Split(profile, "/") {
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
	// Use the validated absolute path for file operations
	var content []byte
	for i := 0; i < 100; i++ {
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
