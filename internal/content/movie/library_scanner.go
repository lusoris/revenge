package movie

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie/adapters"
	"github.com/lusoris/revenge/internal/content/shared/scanner"
	"github.com/lusoris/revenge/internal/util"
)

// VideoExtensions are supported video file extensions.
// Deprecated: Use scanner.VideoExtensions from the shared package.
var VideoExtensions = scanner.VideoExtensions

// isVideoFile checks if a filename has a video extension
func isVideoFile(filename string) bool {
	return scanner.IsVideoFile(filename)
}

// Scanner handles file system scanning for movie files.
// It uses the shared FilesystemScanner internally with a MovieFileParser.
type Scanner struct {
	internal *scanner.FilesystemScanner
	paths    []string
}

// ScanResult represents a discovered movie file.
// This is a movie-specific type that wraps the shared scanner.ScanResult.
type ScanResult struct {
	FilePath    string
	FileName    string
	ParsedTitle string
	ParsedYear  *int
	FileSize    int64
	IsVideo     bool
	Error       error
}

// NewScanner creates a new library scanner using the shared scanner framework
func NewScanner(libraryPaths []string) *Scanner {
	parser := adapters.NewMovieFileParser()
	internal := scanner.NewFilesystemScanner(libraryPaths, parser)

	return &Scanner{
		internal: internal,
		paths:    libraryPaths,
	}
}

// Scan scans all library paths for movie files
func (s *Scanner) Scan(ctx context.Context) ([]ScanResult, error) {
	// Use the shared scanner
	sharedResults, err := s.internal.Scan(ctx)
	if err != nil {
		return nil, err
	}

	// Convert shared results to movie-specific results
	results := make([]ScanResult, 0, len(sharedResults))
	for _, sr := range sharedResults {
		results = append(results, convertScanResult(sr))
	}

	return results, nil
}

// convertScanResult converts a shared ScanResult to a movie ScanResult
func convertScanResult(sr scanner.ScanResult) ScanResult {
	return ScanResult{
		FilePath:    sr.FilePath,
		FileName:    sr.FileName,
		ParsedTitle: sr.ParsedTitle,
		ParsedYear:  sr.GetYear(),
		FileSize:    sr.FileSize,
		IsVideo:     sr.IsMedia,
		Error:       sr.Error,
	}
}

// parseMovieFilename attempts to extract title and year from filename.
// This function delegates to the MovieFileParser for consistent behavior.
// Examples:
// - "The Matrix (1999).mkv" -> "The Matrix", 1999
// - "Inception.2010.1080p.BluRay.mkv" -> "Inception", 2010
// - "The.Dark.Knight.2008.mkv" -> "The Dark Knight", 2008
func parseMovieFilename(filename string) (title string, year *int) {
	parser := adapters.NewMovieFileParser()
	title, metadata := parser.Parse(filename)
	if metadata != nil {
		if y, ok := metadata["year"]; ok {
			if yearInt, ok := y.(int); ok {
				year = &yearInt
			}
		}
	}
	return title, year
}

// cleanTitle cleans up the title string.
// Deprecated: Use scanner.CleanTitle from the shared package.
func cleanTitle(title string) string {
	return scanner.CleanTitle(title)
}

// MovieFileInfo contains detailed information about a movie file
type MovieFileInfo struct {
	Path            string
	Size            int64
	Container       string
	Resolution      string
	ResolutionLabel string
	VideoCodec      string
	VideoProfile    string
	AudioCodec      string
	Languages       []string
	BitrateKbps     int32
	DurationSeconds float64
	Framerate       float64
	DynamicRange    string
	ColorSpace      string
	AudioChannels   int
	AudioLayout     string
	SubtitleLangs   []string
}

// ExtractFileInfo extracts technical details from a video file using go-astiav (FFmpeg)
func ExtractFileInfo(filePath string) (*MovieFileInfo, error) {
	prober := NewMediaInfoProber()
	mediaInfo, err := prober.Probe(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to probe file: %w", err)
	}

	info := &MovieFileInfo{
		Path:            mediaInfo.FilePath,
		Size:            mediaInfo.FileSize,
		Container:       mediaInfo.Container,
		Resolution:      mediaInfo.Resolution,
		ResolutionLabel: mediaInfo.ResolutionLabel,
		VideoCodec:      mediaInfo.VideoCodec,
		VideoProfile:    mediaInfo.VideoProfile,
		BitrateKbps:     util.SafeInt64ToInt32(mediaInfo.BitrateKbps),
		DurationSeconds: mediaInfo.DurationSeconds,
		Framerate:       mediaInfo.Framerate,
		DynamicRange:    mediaInfo.DynamicRange,
		ColorSpace:      mediaInfo.ColorSpace,
		Languages:       mediaInfo.GetAudioLanguages(),
		SubtitleLangs:   mediaInfo.GetSubtitleLanguages(),
	}

	// Get file size from stat if not in mediainfo
	if info.Size == 0 {
		if fileInfo, err := os.Stat(filePath); err == nil {
			info.Size = fileInfo.Size()
		}
	}

	// Primary audio info
	if len(mediaInfo.AudioStreams) > 0 {
		primary := mediaInfo.AudioStreams[0]
		info.AudioCodec = primary.Codec
		info.AudioChannels = primary.Channels
		info.AudioLayout = primary.Layout
	}

	return info, nil
}

// CreateMovieFile creates a domain MovieFile from file info
func CreateMovieFile(movieID uuid.UUID, info *MovieFileInfo) *MovieFile {
	return &MovieFile{
		ID:         uuid.New(),
		MovieID:    movieID,
		FilePath:   info.Path,
		FileSize:   info.Size,
		FileName:   filepath.Base(info.Path),
		Container:  &info.Container,
		Resolution: parseOptionalString(info.Resolution),
		VideoCodec: parseOptionalString(info.VideoCodec),
		AudioCodec: parseOptionalString(info.AudioCodec),
	}
}

func parseOptionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
