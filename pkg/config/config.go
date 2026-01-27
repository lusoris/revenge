package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Cache    CacheConfig    `koanf:"cache"`
	Log      LogConfig      `koanf:"log"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `koanf:"host"`     // PostgreSQL host
	Port     int    `koanf:"port"`     // PostgreSQL port
	User     string `koanf:"user"`     // PostgreSQL user
	Password string `koanf:"password"` // PostgreSQL password
	Name     string `koanf:"name"`     // PostgreSQL database name
	SSLMode  string `koanf:"sslmode"`  // PostgreSQL SSL mode (disable, require, verify-ca, verify-full)
}

// CacheConfig holds cache-related configuration (Dragonfly/Redis)
type CacheConfig struct {
	Addr     string `koanf:"addr"`     // Dragonfly/Redis address (host:port)
	Password string `koanf:"password"` // Dragonfly/Redis password
	DB       int    `koanf:"db"`       // Dragonfly/Redis database number
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `koanf:"level"`  // debug, info, warn, error
	Format string `koanf:"format"` // json or console
}

// New creates a new configuration instance
func New() (*Config, error) {
	k := koanf.New(".")

	// Load defaults
	defaults := Defaults()

	// Set defaults programmatically
	k.Load(nil, nil) // Initialize

	// Load main config file (optional)
	if err := k.Load(file.Provider("configs/config.yaml"), yaml.Parser()); err != nil {
		// Config file is optional, use defaults
	}

	// Load environment-specific config (optional)
	envConfig := os.Getenv("JELLYFIN_ENV")
	if envConfig != "" {
		configPath := fmt.Sprintf("configs/config.%s.yaml", envConfig)
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			// Environment-specific config is optional
		}
	}

	// Load environment variables (highest priority)
	// JELLYFIN_SERVER_PORT=8080 becomes server.port=8080
	if err := k.Load(env.Provider("JELLYFIN_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "JELLYFIN_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load env vars: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		// If unmarshal fails, use defaults
		return defaults, nil
	}

	// Apply defaults for missing values
	if cfg.Server.Host == "" {
		cfg.Server = defaults.Server
	}
	if cfg.Log.Level == "" {
		cfg.Log = defaults.Log
	}

	return &cfg, nil
}

// Defaults returns default configuration values
func Defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8096,
		},
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			User:    "jellyfin",
			Name:    "jellyfin",
			SSLMode: "disable",
		},
		Cache: CacheConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "console",
		},
	}
}
