package database

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	StopSharedPG()
	os.Exit(code)
}
