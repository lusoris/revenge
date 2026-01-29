// Package search provides a Typesense search client.
package search

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)


// Client wraps the Typesense client for search operations.
type Client struct {
	ts     *typesense.Client
	logger *slog.Logger
}

// NewClient creates a new search client.
func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	serverURL := fmt.Sprintf("http://%s:%d", cfg.Search.Host, cfg.Search.Port)
	ts := typesense.NewClient(
		typesense.WithServer(serverURL),
		typesense.WithAPIKey(cfg.Search.APIKey),
	)

	return &Client{
		ts:     ts,
		logger: logger.With(slog.String("component", "search")),
	}, nil
}

// Health checks if the search service is healthy.
func (c *Client) Health(ctx context.Context) error {
	_, err := c.ts.Health(ctx, 5)
	return err
}

// CreateCollection creates a new search collection.
func (c *Client) CreateCollection(ctx context.Context, schema *api.CollectionSchema) (*api.CollectionResponse, error) {
	return c.ts.Collections().Create(ctx, schema)
}

// Index indexes a document in a collection.
func (c *Client) Index(ctx context.Context, collection string, document any) (map[string]any, error) {
	return c.ts.Collection(collection).Documents().Create(ctx, document, nil)
}

// Upsert upserts a document in a collection.
func (c *Client) Upsert(ctx context.Context, collection string, document any) (map[string]any, error) {
	return c.ts.Collection(collection).Documents().Upsert(ctx, document, nil)
}

// Search performs a search query.
func (c *Client) Search(ctx context.Context, collection string, params *api.SearchCollectionParams) (*api.SearchResult, error) {
	return c.ts.Collection(collection).Documents().Search(ctx, params)
}

// Delete removes a document from a collection.
func (c *Client) Delete(ctx context.Context, collection string, id string) (map[string]any, error) {
	return c.ts.Collection(collection).Document(id).Delete(ctx)
}

// Module provides search dependencies for fx.
var Module = fx.Module("search",
	fx.Provide(NewClient),
)
