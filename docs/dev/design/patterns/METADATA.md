# Metadata Enrichment Pattern

> How content gets its metadata: two tiers, priority-based providers, adapters, and caching. Written from code as of 2026-02-06.

---

## Two-Tier Model

Metadata flows through two tiers:

| Tier | Source | Role | Code |
|------|--------|------|------|
| **Tier 1: PRIMARY** | Arr integrations (Radarr, Sonarr) | Source of truth for library state + base metadata | `internal/integration/{radarr,sonarr}/` |
| **Tier 2: SUPPLEMENTARY** | External APIs (TMDb, TVDb) | Enrichment: translations, credits, images, ratings, related content | `internal/service/metadata/` |

**The arr is always authoritative when configured and present.** TMDb/TVDb are fallback + enrichment — they fill gaps (credits, images, translations) and provide data that arr services don't carry.

The priority depends on the content type and what each source can deliver. For movies, TMDb (priority 100) is the richest enrichment source. For TV shows, both TMDb and TVDb (priority 80) contribute, with TMDb preferred.

For the full architecture, provider interfaces, caching details, and fx wiring, see **[Metadata System (Architecture)](../architecture/METADATA_SYSTEM.md)**.

---

## Data Flow

```
Radarr/Sonarr sync
    ↓ base metadata (title, year, IDs, files, quality)
    ↓ stored in content module tables
    ↓
Content Service calls MetadataProvider.Enrich()
    ↓
Adapter (movie or tvshow) translates types
    ↓
metadata.Service dispatches to providers by priority
    ↓
TMDb API (pri 100) → fallback → TVDb API (pri 80)
    ↓
Enriched data: translations, credits, images, ratings
    ↓
Merged into content domain types, stored in DB
```

---

## Key Patterns

### 1. Content Modules Never Call APIs Directly

Content modules (`movie.Service`, `tvshow.Service`) receive a `MetadataProvider` interface via fx. They don't know about TMDb, TVDb, or any specific provider:

```go
type MetadataProvider interface {
    SearchMovies(ctx context.Context, query string, year int, languages []string) ([]Movie, error)
    EnrichMovie(ctx context.Context, movie *Movie) error
    GetMovieCredits(ctx context.Context, tmdbID int) ([]MovieCredit, error)
    GetMovieImages(ctx context.Context, tmdbID int) ([]Image, error)
    // ...
}
```

### 2. Adapters Bridge Types

Adapters in `internal/service/metadata/adapters/{movie,tvshow}/` convert between the shared metadata types and content module domain types:

- `float64` → `decimal.Decimal` (vote average, popularity)
- Release dates → age ratings map (US → MPAA → PG-13)
- Translations → i18n maps (`TitlesI18n`, `TaglinesI18n`, `OverviewsI18n`)

### 3. Provider Priority and Fallback

Providers are sorted by `Priority()` (highest first). If `EnableProviderFallback` is true, the service tries the next provider on failure. This is transparent to callers.

### 4. Multi-Language Fetching

When multiple languages are requested, the first language becomes the base result. Subsequent languages are merged as translations into a `Translations` map. Default: `["en"]`.

### 5. Async Refresh via River

Refresh operations (`RefreshMovie`, `RefreshTVShow`) enqueue River jobs rather than fetching synchronously. Related workers: `MetadataRefreshMovie`, `MetadataRefreshTVShow`, `SeriesRefresh`.

---

## Caching Layers

| Layer | Technology | Scope | TTL |
|-------|-----------|-------|-----|
| L0 | `sync.Map` per provider | HTTP response caching in TMDb/TVDb clients | 24h metadata, 15m search |
| L1 | otter (W-TinyLFU) | In-process, bounded | 5–10 min |
| L2 | rueidis → Dragonfly | Shared across instances | Per-key (see [Cache Strategy](CACHE_STRATEGY.md)) |

**Known issue**: L0 `sync.Map` caching has unbounded memory growth and should migrate to otter. Tracked in `.workingdir3/CODEBASE_TODOS.md` item #17.

---

## Adding a New Metadata Provider

1. Implement the `Provider` base interface + capability interfaces (`MovieProvider`, `TVShowProvider`, etc.)
2. Create API client in `internal/service/metadata/providers/{name}/`
3. Add config section to `config.go`
4. Register in the metadata fx module (conditional on API key being configured)
5. The service auto-sorts by priority and dispatches — no other wiring needed

Reserved provider IDs (not yet implemented): `fanarttv`, `omdb`.

---

## Related Documentation

- **[Metadata System (Architecture)](../architecture/METADATA_SYSTEM.md)** — full provider interfaces, fx wiring, error types, 27-method service interface
- [Radarr Integration](../integrations/servarr/RADARR.md) — Tier 1 provider for movies
- [Sonarr Integration](../integrations/servarr/SONARR.md) — Tier 1 provider for TV shows
- [Cache Strategy](CACHE_STRATEGY.md) — L1/L2 caching infrastructure
- [River Workers](RIVER_WORKERS.md) — Async metadata refresh jobs
