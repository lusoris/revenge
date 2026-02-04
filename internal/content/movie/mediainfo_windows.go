//go:build windows
// +build windows

package movie

import (
	"fmt"
)

// MediaInfoProber is a stub implementation for Windows
// Media probing requires FFmpeg/libav which uses CGO and is not supported on Windows builds
type MediaInfoProber struct{}

// Ensure MediaInfoProber implements Prober
var _ Prober = (*MediaInfoProber)(nil)

// NewMediaInfoProber creates a new Windows stub prober
func NewMediaInfoProber() *MediaInfoProber {
	return &MediaInfoProber{}
}

// Probe returns an error on Windows as FFmpeg/libav is not available
func (p *MediaInfoProber) Probe(filePath string) (*MediaInfo, error) {
	return nil, fmt.Errorf("media probing not supported on Windows: FFmpeg/libav requires CGO which is not available in this build")
}
