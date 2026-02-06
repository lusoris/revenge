package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
)

var vipsOnce sync.Once

// ensureVips initializes libvips on first use.
func ensureVips() {
	vipsOnce.Do(func() {
		vips.LoggingSettings(nil, vips.LogLevelError)
		vips.Startup(nil)
	})
}

// detectImageInfo reads an image and returns its content type and dimensions.
// It reads the file into memory to detect the format and dimensions,
// then wraps it in a new reader that can be read again for storage.
func detectImageInfo(r io.Reader) (contentType string, width, height int, err error) {
	ensureVips()

	// Read the entire file into memory for detection
	data, err := io.ReadAll(r)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read image data: %w", err)
	}

	// Detect content type using http.DetectContentType
	contentType = http.DetectContentType(data)

	// Validate it's an image type
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		// Valid image types
	default:
		return "", 0, 0, fmt.Errorf("unsupported content type: %s", contentType)
	}

	// Decode image with libvips to get dimensions
	img, err := vips.NewImageFromBuffer(data)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}
	defer img.Close()

	return contentType, img.Width(), img.Height(), nil
}

// detectImageInfoWithReader reads an image and returns its content type, dimensions,
// and a new reader containing the original data (since the original reader was consumed).
func detectImageInfoWithReader(r io.Reader) (contentType string, width, height int, newReader io.Reader, err error) {
	ensureVips()

	// Read the entire file into memory for detection
	data, err := io.ReadAll(r)
	if err != nil {
		return "", 0, 0, nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Detect content type using http.DetectContentType
	contentType = http.DetectContentType(data)

	// Validate it's an image type
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		// Valid image types
	default:
		return "", 0, 0, nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	// Decode image with libvips to get dimensions
	img, err := vips.NewImageFromBuffer(data)
	if err != nil {
		return "", 0, 0, nil, fmt.Errorf("failed to decode image: %w", err)
	}
	defer img.Close()

	// Return a new reader for the data
	return contentType, img.Width(), img.Height(), bytes.NewReader(data), nil
}
