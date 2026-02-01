---
sources:
  - name: go-blurhash
    url: ../../../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: google/uuid
    url: ../../../sources/tooling/uuid.md
    note: Auto-resolved from google-uuid
  - name: pgx PostgreSQL Driver
    url: ../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Servarr Wiki
    url: ../../../sources/apis/servarr-wiki.md
    note: Auto-resolved from servarr-wiki
  - name: sqlc
    url: ../../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
  - name: Typesense API
    url: ../../../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../../../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Adult Gallery Module (QAR: Treasures)](#adult-gallery-module-qar-treasures)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Gallery settings](#gallery-settings)
- [Prowlarr integration](#prowlarr-integration)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Adult Gallery Module (QAR: Treasures)


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Scenes, Performers, Studios

> Image gallery management for adult content with performer links and Prowlarr integration

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Database Schema

**Schema**: `qar`

<!-- Schema diagram -->

### Module Structure

```
internal/content/adult_gallery_(qar:_treasures)/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── adult_gallery_(qar:_treasures)_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### File Structure

**Key Files**:
- `internal/content/adult/gallery/service.go` - Gallery management logic
- `internal/content/adult/gallery/scanner/*.go` - Folder scanning and EXIF extraction
- `internal/content/adult/gallery/prowlarr/*.go` - Prowlarr integration
- `migrations/qar/015_galleries.sql` - Database schema


### Key Interfaces

```go
// GalleryService manages adult image galleries
type GalleryService interface {
  // Galleries
  CreateGallery(ctx context.Context, gallery Gallery) (*Gallery, error)
  GetGallery(ctx context.Context, id uuid.UUID) (*Gallery, error)
  ListGalleries(ctx context.Context, filters GalleryFilters) ([]Gallery, int, error)
  UpdateGallery(ctx context.Context, id uuid.UUID, updates GalleryUpdate) (*Gallery, error)
  DeleteGallery(ctx context.Context, id uuid.UUID) error

  // Images
  GetGalleryImages(ctx context.Context, galleryID uuid.UUID) ([]GalleryImage, error)
  GetImage(ctx context.Context, imageID uuid.UUID) (*GalleryImage, error)

  // Scanning
  ScanFolder(ctx context.Context, folderPath string) (*Gallery, error)
  RescanGallery(ctx context.Context, galleryID uuid.UUID) error
  ScanLibrary(ctx context.Context) (int, error) // Returns count of new galleries

  // Search
  SearchGalleries(ctx context.Context, query string, filters GalleryFilters) ([]Gallery, int, error)

  // Performers/Studios
  LinkPerformer(ctx context.Context, galleryID, performerID uuid.UUID) error
  UnlinkPerformer(ctx context.Context, galleryID, performerID uuid.UUID) error
  GetGalleriesByPerformer(ctx context.Context, performerID uuid.UUID) ([]Gallery, error)

  // Prowlarr
  SyncProwlarrIndexers(ctx context.Context) error
  SearchProwlarr(ctx context.Context, query string) ([]ProwlarrResult, error)
  QueueDownload(ctx context.Context, url string, indexerID uuid.UUID) (*Download, error)
}

// GalleryRepository handles database operations
type GalleryRepository interface {
  // CRUD
  Create(ctx context.Context, gallery Gallery) error
  GetByID(ctx context.Context, id uuid.UUID) (*Gallery, error)
  List(ctx context.Context, filters GalleryFilters, limit, offset int) ([]Gallery, int, error)
  Update(ctx context.Context, gallery Gallery) error
  Delete(ctx context.Context, id uuid.UUID) error

  // Images
  AddImage(ctx context.Context, image GalleryImage) error
  GetImages(ctx context.Context, galleryID uuid.UUID) ([]GalleryImage, error)
  UpdateImagePosition(ctx context.Context, imageID uuid.UUID, position int) error

  // Relationships
  LinkPerformer(ctx context.Context, galleryID, performerID uuid.UUID) error
  UnlinkPerformer(ctx context.Context, galleryID, performerID uuid.UUID) error
  GetPerformers(ctx context.Context, galleryID uuid.UUID) ([]Performer, error)
  GetByPerformer(ctx context.Context, performerID uuid.UUID) ([]Gallery, error)

  // Search
  Search(ctx context.Context, query string, filters GalleryFilters, limit, offset int) ([]Gallery, int, error)

  // Downloads
  QueueDownload(ctx context.Context, download Download) error
  GetQueuedDownloads(ctx context.Context, limit int) ([]Download, error)
  UpdateDownloadStatus(ctx context.Context, downloadID uuid.UUID, status DownloadStatus) error
}

// Scanner extracts gallery metadata from filesystem
type Scanner interface {
  // ScanFolder creates gallery from folder of images
  ScanFolder(ctx context.Context, folderPath string) (*Gallery, []GalleryImage, error)

  // ExtractEXIF extracts EXIF metadata from image
  ExtractEXIF(ctx context.Context, imagePath string) (*EXIFData, error)

  // GenerateThumbnail creates thumbnail for image
  GenerateThumbnail(ctx context.Context, imagePath string, maxWidth, maxHeight int) ([]byte, error)

  // GenerateBlurhash creates blurhash for image
  GenerateBlurhash(ctx context.Context, imagePath string) (string, error)
}

// ProwlarrClient interacts with Prowlarr API
type ProwlarrClient interface {
  // GetIndexers fetches all configured indexers
  GetIndexers(ctx context.Context) ([]ProwlarrIndexer, error)

  // Search searches for galleries across indexers
  Search(ctx context.Context, query string, categories []int) ([]ProwlarrResult, error)

  // Download queues a download
  Download(ctx context.Context, downloadURL string) error
}

// Types
type Gallery struct {
  ID              uuid.UUID       `db:"id" json:"id"`
  Title           string          `db:"title" json:"title"`
  Description     *string         `db:"description" json:"description,omitempty"`
  FolderPath      string          `db:"folder_path" json:"folder_path"`
  ImageCount      int             `db:"image_count" json:"image_count"`
  TotalSizeBytes  int64           `db:"total_size_bytes" json:"total_size_bytes"`
  CoverImagePath  *string         `db:"cover_image_path" json:"cover_image_path,omitempty"`
  CoverBlurhash   *string         `db:"cover_blurhash" json:"cover_blurhash,omitempty"`
  DateAdded       time.Time       `db:"date_added" json:"date_added"`
  DateCaptured    *time.Time      `db:"date_captured" json:"date_captured,omitempty"`
  Photographer    *string         `db:"photographer" json:"photographer,omitempty"`
  Tags            []string        `db:"tags" json:"tags"`
  Rating          *int            `db:"rating" json:"rating,omitempty"`
  CreatedAt       time.Time       `db:"created_at" json:"created_at"`
  UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`

  // Relationships (not in DB, loaded separately)
  Performers      []Performer     `json:"performers,omitempty"`
  Studios         []Studio        `json:"studios,omitempty"`
}

type GalleryImage struct {
  ID            uuid.UUID       `db:"id" json:"id"`
  TreasureID    uuid.UUID       `db:"treasure_id" json:"treasure_id"`
  FilePath      string          `db:"file_path" json:"file_path"`
  Filename      string          `db:"filename" json:"filename"`
  Position      int             `db:"position" json:"position"`
  Width         *int            `db:"width" json:"width,omitempty"`
  Height        *int            `db:"height" json:"height,omitempty"`
  FileSizeBytes *int64          `db:"file_size_bytes" json:"file_size_bytes,omitempty"`
  MimeType      *string         `db:"mime_type" json:"mime_type,omitempty"`
  Blurhash      *string         `db:"blurhash" json:"blurhash,omitempty"`
  EXIFData      *EXIFData       `db:"exif_data" json:"exif_data,omitempty"`
  DateTaken     *time.Time      `db:"date_taken" json:"date_taken,omitempty"`
  CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}

type EXIFData struct {
  Camera       string  `json:"camera,omitempty"`
  Lens         string  `json:"lens,omitempty"`
  FocalLength  string  `json:"focal_length,omitempty"`
  Aperture     string  `json:"aperture,omitempty"`
  ShutterSpeed string  `json:"shutter_speed,omitempty"`
  ISO          int     `json:"iso,omitempty"`
  GPSLatitude  float64 `json:"gps_latitude,omitempty"`
  GPSLongitude float64 `json:"gps_longitude,omitempty"`
}

type Download struct {
  ID              uuid.UUID       `db:"id" json:"id"`
  Title           string          `db:"title" json:"title"`
  IndexerID       *uuid.UUID      `db:"indexer_id" json:"indexer_id,omitempty"`
  DownloadURL     string          `db:"download_url" json:"download_url"`
  Status          string          `db:"status" json:"status"`
  ProgressPercent int             `db:"progress_percent" json:"progress_percent"`
  DownloadedBytes int64           `db:"downloaded_bytes" json:"downloaded_bytes"`
  TotalBytes      *int64          `db:"total_bytes" json:"total_bytes,omitempty"`
  TreasureID      *uuid.UUID      `db:"treasure_id" json:"treasure_id,omitempty"`
  ErrorMessage    *string         `db:"error_message" json:"error_message,omitempty"`
  QueuedAt        time.Time       `db:"queued_at" json:"queued_at"`
  StartedAt       *time.Time      `db:"started_at" json:"started_at,omitempty"`
  CompletedAt     *time.Time      `db:"completed_at" json:"completed_at,omitempty"`
}
```


### Dependencies

**Go Packages**:
- `github.com/google/uuid` - UUID handling
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/maypok86/otter` - L1 in-memory cache
- `github.com/riverqueue/river` - Background job queue for downloads
- `go.uber.org/fx` - Dependency injection
- `go.uber.org/zap` - Structured logging
- `github.com/bbrks/go-blurhash` - Blurhash generation
- `github.com/disintegration/imaging` - Image processing
- `github.com/rwcarlsen/goexif` - EXIF extraction
- `github.com/typesense/typesense-go` - Search indexing

**External APIs**:
- Prowlarr API - Gallery indexing and downloads






## Configuration
### Environment Variables

```bash
# Gallery settings
GALLERY_LIBRARY_PATH=/path/to/galleries           # Root path for galleries
GALLERY_THUMBNAIL_WIDTH=400                       # Thumbnail max width
GALLERY_THUMBNAIL_HEIGHT=400                      # Thumbnail max height
GALLERY_SCAN_INTERVAL=24h                         # Auto-scan interval
GALLERY_SUPPORTED_FORMATS=jpg,jpeg,png,webp,gif   # Supported image formats

# Prowlarr integration
PROWLARR_URL=http://localhost:9696                # Prowlarr URL
PROWLARR_API_KEY=your-api-key                     # Prowlarr API key
PROWLARR_ENABLED=true                             # Enable Prowlarr integration
PROWLARR_CATEGORIES=2030,2040                     # Adult image categories
```


### Config Keys

```yaml
gallery:
  library_path: /path/to/galleries
  scan_interval: 24h
  supported_formats:
    - jpg
    - jpeg
    - png
    - webp
    - gif

  thumbnails:
    max_width: 400
    max_height: 400
    quality: 85                    # JPEG quality
    format: webp                   # Thumbnail format

  blurhash:
    x_components: 4
    y_components: 3

  prowlarr:
    enabled: true
    url: http://localhost:9696
    api_key: ${PROWLARR_API_KEY}
    categories:
      - 2030                       # Adult - Images
      - 2040                       # Adult - Magazines
    sync_interval: 1h
    max_concurrent_downloads: 3
```



## API Endpoints

### Content Management
**Endpoints**:
```
GET    /api/v1/legacy/treasures                    # List galleries
GET    /api/v1/legacy/treasures/:id                # Get gallery details
GET    /api/v1/legacy/treasures/:id/images         # Get gallery images
PATCH  /api/v1/legacy/treasures/:id                # Update gallery
DELETE /api/v1/legacy/treasures/:id                # Delete gallery
POST   /api/v1/legacy/treasures/:id/scan           # Rescan gallery

GET    /api/v1/legacy/treasures/search             # Search galleries
GET    /api/v1/legacy/treasures/performers/:id     # Galleries by performer

GET    /api/v1/legacy/treasures/prowlarr/search    # Search Prowlarr
POST   /api/v1/legacy/treasures/download           # Queue download
GET    /api/v1/legacy/treasures/downloads          # List downloads
```

**Request/Response Examples**:

**List Galleries**:
```http
GET /api/v1/legacy/treasures?limit=20&offset=0&sort=date_added&order=desc

Response 200 OK:
{
  "treasures": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Beach Photoshoot 2024",
      "description": "Summer beach session",
      "folder_path": "/galleries/beach-2024",
      "image_count": 45,
      "total_size_bytes": 125829120,
      "cover_image_path": "/galleries/beach-2024/IMG_001.jpg",
      "cover_blurhash": "LKO2?U%2Tw=w]~RBVZRi};RPxuwH",
      "date_added": "2024-06-15T10:00:00Z",
      "date_captured": "2024-06-10T14:30:00Z",
      "photographer": "John Doe",
      "tags": ["outdoor", "beach", "summer"],
      "rating": 5,
      "performers": [
        {
          "id": "660e8400-e29b-41d4-a716-446655440001",
          "name": "Jane Smith"
        }
      ]
    }
  ],
  "total": 156,
  "limit": 20,
  "offset": 0
}
```

**Get Gallery Images**:
```http
GET /api/v1/legacy/treasures/550e8400-e29b-41d4-a716-446655440000/images

Response 200 OK:
{
  "images": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "treasure_id": "550e8400-e29b-41d4-a716-446655440000",
      "file_path": "/galleries/beach-2024/IMG_001.jpg",
      "filename": "IMG_001.jpg",
      "position": 0,
      "width": 4000,
      "height": 3000,
      "file_size_bytes": 2854912,
      "mime_type": "image/jpeg",
      "blurhash": "LKO2?U%2Tw=w]~RBVZRi};RPxuwH",
      "exif_data": {
        "camera": "Canon EOS R5",
        "lens": "RF 24-70mm f/2.8",
        "focal_length": "50mm",
        "aperture": "f/2.8",
        "shutter_speed": "1/500",
        "iso": 100,
        "gps_latitude": 34.0522,
        "gps_longitude": -118.2437
      },
      "date_taken": "2024-06-10T14:30:15Z"
    }
  ]
}
```

**Search Prowlarr**:
```http
GET /api/v1/legacy/treasures/prowlarr/search?q=photographer+name

Response 200 OK:
{
  "results": [
    {
      "title": "Photographer Collection 2024",
      "indexer": "ExampleIndexer",
      "size_bytes": 1073741824,
      "download_url": "magnet:?xt=urn:...",
      "published_date": "2024-06-01T00:00:00Z",
      "seeders": 10,
      "leechers": 2
    }
  ]
}
```

**Queue Download**:
```http
POST /api/v1/legacy/treasures/download
{
  "title": "Photographer Collection 2024",
  "download_url": "magnet:?xt=urn:...",
  "indexer_id": "880e8400-e29b-41d4-a716-446655440003"
}

Response 201 Created:
{
  "id": "990e8400-e29b-41d4-a716-446655440004",
  "title": "Photographer Collection 2024",
  "status": "queued",
  "progress_percent": 0,
  "queued_at": "2024-06-15T15:00:00Z"
}
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [go-blurhash](../../../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [google/uuid](../../../sources/tooling/uuid.md) - Auto-resolved from google-uuid
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [Servarr Wiki](../../../sources/apis/servarr-wiki.md) - Auto-resolved from servarr-wiki
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
- [Typesense API](../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

