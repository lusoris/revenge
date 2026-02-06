package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/riverqueue/river"
)

// JobProgress represents the progress of a running job.
// This is stored as the job's output and can be polled by the frontend.
type JobProgress struct {
	// Phase describes the current phase of the job (e.g., "scanning", "matching", "indexing").
	Phase string `json:"phase"`

	// Current is the number of items processed so far.
	Current int `json:"current"`

	// Total is the total number of items to process (0 if unknown).
	Total int `json:"total,omitempty"`

	// Percent is the completion percentage (0-100). Calculated automatically if Total > 0.
	Percent int `json:"percent"`

	// Message is an optional human-readable status message.
	Message string `json:"message,omitempty"`
}

// ReportProgress persists job progress immediately so it can be polled by the frontend.
// This uses River's JobUpdate to eagerly store the progress as the job's output.
func (c *Client) ReportProgress(ctx context.Context, jobID int64, progress *JobProgress) error {
	if c.client == nil {
		return nil // no-op if client not initialized
	}

	// Auto-calculate percent
	if progress.Total > 0 {
		progress.Percent = (progress.Current * 100) / progress.Total
	}

	_, err := c.client.JobUpdate(ctx, jobID, &river.JobUpdateParams{
		Output: progress,
	})
	if err != nil {
		return fmt.Errorf("failed to report job progress: %w", err)
	}
	return nil
}

// GetJobProgress retrieves the current progress of a running job.
// Returns nil if the job has no progress output.
func (c *Client) GetJobProgress(ctx context.Context, jobID int64) (*JobProgress, error) {
	if c.client == nil {
		return nil, fmt.Errorf("river client not initialized")
	}

	row, err := c.client.JobGet(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	if len(row.Metadata) == 0 {
		return nil, nil
	}

	// Output is stored under the "output" key in metadata
	var metadata map[string]json.RawMessage
	if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
		return nil, nil
	}

	outputRaw, ok := metadata["output"]
	if !ok || len(outputRaw) == 0 {
		return nil, nil
	}

	var progress JobProgress
	if err := json.Unmarshal(outputRaw, &progress); err != nil {
		return nil, nil
	}

	return &progress, nil
}
