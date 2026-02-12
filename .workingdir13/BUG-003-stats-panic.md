# BUG-003: stats_aggregation worker panics when queries is nil
**Status: RESOLVED**

**Severity:** HIGH  
**Category:** River Jobs  
**Status:** RESOLVED

## Description

In `internal/service/analytics/stats_worker.go`, `collectStats()` calls:

```go
if w.queries == nil {
    panic("stats_worker: queries is nil")
}
```

This will crash the entire worker goroutine and potentially the River client. Should return an error instead.

## File

`internal/service/analytics/stats_worker.go` — line ~124

## Fix

Replace `panic()` with `return nil, fmt.Errorf("stats_worker: queries is nil")`.
