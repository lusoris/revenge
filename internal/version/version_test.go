package version

import (
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	info := Info()

	if info == "" {
		t.Error("Info() returned empty string")
	}

	if !strings.Contains(info, Version) {
		t.Errorf("Info() = %q, want to contain Version %q", info, Version)
	}

	if !strings.Contains(info, Commit) {
		t.Errorf("Info() = %q, want to contain Commit %q", info, Commit)
	}
}

func TestDefaultValues(t *testing.T) {
	// Default values should be set
	if Version == "" {
		t.Error("Version should have a default value")
	}
	if Commit == "" {
		t.Error("Commit should have a default value")
	}
	if Date == "" {
		t.Error("Date should have a default value")
	}
}
