# v0.2.0 Questions & Decisions

**Version**: v0.2.0 - Core Backend Services
**Last Updated**: 2026-02-04

---

## Decided Questions

### Q0: Server Settings Schema
**Decision**: Keep `shared.server_settings` (follows existing schema pattern)
**Date**: 2026-02-02

---

### Q1: JWT Token Expiry Times
**Decision**: Access: 24h, Refresh: 7 days (configurable via `auth.jwt_expiry`, `auth.refresh_expiry`)
**Implementation**: `internal/config/config.go:168-172`
**Date**: 2026-02-04

---

### Q2: Password Requirements
**Decision**: Minimum 8 characters (enforced in OpenAPI spec)
**Implementation**: `internal/api/ogen/oas_validators_gen.go`
**Date**: 2026-02-04

---

### Q3: Session Device Limit
**Decision**: Configurable via `maxPerUser`, warns when exceeded
**Implementation**: `internal/service/session/service.go:50-62`
**Date**: 2026-02-04

---

### Q4: Default Roles
**Decision**: 3 roles: `admin` (all), `user` (basic), `guest` (read-only)
**Implementation**: `migrations/000011_create_casbin_rule_table.up.sql`
**Date**: 2026-02-04

---

### Q5: Permission Granularity
**Decision**: Fine-grained permissions
**Format**: `resource:action` (e.g., `movie:list`, `movie:get`, `movie:create`, `movie:update`, `movie:delete`)
**Notes**: Allows better control per role
**Date**: 2026-02-04

---

### Q6: Email Service Provider
**Decision**: SMTP with configurable settings
**Implementation**: `internal/service/notification/agents/email.go`
**Date**: 2026-02-04

---

### Q7: Device Fingerprinting
**Decision**: DEFERRED - Detailed design needed for transcoding integration
**Notes**: Will be addressed in v0.6.0 (Transcoding) for stream session tracking
**Date**: 2026-02-04

---

### Q8: Database Migration Strategy
**Decision**: golang-migrate with manual UP/DOWN
**Implementation**: `internal/infra/database/migrations.go`
**Date**: 2026-02-04

---

### Q9: API Versioning Strategy
**Decision**: URL versioning `/api/v1/*`
**Date**: 2026-02-04

---

### Q10: Pagination Strategy
**Decision**: Both cursor-based and offset-based supported, cursor default
**Notes**: Cursor for performance, offset for compatibility
**Date**: 2026-02-04

---

### Q11: River Queue Priorities
**Decision**: 5 priority levels
- `critical` - Security events, auth failures
- `high` - User-initiated actions, notifications
- `default` - Metadata fetching, sync jobs
- `low` - Cleanup, maintenance
- `bulk` - Library scans, batch operations
**Date**: 2026-02-04

---

### Q12: Cache TTL Strategy
**Decision**: TTLs defined in `internal/infra/cache/keys.go`
- Session: 15 min
- User: 5 min
- RBAC: 10 min
- Settings: 1 hour
- Movie: 30 min
- Search: 5 min
**Date**: 2026-02-04

---

### Q13: Cache Integration Tests
**Decision**: Implement testcontainers-based Redis integration tests
**Notes**: Use existing `internal/testutil/containers.go` infrastructure
**Date**: 2026-02-04

---

### Q14: Session Storage Strategy
**Decision**: Hybrid - Dragonfly L1 + PostgreSQL L2
**Notes**: Dragonfly for fast session lookups, PostgreSQL for persistence
**Date**: 2026-02-04

---

### Q15: Activity Log Retention
**Decision**: 90 days default, configurable via Server Settings
**Notes**: GDPR-compliant with user-adjustable retention
**Date**: 2026-02-04

---

## Open Questions

None - All questions resolved for v0.2.0

---

## Updates Log

| Date | Update |
|------|--------|
| 2026-02-02 | Created initial questions file |
| 2026-02-04 | Verified implementations, moved 9 questions to Decided |
| 2026-02-04 | User decisions: Fine-grained permissions, testcontainers, Hybrid sessions, 5 queue priorities |
