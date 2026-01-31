# Professional Versioning Approaches for MVP

**Date**: 2026-01-31
**Purpose**: Guide versioning strategy decision

---

## SemVer 2.0.0 Basics (Pre-1.0)

**Key principle**: Versions before v1.0.0 are "development phase"
- **v0.x.x** = "API not stable, expect breaking changes"
- **v1.0.0** = "First public API, stability promise begins"
- Pre-1.0: MINOR can have breaking changes
- Post-1.0: Only MAJOR can have breaking changes

---

## Three Professional Approaches

### Approach 1: Lean Startup (MVP = v0.1.0)

**Philosophy**: "Ship fast, iterate, learn from users"

```
v0.1.0  MVP - Minimal but functional
v0.2.0  First iteration based on feedback
v0.3.0  Second iteration
...
v1.0.0  Stable, feature-complete
```

**When to use**:
- ✅ Startups needing user validation
- ✅ When you can iterate weekly
- ✅ Internal tools (forgiving users)
- ✅ Experimental projects

**Risks**:
- ❌ Too minimal = users don't return
- ❌ First impression matters
- ❌ Hard to iterate on media server (infrastructure-heavy)

**Example**: Tailwind CSS v0.1.0 was very minimal, but web framework could iterate fast

---

### Approach 2: Staged Foundation (MVP = v0.3-0.5) ⭐ RECOMMENDED

**Philosophy**: "Build solid foundation, then call it MVP"

```
v0.1.x  Phase 1: Foundation
        - Infrastructure (DB, cache, auth)
        - Core services working
        - No content yet

v0.2.x  Phase 2: First Features
        - One content type (movies OR tv)
        - Basic playback
        - Minimal UI

v0.3.x  Phase 3: MVP ← Actually usable!
        - Movies + TV both working
        - Full SvelteKit UI
        - Metadata integration
        - Direct play + HLS

v0.4.x  Phase 4: Enhancements
        - More content types
        - Better features

v1.0.0  Stable: Full design spec
```

**When to use**:
- ✅ Complex systems (media servers, databases, platforms)
- ✅ When users expect polish
- ✅ Infrastructure-heavy projects
- ✅ When "barely works" isn't good enough

**Benefits**:
- ✅ MVP is actually viable (won't scare users away)
- ✅ Foundation is solid before adding features
- ✅ Each minor version has clear purpose
- ✅ Matches your VERSIONING.md plan

**Example**:
- PostgreSQL 8.0 was "MVP" but 8.x had years of foundation work
- Kubernetes 0.x had many releases before calling anything "production ready"
- Docker 0.1-0.8 were foundation, 0.9+ were "usable for real"

---

### Approach 3: Quality Gates (Alpha/Beta/RC)

**Philosophy**: "Focus on stability level, not features"

```
v0.1.0-alpha.1   Early builds (unstable)
v0.2.0-alpha.5   More features (still unstable)
v0.3.0-beta.1    Feature complete (testing)
v0.4.0-rc.1      Release candidate (almost ready)
v1.0.0           Stable release
```

**When to use**:
- ✅ When quality matters more than features
- ✅ Critical systems (medical, financial)
- ✅ When you need clear signal to users

**Downsides**:
- ❌ Labels can be confusing (what's difference between alpha and beta?)
- ❌ Harder to convey progress
- ❌ Pre-release labels are optional in SemVer

---

## Recommendation for Revenge

**Use Approach 2: Staged Foundation (MVP = v0.3.x)**

### Why?

1. **Media servers need polish**
   - Users compare to Plex/Jellyfin
   - "Barely works" won't get adoption
   - First impression matters

2. **Infrastructure-heavy**
   - Auth, DB, cache, metadata must work well
   - Can't ship "foundation only" as MVP
   - Need content types actually working

3. **Matches existing VERSIONING.md**
   - Already planned as v0.3.x = MVP Complete
   - Makes sense: v0.1 = infra, v0.2 = partial, v0.3 = usable

4. **Aligns with design completeness**
   - Design is 80% done
   - Code is 0% done
   - Need v0.1-v0.2 to build foundation

---

## Proposed Versioning (Aligned with VERSIONING.md)

### v0.1.x - Core Foundation
**Goal**: Infrastructure works, no UI yet

**Deliverables**:
- ✅ PostgreSQL + migrations
- ✅ Dragonfly cache + River jobs
- ✅ Auth, User, Session, RBAC services
- ✅ Library scanner (file detection)
- ✅ Health checks, logging, metrics
- ⚠️ API-only (no frontend)
- ⚠️ No playback yet

**Exit criteria**: Backend services all pass tests, API responds

---

### v0.2.x - First Content Type
**Goal**: One content type end-to-end

**Deliverables**:
- ✅ Movie module complete (OR TV module)
- ✅ TMDB integration
- ✅ Radarr integration (OR Sonarr)
- ✅ Metadata fetching & storage
- ✅ Basic SvelteKit UI (browse movies)
- ✅ Direct play (no transcoding)
- ⚠️ Only one content type
- ⚠️ Limited playback features

**Exit criteria**: Can add movies, browse, and play them

---

### v0.3.x - MVP Complete ⭐
**Goal**: Movies + TV both working, ready for real use

**Deliverables**:
- ✅ Movie module complete
- ✅ TV Show module complete
- ✅ TMDB + TheTVDB integration
- ✅ Radarr + Sonarr integration
- ✅ Full SvelteKit UI
- ✅ Direct play + HLS/DASH
- ✅ Search (Typesense)
- ✅ OIDC authentication
- ✅ User management
- ⚠️ No transcoding (Blackbeard not ready)
- ⚠️ No advanced features (skip intro, etc.)

**Exit criteria**: Can replace Jellyfin for basic movie/TV use

---

### v0.4.x - Third Content Type
**Goal**: Add music or another content type

**Deliverables**:
- ✅ Music module (OR audiobook/podcast)
- ✅ Music metadata providers
- ✅ Lidarr integration (if music)
- ✅ Audio playback

---

### v0.5.x - Transcoding
**Goal**: Blackbeard integration

**Deliverables**:
- ✅ Blackbeard transcoding
- ✅ Format conversion
- ✅ Quality profiles

---

### v0.6-0.9.x - Advanced Features
**Goal**: Polish, advanced features, more content types

**Deliverables**:
- ✅ Skip intro, trickplay, syncplay
- ✅ Photos, Comics, Books modules
- ✅ LiveTV/DVR
- ✅ QAR (adult content)
- ✅ Plugin system
- ✅ Performance optimization

---

### v1.0.0 - Stable Release
**Goal**: Full design spec implemented, production-ready

**Criteria**:
- ✅ All modules from design complete
- ✅ 80%+ test coverage
- ✅ All integrations working
- ✅ Performance benchmarks met
- ✅ Documentation complete
- ✅ API stable (no more breaking changes)

---

## Comparison

| Approach | MVP Version | Time to MVP | MVP Quality | Post-MVP Path |
|----------|-------------|-------------|-------------|---------------|
| Lean Startup | v0.1.0 | 1-2 months | Barely usable | Iterate fast |
| Staged Foundation | v0.3.0 | 4-6 months | Actually good | Measured growth |
| Quality Gates | v0.x.0-beta | Variable | Depends on label | Confusing |

---

## Industry Examples

### Lean Startup (MVP = v0.1.0):
- **Next.js**: v0.1.0 was very minimal, but web frameworks can iterate fast
- **Vite**: v0.1.0 was proof of concept, v1.0.0 came quickly after

### Staged Foundation (MVP = v0.3+):
- **PostgreSQL**: Many 0.x releases before "production ready"
- **Kubernetes**: 0.1-0.20+ before considering it "ready"
- **Prometheus**: 0.1-0.16 before calling it stable at 1.0

### Your Competition:
- **Jellyfin**: Fork of Emby, started at v10.x (inherited maturity)
- **Plex**: Proprietary, not SemVer, but was polished from start

---

## Recommendation

**Adopt VERSIONING.md plan with v0.3.x = MVP**

**Rationale**:
1. Media servers need quality (user expectations)
2. Infrastructure is complex (need foundation)
3. Makes sense: build → test → release
4. Matches industry norms for complex systems
5. Already documented in VERSIONING.md

**Update needed**:
- VERSIONING.md is correct
- Just need to clarify what each version delivers
- Create detailed roadmap: v0.1.x → v0.2.x → v0.3.x MVP → v1.0.0

---

## Alternative If You Want Fast Release

If you absolutely need v0.1.0 to be "MVP":

**Redefine v0.1.0 scope** (much smaller):
- Only Movies (no TV)
- API + basic UI only
- Direct play only (no HLS)
- Single OIDC provider
- No Radarr integration
- "Proof of concept" level

But this risks:
- ❌ Not actually useful
- ❌ Users won't return
- ❌ Negative first impression

**Better**: Keep v0.3.x as MVP, accept it takes longer to get there

---

## Decision

**Which approach do you prefer?**

1. **Keep VERSIONING.md** (MVP = v0.3.x) ← Recommended
2. **Compress to v0.1.0 MVP** (much smaller scope)
3. **Hybrid** (v0.1.0 = alpha, v0.3.0 = beta, v1.0.0 = stable)
