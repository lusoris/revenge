package stashdb

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/lusoris/revenge/internal/content/qar/voyage"
)

// VoyageProvider implements the voyage.MetadataProvider interface for StashDB.
// QAR obfuscation: This provides metadata for adult scenes from StashDB.
type VoyageProvider struct {
	client *Client
	logger *slog.Logger
}

// NewVoyageProvider creates a new StashDB metadata provider for voyages.
func NewVoyageProvider(client *Client, logger *slog.Logger) *VoyageProvider {
	return &VoyageProvider{
		client: client,
		logger: logger.With("provider", "stashdb", "type", "voyage"),
	}
}

// IsAvailable returns true if the provider is configured and reachable.
func (p *VoyageProvider) IsAvailable() bool {
	if p.client == nil {
		return false
	}
	return p.client.IsConfigured()
}

// GetVoyageMetadata retrieves metadata for a voyage by StashDB ID.
// QAR obfuscation: voyage = adult scene, charter = stashdb_id
func (p *VoyageProvider) GetVoyageMetadata(ctx context.Context, charter string) (*voyage.Metadata, error) {
	p.logger.Debug("fetching voyage metadata", "charter", charter)

	scene, err := p.client.GetScene(ctx, charter)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		p.logger.Error("failed to fetch scene", "charter", charter, "error", err)
		return nil, err
	}

	return p.sceneToVoyageMetadata(scene), nil
}

// MatchVoyage matches a voyage by title, date, and/or fingerprints.
// QAR obfuscation: coordinates = phash, oshash = oshash fingerprint
func (p *VoyageProvider) MatchVoyage(ctx context.Context, title string, launchDate *time.Time, coordinates string, oshash string) (*voyage.Metadata, error) {
	p.logger.Debug("matching voyage",
		"title", title,
		"has_launch_date", launchDate != nil,
		"has_coordinates", coordinates != "",
		"has_oshash", oshash != "",
	)

	// Try fingerprint matching first (most accurate)
	// StashDB supports PHASH and OSHASH
	if coordinates != "" || oshash != "" {
		var fingerprints []Fingerprint
		if coordinates != "" {
			fingerprints = append(fingerprints, Fingerprint{Algorithm: "PHASH", Hash: coordinates})
		}
		if oshash != "" {
			fingerprints = append(fingerprints, Fingerprint{Algorithm: "OSHASH", Hash: oshash})
		}

		scenes, err := p.client.FindScenesByFingerprints(ctx, fingerprints)
		if err == nil && len(scenes) > 0 {
			p.logger.Debug("matched by fingerprint", "scene_id", scenes[0].ID)
			return p.sceneToVoyageMetadata(&scenes[0]), nil
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
			// Check date match if specified
			if launchDate != nil && scene.Date != "" {
				sceneDate := parseDate(scene.Date)
				if sceneDate != nil && !sameDay(*sceneDate, *launchDate) {
					continue
				}
			}

			// Basic title similarity check
			if matchesTitle(title, scene.Title) {
				p.logger.Debug("matched by title", "scene_id", scene.ID, "scene_title", scene.Title)
				// Fetch full scene data
				fullScene, err := p.client.GetScene(ctx, scene.ID)
				if err != nil {
					return p.sceneToVoyageMetadata(&scene), nil
				}
				return p.sceneToVoyageMetadata(fullScene), nil
			}
		}
	}

	return nil, ErrNotFound
}

// sceneToVoyageMetadata converts a StashDB scene to voyage metadata.
func (p *VoyageProvider) sceneToVoyageMetadata(scene *Scene) *voyage.Metadata {
	if scene == nil {
		return nil
	}

	meta := &voyage.Metadata{
		Title:    scene.Title,
		Overview: scene.Details,
		Charter:  scene.ID, // StashDB ID
	}

	// Parse launch date
	if scene.ReleaseDate != "" {
		if t := parseDate(scene.ReleaseDate); t != nil {
			meta.LaunchDate = *t
		}
	} else if scene.Date != "" {
		if t := parseDate(scene.Date); t != nil {
			meta.LaunchDate = *t
		}
	}

	// Duration in minutes (distance in QAR terminology)
	if scene.Duration > 0 {
		meta.Distance = scene.Duration / 60
	}

	// Studio name (port in QAR terminology)
	if scene.Studio != nil {
		meta.PortName = scene.Studio.Name
	}

	// Cover URL (first image)
	if len(scene.Images) > 0 {
		meta.CoverURL = scene.Images[0].URL
	}

	return meta
}

// parseDate parses a date string in YYYY-MM-DD format.
func parseDate(date string) *time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	return &t
}

// sameDay checks if two dates are the same day.
func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
