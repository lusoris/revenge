package jobs

import (
	"time"

	"github.com/lusoris/revenge/internal/validate"
	"github.com/riverqueue/river"
)

// Queue names for different priority levels.
const (
	QueueCritical = "critical"
	QueueDefault  = river.QueueDefault
	QueueLow      = "low"
)

// QueueConfig holds configuration for all queues.
type QueueConfig struct {
	Queues map[string]river.QueueConfig
}

// DefaultQueueConfig returns the default queue configuration.
func DefaultQueueConfig() *QueueConfig {
	return &QueueConfig{
		Queues: map[string]river.QueueConfig{
			QueueCritical: {MaxWorkers: 20},
			QueueDefault:  {MaxWorkers: 10},
			QueueLow:      {MaxWorkers: 5},
		},
	}
}

// RetryPolicy holds retry configuration for jobs.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     func(attempt int) time.Duration
}

// DefaultRetryPolicy returns the default retry policy.
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts: 25,
		Backoff:     ExponentialBackoff,
	}
}

// ExponentialBackoff calculates exponential backoff duration.
// Formula: min(base * 2^attempt, max)
func ExponentialBackoff(attempt int) time.Duration {
	const (
		base = 1 * time.Second
		max  = 1 * time.Hour
	)

	if attempt < 0 {
		attempt = 0
	}

	// Cap attempt to prevent overflow (2^30 seconds is already > 1 hour)
	if attempt > 30 {
		return max
	}

	// Safe conversion for bitshift operation
	attemptUint := validate.MustUint(attempt)
	duration := base * (1 << attemptUint) // 2^attempt
	if duration > max {
		return max
	}
	return duration
}

// LinearBackoff calculates linear backoff duration.
// Formula: min(base * attempt, max)
func LinearBackoff(attempt int) time.Duration {
	const (
		base = 30 * time.Second
		max  = 30 * time.Minute
	)

	if attempt < 0 {
		attempt = 0
	}

	duration := base * time.Duration(attempt)
	if duration > max {
		return max
	}
	return duration
}

// QueuePriority returns the queue name for a given priority level.
func QueuePriority(priority int) string {
	switch {
	case priority >= 10:
		return QueueCritical
	case priority <= -10:
		return QueueLow
	default:
		return QueueDefault
	}
}
