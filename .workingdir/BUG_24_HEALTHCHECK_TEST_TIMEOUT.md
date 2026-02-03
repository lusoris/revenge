# BUG #24: Health Check Test Timeout Too Short

## Status
**FOUND** - 2025-02-03

## Severity
Low - Test configuration issue

## Location
- File: `tests/integration/search/search_test.go`
- Test: `TestSearchClientConnection`

## Description
Health check test context timeout (10s) expires before the Typesense client can complete the health check.

## Error
```
Error: Get "http://localhost:8108/health": context deadline exceeded
```

## Root Cause
Test context has 10-second timeout but the health check itself has a 10-second timeout parameter. These overlap and cause the context to expire before the health check completes.

## Fix
Increase test context timeout to 15 seconds to allow the 10-second health check to complete.

## Resolution
Fixed by increasing context timeout from 10s to 15s in TestSearchClientConnection.
