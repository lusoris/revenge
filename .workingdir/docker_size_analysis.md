# Docker Image Size Analysis

**Current Size**: 175MB
**Question**: Is this because of develop build or can we optimize?

---

## Answer: It's FFmpeg, Not Debug Symbols!

The Dockerfile DOES strip debug symbols (`-ldflags "-w -s"`), so it's not a develop vs production issue.

**The culprit is FFmpeg** (Dockerfile line 39) which is **~100-150MB** in Alpine!

---

## Current Dockerfile Analysis

### ✅ GOOD: Multi-Stage Build
```dockerfile
FROM golang:1.25-alpine AS builder  # Build stage
# ... build ...
FROM alpine:latest                   # Runtime stage (only 7MB base!)
```

### ✅ GOOD: Strip Debug Symbols
```dockerfile
-ldflags "-w -s"  # Strips debug info and symbol table
```

### ✅ GOOD: CGO Static Build
```dockerfile
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo
```

### ⚠️ PROBLEM: FFmpeg is MASSIVE

**Dockerfile:39:**
```dockerfile
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \        # <-- 100-150MB!
    tzdata \
```

---

## Size Breakdown (Estimated)

| Component | Size | Reason |
|-----------|------|--------|
| Alpine base | ~7 MB | Minimal Linux |
| Go binary (stripped) | ~10 MB | With CGO, stripped |
| ca-certificates | ~1 MB | SSL certs |
| tzdata | ~3 MB | Timezone data |
| **FFmpeg + deps** | **~100-150 MB** | **Video codecs, libraries, filters** |
| **TOTAL** | **~175 MB** | **Mostly FFmpeg!** |

---

## Questions to Answer

1. **Do we need FFmpeg in the main image?**
   - Is transcoding done by the app directly?
   - Or do we offload to Blackbeard (external transcoder)?

2. **What do we use FFmpeg for?**
   - Video/audio transcoding
   - Metadata extraction (thumbnails, duration, codec info)
   - HLS segment generation
   - Audio normalization

---

## Optimization Options

### Option 1: Remove FFmpeg (if using Blackbeard)

**If using Blackbeard for transcoding:**
```dockerfile
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*
```
**New size**: ~20-30 MB (**85% reduction!**)

### Option 2: FFmpeg Sidecar Container

**Separate transcoding service:**
```yaml
services:
  revenge:
    image: revenge:latest  # 20-30MB without FFmpeg
  transcoder:
    image: lscr.io/linuxserver/ffmpeg:latest  # Pre-built FFmpeg
    # or
    image: revenge-transcoder:latest  # Our FFmpeg + tools
```

**Benefits:**
- Scale transcoding independently
- Upgrade FFmpeg without rebuilding main app
- Main app stays small

### Option 3: Minimal FFmpeg Build

**Only include needed codecs:**
```dockerfile
# Custom FFmpeg with only H.264, H.265, AAC, MP3
# Size: ~30-50MB instead of 150MB
```

### Option 4: Use Distroless for App Only

```dockerfile
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /build/revenge /app/revenge
# Size: ~15-20MB (no FFmpeg)
```

---

## Check TECH_STACK and Design Docs

Let me check what the design says about transcoding...
