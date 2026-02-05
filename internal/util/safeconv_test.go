package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeIntToInt32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected int32
	}{
		{"zero", 0, 0},
		{"positive in range", 100, 100},
		{"negative in range", -100, -100},
		{"max int32", math.MaxInt32, math.MaxInt32},
		{"min int32", math.MinInt32, math.MinInt32},
		{"above max int32", math.MaxInt32 + 1, math.MaxInt32},
		{"below min int32", math.MinInt32 - 1, math.MinInt32},
		{"large positive", math.MaxInt64, math.MaxInt32},
		{"large negative", math.MinInt64, math.MinInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeIntToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeInt64ToInt32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int64
		expected int32
	}{
		{"zero", 0, 0},
		{"positive in range", 100, 100},
		{"negative in range", -100, -100},
		{"max int32", math.MaxInt32, math.MaxInt32},
		{"min int32", math.MinInt32, math.MinInt32},
		{"above max int32", math.MaxInt32 + 1, math.MaxInt32},
		{"below min int32", math.MinInt32 - 1, math.MinInt32},
		{"max int64", math.MaxInt64, math.MaxInt32},
		{"min int64", math.MinInt64, math.MinInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeInt64ToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeUint64ToInt32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    uint64
		expected int32
	}{
		{"zero", 0, 0},
		{"positive in range", 100, 100},
		{"max int32", math.MaxInt32, math.MaxInt32},
		{"above max int32", math.MaxInt32 + 1, math.MaxInt32},
		{"max uint64", math.MaxUint64, math.MaxInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeUint64ToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeInt32ToUint32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int32
		expected uint32
	}{
		{"zero", 0, 0},
		{"positive", 100, 100},
		{"negative", -100, 0},
		{"max int32", math.MaxInt32, math.MaxInt32},
		{"min int32", math.MinInt32, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeInt32ToUint32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeUint32ToInt32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    uint32
		expected int32
	}{
		{"zero", 0, 0},
		{"positive in range", 100, 100},
		{"max int32", math.MaxInt32, math.MaxInt32},
		{"above max int32", math.MaxInt32 + 1, math.MaxInt32},
		{"max uint32", math.MaxUint32, math.MaxInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeUint32ToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeIntToUint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected uint
	}{
		{"zero", 0, 0},
		{"positive", 100, 100},
		{"negative", -100, 0},
		{"large negative", math.MinInt, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeIntToUint(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
