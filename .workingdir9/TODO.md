# Codebase Analyse — Priorisierte Action Items

**Datum:** 2026-02-10

---

## P0 — Muss vor Frontend-Arbeit gefixt werden

### 1. snake_case vs camelCase vereinheitlichen
**Impact:** Jeder einzelne API-Call im Frontend betroffen
**Scope:** OpenAPI Spec + ogen regenerieren + Handler-Converter anpassen

Betroffene Ressourcen die auf camelCase umgestellt werden müssen (oder umgekehrt — Konvention wählen):
- Auth, MFA, Movies, TV Shows (aktuell snake_case)
- OIDC, Activity, Libraries, Radarr/Sonarr, Playback (aktuell camelCase)

**Empfehlung:** snake_case als Standard (REST-Konvention, Go-Backend), camelCase-Ressourcen anpassen.

### 2. List Endpoints mit Envelope wrappen
**Impact:** Ohne `total` kein Pagination-UI möglich
**Scope:** 17+ Endpoints die nackte Arrays zurückgeben

Standard-Envelope:
```json
{
  "items": [...],
  "total": 42,
  "limit": 20,
  "offset": 0
}
```

Betroffene Endpoints:
- `GET /movies/search`, `/continue-watching`, `/watch-history`
- `GET /movies/{id}/files`, `/genres`
- `GET /tvshows/search`, `/continue-watching`
- `GET /tvshows/episodes/recent`, `/upcoming`
- `GET /tvshows/{id}/seasons`, `/episodes`, `/genres`, `/networks`
- `GET /tvshows/seasons/{id}/episodes`
- `GET /tvshows/episodes/{id}/files`
- `GET /collections/{id}/movies`
- `GET /genres`

---

## P1 — Sollte vor Frontend-Arbeit gefixt werden

### 3. Pagination standardisieren
**Impact:** Sonst 3 verschiedene Pagination-Logiken im Frontend
**Action:** `limit`/`offset` als Standard, Typesense-Search intern auf `limit`/`offset` mappen (berechnet intern `page`/`per_page`)

### 4. HTTP Methods angleichen
**Action:** Movie Progress `POST` → `PUT` (oder Episode Progress `PUT` → `POST`). Empfehlung: beide auf `PUT` (idempotent).

### 5. Sort-Parameter vereinheitlichen
**Action:** Alle auf `sort_by` + `sort_order` standardisieren (oder `order_by` + `order`). Enum-Werte konsistent machen.

### 6. Dupliziertes TVShowListResponse/TVSeriesListResponse bereinigen
**Action:** Eins entfernen, Referenzen auf das verbleibende umleiten.

---

## P2 — Sollte zeitnah gefixt werden

### 7. Fehlende OpenAPI Tag-Definitionen
**Action:** 8 fehlende Tags zum Top-Level `tags` Array hinzufügen, Casing vereinheitlichen.

### 8. Refresh-Endpoints Job-ID returnen
**Action:** `POST .../movies/{id}/refresh` und `POST .../tvshows/{id}/refresh` sollen `{ message, job_id }` zurückgeben (wie Search Reindex).

### 9. MFA Inline-Schemas extrahieren
**Action:** Ad-hoc `{ success, message }` Schemas als wiederverwendbare `$ref` Komponente definieren.

### 10. RefreshPersonWorker fixen
**Action:** Entweder `ErrNotImplemented` returnen (damit Jobs im Failed-State landen) oder Worker komplett deregistrieren bis Person-Service existiert.

---

## P3 — Nice to have / Technical Debt

### 11. Ungenutzter DB-Pool in migrate.go entfernen
**Action:** `database.NewPool` Call + Ping in `cmd/revenge/migrate.go` entfernen.

### 12. `panic()` in `river.Subscribe()` durch Error-Return ersetzen
**Action:** `if c.client == nil { return nil, func(){} }` statt `panic`.

### 13. Email Service Interface + Mock hinzufügen
**Action:** `email.Sender` Interface extrahieren, mockery Mock generieren. Verbessert Testbarkeit von auth/user Services.

### 14. Fehlende Service-Mocks generieren
**Scope:** metadata (Service), mfa, rbac, search — je nach Bedarf für Frontend-unterstützende Unit Tests.

### 15. Typesense-Syntax abstrahieren
**Action:** `filter_by` in Search-Endpoints durch strukturierte Filter-Parameter ersetzen (z.B. `genres`, `year_min`, `year_max`), Backend baut Typesense-Query.

### 16. movies.sql aufspalten
**Action:** Monolithische `queries/movie/movies.sql` (666 Zeilen) in domain-spezifische Dateien aufteilen (wie bei tvshow: series.sql, seasons.sql, episodes.sql, etc.).

### 17. Leere Migration-Directories aufräumen
**Action:** `migrations/qar/` und `migrations/tvshow/` entfernen oder dokumentieren warum sie existieren.

### 18. `validate.SafeInt32` vs `util.SafeIntToInt32` dokumentieren
**Action:** Cross-Reference-Comments hinzufügen die den Unterschied (error vs saturate) erklären.

---

## Nicht blockierend (Future)

- [ ] ETag/If-None-Match Support für SvelteKit `load()` Caching
- [ ] Cursor-basierte Pagination für große Libraries
- [ ] Circuit Breaker Integration (gobreaker Dep existiert)
- [ ] Request Coalescing
- [ ] Cache Warming on Startup
- [ ] Search Observability (Typesense Metrics)
- [ ] Person Service implementieren (dann RefreshPersonWorker aktivieren)
