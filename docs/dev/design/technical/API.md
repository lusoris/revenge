# API Reference

<!-- SOURCES: ogen -->

<!-- DESIGN: technical, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
  - [Base URL](#base-url)
  - [OpenAPI Specifications](#openapi-specifications)
  - [Code Generation](#code-generation)
- [Authentication](#authentication)
  - [Token Flow](#token-flow)
- [Health Endpoints](#health-endpoints)
  - [Health Response](#health-response)
- [Content Rating System](#content-rating-system)
- [Error Responses](#error-responses)
- [Adult Content (QAR)](#adult-content-qar)
- [Viewing the API](#viewing-the-api)
  - [Swagger UI](#swagger-ui)
  - [OpenAPI JSON](#openapi-json)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Development](#development)
  - [Adding New Endpoints](#adding-new-endpoints)
  - [Handler Pattern](#handler-pattern)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

This document describes the HTTP API architecture for Revenge.

## Overview

Revenge uses **ogen** for OpenAPI spec-first code generation. The OpenAPI specifications are the source of truth for all API endpoints.

### Base URL

```
http://localhost:8096/api/v1
```

### OpenAPI Specifications

All API endpoints are defined in OpenAPI 3.1 YAML files:

| Spec File | Description |
|-----------|-------------|
| [revenge.yaml](../../../../api/openapi/revenge.yaml) | Main entrypoint (references all other specs) |
| [auth.yaml](../../../../api/openapi/auth.yaml) | Authentication endpoints |
| [users.yaml](../../../../api/openapi/users.yaml) | User management |
| [libraries.yaml](../../../../api/openapi/libraries.yaml) | Library management |
| [movies.yaml](../../../../api/openapi/movies.yaml) | Movie content |
| [tvshows.yaml](../../../../api/openapi/tvshows.yaml) | TV show content |
| [qar.yaml](../../../../api/openapi/qar.yaml) | Adult content (Queen Anne's Revenge) |
| [system.yaml](../../../../api/openapi/system.yaml) | System health and info |
| [components/schemas.yaml](../../../../api/openapi/components/schemas.yaml) | Shared schemas |

### Code Generation

API handlers are generated using ogen:

```bash
go generate ./api/...
```

Generated code is placed in `api/generated/`. Handler implementations are in `internal/api/`.

---

## Authentication

Most endpoints require authentication via Bearer token:

```
Authorization: Bearer <access_token>
```

Or via custom header:

```
X-Revenge-Token: <access_token>
```

### Token Flow

1. **Login**: `POST /auth/login` → Returns `access_token` + `refresh_token`
2. **Refresh**: `POST /auth/refresh` → Returns new tokens
3. **Logout**: `POST /auth/logout` → Invalidates tokens

---

## Health Endpoints

Health endpoints are outside the `/api/v1` prefix for Kubernetes compatibility:

| Endpoint | Description |
|----------|-------------|
| `GET /health/live` | Liveness probe (always returns 200 OK) |
| `GET /health/ready` | Readiness probe (checks all dependencies) |
| `GET /health` | Detailed health status (JSON) |
| `GET /health/db` | Database pool statistics |
| `GET /version` | Build version info |

### Health Response

```json
{
  "status": "healthy",
  "checks": {
    "database": {"healthy": true, "latency_ms": 1},
    "cache": {"healthy": true, "latency_ms": 0},
    "search": {"healthy": true, "latency_ms": 2},
    "jobs": {"healthy": true, "latency_ms": 0}
  }
}
```

Status values: `healthy`, `degraded`, `unhealthy`

---

## Content Rating System

Revenge uses a normalized rating system (0-100 scale):

| Level | Age | Examples |
|-------|-----|----------|
| 0 | 0+ | G, FSK 0, U |
| 25 | 6+ | PG, FSK 6 |
| 50 | 12+ | PG-13, FSK 12 |
| 75 | 16+ | R, FSK 16 |
| 90 | 18+ | NC-17, FSK 18 |
| 100 | 18+ | Adult/XXX |

---

## Error Responses

All errors return a consistent JSON format:

```json
{
  "code": "not_found",
  "message": "Resource not found"
}
```

**HTTP Status Codes:**

| Code | Description |
|------|-------------|
| 400 | Bad Request (validation errors) |
| 401 | Unauthorized (missing/invalid auth) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (duplicate resource) |
| 500 | Internal Server Error |

---

## Adult Content (QAR)

Adult content endpoints use obfuscated terminology per the Queen Anne's Revenge (QAR) specification:

| Public Term | Internal Code |
|-------------|---------------|
| Adult Movie | Expedition |
| Scene | Voyage |
| Performer | Crew |
| Studio | Port |
| Tag | Flag |
| Library | Fleet |

All QAR endpoints require the `adult.browse` permission and are prefixed with `/qar/`.

See [ADULT_CONTENT_SYSTEM.md](../features/adult/ADULT_CONTENT_SYSTEM.md) for full documentation.

---

## Viewing the API

### Swagger UI

When running in development mode, Swagger UI is available at:

```
http://localhost:8096/swagger/
```

### OpenAPI JSON

The compiled OpenAPI spec is available at:

```
http://localhost:8096/api/v1/openapi.json
```


---

## Development

### Adding New Endpoints

1. Define endpoint in appropriate `api/openapi/*.yaml` file
2. Run `go generate ./api/...`
3. Implement handler in `internal/api/`
4. Wire handler method to generated interface

### Handler Pattern

```go
func (h *Handler) GetUser(ctx context.Context, params gen.GetUserParams) (gen.GetUserRes, error) {
    user, err := h.userService.GetByID(ctx, params.UserID)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return &gen.GetUserNotFound{
                Code:    "not_found",
                Message: "User not found",
            }, nil
        }
        return nil, err
    }
    return &gen.User{...}, nil
}
```
