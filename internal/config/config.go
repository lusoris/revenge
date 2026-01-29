// Package config provides configuration management using koanf v2.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Cache    CacheConfig    `koanf:"cache"`
	Search   SearchConfig   `koanf:"search"`
	Auth     AuthConfig     `koanf:"auth"`
	Metadata MetadataConfig `koanf:"metadata"`
	Modules  ModulesConfig  `koanf:"modules"`
	Logging  LoggingConfig  `koanf:"logging"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host              string        `koanf:"host"`
	Port              int           `koanf:"port"`
	BaseURL           string        `koanf:"base_url"`
	ReadTimeout       time.Duration `koanf:"read_timeout"`
	WriteTimeout      time.Duration `koanf:"write_timeout"`
	IdleTimeout       time.Duration `koanf:"idle_timeout"`
	ReadHeaderTimeout time.Duration `koanf:"read_header_timeout"`
	MaxHeaderBytes    int           `koanf:"max_header_bytes"`
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Name     string `koanf:"name"`
	SSLMode  string `koanf:"sslmode"`
	MaxConns int32  `koanf:"max_conns"`
	MinConns int32  `koanf:"min_conns"`
}

// DSN returns the PostgreSQL connection string.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

// CacheConfig holds Dragonfly/Redis connection settings.
type CacheConfig struct {
	Addr     string `koanf:"addr"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`

	// Local cache settings (otter)
	LocalCapacity int `koanf:"local_capacity"`
	LocalTTL      int `koanf:"local_ttl"` // seconds

	// API cache settings (sturdyc)
	APICapacity  int `koanf:"api_capacity"`
	APINumShards int `koanf:"api_num_shards"`
	APITTL       int `koanf:"api_ttl"` // seconds
}

// SearchConfig holds Typesense connection settings.
type SearchConfig struct {
	Host   string `koanf:"host"`
	Port   int    `koanf:"port"`
	APIKey string `koanf:"api_key"`
}

// URL returns the Typesense URL.
func (c SearchConfig) URL() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	JWTSecret       string `koanf:"jwt_secret"`
	SessionDuration int    `koanf:"session_duration"` // hours
}

// MetadataConfig holds metadata provider settings.
type MetadataConfig struct {
	Radarr RadarrConfig `koanf:"radarr"`
	TMDb   TMDbConfig   `koanf:"tmdb"`
}

// RadarrConfig holds Radarr connection settings.
type RadarrConfig struct {
	BaseURL string `koanf:"base_url"`
	APIKey  string `koanf:"api_key"`
}

// TMDbConfig holds TMDb connection settings.
type TMDbConfig struct {
	APIKey     string `koanf:"api_key"`
	BaseURL    string `koanf:"base_url"`
	ImageURL   string `koanf:"image_url"`
	Timeout    int    `koanf:"timeout"`   // seconds
	CacheTTL   int    `koanf:"cache_ttl"` // seconds
	CacheSize  int    `koanf:"cache_size"`
	RetryCount int    `koanf:"retry_count"`
}

// ModulesConfig holds settings for which modules are enabled.
type ModulesConfig struct {
	Movie     bool `koanf:"movie"`
	TVShow    bool `koanf:"tvshow"`
	Music     bool `koanf:"music"`
	Audiobook bool `koanf:"audiobook"`
	Book      bool `koanf:"book"`
	Podcast   bool `koanf:"podcast"`
	Photo     bool `koanf:"photo"`
	LiveTV    bool `koanf:"livetv"`
	Comics    bool `koanf:"comics"`
	Adult     bool `koanf:"adult"` // Explicit opt-in
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level  string `koanf:"level"`  // debug, info, warn, error
	Format string `koanf:"format"` // json, text
}

// Load loads configuration from file and environment variables.
// Environment variables use REVENGE_ prefix (e.g., REVENGE_DATABASE_HOST).
func Load(configPath string) (*Config, error) {
	k := koanf.New(".")

	// Load from YAML file if provided
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("load config file: %w", err)
		}
	}

	// Load from environment variables (overrides file)
	// REVENGE_DATABASE_HOST -> database.host
	if err := k.Load(env.Provider("REVENGE_", ".", func(s string) string {
		return strings.ReplaceAll(
			strings.ToLower(strings.TrimPrefix(s, "REVENGE_")),
			"_", ".",
		)
	}), nil); err != nil {
		return nil, fmt.Errorf("load env vars: %w", err)
	}

	applyEnvAliases(k)

	// Set defaults
	setDefaults(k)

	// Unmarshal into struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadDefault attempts to load configuration from standard paths, then env-only.
func LoadDefault() (*Config, error) {
	paths := []string{
		"config.yaml",
		"configs/config.yaml",
		"/etc/revenge/config.yaml",
	}

	for _, path := range paths {
		cfg, err := Load(path)
		if err == nil {
			return cfg, nil
		}
	}

	return Load("")
}

func applyEnvAliases(k *koanf.Koanf) {
	setIfMissing := func(key string, value any) {
		if !k.Exists(key) {
			_ = k.Set(key, value)
		}
	}

	if v, ok := os.LookupEnv("REVENGE_DB_HOST"); ok {
		setIfMissing("database.host", v)
	}
	if v, ok := os.LookupEnv("REVENGE_DB_PORT"); ok {
		if port, err := strconv.Atoi(v); err == nil {
			setIfMissing("database.port", port)
		}
	}
	if v, ok := os.LookupEnv("REVENGE_DB_USER"); ok {
		setIfMissing("database.user", v)
	}
	if v, ok := os.LookupEnv("REVENGE_DB_PASSWORD"); ok {
		setIfMissing("database.password", v)
	}
	if v, ok := os.LookupEnv("REVENGE_DB_NAME"); ok {
		setIfMissing("database.name", v)
	}
	if v, ok := os.LookupEnv("REVENGE_DB_SSLMODE"); ok {
		setIfMissing("database.sslmode", v)
	}

	if v, ok := os.LookupEnv("REVENGE_CACHE_URL"); ok {
		setIfMissing("cache.addr", v)
	}
	if v, ok := os.LookupEnv("REVENGE_CACHE_PASSWORD"); ok {
		setIfMissing("cache.password", v)
	}

	if v, ok := os.LookupEnv("REVENGE_TYPESENSE_URL"); ok {
		if parsed, err := url.Parse(v); err == nil {
			if host := parsed.Hostname(); host != "" {
				setIfMissing("search.host", host)
			}
			if port := parsed.Port(); port != "" {
				if portInt, err := strconv.Atoi(port); err == nil {
					setIfMissing("search.port", portInt)
				}
			}
		}
	}
	if v, ok := os.LookupEnv("REVENGE_TYPESENSE_API_KEY"); ok {
		setIfMissing("search.api_key", v)
	}

	if v, ok := os.LookupEnv("REVENGE_HOST"); ok {
		setIfMissing("server.host", v)
	}
	if v, ok := os.LookupEnv("REVENGE_PORT"); ok {
		if port, err := strconv.Atoi(v); err == nil {
			setIfMissing("server.port", port)
		}
	}
	if v, ok := os.LookupEnv("REVENGE_BASE_URL"); ok {
		setIfMissing("server.base_url", v)
	}

	if v, ok := os.LookupEnv("REVENGE_LOG_LEVEL"); ok {
		setIfMissing("logging.level", v)
	}
	if v, ok := os.LookupEnv("REVENGE_LOG_FORMAT"); ok {
		setIfMissing("logging.format", v)
	}

	if v := k.String("search.url"); v != "" && !k.Exists("search.host") {
		if parsed, err := url.Parse(v); err == nil {
			if host := parsed.Hostname(); host != "" {
				setIfMissing("search.host", host)
			}
			if port := parsed.Port(); port != "" {
				if portInt, err := strconv.Atoi(port); err == nil {
					setIfMissing("search.port", portInt)
				}
			}
		}
	}
}

// setDefaults sets default configuration values.
func setDefaults(k *koanf.Koanf) {
	defaults := map[string]any{
		// Server
		"server.host":                "0.0.0.0",
		"server.port":                8096,
		"server.base_url":            "/",
		"server.read_timeout":        30 * time.Second,
		"server.write_timeout":       30 * time.Second,
		"server.idle_timeout":        60 * time.Second,
		"server.read_header_timeout": 5 * time.Second,
		"server.max_header_bytes":    1 << 20,

		// Database
		"database.host":      "localhost",
		"database.port":      5432,
		"database.user":      "revenge",
		"database.password":  "",
		"database.name":      "revenge",
		"database.sslmode":   "disable",
		"database.max_conns": 25,
		"database.min_conns": 5,

		// Cache
		"cache.addr":           "localhost:6379",
		"cache.password":       "",
		"cache.db":             0,
		"cache.local_capacity": 10000,
		"cache.local_ttl":      300, // 5 minutes
		"cache.api_capacity":   5000,
		"cache.api_num_shards": 10,
		"cache.api_ttl":        3600, // 1 hour

		// Search
		"search.host":    "localhost",
		"search.port":    8108,
		"search.api_key": "",

		// Auth
		"auth.jwt_secret":       "",
		"auth.session_duration": 24,

		// Metadata - Radarr
		"metadata.radarr.base_url": "",
		"metadata.radarr.api_key":  "",

		// Metadata - TMDb
		"metadata.tmdb.api_key":     "",
		"metadata.tmdb.base_url":    "https://api.themoviedb.org/3",
		"metadata.tmdb.image_url":   "https://image.tmdb.org/t/p",
		"metadata.tmdb.timeout":     30,
		"metadata.tmdb.cache_ttl":   3600,
		"metadata.tmdb.cache_size":  50000,
		"metadata.tmdb.retry_count": 3,

		// Modules (only core modules enabled by default)
		"modules.movie":     true,
		"modules.tvshow":    true,
		"modules.music":     true,
		"modules.audiobook": false,
		"modules.book":      false,
		"modules.podcast":   false,
		"modules.photo":     false,
		"modules.livetv":    false,
		"modules.comics":    false,
		"modules.adult":     false, // Explicit opt-in

		// Logging
		"logging.level":  "info",
		"logging.format": "json",
	}

	for key, value := range defaults {
		if !k.Exists(key) {
			_ = k.Set(key, value)
		}
	}
}

// Module provides configuration dependencies for fx.
var Module = fx.Module("config",
	fx.Provide(func() (*Config, error) {
		// Try config.yaml in current directory, then /etc/revenge/
		paths := []string{
			"config.yaml",
			"configs/config.yaml",
			"/etc/revenge/config.yaml",
		}

		for _, path := range paths {
			cfg, err := Load(path)
			if err == nil {
				return cfg, nil
			}
		}

		// Fall back to env-only config
		return Load("")
	}),
)
