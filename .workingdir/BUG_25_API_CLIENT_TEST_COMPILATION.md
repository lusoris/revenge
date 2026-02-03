# BUG #25: API Client Tests Compilation Errors

## Status
**FOUND** - 2025-02-03

## Severity
High - Tests won't compile

## Location
- File: `tests/integration/api_client_test.go`
- Multiple test functions

## Description
API client tests failing to compile due to ogen interface changes. `ogen.WithClient()` returns `ClientOption` but `ogen.NewClient()` expects `SecuritySource`.

## Errors
```
cannot use ogen.WithClient(ts.HTTPClient) (value of interface type ogen.ClientOption)
as ogen.SecuritySource value in argument to ogen.NewClient:
ogen.ClientOption does not implement ogen.SecuritySource (missing method BearerAuth)
```

## Root Cause
API interface mismatch in generated ogen client code.

## Fix Required
Update API client test code to match current ogen-generated API.
