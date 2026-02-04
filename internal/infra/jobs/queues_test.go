package jobs

import (
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
)

// TestDefaultQueueConfig tests the default queue configuration with 5 priority levels.
func TestDefaultQueueConfig(t *testing.T) {
	cfg := DefaultQueueConfig()

	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.Queues)
	assert.Len(t, cfg.Queues, 5)

	// Critical queue - highest priority, most workers
	assert.Contains(t, cfg.Queues, QueueCritical)
	assert.Equal(t, 20, cfg.Queues[QueueCritical].MaxWorkers)

	// High queue - notifications, user actions
	assert.Contains(t, cfg.Queues, QueueHigh)
	assert.Equal(t, 15, cfg.Queues[QueueHigh].MaxWorkers)

	// Default queue - general tasks
	assert.Contains(t, cfg.Queues, QueueDefault)
	assert.Equal(t, 10, cfg.Queues[QueueDefault].MaxWorkers)

	// Low priority queue - maintenance
	assert.Contains(t, cfg.Queues, QueueLow)
	assert.Equal(t, 5, cfg.Queues[QueueLow].MaxWorkers)

	// Bulk queue - batch operations, fewest workers
	assert.Contains(t, cfg.Queues, QueueBulk)
	assert.Equal(t, 3, cfg.Queues[QueueBulk].MaxWorkers)
}

// TestQueueConstants tests queue name constants.
func TestQueueConstants(t *testing.T) {
	assert.Equal(t, "critical", QueueCritical)
	assert.Equal(t, "high", QueueHigh)
	assert.Equal(t, river.QueueDefault, QueueDefault)
	assert.Equal(t, "low", QueueLow)
	assert.Equal(t, "bulk", QueueBulk)
}

// TestDefaultRetryPolicy tests the default retry policy.
func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy()

	assert.NotNil(t, policy)
	assert.Equal(t, 25, policy.MaxAttempts)
	assert.NotNil(t, policy.Backoff)
}

// TestExponentialBackoff tests exponential backoff calculation.
func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		name     string
		attempt  int
		expected time.Duration
	}{
		{"attempt 0", 0, 1 * time.Second},
		{"attempt 1", 1, 2 * time.Second},
		{"attempt 2", 2, 4 * time.Second},
		{"attempt 3", 3, 8 * time.Second},
		{"attempt 4", 4, 16 * time.Second},
		{"attempt 5", 5, 32 * time.Second},
		{"attempt 10", 10, 1024 * time.Second},
		{"attempt 20", 20, 1 * time.Hour}, // max cap
		{"attempt 30", 30, 1 * time.Hour}, // max cap
		{"negative attempt", -5, 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExponentialBackoff(tt.attempt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLinearBackoff tests linear backoff calculation.
func TestLinearBackoff(t *testing.T) {
	tests := []struct {
		name     string
		attempt  int
		expected time.Duration
	}{
		{"attempt 0", 0, 0 * time.Second},
		{"attempt 1", 1, 30 * time.Second},
		{"attempt 2", 2, 60 * time.Second},
		{"attempt 5", 5, 150 * time.Second},
		{"attempt 10", 10, 300 * time.Second},
		{"attempt 60", 60, 30 * time.Minute},   // max cap
		{"attempt 100", 100, 30 * time.Minute}, // max cap
		{"negative attempt", -5, 0 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LinearBackoff(tt.attempt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestQueuePriority tests queue priority mapping with 5 levels.
func TestQueuePriority(t *testing.T) {
	tests := []struct {
		name     string
		priority int
		expected string
	}{
		// Critical: >= 20
		{"critical priority 20", 20, QueueCritical},
		{"critical priority 25", 25, QueueCritical},
		{"critical priority 100", 100, QueueCritical},
		// High: >= 10
		{"high priority 10", 10, QueueHigh},
		{"high priority 15", 15, QueueHigh},
		{"high priority 19", 19, QueueHigh},
		// Default: >= -9
		{"default priority 0", 0, QueueDefault},
		{"default priority 5", 5, QueueDefault},
		{"default priority 9", 9, QueueDefault},
		{"default priority -5", -5, QueueDefault},
		{"default priority -9", -9, QueueDefault},
		// Low: >= -19
		{"low priority -10", -10, QueueLow},
		{"low priority -15", -15, QueueLow},
		{"low priority -19", -19, QueueLow},
		// Bulk: < -19
		{"bulk priority -20", -20, QueueBulk},
		{"bulk priority -50", -50, QueueBulk},
		{"bulk priority -100", -100, QueueBulk},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QueuePriority(tt.priority)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestQueueConfig_CustomQueues tests custom queue configuration.
func TestQueueConfig_CustomQueues(t *testing.T) {
	cfg := &QueueConfig{
		Queues: map[string]river.QueueConfig{
			"custom_queue": {MaxWorkers: 50},
		},
	}

	assert.NotNil(t, cfg)
	assert.Len(t, cfg.Queues, 1)
	assert.Contains(t, cfg.Queues, "custom_queue")
	assert.Equal(t, 50, cfg.Queues["custom_queue"].MaxWorkers)
}

// TestRetryPolicy_CustomBackoff tests custom retry policy.
func TestRetryPolicy_CustomBackoff(t *testing.T) {
	customBackoff := func(attempt int) time.Duration {
		return time.Duration(attempt*5) * time.Second
	}

	policy := &RetryPolicy{
		MaxAttempts: 10,
		Backoff:     customBackoff,
	}

	assert.NotNil(t, policy)
	assert.Equal(t, 10, policy.MaxAttempts)
	assert.Equal(t, 15*time.Second, policy.Backoff(3))
	assert.Equal(t, 50*time.Second, policy.Backoff(10))
}

// TestExponentialBackoff_EdgeCases tests edge cases.
func TestExponentialBackoff_EdgeCases(t *testing.T) {
	// Very large attempt should cap at max
	result := ExponentialBackoff(100)
	assert.Equal(t, 1*time.Hour, result)

	// Negative attempt should be treated as 0
	result = ExponentialBackoff(-1)
	assert.Equal(t, 1*time.Second, result)
}

// TestLinearBackoff_EdgeCases tests edge cases.
func TestLinearBackoff_EdgeCases(t *testing.T) {
	// Very large attempt should cap at max
	result := LinearBackoff(1000)
	assert.Equal(t, 30*time.Minute, result)

	// Negative attempt should be treated as 0
	result = LinearBackoff(-10)
	assert.Equal(t, 0*time.Second, result)
}
