package config

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	if cfg.Server.Host == "" {
		t.Error("Server.Host should have a default value")
	}

	if cfg.Server.Port == 0 {
		t.Error("Server.Port should have a default value")
	}

	if cfg.Database.URL == "" {
		t.Error("Database.URL should have a default value")
	}
}

func TestDefaultServerConfig(t *testing.T) {
	cfg := Default()

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
}
