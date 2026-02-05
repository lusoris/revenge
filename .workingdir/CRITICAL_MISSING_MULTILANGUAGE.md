# CRITICAL: Multi-Language & International Ratings Support

**Discovered**: 2026-02-05
**Severity**: CRITICAL - Major Design Gap
**Impact**: Affects Movie, TV, and all content modules

---

## Problem

**Aktuelles Design geht davon aus**:
- Radarr/Sonarr liefern vollständige Metadaten
- Ein Titel pro Film/Serie
- Ein Rating-System (vermutlich MPAA)

**Realität**:
- ❌ Radarr/Sonarr holen NUR englische Metadaten
- ❌ Titel ändern sich je nach Sprache (z.B. "The Shawshank Redemption" vs "Die Verurteilten" vs "Les Évadés")
- ❌ Age Ratings unterscheiden sich pro Land (FSK, MPAA, BBFC, PEGI, etc.)
- ❌ Descriptions, Taglines, Keywords sind sprachabhängig

---

## Beispiel: The Shawshank Redemption

### Titel pro Sprache
- 🇺🇸 English: "The Shawshank Redemption"
- 🇩🇪 German: "Die Verurteilten"
- 🇫🇷 French: "Les Évadés"
- 🇪🇸 Spanish: "Cadena perpetua"
- 🇯🇵 Japanese: "ショーシャンクの空に" (Shōshanku no Sora ni)

### Age Ratings pro Land
- 🇺🇸 MPAA: R (Restricted)
- 🇩🇪 FSK: 12 (ab 12 Jahren)
- 🇬🇧 BBFC: 15
- 🇫🇷 CNC: Tous publics (All audiences)

### Overview/Description
Komplett unterschiedlicher Text je nach Sprache.

---

## Aktuelle Tabellen-Struktur (Movie)

```sql
-- internal/infra/database/queries/movie/movies.sql
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,              -- ❌ Nur eine Sprache!
    original_title TEXT,
    tagline TEXT,                     -- ❌ Nur eine Sprache!
    overview TEXT,                    -- ❌ Nur eine Sprache!
    -- ...
);
```

**Problem**: Keine Multi-Language-Unterstützung!

---

## Benötigte Änderungen

### Option A: JSONB Columns (Einfach)

```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    tmdb_id INTEGER,
    imdb_id TEXT,

    -- Multi-language fields as JSONB
    titles JSONB NOT NULL,           -- {"en": "Shawshank...", "de": "Die Verurteilten", ...}
    taglines JSONB,                  -- {"en": "Fear can...", "de": "Angst kann...", ...}
    overviews JSONB,                 -- {"en": "Two imprisoned...", "de": "Zwei Gefangene...", ...}

    -- Ratings per country
    age_ratings JSONB,               -- {"US": "R", "DE": "FSK 12", "GB": "15", ...}

    -- Original language
    original_language TEXT,          -- "en"
    original_title TEXT,             -- "The Shawshank Redemption"

    -- ...
);

-- Index for language lookups
CREATE INDEX idx_movies_titles ON movies USING GIN (titles);
```

**Vorteile**:
- ✅ Einfach zu implementieren
- ✅ Flexibel (neue Sprachen ohne Schema-Änderung)
- ✅ PostgreSQL hat gute JSONB-Unterstützung

**Nachteile**:
- ⚠️ Queries komplexer (`titles->>'de'`)
- ⚠️ Keine referentielle Integrität

---

### Option B: Separate Translation Tables (Normalized)

```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    tmdb_id INTEGER,
    imdb_id TEXT,
    original_language TEXT NOT NULL,
    original_title TEXT NOT NULL,
    -- ...
);

CREATE TABLE movie_translations (
    id UUID PRIMARY KEY,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    language TEXT NOT NULL,          -- "en", "de", "fr", "es", "ja", ...
    title TEXT NOT NULL,
    tagline TEXT,
    overview TEXT,
    UNIQUE(movie_id, language)
);

CREATE TABLE movie_age_ratings (
    id UUID PRIMARY KEY,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    country TEXT NOT NULL,           -- "US", "DE", "GB", "FR", ...
    rating_system TEXT NOT NULL,    -- "MPAA", "FSK", "BBFC", ...
    rating TEXT NOT NULL,            -- "R", "FSK 12", "15", ...
    UNIQUE(movie_id, country, rating_system)
);

CREATE INDEX idx_movie_translations_movie_lang ON movie_translations(movie_id, language);
CREATE INDEX idx_movie_age_ratings_movie_country ON movie_age_ratings(movie_id, country);
```

**Vorteile**:
- ✅ Normalized, saubere Datenstruktur
- ✅ Referentielle Integrität
- ✅ Einfache Queries mit JOINs

**Nachteile**:
- ⚠️ Mehr Tabellen
- ⚠️ JOINs für jede Language-Abfrage

---

### Option C: Hybrid (Empfohlen)

```sql
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    tmdb_id INTEGER,
    imdb_id TEXT,

    -- Default language (from Radarr/Sonarr)
    title TEXT NOT NULL,             -- Default (meist Englisch)
    tagline TEXT,
    overview TEXT,

    -- Multi-language als JSONB (für Performance bei UI)
    titles_i18n JSONB,               -- {"de": "...", "fr": "...", "es": "...", "ja": "..."}
    taglines_i18n JSONB,
    overviews_i18n JSONB,

    -- Original
    original_language TEXT,
    original_title TEXT,

    -- Age Ratings als JSONB (weniger häufig abgefragt)
    age_ratings JSONB,               -- {"US": {"MPAA": "R"}, "DE": {"FSK": "12"}, ...}

    -- ...
);
```

**Vorteile**:
- ✅ Einfache Queries für Default-Sprache
- ✅ Flexibel für Multi-Language
- ✅ Performant

**Nachteile**:
- ⚠️ Leichte Redundanz (title + titles_i18n)

---

## TMDb API Support

**TMDb unterstützt Multi-Language nativ**:

```
GET https://api.themoviedb.org/3/movie/278?language=de-DE
GET https://api.themoviedb.org/3/movie/278?language=fr-FR
GET https://api.themoviedb.org/3/movie/278?language=es-ES
```

**TMDb Age Ratings**:
```
GET https://api.themoviedb.org/3/movie/278/release_dates
```

Returns:
```json
{
  "results": [
    {
      "iso_3166_1": "US",
      "release_dates": [
        {
          "certification": "R",
          "type": 3
        }
      ]
    },
    {
      "iso_3166_1": "DE",
      "release_dates": [
        {
          "certification": "12",
          "type": 3
        }
      ]
    }
  ]
}
```

---

## TheTVDB API Support

**TheTVDB hat auch Multi-Language**:
```
GET https://api4.thetvdb.com/v4/series/76156
Headers: Accept-Language: de
```

---

## Impact auf bestehenden Code

### Movie Module

**Betroffene Dateien**:
- `internal/content/movie/types.go` - Movie struct
- `internal/content/movie/repository_postgres.go` - DB queries
- `internal/content/movie/tmdb_client.go` - API calls
- `internal/content/movie/tmdb_mapper.go` - Mapping
- `internal/infra/database/queries/movie/movies.sql` - SQL
- `migrations/` - Neue Migration nötig

**Beispiel Movie Struct (aktuell)**:
```go
type Movie struct {
    ID       uuid.UUID
    Title    string   // ❌ Nur eine Sprache
    Tagline  string   // ❌ Nur eine Sprache
    Overview string   // ❌ Nur eine Sprache
    // ...
}
```

**Beispiel Movie Struct (neu - Hybrid)**:
```go
type Movie struct {
    ID       uuid.UUID

    // Default language
    Title    string
    Tagline  string
    Overview string

    // Multi-language
    TitlesI18n    map[string]string  // language -> title
    TaglinesI18n  map[string]string
    OverviewsI18n map[string]string

    // Age ratings
    AgeRatings map[string]map[string]string  // country -> system -> rating

    // Original
    OriginalLanguage string
    OriginalTitle    string

    // ...
}
```

---

## Frontend Impact

**User-Einstellungen**:
- User wählt bevorzugte Sprache(n)
- UI zeigt Titel in User-Sprache oder Fallback auf Original

**Beispiel**:
- User-Sprache: Deutsch
- Film: "The Shawshank Redemption"
- Anzeige: "Die Verurteilten" (deutscher Titel)
- Fallback: "The Shawshank Redemption" (wenn keine deutsche Übersetzung)

---

## Migration-Strategie

### Phase 1: Schema-Erweiterung
1. Migration: Add JSONB columns zu movies
2. Bestehende Daten als "en" migrieren
3. Keine Breaking Changes (alte Columns bleiben)

### Phase 2: TMDb Multi-Language
1. TMDb-Client erweitern: Fetch multiple languages
2. Mapper: Parse translations
3. Repository: Store in JSONB

### Phase 3: Frontend
1. User-Einstellung: Bevorzugte Sprache
2. API: Return localized data based on user preference
3. Fallback-Logic

### Phase 4: TV Shows
1. Gleiche Struktur für TV (series, episodes)
2. TheTVDB Multi-Language

---

## Fragen zur Klärung

1. **Welche Sprachen sind Priorität?**
   - [ ] English (en)
   - [ ] German (de)
   - [ ] French (fr)
   - [ ] Spanish (es)
   - [ ] Japanese (ja)
   - [ ] Andere?

2. **Schema-Ansatz?**
   - [ ] Option A: JSONB Columns (einfach)
   - [ ] Option B: Separate Translation Tables (normalized)
   - [ ] Option C: Hybrid (empfohlen)

3. **Age Rating Systems?**
   - [ ] MPAA (US)
   - [ ] FSK (Germany)
   - [ ] BBFC (UK)
   - [ ] PEGI (Europe)
   - [ ] Andere?

4. **Default Language?**
   - [ ] Immer Original-Sprache
   - [ ] User-Einstellung
   - [ ] Radarr/Sonarr Language (meist en)

5. **Fallback-Strategy?**
   - [ ] Original -> English -> User-Language
   - [ ] User-Language -> English -> Original
   - [ ] User-konfigurierbar

---

## Aufwand-Schätzung

**Schema-Änderungen**: 8-16 Stunden
- Migrations schreiben
- Repository anpassen
- Types erweitern

**TMDb Multi-Language**: 16-24 Stunden
- Client erweitern (multiple language requests)
- Age Ratings fetchen
- Mapper anpassen
- Caching-Strategy

**Frontend**: 8-16 Stunden (wenn Frontend existiert)
- User-Settings für Sprache
- Localized display
- Fallback-Logic

**Total**: 32-56 Stunden (1-1.5 Wochen)

---

## Empfehlung

**Priorität**: HOCH (vor TV-Implementation)

**Begründung**:
- Affects alle Content-Module (Movie, TV, Music, etc.)
- Besser jetzt fixen als später migrieren
- TV-Module sollte gleich mit Multi-Language starten

**Reihenfolge**:
1. Design finalisieren (welche Option?)
2. Movie-Modul erweitern (als Proof-of-Concept)
3. Dann TV-Modul mit Multi-Language von Anfang an
