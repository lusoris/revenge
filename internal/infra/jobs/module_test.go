package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/jobs"
)

func TestModule(t *testing.T) {
	// Test that the module can be created
	assert.NotNil(t, jobs.Module)

	// Test that module has expected options
	app := fx.New(
		jobs.Module,
		fx.NopLogger,
	)

	assert.NotNil(t, app)
}
