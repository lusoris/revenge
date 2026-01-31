# Letterboxd API

> Source: https://api-docs.letterboxd.com/
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

Letterboxd is a social film discovery platform. The API provides access to film data, user profiles, lists, and activity.

## API Base URL

```
https://api.letterboxd.com/api/v0
```

## Authentication

OAuth 2.0 authentication for user-specific operations. API key required for all requests.

## Core Endpoints

### Films

| Endpoint | Description |
|----------|-------------|
| `GET /films` | Search films |
| `GET /film/{id}` | Get film details |
| `GET /film/{id}/statistics` | Get film statistics |
| `GET /film/{id}/relationships` | Get related films |

### Members

| Endpoint | Description |
|----------|-------------|
| `GET /member/{id}` | Get member profile |
| `GET /member/{id}/watchlist` | Get member watchlist |
| `GET /member/{id}/lists` | Get member lists |

### Lists

| Endpoint | Description |
|----------|-------------|
| `GET /lists` | Search lists |
| `GET /list/{id}` | Get list details |
| `GET /list/{id}/entries` | Get list entries |

### Activity

| Endpoint | Description |
|----------|-------------|
| `GET /log-entries` | Get activity entries |
| `POST /film/{id}/log` | Log a film |

## Related

- [Letterboxd Scrobbling](../../design/integrations/scrobbling/LETTERBOXD.md)
- [Scrobbling Feature](../../design/features/shared/SCROBBLING.md)
