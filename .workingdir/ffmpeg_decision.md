# FFmpeg Decision: Keep or Remove?

**Current Status**: FFmpeg is in Dockerfile (~150MB) but NOT used by code yet

---

## Reality Check

### Code Analysis
- ✅ No `go-astiav` (FFmpeg bindings) in go.mod
- ✅ No FFmpeg calls in internal/ or cmd/
- ✅ No gohlslib (HLS needs FFmpeg) in go.mod yet
- ⚠️ Tech stack mentions gohlslib for HLS streaming

### Likely Future Needs
1. **Media Probing** - Get duration, resolution, codecs
2. **Thumbnail Generation** - Extract frame for preview
3. **HLS Streaming** - gohlslib will need FFmpeg
4. **Trickplay** - Generate BIF files (thumbnails for seeking)

---

## Options

### Option 1: Remove Now, Add Later (Recommended)
**Pros:**
- 10MB image NOW
- Can add back when actually implementing features
- Forces us to be intentional about dependencies

**Cons:**
- Need to rebuild when adding features

### Option 2: Keep FFmpeg, Use Distroless/Alpine Slim
**Result:** ~50-100MB (FFmpeg) + ~2-10MB (base) = ~60-110MB

**Pros:**
- Ready for features
- One less thing to worry about

**Cons:**
- Still large (3-6x bigger)
- Paying size cost for unused feature

### Option 3: Minimal FFmpeg Build
**Custom FFmpeg with only:**
- H.264/H.265 decoders
- AAC/MP3 decoders
- Thumbnail extraction
- Basic probing

**Result:** ~30-40MB instead of 150MB
**Best of both worlds:** Ready for features, much smaller

### Option 4: FFmpeg Sidecar Pattern
```yaml
services:
  revenge:
    image: revenge:latest  # 10MB, no FFmpeg
  
  media-processor:
    image: revenge/ffmpeg:latest  # 150MB, has FFmpeg
    # Handle transcoding, thumbnails, probing
```

**Pros:**
- Main app stays tiny
- Scale processors independently
- Upgrade FFmpeg independently

**Cons:**
- More complex architecture
- Network overhead for processing

---

## Recommendation

**Phase 1 (Now):** Remove FFmpeg
- Image: 10-12MB
- Add back when implementing HLS/thumbnails

**Phase 2 (When needed):** Minimal FFmpeg Build
- Image: 40-50MB
- Only codecs we actually use

**Phase 3 (Scale):** FFmpeg Sidecar
- Main: 10MB
- Processor: 150MB
- Scale independently

---

## Decision Time

**Question for user:**
Do you want to:
1. Remove FFmpeg now (10MB), add when needed?
2. Keep FFmpeg, switch to distroless base (60-110MB)?
3. Custom minimal FFmpeg build (40-50MB)?
4. Wait and decide when implementing HLS/thumbnails?
