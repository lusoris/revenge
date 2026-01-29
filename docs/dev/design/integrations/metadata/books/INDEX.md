# Book Metadata Providers

> Books and literature metadata

---

## Overview

Book metadata providers supply information for:
- Titles and authors
- Cover artwork
- ISBN/identifiers
- Publication details
- Descriptions and reviews
- Series information

---

## Providers

| Provider | Type | API | Status |
|----------|------|-----|--------|
| [Open Library](OPENLIBRARY.md) | Books | REST | 🟢 Primary |
| [Hardcover](HARDCOVER.md) | Books | GraphQL | 🟡 Secondary |
| Google Books | Books | REST | 🟡 Planned |
| [Goodreads](GOODREADS.md) | Reviews | Scraping | 🟡 Supplementary |
| [Audible](AUDIBLE.md) | Audiobooks | API | 🟡 Supplementary |
| ISBN DB | ISBN | REST | 🟡 Planned |

---

## Provider Details

### Open Library
**Primary provider - open book database**

- ✅ Comprehensive book data
- ✅ Cover images
- ✅ Author information
- ✅ Edition tracking
- ✅ Free, no API key required
- ✅ Links to Internet Archive

### Google Books
**Secondary for additional metadata**

- ✅ Good search capabilities
- ✅ Preview availability
- ✅ Publisher information
- ✅ High quality covers
- ⚠️ API key required

### Goodreads
**Supplementary for reviews and ratings**

- ✅ User ratings
- ✅ Review excerpts
- ✅ Series information
- ⚠️ No official API (deprecated)
- ⚠️ Requires scraping

### ISBN DB
**Fallback for ISBN lookup**

- ✅ ISBN-10/13 lookup
- ✅ Barcode scanning support
- ⚠️ Paid subscription

---

## Data Flow

```
Scan Library
    ↓
Identify via ISBN/filename
    ↓
Fetch from Open Library
    ↓
Fallback to Google Books
    ↓
Enrich with Goodreads ratings
    ↓
Download cover artwork
```

---

## Configuration

```yaml
metadata:
  books:
    primary: openlibrary
    fallback: [googlebooks, isbndb]
    enrichment:
      - goodreads
```

---

## Related Documentation

- [Metadata Overview](../INDEX.md)
- [Audiobooks](../../audiobook/INDEX.md)
- [Comics](../comics/INDEX.md)
