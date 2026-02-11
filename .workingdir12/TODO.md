# Workdir 12 — Playback & Streaming Integration Tests

## Summary

Created comprehensive real-video integration tests for the playback/streaming
subsystem using the Big Buck Bunny 4K test video with actual FFmpeg processing.

## Changes

### New: `internal/playback/integration_test.go`

10 integration test functions (package `playback_test`) using real FFmpeg and
BBB 4K video (`.workingdir11/bbb_sunflower_2160p_30fps_normal.mp4`):

| # | Test | What it covers |
|---|------|----------------|
| 1 | `TestIntegration_ProbeRealVideo` | Probes real file, asserts H.264/4K/2 audio tracks |
| 2 | `TestIntegration_TranscodeDecision` | AnalyzeMedia decides copy vs transcode correctly |
| 3 | `TestIntegration_SessionLifecycle` | Full start→segments→master.m3u8→media.m3u8→seg.ts→audio→stop |
| 4 | `TestIntegration_SeekStart` | Session at 300s, verifies segments from seek point |
| 5 | `TestIntegration_MaxConcurrentSessions` | max=1 enforcement, slot freeing after stop |
| 6 | `TestIntegration_SessionToResponse` | SessionToResponse converter with real probed data |
| 7 | `TestIntegration_PathTraversalSecurity` | 6 CWE-22 path traversal attack vectors blocked |
| 8 | `TestIntegration_MultipleAudioRenditions` | Both MP3+AC3 tracks get separate HLS renditions |
| 9 | `TestIntegration_PipelineCleanup` | FFmpeg processes cleaned up on session stop |
| 10 | `TestIntegration_SessionTouch` | Touch extends ExpiresAt and LastAccessedAt |

All tests use `if testing.Short() { t.Skip() }` — skipped during `make test-short`.

### Fix: `internal/playback/transcode/pipeline_integration_test.go`

Added missing `"os/exec"` import (pre-existing build error: `newSleepCmd()` and
`newTrueCmd()` helpers called `exec.Command()` without importing `os/exec`).

## Validation

- `make test-short` — all packages pass (0 issues)
- `make lint` — 0 issues
- `go test -count=1 -timeout 120s ./internal/playback/` — all 10 integration tests pass
- `go test -count=1 ./internal/playback/transcode/` — passes after import fix

## Test Video

- File: `bbb_sunflower_2160p_30fps_normal.mp4` (BBB 4K)
- Codec: H.264, 3840×2160, 30fps, 634.6s
- Audio: Track 0 MP3 stereo 160kbps, Track 1 AC-3 5.1 320kbps
- Size: ~633 MB
