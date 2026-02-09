# Session 5 Status

## Completed Work

### Profiling, Benchmarking & Metrics Instrumentation (commit 16e2f739)

**Phase 1: Wire Dead Metrics** (8 files)
- Auth attempts: login/register/verify_email with success/failure in `auth/service.go`
- Active sessions gauge: inc/dec on create/revoke in `session/service.go`
- Rate limit hits: per-operation blocked tracking in `middleware/ratelimit.go`
- Job enqueue counters: per job type in `jobs/river.go`
- Search query counters + duration histograms in `search/movie_service.go`
- Library scan duration, files scanned, errors in `moviejobs/library_scan.go`
- DB query duration by operation type + error counters in `database/logger.go`
- Cache size gauge (L1 item count) in `cache/cache.go`

**Phase 2: Periodic Stats Collector** (`observability/collector.go`)
- pgxpool gauges: acquired, idle, total, max, constructing connections
- River queue depth by state (available, running, retryable, scheduled, etc.)
- Runs every 15s via fx lifecycle hook

**Phase 3: Always-On pprof**
- Removed dev-mode gate from `observability/server.go`
- Observability port is internal-only, safe to expose always

**Phase 4: Makefile Targets**
- `make bench` / `bench-cpu` / `bench-mem` for benchmarks
- `make pprof-cpu` / `pprof-heap` / `pprof-goroutine` for live profiling

**Phase 5: Benchmark Tests**
- `cache/cache_bench_test.go`: L1 Get hit ~78ns, miss ~31ns, set ~442ns (zero allocs on reads)
- `observability/metrics_bench_test.go`: all counter/histogram ops ~80-110ns, zero allocs

**Phase 6: Prometheus + Grafana**
- Added to `docker-compose.dev.yml` under `tools` profile
- Prometheus: scrapes revenge:9096/metrics every 15s, 7d retention
- Grafana: pre-provisioned dashboard with 20 panels across 7 categories
- Config in `deploy/prometheus/` and `deploy/grafana/provisioning/`

## Remaining Uncommitted Changes (from previous sessions)

These files are modified but not part of the observability commit:
- `internal/api/handler.go` — GetUserById error handling fix
- `internal/api/handler_metadata.go` — ErrNoProviders/ErrNotFound handling, nil imageService guard
- `internal/api/server.go` — server changes
- `internal/app/module.go` — module additions (playback, notification)
- `internal/infra/image/module.go` — image module changes
- `internal/integration/radarr/*` — radarr slog migration + module additions
- `internal/integration/sonarr/*` — sonarr slog migration + module additions
- `tests/live/smoke_test.go` — live end-to-end smoke tests

## Benchmark Results (i9-14900KS)

```
BenchmarkL1Cache_Set-32              2,479,852    442.0 ns/op    131 B/op    3 allocs/op
BenchmarkL1Cache_Get_Hit-32         16,769,898     77.8 ns/op      0 B/op    0 allocs/op
BenchmarkL1Cache_Get_Miss-32        36,624,476     30.9 ns/op      0 B/op    0 allocs/op
BenchmarkL1Cache_Concurrent-32       8,469,584    133.9 ns/op     46 B/op    2 allocs/op
BenchmarkHTTPRequestsTotal_Inc-32   12,855,046    108.9 ns/op      0 B/op    0 allocs/op
BenchmarkCacheHit_Record-32         16,327,257     80.3 ns/op      0 B/op    0 allocs/op
BenchmarkDBQueryDuration_Observe-32 11,711,965     99.6 ns/op      0 B/op    0 allocs/op
```

## Docker Status

Prometheus + Grafana containers downloading (first pull). Access:
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)
- pprof: http://localhost:9096/debug/pprof/
