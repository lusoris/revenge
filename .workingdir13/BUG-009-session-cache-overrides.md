# BUG-009: Session CachedService missing overrides for RefreshSession and RevokeAllUserSessionsExcept

**Severity:** MEDIUM  
**Category:** Cache  
**Status:** RESOLVED

## Description

The session `CachedService` at `internal/service/session/cached_service.go` properly overrides `CreateSession`, `ValidateSession`, `RevokeSession`, and `RevokeAllUserSessions` — but does NOT override:

1. `RefreshSession` — creates new token pair, but old token hash remains in L1/L2 cache for up to `SessionTTL` (30s)
2. `RevokeAllUserSessionsExcept` — revokes all sessions but one, but doesn't invalidate individual cache entries

This means revoked/refreshed sessions can be served from cache until TTL expires.

## Fix

Override both methods in the CachedService to properly invalidate cache entries.
