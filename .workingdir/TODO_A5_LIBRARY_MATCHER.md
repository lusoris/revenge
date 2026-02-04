# Phase A5: Library Matcher

**Priority**: P2
**Effort**: 4-6h
**Dependencies**: A0, A1

---

## A5.1: Library File Matching

**Affected Files**:
- `internal/content/movie/library_matcher.go:120`
- `internal/content/movie/library_service.go:190`

**Current State**:
```go
return nil, fmt.Errorf("not implemented")
// For now, return placeholder
```

**Tasks**:
- [ ] Implement file → movie matching algorithm
- [ ] Parse filename for title, year, quality
- [ ] Query TMDb for matches
- [ ] Score and rank matches
- [ ] Return best match or candidates
- [ ] Tests

---

## Matching Algorithm

### Step 1: Filename Parsing
Extract from filename:
- Title (cleaned)
- Year (if present)
- Quality indicators (1080p, 4K, etc.)
- Release group (optional)

Example: `The.Matrix.1999.1080p.BluRay.x264-GROUP.mkv`
→ Title: "The Matrix", Year: 1999, Quality: 1080p

### Step 2: TMDb Lookup
- Search TMDb by title
- Filter by year if available
- Get multiple candidates

### Step 3: Confidence Scoring
Score based on:
- Title match (Levenshtein distance)
- Year match (exact = 100%, ±1 year = 80%)
- Existing file in library (avoid duplicates)

### Step 4: Result
- If single high-confidence match (>90%): auto-match
- If multiple candidates: return for manual selection
- If no match: flag for manual lookup
