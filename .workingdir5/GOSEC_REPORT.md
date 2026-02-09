# Gosec Security Scan Report

**Date**: 2026-02-07
**Tool**: gosec v2.22.4 (Go security checker)
**Scope**: `./internal/...` (all authored code)

## Summary

| Severity | Count | Notes |
|----------|-------|-------|
| HIGH     | 156   | 84 G115 + 63 G101 + 9 G109 |
| MEDIUM   | 9     | 4 G301 + 3 G304 + 2 G204 |
| **Total**| **165** | |

**Actionable findings**: ~18 (G115 integer overflows worth fixing, G204/G304/G301 in playback)
**False positives**: ~147 (G101 all in sqlc-generated code, most G115 in metadata adapters are safe)

---

## Findings by Rule

### G115 — Integer Overflow Conversion (84 HIGH)

**CWE-190**: `int64 -> int32` conversions that could overflow.

| File | Count | Risk |
|------|-------|------|
| `internal/content/tvshow/jobs/jobs.go` | 17 | Low — Sonarr IDs, won't exceed int32 |
| `internal/service/metadata/adapters/tvshow/adapter.go` | 12 | Low — TMDb/TVDb IDs |
| `internal/service/metadata/providers/tmdb/mapping.go` | 11 | Low — TMDb API response fields |
| `internal/service/metadata/adapters/movie/adapter.go` | 10 | Low — TMDb/TVDb IDs |
| `internal/service/metadata/providers/tvdb/mapping.go` | 7 | Low — TVDb API response fields |
| `internal/integration/sonarr/mapper.go` | 6 | Low — Sonarr IDs |
| `internal/api/handler_metadata.go` | 5 | Low — API request params |
| `internal/integration/sonarr/service.go` | 4 | Low — Sonarr API fields |
| Other files | 12 | Low |

**Assessment**: These are almost entirely external API integer fields (TMDb, TVDb, Sonarr) being mapped to `int32` database columns. TMDb IDs are currently ~1.2M (well within int32 range of 2.1B). Sonarr IDs are local database IDs. **No real-world overflow risk**.

**Recommendation**: Helpers already exist — `util.SafeInt64ToInt32()` (clamping) and `validate.SafeInt32()` (error-returning) in `internal/util/safeconv.go` and `internal/validate/convert.go`. Some files already use them (API handlers, radarr/sonarr mappers, library matcher). The remaining 84 flagged locations (metadata adapters, tvshow/jobs, tmdb/tvdb providers) should be migrated to use these helpers. Medium priority — straightforward mechanical changes.

### G101 — Potential Hardcoded Credentials (63 HIGH)

**CWE-798**: Strings containing "token", "key", "password" etc.

| File | Count | Notes |
|------|-------|------|
| `internal/infra/database/db/auth_tokens.sql.go` | 25 | sqlc-generated SQL queries |
| `internal/infra/database/db/apikeys.sql.go` | 13 | sqlc-generated SQL queries |
| `internal/infra/database/db/mfa.sql.go` | 12 | sqlc-generated SQL queries |
| `internal/infra/database/db/sessions.sql.go` | 5 | sqlc-generated SQL queries |
| `internal/content/tvshow/db/credits.sql.go` | 4 | sqlc-generated SQL queries |
| `internal/content/movie/db/movies.sql.go` | 2 | sqlc-generated SQL queries |
| `internal/infra/database/db/users.sql.go` | 1 | sqlc-generated SQL queries |
| `internal/service/auth/service.go` | 1 | `tokenHash` variable name |

**Assessment**: **All false positives**. 62/63 are in sqlc-generated code where SQL query strings naturally contain column names like `token_hash`, `api_key`, `hashed_password`. The 1 in auth/service.go is a variable name `tokenHash` (hashed value, not a credential).

**Recommendation**: Add `//nolint:gosec` or configure gosec exclusion for `db/*.sql.go` files. No action needed on auth/service.go.

### G109 — strconv.Atoi Integer Overflow (9 HIGH)

**CWE-190**: `strconv.Atoi` result converted to `int16/int32`.

| File | Count | Notes |
|------|-------|------|
| `internal/service/metadata/providers/tvdb/provider.go` | 3 | TVDb API string→int parsing |
| `internal/service/metadata/providers/tvdb/mapping.go` | 3 | TVDb API string→int parsing |
| `internal/service/metadata/providers/tmdb/provider.go` | 3 | TMDb API string→int parsing |

**Assessment**: Same category as G115 — parsing external API strings to int32. TVDb/TMDb IDs are small numbers. **Low real-world risk**.

**Recommendation**: Could use `strconv.ParseInt(s, 10, 32)` instead of `strconv.Atoi` + cast for explicit bounds checking. Low priority.

### G204 — Subprocess Launched with Variable (2 MEDIUM)

**CWE-78**: Command injection risk.

| File | Notes |
|------|-------|
| `internal/playback/transcode/ffmpeg.go` | FFmpeg exec with config-sourced path |
| `internal/playback/subtitle/extract.go` | FFmpeg exec with config-sourced path |

**Assessment**: The FFmpeg path comes from `PlaybackConfig.FFmpegPath` (server config, not user input). Input file paths come from database records (MovieFile.FilePath, EpisodeFile.FilePath), not direct user input. **Low risk** in practice, but worth validating.

**Recommendation**:
- Validate `FFmpegPath` at startup (must be absolute path, must exist)
- Ensure file paths from DB are not user-controllable (they come from library scanning)
- Consider using `exec.LookPath` validation

### G304 — File Inclusion via Variable (3 MEDIUM)

**CWE-22**: Path traversal risk.

| File | Count | Notes |
|------|-------|------|
| `internal/service/storage/storage.go` | 2 | Serving media files from storage paths |
| `internal/playback/hls/manifest.go` | 1 | Reading FFmpeg-generated playlist from segment dir |

**Assessment**: Storage paths are constructed from config base paths + database file paths. HLS manifest paths are constructed from session segment directories. **Medium risk** — a malicious filename in the database could enable path traversal.

**Recommendation**:
- Add `filepath.Clean` + prefix validation on all file paths before serving
- Ensure served paths stay within configured base directories
- Priority: **Medium** — worth fixing for defense-in-depth

### G301 — Directory Created with Permissions > 0750 (4 MEDIUM)

**CWE-276**: Overly permissive directory creation.

| File | Count | Notes |
|------|-------|------|
| `internal/playback/transcode/pipeline.go` | 2 | Segment directories (0755) |
| `internal/playback/subtitle/extract.go` | 1 | Subtitle output directory (0755) |
| `internal/playback/service.go` | 1 | Session directory (0755) |

**Assessment**: These are temporary segment directories under `/tmp/revenge-segments/`. Using 0755 instead of 0750 means "other" users can read (but not write). In a container deployment, this is irrelevant.

**Recommendation**: Change to `0750` for principle of least privilege. **Low priority** — trivial fix.

---

## Priority Action Items

| Priority | Issue | Action | Effort |
|----------|-------|--------|--------|
| Medium | G304 path traversal (3) | Add filepath.Clean + prefix validation | 1-2h |
| Low | G301 dir permissions (4) | Change 0755 → 0750 | 10min |
| Low | G204 subprocess (2) | Validate FFmpegPath at startup | 30min |
| Low | G115/G109 int overflow (93) | Migrate to existing `util.SafeInt64ToInt32` / `validate.SafeInt32` helpers | 2-3h |
| None | G101 false positives (63) | Configure gosec exclusion for sqlc | 10min |

---

## Gosec Configuration Recommendation

Create `.gosec.conf` or add to CI:

```bash
gosec -exclude=G101 -exclude-dir=internal/infra/database/db -exclude-dir=internal/api/ogen ./internal/...
```

This would reduce noise from 165 to ~39 actionable findings.
