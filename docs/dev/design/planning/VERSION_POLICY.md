# Version Policy - Global "Bleeding Edge / Latest Stable"

> ALL components use latest stable versions (state-of-the-art only)

## Policy

**Revenge adopts LATEST STABLE versions of ALL components**: frontend, backend, dependencies, infrastructure.

**NO exceptions**. If a component is outdated or unmaintained, replace it.

---

## Backend (Go)

### Language & Stdlib

| Component | Policy | Rationale |
|-----------|--------|-----------|
| **Go version** | Latest stable (1.25+) | New features (WaitGroup.Go, Loop, container GOMAXPROCS), security patches |
| **Stdlib** | Current Go version | HTTP routing (1.22+), slog (1.21+), testing improvements (1.24+) |

**Current**: Go 1.25
**Update Cadence**: Every Go minor release (~6 months)
**Automated**: Renovate bot creates PR on new Go version

---

### Core Dependencies

| Dependency | Policy | Current | Update Cadence |
|------------|--------|---------|----------------|
| **uber-go/fx** | Latest stable | v1.24+ | Every minor release |
| **knadh/koanf/v2** | Latest stable | v2.x | Every minor release |
| **jackc/pgx/v5** | Latest stable | v5.x | Every minor release |
| **redis/go-redis/v9** | Latest stable | v9.x | Every minor release |
| **ogen-go/ogen** | Latest stable | Latest | Every minor release |
| **riverqueue/river** | Latest stable | Latest | Every minor release |

**CI/CD**: Dependabot OR Renovate bot (weekly PRs for outdated deps)

---

### Media Processing

| Package | Policy | Purpose |
|---------|--------|---------|
| **go-astiav** | Latest stable | FFmpeg bindings (Blackbeard transcoding) |
| **go-astisub** | Latest stable | Subtitle parsing (.srt, .ass, .webvtt) |
| **gortsplib** | Latest stable | RTSP server/client (Live TV) |
| **mp4ff** | Latest stable | MP4 manipulation |
| **bimg** OR **govips** | Latest stable | Image resizing (libvips wrapper) |
| **flac** | Latest stable | FLAC codec (lossless audio) |

**Update Strategy**: Monitor GitHub releases, update every 2-3 months OR when critical bugs fixed.

---

## Frontend (Web)

### Framework & UI

| Component | Policy | Current | Rationale |
|-----------|--------|---------|-----------|
| **SvelteKit** | Latest stable | v2.x | SSR, routing, modern reactivity (runes) |
| **Svelte** | Latest stable | v5.x | Runes, performance, smaller bundles |
| **Tailwind CSS** | Latest stable | v4.x | Latest features, beta → stable when released |
| **shadcn-svelte** | Latest stable | Latest | UI components (Tailwind + Svelte) |
| **TanStack Query** | Latest stable | v5.x | Server state management |

**Update Cadence**:
- **SvelteKit**: Every minor release (monthly)
- **Tailwind CSS**: Every major/minor (quarterly)
- **shadcn-svelte**: Every release (weekly/monthly)

**Automated**: Renovate bot PRs for npm packages

---

### Build Tools

| Tool | Policy | Purpose |
|------|--------|---------|
| **Vite** | Latest stable | Build tool (bundler, HMR) |
| **TypeScript** | Latest stable | Type safety (5.x+) |
| **ESLint** | Latest stable | Linting (9.x flat config) |
| **Prettier** | Latest stable | Code formatting |

**Update Cadence**: Every minor release (monthly for Vite/TS, quarterly for ESLint)

---

## Infrastructure

### Database

| Component | Policy | Current | Rationale |
|-----------|--------|---------|-----------|
| **PostgreSQL** | Latest stable | 18+ | Performance, ACID, new features (partitioning, JIT) |
| **sqlc** | Latest stable | v1.x | Type-safe SQL code generation |
| **golang-migrate** | Latest stable | v4.x | Database migrations |

**PostgreSQL Update Strategy**:
- Follow PostgreSQL major releases (yearly)
- Upgrade within 6 months of new major release
- Test in staging before production

**CI/CD**: Docker image `postgres:18-alpine` → `postgres:latest` (pinned to major version)

---

### Cache & Search

| Component | Policy | Current | Rationale |
|-----------|--------|---------|-----------|
| **Dragonfly** | Latest stable | Latest | Redis-compatible, high-performance |
| **Typesense** | Latest stable | v0.25+ | Fast search, typo-tolerant |

**Update Cadence**: Every minor release (monthly/quarterly)

---

### Job Queue

| Component | Policy | Current | Rationale |
|-----------|--------|---------|-----------|
| **River** | Latest stable | Latest | PostgreSQL-native, modern Go API |

**Update Strategy**: Monitor GitHub releases, update every 1-2 months

---

## Docker & Deployment

### Base Images

| Image | Policy | Rationale |
|-------|--------|-----------|
| **Go build** | `golang:1.25-alpine` | Minimal, latest Go version |
| **PostgreSQL** | `postgres:18-alpine` | Latest stable PostgreSQL |
| **Dragonfly** | `docker.dragonflydb.io/dragonflydb/dragonfly:latest` | Latest Dragonfly |
| **Typesense** | `typesense/typesense:0.25.2` (→ `:latest` after testing) | Latest Typesense |

**Update Strategy**:
- Pin major versions (e.g., `postgres:18`) NOT `:latest` (avoid breaking changes)
- Test new major versions in staging before production
- Renovate bot creates PRs for new versions

---

### Docker Compose

**Policy**: Use latest Compose file version (3.9 → 4.x when stable)

---

## CI/CD Tools

### Automation

| Tool | Policy | Purpose |
|------|--------|---------|
| **GitHub Actions** | Latest workflows | CI/CD pipelines |
| **Renovate** OR **Dependabot** | Enabled | Automated dependency updates |
| **golangci-lint** | Latest stable | Go linting (v1.62+) |
| **prettier** | Latest stable | Code formatting |

**Update Cadence**:
- **GitHub Actions**: Update workflow syntax every 6 months
- **Renovate**: Weekly dependency PRs (auto-merge minor/patch)
- **golangci-lint**: Every release (monthly)

---

## Update Workflow

### Automated (Renovate Bot)

**Configuration** (`renovate.json`):
```json
{
  "extends": ["config:recommended"],
  "schedule": ["every weekend"],
  "automerge": true,
  "automergeType": "pr",
  "packageRules": [
    {
      "matchUpdateTypes": ["minor", "patch"],
      "automerge": true
    },
    {
      "matchUpdateTypes": ["major"],
      "automerge": false,
      "labels": ["major-update", "needs-review"]
    }
  ]
}
```

**Behavior**:
- **Minor/Patch**: Auto-merge after CI passes
- **Major**: Create PR with `needs-review` label (manual review required)
- **Weekly**: Check for updates every weekend (low-traffic time)

---

### Manual (Major Versions)

**Process**:
1. **PR created** by Renovate with `major-update` label
2. **Review changelog** for breaking changes
3. **Test in staging** environment (Docker Compose)
4. **Update code** if breaking changes (fix imports, API changes)
5. **Merge** after CI passes + manual testing

**Example**: PostgreSQL 18 → 19 (yearly major release)
- Review PostgreSQL 19 release notes
- Test migrations in staging
- Update `docker-compose.yml`: `postgres:18` → `postgres:19`
- Merge PR after validation

---

## Exceptions (When to Downgrade)

### Critical Bug

If latest version has critical bug:
1. **Pin previous stable version** (e.g., `postgres:18.1` → `postgres:18.0`)
2. **Document in `CHANGELOG.md`** with bug reference
3. **Monitor upstream** for fix
4. **Upgrade** when fix released

**Example**: Dragonfly v1.x memory leak → downgrade to v0.9.x until patched

---

### Ecosystem Lag

If ecosystem hasn't caught up to latest version:
1. **Use latest stable that ecosystem supports**
2. **Monitor ecosystem compatibility**
3. **Upgrade** when ecosystem ready

**Example**: Svelte 5 runes (2024) → Wait for shadcn-svelte Svelte 5 support before upgrading

**Current Status**: SvelteKit 2 + Svelte 5 runes are stable (2024+) → Use now

---

## Rationale

### Why "Latest Stable"?

1. **Security**: Latest versions have security patches
2. **Performance**: Modern versions are faster (PostgreSQL 18 vs 12, Go 1.25 vs 1.18)
3. **Features**: New capabilities (Go generics, Svelte runes, PostgreSQL partitioning)
4. **Maintenance**: Easier to maintain bleeding edge than upgrade from ancient versions
5. **Community**: Latest versions have active support, documentation, tutorials
6. **Recruiting**: Developers want to work with modern tech (not legacy stacks)

### Why NOT "Stable LTS"?

**LTS = Technical Debt**:
- **PostgreSQL 12 LTS** (EOL 2024) vs **PostgreSQL 18** (2024):
  - Missing: JIT compilation, partitioning improvements, performance gains
  - Risk: Security vulnerabilities (no patches after EOL)
- **Go 1.18** vs **Go 1.25**:
  - Missing: WaitGroup.Go, Loop, improved generics, slog
  - Risk: Missing performance optimizations, security fixes

**Revenge Philosophy**: Pay update cost incrementally (monthly) vs big-bang migration (yearly).

---

## Monitoring

### Dependency Tracking

**Tools**:
- **Renovate Dashboard**: https://app.renovatebot.com/
- **GitHub Dependency Graph**: Alerts for vulnerable dependencies
- **`go list -m -u all`**: Check outdated Go modules

**Alerts**:
- **Critical**: Security vulnerabilities (GitHub Dependabot alerts)
- **Warning**: Outdated major versions (>6 months behind)
- **Info**: Minor/patch updates available

---

### Version Dashboard

**Create `docs/VERSIONS.md`** (auto-generated):

```markdown
# Component Versions

**Last Updated**: 2026-01-28

## Backend
- Go: 1.25.0
- Echo: v4.12.0
- PostgreSQL: 18.1
- Dragonfly: v1.2.0
- River: v0.14.0

## Frontend
- SvelteKit: 2.10.0
- Svelte: 5.1.0
- Tailwind CSS: 4.0.0
- TypeScript: 5.7.0

## Infrastructure
- Docker Compose: 3.9
- PostgreSQL: 18-alpine
- Typesense: 0.25.2
```

**Automated**: CI job updates `VERSIONS.md` on release

---

## Summary

| Aspect | Policy |
|--------|--------|
| **Go** | Latest stable (1.25+) |
| **PostgreSQL** | Latest stable (18+) |
| **SvelteKit** | Latest stable (2.x) |
| **Svelte** | Latest stable (5.x runes) |
| **Tailwind CSS** | Latest stable (4.x) |
| **Dependencies** | Latest stable (automated via Renovate) |
| **Docker Images** | Pin major versions, test new releases in staging |
| **Update Cadence** | Minor/patch auto-merge, major manual review |
| **Exceptions** | Critical bugs, ecosystem lag (document + monitor) |

**Core Principle**: **State-of-the-art only**. If a component is outdated, replace it. No technical debt.

---

## References

- [PostgreSQL Versioning Policy](https://www.postgresql.org/support/versioning/)
- [Go Release Policy](https://go.dev/doc/devel/release)
- [Renovate Docs](https://docs.renovatebot.com/)
- [SvelteKit Versioning](https://kit.svelte.dev/docs/migrating)
- [Semantic Versioning](https://semver.org/)

