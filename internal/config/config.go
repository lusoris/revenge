// Package config provides configuration management for the Revenge server.
package config

// Config holds the application configuration.
type Config struct {
	// Server configuration
	Server ServerConfig

	// Database configuration
	Database DatabaseConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	// Host is the server bind address.
	Host string

	// Port is the server port.
	Port int
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	// URL is the PostgreSQL connection string.
	URL string
}

// Default returns a Config with default values.
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			URL: "postgres://localhost:5432/revenge?sslmode=disable",
		},
	}
}
