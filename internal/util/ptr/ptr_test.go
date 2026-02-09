package ptr_test

import (
	"testing"

	"github.com/lusoris/revenge/internal/util/ptr"
)

func TestTo(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		p := ptr.To(42)
		if p == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *p != 42 {
			t.Errorf("expected 42, got %d", *p)
		}
	})

	t.Run("string", func(t *testing.T) {
		p := ptr.To("hello")
		if p == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *p != "hello" {
			t.Errorf("expected hello, got %s", *p)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type S struct{ X int }
		p := ptr.To(S{X: 10})
		if p == nil {
			t.Fatal("expected non-nil pointer")
		}
		if p.X != 10 {
			t.Errorf("expected X=10, got X=%d", p.X)
		}
	})
}

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("non-nil", func(t *testing.T) {
		p := ptr.To(42)
		if v := ptr.Value(p); v != 42 {
			t.Errorf("expected 42, got %d", v)
		}
	})

	t.Run("nil returns zero", func(t *testing.T) {
		var p *int
		if v := ptr.Value(p); v != 0 {
			t.Errorf("expected 0, got %d", v)
		}
	})

	t.Run("nil string returns empty", func(t *testing.T) {
		var p *string
		if v := ptr.Value(p); v != "" {
			t.Errorf("expected empty string, got %q", v)
		}
	})
}

func TestValueOr(t *testing.T) {
	t.Parallel()

	t.Run("non-nil returns value", func(t *testing.T) {
		p := ptr.To(42)
		if v := ptr.ValueOr(p, 100); v != 42 {
			t.Errorf("expected 42, got %d", v)
		}
	})

	t.Run("nil returns default", func(t *testing.T) {
		var p *int
		if v := ptr.ValueOr(p, 100); v != 100 {
			t.Errorf("expected 100, got %d", v)
		}
	})
}

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("both nil", func(t *testing.T) {
		var a, b *int
		if !ptr.Equal(a, b) {
			t.Error("expected both nil to be equal")
		}
	})

	t.Run("one nil", func(t *testing.T) {
		a := ptr.To(42)
		var b *int
		if ptr.Equal(a, b) {
			t.Error("expected non-nil and nil to not be equal")
		}
		if ptr.Equal(b, a) {
			t.Error("expected nil and non-nil to not be equal")
		}
	})

	t.Run("same value", func(t *testing.T) {
		a := ptr.To(42)
		b := ptr.To(42)
		if !ptr.Equal(a, b) {
			t.Error("expected same values to be equal")
		}
	})

	t.Run("different values", func(t *testing.T) {
		a := ptr.To(42)
		b := ptr.To(43)
		if ptr.Equal(a, b) {
			t.Error("expected different values to not be equal")
		}
	})
}

func TestClone(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		var p *int
		if c := ptr.Clone(p); c != nil {
			t.Error("expected nil")
		}
	})

	t.Run("clones value", func(t *testing.T) {
		p := ptr.To(42)
		c := ptr.Clone(p)
		if c == nil {
			t.Fatal("expected non-nil")
		}
		if *c != 42 {
			t.Errorf("expected 42, got %d", *c)
		}
		// Verify it's a different pointer
		if p == c {
			t.Error("expected different pointer addresses")
		}
	})
}

func TestCoalesce(t *testing.T) {
	t.Parallel()

	t.Run("all nil", func(t *testing.T) {
		var a, b, c *int
		if p := ptr.Coalesce(a, b, c); p != nil {
			t.Error("expected nil")
		}
	})

	t.Run("first non-nil", func(t *testing.T) {
		a := ptr.To(1)
		b := ptr.To(2)
		c := ptr.To(3)
		if p := ptr.Coalesce(a, b, c); *p != 1 {
			t.Errorf("expected 1, got %d", *p)
		}
	})

	t.Run("middle non-nil", func(t *testing.T) {
		var a *int
		b := ptr.To(2)
		c := ptr.To(3)
		if p := ptr.Coalesce(a, b, c); *p != 2 {
			t.Errorf("expected 2, got %d", *p)
		}
	})

	t.Run("last non-nil", func(t *testing.T) {
		var a, b *int
		c := ptr.To(3)
		if p := ptr.Coalesce(a, b, c); *p != 3 {
			t.Errorf("expected 3, got %d", *p)
		}
	})
}
