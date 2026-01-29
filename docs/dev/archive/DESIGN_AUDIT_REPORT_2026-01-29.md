# Design vs Code Audit Report

> Generated: 2026-01-29
> Scope: All docs under `docs/dev/design/` (source of truth)

---

## Executive Summary

**Total Issues Found**: 67
- **CRITICAL**: 12
- **HIGH**: 28
- **MEDIUM**: 18
- **LOW**: 9

**Overall Code-Design Alignment**: ~55%

---

## CRITICAL Issues (Must Fix Immediately)

### 1. Session Service - Wrong Signature
**File**: `internal/service/session/service.go`
**Design**: `UpdateActivity(ctx, sessionID, ipAddress *netip.Addr)`
**Code**: `UpdateActivity(ctx, sessionID, profileID *uuid.UUID)`
**Action**: Update method signature to match design

### 2. RBAC Service - Missing Core Methods
**File**: `internal/service/rbac/casbin.go`
**Design**: Specifies 7 methods (Enforce, AddRoleForUser, RemoveRoleForUser, etc.)
**Code**: Different structural approach, missing exact methods
**Action**: Implement missing methods per design specification

### 3. Configuration - Wrong Location
**Design**: `internal/config/config.go`
**Code**: `pkg/config/config.go`
**Action**: Move config to internal/config/ per design

### 4. Metadata Service - Missing Central Aggregation
**Design**: Central service orchestrating providers with priority fallback
**Code**: Only individual providers, no aggregation service
**Action**: Implement central MetadataService with provider selection logic

### 5. Radarr Provider - Stub Only
**File**: `internal/service/metadata/radarr/provider.go`
**Design**: Priority 1 provider (Servarr-first principle)
**Code**: Returns `ErrUnavailable` - NOT IMPLEMENTED
**Action**: Implement actual Radarr API integration

### 6. TVShow Module - Not Registered
**File**: `cmd/revenge/main.go`
**Design**: TVShow module should be active
**Code**: `tvshow.ModuleWithRiver` exists but NOT registered in main.go
**Action**: Add `tvshow.ModuleWithRiver` to fx.New()

### 7. Adult Module - Missing Obfuscation Layer
**Design**: Schema `qar` with Queen Anne's Revenge themed obfuscation
**Code**: Schema `c` with plain field names
**Action**: Implement full obfuscation layer per ADULT_CONTENT_SYSTEM.md

### 8. Adult Module - Missing Access Controls
**Design**: Scoped permissions, audit logging, PIN protection
**Code**: Only boolean flag `adultEnabled`
**Action**: Implement full access control framework

### 9. Adult Module - Missing Services
**Design**: WhisparrClient, StashDBClient, StashAppClient, FingerprintService
**Code**: None implemented
**Action**: Implement external data source integrations

### 10. Missing Content Modules (8 of 12)
**Design**: movie, tvshow, music, audiobook, book, podcast, photo, livetv, comics, collection + adult
**Code**: Only movie, tvshow, c.movie, c.scene
**Missing**: music, audiobook, book, podcast, photo, livetv, comics, collection
**Action**: Implement remaining 8 content modules

### 11. os.Exit() in Production Code
**File**: `cmd/revenge/main.go:299`
**Design**: "Panic for errors - Crashes server" is anti-pattern
**Code**: `os.Exit(1)` in server startup goroutine
**Action**: Return error and let Fx handle graceful shutdown

### 12. Missing Activity Data Encryption
**Design**: AES-256-GCM encryption for activity data at rest
**Code**: Activity data stored in plaintext
**Action**: Implement encryption layer per DESIGN_PRINCIPLES.md

---

## HIGH Severity Issues

### Services

| Issue | File | Design | Code | Action |
|-------|------|--------|------|--------|
| OIDC - Missing UpdateProvider | service.go | Method specified | Not implemented | Add UpdateProvider() |
| OIDC - No API bindings | api/ | Exposed via API | No handlers | Wire OIDC endpoints |
| User Service - Incomplete design doc | USER.md | Missing methods | N/A | Update design doc |
| Metadata - Missing `enabled` config | config | Both providers | Missing field | Add enabled boolean |

### Content Modules

| Issue | File | Design | Code | Action |
|-------|------|--------|------|--------|
| TVShow - 16 undocumented tables | TVSHOW_MODULE.md | Basic schema | 16 extra tables | Update design or remove |
| TVShow - 50+ undocumented methods | repository | Basic methods | 50+ extra | Document or remove |
| Movie - Entity field mismatches | MOVIE_MODULE.md | Specific fields | Different fields | Align entity structures |
| c.show directory empty | c/show/ | Adult scene module | 0 files | Implement module |

### Architecture

| Issue | File | Design | Code | Action |
|-------|------|--------|------|--------|
| TVShow handler not in DI | api/module.go | Required | Missing TvshowService | Add to handler deps |
| Adult modules not in main.go | main.go | Should be active | Not registered | Register modules |
| Health checks incomplete | main.go | DB, Cache, Search | Only DB | Enable cache/search checks |

### Tech Stack

| Issue | Library | Design | Code | Action |
|-------|---------|--------|------|--------|
| failsafe-go missing | go.mod | Required | Custom impl | Add library or update doc |
| golang-migrate missing | go.mod | Required | Not found | Add library |
| coder/websocket missing | go.mod | Required | Not found | Add library |
| govips missing | go.mod | Required | Not found | Add library |
| casbin undocumented | TECH_STACK.md | Not listed | Used | Add to design doc |

### Features

| Issue | Feature | Design | Code | Action |
|-------|---------|--------|------|--------|
| Watch Next endpoint missing | tvshows API | Specified | No handler | Implement handler |
| Up Next queue missing | playback | Full feature | Not implemented | Implement feature |
| Cross-device sync missing | playback | WebSocket + polling | Not implemented | Implement feature |
| User preference fields missing | watch progress | 4 fields | None | Add to schema |
| Resource grants missing | RBAC | Full table | Not implemented | Implement feature |
| Metadata audit logging missing | RBAC | Full system | Basic activity log | Redesign schema |
| Request system permissions | RBAC | 15 permissions | Defined but no module | Implement module |

---

## MEDIUM Severity Issues

### Services

| Issue | Details |
|-------|---------|
| OIDC ParseClaimMapping silent errors | Returns defaults instead of errors |
| Activity migration comment wrong | References 000008 instead of 000010 |
| Settings fully compliant | No issues |
| APIKeys fully compliant | No issues |

### Architecture

| Issue | Details |
|-------|---------|
| Cache layers incomplete | Only rueidis, missing otter/sturdyc integration |
| Search implementation partial | Basic wrapper only, no per-module collections |
| River workers incomplete | Some workers optional |
| Config missing modules section | Enable/disable per module not implemented |

### Features

| Issue | Details |
|-------|---------|
| 30-day filter missing | Continue watching has no time window |
| Series progress approach differs | Uses aggregated table vs direct scan |
| Event triggers not modeled | Logic distributed, not centralized |
| Adult PIN protection missing | Design specifies, not implemented |
| Adult async jobs missing | No River jobs for fingerprinting |

---

## LOW Severity Issues

| Issue | Details |
|-------|---------|
| TVShow episode watch history naming | `is_completed` vs `is_watching` |
| typesense-go alpha version | v4.0.0-alpha2 in production |
| fsnotify unused | In go.mod but not imported |
| API route TODO comments | Incomplete awareness markers |
| OpenAPI adult endpoints | Spec integration unclear |
| CollectionInfo type extra | Not in design but implemented |

---

## Module Implementation Status

| Module | Design | Code | Status |
|--------|--------|------|--------|
| movie | Yes | Yes | 85% |
| tvshow | Yes | Yes (not registered) | 75% |
| music | Yes | No | 0% |
| audiobook | Yes | No | 0% |
| book | Yes | No | 0% |
| podcast | Yes | No | 0% |
| photo | Yes | No | 0% |
| livetv | Yes | No | 0% |
| comics | Yes | No | 0% |
| collection | Yes | No | 0% |
| c.movie (adult) | Yes | Partial | 35% |
| c.scene (adult) | Yes | Partial | 35% |

---

## Service Implementation Status

| Service | Design | Code | Status |
|---------|--------|------|--------|
| Auth | Yes | Yes | 100% |
| User | Incomplete | Yes | N/A |
| Session | Yes | Yes | 95% (signature issue) |
| Library | Yes | Yes | 100% |
| RBAC | Yes | Partial | 60% |
| OIDC | Yes | Yes | 90% |
| Activity | Yes | Yes | 95% |
| Settings | Yes | Yes | 100% |
| APIKeys | Yes | Yes | 100% |
| Metadata | Yes | Partial | 40% |

---

## Feature Implementation Status

| Feature | Design | Code | Status |
|---------|--------|------|--------|
| Watch Next / Continue | Yes | Partial | 60% |
| RBAC with Casbin | Yes | Partial | 70% |
| Adult Content System | Yes | Partial | 35% |
| Access Controls | Planning | No | 0% |

---

## Recommended Fix Priority

### Phase 1: Critical Blockers
1. Fix Session.UpdateActivity signature
2. Register TVShow module in main.go
3. Move config to internal/config/
4. Fix os.Exit() → proper error handling
5. Implement Radarr provider (Servarr-first)

### Phase 2: Core Features
6. Implement central MetadataService
7. Implement Watch Next endpoint handlers
8. Implement TVShow continue watching handler
9. Add missing RBAC methods
10. Implement resource_grants table

### Phase 3: Adult Module
11. Implement obfuscation layer (qar schema)
12. Implement access control framework
13. Implement external integrations (Whisparr, StashDB)
14. Implement fingerprinting service

### Phase 4: Missing Modules
15. Implement remaining 8 content modules
16. Add missing tech stack libraries
17. Implement cross-device sync
18. Implement Up Next queue

### Phase 5: Polish
19. Add activity data encryption
20. Complete health checks
21. Fix all LOW severity items
22. Update incomplete design docs

---

## Notes

- Design docs under `docs/dev/design/` are the **source of truth**
- Other documentation may be outdated
- Features marked as "PLANNING" in design docs are not yet required but should be tracked
- Code should be adapted to match design, not vice versa
