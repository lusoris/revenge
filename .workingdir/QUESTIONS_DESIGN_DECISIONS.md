# Design-Entscheidungen & Offene Fragen

**Erstellt**: 2026-02-05
**Kontext**: Fokus auf Clusterfähigkeit, Security-Fixes, und Deduplication ohne Monolith

---

## 1. Clusterfähigkeit

### 1.1 Media-File-Storage

**Problem**: Media-Dateien sind aktuell auf lokalem Filesystem, nicht cluster-tauglich.

**Optionen**:
- **Option A: NFS/CephFS/SMB** (ReadOnlyMany Kubernetes Volume)
  - ✅ Native Filesystem-Performance
  - ✅ Keine Code-Änderungen nötig
  - ❌ Externe NFS-Server-Infrastruktur nötig

- **Option B: Object Storage** (S3/MinIO)
  - ✅ Cloud-native, hochverfügbar
  - ✅ Keine NFS-Infrastruktur
  - ❌ Großer Refactoring-Aufwand (Scanner umschreiben)
  - ❌ Schlechtere Performance für Streaming

- **Option C: Hybrid** (NFS für Media, S3 für Avatars)
  - ✅ Beste Performance wo nötig
  - ❌ Zwei Storage-Systeme zu verwalten

**Frage**: Welche Storage-Lösung bevorzugst du für Media-Dateien?

---

### 1.2 Avatar/User-Content-Storage

**Problem**: Avatar-Uploads aktuell auf lokalem Filesystem.

**Optionen**:
- **Option A: S3/MinIO** (Empfohlen)
  - ✅ Cloud-native, einfach skalierbar
  - ✅ Interface bereits vorbereitet
  - ✅ 8-16h Implementierungsaufwand

- **Option B: Shared Volume** (NFS mit ReadWriteMany)
  - ✅ Einfacher wenn NFS schon da ist
  - ❌ RWX-Volumes komplexer als RO

**Frage**: S3/MinIO für Avatars implementieren?

---

### 1.3 Leader Election für Periodic Jobs

**Problem**: Cleanup-Jobs würden auf jedem Pod laufen (Mehrfachausführung).

**Optionen**:
- **Option A: River Periodic Jobs** (Empfohlen)
  - ✅ Native River-Feature
  - ✅ Einfache Integration
  - ✅ PostgreSQL-basiert (kein extra Service)
  - ✅ 4-8h Aufwand

- **Option B: HashiCorp Raft**
  - ✅ Bereits im go.mod
  - ❌ Komplexer Setup
  - ❌ Overkill für unseren Use-Case

- **Option C: Kubernetes Lease API**
  - ✅ Native Kubernetes-Feature
  - ❌ Nur für K8s, nicht für Docker Swarm

- **Option D: PostgreSQL Advisory Locks**
  - ✅ Einfach, kein extra Service
  - ❌ Manuell implementieren

**Frage**: River Periodic Jobs für Leader Election verwenden?

---

## 2. TV Shows Implementation

### 2.1 Metadata Provider

**Kontext**: TMDb unterstützt sowohl Filme als auch TV-Serien!

**Optionen**:
- **Option A: TMDb für TV Shows nutzen** (Empfohlen)
  - ✅ TMDb-Client bereits vorhanden
  - ✅ Einheitliche API
  - ✅ Weniger Dependencies
  - ✅ Schnellere Implementierung
  - ❌ Möglicherweise weniger TV-spezifische Features als TheTVDB

- **Option B: TheTVDB zusätzlich**
  - ✅ Spezialisiert auf TV-Daten
  - ✅ Community-Daten
  - ❌ Neuer Client nötig
  - ❌ Zwei Systeme zu pflegen

- **Option C: Beide (TMDb primär, TheTVDB fallback)**
  - ✅ Beste Datenqualität
  - ❌ Komplexere Logik

**Frage**: Nur TMDb für TV Shows verwenden oder zusätzlich TheTVDB?

---

### 2.2 TV Shows Timing

**Frage**: Wann soll TV-Implementierung starten?
- Direkt nach Cluster + Security-Fixes?
- Oder erst nach vollständiger Deduplication/Refactoring?

---

## 3. Shared Abstractions & Deduplication

### 3.1 Tabellen-Architektur

**Prinzip**: Module eigenständig, Tabellen NICHT monolithisch!

**Mögliche Ansätze**:

**Option A: Separate Schemas pro Content-Typ** (Empfohlen)
```sql
-- Movies
public.movies
public.movie_files
public.movie_genres
public.movie_credits

-- TV Shows
tvshow.series
tvshow.seasons
tvshow.episodes
tvshow.episode_files

-- Music
music.artists
music.albums
music.tracks
music.track_files
```
✅ Klare Trennung
✅ Unabhängige Migrationen
✅ Keine Kollisionen

**Option B: Prefix-basiert in public Schema**
```sql
public.movies
public.movie_files
public.tvshows
public.tvshow_episodes
public.music_albums
public.music_tracks
```
✅ Einfacher
❌ Weniger Isolation

**Option C: Separate Datenbanken**
```
revenge_movies
revenge_tvshows
revenge_music
```
❌ Overkill
❌ Cross-Module-Queries unmöglich

**Frage**: Separate Schemas (Option A) oder Prefix-basiert (Option B)?

---

### 3.2 Shared Code vs. Eigenständigkeit

**Zu teilen** (OK, kein Monolith):
- ✅ Algorithmen (Levenshtein, Fuzzy Matching)
- ✅ File-Scanner-Framework (mit Adaptern)
- ✅ HTTP-Provider-Framework (mit Adaptern)
- ✅ Background-Job-Boilerplate

**NICHT zu teilen** (würde Monolith erzeugen):
- ❌ Domain-Models (Movie, TVShow, Album bleiben getrennt)
- ❌ Repositories (jedes Modul hat eigenes)
- ❌ Database-Tabellen
- ❌ API-Endpoints

**Frage**: Ist diese Trennung OK für dich?

---

### 3.3 Refactoring-Timing

**Option A: Jetzt sofort refactoren**
- Shared Abstractions extrahieren
- Movie-Modul umbauen
- Dann TV implementieren
- ⏱️ 4-6 Wochen vor TV-Start

**Option B: Schrittweise**
- TV mit aktuellem Pattern implementieren (wie Movie)
- Danach beide refactoren
- ⏱️ TV schneller, aber mehr Duplikation

**Option C: Minimal jetzt, aggressiv später**
- Nur kritische Utils extrahieren (Levenshtein, Patterns)
- TV implementieren
- Großes Refactoring wenn 3+ Module da sind

**Frage**: Wann und wie aggressiv sollen wir refactoren?

---

## 4. Security-Fixes

### 4.1 Transaktionen

**Betroffene Stellen**:
1. User-Registrierung (`CreateUser` + `CreateEmailVerificationToken`)
2. Avatar-Upload (mehrere DB-Operationen)
3. Session-Refresh (`CreateSession` + `RevokeSession`)

**Frage**: Alle auf einmal fixen oder iterativ?

---

### 4.2 Timing-Attack im Login

**Fix**: Dummy-Hash-Vergleich auch wenn User nicht existiert

**Frage**:
- Argon2id-Parameter beibehalten? (memory=64MB, time=3, parallelism=2)
- Oder anpassen für bessere Performance/Security-Balance?

---

### 4.3 Account-Lockout

**Optionen**:
- **Option A: Service-Layer Lockout**
  - DB-Tabelle für Failed-Attempts
  - Lockout nach 5 Versuchen
  - 15-30 Minuten Sperre

- **Option B: Redis-basiert**
  - Schneller
  - Automatisches Expiry
  - Aber: State in Redis

**Frage**: Service-Layer mit DB oder Redis-basiert?

---

## 5. Deployment & Infrastructure

### 5.1 Deployment-Target

**Frage**: Welche Plattformen sind primär?
- [ ] Kubernetes (self-hosted)
- [ ] k3s (Lightweight K8s)
- [ ] Docker Swarm
- [ ] Docker Compose (Development)
- [ ] Bare Metal (direkt auf Server)

---

### 5.2 Monitoring & Observability

**Aktuell**:
- ✅ Prometheus-Metriken
- ✅ Strukturiertes Logging
- 🟡 OpenTelemetry importiert, aber nicht instrumentiert

**Frage**:
- Distributed Tracing jetzt implementieren oder später?
- Welche Monitoring-Stack bevorzugt? (Prometheus/Grafana, ELK, andere?)

---

## 6. Testing

### 6.1 Test-Coverage-Ziel

**Aktuell**: 41% (Movie), aber Ziel ist 80%

**Frage**:
- Erst 80% erreichen, dann neue Features?
- Oder parallel: neue Features mit 80% Coverage schreiben?

---

### 6.2 Integration-Test-Strategy

**Frage**:
- Testcontainers für alle Services? (PostgreSQL, Dragonfly, Typesense)
- Oder Mocks für schnellere Tests?

---

## 7. Migration-Path

### Prioritäten-Reihenfolge

**Vorschlag**:
1. **Woche 1-2**: Security-Fixes (Kritisch)
   - Transaktionen
   - Timing-Attacks
   - Goroutine-Leaks

2. **Woche 3-4**: Cluster-Readiness (Blocker für Production)
   - Media-Storage (NFS)
   - Avatar-Storage (S3)
   - Leader-Election (River)

3. **Woche 5-6**: TODOs in existierenden Modulen
   - IP/User-Agent Extraktion
   - SendGrid API
   - Async Reindex

4. **Woche 7-8**: Shared Abstractions
   - Scanner-Framework
   - Matcher-Framework

5. **Woche 9+**: TV Shows Implementation

**Frage**: Ist diese Reihenfolge OK oder andere Prioritäten?

---

## 8. TV Shows Spezifisch

### 8.1 TMDb TV API Struktur

**TMDb TV Endpoints**:
- `/tv/{id}` - Series Details
- `/tv/{id}/season/{season_number}` - Season Details
- `/tv/{id}/season/{season_number}/episode/{episode_number}` - Episode

**Frage**: Separate Tabellen für Series/Seasons/Episodes oder flachere Struktur?

```sql
-- Option A: Hierarchisch (empfohlen)
tvshow.series (id, tmdb_id, title, ...)
tvshow.seasons (id, series_id, season_number, ...)
tvshow.episodes (id, season_id, episode_number, title, ...)
tvshow.episode_files (id, episode_id, path, ...)

-- Option B: Flacher
tvshow.episodes (id, series_id, season_number, episode_number, ...)
tvshow.episode_files (id, episode_id, ...)
```

---

### 8.2 Filename-Parsing für TV

**Patterns**:
- `Series.Name.S01E05.Episode.Title.1080p.mkv`
- `Series Name - 1x05 - Episode Title.mkv`
- `Series.Name.105.Episode.Title.mkv` (season 1, episode 5)

**Frage**: Welche Naming-Conventions sind wichtig für dich?

---

## 9. Sonarr Integration

**Frage**:
- Sonarr-Integration gleichzeitig mit TV-Shows?
- Oder erst TV-Modul, dann später Sonarr?

---

## Zusammenfassung: Wichtigste Entscheidungen

**Bitte priorisieren (1 = höchste Priorität)**:

1. [ ] **Storage-Lösung** (NFS vs S3 vs Hybrid)
2. [ ] **Leader Election** (River Periodic vs andere)
3. [ ] **TMDb für TV** (ja/nein/zusätzlich TheTVDB)
4. [ ] **Tabellen-Architektur** (Separate Schemas vs Prefix)
5. [ ] **Refactoring-Timing** (jetzt vs schrittweise vs später)
6. [ ] **Security-Fix-Scope** (alle sofort vs iterativ)
7. [ ] **Deployment-Priorität** (K8s vs k3s vs Swarm)
8. [ ] **Test-Coverage** (erst 80% dann Features vs parallel)
9. [ ] **TV-Timing** (nach Cluster-Fixes vs nach Refactoring)
10. [ ] **Sonarr-Timing** (mit TV vs später)
