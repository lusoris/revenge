package api

import (
	"bytes"
	"fmt"
	"image"
	// Register image format decoders
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"

	// WebP decoder
	_ "golang.org/x/image/webp"
)

// detectImageInfo reads an image and returns its content type and dimensions.
// It reads the file into memory to detect the format and dimensions,
// then wraps it in a new reader that can be read again for storage.
func detectImageInfo(r io.Reader) (contentType string, width, height int, err error) {
	// Read the entire file into memory for detection
	// (necessary because image.Decode consumes the reader)
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

	// Decode image to get dimensions
	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}

	return contentType, img.Width, img.Height, nil
}

// detectImageInfoWithReader reads an image and returns its content type, dimensions,
// and a new reader containing the original data (since the original reader was consumed).
func detectImageInfoWithReader(r io.Reader) (contentType string, width, height int, newReader io.Reader, err error) {
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

	// Decode image to get dimensions
	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Return a new reader for the data
	return contentType, img.Width, img.Height, bytes.NewReader(data), nil
}
