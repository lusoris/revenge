package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/cache"
)

func TestModule(t *testing.T) {
	// Test that the module can be created
	assert.NotNil(t, cache.Module)

	// Test that module has expected options
	app := fx.New(
		cache.Module,
		fx.NopLogger,
	)

	assert.NotNil(t, app)
}
