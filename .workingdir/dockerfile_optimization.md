# Dockerfile Optimization - From 175MB to ~10MB

**Current**: 175MB
**Target**: ~10-15MB (94% reduction!)

---

## Issues Found

### Issue 1: FFmpeg (100-150MB)
**Why included**: Transcoding
**Reality**: Design docs show Blackbeard (external transcoder) is planned
**Solution**: Remove FFmpeg, use Blackbeard when ready

### Issue 2: Alpine Base (7MB)
**Why used**: Common practice
**Reality**: Go binaries don't need Alpine!
**Solution**: Use `gcr.io/distroless/static:nonroot` or `scratch`

### Issue 3: CGO Enabled
**Current**: `CGO_ENABLED=1`
**Check needed**: Do we actually use CGO?
**If no**: Disable CGO for truly static binary

---

## Optimized Dockerfile

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go files
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build STATIC binary (no CGO if not needed)
ENV GOEXPERIMENT=greenteagc,jsonv2
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -w -s -extldflags '-static'" \
    -a \
    -o revenge \
    ./cmd/revenge

# Runtime stage - DISTROLESS!
FROM gcr.io/distroless/static:nonroot

# Copy only the binary
COPY --from=builder /build/revenge /app/revenge

# Distroless already has:
# - CA certificates
# - /etc/passwd with nonroot user
# - Timezone data

WORKDIR /app

# Expose port
EXPOSE 8096

# No shell, no package manager, no FFmpeg
# Just the binary!

# Health check (if distroless supports it, otherwise remove)
# HEALTHCHECK needs wget/curl which distroless doesn't have
# Use k8s liveness probe instead

USER nonroot:nonroot

ENTRYPOINT ["/app/revenge"]
```

---

## Size Comparison

| Base Image | Size | Has Shell | Has Package Manager | CA Certs | User Management |
|------------|------|-----------|---------------------|----------|-----------------|
| **alpine:latest** | 7 MB | ✅ sh | ✅ apk | ✅ | ✅ |
| **distroless/static** | 2 MB | ❌ | ❌ | ✅ | ✅ (nonroot) |
| **scratch** | 0 MB | ❌ | ❌ | ❌ | ❌ |

**Recommendation**: Use `distroless/static:nonroot` for CA certs + nonroot user

---

## New Size Breakdown

| Component | Size |
|-----------|------|
| Distroless base | ~2 MB |
| Go binary (stripped, static) | ~8-10 MB |
| **TOTAL** | **~10-12 MB** |

**Reduction**: 175MB → 10MB = **165MB saved (94%)**

---

## What About FFmpeg?

### Option A: Wait for Blackbeard (Recommended)
- Design docs show Blackbeard integration planned
- Remove FFmpeg from Dockerfile now
- Add Blackbeard container later

### Option B: Separate Transcoder Container
```yaml
services:
  revenge:
    image: revenge:latest  # 10MB
  transcoder:
    image: lscr.io/linuxserver/ffmpeg:latest
    # or custom image with our transcoding logic
```

---

## Questions to Answer

1. **Do we use CGO?** Check for `import "C"` in codebase
2. **Is pgx OK without CGO?** (Yes, pgx works fine without CGO)
3. **Health check?** Distroless has no shell, need k8s probes or HTTP check from outside

---

## Implementation Plan

1. [ ] Verify no CGO usage in codebase
2. [ ] Test binary builds with `CGO_ENABLED=0`
3. [ ] Update Dockerfile to use distroless
4. [ ] Remove FFmpeg
5. [ ] Update docker-compose examples
6. [ ] Document Blackbeard integration path
7. [ ] Update health check strategy (k8s probes)

---

**Estimated time to implement**: 30 minutes
**Estimated testing time**: 1 hour
**Total**: ~1.5 hours for 94% size reduction
