// Package subtitle provides subtitle extraction from media files to WebVTT format.
package subtitle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/asticode/go-astiav"
)

// vttCue represents a single WebVTT cue (subtitle event).
type vttCue struct {
	startMs int64
	endMs   int64
	text    string
}

// assTagRegex strips ASS/SSA override tags like {\b1}, {\i0}, {\pos(x,y)}, etc.
var assTagRegex = regexp.MustCompile(`\{\\[^}]*\}`)

// ExtractToWebVTT extracts a subtitle track from a media file and converts it
// to WebVTT format. Supports text-based subtitle codecs: SRT (SubRip),
// ASS/SSA (Advanced SubStation Alpha), and WebVTT passthrough.
//
// The extraction reads raw demuxed packets and converts them to WebVTT cues
// directly, bypassing the FFmpeg WebVTT muxer (which rejects non-WebVTT input
// codecs). This approach is faster and more reliable for text subtitles.
//
// trackIndex is the subtitle stream index (0-based relative to subtitle streams).
func ExtractToWebVTT(ctx context.Context, inputFile, outputDir string, trackIndex int) (string, error) {
	subsDir := filepath.Join(outputDir, "subs")
	if err := os.MkdirAll(subsDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create subtitle output dir: %w", err)
	}

	outputFile := filepath.Join(subsDir, strconv.Itoa(trackIndex)+".vtt")

	// Open input
	inputFmtCtx := astiav.AllocFormatContext()
	if inputFmtCtx == nil {
		return "", errors.New("failed to allocate input format context")
	}
	defer inputFmtCtx.Free()

	if err := inputFmtCtx.OpenInput(inputFile, nil, nil); err != nil {
		return "", fmt.Errorf("failed to open input %q: %w", inputFile, err)
	}
	defer inputFmtCtx.CloseInput()

	if err := inputFmtCtx.FindStreamInfo(nil); err != nil {
		return "", fmt.Errorf("failed to find stream info: %w", err)
	}

	// Find the target subtitle stream by relative index within subtitle streams.
	var subStream *astiav.Stream
	subIdx := 0
	for _, s := range inputFmtCtx.Streams() {
		if s.CodecParameters().MediaType() == astiav.MediaTypeSubtitle {
			if subIdx == trackIndex {
				subStream = s
				break
			}
			subIdx++
		}
	}
	if subStream == nil {
		return "", fmt.Errorf("subtitle track %d not found", trackIndex)
	}

	codecID := subStream.CodecParameters().CodecID()

	// Determine if the source is ASS/SSA (needs dialogue field extraction).
	isASS := codecID == astiav.CodecIDAss || codecID == astiav.CodecIDSsa

	// Read all subtitle packets and convert to WebVTT cues.
	pkt := astiav.AllocPacket()
	if pkt == nil {
		return "", errors.New("failed to allocate packet")
	}
	defer pkt.Free()

	timeBase := subStream.TimeBase()
	var cues []vttCue

	for {
		if ctx.Err() != nil {
			break
		}

		if err := inputFmtCtx.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			return "", fmt.Errorf("failed to read frame: %w", err)
		}

		// Only process packets from our subtitle stream.
		if pkt.StreamIndex() != subStream.Index() {
			pkt.Unref()
			continue
		}

		// Extract timing from packet PTS and duration.
		pts := pkt.Pts()
		dur := pkt.Duration()
		if pts == astiav.NoPtsValue {
			pkt.Unref()
			continue
		}

		startMs := ptsToMillis(pts, timeBase)
		endMs := startMs
		if dur > 0 {
			endMs = ptsToMillis(pts+dur, timeBase)
		}

		// Extract text from packet data.
		data := string(pkt.Data())

		var text string
		if isASS {
			text = extractASSText(data)
		} else {
			// SRT, WebVTT, and other text formats: raw data is the text.
			text = data
		}

		text = strings.TrimSpace(text)
		if text != "" {
			cues = append(cues, vttCue{
				startMs: startMs,
				endMs:   endMs,
				text:    text,
			})
		}

		pkt.Unref()
	}

	// Write WebVTT file directly (no FFmpeg muxer needed).
	f, err := os.Create(outputFile) //nolint:gosec // path is constructed internally
	if err != nil {
		return "", fmt.Errorf("failed to create WebVTT file: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprint(f, "WEBVTT\n\n"); err != nil {
		return "", fmt.Errorf("failed to write WebVTT header: %w", err)
	}

	for i, cue := range cues {
		if _, err := fmt.Fprintf(f, "%d\n%s --> %s\n%s\n\n",
			i+1,
			formatVTTTime(cue.startMs),
			formatVTTTime(cue.endMs),
			cue.text,
		); err != nil {
			return "", fmt.Errorf("failed to write cue %d: %w", i+1, err)
		}
	}

	return outputFile, nil
}

// ptsToMillis converts a PTS value with a given time base to milliseconds.
func ptsToMillis(pts int64, tb astiav.Rational) int64 {
	return pts * int64(tb.Num()) * 1000 / int64(tb.Den())
}

// formatVTTTime formats milliseconds as a WebVTT timestamp: HH:MM:SS.mmm
func formatVTTTime(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	h := ms / 3600000
	ms %= 3600000
	m := ms / 60000
	ms %= 60000
	s := ms / 1000
	ms %= 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

// extractASSText extracts plain text from an ASS/SSA dialogue event line.
// ASS event packet format: ReadOrder,Layer,Style,Name,MarginL,MarginR,MarginV,Effect,Text
func extractASSText(data string) string {
	// The Text field is everything after the 8th comma.
	parts := strings.SplitN(data, ",", 9)
	if len(parts) < 9 {
		// Fallback: strip tags from the whole string.
		return stripASSTags(data)
	}
	return stripASSTags(parts[8])
}

// stripASSTags removes ASS override tags and converts ASS newlines to real newlines.
func stripASSTags(s string) string {
	s = assTagRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\\N", "\n")
	s = strings.ReplaceAll(s, "\\n", "\n")
	return strings.TrimSpace(s)
}
