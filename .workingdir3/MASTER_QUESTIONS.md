# Master Questions — All Decisions Made

> Consolidated from: STUB_QUESTIONS.md, PLANNING_ANALYSIS.md, CODEBASE_TODOS.md
> Generated: 2026-02-06 | **All 27 questions answered.**

---

## Part 1: Pattern Doc Stubs

### Q1 — SERVARR.md Scope
→ **B) Template for all** — prescriptive template for how ALL arr integrations should work, using Radarr as reference impl.

### Q2 — Servarr Webhook Auth
→ **Undecided** — user asked "hmac?" (needs clarification on what HMAC validation means in this context)

### Q3 — Servarr Priority Chain
→ **A) Arr authoritative** — Arr is always authoritative when configured and present. TMDb/TVDb are fallback + enrichment. Priority depends on content type and what sources can deliver. **Note**: Verify METADATA_SYSTEM.md reflects this correctly.

### Q4 — HTTP_CLIENT.md
→ **Keep and rewrite** — should be for external metadata API clients with configurable proxy support (enabled/disabled per-client). Not just VPN/Tor but general HTTP client factory for all external API calls.

### Q5 — METADATA.md
→ **C) Keep as overview** — higher-level pattern overview linking to METADATA_SYSTEM.md for details.

### Q6 — Request Coalescing (Sturdyc)
→ **No Sturdyc** — sync.Map is the bug (see CODE_ISSUES.md #1). Fix is to use the existing otter L1 + rueidis L2 cache infrastructure. No new library needed.

### Q7 — OBSERVABILITY.md Scope
→ **B) Include full vision** — document what exists AND the full planned observability story (RED metrics, distributed tracing, etc.)

### Q8 — Per-Handler RED Metrics
→ **A) Yes** — add to ogen middleware for all endpoints.

### Q9 — TESTING.md Scope
→ **A) Full guide** — comprehensive "how to write tests" guide covering unit tests, integration tests, and test helpers.

### Q10 — WEBHOOKS.md
→ **B) Keep standalone** — build a proper webhook service/pattern. Arrs use it via adapter. Design should consider River for async processing and caching infrastructure.

### Q11 — Outgoing Webhook Support
→ **Yes** — implied by Q10. Webhook service handles both incoming (from arrs) and outgoing (to external services). Notification service agents (Discord/Gotify/email) may integrate with or use the webhook service.

---

## Part 2: Planning Docs

### Q12 — Planning Doc Update Strategy
→ **A) Full rewrite** — update all 13 planning docs with correct statuses, fix links, add missing features.

### Q13 — 00_SOURCE_OF_TRUTH.md References
→ **B) Replace with DESIGN_INDEX.md** — point to DESIGN_INDEX.md as the new entry point.

### Q14 — Unplanned Features
→ **A + C) Both** — add retroactively to appropriate TODO files AND update ROADMAP.md.

---

## Part 3: Codebase TODOs

### Q15 — WebAuthn MFA Verification
→ **A) Implement now** — it's a gap in the auth chain.

### Q16 — OIDC Redirect Handler
→ **A) Fix with middleware** — custom middleware to inject Location header.

### Q17 — TV Show Search Indexing
→ **A) Implement now** — connect to existing Typesense infrastructure.

### Q18 — SendGrid Provider
→ **A) Implement SendGrid** — add actual SendGrid API integration.

### Q19 — Login/PasswordReset IP Extraction
→ **A) Fix now (full)** — extract IP, user agent, AND fingerprint from request headers.

### Q20 — Notification Settings in Profile Update
→ **A) Implement now** — wire up the JSONB notification preferences handling.

### Q21 — Search Reindex Endpoint
→ **A) Implement as River job** — wire to async job queue.

### Q22 — Library Stats
→ **A) Implement actual counts** — query from repository.

### Q23 — Cleanup Job
→ **A) Implement actual cleanup logic** — needs to define what gets cleaned (expired tokens, old activity logs, orphaned files, etc.)

### Q24 — Job MaxAttempts
→ **A) Make configurable** — add to config file.

### Q25 — Stale Comments Cleanup
→ **A) Clean up all** — remove all stale/misleading comments.

---

## Part 4: Architecture / Design

### Q26 — Profile Visibility Enforcement
→ **A) Implement now** — enforce in GetUserById handler.

### Q27 — Config Fields
→ **B) Keep** — LegacyPrivacyConfig (RequirePIN, AuditAllAccess) and TMDbConfig.ProxyURL are needed. If missing from design docs, add them.

---

## Action Summary

### Docs to Write (pattern stubs → full docs)
| Doc | Scope |
|-----|-------|
| SERVARR.md | Template for all arr integrations (Radarr as reference) |
| HTTP_CLIENT.md | External API client factory with configurable proxy |
| METADATA.md | High-level overview linking to METADATA_SYSTEM.md |
| OBSERVABILITY.md | Full vision (implemented + planned) |
| TESTING.md | Full "how to write tests" guide |
| WEBHOOKS.md | Standalone webhook service pattern (incoming + outgoing) |

### Planning Docs to Rewrite
All 13 files: fix statuses, fix links (→ DESIGN_INDEX.md), add missing features (MFA, Storage, etc.), update ROADMAP.md.

### Code to Implement
| Item | File | Priority |
|------|------|----------|
| WebAuthn MFA verification | mfa_integration.go | HIGH |
| OIDC redirect middleware | handler_oidc.go | HIGH |
| TV show search indexing | tvshow/jobs/jobs.go | HIGH |
| SendGrid email provider | email/service.go | HIGH |
| IP/UA/fingerprint extraction | handler.go | HIGH |
| Notification settings handler | handler.go | MEDIUM |
| Search reindex River job | handler_search.go | MEDIUM |
| Library stats from repo | library_service.go | MEDIUM |
| Cleanup job logic | cleanup_job.go | MEDIUM |
| Profile visibility enforcement | handler.go | MEDIUM |
| MaxAttempts configurable | jobs/module.go | LOW |
| Stale comments cleanup | multiple files | LOW |

### Architecture Fixes
| Item | File |
|------|------|
| Replace sync.Map with otter+rueidis | content/shared/metadata/client.go |
| Verify METADATA_SYSTEM.md reflects arr priority | docs/dev/design/architecture/METADATA_SYSTEM.md |
| Verify config fields in design docs | LegacyPrivacyConfig, TMDbConfig.ProxyURL |
