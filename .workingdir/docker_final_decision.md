# Docker Strategy - FINAL DECISION

**Based on design docs analysis:**

## What We Know

1. **Revenge app needs FFmpeg** for:
   - HLS audio streaming (gohlslib + FFmpeg)
   - Audio transcoding (AAC, MP3, FLAC, Opus)
   - Metadata extraction (duration, codecs, resolution)
   - Thumbnail generation
   - Light video processing (go-astiav)

2. **Blackbeard is separate** for:
   - Heavy video transcoding
   - Long-running jobs
   - Offloaded to background workers

## Docker Strategy Options

### Option A: Keep FFmpeg, Optimize Base (RECOMMENDED)
```dockerfile
# Build: golang:1.25-alpine
# Runtime: alpine:latest (not distroless - FFmpeg needs libs)
# Size: ~60-80MB (FFmpeg ~100MB, but alpine has shared libs)
```
**Pros:** Simple, FFmpeg just works
**Cons:** Larger than ideal

### Option B: Distroless + Static FFmpeg
```dockerfile
# Build: golang:1.25-alpine + build FFmpeg static
# Runtime: gcr.io/distroless/base (needs glibc for FFmpeg)
# Size: ~80-100MB (static FFmpeg is bigger)
```
**Pros:** More secure (distroless)
**Cons:** Static FFmpeg is actually LARGER

### Option C: Minimal FFmpeg Build
```dockerfile
# Build custom FFmpeg with only:
# - H.264, H.265 decoders
# - AAC, MP3, FLAC, Opus codecs
# - HLS muxer
# - Thumbnail extraction
# Runtime: alpine:latest
# Size: ~40-50MB (custom FFmpeg ~30MB)
```
**Pros:** Best balance - smaller but has what we need
**Cons:** Custom build maintenance

### Option D: Two-Image Strategy
```dockerfile
# revenge:latest - 10MB, no FFmpeg (API/DB only)
# revenge:media - 80MB, with FFmpeg (for media nodes)
```
**Pros:** Can scale API without media processing
**Cons:** More complex deployment

## RECOMMENDATION

**Start with Option A (Keep FFmpeg + Alpine)**
- Size: ~80MB (acceptable for media server)
- Simple: Just works, no custom builds
- Can optimize later with Option C when needed

**Later (when optimizing):**
- Option C: Custom minimal FFmpeg build → 40-50MB
- Or stick with 80MB if not an issue

## Action Plan

1. **Now:** Keep current Dockerfile as-is (FFmpeg + Alpine)
2. **Add:** Note in Dockerfile why FFmpeg is needed (reference design docs)
3. **Later:** Consider minimal FFmpeg build when we have metrics on actual codec usage

---

**Size Comparison:**
- Current (Alpine + FFmpeg): ~175MB (with dev deps? check actual)
- Optimized (Alpine + FFmpeg, cleaned): ~80MB
- Custom FFmpeg: ~40-50MB
- Distroless no FFmpeg: ~10MB (but missing required features)

**Verdict:** 80MB is fine for a media server. Jellyfin is 400MB+, Plex is 300MB+.
