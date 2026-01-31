# Internal Wiki System

> Modern, integrated knowledge base for users, mods, admins, and devs

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with DB schema, Goldmark rendering, search |
| Sources | ✅ | goldmark, bleve, meilisearch-go documented |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Last Updated**: 2026-01-30
**Location**: `internal/wiki/`

---

## Overview

A fully-integrated wiki/helpdesk system with:
- Modern UI/UX following design principles
- Breadcrumb navigation
- Auto-enrichment from external wikis
- Role-based content visibility
- Full-text search
- Adult content isolation in `c` schema

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                       Wiki System                                │
├──────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Users     │  │    Mods     │  │   Admins    │   Dev Docs   │
│  │   Guides    │  │   Guides    │  │   Guides    │   API Docs   │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├──────────────────────────────────────────────────────────────────┤
│  Full-Text Search (Bleve/Meilisearch) │ Auto-Enrichment         │
├──────────────────────────────────────────────────────────────────┤
│  Markdown Rendering (Goldmark) │ Version History │ RBAC         │
└──────────────────────────────────────────────────────────────────┘
```

---

## Go Packages

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core) for package versions.

Key packages used:
- **goldmark** - Markdown parsing (CommonMark)
- **goldmark-wikilink** - [[wiki-style]] links
- **goldmark-highlighting** - Syntax highlighting
- **goldmark-meta** - YAML frontmatter
- **bleve** - Full-text search (embedded)
- **meilisearch-go** - Full-text search (external)

---

## Content Hierarchy

### Spaces (Top-Level)

| Space | Visibility | Purpose |
|-------|------------|---------|
| `help` | All users | User guides, FAQ, tutorials |
| `mod` | Moderators | Moderation guidelines |
| `admin` | Admins | Admin documentation |
| `dev` | Developers | API docs, integrations |
| `internal` | Staff | Internal processes |

### Adult Content Space (Isolated)

| Space | Visibility | Purpose |
|-------|------------|---------|
| `c.help` | Adult users | Adult module guides |
| `c.mod` | Adult mods | Adult moderation |

---

## Database Schema

```sql
-- Wiki spaces (namespaces)
CREATE TABLE wiki.spaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    visibility VARCHAR(50) DEFAULT 'public', -- public, role_based
    required_roles TEXT[],
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Wiki pages
CREATE TABLE wiki.pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id UUID REFERENCES wiki.spaces(id),
    parent_id UUID REFERENCES wiki.pages(id),
    slug VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    content_markdown TEXT NOT NULL,
    content_html TEXT, -- Pre-rendered
    summary TEXT, -- Auto-generated or manual
    tags TEXT[],

    -- Breadcrumb path (materialized)
    path TEXT NOT NULL, -- /help/getting-started/installation
    depth INT DEFAULT 0,

    -- Metadata
    author_id UUID REFERENCES users(id),
    last_editor_id UUID REFERENCES users(id),
    published BOOLEAN DEFAULT false,

    -- Search optimization
    search_vector TSVECTOR,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(space_id, slug)
);

-- Version history
CREATE TABLE wiki.page_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID REFERENCES wiki.pages(id),
    version INT NOT NULL,
    content_markdown TEXT NOT NULL,
    editor_id UUID REFERENCES users(id),
    change_summary TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(page_id, version)
);

-- Cross-references (auto-detected [[links]])
CREATE TABLE wiki.page_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_page_id UUID REFERENCES wiki.pages(id),
    target_page_id UUID REFERENCES wiki.pages(id),
    link_text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(source_page_id, target_page_id)
);

-- External enrichment cache
CREATE TABLE wiki.enrichment_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID REFERENCES wiki.pages(id),
    source VARCHAR(100), -- wikipedia, wikidata, etc.
    external_url TEXT,
    cached_content JSONB,
    fetched_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

-- Adult wiki (isolated in c schema)
CREATE TABLE c.wiki_spaces ( /* same structure */ );
CREATE TABLE c.wiki_pages ( /* same structure */ );
CREATE TABLE c.wiki_page_versions ( /* same structure */ );
```

---

## Markdown Processing

### Goldmark Configuration

```go
import (
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    highlighting "github.com/yuin/goldmark-highlighting/v2"
    meta "github.com/yuin/goldmark-meta"
    wikilink "go.abhg.dev/goldmark/wikilink"
)

func NewWikiMarkdown() goldmark.Markdown {
    return goldmark.New(
        goldmark.WithExtensions(
            extension.GFM,          // GitHub Flavored Markdown
            extension.Footnote,
            extension.Typographer,
            meta.Meta,              // YAML frontmatter
            highlighting.NewHighlighting(
                highlighting.WithStyle("monokai"),
            ),
            &wikilink.Extender{},   // [[wiki-style]] links
        ),
        goldmark.WithParserOptions(
            parser.WithAutoHeadingID(),
        ),
        goldmark.WithRendererOptions(
            html.WithHardWraps(),
            html.WithXHTML(),
        ),
    )
}
```

### Frontmatter Support

```markdown
---
title: Getting Started
summary: Quick start guide for new users
tags: [beginner, setup, installation]
---

# Getting Started

Welcome to Revenge! This guide will help you...
```

### Wiki Links

```markdown
See [[Installation Guide]] for setup instructions.

Link to specific section: [[Installation Guide#docker]]

Cross-space link: [[admin:User Management]]
```

---

## Auto-Enrichment

### External Wiki Sources

| Source | Type | Use Case |
|--------|------|----------|
| Wikipedia | REST API | General knowledge enrichment |
| Wikidata | SPARQL | Structured data, identifiers |
| Fandom | REST API | Franchise-specific content |

### River Jobs for Enrichment

```go
const (
    JobKindEnrichPage     = "wiki.enrich_page"
    JobKindRefreshEnrich  = "wiki.refresh_enrichment"
)

type EnrichPageArgs struct {
    PageID    uuid.UUID `json:"page_id"`
    Sources   []string  `json:"sources"` // ["wikipedia", "wikidata"]
    ForceRefresh bool   `json:"force_refresh"`
}

// Auto-detects enrichable terms in page content
// Fetches relevant data from external wikis
// Caches in wiki.enrichment_cache
```

### Enrichment Display

```
┌─────────────────────────────────────────────┐
│ Page: The Matrix (1999)                     │
├─────────────────────────────────────────────┤
│ [User Content]                              │
│                                             │
│ ─────────────────────────────────────────── │
│ 📚 From Wikipedia:                          │
│ "The Matrix is a 1999 science fiction..."   │
│ [Read more on Wikipedia →]                  │
└─────────────────────────────────────────────┘
```

---

## Full-Text Search

### Option 1: Bleve (Embedded)

```go
import "github.com/blevesearch/bleve/v2"

type WikiSearchService struct {
    index bleve.Index
}

func (s *WikiSearchService) IndexPage(page *Page) error {
    doc := map[string]interface{}{
        "title":   page.Title,
        "content": page.ContentMarkdown,
        "tags":    page.Tags,
        "path":    page.Path,
        "space":   page.SpaceSlug,
    }
    return s.index.Index(page.ID.String(), doc)
}

func (s *WikiSearchService) Search(ctx context.Context, query string, spaceFilter string) ([]SearchResult, error) {
    q := bleve.NewMatchQuery(query)
    req := bleve.NewSearchRequest(q)
    req.Highlight = bleve.NewHighlight()

    if spaceFilter != "" {
        spaceQuery := bleve.NewTermQuery(spaceFilter)
        spaceQuery.SetField("space")
        q = bleve.NewConjunctionQuery(q, spaceQuery)
    }

    results, err := s.index.Search(req)
    // Convert to SearchResult...
}
```

### Option 2: Meilisearch (External)

```go
import "github.com/meilisearch/meilisearch-go"

type WikiSearchService struct {
    client meilisearch.ServiceManager
}

func (s *WikiSearchService) Search(query string, filters string) (*meilisearch.SearchResponse, error) {
    return s.client.Index("wiki_pages").Search(query, &meilisearch.SearchRequest{
        Filter: filters, // "space = 'help' AND published = true"
        AttributesToHighlight: []string{"title", "content"},
        HighlightPreTag:  "<mark>",
        HighlightPostTag: "</mark>",
    })
}
```

---

## UI/UX Integration

### Breadcrumb Navigation

```
Home > Help > Getting Started > Installation
                                     ↓
                              [Current Page]
```

### Sidebar Navigation

```
📚 Help
  ├── 🏠 Getting Started
  │   ├── Installation
  │   ├── Configuration
  │   └── First Steps
  ├── 📺 Media Management
  │   ├── Adding Libraries
  │   └── Metadata
  └── ❓ FAQ

🛡️ Moderation (mods only)
  └── ...
```

### Features

- **Table of Contents**: Auto-generated from headings
- **Edit Button**: Inline editing for authorized users
- **Version History**: Diff view between versions
- **Related Pages**: Auto-linked from [[wiki links]]
- **Search**: Global search with space filtering
- **Tags**: Tag-based navigation and filtering
- **Print View**: Clean printable version
- **Dark Mode**: Follows system theme

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `wiki.pages.view` | View published pages |
| `wiki.pages.view_draft` | View draft pages |
| `wiki.pages.create` | Create new pages |
| `wiki.pages.edit` | Edit existing pages |
| `wiki.pages.delete` | Delete pages |
| `wiki.pages.publish` | Publish/unpublish pages |
| `wiki.spaces.manage` | Manage wiki spaces |
| `wiki.versions.restore` | Restore previous versions |

---

## API Endpoints

```
# Spaces
GET  /api/v1/wiki/spaces              # List spaces (filtered by role)
GET  /api/v1/wiki/spaces/:slug        # Get space

# Pages
GET  /api/v1/wiki/pages               # List pages
GET  /api/v1/wiki/pages/:path         # Get page by path
POST /api/v1/wiki/pages               # Create page
PUT  /api/v1/wiki/pages/:id           # Update page
DELETE /api/v1/wiki/pages/:id         # Delete page

# Versions
GET  /api/v1/wiki/pages/:id/versions  # List versions
GET  /api/v1/wiki/pages/:id/versions/:v # Get version
POST /api/v1/wiki/pages/:id/restore/:v # Restore version

# Search
GET  /api/v1/wiki/search?q=...        # Full-text search

# Adult (isolated)
GET  /api/v1/legacy/wiki/spaces
GET  /api/v1/legacy/wiki/pages/:path
GET  /api/v1/legacy/wiki/search?q=...
```

---

## Configuration

```yaml
wiki:
  enabled: true

  # Search engine
  search:
    engine: bleve  # bleve (embedded) or meilisearch
    meilisearch:
      url: "http://localhost:7700"
      api_key: "${MEILISEARCH_API_KEY}"

  # Auto-enrichment
  enrichment:
    enabled: true
    sources:
      - wikipedia
      - wikidata
    cache_ttl: 7d

  # Version history
  versions:
    max_per_page: 100
    retention_days: 365

  # Adult wiki (isolated)
  adult:
    enabled: false
    schema: c
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure at `internal/wiki/`
- [ ] Create sub-packages: `spaces/`, `pages/`, `search/`, `enrichment/`
- [ ] Define Space entity (`spaces/entity.go`)
- [ ] Define Page entity (`pages/entity.go`)
- [ ] Define PageVersion entity
- [ ] Define PageLink entity (cross-references)
- [ ] Define EnrichmentCache entity
- [ ] Create repository interfaces
- [ ] Implement PostgreSQL repositories
- [ ] Create fx module (`module.go`)
- [ ] Add configuration structs

### Phase 2: Database
- [ ] Create `wiki` schema
- [ ] Create migration for `wiki.spaces` table
- [ ] Create migration for `wiki.pages` table
- [ ] Create migration for `wiki.page_versions` table
- [ ] Create migration for `wiki.page_links` table
- [ ] Create migration for `wiki.enrichment_cache` table
- [ ] Create adult wiki tables in `c` schema (isolated)
- [ ] Add GIN index on search_vector for full-text search
- [ ] Add indexes for path lookups, space filtering
- [ ] Write sqlc queries for space CRUD
- [ ] Write sqlc queries for page CRUD
- [ ] Write sqlc queries for version history
- [ ] Write sqlc queries for link tracking

### Phase 3: Service Layer
- [ ] Implement SpaceService
  - [ ] Space CRUD with RBAC
  - [ ] Visibility filtering by role
- [ ] Implement PageService
  - [ ] Page CRUD with validation
  - [ ] Path/breadcrumb computation
  - [ ] Version creation on save
  - [ ] Version restoration
- [ ] Implement MarkdownService
  - [ ] Configure Goldmark with extensions (GFM, footnotes, typographer)
  - [ ] Add goldmark-meta for YAML frontmatter
  - [ ] Add goldmark-highlighting for syntax highlighting
  - [ ] Add goldmark-wikilink for [[wiki-style]] links
  - [ ] Implement wiki link resolution
  - [ ] Pre-render HTML on save
- [ ] Implement SearchService
  - [ ] Bleve index management (embedded option)
  - [ ] Meilisearch client (external option)
  - [ ] Index page on create/update
  - [ ] Remove from index on delete
  - [ ] Search with space filtering
- [ ] Add caching for rendered pages

### Phase 4: Background Jobs
- [ ] Create River job for auto-enrichment (Wikipedia, Wikidata)
- [ ] Create River job for enrichment cache refresh
- [ ] Create River job for search index rebuild
- [ ] Create River job for orphaned page detection
- [ ] Create River job for version history pruning
- [ ] Register jobs in fx module

### Phase 5: API Integration
- [ ] Add OpenAPI spec for wiki endpoints
- [ ] Generate ogen handlers
- [ ] Implement space endpoints:
  - [ ] GET /api/v1/wiki/spaces (list, filtered by role)
  - [ ] GET /api/v1/wiki/spaces/{slug}
- [ ] Implement page endpoints:
  - [ ] GET /api/v1/wiki/pages (list)
  - [ ] GET /api/v1/wiki/pages/{path} (by path)
  - [ ] POST /api/v1/wiki/pages (create)
  - [ ] PUT /api/v1/wiki/pages/{id} (update)
  - [ ] DELETE /api/v1/wiki/pages/{id}
- [ ] Implement version endpoints:
  - [ ] GET /api/v1/wiki/pages/{id}/versions
  - [ ] GET /api/v1/wiki/pages/{id}/versions/{v}
  - [ ] POST /api/v1/wiki/pages/{id}/restore/{v}
- [ ] Implement search endpoint:
  - [ ] GET /api/v1/wiki/search?q=...
- [ ] Implement adult wiki endpoints (isolated):
  - [ ] GET /api/v1/legacy/wiki/spaces
  - [ ] GET /api/v1/legacy/wiki/pages/{path}
  - [ ] GET /api/v1/legacy/wiki/search?q=...
- [ ] Add authentication middleware
- [ ] Add RBAC permission checks (wiki.pages.view, create, edit, delete, publish)

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Go io](https://pkg.go.dev/io) | [Local](../../../sources/go/stdlib/io.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related

- [News System](NEWS_SYSTEM.md) - Announcement system
- [RBAC Permissions](RBAC_CASBIN.md) - Permission management
- [Wiki Providers](../integrations/wiki/INDEX.md) - External wiki sources
- [UI/UX Guidelines](../architecture/UI_UX_GUIDELINES.md) - Design patterns
