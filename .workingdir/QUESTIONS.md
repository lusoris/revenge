# v0.1.0 Questions & Decisions Log

**Updated**: 2026-02-02 11:30
**Milestone**: v0.1.0 - Skeleton

---

## Open Questions

### Q1: Should we generate ogen code in CI or commit it?

**Status**: ⏳ Open
**Priority**: High
**Context**: We need to decide whether to commit ogen-generated code or regenerate it in CI/CD.

**Options**:

**Option A: Commit Generated Code**
- ✅ Pros:
  - Faster CI builds (no generation step)
  - Reviewers can see API changes in PRs
  - No risk of generation failures in CI
- ❌ Cons:
  - Large diffs when OpenAPI spec changes
  - Potential merge conflicts in generated files
  - Need to remember to regenerate before commits

**Option B: Generate in CI**
- ✅ Pros:
  - Smaller Git history
  - Always in sync with OpenAPI spec
  - No manual regeneration needed
- ❌ Cons:
  - Slower CI builds
  - Cannot review generated code changes
  - Generation failures break CI

**Recommendation**: ❓ TBD

**Decision Maker**: Tech Lead
**Deadline**: Before Phase 1 completion

---

### Q2: What should the default config values be?

**Status**: ⏳ Open
**Priority**: Medium
**Context**: Need to define sensible defaults for all configuration options.

**Current Status**:
- `config/config.yaml` exists but may need review
- `config/config.test.yaml` exists for testing

**Questions**:
1. Default server port: 8080 or 3000?
2. Default log level: info or debug for development?
3. Database connection pool size defaults?
4. Cache TTL defaults?
5. Health check timeout defaults?

**Action**: Review current config.yaml and document decisions

**Deadline**: Phase 4

---

### Q3: Should we use embedded-postgres or testcontainers for database tests?

**Status**: ⏳ Open
**Priority**: Medium
**Context**: Both options are implemented in testutil package.

**Options**:

**Option A: embedded-postgres**
- ✅ Pros:
  - Faster startup
  - No Docker required
  - Simpler CI setup
- ❌ Cons:
  - Less realistic (different PostgreSQL version)
  - Limited to unit tests
  - No other services (Dragonfly, Typesense)

**Option B: testcontainers**
- ✅ Pros:
  - Real PostgreSQL, Dragonfly, Typesense
  - Full integration testing
  - Production-like environment
- ❌ Cons:
  - Requires Docker
  - Slower startup
  - More complex CI setup

**Option C: Both**
- ✅ Pros:
  - Fast unit tests (embedded-postgres)
  - Real integration tests (testcontainers)
  - Best of both worlds
- ❌ Cons:
  - Maintain two test setups
  - More complexity

**Recommendation**: Option C (Both)
- Use embedded-postgres for fast database unit tests
- Use testcontainers for full integration tests
- Separate with build tags or test names

**Decision**: Pending approval

**Deadline**: Phase 3

---

### Q4: What HTTP framework should we use with ogen?

**Status**: ⏳ Open
**Priority**: High
**Context**: ogen generates server interfaces, but we need an HTTP framework underneath.

**Options**:

**Option A: net/http (stdlib)**
- ✅ Pros:
  - No dependencies
  - Maximum compatibility
  - Well-known and stable
- ❌ Cons:
  - Less convenient API
  - Manual middleware management
  - No context handling helpers

**Option B: chi**
- ✅ Pros:
  - Lightweight
  - Compatible with stdlib
  - Good middleware ecosystem
  - Context handling
- ❌ Cons:
  - Another dependency
  - Learning curve

**Option C: echo**
- ✅ Pros:
  - Feature-rich
  - Good performance
  - Large community
- ❌ Cons:
  - Not stdlib-compatible
  - More opinionated

**Recommendation**: Option A (net/http)
- ogen generates stdlib-compatible code
- Keep dependencies minimal for v0.1.0
- Can add chi/echo later if needed

**Decision**: Tentatively Option A

**Deadline**: Phase 1

---

### Q5: How should we structure SQL queries for sqlc?

**Status**: ⏳ Open
**Priority**: Medium
**Context**: Need to decide on directory structure and naming conventions for SQL files.

**Options**:

**Option A: By Domain**
```
internal/
  user/
    queries/
      user.sql
  session/
    queries/
      session.sql
```

**Option B: Centralized**
```
internal/
  db/
    queries/
      user.sql
      session.sql
```

**Option C: By Table**
```
sql/
  users/
    queries.sql
  sessions/
    queries.sql
```

**Recommendation**: Option A (By Domain)
- Follows clean architecture
- Co-locates queries with domain logic
- Easier to navigate for larger projects

**Decision**: Pending

**Deadline**: Before implementing sqlc queries

---

## Resolved Questions

### ✅ Q-RESOLVED-001: Should we use Go 1.25.6 or downgrade to 1.25.5?

**Status**: ✅ Resolved
**Decision**: Keep Go 1.25.6
**Rationale**: 
- Latest version has bug fixes
- Source of truth document specifies 1.25.6
- Version mismatch is a toolchain issue, not a project issue
- Will fix toolchain separately

**Decided By**: Development Team
**Date**: 2026-02-02

---

## Design Decisions

### DD-001: Use fx for Dependency Injection

**Status**: ✅ Approved
**Decision**: Use go.uber.org/fx for dependency injection
**Rationale**:
- Industry standard in Go
- Good lifecycle management
- Clear module boundaries
- Excellent for testing (can swap dependencies)

**Alternatives Considered**:
- Wire (Google) - Too compile-time, less flexible
- Manual DI - Too error-prone, hard to maintain
- dig - fx is built on dig, better ergonomics

**Date**: Pre-v0.1.0
**References**: [01_ARCHITECTURE.md](../docs/dev/design/architecture/01_ARCHITECTURE.md)

---

### DD-002: Use ogen for OpenAPI Code Generation

**Status**: ✅ Approved
**Decision**: Use github.com/ogen-go/ogen for OpenAPI 3.1 code generation
**Rationale**:
- OpenAPI 3.1 support (latest spec)
- Pure Go implementation
- Generates both server and client
- Type-safe API implementation
- Active development

**Alternatives Considered**:
- oapi-codegen - No OpenAPI 3.1 support
- go-swagger - OpenAPI 2.0 only
- Manual implementation - Too error-prone

**Date**: Pre-v0.1.0
**References**: [API.md](../docs/dev/design/technical/API.md)

---

### DD-003: Use sqlc for Type-Safe SQL

**Status**: ✅ Approved
**Decision**: Use github.com/sqlc-dev/sqlc for SQL code generation
**Rationale**:
- Type-safe database queries
- Write SQL, get Go
- Compile-time safety
- No ORM overhead

**Alternatives Considered**:
- GORM - Too much magic, performance overhead
- sqlx - Less type-safe
- Raw pgx - Too verbose, error-prone

**Date**: Pre-v0.1.0
**References**: [POSTGRESQL.md](../docs/dev/design/integrations/infrastructure/POSTGRESQL.md)

---

### DD-004: Use Structured Logging (slog)

**Status**: ✅ Approved
**Decision**: Use stdlib log/slog with tint (dev) and zap (prod) handlers
**Rationale**:
- Stdlib slog is standard
- tint provides colored output for development
- zap provides high-performance JSON logging for production
- Structured fields enable better observability

**Date**: Pre-v0.1.0
**References**: [OBSERVABILITY.md](../docs/dev/design/technical/OBSERVABILITY.md)

---

## Action Items

1. **IMMEDIATE**: Decide Q4 (HTTP framework) before starting Phase 1
2. **HIGH**: Resolve Q1 (ogen code commit) before first PR
3. **MEDIUM**: Review Q2 (config defaults) in Phase 4
4. **MEDIUM**: Decide Q3 (test strategy) before Phase 3
5. **MEDIUM**: Decide Q5 (SQL structure) before implementing queries

---

## Decision Log Template

Use this template for new decisions:

```markdown
### Q#: Question Title

**Status**: ⏳ Open / ✅ Resolved
**Priority**: High / Medium / Low
**Context**: Brief description of the question

**Options**:
- Option A: Description
  - ✅ Pros:
  - ❌ Cons:
- Option B: Description
  - ✅ Pros:
  - ❌ Cons:

**Recommendation**: Preferred option
**Decision**: TBD / Approved
**Decision Maker**: Who decides
**Deadline**: When decision is needed
**References**: Links to docs
```

---

## Notes

- All major architectural decisions are documented in design docs
- This file tracks implementation-level questions
- Resolved questions moved to "Resolved" section for reference
- Link to design docs for authoritative architectural decisions
