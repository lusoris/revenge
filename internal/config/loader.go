package config

import (
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

// validate validates the configuration using the validator package.
func validate(cfg *Config) error {
	v := validator.New()
	return v.Struct(cfg)
}
