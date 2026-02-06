# Library Management

<!-- DESIGN: features/shared -->

**Package**: `internal/service/library`

> Media library creation, scanning, and management

---

Library management is documented in the service layer:

- **[../../services/LIBRARY.md](../../services/LIBRARY.md)** - Library service (CRUD, scanning, permissions, cleanup)
- **[SHARED_CONTENT.md](SHARED_CONTENT.md)** - Shared content packages (scanner, matcher, library interfaces)
- **[../video/MOVIE_MODULE.md](../video/MOVIE_MODULE.md)** - Movie LibraryService (scan, match, probe)
- **[../video/TVSHOW_MODULE.md](../video/TVSHOW_MODULE.md)** - TV show library scan workers

## Supported Library Types

| Type | Content Module | Status |
|------|---------------|--------|
| `movie` | `internal/content/movie` | Implemented |
| `tvshow` | `internal/content/tvshow` | Implemented |
| `adult` | `internal/content/qar` | Placeholder (v0.3.0) |
| `music` | - | Planned |
| `photo` | - | Planned |
| `book` | - | Planned |
| `audiobook` | - | Planned |
| `comic` | - | Planned |
| `podcast` | - | Planned |
