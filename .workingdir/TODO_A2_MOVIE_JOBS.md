# Phase A2: Movie Jobs Completion

**Priority**: P2
**Effort**: 4-6h
**Dependencies**: A0, A1

---

## A2.1: File Match Job

**Affected File**: `internal/content/movie/moviejobs/file_match.go:58-64`

**Tasks**:
- [ ] Implement `library.Service.MatchFile` method
- [ ] Update file match worker to use it
- [ ] Match logic: filename parsing → TMDb lookup → confidence scoring
- [ ] Tests

**Current Error**:
```
movie file match not implemented: library.Service.MatchFile method not available
```

---

## A2.2: Metadata Refresh Job

**Affected Files**:
- `internal/content/movie/service.go:275-277`
- `internal/content/movie/moviejobs/metadata_refresh.go:90-92`

**Tasks**:
- [ ] Implement `RefreshMetadata(ctx, movieID)` in service
- [ ] Queue River job for async processing
- [ ] Refresh credits during metadata refresh
- [ ] Tests

**Current Error**:
```go
return fmt.Errorf("metadata refresh not implemented yet")
```
