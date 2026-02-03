// Package config provides configuration management for the Revenge server.
package config

import (
	"time"
)

// Config holds the complete application configuration.
type Config struct {
	// Server configuration
	Server ServerConfig `koanf:"server"`

	// Database configuration
	Database DatabaseConfig `koanf:"database"`

	// Cache configuration (Dragonfly/Redis)
	Cache CacheConfig `koanf:"cache"`

	// Search configuration (Typesense)
	Search SearchConfig `koanf:"search"`

	// Jobs configuration (River)
	Jobs JobsConfig `koanf:"jobs"`

	// Logging configuration
	Logging LoggingConfig `koanf:"logging"`

	// Auth configuration
	Auth AuthConfig `koanf:"auth"`

	// RBAC configuration
	RBAC RBACConfig `koanf:"rbac"`

	// Movie module configuration
	Movie MovieConfig `koanf:"movie"`

	// Integrations configuration
	Integrations IntegrationsConfig `koanf:"integrations"`

	// Legacy (QAR) module configuration
	Legacy LegacyConfig `koanf:"legacy"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	// Host is the server bind address.
	Host string `koanf:"host" validate:"required"`

	// Port is the server port.
	Port int `koanf:"port" validate:"required,min=1,max=65535"`

	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration `koanf:"read_timeout"`

	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration `koanf:"write_timeout"`

	// IdleTimeout is the maximum amount of time to wait for the next request.
	IdleTimeout time.Duration `koanf:"idle_timeout"`

	// ShutdownTimeout is the maximum duration to wait for active connections to finish.
	ShutdownTimeout time.Duration `koanf:"shutdown_timeout"`

	// RateLimit configures API rate limiting.
	RateLimit RateLimitConfig `koanf:"rate_limit"`
}

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	// Enabled controls whether rate limiting is active.
	Enabled bool `koanf:"enabled"`

	// Global configures global rate limiting for all endpoints.
	Global RateLimitTier `koanf:"global"`

	// Auth configures stricter rate limiting for auth endpoints.
	Auth RateLimitTier `koanf:"auth"`
}

// RateLimitTier holds rate limiting settings for a specific tier.
type RateLimitTier struct {
	// RequestsPerSecond is the number of requests allowed per second per IP.
	RequestsPerSecond float64 `koanf:"requests_per_second"`

	// Burst is the maximum number of requests allowed in a burst.
	Burst int `koanf:"burst"`
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	// URL is the PostgreSQL connection string.
	URL string `koanf:"url" validate:"required"`

	// MaxConns is the maximum number of connections in the pool.
	MaxConns int `koanf:"max_conns"`

	// MinConns is the minimum number of connections in the pool.
	MinConns int `koanf:"min_conns"`

	// MaxConnLifetime is the maximum lifetime of a connection.
	MaxConnLifetime time.Duration `koanf:"max_conn_lifetime"`

	// MaxConnIdleTime is the maximum idle time of a connection.
	MaxConnIdleTime time.Duration `koanf:"max_conn_idle_time"`

	// HealthCheckPeriod is the duration between health checks.
	HealthCheckPeriod time.Duration `koanf:"health_check_period"`
}

// CacheConfig holds cache (Dragonfly/Redis) configuration.
type CacheConfig struct {
	// URL is the Redis/Dragonfly connection URL.
	URL string `koanf:"url"`

	// Enabled indicates if cache is enabled.
	Enabled bool `koanf:"enabled"`
}

// SearchConfig holds search (Typesense) configuration.
type SearchConfig struct {
	// URL is the Typesense server URL.
	URL string `koanf:"url"`

	// APIKey is the Typesense API key.
	APIKey string `koanf:"api_key"`

	// Enabled indicates if search is enabled.
	Enabled bool `koanf:"enabled"`
}

// JobsConfig holds job queue (River) configuration.
type JobsConfig struct {
	// MaxWorkers is the maximum number of concurrent workers.
	MaxWorkers int `koanf:"max_workers"`

	// FetchCooldown is the duration to wait between fetch attempts.
	FetchCooldown time.Duration `koanf:"fetch_cooldown"`

	// FetchPollInterval is the interval between polling for new jobs.
	FetchPollInterval time.Duration `koanf:"fetch_poll_interval"`

	// RescueStuckJobsAfter is the duration after which stuck jobs are rescued.
	RescueStuckJobsAfter time.Duration `koanf:"rescue_stuck_jobs_after"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	// Level is the minimum log level (debug, info, warn, error).
	Level string `koanf:"level" validate:"oneof=debug info warn error"`

	// Format is the log format (text, json).
	Format string `koanf:"format" validate:"oneof=text json"`

	// Development enables development mode (pretty printing, etc.).
	Development bool `koanf:"development"`
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	// JWTSecret is the secret key for JWT signing.
	// Optional for v0.1.0 (auth not implemented yet)
	JWTSecret string `koanf:"jwt_secret" validate:"omitempty,min=32"`

	// JWTExpiry is the duration for JWT token validity.
	JWTExpiry time.Duration `koanf:"jwt_expiry"`

	// RefreshExpiry is the duration for refresh token validity.
	RefreshExpiry time.Duration `koanf:"refresh_expiry"`
}

// RBACConfig holds RBAC configuration.
type RBACConfig struct {
	// ModelPath is the path to the Casbin model file.
	ModelPath string `koanf:"model_path"`

	// PolicyReloadInterval is the interval to reload policies from database.
	PolicyReloadInterval time.Duration `koanf:"policy_reload_interval"`
}

// LegacyConfig holds QAR (adult content) module configuration.
type LegacyConfig struct {
	// Enabled indicates if the QAR module is enabled.
	Enabled bool `koanf:"enabled"`

	// EncryptionKey is the encryption key for QAR data.
	EncryptionKey string `koanf:"encryption_key"`

	// Privacy settings
	Privacy LegacyPrivacyConfig `koanf:"privacy"`
}

// LegacyPrivacyConfig holds QAR privacy settings.
type LegacyPrivacyConfig struct {
	// RequirePIN requires a PIN to access QAR content.
	RequirePIN bool `koanf:"require_pin"`

	// AuditAllAccess logs all access to QAR content.
	AuditAllAccess bool `koanf:"audit_all_access"`
}

// MovieConfig holds movie module configuration.
type MovieConfig struct {
	// TMDb configuration for metadata
	TMDb TMDbConfig `koanf:"tmdb"`

	// Library configuration for file scanning
	Library LibraryConfig `koanf:"library"`
}

// TMDbConfig holds TMDb API configuration.
type TMDbConfig struct {
	// APIKey is the TMDb API key (required if using TMDb metadata).
	APIKey string `koanf:"api_key"`

	// RateLimit is requests per 10 seconds (default: 40).
	RateLimit int `koanf:"rate_limit"`

	// CacheTTL is how long to cache TMDb responses (default: 5m).
	CacheTTL time.Duration `koanf:"cache_ttl"`

	// ProxyURL is optional SOCKS5/HTTP proxy for TMDb requests.
	ProxyURL string `koanf:"proxy_url"`
}

// LibraryConfig holds movie library configuration.
type LibraryConfig struct {
	// Paths are the directories to scan for movie files.
	Paths []string `koanf:"paths"`

	// ScanInterval is how often to automatically scan libraries (0 = disabled).
	ScanInterval time.Duration `koanf:"scan_interval"`
}

// IntegrationsConfig holds all external integrations configuration.
type IntegrationsConfig struct {
	// Radarr integration configuration
	Radarr RadarrConfig `koanf:"radarr"`
}

// RadarrConfig holds Radarr integration configuration.
type RadarrConfig struct {
	// Enabled indicates if Radarr integration is enabled.
	Enabled bool `koanf:"enabled"`

	// BaseURL is the Radarr server URL (e.g., http://localhost:7878).
	BaseURL string `koanf:"base_url"`

	// APIKey is the Radarr API key.
	APIKey string `koanf:"api_key"`

	// AutoSync enables automatic library sync.
	AutoSync bool `koanf:"auto_sync"`

	// SyncInterval is the interval between automatic syncs (seconds).
	SyncInterval int `koanf:"sync_interval"`
}

// GetRadarrConfig returns the Radarr configuration.
func (c *Config) GetRadarrConfig() RadarrConfig {
	return c.Integrations.Radarr
}

// Defaults returns a map of default configuration values.
func Defaults() map[string]interface{} {
	return map[string]interface{}{
		// Server defaults
		"server.host":             "0.0.0.0",
		"server.port":             8080,
		"server.read_timeout":     "30s",
		"server.write_timeout":    "30s",
		"server.idle_timeout":     "120s",
		"server.shutdown_timeout": "10s",

		// Database defaults
		"database.url":                 "postgres://revenge:changeme@localhost:5432/revenge?sslmode=disable",
		"database.max_conns":           0, // 0 = (CPU * 2) + 1
		"database.min_conns":           2,
		"database.max_conn_lifetime":   "30m",
		"database.max_conn_idle_time":  "5m",
		"database.health_check_period": "30s",

		// Cache defaults

		// Movie defaults
		"movie.tmdb.api_key":    "",
		"movie.tmdb.rate_limit": 40,
		"movie.tmdb.cache_ttl":  "5m",
		"movie.tmdb.proxy_url":  "",
		"movie.library.paths":   []string{},
		"movie.library.scan_interval": "0s", // Disabled by default
		"cache.url":     "",
		"cache.enabled": false,

		// Search defaults
		"search.url":     "",
		"search.api_key": "",
		"search.enabled": false,

		// Jobs defaults
		"jobs.max_workers":             100,
		"jobs.fetch_cooldown":          "200ms",
		"jobs.fetch_poll_interval":     "2s",
		"jobs.rescue_stuck_jobs_after": "30m",

		// Logging defaults
		"logging.level":       "info",
		"logging.format":      "text",
		"logging.development": false,

		// Auth defaults
		"auth.jwt_secret":     "",
		"auth.jwt_expiry":     "24h",
		"auth.refresh_expiry": "168h", // 7 days

		// RBAC defaults
		"rbac.model_path":             "config/casbin_model.conf",
		"rbac.policy_reload_interval": "5m",

		// Legacy defaults
		"legacy.enabled":                  false,
		"legacy.encryption_key":           "",
		"legacy.privacy.require_pin":      true,
		"legacy.privacy.audit_all_access": true,

		// Integrations defaults
		"integrations.radarr.enabled":       false,
		"integrations.radarr.base_url":      "http://localhost:7878",
		"integrations.radarr.api_key":       "",
		"integrations.radarr.auto_sync":     false,
		"integrations.radarr.sync_interval": 300, // 5 minutes
	}
}
