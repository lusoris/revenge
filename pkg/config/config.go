// Package config provides configuration management for Revenge Go.
// It uses koanf v2 for hierarchical configuration from files and environment variables.
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
	Search   SearchConfig   `koanf:"search"`
	Auth     AuthConfig     `koanf:"auth"`
	OIDC     OIDCConfig     `koanf:"oidc"`
	Log      LogConfig      `koanf:"log"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// JWT settings
	JWTSecret            string `koanf:"jwt_secret"`             // Secret key for signing JWTs (min 32 chars)
	AccessTokenDuration  string `koanf:"access_token_duration"`  // Duration string, e.g., "15m"
	RefreshTokenDuration string `koanf:"refresh_token_duration"` // Duration string, e.g., "7d"

	// Password settings
	BcryptCost int `koanf:"bcrypt_cost"` // bcrypt cost factor (10-14 recommended)

	// Session settings
	MaxSessionsPerUser int `koanf:"max_sessions_per_user"` // 0 = unlimited
}

// OIDCConfig holds OIDC/SSO configuration
type OIDCConfig struct {
	Enabled   bool                 `koanf:"enabled"`   // Enable OIDC authentication
	Providers []OIDCProviderConfig `koanf:"providers"` // Static provider configuration
}

// OIDCProviderConfig holds configuration for a single OIDC provider
type OIDCProviderConfig struct {
	Name            string   `koanf:"name"`              // Internal name (e.g., "keycloak")
	DisplayName     string   `koanf:"display_name"`      // UI display name
	IssuerURL       string   `koanf:"issuer_url"`        // OIDC issuer URL
	ClientID        string   `koanf:"client_id"`         // OAuth2 client ID
	ClientSecret    string   `koanf:"client_secret"`     // OAuth2 client secret
	Scopes          []string `koanf:"scopes"`            // OAuth2 scopes (default: openid, profile, email)
	AutoCreateUsers bool     `koanf:"auto_create_users"` // Create users on first login
	DefaultAdmin    bool     `koanf:"default_admin"`     // New users are admins
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

// DatabaseConfig holds PostgreSQL configuration (REQUIRED)
type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Name     string `koanf:"name"`
	SSLMode  string `koanf:"ssl_mode"`
	MaxConns int    `koanf:"max_conns"`
}

// CacheConfig holds Dragonfly configuration (REQUIRED)
type CacheConfig struct {
	Addr     string `koanf:"addr"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

// SearchConfig holds Typesense configuration (REQUIRED)
type SearchConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	APIKey   string `koanf:"api_key"`
	Protocol string `koanf:"protocol"`
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
	_ = k.Load(nil, nil) //nolint:errcheck // Initialize - nil provider always succeeds

	// Load main config file (optional)
	_ = k.Load(file.Provider("configs/config.yaml"), yaml.Parser()) //nolint:errcheck // config file is optional

	// Load environment-specific config (optional)
	envConfig := os.Getenv("REVENGE_ENV")
	if envConfig != "" {
		configPath := fmt.Sprintf("configs/config.%s.yaml", envConfig)
		_ = k.Load(file.Provider(configPath), yaml.Parser()) //nolint:errcheck // env-specific config is optional
	}

	// Load environment variables (highest priority)
	// REVENGE_SERVER_PORT=8080 becomes server.port=8080
	if err := k.Load(env.Provider("REVENGE_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "REVENGE_")), "_", ".")
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
			Host:     "localhost",
			Port:     5432,
			User:     "revenge",
			Name:     "revenge",
			SSLMode:  "disable",
			MaxConns: 25,
		},
		Cache: CacheConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Search: SearchConfig{
			Host:     "localhost",
			Port:     8108,
			Protocol: "http",
		},
		Auth: AuthConfig{
			JWTSecret:            "", // Must be set in production
			AccessTokenDuration:  "15m",
			RefreshTokenDuration: "7d",
			BcryptCost:           12,
			MaxSessionsPerUser:   0, // Unlimited
		},
		OIDC: OIDCConfig{
			Enabled:   false,
			Providers: nil,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "console",
		},
	}
}
