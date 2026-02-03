package validate_test

import (
	"math"
	"testing"

	"github.com/lusoris/revenge/internal/validate"
)

func TestSafeInt32(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    int32
		wantErr bool
	}{
		{"Valid positive", 42, 42, false},
		{"Valid negative", -42, -42, false},
		{"Max int32", math.MaxInt32, math.MaxInt32, false},
		{"Min int32", math.MinInt32, math.MinInt32, false},
		{"Overflow positive", math.MaxInt32 + 1, 0, true},
		{"Overflow negative", math.MinInt32 - 1, 0, true},
		{"Zero", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate.SafeInt32(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SafeInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustInt32(t *testing.T) {
	// Valid conversions
	if got := validate.MustInt32(42); got != 42 {
		t.Errorf("MustInt32(42) = %v, want 42", got)
	}

	// Overflow should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustInt32(overflow) should panic")
		}
	}()
	validate.MustInt32(math.MaxInt32 + 1)
}

func TestSafeUint32(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    uint32
		wantErr bool
	}{
		{"Valid positive", 42, 42, false},
		{"Zero", 0, 0, false},
		{"Max uint32", math.MaxUint32, math.MaxUint32, false},
		{"Negative", -1, 0, true},
		{"Overflow", math.MaxUint32 + 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate.SafeUint32(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SafeUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustUint32(t *testing.T) {
	// Valid conversion
	if got := validate.MustUint32(42); got != 42 {
		t.Errorf("MustUint32(42) = %v, want 42", got)
	}

	// Negative should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustUint32(negative) should panic")
		}
	}()
	validate.MustUint32(-1)
}

func TestSafeUint(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    uint
		wantErr bool
	}{
		{"Valid positive", 42, 42, false},
		{"Zero", 0, 0, false},
		{"Large value", 1000000, 1000000, false},
		{"Negative", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate.SafeUint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SafeUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustUint(t *testing.T) {
	// Valid conversion
	if got := validate.MustUint(42); got != 42 {
		t.Errorf("MustUint(42) = %v, want 42", got)
	}

	// Negative should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustUint(negative) should panic")
		}
	}()
	validate.MustUint(-1)
}

func TestValidateSliceIndex(t *testing.T) {
	tests := []struct {
		name    string
		index   int
		length  int
		wantErr bool
	}{
		{"Valid index", 5, 10, false},
		{"First index", 0, 10, false},
		{"Last index", 9, 10, false},
		{"Negative index", -1, 10, true},
		{"Index equals length", 10, 10, true},
		{"Index exceeds length", 15, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.ValidateSliceIndex(tt.index, tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSliceIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSliceRange(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		end     int
		length  int
		wantErr bool
	}{
		{"Valid range", 2, 5, 10, false},
		{"Full slice", 0, 10, 10, false},
		{"Empty range", 5, 5, 10, false},
		{"Negative start", -1, 5, 10, true},
		{"Start exceeds length", 11, 15, 10, true},
		{"End before start", 5, 2, 10, true},
		{"End exceeds length", 5, 15, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.ValidateSliceRange(tt.start, tt.end, tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSliceRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
