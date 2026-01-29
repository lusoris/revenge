# OpenLibrary Integration

> Open book metadata database - primary metadata provider for books (Goodreads alternative)

**Service**: OpenLibrary (Internet Archive)
**Type**: Metadata Provider (Books)
**API Version**: v1
**Website**: https://openlibrary.org
**API Docs**: https://openlibrary.org/developers/api

---

## Overview

**OpenLibrary** is the **recommended primary book metadata provider** (replaces Goodreads API).

**Why OpenLibrary**:
- Free and open (Internet Archive project)
- Comprehensive book database (20M+ books)
- No authentication required
- Stable RESTful API
- Book covers available
- ISBN, OCLC, LCCN lookup
- Author information
- No rate limits (fair use)

**Use Cases**:
- Book metadata (title, author, ISBN, publication date, pages, subjects)
- Author metadata (biography, birth/death dates, photos)
- Book covers (small, medium, large)
- Edition information (multiple editions of same work)
- Subject/genre classification

---

## Developer Resources

**API Documentation**: https://openlibrary.org/developers/api
**Books API**: https://openlibrary.org/dev/docs/api/books
**Covers API**: https://openlibrary.org/dev/docs/api/covers
**Search API**: https://openlibrary.org/dev/docs/api/search

**Authentication**: None required
**Rate Limit**: No strict limit (respect fair use, ~1 req/s recommended)
**Free Tier**: Unlimited

---

## API Details

### Base URLs
```
https://openlibrary.org/
https://covers.openlibrary.org/
```

### Key Endpoints

#### Get Book by ISBN
```bash
GET /isbn/{ISBN}.json
```

**Example**:
```bash
GET /isbn/9780140328721.json
```

**Response**:
```json
{
  "key": "/books/OL7353617M",
  "title": "Fantastic Mr. Fox",
  "authors": [{"key": "/authors/OL34184A"}],
  "publish_date": "1988",
  "publishers": ["Puffin"],
  "isbn_10": ["0140328726"],
  "isbn_13": ["9780140328721"],
  "number_of_pages": 96,
  "subjects": ["Foxes", "Children's stories"],
  "covers": [6498519],
  "works": [{"key": "/works/OL45804W"}]
}
```

#### Get Work (Canonical Book)
```bash
GET /works/{WORK_ID}.json
```

**Work vs Edition**:
- **Work**: Canonical book (e.g., "1984" by George Orwell)
- **Edition**: Specific publication (e.g., 1984 UK 1st edition, 1984 US paperback)

**Example**:
```bash
GET /works/OL45804W.json
```

**Response**:
```json
{
  "key": "/works/OL45804W",
  "title": "Fantastic Mr. Fox",
  "authors": [{"author": {"key": "/authors/OL34184A"}, "type": {"key": "/type/author_role"}}],
  "description": "Boggis, Bunce and Bean are the meanest...",
  "subjects": ["Foxes", "Stealing", "Fathers and sons", "Children's stories"],
  "covers": [6498519]
}
```

#### Get Author
```bash
GET /authors/{AUTHOR_ID}.json
```

**Example**:
```bash
GET /authors/OL34184A.json
```

**Response**:
```json
{
  "key": "/authors/OL34184A",
  "name": "Roald Dahl",
  "birth_date": "13 September 1916",
  "death_date": "23 November 1990",
  "bio": "Roald Dahl was a British novelist...",
  "photos": [6498519],
  "wikipedia": "https://en.wikipedia.org/wiki/Roald_Dahl"
}
```

#### Search Books
```bash
GET /search.json?q={QUERY}&limit=10
GET /search.json?title=fantastic+mr+fox&author=roald+dahl
GET /search.json?isbn=9780140328721
```

**Response**:
```json
{
  "numFound": 1,
  "docs": [
    {
      "key": "/works/OL45804W",
      "title": "Fantastic Mr. Fox",
      "author_name": ["Roald Dahl"],
      "first_publish_year": 1970,
      "isbn": ["9780140328721", "0140328726"],
      "cover_i": 6498519,
      "publisher": ["Puffin"],
      "language": ["eng"],
      "subject": ["Foxes", "Children's stories"]
    }
  ]
}
```

#### Get Cover Image
```bash
GET https://covers.openlibrary.org/b/isbn/{ISBN}-L.jpg
GET https://covers.openlibrary.org/b/id/{COVER_ID}-L.jpg
```

**Cover Sizes**:
- `-S.jpg`: Small (thumbnail)
- `-M.jpg`: Medium (recommended for lists)
- `-L.jpg`: Large (recommended for detail pages)

**Example**:
```bash
GET https://covers.openlibrary.org/b/isbn/9780140328721-L.jpg
GET https://covers.openlibrary.org/b/id/6498519-L.jpg
```

---

## Implementation Checklist

### API Client (`internal/infra/metadata/provider_openlibrary.go`)
- [ ] Base URL configuration
- [ ] HTTP client with User-Agent
- [ ] Rate limiting (1 req/s recommended, no strict limit)
- [ ] Error handling (404: Book not found, 500: Server error)
- [ ] Response parsing (JSON unmarshalling)

### Book Metadata
- [ ] Fetch book by ISBN
- [ ] Fetch work (canonical book)
- [ ] Search books by title, author, ISBN
- [ ] Extract: title, author, ISBN, publication date, pages, subjects, publishers
- [ ] Store in `books` table

### Author Metadata
- [ ] Fetch author by OpenLibrary ID
- [ ] Extract: name, biography, birth/death dates, photos, Wikipedia link
- [ ] Store in `book_authors` table

### Cover Art Handling
- [ ] Fetch cover by ISBN or cover ID
- [ ] Download large cover (-L.jpg)
- [ ] Generate Blurhash
- [ ] Convert to WebP
- [ ] Store locally (`data/books/covers/`)

### Work vs Edition Handling
- [ ] Store work ID (canonical book)
- [ ] Store edition ID (specific publication)
- [ ] Group editions by work (multiple editions of same book)

### Error Handling
- [ ] Handle 404 (Book not found)
- [ ] Handle 500 (Server error - retry)
- [ ] Log errors (no sensitive data)

---

## Integration Pattern

### Book Metadata Workflow
```go
// Fetch book metadata from OpenLibrary
func (s *BookService) FetchBookMetadata(isbn string) error {
    // 1. Fetch book by ISBN
    edition := s.openlibraryClient.GetBookByISBN(isbn)
    if edition == nil {
        return errors.New("book not found")
    }

    // 2. Fetch work (canonical book)
    work := s.openlibraryClient.GetWork(edition.Works[0].Key)

    // 3. Fetch author
    author := s.openlibraryClient.GetAuthor(edition.Authors[0].Key)

    // 4. Fetch cover
    coverURL := fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn)
    coverPath := s.downloadCover(coverURL)

    // 5. Store in database
    s.db.InsertAuthor(author)
    s.db.InsertBook(map[string]interface{}{
        "openlibrary_id":   edition.Key,
        "work_id":          work.Key,
        "title":            edition.Title,
        "author_id":        author.ID,
        "isbn_10":          edition.ISBN10[0],
        "isbn_13":          edition.ISBN13[0],
        "publish_date":     edition.PublishDate,
        "pages":            edition.NumberOfPages,
        "subjects":         work.Subjects,
        "description":      work.Description,
        "cover_path":       coverPath,
    })

    return nil
}
```

### Readarr Integration
```go
// Readarr webhook → OpenLibrary metadata sync
func (s *BookService) HandleReadarrBookAdded(bookID string) error {
    // 1. Get book from Readarr
    readarrBook := s.readarrClient.GetBook(bookID)
    isbn := readarrBook.ISBN // ISBN-13 or ISBN-10

    // 2. Fetch metadata from OpenLibrary
    s.FetchBookMetadata(isbn)

    return nil
}
```

---

## Related Documentation

- **Book Module**: [MODULE_IMPLEMENTATION_TODO.md](../../../planning/MODULE_IMPLEMENTATION_TODO.md) (Book section)
- **Goodreads Integration**: [GOODREADS.md](GOODREADS.md) (API retired, use OpenLibrary)
- **Readarr Integration**: [../servarr/READARR.md](../servarr/READARR.md)
- **Hardcover Integration**: [HARDCOVER.md](HARDCOVER.md) (social reading platform)

---

## Notes

- **No authentication required**: Public API (no API key)
- **No strict rate limits**: Fair use (recommend 1 req/s to avoid overload)
- **Work vs Edition**: Store both (work = canonical, edition = specific publication)
- **ISBN lookup**: Supports both ISBN-10 and ISBN-13
- **OCLC/LCCN lookup**: Also supported (alternative identifiers)
- **Cover images**: Free CDN (download separately, no rate limit)
- **Subjects**: User-generated tags (similar to genres)
- **Multiple editions**: Same work, different editions (group by work ID)
- **Author photos**: Available (download from `/authors/{ID}/photos`)
- **Wikipedia links**: Many authors have Wikipedia links
- **API stable**: v1 stable, no breaking changes expected
- **Free and open**: Internet Archive project (non-profit)
- **Data quality**: Community-driven (some books have incomplete data)
- **Fallback strategy**: OpenLibrary primary, Google Books API fallback
- **Migration from Goodreads**: Use ISBN as common identifier (cross-reference Goodreads ID → ISBN → OpenLibrary ID)
- **Search syntax**: Lucene-based (use `title:`, `author:`, `isbn:` prefixes)
- **Response formats**: JSON (default), YAML, RDF
