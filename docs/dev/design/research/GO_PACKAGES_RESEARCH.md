# Go Packages Research - Media Server Ecosystem

> Bleeding-edge Go packages to massively reduce development work
> **Last Updated:** 2026-01-28
> **Go Version:** 1.25+

## Research Methodology

- **Primary Source**: GitHub repository analysis, benchmarks, maintenance status
- **Criteria**: Stars, activity, license (prefer MIT/BSD), Go 1.24+ compatibility
- **Focus**: Packages that save weeks of development or provide 10x+ performance

---

## ✅ APPROVED - Ready for Implementation

### Caching Layer (Replaces go-redis for hot paths)

| Package | Version | Stars | License | Key Benefit |
|---------|---------|-------|---------|-------------|
| **rueidis** | v1.0.71 | 2.8k | Apache-2.0 | **14x faster** than go-redis, auto-pipelining, server-assisted client-side caching |
| **otter** | v2.3.0 | 2.4k | Apache-2.0 | W-TinyLFU local cache, used by Grafana/Centrifugo, 50% less memory than ristretto |
| **sturdyc** | v1.1.5 | 1k | MIT | **90% API call reduction**, batch coalescing, early refreshes, stampede protection |

**Migration Path:**
```go
// Replace: github.com/redis/go-redis/v9
// With:    github.com/redis/rueidis

// rueidis is wire-compatible with Dragonfly
client, _ := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
    ClientTrackingOptions: []string{"PREFIX", "cache:", "BCAST"},
})
```

**Rationale:**
- rueidis: Auto-pipelining batches commands, server-assisted caching for read-heavy workloads
- otter: Local cache for hot data (user sessions, config), reduces Dragonfly round-trips
- sturdyc: Wraps metadata API calls, deduplicates concurrent requests to TMDb/MusicBrainz

---

## 🌐 WebSocket

| Package | Version | Stars | License | Recommendation |
|---------|---------|-------|---------|----------------|
| **coder/websocket** | v1.8.14 | 4.9k | ISC | ✅ **USE** - Gorilla replacement |
| gws | v1.8.9 | 1.7k | Apache-2.0 | ⚠️ Alternative (event-driven) |
| nbio | v1.6.8 | 2.7k | MIT | ❌ Overkill (1M+ connections) |

**Winner: `github.com/coder/websocket`**
- Zero dependencies
- Full context.Context support
- Zero-alloc reads/writes
- Wasm compilation support
- Maintained by Coder (active)
- Passes Autobahn test suite

```go
// Usage
c, err := websocket.Accept(w, r, nil)
defer c.CloseNow()

err = wsjson.Read(ctx, c, &msg)
err = wsjson.Write(ctx, c, response)
```

---

## 📁 File Watching

| Package | Version | Stars | License | Recommendation |
|---------|---------|-------|---------|----------------|
| **fsnotify** | v1.9.0 | 10.5k | BSD-3-Clause | ✅ **USE** - De-facto standard |

**Note:** 314k dependents, Go 1.25 CI, cross-platform (Linux/macOS/Windows/BSD)

---

## 🌐 HTTP Client (for Metadata Providers)

| Package | Version | Stars | License | Use Case |
|---------|---------|-------|---------|----------|
| **resty** | v2.16.0 | 11.5k | MIT | ✅ **USE** - REST APIs (TMDb, MusicBrainz) |
| fasthttp | v1.69.0 | 23.2k | MIT | ❌ Server-only (API incompatible) |

**Winner: `github.com/go-resty/resty/v2`**
- Circuit breaker built-in
- Retry with exponential backoff
- SSE (Server-Sent Events) support
- Request/Response middleware
- 23k dependents

```go
client := resty.New().
    SetRetryCount(3).
    SetRetryWaitTime(500 * time.Millisecond).
    SetCircuitBreaker(resty.NewCircuitBreaker())

resp, err := client.R().
    SetResult(&TMDbMovie{}).
    Get("https://api.themoviedb.org/3/movie/550")
```

---

## 🎬 Video/Streaming (Blackbeard + Live TV)

| Package | Version | Stars | Purpose | Recommendation |
|---------|---------|-------|---------|----------------|
| **go-astiav** | n8.0 | 665 | FFmpeg bindings | ✅ **USE** - Cleaner than gmf |
| **gortsplib** | v5 | 881 | RTSP client/server | ✅ **USE** - Live TV |
| **mp4ff** | v0.50.0 | 588 | MP4/fMP4 parsing | ✅ **USE** - Segment creation |
| **gohlslib** | v2 | 161 | HLS client/muxer | ✅ **USE** - LL-HLS support |

**Codec Support (mp4ff):**
- Video: H.264, H.265/HEVC, AV1, VP9
- Audio: AAC, AC-3, E-AC-3, Opus
- Containers: MP4, fMP4, CMAF, DASH, HLS

---

## 🖼️ Image Processing

| Package | Version | Stars | License | Recommendation |
|---------|---------|-------|---------|----------------|
| **bimg** | v1.1.5 | 3k | MIT | ✅ **USE** - libvips bindings |
| govips | v2.16.0 | 1.5k | MIT | ⚠️ Maintainers recommend vipsgen |
| imaging | v1.6.2 | 5.2k | MIT | ⚠️ Pure Go fallback (slower) |

**Note:** bimg last commit 2 years ago but stable. Consider vipsgen when mature.

---

## 📝 Subtitles

| Package | Version | Stars | Formats | Recommendation |
|---------|---------|-------|---------|----------------|
| **go-astisub** | latest | 681 | SRT, STL, TTML, SSA/ASS, WebVTT, Teletext | ✅ **USE** |
| subtitles | v0.2.4 | 48 | SRT, VTT, SSA | ⚠️ Simpler alternative |

---

## 🎨 Blurhash

| Package | Version | Stars | License | Recommendation |
|---------|---------|-------|---------|----------------|
| **bbrks/go-blurhash** | v1.1.1 | 167 | **MIT** | ✅ **USE** |
| buckket/go-blurhash | latest | 218 | ⚠️ GPL-3.0 | ❌ License conflict |

**Important:** Use bbrks version - MIT license compatible with commercial use.

---

## 🎵 Audio Metadata

| Package | Version | Stars | License | Purpose |
|---------|---------|-------|---------|---------|
| **dhowden/tag** | latest | 638 | BSD-2-Clause | ✅ Read MP3/MP4/OGG/FLAC |
| **bogem/id3v2** | v2.1.4 | 360 | MIT | ✅ Read+Write ID3v2 tags |

**Use both:** tag for multi-format reading, id3v2 for MP3 tag writing.

---

## 📊 Current Stack (KEEP)

| Component | Package | Rationale |
|-----------|---------|-----------|
| HTTP | **stdlib net/http** | Go 1.22+ routing patterns, ogen handlers |
| Database | **pgx/v5 + sqlc** | Type-safe, no runtime overhead |
| Job Queue | **River** | PostgreSQL-native, no extra service |
| Search | **typesense-go/v4** | Fast, typo-tolerant |
| DI | **uber-go/fx** | Production-proven |
| Config | **koanf/v2** | Hot reload, multiple sources |
| API | **ogen** | OpenAPI spec-first code generation |

---

## 📦 Final go.mod Additions

```go
require (
    // Caching (APPROVED)
    github.com/redis/rueidis v1.0.71       // Replaces go-redis for Dragonfly
    github.com/maypok86/otter v2.3.0       // Local W-TinyLFU cache
    github.com/viccon/sturdyc v1.1.5       // API response caching

    // WebSocket
    github.com/coder/websocket v1.8.14     // Watch Party, live updates

    // HTTP Client
    github.com/go-resty/resty/v2           // Metadata provider calls

    // Video/Streaming (Blackbeard integration)
    github.com/asticode/go-astiav          // FFmpeg bindings
    github.com/bluenviron/gortsplib/v5     // RTSP (Live TV)
    github.com/Eyevinn/mp4ff v0.50.0       // MP4/fMP4 parsing
    github.com/bluenviron/gohlslib/v2      // HLS muxer

    // Image Processing
    github.com/h2non/bimg v1.1.5           // libvips (posters, fanart)

    // Subtitles
    github.com/asticode/go-astisub         // Multi-format parsing

    // Blurhash
    github.com/bbrks/go-blurhash v1.1.1    // Placeholder generation

    // Audio Metadata
    github.com/dhowden/tag                 // Multi-format reading
    github.com/bogem/id3v2/v2 v2.1.4       // ID3 read/write

    // File Watching
    github.com/fsnotify/fsnotify v1.9.0    // Library scanning
)
```

---

## ⏱️ Time Savings Estimate

| Package | Saves | vs Alternative |
|---------|-------|----------------|
| rueidis | 2 weeks | Custom pipelining |
| sturdyc | 3 weeks | Custom request coalescing |
| go-astiav | 4 weeks | Direct FFmpeg integration |
| go-astisub | 2 weeks | Manual subtitle parsing |
| gortsplib | 4 weeks | RTSP from scratch |
| bimg | 1 week | Pure Go imaging |
| resty | 1 week | Custom HTTP client |
| **TOTAL** | **~17 weeks** | |

---

## 🚫 Packages to AVOID

| Package | Reason |
|---------|--------|
| **go-redis/v9** | Replace with rueidis (14x faster) |
| **ristretto** | Replace with otter (50% less memory) |
| **gorilla/websocket** | Deprecated, use coder/websocket |
| **Fiber** | fasthttp-based, not stdlib-compatible |
| **Echo/Gin** | Using stdlib + ogen instead |
| **GORM** | Using sqlc (type-safe, no reflection) |
| **buckket/go-blurhash** | GPL-3.0 license |
| **govips** | Maintainers recommend vipsgen |

---

## Next Steps

1. ✅ Research complete
2. 🔄 Update `go.mod` with approved packages
3. 🔄 Create `internal/infra/cache/` with rueidis + otter
4. 🔄 Create `pkg/httpclient/` with resty + sturdyc
5. 🔄 Update Blackbeard integration docs with go-astiav
6. 🔄 Add WebSocket handlers with coder/websocket
