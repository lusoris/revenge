---
applyTo: "**/internal/service/library/**/*.go,**/internal/content/**/*.go"
---

# File Watching - fsnotify

> File system event monitoring for library scanning

## Overview

Use `fsnotify` to watch library directories for changes (new files, modifications, deletions). This enables automatic library updates without manual scanning.

**Package**: `github.com/fsnotify/fsnotify`

## Installation

```bash
go get github.com/fsnotify/fsnotify
```

## Basic Usage

### Create Watcher

```go
import "github.com/fsnotify/fsnotify"

watcher, err := fsnotify.NewWatcher()
if err != nil {
    return err
}
defer watcher.Close()

// Add directory to watch
err = watcher.Add("/media/movies")
if err != nil {
    return err
}
```

### Event Loop

```go
func watchLoop(ctx context.Context, watcher *fsnotify.Watcher) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case event, ok := <-watcher.Events:
            if !ok {
                return nil
            }

            switch {
            case event.Op&fsnotify.Create == fsnotify.Create:
                log.Info("file created", "path", event.Name)
                // Queue scan job

            case event.Op&fsnotify.Write == fsnotify.Write:
                log.Info("file modified", "path", event.Name)
                // Queue metadata refresh

            case event.Op&fsnotify.Remove == fsnotify.Remove:
                log.Info("file removed", "path", event.Name)
                // Mark as missing in database

            case event.Op&fsnotify.Rename == fsnotify.Rename:
                log.Info("file renamed", "path", event.Name)
                // Update path in database
            }

        case err, ok := <-watcher.Errors:
            if !ok {
                return nil
            }
            log.Error("watcher error", "error", err)
        }
    }
}
```

## Library Watcher Service

```go
package library

import (
    "context"
    "log/slog"
    "path/filepath"
    "sync"
    "time"

    "github.com/fsnotify/fsnotify"
    "github.com/riverqueue/river"
)

type WatcherService struct {
    watcher   *fsnotify.Watcher
    jobClient *river.Client[pgx.Tx]
    logger    *slog.Logger

    // Debouncing
    pending   map[string]time.Time
    pendingMu sync.Mutex
}

func NewWatcherService(jobClient *river.Client[pgx.Tx], logger *slog.Logger) (*WatcherService, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }

    return &WatcherService{
        watcher:   watcher,
        jobClient: jobClient,
        logger:    logger.With(slog.String("service", "library-watcher")),
        pending:   make(map[string]time.Time),
    }, nil
}

func (s *WatcherService) WatchLibrary(ctx context.Context, libraryID uuid.UUID, path string) error {
    // Watch directory recursively
    err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return s.watcher.Add(p)
        }
        return nil
    })
    if err != nil {
        return fmt.Errorf("failed to walk directory: %w", err)
    }

    s.logger.Info("watching library", "library_id", libraryID, "path", path)
    return nil
}

func (s *WatcherService) Run(ctx context.Context) error {
    // Debounce timer - batch events within 5 seconds
    debounceInterval := 5 * time.Second
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case event, ok := <-s.watcher.Events:
            if !ok {
                return nil
            }
            s.handleEvent(ctx, event)

        case err, ok := <-s.watcher.Errors:
            if !ok {
                return nil
            }
            s.logger.Error("watcher error", "error", err)

        case <-ticker.C:
            // Process debounced events
            s.processPending(ctx, debounceInterval)
        }
    }
}

func (s *WatcherService) handleEvent(ctx context.Context, event fsnotify.Event) {
    // Skip temporary files
    if isTemporaryFile(event.Name) {
        return
    }

    // Only process media files
    if !isMediaFile(event.Name) {
        // But watch new directories
        if event.Op&fsnotify.Create == fsnotify.Create {
            if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
                s.watcher.Add(event.Name)
                s.logger.Debug("watching new directory", "path", event.Name)
            }
        }
        return
    }

    // Debounce - record event time
    s.pendingMu.Lock()
    s.pending[event.Name] = time.Now()
    s.pendingMu.Unlock()

    s.logger.Debug("file event",
        "op", event.Op.String(),
        "path", event.Name,
    )
}

func (s *WatcherService) processPending(ctx context.Context, debounce time.Duration) {
    s.pendingMu.Lock()
    defer s.pendingMu.Unlock()

    now := time.Now()
    for path, eventTime := range s.pending {
        if now.Sub(eventTime) >= debounce {
            // Event has settled, process it
            s.queueScanJob(ctx, path)
            delete(s.pending, path)
        }
    }
}

func (s *WatcherService) queueScanJob(ctx context.Context, path string) {
    _, err := s.jobClient.Insert(ctx, &ScanFileArgs{
        Path: path,
    }, nil)
    if err != nil {
        s.logger.Error("failed to queue scan job", "error", err, "path", path)
    }
}

func (s *WatcherService) Close() error {
    return s.watcher.Close()
}

// Helper functions
func isTemporaryFile(path string) bool {
    base := filepath.Base(path)
    return strings.HasPrefix(base, ".") ||
           strings.HasSuffix(base, ".tmp") ||
           strings.HasSuffix(base, ".part") ||
           strings.Contains(base, "~")
}

func isMediaFile(path string) bool {
    ext := strings.ToLower(filepath.Ext(path))
    mediaExts := map[string]bool{
        ".mkv": true, ".mp4": true, ".avi": true, ".mov": true,
        ".wmv": true, ".m4v": true, ".webm": true,
        ".mp3": true, ".flac": true, ".m4a": true, ".ogg": true,
        ".wav": true, ".aac": true, ".opus": true,
    }
    return mediaExts[ext]
}
```

## River Job for File Scanning

```go
package library

type ScanFileArgs struct {
    Path string `json:"path"`
}

func (ScanFileArgs) Kind() string { return "library.scan_file" }

type ScanFileWorker struct {
    river.WorkerDefaults[ScanFileArgs]
    svc *LibraryService
}

func (w *ScanFileWorker) Work(ctx context.Context, job *river.Job[ScanFileArgs]) error {
    return w.svc.ProcessFile(ctx, job.Args.Path)
}
```

## fx Module Integration

```go
package library

import (
    "context"
    "go.uber.org/fx"
)

var WatcherModule = fx.Module("library-watcher",
    fx.Provide(NewWatcherService),
    fx.Invoke(func(lc fx.Lifecycle, svc *WatcherService) {
        lc.Append(fx.Hook{
            OnStart: func(ctx context.Context) error {
                go svc.Run(context.Background())
                return nil
            },
            OnStop: func(ctx context.Context) error {
                return svc.Close()
            },
        })
    }),
)
```

## Configuration

```yaml
library:
  watcher:
    enabled: true
    debounce_seconds: 5
    exclude_patterns:
      - "*.tmp"
      - "*.part"
      - ".DS_Store"
      - "@eaDir/*" # Synology
      - "#recycle/*"
```

```go
type WatcherConfig struct {
    Enabled          bool          `koanf:"enabled"`
    DebounceSeconds  int           `koanf:"debounce_seconds"`
    ExcludePatterns  []string      `koanf:"exclude_patterns"`
}

var DefaultWatcherConfig = WatcherConfig{
    Enabled:         true,
    DebounceSeconds: 5,
    ExcludePatterns: []string{"*.tmp", "*.part", ".DS_Store"},
}
```

## DO's and DON'Ts

### DO

- ✅ Debounce events (files often trigger multiple events)
- ✅ Watch directories recursively
- ✅ Add new directories when created
- ✅ Filter temporary/partial files
- ✅ Use River jobs for actual processing
- ✅ Close watcher on shutdown

### DON'T

- ❌ Process events immediately (causes duplicate work)
- ❌ Watch individual files (watch directories)
- ❌ Block the event loop with slow operations
- ❌ Ignore watcher errors
- ❌ Watch network paths without testing performance
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
