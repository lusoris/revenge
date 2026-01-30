# Audible Integration

> Audiobook metadata provider - primary source for audiobooks

**Service**: Audible (Amazon)
**Type**: Metadata Provider (Audiobooks)
**API Version**: NO official public API
**Website**: https://www.audible.com
**Alternative**: Web scraping OR Chaptarr metadata

---

## Overview

**Audible** is the primary audiobook platform, but **NO official public API** exists. Alternative approaches required.

**Why Audible (if API existed)**:
- Comprehensive audiobook catalog
- Narrator information
- Audio sample previews
- User ratings and reviews
- Series information
- Publisher information

**Current Status**:
- ❌ **NO official public API**
- ✅ **Web scraping** (fragile, not recommended)
- ✅ **Chaptarr metadata** (uses Goodreads/custom sources)
- ✅ **Unofficial API libraries** (community projects, no guarantees)

**Use Cases**:
- Audiobook metadata (title, author, narrator, duration, publication date)
- Cover art
- Series information
- Sample audio clips
- User ratings

---

## Developer Resources

**API Status**: NO official public API
**Unofficial Libraries**:
- **audible-cli** (Python): https://github.com/mkb79/audible-cli
- **Audnexus** (community API): https://github.com/laxamentumtech/audnexus

**Alternatives**:
1. **Chaptarr** (recommended): Use Chaptarr's metadata sources
2. **Web scraping** (fragile, not recommended)
3. **Audnexus API** (community-driven, unofficial)

---

## Migration Strategy

### Option 1: Chaptarr Metadata (Recommended)
Use **Chaptarr** as audiobook manager, fetch metadata from Chaptarr.

**Advantages**:
- Chaptarr has audiobook support (uses Goodreads, custom sources)
- Webhook integration
- Automated downloads
- Quality management

**See**: [../../servarr/CHAPTARR.md](../../servarr/CHAPTARR.md)

### Option 2: Audnexus API (Community Project)
Use **Audnexus** API (unofficial community project).

**Audnexus**: https://github.com/laxamentumtech/audnexus
**API Docs**: https://api.audnex.us/

**Advantages**:
- RESTful API
- Audiobook metadata (Audible + Goodreads + ASIN lookup)
- Narrator information
- Series information
- Free (self-hosted OR public instance)

**Disadvantages**:
- Unofficial (may break if Audible changes)
- No guarantees
- Community-maintained

### Option 3: Web Scraping (NOT Recommended)
Scrape Audible website for metadata.

**Disadvantages**:
- Legal gray area (ToS violation)
- Fragile (breaks when HTML changes)
- No official support
- CAPTCHA challenges
- Rate limiting required

---

## Audnexus API Details (Unofficial)

### Base URL
```
https://api.audnex.us/
```

### Authentication
None required (public instance)

### Key Endpoints

#### Search Audiobooks
```bash
GET /books?title={TITLE}&author={AUTHOR}
```

#### Get Audiobook by ASIN
```bash
GET /books/{ASIN}
```

**ASIN**: Audible Standard Identification Number (e.g., `B002V1A0WE`)

**Response** (Example):
```json
{
  "asin": "B002V1A0WE",
  "title": "Harry Potter and the Philosopher's Stone",
  "subtitle": "Book 1",
  "authors": [{"asin": "B000APZOQA", "name": "J.K. Rowling"}],
  "narrators": [{"name": "Stephen Fry"}],
  "publisher": "Pottermore Publishing",
  "publisherSummary": "Turning the envelope over...",
  "releaseDate": "2015-11-20",
  "language": "English",
  "runtimeLengthMin": 477,
  "image": "https://m.media-amazon.com/images/I/...",
  "rating": 4.9,
  "ratings_count": 12345,
  "series": [
    {
      "asin": "B017V4IM8M",
      "title": "Harry Potter",
      "position": "1"
    }
  ]
}
```

---

## Implementation Checklist

### Option 1: Chaptarr Integration (Recommended)
- [ ] **Use Chaptarr for audiobook management** (see [../../servarr/CHAPTARR.md](../../servarr/CHAPTARR.md))
- [ ] Fetch audiobook metadata from Chaptarr API
- [ ] Webhook integration (Chaptarr → Revenge)
- [ ] Store in `audiobooks` table

### Option 2: Audnexus API Client (Unofficial)
- [ ] Base URL configuration
- [ ] HTTP client with User-Agent
- [ ] Error handling (404: Audiobook not found, 500: Server error)
- [ ] Response parsing (JSON unmarshalling)
- [ ] Rate limiting (no strict limit, respect fair use ~1 req/s)

### Audiobook Metadata
- [ ] Search audiobooks by title, author
- [ ] Fetch audiobook by ASIN
- [ ] Extract: title, author, narrator, duration, publication date, series, publisher
- [ ] Store in `audiobooks` table

### Narrator Information
- [ ] Fetch narrator details
- [ ] Store in `audiobook_narrators` table
- [ ] Link audiobooks to narrators

### Cover Art Handling
- [ ] Download cover from Audnexus/Audible CDN
- [ ] Generate Blurhash
- [ ] Convert to WebP
- [ ] Store locally (`data/audiobooks/covers/`)

### Series Handling
- [ ] Fetch series information
- [ ] Store in `audiobook_series` table
- [ ] Link audiobooks to series (position in series)

### Error Handling
- [ ] Handle 404 (Audiobook not found)
- [ ] Handle 500 (Server error - retry)
- [ ] Log errors (no sensitive data)

---

## Integration Pattern

### Chaptarr Webhook → Audiobook Metadata Sync
```go
// Webhook: Chaptarr added new audiobook
func (s *AudiobookService) HandleChaptarrAudiobookAdded(audiobookID string) error {
    // 1. Get audiobook from Chaptarr
    chaptarrAudiobook := s.chaptarrClient.GetBook(audiobookID)

    // 2. Check if audiobook (not ebook)
    if !chaptarrAudiobook.IsAudiobook {
        return nil // Skip ebooks
    }

    // 3. Fetch additional metadata from Audnexus (optional)
    var narrator string
    if asin := chaptarrAudiobook.ForeignBookId; asin != "" {
        audnexusBook := s.audnexusClient.GetBookByASIN(asin)
        if audnexusBook != nil {
            narrator = audnexusBook.Narrators[0].Name
        }
    }

    // 4. Store in Revenge database
    s.db.InsertAudiobook(map[string]interface{}{
        "title":            chaptarrAudiobook.Title,
        "author":           chaptarrAudiobook.Author.Name,
        "narrator":         narrator,
        "duration_minutes": chaptarrAudiobook.DurationMinutes,
        "release_date":     chaptarrAudiobook.ReleaseDate,
        "asin":             chaptarrAudiobook.ForeignBookId,
    })

    return nil
}
```

---

## Related Documentation

- **Audiobook Module**: [MODULE_IMPLEMENTATION_TODO.md](../../../planning/MODULE_IMPLEMENTATION_TODO.md) (Audiobook section)
- **Chaptarr Integration**: [../../servarr/CHAPTARR.md](../../servarr/CHAPTARR.md) (recommended approach)
- **Native Audiobook/Podcast**: [../audiobook/INDEX.md](../audiobook/INDEX.md) (native implementation)

---

## Notes

- **NO official API**: Audible does not provide public API
- **Unofficial libraries**: audible-cli, Audnexus (community projects, no guarantees)
- **Audnexus recommended**: Community-driven API aggregator (Audible + Goodreads + ASIN lookup)
- **Chaptarr recommended**: Use Chaptarr for audiobook management (has built-in metadata sources)
- **ASIN**: Audible Standard Identification Number (unique identifier, e.g., B002V1A0WE)
- **Narrator information**: Critical for audiobooks (Audnexus provides, Readarr may not)
- **Series information**: Harry Potter, Lord of the Rings, etc. (position in series)
- **Duration**: Runtime in minutes (important for audiobooks)
- **Cover art**: High-quality covers available (download from Audible CDN via Audnexus)
- **Web scraping**: NOT recommended (fragile, legal gray area, CAPTCHA challenges)
- **Fallback strategy**: Chaptarr primary, Audnexus fallback (narrator info)
- **User privacy**: No user data collected (metadata only)
- **Self-hosted option**: Audnexus can be self-hosted (Docker available)
- **Public instance**: https://api.audnex.us/ (community-maintained, may have downtime)
- **Rate limiting**: No strict limits (Audnexus), respect fair use (~1 req/s)
- **Quality**: Audnexus data quality depends on Audible scraping (may have incomplete data)
- **Alternative**: Use OpenLibrary for book metadata, add narrator manually (not ideal)
