## Table of Contents

- [HTTP Client with Proxy/VPN Pattern](#http-client-with-proxyvpn-pattern)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# HTTP Client with Proxy/VPN Pattern

<!-- DESIGN: patterns, README, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


**Created**: 2026-02-02
**Status**: ✅ Complete
**Category**: pattern


> > Reusable HTTP client factory with optional proxy/VPN routing for external API calls

Provides HTTP clients with configurable proxy/VPN routing:
- **HTTP/HTTPS Proxy**: Standard proxy with authentication
- **SOCKS5 Proxy**: Tor and commercial VPN support
- **VPN Binding**: Route through specific network interface
- **Tor Support**: Convenience alias for SOCKS5 (127.0.0.1:9050)

**CRITICAL**: Proxy/VPN is **OPTIONAL** - must be explicitly enabled


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete HTTP client pattern with proxy/VPN support |
| Sources | ✅ | All proxy types documented |
| Instructions | ✅ | Implementation examples provided |
| Code | 🔴 | To be implemented |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../architecture/METADATA_SYSTEM.md)
- [METADATA](../services/METADATA.md)
- [TMDB](../integrations/metadata/video/TMDB.md)
- [THETVDB](../integrations/metadata/video/THETVDB.md)

### External Sources
- [Go net/http Package](../../sources/go/stdlib/nethttp.md) - Standard library HTTP client
- [Go net/url Package](https://pkg.go.dev/net/url) - URL parsing for proxy configuration
- [golang.org/x/net/proxy](https://pkg.go.dev/golang.org/x/net/proxy) - SOCKS5 proxy support
- [Tor Project](https://www.torproject.org/) - Tor SOCKS5 proxy documentation
- [HTTP Proxy RFC](https://datatracker.ietf.org/doc/html/rfc7231#section-4.3.6) - HTTP CONNECT method specification
- [SOCKS5 RFC](https://datatracker.ietf.org/doc/html/rfc1928) - SOCKS5 protocol specification

