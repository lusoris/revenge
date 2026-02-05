// Package jobs provides shared types and utilities for background job processing.
// It defines common patterns for job arguments, workers, and result handling
// using River as the job processing framework.
package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// JobResult represents the outcome of a job execution.
type JobResult struct {
	// Success indicates if the job completed successfully.
	Success bool

	// ItemsProcessed is the number of items processed.
	ItemsProcessed int

	// ItemsFailed is the number of items that failed.
	ItemsFailed int

	// Duration is how long the job took.
	Duration time.Duration

	// Errors contains any errors that occurred.
	Errors []error

	// Message is an optional human-readable summary.
	Message string
}

// AddError adds an error to the result and marks it as not fully successful.
func (r *JobResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
	r.ItemsFailed++
}

// HasErrors returns true if any errors occurred.
func (r *JobResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// LogSummary logs the job result with zap.
func (r *JobResult) LogSummary(logger *zap.Logger, jobKind string) {
	if r.Success && !r.HasErrors() {
		logger.Info("job completed successfully",
			zap.String("job_kind", jobKind),
			zap.Int("items_processed", r.ItemsProcessed),
			zap.Duration("duration", r.Duration),
		)
	} else if r.Success && r.HasErrors() {
		logger.Warn("job completed with errors",
			zap.String("job_kind", jobKind),
			zap.Int("items_processed", r.ItemsProcessed),
			zap.Int("items_failed", r.ItemsFailed),
			zap.Duration("duration", r.Duration),
			zap.Int("error_count", len(r.Errors)),
		)
	} else {
		logger.Error("job failed",
			zap.String("job_kind", jobKind),
			zap.Int("items_processed", r.ItemsProcessed),
			zap.Int("items_failed", r.ItemsFailed),
			zap.Duration("duration", r.Duration),
			zap.Int("error_count", len(r.Errors)),
		)
	}
}

// LogErrors logs individual errors (up to maxErrors).
func (r *JobResult) LogErrors(logger *zap.Logger, maxErrors int) {
	for i, err := range r.Errors {
		if i >= maxErrors {
			logger.Warn("additional errors truncated",
				zap.Int("total_errors", len(r.Errors)),
				zap.Int("shown_errors", maxErrors),
			)
			break
		}
		logger.Warn("job error",
			zap.Int("error_index", i),
			zap.Error(err),
		)
	}
}

// LibraryScanArgs represents common arguments for library scan jobs.
type LibraryScanArgs struct {
	// Paths are the library paths to scan.
	Paths []string `json:"paths"`

	// Force indicates whether to force a full rescan.
	Force bool `json:"force"`
}

// FileMatchArgs represents common arguments for file match jobs.
type FileMatchArgs struct {
	// FilePath is the path to the file to match.
	FilePath string `json:"file_path"`

	// ForceRematch indicates whether to rematch even if already matched.
	ForceRematch bool `json:"force_rematch"`
}

// MetadataRefreshArgs represents common arguments for metadata refresh jobs.
type MetadataRefreshArgs struct {
	// ContentID is the UUID of the content to refresh.
	ContentID uuid.UUID `json:"content_id"`

	// Force indicates whether to force a refresh even if recently updated.
	Force bool `json:"force"`
}

// SearchIndexArgs represents common arguments for search index jobs.
type SearchIndexArgs struct {
	// ContentID is the optional UUID of specific content to index.
	// If nil, reindex all content.
	ContentID *uuid.UUID `json:"content_id,omitempty"`

	// FullReindex indicates whether to do a full reindex.
	FullReindex bool `json:"full_reindex"`
}

// JobContext wraps context.Context with additional job-specific utilities.
type JobContext struct {
	context.Context
	Logger   *zap.Logger
	JobID    int64
	JobKind  string
	StartTime time.Time
}

// NewJobContext creates a new JobContext.
func NewJobContext(ctx context.Context, logger *zap.Logger, jobID int64, jobKind string) *JobContext {
	return &JobContext{
		Context:   ctx,
		Logger:    logger.With(zap.Int64("job_id", jobID), zap.String("job_kind", jobKind)),
		JobID:     jobID,
		JobKind:   jobKind,
		StartTime: time.Now(),
	}
}

// Elapsed returns the time elapsed since the job started.
func (jc *JobContext) Elapsed() time.Duration {
	return time.Since(jc.StartTime)
}

// LogStart logs the job start.
func (jc *JobContext) LogStart(fields ...zap.Field) {
	allFields := append([]zap.Field{}, fields...)
	jc.Logger.Info("starting job", allFields...)
}

// LogComplete logs job completion with duration.
func (jc *JobContext) LogComplete(fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Duration("duration", jc.Elapsed())}, fields...)
	jc.Logger.Info("job completed", allFields...)
}

// LogError logs an error that occurred during job execution.
func (jc *JobContext) LogError(msg string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err), zap.Duration("elapsed", jc.Elapsed())}, fields...)
	jc.Logger.Error(msg, allFields...)
}

// JobKind generates a standardized job kind string.
// Format: {content_type}_{action}
// Examples: "movie_library_scan", "tvshow_metadata_refresh"
func JobKind(contentType string, action string) string {
	return fmt.Sprintf("%s_%s", contentType, action)
}

// Common job action constants
const (
	ActionLibraryScan     = "library_scan"
	ActionFileMatch       = "file_match"
	ActionMetadataRefresh = "metadata_refresh"
	ActionSearchIndex     = "search_index"
	ActionMediaProbe      = "media_probe"
)
