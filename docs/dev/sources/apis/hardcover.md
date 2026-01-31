# Hardcover API

> Source: https://hardcover.app/docs/api
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

Hardcover is a modern book tracking and discovery platform with a GraphQL API for accessing book metadata and user libraries.

## API Base URL

```
https://api.hardcover.app/v1/graphql
```

## Authentication

OAuth 2.0 authentication required for user-specific operations.

## GraphQL Schema

### Queries

```graphql
type Query {
  book(id: ID!): Book
  books(search: String, limit: Int): [Book]
  author(id: ID!): Author
  user(username: String!): User
}
```

### Types

```graphql
type Book {
  id: ID!
  title: String!
  authors: [Author]
  description: String
  cover_url: String
  published_date: String
  isbn_10: String
  isbn_13: String
  page_count: Int
}
```

## Related

- [Book Module](../../design/features/book/BOOK_MODULE.md)
- [Goodreads Integration](../../design/integrations/metadata/books/GOODREADS.md)
