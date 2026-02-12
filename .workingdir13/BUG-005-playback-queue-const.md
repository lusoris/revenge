# BUG-005: playback_cleanup uses hardcoded "low" instead of QueueLow constant
**Status: RESOLVED**

**Severity:** LOW
**Category:** River Jobs
**Status:** RESOLVED

## Description

In `internal/playback/jobs/cleanup.go`, the `InsertOpts()` for `CleanupArgs` uses the string literal `"low"` instead of the `infrajobs.QueueLow` constant:

```go
func (CleanupArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue: "low",  // should be infrajobs.QueueLow
        ...
    }
}
```

## File

`internal/playback/jobs/cleanup.go` — InsertOpts method

## Fix

Import `infrajobs` and use `infrajobs.QueueLow`.
