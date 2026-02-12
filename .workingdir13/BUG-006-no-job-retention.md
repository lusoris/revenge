# BUG-006: No CompletedJobRetention — river_job table grows unbounded
**Status: RESOLVED**

**Severity:** MEDIUM
**Category:** River Jobs
**Status:** RESOLVED

## Description

The River client in `internal/infra/jobs/river.go` is created without setting `CompletedJobRetention` or `DiscardedJobRetention` in the `river.Config`. This means completed and discarded jobs accumulate in the `river_job` table indefinitely, consuming disk space.

## File

`internal/infra/jobs/river.go` — `NewClient()` function

## Fix

Add retention settings to `river.Config`:
```go
CompletedJobRetention: 24 * time.Hour,
DiscardedJobRetention: 7 * 24 * time.Hour,
```
