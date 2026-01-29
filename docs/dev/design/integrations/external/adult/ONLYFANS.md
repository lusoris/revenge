# OnlyFans Integration

> Adult content platform subscription service - profile link tracking

**Service**: OnlyFans (https://onlyfans.com)
**API**: No public API (web scraping NOT recommended)
**Category**: External Platform (Adult Content - Content Platform)
**Priority**: 🔴 VERY LOW (No public API, link tracking only)
**Status**: 🔴 DESIGN PHASE

---

## Overview

**OnlyFans** is a subscription-based content platform where creators (including adult performers) share exclusive content with paying subscribers. OnlyFans does not provide a public API and actively blocks scrapers.

**Key Features**:
- **Creator profiles**: Performer profiles with subscription options
- **Exclusive content**: Subscription-based content (photos, videos)
- **Direct messaging**: Subscriber communication
- **Tips/PPV**: Pay-per-view content, tipping
- **No public API**: No official public API available

**Use Cases**:
- **Profile link tracking**: Store OnlyFans profile URLs for performers
- **External link enrichment**: Add OnlyFans links to performer profiles
- **Verified status**: Track if performer has OnlyFans account

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`)
- **API namespace**: `/api/v1/c/external/onlyfans/*` (NOT `/api/v1/external/onlyfans/*`)
- **Module location**: `internal/content/c/external/onlyfans/` (NOT `internal/service/external/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

**⚠️ IMPORTANT: Link Tracking ONLY**:
- **NO scraping**: OnlyFans actively blocks scrapers (legal/ToS issues)
- **NO API**: No official public API
- **Link storage only**: Store OnlyFans profile URLs (from FreeOnes, StashDB, user input)
- **Display only**: Show OnlyFans links in UI (external link, opens in browser)

---

## Developer Resources

### API Status
- **Official API**: NONE (no public API)
- **Web Scraping**: NOT RECOMMENDED (actively blocked, legal issues)
- **Content Partners**: Private API exists for OnlyFans partners only

### Integration Approach
- **Link tracking only**: Do NOT scrape OnlyFans
- **Source links from**:
  - StashDB (`performer.urls` field with OnlyFans links)
  - FreeOnes (scrape OnlyFans links from FreeOnes performer pages)
  - User input (users can manually add OnlyFans links)
- **Display links**: Show OnlyFans icon with external link (opens in browser)

---

## Implementation Checklist

### Phase 1: Link Tracking (Adult Content - c schema)
- [ ] Link storage table: `c.performer_external_urls`
  - Fields: `platform='onlyfans'`, `url`, `verified`
- [ ] Import OnlyFans links from StashDB (`performer.urls` field)
- [ ] Import OnlyFans links from FreeOnes scraping
- [ ] Manual link input (user can add OnlyFans link to performer profile)
- [ ] Link validation (verify URL format `https://onlyfans.com/{username}`)
- [ ] **c schema storage**: `c.performers.external_urls` (OnlyFans profile link)

### Phase 2: UI Integration
- [ ] Display OnlyFans links in performer profile (external links section)
- [ ] OnlyFans icon (recognizable OF logo)
- [ ] External link (opens in new tab/browser)
- [ ] Verified badge (if link verified)

### Phase 3: Link Verification
- [ ] User verification (users can verify OnlyFans links)
- [ ] Admin moderation (admins can approve/reject links)
- [ ] Community flagging (flag incorrect/dead links)

---

## Integration Pattern

### Link Tracking Flow
```
OnlyFans link source (StashDB OR FreeOnes OR user input)
        ↓
Extract OnlyFans URL (https://onlyfans.com/{username})
        ↓
Validate URL format
        ↓
Store in c.performers.external_urls:
  - performer_id: UUID
  - source: 'stashdb' OR 'freeones' OR 'user'
  - platform: 'onlyfans'
  - url: 'https://onlyfans.com/username'
  - verified: FALSE (pending verification)
        ↓
Display in UI:
  - Performer profile → External links section
  - OnlyFans icon → Click opens https://onlyfans.com/username in new tab
```

### Link Verification Flow
```
User clicks "Verify Link" on OnlyFans link
        ↓
User confirms: "I verified this OnlyFans link is correct"
        ↓
Update c.performers.external_urls:
  - verified: TRUE
  - verified_by: user_id
  - verified_at: timestamp
        ↓
Display verified badge (✓) next to OnlyFans link
```

---

## Related Documentation

- [FREEONES.md](./FREEONES.md) - FreeOnes performer database (OnlyFans links)
- [PORNHUB.md](./PORNHUB.md) - Pornhub content platform
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata (performer.urls field)

---

## Notes

### NO Scraping (Legal/ToS Issues)
- **OnlyFans ToS**: Prohibits scraping, automation, data collection
- **Legal risks**: Scraping OnlyFans can result in legal action
- **Active blocking**: OnlyFans actively blocks scrapers (CAPTCHA, IP bans)
- **Revenge policy**: Do NOT scrape OnlyFans (link tracking ONLY)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.external_urls` (OnlyFans profile link)
  - NO data in public schema
- **API namespace**: `/api/v1/c/external/onlyfans/*` (isolated)
- **Module location**: `internal/content/c/external/onlyfans/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### Link Sources
1. **StashDB**: `performer.urls` field (OnlyFans links from StashDB)
   - Import during StashDB performer sync
   - Filter `site.name == "OnlyFans"`

2. **FreeOnes**: OnlyFans links on FreeOnes performer pages
   - Scrape OnlyFans links from FreeOnes external links section
   - Validate URL format

3. **User input**: Users can manually add OnlyFans links
   - Input field: "OnlyFans URL"
   - Validation: URL format `https://onlyfans.com/{username}`
   - Moderation: Requires admin approval (prevent spam)

### External URLs Table Schema
```sql
-- Already defined in FREEONES.md, reuse for OnlyFans
CREATE TABLE c.performer_external_urls (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  performer_id UUID NOT NULL REFERENCES c.performers(id) ON DELETE CASCADE,
  source VARCHAR(50) NOT NULL, -- 'stashdb', 'freeones', 'user'
  platform VARCHAR(50) NOT NULL, -- 'onlyfans', 'fansly', 'twitter', etc.
  url TEXT NOT NULL,
  verified BOOLEAN DEFAULT FALSE,
  verified_by UUID REFERENCES users(id),
  verified_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(performer_id, platform, url)
);
```

### URL Validation
- **Format**: `https://onlyfans.com/{username}`
- **Username**: Alphanumeric + underscore (e.g., `performer_name`)
- **Normalization**: Lowercase, remove trailing slash
- **Validation**: Regex `^https://onlyfans\.com/[a-zA-Z0-9_]+$`

### Link Display (UI)
- **External links section**: Performer profile page
- **OnlyFans icon**: Recognizable OF logo (blue/white)
- **External link**: `<a href="https://onlyfans.com/username" target="_blank" rel="noopener noreferrer">`
- **Verified badge**: ✓ if `verified=TRUE`
- **Tooltip**: "OnlyFans: @username (Verified)" OR "OnlyFans: @username (Unverified)"

### Community Moderation
- **Link verification**: Users can verify links (confirm correctness)
- **Admin moderation**: Admins can approve/reject user-submitted links
- **Flagging**: Users can flag incorrect/dead links
- **Removal**: Admins can remove spam/incorrect links

### Priority: VERY LOW
- **OnlyFans**: Link tracking only (no metadata extraction)
- **Use case**: Display OnlyFans profile link in performer external links
- **Implementation**: VERY LOW priority (basic link storage)

### Alternative Platforms
- **Fansly**: Similar to OnlyFans (subscription-based content)
- **ManyVids**: Content platform with creator profiles
- **Both**: Same approach (link tracking only, no scraping)

### Fallback Strategy (Adult Content Platform Links)
- **Order**: StashDB (primary urls field) → FreeOnes (external links scraping) → User input (manual)
- **OnlyFans**: Link tracking only (display external link)
