package search

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"github.com/typesense/typesense-go/v2/typesense"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"go.uber.org/fx"
)

// Module provides search dependencies.
var Module = fx.Module("search",
	fx.Provide(NewClient),
	fx.Invoke(registerHooks),
)

// Client represents the Typesense search client.
type Client struct {
	client *typesense.Client
	config *config.Config
	logger *slog.Logger
}

// NewClient creates a new Typesense search client.
func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	if !cfg.Search.Enabled {
		logger.Info("search disabled, returning nil client")
		return &Client{
			client: nil,
			config: cfg,
			logger: logger,
		}, nil
	}

	// Parse search URL to extract host, port, and protocol
	url := strings.TrimSpace(cfg.Search.URL)
	if url == "" {
		return nil, fmt.Errorf("search URL is required when search is enabled")
	}

	// Parse URL components
	protocol := "http"
	var host string
	port := "8108"

	// Simple URL parsing
	if strings.HasPrefix(url, "https://") {
		protocol = "https"
		url = strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "http://") {
		protocol = "http"
		url = strings.TrimPrefix(url, "http://")
	}

	// Extract host:port
	if strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		host = parts[0]
		port = parts[1]
	} else {
		host = url
	}

	logger.Info("initializing typesense client",
		"host", host,
		"port", port,
		"protocol", protocol,
	)

	// Create Typesense client
	// Note: Remove circuit breaker for now as it might be causing timeout issues
	client := typesense.NewClient(
		typesense.WithServer(protocol+"://"+host+":"+port),
		typesense.WithAPIKey(cfg.Search.APIKey),
	)

	return &Client{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// IsEnabled returns true if search is enabled and client is initialized.
func (c *Client) IsEnabled() bool {
	return c.client != nil && c.config.Search.Enabled
}

// CreateCollection creates a new collection/index in Typesense.
func (c *Client) CreateCollection(ctx context.Context, schema *api.CollectionSchema) error {
	if !c.IsEnabled() {
		return fmt.Errorf("search is disabled")
	}

	_, err := c.client.Collections().Create(ctx, schema)
	return err
}

// DeleteCollection deletes a collection from Typesense.
func (c *Client) DeleteCollection(ctx context.Context, name string) error {
	if !c.IsEnabled() {
		return fmt.Errorf("search is disabled")
	}

	_, err := c.client.Collection(name).Delete(ctx)
	return err
}

// GetCollection retrieves collection information.
func (c *Client) GetCollection(ctx context.Context, name string) (*api.CollectionResponse, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.Collection(name).Retrieve(ctx)
}

// ListCollections returns all collections.
func (c *Client) ListCollections(ctx context.Context) ([]*api.CollectionResponse, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.Collections().Retrieve(ctx)
}

// IndexDocument indexes a single document in a collection.
func (c *Client) IndexDocument(ctx context.Context, collectionName string, document interface{}) (map[string]interface{}, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.Collection(collectionName).Documents().Create(ctx, document)
}

// UpdateDocument updates an existing document.
func (c *Client) UpdateDocument(ctx context.Context, collectionName, documentID string, document interface{}) (map[string]interface{}, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.Collection(collectionName).Document(documentID).Update(ctx, document)
}

// DeleteDocument deletes a document from a collection.
func (c *Client) DeleteDocument(ctx context.Context, collectionName, documentID string) (map[string]interface{}, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.Collection(collectionName).Document(documentID).Delete(ctx)
}

// Search performs a search query against a collection.
func (c *Client) Search(ctx context.Context, collectionName string, params *api.SearchCollectionParams) (*api.SearchResult, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.Collection(collectionName).Documents().Search(ctx, params)
}

// MultiSearch performs a multi-collection search.
func (c *Client) MultiSearch(ctx context.Context, params *api.MultiSearchParams) (*api.MultiSearchResult, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	return c.client.MultiSearch.Perform(ctx, params, api.MultiSearchSearchesParameter{})
}

// ImportDocuments bulk imports documents into a collection.
func (c *Client) ImportDocuments(ctx context.Context, collectionName string, documents []interface{}, action string) ([]*api.ImportDocumentResponse, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("search is disabled")
	}

	batchSize := 100
	params := &api.ImportDocumentsParams{
		Action:    &action,
		BatchSize: &batchSize,
	}

	return c.client.Collection(collectionName).Documents().Import(ctx, documents, params)
}

// HealthCheck performs a health check on the Typesense server.
func (c *Client) HealthCheck(ctx context.Context) error {
	if !c.IsEnabled() {
		return nil // Not an error if disabled
	}

	healthy, err := c.client.Health(ctx, 10*time.Second)
	if err != nil {
		return err
	}
	if !healthy {
		return fmt.Errorf("typesense server reports unhealthy")
	}
	return nil
}

// registerHooks registers lifecycle hooks for the search client.
func registerHooks(lc fx.Lifecycle, client *Client, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if !client.IsEnabled() {
				logger.Info("search disabled, skipping startup")
				return nil
			}

			// Retry health check on startup with backoff
			maxRetries := 5
			for i := 0; i < maxRetries; i++ {
				if err := client.HealthCheck(ctx); err != nil {
					if i < maxRetries-1 {
						logger.Debug("typesense health check failed, retrying...",
							"attempt", i+1,
							"max_retries", maxRetries,
							"error", err)
						time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
						continue
					}
					logger.Warn("typesense health check failed after retries", "error", err)
					// Don't fail startup, just log warning
					return nil
				}
				logger.Info("typesense client connected successfully")
				return nil
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !client.IsEnabled() {
				return nil
			}
			logger.Info("typesense client stopped")
			return nil
		},
	})
}
