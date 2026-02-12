# BUG-007: Notification system fully built but never called

**Severity:** HIGH  
**Category:** Flow Integration  
**Status:** RESOLVED (scan start/done events wired in movie + tvshow workers)

## Description

The entire notification stack is wired and functional:
- `Dispatcher` → `NotificationWorker` → Discord/Gotify/Ntfy/Webhook/Email agents
- Event types defined in `internal/service/notification/notification.go`
- `NotificationWorker` registered via `internal/infra/jobs/module.go`

BUT: **No production code anywhere enqueues `NotificationArgs` jobs.** The services that should emit events never call `Dispatch()` or insert notification jobs.

## Missing Dispatch Points

| Location | Event Type | Should Fire |
|----------|-----------|-------------|
| Library scan completion (movie/tvshow workers) | `EventLibraryScanDone` | After successful scan |
| Library scan start | `EventLibraryScanStarted` | Before scan begins |
| Movie added (during scan) | `EventMovieAdded` | When new movie is discovered |
| Auth login success | `EventLoginSuccess` | After successful login |
| Auth login failure | `EventLoginFailed` | After failed login |
| User created | `EventUserCreated` | After user registration |
| Playback started | `EventPlaybackStarted` | When playback session starts |
| System startup | `EventSystemStartup` | On app start |

## Fix

Add notification dispatch calls at key lifecycle points through the existing `infrajobs.Client.Insert()` using `NotificationArgs`.
