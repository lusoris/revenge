package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/lusoris/revenge/internal/config"
)

// S3Storage implements Storage using S3-compatible object storage (AWS S3, MinIO, etc).
// Suitable for production clustering where multiple instances need shared storage.
type S3Storage struct {
	client   *s3.Client
	uploader *s3manager.Uploader //nolint:staticcheck // transfermanager is pre-v1, using stable manager
	bucket   string
	endpoint string
	logger   *slog.Logger
}

// NewS3Storage creates a new S3-compatible storage backend.
// Supports AWS S3, MinIO, and any S3-compatible storage.
func NewS3Storage(cfg config.S3Config, logger *slog.Logger) (*S3Storage, error) {
	// Create credentials provider
	credsProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"", // session token (not needed)
	)

	// Load AWS config with credentials
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credsProvider),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with custom endpoint for MinIO compatibility
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		if cfg.UsePathStyle {
			o.UsePathStyle = true // Required for MinIO
		}
	})

	storage := &S3Storage{
		client:   client,
		uploader: s3manager.NewUploader(client), //nolint:staticcheck // transfermanager is pre-v1
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		logger:   logger.With("component", "s3-storage"),
	}

	// Verify bucket exists (optional health check)
	if err := storage.verifyBucket(context.Background()); err != nil {
		logger.Warn("Failed to verify S3 bucket", slog.Any("error",err))
	}

	return storage, nil
}

// verifyBucket checks if the configured bucket exists and is accessible.
func (s *S3Storage) verifyBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return fmt.Errorf("bucket not accessible: %w", err)
	}
	return nil
}

// Store uploads a file to S3 and returns its key.
// Uses the S3 upload manager which automatically handles multipart uploads
// for large files (>5 MiB), improving throughput and reliability.
func (s *S3Storage) Store(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	// Sanitize key to prevent path traversal
	key = sanitizeKey(key)

	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{ //nolint:staticcheck // transfermanager is pre-v1
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	s.logger.Info("File stored in S3",
		slog.String("bucket", s.bucket),
		slog.String("key", key),
		slog.String("content_type", contentType))

	return key, nil
}

// Get retrieves a file from S3.
func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	key = sanitizeKey(key)

	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get from S3: %w", err)
	}

	return output.Body, nil
}

// Delete removes a file from S3.
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	key = sanitizeKey(key)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	s.logger.Info("File deleted from S3",
		slog.String("bucket", s.bucket),
		slog.String("key", key))

	return nil
}

// Exists checks if a file exists in S3.
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	key = sanitizeKey(key)

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if error is "not found"
		// AWS SDK v2 doesn't expose error codes directly, so we check the error string
		if isNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check S3 object: %w", err)
	}

	return true, nil
}

// GetURL returns the S3 URL for accessing a file.
// For public buckets, this will be the direct S3 URL.
// For private buckets, consider using pre-signed URLs instead.
func (s *S3Storage) GetURL(key string) string {
	key = sanitizeKey(key)

	// If custom endpoint (MinIO), use that
	if s.endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key)
	}

	// Standard AWS S3 URL
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
}

// isNotFoundError checks if the error is a "not found" error from S3.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	_, noSuchKey := errors.AsType[*types.NoSuchKey](err)
	_, notFound := errors.AsType[*types.NotFound](err)
	return noSuchKey || notFound
}
