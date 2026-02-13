package scanner

import (
	"path/filepath"
	"strings"
)

// VideoExtensions contains common video file extensions used across movies and TV shows.
var VideoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
	".m4v":  true,
	".mpg":  true,
	".mpeg": true,
	".3gp":  true,
	".ts":   true,
	".m2ts": true,
	".vob":  true, // DVD
	".ogv":  true, // Ogg Video
}

// AudioExtensions contains common audio file extensions used for music libraries.
var AudioExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".wav":  true,
	".aac":  true,
	".ogg":  true,
	".opus": true,
	".m4a":  true,
	".wma":  true,
	".alac": true,
	".aiff": true,
	".ape":  true,
	".dsf":  true, // DSD
	".dff":  true, // DSD
}

// ImageExtensions contains common image extensions for artwork and photos.
var ImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".tiff": true,
	".svg":  true,
}

// SubtitleExtensions contains subtitle file extensions.
var SubtitleExtensions = map[string]bool{
	".srt": true,
	".ass": true,
	".ssa": true,
	".sub": true,
	".idx": true,
	".vtt": true,
	".smi": true,
	".pgs": true,
	".sup": true,
}

// IsVideoFile checks if a filename has a video extension
func IsVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return VideoExtensions[ext]
}

// IsAudioFile checks if a filename has an audio extension
func IsAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return AudioExtensions[ext]
}

// IsImageFile checks if a filename has an image extension
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ImageExtensions[ext]
}

// IsSubtitleFile checks if a filename has a subtitle extension
func IsSubtitleFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return SubtitleExtensions[ext]
}

// HasExtension checks if a filename matches any of the given extensions
func HasExtension(filename string, extensions map[string]bool) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return extensions[ext]
}

// ExtensionsToSlice converts an extension map to a slice
func ExtensionsToSlice(extensions map[string]bool) []string {
	result := make([]string, 0, len(extensions))
	for ext := range extensions {
		result = append(result, ext)
	}
	return result
}

// MergeExtensions combines multiple extension maps into one
func MergeExtensions(maps ...map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for _, m := range maps {
		for ext := range m {
			result[ext] = true
		}
	}
	return result
}
