# ThePornDB API

> Source: https://theporndb.net/api
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

ThePornDB provides adult content metadata through a REST API, including scene, performer, and studio information.

## API Base URL

```
https://api.theporndb.net
```

## Authentication

API key authentication via header:
```
Authorization: Bearer {api_key}
```

## Core Endpoints

### Scenes

| Endpoint | Description |
|----------|-------------|
| `GET /scenes` | Search scenes |
| `GET /scenes/{id}` | Get scene by ID |
| `GET /scenes/parse` | Parse filename for metadata |

### Performers

| Endpoint | Description |
|----------|-------------|
| `GET /performers` | Search performers |
| `GET /performers/{id}` | Get performer by ID |

### Studios/Sites

| Endpoint | Description |
|----------|-------------|
| `GET /sites` | List sites/studios |
| `GET /sites/{id}` | Get site by ID |

## Related

- [ThePornDB Integration](../../design/integrations/metadata/adult/THEPORNDB.md)
- [Adult Content System](../../design/features/adult/ADULT_CONTENT_SYSTEM.md)
