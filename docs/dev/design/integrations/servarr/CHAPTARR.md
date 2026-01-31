# Chaptarr Integration

> Book & audiobook management automation (uses Readarr API)

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | |
| Sources | ✅ | |
| Instructions | 🟡 | |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Priority**: 🟡 MEDIUM (Phase 6 - Book Module)
**Type**: Webhook listener + API client for metadata sync

---

## Overview

Chaptarr is a book and audiobook management tool that uses the Readarr API.
It replaces Readarr, which is currently inactive. Revenge integrates with Chaptarr to:
- Receive webhook notifications when books are imported
- Sync author and book metadata
- Monitor Readarr download/import status
- Separate audiobooks from ebooks (distinct modules)

**Integration Points**:
- **Webhook listener**: Process Readarr events (On Import, On Book Added, etc.)
- **API client**: Query books, authors, editions
- **Metadata sync**: Enrich Revenge metadata with Readarr data
- **Module routing**: Audiobooks → Audiobook module, ebooks → Book module

---

## Developer Resources

- 📚 **API Docs**: https://readarr.com/docs/api/
- 🔗 **OpenAPI Spec**: https://github.com/Readarr/Readarr/blob/develop/src/Readarr.Api.V1/openapi.json
- 🔗 **GitHub**: https://github.com/Readarr/Readarr
- 🔗 **Wiki**: https://wiki.servarr.com/readarr

---

## API Details

**Base Path**: `/api/v1/`
**Authentication**: `X-Api-Key` header (API key from Readarr settings)
**Rate Limits**: None (self-hosted)

### Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/book` | GET | List all books |
| `/book/{id}` | GET | Get specific book details |
| `/author` | GET | List all authors |
| `/author/{id}` | GET | Get specific author details |
| `/bookfile` | GET | List book files |
| `/bookfile/{id}` | GET | Get specific book file details |
| `/importlist` | GET | List configured import lists |
| `/metadata` | GET | Get metadata settings |
| `/qualityprofile` | GET | List quality profiles |
| `/metadataprofile` | GET | List metadata profiles |
| `/system/status` | GET | Get Readarr version & status |
| `/health` | GET | Check Readarr health |

---

## Webhook Events

Readarr can send webhooks for the following events:

### On Import (Book Downloaded & Imported)
```json
{
  "eventType": "Download",
  "author": {
    "id": 1,
    "name": "Brandon Sanderson",
    "foreignAuthorId": "goodreads-123456",
    "path": "/media/Books/Brandon Sanderson"
  },
  "book": {
    "id": 123,
    "title": "The Way of Kings",
    "seriesTitle": "The Stormlight Archive",
    "position": "1",
    "releaseDate": "2010-08-31",
    "foreignBookId": "goodreads-7235533",
    "overview": "Epic fantasy novel...",
    "images": [
      {
        "coverType": "cover",
        "url": "https://images.gr-assets.com/books/1388184640l/7235533.jpg"
      }
    ]
  },
  "bookFiles": [
    {
      "id": 456,
      "path": "/media/Books/Brandon Sanderson/The Way of Kings.epub",
      "quality": "EPUB",
      "size": 4194304,
      "mediaInfo": {
        "format": "EPUB",
        "pages": 1007
      }
    }
  ],
  "isAudiobook": false
}
```

### On Book Added (New Book Tracked)
Triggered when Readarr starts monitoring a new book.

### On Book File Delete
Triggered when book file is deleted from Readarr.

### On Author Delete
Triggered when author is removed from Readarr.

### On Rename
Triggered when book files are renamed.

### On Health Issue
Triggered when Readarr detects health issues.

---

## Implementation Checklist

### Phase 1: Client Setup
- [ ] Create client package structure
- [ ] Implement HTTP client with resty
- [ ] Add API key authentication
- [ ] Implement rate limiting

### Phase 2: API Implementation
- [ ] Implement core API methods
- [ ] Add response type definitions
- [ ] Implement error handling

### Phase 3: Service Integration
- [ ] Create service wrapper
- [ ] Add caching layer
- [ ] Implement fx module wiring

### Phase 4: Testing
- [ ] Add unit tests with mocks
- [ ] Add integration tests

---

## Revenge Integration Pattern

```
Readarr imports book (The Way of Kings)
           ↓
Sends webhook to Revenge
           ↓
Revenge processes webhook
           ↓
Detects file type (EPUB = ebook)
           ↓
Stores author/book in PostgreSQL (book_authors, books)
           ↓
Enriches metadata from Goodreads (ratings, description)
           ↓
Updates Typesense search index
           ↓
Book available for reading
```

### Audiobook Routing

```
Readarr imports audiobook (The Way of Kings.m4b)
           ↓
Sends webhook to Revenge
           ↓
Revenge processes webhook
           ↓
Detects file type (.m4b = audiobook)
           ↓
Routes to Audiobook module
           ↓
Stores in PostgreSQL (audiobook_authors, audiobooks)
           ↓
Enriches metadata from Audible
           ↓
Updates Typesense search index
           ↓
Audiobook available for listening
```

### Go Client Example

```go
type ChaptarrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *ReadarrClient) GetAuthor(ctx context.Context, authorID int) (*Author, error) {
    url := fmt.Sprintf("%s/api/v1/author/%d", c.baseURL, authorID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get author: %w", err)
    }
    defer resp.Body.Close()

    var author Author
    json.NewDecoder(resp.Body).Decode(&author)
    return &author, nil
}

func (c *ReadarrClient) GetBooksByAuthor(ctx context.Context, authorID int) ([]Book, error) {
    url := fmt.Sprintf("%s/api/v1/book?authorId=%d", c.baseURL, authorID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get books: %w", err)
    }
    defer resp.Body.Close()

    var books []Book
    json.NewDecoder(resp.Body).Decode(&books)
    return books, nil
}
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Servarr Wiki](https://wiki.servarr.com/) | [Local](../../../sources/apis/servarr-wiki.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Servarr](INDEX.md)

### In This Section

- [Lidarr Integration](LIDARR.md)
- [Radarr Integration](RADARR.md)
- [Sonarr Integration](SONARR.md)
- [Whisparr v3 Integration](WHISPARR.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Book Module](../../architecture/modules/BOOK.md)
- [Audiobook Module](../../architecture/modules/AUDIOBOOK.md)
- [Goodreads Integration](../metadata/books/GOODREADS.md)
- [Audible Integration](../metadata/books/AUDIBLE.md)
- [Native Audiobook/Podcast](../audiobook/INDEX.md)
- [Arr Integration Pattern](../../patterns/arr_integration.md)
- [Webhook Handling](../../patterns/webhook_patterns.md)

---

## Quality Profile Mapping

### Audiobooks

| Readarr Quality | Revenge Quality | Bitrate | Format |
|----------------|-----------------|---------|--------|
| MP3-320 | `high` | 320 kbps | MP3 (CBR) |
| AAC-256 | `high` | 256 kbps | AAC |
| MP3-192 | `medium` | 192 kbps | MP3 (CBR) |
| AAC-192 | `medium` | 192 kbps | AAC |
| MP3-128 | `low` | 128 kbps | MP3 (CBR) |
| AAC-128 | `low` | 128 kbps | AAC |
| Any | `auto` | Varies | Varies |

### Ebooks

| File Extension | Revenge Format |
|----------------|----------------|
| `.epub` | `epub` |
| `.pdf` | `pdf` |
| `.mobi` | `mobi` |
| `.azw3` | `azw3` |
| `.txt` | `txt` |
| `.cbz` | `cbz` (comic book archive) |
| `.cbr` | `cbr` (comic book archive) |

---

## Notes

- **Goodreads is primary metadata source** for ebooks (consistency with Readarr)
- **Audible is primary metadata source** for audiobooks
- Readarr API v1 is stable (widely adopted)
- Self-hosted = no rate limits (unlike cloud APIs)
- Quality profiles are customizable in Readarr (respect user settings)
- Readarr uses Goodreads IDs (`foreignAuthorId`, `foreignBookId`)
- Book series: Readarr tracks series title + position (e.g., "The Stormlight Archive #1")
- Editions: Readarr can track multiple editions of the same book (different covers, ISBNs)
- Metadata profiles: Control which book types to monitor (novels, non-fiction, graphic novels, etc.)
- Release date: Readarr uses earliest release date from Goodreads
- Wanted missing: Readarr tracks "monitored" status (not yet released = no file)
- **Audiobook vs Ebook detection**: Use file extension + Readarr `isAudiobook` flag
- **Module routing**: Audiobooks → Audiobook module, ebooks → Book module (distinct databases)
