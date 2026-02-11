package metadata

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

// TMDb image URL constants
const (
	// TMDbImageBaseURL is the base URL for TMDb images.
	TMDbImageBaseURL = "https://image.tmdb.org/t/p"
)

// Image size constants for TMDb
const (
	// Poster sizes
	PosterSizeSmall  = "w92"
	PosterSizeMedium = "w185"
	PosterSizeLarge  = "w342"
	PosterSizeXLarge = "w500"
	PosterSizeHuge   = "w780"
	PosterSizeOrig   = "original"

	// Backdrop sizes
	BackdropSizeSmall  = "w300"
	BackdropSizeMedium = "w780"
	BackdropSizeLarge  = "w1280"
	BackdropSizeOrig   = "original"

	// Profile (person) sizes
	ProfileSizeSmall  = "w45"
	ProfileSizeMedium = "w185"
	ProfileSizeLarge  = "h632"
	ProfileSizeOrig   = "original"

	// Logo sizes
	LogoSizeSmall  = "w45"
	LogoSizeMedium = "w154"
	LogoSizeLarge  = "w300"
	LogoSizeXLarge = "w500"
	LogoSizeOrig   = "original"

	// Still (episode) sizes
	StillSizeSmall  = "w92"
	StillSizeMedium = "w185"
	StillSizeLarge  = "w300"
	StillSizeOrig   = "original"
)

// ImageURLBuilder constructs full image URLs from paths.
type ImageURLBuilder struct {
	baseURL string
}

// NewImageURLBuilder creates a new ImageURLBuilder with the default TMDb image base URL.
func NewImageURLBuilder() *ImageURLBuilder {
	return &ImageURLBuilder{
		baseURL: TMDbImageBaseURL,
	}
}

// NewImageURLBuilderWithBase creates a new ImageURLBuilder with a custom base URL.
func NewImageURLBuilderWithBase(baseURL string) *ImageURLBuilder {
	return &ImageURLBuilder{
		baseURL: baseURL,
	}
}

// isFullURL returns true if the path is already a complete URL.
func isFullURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// GetURL constructs a full image URL from a path and size.
// If the path is already a full URL (from providers like fanart.tv, Radarr, etc.),
// it is returned unchanged. Returns empty string if path is empty.
func (b *ImageURLBuilder) GetURL(path string, size string) string {
	if path == "" {
		return ""
	}
	// Full URLs from non-TMDb providers should be passed through as-is.
	if isFullURL(path) {
		return path
	}
	return fmt.Sprintf("%s/%s%s", b.baseURL, size, path)
}

// GetPosterURL returns a poster URL with the specified size.
// Uses w500 as default if size is empty.
func (b *ImageURLBuilder) GetPosterURL(path *string, size string) *string {
	if path == nil || *path == "" {
		return nil
	}
	if size == "" {
		size = PosterSizeXLarge
	}
	url := b.GetURL(*path, size)
	return &url
}

// GetBackdropURL returns a backdrop URL with the specified size.
// Uses w1280 as default if size is empty.
func (b *ImageURLBuilder) GetBackdropURL(path *string, size string) *string {
	if path == nil || *path == "" {
		return nil
	}
	if size == "" {
		size = BackdropSizeLarge
	}
	url := b.GetURL(*path, size)
	return &url
}

// GetProfileURL returns a profile image URL with the specified size.
// Uses w185 as default if size is empty.
func (b *ImageURLBuilder) GetProfileURL(path *string, size string) *string {
	if path == nil || *path == "" {
		return nil
	}
	if size == "" {
		size = ProfileSizeMedium
	}
	url := b.GetURL(*path, size)
	return &url
}

// GetLogoURL returns a logo URL with the specified size.
// Uses w154 as default if size is empty.
func (b *ImageURLBuilder) GetLogoURL(path *string, size string) *string {
	if path == nil || *path == "" {
		return nil
	}
	if size == "" {
		size = LogoSizeMedium
	}
	url := b.GetURL(*path, size)
	return &url
}

// GetStillURL returns an episode still URL with the specified size.
// Uses w300 as default if size is empty.
func (b *ImageURLBuilder) GetStillURL(path *string, size string) *string {
	if path == nil || *path == "" {
		return nil
	}
	if size == "" {
		size = StillSizeLarge
	}
	url := b.GetURL(*path, size)
	return &url
}

// ImageDownloader handles downloading images from metadata providers.
type ImageDownloader struct {
	client     *BaseClient
	imgClient  *req.Client
	urlBuilder *ImageURLBuilder
}

// NewImageDownloader creates a new ImageDownloader.
func NewImageDownloader(client *BaseClient) *ImageDownloader {
	// Use a dedicated HTTP client for image downloads (different host than API).
	// Configured once instead of creating a throwaway req.C() per request.
	imgClient := req.C().
		SetTimeout(30 * time.Second).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	return &ImageDownloader{
		client:     client,
		imgClient:  imgClient,
		urlBuilder: NewImageURLBuilder(),
	}
}

// Download downloads an image by path and size.
func (d *ImageDownloader) Download(ctx context.Context, path string, size string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("empty image path")
	}

	if err := d.client.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	url := d.urlBuilder.GetURL(path, size)

	resp, err := d.imgClient.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("download image: %w", err)
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("download image: status %d", resp.StatusCode)
	}

	return resp.Bytes(), nil
}

// DownloadPoster downloads a poster image.
func (d *ImageDownloader) DownloadPoster(ctx context.Context, path string, size string) ([]byte, error) {
	if size == "" {
		size = PosterSizeXLarge
	}
	return d.Download(ctx, path, size)
}

// DownloadBackdrop downloads a backdrop image.
func (d *ImageDownloader) DownloadBackdrop(ctx context.Context, path string, size string) ([]byte, error) {
	if size == "" {
		size = BackdropSizeLarge
	}
	return d.Download(ctx, path, size)
}

// DownloadProfile downloads a profile image.
func (d *ImageDownloader) DownloadProfile(ctx context.Context, path string, size string) ([]byte, error) {
	if size == "" {
		size = ProfileSizeMedium
	}
	return d.Download(ctx, path, size)
}

// GetURLBuilder returns the underlying URL builder for direct URL construction.
func (d *ImageDownloader) GetURLBuilder() *ImageURLBuilder {
	return d.urlBuilder
}
