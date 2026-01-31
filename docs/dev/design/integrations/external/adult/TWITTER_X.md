# Twitter/X Integration

> Performer social media presence

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | REST API v2 spec, data mapping, database schema |
| Sources | ✅ | API docs, pricing tiers linked |
| Instructions | ✅ | Implementation checklist with River jobs |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

Twitter/X integration for adult content performer profiles:
- Official performer accounts
- Promotional content links
- Social media presence tracking
- Profile verification

**Note**: This is supplementary metadata for performer profiles, not content hosting.

---

## Developer Resources

- 📚 **API Docs**: https://developer.twitter.com/en/docs
- 🔗 **API v2**: https://developer.twitter.com/en/docs/twitter-api
- 💰 **Pricing**: https://developer.twitter.com/en/products/twitter-api

---

## API Access

**Tiers**:
| Tier | Requests | Cost |
|------|----------|------|
| Free | 1,500 tweets/month read | Free |
| Basic | 10,000 tweets/month | $100/month |
| Pro | 1M tweets/month | $5,000/month |

**Note**: For performer profile lookup only, Free tier may suffice.

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  adult:
    social:
      twitter:
        enabled: false  # Disabled by default
        api_key: "${TWITTER_API_KEY:}"
        api_secret: "${TWITTER_API_SECRET:}"
        bearer_token: "${TWITTER_BEARER_TOKEN:}"

        # Rate limiting
        rate_limit:
          requests_per_15min: 15  # Free tier

        # Cache settings
        cache_ttl: "7d"  # Performer profiles don't change often
```

---

## Data Mapping

### Performer Social Presence

| Twitter Field | Revenge Field | Notes |
|---------------|---------------|-------|
| `id` | `twitter_id` | Unique identifier |
| `username` | `twitter_handle` | @username |
| `name` | `display_name` | Display name |
| `description` | `bio` | Profile bio |
| `profile_image_url` | `avatar_url` | Profile picture |
| `verified` | `is_verified` | Blue checkmark |
| `public_metrics.followers_count` | `followers` | Follower count |
| `url` | `website_url` | Profile link |

---

## API Endpoints

### Lookup User by Username

```
GET https://api.twitter.com/2/users/by/username/:username
```

**Response**:
```json
{
  "data": {
    "id": "123456789",
    "name": "Display Name",
    "username": "performer_handle",
    "description": "Bio text...",
    "profile_image_url": "https://pbs.twimg.com/...",
    "verified": true,
    "public_metrics": {
      "followers_count": 50000,
      "following_count": 100,
      "tweet_count": 5000
    }
  }
}
```

---

## Implementation

### Twitter Client

```go
// internal/content/c/metadata/twitter.go
package metadata

type TwitterClient struct {
    httpClient  *http.Client
    bearerToken string
    cache       *cache.Cache
    logger      *slog.Logger
}

type TwitterProfile struct {
    ID              string `json:"id"`
    Username        string `json:"username"`
    Name            string `json:"name"`
    Description     string `json:"description"`
    ProfileImageURL string `json:"profile_image_url"`
    Verified        bool   `json:"verified"`
    FollowersCount  int    `json:"followers_count"`
}

func (c *TwitterClient) GetProfile(ctx context.Context, username string) (*TwitterProfile, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("twitter:profile:%s", username)
    if cached, ok := c.cache.Get(cacheKey); ok {
        return cached.(*TwitterProfile), nil
    }

    // Fetch from API
    url := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+c.bearerToken)

    // Add fields
    q := req.URL.Query()
    q.Add("user.fields", "description,profile_image_url,verified,public_metrics")
    req.URL.RawQuery = q.Encode()

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 429 {
        return nil, ErrRateLimited
    }

    var result struct {
        Data TwitterProfile `json:"data"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    // Cache result
    c.cache.Set(cacheKey, &result.Data, 7*24*time.Hour)

    return &result.Data, nil
}
```

---

## Database Schema

```sql
-- In schema 'c' (adult isolated)
CREATE TABLE c.performer_social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    performer_id UUID NOT NULL REFERENCES c.performers(id),
    platform VARCHAR(50) NOT NULL,  -- 'twitter', 'instagram', etc.
    handle VARCHAR(255) NOT NULL,
    external_id VARCHAR(100),
    display_name VARCHAR(255),
    verified BOOLEAN DEFAULT FALSE,
    followers_count INTEGER,
    profile_url TEXT,
    avatar_url TEXT,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(performer_id, platform)
);

CREATE INDEX idx_performer_social_platform
    ON c.performer_social_links(platform);
```

---

## Implementation Checklist

- [ ] **Twitter Client** (`internal/content/c/metadata/twitter.go`)
  - [ ] Bearer token auth
  - [ ] Profile lookup
  - [ ] Rate limiting
  - [ ] Caching

- [ ] **Social Service** (`internal/content/c/service/social.go`)
  - [ ] Link performer to Twitter
  - [ ] Sync profile data
  - [ ] Unlink account

- [ ] **River Jobs** (`internal/content/c/jobs/social.go`)
  - [ ] SyncPerformerSocialsJob
  - [ ] Batch sync

---

## Privacy Considerations

- Only public profile data is stored
- No tweet content is fetched
- Performers can opt-out
- Data is isolated in `c` schema
- Follows adult content isolation patterns

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| 401 | Invalid token | Check bearer token |
| 404 | User not found | Handle gracefully |
| 429 | Rate limited | Backoff and retry |
| 403 | Suspended account | Skip user |

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
- [Instagram Integration](INSTAGRAM.md)
- [OnlyFans Integration](ONLYFANS.md)
- [Pornhub Integration](PORNHUB.md)
- [TheNude Integration](THENUDE.md)

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

- [Instagram](INSTAGRAM.md)
- [FreeOnes](../metadata/adult/FREEONES.md)
- [StashDB](../metadata/adult/STASHDB.md)
