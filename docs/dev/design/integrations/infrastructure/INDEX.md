# Core Infrastructure

> Database, cache, search, and job queue

---

## Overview

Revenge's core infrastructure stack:
- **Database**: PostgreSQL for persistent storage
- **Cache**: Dragonfly (Redis-compatible) for caching
- **Search**: Typesense for full-text search
- **Jobs**: River for background processing

---

## Components

| Component | Technology | Purpose | Status |
|-----------|------------|---------|--------|
| [PostgreSQL](POSTGRESQL.md) | PostgreSQL 18+ | Primary database | 🟢 Required |
| [Dragonfly](DRAGONFLY.md) | Dragonfly | Cache & sessions | 🟢 Required |
| [Typesense](TYPESENSE.md) | Typesense 27+ | Full-text search | 🟢 Required |
| [River](RIVER.md) | River | Job queue | 🟢 Required |

---

## Component Details

### PostgreSQL
**Primary data store**

- ✅ All persistent data
- ✅ User accounts and sessions
- ✅ Media metadata
- ✅ Watch history and progress
- ✅ Adult content in `c` schema

### Dragonfly
**In-memory cache**

- ✅ Session storage
- ✅ API response caching
- ✅ Rate limiting
- ✅ Pub/sub for real-time
- ✅ Redis-compatible

### Typesense
**Search engine**

- ✅ Full-text search
- ✅ Typo tolerance
- ✅ Faceted search
- ✅ Vector search (similarity)

### River
**Background jobs**

- ✅ Library scanning
- ✅ Metadata fetching
- ✅ Image downloads
- ✅ Scheduled tasks
- ✅ PostgreSQL-native

---

## Architecture

```
                    ┌──────────────┐
                    │   Clients    │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │   Revenge    │
                    │    Server    │
                    └──────┬───────┘
           ┌───────────────┼───────────────┐
           │               │               │
    ┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
    │  PostgreSQL │ │  Dragonfly  │ │  Typesense  │
    │  (Primary)  │ │  (Cache)    │ │  (Search)   │
    └─────────────┘ └─────────────┘ └─────────────┘
           │
    ┌──────▼──────┐
    │    River    │
    │   (Jobs)    │
    └─────────────┘
```

---

## Docker Compose

```yaml
services:
  revenge:
    depends_on:
      postgres:
        condition: service_healthy
      dragonfly:
        condition: service_healthy
      typesense:
        condition: service_healthy

  postgres:
    image: postgres:18
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U revenge"]

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]

  typesense:
    image: typesense/typesense:27.0
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8108/health"]
```

---

## Configuration

```yaml
database:
  host: "${DB_HOST:localhost}"
  port: ${DB_PORT:5432}
  name: "revenge"
  user: "${DB_USER:revenge}"
  password: "${DB_PASSWORD}"

cache:
  host: "${CACHE_HOST:localhost}"
  port: ${CACHE_PORT:6379}

search:
  host: "${SEARCH_HOST:localhost}"
  port: ${SEARCH_PORT:8108}
  api_key: "${TYPESENSE_API_KEY}"

jobs:
  workers: 10
  queues:
    default: 5
    high: 10
    low: 2
```

---

## Related Documentation

- [Tech Stack](../../technical/TECH_STACK.md)
- [Setup Guide](../../operations/SETUP.md)
