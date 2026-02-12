# BUG-008: OIDC service uses sync.Map instead of L1Cache

**Severity:** MEDIUM
**Category:** Cache
**Status:** RESOLVED

## Description

The OIDC service at `internal/service/oidc/service.go` uses `sync.Map` for provider caching, which:
- Has unbounded memory growth (no eviction)
- Has no TTL (providers never expire)
- Violates the project rule "Never use `sync.Map` for caching"

## Fix

Replace `providerCache sync.Map` with `cache.L1Cache[string, *oidc.Provider]` from `internal/infra/cache/otter.go`.
