// Package search provides integration tests for Typesense search functionality
//go:build integration
// +build integration

package search

import (
	"context"
	"fmt"
	"testing"
	"time"

	"log/slog"
	"os"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

// Helper functions for creating pointers
func ptr[T any](v T) *T {
	return &v
}

func newTestClient(t *testing.T) *search.Client {
	cfg := &config.Config{
		Search: config.SearchConfig{
			Enabled: true,
			URL:     "http://localhost:8108",
			APIKey:  "dev_api_key",
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	client, err := search.NewClient(cfg, logger)
	require.NoError(t, err, "should create search client")
	require.NotNil(t, client, "client should not be nil")

	return client
}

func TestSearchClientConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client := newTestClient(t)

	// Test health check with retry logic since it might take a moment
	var err error
	for i := 0; i < 3; i++ {
		err = client.HealthCheck(ctx)
		if err == nil {
			break
		}
		if i < 2 {
			time.Sleep(2 * time.Second)
		}
	}
	assert.NoError(t, err, "health check should succeed")
}

func TestSearchCollectionLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := newTestClient(t)

	collectionName := fmt.Sprintf("test_collection_%d", time.Now().UnixNano())

	// Create collection
	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "description", Type: "string"},
			{Name: "year", Type: "int32", Facet: ptr(true)},
		},
	}

	err := client.CreateCollection(ctx, schema)
	require.NoError(t, err, "should create collection")

	// Verify collection exists
	retrieved, err := client.GetCollection(ctx, collectionName)
	require.NoError(t, err, "should retrieve collection")
	assert.Equal(t, collectionName, retrieved.Name)

	// List collections
	collections, err := client.ListCollections(ctx)
	require.NoError(t, err, "should list collections")

	found := false
	for _, c := range collections {
		if c.Name == collectionName {
			found = true
			break
		}
	}
	assert.True(t, found, "created collection should be in list")

	// Clean up
	err = client.DeleteCollection(ctx, collectionName)
	assert.NoError(t, err, "should delete collection")
}

func TestSearchDocumentOperations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := newTestClient(t)

	collectionName := fmt.Sprintf("test_docs_%d", time.Now().UnixNano())

	// Create collection
	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "content", Type: "string"},
		},
	}

	err := client.CreateCollection(ctx, schema)
	require.NoError(t, err)

	defer func() {
		_ = client.DeleteCollection(ctx, collectionName)
	}()

	// Index a document
	doc := map[string]interface{}{
		"id":      "1",
		"title":   "Test Document",
		"content": "This is a test document for search",
	}

	_, err = client.IndexDocument(ctx, collectionName, doc)
	require.NoError(t, err, "should index document")

	// Wait for indexing
	time.Sleep(100 * time.Millisecond)

	// Search for the document
	searchParams := &api.SearchCollectionParams{
		Q:       ptr("test"),
		QueryBy: ptr("title,content"),
	}

	results, err := client.Search(ctx, collectionName, searchParams)
	require.NoError(t, err, "should search successfully")
	assert.NotNil(t, results, "results should not be nil")
	assert.Greater(t, int(*results.Found), 0, "should find at least one document")

	// Update document
	updatedDoc := map[string]interface{}{
		"title":   "Updated Test Document",
		"content": "This content has been updated",
	}

	_, err = client.UpdateDocument(ctx, collectionName, "1", updatedDoc)
	require.NoError(t, err, "should update document")

	// Delete document
	_, err = client.DeleteDocument(ctx, collectionName, "1")
	assert.NoError(t, err, "should delete document")
}

func TestSearchBulkImport(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := newTestClient(t)

	collectionName := fmt.Sprintf("test_bulk_%d", time.Now().UnixNano())

	// Create collection
	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "rating", Type: "float"},
		},
	}

	err := client.CreateCollection(ctx, schema)
	require.NoError(t, err)

	defer func() {
		_ = client.DeleteCollection(ctx, collectionName)
	}()

	// Bulk import documents
	documents := []interface{}{
		map[string]interface{}{"id": "1", "title": "First Movie", "rating": 8.5},
		map[string]interface{}{"id": "2", "title": "Second Movie", "rating": 7.2},
		map[string]interface{}{"id": "3", "title": "Third Movie", "rating": 9.0},
	}

	results, err := client.ImportDocuments(ctx, collectionName, documents, "create")
	require.NoError(t, err, "should import documents")
	assert.Len(t, results, 3, "should return 3 import results")

	// Verify all documents were imported successfully
	for i, result := range results {
		assert.True(t, result.Success, "document %d should import successfully", i)
	}

	// Wait for indexing
	time.Sleep(200 * time.Millisecond)

	// Search to verify
	searchParams := &api.SearchCollectionParams{
		Q:       ptr("movie"),
		QueryBy: ptr("title"),
	}

	searchResults, err := client.Search(ctx, collectionName, searchParams)
	require.NoError(t, err, "should search successfully")
	assert.Equal(t, 3, int(*searchResults.Found), "should find all 3 movies")
}

func TestSearchErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := newTestClient(t)

	// Try to search non-existent collection
	searchParams := &api.SearchCollectionParams{
		Q:       ptr("query"),
		QueryBy: ptr("title"),
	}

	_, err := client.Search(ctx, "nonexistent_collection_12345", searchParams)
	assert.Error(t, err, "should fail on non-existent collection")

	// Try to get non-existent collection
	_, err = client.GetCollection(ctx, "nonexistent_collection_12345")
	assert.Error(t, err, "should fail on non-existent collection")

	// Try to delete non-existent collection
	err = client.DeleteCollection(ctx, "nonexistent_collection_12345")
	assert.Error(t, err, "should fail on non-existent collection")
}

func TestSearchWithFiltersAndSorting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := newTestClient(t)

	collectionName := fmt.Sprintf("test_filters_%d", time.Now().UnixNano())

	// Create collection with faceted fields
	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "year", Type: "int32", Facet: ptr(true)},
			{Name: "rating", Type: "float"},
		},
		DefaultSortingField: ptr("rating"),
	}

	err := client.CreateCollection(ctx, schema)
	require.NoError(t, err)

	defer func() {
		_ = client.DeleteCollection(ctx, collectionName)
	}()

	// Import test data
	documents := []interface{}{
		map[string]interface{}{"id": "1", "title": "Old Movie", "year": 1990, "rating": 7.5},
		map[string]interface{}{"id": "2", "title": "Recent Movie", "year": 2020, "rating": 8.5},
		map[string]interface{}{"id": "3", "title": "New Movie", "year": 2023, "rating": 9.0},
	}

	_, err = client.ImportDocuments(ctx, collectionName, documents, "create")
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	// Search with filter (year > 2000)
	searchParams := &api.SearchCollectionParams{
		Q:        ptr("*"),
		QueryBy:  ptr("title"),
		FilterBy: ptr("year:>2000"),
		SortBy:   ptr("rating:desc"),
	}

	results, err := client.Search(ctx, collectionName, searchParams)
	require.NoError(t, err, "filtered search should succeed")
	assert.Equal(t, 2, int(*results.Found), "should find 2 movies after year 2000")

	// Verify sorting (highest rating first)
	if len(*results.Hits) > 0 {
		firstHit := (*results.Hits)[0]
		rating := (*firstHit.Document)["rating"]
		assert.Equal(t, float64(9.0), rating, "highest rated movie should be first")
	}
}
