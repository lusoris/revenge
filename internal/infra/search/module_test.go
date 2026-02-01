package search_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/search"
)

func TestModule(t *testing.T) {
	// Test that the module can be created
	assert.NotNil(t, search.Module)

	// Test that module has expected options
	app := fx.New(
		search.Module,
		fx.NopLogger,
	)

	assert.NotNil(t, app)
}
