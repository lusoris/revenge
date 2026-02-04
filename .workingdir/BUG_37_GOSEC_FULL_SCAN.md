# Bug #37: gosec Full Codebase Security Scan

**Status**: OPEN
**Date**: 2026-02-04
**Severity**: MIXED (HIGH/MEDIUM/LOW)
**Component**: Multiple

## Summary

Full gosec scan found 117 issues across the codebase:

| Rule | Count | Severity | Description |
|------|-------|----------|-------------|
| G101 | 58 | HIGH | Hardcoded credentials (mostly false positives) |
| G115 | 33 | HIGH | Integer overflow conversions |
| G602 | 20 | LOW | Potential slice bounds issue |
| G301 | 2 | MEDIUM | Directory permissions |
| G112 | 1 | MEDIUM | Potential slowloris attack |
| G304 | 1 | MEDIUM | File path controlled by attacker |
| G306 | 1 | MEDIUM | File write permissions |
| G104 | 1 | LOW | Unhandled error |

## Analysis Required

### G101: Hardcoded Credentials (58 instances)

Most likely false positives in:
- Test files with test credentials
- Constant names containing "password", "secret", "token"
- Configuration examples

**Action**: Review each, add #nosec where appropriate with justification

### G115: Integer Overflow Conversions (33 instances)

**Files affected**:
- `internal/content/movie/mediainfo_types.go`
- `internal/content/movie/mediainfo.go`
- `internal/content/movie/library_scanner.go`
- `internal/content/movie/tmdb_mapper.go`
- `internal/integration/radarr/service.go`
- `internal/integration/radarr/mapper.go`

**Action**: Add safe conversion helpers or bounds checking

### G602: Slice Bounds (20 instances)

Potential slice access out of bounds issues.

**Action**: Review each, add bounds checking where necessary

### G301/G306: File Permissions (3 instances)

Directory and file write permissions may be too permissive.

**Action**: Review and adjust permissions

### G112: Slowloris (1 instance)

Potential slowloris attack via unbounded request handling.

**Action**: Add timeouts/limits

### G304: File Path (1 instance)

File path potentially controlled by attacker input.

**Action**: Add path validation

### G104: Unhandled Error (1 instance)

**Action**: Add error handling

## Priority

1. **HIGH**: G115 (integer overflow) - real security risk
2. **MEDIUM**: G301/G304/G306/G112 - file/network security
3. **LOW**: G101 (mostly false positives), G602, G104

## Question for User

How do you want to proceed?
- A) Fix all issues now (comprehensive)
- B) Fix HIGH severity only (G115)
- C) Fix HIGH + MEDIUM, defer LOW
- D) Just document and defer all to later sprint
