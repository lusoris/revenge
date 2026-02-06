# Patterns

← Back to [Design Docs](../)

---

## Written from Code

| Document | Description | Status |
|----------|-------------|--------|
| [Authentication & Authorization Flow](AUTH_FLOW.md) | JWT, sessions, MFA, RBAC, API keys — full auth chain | ✅ Complete |
| [New Service Checklist](NEW_SERVICE.md) | Step-by-step guide: repo → service → cached service → handler → fx → tests | ✅ Complete |
| [River Worker Guide](RIVER_WORKERS.md) | How to add background job workers, 5-tier queue system, all 17 workers | ✅ Complete |
| [Error Handling](ERROR_HANDLING.md) | Error flow from database → repo → service → handler → HTTP response | ✅ Complete |
| [Cache Strategy](CACHE_STRATEGY.md) | L1 (otter) + L2 (rueidis) caching, key conventions, CachedService pattern | ✅ Complete |
| [Database Transactions](DATABASE_TRANSACTIONS.md) | pool.Begin / txQueries pattern for atomic multi-row operations | ✅ Complete |
| [fx Module Patterns](FX_MODULE_PATTERNS.md) | Three wiring patterns for dependency injection with uber/fx | ✅ Complete |
| [Arr Integration Pattern](SERVARR.md) | 8-layer template for Radarr, Sonarr, and future arr integrations | ✅ Complete |
| [HTTP Client Pattern](HTTP_CLIENT.md) | req/v3 clients: rate limiting, caching, proxy, retry, auth patterns | ✅ Complete |
| [Metadata Enrichment Pattern](METADATA.md) | Two-tier model, provider priority, adapters, caching layers | ✅ Complete |
| [Observability Pattern](OBSERVABILITY.md) | Prometheus metrics, logging, pprof, middleware, planned tracing | ✅ Complete |
| [Testing Patterns](TESTING.md) | Fast DB tests, mocking, table-driven, integration tests, CI | ✅ Complete |
| [Webhook Patterns](WEBHOOKS.md) | Incoming webhook handling, type conversion, async processing | ✅ Complete |

---

## Status Legend

✅ Complete — written from code, accurate as of 2026-02-06
