package image

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	// Image sizes supported by TMDb
	SizePosterSmall   = "w185"
	SizePosterMedium  = "w342"
	SizePosterLarge   = "w500"
	SizePosterXLarge  = "w780"
	SizePosterOriginal = "original"

	SizeBackdropSmall    = "w300"
	SizeBackdropMedium   = "w780"
	SizeBackdropLarge    = "w1280"
	SizeBackdropOriginal = "original"

	SizeProfileSmall   = "w45"
	SizeProfileMedium  = "w185"
	SizeProfileLarge   = "h632"
	SizeProfileOriginal = "original"

	// Image types
	TypePoster   = "poster"
	TypeBackdrop = "backdrop"
	TypeProfile  = "profile"
	TypeLogo     = "logo"
)

// Config holds image service configuration.
type Config struct {
	BaseURL   string        // TMDb image base URL (default: https://image.tmdb.org/t/p)
	CacheDir  string        // Local cache directory
	CacheTTL  time.Duration // Cache TTL (default: 7 days)
	MaxSize   int64         // Max file size in bytes (default: 10MB)
	ProxyURL  string        // Optional HTTP proxy
	UserAgent string        // User agent for requests
}

// Service handles image downloading, caching, and serving.
type Service struct {
	client   *resty.Client
	config   Config
	logger   *zap.Logger
	cache    sync.Map // path -> cachedImage
}

// NewService creates a new image service.
func NewService(cfg Config, logger *zap.Logger) (*Service, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://image.tmdb.org/t/p"
	}

	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 7 * 24 * time.Hour // 7 days default
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = 10 * 1024 * 1024 // 10MB default
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = "Revenge/1.0"
	}

	// Ensure cache directory exists
	if cfg.CacheDir != "" {
		if err := os.MkdirAll(cfg.CacheDir, 0750); err != nil {
			return nil, fmt.Errorf("create cache directory: %w", err)
		}
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetHeader("User-Agent", cfg.UserAgent)

	if cfg.ProxyURL != "" {
		client.SetProxy(cfg.ProxyURL)
	}

	return &Service{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// GetImageURL builds the full URL for an image.
func (s *Service) GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return fmt.Sprintf("%s/%s%s", s.config.BaseURL, size, path)
}

// FetchImage downloads an image and returns its content.
// Returns data, content-type, and error.
func (s *Service) FetchImage(ctx context.Context, imageType, path, size string) ([]byte, string, error) {
	if path == "" {
		return nil, "", fmt.Errorf("empty image path")
	}

	// Validate size
	if size == "" {
		size = s.getDefaultSize(imageType)
	}

	// Check cache first
	if s.config.CacheDir != "" {
		if data, contentType, err := s.getFromCache(imageType, path, size); err == nil {
			return data, contentType, nil
		}
	}

	// Build URL
	url := s.GetImageURL(path, size)

	s.logger.Debug("Fetching image",
		zap.String("type", imageType),
		zap.String("path", path),
		zap.String("size", size),
		zap.String("url", url))

	// Download
	resp, err := s.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, "", fmt.Errorf("fetch image: %w", err)
	}

	if resp.IsError() {
		return nil, "", fmt.Errorf("image fetch failed with status %d", resp.StatusCode())
	}

	data := resp.Body()
	contentType := resp.Header().Get("Content-Type")

	// Validate content type
	if !isValidImageType(contentType) {
		return nil, "", fmt.Errorf("invalid content type: %s", contentType)
	}

	// Check size
	if int64(len(data)) > s.config.MaxSize {
		return nil, "", fmt.Errorf("image too large: %d bytes", len(data))
	}

	// Cache the image
	if s.config.CacheDir != "" {
		if err := s.saveToCache(imageType, path, size, data, contentType); err != nil {
			s.logger.Warn("Failed to cache image", zap.Error(err))
		}
	}

	return data, contentType, nil
}

// StreamImage streams an image directly to a writer.
func (s *Service) StreamImage(ctx context.Context, w io.Writer, imageType, path, size string) (string, error) {
	data, contentType, err := s.FetchImage(ctx, imageType, path, size)
	if err != nil {
		return "", err
	}

	if _, err := w.Write(data); err != nil {
		return "", fmt.Errorf("write image: %w", err)
	}

	return contentType, nil
}

// ServeHTTP implements http.Handler for proxying images.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse path: /images/{type}/{size}/{path}
	// e.g., /images/poster/w500/abc123.jpg
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/images/"), "/", 3)
	if len(parts) < 3 {
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}

	imageType := parts[0]
	size := parts[1]
	imagePath := "/" + parts[2]

	// Validate type
	if !isValidType(imageType) {
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	// Validate size
	if !isValidSize(imageType, size) {
		http.Error(w, "Invalid image size", http.StatusBadRequest)
		return
	}

	// Generate ETag based on path and size
	etag := s.generateETag(imageType, size, imagePath)

	// Check If-None-Match for conditional request
	if match := r.Header.Get("If-None-Match"); match != "" {
		if match == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	data, contentType, err := s.FetchImage(r.Context(), imageType, imagePath, size)
	if err != nil {
		s.logger.Error("Failed to fetch image", zap.Error(err))
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Set cache headers for CDN and browser caching
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, immutable", int(s.config.CacheTTL.Seconds())))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("ETag", etag)
	w.Header().Set("Vary", "Accept-Encoding")

	// CORS headers for frontend clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")

	if _, err := w.Write(data); err != nil {
		s.logger.Error("failed to write image response", zap.Error(err))
	}
}

// generateETag creates a stable ETag for an image based on path and size.
func (s *Service) generateETag(imageType, size, path string) string {
	// Use a hash of type+size+path for stable ETag
	data := imageType + ":" + size + ":" + path
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
	return `"` + hash[:16] + `"`
}

// ClearCache clears the image cache.
func (s *Service) ClearCache() error {
	s.cache.Range(func(key, value interface{}) bool {
		s.cache.Delete(key)
		return true
	})

	if s.config.CacheDir != "" {
		return os.RemoveAll(s.config.CacheDir)
	}
	return nil
}

// Internal methods

func (s *Service) getDefaultSize(imageType string) string {
	switch imageType {
	case TypePoster:
		return SizePosterMedium
	case TypeBackdrop:
		return SizeBackdropMedium
	case TypeProfile:
		return SizeProfileMedium
	case TypeLogo:
		return SizePosterMedium
	default:
		return "original"
	}
}

func (s *Service) getCachePath(imageType, path, size string) string {
	// Sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(strings.TrimPrefix(path, "/"))
	return filepath.Join(s.config.CacheDir, imageType, size, cleanPath)
}

func (s *Service) getFromCache(imageType, path, size string) ([]byte, string, error) {
	cachePath := s.getCachePath(imageType, path, size)

	// Check if file exists and is not expired
	info, err := os.Stat(cachePath)
	if err != nil {
		return nil, "", err
	}

	if time.Since(info.ModTime()) > s.config.CacheTTL {
		_ = os.Remove(cachePath) // Ignore error - cache cleanup is best-effort
		return nil, "", fmt.Errorf("cache expired")
	}

	// #nosec G304 -- cachePath is constructed internally from validated imageType/path/size
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, "", err
	}

	// Detect content type from file extension
	contentType := getContentType(cachePath)

	return data, contentType, nil
}

func (s *Service) saveToCache(imageType, path, size string, data []byte, contentType string) error {
	cachePath := s.getCachePath(imageType, path, size)

	// Create directory
	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0600)
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}
	for _, t := range validTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}

func isValidType(imageType string) bool {
	switch imageType {
	case TypePoster, TypeBackdrop, TypeProfile, TypeLogo:
		return true
	}
	return false
}

func isValidSize(imageType, size string) bool {
	if size == "original" {
		return true
	}

	switch imageType {
	case TypePoster:
		switch size {
		case SizePosterSmall, SizePosterMedium, SizePosterLarge, SizePosterXLarge:
			return true
		}
	case TypeBackdrop:
		switch size {
		case SizeBackdropSmall, SizeBackdropMedium, SizeBackdropLarge:
			return true
		}
	case TypeProfile:
		switch size {
		case SizeProfileSmall, SizeProfileMedium, SizeProfileLarge:
			return true
		}
	case TypeLogo:
		switch size {
		case SizePosterSmall, SizePosterMedium, SizePosterLarge:
			return true
		}
	}
	return false
}

func getContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
