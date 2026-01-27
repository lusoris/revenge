---
applyTo: "**/pkg/config/**/*.go,**/configs/*.yaml"
---

# koanf v2 Configuration Guide

> Modern configuration management for Go

## Installation

```bash
go get github.com/knadh/koanf/v2
go get github.com/knadh/koanf/providers/file
go get github.com/knadh/koanf/providers/env/v2
go get github.com/knadh/koanf/parsers/yaml
```

## Basic Usage

```go
import (
    "github.com/knadh/koanf/v2"
    "github.com/knadh/koanf/providers/file"
    "github.com/knadh/koanf/providers/env/v2"
    "github.com/knadh/koanf/parsers/yaml"
)

// Global instance with "." delimiter
var k = koanf.New(".")

func LoadConfig() error {
    // Load YAML config file
    if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
        return err
    }

    // Merge environment variables (overrides file)
    return k.Load(env.Provider(".", env.Opt{
        Prefix: "JELLYFIN_",
        TransformFunc: func(key, val string) (string, any) {
            // JELLYFIN_DATABASE_HOST → database.host
            key = strings.ToLower(strings.TrimPrefix(key, "JELLYFIN_"))
            key = strings.ReplaceAll(key, "_", ".")
            return key, val
        },
    }), nil)
}
```

## Config Struct Pattern

```go
// Define config structure
type Config struct {
    Server   ServerConfig   `koanf:"server"`
    Database DatabaseConfig `koanf:"database"`
    Cache    CacheConfig    `koanf:"cache"`
    Logging  LoggingConfig  `koanf:"logging"`
}

type ServerConfig struct {
    Host         string        `koanf:"host"`
    Port         int           `koanf:"port"`
    ReadTimeout  time.Duration `koanf:"read_timeout"`
    WriteTimeout time.Duration `koanf:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `koanf:"host"`
    Port     int    `koanf:"port"`
    Name     string `koanf:"name"`
    User     string `koanf:"user"`
    Password string `koanf:"password"`
    SSLMode  string `koanf:"ssl_mode"`
    MaxConns int    `koanf:"max_conns"`
}

// Unmarshal to struct
func GetConfig() (*Config, error) {
    var cfg Config
    if err := k.Unmarshal("", &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

## YAML Config File

```yaml
# config.yaml
server:
  host: "0.0.0.0"
  port: 8096
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 5432
  name: jellyfin
  user: jellyfin
  password: "" # Override with env
  ssl_mode: disable
  max_conns: 25

cache:
  enabled: true
  host: localhost
  port: 6379
  db: 0

logging:
  level: info
  format: json
  output: stdout
```

## Environment Override

```bash
# Override database password via environment
export JELLYFIN_DATABASE_PASSWORD=secretpassword
export JELLYFIN_SERVER_PORT=9000
export JELLYFIN_LOGGING_LEVEL=debug
```

## Multiple Config Sources

```go
func LoadConfig() error {
    // 1. Load defaults first
    if err := k.Load(file.Provider("defaults.yaml"), yaml.Parser()); err != nil {
        return err
    }

    // 2. Load main config (overrides defaults)
    if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
        // Config file is optional
        if !os.IsNotExist(err) {
            return err
        }
    }

    // 3. Load environment-specific config
    env := os.Getenv("JELLYFIN_ENV")
    if env != "" {
        envFile := fmt.Sprintf("config.%s.yaml", env)
        _ = k.Load(file.Provider(envFile), yaml.Parser())
    }

    // 4. Environment variables override everything
    return k.Load(env.Provider(".", env.Opt{
        Prefix: "JELLYFIN_",
        TransformFunc: envKeyTransform,
    }), nil)
}
```

## Getting Values

```go
// Type-safe getters
host := k.String("server.host")        // string
port := k.Int("server.port")           // int
timeout := k.Duration("server.timeout") // time.Duration
enabled := k.Bool("cache.enabled")      // bool

// With defaults
port := k.Int("server.port")
if port == 0 {
    port = 8096
}

// Check existence
if k.Exists("database.password") {
    // ...
}

// Get nested map
dbConfig := k.StringMap("database") // map[string]string
```

## Strict Merge (Type Safety)

```go
var k = koanf.NewWithConf(koanf.Conf{
    Delim:       ".",
    StrictMerge: true, // Error if types don't match
})
```

## Watch for Changes

```go
f := file.Provider("config.yaml")
k.Load(f, yaml.Parser())

// Watch for changes
f.Watch(func(event any, err error) {
    if err != nil {
        slog.Error("config watch error", "error", err)
        return
    }

    slog.Info("config changed, reloading...")

    // Reload config (thread-safe)
    newK := koanf.New(".")
    if err := newK.Load(f, yaml.Parser()); err != nil {
        slog.Error("failed to reload config", "error", err)
        return
    }

    // Atomic swap
    k = newK
})
```

## Provider: confmap (Defaults)

```go
import "github.com/knadh/koanf/providers/confmap"

// Load defaults from map
k.Load(confmap.Provider(map[string]any{
    "server.host": "0.0.0.0",
    "server.port": 8096,
    "database.max_conns": 25,
}, "."), nil)
```

## Provider: structs

```go
import "github.com/knadh/koanf/providers/structs"

// Load defaults from struct
defaults := Config{
    Server: ServerConfig{
        Host: "0.0.0.0",
        Port: 8096,
    },
}
k.Load(structs.Provider(defaults, "koanf"), nil)
```

## Provider: Command Line Flags

```go
import (
    "github.com/knadh/koanf/providers/posflag"
    flag "github.com/spf13/pflag"
)

f := flag.NewFlagSet("config", flag.ContinueOnError)
f.String("config", "config.yaml", "config file path")
f.Int("port", 8096, "server port")
f.Parse(os.Args[1:])

// Load flags (pass k to handle defaults correctly)
k.Load(posflag.Provider(f, ".", k), nil)
```

## Custom Unmarshal

```go
import "github.com/go-viper/mapstructure/v2"

var cfg Config
k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{
    Tag: "koanf",
    DecoderConfig: &mapstructure.DecoderConfig{
        WeaklyTypedInput: true,
        DecodeHook: mapstructure.ComposeDecodeHookFunc(
            mapstructure.StringToTimeDurationHookFunc(),
            mapstructure.StringToSliceHookFunc(","),
        ),
    },
})
```

## fx Integration

```go
func NewConfig() (*Config, error) {
    k := koanf.New(".")

    // Load config...

    var cfg Config
    if err := k.Unmarshal("", &cfg); err != nil {
        return nil, fmt.Errorf("unmarshal config: %w", err)
    }

    return &cfg, nil
}

// In fx
fx.Provide(NewConfig)
```

## Jellyfin Go Pattern

```go
// pkg/config/config.go
package config

type Config struct {
    Server   Server   `koanf:"server"`
    Database Database `koanf:"database"`
    Cache    Cache    `koanf:"cache"`
    Logging  Logging  `koanf:"logging"`
    Media    Media    `koanf:"media"`
}

func Load(paths ...string) (*Config, error) {
    k := koanf.New(".")

    // Defaults
    k.Load(confmap.Provider(defaults, "."), nil)

    // Config files
    for _, p := range paths {
        if err := k.Load(file.Provider(p), yaml.Parser()); err != nil {
            if !os.IsNotExist(err) {
                return nil, err
            }
        }
    }

    // Environment
    k.Load(env.Provider(".", env.Opt{
        Prefix: "JELLYFIN_",
        TransformFunc: transformEnvKey,
    }), nil)

    var cfg Config
    if err := k.Unmarshal("", &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```
