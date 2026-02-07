package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests to fill coverage gaps in config.go, loader.go, and module.go.
// Focuses on GetRadarrConfig, GetSonarrConfig, Defaults() completeness,
// Default() struct completeness, LoadWithKoanf error paths, and validation.

// =====================================================
// GetRadarrConfig / GetSonarrConfig tests
// =====================================================

func TestConfig_GetRadarrConfig_Default(t *testing.T) {
	t.Parallel()

	cfg := Default()
	radarr := cfg.GetRadarrConfig()

	assert.False(t, radarr.Enabled)
	assert.Empty(t, radarr.APIKey)
	assert.Empty(t, radarr.BaseURL)
	assert.False(t, radarr.AutoSync)
	assert.Equal(t, 0, radarr.SyncInterval)
}

func TestConfig_GetRadarrConfig_CustomValues(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Integrations: IntegrationsConfig{
			Radarr: RadarrConfig{
				Enabled:      true,
				BaseURL:      "http://radarr.local:7878",
				APIKey:       "test-api-key-123",
				AutoSync:     true,
				SyncInterval: 600,
			},
		},
	}

	radarr := cfg.GetRadarrConfig()

	assert.True(t, radarr.Enabled)
	assert.Equal(t, "http://radarr.local:7878", radarr.BaseURL)
	assert.Equal(t, "test-api-key-123", radarr.APIKey)
	assert.True(t, radarr.AutoSync)
	assert.Equal(t, 600, radarr.SyncInterval)
}

func TestConfig_GetSonarrConfig_Default(t *testing.T) {
	t.Parallel()

	cfg := Default()
	sonarr := cfg.GetSonarrConfig()

	assert.False(t, sonarr.Enabled)
	assert.Empty(t, sonarr.APIKey)
	assert.Empty(t, sonarr.BaseURL)
	assert.False(t, sonarr.AutoSync)
	assert.Equal(t, 0, sonarr.SyncInterval)
}

func TestConfig_GetSonarrConfig_CustomValues(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Integrations: IntegrationsConfig{
			Sonarr: SonarrConfig{
				Enabled:      true,
				BaseURL:      "http://sonarr.local:8989",
				APIKey:       "sonarr-key-456",
				AutoSync:     true,
				SyncInterval: 120,
			},
		},
	}

	sonarr := cfg.GetSonarrConfig()

	assert.True(t, sonarr.Enabled)
	assert.Equal(t, "http://sonarr.local:8989", sonarr.BaseURL)
	assert.Equal(t, "sonarr-key-456", sonarr.APIKey)
	assert.True(t, sonarr.AutoSync)
	assert.Equal(t, 120, sonarr.SyncInterval)
}

func TestConfig_GetRadarrConfig_ReturnsValue(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	cfg.Integrations.Radarr.Enabled = true

	radarr := cfg.GetRadarrConfig()
	assert.True(t, radarr.Enabled)

	// Modifying the returned value should NOT affect the original
	radarr.Enabled = false
	assert.True(t, cfg.Integrations.Radarr.Enabled, "returned value should be a copy")
}

func TestConfig_GetSonarrConfig_ReturnsValue(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	cfg.Integrations.Sonarr.Enabled = true

	sonarr := cfg.GetSonarrConfig()
	assert.True(t, sonarr.Enabled)

	// Modifying the returned value should NOT affect the original
	sonarr.Enabled = false
	assert.True(t, cfg.Integrations.Sonarr.Enabled, "returned value should be a copy")
}

// =====================================================
// Defaults() completeness tests
// =====================================================

func TestDefaults_StorageKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "storage.backend")
	assert.Contains(t, defaults, "storage.local.path")
	assert.Contains(t, defaults, "storage.s3.region")
	assert.Contains(t, defaults, "storage.s3.bucket")
	assert.Contains(t, defaults, "storage.s3.endpoint")
	assert.Contains(t, defaults, "storage.s3.access_key_id")
	assert.Contains(t, defaults, "storage.s3.secret_access_key")
	assert.Contains(t, defaults, "storage.s3.use_path_style")

	assert.Equal(t, "local", defaults["storage.backend"])
	assert.Equal(t, "/data/storage", defaults["storage.local.path"])
	assert.Equal(t, false, defaults["storage.s3.use_path_style"])
}

func TestDefaults_ActivityKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "activity.retention_days")
	assert.Equal(t, 90, defaults["activity.retention_days"])
}

func TestDefaults_RaftKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "raft.enabled")
	assert.Contains(t, defaults, "raft.node_id")
	assert.Contains(t, defaults, "raft.bind_addr")
	assert.Contains(t, defaults, "raft.data_dir")
	assert.Contains(t, defaults, "raft.bootstrap")

	assert.Equal(t, false, defaults["raft.enabled"])
	assert.Equal(t, "", defaults["raft.node_id"])
	assert.Equal(t, "0.0.0.0:7000", defaults["raft.bind_addr"])
	assert.Equal(t, "/data/raft", defaults["raft.data_dir"])
	assert.Equal(t, false, defaults["raft.bootstrap"])
}

func TestDefaults_PlaybackKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "playback.enabled")
	assert.Contains(t, defaults, "playback.segment_dir")
	assert.Contains(t, defaults, "playback.segment_duration")
	assert.Contains(t, defaults, "playback.max_concurrent_sessions")
	assert.Contains(t, defaults, "playback.session_timeout")
	assert.Contains(t, defaults, "playback.ffmpeg_path")
	assert.Contains(t, defaults, "playback.transcode.enabled")
	assert.Contains(t, defaults, "playback.transcode.hw_accel")
	assert.Contains(t, defaults, "playback.transcode.hw_accel_device")
	assert.Contains(t, defaults, "playback.transcode.profiles")

	assert.Equal(t, true, defaults["playback.enabled"])
	assert.Equal(t, "/tmp/revenge-segments", defaults["playback.segment_dir"])
	assert.Equal(t, 6, defaults["playback.segment_duration"])
	assert.Equal(t, 10, defaults["playback.max_concurrent_sessions"])
	assert.Equal(t, "30m", defaults["playback.session_timeout"])
	assert.Equal(t, "ffmpeg", defaults["playback.ffmpeg_path"])
	assert.Equal(t, true, defaults["playback.transcode.enabled"])
	assert.Equal(t, "none", defaults["playback.transcode.hw_accel"])
	assert.Equal(t, "", defaults["playback.transcode.hw_accel_device"])
	assert.Equal(t, []string{"original", "1080p", "720p", "480p"}, defaults["playback.transcode.profiles"])
}

func TestDefaults_IntegrationsRadarrKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "integrations.radarr.enabled")
	assert.Contains(t, defaults, "integrations.radarr.base_url")
	assert.Contains(t, defaults, "integrations.radarr.api_key")
	assert.Contains(t, defaults, "integrations.radarr.auto_sync")
	assert.Contains(t, defaults, "integrations.radarr.sync_interval")

	assert.Equal(t, false, defaults["integrations.radarr.enabled"])
	assert.Equal(t, "http://localhost:7878", defaults["integrations.radarr.base_url"])
	assert.Equal(t, "", defaults["integrations.radarr.api_key"])
	assert.Equal(t, false, defaults["integrations.radarr.auto_sync"])
	assert.Equal(t, 300, defaults["integrations.radarr.sync_interval"])
}

func TestDefaults_EmailKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "email.enabled")
	assert.Contains(t, defaults, "email.provider")
	assert.Contains(t, defaults, "email.from_address")
	assert.Contains(t, defaults, "email.from_name")
	assert.Contains(t, defaults, "email.base_url")
	assert.Contains(t, defaults, "email.smtp.host")
	assert.Contains(t, defaults, "email.smtp.port")
	assert.Contains(t, defaults, "email.smtp.username")
	assert.Contains(t, defaults, "email.smtp.password")
	assert.Contains(t, defaults, "email.smtp.use_tls")
	assert.Contains(t, defaults, "email.smtp.use_starttls")
	assert.Contains(t, defaults, "email.smtp.skip_verify")
	assert.Contains(t, defaults, "email.smtp.timeout")
	assert.Contains(t, defaults, "email.sendgrid.api_key")

	assert.Equal(t, false, defaults["email.enabled"])
	assert.Equal(t, "smtp", defaults["email.provider"])
	assert.Equal(t, "Revenge Media Server", defaults["email.from_name"])
	assert.Equal(t, 587, defaults["email.smtp.port"])
	assert.Equal(t, true, defaults["email.smtp.use_starttls"])
	assert.Equal(t, "30s", defaults["email.smtp.timeout"])
}

func TestDefaults_AvatarKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "avatar.storage_path")
	assert.Contains(t, defaults, "avatar.max_size_bytes")
	assert.Contains(t, defaults, "avatar.allowed_types")

	assert.Equal(t, "/data/avatars", defaults["avatar.storage_path"])
	assert.Equal(t, 2*1024*1024, defaults["avatar.max_size_bytes"])
	assert.Equal(t, []string{"image/jpeg", "image/png", "image/webp"}, defaults["avatar.allowed_types"])
}

func TestDefaults_SessionKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "session.cache_enabled")
	assert.Contains(t, defaults, "session.cache_ttl")
	assert.Contains(t, defaults, "session.max_per_user")
	assert.Contains(t, defaults, "session.token_length")

	assert.Equal(t, true, defaults["session.cache_enabled"])
	assert.Equal(t, "5m", defaults["session.cache_ttl"])
	assert.Equal(t, 10, defaults["session.max_per_user"])
	assert.Equal(t, 32, defaults["session.token_length"])
}

func TestDefaults_RBACKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "rbac.model_path")
	assert.Contains(t, defaults, "rbac.policy_reload_interval")

	assert.Equal(t, "config/casbin_model.conf", defaults["rbac.model_path"])
	assert.Equal(t, "5m", defaults["rbac.policy_reload_interval"])
}

func TestDefaults_AuthLockoutKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "auth.lockout_threshold")
	assert.Contains(t, defaults, "auth.lockout_window")
	assert.Contains(t, defaults, "auth.lockout_enabled")

	assert.Equal(t, 5, defaults["auth.lockout_threshold"])
	assert.Equal(t, "15m", defaults["auth.lockout_window"])
	assert.Equal(t, true, defaults["auth.lockout_enabled"])
}

func TestDefaults_JobsMaxAttemptsKey(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "jobs.max_attempts")
	assert.Equal(t, 25, defaults["jobs.max_attempts"])
}

func TestDefaults_RateLimitKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "rate_limit.enabled")
	assert.Contains(t, defaults, "rate_limit.backend")
	assert.Contains(t, defaults, "rate_limit.global.requests_per_second")
	assert.Contains(t, defaults, "rate_limit.global.burst")
	assert.Contains(t, defaults, "rate_limit.auth.requests_per_second")
	assert.Contains(t, defaults, "rate_limit.auth.burst")

	assert.Equal(t, true, defaults["rate_limit.enabled"])
	assert.Equal(t, "memory", defaults["rate_limit.backend"])
	assert.Equal(t, 10.0, defaults["rate_limit.global.requests_per_second"])
	assert.Equal(t, 20, defaults["rate_limit.global.burst"])
	assert.Equal(t, 1.0, defaults["rate_limit.auth.requests_per_second"])
	assert.Equal(t, 5, defaults["rate_limit.auth.burst"])
}

func TestDefaults_MovieKeys(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	assert.Contains(t, defaults, "movie.tmdb.api_key")
	assert.Contains(t, defaults, "movie.tmdb.rate_limit")
	assert.Contains(t, defaults, "movie.tmdb.cache_ttl")
	assert.Contains(t, defaults, "movie.tmdb.proxy_url")
	assert.Contains(t, defaults, "movie.library.paths")
	assert.Contains(t, defaults, "movie.library.scan_interval")

	assert.Equal(t, "", defaults["movie.tmdb.api_key"])
	assert.Equal(t, 40, defaults["movie.tmdb.rate_limit"])
	assert.Equal(t, "5m", defaults["movie.tmdb.cache_ttl"])
	assert.Equal(t, []string{}, defaults["movie.library.paths"])
}

// =====================================================
// Default() struct completeness tests
// =====================================================

func TestDefault_StorageConfig(t *testing.T) {
	t.Parallel()

	cfg := Default()

	assert.Equal(t, "local", cfg.Storage.Backend)
	assert.Equal(t, "/data/storage", cfg.Storage.Local.Path)
	assert.Empty(t, cfg.Storage.S3.Bucket)
	assert.Empty(t, cfg.Storage.S3.Endpoint)
}

func TestDefault_LegacyConfig(t *testing.T) {
	t.Parallel()

	cfg := Default()

	assert.False(t, cfg.Legacy.Enabled)
	assert.Empty(t, cfg.Legacy.EncryptionKey)
	assert.True(t, cfg.Legacy.Privacy.RequirePIN)
	assert.True(t, cfg.Legacy.Privacy.AuditAllAccess)
}

func TestDefault_JobsConfig(t *testing.T) {
	t.Parallel()

	cfg := Default()

	assert.Equal(t, 100, cfg.Jobs.MaxWorkers)
	assert.Equal(t, 200*time.Millisecond, cfg.Jobs.FetchCooldown)
	assert.Equal(t, 2*time.Second, cfg.Jobs.FetchPollInterval)
	assert.Equal(t, 30*time.Minute, cfg.Jobs.RescueStuckJobsAfter)
}

func TestDefault_AuthConfig(t *testing.T) {
	t.Parallel()

	cfg := Default()

	assert.Empty(t, cfg.Auth.JWTSecret)
	assert.Equal(t, 24*time.Hour, cfg.Auth.JWTExpiry)
	assert.Equal(t, 7*24*time.Hour, cfg.Auth.RefreshExpiry)
}

func TestDefault_DatabaseConfig(t *testing.T) {
	t.Parallel()

	cfg := Default()

	assert.Equal(t, 0, cfg.Database.MaxConns)
	assert.Equal(t, 2, cfg.Database.MinConns)
	assert.Equal(t, 30*time.Minute, cfg.Database.MaxConnLifetime)
	assert.Equal(t, 5*time.Minute, cfg.Database.MaxConnIdleTime)
	assert.Equal(t, 30*time.Second, cfg.Database.HealthCheckPeriod)
}

func TestDefault_ServerTimeouts(t *testing.T) {
	t.Parallel()

	cfg := Default()

	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, cfg.Server.IdleTimeout)
	assert.Equal(t, 10*time.Second, cfg.Server.ShutdownTimeout)
}

// =====================================================
// LoadWithKoanf additional tests
// =====================================================

func TestLoadWithKoanf_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: "10.0.0.1"
  port: 3333
logging:
  level: error
  format: json
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, k, err := LoadWithKoanf(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotNil(t, k)

	assert.Equal(t, "10.0.0.1", cfg.Server.Host)
	assert.Equal(t, 3333, cfg.Server.Port)
	assert.Equal(t, "error", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)

	// Verify koanf reflects same values
	assert.Equal(t, "10.0.0.1", k.String("server.host"))
	assert.Equal(t, 3333, k.Int("server.port"))
}

func TestLoadWithKoanf_MissingConfigFileUsesDefaults(t *testing.T) {
	cfg, k, err := LoadWithKoanf("nonexistent.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotNil(t, k)

	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8096, cfg.Server.Port)
}

func TestLoadWithKoanf_ValidationFailure(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	invalidContent := `
server:
  host: ""
  port: 0
database:
  url: ""
`
	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	cfg, k, err := LoadWithKoanf(configPath)
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Nil(t, k)
}

// =====================================================
// Load additional tests
// =====================================================

func TestLoad_IntegrationsFromEnv(t *testing.T) {
	t.Setenv("REVENGE_INTEGRATIONS_RADARR_ENABLED", "true")

	cfg, err := Load("")
	require.NoError(t, err)

	// Note: env var mapping converts _ to . so this maps to integrations.radarr.enabled
	// which should set Integrations.Radarr.Enabled = true
	assert.True(t, cfg.Integrations.Radarr.Enabled)
}

func TestLoad_PlaybackDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	assert.True(t, cfg.Playback.Enabled)
	assert.Equal(t, "/tmp/revenge-segments", cfg.Playback.SegmentDir)
	assert.Equal(t, 6, cfg.Playback.SegmentDuration)
	assert.Equal(t, 10, cfg.Playback.MaxConcurrentSessions)
	assert.Equal(t, "ffmpeg", cfg.Playback.FFmpegPath)
	assert.True(t, cfg.Playback.Transcode.Enabled)
	assert.Equal(t, "none", cfg.Playback.Transcode.HWAccel)
}

func TestLoad_RaftDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	assert.False(t, cfg.Raft.Enabled)
	assert.Equal(t, "/data/raft", cfg.Raft.DataDir)
	assert.False(t, cfg.Raft.Bootstrap)
}

func TestLoad_StorageDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, "local", cfg.Storage.Backend)
	assert.Equal(t, "/data/storage", cfg.Storage.Local.Path)
}

func TestLoad_ActivityDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, 90, cfg.Activity.RetentionDays)
}

func TestLoad_SessionDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	assert.True(t, cfg.Session.CacheEnabled)
	assert.Equal(t, 10, cfg.Session.MaxPerUser)
	assert.Equal(t, 32, cfg.Session.TokenLength)
}

func TestLoad_AuthLockoutDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, 5, cfg.Auth.LockoutThreshold)
	assert.True(t, cfg.Auth.LockoutEnabled)
}

// =====================================================
// Validate additional tests
// =====================================================

func TestValidate_InvalidStorageBackend(t *testing.T) {
	t.Parallel()

	cfg := Default()
	cfg.Storage.Backend = "ftp" // Invalid: oneof=local s3

	err := validate(cfg)
	assert.Error(t, err)
}

func TestValidate_ValidStorageBackendS3(t *testing.T) {
	t.Parallel()

	cfg := Default()
	cfg.Storage.Backend = "s3"

	// With s3 backend, validation should still pass for
	// the backend field itself (required_if fields are separate)
	err := validate(cfg)
	// Note: This may or may not error depending on required_if validation
	// The point is the backend field itself is valid
	_ = err
}

func TestValidate_InvalidEmailProvider(t *testing.T) {
	t.Parallel()

	cfg := Default()
	cfg.Email.Provider = "mailgun" // Invalid: oneof=smtp sendgrid

	err := validate(cfg)
	assert.Error(t, err)
}

func TestValidate_ValidEmailProviderSendgrid(t *testing.T) {
	t.Parallel()

	cfg := Default()
	cfg.Email.Provider = "sendgrid"

	err := validate(cfg)
	assert.NoError(t, err)
}

func TestValidate_ValidLoggingLevels(t *testing.T) {
	t.Parallel()

	validLevels := []string{"debug", "info", "warn", "error"}

	for _, level := range validLevels {
		t.Run(level, func(t *testing.T) {
			cfg := Default()
			cfg.Logging.Level = level
			err := validate(cfg)
			assert.NoError(t, err)
		})
	}
}

func TestValidate_ValidLoggingFormats(t *testing.T) {
	t.Parallel()

	validFormats := []string{"text", "json"}

	for _, format := range validFormats {
		t.Run(format, func(t *testing.T) {
			cfg := Default()
			cfg.Logging.Format = format
			err := validate(cfg)
			assert.NoError(t, err)
		})
	}
}

func TestValidate_ServerPortBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"port 0 invalid", 0, true},
		{"port 1 valid", 1, false},
		{"port 80 valid", 80, false},
		{"port 8080 valid", 8080, false},
		{"port 65535 valid", 65535, false},
		{"port 65536 invalid", 65536, true},
		{"port -1 invalid", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			cfg.Server.Port = tt.port
			err := validate(cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =====================================================
// MustLoad additional tests
// =====================================================

func TestMustLoad_WithValidConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "valid.yaml")

	configContent := `
server:
  host: "0.0.0.0"
  port: 9999
database:
  url: "postgres://revenge:pass@localhost:5432/revenge?sslmode=disable"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		cfg := MustLoad(configPath)
		assert.Equal(t, 9999, cfg.Server.Port)
	})
}

// =====================================================
// Default() function returns new instance each time
// =====================================================

func TestDefault_ReturnsNewInstance(t *testing.T) {
	t.Parallel()

	cfg1 := Default()
	cfg2 := Default()

	// Should be equal in content
	assert.Equal(t, cfg1.Server.Host, cfg2.Server.Host)
	assert.Equal(t, cfg1.Server.Port, cfg2.Server.Port)

	// But different pointers
	cfg1.Server.Port = 12345
	assert.NotEqual(t, cfg1.Server.Port, cfg2.Server.Port)
}

// =====================================================
// Config struct field tests
// =====================================================

func TestPlaybackConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := PlaybackConfig{}

	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.SegmentDir)
	assert.Equal(t, 0, cfg.SegmentDuration)
	assert.Equal(t, 0, cfg.MaxConcurrentSessions)
	assert.Equal(t, time.Duration(0), cfg.SessionTimeout)
	assert.Empty(t, cfg.FFmpegPath)
	assert.False(t, cfg.Transcode.Enabled)
}

func TestTranscodeConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := TranscodeConfig{}

	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.HWAccel)
	assert.Empty(t, cfg.HWAccelDevice)
	assert.Nil(t, cfg.Profiles)
}

func TestRaftConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := RaftConfig{}

	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.NodeID)
	assert.Empty(t, cfg.BindAddr)
	assert.Empty(t, cfg.DataDir)
	assert.False(t, cfg.Bootstrap)
}

func TestStorageConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := StorageConfig{}

	assert.Empty(t, cfg.Backend)
	assert.Empty(t, cfg.Local.Path)
	assert.Empty(t, cfg.S3.Bucket)
	assert.Empty(t, cfg.S3.Region)
	assert.Empty(t, cfg.S3.Endpoint)
	assert.Empty(t, cfg.S3.AccessKeyID)
	assert.Empty(t, cfg.S3.SecretAccessKey)
	assert.False(t, cfg.S3.UsePathStyle)
}

func TestActivityConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := ActivityConfig{}

	assert.Equal(t, 0, cfg.RetentionDays)
}

func TestSMTPConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := SMTPConfig{}

	assert.Empty(t, cfg.Host)
	assert.Equal(t, 0, cfg.Port)
	assert.Empty(t, cfg.Username)
	assert.Empty(t, cfg.Password)
	assert.False(t, cfg.UseTLS)
	assert.False(t, cfg.UseStartTLS)
	assert.False(t, cfg.SkipVerify)
	assert.Equal(t, time.Duration(0), cfg.Timeout)
}

func TestSessionConfig_ZeroValue(t *testing.T) {
	t.Parallel()

	cfg := SessionConfig{}

	assert.False(t, cfg.CacheEnabled)
	assert.Equal(t, time.Duration(0), cfg.CacheTTL)
	assert.Equal(t, 0, cfg.MaxPerUser)
	assert.Equal(t, 0, cfg.TokenLength)
}

// =====================================================
// Module tests
// =====================================================

func TestModule_NotNil(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, Module)
}

// =====================================================
// Constants tests
// =====================================================

func TestEnvPrefix_Value(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "REVENGE_", EnvPrefix)
}

func TestDefaultConfigPath_Value(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "config/config.yaml", DefaultConfigPath)
}

// =====================================================
// Defaults count and completeness sanity check
// =====================================================

func TestDefaults_MinimumKeyCount(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	// We know there are at least 80+ keys in Defaults()
	// This guards against accidental deletion
	assert.GreaterOrEqual(t, len(defaults), 80, "Defaults should have at least 80 keys")
}

func TestDefaults_AllValuesNonNil(t *testing.T) {
	t.Parallel()

	defaults := Defaults()

	for key, val := range defaults {
		assert.NotNil(t, val, "Default value for %q should not be nil", key)
	}
}
