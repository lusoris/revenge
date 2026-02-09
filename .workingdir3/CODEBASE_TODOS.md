# Codebase TODO / Stub / Incomplete Scan

> Full scan of all Go source, config, SQL, tests, and CI files. Generated 2026-02-06.

---

## HIGH — Broken or Missing Functionality

### 1. WebAuthn MFA Verification Not Implemented
- **File**: `internal/service/auth/mfa_integration.go:112`
- **Issue**: WebAuthn case returns `errors.New("webauthn verification not yet implemented")`
- **Impact**: FIDO2/WebAuthn MFA method completely broken — users who try it get an error
- **Schema supports it**: Yes (migration 000018 creates `webauthn_credentials` table)

### 2. OIDC Redirect Handler Broken
- **File**: `internal/api/handler_oidc.go:77`
- **Issue**: `OidcAuthorize` computes auth URL but returns empty `ogen.OidcAuthorizeFound{}` — ogen can't set Location header on 302
- **Comment**: `// ogen generates empty struct for 302, no Location header support`
- **Impact**: OIDC OAuth flow completely non-functional

### 3. TV Show Search Indexing is No-Op
- **File**: `internal/content/tvshow/jobs/jobs.go:776-865`
- **Issue**: `SearchIndexWorker` logs "would index series (search service not implemented)" but does nothing
- **Impact**: TV shows not searchable via Typesense (movie search works fine)

### 4. SendGrid Email Provider Not Implemented
- **File**: `internal/service/email/service.go:205-207`
- **Issue**: Logs "Email sent via SendGrid" then returns `fmt.Errorf("SendGrid provider not yet implemented, use SMTP")`
- **Impact**: Misleading log + broken if config selects sendgrid provider. SMTP works fine.

---

## MEDIUM — Incomplete but Functional

### 5. Login Handler Missing IP/UserAgent/Fingerprint Extraction
- **File**: `internal/api/handler.go:713`
- **Comment**: `// TODO: extract IP, user agent, fingerprint from request`
- **Current**: Passes `nil, nil, deviceName, nil` to `authService.Login()`
- **Impact**: No IP logging, no device fingerprint for security tracking. Auth works, security metadata missing.

### 6. Password Reset Missing IP/UserAgent Extraction
- **File**: `internal/api/handler.go:828`
- **Comment**: `// TODO: Extract IP address and user agent from request`
- **Current**: Passes `nil, nil` to `authService.RequestPasswordReset()`
- **Impact**: Password reset events not properly audited

### 7. Profile Update Ignores Notification Settings
- **File**: `internal/api/handler.go:528`
- **Comment**: `// TODO: Handle notification settings (JSONB fields)`
- **Schema**: `user_preferences` table has `email_notifications`, `push_notifications`, `digest_notifications` JSONB columns (migration 000006)
- **Impact**: Notification preferences silently ignored on profile update

### 8. Search Reindex is Synchronous Stub
- **File**: `internal/api/handler_search.go:194`
- **Comment**: `// TODO: This should be an async job via River`
- **Current**: Generates a fake job ID, logs request, returns accepted — does NOT actually reindex
- **Impact**: Reindex endpoint is a no-op

### 9. Library Stats Returns Hardcoded Zeros
- **File**: `internal/content/movie/library_service.go:241`
- **Comment**: `// For now, return placeholder`
- **Current**: Returns `{"total_movies": 0, "total_files": 0}`
- **Impact**: Stats API always shows empty library

### 10. Cleanup Job is Stub
- **File**: `internal/infra/jobs/cleanup_job.go:111`
- **Comment**: `// Actual cleanup logic would go here (database operations)`
- **Current**: Logs "performing cleanup" but no actual deletion
- **Impact**: Old data never cleaned up (job runs but does nothing)

### 11. Localization User Settings Not Implemented
- **File**: `internal/api/localization.go:15`
- **Comment**: `// TODO: Implement when user settings are available`
- **Current**: Commented-out code for reading user preferred language from context
- **Impact**: User language preference in `user_preferences.display_language` ignored, falls back to Accept-Language header

### 12. Movie Repository "Placeholder" Comment Misleading
- **File**: `internal/content/movie/repository_postgres.go:343`
- **Comment**: `// Placeholder implementations for remaining methods / TODO: Implement all repository methods`
- **Actual**: The methods below ARE implemented — the comment is stale
- **Impact**: No functional issue, just misleading comment

---

## LOW — Polish / Config

### 13. Job MaxAttempts Hardcoded
- **File**: `internal/infra/jobs/module.go:39`
- **Comment**: `// TODO: Make configurable`
- **Current**: `MaxAttempts: 25`
- **Impact**: Can't tune without code change

### 14. Config: Stale v0.1.0 Comment
- **File**: `internal/config/config.go:191`
- **Comment**: `// Optional for v0.1.0 (auth not implemented yet)` on JWTSecret
- **Actual**: Auth IS fully implemented. Comment is stale.

### 15. Windows Media Probing Stub (By Design)
- **File**: `internal/content/movie/mediainfo_windows.go:23`
- **Returns**: `fmt.Errorf("media probing not supported on Windows: ...")`
- **Note**: This is intentional — CGO/mediainfo not available on Windows

### 16. API Key Usage Count is Workaround
- **File**: `internal/infra/database/db/querier.go:174-176`
- **Comment**: `// This is a placeholder - actual usage tracking would be in a separate table`
- **Current**: Returns `last_used_at` instead of actual count

---

## SQL Placeholder Queries (sqlc Requirement)

These exist solely to satisfy sqlc's non-empty directory requirement. Not bugs.

| File | Content |
|------|---------|
| `internal/infra/database/queries/movie/placeholder.sql` | `SELECT 1 AS placeholder` |
| `internal/infra/database/queries/tvshow/placeholder.sql` | `SELECT 1 AS placeholder` |
| `internal/infra/database/queries/qar/placeholder.sql` | `SELECT 1 AS placeholder` |

---

## Error Swallowing (Intentional but Worth Noting)

| File | Line | What's Swallowed | Reason |
|------|------|-----------------|--------|
| `internal/service/user/service.go` | 108 | User preferences creation | Don't fail user creation for prefs |
| `internal/content/movie/library_matcher.go` | 184 | Metadata enrichment error | Continue with partial data |
| `internal/config/loader.go` | 39 | Config file load error | Config file is optional |

---

## Summary

| Severity | Count | Key Items |
|----------|-------|-----------|
| **HIGH** | 4 | WebAuthn MFA, OIDC redirect, TV show search, SendGrid |
| **MEDIUM** | 8 | IP extraction, notification prefs, reindex, stats, cleanup |
| **LOW** | 4 | Config hardcoding, stale comments, API key usage |
| **Placeholder** | 3 | sqlc directory requirements |
| **Error Swallowing** | 3 | Intentional graceful degradation |

**Total markers found**: 22 (excluding 3 sqlc placeholders and generated code)

---

## Architecture Issues (from CODE_ISSUES.md)

### 17. metadata.BaseClient uses sync.Map Instead of Proper Cache
- **File**: `internal/content/shared/metadata/client.go`
- **Severity**: Medium (unbounded memory growth)
- **Issue**: Raw `sync.Map` with no size limit, no eviction, expired entries linger. The proper cache infra (otter L1 + Dragonfly L2) exists but isn't used here.
- **Fix**: Replace with `cache.Cache` or bounded otter directly.

### 18. CI/CD Workflows Need Review
- **Location**: `.github/workflows/`
- **Severity**: Medium
- **Issue**: Some workflows may be broken or misconfigured. Needs deeper audit.

---

## Deeper Audit Checklist (Not Yet Scanned)

These architecture checks were flagged during the doc rewrite but haven't been fully audited:

- [ ] Other `sync.Map` caches that should use otter/Dragonfly
- [ ] Inconsistent error handling patterns across modules
- [ ] Services bypassing repository layer (direct DB access)
- [ ] Missing cache invalidation (writes that don't invalidate related reads)
- [ ] Hardcoded values that should come from config (timeouts, limits, URLs)
- [ ] Duplicate logic between movie and tvshow modules that should be in shared/
- [ ] Workers missing progress reporting or proper error handling
- [ ] Rate limiters with wrong values or missing entirely
- [ ] Context propagation gaps (functions that don't accept/pass context)
- [ ] Type conversions that silently lose data (int64 → int32, etc.)

> These are tracked here for completeness. A full architecture audit is a separate task.
