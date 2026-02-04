# Stubs and Unimplemented Functionality Report

**Generated**: 2026-02-04
**Status**: Analysis complete

---

## Summary

This report documents all stubs, placeholders, and unimplemented functionality discovered during a comprehensive codebase scan. Items are categorized by priority and type.

### Quick Stats
| Category | Count | Priority |
|----------|-------|----------|
| Critical Blockers | 3 | P0 |
| API Handler Issues | 9 | P1 |
| Service Stubs | 8 | P1-P2 |
| Repository Placeholders | 25+ | P2 |
| Content Module Stubs | 2 modules | P3 |
| Test Infrastructure | 2 | P2 |
| Future Features | 4 | P4 |

---

## P0: Critical Blockers

These items block core functionality and must be addressed immediately.

### 1. Email Service Not Implemented
**Impact**: User registration, password reset, email verification non-functional

| File | Line | Issue |
|------|------|-------|
| `internal/service/auth/service.go` | 96 | `// TODO: Send verification email (requires email service)` |
| `internal/service/auth/service.go` | 156 | `// TODO: Send verification email` |
| `internal/service/auth/service.go` | 409 | `// TODO: Send reset email` |

**Required**: Implement email service integration (SMTP, SendGrid, or similar)

### 2. User ID Placeholder in Authenticated Handlers
**Impact**: All authenticated API calls use hardcoded/random UUIDs instead of actual user

| File | Line | Issue |
|------|------|-------|
| `internal/api/handler.go` | 190-263 | `// TODO: Get user ID from auth context` - Multiple occurrences |
| `internal/api/handler.go` | 191 | `userID := uuid.New() // Placeholder` |
| `internal/api/handler.go` | 211 | `userID := uuid.New() // Placeholder` |
| `internal/api/handler.go` | 231 | `userID := uuid.New() // Placeholder` |
| `internal/api/handler.go` | 247 | `userID := uuid.New() // Placeholder` |
| `internal/api/handler.go` | 263 | `userID := uuid.New() // Placeholder` |
| `internal/api/handler.go` | 413 | `userID := uuid.MustParse("550e8400-...")` |
| `internal/api/handler.go` | 474 | `userID := uuid.MustParse("550e8400-...")` |
| `internal/api/handler.go` | 499 | `userID := uuid.MustParse("550e8400-...")` |
| `internal/api/handler.go` | 560 | `userID := uuid.MustParse("550e8400-...")` |

**Required**: Implement auth middleware to extract user ID from JWT/session and inject into context

### 3. Session Count Not Implemented
**Impact**: Cannot display active session count to users

| File | Line | Issue |
|------|------|-------|
| `internal/service/session/service.go` | 251 | `return 0, nil // TODO: Return actual count` |

**Required**: Implement `CountUserSessions` query in session repository

---

## P1: API Handler Issues

### Avatar Upload Not Implemented

| File | Line | Issue |
|------|------|-------|
| `internal/api/handler.go` | 562-566 | `// TODO: Parse multipart form and get file metadata` |
| `internal/api/handler.go` | 566 | `return &ogen.UploadAvatarBadRequest{}, fmt.Errorf("avatar upload not yet implemented")` |
| `internal/service/user/service.go` | 326-327 | `// TODO: Actually upload file to storage` - Returns placeholder path |

**Required**:
- Parse multipart form data
- Validate file type/size
- Upload to storage (local/S3)
- Update user avatar URL

### Request Metadata Extraction

| File | Line | Issue |
|------|------|-------|
| `internal/api/handler.go` | 621 | `// Authenticate user (TODO: extract IP, user agent, fingerprint from request)` |
| `internal/api/handler.go` | 735 | `// TODO: Extract IP address and user agent from request` |

**Required**: Middleware to extract and attach request metadata to context

### Search Async Job

| File | Line | Issue |
|------|------|-------|
| `internal/api/handler_search.go` | 194 | `// TODO: This should be an async job via River` |

**Required**: Implement River job for search indexing operations

---

## P1-P2: Service Stubs

### WebAuthn Session Cache

| File | Line | Issue |
|------|------|-------|
| `internal/service/mfa/webauthn.go` | 174 | `// TODO: Store session in cache (Redis/Dragonfly) with 5min TTL` |
| `internal/service/mfa/webauthn.go` | 340 | `// TODO: Store session in cache (Redis/Dragonfly) with 5min TTL` |
| `internal/service/auth/mfa_integration.go` | 112 | `return nil, errors.New("webauthn verification not yet implemented")` |

**Required**: Store WebAuthn challenge sessions in Dragonfly with TTL

### MFA Settings

| File | Line | Issue |
|------|------|-------|
| `internal/service/mfa/manager.go` | 77 | `RememberDeviceEnabled: false, // TODO: Get from user_mfa_settings` |

**Required**: Query user_mfa_settings table for remember device preference

### Notification Webhook Template

| File | Line | Issue |
|------|------|-------|
| `internal/service/notification/agents/webhook.go` | 190 | `// TODO: Support custom PayloadTemplate with Go templates` |

**Required**: Implement Go template parsing for custom webhook payloads

### OIDC Redirect Handling

| File | Line | Issue |
|------|------|-------|
| `internal/api/handler_oidc.go` | 75 | `// TODO: Implement custom redirect middleware or switch to JSON response` |
| `internal/api/handler_oidc.go` | 121 | `// TODO: If IsNewUser, create the user account via user service` |

**Required**:
- Handle OIDC callback redirects properly
- Auto-create user accounts on first OIDC login

---

## P2: Repository Placeholders

### Movie Repository (25+ unimplemented methods)

| File | Lines | Status |
|------|-------|--------|
| `internal/content/movie/repository_postgres.go` | 136-223 | Stub implementations returning `fmt.Errorf("not implemented")` |
| `internal/content/movie/repository_postgres.go` | 354-390 | Additional stub methods |

**Affected Methods**:
- `ListMoviesByIDs`
- `CountMovies`
- `SearchMovies`
- `ListMoviesByGenre`
- `ListRecentMovies`
- `ListPopularMovies`
- `GetMovieWithDetails`
- `GetMovieByExternalID`
- `UpdateMovieMetadata`
- `ListMovieCredits`
- `GetMovieCredit`
- `CreateMovieCredit`
- `DeleteMovieCredit`
- `DeleteMovieCredits`
- (and more)

**Note**: Some of these may be handled by sqlc-generated queries. Need verification.

### Library Matcher

| File | Line | Issue |
|------|------|-------|
| `internal/content/movie/library_matcher.go` | 120 | `return nil, fmt.Errorf("not implemented")` |
| `internal/content/movie/library_service.go` | 190 | `// For now, return placeholder` |

---

## P2: Movie Jobs

### File Match Job

| File | Line | Issue |
|------|------|-------|
| `internal/content/movie/moviejobs/file_match.go` | 58-64 | `// TODO: Implement once library.Service.MatchFile method is available` |

**Error returned**: `"movie file match not implemented: library.Service.MatchFile method not available"`

### Metadata Refresh Job

| File | Line | Issue |
|------|------|-------|
| `internal/content/movie/service.go` | 275 | `// TODO: Implement metadata refresh via River job` |
| `internal/content/movie/service.go` | 277 | `return fmt.Errorf("metadata refresh not implemented yet")` |
| `internal/content/movie/moviejobs/metadata_refresh.go` | 90-92 | `// TODO: Refresh credits if needed` |

---

## P2: Test Infrastructure

### Container Stubs

| File | Line | Issue |
|------|------|-------|
| `internal/testutil/containers.go` | 168-171 | Dragonfly container stub - `t.Skip("Dragonfly container not yet implemented")` |
| `internal/testutil/containers.go` | 182-185 | Typesense container stub - `t.Skip("Typesense container not yet implemented")` |

**Required**: Implement testcontainers for:
- Dragonfly (Redis-compatible cache)
- Typesense (search engine)

---

## P3: Content Module Placeholders

### TV Show Module

| Files | Status |
|-------|--------|
| `internal/content/tvshow/db/querier.go` | Placeholder query only |
| `internal/content/tvshow/db/placeholder.sql.go` | `SELECT 1 AS placeholder` |
| `internal/infra/database/queries/tvshow/placeholder.sql` | Placeholder SQL |

**Note**: Intentional - TV show support is v0.3.0+ feature

### QAR (Adult Content) Module

| Files | Status |
|-------|--------|
| `internal/content/qar/db/querier.go` | Placeholder query only |
| `internal/content/qar/db/placeholder.sql.go` | `SELECT 1 AS placeholder` |
| `internal/infra/database/queries/qar/placeholder.sql` | Placeholder SQL |

**Note**: Intentional - QAR support is v0.3.0+ feature

---

## P4: Future Features / Low Priority

### Job Configuration

| File | Line | Issue |
|------|------|-------|
| `internal/infra/jobs/module.go` | 40 | `MaxAttempts: 25, // TODO: Make configurable` |
| `internal/infra/jobs/cleanup_job.go` | 84 | `// For now, this is a stub that simulates work` |

### Windows MediaInfo Stub

| File | Line | Issue |
|------|------|-------|
| `internal/content/movie/mediainfo_windows.go` | 10-17 | `// MediaInfoProber is a stub implementation for Windows` |

**Note**: Intentional - MediaInfo uses FFmpeg/libav which has different Windows support

---

## Unimplemented API Endpoints (ogen generated)

The following endpoints exist in the OpenAPI spec but have no handler implementation (falling back to `ht.ErrNotImplemented`):

### RBAC Endpoints
- `POST /api/v1/rbac/policies` - AddPolicy
- `POST /api/v1/rbac/users/{userId}/roles` - AssignRole
- `POST /api/v1/rbac/roles` - CreateRole
- `DELETE /api/v1/rbac/roles/{roleId}` - DeleteRole

### OIDC Admin Endpoints
- `POST /api/v1/admin/oidc/providers` - AdminCreateOIDCProvider
- `DELETE /api/v1/admin/oidc/providers/{providerId}` - AdminDeleteOIDCProvider
- `POST /api/v1/admin/oidc/providers/{providerId}/disable` - AdminDisableOIDCProvider
- `POST /api/v1/admin/oidc/providers/{providerId}/enable` - AdminEnableOIDCProvider
- `GET /api/v1/admin/oidc/providers/{providerId}` - AdminGetOIDCProvider
- `GET /api/v1/admin/oidc/providers` - AdminListOIDCProviders
- `POST /api/v1/admin/oidc/providers/{providerId}/default` - AdminSetDefaultOIDCProvider
- `PATCH /api/v1/admin/oidc/providers/{providerId}` - AdminUpdateOIDCProvider

### Radarr Integration Endpoints
- `GET /api/v1/admin/integrations/radarr/quality-profiles` - AdminGetRadarrQualityProfiles
- `GET /api/v1/admin/integrations/radarr/root-folders` - AdminGetRadarrRootFolders
- `GET /api/v1/admin/integrations/radarr/status` - AdminGetRadarrStatus
- `POST /api/v1/admin/integrations/radarr/sync` - AdminTriggerRadarrSync

### Search Endpoints
- `GET /api/v1/search/movies/autocomplete` - AutocompleteMovies

### MFA Endpoints
- `POST /api/v1/auth/change-password` - ChangePassword
- `DELETE /api/v1/mfa` - DisableMFA
- `DELETE /api/v1/mfa/totp` - DisableTOTP
- `POST /api/v1/mfa` - EnableMFA
- `POST /api/v1/auth/forgot-password` - ForgotPassword
- `POST /api/v1/mfa/backup-codes` - GenerateBackupCodes
- `GET /api/v1/mfa/status` - GetMFAStatus

### Library Endpoints
- `POST /api/v1/libraries` - CreateLibrary
- `DELETE /api/v1/libraries/{libraryId}` - DeleteLibrary
- `GET /api/v1/libraries/{libraryId}` - GetLibrary

### API Key Endpoints
- `POST /api/v1/apikeys` - CreateAPIKey
- `GET /api/v1/apikeys/{keyId}` - GetAPIKey

### User Endpoints
- `GET /api/v1/users/me` - GetCurrentUser
- `GET /api/v1/users/me/session` - GetCurrentSession
- `DELETE /api/v1/users/me/settings/{key}` - DeleteUserSetting

### Content Endpoints (many more exist)
- Collection endpoints
- Movie detail endpoints
- Watch progress endpoints
- Activity endpoints

---

## Recommendations

### Immediate (v0.2.1 or hotfix)
1. **Fix user ID extraction** - This is blocking all authenticated endpoints
2. **Implement session count** - Minor fix, improves UX

### Short-term (v0.2.x)
1. **Email service** - Critical for user registration flow
2. **Avatar upload** - User-facing feature
3. **WebAuthn session cache** - Complete MFA functionality
4. **Auth middleware** for IP/user-agent extraction

### Medium-term (v0.3.0)
1. **Movie repository methods** - Many are likely needed for full movie browsing
2. **File match job** - Required for library scanning
3. **Metadata refresh job** - Required for keeping data up-to-date
4. **Test containers** - Required for cache/search integration tests
5. **Implement RBAC endpoints** - Required for permission management

### Long-term (v0.3.0+)
1. **TV Show module** - Planned feature
2. **QAR module** - Planned feature
3. **Radarr integration** - External integration
4. **OIDC admin management** - SSO configuration

---

## Files Scanned

- `internal/service/**/*.go`
- `internal/api/**/*.go`
- `internal/content/**/*.go`
- `internal/infra/**/*.go`
- `internal/testutil/**/*.go`

**Scan Method**: Grep for patterns:
- `TODO:`, `FIXME:`, `HACK:`, `XXX:`
- `not (yet )?implemented`, `placeholder`, `stub`
- `t.Skip(`, `panic("not implemented`, `return nil, nil //`, `return 0, nil //`

---

*Report generated by codebase analysis on 2026-02-04*
