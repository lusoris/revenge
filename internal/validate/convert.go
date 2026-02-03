// Package validate provides validation and safe conversion utilities
package validate

import (
	"fmt"
	"math"
)

// SafeInt32 safely converts an int to int32, returning an error if overflow would occur
func SafeInt32(value int) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("value %d overflows int32 range [%d, %d]", value, math.MinInt32, math.MaxInt32)
	}
	return int32(value), nil
}

// MustInt32 converts int to int32, panicking on overflow
// Use this only when you're certain the value is in range
func MustInt32(value int) int32 {
	if value > math.MaxInt32 || value < math.MinInt32 {
		panic(fmt.Sprintf("value %d overflows int32 range", value))
	}
	return int32(value)
}

// SafeUint32 safely converts an int to uint32, returning an error if overflow would occur
func SafeUint32(value int) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, fmt.Errorf("value %d overflows uint32 range [0, %d]", value, math.MaxUint32)
	}
	return uint32(value), nil
}

// MustUint32 converts int to uint32, panicking on overflow or negative values
func MustUint32(value int) uint32 {
	if value < 0 || value > math.MaxUint32 {
		panic(fmt.Sprintf("value %d overflows uint32 range", value))
	}
	return uint32(value)
}

// SafeUint safely converts an int to uint, returning an error for negative values
func SafeUint(value int) (uint, error) {
	if value < 0 {
		return 0, fmt.Errorf("value %d is negative, cannot convert to uint", value)
	}
	return uint(value), nil
}

// MustUint converts int to uint, panicking on negative values
func MustUint(value int) uint {
	if value < 0 {
		panic(fmt.Sprintf("value %d is negative, cannot convert to uint", value))
	}
	return uint(value)
}

// ValidateSliceIndex checks if an index is within bounds for a slice
func ValidateSliceIndex(index, length int) error {
	if index < 0 || index >= length {
		return fmt.Errorf("index %d out of bounds for slice of length %d", index, length)
	}
	return nil
}

// ValidateSliceRange checks if a range [start:end] is valid for a slice
func ValidateSliceRange(start, end, length int) error {
	if start < 0 || start > length {
		return fmt.Errorf("start index %d out of bounds for slice of length %d", start, length)
	}
	if end < start || end > length {
		return fmt.Errorf("end index %d invalid for slice of length %d (start: %d)", end, length, start)
	}
	return nil
}
