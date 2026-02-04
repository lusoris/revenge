# Phase A3: Test Infrastructure

**Priority**: P2
**Effort**: 3-4h
**Dependencies**: A0

---

## A3.1: Dragonfly Testcontainer

**Affected File**: `internal/testutil/containers.go:168-171`

**Current State**:
```go
t.Skip("Dragonfly container not yet implemented - implement when cache module is needed")
```

**Tasks**:
- [ ] Implement `NewDragonflyContainer(t) (*DragonflyContainer, error)`
- [ ] Use `docker.io/dragonflydb/dragonfly:latest` image
- [ ] Return connection string
- [ ] Cleanup on test completion
- [ ] Test the container helper itself

---

## A3.2: Typesense Testcontainer

**Affected File**: `internal/testutil/containers.go:182-185`

**Current State**:
```go
t.Skip("Typesense container not yet implemented - implement when search module is needed")
```

**Tasks**:
- [ ] Implement `NewTypesenseContainer(t) (*TypesenseContainer, error)`
- [ ] Use `typesense/typesense:latest` image
- [ ] Configure API key
- [ ] Return connection details
- [ ] Cleanup on test completion
- [ ] Test the container helper itself

---

## Notes

Both containers should follow the pattern established by the PostgreSQL testcontainer in the same file.
