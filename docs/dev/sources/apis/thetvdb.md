# TheTVDB API v4

> Source: https://thetvdb.github.io/v4-api/
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

TheTVDB provides comprehensive TV show, movie, and person metadata through their v4 REST API.

## API Base URL

```
https://api4.thetvdb.com/v4
```

## Authentication

JWT-based authentication:

```http
POST /login
Content-Type: application/json

{
  "apikey": "your-api-key"
}
```

Response includes a JWT token for subsequent requests.

## Core Endpoints

### Series

| Endpoint | Description |
|----------|-------------|
| `GET /series/{id}` | Get series by ID |
| `GET /series/{id}/extended` | Get extended series info |
| `GET /series/{id}/episodes` | Get episodes for series |
| `GET /series/{id}/artworks` | Get series artwork |
| `GET /search` | Search for series |

### Episodes

| Endpoint | Description |
|----------|-------------|
| `GET /episodes/{id}` | Get episode by ID |
| `GET /episodes/{id}/extended` | Get extended episode info |

### People

| Endpoint | Description |
|----------|-------------|
| `GET /people/{id}` | Get person by ID |
| `GET /people/{id}/extended` | Get extended person info |

### Artwork

| Endpoint | Description |
|----------|-------------|
| `GET /artwork/{id}` | Get artwork by ID |
| `GET /artwork/types` | List artwork types |

## Artwork Types

| Type | Description |
|------|-------------|
| `poster` | Series poster (vertical) |
| `fanart` | Backdrop/fanart (horizontal) |
| `banner` | Series banner (wide) |
| `icon` | Series icon |
| `clearlogo` | Clear logo |
| `clearart` | Clear art |

## Related

- [TheTVDB Integration](../../design/integrations/metadata/video/THETVDB.md)
- [TV Show Module](../../design/features/video/TVSHOW_MODULE.md)
