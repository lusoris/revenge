// Package metadatafx provides fx dependency injection for the metadata service.
package metadatafx

import (
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/service/metadata"
	movieadapter "github.com/lusoris/revenge/internal/service/metadata/adapters/movie"
	tvshowadapter "github.com/lusoris/revenge/internal/service/metadata/adapters/tvshow"
	"github.com/lusoris/revenge/internal/service/metadata/providers/fanarttv"
	"github.com/lusoris/revenge/internal/service/metadata/providers/tmdb"
	"github.com/lusoris/revenge/internal/service/metadata/providers/tvdb"

	"go.uber.org/fx"
)

// Config contains metadata service configuration.
type Config struct {
	// DefaultLanguages are fetched if no specific languages are requested.
	DefaultLanguages []string

	// EnableProviderFallback enables trying secondary providers on failure.
	EnableProviderFallback bool

	// EnableEnrichment enables merging data from multiple providers.
	EnableEnrichment bool

	// TMDb configuration
	TMDbAPIKey    string
	TMDbProxyURL  string

	// TVDb configuration (optional)
	TVDbAPIKey string
	TVDbPIN    string

	// Fanart.tv configuration (optional)
	FanartTVAPIKey    string
	FanartTVClientKey string
}

// ModuleParams contains parameters for the metadata module.
type ModuleParams struct {
	fx.In

	Config         Config          `optional:"true"`
	TMDbConfig     tmdb.Config     `optional:"true"`
	TVDbConfig     tvdb.Config     `optional:"true"`
	FanartTVConfig fanarttv.Config `optional:"true"`
}

// ModuleResult contains the provided services.
type ModuleResult struct {
	fx.Out

	Service               metadata.Service
	MovieMetadataAdapter  movie.MetadataProvider
	TVShowMetadataAdapter tvshow.MetadataProvider
	TMDbProvider          *tmdb.Provider     `optional:"true"`
	TVDbProvider          *tvdb.Provider     `optional:"true"`
	FanartTVProvider      *fanarttv.Provider `optional:"true"`
}

// NewModule creates a new metadata service with providers.
func NewModule(params ModuleParams) (ModuleResult, error) {
	result := ModuleResult{}

	// Create service config
	serviceConfig := metadata.ServiceConfig{
		DefaultLanguages:       []string{"en"},
		EnableProviderFallback: true,
		EnableEnrichment:       false,
	}

	if len(params.Config.DefaultLanguages) > 0 {
		serviceConfig.DefaultLanguages = params.Config.DefaultLanguages
	}
	serviceConfig.EnableProviderFallback = params.Config.EnableProviderFallback
	serviceConfig.EnableEnrichment = params.Config.EnableEnrichment

	// Create service
	svc := metadata.NewService(serviceConfig)

	// Create and register TMDb provider if configured
	tmdbConfig := params.TMDbConfig
	if tmdbConfig.APIKey == "" && params.Config.TMDbAPIKey != "" {
		tmdbConfig = tmdb.Config{
			APIKey:   params.Config.TMDbAPIKey,
			ProxyURL: params.Config.TMDbProxyURL,
		}
	}

	if tmdbConfig.APIKey != "" {
		tmdbProvider, err := tmdb.NewProvider(tmdbConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(tmdbProvider)
		result.TMDbProvider = tmdbProvider
	}

	// Create and register TVDb provider if configured
	tvdbConfig := params.TVDbConfig
	if tvdbConfig.APIKey == "" && params.Config.TVDbAPIKey != "" {
		tvdbConfig = tvdb.Config{
			APIKey: params.Config.TVDbAPIKey,
			PIN:    params.Config.TVDbPIN,
		}
	}

	if tvdbConfig.APIKey != "" {
		tvdbProvider, err := tvdb.NewProvider(tvdbConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(tvdbProvider)
		result.TVDbProvider = tvdbProvider
	}

	// Create and register Fanart.tv provider if configured
	fanartConfig := params.FanartTVConfig
	if fanartConfig.APIKey == "" && params.Config.FanartTVAPIKey != "" {
		fanartConfig = fanarttv.Config{
			APIKey:    params.Config.FanartTVAPIKey,
			ClientKey: params.Config.FanartTVClientKey,
		}
	}

	if fanartConfig.APIKey != "" {
		fanartProvider, err := fanarttv.NewProvider(fanartConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(fanartProvider)
		result.FanartTVProvider = fanartProvider
	}

	result.Service = svc

	// Create adapters using the service
	result.MovieMetadataAdapter = movieadapter.NewAdapter(svc, serviceConfig.DefaultLanguages)
	result.TVShowMetadataAdapter = tvshowadapter.NewAdapter(svc, serviceConfig.DefaultLanguages)

	return result, nil
}

// Module provides the metadata service and providers via fx.
var Module = fx.Module("metadata",
	fx.Provide(NewModule),
)

// ProvideConfig creates a Config from environment or defaults.
// This can be overridden by the application.
func ProvideConfig() Config {
	return Config{
		DefaultLanguages:       []string{"en", "de", "fr", "es", "ja"},
		EnableProviderFallback: true,
		EnableEnrichment:       false,
	}
}
