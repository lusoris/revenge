// Package app provides the main application module that wires all dependencies together.
package app

import (
	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/lusoris/revenge/internal/infra/search"
	"go.uber.org/fx"
)

// Module is the main application module that includes all sub-modules.
var Module = fx.Module("app",
	// Configuration
	config.Module,

	// Infrastructure
	logging.Module,
	database.Module,
	cache.Module,
	search.Module,
	jobs.Module,
	health.Module,

	// HTTP API Server (ogen-generated)
	api.Module,
)
