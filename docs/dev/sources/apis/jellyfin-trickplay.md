# Jellyfin Trickplay API

> Source: https://jellyfin.org/docs/general/server/trickplay.html
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

Trickplay provides video thumbnail previews for scrubbing/seeking. The API manages trickplay image generation and retrieval.

## Core Endpoints

### Generation

| Endpoint | Description |
|----------|-------------|
| `POST /Items/{id}/Trickplay/Generate` | Generate trickplay images |
| `GET /Items/{id}/Trickplay/Status` | Get generation status |

### Retrieval

| Endpoint | Description |
|----------|-------------|
| `GET /Videos/{id}/Trickplay/{width}/tiles.bif` | Get trickplay BIF file |
| `GET /Videos/{id}/Trickplay/{width}/{index}` | Get individual tile |

## BIF Format

Trickplay uses the BIF (Base Index Frames) format:
- Fixed-size thumbnail grid
- Indexed by timestamp
- Configurable thumbnail width/interval

## Configuration

```json
{
  "TrickplayOptions": {
    "EnableTrickplay": true,
    "TrickplayInterval": 10000,
    "TrickplayWidth": 320,
    "TrickplayQuality": 85
  }
}
```

## Related

- [Trickplay Feature](../../design/features/playback/TRICKPLAY.md)
- [BIF Protocol](../protocols/bif.md)
