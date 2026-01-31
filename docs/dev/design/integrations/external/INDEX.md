# External Services

← Back to [Integrations](../)

> Third-party integrations and data sources

**Source of Truth**: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md)

---

## Overview

External services provide supplementary features beyond core media server functionality:
- Notifications
- Download management
- Sync services
- Webhooks

---

## Categories

### 📱 Notifications
Push notifications and alerts.

*(Planned for future implementation)*

### 🔗 Webhooks
Event-driven integrations.

*(Planned for future implementation)*

### 🔞 [Adult Services](adult/INDEX.md)
Adult content social and metadata services.

| Provider | Type | Status |
|----------|------|--------|
| [Twitter/X](adult/TWITTER_X.md) | Social | 🔴 Planned |
| [Instagram](adult/INSTAGRAM.md) | Social | 🔴 Planned |
| [FreeOnes](../metadata/adult/FREEONES.md) | Metadata | 🟡 Planned |
| [Pornhub](adult/PORNHUB.md) | Metadata | 🔴 Planned |
| [OnlyFans](adult/ONLYFANS.md) | Metadata | 🔴 Planned |
| [TheNude](adult/THENUDE.md) | Metadata | 🔴 Planned |

---

## Service Types

| Type | Description | Examples |
|------|-------------|----------|
| Metadata | Content information | FreeOnes, IAFD |
| Social | Social media links | Twitter, Instagram |
| Notifications | Push alerts | Discord, Telegram |
| Webhooks | Event triggers | Custom endpoints |

---

## Configuration

```yaml
external:
  notifications:
    discord:
      enabled: false
      webhook_url: "${DISCORD_WEBHOOK}"

  webhooks:
    enabled: false
    endpoints: []
```

---

## Related Documentation

- [Metadata Providers](../metadata/INDEX.md)
- [Scrobbling Services](../scrobbling/INDEX.md)
- [Servarr Stack](../servarr/INDEX.md)
