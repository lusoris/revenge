# Logging Infrastructure

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/logging`
**fx Module**: `logging.Module`

> Dual logger setup: slog (standard) + zap (structured), with tint for development

---

## Service Structure

```
internal/infra/logging/
├── logging.go             # Logger factories (slog, zap, test)
└── module.go              # fx module (ProvideSlogLogger, ProvideZapLogger)
```

## Logger Factories

```go
func NewLogger(cfg Config) *slog.Logger       // slog (sets as slog.SetDefault)
func NewZapLogger(cfg Config) *zap.Logger     // zap (sets as zap.ReplaceGlobals)
func NewTestLogger() *slog.Logger             // Discards output, debug level
```

Both loggers are provided via fx to all services. Most services use one or the other:

| Logger | Used By |
|--------|---------|
| `slog` | Database, cache, search infra, notification, health |
| `zap` | Auth, user, RBAC, library, settings, activity, MFA, OIDC, API keys, image |

## Modes

### Development Mode

Uses **tint** handler (colorized, human-readable output):
- Source location included (`AddSource: true`)
- Colored log levels
- Human-readable timestamps
- Caller information

### Production Mode

Uses **JSON handler** (structured, machine-parseable):
- ISO8601 timestamps
- No source location (performance)
- Stack traces for error level

## Configuration

From `config.go` `LoggingConfig` (koanf namespace `logging.*`):
```yaml
logging:
  level: info              # debug, info, warn, error
  format: text             # text or json (auto: text in dev, json in prod)
  development: true        # Enables tint handler + pprof
```

## Dependencies

- `log/slog` - Standard library structured logging
- `go.uber.org/zap` - High-performance structured logging
- `github.com/lmittmann/tint` - Colorized slog handler for development

## Related Documentation

- [OBSERVABILITY.md](OBSERVABILITY.md) - Metrics and pprof (enabled via `development` flag)
- [DATABASE.md](DATABASE.md) - Query logger uses slog
