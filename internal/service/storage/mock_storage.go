package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"sync"
)

// MockStorage implements Storage interface for testing.
type MockStorage struct {
	mu    sync.RWMutex
	files map[string][]byte
}

// NewMockStorage creates a new mock storage for testing.
func NewMockStorage() *MockStorage {
	return &MockStorage{
		files: make(map[string][]byte),
	}
}

// Store saves a file to in-memory storage.
func (s *MockStorage) Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var data []byte
	if reader != nil {
		var err error
		data, err = io.ReadAll(reader)
		if err != nil {
			return "", fmt.Errorf("failed to read data: %w", err)
		}
	}

	s.files[key] = data
	return key, nil
}

// Get retrieves a file from in-memory storage.
func (s *MockStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.files[key]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", key)
	}

	return io.NopCloser(bytes.NewReader(data)), nil
}

// Delete removes a file from in-memory storage.
func (s *MockStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.files, key)
	return nil
}

// Exists checks if a file exists in in-memory storage.
func (s *MockStorage) Exists(ctx context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.files[key]
	return ok, nil
}

// GetURL returns a mock URL for the file.
func (s *MockStorage) GetURL(key string) string {
	return "/api/v1/files/" + key
}

// GetStoredFiles returns all stored files (for test assertions).
func (s *MockStorage) GetStoredFiles() map[string][]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string][]byte)
	maps.Copy(result, s.files)
	return result
}
