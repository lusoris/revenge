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
	Type     string `koanf:"type"` // sqlite or postgres
	Path     string `koanf:"path"` // for SQLite
	Host     string `koanf:"host"` // for PostgreSQL
	Port     int    `koanf:"port"` // for PostgreSQL
	User     string `koanf:"user"` // for PostgreSQL
	Password string `koanf:"password"` // for PostgreSQL
	Name     string `koanf:"name"` // for PostgreSQL
}

// CacheConfig holds cache-related configuration
type CacheConfig struct {
	Type     string `koanf:"type"` // memory or redis
	Addr     string `koanf:"addr"` // for Redis
	Password string `koanf:"password"` // for Redis
	DB       int    `koanf:"db"` // for Redis
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
			Type: "sqlite",
			Path: "./data/jellyfin.db",
			Host: "localhost",
			Port: 5432,
			User: "jellyfin",
			Name: "jellyfin",
		},
		Cache: CacheConfig{
			Type: "memory",
			Addr: "localhost:6379",
			DB:   0,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "console",
		},
	}
}
