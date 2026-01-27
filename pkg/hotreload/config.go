// Package hotreload provides configuration hot reloading.
package hotreload

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// ReloadableConfig is a configuration that can be reloaded.
type ReloadableConfig interface {
	// Load loads the configuration.
	Load() error
	// Validate validates the configuration.
	Validate() error
}

// WatcherConfig configures the config watcher.
type WatcherConfig struct {
	// Files to watch for changes.
	Files []string

	// PollInterval is how often to check for changes.
	PollInterval time.Duration

	// OnReload is called when config is reloaded.
	OnReload func(err error)

	// Debounce prevents rapid reloads.
	Debounce time.Duration
}

// DefaultWatcherConfig returns sensible defaults.
func DefaultWatcherConfig(files ...string) WatcherConfig {
	return WatcherConfig{
		Files:        files,
		PollInterval: 5 * time.Second,
		Debounce:     time.Second,
	}
}

// ConfigWatcher watches config files for changes.
type ConfigWatcher struct {
	config     WatcherConfig
	loader     ReloadableConfig
	logger     *slog.Logger
	modTimes   map[string]time.Time
	mu         sync.Mutex
	stopCh     chan struct{}
	reloading  atomic.Bool
}

// NewConfigWatcher creates a new config watcher.
func NewConfigWatcher(cfg WatcherConfig, loader ReloadableConfig, logger *slog.Logger) *ConfigWatcher {
	return &ConfigWatcher{
		config:   cfg,
		loader:   loader,
		logger:   logger.With(slog.String("component", "config-watcher")),
		modTimes: make(map[string]time.Time),
		stopCh:   make(chan struct{}),
	}
}

// Start begins watching for config changes.
func (w *ConfigWatcher) Start(ctx context.Context) {
	// Get initial mod times
	w.updateModTimes()

	go w.watchLoop(ctx)
	w.logger.Info("config watcher started", "files", w.config.Files)
}

// Stop stops watching.
func (w *ConfigWatcher) Stop() {
	close(w.stopCh)
}

// Reload manually triggers a reload.
func (w *ConfigWatcher) Reload() error {
	return w.doReload()
}

func (w *ConfigWatcher) watchLoop(ctx context.Context) {
	ticker := time.NewTicker(w.config.PollInterval)
	defer ticker.Stop()

	var debounceTimer *time.Timer
	var pendingReload atomic.Bool

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			if w.checkForChanges() {
				if w.config.Debounce > 0 {
					// Debounce: wait for changes to settle
					if !pendingReload.Load() {
						pendingReload.Store(true)
						debounceTimer = time.AfterFunc(w.config.Debounce, func() {
							pendingReload.Store(false)
							w.doReload()
						})
					} else if debounceTimer != nil {
						// Reset debounce timer
						debounceTimer.Reset(w.config.Debounce)
					}
				} else {
					w.doReload()
				}
			}
		}
	}
}

func (w *ConfigWatcher) checkForChanges() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	changed := false
	for _, file := range w.config.Files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		modTime := info.ModTime()
		if lastMod, ok := w.modTimes[file]; ok {
			if modTime.After(lastMod) {
				w.logger.Debug("config file changed", "file", file)
				changed = true
			}
		}
		w.modTimes[file] = modTime
	}

	return changed
}

func (w *ConfigWatcher) updateModTimes() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, file := range w.config.Files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		w.modTimes[file] = info.ModTime()
	}
}

func (w *ConfigWatcher) doReload() error {
	if w.reloading.Load() {
		return nil
	}
	w.reloading.Store(true)
	defer w.reloading.Store(false)

	w.logger.Info("reloading configuration")

	if err := w.loader.Load(); err != nil {
		w.logger.Error("failed to load config", "error", err)
		if w.config.OnReload != nil {
			w.config.OnReload(err)
		}
		return err
	}

	if err := w.loader.Validate(); err != nil {
		w.logger.Error("config validation failed", "error", err)
		if w.config.OnReload != nil {
			w.config.OnReload(err)
		}
		return err
	}

	w.updateModTimes()
	w.logger.Info("configuration reloaded successfully")

	if w.config.OnReload != nil {
		w.config.OnReload(nil)
	}

	return nil
}

// AtomicValue provides atomic access to a value (for hot-swappable config).
type AtomicValue[T any] struct {
	value atomic.Pointer[T]
}

// NewAtomicValue creates a new atomic value.
func NewAtomicValue[T any](initial T) *AtomicValue[T] {
	av := &AtomicValue[T]{}
	av.value.Store(&initial)
	return av
}

// Load returns the current value.
func (av *AtomicValue[T]) Load() T {
	return *av.value.Load()
}

// Store sets a new value.
func (av *AtomicValue[T]) Store(val T) {
	av.value.Store(&val)
}

// Swap atomically swaps the value.
func (av *AtomicValue[T]) Swap(val T) T {
	old := av.value.Swap(&val)
	return *old
}

// FeatureFlagConfig represents a feature flag.
type FeatureFlagConfig struct {
	Name        string
	Enabled     bool
	Percentage  int // 0-100 for gradual rollout
	Description string
}

// FeatureFlags manages runtime feature flags.
type FeatureFlags struct {
	mu    sync.RWMutex
	flags map[string]*FeatureFlagConfig
}

// NewFeatureFlags creates a new feature flag manager.
func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		flags: make(map[string]*FeatureFlagConfig),
	}
}

// Set sets a feature flag.
func (ff *FeatureFlags) Set(flag FeatureFlagConfig) {
	ff.mu.Lock()
	defer ff.mu.Unlock()
	ff.flags[flag.Name] = &flag
}

// IsEnabled checks if a feature is enabled.
func (ff *FeatureFlags) IsEnabled(name string) bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	flag, ok := ff.flags[name]
	if !ok {
		return false
	}

	if flag.Percentage > 0 && flag.Percentage < 100 {
		// Simple hash-based percentage check
		hash := 0
		for _, c := range name {
			hash = 31*hash + int(c)
		}
		return hash%100 < flag.Percentage
	}

	return flag.Enabled
}

// IsEnabledForUser checks if feature is enabled for specific user.
func (ff *FeatureFlags) IsEnabledForUser(name string, userID string) bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	flag, ok := ff.flags[name]
	if !ok {
		return false
	}

	if flag.Percentage > 0 && flag.Percentage < 100 {
		// Consistent hash per user
		hash := 0
		for _, c := range userID {
			hash = 31*hash + int(c)
		}
		return (hash%100)+1 <= flag.Percentage
	}

	return flag.Enabled
}

// All returns all feature flags.
func (ff *FeatureFlags) All() []FeatureFlagConfig {
	ff.mu.RLock()
	defer ff.mu.RUnlock()

	flags := make([]FeatureFlagConfig, 0, len(ff.flags))
	for _, f := range ff.flags {
		flags = append(flags, *f)
	}
	return flags
}

// DirWatcher watches a directory for file changes.
type DirWatcher struct {
	dir        string
	pattern    string
	logger     *slog.Logger
	onChange   func(path string)
	files      map[string]time.Time
	mu         sync.Mutex
	stopCh     chan struct{}
	interval   time.Duration
}

// NewDirWatcher creates a directory watcher.
func NewDirWatcher(dir, pattern string, interval time.Duration, onChange func(path string), logger *slog.Logger) *DirWatcher {
	return &DirWatcher{
		dir:      dir,
		pattern:  pattern,
		logger:   logger.With(slog.String("component", "dir-watcher")),
		onChange: onChange,
		files:    make(map[string]time.Time),
		stopCh:   make(chan struct{}),
		interval: interval,
	}
}

// Start begins watching.
func (dw *DirWatcher) Start(ctx context.Context) {
	dw.scan() // Initial scan

	go func() {
		ticker := time.NewTicker(dw.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-dw.stopCh:
				return
			case <-ticker.C:
				dw.scan()
			}
		}
	}()
}

// Stop stops watching.
func (dw *DirWatcher) Stop() {
	close(dw.stopCh)
}

func (dw *DirWatcher) scan() {
	matches, err := filepath.Glob(filepath.Join(dw.dir, dw.pattern))
	if err != nil {
		dw.logger.Error("glob failed", "error", err)
		return
	}

	dw.mu.Lock()
	defer dw.mu.Unlock()

	// Check for new/modified files
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		modTime := info.ModTime()
		if lastMod, ok := dw.files[path]; !ok || modTime.After(lastMod) {
			dw.files[path] = modTime
			if dw.onChange != nil {
				go dw.onChange(path)
			}
		}
	}

	// Check for deleted files
	currentFiles := make(map[string]bool)
	for _, path := range matches {
		currentFiles[path] = true
	}
	for path := range dw.files {
		if !currentFiles[path] {
			delete(dw.files, path)
		}
	}
}
