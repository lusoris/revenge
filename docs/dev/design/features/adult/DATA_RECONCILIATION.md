---
sources:
  - name: Casbin
    url: ../../../sources/security/casbin.md
    note: Auto-resolved from casbin
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Adult Data Reconciliation](#adult-data-reconciliation)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Reconciliation settings](#reconciliation-settings)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Adult Data Reconciliation


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Scenes, Performers, Studios

> Fuzzy logic and confidence scoring for conflicting metadata

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | 🟡 | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Database Schema

**Schema**: `qar`

<!-- Schema diagram -->

### Module Structure

```
internal/content/adult_data_reconciliation/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── adult_data_reconciliation_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### File Structure

**Key Files**:
- `internal/content/adult/reconciliation/service.go` - Core reconciliation logic
- `internal/content/adult/reconciliation/fuzzy/*.go` - Fuzzy matching algorithms
- `internal/content/adult/reconciliation/repository.go` - Metadata source tracking
- `migrations/qar/014_reconciliation.sql` - Database schema


### Key Interfaces

```go
// ReconciliationService handles metadata conflict resolution
type ReconciliationService interface {
  // ReconcileScene merges scene metadata from multiple sources
  ReconcileScene(ctx context.Context, sceneID uuid.UUID, sources []MetadataSource) (*Scene, error)

  // ReconcilePerformer merges performer metadata from multiple sources
  ReconcilePerformer(ctx context.Context, performerID uuid.UUID, sources []MetadataSource) (*Performer, error)

  // ReconcileStudio merges studio metadata from multiple sources
  ReconcileStudio(ctx context.Context, studioID uuid.UUID, sources []MetadataSource) (*Studio, error)

  // DetectConflicts finds fields with conflicting values
  DetectConflicts(ctx context.Context, contentType string, contentID uuid.UUID) ([]Conflict, error)

  // AutoResolve automatically resolves conflicts based on rules
  AutoResolve(ctx context.Context, conflictID uuid.UUID) error

  // ApplyManualOverride sets user-corrected value for a field
  ApplyManualOverride(ctx context.Context, override ManualOverride) error

  // GetMetadataHistory shows all source values for a field
  GetMetadataHistory(ctx context.Context, contentType string, contentID uuid.UUID, field string) ([]MetadataSource, error)
}

// ReconciliationRepository handles database operations
type ReconciliationRepository interface {
  // Rules
  GetRules(ctx context.Context, contentType string) ([]ReconciliationRule, error)
  CreateRule(ctx context.Context, rule ReconciliationRule) error
  UpdateRule(ctx context.Context, rule ReconciliationRule) error

  // Source tracking
  TrackMetadataSource(ctx context.Context, source MetadataSource) error
  GetMetadataSources(ctx context.Context, contentType string, contentID uuid.UUID, field string) ([]MetadataSource, error)
  SetActiveSource(ctx context.Context, sourceID uuid.UUID) error

  // Conflicts
  FindConflicts(ctx context.Context, contentType string, contentID uuid.UUID) ([]Conflict, error)
  CreateConflict(ctx context.Context, conflict Conflict) error
  ResolveConflict(ctx context.Context, conflictID uuid.UUID, resolution ConflictResolution) error

  // Manual overrides
  GetManualOverride(ctx context.Context, contentType string, contentID uuid.UUID, field string) (*ManualOverride, error)
  SetManualOverride(ctx context.Context, override ManualOverride) error
  DeleteManualOverride(ctx context.Context, overrideID uuid.UUID) error

  // Fuzzy matching cache
  GetFuzzyMatch(ctx context.Context, a, b string, algorithm string) (*float64, error)
  CacheFuzzyMatch(ctx context.Context, a, b string, algorithm string, score float64) error
}

// FuzzyMatcher performs string similarity matching
type FuzzyMatcher interface {
  // Levenshtein calculates edit distance (0.0-1.0)
  Levenshtein(a, b string) float64

  // JaroWinkler calculates Jaro-Winkler similarity (0.0-1.0)
  JaroWinkler(a, b string) float64

  // Trigram calculates trigram similarity using PostgreSQL pg_trgm (0.0-1.0)
  Trigram(ctx context.Context, a, b string) (float64, error)

  // BestMatch finds the best matching string from candidates
  BestMatch(ctx context.Context, target string, candidates []string, threshold float64) (*string, float64, error)
}

// Types
type ReconciliationRule struct {
  ID              uuid.UUID       `db:"id" json:"id"`
  ContentType     string          `db:"content_type" json:"content_type"`     // 'scene', 'performer', 'studio'
  FieldName       string          `db:"field_name" json:"field_name"`
  SourcePriority  []string        `db:"source_priority" json:"source_priority"` // ['stashdb', 'whisparr', 'stash', 'manual']
  FuzzyThreshold  float64         `db:"fuzzy_threshold" json:"fuzzy_threshold"`
  AutoMerge       bool            `db:"auto_merge" json:"auto_merge"`
  CreatedAt       time.Time       `db:"created_at" json:"created_at"`
  UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

type MetadataSource struct {
  ID              uuid.UUID       `db:"id" json:"id"`
  ContentType     string          `db:"content_type" json:"content_type"`
  ContentID       uuid.UUID       `db:"content_id" json:"content_id"`
  FieldName       string          `db:"field_name" json:"field_name"`
  SourceName      string          `db:"source_name" json:"source_name"`       // 'stashdb', 'whisparr', 'stash', 'manual'
  SourceID        *string         `db:"source_id" json:"source_id,omitempty"`
  FieldValue      json.RawMessage `db:"field_value" json:"field_value"`
  ConfidenceScore *float64        `db:"confidence_score" json:"confidence_score,omitempty"`
  IsActive        bool            `db:"is_active" json:"is_active"`
  CreatedAt       time.Time       `db:"created_at" json:"created_at"`
  UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

type Conflict struct {
  ID              uuid.UUID       `db:"id" json:"id"`
  ContentType     string          `db:"content_type" json:"content_type"`
  ContentID       uuid.UUID       `db:"content_id" json:"content_id"`
  FieldName       string          `db:"field_name" json:"field_name"`
  ConflictValues  json.RawMessage `db:"conflict_values" json:"conflict_values"` // [{source, value, confidence}]
  ResolutionStatus string         `db:"resolution_status" json:"resolution_status"`
  ResolvedValue   *json.RawMessage `db:"resolved_value" json:"resolved_value,omitempty"`
  ResolvedBy      *uuid.UUID      `db:"resolved_by" json:"resolved_by,omitempty"`
  ResolvedAt      *time.Time      `db:"resolved_at" json:"resolved_at,omitempty"`
  CreatedAt       time.Time       `db:"created_at" json:"created_at"`
}

type ManualOverride struct {
  ID            uuid.UUID       `db:"id" json:"id"`
  ContentType   string          `db:"content_type" json:"content_type"`
  ContentID     uuid.UUID       `db:"content_id" json:"content_id"`
  FieldName     string          `db:"field_name" json:"field_name"`
  OverrideValue json.RawMessage `db:"override_value" json:"override_value"`
  Reason        *string         `db:"reason" json:"reason,omitempty"`
  CreatedBy     uuid.UUID       `db:"created_by" json:"created_by"`
  CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid` - UUID handling
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/maypok86/otter` - L1 in-memory cache
- `github.com/riverqueue/river` - Background job queue
- `go.uber.org/fx` - Dependency injection
- `go.uber.org/zap` - Structured logging
- `github.com/agnivade/levenshtein` - Edit distance algorithm
- `github.com/xrash/smetrics` - Jaro-Winkler similarity

**PostgreSQL Extensions**:
- `pg_trgm` - Trigram similarity matching

**External APIs**:
- None (internal reconciliation logic)






## Configuration
### Environment Variables

```bash
# Reconciliation settings
RECONCILIATION_DEFAULT_FUZZY_THRESHOLD=0.85  # Default similarity threshold
RECONCILIATION_AUTO_MERGE=false               # Enable automatic merging
RECONCILIATION_CACHE_SIZE=10000               # Fuzzy match cache entries
RECONCILIATION_CACHE_TTL=24h                  # Cache expiration
```


### Config Keys
```yaml
reconciliation:
  default_fuzzy_threshold: 0.85        # Default similarity threshold (0.0-1.0)
  auto_merge: false                    # Enable automatic conflict resolution
  cache:
    size: 10000                        # Fuzzy match cache entries
    ttl: 24h                           # Cache expiration

  # Default source priority (highest to lowest)
  source_priority:
    scene:
      title: ["manual", "stashdb", "whisparr", "stash"]
      performers: ["manual", "stashdb", "stash", "whisparr"]
      studio: ["manual", "stashdb", "whisparr", "stash"]
      release_date: ["whisparr", "stashdb", "stash", "manual"]
      duration: ["stash", "whisparr", "manual"]
    performer:
      name: ["manual", "stashdb", "whisparr", "stash"]
      birthdate: ["stashdb", "manual", "stash"]
      aliases: ["stashdb", "stash", "manual"]
    studio:
      name: ["manual", "stashdb", "whisparr", "stash"]
      url: ["stashdb", "manual", "stash"]

  # Fuzzy matching algorithms
  algorithms:
    - levenshtein      # Edit distance
    - jaro_winkler     # Jaro-Winkler similarity
    - trigram          # PostgreSQL pg_trgm
```



## API Endpoints

### Content Management
**Endpoints**:
```
GET    /api/v1/legacy/reconciliation/conflicts                 # List pending conflicts
GET    /api/v1/legacy/reconciliation/conflicts/:id             # Get conflict details
POST   /api/v1/legacy/reconciliation/conflicts/:id/resolve     # Resolve conflict
POST   /api/v1/legacy/reconciliation/conflicts/:id/ignore      # Ignore conflict

GET    /api/v1/legacy/reconciliation/sources/:type/:id         # Get metadata sources for content
GET    /api/v1/legacy/reconciliation/sources/:type/:id/:field  # Get field history

GET    /api/v1/legacy/reconciliation/overrides/:type/:id       # List overrides for content
POST   /api/v1/legacy/reconciliation/overrides                 # Create manual override
DELETE /api/v1/legacy/reconciliation/overrides/:id             # Remove override

GET    /api/v1/legacy/reconciliation/rules                     # List reconciliation rules
POST   /api/v1/legacy/reconciliation/rules                     # Create rule
PUT    /api/v1/legacy/reconciliation/rules/:id                 # Update rule
```

**Request/Response Examples**:

**List Pending Conflicts**:
```http
GET /api/v1/legacy/reconciliation/conflicts?status=pending&limit=20

Response 200 OK:
{
  "conflicts": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "content_type": "scene",
      "content_id": "660e8400-e29b-41d4-a716-446655440001",
      "field_name": "title",
      "conflict_values": [
        {
          "source": "stashdb",
          "value": "Amazing Scene 1",
          "confidence": 0.95
        },
        {
          "source": "whisparr",
          "value": "Amazing Scene One",
          "confidence": 0.88
        }
      ],
      "resolution_status": "pending",
      "created_at": "2026-01-31T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

**Resolve Conflict**:
```http
POST /api/v1/legacy/reconciliation/conflicts/550e8400-e29b-41d4-a716-446655440000/resolve
{
  "resolved_value": "Amazing Scene 1",
  "source": "stashdb"
}

Response 200 OK:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "resolution_status": "resolved",
  "resolved_value": "Amazing Scene 1",
  "resolved_by": "770e8400-e29b-41d4-a716-446655440002",
  "resolved_at": "2026-01-31T10:30:00Z"
}
```

**Create Manual Override**:
```http
POST /api/v1/legacy/reconciliation/overrides
{
  "content_type": "performer",
  "content_id": "880e8400-e29b-41d4-a716-446655440003",
  "field_name": "birthdate",
  "override_value": "1990-05-15",
  "reason": "Corrected from official source"
}

Response 201 Created:
{
  "id": "990e8400-e29b-41d4-a716-446655440004",
  "content_type": "performer",
  "content_id": "880e8400-e29b-41d4-a716-446655440003",
  "field_name": "birthdate",
  "override_value": "1990-05-15",
  "reason": "Corrected from official source",
  "created_by": "770e8400-e29b-41d4-a716-446655440002",
  "created_at": "2026-01-31T11:00:00Z"
}
```

**Get Field History**:
```http
GET /api/v1/legacy/reconciliation/sources/scene/660e8400-e29b-41d4-a716-446655440001/title

Response 200 OK:
{
  "field_name": "title",
  "sources": [
    {
      "id": "aa0e8400-e29b-41d4-a716-446655440005",
      "source_name": "stashdb",
      "source_id": "stashdb_12345",
      "field_value": "Amazing Scene 1",
      "confidence_score": 0.95,
      "is_active": true,
      "created_at": "2026-01-31T09:00:00Z"
    },
    {
      "id": "bb0e8400-e29b-41d4-a716-446655440006",
      "source_name": "whisparr",
      "source_id": "whisparr_67890",
      "field_value": "Amazing Scene One",
      "confidence_score": 0.88,
      "is_active": false,
      "created_at": "2026-01-31T09:30:00Z"
    }
  ]
}
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Casbin](../../../sources/security/casbin.md) - Auto-resolved from casbin
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config

