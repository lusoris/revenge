package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	Log      LogConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type     string // sqlite or postgres
	Path     string // for SQLite
	Host     string // for PostgreSQL
	Port     int    // for PostgreSQL
	User     string // for PostgreSQL
	Password string // for PostgreSQL
	Name     string // for PostgreSQL
}

// CacheConfig holds cache-related configuration
type CacheConfig struct {
	Type  string // memory or redis
	Addr  string // for Redis
	Password string // for Redis
	DB    int    // for Redis
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Environment variables
	viper.SetEnvPrefix("JELLYFIN")
	viper.AutomaticEnv()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8096)

	// Database defaults
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.path", "./data/jellyfin.db")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "jellyfin")
	viper.SetDefault("database.name", "jellyfin")

	// Cache defaults
	viper.SetDefault("cache.type", "memory")
	viper.SetDefault("cache.addr", "localhost:6379")
	viper.SetDefault("cache.db", 0)

	// Log defaults
	viper.SetDefault("log.level", "info")
}
