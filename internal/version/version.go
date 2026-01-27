// Package version provides version information for Jellyfin Go.
package version

import (
	"fmt"
	"runtime"
)

// Build-time variables (set via ldflags)
var (
	// Version is the semantic version of the application
	Version = "dev"

	// GitCommit is the git commit SHA
	GitCommit = "unknown"

	// BuildTime is the build timestamp
	BuildTime = "unknown"

	// GoVersion is the Go version used to build
	GoVersion = runtime.Version()
)

// Info contains version information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// GetInfo returns the version information
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf(
		"Jellyfin Go %s (commit: %s, built: %s, %s, %s/%s)",
		i.Version,
		i.GitCommit,
		i.BuildTime,
		i.GoVersion,
		i.OS,
		i.Arch,
	)
}

// Short returns a short version string
func Short() string {
	if GitCommit != "unknown" && len(GitCommit) > 7 {
		return fmt.Sprintf("%s (%s)", Version, GitCommit[:7])
	}
	return Version
}

// Full returns the full version string
func Full() string {
	return GetInfo().String()
}
