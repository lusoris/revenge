package jobs

import (
	"time"

	"github.com/lusoris/revenge/internal/validate"
	"github.com/riverqueue/river"
)

// Queue names for different priority levels (5-level system).
// Higher priority queues are processed before lower priority ones.
const (
	// QueueCritical handles security events, auth failures, urgent system tasks.
	// Highest priority - always processed first.
	QueueCritical = "critical"

	// QueueHigh handles user-initiated actions, notifications, webhooks.
	// High priority for responsive user experience.
	QueueHigh = "high"

	// QueueDefault handles metadata fetching, sync operations, general tasks.
	// Standard priority for most background work.
	QueueDefault = river.QueueDefault

	// QueueLow handles cleanup, maintenance, session pruning, expired token cleanup.
	// Low priority tasks that can wait.
	QueueLow = "low"

	// QueueBulk handles library scans, batch operations, search reindexing.
	// Lowest priority for resource-intensive batch operations.
	QueueBulk = "bulk"
)

// QueueConfig holds configuration for all queues.
type QueueConfig struct {
	Queues map[string]river.QueueConfig
}

// DefaultQueueConfig returns the default queue configuration with 5 priority levels.
// Worker allocation reflects priority: critical gets most workers, bulk gets fewest.
func DefaultQueueConfig() *QueueConfig {
	return &QueueConfig{
		Queues: map[string]river.QueueConfig{
			QueueCritical: {MaxWorkers: 20}, // Security events, auth failures
			QueueHigh:     {MaxWorkers: 15}, // Notifications, webhooks, user actions
			QueueDefault:  {MaxWorkers: 10}, // Metadata fetch, sync operations
			QueueLow:      {MaxWorkers: 5},  // Cleanup, maintenance
			QueueBulk:     {MaxWorkers: 3},  // Library scans, batch operations
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
// Priority mapping:
//
//	priority >= 20: critical (security, urgent)
//	priority >= 10: high (user actions, notifications)
//	priority >= -9: default (general tasks)
//	priority >= -19: low (maintenance)
//	priority < -19: bulk (batch operations)
func QueuePriority(priority int) string {
	switch {
	case priority >= 20:
		return QueueCritical
	case priority >= 10:
		return QueueHigh
	case priority >= -9:
		return QueueDefault
	case priority >= -19:
		return QueueLow
	default:
		return QueueBulk
	}
}
