# Adult Data Reconciliation

<!-- SOURCES: casbin, river, sqlc, sqlc-config -->

<!-- DESIGN: features/adult, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Fuzzy logic and confidence scoring for conflicting metadata


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Problem Examples](#problem-examples)
  - [Measurements](#measurements)
  - [Career Status](#career-status)
- [Solution Architecture](#solution-architecture)
- [Source Trust Scoring](#source-trust-scoring)
  - [Trust Modifiers](#trust-modifiers)
- [Fuzzy Logic Rules](#fuzzy-logic-rules)
  - [Measurements (Bust, Waist, Hips)](#measurements-bust-waist-hips)
  - [Cup Size Normalization](#cup-size-normalization)
  - [Career Status Rules](#career-status-rules)
- [Confidence Scoring](#confidence-scoring)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [UI Integration](#ui-integration)
  - [Performer Profile](#performer-profile)
  - [Source Comparison Modal](#source-comparison-modal)
- [Manual Override](#manual-override)
- [Configuration](#configuration)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Schema & Database](#phase-1-schema-database)
  - [Phase 2: Core Reconciliation Engine](#phase-2-core-reconciliation-engine)
  - [Phase 3: Fuzzy Logic Rules](#phase-3-fuzzy-logic-rules)
  - [Phase 4: Source Integration](#phase-4-source-integration)
  - [Phase 5: Reconciliation Pipeline](#phase-5-reconciliation-pipeline)
  - [Phase 6: Manual Override Support](#phase-6-manual-override-support)
  - [Phase 7: UI Components](#phase-7-ui-components)
  - [Phase 8: Reconciliation Scheduling](#phase-8-reconciliation-scheduling)
  - [Phase 9: Configuration](#phase-9-configuration)
  - [Phase 10: Testing](#phase-10-testing)
  - [Phase 11: Documentation](#phase-11-documentation)
- [Go Packages](#go-packages)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

**⚠️ Adult Content**: All processing isolated in `qar` schema.

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | 🟡 |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |---

## Overview

Adult performer metadata from multiple sources often contains:
- **Conflicting data** (different measurements, birth dates)
- **Outdated data** (old aliases, retired status)
- **Incomplete data** (missing fields)
- **Incorrect data** (typos, wrong attributions)

This system uses **fuzzy logic** and **confidence scoring** to reconcile data intelligently.

---

## Problem Examples

### Measurements

| Source | Bust | Waist | Hips |
|--------|------|-------|------|
| FreeOnes | 34D | 24 | 34 |
| Babepedia | 34DD | 25 | 35 |
| IAFD | 34D | 24 | 34 |
| StashDB | 34D | 24 | 35 |

**Question**: Which is correct? How do we decide?

### Career Status

| Source | Status | Last Updated |
|--------|--------|--------------|
| FreeOnes | Active | 2024-01 |
| Babepedia | Retired | 2022-06 |
| IAFD | Active | 2024-03 |

**Question**: Is performer active or retired?

---

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                   Data Reconciliation Engine                    │
├─────────────────────────────────────────────────────────────────┤
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐   │
│  │ FreeOnes  │  │ Babepedia │  │   IAFD    │  │  StashDB  │   │
│  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘   │
│        │              │              │              │          │
│        └──────────────┴──────────────┴──────────────┘          │
│                              │                                  │
│                    ┌─────────▼─────────┐                       │
│                    │  Field Analyzer   │                       │
│                    └─────────┬─────────┘                       │
│                              │                                  │
│  ┌───────────────────────────┼───────────────────────────┐     │
│  │                           │                           │     │
│  ▼                           ▼                           ▼     │
│ ┌────────────┐     ┌────────────────┐     ┌────────────────┐  │
│ │  Majority  │     │   Recency      │     │   Source       │  │
│ │  Voting    │     │   Weighting    │     │   Trust Score  │  │
│ └────────────┘     └────────────────┘     └────────────────┘  │
│        │                    │                    │             │
│        └────────────────────┴────────────────────┘             │
│                             │                                   │
│                   ┌─────────▼─────────┐                        │
│                   │  Fuzzy Decision   │                        │
│                   │     Engine        │                        │
│                   └─────────┬─────────┘                        │
│                             │                                   │
│                   ┌─────────▼─────────┐                        │
│                   │  Final Value +    │                        │
│                   │  Confidence Score │                        │
│                   └───────────────────┘                        │
└─────────────────────────────────────────────────────────────────┘
```

---

## Source Trust Scoring

Each source gets a base trust score (configurable):

| Source | Base Trust | Recency Weight | Notes |
|--------|------------|----------------|-------|
| StashDB | 0.85 | High | Community curated, frequently updated |
| FreeOnes | 0.80 | High | Professional database, good accuracy |
| IAFD | 0.75 | Medium | Long history, sometimes outdated |
| Babepedia | 0.70 | Medium | Wiki-style, variable quality |
| Boobpedia | 0.65 | Low | Community edited, less reliable |
| TheNude | 0.70 | Medium | Good for EU performers |
| Indexxx | 0.75 | Medium | Cross-reference useful |

### Trust Modifiers

```go
type TrustModifier struct {
    RecencyBonus    float64 // +0.1 if updated < 6 months
    RecencyPenalty  float64 // -0.1 if updated > 2 years
    VerifiedBonus   float64 // +0.1 if verified/official
    ConsistencyMod  float64 // +/- based on historical accuracy
}
```

---

## Fuzzy Logic Rules

### Measurements (Bust, Waist, Hips)

```go
type MeasurementReconciliation struct {
    // Tolerance ranges (in inches)
    BustTolerance   float64 = 1.0   // 34D vs 34DD = close enough
    WaistTolerance  float64 = 1.0
    HipsTolerance   float64 = 1.0
    WeightTolerance float64 = 5.0   // lbs
    HeightTolerance float64 = 1.0   // inches
}

func ReconcileMeasurement(values []SourceValue) (string, float64) {
    // 1. Group similar values (within tolerance)
    clusters := clusterSimilarValues(values, tolerance)

    // 2. Calculate weighted score per cluster
    for _, cluster := range clusters {
        cluster.Score = 0
        for _, v := range cluster.Values {
            weight := v.Source.TrustScore * recencyWeight(v.UpdatedAt)
            cluster.Score += weight
        }
    }

    // 3. Select highest scoring cluster
    winner := selectHighestScoring(clusters)

    // 4. Calculate confidence
    confidence := calculateConfidence(winner, clusters)

    return winner.RepresentativeValue, confidence
}
```

### Cup Size Normalization

Different regions use different cup size systems:

```go
var CupSizeMap = map[string]int{
    // US/UK sizes
    "AA": 1, "A": 2, "B": 3, "C": 4, "D": 5, "DD": 6, "DDD": 7, "E": 7,
    "F": 8, "G": 9, "H": 10, "I": 11, "J": 12,
    // EU equivalents
    "70A": 2, "70B": 3, "75C": 4, "80D": 5,
}

func NormalizeCupSize(input string) (string, int) {
    // Convert to standardized US size
    // Return (normalized_string, numeric_value)
}
```

### Career Status Rules

```go
func ReconcileCareerStatus(statuses []SourceStatus) (string, float64) {
    rules := []FuzzyRule{
        // Most recent data wins (if recent enough)
        {
            Condition: func(s []SourceStatus) bool {
                recent := getNewestStatus(s)
                return recent.UpdatedAt.After(time.Now().AddDate(0, -6, 0))
            },
            Action: func(s []SourceStatus) string {
                return getNewestStatus(s).Status
            },
            Weight: 0.8,
        },
        // Majority voting for older data
        {
            Condition: func(s []SourceStatus) bool {
                return len(s) >= 3
            },
            Action: func(s []SourceStatus) string {
                return majorityVote(s)
            },
            Weight: 0.6,
        },
        // Trust highest-ranked source
        {
            Condition: func(s []SourceStatus) bool {
                return true // fallback
            },
            Action: func(s []SourceStatus) string {
                return getHighestTrustSource(s).Status
            },
            Weight: 0.4,
        },
    }

    return evaluateFuzzyRules(statuses, rules)
}
```

---

## Confidence Scoring

Each reconciled field gets a confidence score (0.0 - 1.0):

| Confidence | Meaning | UI Indicator |
|------------|---------|--------------|
| 0.9 - 1.0 | High confidence | ✅ |
| 0.7 - 0.9 | Good confidence | 🟢 |
| 0.5 - 0.7 | Moderate confidence | 🟡 |
| 0.3 - 0.5 | Low confidence | 🟠 |
| 0.0 - 0.3 | Very low / conflicting | 🔴 |

```go
func CalculateConfidence(winner Cluster, allClusters []Cluster) float64 {
    // Factors:
    // 1. Source agreement (how many sources agree)
    agreement := float64(len(winner.Values)) / float64(totalValues)

    // 2. Trust score sum of agreeing sources
    trustSum := sumTrustScores(winner.Values)
    maxTrust := sumTrustScores(allValues)
    trustRatio := trustSum / maxTrust

    // 3. Spread of other clusters (if close competitors, lower confidence)
    spread := 1.0 - (secondHighest.Score / winner.Score)

    // Weighted combination
    confidence := (agreement * 0.3) + (trustRatio * 0.4) + (spread * 0.3)

    return math.Min(confidence, 1.0)
}
```

---

## Database Schema

```sql
-- Reconciled crew data (stored in qar schema)
CREATE TABLE qar.crew_reconciled (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crew_id UUID REFERENCES qar.crew(id),

    -- Reconciled values with confidence
    height_cm INT,
    height_confidence DECIMAL(3,2),

    weight_kg INT,
    weight_confidence DECIMAL(3,2),

    bust_size VARCHAR(10),
    bust_confidence DECIMAL(3,2),

    waist_size INT,
    waist_confidence DECIMAL(3,2),

    hip_size INT,
    hip_confidence DECIMAL(3,2),

    christening DATE,                        -- birth_date (obfuscated)
    christening_confidence DECIMAL(3,2),

    career_status VARCHAR(50),
    career_status_confidence DECIMAL(3,2),

    -- Source data for auditing
    source_data JSONB, -- Raw data from all sources

    reconciled_at TIMESTAMPTZ DEFAULT NOW(),
    last_source_update TIMESTAMPTZ
);

-- Raw source data (for re-reconciliation)
CREATE TABLE qar.crew_source_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crew_id UUID REFERENCES qar.crew(id),
    source VARCHAR(100) NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    field_value TEXT,
    source_url TEXT,
    fetched_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(crew_id, source, field_name)
);
```

---

## River Jobs

```go
const (
    JobKindReconcilePerformer = "adult.reconcile_performer"
    JobKindReconcileAll       = "adult.reconcile_all_performers"
    JobKindFetchSourceData    = "adult.fetch_source_data"
)

// Triggered when:
// 1. New source data is fetched
// 2. Manually requested
// 3. Periodic re-reconciliation (monthly)
```

---

## UI Integration

### Performer Profile

```
┌─────────────────────────────────────────┐
│ Measurements                            │
├─────────────────────────────────────────┤
│ Height: 5'7" (170cm)          ✅ 0.95  │
│ Bust: 34D                     🟢 0.82  │
│ Waist: 24"                    🟢 0.78  │
│ Hips: 35"                     🟡 0.65  │ ← Click to see sources
│                                         │
│ [View all sources] [Report error]       │
└─────────────────────────────────────────┘
```

### Source Comparison Modal

```
┌─────────────────────────────────────────────────────────────┐
│ Hips Measurement - Sources                                  │
├─────────────────────────────────────────────────────────────┤
│ Source      │ Value │ Updated    │ Trust │ Selected        │
├─────────────┼───────┼────────────┼───────┼─────────────────┤
│ FreeOnes    │ 34"   │ 2024-01-15 │ 0.80  │                 │
│ StashDB     │ 35"   │ 2024-03-01 │ 0.85  │ ✓ (selected)    │
│ IAFD        │ 34"   │ 2023-06-10 │ 0.75  │                 │
│ Babepedia   │ 35"   │ 2022-08-20 │ 0.70  │                 │
├─────────────┴───────┴────────────┴───────┴─────────────────┤
│ Decision: 35" (StashDB + Babepedia agreement + recency)    │
│ Confidence: 0.65 (moderate - 50/50 split)                  │
│                                                             │
│ [Override to 34"] [Report source error] [Close]            │
└─────────────────────────────────────────────────────────────┘
```

---

## Manual Override

Admins/mods can override reconciled values:

```sql
CREATE TABLE qar.crew_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crew_id UUID REFERENCES qar.crew(id),
    field_name VARCHAR(100) NOT NULL,
    override_value TEXT NOT NULL,
    reason TEXT,
    override_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(crew_id, field_name)
);
```

RBAC Permission: `adult.metadata.override`

---

## Configuration

```yaml
reconciliation:
  adult:
    enabled: true

    # Source trust scores (override defaults)
    source_trust:
      stashdb: 0.85
      freeones: 0.80
      iafd: 0.75
      babepedia: 0.70

    # Recency weights
    recency:
      fresh_bonus: 0.1      # < 6 months
      stale_penalty: -0.1   # > 2 years
      ancient_penalty: -0.2 # > 5 years

    # Confidence thresholds
    confidence:
      min_for_display: 0.3  # Don't show if below this
      manual_review: 0.5    # Flag for review if below

    # Re-reconciliation schedule
    schedule:
      full_reconcile: "0 0 1 * *"  # Monthly
      on_new_data: true
```

---

## Implementation Checklist

### Phase 1: Schema & Database
- [ ] Create `qar.crew_reconciled` table for reconciled values
- [ ] Create `qar.crew_source_data` table for raw source data
- [ ] Create `qar.crew_overrides` table for admin overrides
- [ ] Create indexes on `crew_id`, `source`, `field_name`
- [ ] Implement sqlc queries for reconciliation data access

### Phase 2: Core Reconciliation Engine
- [ ] Implement `internal/content/qar/reconcile/entity.go` with types
- [ ] Implement `internal/content/qar/reconcile/repository.go` interface
- [ ] Implement `internal/content/qar/reconcile/repository_pg.go`
- [ ] Implement `internal/content/qar/reconcile/service.go` base logic

### Phase 3: Fuzzy Logic Rules
- [ ] Implement measurement reconciliation (bust, waist, hips, height, weight)
- [ ] Implement cup size normalization (US/UK/EU conversions)
- [ ] Implement career status rules (recency weighting, majority voting)
- [ ] Implement trust score calculation based on source reliability
- [ ] Implement recency bonus/penalty modifiers
- [ ] Implement confidence scoring calculation

### Phase 4: Source Integration
- [ ] Create `internal/integrations/freeones/client.go`
- [ ] Create `internal/integrations/babepedia/client.go` (if applicable)
- [ ] Create `internal/integrations/iafd/client.go` (if applicable)
- [ ] Create `internal/integrations/stashdb/performer_sync.go` extension
- [ ] Implement data fetching and normalization for each source
- [ ] Store raw source data in `qar.crew_source_data`

### Phase 5: Reconciliation Pipeline
- [ ] Implement `internal/content/qar/reconcile/reconciler.go`
- [ ] Implement measurement field reconciliation
- [ ] Implement career status reconciliation
- [ ] Implement demographic field reconciliation
- [ ] Implement confidence calculation for each field
- [ ] Create River job for reconciliation processing

### Phase 6: Manual Override Support
- [ ] Implement override repository methods
- [ ] Implement RBAC permission check (`adult.metadata.override`)
- [ ] Create River job for override processing
- [ ] Implement audit logging for overrides
- [ ] Implement API endpoints for admin overrides

### Phase 7: UI Components
- [ ] Create `internal/web/handlers/qar/performer_reconcile.go`
- [ ] Implement performer profile endpoint with confidence indicators
- [ ] Implement source comparison modal data endpoint
- [ ] Implement "view all sources" functionality
- [ ] Implement override UI (admin only)
- [ ] Implement "report error" submission

### Phase 8: Reconciliation Scheduling
- [ ] Create cron job for monthly full reconciliation
- [ ] Create River job triggered on new source data
- [ ] Implement manual reconciliation trigger endpoint
- [ ] Create re-reconciliation on source data update

### Phase 9: Configuration
- [ ] Add reconciliation section to config.yaml
- [ ] Implement source trust score configuration
- [ ] Implement recency weight configuration
- [ ] Implement confidence threshold configuration
- [ ] Implement schedule configuration

### Phase 10: Testing
- [ ] Write unit tests for fuzzy logic rules
- [ ] Write unit tests for confidence scoring
- [ ] Write integration tests for source clients
- [ ] Write tests for edge cases (conflicting data, missing fields)
- [ ] Write tests for override functionality
- [ ] Test different tolerance ranges

### Phase 11: Documentation
- [ ] Document source trust scores
- [ ] Document fuzzy logic rules
- [ ] Document confidence scoring algorithm
- [ ] Create performer metadata audit guide
- [ ] Document override process for admins

---

## Go Packages

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core) for package versions.

Key packages used:
- **fuzzy** - Fuzzy string matching (github.com/lithammer/fuzzysearch)
- **go-levenshtein** - Edit distance (github.com/agnivade/levenshtein)

---


## Related Documentation

- [Adult Metadata Providers](../../integrations/metadata/adult/INDEX.md)
- [Adult Wiki Providers](../../integrations/wiki/adult/INDEX.md)
- [River Job Queue Patterns](../../00_SOURCE_OF_TRUTH.md#river-job-queue-patterns)
- [RBAC Permissions](RBAC_CASBIN.md)
