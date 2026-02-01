# Questions - v0.1.0 Implementation

> Collected questions during implementation

---

## Open Questions

<!-- No open questions -->

---

## Resolved Questions

### [Q001] otter package version mismatch ✅ RESOLVED
**SOURCE_OF_TRUTH said**: `github.com/maypok86/otter/v2 v2.x`
**Reality**: Latest version is `v1.2.4`, no v2 exists yet
**Action taken**: Using v1.2.4 in go.mod
**Resolution**: Updated SOURCE_OF_TRUTH to `github.com/maypok86/otter v1.2.4`

### [Q002] zap package version mismatch ✅ RESOLVED
**SOURCE_OF_TRUTH said**: `go.uber.org/zap v1.28.0`
**Reality**: Latest version is `v1.27.1`, v1.28.0 doesn't exist yet
**Action taken**: Using v1.27.1 in go.mod
**Resolution**: Updated SOURCE_OF_TRUTH to `go.uber.org/zap v1.27.1`
