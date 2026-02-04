package search_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/search"
)

func TestModule(t *testing.T) {
	// Test that the module can be created
	assert.NotNil(t, search.Module)

	// Test that module has expected options
	app := fx.New(
		search.Module,
		fx.NopLogger,
	)

	assert.NotNil(t, app)
}

func TestNewClient_Disabled(t *testing.T) {
	cfg := &config.Config{
		Search: config.SearchConfig{
			Enabled: false,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	client, err := search.NewClient(cfg, logger)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.False(t, client.IsEnabled())
}

func TestNewClient_EnabledEmptyURL(t *testing.T) {
	cfg := &config.Config{
		Search: config.SearchConfig{
			Enabled: true,
			URL:     "",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	_, err := search.NewClient(cfg, logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "search URL is required")
}

func TestNewClient_EnabledWithURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"HTTP URL with port", "http://localhost:8108"},
		{"HTTPS URL with port", "https://search.example.com:443"},
		{"URL without scheme", "localhost:8108"},
		{"URL with just host", "localhost"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Search: config.SearchConfig{
					Enabled: true,
					URL:     tt.url,
					APIKey:  "test-api-key",
				},
			}
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

			client, err := search.NewClient(cfg, logger)
			require.NoError(t, err)
			assert.NotNil(t, client)
			assert.True(t, client.IsEnabled())
		})
	}
}

func TestClient_DisabledOperations(t *testing.T) {
	cfg := &config.Config{
		Search: config.SearchConfig{
			Enabled: false,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	client, err := search.NewClient(cfg, logger)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("CreateCollection", func(t *testing.T) {
		err := client.CreateCollection(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("DeleteCollection", func(t *testing.T) {
		err := client.DeleteCollection(ctx, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("GetCollection", func(t *testing.T) {
		_, err := client.GetCollection(ctx, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("ListCollections", func(t *testing.T) {
		_, err := client.ListCollections(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("IndexDocument", func(t *testing.T) {
		_, err := client.IndexDocument(ctx, "test", map[string]interface{}{"id": "1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("UpdateDocument", func(t *testing.T) {
		_, err := client.UpdateDocument(ctx, "test", "1", map[string]interface{}{"id": "1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("DeleteDocument", func(t *testing.T) {
		_, err := client.DeleteDocument(ctx, "test", "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("Search", func(t *testing.T) {
		_, err := client.Search(ctx, "test", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("MultiSearch", func(t *testing.T) {
		_, err := client.MultiSearch(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("ImportDocuments", func(t *testing.T) {
		_, err := client.ImportDocuments(ctx, "test", nil, "create")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "search is disabled")
	})

	t.Run("HealthCheck disabled returns nil", func(t *testing.T) {
		err := client.HealthCheck(ctx)
		assert.NoError(t, err) // No error when disabled
	})
}

func TestClient_IsEnabled(t *testing.T) {
	t.Run("disabled config", func(t *testing.T) {
		cfg := &config.Config{
			Search: config.SearchConfig{
				Enabled: false,
			},
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		client, err := search.NewClient(cfg, logger)
		require.NoError(t, err)
		assert.False(t, client.IsEnabled())
	})

	t.Run("enabled config", func(t *testing.T) {
		cfg := &config.Config{
			Search: config.SearchConfig{
				Enabled: true,
				URL:     "http://localhost:8108",
				APIKey:  "test",
			},
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		client, err := search.NewClient(cfg, logger)
		require.NoError(t, err)
		assert.True(t, client.IsEnabled())
	})
}
