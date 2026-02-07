package transcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVideoOnlyCommand_Copy(t *testing.T) {
	pd := ProfileDecision{
		Name:       "original",
		VideoCodec: "copy",
	}

	cmd := BuildVideoOnlyCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/original", pd, 6, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "-i")
	assert.Contains(t, args, "/media/movie.mkv")
	assert.Contains(t, args, "0:v:0")
	assert.Contains(t, args, "-an", "should exclude audio")
	assert.Contains(t, args, "-c:v")
	assert.Contains(t, args, "copy")
	assert.Contains(t, args, "-f")
	assert.Contains(t, args, "hls")
	assert.Contains(t, args, "-hls_time")
	assert.Contains(t, args, "6")

	// Should NOT contain any audio codec flags
	assert.NotContains(t, args, "-c:a")
	assert.NotContains(t, args, "-b:a")
}

func TestBuildVideoOnlyCommand_Transcode(t *testing.T) {
	pd := ProfileDecision{
		Name:         "720p",
		Width:        1280,
		Height:       720,
		VideoBitrate: 2800,
		VideoCodec:   "libx264",
	}

	cmd := BuildVideoOnlyCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/720p", pd, 6, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "-c:v")
	assert.Contains(t, args, "libx264")
	assert.Contains(t, args, "-preset")
	assert.Contains(t, args, "veryfast")
	assert.Contains(t, args, "-vf")
	assert.Contains(t, args, "scale=-2:720")
	assert.Contains(t, args, "-an", "should exclude audio")
	assert.NotContains(t, args, "-c:a")
}

func TestBuildVideoOnlyCommand_WithSeek(t *testing.T) {
	pd := ProfileDecision{
		Name:       "original",
		VideoCodec: "copy",
	}

	cmd := BuildVideoOnlyCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/original", pd, 6, 120)
	args := CommandArgs(cmd)

	ssIdx := -1
	iIdx := -1
	for i, arg := range args {
		if arg == "-ss" {
			ssIdx = i
		}
		if arg == "-i" {
			iIdx = i
		}
	}

	assert.Greater(t, ssIdx, -1, "-ss should be present")
	assert.Greater(t, iIdx, -1, "-i should be present")
	assert.Less(t, ssIdx, iIdx, "-ss should come before -i for fast seeking")
	assert.Equal(t, "120", args[ssIdx+1])
}

func TestBuildAudioRenditionCommand_Copy(t *testing.T) {
	cmd := BuildAudioRenditionCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/audio/0", 0, "copy", 0, 6, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "-i")
	assert.Contains(t, args, "0:a:0")
	assert.Contains(t, args, "-vn", "should exclude video")
	assert.Contains(t, args, "-c:a")
	assert.Contains(t, args, "copy")
	assert.Contains(t, args, "-f")
	assert.Contains(t, args, "hls")

	// Should NOT contain any video codec flags
	assert.NotContains(t, args, "-c:v")
	assert.NotContains(t, args, "-vf")
}

func TestBuildAudioRenditionCommand_Transcode(t *testing.T) {
	cmd := BuildAudioRenditionCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/audio/1", 1, "aac", 256, 6, 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "0:a:1")
	assert.Contains(t, args, "-vn")
	assert.Contains(t, args, "-c:a")
	assert.Contains(t, args, "aac")
	assert.Contains(t, args, "-b:a")
	assert.Contains(t, args, "256k")
}

func TestBuildAudioRenditionCommand_WithSeek(t *testing.T) {
	cmd := BuildAudioRenditionCommand("ffmpeg", "/media/movie.mkv", "/tmp/segments/audio/0", 0, "copy", 0, 6, 300)
	args := CommandArgs(cmd)

	ssIdx := -1
	iIdx := -1
	for i, arg := range args {
		if arg == "-ss" {
			ssIdx = i
		}
		if arg == "-i" {
			iIdx = i
		}
	}

	assert.Greater(t, ssIdx, -1, "-ss should be present")
	assert.Less(t, ssIdx, iIdx, "-ss should come before -i")
	assert.Equal(t, "300", args[ssIdx+1])
}

func TestBuildSubtitleExtractCommand(t *testing.T) {
	cmd := BuildSubtitleExtractCommand("ffmpeg", "/media/movie.mkv", "/tmp/subs/0.vtt", 0)
	args := CommandArgs(cmd)

	assert.Contains(t, args, "-i")
	assert.Contains(t, args, "/media/movie.mkv")
	assert.Contains(t, args, "-map")
	assert.Contains(t, args, "0:s:0")
	assert.Contains(t, args, "-f")
	assert.Contains(t, args, "webvtt")
	assert.Contains(t, args, "/tmp/subs/0.vtt")
}
