# Instagram Integration

<!-- SOURCES: sqlc, sqlc-config -->

<!-- DESIGN: integrations/metadata/adult, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Performer social media presence


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Access Challenges](#api-access-challenges)
- [Configuration](#configuration)
- [Data Model](#data-model)
  - [Performer Social Link](#performer-social-link)
- [Implementation](#implementation)
  - [Manual Link Management](#manual-link-management)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
  - [Add Instagram Link](#add-instagram-link)
  - [Get Social Links](#get-social-links)
- [Implementation Checklist](#implementation-checklist)
- [Future Considerations](#future-considerations)
- [Privacy Considerations](#privacy-considerations)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | API access challenges documented, manual link approach defined |
| Sources | ✅ | Basic Display API, Graph API docs linked |
| Instructions | ✅ | Implementation checklist with sqlc queries |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |---

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


## Related Documentation

- [Twitter/X](TWITTER_X.md)
- [FreeOnes](../metadata/adult/FREEONES.md)
- [StashDB](../metadata/adult/STASHDB.md)
