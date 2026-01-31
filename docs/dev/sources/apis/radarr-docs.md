# Radarr API Documentation

> Source: https://radarr.video/docs/api/
> Type: stub
> Status: Placeholder - See radarr-openapi.json for full spec

---

## Overview

Radarr is a movie collection manager. The API provides access to movie management, quality profiles, indexers, and download clients.

## API Base URL

```
https://{server}/api/v3
```

## Authentication

API key header:
```
X-Api-Key: {api_key}
```

## Core Endpoints

### Movies

| Endpoint | Description |
|----------|-------------|
| `GET /movie` | List all movies |
| `GET /movie/{id}` | Get movie by ID |
| `POST /movie` | Add movie |
| `PUT /movie/{id}` | Update movie |
| `DELETE /movie/{id}` | Delete movie |

### Queue

| Endpoint | Description |
|----------|-------------|
| `GET /queue` | Get download queue |
| `DELETE /queue/{id}` | Remove from queue |

### Calendar

| Endpoint | Description |
|----------|-------------|
| `GET /calendar` | Get upcoming releases |

### Commands

| Endpoint | Description |
|----------|-------------|
| `POST /command` | Execute command |
| `GET /command` | List commands |

## Webhook Events

- `OnGrab`
- `OnDownload`
- `OnMovieFileDelete`
- `OnMovieDelete`
- `OnRename`
- `OnHealthIssue`

## Related

- [Radarr OpenAPI Spec](radarr-openapi.json)
- [Radarr Integration](../../design/integrations/servarr/RADARR.md)
- [Movie Module](../../design/features/video/MOVIE_MODULE.md)
