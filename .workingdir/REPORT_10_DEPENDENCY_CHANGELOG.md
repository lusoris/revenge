# Dependency Changelog Report

**Date**: 2026-02-06
**Source**: Dependabot PR #26 - 9 dependency updates
**Purpose**: Analyze changes between current and proposed versions for (a) required code changes and (b) new features worth adopting.

---

## Summary

| Package | Current | Target | Risk | Action Required |
|---------|---------|--------|------|-----------------|
| `riverqueue/river` (+driver, +type) | 0.26.0 | 0.30.2 | Medium | No breaking changes, but big feature jump - review new APIs |
| `go.opentelemetry.io/otel` (+metric, +trace) | 1.39.0 | 1.40.0 | Low-Medium | Check Prometheus exporter error handling if used |
| `hashicorp/go-hclog` | 1.6.2 | 1.6.3 | None | Drop-in patch |
| `lmittmann/tint` | 1.1.2 | 1.1.3 | None | Drop-in patch |
| `golang.org/x/time` | 0.12.0 | 0.14.0 | None | Drop-in, maintenance only |

---

## Detailed Analysis

### 1. riverqueue/river 0.26.0 -> 0.30.2

**Breaking Changes**: None.

**Our Usage** (extensive):
- 5-tier priority queue system (critical/high/default/low/bulk)
- ~20 worker types across media processing, metadata, library scanning, notifications
- Unique jobs, batch insert, periodic jobs, leader election
- Custom middleware (logging, metrics, error handling)

**New Features Worth Adopting**:

#### a) Stuck Job Detection (v0.27.0)
Worker-level `Timeout()` method that marks jobs as stuck if they exceed the timeout, then retries them. Currently we rely only on River's global job timeout.

```go
// Workers can now declare their own timeout:
func (w *MediaScanWorker) Timeout(job *river.Job[MediaScanArgs]) time.Duration {
    return 10 * time.Minute
}
```

**Recommendation**: Add `Timeout()` to long-running workers (media transcoding, library scan, metadata fetch). Prevents zombie jobs from blocking queues.

#### b) Client.JobUpdate (v0.28.0)
Persist job state changes mid-execution (progress tracking, partial results).

```go
// Update job metadata during execution:
job.Args.Progress = 75
client.JobUpdate(ctx, job)
```

**Recommendation**: Use for media transcoding progress and library scan progress reporting. Currently we have no mid-job progress visibility.

#### c) Periodic Job Removal by ID (v0.29.0)
Remove specific periodic jobs without restarting the client.

**Recommendation**: Low priority. Useful if we add user-configurable scheduled tasks later.

#### d) HookPeriodicJobsStart (v0.29.0)
Hook that fires when periodic jobs are enqueued.

**Recommendation**: Could use for metrics/logging of scheduled job execution. Low priority.

#### e) River Pro Features (informational)
Batch jobs, workflow orchestration, and rate limiting are Pro-only. Not applicable to us (open source).

**Migration Effort**: Drop-in upgrade. No code changes required. New features are opt-in.

---

### 2. go.opentelemetry.io/otel 1.39.0 -> 1.40.0

**Breaking Changes**:

#### a) Prometheus Exporter Error Handling (v1.39.0)
If we use the Prometheus exporter, it now emits `NewInvalidMetric` on data loss, causing HTTP 500 on scrapes. Fix:
```go
promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}
```

**Action**: Check if our Prometheus setup uses default handler options. If so, add `ContinueOnError`.

#### b) Go 1.23 Support Dropped
Not relevant - we're on Go 1.25.

#### c) Metric Attribute Rename
`rpc.grpc.status_code` -> `rpc.response.status_code`

**Action**: Check if we use gRPC-specific metric attributes. If not, no action needed.

**New Features Worth Adopting**:

#### a) Enabled() Method for Instruments (v1.40.0)
Check if an instrument is active before doing expensive work:

```go
counter, _ := meter.Int64Counter("expensive_metric")
if counter.Enabled(ctx) {
    value := expensiveComputation()
    counter.Add(ctx, value)
}
```

**Recommendation**: Use for any metrics that require expensive computation (e.g., cache hit ratios, queue depth calculations). Medium priority.

#### b) AlwaysRecord Sampler (v1.40.0)
Records all spans regardless of parent sampling decision. Useful for metrics-from-traces.

**Recommendation**: Low priority unless we need 100% trace recording for specific paths.

#### c) Zipkin Exporter Deprecated
We don't use Zipkin, so no impact.

**Migration Effort**: Mostly drop-in. Check Prometheus exporter config if applicable.

---

### 3. hashicorp/go-hclog 1.6.2 -> 1.6.3

**Changes**: Adds optional JSON escaping for log output fields.
**Breaking Changes**: None.
**Action**: Drop-in upgrade. We use hclog as a transitive dependency for HashiCorp Raft.
**Features to Adopt**: None relevant.

---

### 4. lmittmann/tint 1.1.2 -> 1.1.3

**Changes**: Fixes color reset bug in terminal output (ANSI escape sequences not properly reset in some edge cases).
**Breaking Changes**: None.
**Action**: Drop-in upgrade. We use tint for colored slog output.
**Features to Adopt**: None - bugfix only.

---

### 5. golang.org/x/time 0.12.0 -> 0.14.0

**Changes**:
- Internal fix: uses `time.Time.Equal()` instead of `==` for time comparison (correctness improvement)
- Go directive bumped to 1.24
**Breaking Changes**: None.
**Action**: Drop-in upgrade. We use this for rate limiting.
**Features to Adopt**: None - maintenance only.

---

## Recommendations

### Must Do (before merging PR #26)
1. **Check Prometheus exporter config** - ensure error handling is set to `ContinueOnError` if we use `promhttp`
2. **Run full test suite** after upgrade to catch any subtle behavior changes

### Should Do (after merge, new tasks)
1. **Add `Timeout()` to long-running River workers** - prevents stuck/zombie jobs
   - `MediaTranscodeWorker`: 30min timeout
   - `LibraryScanWorker`: 15min timeout
   - `MetadataFetchWorker`: 5min timeout
   - Other workers: 2-5min as appropriate
2. **Implement `JobUpdate` for progress tracking** on media transcoding and library scan jobs
3. **Use OTel `Enabled()` check** before expensive metric computations

### Nice to Have (low priority)
1. Periodic job removal by ID (useful for future user-configurable schedules)
2. `HookPeriodicJobsStart` for scheduled job metrics
3. `AlwaysRecord` sampler for specific trace paths

---

## Merge Decision

**Safe to merge PR #26**: Yes, with the Prometheus exporter check.

All upgrades are backwards-compatible. The River jump (4 minor versions) looks large but has no breaking changes. The OTel changes are mostly additive with one behavioral change in Prometheus error handling.
