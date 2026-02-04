// Package storage provides file storage abstraction for avatars and other media.
// It supports multiple backends (local filesystem, S3, etc.) for clustering-ready deployments.
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/zap"
)

// Storage defines the interface for file storage operations.
// Implementations can use local filesystem, S3, or other backends.
type Storage interface {
	// Store saves a file and returns its path/URL
	Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
	// Get retrieves a file by key
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete removes a file by key
	Delete(ctx context.Context, key string) error
	// Exists checks if a file exists
	Exists(ctx context.Context, key string) (bool, error)
	// GetURL returns the public URL for a file
	GetURL(key string) string
}

// LocalStorage implements Storage using the local filesystem.
// For production clustering, consider using S3-compatible storage.
type LocalStorage struct {
	basePath string
	baseURL  string
	logger   *zap.Logger
	mu       sync.RWMutex
}

// NewLocalStorage creates a new local filesystem storage.
func NewLocalStorage(cfg config.AvatarConfig, logger *zap.Logger) (*LocalStorage, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(cfg.StoragePath, 0750); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: cfg.StoragePath,
		baseURL:  "/api/v1/files", // Served via API endpoint
		logger:   logger.Named("storage"),
	}, nil
}

// Store saves a file to the local filesystem.
func (s *LocalStorage) Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Sanitize key to prevent path traversal
	key = sanitizeKey(key)

	fullPath := filepath.Join(s.basePath, key)

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Copy content
	written, err := io.Copy(file, reader)
	if err != nil {
		// Cleanup on error
		_ = os.Remove(fullPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	s.logger.Info("File stored",
		zap.String("key", key),
		zap.Int64("size", written),
		zap.String("content_type", contentType))

	return key, nil
}

// Get retrieves a file from the local filesystem.
func (s *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key = sanitizeKey(key)
	fullPath := filepath.Join(s.basePath, key)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes a file from the local filesystem.
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key = sanitizeKey(key)
	fullPath := filepath.Join(s.basePath, key)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Info("File deleted", zap.String("key", key))
	return nil
}

// Exists checks if a file exists in the local filesystem.
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key = sanitizeKey(key)
	fullPath := filepath.Join(s.basePath, key)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file: %w", err)
	}

	return true, nil
}

// GetURL returns the API URL for accessing a file.
func (s *LocalStorage) GetURL(key string) string {
	key = sanitizeKey(key)
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}

// GenerateAvatarKey generates a unique storage key for an avatar.
func GenerateAvatarKey(userID uuid.UUID, filename string) string {
	// Extract extension
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png" // Default extension
	}

	// Generate unique filename with user prefix for organization
	uniqueID := uuid.New()
	return fmt.Sprintf("avatars/%s/%s%s", userID.String(), uniqueID.String(), ext)
}

// sanitizeKey prevents path traversal attacks.
func sanitizeKey(key string) string {
	// Remove leading slashes
	key = strings.TrimPrefix(key, "/")

	// Clean the path
	key = filepath.Clean(key)

	// Remove any parent directory references
	parts := strings.Split(key, string(filepath.Separator))
	var safe []string
	for _, part := range parts {
		if part != ".." && part != "." && part != "" {
			safe = append(safe, part)
		}
	}

	return strings.Join(safe, string(filepath.Separator))
}
