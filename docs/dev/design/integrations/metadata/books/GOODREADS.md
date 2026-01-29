# Goodreads Integration

> Book metadata and social reading platform - ratings, reviews, book lists

**Service**: Goodreads (Amazon)
**Type**: Metadata Provider (Books)
**API Version**: N/A (API retired 2020)
**Website**: https://www.goodreads.com
**Alternative**: Web scraping OR OpenLibrary

---

## Overview

**Goodreads** was the primary book metadata provider, but the **API was retired in December 2020**. Alternative approaches required.

**Why Goodreads (historical)**:
- Comprehensive book database
- User ratings and reviews
- Book lists and recommendations
- Author information
- Reading challenges

**Current Status**:
- ❌ **API retired** (December 2020)
- ✅ **Web scraping** (legal gray area, fragile)
- ✅ **OpenLibrary** (open alternative, recommended)

**Use Cases**:
- Book metadata (title, author, ISBN, publication date)
- User ratings and reviews
- Book covers
- Author information
- Book recommendations

---

## Developer Resources

**API Status**: RETIRED (December 2020)
**Official Announcement**: https://www.goodreads.com/api

**Alternatives**:
1. **OpenLibrary** (recommended): https://openlibrary.org/developers/api
2. **Web scraping** (fragile, not recommended)
3. **Google Books API**: https://developers.google.com/books

---

## Migration Strategy

### Option 1: OpenLibrary (Recommended)
Use **OpenLibrary** as primary book metadata source.

**Advantages**:
- Free and open API
- Comprehensive book database (Open Library, Internet Archive)
- Book covers available
- No authentication required
- Stable API

**See**: [OPENLIBRARY.md](OPENLIBRARY.md)

### Option 2: Google Books API
Use **Google Books API** as fallback.

**Advantages**:
- Official Google API (stable)
- Book metadata, covers, previews
- Free tier (1000 requests/day)

**Disadvantages**:
- Limited to Google's catalog
- Requires API key

### Option 3: Web Scraping (NOT Recommended)
Scrape Goodreads website for metadata.

**Disadvantages**:
- Legal gray area (ToS violation)
- Fragile (breaks when HTML changes)
- No official support
- Rate limiting required (respect robots.txt)
- CAPTCHA challenges

---

## Historical API Details (Archived)

### Base URL (RETIRED)
```
https://www.goodreads.com/
```

### Authentication (RETIRED)
- API Key (registration closed)
- OAuth 1.0a (deprecated)

### Key Endpoints (RETIRED)
```bash
# Search books (RETIRED)
GET /search/index.xml?key={API_KEY}&q={QUERY}

# Get book (RETIRED)
GET /book/show/{BOOK_ID}.xml?key={API_KEY}

# Get author (RETIRED)
GET /author/show/{AUTHOR_ID}.xml?key={API_KEY}
```

---

## Implementation Checklist

### Migration to OpenLibrary
- [ ] **REMOVE Goodreads API client** (API retired)
- [ ] **Implement OpenLibrary client** (see [OPENLIBRARY.md](OPENLIBRARY.md))
- [ ] Update book search to use OpenLibrary
- [ ] Update book metadata fetching to use OpenLibrary
- [ ] Migrate existing Goodreads IDs to ISBNs (cross-reference)

### Web Scraping Fallback (Optional, NOT Recommended)
- [ ] Implement web scraper (BeautifulSoup, Playwright)
- [ ] Respect robots.txt
- [ ] Rate limiting (1 req/5s minimum)
- [ ] User-Agent header
- [ ] CAPTCHA handling (fail silently)
- [ ] Graceful degradation (OpenLibrary fallback)

### Google Books API Fallback (Optional)
- [ ] Implement Google Books API client
- [ ] API Key configuration
- [ ] Search books by ISBN, title, author
- [ ] Fetch book metadata, covers
- [ ] Rate limiting (1000 req/day)

---

## Integration Pattern (Historical)

### Book Metadata Workflow (Archived)
```go
// ARCHIVED: Goodreads API retired
func (s *BookService) FetchBookMetadata(isbn string) error {
    // OPTION 1: OpenLibrary (recommended)
    book := s.openlibraryClient.GetBookByISBN(isbn)
    if book != nil {
        s.db.InsertBook(book)
        return nil
    }

    // OPTION 2: Google Books API (fallback)
    book = s.googlebooksClient.GetBookByISBN(isbn)
    if book != nil {
        s.db.InsertBook(book)
        return nil
    }

    // OPTION 3: Web scraping (NOT recommended)
    // book = s.goodreadsScraper.GetBookByISBN(isbn)

    return errors.New("book not found")
}
```

---

## Related Documentation

- **Book Module**: [MODULE_IMPLEMENTATION_TODO.md](../../../planning/MODULE_IMPLEMENTATION_TODO.md) (Book section)
- **OpenLibrary Integration**: [OPENLIBRARY.md](OPENLIBRARY.md) (recommended alternative)
- **Chaptarr Integration**: [../../servarr/CHAPTARR.md](../../servarr/CHAPTARR.md)
- **Hardcover Integration**: [HARDCOVER.md](HARDCOVER.md) (social reading platform)

---

## Notes

- **API retired**: December 2020 (Amazon decision)
- **No new API keys**: Registration closed
- **Existing API keys**: Still work (for now, may be disabled in future)
- **Migration required**: Use OpenLibrary OR Google Books API
- **Web scraping**: Legal gray area, fragile, NOT recommended
- **OpenLibrary recommended**: Free, open, stable API
- **Goodreads IDs**: Store in `books.goodreads_id` (historical reference, cross-reference with ISBNs)
- **User data migration**: Export Goodreads reading lists via CSV (user-initiated)
- **Alternative platforms**: Hardcover, LibraryThing, StoryGraph (API availability varies)
- **Amazon ownership**: Goodreads owned by Amazon (2013), API sunset likely strategic decision
- **Community frustration**: Developers migrated to OpenLibrary, Hardcover, LibraryThing
- **Future**: Monitor for API resurrection (unlikely), continue with OpenLibrary
- **Recommendation**: **Use OpenLibrary as primary, Google Books as fallback**
