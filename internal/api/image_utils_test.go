package api

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestPNG creates a simple PNG image for testing
func createTestPNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a solid color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestDetectImageInfo(t *testing.T) {
	t.Run("valid PNG", func(t *testing.T) {
		imgData := createTestPNG(100, 50)
		contentType, width, height, err := detectImageInfo(bytes.NewReader(imgData))

		require.NoError(t, err)
		assert.Equal(t, "image/png", contentType)
		assert.Equal(t, 100, width)
		assert.Equal(t, 50, height)
	})

	t.Run("invalid data", func(t *testing.T) {
		_, _, _, err := detectImageInfo(bytes.NewReader([]byte("not an image")))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported content type")
	})

	t.Run("empty data", func(t *testing.T) {
		_, _, _, err := detectImageInfo(bytes.NewReader([]byte{}))
		assert.Error(t, err)
	})
}

func TestDetectImageInfoWithReader(t *testing.T) {
	t.Run("valid PNG returns new reader", func(t *testing.T) {
		imgData := createTestPNG(200, 150)
		contentType, width, height, newReader, err := detectImageInfoWithReader(bytes.NewReader(imgData))

		require.NoError(t, err)
		assert.Equal(t, "image/png", contentType)
		assert.Equal(t, 200, width)
		assert.Equal(t, 150, height)
		assert.NotNil(t, newReader)

		// Verify the new reader can be read
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(newReader)
		require.NoError(t, err)
		assert.Equal(t, imgData, buf.Bytes())
	})

	t.Run("invalid data", func(t *testing.T) {
		_, _, _, newReader, err := detectImageInfoWithReader(bytes.NewReader([]byte("invalid")))
		assert.Error(t, err)
		assert.Nil(t, newReader)
	})
}
