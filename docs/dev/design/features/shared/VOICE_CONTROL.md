# Voice Control

> Voice assistant integration (Alexa, Google Assistant)

**Status**: 🔴 PLANNING
**Priority**: 🔵 LOW (Nice to have - Emby has this)
**Inspired By**: Emby Voice Control

---

## Overview

Voice control allows users to control playback and browse content using voice commands through Amazon Alexa and Google Assistant.

---

## Supported Platforms

| Platform | Integration Type | Status |
|----------|-----------------|--------|
| Amazon Alexa | Smart Home Skill | 🟡 Planned |
| Google Assistant | Actions on Google | 🟡 Planned |
| Apple Siri | HomeKit / Shortcuts | 🔴 Future |

---

## Features

### Playback Control

| Command | Example |
|---------|---------|
| Play | "Play The Matrix" |
| Pause | "Pause" |
| Resume | "Resume" |
| Stop | "Stop" |
| Next | "Next episode" |
| Previous | "Previous episode" |
| Seek | "Skip forward 30 seconds" |

### Content Navigation

| Command | Example |
|---------|---------|
| Browse | "Show my movies" |
| Search | "Search for action movies" |
| Play Series | "Play Breaking Bad" |
| Continue | "Continue watching" |
| Recommendations | "What should I watch?" |

### Information

| Command | Example |
|---------|---------|
| What's Playing | "What's playing?" |
| Episode Info | "What episode is this?" |
| Duration | "How long is this movie?" |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Voice Control Flow                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   User Voice ──► Alexa/Google ──► Revenge Skill ──► API        │
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│  │   "Play     │───►│   Amazon    │───►│  Revenge    │        │
│  │ The Matrix" │    │   Lambda    │    │   API       │        │
│  └─────────────┘    └─────────────┘    └─────────────┘        │
│                            │                  │                 │
│                            │           ┌──────▼──────┐         │
│                            │           │   Search    │         │
│                            │           │   Service   │         │
│                            │           └──────┬──────┘         │
│                            │                  │                 │
│                            │           ┌──────▼──────┐         │
│                            └──────────►│   Device    │         │
│                                        │  Playback   │         │
│                                        └─────────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Alexa Skill

### Skill Definition

```json
{
    "interactionModel": {
        "languageModel": {
            "invocationName": "revenge",
            "intents": [
                {
                    "name": "PlayMediaIntent",
                    "slots": [
                        {
                            "name": "title",
                            "type": "AMAZON.VideoGame"
                        }
                    ],
                    "samples": [
                        "play {title}",
                        "watch {title}",
                        "start {title}"
                    ]
                },
                {
                    "name": "AMAZON.PauseIntent"
                },
                {
                    "name": "AMAZON.ResumeIntent"
                },
                {
                    "name": "AMAZON.NextIntent"
                },
                {
                    "name": "AMAZON.PreviousIntent"
                },
                {
                    "name": "ContinueWatchingIntent",
                    "samples": [
                        "continue watching",
                        "resume my show",
                        "what was I watching"
                    ]
                }
            ]
        }
    }
}
```

### Lambda Handler (Go)

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/revenge/voice/alexa"
)

type AlexaRequest struct {
    Version string `json:"version"`
    Request struct {
        Type   string `json:"type"`
        Intent struct {
            Name  string `json:"name"`
            Slots map[string]struct {
                Value string `json:"value"`
            } `json:"slots"`
        } `json:"intent"`
    } `json:"request"`
}

type AlexaResponse struct {
    Version  string `json:"version"`
    Response struct {
        OutputSpeech struct {
            Type string `json:"type"`
            Text string `json:"text"`
        } `json:"outputSpeech"`
        ShouldEndSession bool `json:"shouldEndSession"`
    } `json:"response"`
}

func handler(ctx context.Context, req AlexaRequest) (*AlexaResponse, error) {
    switch req.Request.Intent.Name {
    case "PlayMediaIntent":
        title := req.Request.Intent.Slots["title"].Value
        return handlePlayMedia(ctx, title)
    case "AMAZON.PauseIntent":
        return handlePause(ctx)
    case "AMAZON.ResumeIntent":
        return handleResume(ctx)
    case "ContinueWatchingIntent":
        return handleContinueWatching(ctx)
    default:
        return unknownIntent()
    }
}

func handlePlayMedia(ctx context.Context, title string) (*AlexaResponse, error) {
    // Search for content
    results, err := searchContent(ctx, title)
    if err != nil || len(results) == 0 {
        return speak("I couldn't find " + title), nil
    }

    // Send play command to active device
    err = playOnDevice(ctx, results[0].ID)
    if err != nil {
        return speak("I had trouble starting playback"), nil
    }

    return speak("Now playing " + results[0].Title), nil
}

func speak(text string) *AlexaResponse {
    return &AlexaResponse{
        Version: "1.0",
        Response: struct {
            OutputSpeech struct {
                Type string `json:"type"`
                Text string `json:"text"`
            } `json:"outputSpeech"`
            ShouldEndSession bool `json:"shouldEndSession"`
        }{
            OutputSpeech: struct {
                Type string `json:"type"`
                Text string `json:"text"`
            }{
                Type: "PlainText",
                Text: text,
            },
            ShouldEndSession: true,
        },
    }
}

func main() {
    lambda.Start(handler)
}
```

---

## Google Assistant

### Actions Definition

```yaml
# actions.yaml
actions:
  - name: actions.intent.PLAY_MEDIA
    fulfillment:
      conversationName: play_media
  - name: actions.intent.MEDIA_PAUSE
    fulfillment:
      conversationName: media_control
  - name: actions.intent.MEDIA_RESUME
    fulfillment:
      conversationName: media_control

conversations:
  play_media:
    name: Play Media
    url: https://revenge.example.com/api/v1/voice/google
    fulfillmentApiVersion: 2

  media_control:
    name: Media Control
    url: https://revenge.example.com/api/v1/voice/google
    fulfillmentApiVersion: 2
```

---

## Database Schema

```sql
-- Voice-linked devices
CREATE TABLE voice_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL, -- alexa, google, siri
    device_id VARCHAR(200) NOT NULL,
    device_name VARCHAR(200),

    -- Linked playback device
    playback_device_id UUID REFERENCES devices(id) ON DELETE SET NULL,

    -- Auth
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, platform, device_id)
);

-- Voice command logs
CREATE TABLE voice_commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voice_device_id UUID REFERENCES voice_devices(id) ON DELETE SET NULL,

    platform VARCHAR(50) NOT NULL,
    intent VARCHAR(100) NOT NULL,
    raw_text TEXT,
    slots JSONB,

    -- Result
    success BOOLEAN,
    response_text TEXT,
    action_taken TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_voice_devices_user ON voice_devices(user_id);
CREATE INDEX idx_voice_commands_user ON voice_commands(user_id);
```

---

## API Endpoints

```
# Device linking
GET  /api/v1/voice/devices              # List linked devices
POST /api/v1/voice/devices/link         # Link new device
DELETE /api/v1/voice/devices/:id        # Unlink device
PUT  /api/v1/voice/devices/:id          # Update device settings

# Voice handlers (called by Alexa/Google)
POST /api/v1/voice/alexa                # Alexa skill webhook
POST /api/v1/voice/google               # Google Actions webhook

# OAuth endpoints (for account linking)
GET  /api/v1/voice/oauth/authorize      # OAuth authorize
POST /api/v1/voice/oauth/token          # OAuth token
```

---

## Configuration

```yaml
voice_control:
  enabled: true

  alexa:
    enabled: true
    skill_id: "amzn1.ask.skill.xxx"
    client_id: "${ALEXA_CLIENT_ID}"
    client_secret: "${ALEXA_CLIENT_SECRET}"

  google:
    enabled: true
    project_id: "revenge-voice"
    service_account_key: "/config/google-service-account.json"

  # Default device selection
  device_selection:
    prefer_active: true  # Use currently active device
    fallback_to_last: true  # Use last used device
```

---

## Account Linking

### Alexa Account Linking Flow

```
1. User enables Revenge skill in Alexa app
2. Redirected to Revenge OAuth authorize endpoint
3. User logs in to Revenge
4. Authorization code sent to Alexa
5. Alexa exchanges code for tokens
6. Skill linked to user account
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `voice.control` | Use voice commands |
| `voice.link` | Link voice devices |
| `voice.admin` | View all voice devices |

---

## Related Documentation

- [Client Support](CLIENT_SUPPORT.md)
- [Go Packages](../architecture/GO_PACKAGES.md)
