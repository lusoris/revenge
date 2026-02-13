// Package ptr provides generic pointer utilities for reducing boilerplate
// when working with optional/pointer fields.
package ptr

// To returns a pointer to the given value.
// Use when constructing structs with optional pointer fields.
//
//	movie := Movie{
//	    Title:    "Inception",
//	    Year:     ptr.To(2010),
//	    Director: ptr.To("Christopher Nolan"),
//	}
//
//go:fix inline
func To[T any](v T) *T {
	return new(v)
}

// Value returns the value pointed to by p, or the zero value if p is nil.
//
//	year := ptr.Value(movie.Year) // 0 if Year is nil
func Value[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// ValueOr returns the value pointed to by p, or the default if p is nil.
//
//	year := ptr.ValueOr(movie.Year, 1900) // 1900 if Year is nil
func ValueOr[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// Equal returns true if both pointers are nil, or both point to equal values.
// For comparable types only.
//
//	ptr.Equal(a.Year, b.Year) // true if both nil or both point to same value
func Equal[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Clone returns a pointer to a copy of the value, or nil if p is nil.
// Useful when you need to copy a pointer field without sharing memory.
//
//	movie.Year = ptr.Clone(original.Year)
func Clone[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// Coalesce returns the first non-nil pointer, or nil if all are nil.
//
//	year := ptr.Coalesce(movie.Year, defaultYear, ptr.To(1900))
func Coalesce[T any](ptrs ...*T) *T {
	for _, p := range ptrs {
		if p != nil {
			return p
		}
	}
	return nil
}
