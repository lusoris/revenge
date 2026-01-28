# Hardcover Integration

> Social reading platform - Goodreads alternative with API

**Service**: Hardcover
**Type**: Social Reading + Metadata Provider (Books)
**API Version**: GraphQL
**Website**: https://hardcover.app
**API Docs**: https://hardcover.app/docs/api

---

## Overview

**Hardcover** is a modern social reading platform with an **official GraphQL API** (Goodreads alternative).

**Why Hardcover**:
- **Official GraphQL API** (unlike Goodreads)
- Social reading features (reading lists, reviews, ratings)
- Book metadata (title, author, ISBN, publication date)
- User reading challenges
- Book recommendations
- Active development (indie project)

**Use Cases**:
- Book metadata (fallback to OpenLibrary)
- Social features (user reading lists, reviews, ratings)
- Reading challenges integration
- Book recommendations
- User profiles (reading history, favorites)

**Status**:
- ✅ **Official GraphQL API** (active development)
- ✅ **Free tier** (generous limits)
- ✅ **Open source roadmap** (community-driven)

---

## Developer Resources

**Website**: https://hardcover.app
**API Documentation**: https://hardcover.app/docs/api
**GraphQL Playground**: https://api.hardcover.app/graphiql

**Authentication**: API Key (registration required)
**Rate Limit**: Generous (no strict public limit)
**Free Tier**: Unlimited (API key required)

---

## API Details

### Base URL
```
https://api.hardcover.app/v1/graphql
```

### Authentication
API Key (Header):
```
Authorization: Bearer {API_KEY}
```

**API Key Registration**: https://hardcover.app/settings/api

### GraphQL Schema

#### Search Books
```graphql
query SearchBooks($query: String!, $limit: Int) {
  search(query: $query, limit: $limit) {
    books {
      id
      title
      subtitle
      description
      isbn_10
      isbn_13
      pages
      release_date
      image
      authors {
        id
        name
        bio
        image
      }
      contributions {
        author {
          name
        }
        role
      }
      user_book {
        status
        rating
      }
    }
  }
}
```

#### Get Book by ID
```graphql
query GetBook($id: Int!) {
  book(id: $id) {
    id
    title
    subtitle
    description
    isbn_10
    isbn_13
    pages
    release_date
    image
    authors {
      id
      name
      bio
      image
    }
    series {
      id
      name
      position
    }
    editions {
      id
      title
      isbn_13
      release_date
    }
    ratings {
      average
      count
    }
    reviews {
      id
      body
      rating
      user {
        name
        image
      }
    }
  }
}
```

#### Get Book by ISBN
```graphql
query GetBookByISBN($isbn: String!) {
  books(where: {isbn_13: {_eq: $isbn}}) {
    id
    title
    # ... same fields as GetBook
  }
}
```

#### Get User Reading List
```graphql
query GetUserBooks($userId: Int!, $status: String) {
  user_books(where: {user_id: {_eq: $userId}, status: {_eq: $status}}) {
    book {
      id
      title
      image
    }
    status
    rating
    started_at
    finished_at
  }
}
```

**Status Values**:
- `WANT_TO_READ`
- `CURRENTLY_READING`
- `READ`
- `DID_NOT_FINISH`

#### Add Book to User List
```graphql
mutation AddBookToList($bookId: Int!, $status: String!, $rating: Int) {
  insert_user_books_one(object: {book_id: $bookId, status: $status, rating: $rating}) {
    id
    status
    rating
  }
}
```

---

## Implementation Checklist

### API Client (`internal/infra/metadata/provider_hardcover.go`)
- [ ] Base URL configuration
- [ ] API Key configuration (Authorization header)
- [ ] GraphQL client (HTTP POST to /v1/graphql)
- [ ] Error handling (401: Invalid API key, 404: Book not found)
- [ ] Response parsing (JSON unmarshalling)

### Book Metadata (Fallback)
- [ ] Search books by title, author, ISBN
- [ ] Fetch book by ID or ISBN
- [ ] Extract: title, subtitle, description, ISBN, pages, release date, authors
- [ ] Store in `books` table (fallback to OpenLibrary)

### Social Features
- [ ] Fetch user reading list (Want to Read, Currently Reading, Read)
- [ ] Add book to user list
- [ ] Update reading status
- [ ] Rate book
- [ ] Write review (optional)

### Author Metadata
- [ ] Fetch author by ID
- [ ] Extract: name, bio, image
- [ ] Store in `book_authors` table

### Series Handling
- [ ] Fetch series information
- [ ] Store in `book_series` table
- [ ] Link books to series (position in series)

### Error Handling
- [ ] Handle 401 (Invalid API key - check configuration)
- [ ] Handle 404 (Book not found)
- [ ] Handle GraphQL errors (parse `errors` array in response)
- [ ] Log errors (no sensitive data)

---

## Integration Pattern

### Sync User Reading List
```go
// Sync user's Hardcover reading list to Revenge
func (s *BookService) SyncHardcoverReadingList(ctx context.Context, userID uuid.UUID) error {
    // 1. Get user's Hardcover API key
    hardcoverKey := s.db.GetUserHardcoverKey(userID)
    if hardcoverKey == "" {
        return errors.New("user not connected to Hardcover")
    }

    // 2. Fetch user's reading list from Hardcover
    hardcoverUserID := s.hardcoverClient.GetUserID(hardcoverKey)
    readingList := s.hardcoverClient.GetUserBooks(hardcoverUserID, "CURRENTLY_READING")

    // 3. Store in Revenge database
    for _, userBook := range readingList {
        book := userBook.Book

        // Ensure book exists in database
        if !s.db.BookExists(book.ISBN13) {
            s.db.InsertBook(map[string]interface{}{
                "hardcover_id": book.ID,
                "title":        book.Title,
                "isbn_13":      book.ISBN13,
                "release_date": book.ReleaseDate,
            })
        }

        // Add to user's reading list
        s.db.InsertUserReadingList(userID, book.ISBN13, map[string]interface{}{
            "status":      "currently_reading",
            "rating":      userBook.Rating,
            "started_at":  userBook.StartedAt,
        })
    }

    return nil
}
```

### Add Book to Hardcover List
```go
// User marks book as "Want to Read" → Add to Hardcover
func (s *BookService) AddToHardcoverWantToRead(userID uuid.UUID, bookID uuid.UUID) error {
    // 1. Get user's Hardcover API key
    hardcoverKey := s.db.GetUserHardcoverKey(userID)
    if hardcoverKey == "" {
        return nil // User not connected, skip
    }

    // 2. Get book
    book := s.db.GetBook(bookID)

    // 3. Search Hardcover for book (by ISBN)
    hardcoverBook := s.hardcoverClient.GetBookByISBN(hardcoverKey, book.ISBN13)
    if hardcoverBook == nil {
        return errors.New("book not found on Hardcover")
    }

    // 4. Add to user's list
    s.hardcoverClient.AddBookToList(hardcoverKey, hardcoverBook.ID, "WANT_TO_READ", nil)

    return nil
}
```

---

## Related Documentation

- **Book Module**: [docs/MODULE_IMPLEMENTATION_TODO.md](../../MODULE_IMPLEMENTATION_TODO.md) (Book section)
- **OpenLibrary Integration**: [OPENLIBRARY.md](OPENLIBRARY.md) (primary metadata)
- **Goodreads Integration**: [GOODREADS.md](GOODREADS.md) (API retired)
- **Readarr Integration**: [../servarr/READARR.md](../servarr/READARR.md)

---

## Notes

- **Official GraphQL API**: Modern alternative to Goodreads (which has retired API)
- **API Key required**: Register at https://hardcover.app/settings/api
- **Free tier**: Generous limits (no strict public limit, fair use)
- **GraphQL**: Single endpoint, flexible queries (specify exactly what you need)
- **Social features**: Reading lists, reviews, ratings, reading challenges
- **Metadata quality**: Community-driven (some books may have incomplete data)
- **OpenLibrary primary**: Use Hardcover for social features, OpenLibrary for comprehensive metadata
- **User connection**: Users connect Hardcover account via OAuth or API key
- **Sync strategy**: Two-way sync (Revenge ↔ Hardcover reading lists)
- **Reading statuses**: Want to Read, Currently Reading, Read, Did Not Finish
- **Ratings**: 1-5 stars (half-stars supported)
- **Reviews**: Text reviews (optional, can be synced)
- **Series support**: Books linked to series (position in series)
- **Editions**: Multiple editions of same book (group by work)
- **Author bios**: Available (store in `book_authors` table)
- **Cover images**: Available (download and store locally)
- **GraphQL playground**: https://api.hardcover.app/graphiql (test queries)
- **Error handling**: GraphQL errors in `errors` array (parse separately from `data`)
- **Indie project**: Developed by small team, active development
- **Community**: Growing community (Reddit, Discord)
- **Future**: Potential for deeper integration (challenges, recommendations, friends)
