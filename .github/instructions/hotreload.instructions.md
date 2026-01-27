---
applyTo: "**/pkg/hotreload/**/*.go,**/pkg/config/**/*.go"
---

# Hot Reload & Feature Flags Instructions

## Overview

`pkg/hotreload` provides:

- Configuration file watching and reloading
- Atomic configuration swapping
- Feature flags for gradual rollouts

## Configuration Hot Reload

### When to Use

- Non-critical settings that can change at runtime
- Rate limits, feature toggles, logging levels
- Avoiding restarts for config changes

### Basic Pattern

```go
// Implement ReloadableConfig interface
type Config struct {
    path string
    data atomic.Pointer[ConfigData]
}

func (c *Config) Load() error {
    data, err := loadFromFile(c.path)
    if err != nil {
        return err
    }
    c.data.Store(data)
    return nil
}

func (c *Config) Validate() error {
    data := c.data.Load()
    // Validate data...
    return nil
}

// Create watcher
watcher := hotreload.NewConfigWatcher(
    hotreload.WatcherConfig{
        Files:        []string{"config.yaml", "config.local.yaml"},
        PollInterval: 5 * time.Second,
        Debounce:     time.Second,
        OnReload: func(err error) {
            if err != nil {
                logger.Error("config reload failed", "error", err)
            } else {
                logger.Info("config reloaded")
                // Notify dependent services
            }
        },
    },
    config,
    logger,
)

watcher.Start(ctx)
defer watcher.Stop()
```

### Atomic Value for Thread-Safe Access

```go
type RuntimeConfig struct {
    LogLevel     string
    RateLimit    int
    FeatureFlags map[string]bool
}

var runtimeConfig = hotreload.NewAtomicValue(RuntimeConfig{
    LogLevel:  "info",
    RateLimit: 100,
})

// Read (lock-free)
cfg := runtimeConfig.Load()
if cfg.RateLimit > 0 {
    // apply rate limit
}

// Update (from reload callback)
runtimeConfig.Store(newConfig)
```

### What Can Be Hot-Reloaded

| ✅ Safe to Reload | ❌ Requires Restart |
| ----------------- | ------------------- |
| Log level         | Database URL        |
| Rate limits       | Listen port         |
| Feature flags     | TLS certificates    |
| Timeouts          | Worker count        |
| Cache sizes       | Pool sizes          |

## Feature Flags

### When to Use

- Gradual feature rollouts
- A/B testing
- Kill switches for new features
- User-specific features

### Basic Pattern

```go
flags := hotreload.NewFeatureFlags()

// Set flags (from config or API)
flags.Set(hotreload.FeatureFlagConfig{
    Name:        "new-player",
    Enabled:     true,
    Percentage:  10,  // 10% rollout
    Description: "New video player UI",
})

flags.Set(hotreload.FeatureFlagConfig{
    Name:    "experimental-codec",
    Enabled: false, // Kill switch
})
```

### Check Feature Enabled

```go
// Simple check (global)
if flags.IsEnabled("new-player") {
    useNewPlayer()
}

// Per-user check (consistent per user)
if flags.IsEnabledForUser("new-player", userID) {
    useNewPlayer()
}
```

### Percentage Rollout

```go
flags.Set(hotreload.FeatureFlagConfig{
    Name:       "new-transcoder",
    Enabled:    true,
    Percentage: 25,  // 25% of users
})

// Same user always gets same result (consistent hashing)
user1Enabled := flags.IsEnabledForUser("new-transcoder", "user-1")
user1Enabled2 := flags.IsEnabledForUser("new-transcoder", "user-1") // Same!
```

### Flag Configuration from YAML

```yaml
feature_flags:
  - name: new-player
    enabled: true
    percentage: 100

  - name: experimental-transcoding
    enabled: true
    percentage: 10

  - name: beta-ui
    enabled: false
```

```go
func (c *Config) Load() error {
    // Load config file...

    // Update feature flags
    for _, flag := range c.data.FeatureFlags {
        flags.Set(flag)
    }

    return nil
}
```

## Directory Watching

Watch for new/changed files in a directory:

```go
watcher := hotreload.NewDirWatcher(
    "/etc/revenge/plugins",
    "*.yaml",
    10 * time.Second,
    func(path string) {
        logger.Info("plugin config changed", "path", path)
        reloadPlugin(path)
    },
    logger,
)

watcher.Start(ctx)
```

## Integration with fx

```go
func NewConfigWatcher(
    lc fx.Lifecycle,
    config *Config,
    logger *slog.Logger,
) *hotreload.ConfigWatcher {
    watcher := hotreload.NewConfigWatcher(
        hotreload.DefaultWatcherConfig(config.Path()),
        config,
        logger,
    )

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            watcher.Start(ctx)
            return nil
        },
        OnStop: func(ctx context.Context) error {
            watcher.Stop()
            return nil
        },
    })

    return watcher
}
```

## DO's

- ✅ Debounce rapid file changes
- ✅ Validate config before applying
- ✅ Log all config changes
- ✅ Use atomic values for thread safety
- ✅ Keep feature flag names consistent

## DON'Ts

- ❌ Hot-reload database connections
- ❌ Change pool sizes at runtime
- ❌ Forget validation in Load()
- ❌ Use feature flags for permanent config
- ❌ Change TLS config at runtime

## Testing Feature Flags

```go
func TestFeatureFlag(t *testing.T) {
    flags := hotreload.NewFeatureFlags()
    flags.Set(hotreload.FeatureFlagConfig{
        Name:    "test-feature",
        Enabled: true,
    })

    if !flags.IsEnabled("test-feature") {
        t.Error("expected feature to be enabled")
    }
}

func TestPercentageRollout(t *testing.T) {
    flags := hotreload.NewFeatureFlags()
    flags.Set(hotreload.FeatureFlagConfig{
        Name:       "gradual-feature",
        Enabled:    true,
        Percentage: 50,
    })

    // Same user should get consistent result
    user := "test-user-123"
    result1 := flags.IsEnabledForUser("gradual-feature", user)
    result2 := flags.IsEnabledForUser("gradual-feature", user)

    if result1 != result2 {
        t.Error("expected consistent result for same user")
    }
}
```
