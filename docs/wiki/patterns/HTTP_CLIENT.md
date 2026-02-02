## Table of Contents

- [HTTP Client with Proxy/VPN Pattern](#http-client-with-proxyvpn-pattern)
  - [Features](#features)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# HTTP Client with Proxy/VPN Pattern




> HTTP client factory with optional proxy/VPN routing for external API calls


The HTTP Client pattern provides a centralized factory for creating HTTP clients with
optional proxy/VPN routing. This enables external API calls (TMDb, TheTVDB, etc.) to be
routed through privacy-preserving proxies or VPN tunnels when needed.

Key features:
- Service-specific client configuration
- Multiple proxy types (HTTP, SOCKS5, Tor)
- VPN interface binding
- Connection pooling and timeout management
- Optional by default - must be explicitly enabled


---






## Features
<!-- Feature list placeholder -->
## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Go net/http Package](../../sources/go/stdlib/nethttp.md)
- [Go net/url Package](https://pkg.go.dev/net/url)
- [golang.org/x/net/proxy](https://pkg.go.dev/golang.org/x/net/proxy)
- [Tor Project](https://www.torproject.org/)
- [HTTP Proxy RFC](https://datatracker.ietf.org/doc/html/rfc7231#section-4.3.6)
- [SOCKS5 RFC](https://datatracker.ietf.org/doc/html/rfc1928)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)