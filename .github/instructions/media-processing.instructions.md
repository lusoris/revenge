---
applyTo: "**/internal/service/media/**/*.go,**/internal/content/**/scanner*.go"
---

# Media Processing Instructions

> FFmpeg bindings, image processing, audio metadata extraction

## Overview

Media processing is delegated to external service (Blackbeard) for transcoding, but Revenge handles:

- **Media probing** - Extract metadata from files
- **Image processing** - Thumbnails, blurhash generation
- **Audio metadata** - ID3 tags, embedded artwork
- **Subtitle extraction** - Container subtitles to external files

## Packages

| Package       | Purpose                               | Live Docs                                                     |
| ------------- | ------------------------------------- | ------------------------------------------------------------- |
| `go-astiav`   | FFmpeg bindings (probing, extraction) | [go-astiav.md](../../docs/dev/sources/media/go-astiav.md)     |
| `bimg`        | Image processing (libvips)            | [bimg.md](../../docs/dev/sources/media/bimg.md)               |
| `dhowden/tag` | Audio metadata (ID3, Vorbis)          | [dhowden-tag.md](../../docs/dev/sources/media/dhowden-tag.md) |
| `bogem/id3v2` | ID3v2 tag writing                     | [bogem-id3v2.md](../../docs/dev/sources/media/bogem-id3v2.md) |
| `go-astisub`  | Subtitle parsing/conversion           | [go-astisub.md](../../docs/dev/sources/media/go-astisub.md)   |
| `go-blurhash` | Blurhash generation                   | [go-blurhash.md](../../docs/dev/sources/media/go-blurhash.md) |

---

## Media Probing (go-astiav)

Extract video/audio stream info without full decode:

```go
import "github.com/asticode/go-astiav"

type MediaInfo struct {
    Duration    time.Duration
    VideoCodec  string
    AudioCodec  string
    Width       int
    Height      int
    Bitrate     int64
    AudioTracks []AudioTrack
    Subtitles   []SubtitleTrack
}

func ProbeMedia(path string) (*MediaInfo, error) {
    // Open input
    inputCtx := astiav.AllocFormatContext()
    defer inputCtx.Free()

    if err := inputCtx.OpenInput(path, nil, nil); err != nil {
        return nil, fmt.Errorf("open input: %w", err)
    }
    defer inputCtx.CloseInput()

    // Find stream info
    if err := inputCtx.FindStreamInfo(nil); err != nil {
        return nil, fmt.Errorf("find stream info: %w", err)
    }

    info := &MediaInfo{
        Duration: time.Duration(inputCtx.Duration()) * time.Microsecond / time.Duration(astiav.TimeBase),
    }

    // Extract stream info
    for _, stream := range inputCtx.Streams() {
        codecParams := stream.CodecParameters()

        switch codecParams.MediaType() {
        case astiav.MediaTypeVideo:
            info.VideoCodec = codecParams.CodecID().String()
            info.Width = codecParams.Width()
            info.Height = codecParams.Height()
            info.Bitrate = codecParams.BitRate()

        case astiav.MediaTypeAudio:
            info.AudioTracks = append(info.AudioTracks, AudioTrack{
                Codec:      codecParams.CodecID().String(),
                Channels:   codecParams.ChannelLayout().NbChannels(),
                SampleRate: codecParams.SampleRate(),
                Language:   stream.Metadata().Get("language", nil, 0).Value(),
            })

        case astiav.MediaTypeSubtitle:
            info.Subtitles = append(info.Subtitles, SubtitleTrack{
                Codec:    codecParams.CodecID().String(),
                Language: stream.Metadata().Get("language", nil, 0).Value(),
                Forced:   stream.Disposition()&astiav.StreamDispositionForced != 0,
            })
        }
    }

    return info, nil
}
```

---

## Image Processing (bimg)

Generate thumbnails and poster images:

```go
import "github.com/h2non/bimg"

// Generate thumbnail
func GenerateThumbnail(input []byte, width, height int) ([]byte, error) {
    return bimg.NewImage(input).Process(bimg.Options{
        Width:   width,
        Height:  height,
        Crop:    true,
        Quality: 85,
        Type:    bimg.WEBP,
    })
}

// Resize maintaining aspect ratio
func ResizeImage(input []byte, maxWidth int) ([]byte, error) {
    img := bimg.NewImage(input)

    size, err := img.Size()
    if err != nil {
        return nil, err
    }

    // Calculate proportional height
    ratio := float64(maxWidth) / float64(size.Width)
    newHeight := int(float64(size.Height) * ratio)

    return img.Process(bimg.Options{
        Width:   maxWidth,
        Height:  newHeight,
        Quality: 90,
        Type:    bimg.WEBP,
    })
}

// Extract frame from video (via FFmpeg, then process with bimg)
func ExtractPosterFrame(videoPath string, timestamp time.Duration) ([]byte, error) {
    // Use go-astiav to extract frame, then process with bimg
    // ... see go-astiav docs for frame extraction
}
```

---

## Blurhash Generation

Generate blurhash for placeholder images:

```go
import (
    "image"
    _ "image/jpeg"
    _ "image/png"

    "github.com/buckket/go-blurhash"
)

func GenerateBlurhash(img image.Image) (string, error) {
    // Components: 4x3 gives good balance of quality/size
    hash, err := blurhash.Encode(4, 3, img)
    if err != nil {
        return "", err
    }
    return hash, nil
}

// From file
func BlurhashFromFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        return "", err
    }

    return GenerateBlurhash(img)
}
```

---

## Audio Metadata (dhowden/tag)

Read audio file metadata:

```go
import "github.com/dhowden/tag"

type AudioMetadata struct {
    Title       string
    Artist      string
    Album       string
    AlbumArtist string
    Year        int
    Track       int
    TotalTracks int
    Genre       string
    Artwork     []byte
}

func ReadAudioMetadata(path string) (*AudioMetadata, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    m, err := tag.ReadFrom(f)
    if err != nil {
        return nil, err
    }

    track, totalTracks := m.Track()

    meta := &AudioMetadata{
        Title:       m.Title(),
        Artist:      m.Artist(),
        Album:       m.Album(),
        AlbumArtist: m.AlbumArtist(),
        Year:        m.Year(),
        Track:       track,
        TotalTracks: totalTracks,
        Genre:       m.Genre(),
    }

    // Extract embedded artwork
    if pic := m.Picture(); pic != nil {
        meta.Artwork = pic.Data
    }

    return meta, nil
}
```

---

## Subtitle Handling (go-astisub)

Parse and convert subtitles:

```go
import "github.com/asticode/go-astisub"

// Read subtitles from file
func ReadSubtitles(path string) (*astisub.Subtitles, error) {
    return astisub.OpenFile(path)
}

// Convert SRT to WebVTT
func ConvertToWebVTT(input, output string) error {
    subs, err := astisub.OpenFile(input)
    if err != nil {
        return err
    }
    return subs.Write(output)
}

// Extract embedded subtitles (via FFmpeg)
func ExtractSubtitles(videoPath string, trackIndex int, outputPath string) error {
    // Use go-astiav to extract subtitle track
    // ... see go-astiav docs
}
```

---

## Integration with Scanner

```go
type MediaScanner struct {
    logger *slog.Logger
}

func (s *MediaScanner) ScanFile(ctx context.Context, path string) (*ScannedMedia, error) {
    ext := strings.ToLower(filepath.Ext(path))

    switch ext {
    case ".mkv", ".mp4", ".avi", ".webm":
        return s.scanVideo(ctx, path)
    case ".mp3", ".flac", ".m4a", ".ogg":
        return s.scanAudio(ctx, path)
    case ".jpg", ".jpeg", ".png", ".webp":
        return s.scanImage(ctx, path)
    default:
        return nil, fmt.Errorf("unsupported format: %s", ext)
    }
}

func (s *MediaScanner) scanVideo(ctx context.Context, path string) (*ScannedMedia, error) {
    info, err := ProbeMedia(path)
    if err != nil {
        return nil, err
    }

    return &ScannedMedia{
        Path:        path,
        Type:        MediaTypeVideo,
        Duration:    info.Duration,
        VideoCodec:  info.VideoCodec,
        AudioCodec:  info.AudioCodec,
        Resolution:  fmt.Sprintf("%dx%d", info.Width, info.Height),
        AudioTracks: info.AudioTracks,
        Subtitles:   info.Subtitles,
    }, nil
}
```

---

## fx Module

```go
var MediaProcessingModule = fx.Options(
    fx.Provide(NewMediaProber),
    fx.Provide(NewImageProcessor),
    fx.Provide(NewAudioMetadataReader),
    fx.Provide(NewSubtitleExtractor),
)
```

---

## DO's and DON'Ts

### DO ✅

- Use go-astiav for probing (no FFmpeg shell calls)
- Use bimg for image processing (faster than ImageMagick)
- Generate blurhash for all poster images
- Cache extracted metadata in database
- Extract subtitles during library scan

### DON'T ❌

- Shell out to FFmpeg for probing (use go-astiav)
- Transcode in Revenge (use Blackbeard)
- Store extracted images in database (use filesystem)
- Block on image processing (use River jobs)

---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index
- [OFFLOADING.md](../../docs/dev/design/technical/OFFLOADING.md) - Blackbeard transcoding
- [river-job-queue.instructions.md](river-job-queue.instructions.md) - Background processing
- [fsnotify-file-watching.instructions.md](fsnotify-file-watching.instructions.md) - File scanning
