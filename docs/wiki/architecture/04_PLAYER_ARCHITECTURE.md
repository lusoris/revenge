## Table of Contents

- [Revenge - Player Architecture](#revenge-player-architecture)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Revenge - Player Architecture




> How Revenge plays your media across all devices

The player architecture handles media playback across web, mobile, and TV clients. Videos stream via HLS for adaptive quality based on your connection. The backend generates HLS manifests and handles transcoding when needed (using FFmpeg). The web player (Vidstack) supports chapters, subtitles, skip intro, and trickplay thumbnails. Chromecast lets you cast to TV, and SyncPlay lets multiple users watch together in sync.

---





---


## Features
<!-- Feature list placeholder -->
## Configuration
## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [FFmpeg Documentation](../../sources/media/ffmpeg.md)
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md)
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md)
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md)
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md)
- [gohlslib (HLS)](../../sources/media/gohlslib.md)
- [M3U8 Extended Format](../../sources/protocols/m3u8.md)
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md)
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md)
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)