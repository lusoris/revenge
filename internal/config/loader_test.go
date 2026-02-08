package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithDefaults(t *testing.T) {
	// Load with non-existent config file should still work with defaults
	cfg, err := Load("nonexistent.yaml")
	require.NoError(t, err, "Load should not error on missing config file")
	require.NotNil(t, cfg, "Config should not be nil")

	// Verify defaults are applied
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8096, cfg.Server.Port)
	assert.Contains(t, cfg.Database.URL, "postgres://")
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "text", cfg.Logging.Format)
}

func TestLoad_WithConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "127.0.0.1"
  port: 9090
database:
  url: "postgres://test:test@localhost:5432/test?sslmode=disable"
logging:
  level: debug
  format: json
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err, "Failed to write test config file")

	cfg, err := Load(configPath)
	require.NoError(t, err, "Load should succeed with valid config file")

	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Contains(t, cfg.Database.URL, "test:test@localhost")
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestLoad_WithEnvOverride(t *testing.T) {
	// Set environment variables
	t.Setenv("REVENGE_SERVER_PORT", "3000")
	t.Setenv("REVENGE_LOGGING_LEVEL", "warn")

	cfg, err := Load("")
	require.NoError(t, err, "Load should succeed with env vars")

	// Env vars should override defaults
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "warn", cfg.Logging.Level)
}

func TestLoad_EnvOverridesConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  port: 9090
logging:
  level: debug
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set env var that should override config file
	t.Setenv("REVENGE_SERVER_PORT", "4000")

	cfg, err := Load(configPath)
	require.NoError(t, err)

	// Env var wins over config file
	assert.Equal(t, 4000, cfg.Server.Port, "Env var should override config file")
	// Config file value should still be used where no env var exists
	assert.Equal(t, "debug", cfg.Logging.Level, "Config file value should be used")
}

func TestLoad_WithEmptyPath(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err, "Load with empty path should use defaults")
	require.NotNil(t, cfg)

	// Should have default values
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8096, cfg.Server.Port)
}

func TestLoadWithKoanf(t *testing.T) {
	cfg, k, err := LoadWithKoanf("")
	require.NoError(t, err, "LoadWithKoanf should succeed")
	require.NotNil(t, cfg, "Config should not be nil")
	require.NotNil(t, k, "Koanf instance should not be nil")

	// Verify we can access raw values through koanf
	assert.Equal(t, "0.0.0.0", k.String("server.host"))
	assert.Equal(t, 8096, k.Int("server.port"))
	assert.Equal(t, "info", k.String("logging.level"))
}

func TestLoadWithKoanf_WithEnvVars(t *testing.T) {
	t.Setenv("REVENGE_SERVER_HOST", "192.168.1.1")
	t.Setenv("REVENGE_SERVER_PORT", "9999")

	cfg, k, err := LoadWithKoanf("")
	require.NoError(t, err)

	// Config struct should have env values
	assert.Equal(t, "192.168.1.1", cfg.Server.Host)
	assert.Equal(t, 9999, cfg.Server.Port)

	// Koanf should also reflect env values
	assert.Equal(t, "192.168.1.1", k.String("server.host"))
	assert.Equal(t, 9999, k.Int("server.port"))
}

func TestMustLoad_Success(t *testing.T) {
	// MustLoad should not panic with valid/default config
	assert.NotPanics(t, func() {
		cfg := MustLoad("")
		assert.NotNil(t, cfg)
	})
}

func TestMustLoad_Panic(t *testing.T) {
	// MustLoad should panic when Load returns an error
	// The easiest way to trigger this is validation failure
	// We need to create a config that passes YAML parsing but fails validation

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// Config with invalid port (0 fails min=1 validation)
	invalidContent := `
server:
  host: ""
  port: 0
database:
  url: ""
`
	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	// MustLoad should panic on validation failure
	assert.Panics(t, func() {
		MustLoad(configPath)
	}, "MustLoad should panic on validation failure")
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := Default()

	// Default config should be valid
	err := validate(cfg)
	assert.NoError(t, err, "Default config should be valid")
}

func TestValidate_InvalidServerPort(t *testing.T) {
	cfg := Default()
	cfg.Server.Port = 0 // Invalid: min=1

	err := validate(cfg)
	assert.Error(t, err, "Port 0 should fail validation")
}

func TestValidate_InvalidServerPortMax(t *testing.T) {
	cfg := Default()
	cfg.Server.Port = 70000 // Invalid: max=65535

	err := validate(cfg)
	assert.Error(t, err, "Port > 65535 should fail validation")
}

func TestValidate_EmptyServerHost(t *testing.T) {
	cfg := Default()
	cfg.Server.Host = "" // Invalid: required

	err := validate(cfg)
	assert.Error(t, err, "Empty host should fail validation")
}

func TestValidate_EmptyDatabaseURL(t *testing.T) {
	cfg := Default()
	cfg.Database.URL = "" // Invalid: required

	err := validate(cfg)
	assert.Error(t, err, "Empty database URL should fail validation")
}

func TestValidate_InvalidLoggingLevel(t *testing.T) {
	cfg := Default()
	cfg.Logging.Level = "invalid" // Invalid: oneof=debug info warn error

	err := validate(cfg)
	assert.Error(t, err, "Invalid logging level should fail validation")
}

func TestValidate_InvalidLoggingFormat(t *testing.T) {
	cfg := Default()
	cfg.Logging.Format = "xml" // Invalid: oneof=text json

	err := validate(cfg)
	assert.Error(t, err, "Invalid logging format should fail validation")
}

func TestValidate_JWTSecretTooShort(t *testing.T) {
	cfg := Default()
	cfg.Auth.JWTSecret = "short" // Invalid: min=32 if set

	err := validate(cfg)
	assert.Error(t, err, "Short JWT secret should fail validation")
}

func TestValidate_JWTSecretValid(t *testing.T) {
	cfg := Default()
	cfg.Auth.JWTSecret = "this-is-a-valid-secret-key-32chars" // 35 chars, valid

	err := validate(cfg)
	assert.NoError(t, err, "Valid JWT secret should pass validation")
}

func TestEnvPrefix(t *testing.T) {
	assert.Equal(t, "REVENGE_", EnvPrefix, "EnvPrefix should be REVENGE_")
}

func TestDefaultConfigPath(t *testing.T) {
	assert.Equal(t, "config/config.yaml", DefaultConfigPath)
}

func TestDefaults_AllKeysExist(t *testing.T) {
	defaults := Defaults()

	// Server
	assert.Contains(t, defaults, "server.host")
	assert.Contains(t, defaults, "server.port")
	assert.Contains(t, defaults, "server.read_timeout")
	assert.Contains(t, defaults, "server.write_timeout")
	assert.Contains(t, defaults, "server.idle_timeout")
	assert.Contains(t, defaults, "server.shutdown_timeout")

	// Database
	assert.Contains(t, defaults, "database.url")
	assert.Contains(t, defaults, "database.max_conns")
	assert.Contains(t, defaults, "database.min_conns")
	assert.Contains(t, defaults, "database.max_conn_lifetime")
	assert.Contains(t, defaults, "database.max_conn_idle_time")
	assert.Contains(t, defaults, "database.health_check_period")

	// Cache
	assert.Contains(t, defaults, "cache.url")
	assert.Contains(t, defaults, "cache.enabled")

	// Search
	assert.Contains(t, defaults, "search.url")
	assert.Contains(t, defaults, "search.api_key")
	assert.Contains(t, defaults, "search.enabled")

	// Jobs
	assert.Contains(t, defaults, "jobs.max_workers")
	assert.Contains(t, defaults, "jobs.fetch_cooldown")
	assert.Contains(t, defaults, "jobs.fetch_poll_interval")
	assert.Contains(t, defaults, "jobs.rescue_stuck_jobs_after")

	// Logging
	assert.Contains(t, defaults, "logging.level")
	assert.Contains(t, defaults, "logging.format")
	assert.Contains(t, defaults, "logging.development")

	// Auth
	assert.Contains(t, defaults, "auth.jwt_secret")
	assert.Contains(t, defaults, "auth.jwt_expiry")
	assert.Contains(t, defaults, "auth.refresh_expiry")

	// Legacy
	assert.Contains(t, defaults, "legacy.enabled")
	assert.Contains(t, defaults, "legacy.encryption_key")
	assert.Contains(t, defaults, "legacy.privacy.require_pin")
	assert.Contains(t, defaults, "legacy.privacy.audit_all_access")
}

func TestDefaults_ValuesAreCorrectTypes(t *testing.T) {
	defaults := Defaults()

	// String types
	_, ok := defaults["server.host"].(string)
	assert.True(t, ok, "server.host should be string")

	_, ok = defaults["database.url"].(string)
	assert.True(t, ok, "database.url should be string")

	_, ok = defaults["logging.level"].(string)
	assert.True(t, ok, "logging.level should be string")

	// Int types
	_, ok = defaults["server.port"].(int)
	assert.True(t, ok, "server.port should be int")

	_, ok = defaults["jobs.max_workers"].(int)
	assert.True(t, ok, "jobs.max_workers should be int")

	// Bool types
	_, ok = defaults["cache.enabled"].(bool)
	assert.True(t, ok, "cache.enabled should be bool")

	_, ok = defaults["logging.development"].(bool)
	assert.True(t, ok, "logging.development should be bool")
}

func TestDefault_ReturnsFullConfig(t *testing.T) {
	cfg := Default()
	require.NotNil(t, cfg)

	// Server
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8096, cfg.Server.Port)
	assert.NotZero(t, cfg.Server.ReadTimeout)
	assert.NotZero(t, cfg.Server.WriteTimeout)
	assert.NotZero(t, cfg.Server.IdleTimeout)
	assert.NotZero(t, cfg.Server.ShutdownTimeout)

	// Database
	assert.Contains(t, cfg.Database.URL, "postgres://")
	assert.Equal(t, 0, cfg.Database.MaxConns) // Auto
	assert.Equal(t, 2, cfg.Database.MinConns)
	assert.NotZero(t, cfg.Database.MaxConnLifetime)
	assert.NotZero(t, cfg.Database.MaxConnIdleTime)
	assert.NotZero(t, cfg.Database.HealthCheckPeriod)

	// Cache (disabled by default)
	assert.False(t, cfg.Cache.Enabled)
	assert.Empty(t, cfg.Cache.URL)

	// Search (disabled by default)
	assert.False(t, cfg.Search.Enabled)
	assert.Empty(t, cfg.Search.URL)
	assert.Empty(t, cfg.Search.APIKey)

	// Jobs
	assert.Equal(t, 100, cfg.Jobs.MaxWorkers)
	assert.NotZero(t, cfg.Jobs.FetchCooldown)
	assert.NotZero(t, cfg.Jobs.FetchPollInterval)
	assert.NotZero(t, cfg.Jobs.RescueStuckJobsAfter)

	// Logging
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "text", cfg.Logging.Format)
	assert.False(t, cfg.Logging.Development)

	// Auth
	assert.Empty(t, cfg.Auth.JWTSecret) // Must be set by user
	assert.NotZero(t, cfg.Auth.JWTExpiry)
	assert.NotZero(t, cfg.Auth.RefreshExpiry)

	// Legacy
	assert.False(t, cfg.Legacy.Enabled)
	assert.Empty(t, cfg.Legacy.EncryptionKey)
	assert.True(t, cfg.Legacy.Privacy.RequirePIN)
	assert.True(t, cfg.Legacy.Privacy.AuditAllAccess)
}

func TestProvideConfig_Error(t *testing.T) {
	// ProvideConfig uses DefaultConfigPath which may not exist
	// But it should still return config with defaults
	cfg, err := ProvideConfig()
	require.NoError(t, err, "ProvideConfig should not error (uses defaults)")
	require.NotNil(t, cfg)
}

func TestLoad_NestedEnvVars(t *testing.T) {
	// Test that nested config keys work with env vars
	t.Setenv("REVENGE_LOGGING_LEVEL", "debug")
	t.Setenv("REVENGE_LOGGING_FORMAT", "json")
	t.Setenv("REVENGE_LOGGING_DEVELOPMENT", "true")

	cfg, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
	assert.True(t, cfg.Logging.Development)
}

func TestLoad_CacheConfig(t *testing.T) {
	t.Setenv("REVENGE_CACHE_ENABLED", "true")
	t.Setenv("REVENGE_CACHE_URL", "redis://localhost:6379")

	cfg, err := Load("")
	require.NoError(t, err)

	assert.True(t, cfg.Cache.Enabled)
	assert.Equal(t, "redis://localhost:6379", cfg.Cache.URL)
}

func TestLoad_SearchConfig(t *testing.T) {
	t.Setenv("REVENGE_SEARCH_ENABLED", "true")
	t.Setenv("REVENGE_SEARCH_URL", "http://localhost:8108")
	t.Setenv("REVENGE_SEARCH_API_KEY", "test-typesense-key")

	cfg, err := Load("")
	require.NoError(t, err)

	assert.True(t, cfg.Search.Enabled)
	assert.Equal(t, "http://localhost:8108", cfg.Search.URL)
	assert.Equal(t, "test-typesense-key", cfg.Search.APIKey)
}

func TestLoad_LegacyConfig(t *testing.T) {
	t.Setenv("REVENGE_LEGACY_ENABLED", "true")
	t.Setenv("REVENGE_LEGACY_ENCRYPTION_KEY", "test-encryption-key")

	cfg, err := Load("")
	require.NoError(t, err)

	assert.True(t, cfg.Legacy.Enabled)
	assert.Equal(t, "test-encryption-key", cfg.Legacy.EncryptionKey)
}

func TestLoad_CompoundEnvVarMapping(t *testing.T) {
	// Verify that compound field names (containing underscores) are mapped correctly.
	// The simple _→. transform would break these: SEARCH_API_KEY → search.api.key
	// instead of search.api_key. The applyCompoundEnvOverrides fix handles this.

	t.Run("search_api_key", func(t *testing.T) {
		t.Setenv("REVENGE_SEARCH_API_KEY", "ts-key-123")
		cfg, err := Load("")
		require.NoError(t, err)
		assert.Equal(t, "ts-key-123", cfg.Search.APIKey)
	})

	t.Run("radarr_api_key", func(t *testing.T) {
		t.Setenv("REVENGE_INTEGRATIONS_RADARR_API_KEY", "radarr-key-456")
		cfg, err := Load("")
		require.NoError(t, err)
		assert.Equal(t, "radarr-key-456", cfg.Integrations.Radarr.APIKey)
	})

	t.Run("sonarr_api_key", func(t *testing.T) {
		t.Setenv("REVENGE_INTEGRATIONS_SONARR_API_KEY", "sonarr-key-789")
		cfg, err := Load("")
		require.NoError(t, err)
		assert.Equal(t, "sonarr-key-789", cfg.Integrations.Sonarr.APIKey)
	})

	t.Run("radarr_base_url", func(t *testing.T) {
		t.Setenv("REVENGE_INTEGRATIONS_RADARR_BASE_URL", "http://radarr:7878")
		cfg, err := Load("")
		require.NoError(t, err)
		assert.Equal(t, "http://radarr:7878", cfg.Integrations.Radarr.BaseURL)
	})
}
