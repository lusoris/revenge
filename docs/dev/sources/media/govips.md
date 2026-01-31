# GoVips Package

> Source: https://pkg.go.dev/github.com/davidbyttow/govips/v2/vips
> Fetched: 2026-01-31
> Content-Hash: auto-generated
> Type: html

---

## Overview

The `vips` package provides Go bindings for **libvips**, a fast, low-memory image processing library.

**Module:** `github.com/davidbyttow/govips/v2`
**Version:** v2.16.0
**License:** MIT

## Key Features

- **Fast image processing** with minimal memory footprint
- Support for multiple image formats (JPEG, PNG, WebP, GIF, TIFF, HEIF, AVIF, JP2K, JXL, etc.)
- Comprehensive image manipulation operations
- ICC color profile management
- EXIF metadata handling
- Thread-safe operations with configurable concurrency

## Initialization

```go
import "github.com/davidbyttow/govips/v2/vips"

func init() {
    vips.Startup(&vips.Config{
        ConcurrencyLevel: 4,
        MaxCacheMem:      50,
        MaxCacheSize:     100,
    })
}

defer vips.Shutdown()
```

## Loading Images

```go
// From file
img, err := vips.NewImageFromFile("image.jpg")

// From buffer
img, err := vips.NewImageFromBuffer(data)

// From reader
img, err := vips.NewImageFromReader(reader)
```

## Image Operations

### Resizing & Scaling

```go
// Resize with uniform scale
err := img.Resize(0.5, vips.KernelAuto)

// Thumbnail (faster, with cropping)
err := img.Thumbnail(200, 200, vips.InterestingCentre)
```

### Rotation & Flipping

```go
// Auto-rotate based on EXIF
err := img.AutoRotate()

// Rotate by specific angle
err := img.Rotate(vips.Angle90)

// Flip
err := img.Flip(vips.DirectionHorizontal)
```

### Color Operations

```go
// Convert color space
err := img.ToColorSpace(vips.InterpretationSRGB)

// Modulate (brightness, saturation, hue)
err := img.Modulate(1.2, 1.0, 0)

// Add alpha channel
err := img.AddAlpha()
```

### Filters & Effects

```go
// Gaussian blur
err := img.GaussianBlur(2.0, 2.0)

// Sharpen
err := img.Sharpen(sigma, x1, m2)
```

## Export Operations

```go
// JPEG
params := vips.NewJpegExportParams()
params.Quality = 85
data, metadata, err := img.ExportJpeg(params)

// WebP
params := vips.NewWebpExportParams()
params.Quality = 80
data, metadata, err := img.ExportWebp(params)

// AVIF
params := vips.NewAvifExportParams()
params.Quality = 85
data, metadata, err := img.ExportAvif(params)
```

## Supported Image Types

- `ImageTypeGIF`, `ImageTypeJPEG`, `ImageTypePNG`, `ImageTypeWEBP`
- `ImageTypeTIFF`, `ImageTypeHEIF`, `ImageTypeAVIF`
- `ImageTypeJP2K`, `ImageTypeJXL`, `ImageTypeSVG`, `ImageTypePDF`

## Kernels (for resizing)

- `KernelAuto`, `KernelCubic`, `KernelLinear`
- `KernelNearest`, `KernelLanczos2`, `KernelLanczos3`

## Best Practices

1. Always call `Startup()` before using vips functions
2. Call `Shutdown()` when done (use defer)
3. Close images manually in high-volume applications
4. Use format-specific export params for better control
5. Configure concurrency based on system resources
6. Use `ShutdownThread()` when using vips in goroutines

## Resources

- **GitHub:** https://github.com/davidbyttow/govips
- **Go Reference:** https://pkg.go.dev/github.com/davidbyttow/govips/v2/vips
- **libvips:** https://www.libvips.org/
