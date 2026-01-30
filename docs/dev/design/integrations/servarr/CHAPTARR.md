# Chaptarr Integration

> Book & audiobook management automation (uses Readarr API)

**Status**: 🟡 PLANNED
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

- [ ] **API Client** (`internal/service/metadata/provider_readarr.go`)
  - [ ] Author listing & detail fetching
  - [ ] Book listing & detail fetching
  - [ ] Book file listing & detail fetching
  - [ ] Quality profile mapping
  - [ ] Metadata profile mapping
  - [ ] Health check integration

- [ ] **Webhook Handler** (`internal/api/handlers/webhook_readarr.go`)
  - [ ] Parse webhook payload (On Download event)
  - [ ] Extract author + book + file metadata
  - [ ] Detect audiobook vs ebook (based on file format OR `isAudiobook` flag)
  - [ ] Route to appropriate module (audiobook OR book)
  - [ ] Trigger metadata enrichment (Goodreads, Audible)
  - [ ] Store in PostgreSQL (`book_authors`, `books` OR `audiobooks`)
  - [ ] Update Typesense search index

- [ ] **Metadata Sync**
  - [ ] Map Readarr authors → Revenge `book_authors` table
  - [ ] Map Readarr books → Revenge `books` table (ebooks)
  - [ ] Map Readarr audiobooks → Revenge `audiobooks` table
  - [ ] Map Readarr quality profiles → Revenge quality tiers
  - [ ] Handle book series (series title + position)

- [ ] **Audiobook vs Ebook Detection**
  - [ ] Check file extension: `.mp3`, `.m4b`, `.m4a` → audiobook module
  - [ ] Check file extension: `.epub`, `.pdf`, `.mobi`, `.azw3` → book module
  - [ ] Fallback: Use Readarr `isAudiobook` flag (if available)

- [ ] **Quality Profile Mapping**
  - [ ] **Audiobooks**:
    - [ ] High (320kbps MP3) → `quality='high'`, `bitrate=320`
    - [ ] Medium (192kbps MP3) → `quality='medium'`, `bitrate=192`
    - [ ] Low (128kbps MP3) → `quality='low'`, `bitrate=128`
  - [ ] **Ebooks**:
    - [ ] EPUB → `format='epub'`
    - [ ] PDF → `format='pdf'`
    - [ ] MOBI → `format='mobi'`
    - [ ] AZW3 → `format='azw3'`

- [ ] **Error Handling**
  - [ ] Retry failed API calls (circuit breaker)
  - [ ] Log webhook failures
  - [ ] Handle missing books (not yet released)

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
