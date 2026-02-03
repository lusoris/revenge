# User Service Bugs

## Bug 1: ListUsers filters by default values instead of ignoring nil filters

**Location**: `internal/service/user/repository_pg.go:60-90`

**Problem**:
Wenn `UserFilters.IsActive` oder `UserFilters.IsAdmin` nil sind, werden sie als `false` behandelt und dann zum Filtern verwendet:

```go
var isActive, isAdmin bool  // ❌ Defaults to false!
if filters.IsActive != nil {
    isActive = *filters.IsActive
}
if filters.IsAdmin != nil {
    isAdmin = *filters.IsAdmin
}

count, err := r.queries.CountUsers(ctx, db.CountUsersParams{
    Column1: isActive,  // ❌ Always filters by is_active=false when nil
    Column2: isAdmin,   // ❌ Always filters by is_admin=false when nil
})
```

**Expected**:
Wenn Filter nil sind, sollten ALLE Users zurückgegeben werden (kein Filter)

**Actual**:
Wenn Filter nil sind, werden nur Users mit `is_active=false AND is_admin=false` zurückgegeben

**Impact**:
- ListUsers ohne Filter gibt nur inaktive, nicht-admin Users zurück
- API calls die alle Users wollen funktionieren nicht
- Tests mussten workaround mit allen Users als inactive/non-admin

**Fix Options**:

### Option 1: SQL Query mit optionalen Filtern umschreiben
Braucht Änderungen in sqlc queries um `WHERE (is_active = $1 OR $1 IS NULL)` Pattern zu nutzen

### Option 2: Unterschiedliche Queries je nach Filter
```go
if filters.IsActive == nil && filters.IsAdmin == nil {
    users, err = r.queries.ListAllUsers(ctx, ...)
} else if filters.IsActive != nil && filters.IsAdmin != nil {
    users, err = r.queries.ListUsersWithBothFilters(ctx, ...)
} else ...
```

### Option 3: Repository-Level filtering (nicht empfohlen)
Load all, filter in Go - schlecht für Performance

**Recommended Fix**: Option 1 - SQL Query Änderungen
