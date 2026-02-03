# BUG #23: Type Assertion Errors in Search Integration Tests

## Status
**FOUND** - 2025-02-03

## Severity
Medium - Test failures due to type mismatches

## Location
- File: `tests/integration/search/search_test.go`
- Multiple test functions

## Description
Integration tests failing because of type comparison errors between `int32` and `int`.

## Errors
```
TestSearchDocumentOperations:
  Error: Elements should be the same type (comparing int32 with int)

TestSearchBulkImport:
  Error: Not equal: expected int32(3), actual int(3)

TestSearchWithFiltersAndSorting:
  Error: Not equal: expected int32(2), actual int(2)
```

## Root Cause
Typesense API returns `*int32` for `Found` field, but test assertions are comparing with plain `int` literals.

## Fix
Convert test expectations to `int32` type for all assertions comparing with `results.Found`.

## Resolution
Fixed by changing all comparison values from `int` to `int32` in test assertions.
