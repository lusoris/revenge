// Package sonarr provides a client for the Sonarr API v3.
//
// Sonarr is a PRIMARY metadata provider for TV shows in the Revenge media system.
// It manages TV series, episodes, and downloads, providing automated episode
// tracking and download management.
//
// # API Version
//
// This client is compatible with Sonarr API v3, which is used by Sonarr v3.x and v4.x.
//
// # Configuration
//
// The client requires a base URL and API key:
//
//	client := sonarr.NewClient(sonarr.Config{
//	    BaseURL: "http://localhost:8989",
//	    APIKey:  "your-api-key",
//	})
//
// # Rate Limiting
//
// The client includes built-in rate limiting to avoid overwhelming the Sonarr server.
// By default, it allows 10 requests per second with a burst of 20.
//
// # Caching
//
// The client implements a simple in-memory cache with a configurable TTL (default 5 minutes).
// This reduces load on the Sonarr server for frequently accessed data like series lists
// and quality profiles.
//
// # Error Handling
//
// The client returns specific errors for common scenarios:
//   - ErrSeriesNotFound: When a series is not found (404)
//   - ErrEpisodeNotFound: When an episode is not found (404)
//   - ErrEpisodeFileNotFound: When an episode file is not found (404)
//
// # Webhooks
//
// Sonarr can send webhooks for various events. Use the WebhookPayload type to parse
// incoming webhook payloads. Common events include:
//   - Grab: When a release is grabbed for download
//   - Download: When a download completes and is imported
//   - Rename: When files are renamed
//   - SeriesAdd/SeriesDelete: When series are added or deleted
//
// # Example Usage
//
//	// Get all series
//	series, err := client.GetAllSeries(ctx)
//
//	// Get episodes for a series
//	episodes, err := client.GetEpisodes(ctx, seriesID)
//
//	// Get calendar (upcoming episodes)
//	calendar, err := client.GetCalendar(ctx, time.Now(), time.Now().AddDate(0, 0, 7))
//
//	// Trigger a series refresh
//	cmd, err := client.RefreshSeries(ctx, seriesID)
package sonarr
