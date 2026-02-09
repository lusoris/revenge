# Pattern Stub Questions

These 6 pattern docs in `docs/dev/design/patterns/` are auto-generated stubs with no real content. Each needs a decision: rewrite from code, move to planned, or delete.

---

## 1. SERVARR.md — Arr Integration Pattern

**Current content**: Claims webhook-based integration with Radarr, Sonarr, Lidarr, Whisparr. Priority chain: Arr metadata > Internal > External APIs.

**What actually exists in code**:
- Radarr: Full client (`internal/integration/radarr/`) — client.go (543 LOC), mapper.go, jobs.go (RadarrSyncWorker, RadarrWebhookWorker)
- Sonarr: Full client (`internal/integration/sonarr/`) — client.go (724 LOC), mapper.go, jobs.go (SonarrSyncWorker, SonarrWebhookWorker)
- Lidarr: **No code exists** (planned for music module)
- Whisparr: **No code exists** (planned for QAR/adult module)
- Chaptarr/Readarr: **No code exists** (planned for book/comic module)

**Questions**:

### Q1.1: What is the scope of this doc?
You said Servarr = "all arr master metadata sources" (Radarr, Whisparr, Lidarr, Sonarr, Chaptarr/Readarr).
- **A)** Rewrite to document the **existing** Radarr+Sonarr patterns, with a section noting planned arr integrations (Lidarr, Whisparr, Chaptarr)
- **B)** Write it as a **prescriptive template** for how ALL arr integrations should work (client, mapper, sync worker, webhook worker), using Radarr as the reference implementation
- **C)** Both: document Radarr/Sonarr patterns + template for future arr integrations

### Q1.2: Webhook handling
The Radarr/Sonarr handlers exist in `internal/api/handler_radarr.go` and `handler_sonarr.go`. They receive webhooks and queue jobs. Is HMAC validation planned or is the current approach (API key auth) the final design?
- **A)** Current approach (API key + queueing) is final
- **B)** HMAC validation should be added
- **C)** Not decided yet

### Q1.3: Arr priority chain
The stub mentions "Priority Chain: Arr metadata > Internal > External APIs". Is this still the intended priority for metadata resolution?
- **A)** Yes — Arr metadata is authoritative, TMDb/TVDb fill gaps
- **B)** No — TMDb/TVDb are primary, Arr provides sync/library management only
- **C)** Depends on the content type

---

## 2. HTTP_CLIENT.md — HTTP Client with Proxy/VPN Pattern

**Current content**: Claims reusable HTTP client factory with proxy/VPN/Tor/SOCKS5 support. Status: Code 🔴

**What actually exists in code**: Standard `http.Client` usage in integration clients (TMDb, TVDb, Radarr, Sonarr). No proxy/VPN/Tor/SOCKS5 code exists.

**Questions**:

### Q2.1: Is proxy/VPN support still planned?
- **A)** Yes — needed for accessing metadata providers in restricted regions
- **B)** Not a priority — delete stub or move to planned/
- **C)** Defer — will revisit when needed

### Q2.2: If yes, should it be a shared HTTP client factory?
- **A)** Yes — all external API clients should use a configurable factory
- **B)** Per-client configuration is fine
- **C)** Not decided yet

---

## 3. METADATA.md — Metadata Enrichment Pattern

**Current content**: Claims multi-tier enrichment with request coalescing (Sturdyc). Priority chain: Cache → Arr → Internal → External → Background.

**What actually exists in code**:
- Full metadata service (`internal/service/metadata/`) — 18 files, 8,087 LOC
- TMDb provider, TVDb provider, adapter pattern, ClearCache
- sync.Map caching in providers (L0), otter+rueidis (L1+L2) for services
- Background refresh via River workers
- **No Sturdyc (request coalescing)** — the doc references it but it's not in the code
- **No "priority chain" as described** — the actual chain is: providers with priority scores (TMDb=100, TVDb=80), higher priority wins

**Questions**:

### Q3.1: Rewrite from code or merge with METADATA_SYSTEM.md?
The architecture doc `METADATA_SYSTEM.md` already documents the actual metadata system thoroughly. This stub overlaps.
- **A)** Delete this stub — METADATA_SYSTEM.md is sufficient
- **B)** Rewrite as a "how to add a new metadata provider" guide (practical, not architectural)
- **C)** Keep as a higher-level pattern overview, link to METADATA_SYSTEM.md for details

### Q3.2: Is request coalescing (Sturdyc) still planned?
- **A)** Yes — needed for deduplicating concurrent metadata requests
- **B)** Not needed — sync.Map caching handles the hot path
- **C)** Not decided yet

---

## 4. OBSERVABILITY.md — Observability Pattern

**Current content**: Claims Prometheus metrics, OpenTelemetry tracing, structured logging. Status: Code 🔴

**What actually exists in code**:
- `internal/infra/observability/` — 5 files (548 LOC):
  - pprof endpoints
  - Prometheus metrics (cache hit/miss counters, operation duration histograms)
  - OpenTelemetry setup (OTLP exporter)
- `internal/infra/logging/` — 2 files:
  - slog (dev mode) + zap (prod mode)
  - Structured logging throughout codebase
- Cache operations record metrics via `observability.RecordCacheHit()` etc.

**Questions**:

### Q4.1: Rewrite from code?
Code exists for observability. Should the stub be rewritten to document:
- **A)** What's implemented (pprof, Prometheus cache metrics, OTLP, slog/zap) — a practical "how to add metrics" guide
- **B)** The full observability vision (including things not yet done like per-endpoint RED metrics, distributed tracing spans)
- **C)** Just what's implemented, with a "planned" section for future work

### Q4.2: Are per-handler metrics (RED: Rate, Errors, Duration) planned?
- **A)** Yes — should be added to the ogen middleware
- **B)** No — Prometheus scraping of pprof is enough for now
- **C)** Not decided yet

---

## 5. TESTING.md — Testing Patterns

**Current content**: Claims table-driven tests, testify, mockery, embedded Postgres, testcontainers.

**What actually exists in code**:
- 144 test files, 62,906 LOC of tests
- `testutil` package with `NewFastTestDB(t)` (shared PostgreSQL via testcontainers)
- Table-driven tests throughout
- testify (assert/require/mock) everywhere
- mockery-generated mocks for repository interfaces
- Integration tests in `tests/integration/` (13 files, testcontainers)
- **No embedded Postgres** — uses testcontainers PostgreSQL instead

**Questions**:

### Q5.1: Rewrite from code?
This has the most code backing it. Should be rewritten to document actual test patterns:
- **A)** Full rewrite as a "how to write tests" guide (unit tests with mocks, integration tests with testcontainers, test helpers)
- **B)** Brief overview — the test files themselves are the best documentation
- **C)** Focus on the testutil package and testcontainers setup (the non-obvious parts)

### Q5.2: Embedded Postgres?
The stub mentions it but it doesn't exist. The codebase uses testcontainers PostgreSQL for all DB tests.
- **A)** Remove reference — testcontainers is the approach
- **B)** Embedded Postgres is planned as a faster alternative
- **C)** Not decided

---

## 6. WEBHOOKS.md — Webhook Patterns

**Current content**: Claims HMAC validation, async processing, deduplication, retry logic, IP whitelisting.

**What actually exists in code**:
- Radarr webhook handler (`internal/api/handler_radarr.go`) — receives webhook, queues job
- Sonarr webhook handler (`internal/api/handler_sonarr.go`) — same pattern
- RadarrWebhookWorker, SonarrWebhookWorker — process queued webhooks
- **No HMAC validation** — webhooks are received on authenticated endpoints
- **No deduplication** — processed as-is
- **No IP whitelisting** — standard auth applies
- **No outgoing webhooks** — only incoming from arr stack

**Questions**:

### Q6.1: Is this about incoming or outgoing webhooks?
- **A)** Incoming only (from Radarr/Sonarr) — merge into SERVARR.md
- **B)** Outgoing only (sending events to external services) — this is planned but not built
- **C)** Both — document incoming (from arr) and plan for outgoing

### Q6.2: Is HMAC validation/deduplication/IP whitelisting planned?
- **A)** Yes — needed for security
- **B)** No — current auth-based approach is sufficient
- **C)** Only for outgoing webhooks

---

## Summary: Recommended Actions

| Stub | Code exists? | Recommendation |
|------|-------------|---------------|
| SERVARR.md | Yes (Radarr+Sonarr) | Rewrite from code + template for future arr integrations |
| HTTP_CLIENT.md | No | Move to planned/ or delete — no code to document |
| METADATA.md | Yes but METADATA_SYSTEM.md covers it | Delete or convert to "add a provider" guide |
| OBSERVABILITY.md | Partial | Rewrite from code (metrics, logging, pprof) |
| TESTING.md | Yes (extensive) | Rewrite from code (test patterns, testutil, testcontainers) |
| WEBHOOKS.md | Partial (incoming only) | Merge incoming into SERVARR.md, move outgoing to planned/ |

**Waiting for your answers before proceeding with any rewrites.**
