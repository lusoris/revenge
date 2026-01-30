// Package search provides Typesense search collections for QAR (adult content).
// These collections are isolated from main content search.
package search

import (
	"context"
	"log/slog"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// QAR collection names (isolated from main content).
const (
	CollectionQARExpeditions = "qar_expeditions" // Movies
	CollectionQARVoyages     = "qar_voyages"     // Scenes
	CollectionQARCrew        = "qar_crew"        // Performers
	CollectionQARPorts       = "qar_ports"       // Studios
	CollectionQARFlags       = "qar_flags"       // Tags
)

// QARCollectionSchemas returns all QAR search collection schemas.
func QARCollectionSchemas() []*api.CollectionSchema {
	return []*api.CollectionSchema{
		QARExpeditionsSchema(),
		QARVoyagesSchema(),
		QARCrewSchema(),
		QARPortsSchema(),
		QARFlagsSchema(),
	}
}

// QARExpeditionsSchema returns the schema for adult movies (expeditions).
func QARExpeditionsSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: CollectionQARExpeditions,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string", Facet: ptr(false), Index: ptr(true)},
			{Name: "sort_title", Type: "string", Optional: ptr(true)},
			{Name: "overview", Type: "string", Optional: ptr(true)},
			{Name: "launch_year", Type: "int32", Optional: ptr(true), Facet: ptr(true)},
			{Name: "distance", Type: "int32", Optional: ptr(true)}, // runtime_minutes
			{Name: "port_id", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "port_name", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "crew_names", Type: "string[]", Optional: ptr(true), Facet: ptr(true)},
			{Name: "crew_ids", Type: "string[]", Optional: ptr(true)},
			{Name: "flag_names", Type: "string[]", Optional: ptr(true), Facet: ptr(true)},
			{Name: "flag_ids", Type: "string[]", Optional: ptr(true)},
			{Name: "charter", Type: "string", Optional: ptr(true)},  // stashdb_id
			{Name: "registry", Type: "string", Optional: ptr(true)}, // tpdb_id
			{Name: "cover_path", Type: "string", Optional: ptr(true)},
			{Name: "created_at", Type: "int64"}, // Unix timestamp for sorting
		},
		DefaultSortingField: ptr("created_at"),
		TokenSeparators:     &[]string{"-", "_", "."},
		SymbolsToIndex:      &[]string{"+"},
	}
}

// QARVoyagesSchema returns the schema for adult scenes (voyages).
func QARVoyagesSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: CollectionQARVoyages,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string", Facet: ptr(false), Index: ptr(true)},
			{Name: "overview", Type: "string", Optional: ptr(true)},
			{Name: "launch_year", Type: "int32", Optional: ptr(true), Facet: ptr(true)},
			{Name: "distance", Type: "int32", Optional: ptr(true)}, // runtime_minutes
			{Name: "port_id", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "port_name", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "crew_names", Type: "string[]", Optional: ptr(true), Facet: ptr(true)},
			{Name: "crew_ids", Type: "string[]", Optional: ptr(true)},
			{Name: "flag_names", Type: "string[]", Optional: ptr(true), Facet: ptr(true)},
			{Name: "flag_ids", Type: "string[]", Optional: ptr(true)},
			{Name: "charter", Type: "string", Optional: ptr(true)},     // stashdb_id
			{Name: "registry", Type: "string", Optional: ptr(true)},    // tpdb_id
			{Name: "coordinates", Type: "string", Optional: ptr(true)}, // phash
			{Name: "oshash", Type: "string", Optional: ptr(true)},
			{Name: "cover_path", Type: "string", Optional: ptr(true)},
			{Name: "created_at", Type: "int64"}, // Unix timestamp for sorting
		},
		DefaultSortingField: ptr("created_at"),
		TokenSeparators:     &[]string{"-", "_", "."},
		SymbolsToIndex:      &[]string{"+"},
	}
}

// QARCrewSchema returns the schema for adult performers (crew).
func QARCrewSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: CollectionQARCrew,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string", Facet: ptr(false), Index: ptr(true)},
			{Name: "aliases", Type: "string[]", Optional: ptr(true)},
			{Name: "disambiguation", Type: "string", Optional: ptr(true)},
			{Name: "gender", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "origin", Type: "string", Optional: ptr(true), Facet: ptr(true)},      // ethnicity
			{Name: "nationality", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "rigging", Type: "string", Optional: ptr(true), Facet: ptr(true)},     // hair_color
			{Name: "compass", Type: "string", Optional: ptr(true), Facet: ptr(true)},     // eye_color
			{Name: "height_cm", Type: "int32", Optional: ptr(true)},
			{Name: "christening_year", Type: "int32", Optional: ptr(true)},              // birth_year
			{Name: "maiden_voyage", Type: "int32", Optional: ptr(true), Facet: ptr(true)}, // career_start
			{Name: "last_port", Type: "int32", Optional: ptr(true)},                      // career_end
			{Name: "voyage_count", Type: "int32", Optional: ptr(true)},                   // scene_count
			{Name: "charter", Type: "string", Optional: ptr(true)},                       // stashdb_id
			{Name: "registry", Type: "string", Optional: ptr(true)},                      // tpdb_id
			{Name: "image_path", Type: "string", Optional: ptr(true)},
			{Name: "created_at", Type: "int64"}, // Unix timestamp for sorting
		},
		DefaultSortingField: ptr("created_at"),
		TokenSeparators:     &[]string{"-", "_", "."},
		SymbolsToIndex:      &[]string{"+"},
	}
}

// QARPortsSchema returns the schema for adult studios (ports).
func QARPortsSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: CollectionQARPorts,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string", Facet: ptr(false), Index: ptr(true)},
			{Name: "parent_id", Type: "string", Optional: ptr(true)},
			{Name: "parent_name", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "stashdb_id", Type: "string", Optional: ptr(true)},
			{Name: "tpdb_id", Type: "string", Optional: ptr(true)},
			{Name: "url", Type: "string", Optional: ptr(true)},
			{Name: "logo_path", Type: "string", Optional: ptr(true)},
			{Name: "voyage_count", Type: "int32", Optional: ptr(true)},     // scene_count
			{Name: "expedition_count", Type: "int32", Optional: ptr(true)}, // movie_count
			{Name: "created_at", Type: "int64"}, // Unix timestamp for sorting
		},
		DefaultSortingField: ptr("created_at"),
		TokenSeparators:     &[]string{"-", "_", "."},
		SymbolsToIndex:      &[]string{"+"},
	}
}

// QARFlagsSchema returns the schema for adult tags (flags).
func QARFlagsSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: CollectionQARFlags,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string", Facet: ptr(false), Index: ptr(true)},
			{Name: "description", Type: "string", Optional: ptr(true)},
			{Name: "parent_id", Type: "string", Optional: ptr(true)},
			{Name: "parent_name", Type: "string", Optional: ptr(true), Facet: ptr(true)},
			{Name: "waters", Type: "string", Optional: ptr(true), Facet: ptr(true)}, // category
			{Name: "stashdb_id", Type: "string", Optional: ptr(true)},
			{Name: "usage_count", Type: "int32", Optional: ptr(true)}, // How many items have this tag
			{Name: "created_at", Type: "int64"}, // Unix timestamp for sorting
		},
		DefaultSortingField: ptr("created_at"),
		TokenSeparators:     &[]string{"-", "_", "."},
		SymbolsToIndex:      &[]string{"+"},
	}
}

// InitQARCollections creates all QAR search collections if they don't exist.
func (c *Client) InitQARCollections(ctx context.Context) error {
	schemas := QARCollectionSchemas()

	for _, schema := range schemas {
		_, err := c.CreateCollection(ctx, schema)
		if err != nil {
			// Ignore "already exists" errors
			c.logger.Debug("collection creation result",
				slog.String("collection", schema.Name),
				slog.String("error", err.Error()),
			)
		} else {
			c.logger.Info("created QAR search collection",
				slog.String("collection", schema.Name),
			)
		}
	}

	return nil
}

// ptr is a helper to create pointers for optional fields.
func ptr[T any](v T) *T {
	return &v
}

// =============================================================================
// QAR Search Documents
// =============================================================================

// QARExpeditionDoc represents an expedition (movie) in the search index.
type QARExpeditionDoc struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	SortTitle  string   `json:"sort_title,omitempty"`
	Overview   string   `json:"overview,omitempty"`
	LaunchYear int32    `json:"launch_year,omitempty"`
	Distance   int32    `json:"distance,omitempty"` // runtime_minutes
	PortID     string   `json:"port_id,omitempty"`
	PortName   string   `json:"port_name,omitempty"`
	CrewNames  []string `json:"crew_names,omitempty"`
	CrewIDs    []string `json:"crew_ids,omitempty"`
	FlagNames  []string `json:"flag_names,omitempty"`
	FlagIDs    []string `json:"flag_ids,omitempty"`
	Charter    string   `json:"charter,omitempty"`  // stashdb_id
	Registry   string   `json:"registry,omitempty"` // tpdb_id
	CoverPath  string   `json:"cover_path,omitempty"`
	CreatedAt  int64    `json:"created_at"`
}

// QARVoyageDoc represents a voyage (scene) in the search index.
type QARVoyageDoc struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Overview    string   `json:"overview,omitempty"`
	LaunchYear  int32    `json:"launch_year,omitempty"`
	Distance    int32    `json:"distance,omitempty"` // runtime_minutes
	PortID      string   `json:"port_id,omitempty"`
	PortName    string   `json:"port_name,omitempty"`
	CrewNames   []string `json:"crew_names,omitempty"`
	CrewIDs     []string `json:"crew_ids,omitempty"`
	FlagNames   []string `json:"flag_names,omitempty"`
	FlagIDs     []string `json:"flag_ids,omitempty"`
	Charter     string   `json:"charter,omitempty"`     // stashdb_id
	Registry    string   `json:"registry,omitempty"`    // tpdb_id
	Coordinates string   `json:"coordinates,omitempty"` // phash
	Oshash      string   `json:"oshash,omitempty"`
	CoverPath   string   `json:"cover_path,omitempty"`
	CreatedAt   int64    `json:"created_at"`
}

// QARCrewDoc represents a crew member (performer) in the search index.
type QARCrewDoc struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Aliases         []string `json:"aliases,omitempty"`
	Disambiguation  string   `json:"disambiguation,omitempty"`
	Gender          string   `json:"gender,omitempty"`
	Origin          string   `json:"origin,omitempty"` // ethnicity
	Nationality     string   `json:"nationality,omitempty"`
	Rigging         string   `json:"rigging,omitempty"` // hair_color
	Compass         string   `json:"compass,omitempty"` // eye_color
	HeightCM        int32    `json:"height_cm,omitempty"`
	ChristeningYear int32    `json:"christening_year,omitempty"` // birth_year
	MaidenVoyage    int32    `json:"maiden_voyage,omitempty"`    // career_start
	LastPort        int32    `json:"last_port,omitempty"`        // career_end
	VoyageCount     int32    `json:"voyage_count,omitempty"`     // scene_count
	Charter         string   `json:"charter,omitempty"`          // stashdb_id
	Registry        string   `json:"registry,omitempty"`         // tpdb_id
	ImagePath       string   `json:"image_path,omitempty"`
	CreatedAt       int64    `json:"created_at"`
}

// QARPortDoc represents a port (studio) in the search index.
type QARPortDoc struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ParentID        string `json:"parent_id,omitempty"`
	ParentName      string `json:"parent_name,omitempty"`
	StashDBID       string `json:"stashdb_id,omitempty"`
	TPDBID          string `json:"tpdb_id,omitempty"`
	URL             string `json:"url,omitempty"`
	LogoPath        string `json:"logo_path,omitempty"`
	VoyageCount     int32  `json:"voyage_count,omitempty"`     // scene_count
	ExpeditionCount int32  `json:"expedition_count,omitempty"` // movie_count
	CreatedAt       int64  `json:"created_at"`
}

// QARFlagDoc represents a flag (tag) in the search index.
type QARFlagDoc struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	ParentName  string `json:"parent_name,omitempty"`
	Waters      string `json:"waters,omitempty"` // category
	StashDBID   string `json:"stashdb_id,omitempty"`
	UsageCount  int32  `json:"usage_count,omitempty"`
	CreatedAt   int64  `json:"created_at"`
}
