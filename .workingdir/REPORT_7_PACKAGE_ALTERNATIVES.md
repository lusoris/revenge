# Package Alternatives Analysis Report (2026-02-06)

**Purpose**: Identify packages that may have faster/better alternatives
**Scope**: Application-level packages (excluding pinned: db, cache, job queue)

---

## Executive Summary

| Package | Current | Alternative | Verdict |
|---------|---------|-------------|---------|
| HTTP Client | go-resty/resty v2 | imroc/req v3 | **CONSIDER** |
| Validation | go-playground/validator | - | KEEP |
| Logging | zap | slog (stdlib) | KEEP (zap) |
| UUID | google/uuid | gofrs/uuid or rs/xid | **CONSIDER** |
| Decimal | shopspring/decimal | govalues/decimal | **CONSIDER** |
| Image | golang.org/x/image | libvips-based | **CONSIDER** |
| Config | knadh/koanf | - | KEEP |

---

## Detailed Analysis

### 1. HTTP Client: go-resty/resty v2

**Current Usage**: TMDb, TVDb, Radarr, Sonarr API clients

**Alternatives Researched**:

| Library | Performance | Features | Notes |
|---------|-------------|----------|-------|
| **go-resty/resty** | Standard (net/http wrapper) | Rich middleware, retry, auth | Currently used |
| **imroc/req v3** | Standard (net/http wrapper) | HTTP/1.1, HTTP/2, HTTP/3, auto-detect, better debug | More modern API |
| **valyala/fasthttp** | 4x faster than net/http | Raw performance | Breaking changes, different API |

**Recommendation**: **CONSIDER imroc/req v3**
- Same performance tier as resty (both wrap net/http)
- Native HTTP/2 and HTTP/3 support with auto-detection
- Better debugging/tracing built-in
- Actively maintained (Dec 2025 release)
- Migration would be moderate effort

**Migration Risk**: Medium - Different API but similar patterns

**Sources**:
- [Web Scraping FYI - Resty vs Req](https://webscraping.fyi/lib/compare/go-req-vs-go-resty/)
- [imroc/req GitHub](https://github.com/imroc/req)
- [FastHTTP Discussion](https://github.com/go-resty/resty/discussions/526)

---

### 2. Validation: go-playground/validator v10

**Current Usage**: Struct validation across API handlers

**Performance**:
- Field success: 27.88 ns/op, 0 allocs
- Field failure: 121.3 ns/op, 4 allocs/op

**Alternatives**:
| Library | Approach | Best For |
|---------|----------|----------|
| **go-playground/validator** | Struct tags | REST APIs, declarative |
| **ozzo-validation** | Programmatic | Complex business rules |

**Recommendation**: **KEEP go-playground/validator**
- Excellent performance with zero allocations on success
- Industry standard for REST APIs
- Tag-based validation is perfect for our ogen-generated structs
- No compelling reason to change

**Sources**:
- [Leapcell - Go Validation Libraries](https://leapcell.io/blog/exploring-golang-s-validation-libraries)
- [Dalton Tan - Comparison](https://daltontan.com/comparison-of-golang-input-validator-libraries/29/)

---

### 3. Logging: zap

**Current Usage**: Structured logging throughout application

**2025 Benchmarks**:
| Library | ns/op | B/op | allocs/op |
|---------|-------|------|-----------|
| zerolog | Fastest | 40 | ~1 |
| zap | Very Fast | 168 | 3 |
| slog | Fast | 40 | ~1 |

**Recommendation**: **KEEP zap**
- Already integrated with fx, otel, and other infrastructure
- Production-proven, battle-tested
- slog is an option for new projects but migration cost is high
- zerolog is marginally faster but zap's ecosystem integration is stronger

**Consideration**: For new services, slog (stdlib) is worth considering for reduced dependencies

**Sources**:
- [Better Stack - Go Logging Benchmarks](https://github.com/betterstack-community/go-logging-benchmarks)
- [Dwarves Foundation - slog Benchmarks](https://dwarvesf.hashnode.dev/go-1-21-release-slog-with-benchmarks-zerolog-and-zap)
- [Dash0 - Best Go Logging 2025](https://www.dash0.com/faq/best-go-logging-tools-in-2025-a-comprehensive-guide)

---

### 4. UUID: google/uuid

**Current Usage**: Entity IDs throughout the application

**Alternatives**:
| Library | Type | Size | Sortable | Performance |
|---------|------|------|----------|-------------|
| **google/uuid** | UUIDv4 | 16 bytes | No | Standard |
| **gofrs/uuid** | UUIDv4/v7 | 16 bytes | v7 yes | Standard |
| **rs/xid** | XID | 12 bytes | Yes | Faster |
| **segmentio/ksuid** | KSUID | 20 bytes | Yes | Fast |

**Recommendation**: **CONSIDER gofrs/uuid or rs/xid**

For new entity IDs:
- **gofrs/uuid v7**: Timestamp-sortable UUIDs, RFC-9562 compliant, drop-in compatible
- **rs/xid**: 12-byte, globally unique, sortable, optimized for performance and indexing

**Benefits of sortable IDs**:
- Better database index performance (B-tree friendly)
- Natural chronological ordering
- Reduced index fragmentation

**Migration Strategy**: Could adopt for new entities while keeping google/uuid for existing

**Sources**:
- [Generating Good Unique IDs in Go](https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html)
- [XID: The GUID Alternative](https://lumochift.org/blog/xid-the-guid-alternative)
- [gofrs/uuid GitHub](https://github.com/gofrs/uuid)

---

### 5. Decimal: shopspring/decimal

**Current Usage**: Financial calculations, prices

**Benchmarks**:
| Library | Add (ns/op) | Precision | Allocations |
|---------|-------------|-----------|-------------|
| **shopspring/decimal** | 35.58 | Arbitrary | Higher |
| **alpacahq/alpacadecimal** | 1.385 | 12 digits | Lower |
| **govalues/decimal** | Fast | 19 digits | Zero |
| **cockroachdb/apd** | Fast | Arbitrary | Mutable API |

**Recommendation**: **CONSIDER govalues/decimal**
- Zero-allocation design
- 19 digits of precision (sufficient for most financial use cases)
- Immutable API (safer)
- Significant performance improvement over shopspring

**Note**: If arbitrary precision is required, keep shopspring/decimal

**Sources**:
- [govalues/decimal GitHub](https://github.com/govalues/decimal)
- [alpacahq/alpacadecimal GitHub](https://github.com/alpacahq/alpacadecimal)
- [shopspring/decimal GitHub](https://github.com/shopspring/decimal)

---

### 6. Image Processing: golang.org/x/image

**Current Usage**: WebP decoding, image dimension detection

**Performance Comparison**:
| Library | Speed vs ImageMagick | Memory | Dependencies |
|---------|---------------------|--------|--------------|
| **x/image** | Standard Go | Moderate | None |
| **bimg** | 4-8x faster | Low | libvips (C) |
| **govips** | 4-8x faster | Low | libvips (C) |
| **imagor** | 4-8x faster | Low | libvips (C) |

**Current Implementation**: Only doing format detection and dimension reading, not heavy processing.

**Recommendation**: **CONSIDER if heavy processing is added**
- For current use (decode config only): x/image is adequate
- For thumbnails, resizing, format conversion: Switch to govips or bimg
- libvips-based solutions are 4-8x faster with lower memory

**Trade-off**: libvips requires CGO and system dependencies

**Sources**:
- [bimg GitHub](https://github.com/h2non/bimg)
- [govips GitHub](https://github.com/davidbyttow/govips)
- [imagor GitHub](https://github.com/cshum/imagor)
- [Transloadit - libvips + Go](https://transloadit.com/devtips/vips-in-combination-with-go/)

---

### 7. Configuration: knadh/koanf v2

**Current Usage**: Application configuration management

**Comparison with Viper**:
| Aspect | koanf | viper |
|--------|-------|-------|
| Binary size | 1x (baseline) | 3.13x larger |
| Dependencies | Minimal | Many |
| Key handling | Preserves case | Forces lowercase |
| Modularity | Provider/Parser plugins | Monolithic |

**Recommendation**: **KEEP koanf**
- Already the better choice over viper
- Lightweight, modular, fewer dependencies
- Actively maintained
- No compelling alternatives

**Sources**:
- [koanf Wiki - Viper Comparison](https://github.com/knadh/koanf/wiki/Comparison-with-spf13-viper)
- [ITNEXT - Viper vs Koanf](https://itnext.io/golang-configuration-management-library-viper-vs-koanf-eea60a652a22)
- [Three Dots Labs - Recommended Libraries](https://threedots.tech/post/list-of-recommended-libraries/)

---

## Action Items

### Priority 1 - Quick Wins (Low Risk, High Value)

1. [ ] **Evaluate govalues/decimal** for financial calculations
   - Run benchmarks with actual use cases
   - Verify 19-digit precision is sufficient
   - Estimate: 2-4 hours

2. [ ] **Evaluate gofrs/uuid v7** for new entity types
   - Test sortable UUID generation
   - Benchmark against google/uuid
   - Estimate: 1-2 hours

### Priority 2 - Medium Term (Moderate Effort)

3. [ ] **Evaluate imroc/req v3** for HTTP clients
   - Create prototype TMDb client with req
   - Compare API ergonomics
   - Test HTTP/2 and HTTP/3 features
   - Estimate: 1 day

### Priority 3 - Future (When Needed)

4. [ ] **Add govips/bimg** when image processing is added
   - Only if thumbnailing, resizing, or format conversion is needed
   - Requires libvips system dependency

---

## Packages That Are Already Optimal

| Package | Reason to Keep |
|---------|---------------|
| pgx/v5 | Best PostgreSQL driver for Go |
| rueidis | Fastest Redis client for Go |
| otter | Best in-memory cache |
| river | Native pgx job queue |
| fx | Industry-standard DI |
| ogen | Type-safe OpenAPI codegen |
| zap | Proven logging with excellent ecosystem |
| koanf | Lightweight config, better than viper |
| go-playground/validator | Standard, fast, tag-based |

---

## Notes

- This analysis excludes pinned packages (database, cache, job queue)
- Performance benchmarks are from 2025-2026 sources
- Trade-offs between performance and API ergonomics were considered
- Migration costs were factored into recommendations

---

*Generated: 2026-02-06*
