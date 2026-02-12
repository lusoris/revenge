# BUG-002: movie_file_match retries 25x on deterministic failure
**Status: RESOLVED**

**Severity:** HIGH  
**Category:** River Jobs  
**Status:** RESOLVED  

## Description

In `internal/content/movie/moviejobs/file_match.go`, when a file can't be matched to a movie, the worker returns:

```go
return errors.New("file could not be matched to any movie")
```

This is a deterministic failure—the file simply doesn't match anything. Retrying 25 times won't change the result. Each retry hits the library service, metadata providers, and potentially TMDb.

## File

`internal/content/movie/moviejobs/file_match.go` — line ~97

## Fix

Return `nil` and log a warning instead of returning an error. The file is simply unmatched, which is a valid outcome, not a job failure.
