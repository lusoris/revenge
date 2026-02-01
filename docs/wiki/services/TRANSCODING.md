## Table of Contents

- [Transcoding Service](#transcoding-service)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Transcoding Service




> High-performance media transcoding with hardware acceleration

The Transcoding Service converts media files on-demand to ensure compatibility across all devices. **INTERNAL transcoding** uses go-astiav (FFmpeg Go bindings) with optional hardware acceleration (NVENC, QSV, VAAPI) and is always available. For heavy workloads, users can optionally configure **EXTERNAL offloading** to a Blackbeard service (third-party, not developed by us). Generates HLS adaptive streams with multiple quality levels, caching results for faster subsequent playback.

---





---


## How It Works

<!-- How it works -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [go-astiav (FFmpeg)](../../sources/media/go-astiav.md)
- [gohlslib](../../sources/media/gohlslib.md)
- [River Job Queue](../../sources/tooling/river.md)
- [Uber fx](../../sources/tooling/fx.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)