package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"movie.mp4", true},
		{"movie.mkv", true},
		{"movie.avi", true},
		{"movie.mov", true},
		{"movie.wmv", true},
		{"movie.webm", true},
		{"movie.m4v", true},
		{"movie.ts", true},
		{"movie.m2ts", true},
		{"movie.MP4", true},  // Case insensitive
		{"movie.MKV", true},  // Case insensitive
		{"movie.txt", false}, // Not a video
		{"movie.mp3", false}, // Audio, not video
		{"movie.jpg", false}, // Image, not video
		{"movie", false},     // No extension
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsVideoFile(tt.filename))
		})
	}
}

func TestIsAudioFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"song.mp3", true},
		{"song.flac", true},
		{"song.wav", true},
		{"song.aac", true},
		{"song.ogg", true},
		{"song.opus", true},
		{"song.m4a", true},
		{"song.MP3", true},  // Case insensitive
		{"song.FLAC", true}, // Case insensitive
		{"song.txt", false}, // Not audio
		{"song.mp4", false}, // Video, not audio
		{"song.jpg", false}, // Image, not audio
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsAudioFile(tt.filename))
		})
	}
}

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"photo.jpg", true},
		{"photo.jpeg", true},
		{"photo.png", true},
		{"photo.gif", true},
		{"photo.webp", true},
		{"photo.bmp", true},
		{"photo.JPG", true},  // Case insensitive
		{"photo.PNG", true},  // Case insensitive
		{"photo.txt", false}, // Not image
		{"photo.mp4", false}, // Video
		{"photo.mp3", false}, // Audio
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsImageFile(tt.filename))
		})
	}
}

func TestIsSubtitleFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"movie.srt", true},
		{"movie.ass", true},
		{"movie.ssa", true},
		{"movie.sub", true},
		{"movie.vtt", true},
		{"movie.SRT", true},  // Case insensitive
		{"movie.txt", false}, // Not subtitle
		{"movie.mp4", false}, // Video
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsSubtitleFile(tt.filename))
		})
	}
}

func TestHasExtension(t *testing.T) {
	customExtensions := map[string]bool{
		".custom": true,
		".test":   true,
	}

	tests := []struct {
		filename string
		expected bool
	}{
		{"file.custom", true},
		{"file.test", true},
		{"file.CUSTOM", true}, // Case insensitive
		{"file.other", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, HasExtension(tt.filename, customExtensions))
		})
	}
}

func TestExtensionsToSlice(t *testing.T) {
	extensions := map[string]bool{
		".mp4": true,
		".mkv": true,
		".avi": true,
	}

	slice := ExtensionsToSlice(extensions)
	assert.Len(t, slice, 3)
	assert.Contains(t, slice, ".mp4")
	assert.Contains(t, slice, ".mkv")
	assert.Contains(t, slice, ".avi")
}

func TestMergeExtensions(t *testing.T) {
	map1 := map[string]bool{".mp4": true, ".mkv": true}
	map2 := map[string]bool{".avi": true, ".mov": true}
	map3 := map[string]bool{".mkv": true, ".wmv": true} // .mkv overlaps

	merged := MergeExtensions(map1, map2, map3)

	assert.Len(t, merged, 5) // 6 unique extensions after dedup
	assert.True(t, merged[".mp4"])
	assert.True(t, merged[".mkv"])
	assert.True(t, merged[".avi"])
	assert.True(t, merged[".mov"])
	assert.True(t, merged[".wmv"])
}

func TestVideoExtensionsCompleteness(t *testing.T) {
	// Ensure common video formats are included
	expectedFormats := []string{
		".mp4", ".mkv", ".avi", ".mov", ".wmv",
		".flv", ".webm", ".m4v", ".mpg", ".mpeg",
		".ts", ".m2ts",
	}

	for _, ext := range expectedFormats {
		assert.True(t, VideoExtensions[ext], "Missing video extension: %s", ext)
	}
}

func TestAudioExtensionsCompleteness(t *testing.T) {
	// Ensure common audio formats are included
	expectedFormats := []string{
		".mp3", ".flac", ".wav", ".aac", ".ogg",
		".opus", ".m4a", ".wma",
	}

	for _, ext := range expectedFormats {
		assert.True(t, AudioExtensions[ext], "Missing audio extension: %s", ext)
	}
}
