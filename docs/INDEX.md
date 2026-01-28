# Revenge Documentation

> Central router for all documentation

## Navigation

| Section | Description | Path |
|---------|-------------|------|
| **Development** | Architecture, features, planning, research | [/docs/dev/](dev/INDEX.md) |
| **API Reference** | OpenAPI specs | [/api/openapi/](../api/openapi/) |
| **Agent Instructions** | Coding agent rules | [AGENTS.md](../AGENTS.md) |
| **Contributing** | Contribution guidelines | [CONTRIBUTING.md](../CONTRIBUTING.md) |

## Quick Links

- [Setup Guide](dev/design/operations/SETUP.md)
- [Development Guide](dev/design/operations/DEVELOPMENT.md)
- [Architecture Overview](dev/design/architecture/ARCHITECTURE_V2.md)
- [Module Implementation Roadmap](dev/design/planning/MODULE_IMPLEMENTATION_TODO.md)

## Documentation Structure

```
/docs/
├── INDEX.md                    # You are here
└── dev/
    ├── INDEX.md                # Developer documentation hub
    ├── MIGRATION_MANIFEST.md   # Migration tracking
    ├── design/                 # 🔒 Project design (protected)
    │   ├── architecture/       # System architecture
    │   ├── features/           # Feature specifications
    │   ├── integrations/       # Integration designs
    │   ├── operations/         # Operations guides
    │   ├── planning/           # Roadmaps & planning
    │   ├── research/           # Research & analysis
    │   └── technical/          # Technical documentation
    └── sources/                # 🔄 External sources (auto-fetched)
        ├── SOURCES.yaml        # Source registry
        └── {categories}/       # Fetched documentation
```

## Source Categories

External documentation auto-fetched from upstream sources:

| Category | Examples |
|----------|----------|
| `go` | Go stdlib, fx, koanf, sqlc |
| `apis` | TMDb, MusicBrainz, Trakt |
| `protocols` | HTTP Range, HLS, DASH |
| `database` | PostgreSQL, sqlc patterns |
| `frontend` | Svelte 5, TanStack Query |
| `tooling` | ogen, river, slog |
| `media` | FFmpeg, codecs, containers |
| `security` | OIDC, OAuth 2.0, PKCE |
| `testing` | Go testing, testify |
| `observability` | Prometheus, OpenTelemetry |
| `infrastructure` | Dragonfly, Typesense |

---

## Related Resources

- [Agent Instructions](../AGENTS.md) - Automated coding agent rules
- [Copilot Instructions](../.github/copilot-instructions.md) - GitHub Copilot rules
- [Instruction Files](../.github/instructions/) - 23 pattern-specific instructions
- [TODO List](../TODO.md) - Project backlog
- [Contributing Guide](../CONTRIBUTING.md) - Contribution guidelines
