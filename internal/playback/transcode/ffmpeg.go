package transcode

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
)

// BuildVideoOnlyCommand builds an FFmpeg command that outputs video-only HLS segments.
// Audio is excluded (-an) so each audio track can be served as a separate rendition.
func BuildVideoOnlyCommand(ffmpegPath, inputFile, segmentDir string, pd ProfileDecision, segmentDuration, seekSeconds int) *exec.Cmd {
	outputPlaylist := filepath.Join(segmentDir, "index.m3u8")
	segmentPattern := filepath.Join(segmentDir, "seg-%05d.ts")

	args := []string{"-hide_banner", "-loglevel", "warning"}

	// Input seeking (fast, before -i)
	if seekSeconds > 0 {
		args = append(args, "-ss", strconv.Itoa(seekSeconds))
	}

	args = append(args, "-i", inputFile)

	// Map video only
	args = append(args, "-map", "0:v:0", "-an")

	// Video codec
	if pd.VideoCodec == "copy" {
		args = append(args, "-c:v", "copy")
	} else {
		args = append(args,
			"-c:v", pd.VideoCodec,
			"-preset", "veryfast",
			"-crf", "23",
		)
		if pd.Height > 0 {
			args = append(args, "-vf", fmt.Sprintf("scale=-2:%d", pd.Height))
		}
		if pd.VideoBitrate > 0 {
			args = append(args,
				"-maxrate", fmt.Sprintf("%dk", pd.VideoBitrate),
				"-bufsize", fmt.Sprintf("%dk", pd.VideoBitrate*2),
			)
		}
	}

	// HLS output
	args = append(args,
		"-f", "hls",
		"-hls_time", strconv.Itoa(segmentDuration),
		"-hls_playlist_type", "event",
		"-hls_segment_filename", segmentPattern,
		"-start_number", "0",
		outputPlaylist,
	)

	return exec.Command(ffmpegPath, args...)
}

// BuildAudioRenditionCommand builds an FFmpeg command that outputs audio-only HLS segments
// for a single audio track. Each track is a separate rendition that HLS.js can switch
// between instantly without restarting the video stream.
func BuildAudioRenditionCommand(ffmpegPath, inputFile, segmentDir string, trackIndex int, codec string, bitrate, segmentDuration, seekSeconds int) *exec.Cmd {
	outputPlaylist := filepath.Join(segmentDir, "index.m3u8")
	segmentPattern := filepath.Join(segmentDir, "seg-%05d.ts")

	args := []string{"-hide_banner", "-loglevel", "warning"}

	if seekSeconds > 0 {
		args = append(args, "-ss", strconv.Itoa(seekSeconds))
	}

	args = append(args, "-i", inputFile)

	// Map single audio track, no video
	args = append(args, "-map", fmt.Sprintf("0:a:%d", trackIndex), "-vn")

	// Audio codec
	if codec == "copy" {
		args = append(args, "-c:a", "copy")
	} else {
		args = append(args, "-c:a", codec)
		if bitrate > 0 {
			args = append(args, "-b:a", fmt.Sprintf("%dk", bitrate))
		}
	}

	// HLS output
	args = append(args,
		"-f", "hls",
		"-hls_time", strconv.Itoa(segmentDuration),
		"-hls_playlist_type", "event",
		"-hls_segment_filename", segmentPattern,
		"-start_number", "0",
		outputPlaylist,
	)

	return exec.Command(ffmpegPath, args...)
}

// BuildSubtitleExtractCommand builds an FFmpeg command to extract subtitles to WebVTT.
func BuildSubtitleExtractCommand(ffmpegPath, inputFile, outputFile string, trackIndex int) *exec.Cmd {
	args := []string{
		"-hide_banner", "-loglevel", "warning",
		"-i", inputFile,
		"-map", fmt.Sprintf("0:s:%d", trackIndex),
		"-f", "webvtt",
		"-y",
		outputFile,
	}

	return exec.Command(ffmpegPath, args...)
}

// CommandArgs returns the argument list for an exec.Cmd (for testing/logging).
func CommandArgs(cmd *exec.Cmd) []string {
	return cmd.Args
}
