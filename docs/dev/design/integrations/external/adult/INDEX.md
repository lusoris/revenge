# Adult External Services

> Social media and supplementary services for adult content

**⚠️ Adult Content**: All data stored in isolated PostgreSQL schema `c`.
API namespace: `/api/v1/c/`

---

## Overview

Adult external services provide:
- Performer social media links
- Supplementary metadata
- Platform presence tracking

---

## Services

### Social Media

| Provider | Type | Status |
|----------|------|--------|
| [Twitter/X](TWITTER_X.md) | Social | 🔴 Planned |
| [Instagram](INSTAGRAM.md) | Social | 🔴 Planned |

### Metadata (Supplementary)

| Provider | Type | Status |
|----------|------|--------|
| [Pornhub](PORNHUB.md) | Platform | 🔴 Planned |
| [OnlyFans](ONLYFANS.md) | Platform | 🔴 Planned |
| [TheNude](THENUDE.md) | Metadata | 🔴 Planned |

---

## Service Details

### Twitter/X
**Performer social presence**

- Official accounts
- Promotional links
- Verification status
- Limited API access (Free tier)

### Instagram
**Performer social presence**

- Profile links only
- No API integration (deprecated)
- Manual link management

### Pornhub
**Platform metadata**

- Channel information
- View statistics
- ⚠️ Unofficial API

### OnlyFans
**Creator platform**

- Profile links
- No API available
- Manual entry only

### TheNude
**Performer metadata**

- Biographical data
- Career information
- ⚠️ Web scraping required

---

## Data Isolation

All data stored in `c` schema:

```sql
-- Performer social links
c.performer_social_links (
    performer_id,
    platform,
    handle,
    profile_url,
    verified
)
```

---

## Configuration

```yaml
integrations:
  adult:
    social:
      twitter:
        enabled: false
      instagram:
        enabled: false
    metadata:
      pornhub:
        enabled: false
      onlyfans:
        enabled: false
```

---

## Privacy Notes

- Only public profile URLs stored
- No content scraping
- User-provided links primarily
- Performers can opt-out

---

## Related Documentation

- [External Overview](../INDEX.md)
- [Adult Metadata](../../metadata/adult/INDEX.md)
- [Adult Wiki](../../wiki/adult/INDEX.md)
