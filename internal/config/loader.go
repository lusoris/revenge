package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	// EnvPrefix is the prefix for environment variables.
	EnvPrefix = "REVENGE_"

	// DefaultConfigPath is the default path to the configuration file.
	DefaultConfigPath = "config/config.yaml"
)

// Load loads the configuration from files and environment variables.
// Priority (highest to lowest):
// 1. Environment variables (REVENGE_*)
// 2. Config file
// 3. Default values
func Load(configPath string) (*Config, error) {
	k := koanf.New(".")

	// Load defaults first
	if err := k.Load(confmap.Provider(Defaults(), "."), nil); err != nil {
		return nil, err
	}

	// Load config file if it exists
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			// Config file is optional, only log if debug
			_ = err
		}
	}

	// Load environment variables (highest priority)
	// REVENGE_SERVER_PORT -> server.port
	if err := k.Load(env.Provider(EnvPrefix, ".", func(s string) string {
		return strings.ReplaceAll(
			strings.ToLower(strings.TrimPrefix(s, EnvPrefix)),
			"_",
			".",
		)
	}), nil); err != nil {
		return nil, err
	}

	// Fix compound field names where the simple _→. transform is wrong.
	// e.g. REVENGE_SEARCH_API_KEY → "search.api.key" but should be "search.api_key"
	applyCompoundEnvOverrides(k)

	// Unmarshal into Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadWithKoanf loads configuration and returns both Config and the underlying koanf instance.
// This is useful for accessing raw configuration values or watching for changes.
func LoadWithKoanf(configPath string) (*Config, *koanf.Koanf, error) {
	k := koanf.New(".")

	// Load defaults first
	if err := k.Load(confmap.Provider(Defaults(), "."), nil); err != nil {
		return nil, nil, err
	}

	// Load config file if it exists
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			// Config file is optional
			_ = err
		}
	}

	// Load environment variables (highest priority)
	if err := k.Load(env.Provider(EnvPrefix, ".", func(s string) string {
		return strings.ReplaceAll(
			strings.ToLower(strings.TrimPrefix(s, EnvPrefix)),
			"_",
			".",
		)
	}), nil); err != nil {
		return nil, nil, err
	}

	// Fix compound field names (same as Load)
	applyCompoundEnvOverrides(k)

	// Unmarshal into Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, nil, err
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, nil, err
	}

	return &cfg, k, nil
}

// MustLoad loads configuration and panics on error.
// Use this in main() or during application initialization.
func MustLoad(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		panic("failed to load configuration: " + err.Error())
	}
	return cfg
}

// compoundEnvVars maps environment variables with compound field names
// (containing underscores) to their correct koanf keys. The simple _→.
// transform converts ALL underscores to dots, which breaks field names
// like "api_key" → "api.key" instead of keeping them as "api_key".
var compoundEnvVars = map[string]string{
	"REVENGE_SEARCH_API_KEY":               "search.api_key",
	"REVENGE_MOVIE_TMDB_API_KEY":           "movie.tmdb.api_key",
	"REVENGE_MOVIE_TMDB_RATE_LIMIT":        "movie.tmdb.rate_limit",
	"REVENGE_MOVIE_TMDB_CACHE_TTL":         "movie.tmdb.cache_ttl",
	"REVENGE_INTEGRATIONS_RADARR_API_KEY":  "integrations.radarr.api_key",
	"REVENGE_INTEGRATIONS_RADARR_BASE_URL": "integrations.radarr.base_url",
	"REVENGE_INTEGRATIONS_SONARR_API_KEY":  "integrations.sonarr.api_key",
	"REVENGE_INTEGRATIONS_SONARR_BASE_URL": "integrations.sonarr.base_url",
	"REVENGE_EMAIL_SENDGRID_API_KEY":       "email.sendgrid.api_key",
	"REVENGE_EMAIL_SMTP_SKIP_VERIFY":       "email.smtp.skip_verify",
	"REVENGE_LEGACY_ENCRYPTION_KEY":        "legacy.encryption_key",
	"REVENGE_SERVER_RATE_LIMIT_ENABLED":                    "server.rate_limit.enabled",
	"REVENGE_SERVER_RATE_LIMIT_BACKEND":                    "server.rate_limit.backend",
	"REVENGE_SERVER_RATE_LIMIT_GLOBAL_REQUESTS_PER_SECOND": "server.rate_limit.global.requests_per_second",
	"REVENGE_SERVER_RATE_LIMIT_GLOBAL_BURST":               "server.rate_limit.global.burst",
	"REVENGE_SERVER_RATE_LIMIT_AUTH_REQUESTS_PER_SECOND":   "server.rate_limit.auth.requests_per_second",
	"REVENGE_SERVER_RATE_LIMIT_AUTH_BURST":                  "server.rate_limit.auth.burst",
}

// applyCompoundEnvOverrides fixes env vars that the simple _→. transform
// maps incorrectly due to compound field names containing underscores.
func applyCompoundEnvOverrides(k *koanf.Koanf) {
	for envVar, configKey := range compoundEnvVars {
		if v := os.Getenv(envVar); v != "" {
			_ = k.Set(configKey, v)
		}
	}
}

// validate validates the configuration using the validator package.
func validate(cfg *Config) error {
	v := validator.New()
	return v.Struct(cfg)
}
