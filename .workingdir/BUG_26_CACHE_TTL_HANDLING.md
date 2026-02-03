# BUG #26: Cache TTL Tests Failing

## Status
**FOUND** - 2025-02-03

## Severity
Medium - Test failures reveal TTL handling issues

## Location
- File: `tests/integration/cache/cache_advanced_test.go`
- Test: `TestTTLAccuracy`

## Description
TTL (Time-To-Live) tests failing due to invalid expire time errors from Dragonfly.

## Errors
```
TestTTLAccuracy/500ms:
  L2 cache set failed: invalid expire time in 'set' command

TestTTLAccuracy/1s:
  Key should have expired (but didn't)

TestTTLAccuracy/2s:
  Key should have expired (but didn't)
```

## Root Cause
1. Sub-second TTLs (500ms) are being rejected by Dragonfly
2. TTL expiration timing is not accurate - keys persisting longer than expected

## Fix Required
1. Adjust TTL handling to use proper format for Dragonfly
2. Fix timing issues in expiration tests
