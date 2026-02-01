## Table of Contents

- [Voice Control](#voice-control)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
- [Voice webhook endpoints (public, verified)](#voice-webhook-endpoints-public-verified)
- [User management](#user-management)
- [Command history](#command-history)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Voice Control


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Voice assistant integration (Alexa, Google Assistant)

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | 🟡 | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```
  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐
  │   Client    │────▶│  API Handler │────▶│   Service   │
  │  (Web/App)  │◀────│   (ogen)     │◀────│   (Logic)   │
  └─────────────┘     └──────────────┘     └──────┬──────┘
                                                   │
                            ┌──────────────────────┼────────────┐
                            ▼                      ▼            ▼
                      ┌──────────┐          ┌───────────┐  ┌────────┐
                      │Repository│          │ Metadata  │  │  Cache │
                      │  (sqlc)  │          │  Service  │  │(otter) │
                      └────┬─────┘          └─────┬─────┘  └────────┘
                           │                      │
                           ▼                      ▼
                    ┌─────────────┐        ┌──────────┐
                    │ PostgreSQL  │        │ External │
                    │   (pgx)     │        │   APIs   │
                    └─────────────┘        └──────────┘
  ```

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/voice_control/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── voice_control_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### Key Interfaces

```go
type VoiceService interface {
  // Connections
  ConnectAssistant(ctx context.Context, userID uuid.UUID, assistantType string, tokens OAuth2Tokens) error
  DisconnectAssistant(ctx context.Context, connectionID uuid.UUID) error
  GetConnections(ctx context.Context, userID uuid.UUID) ([]VoiceConnection, error)

  // Command handling
  ProcessAlexaRequest(ctx context.Context, request AlexaRequest) (*AlexaResponse, error)
  ProcessGoogleRequest(ctx context.Context, request GoogleActionRequest) (*GoogleActionResponse, error)

  // Intent parsing
  ParseIntent(ctx context.Context, rawCommand string) (*VoiceIntent, error)
  ExecuteIntent(ctx context.Context, userID uuid.UUID, intent VoiceIntent) (*IntentResult, error)
}

type VoiceIntent struct {
  Intent   string                 `json:"intent"`    // 'play_media', 'pause', 'skip'
  Entities map[string]interface{} `json:"entities"`  // {"title": "Inception", "type": "movie"}
  Raw      string                 `json:"raw"`
}

type AlexaRequest struct {
  Version string      `json:"version"`
  Session AlexaSession `json:"session"`
  Request struct {
    Type   string `json:"type"`
    Intent struct {
      Name  string               `json:"name"`
      Slots map[string]AlexaSlot `json:"slots"`
    } `json:"intent"`
  } `json:"request"`
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/aws/aws-sdk-go-v2` - Alexa skill verification
- `google.golang.org/api/dialogflow/v2` - Google Assistant integration
- `go.uber.org/fx`

**External APIs**:
- Amazon Alexa Skills Kit (ASK)
- Google Assistant SDK / Dialogflow
- Apple HomeKit (future)







## Configuration

### Environment Variables

```bash
VOICE_ALEXA_SKILL_ID=amzn1.ask.skill.xxxxx
VOICE_GOOGLE_PROJECT_ID=revenge-voice-xxxxx
VOICE_ENABLED=true
```


### Config Keys
```yaml
voice:
  enabled: true
  alexa:
    skill_id: amzn1.ask.skill.xxxxx
    verification_enabled: true
  google:
    project_id: revenge-voice-xxxxx
    credentials_file: /config/google-service-account.json
```



## API Endpoints

### Content Management
```
# Voice webhook endpoints (public, verified)
POST /api/v1/voice/alexa            # Alexa skill webhook
POST /api/v1/voice/google           # Google Assistant webhook

# User management
GET  /api/v1/voice/connections      # List connections
DELETE /api/v1/voice/connections/:id # Disconnect assistant

# Command history
GET  /api/v1/voice/commands         # Get command history
```

**Example Alexa Request**:
```json
{
  "version": "1.0",
  "session": {
    "user": {"userId": "amzn1.ask.account.XXXXX"}
  },
  "request": {
    "type": "IntentRequest",
    "intent": {
      "name": "PlayMediaIntent",
      "slots": {
        "title": {"value": "Inception"},
        "type": {"value": "movie"}
      }
    }
  }
}
```

**Example Alexa Response**:
```json
{
  "version": "1.0",
  "response": {
    "outputSpeech": {
      "type": "PlainText",
      "text": "Playing Inception"
    },
    "shouldEndSession": true
  }
}
```








## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config

