# Servarr Wiki

> Source: https://wiki.servarr.com/
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

The Servarr Wiki provides documentation for the *arr stack: Radarr, Sonarr, Lidarr, Readarr, Prowlarr, and related applications.

## Common API Patterns

All *arr applications share common API patterns:

### Authentication

API key header:
```
X-Api-Key: {api_key}
```

### Base Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/v3/system/status` | System status |
| `GET /api/v3/health` | Health check |
| `GET /api/v3/queue` | Download queue |
| `GET /api/v3/history` | Activity history |

### Webhooks

All *arr apps support webhooks for events:
- `Grab` - Item grabbed
- `Download` - Download complete
- `Rename` - Files renamed
- `Delete` - Item deleted
- `Health` - Health status change

### Common Data Structures

```json
{
  "id": 123,
  "title": "...",
  "path": "/media/...",
  "qualityProfileId": 1,
  "monitored": true,
  "added": "2024-01-01T00:00:00Z"
}
```

## Resources

- [Radarr Wiki](https://wiki.servarr.com/radarr)
- [Sonarr Wiki](https://wiki.servarr.com/sonarr)
- [Lidarr Wiki](https://wiki.servarr.com/lidarr)
- [Readarr Wiki](https://wiki.servarr.com/readarr)

## Related

- [Arr Integration Pattern](../../design/patterns/ARR_INTEGRATION.md)
- [Servarr Integrations](../../design/integrations/servarr/)
