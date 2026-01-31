# Instagram Integration

> Performer social media presence

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | API access challenges documented, manual link approach defined |
| Sources | ✅ | Basic Display API, Graph API docs linked |
| Instructions | ✅ | Implementation checklist with sqlc queries |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

Instagram integration for adult content performer profiles:
- Official performer accounts
- Social media presence tracking
- Profile verification status
- Follower metrics

**Note**: This is supplementary metadata for performer profiles. Instagram's official API has strict limitations.

---

## Developer Resources

- 📚 **API Docs**: https://developers.facebook.com/docs/instagram-basic-display-api/
- 🔗 **Graph API**: https://developers.facebook.com/docs/instagram-api/
- ⚠️ **Limitations**: Basic Display API deprecated, Graph API requires business accounts

---

## API Access Challenges

Instagram has significantly restricted API access:

| Approach | Viability | Notes |
|----------|-----------|-------|
| Basic Display API | ❌ Deprecated | Shut down Dec 2024 |
| Instagram Graph API | ⚠️ Limited | Business/Creator accounts only |
| Scraping | ❌ Blocked | ToS violation, rate limited |
| Manual Entry | ✅ Viable | User provides link |

**Recommendation**: Store Instagram links as user-provided metadata without API verification.

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  adult:
    social:
      instagram:
        enabled: false  # Manual links only
        # No API integration due to restrictions
        validation:
          verify_url_format: true
```

---

## Data Model

### Performer Social Link

| Field | Type | Notes |
|-------|------|-------|
| `performer_id` | UUID | FK to performer |
| `platform` | string | "instagram" |
| `handle` | string | @username |
| `profile_url` | string | Full URL |
| `verified` | bool | Manual flag |
| `followers_count` | int | Manual entry |
| `last_updated` | timestamp | When manually updated |

---

## Implementation

### Manual Link Management

```go
// internal/content/c/service/social.go
package service

type SocialLinkService struct {
    queries *db.Queries
    logger  *slog.Logger
}

type InstagramLink struct {
    Handle     string `json:"handle"`
    ProfileURL string `json:"profile_url"`
    Verified   bool   `json:"verified"`
}

// AddInstagramLink adds an Instagram profile to a performer
func (s *SocialLinkService) AddInstagramLink(ctx context.Context, performerID uuid.UUID, link InstagramLink) error {
    // Validate URL format
    if !isValidInstagramURL(link.ProfileURL) {
        return ErrInvalidInstagramURL
    }

    // Extract handle from URL if not provided
    if link.Handle == "" {
        link.Handle = extractInstagramHandle(link.ProfileURL)
    }

    return s.queries.UpsertSocialLink(ctx, db.UpsertSocialLinkParams{
        PerformerID: performerID,
        Platform:    "instagram",
        Handle:      link.Handle,
        ProfileURL:  link.ProfileURL,
        Verified:    link.Verified,
    })
}

func isValidInstagramURL(url string) bool {
    pattern := `^https?://(www\.)?instagram\.com/[a-zA-Z0-9_.]+/?$`
    matched, _ := regexp.MatchString(pattern, url)
    return matched
}

func extractInstagramHandle(url string) string {
    // Extract username from instagram.com/username/
    pattern := `instagram\.com/([a-zA-Z0-9_.]+)`
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(url)
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}
```

---

## Database Schema

```sql
-- Uses same table as Twitter (in schema 'c')
-- See TWITTER_X.md for full schema

-- Instagram-specific data stored in performer_social_links
-- platform = 'instagram'

-- Example insert
INSERT INTO c.performer_social_links (
    performer_id,
    platform,
    handle,
    profile_url,
    verified
) VALUES (
    'performer-uuid',
    'instagram',
    'performer_handle',
    'https://instagram.com/performer_handle',
    false
);
```

---

## API Endpoints

### Add Instagram Link

```
POST /api/v1/legacy/performers/{id}/social-links
Content-Type: application/json

{
    "platform": "instagram",
    "handle": "performer_handle",
    "profile_url": "https://instagram.com/performer_handle"
}
```

### Get Social Links

```
GET /api/v1/legacy/performers/{id}/social-links

{
    "social_links": [
        {
            "platform": "instagram",
            "handle": "performer_handle",
            "profile_url": "https://instagram.com/performer_handle",
            "verified": false
        },
        {
            "platform": "twitter",
            "handle": "performer_handle",
            "profile_url": "https://x.com/performer_handle",
            "verified": true
        }
    ]
}
```

---

## Implementation Checklist

- [ ] **Social Link Service** (`internal/content/c/service/social.go`)
  - [ ] Add link
  - [ ] Remove link
  - [ ] URL validation
  - [ ] Handle extraction

- [ ] **API Handler** (`internal/content/c/handler/social.go`)
  - [ ] POST social link
  - [ ] DELETE social link
  - [ ] GET performer social links

- [ ] **sqlc Queries** (`internal/infra/database/queries/c/social.sql`)
  - [ ] UpsertSocialLink
  - [ ] DeleteSocialLink
  - [ ] GetPerformerSocialLinks

---

## Future Considerations

If Instagram API access becomes available:

```go
// Future: Instagram Graph API client
type InstagramClient struct {
    accessToken string
    httpClient  *http.Client
}

// Would require:
// 1. Facebook App approval
// 2. Business account verification
// 3. User OAuth consent
func (c *InstagramClient) GetBusinessProfile(ctx context.Context, username string) (*InstagramProfile, error) {
    // Graph API endpoint for business discovery
    // GET /ig_username?fields=biography,id,username,website,profile_picture_url
    return nil, ErrNotImplemented
}
```

---

## Privacy Considerations

- Only public profile URLs stored
- No scraping of private data
- User-provided links only
- Follows adult content isolation (`c` schema)
- Performers can request removal

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [FreeOnes Integration](FREEONES.md)
- [OnlyFans Integration](ONLYFANS.md)
- [Pornhub Integration](PORNHUB.md)
- [TheNude Integration](THENUDE.md)
- [Twitter/X Integration](TWITTER_X.md)

### Related Topics

- [Revenge - Architecture v2](../../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Twitter/X](TWITTER_X.md)
- [FreeOnes](../metadata/adult/FREEONES.md)
- [StashDB](../metadata/adult/STASHDB.md)
