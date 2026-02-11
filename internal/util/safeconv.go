// Package util provides common utility functions
package util

import "math"

// SafeIntToInt32 converts int to int32, capping at int32 bounds
func SafeIntToInt32(v int) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v) // #nosec G115 -- bounds checked above
}

// SafeInt64ToInt32 converts int64 to int32, capping at int32 bounds
func SafeInt64ToInt32(v int64) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v) // #nosec G115 -- bounds checked above
}

// SafeUint64ToInt32 converts uint64 to int32, capping at int32 max
func SafeUint64ToInt32(v uint64) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(v) // #nosec G115 -- bounds checked above
}

// SafeInt32ToUint32 converts int32 to uint32, treating negatives as 0
func SafeInt32ToUint32(v int32) uint32 {
	if v < 0 {
		return 0
	}
	return uint32(v) // #nosec G115 -- bounds checked above
}

// SafeUint32ToInt32 converts uint32 to int32, capping at int32 max
func SafeUint32ToInt32(v uint32) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(v) // #nosec G115 -- bounds checked above
}

// SafeIntToUint converts int to uint, treating negatives as 0
func SafeIntToUint(v int) uint {
	if v < 0 {
		return 0
	}
	return uint(v) // #nosec G115 -- bounds checked above
}

// SafeUintToInt converts uint to int, capping at math.MaxInt
func SafeUintToInt(v uint) int {
	if v > uint(math.MaxInt) {
		return math.MaxInt
	}
	return int(v) // #nosec G115 -- bounds checked above
}
