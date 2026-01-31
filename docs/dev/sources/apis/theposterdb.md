# ThePosterDB API

> Source: https://theposterdb.com/api
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

ThePosterDB provides curated custom poster artwork for movies and TV shows, with matching sets and collections.

## API Base URL

```
https://theposterdb.com/api
```

## Authentication

API key authentication via header or query parameter.

## Core Endpoints

### Posters

| Endpoint | Description |
|----------|-------------|
| `GET /posters` | Search posters |
| `GET /posters/{id}` | Get poster by ID |
| `GET /posters/set/{id}` | Get poster set |

### Sets

| Endpoint | Description |
|----------|-------------|
| `GET /sets` | List poster sets |
| `GET /sets/{id}` | Get set details |

### Users

| Endpoint | Description |
|----------|-------------|
| `GET /users/{username}` | Get user profile |
| `GET /users/{username}/posters` | Get user's posters |

## Poster Types

| Type | Description |
|------|-------------|
| **Textless** | Posters without text/logos |
| **Minimal** | Clean, minimalist designs |
| **4K** | 4K Ultra HD branding |
| **Collection** | Matching set posters |

## Related

- [ThePosterDB Integration](../../design/integrations/metadata/video/THEPOSTERDB.md)
- [Movie Module](../../design/features/video/MOVIE_MODULE.md)
