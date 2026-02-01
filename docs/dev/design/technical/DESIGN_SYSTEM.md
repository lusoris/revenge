

---
sources:
  - name: shadcn-svelte
    url: ../../sources/frontend/shadcn-svelte.md
    note: Component library
  - name: Bits UI
    url: ../../sources/frontend/bits-ui.md
    note: Headless UI primitives
  - name: Tailwind CSS 4
    url: ../../sources/frontend/tailwindcss.md
    note: Utility-first CSS
  - name: Lucide Icons
    url: https://lucide.dev/icons/
    note: Icon library (standard mode)
  - name: Vidstack Player
    url: https://www.vidstack.io/docs/player
    note: Video player components
  - name: WCAG 2.1 Guidelines
    url: https://www.w3.org/WAI/WCAG21/quickref/
    note: Accessibility standards (targeting AAA)
  - name: Playfair Display
    url: https://fonts.google.com/specimen/Playfair+Display
    note: Serif accent font
  - name: Inter
    url: https://fonts.google.com/specimen/Inter
    note: UI body font
design_refs:
  - title: FRONTEND
    path: ../technical/FRONTEND.md
  - title: 04_PLAYER_ARCHITECTURE
    path: ../architecture/04_PLAYER_ARCHITECTURE.md
  - title: USER_EXPERIENCE_FEATURES
    path: ../features/shared/USER_EXPERIENCE_FEATURES.md
---

## Table of Contents

- [Design System](#design-system)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Implementation Checklist](#implementation-checklist)
    - [Phase 1: Foundation](#phase-1-foundation)
    - [Phase 2: Components](#phase-2-components)
    - [Phase 3: Navigation](#phase-3-navigation)
    - [Phase 4: Motion](#phase-4-motion)
    - [Phase 5: Accessibility](#phase-5-accessibility)
    - [Phase 6: Polish](#phase-6-polish)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Design System


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Queen Anne's Revenge themed design system for Revenge frontend

Design system index with pirate-themed elegance:
- **Theme**: Sleek & Professional with Queen Anne's Revenge styling
- **Colors**: Arr-aligned module colors with full theme shift
- **Typography**: Inter body + serif accent titles (Playfair)
- **Motion**: Expressive micro-interactions, delightful details
- **Accessibility**: WCAG 2.1 AAA compliant
- **Easter Eggs**: Achievement-unlocked pirate mode


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete design system from discovery session |
| Sources | ✅ | All design resources documented |
| Instructions | ✅ | Implementation guidelines in sub-documents |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies

<!-- Dependency list -->





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys

<!-- Configuration keys -->




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**




## Implementation Checklist
### Phase 1: Foundation
- [ ] Set up SvelteKit project with TypeScript
- [ ] Configure Tailwind CSS 4 with custom theme
- [ ] Install Inter + Playfair Display fonts
- [ ] Set up CSS variables for module theming
- [ ] Implement dark/light mode toggle
### Phase 2: Components
- [ ] Install and configure shadcn-svelte
- [ ] Build MediaCard with hover overlay
- [ ] Create HeroBanner with personalization
- [ ] Configure Vidstack player
- [ ] Build collapsible rail sidebar
### Phase 3: Navigation
- [ ] Implement command palette (Ctrl+K)
- [ ] Build floating filter panel
- [ ] Mobile bottom navigation
- [ ] TV spatial navigation
- [ ] Keyboard shortcuts
### Phase 4: Motion
- [ ] Implement loading states (skeleton, shimmer, blurhash)
- [ ] Add micro-interactions
- [ ] Page transitions
- [ ] Success/error feedback
### Phase 5: Accessibility
- [ ] Implement skip links
- [ ] Add ARIA labels
- [ ] Subtitle customization
- [ ] Screen reader testing
- [ ] Contrast verification
### Phase 6: Polish
- [ ] Achievement system
- [ ] Pirate icons
- [ ] Pirate-speak translation
- [ ] Celebration animations



## Related Documentation
### Design Documents
- [FRONTEND](../technical/FRONTEND.md)
- [04_PLAYER_ARCHITECTURE](../architecture/04_PLAYER_ARCHITECTURE.md)
- [USER_EXPERIENCE_FEATURES](../features/shared/USER_EXPERIENCE_FEATURES.md)

### External Sources
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md) - Component library
- [Bits UI](../../sources/frontend/bits-ui.md) - Headless UI primitives
- [Tailwind CSS 4](../../sources/frontend/tailwindcss.md) - Utility-first CSS
- [Lucide Icons](https://lucide.dev/icons/) - Icon library (standard mode)
- [Vidstack Player](https://www.vidstack.io/docs/player) - Video player components
- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/) - Accessibility standards (targeting AAA)
- [Playfair Display](https://fonts.google.com/specimen/Playfair+Display) - Serif accent font
- [Inter](https://fonts.google.com/specimen/Inter) - UI body font

