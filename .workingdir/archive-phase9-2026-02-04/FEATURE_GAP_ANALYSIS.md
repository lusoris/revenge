# Feature Gap Analysis Report

**Created**: 2026-02-04
**Scope**: v0.3.0 MVP + Future Versions
**Comparison**: Jellyfin, Plex, Overseerr, Tautulli, Navidrome, Audiobookshelf, Kavita, Immich

---

## Executive Summary

Nach umfassender Analyse des aktuellen Implementierungsstands wurden folgende kritische Lücken identifiziert:

| Kategorie | Kritisch | Wichtig | Nice-to-Have |
|-----------|----------|---------|--------------|
| Library Scanning | 2 | 1 | 1 |
| RBAC/Permissions | 2 | 2 | 1 |
| Notifications | 1 | 2 | 3 |
| Cluster/Infrastructure | 0 | 3 | 2 |
| Caching | 0 | 2 | 2 |
| Security | 1 | 2 | 1 |
| Profiling/Monitoring | 0 | 2 | 2 |
| **Total** | **6** | **14** | **12** |

---

## 1. Library Scanning Gaps

### 1.1 ❌ KRITISCH: FFprobe Integration fehlt → **go-astiav stattdessen**

**Aktueller Stand**:
- `ExtractFileInfo()` in `library_scanner.go` ist ein Stub
- Nur `file.Size()` und Container-Extension werden extrahiert
- Schema (`movie_files`) hat alle Felder, aber sie bleiben leer

**Entscheidung**: **go-astiav** statt FFprobe CLI verwenden!

**Warum go-astiav?**:
- Bereits für Transcoding (v0.6.0) geplant → keine zusätzliche Dependency
- Native Go Bindings → keine Exec-Calls, keine JSON-Parsing
- Typed API → sicherer
- CGO erforderlich, aber das brauchen wir eh für HW-Acceleration

**go-astiav liefert**:
```go
FormatContext:
  - Duration()     // Gesamtdauer in µs
  - BitRate()      // Gesamt-Bitrate
  - Metadata()     // Title, Artist, etc.
  - Streams()      // Alle Streams

Stream.CodecParameters():
  - CodecID()      // H264, HEVC, AAC, etc.
  - BitRate()      // Stream-Bitrate
  - Width/Height() // Resolution
  - PixelFormat()  // YUV420P, etc.
  - Profile()      // H264 Main, High, etc.
  - ColorSpace()   // BT709, BT2020
  - ColorRange()   // Full/Limited

Stream:
  - AvgFrameRate() // 23.976, 24, 25, etc.
  - Duration()     // Stream-Dauer
  - Metadata()     // language, title, etc.
```

**Workaround aktuell**: Radarr-Sync füllt MediaInfo via Radarr API

**Lösung**: `internal/content/movie/mediainfo.go` mit go-astiav

### 1.2 ❌ KRITISCH: Kein Real-time File Watching

**Aktueller Stand**:
- `github.com/fsnotify/fsnotify` ist in `go.mod` (Dependency)
- ABER: Nirgends im Code verwendet!
- Library-Scan nur manuell oder via River Job

**Workaround**: Radarr Webhooks decken Import-Events ab

**Fehlende Funktionalität**:
- Automatische Erkennung neuer Dateien
- Automatische Erkennung gelöschter Dateien
- Automatische Erkennung umbenannter Dateien
- Trigger für Re-Scan bei Änderungen

**Priorität**: v0.4.0 (nicht MVP-kritisch wegen Radarr)

### 1.3 ⚠️ WICHTIG: Keine Duplikat-Erkennung

**Aktueller Stand**: Keine Implementierung

**Fehlend**:
- Perceptual Hashing für Video-Duplikate
- Fingerprinting für Audio-Duplikate
- UI für Duplikat-Management

### 1.4 💡 NICE-TO-HAVE: Keine NFO/Sidecar Support

**Fehlend**:
- Jellyfin/Kodi NFO file parsing
- External subtitle file detection (.srt, .ass, .sub)
- External artwork detection (poster.jpg, fanart.jpg)

---

## 2. RBAC/Permissions Gaps

### 2.1 ❌ KRITISCH: Keine Custom Groups/Roles

**Aktueller Stand**:
```conf
# casbin_model.conf - Nur basic RBAC
[role_definition]
g = _, _  # user, role
```

```sql
-- migrations/000011 - Nur 2 hardcoded Rollen
INSERT INTO shared.casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'admin', '*', '*'),      -- Admin: Alles
('p', 'user', 'movies', 'read'); -- User: Movies lesen
```

**Problem**:
- Admins können KEINE eigenen Gruppen erstellen
- Keine "Moderator" Rolle vordefiniert
- Keine granularen Permissions pro Library/Content

**Fehlende Features**:
- [ ] Create/Edit/Delete custom roles via API
- [ ] Role hierarchy (admin > moderator > user > guest)
- [ ] Per-library permissions
- [ ] Per-user restrictions (age ratings, content types)

### 2.2 ❌ KRITISCH: Keine dynamischen Settings per Role

**Aktueller Stand**:
- `server_settings` - Server-weite Settings
- `user_settings` - User-spezifische Settings
- ABER: Keine Role-based Settings!

**Fehlend**:
- Moderator-Einstellungen (z.B. Request-Approval-Rechte)
- Role-spezifische Limits (Transcode-Limits, Download-Limits)
- Role-spezifische Features (Skip-Intro aktiviert?)

### 2.3 ⚠️ WICHTIG: Keine vordefinierte Moderator-Rolle

**Best Practice** (Jellyfin/Plex):
- Admin: Full control
- Moderator: Manage users, approve requests, moderate content
- User: Normal access
- Guest: Limited access

**Aktuell nur**: admin, user

### 2.4 ⚠️ WICHTIG: Keine Library-Level Permissions

**Fehlend**:
- User X kann nur Library A sehen
- User Y kann Library B+C sehen
- Kindersicherung per Library

---

## 3. Notification Gaps

### 3.1 ❌ KRITISCH: Notification Service fehlt komplett

**Aktueller Stand**:
```sql
-- migrations/000006 - Schema existiert
email_notifications JSONB DEFAULT '{"enabled": true, ...}'
push_notifications JSONB DEFAULT '{"enabled": false}'
digest_notifications JSONB DEFAULT '{"enabled": true, ...}'
```

**ABER**: Keine Implementation! Nur DB-Schema.

**Fehlende Komponenten**:
- [ ] Notification Service Interface
- [ ] Notification Event Types (enum)
- [ ] Notification Queue (River Job)
- [ ] Notification Agents (Discord, Webhook, Email, etc.)
- [ ] User Notification Preferences API

### 3.2 ⚠️ WICHTIG: Keine Notification Agents

**Priorität für v0.3.0**:
1. Webhook Agent (generisch)
2. Discord Agent (populär)

**Später**:
- Telegram
- Email (SMTP)
- Gotify/ntfy
- Apprise Integration (100+ Services)
- Push Notifications (PWA)

### 3.3 ⚠️ WICHTIG: Keine Event-Typen definiert

**Typische Events** (Tautulli/Overseerr):
- `movie.added` - Neuer Film in Library
- `movie.available` - Angeforderter Film verfügbar
- `request.created` - Neue Anfrage
- `request.approved` - Anfrage genehmigt
- `user.created` - Neuer User
- `playback.started` - Wiedergabe gestartet
- `playback.paused` - Wiedergabe pausiert
- `transcode.started` - Transcoding gestartet

---

## 4. Cluster/Infrastructure Gaps

### 4.1 ⚠️ WICHTIG: Keine Session Sharing für Cluster

**Aktueller Stand**:
- Sessions in PostgreSQL gespeichert ✅
- Rate Limiter: In-Memory `sync.Map` ❌

**Problem**: Rate Limiter funktioniert nicht in Multi-Instance

**Lösung**: Rate Limiter State in Dragonfly/Redis speichern

### 4.2 ⚠️ WICHTIG: Keine Graceful Shutdown Propagation

**Aktueller Stand**:
- Basic graceful shutdown vorhanden
- Keine Readiness-Probe während Shutdown

**Fehlend**:
- Pre-stop hook für K8s (z.B. 5s delay)
- Drain connections vor Shutdown
- Health-Endpoints aktualisieren auf "not ready"

### 4.3 ⚠️ WICHTIG: Keine Horizontal Scaling Vorbereitung

**Potenzielle Probleme**:
- River Queue: Shared via PostgreSQL ✅ (OK)
- Casbin: Adapter-based ✅ (OK mit pgx-adapter)
- Cron Jobs: Könnten parallel laufen ❌
- WebSockets: Keine Redis PubSub ❌

### 4.4 💡 NICE-TO-HAVE: Kein Leader Election

**Für Singleton-Tasks** (Scheduled Jobs):
- Aktuell: River handles das
- Aber: Eigene Cron-ähnliche Tasks brauchen Leader Election

### 4.5 💡 NICE-TO-HAVE: Keine Helm Chart Values vollständig

**Aktuell**: Basic Helm Chart existiert in `charts/revenge/`

**Fehlend**:
- HPA (Horizontal Pod Autoscaler) config
- PDB (Pod Disruption Budget) config
- Service Mesh annotations (Istio/Linkerd)
- External Secrets Operator integration

---

## 5. Caching Gaps

### 5.1 ⚠️ WICHTIG: Cache Strategy nicht dokumentiert

**Aktueller Stand**:
- L1: Otter (In-Memory)
- L2: Rueidis → Dragonfly

**Gut implementiert** ✅:
- TMDb Client: Caching mit sync.Map + TTL
- Session: DB-basiert (kein Cache nötig)

**Nicht gecached** (sollte es aber):
- [ ] Movie Details (häufig abgefragt)
- [ ] User Permissions/Roles (bei jedem Request geprüft)
- [ ] Library Stats (Continue Watching, Recently Added)
- [ ] Search Results (Typesense-Responses)

### 5.2 ⚠️ WICHTIG: Cache Invalidation nicht konsistent

**Problem**: Wenn Movie in DB aktualisiert wird, bleibt Cache stale

**Fehlend**:
- Cache-Tags für gruppenweise Invalidierung
- Event-basierte Invalidierung (via River?)
- Cache-Warming nach Invalidierung

### 5.3 💡 NICE-TO-HAVE: Kein Cache Metrics

**Fehlend**:
- Hit/Miss Ratio tracking
- Cache Size monitoring
- TTL distribution

### 5.4 💡 NICE-TO-HAVE: Kein sturdyc Integration

**Design referenziert sturdyc** aber nicht implementiert:
- Stampede protection
- Negative caching
- Background refresh

---

## 6. Security Gaps

### 6.1 ❌ KRITISCH: Kein Audit Logging implementiert

**Aktueller Stand**:
```sql
-- migrations/000014 - Schema existiert
CREATE TABLE shared.activity_log (
    id UUID PRIMARY KEY,
    user_id UUID,
    action VARCHAR(100),
    resource_type VARCHAR(100),
    resource_id UUID,
    ...
);
```

**ABER**: Keine Logging-Calls im Code!

**Fehlend**:
- [ ] Activity Logger Service
- [ ] Automatische Logs für Auth-Events
- [ ] Automatische Logs für Admin-Actions
- [ ] Automatische Logs für Content-Changes
- [ ] API zum Abrufen von Logs

### 6.2 ⚠️ WICHTIG: Keine IP-basierte Zugriffskontrolle

**Fehlend**:
- IP Allowlist/Blocklist
- Geo-blocking
- Admin-only IPs

### 6.3 ⚠️ WICHTIG: Kein Account Lockout nach Failed Attempts

**Aktueller Stand**: Rate Limiting existiert, aber kein permanentes Lockout

**Best Practice**:
- Nach 5 fehlerhaften Logins: 15min Lockout
- Nach 10: 1h Lockout
- Nach 20: Account locked, Admin-Unlock erforderlich

### 6.4 💡 NICE-TO-HAVE: Keine Session Device Management

**Fehlend**:
- "Auf allen Geräten abmelden"
- Session Device Info (Browser, OS, IP)
- Trusted Devices für MFA-Bypass

---

## 7. Profiling/Monitoring Gaps

### 7.1 ⚠️ WICHTIG: Kein pprof Endpoint

**Aktueller Stand**: Nicht vorhanden

**Empfohlen**:
```go
import _ "net/http/pprof"
// Endpoint: /debug/pprof/*
```

**Nur in Dev/Debug Mode aktivieren!**

### 7.2 ⚠️ WICHTIG: Keine Prometheus Metrics

**Aktueller Stand**:
- OpenTelemetry Traces: Dependency vorhanden (`go.opentelemetry.io/otel`)
- ABER: Keine Metriken-Endpoints!

**Fehlend**:
- `/metrics` endpoint (Prometheus format)
- HTTP Request Latency histograms
- Database query latency
- Cache hit/miss ratios
- Active sessions count
- Concurrent transcodes

### 7.3 💡 NICE-TO-HAVE: Keine Benchmark Tests für kritische Pfade

**Vorhanden** ✅:
- `password_bench_test.go` - Password hashing benchmarks

**Fehlend**:
- Database query benchmarks
- Cache operation benchmarks
- API endpoint benchmarks
- Serialization benchmarks

### 7.4 💡 NICE-TO-HAVE: Kein Distributed Tracing aktiviert

**Aktuell**: Dependency da, nicht konfiguriert

**Fehlend**:
- Jaeger/Tempo integration
- Trace ID propagation
- Span creation für wichtige Operations

---

## 8. Additional Missing Features

### 8.1 Request System (Overseerr-Alternative)

**Roadmap**: v0.8.0 (Intelligence)
**Empfehlung**: Vorziehen auf v0.4.0

**Scope**:
- Request table + SQLC
- Request service
- Radarr/Sonarr integration
- Approval workflow
- Notifications

### 8.2 Transcoding/Playback (v0.6.0)

**Design existiert**, keine Implementierung:
- HLS manifest generation (gohlslib)
- FFmpeg transcoding
- Hardware acceleration (VAAPI, NVENC, QSV)
- Quality selection
- Subtitle burning

### 8.3 Analytics Dashboard

**Fehlend**:
- Watch statistics per user
- Most watched content
- Concurrent streams
- Bandwidth usage
- Library growth over time

---

## Prioritäts-Matrix für v0.3.0

### ALLES VOR FRONTEND IMPLEMENTIEREN:

| # | Item | Aufwand | Status |
|---|------|---------|--------|
| 1 | MediaInfo mit go-astiav | 4-6h | ⬜ |
| 2 | Notification Service (alle 4 Agents) | 6-8h | ⬜ |
| 3 | Audit Logging aktivieren | 3-4h | ⬜ |
| 4 | RBAC: alle 4 Rollen + Custom Role API | 6-8h | ⬜ |
| 5 | Rate Limiter zu Dragonfly | 2-3h | ⬜ |
| 6 | Cache für Hot Paths | 4-6h | ⬜ |
| 7 | Observability (pprof + Prometheus) | 4-5h | ⬜ |
| 8 | Test Coverage 80%+ | 8-12h | ⬜ |
| **Total** | | **37-52h** | |

### CAN DEFER (nach v0.3.0):

| # | Item | Target Version |
|---|------|----------------|
| Fsnotify File Watching | v0.4.0 |
| Request System | v0.4.0 |
| Transcoding | v0.6.0 |
| Prometheus Metrics | v0.4.0 |
| Distributed Tracing | v0.5.0 |
| pprof Endpoint | v0.4.0 |

---

## Architektur-Konformität

### ✅ Cluster-Ready

- **PostgreSQL**: Shared state ✅
- **Dragonfly**: Shared cache ✅
- **River**: Distributed job queue ✅
- **Casbin**: DB-backed policies ✅
- **Stateless API**: Keine lokalen Files ✅

### ⚠️ Benötigt Anpassung für Multi-Instance

- **Rate Limiter**: In-Memory → Redis/Dragonfly
- **WebSockets**: Keine PubSub (für SyncPlay später relevant)
- **Scheduled Tasks**: Leader Election oder River-basiert

### ✅ Kubernetes-Ready

- Health Endpoints: `/healthz`, `/readyz`, `/startupz` ✅
- Graceful Shutdown: Basic vorhanden ✅
- Environment Config: Koanf unterstützt ENV vars ✅
- Helm Chart: Vorhanden ✅

---

## Fazit

Der aktuelle Stand ist **solide für ein MVP**, aber es wurden **8 Phasen** identifiziert die vor dem Frontend-Start abgeschlossen werden müssen:

1. **MediaInfo (go-astiav)** - Ohne das keine echten MediaInfo bei Self-Scan
2. **Notification Service** - Alle 4 Agents: Webhook, Discord, Email, Gotify
3. **Audit Logging** - Schema existiert, muss aktiviert werden
4. **RBAC Erweiterungen** - 4 Rollen + Custom Role API
5. **Rate Limiter Migration** - Dragonfly für Multi-Instance
6. **Cache Hot Paths** - Performance-kritische Endpoints
7. **Observability** - pprof (Dev) + Prometheus Metrics
8. **Test Coverage 80%** - Quality Gate

**Geschätzter Aufwand**: ~45-55 Stunden (~1.5 Wochen Vollzeit)

**Entscheidung**: ALLES vor Frontend implementieren, um spätere Loops zu vermeiden.
