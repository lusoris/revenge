// Package metadata provides a unified metadata service that abstracts external metadata providers.
//
// This package implements Option B architecture where content modules (movie, tvshow, etc.)
// call the metadata service without knowing about external providers (TMDb, TVDb, etc.).
//
// # Architecture
//
// The service acts as a facade that:
//   - Aggregates multiple metadata providers (TMDb, TVDb, Fanart.tv, etc.)
//   - Provides centralized rate limiting and caching
//   - Handles provider fallback and error recovery
//   - Manages async metadata refresh via River jobs
//
// Content modules simply call methods like RefreshMovie() or GetSeriesMetadata()
// without any knowledge of where the data comes from.
//
// # Provider Interface
//
// All providers implement the Provider interface which defines methods for:
//   - Movie metadata (search, details, credits, images)
//   - TV series metadata (search, details, seasons, episodes)
//   - Person metadata
//   - Image URLs and downloads
//
// # Usage
//
//	// In content module (movie, tvshow):
//	err := metadataService.RefreshMovie(ctx, movieID)  // Triggers async job
//	metadata, err := metadataService.GetMovieMetadata(ctx, tmdbID, []string{"en", "de"})
//
//	// Service handles provider selection, fallback, caching internally
//
// # Supported Providers
//
//   - TMDb (primary): Movies, TV shows, people, images
//   - TVDb: TV shows, episodes (better episode ordering)
//   - More providers can be added (Fanart.tv, OMDb, MusicBrainz)
//
// # Packages
//
//   - providers/tmdb: TMDb provider implementation
//   - providers/tvdb: TVDb provider implementation
//   - jobs: River job workers for async metadata updates
package metadata
