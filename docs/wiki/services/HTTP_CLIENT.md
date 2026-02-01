## Table of Contents

- [HTTP Client & Proxy](#http-client-proxy)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# HTTP Client & Proxy




> Route external API calls through proxy or VPN for privacy

The HTTP client service provides centralized HTTP client management with proxy and VPN support. External metadata API calls (TMDb, TheTVDB, MusicBrainz, etc.) route through configured proxies or VPN tunnels for privacy and geo-unblocking. Local services (Arr stack, Stash) bypass proxy. Supports HTTP/HTTPS proxies, SOCKS5 (Tor), and VPN interface binding. Per-service routing rules, health checking, and automatic failover.

---





---


## How It Works

<!-- How it works -->
## Features
<!-- Feature list placeholder -->
## Configuration
## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Go context](../../sources/go/stdlib/context.md)
- [Go net/http](../../sources/go/stdlib/nethttp.md)
- [Go net/http.Transport](https://pkg.go.dev/net/http#Transport)
- [golang.org/x/net/proxy](https://pkg.go.dev/golang.org/x/net/proxy)
- [koanf](../../sources/tooling/koanf.md)
- [Uber fx](../../sources/tooling/fx.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)