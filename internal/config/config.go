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

	// Session configuration
	Session SessionConfig `koanf:"session"`

	// RBAC configuration
	RBAC RBACConfig `koanf:"rbac"`

	// Movie module configuration
	Movie MovieConfig `koanf:"movie"`

	// Integrations configuration
	Integrations IntegrationsConfig `koanf:"integrations"`

	// Legacy (QAR) module configuration
	Legacy LegacyConfig `koanf:"legacy"`

	// Email configuration
	Email EmailConfig `koanf:"email"`

	// Avatar configuration
	Avatar AvatarConfig `koanf:"avatar"`

	// Storage configuration (for avatars and user-generated content)
	Storage StorageConfig `koanf:"storage"`

	// Activity configuration
	Activity ActivityConfig `koanf:"activity"`

	// Raft configuration (leader election for clusters)
	Raft RaftConfig `koanf:"raft"`
}

// ActivityConfig holds activity log configuration.
type ActivityConfig struct {
	// RetentionDays is the number of days to retain activity logs.
	// Logs older than this will be automatically deleted by cleanup jobs.
	// Default: 90 days.
	RetentionDays int `koanf:"retention_days"`
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

	// Backend specifies the rate limiting backend: "memory" or "redis".
	// When "redis" is selected but unavailable, falls back to "memory".
	Backend string `koanf:"backend"`

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

	// MaxAttempts is the maximum number of retry attempts for failed jobs.
	MaxAttempts int `koanf:"max_attempts"`
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
	JWTSecret string `koanf:"jwt_secret" validate:"omitempty,min=32"`

	// JWTExpiry is the duration for JWT token validity.
	JWTExpiry time.Duration `koanf:"jwt_expiry"`

	// RefreshExpiry is the duration for refresh token validity.
	RefreshExpiry time.Duration `koanf:"refresh_expiry"`

	// LockoutThreshold is the number of failed login attempts before account lockout.
	// Default: 5 attempts
	LockoutThreshold int `koanf:"lockout_threshold"`

	// LockoutWindow is the time window for counting failed login attempts.
	// Failed attempts older than this are ignored.
	// Default: 15 minutes
	LockoutWindow time.Duration `koanf:"lockout_window"`

	// LockoutEnabled controls whether account lockout is enabled.
	// Default: true
	LockoutEnabled bool `koanf:"lockout_enabled"`
}

// RBACConfig holds RBAC configuration.
type RBACConfig struct {
	// ModelPath is the path to the Casbin model file.
	ModelPath string `koanf:"model_path"`

	// PolicyReloadInterval is the interval to reload policies from database.
	PolicyReloadInterval time.Duration `koanf:"policy_reload_interval"`
}

// SessionConfig holds session management configuration.
type SessionConfig struct {
	// CacheEnabled indicates if session caching is enabled (Dragonfly/Redis L1).
	CacheEnabled bool `koanf:"cache_enabled"`

	// CacheTTL is the TTL for cached sessions.
	CacheTTL time.Duration `koanf:"cache_ttl"`

	// MaxPerUser is the maximum number of active sessions per user.
	MaxPerUser int `koanf:"max_per_user"`

	// TokenLength is the length of generated session tokens in bytes.
	TokenLength int `koanf:"token_length"`
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

	// Sonarr integration configuration
	Sonarr SonarrConfig `koanf:"sonarr"`
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

// SonarrConfig holds Sonarr integration configuration.
type SonarrConfig struct {
	// Enabled indicates if Sonarr integration is enabled.
	Enabled bool `koanf:"enabled"`

	// BaseURL is the Sonarr server URL (e.g., http://localhost:8989).
	BaseURL string `koanf:"base_url"`

	// APIKey is the Sonarr API key.
	APIKey string `koanf:"api_key"`

	// AutoSync enables automatic library sync.
	AutoSync bool `koanf:"auto_sync"`

	// SyncInterval is the interval between automatic syncs (seconds).
	SyncInterval int `koanf:"sync_interval"`
}

// EmailConfig holds email/SMTP configuration for transactional emails.
type EmailConfig struct {
	// Enabled indicates if email sending is enabled.
	Enabled bool `koanf:"enabled"`

	// Provider is the email provider: "smtp" or "sendgrid".
	Provider string `koanf:"provider" validate:"omitempty,oneof=smtp sendgrid"`

	// FromAddress is the sender email address.
	FromAddress string `koanf:"from_address"`

	// FromName is the sender display name.
	FromName string `koanf:"from_name"`

	// BaseURL is the application base URL for email links (e.g., https://myserver.com).
	BaseURL string `koanf:"base_url"`

	// SMTP configuration (when provider=smtp)
	SMTP SMTPConfig `koanf:"smtp"`

	// SendGrid configuration (when provider=sendgrid)
	SendGrid SendGridConfig `koanf:"sendgrid"`
}

// SMTPConfig holds SMTP server configuration.
type SMTPConfig struct {
	// Host is the SMTP server hostname.
	Host string `koanf:"host"`

	// Port is the SMTP server port.
	Port int `koanf:"port"`

	// Username for SMTP authentication.
	Username string `koanf:"username"`

	// Password for SMTP authentication.
	Password string `koanf:"password"`

	// UseTLS enables TLS from the start (port 465).
	UseTLS bool `koanf:"use_tls"`

	// UseStartTLS enables STARTTLS upgrade (port 587).
	UseStartTLS bool `koanf:"use_starttls"`

	// SkipVerify skips TLS certificate verification (for self-signed certs).
	SkipVerify bool `koanf:"skip_verify"`

	// Timeout is the connection timeout.
	Timeout time.Duration `koanf:"timeout"`
}

// SendGridConfig holds SendGrid API configuration.
type SendGridConfig struct {
	// APIKey is the SendGrid API key.
	APIKey string `koanf:"api_key"`
}

// AvatarConfig holds avatar upload configuration.
type AvatarConfig struct {
	// StoragePath is the local directory for avatar storage.
	StoragePath string `koanf:"storage_path"`

	// MaxSizeBytes is the maximum allowed avatar file size.
	MaxSizeBytes int64 `koanf:"max_size_bytes"`

	// AllowedTypes are the allowed MIME types for avatars.
	AllowedTypes []string `koanf:"allowed_types"`
}

// StorageConfig holds configuration for file storage backend.
type StorageConfig struct {
	// Backend specifies the storage backend: "local" or "s3"
	Backend string `koanf:"backend" validate:"oneof=local s3"`

	// Local configuration (used when Backend is "local")
	Local LocalStorageConfig `koanf:"local"`

	// S3 configuration (used when Backend is "s3")
	S3 S3Config `koanf:"s3"`
}

// LocalStorageConfig holds configuration for local filesystem storage.
type LocalStorageConfig struct {
	// Path is the base directory for file storage
	Path string `koanf:"path" validate:"required_if=Backend local"`
}

// S3Config holds configuration for S3-compatible storage (AWS S3, MinIO, etc).
type S3Config struct {
	// Endpoint is the S3 endpoint URL (for MinIO: "http://minio:9000")
	// Leave empty for AWS S3
	Endpoint string `koanf:"endpoint"`

	// Region is the S3 region (e.g., "us-east-1")
	Region string `koanf:"region" validate:"required_if=Backend s3"`

	// Bucket is the S3 bucket name
	Bucket string `koanf:"bucket" validate:"required_if=Backend s3"`

	// AccessKeyID is the S3 access key ID
	AccessKeyID string `koanf:"access_key_id" validate:"required_if=Backend s3"`

	// SecretAccessKey is the S3 secret access key
	SecretAccessKey string `koanf:"secret_access_key" validate:"required_if=Backend s3"`

	// UsePathStyle enables path-style S3 URLs (required for MinIO)
	UsePathStyle bool `koanf:"use_path_style"`
}

// RaftConfig holds configuration for Raft leader election in cluster deployments.
type RaftConfig struct {
	// Enabled controls whether Raft leader election is active
	Enabled bool `koanf:"enabled"`

	// NodeID is the unique identifier for this node (hostname or UUID)
	// If empty, hostname will be used automatically
	NodeID string `koanf:"node_id"`

	// BindAddr is the address for Raft communication (e.g., "0.0.0.0:7000")
	BindAddr string `koanf:"bind_addr" validate:"required_if=Enabled true"`

	// DataDir is the directory for Raft data storage
	DataDir string `koanf:"data_dir" validate:"required_if=Enabled true"`

	// Bootstrap should be true only for the first node to initialize the cluster
	Bootstrap bool `koanf:"bootstrap"`
}

// GetRadarrConfig returns the Radarr configuration.
func (c *Config) GetRadarrConfig() RadarrConfig {
	return c.Integrations.Radarr
}

// GetSonarrConfig returns the Sonarr configuration.
func (c *Config) GetSonarrConfig() SonarrConfig {
	return c.Integrations.Sonarr
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
		"cache.url":     "",
		"cache.enabled": false,

		// Rate limit defaults
		"rate_limit.enabled":                  true,
		"rate_limit.backend":                  "memory", // "memory" or "redis"
		"rate_limit.global.requests_per_second": 10.0,
		"rate_limit.global.burst":             20,
		"rate_limit.auth.requests_per_second": 1.0,
		"rate_limit.auth.burst":               5,

		// Movie defaults
		"movie.tmdb.api_key":    "",
		"movie.tmdb.rate_limit": 40,
		"movie.tmdb.cache_ttl":  "5m",
		"movie.tmdb.proxy_url":  "",
		"movie.library.paths":   []string{},
		"movie.library.scan_interval": "0s", // Disabled by default

		// Search defaults
		"search.url":     "",
		"search.api_key": "",
		"search.enabled": false,

		// Jobs defaults
		"jobs.max_workers":             100,
		"jobs.fetch_cooldown":          "200ms",
		"jobs.fetch_poll_interval":     "2s",
		"jobs.rescue_stuck_jobs_after": "30m",
		"jobs.max_attempts":            25,

		// Logging defaults
		"logging.level":       "info",
		"logging.format":      "text",
		"logging.development": false,

		// Auth defaults
		"auth.jwt_secret":        "",
		"auth.jwt_expiry":        "24h",
		"auth.refresh_expiry":    "168h", // 7 days
		"auth.lockout_threshold": 5,      // 5 failed attempts
		"auth.lockout_window":    "15m",  // 15 minutes
		"auth.lockout_enabled":   true,

		// RBAC defaults
		"rbac.model_path":             "config/casbin_model.conf",
		"rbac.policy_reload_interval": "5m",

		// Session defaults
		"session.cache_enabled": true,
		"session.cache_ttl":     "5m",
		"session.max_per_user":  10,
		"session.token_length":  32,

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

		// Email defaults
		"email.enabled":          false,
		"email.provider":         "smtp",
		"email.from_address":     "",
		"email.from_name":        "Revenge Media Server",
		"email.base_url":         "http://localhost:8080",
		"email.smtp.host":        "",
		"email.smtp.port":        587,
		"email.smtp.username":    "",
		"email.smtp.password":    "",
		"email.smtp.use_tls":     false,
		"email.smtp.use_starttls": true,
		"email.smtp.skip_verify": false,
		"email.smtp.timeout":     "30s",
		"email.sendgrid.api_key": "",

		// Avatar defaults
		"avatar.storage_path":   "/data/avatars",
		"avatar.max_size_bytes": 2 * 1024 * 1024, // 2MB
		"avatar.allowed_types":  []string{"image/jpeg", "image/png", "image/webp"},

		// Storage defaults
		"storage.backend":            "local", // Use local storage by default
		"storage.local.path":         "/data/storage",
		"storage.s3.region":          "us-east-1",
		"storage.s3.bucket":          "",
		"storage.s3.endpoint":        "",
		"storage.s3.access_key_id":   "",
		"storage.s3.secret_access_key": "",
		"storage.s3.use_path_style":  false,

		// Activity defaults
		"activity.retention_days": 90, // 90 days default retention

		// Raft defaults (disabled by default for single-node deployments)
		"raft.enabled":   false,
		"raft.node_id":   "", // Auto-detect from hostname
		"raft.bind_addr": "0.0.0.0:7000",
		"raft.data_dir":  "/data/raft",
		"raft.bootstrap": false,
	}
}
