//go:build tools

// Package tools imports packages that are used for code generation and testing.
// This ensures they are tracked in go.mod even though they're not directly imported
// in the main codebase.
package tools

import (
	_ "github.com/fergusstrange/embedded-postgres" // v1.30.0
	_ "github.com/ogen-go/ogen"                     // v1.18.0
	_ "github.com/stretchr/testify"                 // v1.11.1
)
