# API Reference

This document describes the HTTP API endpoints implemented in revenge.

## Overview

revenge provides a REST API that is fully compatible with the original Revenge API. All responses use JSON format.

### Base URL

```
http://localhost:8096
```

### Authentication

Most endpoints require authentication via Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

Or via custom header:

```
X-Revenge-Token: <access_token>
```

---

## Health Endpoints

### GET /health/live

Liveness check endpoint.

**Authentication:** Not required

**Response:**
```
OK
```

### GET /health/ready

Readiness check endpoint. Verifies database connectivity.

**Authentication:** Not required

**Response:**
```
Ready
```

### GET /health/db

Returns database connection pool statistics.

**Authentication:** Not required

**Response:**
```json
{
  "total_connections": 10,
  "idle_connections": 8,
  "active_connections": 2
}
```

---

## Authentication Endpoints

### POST /Users/AuthenticateByName

Authenticate a user with username and password.

**Authentication:** Not required

**Request:**
```json
{
  "Username": "admin",
  "Pw": "password"
}
```

**Response:**
```json
{
  "User": {
    "Id": "uuid",
    "Name": "admin",
    "ServerId": "uuid",
    "HasPassword": true,
    "HasConfiguredPassword": true,
    "Configuration": {},
    "Policy": {}
  },
  "AccessToken": "jwt_token",
  "ServerId": "uuid"
}
```

### POST /Sessions/Logout

End the current session.

**Authentication:** Required

**Response:** 204 No Content

### POST /Auth/Refresh

Refresh an access token.

**Authentication:** Required

**Request:**
```json
{
  "RefreshToken": "refresh_token"
}
```

**Response:**
```json
{
  "AccessToken": "new_jwt_token",
  "RefreshToken": "new_refresh_token",
  "ExpiresAt": "2024-01-15T10:30:00Z"
}
```

### POST /Users/{userId}/Password

Change user password.

**Authentication:** Required

**Request:**
```json
{
  "CurrentPw": "old_password",
  "NewPw": "new_password"
}
```

**Response:** 204 No Content

---

## User Endpoints

### GET /Users/Me

Get current authenticated user.

**Authentication:** Required

**Response:**
```json
{
  "Id": "uuid",
  "Name": "username",
  "ServerId": "uuid",
  "HasPassword": true,
  "HasConfiguredPassword": true,
  "Configuration": {},
  "Policy": {
    "IsAdministrator": false
  }
}
```

### GET /Users

List all users. Admin users see all users, non-admin users see limited info.

**Authentication:** Required

**Response:**
```json
[
  {
    "Id": "uuid",
    "Name": "username",
    ...
  }
]
```

### GET /Users/{userId}

Get a specific user by ID.

**Authentication:** Required

**Response:** User object

### POST /Users/New

Create a new user.

**Authentication:** Admin required

**Request:**
```json
{
  "Name": "newuser",
  "Password": "password"
}
```

**Response:** Created user object

### POST /Users/{userId}

Update a user.

**Authentication:** Required (admin for other users)

**Request:**
```json
{
  "Name": "updated_name"
}
```

**Response:** 204 No Content

### DELETE /Users/{userId}

Delete a user.

**Authentication:** Admin required

**Response:** 204 No Content

---

## Library Endpoints

### GET /Library/VirtualFolders

List all libraries. Non-admin users get filtered results based on permissions and adult content settings.

**Authentication:** Required

**Response:**
```json
[
  {
    "ItemId": "uuid",
    "Name": "Movies",
    "CollectionType": "movies",
    "Locations": ["/media/movies"],
    "LibraryOptions": {
      "EnableRealtimeMonitor": true,
      "EnableInternetProviders": true,
      "AutomaticRefreshIntervalDays": 1,
      "PathInfos": [
        {"Path": "/media/movies"}
      ]
    }
  }
]
```

### GET /Library/VirtualFolders/{libraryId}

Get a specific library.

**Authentication:** Required

**Response:** Library object

### POST /Library/VirtualFolders

Create a new library.

**Authentication:** Admin required

**Request:**
```json
{
  "Name": "Movies",
  "CollectionType": "movies",
  "Paths": ["/media/movies"],
  "RefreshLibrary": true,
  "LibraryOptions": {
    "AutomaticRefreshIntervalDays": 1
  }
}
```

**Collection Types:**
- `movies` - Movie library
- `tvshows` - TV show library
- `music` - Music library
- `musicvideos` - Music video library
- `photos` - Photo library
- `homevideos` - Home video library
- `boxsets` - Collections
- `livetv` - Live TV
- `playlists` - Playlists
- `mixed` - Mixed content
- `books` - Book library
- `audiobooks` - Audiobook library
- `podcasts` - Podcast library
- `adult_movie` - Adult movie library
- `adult_scene` - Adult scene library

**Response:** 201 Created with library object

### POST /Library/VirtualFolders/{libraryId}

Update a library.

**Authentication:** Admin required

**Request:**
```json
{
  "Name": "Updated Name",
  "Paths": ["/new/path"]
}
```

**Response:** 204 No Content

### DELETE /Library/VirtualFolders/{libraryId}

Delete a library.

**Authentication:** Admin required

**Response:** 204 No Content

### POST /Library/VirtualFolders/{libraryId}/Refresh

Trigger a library scan.

**Authentication:** Admin required

**Response:** 204 No Content

---

## Content Rating Endpoints

### GET /Ratings/Systems

List all rating systems.

**Authentication:** Required

**Query Parameters:**
- `country` (optional) - Filter by country code (e.g., "US", "DE")

**Response:**
```json
[
  {
    "Id": "uuid",
    "Code": "mpaa",
    "Name": "Motion Picture Association",
    "CountryCodes": ["US", "CA"],
    "IsActive": true
  },
  {
    "Id": "uuid",
    "Code": "fsk",
    "Name": "Freiwillige Selbstkontrolle",
    "CountryCodes": ["DE", "AT"],
    "IsActive": true
  }
]
```

### GET /Ratings/Systems/{code}

Get a rating system by code.

**Authentication:** Required

**Response:** Rating system object

### GET /Ratings/Systems/{systemId}/Ratings

List all ratings in a rating system.

**Authentication:** Required

**Response:**
```json
[
  {
    "Id": "uuid",
    "Code": "PG-13",
    "Name": "Parents Strongly Cautioned",
    "Description": "Some material may be inappropriate for children under 13",
    "MinAge": 13,
    "NormalizedLevel": 50,
    "IsAdult": false,
    "SystemCode": "mpaa",
    "SystemName": "Motion Picture Association"
  }
]
```

### GET /Items/{itemId}/Ratings

Get all ratings for a content item.

**Authentication:** Required

**Query Parameters:**
- `contentType` (optional) - Content type filter (default: "media_item")

**Response:**
```json
[
  {
    "Id": "uuid",
    "ContentId": "uuid",
    "ContentType": "media_item",
    "Source": "tmdb",
    "Rating": {
      "Id": "uuid",
      "Code": "PG-13",
      ...
    }
  }
]
```

### GET /Items/{itemId}/Rating

Get the display rating for a content item (single best rating).

**Authentication:** Required

**Query Parameters:**
- `contentType` (optional) - Content type filter
- `preferredSystem` (optional) - Preferred rating system code (default: "mpaa")

**Response:** Content rating object or null

### POST /Items/{itemId}/Ratings

Add a rating to a content item.

**Authentication:** Admin required

**Request:**
```json
{
  "RatingId": "uuid",
  "ContentType": "media_item",
  "Source": "manual"
}
```

**Response:** 201 Created with content rating object

### DELETE /Items/{itemId}/Ratings/{ratingId}

Remove a rating from a content item.

**Authentication:** Admin required

**Response:** 204 No Content

---

## Media Item Endpoints

### GET /Items

List media items with filtering and pagination.

**Authentication:** Required

**Query Parameters:**
- `parentId` - Parent item ID
- `libraryId` - Library ID filter
- `includeItemTypes` - Comma-separated list of types (e.g., "Movie,Episode")
- `genres` - Comma-separated genre filter
- `years` - Comma-separated year filter
- `tags` - Comma-separated tag filter
- `sortBy` - Sort field (e.g., "Name", "DateCreated", "ProductionYear")
- `sortOrder` - "asc" or "desc"
- `limit` - Max items to return (default: 50)
- `startIndex` - Pagination offset

**Response:**
```json
{
  "Items": [
    {
      "Id": "uuid",
      "Name": "Movie Title",
      "Type": "Movie",
      "MediaType": "Video",
      "ProductionYear": 2024,
      "RunTimeTicks": 72000000000,
      "CommunityRating": 8.5,
      "Genres": ["Action", "Comedy"],
      "IsFolder": false,
      ...
    }
  ],
  "TotalRecordCount": 100,
  "StartIndex": 0
}
```

### GET /Items/{itemId}

Get a specific media item.

**Authentication:** Required

**Response:** Media item object

### GET /Items/{itemId}/Similar

Get similar items.

**Authentication:** Required

**Response:** Items response

### GET /Items/{itemId}/ThemeSongs

Get theme songs for an item.

**Authentication:** Required

**Response:** Items response

### GET /Items/{itemId}/ThemeVideos

Get theme videos for an item.

**Authentication:** Required

**Response:** Items response

### GET /Users/{userId}/Items

List items for a specific user (includes user-specific data).

**Authentication:** Required

**Response:** Items response with UserData

### GET /Users/{userId}/Items/{itemId}

Get a specific item for a user.

**Authentication:** Required

**Response:** Item with UserData

### GET /Users/{userId}/Items/Latest

Get the latest added items for a user.

**Authentication:** Required

**Response:** Array of item objects

### GET /Users/{userId}/Items/Resume

Get items that are in progress (resume watching/listening).

**Authentication:** Required

**Response:** Items response

### GET /Search/Hints

Search for items.

**Authentication:** Required

**Query Parameters:**
- `searchTerm` (required) - Search query
- `limit` (optional) - Max results (default: 20)

**Response:**
```json
{
  "SearchHints": [
    {
      "ItemId": "uuid",
      "Id": "uuid",
      "Name": "Movie Title",
      "Type": "Movie",
      "MediaType": "Video",
      "ProductionYear": 2024
    }
  ],
  "TotalRecordCount": 10
}
```

### GET /Movies/Recommendations

Get movie recommendations.

**Authentication:** Required

**Response:** Array of recommendation groups

### GET /Shows/NextUp

Get next episodes to watch.

**Authentication:** Required

**Response:** Items response

---

## Content Rating System

revenge uses a normalized rating system that maps various international rating systems to a 0-100 scale:

| Level | Age | Description | Examples |
|-------|-----|-------------|----------|
| 0 | 0+ | All Ages | G, FSK 0, U |
| 25 | 6+ | Parental Guidance | PG, FSK 6 |
| 50 | 12+ | Teens | PG-13, FSK 12 |
| 75 | 16+ | Mature | R, FSK 16 |
| 90 | 18+ | Adults Only | NC-17, FSK 18 |
| 100 | 18+ | Adult/XXX | R18, X18+ |

### User Rating Settings

Users can configure:
- `max_rating_level` - Maximum allowed rating (0-100)
- `adult_content_enabled` - Whether to show adult (level 100) content
- `preferred_rating_system` - Preferred system for display (e.g., "fsk", "mpaa")

Content is automatically filtered based on user settings.

---

## Error Responses

All errors return a JSON object:

```json
{
  "error": "Error message description"
}
```

**Status Codes:**
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid auth)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `409` - Conflict (e.g., duplicate name)
- `500` - Internal Server Error
