# TODO — Working Directory 10

## HIGH Priority

- [x] **OIDC nonce validation** — Add nonce to auth request, verify in ID token via `oidc.Config{Nonce: ...}` → `internal/service/oidc/service.go`
- [x] **TOTP replay protection** — Track last used code/timestamp, reject reuse within validity window → `internal/service/mfa/totp.go:132`
- [x] **Migrate dirty state** — Check dirty flag, handle with `m.Force(version)` or return error → `internal/infra/database/migrate.go:34,68`
- [x] **Raft cluster join** — Implement `AddVoter` endpoint/mechanism for multi-node → `internal/infra/raft/election.go`
- [x] **River double leader** — Remove Raft leader check from River workers (River has built-in leader election) → `internal/infra/jobs/cleanup_job.go:97-106`, `service/activity/cleanup.go:42`, `service/library/cleanup.go:43`, `service/analytics/stats_worker.go:51`
- [x] **rueidis DoCache()** — Switch `client.Do()` → `client.DoCache()` for GET reads in cache → `internal/infra/cache/cache.go`

## MEDIUM Priority

### rueidis
- [x] **NewLuaScript for rate limiter** — ~~Done (commit 3589e5fc) — rueidis.NewLuaScript for EVALSHA~~
- [x] **IsRedisNil check** — ~~Done as part of DoCache() fix (commit ab2c5435)~~

### otter
- [ ] **Close() on provider caches** — Add `Close()` method + fx lifecycle hook to all 14 metadata/integration clients → all provider `client.go` files
- [ ] **Pipeline cache TTL bug** — `ttl=0` falls through to 5min default, fix `NewL1Cache` to allow zero TTL → `internal/playback/transcode/pipeline.go`, `internal/infra/cache/otter.go`
- [ ] **OnDeletion for transcode cache** — Add deletion listener to kill evicted FFmpeg processes → `internal/playback/transcode/pipeline.go`

### Raft
- [x] **Structured logging for transport/snapshots** — ~~Done (commit faac9188) — slogWriter for transport & snapshots~~
- [x] **hcLogAdapter args** — ~~Done (commit faac9188) — forward key-value args~~
- [x] **Close BoltDB stores** — ~~Done (commit faac9188) — store refs, close in Close()~~

### River
- [x] **UniqueOpts on notification jobs** — ~~Done (commit 97864681) — ByArgs + ByPeriod: 1h~~

### Casbin
- [x] **SavePolicy atomicity** — ~~Already uses tx.Begin → DELETE → INSERT → tx.Commit (no fix needed)~~
- [x] **RemovePolicy empty fields** — ~~Empty fields as wildcards is standard Casbin adapter semantics (by design)~~

### OIDC
- [x] **Cache OIDC provider/discovery** — ~~Done (commit 97864681) — sync.Map cache~~
- [x] **Fix fallback endpoints** — ~~Done (commit 97864681) — use oidcProvider.Endpoint()~~

### TOTP
- [x] **Re-enrollment guard** — ~~Done (commit 5bb597cb) — refuse GenerateSecret when TOTP verified~~

### WebAuthn
- [x] **Counter=0 handling** — ~~Done (commit 5bb597cb) — skip when both counters are 0~~

### S3
- [x] **SDK error types** — ~~Done (commit de5108a2) — errors.As with types.NoSuchKey/types.NotFound~~
- [ ] **Multipart upload** — Use `s3.UploadManager` for large files → `internal/service/storage/s3.go:89-95`

### req/v3
- [x] **Fix ad-hoc clients** — ~~Done (commit cf52c6fc) — dedicated imgClient in both files~~

### govips
- [x] **vips.Shutdown()** — ~~Done (commit de5108a2) — added to server OnStop lifecycle hook~~

### Typesense
- [x] **Bulk index error propagation** — ~~Done (commit cf52c6fc) — return error with count~~

## LOW Priority

### rueidis
- [ ] **rueidisotel integration** — Add built-in OTel metrics for `do_cache_miss/hits` → `internal/infra/cache/module.go`
- [ ] **Fix misleading comment** — `DisableAutoPipelining: false` comment says "disable" → `internal/infra/cache/module.go:107`

### otter
- [ ] **ExpiryAccessing for metadata caches** — Use `ExpiryAccessing` instead of `ExpiryWriting` for read-heavy provider caches
- [ ] **Typed generics** — Replace `L1Cache[string, any]` with concrete types where feasible

### Raft
- [ ] **Merge BoltDB stores** — Use single BoltStore for both log and stable → `internal/infra/raft/election.go:87-96`

### River
- [ ] **Cleanup retry count** — Override `MaxAttempts` for cleanup jobs (currently uses global 25) → `internal/infra/jobs/cleanup_job.go`

### Casbin
- [ ] **Deny policies** — Consider adding `eft` field for explicit deny support → `config/casbin_model.conf`

### ogen
- [ ] **AdminListUsers error type** — 400 validation errors returned as 403 type → `internal/api/handler_admin_users.go:54-61` (may need OpenAPI spec fix)

### req/v3
- [ ] **Retry condition filter** — Add `SetCommonRetryCondition` to skip retries on 4xx → all 17 clients

### govips
- [ ] **StartupConfig tuning** — Set explicit `ConcurrencyLevel` and `MaxCacheSize` → `internal/api/image_utils.go:19`

### S3
- [x] **Remove custom contains** — ~~Done (commit de5108a2) — removed contains/containsInner, replaced with errors.As~~

### Typesense
- [ ] **URL parsing** — Use `net/url.Parse` instead of manual parsing → `internal/infra/search/module.go:44-69`
- [ ] **Client timeout** — Add `WithConnectionTimeout` to Typesense client → `internal/infra/search/module.go:78`

### Prometheus
- [x] **Fix dragonflyCommandsProcessed** — ~~Gauge with .Set() is correct for externally scraped monotonic values (Counter only supports .Inc/.Add)~~
- [ ] **Remove duplicate queue metric** — Deduplicate `JobsQueueSize` vs `riverQueueSize` → `metrics.go:144`, `collector.go:19`

### Argon2id
- [ ] **Param consistency** — Align test parallelism (p=4) with production (p=2) → `internal/service/user/service_unit_test.go:70`

## DONE ✅

- [x] **Dragonfly cluster emulation** — Added `--cluster_mode=emulated` to containers.go, docker-compose.yml, docker-compose.dev.yml
- [x] **rueidis DoCache()** — commit ab2c5435
- [x] **River double leader** — commit 28cd21e9 (removed from 4 workers + auth module + 5 test files)\n- [x] **Migrate dirty state** — commit c98d2ed0 (MigrateUp/Down/To + 2 integration tests)\n- [x] **OIDC nonce validation** — commit ab43d225 (migration 38 + sqlc + service nonce generate/verify)
- [x] **TOTP replay protection** — commit 60ce003b (migration 39 + last_used_code + VerifyCode check)
- [x] **Raft cluster join** — commit eb2d0465 (AddVoter/RemoveServer/GetClusterMembers + 14 tests)
- [x] **S3 SDK error types + vips.Shutdown** — commit de5108a2 (errors.As for S3, vips.Shutdown in OnStop)
- [x] **IsRedisNil** — part of DoCache() commit ab2c5435
- [x] **Casbin SavePolicy** — already uses transactions (no fix needed)
- [x] **Casbin RemovePolicy** — empty fields as wildcards is by design
- [x] **Prometheus counter** — Gauge with .Set() is correct for external values
- [x] **S3 remove contains** — removed in commit de5108a2
- [x] **NewLuaScript rate limiter** — commit 3589e5fc
- [x] **Raft BoltDB/logging/args** — commit faac9188
- [x] **UniqueOpts notifications** — commit 97864681
- [x] **OIDC caching + endpoints** — commit 97864681
- [x] **MFA re-enrollment + counter=0** — commit 5bb597cb
- [x] **Ad-hoc clients + bulk errors** — commit cf52c6fc
