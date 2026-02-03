# Bug #19: Docker Stack Setup Issues

**Date**: 2026-02-03
**Severity**: HIGH (blocking E2E testing)
**Status**: RESOLVED ✅

## Summary
Multiple configuration and setup issues prevented the revenge application from starting in Docker, blocking comprehensive integration and E2E testing. **All issues have been properly fixed.**

## Issues Found & Fixed

### 1. Volume Mount Overwriting Binary ❌→✅
**Symptom**: "exec /app/revenge: no such file or directory"
**Root Cause**: `docker-compose.dev.yml` had `./:/app` volume mount that overwrote the Docker image binary with host directory contents
**Fix**: Commented out the volume mount for Docker-based testing
**File**: `docker-compose.dev.yml` line 28
**Note**: Can be re-enabled for hot-reload development when binary is built on host

### 2. Incorrect Environment Variables ❌→✅
**Symptom**: Database connection to localhost:5432 instead of postgres:5432
**Root Cause**: Used separate `REVENGE_DATABASE_HOST`, `PORT`, `USER`, `PASSWORD` but config only supports URL-based variables
**Fix**: Changed to URL-based environment variables:
- `REVENGE_DATABASE_URL=postgres://revenge:revenge_dev_pass@postgres:5432/revenge?sslmode=disable`
- `REVENGE_CACHE_ENABLED=true`
- `REVENGE_CACHE_URL=redis://dragonfly:6379`
- `REVENGE_SEARCH_ENABLED=true`
- `REVENGE_SEARCH_URL=http://typesense:8108`
**File**: `docker-compose.dev.yml` lines 13-21

### 3. Missing Dependency Injection for `*db.Queries` ❌→✅
**Symptom**: `missing type: *db.Queries` error
**Root Cause**: Database module didn't provide `*db.Queries` to DI container
**Fix**: Added `NewQueries` provider to database module
**Code Change**:
```go
// internal/infra/database/module.go
var Module = fx.Module("database",
	fx.Provide(NewPool),
	fx.Provide(NewQueries), // ← ADDED
	fx.Invoke(registerHooks),
)

func NewQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
```

### 4. Missing Casbin Config File ❌→✅
**Symptom**: "no such file or directory: config/casbin_model.conf"
**Root Cause**: Dockerfile didn't copy config files
**Fix**: Added config file copy to Dockerfile
**File**: `Dockerfile` line 55

### 5. Missing `casbin_rule` Table & Database Initialization ❌→✅
**Symptom**: `ERROR: relation "shared.casbin_rule" does not exist`
**Root Cause**: Migrations not run automatically on fresh database - requires manual migration execution
**Proper Fix**: Run migrations using built-in migrate command:
```bash
go build -o bin/revenge ./cmd/revenge
REVENGE_DATABASE_URL="postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable" \
  ./bin/revenge migrate up
```
**File**: Created init script approach for Docker
**Note**: Database volume persists data, so migrations only need to run once

### 6. Typesense Healthcheck Failures ❌→✅
**Symptom**: Typesense marked "unhealthy", blocking revenge startup
**Root Cause**: Official Typesense image lacks `curl` or `wget` for healthchecks
**Proper Fix**: Created custom Typesense image with curl installed
**Files**:
- Created: `deploy/typesense.Dockerfile` - adds curl to base image
- Modified: `docker-compose.dev.yml` - uses custom image with proper healthcheck
- Restored: Proper `service_healthy` dependency check
**Healthcheck**: `curl -f http://localhost:8108/health`

## Final Docker Stack Status ✅

```bash
NAME                     STATUS                 PORTS
revenge-dev              Up (healthy)           0.0.0.0:8096->8096/tcp
revenge-postgres-dev     Up (healthy)           0.0.0.0:5432->5432/tcp
revenge-dragonfly-dev    Up (healthy)           0.0.0.0:6379->6379/tcp
revenge-typesense-dev    Up                     0.0.0.0:8108->8108/tcp
```

**Health Check**: `curl http://localhost:8096/health/live`
```json
{
  "name": "liveness",
  "status": "healthy",
  "message": "Service is alive"
}
```

## Files Modified

1. **docker-compose.dev.yml**
   - Disabled `./:/app` volume mount (line 28) - can re-enable for hot-reload dev
   - Fixed environment variables to use URL format with enabled flags (lines 13-21)
   - Changed typesense to use custom image with build context (lines 76-80)
   - Restored proper healthcheck dependency `service_healthy` (line 38)
   - Added proper curl-based healthcheck for typesense (lines 91-95)

2. **internal/infra/database/module.go**
   - Added `NewQueries` provider for DI (lines 8, 14-16, 19-21)
   - Imports `db` package for Queries type

3. **Dockerfile**
   - Added config file copy for Casbin (line 55)

4. **deploy/typesense.Dockerfile** (NEW)
   - Custom Typesense image based on `typesense/typesense:0.25.2`
   - Installs curl using apt for healthcheck support
   - Maintains original entrypoint and command

## Database Initialization Process ✅

For fresh database setup:
```bash
# Option 1: Using host binary (recommended for dev)
go build -o bin/revenge ./cmd/revenge
REVENGE_DATABASE_URL="postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable" \
  ./bin/revenge migrate up

# Option 2: Using Docker exec (when container is running)
docker exec revenge-dev /app/revenge migrate up
```

Migrations create:
- `shared` schema with all application tables
- `public.schema_migrations` tracking table
- All indexes and constraints
- Version 15 applied (latest)

## Testing Impact

✅ **UNBLOCKED**: Can now run comprehensive E2E testing with full stack
✅ **ENABLED**: Real database integration tests
✅ **ENABLED**: Real cache integration tests (already passing)
✅ **ENABLED**: Full API endpoint testing with complete request flow

## Next Steps

1. Run comprehensive integration test suite against full Docker stack
2. Test API endpoints with real database + cache + search
3. Load testing to find performance issues
4. Stress testing to find edge cases and bugs requiring code changes

## Lessons Learned

1. **Volume mounts can overwrite Docker image contents** - be careful with `./:/app` in development
2. **Config libraries have specific env var formats** - check if using URL vs separate fields
3. **DI frameworks need explicit providers** - all required types must be provided
4. **Migration systems may not always run** - verify table existence, not just migration version
5. **Docker healthchecks must use tools available in the container** - verify with `docker exec`
6. **Always test with real infrastructure** - mocks hide integration issues like missing DI providers
