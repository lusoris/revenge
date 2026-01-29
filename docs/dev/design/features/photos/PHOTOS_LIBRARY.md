# Photos Library

> Photo organization, viewing, and management

**Status**: 🔴 PLANNING
**Priority**: 🟢 HIGH (Critical Gap - Jellyfin/Plex/Emby have this)
**Inspired By**: Plex Photos, Google Photos, Apple Photos

---

## Overview

Photos Library provides organization, viewing, and management of photo collections with features like albums, face recognition (optional), and slideshows.

---

## Features

### Core

| Feature | Description |
|---------|-------------|
| Photo Import | Scan and import photos from directories |
| Album Management | Create, edit, delete albums |
| Timeline View | Photos organized by date |
| Slideshow | Automatic photo display with transitions |
| Sharing | Share albums via links |
| Favorites | Mark photos as favorites |

### Advanced

| Feature | Description |
|---------|-------------|
| Face Recognition | Optional face detection and grouping |
| Location View | Map-based photo browsing (GPS EXIF) |
| Auto Albums | Smart albums (e.g., "Best of 2024") |
| Duplicate Detection | Find and manage duplicate photos |
| RAW Support | View RAW camera formats |
| Video Clips | Short video clips alongside photos |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      Photos Pipeline                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Photo File ──► EXIF Extract ──► Thumbnail ──► Face Detect     │
│                                                                 │
│  ┌───────────┐   ┌───────────┐   ┌───────────┐   ┌──────────┐ │
│  │  Scanner  │──►│   EXIF    │──►│   bimg    │──►│  GoFace  │ │
│  │           │   │  Parser   │   │  Resize   │   │ (optional)│ │
│  └───────────┘   └───────────┘   └───────────┘   └──────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │                     PostgreSQL                            │ │
│  │  photos │ albums │ faces │ photo_albums │ photo_faces    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Go Packages

| Package | Purpose | URL |
|---------|---------|-----|
| **h2non/bimg** | Thumbnail generation | github.com/h2non/bimg |
| **rwcarlsen/goexif** | EXIF metadata parsing | github.com/rwcarlsen/goexif |
| **bbrks/go-blurhash** | Placeholder images | github.com/bbrks/go-blurhash |
| **Kagami/go-face** | Face recognition (optional) | github.com/Kagami/go-face |
| **golang/image** | Image processing | golang.org/x/image |

---

## Database Schema

```sql
-- Photos
CREATE TABLE photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID REFERENCES libraries(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    file_size_bytes BIGINT,
    file_hash VARCHAR(64), -- SHA256 for duplicate detection

    -- Image properties
    width INT,
    height INT,
    format VARCHAR(20), -- jpeg, png, heic, raw, etc.
    orientation INT,

    -- Thumbnails
    thumbnail_path TEXT,
    blurhash VARCHAR(50),

    -- EXIF metadata
    taken_at TIMESTAMPTZ,
    camera_make VARCHAR(100),
    camera_model VARCHAR(100),
    lens VARCHAR(200),
    focal_length_mm DECIMAL(6,2),
    aperture DECIMAL(4,2),
    shutter_speed VARCHAR(20),
    iso INT,

    -- Location (GPS)
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    altitude_m DECIMAL(8,2),
    location_name VARCHAR(500), -- Reverse geocoded

    -- User data
    is_favorite BOOLEAN DEFAULT false,
    is_hidden BOOLEAN DEFAULT false,
    rating INT CHECK (rating BETWEEN 0 AND 5),
    description TEXT,

    -- Processing status
    processed BOOLEAN DEFAULT false,
    faces_detected BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(library_id, file_path)
);

-- Albums
CREATE TABLE photo_albums (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    cover_photo_id UUID REFERENCES photos(id) ON DELETE SET NULL,
    is_smart BOOLEAN DEFAULT false,
    smart_rules JSONB, -- For smart albums
    is_public BOOLEAN DEFAULT false,
    share_code VARCHAR(20) UNIQUE,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Photo-Album relationship
CREATE TABLE photo_album_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    album_id UUID REFERENCES photo_albums(id) ON DELETE CASCADE,
    photo_id UUID REFERENCES photos(id) ON DELETE CASCADE,
    sort_order INT DEFAULT 0,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(album_id, photo_id)
);

-- Faces (optional)
CREATE TABLE photo_faces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    photo_id UUID REFERENCES photos(id) ON DELETE CASCADE,
    person_id UUID REFERENCES photo_people(id) ON DELETE SET NULL,

    -- Face location in image
    x INT NOT NULL,
    y INT NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,

    -- Face embedding for recognition
    embedding VECTOR(128), -- Face descriptor

    -- Confidence
    detection_confidence DECIMAL(5,4),
    recognition_confidence DECIMAL(5,4),

    is_verified BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- People (for face recognition)
CREATE TABLE photo_people (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200),
    thumbnail_face_id UUID REFERENCES photo_faces(id),
    photo_count INT DEFAULT 0,
    is_favorite BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_photos_library ON photos(library_id);
CREATE INDEX idx_photos_taken_at ON photos(taken_at DESC);
CREATE INDEX idx_photos_location ON photos USING GIST (
    ll_to_earth(latitude, longitude)
) WHERE latitude IS NOT NULL;
CREATE INDEX idx_photos_hash ON photos(file_hash);
CREATE INDEX idx_photo_album_items_album ON photo_album_items(album_id);
CREATE INDEX idx_photo_faces_photo ON photo_faces(photo_id);
CREATE INDEX idx_photo_faces_person ON photo_faces(person_id);
```

---

## River Jobs

```go
const (
    JobKindScanPhotos      = "photos.scan"
    JobKindProcessPhoto    = "photos.process"
    JobKindDetectFaces     = "photos.detect_faces"
    JobKindClusterFaces    = "photos.cluster_faces"
    JobKindGenerateThumbs  = "photos.generate_thumbnails"
    JobKindReverseGeocode  = "photos.reverse_geocode"
)

type ProcessPhotoArgs struct {
    PhotoID   uuid.UUID `json:"photo_id"`
    PhotoPath string    `json:"photo_path"`
    Tasks     []string  `json:"tasks"` // exif, thumbnail, blurhash, faces
}
```

---

## Go Implementation

```go
// internal/content/photos/

type Service struct {
    repo     PhotoRepository
    albumRepo AlbumRepository
    river    *river.Client[pgx.Tx]
    face     *FaceDetector // Optional
}

type PhotoProcessor struct {
    thumbWidth  int
    thumbHeight int
}

func (p *PhotoProcessor) Process(ctx context.Context, photoPath string) (*ProcessResult, error) {
    // Read image
    data, err := os.ReadFile(photoPath)
    if err != nil {
        return nil, err
    }

    img := bimg.NewImage(data)

    // Extract EXIF
    exif, _ := p.extractEXIF(data)

    // Generate thumbnail
    thumbnail, _ := img.Thumbnail(p.thumbWidth)

    // Generate blurhash
    blurhash, _ := p.generateBlurhash(thumbnail)

    return &ProcessResult{
        EXIF:      exif,
        Thumbnail: thumbnail,
        Blurhash:  blurhash,
    }, nil
}

func (p *PhotoProcessor) extractEXIF(data []byte) (*EXIFData, error) {
    reader := bytes.NewReader(data)
    x, err := exif.Decode(reader)
    if err != nil {
        return nil, err
    }

    // Extract common fields
    var result EXIFData

    if dt, err := x.DateTime(); err == nil {
        result.TakenAt = dt
    }
    if make, err := x.Get(exif.Make); err == nil {
        result.CameraMake = make.StringVal()
    }
    if lat, lon, err := x.LatLong(); err == nil {
        result.Latitude = lat
        result.Longitude = lon
    }

    return &result, nil
}

func (p *PhotoProcessor) generateBlurhash(thumbnail []byte) (string, error) {
    img, _, err := image.Decode(bytes.NewReader(thumbnail))
    if err != nil {
        return "", err
    }
    return blurhash.Encode(4, 3, img)
}
```

---

## API Endpoints

```
# Photos
GET  /api/v1/photos                  # List photos (paginated, filtered)
GET  /api/v1/photos/:id              # Get photo details
GET  /api/v1/photos/:id/image        # Get full image
GET  /api/v1/photos/:id/thumbnail    # Get thumbnail
PUT  /api/v1/photos/:id              # Update metadata
DELETE /api/v1/photos/:id            # Delete photo

# Timeline
GET  /api/v1/photos/timeline         # Photos grouped by date

# Albums
GET  /api/v1/photos/albums           # List albums
POST /api/v1/photos/albums           # Create album
GET  /api/v1/photos/albums/:id       # Get album
PUT  /api/v1/photos/albums/:id       # Update album
DELETE /api/v1/photos/albums/:id     # Delete album
POST /api/v1/photos/albums/:id/photos # Add photos to album

# Sharing
GET  /api/v1/photos/shared/:code     # View shared album (public)
POST /api/v1/photos/albums/:id/share # Generate share link

# Faces (optional)
GET  /api/v1/photos/people           # List recognized people
GET  /api/v1/photos/people/:id       # Get person's photos
PUT  /api/v1/photos/people/:id       # Name/merge person
POST /api/v1/photos/faces/:id/identify # Manually identify face

# Map
GET  /api/v1/photos/map              # Photos with location for map view

# Slideshow
GET  /api/v1/photos/slideshow/:album_id # Get slideshow config
```

---

## Configuration

```yaml
photos:
  enabled: true

  processing:
    thumbnail:
      width: 300
      height: 300
      quality: 80
    blurhash:
      enabled: true
      x_components: 4
      y_components: 3

  faces:
    enabled: false  # Requires go-face + dlib
    min_confidence: 0.6
    clustering:
      min_samples: 3
      distance_threshold: 0.4

  geocoding:
    enabled: true
    provider: nominatim  # or google
    rate_limit: 1/second

  formats:
    - jpeg
    - jpg
    - png
    - heic
    - webp
    - gif
    - raw
    - cr2
    - nef
    - arw

  video_clips:
    enabled: true
    max_duration_seconds: 30
    formats: [mp4, mov]
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `photos.view` | View photos |
| `photos.upload` | Upload photos |
| `photos.edit` | Edit photo metadata |
| `photos.delete` | Delete photos |
| `photos.albums.create` | Create albums |
| `photos.albums.share` | Share albums |
| `photos.faces.manage` | Manage face recognition |

---

## Related Documentation

- [Library Types](LIBRARY_TYPES.md)
- [Go Packages](../architecture/GO_PACKAGES.md)
- [Client Support](CLIENT_SUPPORT.md)
