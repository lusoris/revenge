# Go Packages Reference

> Reusable Go packages for Revenge development

---

## Overview

This document lists all recommended Go packages organized by functionality. Using established packages reduces code, improves reliability, and speeds development.

---

## Core Framework

| Package | Purpose | URL |
|---------|---------|-----|
| **go-chi/chi** | HTTP router | github.com/go-chi/chi/v5 |
| **uber-go/fx** | Dependency injection | go.uber.org/fx |
| **riverqueue/river** | Background jobs | github.com/riverqueue/river |
| **jackc/pgx** | PostgreSQL driver | github.com/jackc/pgx/v5 |

---

## Authentication & Authorization

| Package | Purpose | URL |
|---------|---------|-----|
| **casbin/casbin** | RBAC/ABAC authorization | github.com/casbin/casbin/v2 |
| **golang-jwt/jwt** | JWT tokens | github.com/golang-jwt/jwt/v5 |
| **gorilla/sessions** | Session management | github.com/gorilla/sessions |
| **markbates/goth** | OAuth providers | github.com/markbates/goth |

---

## Database & ORM

| Package | Purpose | URL |
|---------|---------|-----|
| **sqlc-dev/sqlc** | Type-safe SQL | github.com/sqlc-dev/sqlc |
| **pressly/goose** | Database migrations | github.com/pressly/goose/v3 |
| **Masterminds/squirrel** | SQL query builder | github.com/Masterminds/squirrel |

---

## Content & Markdown

| Package | Purpose | URL |
|---------|---------|-----|
| **yuin/goldmark** | Markdown parser (CommonMark) | github.com/yuin/goldmark |
| **goldmark-highlighting** | Syntax highlighting | github.com/yuin/goldmark-highlighting |
| **goldmark-meta** | YAML frontmatter | github.com/yuin/goldmark-meta |
| **goldmark-wikilink** | [[wiki]] links | go.abhg.dev/goldmark/wikilink |

---

## Full-Text Search

| Package | Purpose | URL |
|---------|---------|-----|
| **blevesearch/bleve** | Embedded search (Go native) | github.com/blevesearch/bleve/v2 |
| **meilisearch-go** | Meilisearch client | github.com/meilisearch/meilisearch-go |
| **typesense-go** | Typesense client | github.com/typesense/typesense-go |

---

## News & Feeds

| Package | Purpose | URL |
|---------|---------|-----|
| **mmcdole/gofeed** | RSS/Atom/JSON parsing | github.com/mmcdole/gofeed |
| **gorilla/feeds** | Feed generation | github.com/gorilla/feeds |

---

## Real-Time & Notifications

| Package | Purpose | URL |
|---------|---------|-----|
| **gorilla/websocket** | WebSocket server | github.com/gorilla/websocket |
| **r3labs/sse** | Server-Sent Events | github.com/r3labs/sse/v2 |
| **centrifugal/centrifuge** | Real-time messaging | github.com/centrifugal/centrifuge |

---

## Voting & Polls

| Package | Purpose | URL |
|---------|---------|-----|
| **Sam-Izdat/govote** | Voting algorithms | github.com/Sam-Izdat/govote |
| **FabianWe/gopolls** | Poll procedures | github.com/FabianWe/gopolls |

> Note: For production, custom implementation with River jobs recommended.

---

## Image Processing

| Package | Purpose | URL |
|---------|---------|-----|
| **h2non/bimg** | Fast processing (libvips) | github.com/h2non/bimg |
| **disintegration/imaging** | Pure Go imaging | github.com/disintegration/imaging |
| **bbrks/go-blurhash** | Blurhash generation | github.com/bbrks/go-blurhash |
| **chai2010/webp** | WebP encode/decode | github.com/chai2010/webp |
| **kolesa-team/go-webp** | WebP processing | github.com/kolesa-team/go-webp |

---

## Media Processing

| Package | Purpose | URL |
|---------|---------|-----|
| **u2takey/ffmpeg-go** | FFmpeg bindings | github.com/u2takey/ffmpeg-go |
| **3d0c/gmf** | FFmpeg (CGO) | github.com/3d0c/gmf |
| **faiface/beep** | Audio processing | github.com/faiface/beep |

---

## HTTP Clients & APIs

| Package | Purpose | URL |
|---------|---------|-----|
| **go-resty/resty** | REST client | github.com/go-resty/resty/v2 |
| **hasura/go-graphql-client** | GraphQL client | github.com/hasura/go-graphql-client |
| **PuerkitoBio/goquery** | HTML scraping | github.com/PuerkitoBio/goquery |
| **chromedp/chromedp** | Headless Chrome | github.com/chromedp/chromedp |

---

## Caching

| Package | Purpose | URL |
|---------|---------|-----|
| **dgraph-io/ristretto** | In-memory cache | github.com/dgraph-io/ristretto |
| **redis/go-redis** | Redis client | github.com/redis/go-redis/v9 |
| **allegro/bigcache** | Fast in-memory cache | github.com/allegro/bigcache/v3 |

---

## Validation & Serialization

| Package | Purpose | URL |
|---------|---------|-----|
| **go-playground/validator** | Struct validation | github.com/go-playground/validator/v10 |
| **goccy/go-json** | Fast JSON | github.com/goccy/go-json |
| **mitchellh/mapstructure** | Map to struct | github.com/mitchellh/mapstructure |

---

## Logging & Observability

| Package | Purpose | URL |
|---------|---------|-----|
| **rs/zerolog** | Structured logging | github.com/rs/zerolog |
| **prometheus/client_golang** | Metrics | github.com/prometheus/client_golang |
| **open-telemetry/opentelemetry-go** | Tracing | go.opentelemetry.io/otel |

---

## Testing

| Package | Purpose | URL |
|---------|---------|-----|
| **stretchr/testify** | Assertions & mocks | github.com/stretchr/testify |
| **jarcoal/httpmock** | HTTP mocking | github.com/jarcoal/httpmock |
| **testcontainers/testcontainers-go** | Container testing | github.com/testcontainers/testcontainers-go |

---

## Utilities

| Package | Purpose | URL |
|---------|---------|-----|
| **spf13/viper** | Configuration | github.com/spf13/viper |
| **google/uuid** | UUID generation | github.com/google/uuid |
| **samber/lo** | Generics utilities | github.com/samber/lo |
| **hashicorp/go-multierror** | Error aggregation | github.com/hashicorp/go-multierror |
| **cenkalti/backoff** | Retry with backoff | github.com/cenkalti/backoff/v4 |

---

## CMS & Wiki (Optional)

| Package | Purpose | URL |
|---------|---------|-----|
| **gouniverse/cms** | Embeddable CMS | github.com/gouniverse/cms |
| **goliatone/go-cms** | Headless CMS toolkit | github.com/goliatone/go-cms |

> Note: Custom implementation recommended for full control.

---

## Package Selection Guidelines

### Prefer

1. **Actively maintained** - Recent commits, responsive maintainers
2. **Well-tested** - Good test coverage
3. **Documented** - Clear API docs
4. **Widely used** - Many imports on pkg.go.dev
5. **Pure Go** - No CGO when possible (easier deployment)

### Avoid

1. Abandoned projects (no commits in 2+ years)
2. Security vulnerabilities (check Dependabot)
3. Overly complex for simple needs
4. CGO dependencies unless necessary

---

## Related Documentation

- [Architecture Overview](OVERVIEW.md)
- [River Jobs](RIVER_JOBS.md)
- [Database Schema](DATABASE_SCHEMA.md)
