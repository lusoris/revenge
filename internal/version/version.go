// Package version provides build version information for the Revenge server.
package version

// Build information set via ldflags at compile time.
var (
	// Version is the semantic version (e.g., "1.0.0").
	Version = "dev"

	// Commit is the git commit hash.
	Commit = "unknown"

	// Date is the build date.
	Date = "unknown"
)

// Info returns version information as a formatted string.
func Info() string {
	return Version + " (" + Commit + ") built " + Date
}
