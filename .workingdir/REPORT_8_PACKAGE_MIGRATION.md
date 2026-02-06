# Package Migration Report (2026-02-06)

**Purpose**: Migrate suboptimal packages identified in REPORT_7_PACKAGE_ALTERNATIVES.md
**Scope**: 4 package changes across 87 files

---

## Summary

| Change | Files | Commit | Status |
|--------|-------|--------|--------|
| UUID v4 → v7 (google/uuid kept) | 65 | `5bdcbe5a4b` | Done |
| go-resty/resty v2 → imroc/req v3 | 7 | `4186ae1aab` | Done |
| shopspring/decimal → govalues/decimal | 14 | `8adc754b08` | Done |
| x/image → govips (davidbyttow/govips/v2) | 1 | `be0a507950` | Done |

**Build**: Clean (`GOEXPERIMENT=greenteagc,jsonv2 go build ./...`)

---

## 1. UUID v4 → v7 (65 files)

**What**: Changed `uuid.New()` (random v4) to `uuid.Must(uuid.NewV7())` (time-sortable v7)

**Why**:
- RFC 9562 compliant timestamp-sortable UUIDs
- Better PostgreSQL B-tree index performance (sequential inserts)
- Reduced index fragmentation
- Natural chronological ordering

**Approach**: Kept `google/uuid` v1.6.0 (already supports v7). No library change needed.
All 65 non-generated files changed with sed, ogen-generated code left as-is.

**Risk**: None - v7 UUIDs are drop-in compatible with v4 UUID columns.

---

## 2. go-resty v2 → imroc/req v3 (7 files)

**What**: Replaced HTTP client library across all API clients

**Files changed**:
- `internal/service/metadata/providers/tmdb/client.go` - TMDb API client
- `internal/service/metadata/providers/tvdb/client.go` - TVDb API client
- `internal/integration/radarr/client.go` - Radarr API client
- `internal/integration/sonarr/client.go` - Sonarr API client
- `internal/content/shared/metadata/client.go` - BaseClient shared
- `internal/content/shared/metadata/images.go` - Image downloader
- `internal/infra/image/service.go` - Image service

**API mapping applied**:
| resty v2 | req v3 |
|----------|--------|
| `resty.New()` | `req.C()` |
| `.SetRetryCount(n)` | `.SetCommonRetryCount(n)` |
| `.SetRetryWaitTime(d).SetRetryMaxWaitTime(d2)` | `.SetCommonRetryBackoffInterval(d, d2)` |
| `.SetProxy(url)` | `.SetProxyURL(url)` |
| `.SetHeader(k,v)` | `.SetCommonHeader(k,v)` |
| `.SetResult(&r)` | `.SetSuccessResult(&r)` |
| `.SetError(&e)` | `.SetErrorResult(&e)` |
| `.SetAuthToken(t)` | `.SetBearerAuthToken(t)` |
| `resp.IsError()` | `resp.IsErrorState()` |
| `resp.StatusCode()` | `resp.StatusCode` (field) |
| `resp.Status()` | `resp.Status` (field) |
| `resp.Body()` | `resp.Bytes()` |
| `resp.Header()` | `resp.Header` (field) |

**Benefits**: Native HTTP/2 and HTTP/3 with auto-detection, better debugging/tracing.

**Result**: resty completely removed from go.mod (was direct dep, now gone).

---

## 3. shopspring/decimal → govalues/decimal (14 files)

**What**: Replaced decimal library for financial/rating calculations

**Files changed**:
- `internal/content/movie/types.go` - Movie type definitions
- `internal/content/tvshow/types.go` - TV show type definitions
- `internal/content/movie/repository_postgres.go` - DB conversion
- `internal/content/tvshow/repository_postgres.go` - DB conversion
- `internal/content/movie/library_matcher.go` - Library matching logic
- `internal/content/movie/moviejobs/metadata_refresh.go` - Metadata job
- `internal/integration/radarr/mapper.go` - Radarr data mapping
- `internal/integration/sonarr/mapper.go` - Sonarr data mapping
- `internal/service/metadata/adapters/movie/adapter.go` - Movie adapter
- `internal/service/metadata/adapters/tvshow/adapter.go` - TV adapter
- 4 test files

**API mapping applied**:
| shopspring | govalues |
|-----------|----------|
| `decimal.NewFromFloat(f)` | `decimal.NewFromFloat64(f)` (returns error) |
| `decimal.NewFromString(s)` | `decimal.Parse(s)` |
| `decimal.NewFromInt(i)` | `decimal.MustNew(i, 0)` |
| `decimal.Zero` | `decimal.Decimal{}` |
| `d.IsZero()` | `d.IsZero()` (same) |
| `d.String()` | `d.String()` (same) |
| `d.Float64()` | `d.Float64()` (same) |

**Benefits**: Zero-allocation design, 19-digit precision, immutable API, significantly
faster arithmetic (25x faster Add, lower allocs).

**Note**: shopspring remains as indirect dep (used by another package).

---

## 4. x/image → govips (1 file)

**What**: Replaced Go stdlib image decoding with libvips via govips

**File changed**: `internal/api/image_utils.go`

**Before**: Used `image.DecodeConfig()` with registered format decoders
(`image/jpeg`, `image/png`, `image/gif`, `golang.org/x/image/webp`)

**After**: Uses `vips.NewImageFromBuffer(data)` which handles all formats natively
through libvips. Added `sync.Once` initialization for libvips startup.

**Benefits**: 4-8x faster image processing, lower memory, all formats handled natively.

**Requirement**: libvips system dependency (CGO required).

---

## Dependency Changes Summary

### Added
- `github.com/imroc/req/v3` v3.57.0 (HTTP client)
- `github.com/govalues/decimal` v0.1.36 (decimal arithmetic)
- `github.com/davidbyttow/govips/v2` v2.16.0 (image processing)

### Removed (direct)
- `github.com/go-resty/resty/v2` v2.17.1

### Unchanged (still used, kept as-is)
- `github.com/google/uuid` v1.6.0 (now generating v7 instead of v4)

### Remains as indirect
- `github.com/shopspring/decimal` v1.4.0 (used by other dep)
- `golang.org/x/image` v0.35.0 (used by other dep)

---

## Verification

All changes verified with:
```bash
GOEXPERIMENT=greenteagc,jsonv2 go build ./...  # Clean build
go mod tidy                                     # Dependencies resolved
```

---

*Generated: 2026-02-06*
