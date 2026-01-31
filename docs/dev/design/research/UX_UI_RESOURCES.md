# UX/UI Design & Frontend Resources

<!-- SOURCES: shadcn-svelte, svelte-runes, svelte5, sveltekit -->

<!-- DESIGN: research, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> **Status**: ✅ ALL 14 SOURCES FETCHED (2026-01-28)
>
> These authoritative sources inform frontend instruction files, component design, and user experience patterns for Revenge.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [1. Accessibility Standards](#1-accessibility-standards)
  - [W3C WCAG 2.2 (Web Content Accessibility Guidelines)](#w3c-wcag-22-web-content-accessibility-guidelines)
  - [ISO 9241-11:2018 (Ergonomics of Human-System Interaction)](#iso-9241-112018-ergonomics-of-human-system-interaction)
- [2. Usability Heuristics & Laws](#2-usability-heuristics-laws)
  - [Nielsen Norman Group - 10 Usability Heuristics](#nielsen-norman-group---10-usability-heuristics)
  - [Laws of UX](#laws-of-ux)
- [3. Design Systems & Component Libraries](#3-design-systems-component-libraries)
  - [Google Material Design 3 (M3 Expressive)](#google-material-design-3-m3-expressive)
  - [Apple Human Interface Guidelines (HIG)](#apple-human-interface-guidelines-hig)
  - [Microsoft Fluent Design System (Fluent 2)](#microsoft-fluent-design-system-fluent-2)
  - [Atlassian Design System](#atlassian-design-system)
  - [IBM Carbon Design System](#ibm-carbon-design-system)
- [4. Government & Institutional Standards](#4-government-institutional-standards)
  - [UK Government Design Principles (GDS)](#uk-government-design-principles-gds)
- [5. UX Research & Best Practices](#5-ux-research-best-practices)
  - [Interaction Design Foundation (IDF)](#interaction-design-foundation-idf)
  - [Baymard Institute (E-commerce UX Research)](#baymard-institute-e-commerce-ux-research)
  - [Smashing Magazine (UX Design)](#smashing-magazine-ux-design)
- [6. Web Standards & Patterns](#6-web-standards-patterns)
  - [web.dev Patterns (Google)](#webdev-patterns-google)
- [7. Summary: Application to Revenge](#7-summary-application-to-revenge)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [8. Quality Gates for Frontend Development](#8-quality-gates-for-frontend-development)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

---

## 1. Accessibility Standards

### W3C WCAG 2.2 (Web Content Accessibility Guidelines)
- **URL**: https://www.w3.org/WAI/standards-guidelines/wcag/
- **Version**: WCAG 2.2 (published October 2023, updated December 2024)
- **Status**: International Standard, ISO/IEC 40500:2025
- **Scope**: 13 guidelines under 4 principles (Perceivable, Operable, Understandable, Robust)
- **Conformance Levels**: A, AA, AAA (9 new success criteria in 2.2)
- **Key Changes**:
  - Added mobile accessibility
  - Cognitive accessibility improvements
  - 4.1.1 Parsing obsolete
  - Focus indicators, dragging movements, target size
- **Resources**:
  - Quick Reference: https://www.w3.org/WAI/WCAG22/quickref/
  - WCAG 2.2 Standard: https://www.w3.org/TR/WCAG22/
  - Understanding WCAG 2: Detailed guidance
  - Techniques for WCAG 2: Implementation patterns
  - Supplemental Guidance: Beyond baseline
- **Revenge Application**:
  - Level AA compliance target
  - Keyboard navigation for all controls
  - Screen reader optimization
  - Color contrast 4.5:1 (text), 3:1 (UI)
  - Focus visible indicators
  - Alt text for all images
  - Skip links for navigation
  - ARIA landmarks and roles

### ISO 9241-11:2018 (Ergonomics of Human-System Interaction)
- **URL**: https://www.iso.org/standard/63500.html
- **Edition**: 2 (March 2018, confirmed 2023)
- **Scope**: Usability definitions and concepts
- **Key Concepts**:
  - Usability as outcome of use (not inherent property)
  - Effectiveness, efficiency, satisfaction in context
  - User needs drive design
  - System, product, service applicability
- **Revenge Application**:
  - User research validation
  - Usability testing metrics (task success, time, satisfaction)
  - Context of use analysis (device types, environments)

---

## 2. Usability Heuristics & Laws

### Nielsen Norman Group - 10 Usability Heuristics
- **URL**: https://www.nngroup.com/articles/ten-usability-heuristics/
- **Author**: Jakob Nielsen (1994, updated 2020)
- **Status**: Industry standard for heuristic evaluation
- **The 10 Heuristics**:
  1. **Visibility of System Status** - Feedback within reasonable time
  2. **Match System & Real World** - User language, natural mapping
  3. **User Control & Freedom** - Undo/redo, emergency exits
  4. **Consistency & Standards** - Platform conventions, Jakob's Law
  5. **Error Prevention** - Eliminate error-prone conditions, confirm destructive actions
  6. **Recognition > Recall** - Visible options, minimize memory load
  7. **Flexibility & Efficiency** - Shortcuts for experts, personalization
  8. **Aesthetic & Minimalist Design** - No irrelevant information
  9. **Error Recognition & Recovery** - Plain language errors, solutions
  10. **Help & Documentation** - Easy to search, contextual, concrete steps
- **Resources**:
  - Free posters: https://media.nngroup.com/media/articles/attachments/Jakob's10UsabilityHeuristics_AllPosters_5.zip
  - Video explanations (2-3 min each)
  - Application examples (complex apps, VR, video games)
- **Revenge Application**:
  - Heuristic evaluation of all UIs
  - Progress indicators for transcoding
  - Undo/cancel for destructive actions
  - Consistent navigation across modules
  - Error messages with solutions
  - Contextual help tooltips

### Laws of UX
- **URL**: https://lawsofux.com/
- **Author**: Jon Yablonski (O'Reilly book)
- **Scope**: 26 laws/principles for UI design
- **Key Laws for Revenge**:
  - **Aesthetic-Usability Effect**: Beautiful design perceived as more usable
  - **Fitts's Law**: Target acquisition time = f(distance, size) → Larger touch targets
  - **Hick's Law**: Decision time increases with choices → Limit navigation options
  - **Jakob's Law**: Users expect site to work like others → Use conventions
  - **Miller's Law**: Working memory 7±2 items → Chunk information
  - **Serial Position Effect**: Remember first & last items → Key actions at ends
  - **Von Restorff Effect**: Distinctive items remembered → Highlight primary CTA
  - **Doherty Threshold**: <400ms response time boosts productivity
  - **Goal-Gradient Effect**: Motivation increases near goal → Show progress
  - **Peak-End Rule**: Experiences judged by peak & end → Optimize critical moments
  - **Tesler's Law**: Complexity conservation → Simplify user-facing, accept backend complexity
  - **Pareto Principle**: 80% effects from 20% causes → Focus on common tasks
  - **Postel's Law**: Be liberal in inputs, conservative in outputs
  - **Zeigarnik Effect**: Unfinished tasks remembered → Save drafts, show incomplete
  - **Choice Overload**: Too many options → paralysis
  - **Cognitive Load**: Mental resources to understand UI → Minimize
  - **Flow**: Immersed energized focus → Remove friction
- **Revenge Application**:
  - Large tap targets (48×48px minimum) for mobile
  - Limit main navigation to 5-7 items
  - Chunk settings into logical groups
  - Primary actions (Play, Add) visually prominent
  - Progress bars for uploads/transcoding
  - <400ms API response target
  - Autosave for playlist editing
  - Conservative validation (accept variations), strict output

---

## 3. Design Systems & Component Libraries

### Google Material Design 3 (M3 Expressive)
- **URL**: https://m3.material.io/
- **Status**: Latest evolution (I/O 2025 update)
- **Philosophy**: Emotion-driven UX with vibrant colors, intuitive motion, adaptive components
- **Key Features**:
  - **M3 Expressive update** (2025):
    - Vibrant color system (extended palettes)
    - Motion physics (easier-to-implement, token-powered transitions)
    - Shape library (35 shapes with built-in morph motion)
    - Flexible typography
  - **New Components**: Toolbars, Split buttons, Progress indicators, Button groups
  - **Updated Components**: 14 total (existing components refreshed)
- **Libraries**: Web, Android, Flutter, Figma
- **Revenge Application**:
  - Inspiration for component animations (FAB transitions, button states)
  - Motion design for quality switching
  - Progress indicators for transcoding
  - Color palette generation
  - **NOT using Material directly** (using shadcn-svelte), but adopting motion principles

### Apple Human Interface Guidelines (HIG)
- **URL**: https://developer.apple.com/design/human-interface-guidelines/
- **Platforms**: iOS, iPadOS, macOS, watchOS, tvOS, visionOS
- **Core Principles**: Hierarchy, Harmony, Consistency
- **Design Fundamentals**: App icons, color, materials, layout, typography, accessibility
- **Revenge Application**:
  - iOS/macOS native app patterns (future)
  - Touch gesture conventions
  - System integration (PiP, AirPlay)
  - SF Symbols icon style inspiration

### Microsoft Fluent Design System (Fluent 2)
- **URL**: https://fluent2.microsoft.design/
- **Platforms**: Web (React), iOS, Android, Windows
- **Philosophy**: Let creativity flow, accessible & inclusive
- **Revenge Application**:
  - Cross-platform component patterns
  - Windows native app (future)
  - **NOT using Fluent directly**, but reference for Windows UX conventions

### Atlassian Design System
- **URL**: https://atlassian.design/
- **Philosophy**: Better teamwork by design, unified design language
- **Key Features**: Rovo AI patterns, Token-based theming, Component composition
- **Revenge Application**:
  - AI chat patterns for future AI features
  - Token-based theming
  - Component composition patterns

### IBM Carbon Design System
- **URL**: https://www.carbondesignsystem.com/
- **Philosophy**: Adaptable system, best practices of UI design, open-source
- **Libraries**: Web Components, React, Angular, Vue, Svelte
- **Key Features**: AI Chat v1, Comprehensive component catalog, Accessibility-first
- **Revenge Application**:
  - Svelte component patterns
  - AI chat interface (future)
  - Design token structure

---

## 4. Government & Institutional Standards

### UK Government Design Principles (GDS)
- **URL**: https://www.gov.uk/guidance/government-design-principles
- **Published**: April 2012, updated April 2025
- **The 11 Principles**:
  1. **Start with user needs** - Research, don't assume
  2. **Do less** - Reusable platforms, link to others
  3. **Design with data** - Analytics built-in, always on
  4. **Do hard work to make it simple** - Simplicity > "always been that way"
  5. **Iterate. Then iterate again** - MVP, alpha→beta→live, learn from failures
  6. **This is for everyone** - Accessible design = good design
  7. **Understand context** - Not designing for screen, for people
  8. **Build digital services, not websites** - Connect real world
  9. **Be consistent, not uniform** - Shared language/patterns, but improve when needed
  10. **Make things open** - Share code, designs, ideas, failures
  11. **Minimise environmental impact** - Reduce energy, water, materials (NEW 2025)
- **Revenge Application**:
  - User research-driven features
  - Iterative development (alpha modules first)
  - Accessibility priority
  - Open-source ethos
  - Energy-efficient transcoding (Blackbeard optimization)

---

## 5. UX Research & Best Practices

### Interaction Design Foundation (IDF)
- **URL**: https://www.interaction-design.org/literature/topics/ux-design
- **Scope**: World's largest UX education community (1.2M+ enrollments)
- **Key Content**: User Experience definition, UX vs UI, ISO 9241-210, User-centered design process
- **Revenge Application**:
  - User research methodology
  - Persona creation for target users (media enthusiasts, families, power users)
  - Iterative design process
  - Usability testing protocols

### Baymard Institute (E-commerce UX Research)
- **URL**: https://baymard.com/blog
- **Scope**: 200,000+ hours of UX research, 18,000+ design examples
- **Focus**: E-commerce usability (applicable to media libraries)
- **Research Topics**: Homepage navigation, Product lists, Mobile UX, Search & filtering
- **Revenge Application**:
  - Library navigation patterns
  - Filtering & sorting best practices
  - Mobile UX optimization
  - Search interface design

### Smashing Magazine (UX Design)
- **URL**: https://www.smashingmagazine.com/category/ux-design/
- **Scope**: Professional web design & UX articles (57+ UX design articles)
- **Key Topics**: Design patterns, AI in design, Infinite scroll, Accessibility
- **Revenge Application**:
  - Infinite scroll for media libraries
  - Design system collaboration patterns
  - Research-driven feature development

---

## 6. Web Standards & Patterns

### web.dev Patterns (Google)
- **URL**: https://web.dev/patterns/
- **Scope**: Modern web API patterns with browser support (Baseline)
- **Pattern Categories**:
  - **Animation**: CSS/JS animations with accessibility, user preferences
  - **Clipboard**: Copy/paste patterns
  - **Components**: Cross-browser UI components
  - **Files & Directories**: File upload, drag-drop, directory access
  - **Layout**: Modern CSS (Grid, Flexbox, Container Queries)
  - **Media**: Video, audio, images (lazy loading, responsive)
  - **Theming**: Color management, dark mode, CSS custom properties
  - **Web Apps**: PWA patterns (service workers, manifest, offline)
- **Revenge Application**:
  - Responsive layout patterns (media grids)
  - Dark/light mode theming
  - File upload for custom artwork
  - Video player controls
  - PWA offline support
  - Animation performance (GPU-accelerated)

---

## 7. Summary: Application to Revenge

| Resource | Primary Use in Revenge |
|----------|------------------------|
| **WCAG 2.2** | Accessibility compliance (AA level), keyboard navigation, screen readers |
| **ISO 9241-11** | Usability testing framework, effectiveness/efficiency metrics |
| **Nielsen's Heuristics** | Heuristic evaluation of all UIs, error handling, consistency |
| **Laws of UX** | Touch target sizing, navigation chunking, progress visualization |
| **Material Design 3** | Motion design inspiration, progress indicators, animation principles |
| **Apple HIG** | iOS/macOS app patterns, touch gestures, system integration |
| **Fluent 2** | Cross-platform component patterns, Windows UX |
| **Atlassian** | AI chat patterns, token-based theming, component composition |
| **Carbon** | Svelte component patterns, AI features, design tokens |
| **GDS Principles** | User research-driven development, iterative design, accessibility |
| **IDF** | User research methodology, persona creation, iterative design |
| **Baymard** | Library navigation, filtering/sorting, mobile UX, search |
| **Smashing** | Infinite scroll, design patterns, accessibility |
| **web.dev** | Responsive layouts, theming, PWA, file handling, animation |


---

## 8. Quality Gates for Frontend Development

Before ANY frontend component is merged:

- [ ] WCAG 2.2 Level AA compliant (axe-core passes)
- [ ] Keyboard navigable (tab order, focus visible)
- [ ] Screen reader tested (NVDA/JAWS/VoiceOver)
- [ ] Touch targets ≥48×48px (Fitts's Law)
- [ ] Color contrast ≥4.5:1 text, ≥3:1 UI
- [ ] Dark/light mode support
- [ ] Mobile-first responsive
- [ ] `prefers-reduced-motion` respected
- [ ] Error messages with solutions (Nielsen #9)
- [ ] Undo/cancel for destructive actions (Nielsen #3)
- [ ] Consistent with design system
- [ ] i18n keys used (no hardcoded strings)
- [ ] Documented in Storybook (component catalog)
