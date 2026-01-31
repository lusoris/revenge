# Sonarr API Documentation

> Source: https://sonarr.tv/docs/api/
> Type: stub
> Status: Placeholder - See sonarr-openapi.json for full spec

---

## Overview

Sonarr is a TV series collection manager. The API provides access to series management, episode tracking, quality profiles, and download clients.

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

### Series

| Endpoint | Description |
|----------|-------------|
| `GET /series` | List all series |
| `GET /series/{id}` | Get series by ID |
| `POST /series` | Add series |
| `PUT /series/{id}` | Update series |
| `DELETE /series/{id}` | Delete series |

### Episodes

| Endpoint | Description |
|----------|-------------|
| `GET /episode` | List episodes |
| `GET /episode/{id}` | Get episode |
| `PUT /episode/{id}` | Update episode |

### Queue

| Endpoint | Description |
|----------|-------------|
| `GET /queue` | Get download queue |
| `DELETE /queue/{id}` | Remove from queue |

### Calendar

| Endpoint | Description |
|----------|-------------|
| `GET /calendar` | Get upcoming episodes |

## Webhook Events

- `OnGrab`
- `OnDownload`
- `OnEpisodeFileDelete`
- `OnSeriesDelete`
- `OnRename`
- `OnHealthIssue`

## Related

- [Sonarr OpenAPI Spec](sonarr-openapi.json)
- [Sonarr Integration](../../design/integrations/servarr/SONARR.md)
- [TV Show Module](../../design/features/video/TVSHOW_MODULE.md)
