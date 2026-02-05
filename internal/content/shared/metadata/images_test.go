package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImageURLBuilder(t *testing.T) {
	builder := NewImageURLBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, TMDbImageBaseURL, builder.baseURL)
}

func TestNewImageURLBuilderWithBase(t *testing.T) {
	customBase := "https://custom.example.com/images"
	builder := NewImageURLBuilderWithBase(customBase)
	assert.NotNil(t, builder)
	assert.Equal(t, customBase, builder.baseURL)
}

func TestImageURLBuilderGetURL(t *testing.T) {
	builder := NewImageURLBuilder()

	tests := []struct {
		name     string
		path     string
		size     string
		expected string
	}{
		{
			name:     "poster path",
			path:     "/abc123.jpg",
			size:     "w500",
			expected: "https://image.tmdb.org/t/p/w500/abc123.jpg",
		},
		{
			name:     "backdrop path",
			path:     "/backdrop.jpg",
			size:     "w1280",
			expected: "https://image.tmdb.org/t/p/w1280/backdrop.jpg",
		},
		{
			name:     "original size",
			path:     "/image.png",
			size:     "original",
			expected: "https://image.tmdb.org/t/p/original/image.png",
		},
		{
			name:     "empty path",
			path:     "",
			size:     "w500",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.GetURL(tt.path, tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestImageURLBuilderGetPosterURL(t *testing.T) {
	builder := NewImageURLBuilder()

	t.Run("with valid path", func(t *testing.T) {
		path := "/poster.jpg"
		result := builder.GetPosterURL(&path, "w342")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w342/poster.jpg", *result)
	})

	t.Run("with default size", func(t *testing.T) {
		path := "/poster.jpg"
		result := builder.GetPosterURL(&path, "")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w500/poster.jpg", *result) // Default is w500
	})

	t.Run("with nil path", func(t *testing.T) {
		result := builder.GetPosterURL(nil, "w500")
		assert.Nil(t, result)
	})

	t.Run("with empty path", func(t *testing.T) {
		path := ""
		result := builder.GetPosterURL(&path, "w500")
		assert.Nil(t, result)
	})
}

func TestImageURLBuilderGetBackdropURL(t *testing.T) {
	builder := NewImageURLBuilder()

	t.Run("with valid path", func(t *testing.T) {
		path := "/backdrop.jpg"
		result := builder.GetBackdropURL(&path, "w780")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w780/backdrop.jpg", *result)
	})

	t.Run("with default size", func(t *testing.T) {
		path := "/backdrop.jpg"
		result := builder.GetBackdropURL(&path, "")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w1280/backdrop.jpg", *result) // Default is w1280
	})

	t.Run("with nil path", func(t *testing.T) {
		result := builder.GetBackdropURL(nil, "w1280")
		assert.Nil(t, result)
	})
}

func TestImageURLBuilderGetProfileURL(t *testing.T) {
	builder := NewImageURLBuilder()

	t.Run("with valid path", func(t *testing.T) {
		path := "/profile.jpg"
		result := builder.GetProfileURL(&path, "w45")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w45/profile.jpg", *result)
	})

	t.Run("with default size", func(t *testing.T) {
		path := "/profile.jpg"
		result := builder.GetProfileURL(&path, "")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w185/profile.jpg", *result) // Default is w185
	})

	t.Run("with nil path", func(t *testing.T) {
		result := builder.GetProfileURL(nil, "w185")
		assert.Nil(t, result)
	})
}

func TestImageURLBuilderGetLogoURL(t *testing.T) {
	builder := NewImageURLBuilder()

	t.Run("with valid path", func(t *testing.T) {
		path := "/logo.png"
		result := builder.GetLogoURL(&path, "w300")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w300/logo.png", *result)
	})

	t.Run("with default size", func(t *testing.T) {
		path := "/logo.png"
		result := builder.GetLogoURL(&path, "")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w154/logo.png", *result) // Default is w154
	})

	t.Run("with nil path", func(t *testing.T) {
		result := builder.GetLogoURL(nil, "w154")
		assert.Nil(t, result)
	})
}

func TestImageURLBuilderGetStillURL(t *testing.T) {
	builder := NewImageURLBuilder()

	t.Run("with valid path", func(t *testing.T) {
		path := "/still.jpg"
		result := builder.GetStillURL(&path, "w185")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w185/still.jpg", *result)
	})

	t.Run("with default size", func(t *testing.T) {
		path := "/still.jpg"
		result := builder.GetStillURL(&path, "")

		assert.NotNil(t, result)
		assert.Equal(t, "https://image.tmdb.org/t/p/w300/still.jpg", *result) // Default is w300
	})

	t.Run("with nil path", func(t *testing.T) {
		result := builder.GetStillURL(nil, "w300")
		assert.Nil(t, result)
	})
}

func TestImageSizeConstants(t *testing.T) {
	// Verify the size constants are as expected (TMDb specific sizes)
	assert.Equal(t, "w92", PosterSizeSmall)
	assert.Equal(t, "w185", PosterSizeMedium)
	assert.Equal(t, "w342", PosterSizeLarge)
	assert.Equal(t, "w500", PosterSizeXLarge)
	assert.Equal(t, "w780", PosterSizeHuge)
	assert.Equal(t, "original", PosterSizeOrig)

	assert.Equal(t, "w300", BackdropSizeSmall)
	assert.Equal(t, "w780", BackdropSizeMedium)
	assert.Equal(t, "w1280", BackdropSizeLarge)
	assert.Equal(t, "original", BackdropSizeOrig)

	assert.Equal(t, "w45", ProfileSizeSmall)
	assert.Equal(t, "w185", ProfileSizeMedium)
	assert.Equal(t, "h632", ProfileSizeLarge)
	assert.Equal(t, "original", ProfileSizeOrig)
}

func TestNewImageDownloader(t *testing.T) {
	config := ClientConfig{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
	}
	client := NewBaseClient(config)
	downloader := NewImageDownloader(client)

	assert.NotNil(t, downloader)
	assert.NotNil(t, downloader.GetURLBuilder())
}
