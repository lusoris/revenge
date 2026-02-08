// Package metadatafx provides fx dependency injection for the metadata service.
package metadatafx

import (
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/service/metadata"
	movieadapter "github.com/lusoris/revenge/internal/service/metadata/adapters/movie"
	tvshowadapter "github.com/lusoris/revenge/internal/service/metadata/adapters/tvshow"
	"github.com/lusoris/revenge/internal/service/metadata/providers/anidb"
	"github.com/lusoris/revenge/internal/service/metadata/providers/anilist"
	"github.com/lusoris/revenge/internal/service/metadata/providers/fanarttv"
	"github.com/lusoris/revenge/internal/service/metadata/providers/kitsu"
	"github.com/lusoris/revenge/internal/service/metadata/providers/letterboxd"
	"github.com/lusoris/revenge/internal/service/metadata/providers/mal"
	"github.com/lusoris/revenge/internal/service/metadata/providers/omdb"
	"github.com/lusoris/revenge/internal/service/metadata/providers/simkl"
	"github.com/lusoris/revenge/internal/service/metadata/providers/tmdb"
	"github.com/lusoris/revenge/internal/service/metadata/providers/trakt"
	"github.com/lusoris/revenge/internal/service/metadata/providers/tvdb"
	"github.com/lusoris/revenge/internal/service/metadata/providers/tvmaze"

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

	// OMDb configuration (optional)
	OMDbAPIKey string

	// TVmaze configuration (optional, no API key needed)
	TVmazeEnabled bool

	// AniList configuration (optional, no API key needed)
	AniListEnabled bool

	// Kitsu configuration (optional, no API key needed)
	KitsuEnabled bool

	// AniDB configuration (optional)
	AniDBEnabled       bool
	AniDBClientName    string
	AniDBClientVersion int

	// MAL configuration (optional)
	MALEnabled  bool
	MALClientID string

	// Trakt configuration (optional)
	TraktEnabled  bool
	TraktClientID string

	// Simkl configuration (optional)
	SimklEnabled  bool
	SimklClientID string

	// Letterboxd configuration (optional)
	LetterboxdEnabled   bool
	LetterboxdAPIKey    string
	LetterboxdAPISecret string
}

// ModuleParams contains parameters for the metadata module.
type ModuleParams struct {
	fx.In

	Config         Config          `optional:"true"`
	TMDbConfig     tmdb.Config     `optional:"true"`
	TVDbConfig     tvdb.Config     `optional:"true"`
	FanartTVConfig fanarttv.Config `optional:"true"`
	OMDbConfig     omdb.Config     `optional:"true"`
	TVmazeConfig   tvmaze.Config   `optional:"true"`
	AniListConfig  anilist.Config  `optional:"true"`
	KitsuConfig    kitsu.Config    `optional:"true"`
	AniDBConfig    anidb.Config    `optional:"true"`
	MALConfig      mal.Config      `optional:"true"`
	TraktConfig       trakt.Config       `optional:"true"`
	SimklConfig       simkl.Config       `optional:"true"`
	LetterboxdConfig  letterboxd.Config  `optional:"true"`
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
	OMDbProvider          *omdb.Provider     `optional:"true"`
	TVmazeProvider        *tvmaze.Provider   `optional:"true"`
	AniListProvider       *anilist.Provider   `optional:"true"`
	KitsuProvider         *kitsu.Provider     `optional:"true"`
	AniDBProvider         *anidb.Provider     `optional:"true"`
	MALProvider           *mal.Provider       `optional:"true"`
	TraktProvider         *trakt.Provider      `optional:"true"`
	SimklProvider         *simkl.Provider      `optional:"true"`
	LetterboxdProvider    *letterboxd.Provider `optional:"true"`
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

	// Create and register OMDb provider if configured
	omdbConfig := params.OMDbConfig
	if omdbConfig.APIKey == "" && params.Config.OMDbAPIKey != "" {
		omdbConfig = omdb.Config{
			APIKey: params.Config.OMDbAPIKey,
		}
	}

	if omdbConfig.APIKey != "" {
		omdbProvider, err := omdb.NewProvider(omdbConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(omdbProvider)
		result.OMDbProvider = omdbProvider
	}

	// Create and register TVmaze provider if enabled (no API key needed)
	if params.Config.TVmazeEnabled {
		tvmazeProvider, err := tvmaze.NewProvider(params.TVmazeConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(tvmazeProvider)
		result.TVmazeProvider = tvmazeProvider
	}

	// Create and register AniList provider if enabled (no API key needed)
	if params.AniListConfig.Enabled || params.Config.AniListEnabled {
		anilistConfig := params.AniListConfig
		if !anilistConfig.Enabled {
			anilistConfig.Enabled = true
		}
		anilistProvider, err := anilist.NewProvider(anilistConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(anilistProvider)
		result.AniListProvider = anilistProvider
	}

	// Create and register Kitsu provider if enabled (no API key needed)
	if params.KitsuConfig.Enabled || params.Config.KitsuEnabled {
		kitsuConfig := params.KitsuConfig
		if !kitsuConfig.Enabled {
			kitsuConfig.Enabled = true
		}
		kitsuProvider, err := kitsu.NewProvider(kitsuConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(kitsuProvider)
		result.KitsuProvider = kitsuProvider
	}

	// Create and register AniDB provider if configured
	anidbConfig := params.AniDBConfig
	if anidbConfig.ClientName == "" && params.Config.AniDBClientName != "" {
		anidbConfig = anidb.Config{
			Enabled:       true,
			ClientName:    params.Config.AniDBClientName,
			ClientVersion: params.Config.AniDBClientVersion,
		}
	}
	if anidbConfig.ClientName != "" || params.Config.AniDBEnabled {
		if !anidbConfig.Enabled {
			anidbConfig.Enabled = true
		}
		anidbProvider, err := anidb.NewProvider(anidbConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(anidbProvider)
		result.AniDBProvider = anidbProvider
	}

	// Create and register MAL provider if configured
	malConfig := params.MALConfig
	if malConfig.ClientID == "" && params.Config.MALClientID != "" {
		malConfig = mal.Config{
			Enabled:  true,
			ClientID: params.Config.MALClientID,
		}
	}
	if malConfig.ClientID != "" || params.Config.MALEnabled {
		if !malConfig.Enabled {
			malConfig.Enabled = true
		}
		malProvider, err := mal.NewProvider(malConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(malProvider)
		result.MALProvider = malProvider
	}

	// Create and register Trakt provider if configured
	traktConfig := params.TraktConfig
	if traktConfig.ClientID == "" && params.Config.TraktClientID != "" {
		traktConfig = trakt.Config{
			Enabled:  true,
			ClientID: params.Config.TraktClientID,
		}
	}
	if traktConfig.ClientID != "" || params.Config.TraktEnabled {
		if !traktConfig.Enabled {
			traktConfig.Enabled = true
		}
		traktProvider, err := trakt.NewProvider(traktConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(traktProvider)
		result.TraktProvider = traktProvider
	}

	// Create and register Simkl provider if configured
	simklConfig := params.SimklConfig
	if simklConfig.ClientID == "" && params.Config.SimklClientID != "" {
		simklConfig = simkl.Config{
			Enabled:  true,
			ClientID: params.Config.SimklClientID,
		}
	}
	if simklConfig.ClientID != "" || params.Config.SimklEnabled {
		if !simklConfig.Enabled {
			simklConfig.Enabled = true
		}
		simklProvider, err := simkl.NewProvider(simklConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(simklProvider)
		result.SimklProvider = simklProvider
	}

	// Create and register Letterboxd provider if configured
	letterboxdConfig := params.LetterboxdConfig
	if letterboxdConfig.APIKey == "" && params.Config.LetterboxdAPIKey != "" {
		letterboxdConfig = letterboxd.Config{
			Enabled:   true,
			APIKey:    params.Config.LetterboxdAPIKey,
			APISecret: params.Config.LetterboxdAPISecret,
		}
	}
	if (letterboxdConfig.APIKey != "" && letterboxdConfig.APISecret != "") || params.Config.LetterboxdEnabled {
		if !letterboxdConfig.Enabled {
			letterboxdConfig.Enabled = true
		}
		letterboxdProvider, err := letterboxd.NewProvider(letterboxdConfig)
		if err != nil {
			return ModuleResult{}, err
		}
		svc.RegisterProvider(letterboxdProvider)
		result.LetterboxdProvider = letterboxdProvider
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
