// Package radarr provides integration with Radarr for movie management.
//
// Radarr is a PRIMARY metadata provider in the Revenge media system.
// This means data from Radarr is preferred over external APIs like TMDb.
//
// Architecture:
//   - Client: HTTP client for Radarr API v3
//   - Types: Radarr API response types
//   - Mapper: Converts Radarr types to domain types
//   - Service: Business logic for syncing with local database
//
// Usage:
//
//	client := radarr.NewClient(radarr.Config{
//	    BaseURL: "http://localhost:7878",
//	    APIKey:  "your-api-key",
//	})
//
//	movies, err := client.GetAllMovies(ctx)
//	if err != nil {
//	    return err
//	}
//
// See also:
//   - docs/dev/design/integrations/servarr/RADARR.md
//   - docs/dev/design/architecture/03_METADATA_SYSTEM.md
package radarr
