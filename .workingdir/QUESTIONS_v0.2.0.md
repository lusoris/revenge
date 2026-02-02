# v0.2.0 Questions & Decisions

**Version**: v0.2.0 - Core Backend Services
**Last Updated**: 2026-02-02

## Open Questions

### Authentication & Security

#### Q1: JWT Token Expiry Times
**Question**: What should be the default expiry times for access and refresh tokens?

**Options**:
- Access: 15min, Refresh: 7 days (high security)
- Access: 1h, Refresh: 30 days (balanced)
- Access: 24h, Refresh: 90 days (convenience)

**Decision**: TBD
**Notes**: Should be configurable via settings service

---

#### Q2: Password Requirements
**Question**: What password complexity requirements?

**Options**:
- Min 8 chars, any characters (low)
- Min 10 chars, must include uppercase, lowercase, number (medium)
- Min 12 chars, must include uppercase, lowercase, number, special (high)
- Configurable via settings

**Decision**: TBD
**Notes**: Consider using zxcvbn for password strength

---

#### Q3: Session Device Limit
**Question**: Should we limit the number of active sessions per user?

**Options**:
- No limit
- 10 sessions max
- 5 sessions max
- Configurable

**Decision**: TBD
**Notes**: Need to consider family sharing scenarios

---

### RBAC & Permissions

#### Q4: Default Roles
**Question**: What default roles should ship with the system?

**Proposed**:
- `admin` - Full system access
- `user` - Standard user
- `guest` - Read-only access
- `legacy:read` - QAR read-only (for v0.9.0)

**Decision**: TBD
**Notes**: Should these be deletable/editable?

---

#### Q5: Permission Granularity
**Question**: How granular should permissions be?

**Options**:
- Coarse: movie:read, movie:write, movie:delete
- Fine: movie:list, movie:get, movie:create, movie:update, movie:delete
- Very fine: movie:read:own, movie:read:library, movie:read:all

**Decision**: TBD
**Notes**: Need to balance flexibility vs complexity

---

### Database & Performance

#### Q6: Session Storage Strategy
**Question**: Where to store sessions?

**Options**:
- PostgreSQL only (simple, transactional)
- Dragonfly only (fast, but not persistent across restarts)
- Hybrid: Dragonfly L1, PostgreSQL L2 (best of both)

**Decision**: TBD
**Notes**: Hybrid recommended in design docs

---

#### Q7: Activity Log Retention
**Question**: How long should we keep activity logs?

**Options**:
- 30 days
- 90 days
- 1 year
- Configurable via settings

**Decision**: TBD
**Notes**: Need to consider GDPR requirements

---

#### Q8: Database Migration Strategy
**Question**: How to handle database migrations in production?

**Options**:
- golang-migrate with manual UP/DOWN
- Auto-migrate on startup (risky)
- Separate migration service/step in deployment

**Decision**: TBD
**Notes**: Production best practices?

---

### API Design

#### Q9: API Versioning Strategy
**Question**: How to version the API?

**Options**:
- URL versioning: `/api/v1/users`
- Header versioning: `Accept: application/vnd.revenge.v1+json`
- No versioning (breaking changes require new major version)

**Decision**: TBD
**Notes**: Using `/api/v1` in OpenAPI spec already

---

#### Q10: Pagination Strategy
**Question**: Default pagination approach?

**Options**:
- Cursor-based only (recommended in SOURCE_OF_TRUTH)
- Offset-based only
- Both supported (compatibility)

**Decision**: TBD (SOURCE_OF_TRUTH says both, cursor default)
**Notes**: Cursor is more performant for large datasets

---

### Infrastructure

#### Q11: River Queue Names
**Question**: How to organize job queues?

**Options**:
- Single queue for all jobs
- Priority-based: `critical`, `default`, `low`
- Function-based: `email`, `cleanup`, `metadata`, etc.

**Decision**: TBD
**Notes**: Priority-based seems most flexible

---

#### Q12: Cache TTL Strategy
**Question**: Default TTL values for different data types?

**Proposed**:
- Session tokens: 15 minutes
- User profiles: 5 minutes
- Settings: 1 hour
- Metadata: 24 hours

**Decision**: TBD
**Notes**: Should be configurable?

---

## Decided Questions

No decisions yet.

---

## Question Template

When adding questions:

```markdown
#### QX: Question Title
**Question**: What is the question?

**Options**:
- Option 1
- Option 2
- Option 3

**Decision**: TBD / [Decided]
**Notes**: Additional context
```

---

## Updates Log

| Date | Update |
|------|--------|
| 2026-02-02 | Created initial questions file with 12 open questions |
