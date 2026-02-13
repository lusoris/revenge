package api

import (
	"github.com/govalues/decimal"
)

// setOpt sets an ogen optional field from a pointer value of the same type.
// If v is nil, the optional field remains unset.
//
// Usage:
//
//	setOpt(&o.Overview, m.Overview) // *string → OptString
//	setOpt(&o.ReleaseDate, m.ReleaseDate) // *time.Time → OptDateTime
func setOpt[T any](opt interface{ SetTo(T) }, v *T) {
	if v != nil {
		opt.SetTo(*v)
	}
}

// setOptConv sets an ogen optional field from a pointer value with a type
// conversion function. If v is nil, the optional field remains unset.
//
// Usage:
//
//	setOptConv(&o.TmdbID, m.TMDbID, int32ToInt) // *int32 → OptInt
//	setOptConv(&o.Budget, m.Budget, int64ToInt64) // *int64 → OptInt64
func setOptConv[S, T any](opt interface{ SetTo(T) }, v *S, conv func(S) T) {
	if v != nil {
		opt.SetTo(conv(*v))
	}
}

// setOptDecimalFloat32 sets an ogen OptFloat32 from a *decimal.Decimal.
// This handles the common decimal → float32 conversion used for vote averages,
// popularity scores, and other decimal fields.
func setOptDecimalFloat32(opt interface{ SetTo(float32) }, v *decimal.Decimal) {
	if v != nil {
		f, _ := v.Float64()
		opt.SetTo(float32(f))
	}
}

// Common conversion functions for use with setOptConv.

func int32ToInt(v int32) int     { return int(v) }
func int64ToInt64(v int64) int64 { return v } // identity — for documentation clarity
