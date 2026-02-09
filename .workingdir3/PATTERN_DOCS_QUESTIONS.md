# Questions for Pattern Documentation

Before writing the missing pattern docs, these decisions need to be clarified.
Answers drive the content of 5-8 new docs in `docs/dev/design/patterns/`.

---

## 1. Scope & Approach

### Q1.1: Write from code or design intent?
The existing pattern docs in `patterns/` (SERVARR.md, METADATA.md, etc.) are ~60-line auto-generated stubs. Should we:
- **A)** Rewrite them from actual code (document what IS, like Steps 5-10)
- **B)** Write them as prescriptive guides (document what SHOULD BE, for future dev)
- **C)** Both — describe current patterns + add "how to extend" sections

### Q1.2: Target audience?
- **A)** AI agents (Claude Code) building features — needs exact file paths, copy-paste templates
- **B)** Human developers joining the project — needs rationale, why-not alternatives
- **C)** Both — but prioritize AI-consumable format (code examples > prose)

### Q1.3: Rewrite existing pattern stubs or create new files?
The 6 existing files in `patterns/` are stubs. Should we:
- **A)** Delete the stubs, write fresh pattern docs for the actual gaps
- **B)** Keep stubs, add new docs alongside them
- **C)** Rewrite stubs + add new docs

---

## 2. Database Transaction Patterns

### Q2.1: Transaction boundary — where do transactions live?
Looking at the code, transactions appear in service layer (auth.Register uses pool.Begin). Is this the intended pattern?
- **A)** Service layer owns transactions (repo methods take `pgx.Tx` or `pgxpool.Pool`)
- **B)** Repository layer owns transactions (service calls `repo.WithTx(func)`)
- **C)** Mixed — depends on complexity (simple = repo, cross-entity = service)

### Q2.2: Should new services use the same pgxpool.Pool transaction pattern?
Auth was recently refactored to take `*pgxpool.Pool` directly. Is this the go-forward pattern for all services that need transactions?
- **A)** Yes — services that need transactions get pool injected
- **B)** No — introduce a transaction manager abstraction
- **C)** Case by case

### Q2.3: Isolation levels — do you care?
- **A)** Default (ReadCommitted) everywhere, don't document
- **B)** Document when to use Serializable (e.g., financial-like operations)
- **C)** Up to me to figure out from code

---

## 3. River Worker Patterns

### Q3.1: Worker registration pattern?
Currently workers are registered in content module fx.Modules (movie/module.go, tvshow/jobs/module.go). Should the pattern doc prescribe:
- **A)** Workers always live in the content module they belong to
- **B)** Workers can live anywhere but must register via fx
- **C)** Document what exists, don't prescribe

### Q3.2: Queue strategy — single queue or multi-queue?
JOBS.md currently (incorrectly) describes 5 priority queues. Reality is 1 default queue. What's the intent?
- **A)** Single queue is correct — keep it simple
- **B)** Multi-queue is planned — document the future design
- **C)** Single queue for now, note that multi-queue may come later

### Q3.3: Error handling in workers — retry strategy?
- **A)** Workers use River's built-in retry (MaxAttempts, exponential backoff)
- **B)** Custom retry logic per worker type
- **C)** Document what exists in code

---

## 4. Auth/JWT/Session Flow

### Q4.1: JWT vs Session — when to use which?
The codebase has both JWT tokens and database sessions. What's the model?
- **A)** JWT for API auth, sessions for state tracking (both always active)
- **B)** JWT is the primary auth, sessions are optional device tracking
- **C)** Sessions are primary, JWT is a convenience token derived from session

### Q4.2: How should new handlers check permissions?
- **A)** Handler calls `rbac.CheckPermission(ctx, resource, action)` directly
- **B)** Middleware handles RBAC based on route annotations/tags
- **C)** Mix — middleware for role checks, handler for resource-level checks
- **D)** Document what existing handlers do

### Q4.3: MFA integration point?
- **A)** MFA is checked during login only (JWT issued after MFA)
- **B)** MFA can be required for specific sensitive operations
- **C)** Document what exists

---

## 5. Error Handling Flow

### Q5.1: Error response format — is the current format final?
ogen generates error responses from the OpenAPI spec. Is the current error schema what you want?
- **A)** Yes — document the ogen-generated error format
- **B)** No — want custom error format (specify what)
- **C)** Document what exists, note improvements

### Q5.2: Validation errors — field-level detail?
- **A)** Return field-level validation errors (like `{"field": "email", "message": "invalid"}`)
- **B)** Return generic validation error messages
- **C)** Document what ogen does by default

### Q5.3: Internal errors — how much to expose?
- **A)** Never expose internal details (generic "internal server error")
- **B)** Expose error codes but not messages in production
- **C)** Document what exists

---

## 6. Frontend Patterns (if applicable now)

### Q6.1: Is frontend work happening soon?
- **A)** Yes — need SvelteKit patterns doc now
- **B)** Not yet — skip frontend doc, focus on backend patterns
- **C)** Just a stub for now, flesh out later

### Q6.2: API client generation?
- **A)** Generate TypeScript client from OpenAPI spec (which tool?)
- **B)** Hand-write fetch wrappers
- **C)** TanStack Query with manual fetch
- **D)** Not decided yet

### Q6.3: Auth in frontend?
- **A)** JWT in httpOnly cookie (SSR-friendly)
- **B)** JWT in localStorage (SPA-style)
- **C)** Session cookie (server-managed)
- **D)** Not decided yet

---

## 7. Cache Patterns

### Q7.1: Cache key naming convention?
- **A)** `{service}:{entity}:{id}` (e.g., `movie:detail:123`)
- **B)** Already exists in code — just document it
- **C)** No convention yet — need to establish one

### Q7.2: When should new services add caching?
- **A)** All read-heavy services get a CachedService wrapper
- **B)** Only when performance requires it
- **C)** Document the existing CachedService pattern and let devs decide

---

## 8. New Service Checklist

### Q8.1: Should there be a "create new service" template/checklist?
- **A)** Yes — step-by-step: entity → repo interface → pg impl → service → cached service → handler → fx module → tests
- **B)** No — the existing service docs are enough examples
- **C)** Yes but keep it short (just a checklist, not a tutorial)

---

## 9. Doc Fixes

### Q9.1: Fix the 3 known doc issues now or later?
- JOBS.md: wrong worker count (17→11) and queue system (5→1)
- ARCHITECTURE.md: missing 3 fx modules
- METADATA_SYSTEM.md: method count (31→27)

- **A)** Fix now as part of this batch
- **B)** Fix later, separate task
- **C)** Fix JOBS.md now (it's the worst), defer the minor ones

---

## Priority Order

If we can only write 5 docs, which 5 matter most? My suggestion:

1. **New Service Checklist** (patterns/NEW_SERVICE.md) — covers transactions, fx, testing
2. **Auth & Permissions Flow** (patterns/AUTH_FLOW.md) — JWT, session, RBAC, MFA
3. **River Worker Guide** (patterns/RIVER_WORKERS.md) — add worker step-by-step
4. **Error Handling** (patterns/ERROR_HANDLING.md) — repo → service → handler → response
5. **Cache Strategy** (patterns/CACHE_STRATEGY.md) — L1/L2, invalidation, CachedService

Frontend patterns deferred unless Q6.1 = "Yes".
