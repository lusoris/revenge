## Phase 0 - Foundation ✅ COMPLETE

See [docs/PHASE0_FOUNDATION.md](docs/PHASE0_FOUNDATION.md)

## Phase 1 - Core MVP (Current)

See [docs/PHASE1_TODO.md](docs/PHASE1_TODO.md) for detailed tasks.

### Next Steps
1. Database migrations (users, sessions, oidc_providers)
2. User & Auth (JWT + OIDC)
3. Libraries & Scanner
4. Media Items & Playback

---

## Go 1.25 Migration Tasks (August 2025)

### New Features to Adopt

- [ ] **sync.WaitGroup.Go** - Replace manual `wg.Add(1); go func() { defer wg.Done(); ... }()` patterns
- [ ] **testing/synctest** - Use for concurrent code testing (time virtualization)
- [ ] **net/http.CrossOriginProtection** - Replace custom CSRF middleware with stdlib
- [ ] **slog.GroupAttrs** - Use for cleaner grouped logging
- [ ] **runtime/trace.FlightRecorder** - Add to debugging/observability tooling
- [ ] **os.Root extended methods** - Use Chmod, Chown, MkdirAll, ReadFile, Symlink, etc.
- [ ] **hash.Cloner** - Use for cloning hash states where needed
- [ ] **io/fs.ReadLinkFS** - Use for symlink support in filesystem operations
- [ ] **reflect.TypeAssert** - Use for zero-allocation type assertions in hot paths

### Experimental Features (Evaluate)

- [ ] **GOEXPERIMENT=greenteagc** - Test new garbage collector (10-40% GC reduction)
- [ ] **GOEXPERIMENT=jsonv2** - Test new JSON implementation (faster decoding)

### Deprecations to Address

- [ ] **go/ast.FilterPackage, PackageExports, MergePackageFiles** - Deprecated (if used)
- [ ] **go/parser.ParseDir** - Deprecated (if used)

### Breaking Changes to Check

- [ ] **Nil pointer bug fix** - Audit code for patterns like `f, err := os.Open(); name := f.Name(); if err != nil {...}`
- [ ] **SHA-1 in TLS 1.2** - Now disallowed by default (check if any legacy systems need `GODEBUG=tlssha1=1`)
- [ ] **testing.AllocsPerRun** - Now panics if parallel tests running

### Removed Dependencies

- [x] **automaxprocs** - No longer needed! Go 1.25 has container-aware GOMAXPROCS built-in
  - Delete any `go.uber.org/automaxprocs` imports
  - Runtime now respects cgroup CPU limits automatically

### Performance Improvements (Automatic)

- ✅ Container-aware GOMAXPROCS (cgroup limits)
- ✅ Dynamic GOMAXPROCS updates
- ✅ Faster slices (more stack allocations)
- ✅ DWARF5 debug info (smaller binaries, faster linking)
- ✅ Faster crypto/rsa key generation (3x)
- ✅ Faster crypto/sha1 on amd64 with SHA-NI (2x)
- ✅ Faster crypto/sha3 on Apple M processors (2x)
- ✅ Parallel runtime.AddCleanup execution
- ✅ Better unique package memory reclamation

### New go vet Analyzers

- [ ] **waitgroup** - Check for misplaced `sync.WaitGroup.Add` calls
- [ ] **hostport** - Check for IPv6-incompatible address formatting (use `net.JoinHostPort`)

### New go.mod Features

- [ ] **`ignore` directive** - Can be used to ignore directories (e.g., C# upstream code)
  ```
  ignore Jellyfin.*
  ignore MediaBrowser.*
  ignore Emby.*
  ```

### New Tooling

- [ ] **`go doc -http`** - Starts local documentation server
- [ ] **`go version -m -json`** - JSON output of BuildInfo

