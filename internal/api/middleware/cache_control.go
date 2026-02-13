package middleware

import (
	"net/http"

	"github.com/ogen-go/ogen/middleware"
)

// cacheableOperations lists operation IDs whose responses contain catalog data
// that is identical for all authorised users and can be briefly cached by
// the browser.  Everything else (user state, admin, auth, mutations) gets
// no-store.
var cacheableOperations = map[string]struct{}{
	// Movies
	"listMovies": {}, "searchMovies": {}, "autocompleteMovies": {},
	"getMovie": {}, "getMovieFiles": {}, "getMovieCast": {},
	"getMovieCrew": {}, "getMovieGenres": {}, "getMovieExternalIDs": {},
	"getMovieCollection": {}, "getSimilarMovies": {},
	"getCollection": {}, "getCollectionMovies": {},

	// TV Shows
	"listTVShows": {}, "searchTVShows": {}, "autocompleteTVShows": {},
	"getTVShow": {}, "getTVShowSeasons": {}, "getTVShowEpisodes": {},
	"getTVShowCast": {}, "getTVShowCrew": {}, "getTVShowGenres": {},
	"getTVShowContentRatings": {}, "getTVShowNetworks": {}, "getTVShowExternalIDs": {},
	"getTVSeason": {}, "getTVSeasonEpisodes": {},
	"getTVEpisode": {}, "getTVEpisodeFiles": {},

	// Metadata (external provider data)
	"getMovieMetadata": {}, "getMovieMetadataCredits": {}, "getMovieMetadataImages": {},
	"getMovieRecommendationsMetadata": {}, "getSimilarMoviesMetadata": {},
	"getTVShowMetadata": {}, "getTVShowMetadataCredits": {}, "getTVShowMetadataImages": {},
	"getSeasonMetadata": {}, "getSeasonMetadataImages": {},
	"getEpisodeMetadata": {}, "getEpisodeMetadataImages": {},
	"getCollectionMetadata": {},
	"getPersonMetadata":     {}, "getPersonMetadataCredits": {}, "getPersonMetadataImages": {},
	"searchPersonMetadata": {}, "searchMoviesMetadata": {}, "searchTVShowsMetadata": {},

	// Search facets & libraries
	"getSearchFacets": {}, "getTVShowSearchFacets": {},
	"listLibraries": {}, "getLibrary": {}, "listGenres": {},
	"searchLibraryMovies": {}, "searchLibraryTVShows": {},

	// Reference data
	"listMetadataProviders": {},

	// Proxied images
	"getProxiedImage": {},
}

// healthOperations are health probe endpoints that should be fresh
// but may use conditional caching (304 Not Modified).
var healthOperations = map[string]struct{}{
	"getLiveness":  {},
	"getReadiness": {},
	"getStartup":   {},
}

// CacheControlMiddleware sets Cache-Control headers based on the operation type:
//   - Catalog/content endpoints: private, max-age=60  (browser-cacheable)
//   - Health probes:             no-cache              (revalidate every request)
//   - User-specific GETs:       private, no-store
//   - Mutations (POST/PUT/…):   no-store
func CacheControlMiddleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		w, ok := GetResponseWriter(req.Context)
		if !ok {
			return next(req)
		}

		var value string
		switch {
		case req.Raw.Method != http.MethodGet:
			value = "no-store"
		case isCacheable(req.OperationID):
			value = "private, max-age=60"
		case isHealthProbe(req.OperationID):
			value = "no-cache"
		default:
			value = "private, no-store"
		}

		w.Header().Set("Cache-Control", value)

		return next(req)
	}
}

func isCacheable(operationID string) bool {
	_, ok := cacheableOperations[operationID]
	return ok
}

func isHealthProbe(operationID string) bool {
	_, ok := healthOperations[operationID]
	return ok
}
