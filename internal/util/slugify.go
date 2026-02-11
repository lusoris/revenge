package util

import (
	"regexp"
	"strings"
)

var nonAlphanumRegex = regexp.MustCompile(`[^a-z0-9]+`)
var trailingHyphenRegex = regexp.MustCompile(`-+$`)
var leadingHyphenRegex = regexp.MustCompile(`^-+`)

// Slugify converts a string to a URL-friendly slug.
// It lowercases, replaces non-alphanumeric characters with hyphens,
// and trims leading/trailing hyphens.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlphanumRegex.ReplaceAllString(s, "-")
	s = trailingHyphenRegex.ReplaceAllString(s, "")
	s = leadingHyphenRegex.ReplaceAllString(s, "")
	return s
}
