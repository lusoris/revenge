# BUG-004: TriggerLibraryScan only enqueues movie scan jobs
**Status: RESOLVED**

**Severity:** HIGH
**Category:** Flow Integration
**Status:** RESOLVED

## Description

In `internal/api/handler_library.go`, `TriggerLibraryScan` always enqueues `moviejobs.MovieLibraryScanArgs` regardless of the library type. TV show libraries can't be scanned via the API.

```go
_, insertErr := h.riverClient.Insert(ctx, moviejobs.MovieLibraryScanArgs{
    ScanID:    scan.ID.String(),
    LibraryID: params.LibraryId.String(),
    Paths:     lib.Paths,
    Force:     scanType == "full",
}, nil)
```

The TV show scan worker (`tvshowjobs.LibraryScanWorker`) exists and is registered, but there's no code path to enqueue `tvshowjobs.LibraryScanArgs`.

## File

`internal/api/handler_library.go` — TriggerLibraryScan method (~line 334)

## Fix

Check `lib.Type` and enqueue the appropriate job type:
- `library.TypeMovie` → `moviejobs.MovieLibraryScanArgs`
- `library.TypeTVShow` → `tvshowjobs.LibraryScanArgs`
