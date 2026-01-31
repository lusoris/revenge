---
name: check-sources
description: Check external documentation sources and cross-references
argument-hint: "[category|id]"
disable-model-invocation: true
allowed-tools: Bash(*), Read(*), Write(*)
---

# Check Sources

Analyze external documentation sources, check for updates, and verify cross-references.

## Usage

```
/check-sources                     # Check all sources
/check-sources apis                # Check specific category
/check-sources --id tmdb           # Check specific source
/check-sources --stale             # Find sources that may need refresh
```

## Arguments

- `$0`: Category name, `--id <source-id>`, or `--stale`

## Task

### Default: Source overview

1. Read `docs/dev/sources/INDEX.yaml` for fetch status
2. Read `docs/dev/sources/SOURCES.yaml` for source definitions
3. Report:
   - Total sources defined
   - Successfully fetched
   - Failed to fetch
   - Referenced vs unreferenced

### With category: Category details

1. Filter sources by category
2. Show each source with:
   - Name and URL
   - Last fetch date
   - Content hash
   - Design docs that reference it

### With --id: Source details

1. Find the specific source
2. Show:
   - Full configuration
   - Fetch history
   - All design docs referencing it
   - Content preview (if exists)

### With --stale: Find outdated sources

1. Check last fetch dates
2. Identify sources older than 30 days
3. Suggest refresh commands

## Source Categories

Current categories in SOURCES.yaml:
- apis (TMDB, TVDB, etc.)
- database (pgx, sqlc)
- frontend (Svelte, Tanstack)
- tooling (fx, koanf, ogen)
- protocols (HLS, WebRTC)
- And more...

## Refresh Sources

To refresh sources after checking:
```bash
./scripts/source-pipeline.sh --apply
```

To force refresh a specific source:
```bash
python3 scripts/source-pipeline/01-fetch.py --id tmdb --force --apply
```
