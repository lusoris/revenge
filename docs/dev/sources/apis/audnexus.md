# Audnexus API

> Source: https://docs.audnex.us/
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

Audnexus provides metadata for audiobooks, aggregating data from multiple sources including Audible.

## API Base URL

```
https://api.audnex.us
```

## Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/books/{asin}` | GET | Get book metadata by ASIN |
| `/authors/{asin}` | GET | Get author metadata |
| `/chapters/{asin}` | GET | Get chapter data |

## Authentication

API key authentication via header:
```
Authorization: Bearer {api_key}
```

## Related

- [Audiobook Module](../../design/features/audiobook/AUDIOBOOK_MODULE.md)
- [Audible Integration](../../design/integrations/metadata/books/AUDIBLE.md)
