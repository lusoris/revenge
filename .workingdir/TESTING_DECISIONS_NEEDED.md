# Testing - Offene Fragen & Entscheidungen

**Datum**: 2026-02-04
**Status**: Vor Test-Writing - Klärung benötigt

---

## 🎯 Test-Strategie Fragen

### Q1: Mock-Generierung
**Frage**: Wie sollen Mocks erstellt werden?

**Optionen**:
- **A) mockery** - Automatische Mock-Generierung aus Interfaces
  - ✅ Weniger Arbeit
  - ✅ Type-safe
  - ❌ Extra Dependency
  - ❌ Generierte Dateien im Repo

- **B) Manuelle Mocks** - Per Hand in `*_test.go`
  - ✅ Keine Dependencies
  - ✅ Volle Kontrolle
  - ❌ Mehr Arbeit
  - ❌ Fehleranfälliger

- **C) gomock** - Google's mock framework
  - ✅ Sehr verbreitet
  - ❌ Komplexere Syntax
  - ❌ Extra Dependency

**Empfehlung**: mockery (bereits in TODO_v0.3.0.md erwähnt)

**Deine Entscheidung**: [PENDING]

---

### Q2: Integration Test Setup
**Frage**: Wie sollen Integration Tests laufen?

**Optionen**:
- **A) testcontainers-go** - Docker Container pro Test
  - ✅ Isoliert
  - ✅ Keine externe Abhängigkeit
  - ❌ Langsam (Container-Start)
  - ❌ Docker muss laufen

- **B) Shared Test Database** - Ein PostgreSQL für alle Tests
  - ✅ Schnell
  - ✅ Einfacher
  - ❌ Tests können sich beeinflussen
  - ❌ Muss manuell gestartet werden

- **C) In-Memory (pgx/pgmock)** - Keine echte DB
  - ✅ Sehr schnell
  - ❌ Nicht 100% realistisch
  - ❌ Komplexes Mocking

**Empfehlung**: testcontainers-go (beste Balance)

**Deine Entscheidung**: [PENDING]

---

### Q3: Test-Parallelisierung
**Frage**: Soll ich mehrere Packages gleichzeitig testen?

**Optionen**:
- **A) Sequenziell** - Ein Package nach dem anderen
  - ✅ Klarer Progress
  - ✅ Einfacher zu debuggen
  - ❌ Langsamer

- **B) Parallel** - Mehrere Packages gleichzeitig
  - ✅ Schneller
  - ❌ Schwerer zu tracken
  - ❌ Mehr Context Switching

**Empfehlung**: Sequenziell (bessere Qualität)

**Deine Entscheidung**: [PENDING]

---

### Q4: Test-Tiefe
**Frage**: Wie tief sollen Tests gehen?

**Optionen**:
- **A) Nur Happy Path** - Minimale Tests für Coverage
  - ✅ Schnell 80% erreicht
  - ❌ Schlechte Qualität
  - ❌ Viele Edge Cases ungecovered

- **B) Happy + Error Cases** - Standard Testing
  - ✅ Gute Balance
  - ✅ Die meisten Bugs gefangen
  - ⚠️ Moderate Zeit

- **C) Exhaustive** - Alle Edge Cases + Error Paths
  - ✅ Maximale Qualität
  - ❌ Sehr zeitaufwendig
  - ❌ Overkill für manche Funktionen

**Empfehlung**: Option B (Happy + Error Cases)

**Deine Entscheidung**: [PENDING]

---

## 📦 Package-Priorität

### Q5: Kritischste Packages
**Frage**: Falls Zeit knapp wird - welche Packages sind MUST-HAVE für 80%?

**Meine Analyse**:
1. **CRITICAL** (muss >80% sein):
   - Session Service (Auth-Foundation)
   - Auth Service (Security)
   - RBAC Service (Authorization)
   - User Service (Core Entity)

2. **HIGH** (sollte >70% sein):
   - Movie Service (Main Feature)
   - Library Service (Main Feature)
   - Settings Service (Config)

3. **MEDIUM** (kann >50% sein):
   - Activity Service (Audit)
   - Search Service (Feature)
   - Notification Service (Feature)

4. **LOW** (kann >30% sein):
   - API Handlers (werden durch Service-Tests teilweise getestet)
   - Integration Services (Radarr, TMDb)

**Deine Priorisierung**: [PENDING - Zustimmung oder Änderungen?]

---

## 🔧 Tool-Fragen

### Q6: Test-Helper-Library
**Frage**: Soll ich eine zentrale Test-Helper-Library erstellen?

**Was könnte drin sein**:
- Common fixtures (test users, movies, etc.)
- Helper functions (CreateTestUser, CreateTestMovie)
- Assertion helpers
- Mock builders

**Optionen**:
- **A) Ja** - `internal/testutil/` Package
  - ✅ DRY
  - ✅ Konsistenz
  - ❌ Upfront-Arbeit

- **B) Nein** - Jedes Package hat eigene Helpers
  - ✅ Einfacher Start
  - ❌ Code-Duplikation

**Empfehlung**: Ja (testutil Package)

**Deine Entscheidung**: [PENDING]

---

### Q7: Table-Driven Tests
**Frage**: Sollen alle Tests table-driven sein?

**Beispiel**:
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid", "test", "result", false},
    {"invalid", "", "", true},
}
```

**Optionen**:
- **A) Ja, immer** - Go best practice
  - ✅ Konsistent
  - ✅ Einfach zu erweitern
  - ⚠️ Mehr Boilerplate

- **B) Nur für repetitive Tests**
  - ✅ Flexibler
  - ❌ Inkonsistent

**Empfehlung**: Ja (Go Convention)

**Deine Entscheidung**: [PENDING]

---

## 🐛 Fehler-Behandlung

### Q8: Was bei Test-Failures?
**Frage**: Wie gehen wir mit Test-Failures um?

**Optionen**:
- **A) Fix sofort** - Test schreiben → Fehler fixen → weiter
  - ✅ Sauberer Code
  - ❌ Langsamer Progress

- **B) Log + Continue** - Fehler dokumentieren, später fixen
  - ✅ Schneller Progress
  - ❌ Technische Schulden

- **C) Skip Failing** - Mit `t.Skip()` markieren
  - ✅ Tests laufen durch
  - ❌ Coverage-Zahlen verzerrt

**Empfehlung**: Option A (Fix sofort)

**Deine Entscheidung**: [PENDING]

---

## 📊 Coverage-Reporting

### Q9: Coverage-Tracking
**Frage**: Wie tracken wir Progress?

**Optionen**:
- **A) Nach jedem Package** - Coverage-Report nach jedem abgeschlossenen Package
  - ✅ Klarer Progress
  - ⚠️ Viele Reports

- **B) Nach jeder Phase** - Report nach Phase 1, 2, 3, etc.
  - ✅ Weniger Overhead
  - ❌ Weniger Granularität

- **C) Continuous** - Coverage-Report im .workingdir updaten nach jedem Test-File
  - ✅ Maximale Transparenz
  - ❌ Viel File-IO

**Empfehlung**: Option B (Phase-basiert)

**Deine Entscheidung**: [PENDING]

---

## 🚀 Execution Plan

### Q10: Wie starten?
**Frage**: Womit genau soll ich anfangen?

**Mein Vorschlag**:
1. **Setup testutil Package** (30min)
   - Common fixtures
   - Mock generators
   - Helper functions

2. **Session Service Tests** (2-3h)
   - ValidateSession
   - CreateSession
   - RevokeSession
   - Cleanup expired

3. **Erste Coverage-Messung**
   - Sehen ob Strategie funktioniert
   - Bei Bedarf adjustieren

**Alternative**:
- Direkt mit Session Service starten, ohne testutil

**Deine Entscheidung**: [PENDING]

---

## 📝 Zusammenfassung

**Quick-Decisions (falls du zustimmst)**:
- ✅ mockery für Mocks
- ✅ testcontainers für Integration Tests
- ✅ Sequenziell (ein Package nach dem anderen)
- ✅ Happy + Error Cases
- ✅ Priorisierung wie oben
- ✅ testutil Package erstellen
- ✅ Table-driven Tests
- ✅ Fehler sofort fixen
- ✅ Phase-basiertes Coverage-Tracking
- ✅ Start mit testutil → Session Service

**Wenn du allem zustimmst**: Sage einfach "go" und ich starte mit testutil + Session Service Tests

**Wenn Änderungen**: Sag mir welche Punkte du anders haben willst

---

**Status**: ⏸️ Warte auf User-Feedback
