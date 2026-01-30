package stashdb

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/content/qar/expedition"
)

// Provider implements the expedition.MetadataProvider interface for StashDB.
// QAR obfuscation: This provides metadata for adult content from StashDB.
type Provider struct {
	client *Client
	logger *slog.Logger
}

// NewProvider creates a new StashDB metadata provider.
func NewProvider(client *Client, logger *slog.Logger) *Provider {
	return &Provider{
		client: client,
		logger: logger.With("provider", "stashdb"),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "stashdb"
}

// Priority returns the provider priority (lower = higher priority).
func (p *Provider) Priority() int {
	return 10 // Primary provider for QAR content
}

// IsAvailable returns true if the provider is configured and reachable.
func (p *Provider) IsAvailable() bool {
	if p.client == nil {
		return false
	}
	return p.client.IsConfigured()
}

// Ping checks if the provider is reachable.
func (p *Provider) Ping(ctx context.Context) error {
	return p.client.Ping(ctx)
}

// GetExpeditionMetadata retrieves metadata for an expedition by StashDB ID.
// QAR obfuscation: expedition = adult movie, charter = stashdb_id
func (p *Provider) GetExpeditionMetadata(ctx context.Context, charter string) (*expedition.Metadata, error) {
	p.logger.Debug("fetching expedition metadata", "charter", charter)

	scene, err := p.client.GetScene(ctx, charter)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		p.logger.Error("failed to fetch scene", "charter", charter, "error", err)
		return nil, err
	}

	return p.sceneToExpeditionMetadata(scene), nil
}

// MatchExpedition matches an expedition by title, year, and/or fingerprint.
// QAR obfuscation: coordinates = phash fingerprint
func (p *Provider) MatchExpedition(ctx context.Context, title string, launchYear int, coordinates string) (*expedition.Metadata, error) {
	p.logger.Debug("matching expedition",
		"title", title,
		"launch_year", launchYear,
		"has_coordinates", coordinates != "",
	)

	// Try fingerprint matching first (most accurate)
	if coordinates != "" {
		scenes, err := p.client.FindSceneByFingerprint(ctx, "PHASH", coordinates, 0)
		if err == nil && len(scenes) > 0 {
			p.logger.Debug("matched by fingerprint", "scene_id", scenes[0].ID)
			return p.sceneToExpeditionMetadata(&scenes[0]), nil
		}
		if err != nil && !errors.Is(err, ErrNotFound) {
			p.logger.Warn("fingerprint match failed", "error", err)
		}
	}

	// Fall back to title search
	if title != "" {
		results, err := p.client.SearchScenes(ctx, title, 1, 10)
		if err != nil {
			p.logger.Warn("title search failed", "error", err)
			return nil, err
		}

		// Find best match
		for _, scene := range results.Data {
			// Check year match if specified
			if launchYear > 0 && scene.Date != "" {
				sceneYear := extractYear(scene.Date)
				if sceneYear > 0 && sceneYear != launchYear {
					continue
				}
			}

			// Basic title similarity check
			if matchesTitle(title, scene.Title) {
				p.logger.Debug("matched by title", "scene_id", scene.ID, "scene_title", scene.Title)
				// Fetch full scene data
				fullScene, err := p.client.GetScene(ctx, scene.ID)
				if err != nil {
					return p.sceneToExpeditionMetadata(&scene), nil
				}
				return p.sceneToExpeditionMetadata(fullScene), nil
			}
		}
	}

	return nil, ErrNotFound
}

// sceneToExpeditionMetadata converts a StashDB scene to expedition metadata.
func (p *Provider) sceneToExpeditionMetadata(scene *Scene) *expedition.Metadata {
	if scene == nil {
		return nil
	}

	meta := &expedition.Metadata{
		Title:    scene.Title,
		Overview: scene.Details,
		Director: scene.Director,
		Charter:  scene.ID, // StashDB ID
	}

	// Parse launch date
	if scene.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", scene.ReleaseDate); err == nil {
			meta.LaunchDate = t
		}
	} else if scene.Date != "" {
		if t, err := time.Parse("2006-01-02", scene.Date); err == nil {
			meta.LaunchDate = t
		}
	}

	// Convert duration (seconds to ticks: 1 tick = 100ns, 1 second = 10_000_000 ticks)
	if scene.Duration > 0 {
		meta.RuntimeTicks = int64(scene.Duration) * 10_000_000
	}

	// Studio name (port in QAR terminology)
	if scene.Studio != nil {
		meta.PortName = scene.Studio.Name
	}

	// Series (from code/studio naming pattern)
	if scene.Code != "" {
		meta.Series = scene.Code
	}

	// Poster URL (first image)
	if len(scene.Images) > 0 {
		meta.PosterURL = scene.Images[0].URL
		if len(scene.Images) > 1 {
			meta.BackdropURL = scene.Images[1].URL
		}
	}

	return meta
}

// extractYear extracts the year from a date string (YYYY-MM-DD format).
func extractYear(date string) int {
	if len(date) < 4 {
		return 0
	}
	year := 0
	for i := 0; i < 4 && i < len(date); i++ {
		c := date[i]
		if c < '0' || c > '9' {
			return 0
		}
		year = year*10 + int(c-'0')
	}
	return year
}

// matchesTitle checks if two titles are similar enough to be a match.
func matchesTitle(query, candidate string) bool {
	// Normalize both strings
	q := strings.ToLower(strings.TrimSpace(query))
	c := strings.ToLower(strings.TrimSpace(candidate))

	// Exact match
	if q == c {
		return true
	}

	// One contains the other
	if strings.Contains(c, q) || strings.Contains(q, c) {
		return true
	}

	return false
}
