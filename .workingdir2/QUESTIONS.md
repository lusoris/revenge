# Questions

## Q1: ClearCache scope (RESOLVED)
**Question**: The adapters' `ClearCache()` is a no-op. Should we implement actual cache clearing?
**User answer**: "Full cache clear" - add to provider interfaces, implement actual cache invalidation.
**Implementation**:
- Added `ClearCache()` to `metadata.Provider` base interface
- Added `ClearCache()` to `metadata.Service` interface (delegates to all registered providers)
- TMDb and TVDb providers already had `ClearCache()` methods that clear their `sync.Map` caches
- Adapters delegate `ClearCache()` to `a.service.ClearCache()`
- When `Force=true` in `MetadataRefreshOptions`, adapters call `a.service.ClearCache()` before fetching

## Q2: Languages architecture (RESOLVED)
**Question**: How should per-request languages work?
**User discussion**:
- "der user setzt ja nur seine interface language und seine preferred metadata language"
- "der server hält die auswahl ja vor oder nicht? sprich wenn jemand koreanisch will und das ist nicht im server eingestellt -> gibts nicht"
- "das enrichment muss ja beim updaten trotzdem lücken füllen können"

**Decision**: Languages are server-config-driven via `DefaultLanguages` in the adapter. The `MetadataRefreshOptions.Languages` field allows overriding for specific scenarios (e.g., admin force-refresh in a specific language, or adding a newly configured language). If empty, falls back to the adapter's configured languages.

**Implementation**:
- `MetadataRefreshOptions{Force bool, Languages []string}` added to both content modules
- Adapters: `if len(opts[0].Languages) > 0 { languages = opts[0].Languages }`
- Workers pass `Languages` from job args through to service

## Q3: Handler backward compatibility (RESOLVED)
**Question**: `handler.go` calls `RefreshMovieMetadata(ctx, movieID)` with no opts - is this OK?
**Answer**: Yes, using variadic `opts ...MetadataRefreshOptions` means zero-arg calls still work. Handler uses default behavior (no force, default languages).
