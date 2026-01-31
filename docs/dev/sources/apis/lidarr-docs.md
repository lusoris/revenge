# Lidarr API Documentation

> Source: https://lidarr.audio/docs/api/
> Type: stub
> Status: Placeholder - See lidarr-openapi.json for full spec

---

## Overview

Lidarr is a music collection manager. The API provides access to artist/album management, quality profiles, and download clients.

## API Base URL

```
https://{server}/api/v1
```

## Authentication

API key header:
```
X-Api-Key: {api_key}
```

## Core Endpoints

### Artists

| Endpoint | Description |
|----------|-------------|
| `GET /artist` | List all artists |
| `GET /artist/{id}` | Get artist by ID |
| `POST /artist` | Add artist |
| `PUT /artist/{id}` | Update artist |
| `DELETE /artist/{id}` | Delete artist |

### Albums

| Endpoint | Description |
|----------|-------------|
| `GET /album` | List albums |
| `GET /album/{id}` | Get album |
| `PUT /album/{id}` | Update album |

### Tracks

| Endpoint | Description |
|----------|-------------|
| `GET /track` | List tracks |
| `GET /track/{id}` | Get track |

### Queue

| Endpoint | Description |
|----------|-------------|
| `GET /queue` | Get download queue |
| `DELETE /queue/{id}` | Remove from queue |

## Webhook Events

- `OnGrab`
- `OnReleaseImport`
- `OnTrackFileDelete`
- `OnAlbumDelete`
- `OnRename`
- `OnHealthIssue`

## Related

- [Lidarr OpenAPI Spec](lidarr-openapi.json)
- [Lidarr Integration](../../design/integrations/servarr/LIDARR.md)
- [Music Module](../../design/features/music/MUSIC_MODULE.md)
