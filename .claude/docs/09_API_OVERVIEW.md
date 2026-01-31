# Claude API Overview

> Source: https://platform.claude.com/docs/en/api/overview
> Fetched: 2026-01-31
> Type: html

---

## Endpoint

RESTful API at `https://api.anthropic.com`

Primary endpoint: `POST /v1/messages` for conversational interactions

---

## Prerequisites

- [Anthropic Console account](https://platform.claude.com)
- [API key](/settings/keys)

---

## Available APIs

### General Availability

| API | Endpoint | Description |
|-----|----------|-------------|
| **Messages** | `POST /v1/messages` | Send messages for conversational interactions |
| **Message Batches** | `POST /v1/messages/batches` | Async processing with 50% cost reduction |
| **Token Counting** | `POST /v1/messages/count_tokens` | Count tokens before sending |
| **Models** | `GET /v1/models` | List available models |

### Beta

| API | Endpoint | Description |
|-----|----------|-------------|
| **Files** | `POST /v1/files` | Upload and manage files |
| **Skills** | `POST /v1/skills` | Create and manage custom agent skills |

---

## Authentication

Required headers for all requests:

| Header | Value | Required |
|--------|-------|----------|
| `x-api-key` | Your API key | Yes |
| `anthropic-version` | e.g., `2023-06-01` | Yes |
| `content-type` | `application/json` | Yes |

SDKs handle these automatically.

---

## Client SDKs

**Benefits:**
- Automatic header management
- Type-safe request/response handling
- Built-in retry logic and error handling
- Streaming support
- Request timeouts and connection management

**Python Example:**
```python
from anthropic import Anthropic

client = Anthropic()  # Reads ANTHROPIC_API_KEY from environment
message = client.messages.create(
    model="claude-sonnet-4-5",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Hello, Claude"}]
)
```

---

## Third-Party Platforms

| Platform | Provider | Best For |
|----------|----------|----------|
| Amazon Bedrock | AWS | Existing AWS commitments |
| Vertex AI | Google Cloud | GCP integration |
| Azure AI | Microsoft Azure | Azure ecosystem |

---

## Request Size Limits

| Endpoint | Maximum Size |
|----------|--------------|
| Standard (Messages, Token Counting) | 32 MB |
| Batch API | 256 MB |
| Files API | 500 MB |

---

## Response Headers

- `request-id`: Globally unique request identifier
- `anthropic-organization-id`: Organization ID for the API key

---

## Basic Example

```bash
curl https://api.anthropic.com/v1/messages \
  --header "x-api-key: $ANTHROPIC_API_KEY" \
  --header "anthropic-version: 2023-06-01" \
  --header "content-type: application/json" \
  --data '{
    "model": "claude-sonnet-4-5",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Hello, Claude"}
    ]
  }'
```

**Response:**
```json
{
  "id": "msg_01XFDUDYJgAACzvnptvVoYEL",
  "type": "message",
  "role": "assistant",
  "content": [
    {
      "type": "text",
      "text": "Hello! How can I assist you today?"
    }
  ],
  "model": "claude-sonnet-4-5",
  "stop_reason": "end_turn",
  "usage": {
    "input_tokens": 12,
    "output_tokens": 8
  }
}
```

---

## SDK Examples

### Python

```python
import anthropic

client = anthropic.Anthropic()

message = client.messages.create(
    model="claude-sonnet-4-5",
    max_tokens=1000,
    messages=[
        {"role": "user", "content": "What should I search for?"}
    ]
)
print(message.content)
```

### TypeScript

```typescript
import Anthropic from "@anthropic-ai/sdk";

async function main() {
  const anthropic = new Anthropic();
  const msg = await anthropic.messages.create({
    model: "claude-sonnet-4-5",
    max_tokens: 1000,
    messages: [
      { role: "user", content: "Hello" }
    ]
  });
  console.log(msg);
}

main().catch(console.error);
```

### Java

```java
import com.anthropic.client.AnthropicClient;
import com.anthropic.client.okhttp.AnthropicOkHttpClient;
import com.anthropic.models.messages.Message;
import com.anthropic.models.messages.MessageCreateParams;

public class QuickStart {
    public static void main(String[] args) {
        AnthropicClient client = AnthropicOkHttpClient.fromEnv();

        MessageCreateParams params = MessageCreateParams.builder()
            .model("claude-sonnet-4-5-20250929")
            .maxTokens(1000)
            .addUserMessage("Hello")
            .build();

        Message message = client.messages().create(params);
        System.out.println(message.content());
    }
}
```

---

## Rate Limits

Organized into usage tiers that increase automatically:
- **Spend limits**: Maximum monthly cost
- **Rate limits**: RPM (requests/min) and TPM (tokens/min)

View limits in Console at `/settings/limits`.
