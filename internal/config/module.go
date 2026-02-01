package config

import (
	"go.uber.org/fx"
)

// Module provides configuration dependencies.
var Module = fx.Module("config",
	fx.Provide(ProvideConfig),
)

// ProvideConfig loads and provides the application configuration.
func ProvideConfig() (*Config, error) {
	// Load configuration from default path
	// Environment variables will override file values
	cfg, err := Load(DefaultConfigPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Default returns a Config with default values.
// This is useful for testing or when no config file exists.
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			ReadTimeout:     30000000000,  // 30s
			WriteTimeout:    30000000000,  // 30s
			IdleTimeout:     120000000000, // 120s
			ShutdownTimeout: 10000000000,  // 10s
		},
		Database: DatabaseConfig{
			URL:               "",
			MaxConns:          0, // Auto: (CPU * 2) + 1
			MinConns:          2,
			MaxConnLifetime:   1800000000000, // 30m
			MaxConnIdleTime:   300000000000,  // 5m
			HealthCheckPeriod: 30000000000,   // 30s
		},
		Cache: CacheConfig{
			URL:     "",
			Enabled: false,
		},
		Search: SearchConfig{
			URL:     "",
			APIKey:  "",
			Enabled: false,
		},
		Jobs: JobsConfig{
			MaxWorkers:           100,
			FetchCooldown:        200000000,   // 200ms
			FetchPollInterval:    2000000000,  // 2s
			RescueStuckJobsAfter: 1800000000000, // 30m
		},
		Logging: LoggingConfig{
			Level:       "info",
			Format:      "text",
			Development: false,
		},
		Auth: AuthConfig{
			JWTSecret:     "",
			JWTExpiry:     86400000000000,  // 24h
			RefreshExpiry: 604800000000000, // 7 days
		},
		Legacy: LegacyConfig{
			Enabled:       false,
			EncryptionKey: "",
			Privacy: LegacyPrivacyConfig{
				RequirePIN:      true,
				AuditAllAccess: true,
			},
		},
	}
}
