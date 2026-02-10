// Package subtitle provides subtitle extraction from media files to WebVTT format.
package subtitle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// ExtractToWebVTT extracts a subtitle track from a media file to WebVTT format.
// trackIndex is the subtitle stream index (0-based within subtitle streams).
func ExtractToWebVTT(ctx context.Context, ffmpegPath, inputFile, outputDir string, trackIndex int) (string, error) {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create subtitle output dir: %w", err)
	}

	outputFile := filepath.Join(outputDir, strconv.Itoa(trackIndex)+".vtt")

	args := []string{
		"-hide_banner", "-loglevel", "warning",
		"-i", inputFile,
		"-map", fmt.Sprintf("0:s:%d", trackIndex),
		"-f", "webvtt",
		"-y",
		outputFile,
	}

	cmd := exec.CommandContext(ctx, ffmpegPath, args...) // #nosec G204 -- ffmpegPath is from validated config
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("subtitle extraction failed for track %d: %w (output: %s)", trackIndex, err, string(output))
	}

	return outputFile, nil
}
