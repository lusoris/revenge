# Go Packages Research - Media Server Ecosystem

> State-of-the-art Go packages to massively reduce development work

## Research Source

- **Awesome-Go**: https://github.com/avelino/awesome-go (29.5k stars, 188 contributors)
- **Awesome-Go.com**: https://awesome-go.com/ (curated catalog)
- **Categories**: 40+ categories, 500+ packages analyzed

---

## Current Stack vs Alternatives

### 1. HTTP Routing / Web Framework

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **Echo** (current) | 29k | High-performance HTTP framework | Middleware-rich, excellent docs, HTTP/2, WebSocket, auto-cert | Learning curve, some magic | ✅ **KEEP** (good balance) |
| **Chi** | 18k | Lightweight router | Stdlib-compatible, minimal, fast, no dependencies | Less batteries-included | ⚠️ Consider if Echo overhead issues |
| **Gin** | 77k | Fastest HTTP framework | Fastest benchmarks, huge community, JSON validation | Less idiomatic, Chinese-heavy docs | ⚠️ Consider for high-traffic |
| **Fiber** | 33k | Express-inspired API | Express-like (Node.js familiarity), fast, zero allocation | Non-standard (fasthttp, not net/http) | ❌ Avoid (stdlib compatibility) |
| **fasthttp** | 21k | Low-level HTTP | 10x faster than net/http, zero allocation | Not stdlib-compatible, fragile API | ❌ Overkill (premature optimization) |

**Recommendation**: **KEEP Echo** (good docs, middleware ecosystem, HTTP/2, WebSocket built-in)

---

### 2. ORM / Database

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **GORM** (current) | 36k | Feature-rich ORM | Associations, hooks, auto-migration, preloading | Reflection overhead, magic queries | ✅ **KEEP** (unless complexity issues) |
| **sqlc** | 13k | Type-safe SQL code-gen | No runtime overhead, compile-time safety, plain SQL | No associations, manual migrations | ✅ **ALREADY USING** (perfect fit) |
| **ent** | 15k | Graph-based ORM | Schema-as-code (Facebook), type-safe, code-gen | Steep learning curve, complex setup | ❌ Overkill for current needs |
| **Bun** | 4k | Fast SQL builder | 2x faster than GORM, migrations, fixtures | Smaller community, less mature | ⚠️ Consider if GORM bottleneck |

**Recommendation**: **KEEP sqlc** (type-safe, no runtime overhead, plain SQL control)

---

### 3. Video Processing

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **gmf** | 900 | FFmpeg av* bindings | Mature, CGo bindings, all FFmpeg features | CGo dependency, manual memory | ⚠️ Evaluate for Blackbeard |
| **go-astiav** | 400 | Better FFmpeg bindings | Cleaner API than gmf, maintained | Newer (less battle-tested) | ✅ **USE in Blackbeard** (cleaner API) |
| **go-astisub** | 600 | Subtitle parsing | .srt/.stl/.ttml/.webvtt/.ssa/.ass support | Limited to parsing (no rendering) | ✅ **USE** (subtitle support) |
| **gortsplib** | 800 | RTSP server/client | Live streaming, RTSP/RTP/RTCP | Complex for simple use cases | ✅ **USE** (Live TV module) |
| **libvlc-go** | 400 | VLC 2.X/3.X/4.X bindings | Full VLC features, hardware decoding | CGo, VLC dependency | ⚠️ Consider for player fallback |
| **mp4ff** | 500 | MP4 tools | MP4 parsing/manipulation, no CGo | Limited to MP4 only | ✅ **USE** (MP4 manipulation) |

**Recommendation**: 
- **go-astiav** for Blackbeard transcoding (cleaner than gmf)
- **go-astisub** for subtitle parsing
- **gortsplib** for Live TV
- **mp4ff** for MP4 manipulation

---

### 4. Audio Processing

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **beep** | 2.1k | Audio playback | Simple API, multiple formats, streaming | Limited effects, no recording | ⚠️ Player library reference |
| **flac** | 300 | FLAC encoder/decoder | Pure Go, no CGo | FLAC only | ✅ **USE** (lossless audio) |
| **GoAudio** | 700 | Audio processing | Waveform analysis, effects | Heavy, complex API | ❌ Not needed (Blackbeard handles) |
| **malgo** | 700 | Mini audio library | Cross-platform, low-level | Manual audio pipeline | ❌ Too low-level |
| **Oto** | 400 | Low-level playback | Cross-platform, simple | Requires manual mixing | ❌ Too low-level |
| **PortAudio** | 700 | Audio I/O | Professional, cross-platform | CGo dependency | ❌ Not needed (web client) |

**Recommendation**: 
- **flac** for FLAC support (music module)
- Others not needed (Blackbeard transcodes, web client plays)

---

### 5. Image Processing

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **bimg** | 2.7k | libvips wrapper | Lightning fast, low memory, resize/crop | CGo, libvips dependency | ✅ **USE** (poster resizing) |
| **gocv** | 6.5k | OpenCV 3.3+ bindings | Full OpenCV features, ML | Heavy (OpenCV 4+), CGo | ❌ Overkill (no ML needed) |
| **govips** | 1.3k | libvips wrapper | 3-5x faster than imaging, low memory | CGo, libvips dependency | ✅ **Alternative to bimg** |
| **imaginary** | 5.5k | HTTP microservice | Fast, HTTP API, Docker-ready | Separate service, overhead | ⚠️ Consider for external service |
| **imaging** | 5.2k | Simple processing | Pure Go, no dependencies, easy API | Slower than libvips | ⚠️ Fallback if CGo issues |
| **imagor** | 3.5k | Fast+secure processing | libvips + Thumbor API, cache | Complex setup, separate service | ❌ Overkill |

**Recommendation**: 
- **bimg** OR **govips** for poster/fanart resizing (libvips fast)
- **imaging** as pure Go fallback (no CGo)

---

### 6. Distributed Systems / Orchestration

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **Temporal** | 11k | Durable execution | Workflows survive crashes, retries, timers | Complex setup, separate service | ⚠️ Consider for complex workflows |
| **Kratos** | 23k | Microservices framework | Bilibili production, gRPC, HTTP | Microservice-focused (not monolith) | ❌ Not needed (monolith design) |
| **Kitex** | 7k | High-perf RPC | ByteDance production, Thrift/Protobuf | CloudWeGo ecosystem lock-in | ❌ Not needed (no RPC) |
| **NATS** | 15k | Messaging | Lightweight, clustering, JetStream | Separate service, overhead | ⚠️ Consider if messaging needed |

**Recommendation**: 
- **NONE currently** (monolith design, River handles jobs)
- **Temporal** if complex workflows emerge (e.g., multi-step metadata enrichment)

---

### 7. Job Queues

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **River** (current) | 3.2k | PostgreSQL-native queue | No extra service, ACID, modern Go | Newer (less battle-tested) | ✅ **KEEP** (perfect fit) |
| **Asynq** | 9.5k | Redis-based queue | Battle-tested, cron, retries, monitoring | Requires Redis (we use Dragonfly) | ⚠️ Alternative if River issues |
| **Machinery** | 7.5k | Async task queue | Multiple brokers (Redis, AMQP, SQS) | Complex setup, overhead | ❌ Too complex |

**Recommendation**: **KEEP River** (PostgreSQL-native, no extra service, modern API)

---

### 8. Serialization / JSON

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **sonic** | 7k | Blazingly fast JSON | Bytedance, JIT compilation, 3x faster | x86-64 only (no ARM fallback) | ⚠️ Evaluate for API responses |
| **jsoniter** | 13k | 100% compatible JSON | Drop-in replacement, 2x faster | Less dramatic speedup | ⚠️ Consider if JSON bottleneck |
| **protobuf** | 1.4k | Protocol Buffers | Type-safe, compact, gRPC | Requires .proto files, complex | ❌ Not needed (REST API) |

**Recommendation**: 
- **stdlib encoding/json** (good enough for now)
- **sonic** if JSON serialization becomes bottleneck (API responses)

---

### 9. Validation

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **go-playground/validator** | 16k | Struct validation | Comprehensive tags, custom validators | Tag-heavy, magic strings | ✅ **USE** (API request validation) |
| **ozzo-validation** | 3.7k | Fluent validation | Type-safe, code-based, clear errors | More verbose than tags | ⚠️ Alternative if tag issues |

**Recommendation**: **go-playground/validator** (industry standard, Echo integrates)

---

### 10. Authentication / Authorization

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **ory/hydra** | 15k | OAuth2/OIDC server | Production-ready, OAuth2 flows | Separate service, complex | ⚠️ Consider for OIDC provider |
| **ory/kratos** | 11k | Identity management | User management, MFA, passwordless | Separate service, learning curve | ❌ Not needed (custom user service) |
| **casbin** | 17k | RBAC/ABAC | Flexible, adapters, cloud-ready | Complex policies, overhead | ⚠️ Consider for advanced RBAC |

**Recommendation**: 
- **NONE currently** (custom JWT + OIDC client sufficient)
- **casbin** if complex RBAC emerges (permissions, roles, policies)

---

### 11. Caching

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **go-redis/v9** (current) | 20k | Redis/Dragonfly client | Industry standard, Dragonfly compatible | N/A | ✅ **KEEP** (Dragonfly) |
| **ristretto** | 5.5k | In-memory cache | High-perf, cost-based eviction | In-process only (no distributed) | ⚠️ Consider for local cache |
| **freecache** | 7k | Zero GC overhead | No GC pauses, fast | Limited features, in-process | ⚠️ Consider for hot paths |

**Recommendation**: 
- **KEEP go-redis/v9** (Dragonfly cache)
- **ristretto** for in-memory hot cache (e.g., session tokens)

---

### 12. Search

| Package | Stars | Purpose | Pros | Cons | Revenge Fit |
|---------|-------|---------|------|------|-------------|
| **typesense-go** (current) | 100+ | Typesense client | Fast, typo-tolerant, facets | Smaller community | ✅ **KEEP** (Typesense) |
| **elastic/go-elasticsearch** | 5.5k | Elasticsearch client | Official, battle-tested | Elasticsearch complexity | ❌ Not needed (Typesense simpler) |
| **meilisearch-go** | 400+ | Meilisearch client | Fast, easy, Rust-based | Smaller ecosystem | ⚠️ Alternative to Typesense |

**Recommendation**: **KEEP typesense-go** (fast, simple, typo-tolerant)

---

## Recommended Additions

### High Priority (Immediate Use)

| Package | Purpose | Use Case |
|---------|---------|----------|
| **go-playground/validator** | Struct validation | API request validation |
| **go-astiav** | FFmpeg bindings | Blackbeard transcoding |
| **go-astisub** | Subtitle parsing | Subtitle support (.srt, .ass, .webvtt) |
| **gortsplib** | RTSP server/client | Live TV module |
| **mp4ff** | MP4 manipulation | MP4 file processing |
| **bimg** OR **govips** | Image processing | Poster/fanart resizing |
| **flac** | FLAC codec | Music module lossless support |

### Medium Priority (Evaluate)

| Package | Purpose | Use Case |
|---------|---------|----------|
| **sonic** | Fast JSON | API response serialization (if bottleneck) |
| **ristretto** | In-memory cache | Hot cache (session tokens, user data) |
| **Temporal** | Durable workflows | Complex multi-step workflows |
| **casbin** | RBAC | Advanced permission system |

### Low Priority (Monitor)

| Package | Purpose | Use Case |
|---------|---------|----------|
| **Asynq** | Redis job queue | Alternative to River (if issues) |
| **Chi** | Lightweight router | Alternative to Echo (if overhead) |
| **jsoniter** | Fast JSON | Alternative to sonic (ARM support) |
| **ozzo-validation** | Validation | Alternative to go-playground |

---

## Packages to AVOID

| Package | Reason |
|---------|--------|
| **Fiber** | Not stdlib-compatible (fasthttp), fragile |
| **fasthttp** | Premature optimization, fragile API |
| **ent** | Overkill for current schema complexity |
| **Kratos/Kitex** | Microservice-focused (not monolith) |
| **Machinery** | Too complex for simple job queue |
| **ory/kratos** | User management separate service (we have custom) |
| **gocv** | OpenCV overkill (no ML needed) |
| **GoAudio** | Too heavy (Blackbeard handles audio) |
| **PortAudio** | CGo, not needed (web client) |

---

## Summary

### KEEP Current Stack
- ✅ Echo (HTTP framework)
- ✅ sqlc (type-safe SQL)
- ✅ River (PostgreSQL job queue)
- ✅ go-redis/v9 (Dragonfly cache)
- ✅ typesense-go (search client)

### ADD Immediately
- ✅ **go-playground/validator** (API validation)
- ✅ **go-astiav** (Blackbeard FFmpeg)
- ✅ **go-astisub** (subtitles)
- ✅ **gortsplib** (Live TV)
- ✅ **mp4ff** (MP4 tools)
- ✅ **bimg/govips** (image resizing)
- ✅ **flac** (lossless audio)

### EVALUATE Later
- ⚠️ **sonic** (fast JSON if bottleneck)
- ⚠️ **ristretto** (in-memory cache)
- ⚠️ **Temporal** (complex workflows)
- ⚠️ **casbin** (advanced RBAC)

### Time Savings Estimate
- **FFmpeg bindings** (go-astiav): Save 2-4 weeks (vs writing from scratch)
- **Subtitle parsing** (go-astisub): Save 1-2 weeks (complex format support)
- **RTSP** (gortsplib): Save 3-4 weeks (Live TV streaming)
- **Image processing** (bimg): Save 1 week (vs pure Go)
- **Validation** (go-playground/validator): Save 3-5 days (vs manual)
- **TOTAL**: ~8-15 weeks saved

---

## Next Steps

1. ✅ Add go-playground/validator to `go.mod`
2. ✅ Prototype Blackbeard with go-astiav
3. ✅ Test go-astisub for subtitle parsing
4. ✅ Evaluate bimg vs govips (CGo install)
5. ⚠️ Benchmark sonic vs stdlib JSON (API responses)

