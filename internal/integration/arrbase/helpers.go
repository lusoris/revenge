package arrbase

// ResolutionToString converts a resolution value to a human-readable string.
func ResolutionToString(res int) string {
	switch {
	case res >= 2160:
		return "4K"
	case res >= 1080:
		return "1080p"
	case res >= 720:
		return "720p"
	case res >= 480:
		return "480p"
	default:
		return "SD"
	}
}

// PtrString returns a pointer to a string, or nil if empty.
func PtrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
