// Package fingerprint provides video fingerprinting services for scene identification.
// Used by QAR (adult content) module to generate oshash, pHash, and MD5 for StashDB matching.
package fingerprint

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Minimum file size for oshash (64KB * 2 = 128KB).
const minOshashSize = 65536 * 2

// Errors returned by the fingerprint service.
var (
	ErrFileTooSmall  = errors.New("file too small for fingerprinting")
	ErrFileNotFound  = errors.New("file not found")
	ErrFFProbeNotFound = errors.New("ffprobe not found in PATH")
	ErrInvalidVideo  = errors.New("invalid or corrupt video file")
)

// Result contains the generated fingerprints for a video file.
// Compatible with voyage.FingerprintResult interface.
type Result struct {
	Coordinates string // pHash (perceptual hash)
	Oshash      string // OpenSubtitles hash
	MD5         string // MD5 file hash
	Duration    int    // Video duration in seconds
	Resolution  string // e.g., "1920x1080"
	Codec       string // Video codec
}

// Config contains configuration for the fingerprint service.
type Config struct {
	FFProbePath   string        // Path to ffprobe binary (default: search PATH)
	FFMpegPath    string        // Path to ffmpeg binary (default: search PATH)
	GeneratePHash bool          // Whether to generate pHash (slower, requires frame extraction)
	GenerateMD5   bool          // Whether to generate MD5 hash
	Timeout       time.Duration // Timeout for ffprobe/ffmpeg operations
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		FFProbePath:   "ffprobe",
		FFMpegPath:    "ffmpeg",
		GeneratePHash: true,
		GenerateMD5:   true,
		Timeout:       30 * time.Second,
	}
}

// Service provides video fingerprinting capabilities.
type Service struct {
	ffprobePath   string
	ffmpegPath    string
	generatePHash bool
	generateMD5   bool
	timeout       time.Duration
	logger        *slog.Logger
}

// NewService creates a new fingerprint service.
func NewService(cfg Config, logger *slog.Logger) *Service {
	if cfg.FFProbePath == "" {
		cfg.FFProbePath = "ffprobe"
	}
	if cfg.FFMpegPath == "" {
		cfg.FFMpegPath = "ffmpeg"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	return &Service{
		ffprobePath:   cfg.FFProbePath,
		ffmpegPath:    cfg.FFMpegPath,
		generatePHash: cfg.GeneratePHash,
		generateMD5:   cfg.GenerateMD5,
		timeout:       cfg.Timeout,
		logger:        logger.With("service", "fingerprint"),
	}
}

// GenerateFingerprints generates all fingerprints for a video file.
// Implements the voyage.Fingerprinter interface.
func (s *Service) GenerateFingerprints(ctx context.Context, filePath string) (*Result, error) {
	// Check file exists
	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	if stat.Size() < minOshashSize {
		return nil, ErrFileTooSmall
	}

	result := &Result{}

	// Always generate oshash (fast, file-based)
	oshash, err := s.calculateOshash(filePath, stat.Size())
	if err != nil {
		s.logger.Warn("failed to calculate oshash", "path", filePath, "error", err)
	} else {
		result.Oshash = oshash
	}

	// Get video metadata via ffprobe
	probe, err := s.probeFile(ctx, filePath)
	if err != nil {
		s.logger.Warn("failed to probe video", "path", filePath, "error", err)
	} else {
		result.Duration = probe.Duration
		result.Resolution = probe.Resolution
		result.Codec = probe.VideoCodec
	}

	// Generate pHash if enabled (slower, requires frame extraction)
	if s.generatePHash {
		phash, err := s.calculatePHash(ctx, filePath)
		if err != nil {
			s.logger.Warn("failed to calculate phash", "path", filePath, "error", err)
		} else {
			result.Coordinates = phash
		}
	}

	// Generate MD5 if enabled (can be slow for large files)
	if s.generateMD5 {
		md5hash, err := s.calculateMD5(filePath)
		if err != nil {
			s.logger.Warn("failed to calculate md5", "path", filePath, "error", err)
		} else {
			result.MD5 = md5hash
		}
	}

	// At minimum we need oshash to be useful
	if result.Oshash == "" {
		return nil, fmt.Errorf("failed to generate any fingerprints for %s", filePath)
	}

	s.logger.Debug("fingerprints generated",
		"path", filePath,
		"oshash", result.Oshash,
		"phash", result.Coordinates,
		"md5", result.MD5,
		"duration", result.Duration,
	)

	return result, nil
}

// calculateOshash implements the OpenSubtitles hash algorithm.
// Algorithm: hash = filesize + first 64KB + last 64KB (as uint64 sums).
func (s *Service) calculateOshash(filePath string, fileSize int64) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	// oshash = first 64KB + last 64KB + filesize
	buf := make([]byte, 65536)
	var hash uint64 = uint64(fileSize)

	// Read first 64KB
	n, err := io.ReadFull(file, buf)
	if err != nil {
		return "", fmt.Errorf("read first chunk: %w", err)
	}
	if n != 65536 {
		return "", ErrFileTooSmall
	}

	// Sum uint64 values from first chunk
	for i := 0; i < 65536; i += 8 {
		hash += binary.LittleEndian.Uint64(buf[i:])
	}

	// Seek to last 64KB
	_, err = file.Seek(-65536, io.SeekEnd)
	if err != nil {
		return "", fmt.Errorf("seek to end: %w", err)
	}

	// Read last 64KB
	n, err = io.ReadFull(file, buf)
	if err != nil {
		return "", fmt.Errorf("read last chunk: %w", err)
	}
	if n != 65536 {
		return "", ErrFileTooSmall
	}

	// Sum uint64 values from last chunk
	for i := 0; i < 65536; i += 8 {
		hash += binary.LittleEndian.Uint64(buf[i:])
	}

	return fmt.Sprintf("%016x", hash), nil
}

// ProbeResult contains video metadata from ffprobe.
type ProbeResult struct {
	Duration   int    // Duration in seconds
	Resolution string // e.g., "1920x1080"
	VideoCodec string // e.g., "h264"
	AudioCodec string // e.g., "aac"
	Width      int
	Height     int
	Bitrate    int // bits per second
}

// probeFile extracts video metadata using ffprobe.
func (s *Service) probeFile(ctx context.Context, filePath string) (*ProbeResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Check ffprobe is available
	ffprobePath, err := exec.LookPath(s.ffprobePath)
	if err != nil {
		return nil, ErrFFProbeNotFound
	}

	// Run ffprobe with JSON output
	cmd := exec.CommandContext(ctx, ffprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-select_streams", "v:0", // First video stream
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	return s.parseProbeOutput(output)
}

// parseProbeOutput parses ffprobe JSON output.
func (s *Service) parseProbeOutput(output []byte) (*ProbeResult, error) {
	// Simple parsing without full JSON unmarshaling
	// Extract duration, resolution, and codec
	result := &ProbeResult{}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Duration (format section)
		if strings.Contains(line, `"duration"`) {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				durStr := strings.Trim(parts[len(parts)-1], `", `)
				if dur, err := strconv.ParseFloat(durStr, 64); err == nil {
					result.Duration = int(dur)
				}
			}
		}

		// Width
		if strings.Contains(line, `"width"`) {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				widthStr := strings.Trim(parts[len(parts)-1], `", `)
				if w, err := strconv.Atoi(widthStr); err == nil {
					result.Width = w
				}
			}
		}

		// Height
		if strings.Contains(line, `"height"`) {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				heightStr := strings.Trim(parts[len(parts)-1], `", `)
				if h, err := strconv.Atoi(heightStr); err == nil {
					result.Height = h
				}
			}
		}

		// Video codec
		if strings.Contains(line, `"codec_name"`) {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				result.VideoCodec = strings.Trim(parts[len(parts)-1], `", `)
			}
		}

		// Bitrate
		if strings.Contains(line, `"bit_rate"`) {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				brStr := strings.Trim(parts[len(parts)-1], `", `)
				if br, err := strconv.Atoi(brStr); err == nil {
					result.Bitrate = br
				}
			}
		}
	}

	if result.Width > 0 && result.Height > 0 {
		result.Resolution = fmt.Sprintf("%dx%d", result.Width, result.Height)
	}

	return result, nil
}

// calculatePHash generates a perceptual hash for StashDB matching.
// This uses ffmpeg to extract frames and compute a hash.
// Note: Full pHash implementation requires image processing libraries.
// This is a simplified version that extracts a thumbnail and uses a basic hash.
func (s *Service) calculatePHash(ctx context.Context, filePath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout*2) // pHash takes longer
	defer cancel()

	// Check ffmpeg is available
	ffmpegPath, err := exec.LookPath(s.ffmpegPath)
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found: %w", err)
	}

	// Extract a frame at 10% into the video
	// Use video filter to resize to 8x8 and convert to grayscale
	// Then output raw pixel values
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-ss", "10", // Seek to 10 seconds
		"-i", filePath,
		"-vf", "thumbnail,scale=8:8,format=gray", // 8x8 grayscale
		"-frames:v", "1",
		"-f", "rawvideo",
		"-",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg frame extraction failed: %w", err)
	}

	if len(output) < 64 {
		return "", fmt.Errorf("insufficient frame data: got %d bytes", len(output))
	}

	// Calculate average value
	var sum int
	for _, b := range output[:64] {
		sum += int(b)
	}
	avg := byte(sum / 64)

	// Generate hash: 1 if pixel > average, 0 otherwise
	var hash uint64
	for i := 0; i < 64 && i < len(output); i++ {
		if output[i] > avg {
			hash |= 1 << uint(63-i)
		}
	}

	return fmt.Sprintf("%016x", hash), nil
}

// calculateMD5 calculates the MD5 hash of a file.
func (s *Service) calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// IsAvailable checks if the service can generate fingerprints.
func (s *Service) IsAvailable() bool {
	_, err := exec.LookPath(s.ffprobePath)
	return err == nil
}
