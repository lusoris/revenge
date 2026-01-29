# Video Module Documentation

> Movies and TV shows content management

---

## Modules

| Module | Location | Description |
|--------|----------|-------------|
| [Movie Module](MOVIE_MODULE.md) | `internal/content/movie/` | Movie library management |
| [TV Show Module](TVSHOW_MODULE.md) | `internal/content/tvshow/` | Series, seasons, episodes |

---

## Architecture

Both modules follow the same structure:

```
internal/content/{module}/
├── entity.go              # Domain entities
├── repository.go          # Repository interface
├── repository_pg.go       # PostgreSQL implementation
├── repository_pg_user_data.go   # User ratings, favorites
├── repository_pg_relations.go   # Cast, crew, genres
├── service.go             # Business logic + caching
├── jobs.go                # River background jobs
├── metadata_provider.go   # TMDb interface
└── module.go              # fx DI module
```

---

## Metadata Sources

| Content | Primary | Fallback |
|---------|---------|----------|
| Movies | Radarr | TMDb |
| TV Shows | Sonarr | TMDb, TheTVDB |

**Servarr-First Principle**: Use Servarr (Radarr/Sonarr) as primary metadata source. External APIs (TMDb) are only for enrichment via background jobs.

---

## Related

- [Playback Features](../playback/) - Playback and streaming
- [Library Service](../../services/LIBRARY.md) - Library management
- [Metadata Service](../../services/METADATA.md) - Providers
