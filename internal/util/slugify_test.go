package util

import "testing"

func TestSlugify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"Action", "action"},
		{"Science Fiction", "science-fiction"},
		{"Sci-Fi & Fantasy", "sci-fi-fantasy"},
		{"Action & Adventure", "action-adventure"},
		{"War & Politics", "war-politics"},
		{"TV Movie", "tv-movie"},
		{"  Slice of Life  ", "slice-of-life"},
		{"Mahou Shoujo", "mahou-shoujo"},
		{"", ""},
		{"Horror", "horror"},
		{"romance", "romance"},
		{"Sci-Fi", "sci-fi"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
